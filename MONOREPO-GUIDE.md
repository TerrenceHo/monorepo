# Polyglot Bazel Monorepo: A Modern Architecture Guide

## Guiding Principles

A polyglot monorepo in Bazel succeeds or fails on three axes: hermetic builds (every action runs identically everywhere), a single version policy (one version of every external dependency, enforced repo-wide), and Bzlmod-first configuration (MODULE.bazel replaces the legacy WORKSPACE file entirely). Everything below flows from those three ideas.

---

## 1. Repository Layout

```
/
├── MODULE.bazel              # Central dependency registry (Bzlmod)
├── MODULE.bazel.lock          # Lockfile — commit this
├── .bazelrc                   # Build flags, platform configs, remote cache
├── .bazelversion              # Pin Bazel version (e.g. 8.x)
├── BUILD.bazel                # Root BUILD (usually empty or alias-only)
├── toolchains/
│   └── BUILD.bazel            # Custom toolchain registrations
├── platforms/
│   └── BUILD.bazel            # Platform definitions (linux-x86, darwin-arm64…)
├── proto/
│   ├── BUILD.bazel
│   └── api/
│       ├── v1/
│       │   ├── BUILD.bazel
│       │   └── service.proto
│       └── types/
│           ├── BUILD.bazel
│           └── common.proto
├── libs/                      # Shared internal libraries, per-language
│   ├── py/
│   │   └── logging/
│   │       ├── BUILD.bazel
│   │       └── logger.py
│   ├── go/
│   │   └── httputil/
│   │       ├── BUILD.bazel
│   │       └── client.go
│   ├── rs/
│   │   └── config/
│   │       ├── BUILD.bazel
│   │       └── lib.rs
│   ├── cc/
│   │   └── base/
│   │       ├── BUILD.bazel
│   │       └── status.h
│   ├── ts/
│   │   └── sdk/
│   │       ├── BUILD.bazel
│   │       └── index.ts
│   └── ml/
│       └── codec/
│           ├── BUILD.bazel
│           └── codec.ml
├── services/
│   ├── gateway/               # Example Go service
│   │   ├── BUILD.bazel
│   │   ├── main.go
│   │   └── gateway_test.go
│   ├── inference/             # Example Python service
│   │   ├── BUILD.bazel
│   │   ├── server.py
│   │   └── requirements.in
│   ├── ingest/                # Example Rust service
│   │   ├── BUILD.bazel
│   │   └── main.rs
│   └── renderer/              # Example C++ service
│       ├── BUILD.bazel
│       └── main.cc
├── apps/
│   └── dashboard/             # Example TypeScript frontend
│       ├── BUILD.bazel
│       ├── package.json
│       └── src/
├── tools/                     # Developer tooling, custom rules, macros
│   ├── bazel/
│   │   ├── macros/
│   │   │   └── service.bzl    # Macro: "service bundle" (binary + image + deploy)
│   │   └── transitions/
│   └── scripts/
│       └── refresh_deps.sh
├── third_party/               # Vendored or manually managed deps
│   └── BUILD.bazel
└── gazelle_setup/
    ├── BUILD.bazel
    └── extensions.bzl         # Gazelle language extension wiring
```

The key structural decisions:

- `proto/` is a standalone tree. Every language generates bindings *from* it; no language *owns* it. This prevents proto files from getting buried inside a Go or Python package and becoming invisible to other languages.
- `libs/` holds shared internal code split by language. This avoids the "utils" junk-drawer. Each sub-directory is a proper Bazel package with a clear API surface.
- `services/` and `apps/` hold deployable units. A service can depend on anything in `libs/` or `proto/` but never on another service directly — communicate through proto-defined APIs.

---

## 2. MODULE.bazel — The Central Dependency Registry

