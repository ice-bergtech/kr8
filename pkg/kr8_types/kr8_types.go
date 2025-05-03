// Package kr8_types defines the structure for kr8+ cluster and component resources.
package kr8_types

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/ice-bergtech/kr8/pkg/types"
	"github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

// An object that stores cluster-level variables that can be referenced by components.
type Kr8Cluster struct {
	// The name of the cluster.
	// Derived from folder containing the cluster.jsonnet.
	// Not read from config.
	Name string `json:"-"`
	// Path to the cluster folder.
	// Not read from config.
	Path string `json:"-"`
}

// The specification for a clusters.jsonnet file.
// This describes configuration for a cluster that kr8 should process.
type Kr8ClusterJsonnet struct {
	// kr8+ configuration for how to process the cluster
	ClusterSpec Kr8ClusterSpec `json:"_kr8_spec"`
	// Cluster Level configuration that components can reference
	Cluster Kr8Cluster `json:"_cluster"`
	// Distinctly named components.
	Components map[string]Kr8ClusterComponentRef `json:"_components"`
}

// A reference to a component folder that contains a params.jsonnet file.
// This is used in the cluster jsonnet file to reference components.
type Kr8ClusterComponentRef struct {
	// The path to a component folder that contains a params.jsonnet file
	Path string `json:"path"  jsonschema:"example=components/service"`
}

// The specification for how to process a cluster.
// This is used in the cluster jsonnet file to configure how kr8 should process the cluster.
type Kr8ClusterSpec struct {
	// A jsonnet function that each output entry is processed through. Default `function(input) input`
	PostProcessor string `json:"postprocessor,omitempty" jsonschema:"default=function(input) input"`
	// The name of the root generate directory. Default `generated`
	GenerateDir string `json:"generate_dir,omitempty" jsonschema:"default=generated"`
	// If true, we don't use the full file path to generate output file names
	GenerateShortNames bool `json:"generate_short_names,omitempty" jsonschema:"default=false"`
	// If true, we prune component parameters
	PruneParams bool `json:"prune_params,omitempty" jsonschema:"default=false"`
	// If true, kr8 will store and reference a cache file for the cluster.
	EnableCache bool `json:"cache_enable,omitempty" jsonschema:"default=false"`
	// If true, kr8 will compress the cache in a gzip file instead of raw json.
	CompressCache bool `json:"cache_compress,omitempty"  jsonschema:"default=true"`
	// The name of the cluster
	// Not read from config.
	Name string `json:"-"`
	// Cluster output directory
	// Not read from config.
	ClusterOutputDir string `json:"-"`
}

