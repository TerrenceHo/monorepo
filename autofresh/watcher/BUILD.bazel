load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "watcher",
    srcs = ["recursive_watcher.go"],
    importpath = "github.com/TerrenceHo/monorepo/autofresh/watcher",
    visibility = ["//visibility:public"],
    deps = [
        "//utils-go/stackerrors",
        "@com_github_fsnotify_fsnotify//:fsnotify",
    ],
)

go_test(
    name = "watcher_test",
    size = "small",
    srcs = ["recursive_watcher_test.go"],
    embed = [":watcher"],
    deps = [
        "@com_github_fsnotify_fsnotify//:fsnotify",
        "@com_github_stretchr_testify//assert",
    ],
)
