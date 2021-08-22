load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "autofresh_lib",
    srcs = ["main.go"],
    importpath = "github.com/TerrenceHo/monorepo/autofresh/cmd/autofresh",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "autofresh",
    embed = [":autofresh_lib"],
    visibility = ["//visibility:public"],
)
