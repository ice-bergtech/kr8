// Package types contains shared types used across kr8+ packages.
package types

import (
	"fmt"
)

// An object that stores variables that can be referenced by components.
type Kr8Cluster struct {
	Name string `json:"name"`
	Path string `json:"-"`
}

type Kr8Opts struct {
	// Base directory of kr8 configuration
	BaseDir string
	// Directory where component definitions are stored
	ComponentDir string
	// Directory where cluster configurations are stored
	ClusterDir string
}

// Options for running the jsonnet command.
type CmdJsonnetOptions struct {
	Prune         bool
	Cluster       string
	ClusterParams string
	Component     string
	Format        string
	Color         bool
}

// VMConfig describes configuration to initialize the Jsonnet VM with.
type VMConfig struct {
	// Jpaths is a list of paths to search for Jsonnet libraries (libsonnet files)
	Jpaths []string `json:"jpath" yaml:"jpath"`
	// ExtVars is a list of external variables to pass to Jsonnet VM
	ExtVars []string `json:"ext_str_file" yaml:"ext_str_files"`
	// base directory for the project
	BaseDir string `json:"base_dir" yaml:"base_dir"`
}

type Kr8Error struct {
	Message string
	Value   interface{}
}

// Error implements error.
func (e Kr8Error) Error() string {
	return fmt.Sprintf("%s: %v", e.Message, e.Value)
}
