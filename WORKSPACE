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

##### Python Dependencies

# load the python interpreter first
load("//bazel/python/interpreter:setup_python_interpreter.bzl", "setup_python_interpreter")

setup_python_interpreter()

load("//bazel/python:deps.bzl", "fetch_python_deps")

fetch_python_deps()

register_toolchains(
    "//bazel/python/interpreter:monorepo_py_toolchain",
)

##### Protobuf Dependencies

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

##### gRPC Dependencies

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@com_github_grpc_grpc//bazel:grpc_extra_deps.bzl", "grpc_extra_deps")

grpc_extra_deps()
