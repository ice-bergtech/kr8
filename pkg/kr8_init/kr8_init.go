package kr8_init

// Kr8InitOptions defines the options used by the init subcommands.
type Kr8InitOptions struct {
	// URL to fetch the skeleton directory from
	InitUrl string
	// Name of the cluster to initialize
	ClusterName string
	// Name of the component to initialize
	ComponentName string
	// Type of component to initialize (e.g. jsonnet, yml, chart, compose)
	ComponentType string
	// Determines whether to run in interactive mode
	Interactive bool
	// Determines whether to fetch remote resources
	Fetch bool
}
