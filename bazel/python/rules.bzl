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
    "@rules_python//gazelle/manifest:defs.bzl",
    _gazelle_python_manifest = "gazelle_python_manifest",
)
load(
    "@rules_python//gazelle/modules_mapping:def.bzl",
    _modules_mapping = "modules_mapping",
)

py_binary = _py_binary
py_library = _py_library
py_test = _py_test

compile_pip_requirements = _compile_pip_requirements
gazelle_python_manifest = _gazelle_python_manifest
modules_mapping = _modules_mapping
