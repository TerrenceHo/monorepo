load(
    "@io_bazel_rules_docker//container:container.bzl",
    _container_bundle = "container_bundle",
    _container_image = "container_image",
    _container_layer = "container_layer",
    _container_pull = "container_pull",
    _container_push = "container_push",
)

container_image = _container_image
container_layer = _container_layer
container_pull = _container_pull
container_bundle = _container_bundle

# Redefine container_push with some sane defaults
def container_push(**kwargs):
    format = kwargs.pop("format", "Docker")
    registry = kwargs.pop("registry", "ghcr.io")
    skip_unchanged_digest = kwargs.pop("skip_unchanged_digest", True)
    stamp = kwargs.pop("stamp", "@io_bazel_rules_docker//stamp:always")
    tag = kwargs.pop("tag", "{STABLE_GIT_COMMIT}")
    _container_push(
        format = format,
        registry = registry,
        skip_unchanged_digest = skip_unchanged_digest,
        stamp = stamp,
        tag = tag,
        **kwargs
    )
