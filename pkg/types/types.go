package types

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

// An object that stores variables that can be referenced by components
type Kr8Cluster struct {
	Name string `json:"name"`
	Path string `json:"-"`
}

// The specification for a clusters.jsonnet file
// This file contains configuration for clusters, including
type Kr8ClusterJsonnet struct {
	// kr8 configuration for how to process the cluster
	ClusterSpec Kr8ClusterSpec `json:"_kr8_spec"`
	// Cluster Level configuration that components can reference
	Cluster Kr8Cluster `json:"_cluster"`
	// Distictly named components.
	Components map[string]Kr8ClusterComponentRef `json:"_components"`
}

// A reference to a component folder that contains a params.jsonnet file. This is used in the cluster jsonnet file to reference components.
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
	// The root directory for the cluster. Default `clusters`
	ClusterDir string `json:"-"`
}

// This function creates a cluster spec from the given cluster name, spec, base directory, and generate directory override.
// If genDirOverride is empty, the value of generate_dir from the spec is used.
func CreateClusterSpec(clusterName string, spec gjson.Result, baseDir string, genDirOverride string) (Kr8ClusterSpec, error) {
	// First determine the value of generate_dir from the command line args or spec.
	clGenerateDir := genDirOverride
	if clGenerateDir == "" {
		clGenerateDir = spec.Get("generate_dir").String()
	}
	if clGenerateDir == "" {
		return Kr8ClusterSpec{}, fmt.Errorf("_kr8_spec.generate_dir must be set in parameters or passed as generate-dir flag")
	}
	// if generateDir does not start with /, then it goes in baseDir
	if !strings.HasPrefix(clGenerateDir, "/") {
		clGenerateDir = baseDir + "/" + clGenerateDir
	}
	clusterDir := clGenerateDir + "/" + clusterName
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

// The specification for component's params.jsonnet file
// It contains all the configuration and variables used to generate component resources
// This configuration is often modified from the cluster config to add cluster-specific configuration
type Kr8ComponentJsonnet struct {
	// Component-specific configuration for how kr8 should process the component (required)
	Kr8Spec Kr8ComponentSpec `json:"kr8_spec"`
	// The default namespace to deploy the component to
	Namespace string `json:"namespace"`
	// A unique name for the component
	ReleaseName string `json:"release_name"`
	// Component version string (optional)
	Version string `json:"version"`
}

// The kr8_spec object in a cluster config file
// This configures how kr8 processes the component
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

// Extracts a component spec from a jsonnet object.
func CreateComponentSpec(spec gjson.Result) (Kr8ComponentSpec, error) {
	specM := spec.Map()
	// spec is missing?
	if len(specM) == 0 {
		log.Fatal().Msg("Component has no `kr8_spec` object")
	}

	componentSpec := Kr8ComponentSpec{
		Kr8_allparams:         spec.Get("enable_kr8_allparams").Bool(),
		Kr8_allclusters:       spec.Get("enable_kr8_allclusters").Bool(),
		DisableOutputDirClean: spec.Get("disable_output_clean").Bool(),
		ExtFiles:              ExtFileVar{},
		JPaths:                []string{},
		Includes:              []interface{}{},
	}

	for k, v := range spec.Get("extfiles").Map() {
		if v.Type == gjson.String {
			componentSpec.ExtFiles[k] = v.String()
		}
	}

	jPaths := spec.Get("jpaths").Array()
	componentSpec.JPaths = make([]string, len(jPaths))
	for i, p := range jPaths {
		componentSpec.JPaths[i] = p.String()
	}

	incl := spec.Get("includes")
	componentSpec.Includes = make([]interface{}, len(incl.Array()))
	for i, item := range incl.Array() {
		if item.Type == gjson.JSON {
			var include Kr8ComponentSpecIncludeObject
			err := json.Unmarshal([]byte(item.String()), &include)
			if err != nil {
				return componentSpec, fmt.Errorf("error unmarshalling includes: %w", err)
			}
			componentSpec.Includes[i] = include
		} else if item.Type == gjson.String {
			componentSpec.Includes[i] = item.String()
		}
	}

	return componentSpec, nil
}

// Map of external files to load into jsonnet vm as external variables
// Keys are the variable names, values are the paths to the files to load as strings into the jsonnet vm
// To reference the variable in jsonnet code, use std.extvar("variable_name")
type ExtFileVar map[string]string

// A struct describing an included file that will be processed to produce a file
type Kr8ComponentSpecIncludeFile interface {
	string
	Kr8ComponentSpecIncludeObject
}

// An includes object which configures how kr8 includes an object
// It allows configuring the included file's destination directory and file name
// The input file will be processed differently depending on the filetype
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

type CmdJsonnetOptions struct {
	Prune         bool
	Cluster       string
	ClusterParams string
	Component     string
	Format        string
	Color         bool
}

// VMConfig describes configuration to initialize Jsonnet VM with
type VMConfig struct {
	// Jpaths is a list of paths to search for Jsonnet libraries (libsonnet files)
	Jpaths []string `json:"jpath" yaml:"jpath"`
	// ExtVars is a list of external variables to pass to Jsonnet VM
	ExtVars []string `json:"ext_str_file" yaml:"ext_str_files"`
	// base directory for the project
	BaseDir string `json:"base_dir" yaml:"base_dir"`
}

// A struct describing a compose file that will be processed by kompose to produce kubernetes manifests
// Based on https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
type KomposeConvertOptions struct {
	// Kubernetes: Create a Helm chart for converted objects
	CreateChart bool
	// Kubernetes: Set the output controller ("deployment"|"daemonSet"|"replicationController")
	Controller       string
	IsDaemonSetFlag  bool
	IsDeploymentFlag bool

	// Print converted objects to stdout
	ToStdout bool

	// Openshift: Specify source repository for buildconfig (default remote origin)
	BuildRepo string
	// Openshift: Specify repository branch to use for buildconfig (default master)
	BuildBranch string
	// Set the type of build ("local"|"build-config"(OpenShift only)|"none")
	Build string

	// convertCmd.Flags().BoolVar(&ConvertDeploymentConfig, "deployment-config", true, "Generate an OpenShift deploymentconfig object")

	Profiles  []string
	PushImage bool
	// Specify registry for pushing image, which will override registry from image name
	PushImageRegistry string
	// Generate resource files into YAML format
	GenerateYaml bool
	// Generate resource files into JSON format
	GenerateJSON  bool
	StoreManifest bool
	// Use Empty Volumes. Do not generate PVCs
	EmptyVols bool
	// Volumes to be generated ("persistentVolumeClaim"|"emptyDir"|"hostPath" | "configMap")
	Volumes string
	// Specify the size of pvc storage requests in the generated resource spec
	PVCRequestSize string
	// Use an insecure Docker repository for OpenShift ImageStream
	InsecureRepository bool
	// Specify the number of replicas in the generated resource spec
	Replicas   int
	InputFiles []string
	// Specify a file name or directory to save objects to (if path does not exist, a file will be created)
	OutFile  string
	Provider string
	// Specify the namespace of the generated resources`)
	Namespace string

	// convertCmd.Flags().BoolVar(&ConvertPushImage, "push-image", false, "If we should push the docker image we built")
	// convertCmd.Flags().StringVar(&BuildCommand, "build-command", "", `Set the command used to build the container image, which will override the docker build command. Should be used in conjuction with --push-command flag.`)
	// convertCmd.Flags().StringVar(&PushCommand, "push-command", "", `Set the command used to push the container image. override the docker push command. Should be used in conjuction with --build-command flag.`)

	IsReplicationControllerFlag bool
	IsReplicaSetFlag            bool
	IsDeploymentConfigFlag      bool
	IsNamespaceFlag             bool

	BuildCommand string
	PushCommand  string

	Server string

	// Spaces length to indent generated yaml files
	YAMLIndent int

	// Add kompose annotations to generated resource
	WithKomposeAnnotation bool

	// Create multiple containers grouped by 'kompose.service.group' label
	MultipleContainerMode bool
	// Group multiple service to create single workload by `label`(`kompose.service.group`) or `volume`(shared volumes)
	ServiceGroupMode string
	// Using with --service-group-mode=volume to specific a final service name for the group
	ServiceGroupName string

	// Convert docker-compose secrets into files instead of symlinked directories
	SecretsAsFiles bool
	// Specify whether to generate network policies or not
	GenerateNetworkPolicies bool
}
