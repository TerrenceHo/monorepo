#!/usr/bin/env python3
"""
bcr-update — find `bazel_dep(...)` declarations across MODULE.bazel files,
look up the latest version in the Bazel Central Registry (BCR), and either
report what's out of date or rewrite the versions in place.

Why Python (not bash): MODULE.bazel is written in Starlark, a dialect of
Python. The subset used for module directives (`bazel_dep`, `module`,
`use_repo`, ...) is valid Python syntax, so the standard-library `ast` module
parses it into a real syntax tree with exact source positions — no third-party
dependencies, and far more robust than regex. We use the AST to locate each
version string literal precisely, then do a surgical text replacement that
preserves all surrounding formatting and comments. If a file ever contains
Starlark that `ast` can't parse, we fall back to a tolerant regex scan so the
tool still works.

Usage:
    bcr_update.py [PATHS...]            # check mode: report out-of-date deps
    bcr_update.py -u [PATHS...]         # update mode: rewrite versions in place

PATHS may be directories (scanned recursively for MODULE.bazel) or individual
files. Defaults to the current directory.
"""

from __future__ import annotations

import argparse
import ast
import json
import os
import re
import sys
import urllib.error
import urllib.parse
import urllib.request
from concurrent.futures import ThreadPoolExecutor
from dataclasses import dataclass
from functools import cmp_to_key

DEFAULT_REGISTRY = "https://bcr.bazel.build"


# --------------------------------------------------------------------------- #
# Data model
# --------------------------------------------------------------------------- #
@dataclass
class Dep:
    file: str
    name: str
    version: str | None  # None => bazel_dep with no registry version
    # Position of the version string literal (for in-place edits). 1-based line,
    # 0-based columns, matching the ast module. None when version is None.
    line: int | None = None
    col: int | None = None
    end_line: int | None = None
    end_col: int | None = None


# --------------------------------------------------------------------------- #
# Discovery
# --------------------------------------------------------------------------- #
def find_module_files(paths: list[str]) -> list[str]:
    found: list[str] = []
    seen: set[str] = set()

    def add(p: str) -> None:
        rp = os.path.realpath(p)
        if rp not in seen:
            seen.add(rp)
            found.append(p)

    for path in paths:
        if os.path.isfile(path):
            add(path)
            continue
        for dirpath, dirnames, filenames in os.walk(path):
            # Don't descend into Bazel's symlinked output trees or VCS dirs.
            dirnames[:] = [
                d
                for d in dirnames
                if not d.startswith("bazel-") and d not in (".git", ".hg", ".svn")
            ]
            if "MODULE.bazel" in filenames:
                add(os.path.join(dirpath, "MODULE.bazel"))
    return found


# --------------------------------------------------------------------------- #
# Parsing bazel_dep declarations
# --------------------------------------------------------------------------- #
def _const_str(node: ast.AST | None) -> str | None:
    if isinstance(node, ast.Constant) and isinstance(node.value, str):
        return node.value
    return None


def parse_deps(path: str, src: str) -> list[Dep]:
    """Extract bazel_dep declarations. Tries a real AST parse first; falls back
    to a tolerant regex scan if the file isn't valid Python-style Starlark."""
    try:
        tree = ast.parse(src, filename=path)
    except SyntaxError:
        return _parse_deps_regex(path, src)

    deps: list[Dep] = []
    for node in ast.walk(tree):
        if not (
            isinstance(node, ast.Call)
            and isinstance(node.func, ast.Name)
            and node.func.id == "bazel_dep"
        ):
            continue
        kw = {k.arg: k.value for k in node.keywords if k.arg}
        name = _const_str(kw.get("name"))
        if name is None:
            continue
        vnode = kw.get("version")
        vstr = _const_str(vnode)
        if vnode is None or vstr is None:
            # e.g. paired with a non-registry override; nothing to bump.
            deps.append(Dep(file=path, name=name, version=None))
            continue
        deps.append(
            Dep(
                file=path,
                name=name,
                version=vstr,
                line=vnode.lineno,
                col=vnode.col_offset,
                end_line=vnode.end_lineno,
                end_col=vnode.end_col_offset,
            )
        )
    return deps


