# Monorepo

This monorepo is where I store all of my code. It uses Bazel as a unified build
system across all languages and platforms.

## Projects

Go

- example-go an example go project in the monorepo.
- fib: efficient fibonacci calculations in go.
- utils-go: library of useful tools, can be depended upon in the monorepo.

Proto

- proto: declaration of all proto files.

Python

- example_python: example python project in the monorepo.

## Languages and Tools

While Bazel is language agnostic, it does require that rule sets and toolchains
be configured properly. Currently the following languages and tools are
supported:

- Go
- Python
- Protobuf and gRPC
- Docker

Languages I would like to support in the future include:
- C/C++ (Bazel has these toolchains built in, but including a C/C++ compiler in
  the dependencies enforces hermicity).
- [JavaScript/TypeScript/NodeJS](https://github.com/bazelbuild/rules_nodejs)
- [Rust](https://github.com/bazelbuild/rules_rust)
- [Java](https://github.com/bazelbuild/rules_java)
- [Scala](https://github.com/bazelbuild/rules_scala)
- [Haskell](https://github.com/tweag/rules_haskell)

Tools I would like to support in the future include:
- [K8s](https://github.com/bazelbuild/rules_k8s)
- [JSONNet](https://github.com/bazelbuild/rules_jsonnet)
- [Pandoc](https://github.com/ProdriveTechnologies/bazel-pandoc)

In the future, I would also like to bundle all supported tooling, like linters,
into Bazel itself and have Bazel lint code automatically.

### Go

Current version: 1.16.3

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

Current version: 3.8.3

Bazel is configured to install its own version of Python, to guarantee the
hermiticity/reproducibility of builds. The Python toolchain is all configured
using this toolchain. This is setup in
`//bazel/python/interpreter`. 

Bazel also manages pip installs for Python dependencies. Each Python project is
denoted by a `requirements.txt` file, which lists all dependencies for that
project. All of these dependencies are fetched as part of the WORKSPACE setup.
The dependencies are setup as a pseudo-WORKSPACE, with the project path name
used to create the Python project dependency WORKSPACE name. (See
`//bazel/python:deps.bzl` for more details. From these project WORKSPACEs, you
can use the `requirement` rule to import the dependencies.

The convention is that one top level directory is a project (if it has Python
code). An example Python project is `//bazel/example-python`.

Future work:
- `mypy` typechecking support
- pip-compile generation (automatically generate all transitive dependencies for
  `pip_install`.
    - Currently, only top level dependencies are listed in the requirements.txt
      files. Ideally, being able to automatically query and create a list of all
      transitive dependencies from top level imports would ensure more
      reproducibility and pins down package versions.
- Gazelle support
    - Gazelle has extensions for other languages. Adding a Python extension in
      Gazelle would enable us to automatically generate BUILD files for
      Python, saves on a lot of work.

### Protobuf and gRPC

While Protobuf and gRPC are technically different technologies, they are most
often used together, so we will discuss their usage together.

Currently, there is a bug with `go_proto_library`. See
[this issue](https://github.com/TerrenceHo/monorepo/issues/11) for more context.

### Docker

TODO: fill this out once `docker_rules` has been implemented.

## CI

The monorepo uses GitHub Actions as a CI mechanism. When the Action runs,

Currently, no remote cache is available to speed up execution of the dependency
graph; Bazel has to download and rebuild many artifacts it has built before,
even though Bazel was optimized for incremental rebuilds. GitHub Actions do have
a cache, but there is no way to persist the cache between builds if the cache
was used during the run. My current solution is to just have some commits that
invalidate the cache key for the bazel cache directory, have Bazel build the
[entire](entire) graph once, then subsequent builds will have a more complete cache.

Future work:
- As mentioned above, remote cache.
- Artifact extraction: when builds succeed, extract all the artifacts Bazel has
  created and upload them to some artifact storage (Artifactory is the most
  likely candidate).
- Automatic image uploading: When Bazel succeeds in building all the images in
  CI, upload them automatically to the image registry.
    - Should figure out how to avoid uploading duplicate images.

## Cache

TODO: This costs money, but could be done cheaply; either in some GCS storage
bucket or just run [`bazel-remote`](https://github.com/buchgr/bazel-remote) on
some machine somewhere.

## Remote Execution

TODO: This costs money.

## Artifact Storage

TODO: Implement an artifact storage for all dependencies that Bazel relies
on (including `http_archive` and `git_repository` dependencies).
