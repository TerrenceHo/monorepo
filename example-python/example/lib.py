import subprocess


def cmd(args):
    process = subprocess.Popen(args, stdout=subprocess.PIPE)
    out, _ = process.communicate()
    return out.decode("ascii").strip()


def func(a: int) -> int:
    return a + a