_REGEX_DEP = re.compile(r"bazel_dep\s*\((?P<args>.*?)\)", re.DOTALL)
_REGEX_NAME = re.compile(r"""name\s*=\s*["']([^"']+)["']""")
_REGEX_VERSION = re.compile(r"""version\s*=\s*["']([^"']+)["']""")


def _parse_deps_regex(path: str, src: str) -> list[Dep]:
    deps: list[Dep] = []
    for m in _REGEX_DEP.finditer(src):
        args = m.group("args")
        nm = _REGEX_NAME.search(args)
        if not nm:
            continue
        vm = _REGEX_VERSION.search(args)
        if not vm:
            deps.append(Dep(file=path, name=nm.group(1), version=None))
            continue
        # Compute line/col of the version literal for in-place edits.
        start = m.start("args") + vm.start(1) - 1  # include opening quote
        end = m.start("args") + vm.end(1) + 1  # include closing quote
        line = src.count("\n", 0, start) + 1
        col = start - (src.rfind("\n", 0, start) + 1)
        end_line = src.count("\n", 0, end) + 1
        end_col = end - (src.rfind("\n", 0, end) + 1)
        deps.append(
            Dep(
                file=path,
                name=nm.group(1),
                version=vm.group(1),
                line=line,
                col=col,
                end_line=end_line,
                end_col=end_col,
            )
        )
    return deps


# --------------------------------------------------------------------------- #
# Version comparison (SemVer / Bazel module version, build metadata ignored)
# --------------------------------------------------------------------------- #
def _split_version(v: str) -> tuple[list[str], list[str] | None]:
    v = v.split("+", 1)[0]  # drop build metadata
    if "-" in v:
        rel, pre = v.split("-", 1)
        return rel.split("."), pre.split(".")
    return v.split("."), None


def _cmp_ident(x: str, y: str) -> int:
    xd, yd = x.isdigit(), y.isdigit()
    if xd and yd:
        ix, iy = int(x), int(y)
        return (ix > iy) - (ix < iy)
    if xd != yd:  # numeric ranks below alphanumeric
        return -1 if xd else 1
    return (x > y) - (x < y)


def _cmp_segments(a: list[str], b: list[str]) -> int:
    for i in range(max(len(a), len(b))):
        if i >= len(a):
            return -1
        if i >= len(b):
            return 1
        c = _cmp_ident(a[i], b[i])
        if c:
            return c
    return 0


def compare_version(a: str, b: str) -> int:
    ar, ap = _split_version(a)
    br, bp = _split_version(b)
    c = _cmp_segments(ar, br)
    if c:
        return c
    # A version with no prerelease outranks the same version with one.
    if ap is None and bp is None:
        return 0
    if ap is None:
        return 1
    if bp is None:
        return -1
    return _cmp_segments(ap, bp)


def is_prerelease(v: str) -> bool:
    return "-" in v.split("+", 1)[0]


# --------------------------------------------------------------------------- #
# BCR lookup
# --------------------------------------------------------------------------- #
def fetch_latest(
    name: str, registry: str, include_prerelease: bool, timeout: float
) -> tuple[str, str | None]:
    """Return (status, latest_version_or_message).
    status in {"ok", "not_found", "no_versions", "error"}."""
    url = f"{registry.rstrip('/')}/modules/{urllib.parse.quote(name)}/metadata.json"
    try:
        req = urllib.request.Request(url, headers={"User-Agent": "bcr-update"})
        with urllib.request.urlopen(req, timeout=timeout) as resp:
            data = json.load(resp)
    except urllib.error.HTTPError as e:
        if e.code == 404:
            return ("not_found", None)
        return ("error", f"HTTP {e.code}")
    except (urllib.error.URLError, TimeoutError, json.JSONDecodeError) as e:
        return ("error", str(e))

    versions = data.get("versions") or []
    yanked = set((data.get("yanked_versions") or {}).keys())
    candidates = [v for v in versions if v and v not in yanked]
    if not include_prerelease:
        stable = [v for v in candidates if not is_prerelease(v)]
        if stable:  # only fall back to prereleases if no stable
            candidates = stable
    if not candidates:
        return ("no_versions", None)
    return ("ok", max(candidates, key=cmp_to_key(compare_version)))


