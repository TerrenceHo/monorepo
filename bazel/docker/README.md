# Bazel Docker

## Pushing Images

If you want to push an image, import the `container_push` rule that is
reexported from `//bazel/docker/rules.bzl`. The exported rule has some sane
defaults set, but these defaults can be overridden. An example of using the
`container_push` rule is below.

Example:

```
load("//bazel/docker:rules.bzl", "container_push")

container_push(
    name = "push",
    image = "<image rule>",
    repository = "<user>/<repository name of image>",
    visibility = ["//visibility:public"],
)
```

A new image will not push if the image SHA256 digest is not changed (saving push
time and image storage space).

## CI Push

If a new image is created and needs to be build/pushed with all commits, then
add the image push rule over to the `multirun` rule in the root BIULD.bazel
file. When the master CI passes in the GitHub actions after a PR lands, a job
will invoke the `//:push_all` rule and push all images to their respective
repositories.

For pushes to succeed, you must allow the image repository to be pushed from the
monorepo repository. Go to
`https://github.com/users/<username>/packages/container/<image-name>/settings`
and allow the monorepo to write to the image repository.
