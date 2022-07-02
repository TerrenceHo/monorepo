load("@rules_python//python:pip.bzl", "pip_parse")

# Each python project (or repo) will have its own requirements repository. This
# results in more package duplication, but results in fewer package requirements
# incompatibilities.
PYTHON_REPOS = [
    "//example-python:requirements.txt.lock",
]

# for each requirements file in PYTHON_REPOS, run pip_fetch. The name of the
# requirements WORKSPACE is generated from the path to the requirements file.
# Something like "//example-project/lib:requirements.txt" would be translated to
# "@example-project_lib_pip_deps"
#
# pip_fetch creates a new separate WORKSPACe for each of our resolved Python
# dependencies. We can later invoke these workspaces to install the dependencies
# lazily.
def setup_pip_repositories():
    for requirements in PYTHON_REPOS:
        name = "{}_pip_deps".format(
            requirements[2:].split(":")[0].replace("/", "_"),
        )
        pip_parse(
            name = name,
            requirements_lock = requirements,
        )
