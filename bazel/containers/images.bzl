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
