# kr8+

[![CI status](https://github.com/ice-bergtech/kr8/workflows/CI/badge.svg)](https://github.com/ice-bergtech/kr8/actions?query=workflow%3ACI)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.24-61CFDD.svg?style=flat-square)

**Kr8+** is an opinionated Kubernetes cluster configuration management tool designed to simplify and standardize the process of managing Kubernetes clusters. By leveraging best practices and providing a structured approach, **Kr8+** helps DevOps teams streamline their workflows and maintain consistency across multiple environments.

**Kr8+** is `pre-1.0`.
This means that breaking changes will still happen from time to time, but it's stable enough for both scripting and interactive use.

## Key Features

* **Cluster Configuration Management**: Manage Kubernetes clusters across environments, regions and platforms with a declarative and centralized approach.
* **Opinionated Structure**: Enforces best practices for consistent and reliable cluster configurations.
* **Jsonnet Native Funcitons**: Use jsonnet to render and override component config from multiple sources, such as templates, docker-compose files, Kustomize, and Helm.
* **Extensibility**: Easily extensible to meet the needs of diverse Kubernetes environments.
* **CI/CD Friendly**: Statically define all your configuration in a single source of truth, making it easy to integrate with CI/CD pipelines and deployment automation like ArgoCD.

## Technical Overview

**kr8+** consists of:

- **kr8+** - a Go binary for rendering manifests
- [go-jsonnet](https://pkg.go.dev/github.com/google/go-jsonnet) `v0.20.0`
- [ghodss/yaml](https://github.com/ghodss/yaml) `v1.0.0`
- [Grafana/tanka helm](https://github.com/grafana/tanka/pkg/helm) `v0.27.1`
- [kompose](https://github.com/kubernetes/kompose) `v1.35.0`
- [Masterminds/sprig v3 Template Library](https://pkg.go.dev/github.com/Masterminds/sprig#section-readme) - [Template Documentation](https://masterminds.github.io/sprig/) `v3.2.3`

## Why use Kr8+?

* **Standardization**: Ensures consistency across Kubernetes clusters, reducing errors and improving maintainability.
* **Simplicity**: Provides a straightforward approach to complex Kubernetes configurations, making it easier for teams to adopt.
* **Scalability**: Designed to support clusters of all sizes, from simple single-node setups to large-scale production environments.


## Getting Started

An example of a repo can be found in the [example](./example) folder.

### Installation

The latest version is available for download from the [Github releases page](https://github.com/ice-bergtech/kr8/releases).

Once installed, you can use `kr8 init` commands to setup the initial structure and configurations.

### Configuration

All configuration for **kr8+** is written in [Jsonnet](https://jsonnet.org/).
Jsonnet was chosen because it allows us to use code for configuration, while staying as close to JSON as possible.

A typical repo that uses **kr8+** will have the following parts:

* Cluster Configurations
* Component Configurations
* Jsonnet Libraries

#### Clusters Configurations

A cluster is a deployment environment, organized in folders as a tree of configuration.
Configuration the folders is layered on the parent folder's configuration, allowing you to override or extend configurations.

Cluster Spec: [docs/godoc/kr8-types.md#Kr8ClusterJsonnet]

More information: [Managing Clusters](docs/concepts/clusters.md)

#### Conponents Configurations

A component is a deployable unit that you wish to install in one or more clusters.
Components can be declared multiple times within a cluster, as long as they are named distinctly when loaded.

Component Spec: [docs/godoc/kr8-types.md#Kr8ComponentJsonnet]

More information: [Managing Components](docs/concepts/components.md)

#### Jsonnet Libraries

Jsonnet libraries are reusable code that can be imported into your Jsonnet files.
They allow you to write modular and maintainable configuration.

Common libraries include:

* [kr8-libesonnet](https://github.com/ice-bergtech/kr8-libsonnet)
* [kube-libsonnet](https://github.com/kube-libsonnet/kube-libsonnet)

More information: [Jsonnet Libraries](docs/concepts/jsonnetlibs.md)

### Deployment

To generate the final configured manifests, just run `kr8 generate`.
This will read the configuration files and generate the final manifests based on the parameters provided.

These manifests can then be checking into source control, and ingested by tools like ArgoCD, Portainer, Rancher etc.

### Further Information

* [Command Documentation](docs/cmd/kr8.md)
* **kr8+**
  * [Concepts](docs/concepts/overview.md)
  * [Managing Clusters](docs/concepts/clusters.md)
  * [Creating Components](docs/concepts/components.md)
  * [Native Functions](docs/components/nativefuncs.md)
* [Code Documentation](docs/godoc)

## Development

**kr8+** is coded in [Golang](https://golang.org/).
Currently, version `1.24.2` is used.

Common tasks are described in [Taskfile.yml](Taskfile.yml), and can be executed with `go-task`.

### Dependencies

- Golang version  `1.24` or later ([installation](https://go.dev/doc/install))
- `git` for cloning submodules
- [go-task](https://taskfile.dev/installation/) for task automation
- [golangci-lint](https://golangci-lint.run/welcome/install/) for linting
- [Bats](https://bats-core.readthedocs.io/en/stable/installation.html) for binary testing
- [Goreleaser](https://goreleaser.com/intro/) for release packaging

### Setup

Once `go-task` is [installed](https://taskfile.dev/installation/), you can easily your environment by running:

```sh
# Install dev tools
task setup
```

### Running Tasks

```sh
# View available tasks
task -l
task: Available tasks for this project:
* 01_setup:                     Instal dev tools                                         (aliases: setup, s)
* 01_setup-bats:                Install bats testing tools                               (aliases: setup-bats)
* 02_build:                     Build kr8+ for your local system                         (aliases: build, b)
* 03_build-snapshot:            Build a snapshot for all platforms using goreleaser      (aliases: build-snapshot, bs)
* 03_generate-bats-tests:       Generate resources to test against                       (aliases: gt)
* 03_test-go:                   Tesk kr8+ for your local system                          (aliases: test, t)
* 03_test-package:              Test compiled kr8+ binary against test inputs            (aliases: test-package, tp)
* 04_generate-examples:         Generate example clusters and components with kr8+       (aliases: ge, gen)
```

#### Examples

```sh

# Build kr8+ for your local system
task build

# Run tests
task test

# Build snapshot
task build-snapshot
```

### Tests

There are few sets of tests:

- Unit Tests: `go test ./...` or `task test`
- Integration Tests using `bats`: `task test-package`
- Generate examples: `./kr8 generate -B example` or `task gen`

### Build Troubleshooting

* Dependencies download fail: There is a big number of reasons this could fail but the most important might be:
   * Networking problems: Check your connection to: github.com, golang.org and k8s.io.
   * Disk space: If no space is available on the disk, this step might fail.
* The comand `go build` does not start the build:
   * Confirm you are in the correct project directory
   * Make sure your go installation works: `go --version`

## Contributing

We welcome contributions from the community to enhance **kr8+**
Fork the repo in github and send a merge request!

Pull requests are very welcome.

## License

The project is licensed under the [MIT license](LICENSE).

Parts of the code are derived from:

* [kr8](https://github.com/apptio/kr8) - [MIT License](LICENSE-apptio)
* [Lee Briggs](https://leebriggs.co.uk/) - [MIT License](LICENSE-briggs)
* [kubecfg](https://github.com/kubecfg/kubecfg) - [Apache 2.0](LICENSE-kubecfg)

## References and Additional Resources

**kr8+** is a fork of [kr8](https://github.com/apptio/kr8) with some additional features and improvements.
**kr8** was used in production to great success at Apptio for managing components of multiple Kubernetes clusters.

* [Jsonnet Standard Library](https://jsonnet.org/ref/stdlib.html)
* [Jsonnet Language Reference](https://jsonnet.org/ref/language.html)
* [Sprig Template Documentation](https://masterminds.github.io/sprig/)
* [The growing need for Kubernetes Configuration Management](https://leebriggs.co.uk/blog/2018/05/08/kubernetes-config-mgmt.html)
