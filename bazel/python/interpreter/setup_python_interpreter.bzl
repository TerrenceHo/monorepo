# Sourced from https://testdriven.io/blog/bazel-builds/ and
# https://thethoughtfulkoala.com/posts/2020/05/16/bazel-hermetic-python.html
# Special logic for building python interpreter with OpenSSL from homebrew.
# See https://devguide.python.org/setup/#macos-and-os-x

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

_py_configure = """
if [[ "$OSTYPE" == "darwin"* ]]; then
    ./configure --prefix=$(pwd)/bazel_install --with-openssl=$(brew --prefix openssl)
else
    ./configure --prefix=$(pwd)/bazel_install
fi
"""

_py_make = """
function num_cores() {
    platform=$(uname -s)
    if [[ "$platform" == "Linux" ]]; then
        nproc --all
    elif [[ "$platform" == "Darwin" ]]; then
        sysctl -n hw.physicalcpu
    else
        echo "1"
    fi
}
make --jobs=$(num_cores)
make install
"""

_py_build = """
exports_files(["bazel_install/bin/python3"])
filegroup(
    name = "files",
    srcs = glob(["bazel_install/**"], exclude = ["**/* *", "**/__pycache__/*"]),
    visibility = ["//visibility:public"],
)
filegroup(
    name = "includes",
    srcs = glob(["bazel_install/include/python3.8/**/*.h"]),
)
cc_library(
    name = "python_headers",
    hdrs = [":includes"],
    includes = ["bazel_install/include/python3.8"],
    visibility = ["//visibility:public"],
)
"""

def setup_python_interpreter():
    http_archive(
        name = "python3_interpreter",
        urls = ["https://www.python.org/ftp/python/3.8.3/Python-3.8.3.tar.xz"],
        sha256 = "dfab5ec723c218082fe3d5d7ae17ecbdebffa9a1aea4d64aa3a2ecdd2e795864",
        strip_prefix = "Python-3.8.3",
        patch_cmds = [
            "mkdir $(pwd)/bazel_install",
            _py_configure,
            _py_make,
        ],
        build_file_content = _py_build,
    )