```starlark
module(
    name = "mymonorepo",
    version = "0.0.0",
)

# ─── Bazel core ───
bazel_dep(name = "bazel_skylib", version = "1.7.1")
bazel_dep(name = "platforms", version = "0.0.10")

# ─── Protobuf ───
bazel_dep(name = "protobuf", version = "29.3")
bazel_dep(name = "rules_proto", version = "7.1.0")
bazel_dep(name = "grpc", version = "1.70.1")  # if using gRPC

# ─── Go ───
bazel_dep(name = "rules_go", version = "0.53.0")
bazel_dep(name = "gazelle", version = "0.42.0")

go_deps = use_extension("@gazelle//:extensions.bzl", "go_deps")
go_deps.from_file(go_mod = "//:go.mod")
use_repo(go_deps, "com_github_some_dependency", ...)

# ─── Python ───
bazel_dep(name = "rules_python", version = "1.4.1")

python = use_extension("@rules_python//python/extensions:python.bzl", "python")
python.toolchain(python_version = "3.12")

pip = use_extension("@rules_python//python/extensions:pip.bzl", "pip")
pip.parse(
    hub_name = "pip",
    python_version = "3.12",
    requirements_lock = "//:requirements_lock.txt",
)
use_repo(pip, "pip")

# ─── Rust ───
bazel_dep(name = "rules_rust", version = "0.59.2")

crate = use_extension("@rules_rust//crate_universe:extension.bzl", "crate")
crate.from_cargo(
    name = "crate_index",
    cargo_lockfile = "//:Cargo.lock",
    manifests = [
        "//:Cargo.toml",
        "//services/ingest:Cargo.toml",
    ],
)
use_repo(crate, "crate_index")

# ─── C/C++ ───
# rules_cc is bundled with Bazel, but pin explicitly for Bzlmod:
bazel_dep(name = "rules_cc", version = "0.1.1")
# For external C/C++ deps, use a module extension or third_party/

# ─── TypeScript / JavaScript ───
bazel_dep(name = "aspect_rules_js", version = "2.3.7")
bazel_dep(name = "aspect_rules_ts", version = "3.6.0")

npm = use_extension("@aspect_rules_js//npm:extensions.bzl", "npm")
npm.npm_translate_lock(
    name = "npm",
    pnpm_lock = "//:pnpm-lock.yaml",
)
use_repo(npm, "npm")

# ─── OCaml ───
bazel_dep(name = "rules_ocaml", version = "2.3.5")

# ─── Gazelle ───
# Already declared above with gazelle. Gazelle extensions for
# non-Go languages are configured below in gazelle_setup/.

# ─── Remote execution (optional) ───
bazel_dep(name = "toolchains_llvm", version = "1.4.0")  # hermetic CC toolchain
```

### Why this matters

Every external dependency is declared once, here. No `WORKSPACE` file. No `http_archive` scattered through the tree. The lockfile (`MODULE.bazel.lock`) is committed and ensures determinism.

---

## 3. Language-by-Language Deep Dive

### 3.1 Go

**Ruleset:** `rules_go` + `gazelle`

**Dependency management:** A single `/go.mod` at the repo root. Gazelle's `go_deps` module extension reads it and creates repos for every transitive dependency. Run `go mod tidy` normally, then `bazel mod tidy` to sync.

**BUILD file generation:** Gazelle is the gold standard here — it was literally designed for Go. Run:
```bash
bazel run @gazelle//:gazelle
```
It auto-generates `go_library`, `go_binary`, `go_test` targets from your source tree. Configure it via directives in BUILD files:
```starlark
# gazelle:prefix github.com/myorg/mymonorepo
# gazelle:resolve go github.com/myorg/mymonorepo/libs/go/httputil //libs/go/httputil
```

**Central libraries pattern:**
```starlark
# libs/go/httputil/BUILD.bazel
load("@rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "httputil",
    srcs = ["client.go"],
    importpath = "github.com/myorg/mymonorepo/libs/go/httputil",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_hashicorp_go_retryablehttp//:go-retryablehttp",
    ],
)

go_test(
    name = "httputil_test",
    srcs = ["client_test.go"],
    embed = [":httputil"],
)
```

**Proto integration:**
```starlark
load("@rules_go//proto:def.bzl", "go_proto_library")

go_proto_library(
    name = "api_v1_go_proto",
    importpath = "github.com/myorg/mymonorepo/proto/api/v1",
    proto = "//proto/api/v1:service_proto",
    visibility = ["//visibility:public"],
    deps = ["//proto/api/types:common_go_proto"],
)
```

**IDE integration:** Use `gopls`. It works with Bazel if you maintain the `go.mod` at the root — `gopls` reads that while Bazel reads the generated repos. Keep them in sync. For tighter integration, the `gopackagesdriver` from rules_go lets `gopls` query Bazel directly:
```bash
export GOPACKAGESDRIVER=bazel run @rules_go//go/tools/gopackagesdriver -- 
```

