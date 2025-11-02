load(
    "@rules_oci//oci:defs.bzl",
    _oci_image = "oci_image",
    _oci_load = "oci_load",
    _oci_push = "oci_push",
)

BASE_GIT_REPO = "https://github.com/TerrenceHo/monorepo"
IMAGE_SOURCE_KEY = "org.opencontainers.image.source"

oci_load = _oci_load

# TODO: prefer annotations over labels
# https://github.com/opencontainers/image-spec/blob/main/annotations.md#back-compatibility-with-label-schema
def oci_image(**kwargs):
    labels = kwargs.pop("labels", {})
    if IMAGE_SOURCE_KEY not in labels:
        labels[IMAGE_SOURCE_KEY] = BASE_GIT_REPO
    _oci_image(
        labels = labels,
        **kwargs
    )

def oci_push(**kwargs):
    registry = kwargs.pop("registry", "ghcr.io")
    repository = kwargs.pop("repository")
    remote_tags = kwargs.pop("remote_tags", ["{STABLE_GIT_COMMIT}"])
    _oci_push(
        repository = "{registry}/{repository}".format(registry = registry, repository = repository),
        remote_tags = remote_tags,
        **kwargs
    )
