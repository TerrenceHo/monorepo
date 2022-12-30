def _local_install_impl(ctx):
    target = ctx.attr.target
    shell_commands = ""

    for s in ctx.files.srcs:
        shell_commands += "echo Copying %s to %s\n" % (s.short_path, target)
        shell_commands += "sudo cp %s %s\n" % (s.short_path, target)

    ctx.actions.write(
        output = ctx.outputs.executable,
        is_executable = True,
        content = shell_commands,
    )
    runfiles = ctx.runfiles(files = ctx.files.srcs)
    return DefaultInfo(
        executable = ctx.outputs.executable,
        runfiles = runfiles,
    )

local_install = rule(
    executable = True,
    implementation = _local_install_impl,
    attrs = {
        "srcs": attr.label_list(allow_files = True),
        "target": attr.string(default = "/usr/local/bin", doc = "Local install target directory"),
    },
)
