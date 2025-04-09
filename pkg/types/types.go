package types

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
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

// The specification for a clusters.jsonnet file.
// This describes configuration for a cluster that kr8 should process.
type Kr8ClusterJsonnet struct {
	// kr8 configuration for how to process the cluster
	ClusterSpec Kr8ClusterSpec `json:"_kr8_spec"`
	// Cluster Level configuration that components can reference
	Cluster Kr8Cluster `json:"_cluster"`
	// Distictly named components.
	Components map[string]Kr8ClusterComponentRef `json:"_components"`
}

// A reference to a component folder that contains a params.jsonnet file.
// This is used in the cluster jsonnet file to reference components.
type Kr8ClusterComponentRef struct {
	// The path to a component folder that contains a params.jsonnet file
	Path string `json:"path"`
}

// The specification for how to process a cluster.
// This is used in the cluster jsonnet file to configure how kr8 should process the cluster.
type Kr8ClusterSpec struct {
	// The name of the cluster
	Name string `json:"-"`
	// A jsonnet function that each output entry is processed through. Default `function(input) input`
	PostProcessor string `json:"postprocessor"`
	// The name of the root generate directory. Default `generated`
	GenerateDir string `json:"generate_dir"`
	// if this is true, we don't use the full file path to generate output file names
	GenerateShortNames bool `json:"generate_short_names"`
	// if this is true, we prune component parameters
	PruneParams bool `json:"prune_params"`
	// Additional information used to process the cluster that is not stored with it.
	// Cluster output directory
	ClusterDir string `json:"-"`
}

// This function creates a Kr8ClusterSpec from passed params.
// If genDirOverride is empty, the value of generate_dir from the spec is used.
func CreateClusterSpec(
	clusterName string,
	spec gjson.Result,
	kr8Opts Kr8Opts,
	genDirOverride string,
) (Kr8ClusterSpec, error) {
	// First determine the value of generate_dir from the command line args or spec.
	clGenerateDir := genDirOverride
	if clGenerateDir == "" {
		clGenerateDir = spec.Get("generate_dir").String()
	}
	if clGenerateDir == "" {
		log.Warn().Msg("generate_dir should be set in cluster parameters or passed as generate-dir flag.  Defaulting to ./")
		clGenerateDir = "generated"
	}
	// if generateDir does not start with /, then it goes in baseDir
	if !strings.HasPrefix(clGenerateDir, "/") {
		clGenerateDir = filepath.Join(kr8Opts.BaseDir, clGenerateDir)
	}
	clusterDir := filepath.Join(clGenerateDir, clusterName)
	log.Debug().Str("cluster", clusterName).Msg("output directory: " + clusterDir)

	return Kr8ClusterSpec{
		PostProcessor:      spec.Get("postprocessor").String(),
		GenerateDir:        clGenerateDir,
		GenerateShortNames: spec.Get("generate_short_names").Bool(),
		PruneParams:        spec.Get("prune_params").Bool(),
		ClusterDir:         clusterDir,
		Name:               clusterName,
	}, nil
}

// The specification for component's params.jsonnet file.
// It contains all the configuration and variables used to generate component resources.
// This configuration is often modified from the cluster config to add cluster-specific configuration.
type Kr8ComponentJsonnet struct {
	// Component-specific configuration for how kr8 should process the component (required)
	Kr8Spec Kr8ComponentSpec `json:"kr8_spec"`
	// The default namespace to deploy the component to
	Namespace string `json:"namespace"`
	// A unique name for the component
	ReleaseName string `json:"release_name"`
	// Component version string (optional)
	Version string `json:"version"`
	// Relative directory where the component's resources are located (required).
	// Usually std.thisFile.
	CalledFrom string `json:"called_from"`
}

