load("@bazel_gazelle//:deps.bzl", "go_repository")

# fetch_go_deps is a macro that is modified by gazelle. DO NOT MODIFY MANUALLY
def fetch_go_deps():
    go_repository(
        name = "com_github_google_uuid",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/uuid",
        sum = "h1:qJYtXnJRWmpe7m/3XlyhrsLrEURqHRM2kxzoxXqyUDs=",
        version = "v1.2.0",
    )
