# Building kr8+

kr8+ is coded in [Golang](https://golang.org/).
Currently, version `1.23.0` is used.

## Prerequisites

1. Install and configure Go: https://golang.org/doc/install
2. Get familiar with Golang: https://golang.org/doc/
3. If you are fully testing the build you need to install: https://github.com/go-task/task

----

## Building the executable

On the project root:

### Quick build

```sh
task build
```

or

```sh
go build
```

### Snapshot build

```sh
task build-snapshot
```

or

```sh
goreleaser release --snapshot --clean --skip=homebrew,docker
```

## Troubleshooting the process

1. Dependencies download fail: There is a big number of reasons this could fail but the most important might be:
   * Networking problems: Check your connection to: github.com, golang.org and k8s.io.
   * Disk space: If no space is available on the disk, this step might fail.
2. The comand `go build` does not start the build:
   * Confirm you are in the correct project directory
   * Make sure your go installation works: `go --version`

----
