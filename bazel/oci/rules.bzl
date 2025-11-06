load("@bazel_lib//lib:expand_template.bzl", "expand_template")
load(
    "@rules_oci//oci:defs.bzl",
    _oci_image = "oci_image",
    _oci_load = "oci_load",
    _oci_push = "oci_push",
    _oci_push_rule = "oci_push_rule",
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
    name = kwargs.pop("name")
    registry = kwargs.pop("registry", "ghcr.io")
    repository = kwargs.pop("repository")
    remote_tags = kwargs.pop("remote_tags", ["${STABLE_GIT_COMMIT}"])
    tags_name = "{}_tags".format(name)

    native.genrule(
        name = tags_name,
        srcs = [],
        cmd = "grep 'STABLE_GIT_COMMIT' bazel-out/stable-status.txt | awk '{print $$2}' >> $@",
        stamp = -1,
        outs = ["{}_file.txt".format(tags_name)],
    )

    # expand_template(
    #     name = tags_name,
    #     template = ["STABLE_GIT_COMMIT"],
    #     substitutions = {"STABLE_GIT_COMMIT": "$(STABLE_GIT_COMMIT)"},
    # )
    _oci_push_rule(
        name = name,
        repository = "{registry}/{repository}".format(registry = registry, repository = repository),
        remote_tags = ":{}".format(tags_name),
        **kwargs
    )
