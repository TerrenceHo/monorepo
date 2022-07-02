"""
Loads all the repositories created in `//bazel/python:deps.bzl`, and creates an
external repository for each pip dependency package. This does not fetch the
actual dependencies, which only happen when the target is built/tested.

This is located in a separate file because the repositories don't exist when the WORKSPACE file is first loaded, and load statements must go at the top level of a file.
"""

load("@example-python_pip_deps//:requirements.bzl", _example_python_deps = "install_deps")

def load_pip_repositories():
    _example_python_deps()
