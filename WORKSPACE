workspace(
    name = "com_github_terrenceho_monorepo",
)

load("//bazel:workspace_deps.bzl", "fetch_deps")

fetch_deps()

# Needed for skylark unit testing
load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

##### Go Dependencies

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.19")

load("@io_bazel_rules_go//extras:embed_data_deps.bzl", "go_embed_data_dependencies")

go_embed_data_dependencies()

load("//bazel/go:deps.bzl", "fetch_go_deps")

# gazelle:repository_macro bazel/go/deps.bzl%fetch_go_deps
fetch_go_deps()

## Golang Gazelle Setup
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

##### Python Dependencies

load("@rules_python//python:repositories.bzl", "python_register_toolchains")

python_register_toolchains(
    name = "python3_10",
    python_version = "3.10",
)

load("@rules_python//python/pip_install:repositories.bzl", "pip_install_dependencies")

# NOTE: Because `compile_pip_dependencies()` does not install its own
# dependencies, this installs the work around for us.
# Ref: https://github.com/bazelbuild/rules_python/issues/497
pip_install_dependencies()

load("//bazel/python:deps.bzl", "setup_pip_repositories")

setup_pip_repositories()

load("//bazel/python:load_pip_repositories.bzl", "load_pip_repositories")

load_pip_repositories()

## Python Gazelle Setup
load("@rules_python//gazelle:deps.bzl", _py_gazelle_deps = "gazelle_deps")

_py_gazelle_deps()

##### GCC Toolchain Setup

load("@aspect_gcc_toolchain//toolchain:repositories.bzl", "gcc_toolchain_dependencies")

gcc_toolchain_dependencies()

load("@aspect_gcc_toolchain//toolchain:defs.bzl", "ARCHS", "gcc_register_toolchain")

gcc_register_toolchain(
    name = "gcc_toolchain_x86_64",
    sysroot_variant = "x86_64-X11",
    target_arch = ARCHS.x86_64,
)

load("@rules_foreign_cc//foreign_cc:repositories.bzl", "rules_foreign_cc_dependencies")

rules_foreign_cc_dependencies()

##### Docker Dependencies
load(
    "@io_bazel_rules_docker//toolchains/docker:toolchain.bzl",
    docker_toolchain_configure = "toolchain_configure",
)

docker_toolchain_configure(
    name = "docker_config",
    docker_flags = [
        # "--tls",
        "--log-level=info",
    ],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

load("//bazel/containers:images.bzl", "fetch_images")

fetch_images()

##### Extra Tool Setup

load("@com_github_ash2k_bazel_tools//multirun:deps.bzl", "multirun_dependencies")

multirun_dependencies()
