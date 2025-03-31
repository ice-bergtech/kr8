# kr8+

kr8+ is a fork of [kr8](https://github.com/apptio/kr8) with some additional features and improvements.
kr8 was used in production to great success at Apptio for managing components of multiple Kubernetes clusters.

kr8+ is a very opinionated tool used for rendering [jsonnet](http://jsonnet.org) manifests for multiple Kubernetes clusters.

It has been designed to work as a simple configuration management framework, allowing operators to specify configurations at different cluster context levels to generate component manifests across multiple clusters.

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

A cluster cluster configuration file will include:

|key|description|example|
|---|---|---|
|`_kr8_spec`| Configuration parameters for kr8 | \
```
_kr8_spec: {
  generate_dir: 'generated',
  generate_short_names: true,
}
```|
| `_cluster` | Cluster configuration, which can be used as part of the Jsonnet configuration later. This consists of things like the cluster name, type, region, and other cluster specific configuration etc. |```json
_cluster: {
    cluster_name: 'ue1-prod',
    cluster_type: 'k8',
    region_name: 'us-east-1',
    tier: 'prod',
}
``` |
| `_components` | An object with a field for each component you wish into install in a cluster. |```
_components: {
  prometheus: { path: 'components/monitoring/prometheus' }
  argocd: { path: 'components/ci/argocd'}
  argo_apps: { path: 'components/ci/argo-apps'},
}
```|
| `<component_name>` | Component configuration, which is modifications to a component which are specific to a cluster. An example of this might be the filename of an SSL certificate for the nginx-ingress controller, which may be different across cloud providers |```
prometheus+: {
  chart_version: '1.4.1',
}
``` |


### Component

A component is something you install in your cluster to make it function and work as you expect. Some examples of components might be:

 - [cert-manager](https://github.com/jetstack/cert-manager)
 - [nginx-ingress](https://github.com/kubernetes/ingress-nginx)
 - [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets)

Components are _not_ the applications you want to run in your cluster.
Components are generally applications you'd run in your cluster to make those applications function and work as expected.
Individual applications are usually configured and added separately, through automation such as argo.

| Field | Description | Example |
|---|---|---|
| `kr8_spec` | | ```
kr8_spec: {
  includes: ['some-file.jsonnet', 'some-file.yml']
}
``` |
| `release_name` |
| `namespace` | 

### Jsonnet

All configuration for kr8 is written in [Jsonnet](https://jsonnet.org/). 
Jsonnet was chosen because it allows us to use code for configuration, while staying as close to JSON as possible.

## Building

See the [Building](docs/building.md) documentation.

## Contributing

Fork the repo in github and send a merge request!

## Caveats

There are currently no tests, and the code is not very [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself).

This was (one of) Apptio's first exercise in Go, and pull requests are very welcome.
