workspace(
    name = "com_github_terrenceho_monorepo",
)

load("//bazel:workspace_deps.bzl", "fetch_deps")

fetch_deps()

# Needed for skylark unit testing
load("@bazel_skylib//:workspace.bzl", "bazel_skylib_workspace")

bazel_skylib_workspace()

##### Nix Dependencies

load("@io_tweag_rules_nixpkgs//nixpkgs:repositories.bzl", "rules_nixpkgs_dependencies")

rules_nixpkgs_dependencies()

load("@io_tweag_rules_nixpkgs//nixpkgs:nixpkgs.bzl", "nixpkgs_git_repository", "nixpkgs_python_configure")

nixpkgs_git_repository(
    name = "nixpkgs",
    revision = "21.05",  # Any tag or commit hash
    sha256 = "",  # optional sha to verify package integrity!
)

nixpkgs_python_configure(
    repository = "@nixpkgs//:default.nix",
)

##### Go Dependencies

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.18.3")

load("@io_bazel_rules_go//extras:embed_data_deps.bzl", "go_embed_data_dependencies")

go_embed_data_dependencies()

load("//bazel/go:deps.bzl", "fetch_go_deps")

# gazelle:repository_macro bazel/go/deps.bzl%fetch_go_deps
fetch_go_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

##### Python Dependencies

load("//bazel/python:deps.bzl", "fetch_python_deps")

fetch_python_deps()

##### Docker Dependencies
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
