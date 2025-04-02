# kr8+

[![CI status](https://github.com/ice-bergtech/kr8/workflows/CI/badge.svg)](https://github.com/ice-bergtech/kr8/actions?query=workflow%3ACI)

kr8+ is a fork of [kr8](https://github.com/apptio/kr8) with some additional features and improvements.
kr8 was used in production to great success at Apptio for managing components of multiple Kubernetes clusters.

kr8+ is a very opinionated tool used for rendering [jsonnet](http://jsonnet.org) manifests for multiple Kubernetes clusters.

It has been designed to work as a simple configuration management framework, allowing operators to specify configurations at different cluster context levels to generate component manifests across multiple clusters.

Kr8+ is `pre-1.0`.
This means that breaking changes will still happen from time to time, but it's stable enough for both scripting and interactive use.

## Features

- Generate and customize component configuration for Kubernetes clusters across environments, regions and platforms
- Opinionated config, flexible deployment. kr8+ simply generates manifests for you, you decide how to deploy them
- Render and override component config from multiple sources, such as Helm, Kustomize and static manifests
- CI/CD friendly

For more information about the inspiration and the problem kr8+ solves, check out this [blog post](https://leebriggs.co.uk/blog/2018/05/08/kubernetes-config-mgmt.html).

kr8+ consists of:

- kr8+ - a Go binary for rendering manifests
- jsonnet - [go-jsonnet](https://pkg.go.dev/github.com/google/go-jsonnet) `v0.20.0`
- template - [text/template](https://pkg.go.dev/text/template#hdr-Text_and_spaces)

kr8+ is not designed to be a tool to help you install and deploy applications.
It's specifically designed to manage and maintain configuration for the cluster level services.
For more information, see the [components](docs/components) section.

In order to use kr8, you'll need a configuration repository to go with this binary. 
See the [example](./example/) directory for more information.

## Concepts & Tools

Configuration layering and updating jsonnet with `+`.

A typical repo that uses kr8 will have the following parts:

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
- argo applications: generate argo cd applications for manageing applying cluster configuration to live nodes
- application: a single application that you want to run in your cluster. This is usually a web application, but it can also be a database, cron job, or documentation.

Components are applications you want to run in your cluster.
Components are generally applications you'd run in your cluster to make those applications function and work as expected.

See the [Components](docs/concepts/components.md) documentation.

### Jsonnet

All configuration for kr8 is written in [Jsonnet](https://jsonnet.org/). 
Jsonnet was chosen because it allows us to use code for configuration, while staying as close to JSON as possible.

## Building

```sh
go build
# or
task build-snapshot
```

See the [Building](docs/building.md) documentation.

## Testing

```sh
git submodule init
git submodule update --remote --init
```

```sh
go task test
```

## Contributing

Fork the repo in github and send a merge request!

## Caveats

There are currently no tests, and the code is not very [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself).

This was (one of) Apptio's first exercise in Go, and pull requests are very welcome.
