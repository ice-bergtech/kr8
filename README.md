# kr8+

[![CI status](https://github.com/ice-bergtech/kr8/workflows/CI/badge.svg)](https://github.com/ice-bergtech/kr8/actions?query=workflow%3ACI)
![Go Version](https://img.shields.io/badge/go%20version-%3E=1.23-61CFDD.svg?style=flat-square)

kr8+ is a fork of [kr8](https://github.com/apptio/kr8) with some additional features and improvements.
kr8 was used in production to great success at Apptio for managing components of multiple Kubernetes clusters.

kr8+ is a very opinionated tool used for rendering [jsonnet](http://jsonnet.org) manifests for multiple Kubernetes clusters.

It has been designed to work as a simple configuration management framework, allowing operators to specify configurations at different cluster context levels to generate component manifests across multiple clusters.

Kr8+ is `pre-1.0`.
This means that breaking changes will still happen from time to time, but it's stable enough for both scripting and interactive use.

# Installation

The latest version is available for download from the [Github releases page](https://github.com/ice-bergtech/kr8/releases).

## Features

- Generate and customize component configuration for Kubernetes clusters across environments, regions and platforms
- Opinionated config, flexible deployment. kr8+ simply generates manifests for you, you decide how to deploy them
- Render and override component config from multiple sources, such as Helm, Kustomize and static manifests
- CI/CD friendly

For more information about the inspiration and the problem kr8+ solves, check out this [blog post](https://leebriggs.co.uk/blog/2018/05/08/kubernetes-config-mgmt.html).

kr8+ consists of:

- kr8+ - a Go binary for rendering manifests
- [go-jsonnet](https://pkg.go.dev/github.com/google/go-jsonnet) `v0.20.0`
- [ghodss/yaml](https://github.com/ghodss/yaml) `v1.0.0`
- [Grafana/tanka helm](https://github.com/grafana/tanka/pkg/helm) `v0.27.1`
- [kompose](https://github.com/kubernetes/kompose) `v1.35.0`
- [Masterminds/sprig v3 Template Library](https://pkg.go.dev/github.com/Masterminds/sprig#section-readme) - [Template Documentation](https://masterminds.github.io/sprig/) `v3.2.3`

kr8+ is not designed to be a tool to help you install and deploy applications.
It's specifically designed to manage and maintain configuration for the cluster level services.
For more information, see the [components](docs/components) section.

In order to use kr8+, you'll need a configuration repository to go with this binary. 
See the [example](./example/) directory for more information.

## Concepts & Tools

Configuration layering and updating jsonnet with `+`.

A typical repo that uses kr8+ will have the following parts:

* Cluster Configurations
* Component Configurations
* Jsonnet Libraries

### Cluster Config

A cluster is a Kubernetes cluster running in a cloud provider, datacenter or elsewhere.
You will more than likely have multiple clusters across multiple environments and regions.

See the [Clusters](docs/concepts/clusters.md) documentation.

### Component

A component is something you install in your cluster to make it function and work as you expect.
Some examples of components might be:

- cluster core resources: [cert-manager](https://github.com/jetstack/cert-manager) or [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets)
- argo applications: generate argo cd applications for managing applying cluster configuration to live nodes
- application: a single application that you want to run in your cluster. This is usually a web application, but it can also be a database, cron job, or documentation.

Components are applications you want to run in your cluster.
Components are generally applications you'd run in your cluster to make those applications function and work as expected.

See the [Components](docs/concepts/components.md) documentation.

### Jsonnet

All configuration for kr8+ is written in [Jsonnet](https://jsonnet.org/). 
Jsonnet was chosen because it allows us to use code for configuration, while staying as close to JSON as possible.

## Development

kr8+ is coded in [Golang](https://golang.org/).
Currently, version `1.23.0` is used.

Common tasks can be executed with `go-task`.

Tasks are described in [Taskfile.yml](Taskfile.yml).

### Prerequisites

- Go `1.23` or later
- `git` for cloning submodules
- [go-task](https://taskfile.dev/installation/) for task automation
- [golangci-lint](https://golangci-lint.run/welcome/install/) for linting
- [Bats](https://bats-core.readthedocs.io/en/stable/installation.html) for binary testing
- [Goreleaser]() for release packaging

Once `go-task` is installed, you can setup your environment by running:

```sh
# Install dev tools
task setup

# View available tasks
task -l

# Build kr8+ for your local system
task build

# Run tests
task test

# Build snapshot
task build-snapshot
```

### Build Troubleshooting

1. Dependencies download fail: There is a big number of reasons this could fail but the most important might be:
   * Networking problems: Check your connection to: github.com, golang.org and k8s.io.
   * Disk space: If no space is available on the disk, this step might fail.
2. The comand `go build` does not start the build:
   * Confirm you are in the correct project directory
   * Make sure your go installation works: `go --version`

## Contributing

Fork the repo in github and send a merge request!

## Caveats

There are currently no tests.

Pull requests are very welcome.

## License

The project is licensed under the [MIT license](LICENSE).

Parts of the code are derived from:

* [kr8](https://github.com/apptio/kr8) - [MIT License](LICENSE-apptio)
* [Lee Briggs](https://leebriggs.co.uk/) - [MIT License](LICENSE-briggs)
* [kubecfg](https://github.com/kubecfg/kubecfg) - [Apache 2.0](LICENSE-kubecfg)
