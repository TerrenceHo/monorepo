load("//bazel/python:pip_install.bzl", "pip_install")

# Each python project (or repo) will have its own requirements repository. This
# results in more package duplication, but results in fewer package requirements
# incompatibilities.
PYTHON_REPOS = [
    "//example_python:requirements.txt",
]

# for each requirements file in PYTHON_REPOS, run pip_install. The name of the
# requirements WORKSPACE is generated from the path to the requirements file.
# Something like "//example-project/lib:requirements.txt" would be translated to
# "@example-project_lib_pip_deps"
def fetch_python_deps():
    for requirements in PYTHON_REPOS:
        name = "{}_pip_deps".format(
            requirements[2:].split(":")[0].replace("/", "_"),
        )
        pip_install(
            name = name,
            requirements = requirements,
        )
