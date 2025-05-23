# Native Functions

Additional functions have been added to the jsonnet vm to add functionality.
They are able to be called from jsonnet using `std.native('funcName')`, where `funcName` is the name of the function.

Native function definitions: [pkg/nativefuncs](https://github.com/ice-bergtech/kr8/blob/master/pkg/nativefuncs/nativefuncs.go).

## template

Templates the passed in input `str` using the json string `config`.
Config is unmarshaled into a json object and passed to the template engine.
The template engine used is sprig - [Template Documentation](https://masterminds.github.io/sprig/).
The resulting string is returned.

Usage:

```go
std.native("template")(config json, str string) (string)
```

Example:

```go
local templateOutput = std.native("template")(config.data, "Hello {{ .Name }}");
```


## helmTemplate

Provides the same `Helm.Template` functionality as the `grafana/tanka` package. 
Charts are required to be present on the local filesystem, at a relative location to the file that calls `helm.template()` / `std.native('helmTemplate')`. 
This guarantees hermeticity.
Does not use sprig for templating.

Usage:

```go
std.native("helmTemplate")(name string, chart string, opts TemplateOpts) (manifest.List)
```

Example:

```go
local helm_template = std.native("helmTemplate")(config.release_name, "./vendor/"+config.chart_version, {
    calledFrom: std.thisFile,
    namespace: config.namespace,
    values: config.helm_values,
});

[
    object
    for object in std.objectValues(helm_template)
    if "kind" in object && object.kind != "Secret"
]
```

* Template Opts: [godocs grafana/tanka](https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L65)
* Function Source: [grafana/tanka v0.27.1](https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23)


## escapeStringRegex

Uses `regexp.QuoteMeta` to escape a string for use in a regular expression.

Usage:

```go
std.native("escapeStringRegex")(str string) (string)
```

Example:

```go
local clean_string = std.native("escapeStringRegex")(config.knarly_string);
```

## regexMatch

Uses `regexp.MatchString` to check if a string matches a regular expression.

Usage:

```go
std.native("regexMatch")(regex string, str string) (bool)
```

Example:

```go
// check if a string is numbers
if std.native("regexMatch")("\d+", config.thing) then config.thing else ""
```

## regexSubst

Uses `regexp.ReplaceAllString` to replace all occurrences of a regular expression in a string.

Usage:

```go
std.native("regexSubst")(regex string, src string, repl string) (string)
```

Example:

```go
local stringVar = std.native("regexSubst")("\d", config.thing, "<num>");
```
