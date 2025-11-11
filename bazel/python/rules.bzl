# Collection of all Python rules, for a single place to import rules.

load("@rules_python//python:py_binary.bzl", _py_binary = "py_binary")
load("@rules_python//python:py_library.bzl", _py_library = "py_library")
load("@rules_python//python:py_test.bzl", _py_test = "py_test")

py_binary = _py_binary
py_library = _py_library
py_test = _py_test

# load(
#     "@rules_python//python/pip_install:requirements.bzl",
#     _compile_pip_requirements = "compile_pip_requirements",
# )
# load(
#     "@rules_python//gazelle/manifest:defs.bzl",
#     _gazelle_python_manifest = "gazelle_python_manifest",
# )
# load(
#     "@rules_python//gazelle/modules_mapping:def.bzl",
#     _modules_mapping = "modules_mapping",
# )
# compile_pip_requirements = _compile_pip_requirements
# gazelle_python_manifest = _gazelle_python_manifest
# modules_mapping = _modules_mapping