---

### 3.2 Python

**Ruleset:** `rules_python`

**Dependency management:** Maintain a `requirements.in` at the root (or per-service if isolation is needed). Compile it to a lock file:
```bash
bazel run @rules_python//tools:pip_compile -- \
    --requirements_file=requirements.in \
    --output_file=requirements_lock.txt
```
The `pip.parse` extension in MODULE.bazel creates a `@pip` hub repo. Reference deps as `@pip//pypi__numpy` (or the generated alias).

**BUILD file generation:** Gazelle has a Python extension (`gazelle_python`):
```starlark
# gazelle_setup/BUILD.bazel
load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle_bin",
    languages = [
        "@rules_python_gazelle_plugin//python:python",
        "@gazelle//language/go",       # keep Go support too
        "@gazelle//language/proto",
    ],
)

gazelle(
    name = "gazelle",
    gazelle = ":gazelle_bin",
)
```

**Central libraries pattern:**
```starlark
# libs/py/logging/BUILD.bazel
load("@rules_python//python:defs.bzl", "py_library")

py_library(
    name = "logging",
    srcs = ["logger.py"],
    visibility = ["//visibility:public"],
    deps = [
        "@pip//structlog",
    ],
)
```

**Proto integration:**
```starlark
load("@rules_python//python:proto.bzl", "py_proto_library")

py_proto_library(
    name = "api_v1_py_proto",
    deps = ["//proto/api/v1:service_proto"],
    visibility = ["//visibility:public"],
)
```

**IDE integration:** Use a `.pth` file or configure your virtualenv to include bazel-bin paths, or just maintain a parallel `venv` from the same `requirements_lock.txt`. Most teams do:
```bash
python -m venv .venv
pip install -r requirements_lock.txt
```
Point your IDE (VS Code / PyCharm) at `.venv`. The source tree is the same; only external deps need the venv. For more advanced setups, `rules_python` can generate a runfiles-based IDE helper.

---

### 3.3 C / C++

**Ruleset:** `rules_cc` (bundled, but pin via Bzlmod)

**Dependency management:** C/C++ has no standard package manager, so Bazel fills the gap. Options:

1. **Bzlmod registries** — if your deps publish Bazel modules (abseil, googletest, protobuf, grpc, boringssl all do), just `bazel_dep()` them.
2. **`http_archive` via module extensions** — for deps without Bazel module support, write a small module extension that fetches and patches them.
3. **`third_party/` vendoring** — for small or forked deps, vendor the source and write BUILD files by hand.

For a hermetic C++ toolchain (critical for remote execution), use `toolchains_llvm`:
```starlark
llvm = use_extension("@toolchains_llvm//toolchain/extensions:llvm.bzl", "llvm")
llvm.toolchain(
    name = "llvm_toolchain",
    llvm_version = "18.1.8",
)
use_repo(llvm, "llvm_toolchain")
register_toolchains("@llvm_toolchain//:all")
```

**Central libraries pattern:**
```starlark
# libs/cc/base/BUILD.bazel
load("@rules_cc//cc:defs.bzl", "cc_library")

cc_library(
    name = "status",
    hdrs = ["status.h"],
    srcs = ["status.cc"],
    visibility = ["//visibility:public"],
    deps = [
        "@abseil-cpp//absl/status",
        "@abseil-cpp//absl/strings",
    ],
)
```

**Proto integration:** This is native to protobuf's own Bazel support:
```starlark
load("@protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")

cc_proto_library(
    name = "api_v1_cc_proto",
    deps = ["//proto/api/v1:service_proto"],
    visibility = ["//visibility:public"],
)
```

**Gazelle:** There's no mature Gazelle extension for C++. Most teams hand-write BUILD files or use `buildifier` + custom scripts. For large C++ codebases, some teams write a custom Gazelle extension that scans `#include` directives.

**IDE integration:** Generate a `compile_commands.json`:
```bash
# Use the hedronvision extractor (most popular approach)
bazel_dep(name = "hedron_compile_commands", version = "...")
# Then:
bazel run @hedron_compile_commands//:refresh_all
```
Point `clangd` or VS Code's C++ extension at the generated `compile_commands.json`. This gives full autocomplete, go-to-definition, and diagnostics.

