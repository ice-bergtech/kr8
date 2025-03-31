# Components

A component is a deployable unit that you wish to install in one or more clusters.

Your component might begin life before kr8 in one of a few ways:

  - a [Helm Chart](https://github.com/helm/charts/tree/master/stable)
  - a static [YAML manifest](https://github.com/kubernetes/examples/blob/master/guestbook/all-in-one/guestbook-all-in-one.yaml)
  - some [Jsonnet](https://github.com/coreos/prometheus-operator/tree/master/jsonnet/prometheus-operator)

but they all have something in common - you need to deploy it to multiple clusters with slight differences in configuration.

The root of your component will contain a file named `params.jsonnet` which contains configuration parameters.
Additional files can be stored alongside or in folders.
This is often done to deploy multiple versions of a component at once.

## Params

kr8's most useful feature is the ability to configure _parameters_ for a specific cluster.
It does that by specifying a `params.jsonnet` in each component.

Kr8 extracts core configuration parameters from the `kr8_spec` key.

Fields:

  - `namespace`: String. Required. The namespace the component should be installed in
  - `release_name`: String. Required. Analogous to a helm release - what the component should be called when it's installed into a cluster
  - `enable_kr8_allparams`: Bool. Optional. Includes a full render of all component params.  Used for components that reflect properties of other components.
  - `enable_kr8_allclusters`: Bool. Optional. Includes a full render of all cluster params.  Used for components that reflect properties of other clusters.
  - `disable_output_clean`: Bool. Optional. Purge any yaml files in the output dir that were not generated
  - `jpaths`: List[filename]. Optional. Add additional libjsonnet paths.
  - `extfiles`: List[obj]. Optional.  Add additional files to load as jsonnet `ExtVar`s.  Example: `extfiles: [identifier: "filename.jsonnet", otherfile: "filename2.jsonnet" ]`
  - `includes`: List[string or obj]. Optional. Include and process additional files.  Described more below.

## Referencing files and data

When generating a component, multiple types of files can be combined to generate the final component output.

* If the input is meant to have an output file generated, use `includes`
* If the input is data consumed by the component, use `extfiles`, 
* If it's additional jsonnet libs, use `jpaths`


### includes

The `includes` field allows you to include and process additional files. Each item in the list can be either a string (filename) or an object with specific properties.

When the item is a string, it's treated as a filename to include. The output will be placed in the `generate_dir` with the same name and `.yaml` extension.

When the item is an object, it allows for more customization.
There are the following fields:

* `file`: The filename to include. Required. Allowed extensions: [`jsonnet`, `yaml`, `yml`, `tmpl`, `tpl`]
* `dest_dir`: The directory where the output should be placed. Optional.
* `dest_name`: The name of the output file (without extension). Optional.
* `dest_ext`: The extension of the output file. Optional.

The `file` value must be a `jsonnet`, `yaml`, or `tpl` filename.

<detail>
<summary>Examples of `includes` entries</summary>
  ```jsonnet
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
  Will generate the following files:
  ```
  generate_dir:
    filename.jsonnet
    altname2.jsonnet
    altname3.txt
    altDir:
      altname1.txt
  ``` 
</detail>

### extfiles

`kr8_spec.extfiles: [var_name:"filename.jsonnet"]`

This will load the specified file into the jsonnet vm external vars.
These files can then be referenced in your jsonnet code using the function `std.extVar("var_name")` variable.

It will be availble to be used in component jsonnet as a string, but functions like 

### jpaths

`kr8_spec.jpaths: ["path/to/dir/"]`

The `jpaths` parameter allows you to specify additional paths that kr8 should search for components.
This is useful for component-specific jsonnet libraries.
In most cases, it is better to have a shared library that all components can use, but sometimes it is necessary to have a custom library for a specific component.

Each directory string will be passed to the jsonnet vm during processing.

## Taskfile

A taskfile  within the component directory can be extermely usefule for version management and updates.

  - fetch - a manually run task which downloads all the components' dependencies (for example, helm chart or static manifest)
  - generate - this is the task that's run when kr8 generates the manifest for the cluster


These tasks will be highly dependent on the particular component - for example, a component using a helm chart will generally have a different set of fetch and generate tasks to a component using a static manifest.

An example Taskfile might look like this:

```yaml
version: 2

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

  generate:
    desc: "generate"
    cmds:
      - KR8_COMPONENT={{.KR8_COMPONENT}} kr8-helpers clean-output # clean the tmp directories each time we generate
      - KR8_COMPONENT={{.KR8_COMPONENT}} kr8-helpers helm-render-with-patch "{{.KR8_COMPONENT}}" patches.jsonnet # our generate command, which in this case is a helm-render with some patches in a jsonnet file
```
