# Roadmap

**kr8+** is mostly feature complete, and should be able to be used in production environments.

However, there are always things to improve:

* Refinement/documentation of kr8 configuration
* Improve performance of rendering components
* Improve user input error handling and sharing
* Refine documentation for better understanding and usage, especially around cluster `kr8_spec` and component `kr8_component_spec`
* Enable configuration of formatting
* Add additional linting to format command
* Add ability to fetch or coordinate fetching remote resources (e.g. Helm charts, CRDs etc.)
* Improve examples and tutorials for better onboarding
* Identify common tasks and integrate into **kr8+** or a jsonnet libsonnet library
* integrate `go-task` for preparing remote resource definition fetching scripts for user
* Build caching
* Create configuration plan before generating clusters. Create a tree of based on cluster config files, and deep copy jvm each layer - difficult because of how JVM is setup and golang deep copy abilities.  Perform breadth-first "search" until all leafs (`cluster.jsonnet` files) are generated.  Need different solution for component jvm management.  Only thread pool component generation but process cluster serially.
* Add way to output cluster/component jsonnet AST for outside analysis
