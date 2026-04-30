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
        digest = "sha256:668c3ca24a6f511f20cb458619f83ef8a14b27de0425340ca76a3dff0cd1e6e8",
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
