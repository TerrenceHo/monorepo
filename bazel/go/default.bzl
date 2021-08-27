load("@io_bazel_rules_docker//go:image.bzl", _go_image = "go_image")
load(
    "@io_bazel_rules_go//go:def.bzl",
    _go_binary = "go_binary",
    _go_library = "go_library",
    _go_test = "go_test",
)

def go_image(goarch = "amd64", goos = "linux", **kwargs):
    _go_image(goarch = goarch, goos = goos, **kwargs)

def go_binary(name, **kwargs):
    """
    Redfine go_binary to produce both the normal binary and a linux binary for
    container images to use.
    """
    _go_binary(
        name = name,
        **kwargs
    )

    _go_binary(
        name = "linux/" + name,
        goarch = "amd64",
        goos = "linux",
        **kwargs
    )

def go_test(name, **kwargs):
    """
    Redfine go_test to produce test in both the normal and linux environment.
    """

    _go_test(
        name = name,
        **kwargs
    )

    # _go_test(
    #     name = "linux/" + name,
    #     goarch = "amd64",
    #     goos = "linux",
    #     **kwargs
    # )

# Reexports
go_library = _go_library
