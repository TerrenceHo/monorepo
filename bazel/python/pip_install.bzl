load("@rules_python//python:pip.bzl", _pip_install = "pip_install")

def pip_install(**kwargs):
    extra_pip_args = kwargs.pop("extra_pip_args", []) + [
        "--disable-pip-version-check",
    ]

    # python_interpreter_target = kwargs.pop(
    #     "python_interpreter_target",
    #     "@python3_interpreter//:bazel_install/bin/python3",
    # )
    quiet = kwargs.pop("quiet", True)

    _pip_install(
        extra_pip_args = extra_pip_args,
        # python_interpreter_target = python_interpreter_target,
        quiet = quiet,
        **kwargs
    )