---

### 3.4 Rust

**Ruleset:** `rules_rust` + `crate_universe`

**Dependency management:** Maintain a workspace `Cargo.toml` at the root that lists all crates as path members. `crate_universe` reads `Cargo.toml` + `Cargo.lock` and generates Bazel repos for every crate:

```toml
# /Cargo.toml (workspace root)
[workspace]
members = [
    "services/ingest",
    "libs/rs/config",
]

[workspace.dependencies]
tokio = { version = "1.43", features = ["full"] }
serde = { version = "1", features = ["derive"] }
prost = "0.13"
```

The `crate.from_cargo` extension in MODULE.bazel handles the rest. Repin with:
```bash
CARGO_BAZEL_REPIN=1 bazel sync --only=crate_index
```

**Central libraries pattern:**
```starlark
# libs/rs/config/BUILD.bazel
load("@rules_rust//rust:defs.bzl", "rust_library", "rust_test")

rust_library(
    name = "config",
    srcs = ["lib.rs"],
    visibility = ["//visibility:public"],
    deps = [
        "@crate_index//:serde",
        "@crate_index//:toml",
    ],
)

rust_test(
    name = "config_test",
    crate = ":config",
)
```

**Proto integration:** Use `prost` (the Rust ecosystem standard) via rules_rust's proto support:
```starlark
load("@rules_rust//proto/prost:defs.bzl", "rust_prost_library")

rust_prost_library(
    name = "api_v1_rs_proto",
    proto = "//proto/api/v1:service_proto",
    visibility = ["//visibility:public"],
)
```

**Gazelle:** There's a community Gazelle extension for Rust (`nicholasgasior/gazelle_rust` or similar), but it's less mature than Go's. Many Rust-in-Bazel teams hand-write BUILD files — rust_library targets map 1:1 to crate roots, so there's less boilerplate than C++.

**IDE integration:** `rust-analyzer` needs a `rust-project.json`. Generate it:
```bash
bazel run @rules_rust//tools/rust_analyzer:gen_rust_project
```
This produces a `rust-project.json` at the workspace root. Point rust-analyzer at it in your IDE settings. You get full type inference, completions, and macro expansion, all resolving through Bazel's dependency graph.

---

### 3.5 TypeScript / JavaScript

**Ruleset:** `aspect_rules_js` + `aspect_rules_ts`

**Dependency management:** Use `pnpm` as the package manager. Maintain a single `pnpm-lock.yaml` at the root. The `npm_translate_lock` extension creates a `@npm` repo from the lockfile. For a multi-package workspace, each package has its own `package.json`, but the lockfile is unified:

```
/
├── package.json            # Root workspace
├── pnpm-lock.yaml
├── pnpm-workspace.yaml
├── apps/dashboard/package.json
└── libs/ts/sdk/package.json
```

**BUILD file generation:** Gazelle has a JS/TS extension in `aspect_rules_js`:
```bash
bazel run @aspect_rules_js//js/gazelle -- fix
```

**Central libraries pattern:**
```starlark
# libs/ts/sdk/BUILD.bazel
load("@aspect_rules_ts//ts:defs.bzl", "ts_project")

ts_project(
    name = "sdk",
    srcs = glob(["src/**/*.ts"]),
    declaration = True,
    tsconfig = ":tsconfig.json",
    visibility = ["//visibility:public"],
    deps = [
        "@npm//zod",
        "@npm//axios",
    ],
)
```

**Proto integration:** Use `@aspect_rules_ts`'s ts_proto_library, or use `buf` with `connect-es`:
```starlark
load("@aspect_rules_ts//ts:proto.bzl", "ts_proto_library")

ts_proto_library(
    name = "api_v1_ts_proto",
    proto = "//proto/api/v1:service_proto",
    visibility = ["//visibility:public"],
)
```

Many teams prefer a `buf generate` step wrapped in a `genrule` for more control over the TypeScript protobuf plugin (protobuf-es, connect-es, ts-proto, etc.).

