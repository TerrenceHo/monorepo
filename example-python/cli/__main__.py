import platform
import sys

# use flask as a sample dependency
from flask import Flask

# importing local dependencies with py_library
from example.lib import cmd

app = Flask(__name__)


@app.route("/")
def hello():
    return "Hello World!"


if __name__ == "__main__":
    print("Bazel python executable is", sys.executable)
    print("Bazel python version is", platform.python_version())

    print("Host python executable is", cmd(["which", "python3"]))
    print(
        "Host python version is",
        cmd(["python3", "-c", "import platform; print(platform.python_version())"]),
    )

    app.run(port=12345)
