load("@bazel_lib//lib:expand_template.bzl", "expand_template")
load("@bazel_lib//lib:stamping.bzl", "STAMP_ATTRS", "maybe_stamp")
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

# # takes in a file or string, replaces the file where {}
# def _stamp_substitute_impl(ctx):
#     output = ctx.outputs.out
#     if not output:
#         if ctx.file.template and ctx.file.template.is_source:
#             output = ctx.actions.declare_file(ctx.file.template.basename, sibling = ctx.file.template)
#         else:
#             output = ctx.actions.declare_file(ctx.attr.name + ".txt")
#
#     stamp = maybe_stamp(ctx)
#     if stamp:
#         files = []
#
#         args = ctx.actions.args()
#         args.add_all([
#             ctx.file.template,
#             output,
#             stamp.stable_status_file,
#             stamp.volatile_status_file,
#         ])
#
#         # datadog shell script
#         ctx.actions.run_shell(
#             inputs = [
#                 ctx.file.template,
#                 stamp.stable_status_file,
#                 stamp.volatile_status_file,
#             ],
#             outputs = [
#                 output,
#             ],
#             arguments = [args],
#             command = """#!/usr/bin/env bash
#             scratch=$(cat $1)
#             shift
#
#             out=$1
#             shift
#
#             for file in $@
#             do
#                 while read -r key value
#                 do
#                     # Replace the keys with their corresponding values in the scratch output
#                     scratch=${scratch//\\{$key\\}/$value}
#                 done <$file
#             done
#
#             >$out echo -n "$scratch"
#             """,
#         )
#
#     else:
#         # glorified copy with no substitutions, write to output file
#         ctx.actions.expand_template(
#             template = ctx.file,
#             output = output,
#             subtitutions = {},
#             # is_executable = ctx.attr.is_executable,
#         )
#
#     all_outs = [output]
#
#     runfiles = ctx.runfiles(
#         files = all_outs,
#         transitive_files = depset(transitive = [
#             target[DefaultInfo].files
#             for target in ctx.attr.data
#         ]),
#     )
#     return [DefaultInfo(
#         files = depset(all_outs),
#         runfiles = runfiles.merge_all([
#             target[DefaultInfo].default_runfiles
#             for target in ctx.attr.data
#         ]),
#     )]
#
# stamp_substitute_rule = rule(
#     doc = """Substitutes values in file with workspace status templates
#
#     If --stamp is not enabled, then the file will not be templated
#     """,
#     implementation = _stamp_substitute_impl,
#     is_executable = true,
#     attrs = dict({
#         "template": attr.label(
#             doc = "The template file to substitute with values from workspace status",
#             mandatory = true,
#             allow_single_file = true,
#         ),
#         "out": attr.output(
#             doc = """Where to write the expanded file.
#
#             If the `template` is a source file, then `out` defaults to
#             be named the same as the template file and outputted to the same
#             workspace-relative path. In this case there will be no pre-declared
#             label for the output file. It can be referenced by the target label
#             instead. This pattern is similar to `copy_to_bin` but with substitutions on
#             the copy.
#
#             Otherwise, `out` defaults to `[name].txt`.
#             """,
#         ),
#         is_executable: attr.bool(
#             doc = "Whether to mark the output file as executable",
#         ),
#     }, **STAMP_ATTRS),
# )
#
# def stamp_substitute(name, template, **kwargs):
#     """Wrapper macro for `stamp_substitute_rule`.
#
#     Args:
#         name: name of resulting rule
#         template: label of template file, or list of strings which represent the content of the template file to be created
#         **kwargs: other named parameters for `stamp_substitute_rule`.
#     """
#     if types.is_list(template):
#         write_target = "{}_tmpl".format(name)
#         write_file(
#             name = write_target,
#             out = "{}.txt".format(write_target),
#             content = template,
#         )
#         template = write_target
#
#     stamp_substitute_rule(name = name, template = template, **kwargs)

def oci_push(**kwargs):
    name = kwargs.pop("name")
    registry = kwargs.pop("registry", "ghcr.io")
    repository = kwargs.pop("repository")
    tag = kwargs.pop("tag", "dev")

    expand_template(
        name = "{}_stamped".format(name),
        template = [tag],
        stamp_substitutions = {tag: "{{STABLE_GIT_COMMIT}}"},
        out = "{}_stamped.tags.txt".format(name),
    )

    _oci_push(
        name = name,
        repository = "{}/{}".format(registry, repository),
        remote_tags = ":{}_stamped".format(name),
        **kwargs
    )
