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

go_register_toolchains(version = "1.16.3")

load("@io_bazel_rules_go//extras:embed_data_deps.bzl", "go_embed_data_dependencies")

go_embed_data_dependencies()

load("//bazel/go:deps.bzl", "fetch_go_deps")

# gazelle:repository_macro bazel/go/deps.bzl%fetch_go_deps
fetch_go_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
