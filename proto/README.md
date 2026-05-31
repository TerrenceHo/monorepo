# Proto Setup with gRPC

This commits generated source files into the repository, to satisfy tools like `go mod tidy` or LSP servers like `gopls` that having source files on disk. To generate the source files, you must run 
```
bazel run //proto/api/v1/<service>:write_generated
```

There will soon be an updated multi-target to automatically generate all sources into the tree
