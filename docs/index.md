# kr8+

kr8+ is a configuration management tool for Kubernetes clusters, designed to generate deployable manifests for the components required to make your clusters usable.

Its main function is to manipulate JSON and YAML without using a templating engine.
It does this using [jsonnet](http://jsonnet.org)

For documentation of the code, see the godoc directory:

* [cmd](godoc/cmd.md) - how kr8+ processes commands and flags
* [pkg/jvm](godoc/kr8-jsonnet.md) - how kr8+ processes jsonnet
* [pkg/types](godoc/kr8-types.md) - standard types used by kr8+
* [pkg/util](godoc/kr8-util.md) - utility functions used by kr8+

Index

* [Installation](installation.md)
* [Examples](../example)
* [Concepts](concepts/overview.md)
* [Clusters](concepts/clusters.md)
* [Components](concepts/components.md)
* [Helpers](helpers.md) | [Script](../scripts/kr8-helpers)
