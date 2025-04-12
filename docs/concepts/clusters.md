# Clusters

A cluster is a deployment environment, organized as a tree of configuration.
Clusters are unique and singular.

They have a name which is specified via the directory structure under  `./clusters` by default, or the directory specified by the `--clusterdir`, `-D` flags.

The cluster spec is defined as a golang struct: [docs/godoc/kr8-types.md#Kr8ClusterJsonnet]

Cluster configurations are defined in a hierarchical system which is loosely inspired by [Hiera](https://puppet.com/docs/hiera/3.3/index.html) from the Puppet ecosystem.

An example definition for clusters might look like this:

```bash
clusters
├── development
│   └── dev1
│       └── cluster.jsonnet
├── params.jsonnet
└── production
    ├── params.jsonnet
    ├── prod1
    │   └── cluster.jsonnet
    └── prod2
        └── cluster.jsonnet
```

The clusters are given names and then grouped together inside a directory tree.

There are two jsonnet files you'll notice here:

 - `cluster.jsonnet` - this defines a named cluster, kr8+ stops going down the directory tree when it finds this file
 - `params.jsonnet` - this file contains components and parameters defined for clusters. We'll go into more detail about this shortly.

## Components in clusters

Adding a component to a cluster involves add a component key to the `_components` key inside _either_ the `cluster.jsonnet` (for cluster specific components) or the `params.jsonnet` (for deploy components).

Here's an example:

```yaml
_components+: {
  sealed_secrets: { path: 'components/sealed_secrets' },
},
```

Notice we're using the jsonnet `+` operator to append the `sealed_secrets` field to the `_components` object.

## Cluster parameters

Once you've initialized a component for a cluster, you can then start to override parameters for that component.
You do this by simply defining a jsonnet key with the named parameters in it.
Here's an example:

```yaml
external_dns+: {
    provider: 'cloudflare',
    txtPrefix: 'dev1',
    txtOwnerId: 'dev1-',
    domainFilters: [
      'example.com',
    ],
    tolerateMasters: false,
  },
```

## Hierarchy System

The hierarchy system is a very powerful part or kr8+.
It allows you to remove duplication of parameter and component definitions.
Take the previous cluster layout as an example:

```sh
clusters
├── development
│   └── dev1
│       └── cluster.jsonnet
├── params.jsonnet
└── production
    ├── params.jsonnet
    ├── prod1
    │   └── cluster.jsonnet
    └── prod2
        └── cluster.jsonnet
```

You can use the hierarchy system to ensure you have components installed in all clusters.
Let's assume we want to make sure that we want to install the `sealed_secrets` component in _all_ our clusters.
We'd put it in this file:

```sh
clusters
├── development
│   └── dev1
│       └── cluster.jsonnet
├── params.jsonnet <---- place component here
└── production
    ├── params.jsonnet
    ├── prod1
    │   └── cluster.jsonnet
    └── prod2
        └── cluster.jsonnet
```

Alongside that, let's assume all our production clusters are using the same external_dns domain name.
We can define that like so:

```sh
clusters
├── development
│   └── dev1
│       └── cluster.jsonnet
├── params.jsonnet
└── production
    ├── params.jsonnet <--- place external_dns configuration here
    ├── prod1
    │   └── cluster.jsonnet
    └── prod2
        └── cluster.jsonnet
```

kr8+ looks for the smallest unit of configuration, so if you want one cluster to be slightly different inside a hierarchy unit, you can continue to override components and parameters inside a clusters' `cluster.jsonnet` file.

---

A cluster cluster configuration file includes:

`_kr8_spec` Configuration parameters for kr8+
| key                    | description | example     |
| ---------------------- | ----------- | ----------- |
| `generate_dir`         |             | 'generated' |
| `generate_short_names` |             | true        |

| `_cluster` | Cluster configuration, which can be used as part of the Jsonnet configuration later. This consists of things like the cluster name, type, region, and other cluster specific configuration etc. |```json

| key            | description | example     |
| -------------- | ----------- | ----------- |
| `cluster_name` |             | 'ue1-prod'  |
| `cluster_type` |             | 'k8'        |
| `region_name`  |             | 'us-east-1' |
| `tier`         |             | 'prod'      |

