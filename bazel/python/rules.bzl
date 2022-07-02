# Collection of all Python rules, for a single place to import rules.

load(
    "@rules_python//python:defs.bzl",
    _py_binary = "py_binary",
    _py_library = "py_library",
    _py_test = "py_test",
)
load(
    "@rules_python//python/pip_install:requirements.bzl",
    _compile_pip_requirements = "compile_pip_requirements",
)
load(
    "@rules_python//python/pip_install:repositories.bzl",
    _pip_install_dependencies = "pip_install_dependencies",
)

py_binary = _py_binary
py_library = _py_library
py_test = _py_test
compile_pip_requirements = _compile_pip_requirements
pip_install_dependencies = _pip_install_dependencies
