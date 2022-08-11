# Monorepo

This monorepo is where I store all of my code. It uses Bazel as a unified build
system across all languages and platforms.

## Projects

Go

- `example-go`: an example go project in the monorepo.
- `fastlinks`: a golinks implementation, but for local daemon usage.
- `fib`: efficient fibonacci calculations in go.
- `utils-go`: library of useful tools, can be depended upon in the monorepo.

Python

- `//example_python`: example python project in the monorepo.

## Languages and Tools

While Bazel is language agnostic, it does require that rule sets and toolchains
be configured properly. Currently the following languages and tools are
supported:

- Go
- Python
- Gazelle (for BUILD file generation)
- [Docker](https://github.com/bazelbuild/rules_docker) for image building.


Languages I would like to support in the future include:
- C/C++ (Bazel has these toolchains built in, but including a C/C++ compiler in
  the dependencies enforces hermicity).
  - A possible hermetic replacement for C++ would be
    [zig](https://andrewkelley.me/post/zig-cc-powerful-drop-in-replacement-gcc-clang.html),
    with a setup like [this](https://sr.ht/~motiejus/bazel-zig-cc/).
- [JavaScript/TypeScript/NodeJS](https://github.com/bazelbuild/rules_nodejs)
- [Rust](https://github.com/bazelbuild/rules_rust)
- [Java](https://github.com/bazelbuild/rules_java)
- [Scala](https://github.com/bazelbuild/rules_scala)
- [Haskell](https://github.com/tweag/rules_haskell)

Tools I would like to support in the future include:
- [Protobuf](https://github.com/bazelbuild/rules_proto)
- [gRPC](https://github.com/rules-proto-grpc/rules_proto_grpc)
- [K8s](https://github.com/bazelbuild/rules_k8s)
- [JSONNet](https://github.com/bazelbuild/rules_jsonnet)
- [Pandoc](https://github.com/ProdriveTechnologies/bazel-pandoc)

In the future, I would also like to bundle all supported tooling, like linters,
into Bazel itself and have Bazel lint code automatically.

### Go

Current version: 1.18.3

#### Go Dependencies

All go dependencies for all projects are tracked by a single root `go.mod` and
`go.sum` file. If a dependency is needed or updated, running `go get <module>`
will update this and update the dependencies for all projects. This is to ensure
consistent library upgrades across all projects.

#### Go Gazelle
[Gazelle](https://github.com/bazelbuild/bazel-gazelle) automatically generates
BUILD files for Go, Protobuf, and (with extensions) other languages. To help Go
developers, `gazelle` has been made available as rule in the root BUILD file.
There are also several convenience commands available to manage BUILD files.

Commands:
- `bazel run //:gazelle`: Run the `gazelle` itself.
- `bazel run //:update_build_files`: Update all BUILD files that gazelle
  manages.
- `bazel run //:check_build_files`: Dry to of which BUILD files will be updated.
- `bazel run //:update_go_deps`: Updates all dependencies in
  `//bazel/go/deps.bzl%fetch_go_deps`.
  
All go dependencies are listed under `//bazel/go/deps.bzl`, where `gazelle` will
automatically update the macro with new dependencies when they are detected.
`gazelle` will get its dependencies from the root `go.mod` and `go.sum` files;
they let `gazelle` know which exact dependencies to get. This does mean that
developers will have to update the go.mod in the root. You can do this manually,
but its probably best to just use the regular `go get package/path...` with go
modules and allow the go tool to update all checksums and transitive
dependencies.

Future Work:
- Substitute commands from `rules_go` for custom macro implementations. Allows
  more future flexibility.
- Think about pregenerated code:
    - Check into repository? Less hermetic (if generated code isn't kept up to
      date), but more developer friendly.
    - Let Bazel generate with genrules and such? More hermetic, but doesn't
      allow developers to see generated files, IDE/LSP concerns.

### Python

Current version: 3.10.4

This Bazel Python setup installs the
[toolchain](https://github.com/bazelbuild/rules_python#toolchain-registration)
native to `rules_python`. This Python toolchain should be relatively hermetic,
subject to the quirks of
[python-build-standalone](https://github.com/indygreg/python-build-standalone). 

Bazel also manages pip installs for Python dependencies. Each Python project is
denoted by a `requirements.txt` file, which lists all dependencies for that
project. All of these dependencies are fetched as part of the WORKSPACE setup.
The dependencies are setup as a pseudo-WORKSPACE, with the project path name
used to create the Python project dependency WORKSPACE name. (See
`//bazel/python:deps.bzl` for more details. From these project WORKSPACEs, you
can use the `requirement` rule to import the dependencies.

#### Python Dependencies

Any top level dependencies needed are added in a `requirements.txt` file. Then
create a `compile_pip_requirements` Bazel rule that generates a requirements
lock file.

Python build dependencies are downloaded/loaded lazily using
[pip_parse](https://github.com/bazelbuild/rules_python/blob/main/docs/pip.md#pip_parse).
To enable this in a monorepo Python project, add your requirements lock file to
`//bazel/python/deps.bzl` and add the resulting pip dependency workspace to
`//bazel/python/load_pip_repositories.bzl`.

This means that unlike go dependency management, each individual Python project
can have its own list of dependencies. May move to a single shared requirements
in the future.

#### Python Gazelle

Gazelle support has been added to manage your BUILD.bazel files automatically.
Add a `gazelle_python.yaml` file to the project directory and follow the docs
for the [Python gazelle
extension](https://github.com/bazelbuild/rules_python/tree/main/gazelle).

If you need to import a monorepo dependency, make sure that dependency has a 
`#gazelle:python_root` directive somewhere, indicating it is a Python root. This
sets up your imports correctly.

The convention is that one top level directory is a project (if it has Python
code). An example Python project is `//bazel/example-python`.

Future work:
- `mypy` typechecking support

### Docker

TODO: fill this out once `docker_rules` has been implemented.

## CI

Currently, no remote cache is available to speed up execution of the dependency
graph; instead, the outputs of the Bazel tree are saved across builds; there are
two caches:
- Master cache: After a commit is landed into master, the cache is saved with
  the build output of master.
- Branch cache: Each branch can save its own cache. 

When a PR/branch is first created, it doesn't have a cache, and so CI
builds/tests would be slow. To remedy this, the CI can read the master cache to
bootstrap its cache, with the idea that most precomputed artifacts can be used,
and the tests should only have to build what is new in the PR. After a PR CI is
finished, it will save the Bazel cache with its own cache name, that only it can
read.

Future work:
- As mentioned above, remote cache.
- Artifact extraction: when builds succeed, extract all the artifacts Bazel has
  created and upload them to some artifact storage (Artifactory is the most
  likely candidate).
- Automatic image uploading: When Bazel succeeds in building all the images in
  CI, upload them automatically to the image registry.
    - Should figure out how to avoid uploading duplicate images.

## Cache

Cache for CI is saved on github actions using `@TerrenceHo/cache-always` (a fork
of GitHubs `cache-always`), which does not check to see if a cache hit was used
during saving of caches, which allows the bazel cache to always be up to date.

## Remote Execution

TODO: This costs money.

## Artifact Storage

TODO: Implement an artifact storage for all dependencies that Bazel relies
on (including `http_archive` and `git_repository` dependencies).