# --------------------------------------------------------------------------- #
# In-place editing
# --------------------------------------------------------------------------- #
def apply_updates(path: str, src: str, updates: list[tuple[Dep, str]]) -> str:
    """Replace each version literal with its new value, preserving everything
    else. `updates` is a list of (dep, new_version)."""
    lines = src.splitlines(keepends=True)
    # Group edits by line so multiple edits on one line don't corrupt offsets.
    by_line: dict[int, list[tuple[Dep, str]]] = {}
    for dep, new in updates:
        if dep.line is None or dep.end_line != dep.line:
            continue  # multi-line/unknown span: skip rather than risk damage
        by_line.setdefault(dep.line, []).append((dep, new))

    for lineno, edits in by_line.items():
        idx = lineno - 1
        line = lines[idx]
        # Apply right-to-left so earlier columns stay valid.
        for dep, new in sorted(edits, key=lambda e: e[0].col, reverse=True):
            old_literal = line[dep.col : dep.end_col]
            quote = old_literal[0] if old_literal[:1] in ("'", '"') else '"'
            line = line[: dep.col] + f"{quote}{new}{quote}" + line[dep.end_col :]
        lines[idx] = line
    return "".join(lines)


# --------------------------------------------------------------------------- #
# Pretty output
# --------------------------------------------------------------------------- #
class C:
    def __init__(self, enabled: bool):
        self.e = enabled

    def _w(self, code: str, s: str) -> str:
        return f"\033[{code}m{s}\033[0m" if self.e else s

    def bold(self, s):
        return self._w("1", s)

    def dim(self, s):
        return self._w("2", s)

    def green(self, s):
        return self._w("32", s)

    def yellow(self, s):
        return self._w("33", s)

    def red(self, s):
        return self._w("31", s)

    def cyan(self, s):
        return self._w("36", s)


