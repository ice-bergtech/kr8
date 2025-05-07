# Features

* **Cluster Configuration Management**: Manage Kubernetes clusters across environments, regions and platforms with a declarative and centralized approach.
  * Enabled easy management of multiple (10+) clusters with various different configurations
  * Configuration is flattened from cluster directory leaf to root into a single config file, enabling easy application of DRY principals
  * Able to provide all component configs to single component.  Useful for components that generate something for each cluster component, such as network policy, argo applications, or documentation.
  * Able to provide all cluster configs to components.  Useful for cross-cluster monitoring.
  * Generate static configuration in a repeatable, idempotent manner.
* **Opinionated Structure**: Enforces best practices for consistent and reliable cluster configurations.
  * Components are stored separately in their own folders, allowing for easy management and version control.
  * Able to define multiple versions of components with different configuration
  * Configuration files are written in YAML or JSON, ensuring compatibility with existing tools and workflows.
  * Ensures consistency across Kubernetes clusters, reducing errors and improving maintainability.
  * Provides a straightforward approach to complex Kubernetes configurations, making it easier for teams to adopt.
  * Designed to support clusters of all sizes, from simple single-node setups to large-scale production environments.
* **Jsonnet Native Functions**: Use jsonnet to render and override component config from multiple sources
  * Go-templates: Able to output text files templated based off of component configuration.  Integrated with sprig templating functions
  * Docker-compose: Able to process docker-compose as yaml, or through [kubernetes/kompose]() to output kubernetes resources
  * Kustomize: Able to process kustomize files and output kubernetes resources
  * Helm: Able to process locally stored helm charts and output kubernetes resources in deterministic way.
  * URL Parsing: Able to parse URLs into objects that can be used in component configuration.
  * IP Address Manipulation: Able to manipulate IP addresses and CIDRs in component configuration.
  * Jsonnet Std.lib: Use jsonnet std.lib functions to manipulate data and perform operations on component configuration.
* **Extensibility**: Easily extensible to meet the needs of diverse Kubernetes environments.
  * Use jsonnet libraries the same way you would use any other jsonnet library.
  * Output a variety of structured and unstructured files.
* **Build Cache**
  * Keeps cache of cluster-level config, which if invalidated triggers a full generate
  * Keeps cache of cluster-component config. If cluster-level component config changes the component is re-generated.
  * Keeps cache of per-component file hashes.
  * Stores cache per-cluster, in a json file named `.kr8_cache`
* **CI/CD Friendly**: Statically define all your configuration in a single source of truth, making it easy to integrate with CI/CD pipelines and deployment automation like ArgoCD.
  * Easily create reproducible builds by using the same configuration across different environments.
  * Fully version control charts and normally-remote configurations in a single place.
  * Store generated secrets via [SealedSecrets](https://github.com/bitnami-labs/sealed-secrets)