// The kr8_spec object in a cluster config file.
// This configures how kr8 processes the component.
type Kr8ComponentSpec struct {
	// If true, includes the parameters of the current cluster when generating this component
	Kr8_allparams bool `json:"enable_kr8_allparams"`
	// If true, includes the parameters of all other clusters when generating this component
	Kr8_allclusters bool `json:"enable_kr8_allclusters"`
	// If false, all non-generated files present in the output directory will be removed
	DisableOutputDirClean bool `json:"disable_output_clean"`
	// A list of filenames to include as jsonnet vm external vars
	ExtFiles ExtFileVar `json:"extfiles"`
	// Additional jsonnet libs to the jsonnet vm, component-path scoped
	JPaths []string `json:"jpaths"`
	// A list of filenames to include and output as files
	Includes []interface{} `json:"includes"`
}

// Extract jsonnet extVar defintions from spec.
func ExtractExtFiles(spec gjson.Result) map[string]string {
	result := make(map[string]string)
	for k, v := range spec.Get("extfiles").Map() {
		if v.Type == gjson.String {
			result[k] = v.String()
		}
	}

	return result
}

// Extract jsonnet lib paths from spec.
func ExtractJpaths(spec gjson.Result) []string {
	jPathsInput := spec.Get("jpaths").Array()
	jPathsOutput := make([]string, len(jPathsInput))
	for i, p := range jPathsInput {
		jPathsOutput[i] = p.String()
	}

	return jPathsOutput
}

// Extract jsonnet includes filenames or objects from spec.
func ExtractIncludes(spec gjson.Result) []interface{} {
	incl := spec.Get("includes")
	includes := make([]interface{}, len(incl.Array()))
	for idx, item := range incl.Array() {
		switch item.Type {
		case gjson.JSON:
			var include Kr8ComponentSpecIncludeObject
			err := json.Unmarshal([]byte(item.String()), &include)
			if err != nil {
				log.Info().Msg("Error unmarshalling include object: " + item.String() + " skipping.")

				continue
			}
			includes[idx] = include
		case gjson.String:
			includes[idx] = item.String()
		case gjson.True:
		case gjson.False:
		case gjson.Number:
		case gjson.Null:
		default:
			log.Error().Msg("Includes list item is not a string or json object, skipping")

			continue
		}
	}

	return includes
}

// Extracts a component spec from a jsonnet object.
func CreateComponentSpec(spec gjson.Result) (Kr8ComponentSpec, error) {
	specM := spec.Map()
	log.Debug().Msg(spec.String())
	// spec is missing?
	if len(specM) == 0 {
		log.Error().Msg("Component has no `kr8_spec` object")
		// intetionally create an error to return
		return Kr8ComponentSpec{},
			Kr8Error{Message: "Component has no `kr8_spec` object", Value: ""}
	}

	log.Debug().Msg("Component spec: " + spec.Str)

	componentSpec := Kr8ComponentSpec{
		Kr8_allparams:         spec.Get("enable_kr8_allparams").Bool(),
		Kr8_allclusters:       spec.Get("enable_kr8_allclusters").Bool(),
		DisableOutputDirClean: spec.Get("disable_output_clean").Bool(),
		ExtFiles:              ExtractExtFiles(spec),
		JPaths:                ExtractJpaths(spec),
		Includes:              ExtractIncludes(spec),
	}

	return componentSpec, nil
}

// Map of external files to load into jsonnet vm as external variables.
// Keys are the variable names, values are the paths to the files to load as strings into the jsonnet vm.
// To reference the variable in jsonnet code, use std.extvar("variable_name").
type ExtFileVar map[string]string

// A struct describing an included file that will be processed to produce a file.
type Kr8ComponentSpecIncludeFile interface {
	string
	Kr8ComponentSpecIncludeObject
}

// An includes object which configures how kr8 includes an object.
// It allows configuring the included file's destination directory and file name.
// The input file will be processed differently depending on the filetype.
type Kr8ComponentSpecIncludeObject struct {
	// an input file to process
	// accepted filetypes: .jsonnet .yml .yaml .tmpl .tpl
	File string `json:"file"`
	// handle alternate output directory for file
	DestDir string `json:"dest_dir,omitempty"`
	// override destination file name
	DestName string `json:"dest_name,omitempty"`
	// override destination file extension
	DestExt string `json:"dest_ext,omitempty"`
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
