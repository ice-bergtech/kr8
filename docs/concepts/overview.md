# Concepts

kr8+ has two main concepts you should be aware of before you get started:

  - [components](components.md) - one or more application packaged together
  - [clusters](clusters.md) - a deployment environment, organized as a tree of configuration

The relationship between components and clusters are simple: components are installed on clusters.
You will have:

* components that are installed on all clusters (auth, cert management, secrets, monitoring)
* components that are only installed on _some_ clusters (services, hardware dependent workloads)
* components that have multiple versions/deployments installed on a single cluster (upgrades, namespacing)


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




Components can be declared multiple times within a cluster, as long as they are named distinctly.
The are defined in `<baseDir>/components` by default, or the directory specified by the `--componentdir`, `-X` flags.

Clusters are unique and singular.
They have a name which is specified via the directory structure under  `<baseDir>/clusters` by default, or the directory specified by the `--clusterdir`, `-D` flags.

`baseDir` is `.` by default.

Read more about [components](components.md) and [clusters](clusters.md)
