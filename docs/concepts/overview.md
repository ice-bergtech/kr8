# Concepts

kr8+ has two main concepts you should be aware of before you get started:

  - [components](components.md) - one or more application packaged together
  - [clusters](clusters.md) - a deployment environment, organized as a tree of configuration

The relationship between components and clusters are simple: components are installed on clusters.
You will have:

* components that are installed on all clusters (auth, cert management, secrets, monitoring)
* components that are only installed on _some_ clusters (services, hardware dependent workloads)
* components that have multiple versions/deployments installed on a single cluster (upgrades, namespacing)

Components can be declared multiple times within a cluster, as long as they are named distinctly.
The are defined in `<baseDir>/components` by default, or the directory specified by the `--componentdir`, `-X` flags.

Clusters are unique and singular.
They have a name which is specified via the directory structure under  `<baseDir>/clusters` by default, or the directory specified by the `--clusterdir`, `-D` flags.

`baseDir` is `.` by default.

Read more about [components](components.md) and [clusters](clusters.md)
