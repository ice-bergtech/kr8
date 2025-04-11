# Components

A component is a deployable unit that you wish to install in one or more clusters.
Components can be declared multiple times within a cluster, as long as they are named distinctly when loaded.

In a kr8+ project, components are defined in `./components` by default, or the directory specified by the `--componentdir`, `-X` flags.

The component spec is defined as a golang struct: [docs/godoc/kr8-types.md#Kr8ComponentJsonnet]

Your component might begin life before kr8+ in one of a few ways:

  - a [Helm Chart](https://github.com/helm/charts/tree/master/stable)
  - a static [Kubernetes YAML/Json manifest](https://github.com/kubernetes/examples/blob/master/guestbook/all-in-one/guestbook-all-in-one.yaml)
  - [Jsonnet](https://github.com/coreos/prometheus-operator/tree/master/jsonnet/prometheus-operator) describing a deployable unit
  - template files, to generate arbitrary files from a golang-style template.
  - a docker image or script to deploy

The root directory of a component contain a file named `params.jsonnet`.
containing configuration parameters consumed by kr8+ and passed to your Jsonnet code.
Additional files can be stored alongside or in folders.
This is often done to deploy multiple versions of a component at once.

At a minimum a component directory contains:

* **params.jsonnet**: Contains component configuration.
* **includes files**: Files that are processed and placed in the generated output directory

Often there will be additional files

* **vendor/**: Contains versioned vendor files, with each version in it's own directory
* **Taskfile.yml**: Automate common component tasks to be ran to prepare component for being referenced by a cluster.


## params.jsonnet

kr8+'s most useful feature is the ability to easily layer _parameters_ to generate a resource.
The `params.jsonnet` file in each component can be updated at the cluster level, making it simple to customize the behavior of your component across different environments.

```yaml
# params.jsonnet
{
  release_name: 'cert-manager',
  namespace: 'cert-manager',
  kr8_spec: { ... },
  helm_values: {
    webhook: { enabled: false },  # this is a value for the helm chart
  },
}
```


kr8+ extracts core configuration parameters from the `kr8_spec` key.
The `kr8_spec` fields are marshaled into the [Kr8ComponentSpec](../godoc/kr8-types.md#Kr8ComponentSpec) struct.

| Field                    | Description                                                                                                                                                        | Example                                                                                                               |
| ------------------------ | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------------------- |
| `namespace`              | String. Required. The primary namespace the component should be installed in                                                                                       | `'default'`, `'argocd'`                                                                                               |
| `release_name`           | String. Required. Analogous to a helm release - what the component should be called when installed into a cluster                                                  | `'argo-workflows'`                                                                                           |
| `enable_kr8_allparams`   | Bool. Optional, default `False`. Includes a full render of all component params during generate.  Used for components that reflect properties of other components. | `False`, `True`                                                                                                       |
| `enable_kr8_allclusters` | Bool. Optional, default `False`. Includes a full render of all cluster params during generate.  Used for components that reflect properties of other clusters.     | `False`, `True`                                                                                                       |
| `disable_output_clean`   | Bool. Optional, default `False`. If true, stops kr8+ from removing all yaml files in the output dir that were not generated                                         | `False`, `True`                                                                                                       |
| `includes`               | List[string or obj]. Optional, default `[]`. Include and process additional files.  Described more below.                                                          | `["kube.jsonnet", {file: "resource.yaml", dest_name: "asdf"}, {file: "docs.tpl", dest_dir: "docs", dest_ext: ".md"}]` |
| `extfiles`               | {fields}. Optional, default `{}`.  Add additional files to load as jsonnet `ExtVar`s.  The field key is used as the variable name, and the value is the file path. | `{identifier: "filename.txt", otherfile: "filename2.json" }`                                                          |
| `jpaths`                 | List[string]. Optional, default `[]`. Add additional libjsonnet paths with base dir `/baseDir/componentPath/`. The path `baseDir + "/lib"` is always included.     | `["vendor/argo-libsonnet/"]`                                                                                          |


## Referencing files and data

When generating a component, multiple types of files can be combined to generate the final component output.

* If the input is meant to have an output file generated, use `includes`
* If the input is data consumed by the component, use `extfiles`, 
* If it's additional jsonnet libs, use `jpaths`

## includes

The `includes` field allows you to include and process additional files.
Each item in the list can be either a string (filename) or an object with specific properties.

When the item is a string, it's treated as a filename to include.
The output will be placed in the `generate_dir` with the same name and `.yaml` extension.

When the item is an object, it allows for more customization.
The item is marshaled into a [Kr8ComponentSpecIncludeObject](../godoc/kr8-types.md#Kr8ComponentSpecIncludeObject) struct.

There are the following fields:

* `file`: The filename to include. Required. Allowed extensions: [`jsonnet`, `yaml`, `yml`, `tmpl`, `tpl`]
* `dest_dir`: The directory where the output should be placed. Optional.
* `dest_name`: The name of the output file (without extension). Optional.
* `dest_ext`: The extension of the output file. Optional.

The `file` value must be a `jsonnet`, `yaml`, or `tpl` type file.

For example, the `includes` entries:
  
```yaml
includes: [
  "filename.jsonnet",
  {
    file: "filename.jsonnet",
    dest_dir: "altDir",
    dest_name: "altname1",
    dest_ext: "txt"
  },
  {
    file: "filename.jsonnet",
    dest_name: "altname2",
  },
  {
    file: "filename.jsonnet",
    dest_name: "altname3",
    dest_ext: "txt"
  },
]
```

Generates the files:

```sh
generate_dir:
├── filename.jsonnet
├── altname2.jsonnet
├── altname3.txt
└── altDir:
    └── altname1.txt
``` 

### extfiles

`kr8_spec.extfiles: [var_name:"filename.jsonnet"]`

This will load the specified file into the jsonnet vm external vars.
These files can then be referenced in your jsonnet code using the function `std.extVar("var_name")` variable.

It will be available to be used in component jsonnet as a string, but functions like 

### jpaths

`kr8_spec.jpaths: ["path/to/dir/"]`

The `jpaths` parameter allows you to specify additional paths that kr8+ should search for components.
This is useful for component-specific jsonnet libraries.
In most cases, it is better to have a shared library that all components can use, but sometimes it is necessary to have a custom library for a specific component.

Each directory string will be passed to the jsonnet vm during processing.

## Taskfile

A taskfile within the component directory can help manage the lifecycle of components, especially when dealing with dependencies and version management for more complex updates.
A common practice is to create a `fetch` task that downloads all dependencies (e.g., Helm charts or static manifests).
It can also perform other preparation steps like removing files or resources that are not required.

These tasks are highly dependent on the particular component - for example, a component using a helm chart generally have a different set of fetch and generate tasks to a component using a static manifest.

An example Taskfile might look like this:

```yaml
version: 3

vars:
  KR8_COMPONENT: kubemonkey

tasks:
  fetch:
    desc: "fetch component kubemonkey"
    cmds:
      - curl -L https://github.com/asobti/kube-monkey/tarball/master > kubemonkey.tar.gz # download the local helm chart from the git repo
      - tar --strip-components=2 -xzvf kubemonkey.tar.gz asobti-kube-monkey-{{.GIT_COMMIT}}/helm # extract it
      - mv kubemonkey charts # place it in a charts directory
      - rm -fr *.tar.gz # remove the tar.gz from the repo
    vars:
      GIT_COMMIT:
        sh: curl -s https://api.github.com/repos/asobti/kube-monkey/commits/master | jq .sha -r | xargs git rev-parse --short
```

## Common Component Types

### Jsonnet

A very simple component might just be a few lines of jsonnet.

Consider the situation whereby you might have two clusters, one in AWS and one in DigitalOcean.
You need to set a default storageclass.
You could do this with jsonnet.

This example is located in [kr8-examples jsonnetStorageClasses](https://github.com/ice-bergtech/kr8-examples/components/doc-concepts/jsonnetStorageClasses/).

Your jsonnet component would look like this:

```bash
components/doc-conepts/jsonnetStorageClasses
├── params.jsonnet
└── storageclasses.jsonnet
```

#### Params

As a reminder, every component requires a params file.
We need to set a namespace for the component, even though it's a cluster level resource - namespace is a required paramater for kr8+.
We also need to tell kr8+ what files to include for the component:

```yaml
{
  namespace: 'kube-system',
  release_name: 'storageclasses',
  kr8_spec: {
    includes: [ "storageClasses.jsonnet" ],
    extfiles: [
      {"echoManifest": "./vendor/" + version}
    ]
  },
}
```

#### Jsonnet Manifest

Your jsonnet manifest looks like this:

```go
// a jsonnet external variable from kr8 that gets contains cluster-level configuration
local kr8_cluster = std.extVar("kr8_cluster");

// a jsonnet function for creating a storageclass
local StorageClass(name, type, default=false) = {
  apiVersion: "storage.k8s.io/v1",
  kind: "StorageClass",
  metadata: {
    name: name,
    annotations: {
      "storageclass.kubernetes.io/is-default-class": if default then "true" else "false",
    },
  },
  parameters: {
    type: type,
  },
};

// check the cluster configuration for a type, if it's AWS make a gp2 type storageclass
if kr8_cluster.cluster_type == "aws" then std.objectValues(
  {
    // default gp2 storage class, not tied to a zone
    ebs_gp2: StorageClass("gp2", "gp2", true) {},
  }
) else [] // don't make a storageclass if it's not AWS
```

### YAML Component

kr8+ can use a static k8s manifest as a source input.
You can then manipulate the structure of that YAML using Jsonnet.
kr8+ takes care of the heavy lifting for you.

This is useful when loading manifests from remote source, and you want to manipulate the resources before deploying into a cluster.

An example can be found here: [kr8-examples echo-test](https://github.com/ice-bergtech/kr8-examples/components/doc-concepts/echo-test/)

```bash
components/doc-conepts/jsonnetStorageClasses
├── Taskfile.yml
├── params.jsonnet
├── echo.jsonnet
└── vendor
    └── {{.Version}}
        └── echo.yml
```

#### Taskfile

You'll need a taskfile that downloads the original manifests for you in the `fetch` task.
Here's an example:

```yaml
# https://taskfile.dev/usage
version: '3'

vars:
  NAME: external-dns
  REPO: 'https://kubernetes-sigs.github.io/external-dns/'

tasks:
  default:
    cmds:
      - task: fetch-1.15

  fetch-1.15:
    desc: "fetch component dependencies"
    vars:
      THING_VER: v1.15.0
    cmds:
      - task: fetch-chart
        vars: { VER: '{{.THING_VER}}'}

  fetch-1.14:
    desc: "fetch component dependencies"
    vars:
      THING_VER: v1.14.0
    cmds:
      - task: fetch-chart
        vars: { VER: '{{.THING_VER}}'}

  fetch-chart:
    desc: "fetch a helm chart"
    vars:
      VER: '{{default "unset" .VER}}'
    cmds:
      - mkdir -p ./vendor/{{.NAME}}-{{.VER}} && rm -rf ./vendor/{{.NAME}}-{{.VER}}/*
      - do-fetch and place resulting files in ./vendor/{{.NAME}}-{{.VER}}/ 
```

### Helm

**kr8+** can render helm charts locally and inject parameters as helm values.
This provides a great degree of flexibility when installing components into clusters.

Working examples of using helm charts with kr8+ can be found in the [kr8-examples](https://github.com/ice-bergtech/kr8-examples/) repo.

### Processing a chart

#### Params and chart config

The `params.jsonnet` for a helm chart directory should include the helm values you want to use.
A values file is usually a required file for a helm component.
The values can be stored in a field, usually named `helm_values`.

Here's an example:

```yaml
# params.jsonnet
{
  ...
  # Defined a field to contains helm values
  helm_values: {
    webhook: { enabled: false },  # this is a value for the helm chart
  },
}
```

The chart will be referenced and processed through 

```go
local config = std.extVar("kr8");

local helm_template = std.native("helmTemplate")(config.release_name, "./vendor/"+"external-dns-"+config.chart_version, {
    calledFrom: std.thisFile,
    namespace: config.namespace,
    values:  if "helm_values" in config then config.helm_values else {},
});

std.objectValues(helm_template)
```

The values can also be stored in a separate file, and referenced via `includes` or `extfiles` configuration.

#### Patches

There are certain situations where a configuration option is not available for a helm chart.
For example, you might want to add an argument that hasn't quite made it into the helm chart yet, or add something like [pod affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#inter-pod-affinity-and-anti-affinity) where it isn't actually a value option in a helm chart.

kr8+ helps you in this situation by providing a mechanism to patch a helm chart.
After loading the chart through the 

Here's an example `patches.jsonnet` for [external-dns](https://github.com/kubernetes-incubator/external-dns)

```go
local config = std.extVar("kr8");

local helm_template = std.native("helmTemplate")(config.release_name, "./vendor/"+"external-dns-"+config.chart_version, {
    calledFrom: std.thisFile,
    namespace: config.namespace,
    values:  if "helm_values" in config then config.helm_values else {},
});

[
    object { metadata+: { labels+: note: "hello", other_thing: "yes" } }
    for object in std.objectValues(helm_template)
    if "kind" in object && object.kind != "Secret" // stop from committing secrets to repo
]
```

#### Taskfile

This is used to codify the initial chart fetch, which needs to be performed to fetch new updates and changes.
An example taskfile for a helm chart might look like this:

```yaml
# https://taskfile.dev/usage
version: '3'

vars:
  CHART_NAME: external-dns
  CHART_REPO: 'https://kubernetes-sigs.github.io/external-dns/'

tasks:
  default:
    cmds:
      - task: fetch-1.15

  fetch-1.15:
    desc: "fetch component dependencies"
    vars:
      CHART_VER: v1.15.0
    cmds:
      - task: fetch-chart
        vars: { VER: '{{.CHART_VER}}'}

  fetch-1.14:
    desc: "fetch component dependencies"
    vars:
      CHART_VER: v1.14.0
    cmds:
      - task: fetch-chart
        vars: { VER: '{{.CHART_VER}}'}

  fetch-chart:
    desc: "fetch a helm chart"
    vars:
      VER: '{{default "unset" .VER}}'
    cmds:
      - mkdir -p ./vendor/{{.CHART_NAME}}-{{.VER}} && rm -rf ./vendor/{{.CHART_NAME}}-{{.VER}}/*
      - mkdir -p ./vendor/tmp && rm -rf ./vendor/tmp/*
      # add the helm repo and fetch it locally into vendor directory
      - helm fetch --repo {{.CHART_REPO}} --untar --untardir ./vendor/tmp --version "{{.VER}}" "{{.CHART_NAME}}"
      - mv ./vendor/tmp/{{.CHART_NAME}}/* ./vendor/{{.CHART_NAME}}-{{.VER}}/ && rm -rf ./vendor/tmp
```