**IDE integration:** TypeScript IDE support is straightforward because `ts_project` uses the standard `tsc` compiler and reads `tsconfig.json` normally. Ensure your `tsconfig.json` paths resolve correctly:
```json
{
  "compilerOptions": {
    "paths": {
      "@myorg/sdk": ["libs/ts/sdk/src"],
      "@myorg/proto/*": ["bazel-bin/proto/*"]
    }
  }
}
```
VS Code's TypeScript server picks this up natively. For generated proto types, you may need to run an initial build so the `.d.ts` files exist in `bazel-bin/`.

---

### 3.6 OCaml

**Ruleset:** `rules_ocaml`

OCaml in Bazel is the least mainstream of these languages, but `rules_ocaml` (v2+) is functional and actively maintained.

**Dependency management:** `rules_ocaml` can work with `opam` packages. Use a module extension to resolve opam deps:
```starlark
# In MODULE.bazel
bazel_dep(name = "rules_ocaml", version = "2.3.5")

ocaml = use_extension("@rules_ocaml//extensions:ocaml.bzl", "ocaml")
ocaml.toolchain(version = "5.2.1")

opam = use_extension("@rules_ocaml//extensions:opam.bzl", "opam")
opam.dependency(name = "core", version = "0.17.1")
opam.dependency(name = "yojson", version = "2.2.2")
opam.dependency(name = "lwt", version = "5.8.0")
```

Alternatively, vendor opam packages into `third_party/ocaml/` and write BUILD files. This gives full hermeticity.

**Central libraries pattern:**
```starlark
# libs/ml/codec/BUILD.bazel
load("@rules_ocaml//ocaml:defs.bzl", "ocaml_library")

ocaml_library(
    name = "codec",
    srcs = ["codec.ml"],
    visibility = ["//visibility:public"],
    deps = [
        "@opam//yojson",
    ],
)
```

**Proto integration:** There's no first-class `ocaml_proto_library` rule. The standard approach is to use `ocaml-protoc` (or `ocaml-protoc-plugin`) wrapped in a custom rule or `genrule`:
```starlark
genrule(
    name = "gen_ocaml_proto",
    srcs = ["//proto/api/v1:service.proto"],
    outs = ["service_pb.ml", "service_pb.mli"],
    cmd = "$(execpath @ocaml-protoc//:ocaml-protoc) -ml_out=$(RULEDIR) $<",
    tools = ["@ocaml-protoc//:ocaml-protoc"],
)

ocaml_library(
    name = "api_v1_ml_proto",
    srcs = [":gen_ocaml_proto"],
    deps = ["@opam//ocaml-protoc-plugin"],
    visibility = ["//visibility:public"],
)
```

**Gazelle:** No Gazelle extension for OCaml exists. BUILD files are hand-written. OCaml projects tend to have clean module structures (one `.ml` file = one module), so the mapping is direct.

**IDE integration:** Use `ocaml-lsp-server` (via `merlin`). The key is generating a `.merlin` or `dune-project`-like config that points to Bazel's output paths. Some teams maintain a lightweight `dune` setup in parallel purely for IDE support, while the CI uses Bazel. Alternatively, add `bazel-bin` paths to your `.merlin` file:
```
S libs/ml/codec
B bazel-bin/libs/ml/codec
PKG yojson lwt
```

---

## 4. Protobuf: The Cross-Language Spine

The proto tree is the lingua franca. Structure it independently from any language:

```starlark
# proto/api/v1/BUILD.bazel
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "service_proto",
    srcs = ["service.proto"],
    visibility = ["//visibility:public"],
    deps = [
        "//proto/api/types:common_proto",
        "@protobuf//:timestamp_proto",
    ],
)
```

Each language then wraps this with its own `*_proto_library` rule (shown in sections above). The proto_library target is the single source of truth; language bindings are derived artifacts.

**Buf integration:** If you use Buf for linting and breaking-change detection:
```starlark
# proto/BUILD.bazel — run as: bazel run //proto:buf_lint
genrule(
    name = "buf_lint",
    srcs = glob(["**/*.proto"]),
    outs = ["buf_lint.log"],
    cmd = "$(execpath @buf//:buf) lint $(SRCS) > $@",
    tools = ["@buf//:buf"],
)
```

