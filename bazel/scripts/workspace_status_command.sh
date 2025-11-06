#!/bin/bash

# WORKSPACE variables that can be passed into builds
# Docs: https://docs.bazel.build/versions/main/user-manual.html#workspace_status

echo "CURRENT_TIME $(date +%s)"
echo "STABLE_GIT_COMMIT $(git rev-parse HEAD)"
echo "writing" >>/home/tho/bazel_pres
