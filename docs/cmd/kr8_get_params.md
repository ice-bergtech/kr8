## kr8 get params

Get parameter for components and clusters

### Synopsis

Get parameters assigned to clusters and components in the kr8+ config hierarchy

```
kr8 get params [flags]
```

### Options

```
  -C, --cluster string     get components for cluster
  -c, --component string   component to render params for
  -h, --help               help for params
  -P, --param string       return value of json param from supplied path
```

### Options inherited from parent commands

```
  -B, --base string             kr8+ root configuration directory (default "./")
  -D, --clusterdir string       kr8+ cluster directory
  -p, --clusterparams string    provide cluster params as single file - can be combined with --cluster to override cluster
      --color                   enable colorized output (default true)
  -d, --componentdir string     kr8+ component directory
      --config string           a single config file with kr8+ configuration
      --debug                   log additional information about what kr8+ is doing. Overrides --loglevel
      --ext-str-file key=file   set comma-separated jsonnet extVars from file contents in the format key=file
  -J, --jpath stringArray       additional jsonnet library directories
  -L, --loglevel string         set zerolog log level (default "info")
      --parallel int            parallelism - defaults to runtime.GOMAXPROCS(0) (default -1)
      --profiledir string       directory to write pprof profile data to
```

### SEE ALSO

* [kr8 get](kr8_get.md)	 - Display one or many kr8+ resources

###### Auto generated by spf13/cobra on 7-May-2025
