load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

def fetch_images():
    # Docs for distroless static and base (and debug versions):
    # https://github.com/GoogleContainerTools/distroless/tree/main/base
    container_pull(
        name = "distroless-static",
        registry = "gcr.io",
        repository = "distroless/static-debian10",
    )

    container_pull(
        name = "distroless-static-debug",
        registry = "gcr.io",
        repository = "distroless/static-debian10",
        tag = "debug",
        digest = "sha256:29a7c9bd164728e83937518f7e339a51ae096292ba6410cc2bcc07c4e99533f8",
    )

    container_pull(
        name = "distroless-base",
        registry = "gcr.io",
        repository = "distroless/base",
    )

    container_pull(
        name = "distroless-base-debug",
        registry = "gcr.io",
        repository = "distroless/base",
        tag = "debug",
    )
