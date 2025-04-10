# Concepts

kr8+ is used to define and generate cluster config.
It is designed to easily layer config from multiple sources.

A deployment consists of 2 parts:

  - [components](./components.md) - one or more applications packaged together
  - [clusters](./clusters.md) - a deployment environment, organized as a tree of configuration

The relationship between components and clusters are simple: components are installed on clusters.
You will have:

* components that are installed on all clusters (auth, cert management, secrets, monitoring)
* components that are only installed on _some_ clusters (services, hardware dependent workloads)
* components that have multiple versions/deployments installed on a single cluster (upgrades, namespacing)

### Components

A component is something you install in your cluster to make it function and work as you expect.
Some examples of components might be:

- cluster core resources: [cert-manager](https://github.com/jetstack/cert-manager) or [sealed-secrets](https://github.com/bitnami-labs/sealed-secrets)
- argo applications: generate argo cd applications for managing applying cluster configuration to live nodes
- application: a single application that you want to run in your cluster. This is usually a web application, but it can also be a database, cron job, or documentation.

for more information on components see the [Components](./components.md) documentation.

### Cluster Config

A cluster is a Kubernetes cluster running in a cloud provider, datacenter or elsewhere.
You will more than likely have multiple clusters across multiple environments and regions.

By design, configuration is able to be layered and overridden at different levels.
This is the strength of jsonnet and allows for a lot of flexibility in managing your cluster configurations.

Cluster names are based on the directory structure under `./clusters` by default, or the directory specified by the `--clusterdir`, `-D` flags.

A typical cluster config folder layout may look like:

```sh
clusters
├── params.jsonnet  # top level config; used as defaults for all clusters
├── development
│   ├── dev-test
│   │   └── cluster.jsonnet  # dev configs
│   └── dev-staging
│       └── cluster.jsonnet  # staging configs
└── production
    ├── params.jsonnet # standard prod configs for cluster and components
    ├── region-1
    │   ├── params.jsonnet # region specific params
    │   ├── pre-prod
    │   │   └── cluster.jsonnet
    │   └── prod
    │       ├── cluster.jsonnet # prod-level configs
    │       ├── workloads-1
    │       │   └── cluster.jsonnet
    │       └── workloads-2
    │           └── cluster.jsonnet
    └── region-2
        ├── params.jsonnet # region specific params
        ├── pre-prod
        │   └── cluster.jsonnet
        ├── prod-1
        │   └── cluster.jsonnet
        └── workloads-1
            └── cluster.jsonnet
```

See the [Clusters](./clusters.md) documentation.