// This function creates a Kr8ClusterSpec from passed params.
// If genDirOverride is empty, the value of generate_dir from the spec is used.
func CreateClusterSpec(
	clusterName string,
	spec gjson.Result,
	kr8Opts types.Kr8Opts,
	genDirOverride string,
	logger zerolog.Logger,
) (Kr8ClusterSpec, error) {
	// First determine the value of generate_dir from the command line args or spec.
	clGenerateDir := genDirOverride
	if clGenerateDir == "" {
		clGenerateDir = spec.Get("generate_dir").String()
	}
	if clGenerateDir == "" {
		logger.Warn().
			Msg("`generate_dir` should be set in cluster parameters or passed as generate-dir flag." +
				"Defaulting to ./",
			)
		clGenerateDir = "generated"
	}
	// if generateDir does not start with /, then it goes in baseDir
	if !strings.HasPrefix(clGenerateDir, "/") {
		clGenerateDir = filepath.Join(kr8Opts.BaseDir, clGenerateDir)
	}
	clusterDir := filepath.Join(clGenerateDir, clusterName)
	logger.Debug().Str("cluster", clusterName).Msg("output directory: " + clusterDir)

	// Default to compressing the cache
	compress := true
	compressVar := spec.Get("cache_compress")
	if compressVar.Exists() {
		compress = compressVar.Bool()
	}

	return Kr8ClusterSpec{
		PostProcessor:      spec.Get("postprocessor").String(),
		GenerateDir:        clGenerateDir,
		GenerateShortNames: spec.Get("generate_short_names").Bool(),
		PruneParams:        spec.Get("prune_params").Bool(),
		EnableCache:        spec.Get("cache_enable").Bool(),
		CompressCache:      compress,
		ClusterOutputDir:   clGenerateDir + "/" + clusterName,
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
	Version string `json:"version,omitempty"`
}

// The kr8_spec object in a cluster config file.
// This configures how kr8 processes the component.
type Kr8ComponentSpec struct {
	// If true, includes the parameters of the current cluster when generating this component
	Kr8_allParams bool `json:"enable_kr8_allparams,omitempty"`
	// If true, includes the parameters of all other clusters when generating this component
	Kr8_allClusters bool `json:"enable_kr8_allclusters,omitempty"`
	// If false, all non-generated files present in the output directory are removed
	DisableOutputDirClean bool `json:"disable_output_clean,omitempty"`
	// If true, component will not be cached if cluster caching is enabled.
	DisableCache bool `json:"disable_cache,omitempty"`
	// A list of filenames to include as jsonnet vm external vars
	ExtFiles ExtFileVar `json:"extfiles,omitempty"`
	// Additional jsonnet libs to the jsonnet vm, component-path scoped
	JPaths []string `json:"jpaths,omitempty"`
	// A list of filenames to include and output as files
	Includes Kr8ComponentSpecIncludes `json:"includes"`
}

// Extract jsonnet extVar definitions from spec.
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
func ExtractIncludes(spec gjson.Result) (Kr8ComponentSpecIncludes, error) {
	incl := spec.Get("includes")
	includes := Kr8ComponentSpecIncludes{}
	if incl.String() == "" {
		return includes, nil
	}

	err := json.Unmarshal([]byte(incl.String()), &includes)

	return includes, util.ErrorIfCheck("Error unmarshaling include object: "+incl.String(), err)
}

// Extracts a component spec from a jsonnet object.
func CreateComponentSpec(spec gjson.Result, logger zerolog.Logger) (Kr8ComponentSpec, error) {
	specM := spec.Map()
	logger.Debug().Msg(spec.String())
	// spec is missing?
	if len(specM) == 0 {
		// create an error to return
		return Kr8ComponentSpec{},
			types.Kr8Error{Message: "Component has no `kr8_spec` object", Value: ""}
	}

	logger.Debug().Msg("Component spec: " + spec.Str)

	includes, err := ExtractIncludes(spec)
	if err != nil {
		return Kr8ComponentSpec{},
			types.Kr8Error{Message: "Component includes are malformed", Value: err}
	}

	componentSpec := Kr8ComponentSpec{
		Kr8_allParams:         spec.Get("enable_kr8_allparams").Bool(),
		Kr8_allClusters:       spec.Get("enable_kr8_allclusters").Bool(),
		DisableOutputDirClean: spec.Get("disable_output_clean").Bool(),
		ExtFiles:              ExtractExtFiles(spec),
		JPaths:                ExtractJpaths(spec),
		Includes:              includes,
		DisableCache:          false,
	}

	return componentSpec, nil
}

// Map of external files to load into jsonnet vm as external variables.
// Keys are the variable names, values are the paths to the files to load as strings into the jsonnet vm.
// To reference the variable in jsonnet code, use std.extvar("variable_name").
type ExtFileVar map[string]string

// An includes object which configures how kr8 includes an object.
// It allows configuring the included file's destination directory and file name.
// The input files are processed differently depending on the filetype.
type Kr8ComponentSpecIncludeObject struct {
	// An input file to process.
	// Accepted filetypes: .jsonnet .yml .yaml .tmpl .tpl
	File string `json:"file" jsonschema:"example=file.jsonnet,example=.yml,example=template.tpl"`
	// Handle alternate output directory for file.
	// Relative from component output dir.
	DestDir string `json:"dest_dir,omitempty"`
	// Override destination file name
	DestName string `json:"dest_name,omitempty" jsonschema:"default=File field"`
	// Override destination file extension
	// Useful for setting template file extension.
	DestExt string `json:"dest_ext,omitempty" jsonschema:"example=md,example=txt,default=yml"`
	// Override config passed to the includes template file processing.
	// Useful for generating a list of includes in a loop:
	// `[{File: f, Config: data[f]} for f in list]`
	Config string `json:"config,omitempty"`
}

// Define Kr8ComponentSpecIncludes to handle dynamic decoding.
type Kr8ComponentSpecIncludes []Kr8ComponentSpecIncludeObject

// Implement custom unmarshaling for dynamic decoding.
func (k *Kr8ComponentSpecIncludes) UnmarshalJSON(data []byte) error {
	// Check if the data is a single string
	if data[0] == '"' { // JSON strings start with a double quote
		var file string
		if err := json.Unmarshal(data, &file); err != nil {
			return err
		}
		// strip extension from file
		ext := filepath.Ext(file)
		fileName := strings.TrimSuffix(file, ext)
		// Add a default Kr8ComponentSpecIncludeObject using the string as the file
		*k = append(*k, Kr8ComponentSpecIncludeObject{
			File:     file,
			DestExt:  "yaml",
			DestName: fileName,
			DestDir:  "",
			Config:   "",
		})

		return nil
	}

	// Otherwise, expect an array of objects or strings
	var rawIncludes []json.RawMessage
	if err := json.Unmarshal(data, &rawIncludes); err != nil {
		return err
	}

	for _, raw := range rawIncludes {
		if raw[0] == '"' { // Check if it's a string
			var file string
			if err := json.Unmarshal(raw, &file); err != nil {
				return err
			}
			// strip extension from file
			ext := filepath.Ext(file)
			fileName := strings.TrimSuffix(file, ext)
			*k = append(*k, Kr8ComponentSpecIncludeObject{
				File:     file,
				DestExt:  "yaml",
				DestName: fileName,
				DestDir:  "",
				Config:   "",
			})
		} else { // Otherwise, it's an object
			var include Kr8ComponentSpecIncludeObject
			if err := json.Unmarshal(raw, &include); err != nil {
				return err
			}
			*k = append(*k, include)
		}
	}

	return nil
}