Or more practically, run `buf lint` and `buf breaking` as test targets or in CI outside Bazel (Buf's own CLI is fast enough that wrapping it in Bazel adds little value).

**gRPC services:** For each language that needs gRPC stubs, use the corresponding `*_grpc_library` rule. The gRPC Bazel module provides `cc_grpc_library`; rules_go provides `go_grpc_library`; Python uses `grpcio-tools` through `py_proto_library` with a gRPC plugin; etc.

---

## 5. Gazelle Configuration for Polyglot Repos

Gazelle supports multiple languages in a single run. Wire them together:

```starlark
# /BUILD.bazel
load("@gazelle//:def.bzl", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle_bin",
    languages = [
        "@gazelle//language/go",
        "@gazelle//language/proto",
        "@rules_python_gazelle_plugin//python:python",
        # Add community extensions as they mature:
        # "@aspect_rules_js//js/gazelle:js",
    ],
)

gazelle(
    name = "gazelle",
    gazelle = ":gazelle_bin",
    args = ["-go_prefix=github.com/myorg/mymonorepo"],
)

# Separate command for fix mode
gazelle(
    name = "gazelle_fix",
    gazelle = ":gazelle_bin",
    args = ["-go_prefix=github.com/myorg/mymonorepo", "-mode=fix"],
)
```

Directives in BUILD files control per-directory behavior:
```starlark
# gazelle:exclude vendor
# gazelle:python_root
# gazelle:go_naming_convention import
# gazelle:proto_group go_package
```

For languages without Gazelle support (C++, Rust, OCaml), use `buildifier` to keep BUILD files consistently formatted and add a CI check that BUILD files are up to date.

---

## 6. Remote Cache and Remote Execution

### .bazelrc Configuration

```bash
# ─── .bazelrc ───

# ─ Common ─
common --enable_bzlmod
build --incompatible_strict_action_env
build --java_runtime_version=remotejdk_21

# ─ Remote Cache (BuildBuddy example) ─
build:remote-cache --remote_cache=grpcs://remote.buildbuddy.io
build:remote-cache --remote_header=x-buildbuddy-api-key=YOUR_KEY
build:remote-cache --remote_upload_local_results=true
build:remote-cache --remote_timeout=3600

# ─ Remote Execution ─
build:remote-exec --remote_executor=grpcs://remote.buildbuddy.io
build:remote-exec --remote_header=x-buildbuddy-api-key=YOUR_KEY
build:remote-exec --jobs=200
build:remote-exec --remote_instance_name=myorg/default

# Platform flags for remote execution
build:remote-exec --extra_execution_platforms=//platforms:linux_x86_64
build:remote-exec --host_platform=//platforms:linux_x86_64

# ─ CI ─
build:ci --config=remote-cache
build:ci --bes_results_url=https://app.buildbuddy.io/invocation/
build:ci --bes_backend=grpcs://remote.buildbuddy.io
build:ci --build_metadata=ROLE=CI

# ─ Developer ─
build:dev --config=remote-cache
build:dev --remote_upload_local_results=false  # devs read cache, don't pollute it
```

### Platform Definitions

```starlark
# platforms/BUILD.bazel
platform(
    name = "linux_x86_64",
    constraint_values = [
        "@platforms//os:linux",
        "@platforms//cpu:x86_64",
    ],
    exec_properties = {
        "OSFamily": "Linux",
        "container-image": "docker://gcr.io/myorg/bazel-remote:latest",
    },
)
```

### Cache Hygiene Tips

Remote caching gives you the biggest bang-for-buck improvement (10x+ CI speedup for most repos). Remote execution is harder to set up but eliminates local machine variance entirely.

Key practices: keep all toolchains hermetic (use `toolchains_llvm` for C++, managed toolchains for Go/Python/Rust), never depend on things installed on the host, and tag non-hermetic actions with `no-remote` so they don't poison the cache.

---

## 7. IDE Integration Summary

| Language   | IDE/Editor        | Mechanism                                              |
|------------|-------------------|--------------------------------------------------------|
| Go         | VS Code / GoLand  | `gopls` + `gopackagesdriver` from rules_go             |
| Python     | VS Code / PyCharm | Parallel venv from same `requirements_lock.txt`        |
| C/C++      | VS Code / CLion   | `compile_commands.json` via hedron extractor → clangd  |
| Rust       | VS Code / RustRover| `rust-project.json` via `gen_rust_project`            |
| TypeScript | VS Code           | Native `tsconfig.json` paths; build once for gen'd types|
| OCaml      | VS Code + ocamllsp| `.merlin` pointing at `bazel-bin/` paths               |

The universal principle: Bazel is the build system of record, but IDEs need files on disk in expected locations. Each language has a different bridge — generate the IDE config artifact from Bazel's dependency graph, and you get the best of both worlds.

---

## 8. Useful Custom Macros

**Service bundle macro** — wraps a binary, container image, and deploy config:

```starlark
# tools/bazel/macros/service.bzl
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_push")

def service_bundle(name, binary, base_image = "@distroless_base", **kwargs):
    """Bundles a service binary into a container image with a push target."""
    oci_image(
        name = name + "_image",
        base = base_image,
        entrypoint = ["/app/" + name],
        tars = [binary + "_layer"],
        **kwargs
    )

    oci_push(
        name = name + "_push",
        image = ":" + name + "_image",
        repository = "gcr.io/myorg/" + name,
    )
```

**Multi-language proto macro** — generates bindings for all languages at once:

```starlark
# tools/bazel/macros/proto.bzl
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@rules_go//proto:def.bzl", "go_proto_library")
load("@rules_python//python:proto.bzl", "py_proto_library")
load("@protobuf//bazel:cc_proto_library.bzl", "cc_proto_library")

def polyglot_proto_library(name, srcs, deps = [], **kwargs):
    """Generates proto_library + Go, Python, and C++ bindings."""
    proto_library(
        name = name + "_proto",
        srcs = srcs,
        deps = deps,
        visibility = ["//visibility:public"],
    )

    go_proto_library(
        name = name + "_go_proto",
        proto = ":" + name + "_proto",
        importpath = native.package_name().replace("/", "."),
        visibility = ["//visibility:public"],
    )

    py_proto_library(
        name = name + "_py_proto",
        deps = [":" + name + "_proto"],
        visibility = ["//visibility:public"],
    )

    cc_proto_library(
        name = name + "_cc_proto",
        deps = [":" + name + "_proto"],
        visibility = ["//visibility:public"],
    )
```

---

## 9. CI / Developer Workflow Cheat Sheet

```bash
# ─── Day-to-day commands ───

# Build everything
bazel build //...

# Test everything
bazel test //... --config=dev

# Regenerate BUILD files (Go + Python + Proto)
bazel run //:gazelle

# Update Go deps after editing go.mod
go mod tidy && bazel mod tidy

# Update Python deps after editing requirements.in
bazel run @rules_python//tools:pip_compile -- ...
bazel mod tidy

# Repin Rust crates after editing Cargo.toml
CARGO_BAZEL_REPIN=1 bazel sync --only=crate_index

# Update JS/TS deps
pnpm install --lockfile-only
bazel mod tidy

# Refresh IDE configs
bazel run @hedron_compile_commands//:refresh_all          # C++
bazel run @rules_rust//tools/rust_analyzer:gen_rust_project # Rust

# Query the dependency graph
bazel query 'deps(//services/gateway:gateway)' --output=graph | dot -Tpng > deps.png

# Find what depends on a library (reverse deps)
bazel query 'rdeps(//..., //libs/go/httputil:httputil)'

# CI build
bazel test //... --config=ci
```

---

## 10. Migration & Scaling Advice

**Start with the cache.** Before anything else, get remote caching working (`--config=remote-cache`). Even with zero other changes, this slashes CI times by caching unchanged targets across branches and developers.

**Adopt Gazelle early.** The biggest maintenance burden in a monorepo is BUILD file hygiene. Every language that has a Gazelle extension should use it from day one. Lint CI against `gazelle --mode=diff` to catch drift.

**Enforce the single version policy.** Never allow two versions of the same dependency in the repo. Bzlmod helps — it resolves to a single version via MVS (minimum version selection) for each module. For pip and npm, the lockfile is your enforcement mechanism.

**Use `bazel mod tidy` religiously.** It updates your `use_repo` statements after any MODULE.bazel change. Treat a dirty `bazel mod tidy` diff the same as a dirty `go mod tidy` — fail CI on it.

**Make proto the contract.** If service A talks to service B, the proto in `proto/` is the only thing they share. Services never import each other's Go/Python/Rust code directly. This keeps the dependency graph clean and prevents circular dependencies from sneaking in as the repo grows.
