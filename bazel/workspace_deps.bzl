load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

def fetch_deps():
    http_archive(
        name = "bazel_skylib",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.3.0/bazel-skylib-1.3.0.tar.gz",
            "https://github.com/bazelbuild/bazel-skylib/releases/download/1.3.0/bazel-skylib-1.3.0.tar.gz",
        ],
        sha256 = "74d544d96f4a5bb630d465ca8bbcfe231e3594e5aae57e1edbf17a6eb3ca2506",
    )

    http_archive(
        name = "io_bazel_rules_go",
        sha256 = "16e9fca53ed6bd4ff4ad76facc9b7b651a89db1689a2877d6fd7b82aa824e366",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.34.0/rules_go-v0.34.0.zip",
            "https://github.com/bazelbuild/rules_go/releases/download/v0.34.0/rules_go-v0.34.0.zip",
        ],
    )

    http_archive(
        name = "bazel_gazelle",
        sha256 = "5982e5463f171da99e3bdaeff8c0f48283a7a5f396ec5282910b9e8a49c0dd7e",
        urls = [
            "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
            "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.25.0/bazel-gazelle-v0.25.0.tar.gz",
        ],
        # temporary patch until this patch makes it into an official release
        patches = ["//bazel/patches:gazelle.patch"],
        patch_args = ["-p1"],
    )

    http_archive(
        name = "rules_python",
        sha256 = "2f5c284fbb4e86045c2632d3573fc006facbca5d1fa02976e89dc0cd5488b590",
        strip_prefix = "rules_python-1.6.3",
        url = "https://github.com/bazel-contrib/rules_python/releases/download/1.6.3/rules_python-1.6.3.tar.gz",
    )

    http_archive(
        name = "io_bazel_rules_docker",
        sha256 = "b1e80761a8a8243d03ebca8845e9cc1ba6c82ce7c5179ce2b295cd36f7e394bf",
        urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.25.0/rules_docker-v0.25.0.tar.gz"],
    )

    git_repository(
        name = "com_github_ash2k_bazel_tools",
        commit = "6e6c3717bb879f27d473a80fe8a43f2cc0e25298",
        remote = "https://github.com/ash2k/bazel-tools.git",
        shallow_since = "1655981442 +1000",
    )