# --------------------------------------------------------------------------- #
# Main
# --------------------------------------------------------------------------- #
def main(argv: list[str] | None = None) -> int:
    p = argparse.ArgumentParser(
        prog="bcr-update",
        description="Check or update bazel_dep versions against the Bazel " "Central Registry.",
    )
    p.add_argument(
        "paths", nargs="*", default=["."], help="MODULE.bazel files or dirs to scan (default: .)"
    )
    p.add_argument(
        "-u", "--update", action="store_true", help="rewrite out-of-date versions in place"
    )
    p.add_argument(
        "--registry",
        default=DEFAULT_REGISTRY,
        help=f"registry base URL (default: {DEFAULT_REGISTRY})",
    )
    p.add_argument(
        "--include-prerelease",
        action="store_true",
        help="consider prerelease versions as upgrade candidates",
    )
    p.add_argument("--jobs", type=int, default=8, help="parallel registry requests (default: 8)")
    p.add_argument(
        "--timeout", type=float, default=15.0, help="per-request timeout in seconds (default: 15)"
    )
    p.add_argument(
        "--json", action="store_true", help="emit machine-readable JSON (implies check mode)"
    )
    p.add_argument(
        "--exit-code",
        action="store_true",
        help="in check mode, exit 1 if any upgrades are available",
    )
    p.add_argument("--no-color", action="store_true", help="disable color")
    args = p.parse_args(argv)

    c = C(enabled=sys.stdout.isatty() and not args.no_color and not args.json)

    files = find_module_files(args.paths)
    if not files:
        print("No MODULE.bazel files found.", file=sys.stderr)
        return 0

    # Parse every file; keep source text for potential edits.
    sources: dict[str, str] = {}
    all_deps: list[Dep] = []
    for f in files:
        try:
            with open(f, encoding="utf-8") as fh:
                src = fh.read()
        except OSError as e:
            print(f"warning: cannot read {f}: {e}", file=sys.stderr)
            continue
        sources[f] = src
        all_deps.extend(parse_deps(f, src))

    # Unique module names that actually carry a registry version.
    names = sorted({d.name for d in all_deps if d.version is not None})

    latest: dict[str, tuple[str, str | None]] = {}
    if names:
        with ThreadPoolExecutor(max_workers=max(1, args.jobs)) as ex:
            results = ex.map(
                lambda n: (
                    n,
                    fetch_latest(n, args.registry, args.include_prerelease, args.timeout),
                ),
                names,
            )
            latest = dict(results)

    # Decide what needs upgrading and collect problems.
    upgrades: list[tuple[Dep, str]] = []  # (dep, latest)
    problems: list[tuple[Dep, str]] = []  # (dep, reason)
    for d in all_deps:
        if d.version is None:
            continue
        status, value = latest.get(d.name, ("error", "not queried"))
        if status != "ok":
            reason = {
                "not_found": "not found in registry",
                "no_versions": "no usable versions in registry",
            }.get(status, f"lookup failed: {value}")
            problems.append((d, reason))
            continue
        if compare_version(d.version, value) < 0:
            upgrades.append((d, value))

    if args.json:
        out = {
            "upgrades": [
                {"file": d.file, "name": d.name, "current": d.version, "latest": v}
                for d, v in upgrades
            ],
            "problems": [
                {"file": d.file, "name": d.name, "current": d.version, "reason": r}
                for d, r in problems
            ],
        }
        print(json.dumps(out, indent=2))
        return 1 if (args.exit_code and upgrades) else 0

    # ---- Update mode -------------------------------------------------------
    if args.update:
        by_file: dict[str, list[tuple[Dep, str]]] = {}
        for d, v in upgrades:
            by_file.setdefault(d.file, []).append((d, v))
        changed = 0
        for f, ups in by_file.items():
            new_src = apply_updates(f, sources[f], ups)
            if new_src != sources[f]:
                with open(f, "w", encoding="utf-8") as fh:
                    fh.write(new_src)
                print(c.bold(f))
                for d, v in sorted(ups, key=lambda e: e[0].name):
                    print(f"  {c.cyan(d.name)} " f"{c.red(d.version)} -> {c.green(v)}")
                    changed += 1
        if changed:
            print(c.bold(f"\nUpdated {changed} " f"dependenc{'y' if changed == 1 else 'ies'}."))
        else:
            print(c.green("All dependencies are up to date."))
        _print_problems(problems, c)
        return 0

    # ---- Check mode --------------------------------------------------------
    if upgrades:
        by_file: dict[str, list[tuple[Dep, str]]] = {}
        for d, v in upgrades:
            by_file.setdefault(d.file, []).append((d, v))
        # Column-align the name -> version arrows per file.
        for f in sorted(by_file):
            print(c.bold(f))
            ups = sorted(by_file[f], key=lambda e: e[0].name)
            w = max(len(d.name) for d, _ in ups)
            for d, v in ups:
                print(f"  {c.cyan(d.name.ljust(w))}  " f"{c.red(d.version)} -> {c.green(v)}")
        n = len(upgrades)
        print(
            c.bold(
                f"\n{n} dependenc{'y' if n == 1 else 'ies'} can be "
                f"upgraded. Run with -u to apply."
            )
        )
    else:
        print(c.green("All dependencies are up to date."))

    _print_problems(problems, c)
    return 1 if (args.exit_code and upgrades) else 0


def _print_problems(problems: list[tuple[Dep, str]], c: C) -> None:
    if not problems:
        return
    # De-duplicate by (name, reason); a module may appear in several files.
    seen: set[tuple[str, str]] = set()
    lines: list[str] = []
    for d, r in problems:
        key = (d.name, r)
        if key not in seen:
            seen.add(key)
            lines.append(f"  {c.yellow(d.name)}: {r}")
    if lines:
        print(c.dim("\nCould not check:"))
        print("\n".join(lines))


if __name__ == "__main__":
    sys.exit(main())
