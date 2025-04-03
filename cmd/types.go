package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"
)

// init a struct for a single item
type Cluster struct {
	Name string `json:"name"`
	Path string `json:"-"`
}

// init a grouping struct
type Clusters struct {
	Cluster []Cluster
}

type ClusterJsonnet struct {
	// kr8 configuration for how to process the cluster
	ClusterSpec ClusterSpec `json:"_kr8_spec"`
	// Cluster Level configuration that components can reference
	Cluster Cluster `json:"_cluster"`
	// Distictly named components.
	Components map[string]ClusterComponent `json:"_components"`
}

type ClusterComponent struct {
	// The path to a component folder that contains a params.jsonnet file
	Path string `json:"path"`
}

type ClusterSpec struct {
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
	// The name of the current cluster
	Name string `json:"-"`
}

func CreateClusterSpec(clusterName string, spec gjson.Result, baseDir string, genDirOverride string) (ClusterSpec, error) {
	// First determine the value of generate_dir from the command line args or spec.
	clGenerateDir := genDirOverride
	if clGenerateDir == "" {
		clGenerateDir = spec.Get("generate_dir").String()
	}
	if clGenerateDir == "" {
		return ClusterSpec{}, fmt.Errorf("_kr8_spec.generate_dir must be set in parameters or passed as generate-dir flag")
	}
	// if generateDir does not start with /, then it goes in baseDir
	if !strings.HasPrefix(clGenerateDir, "/") {
		clGenerateDir = baseDir + "/" + clGenerateDir
	}
	clusterDir := clGenerateDir + "/" + clusterName
	log.Debug().Str("cluster", clusterName).Msg("output directory: " + clusterDir)
	return ClusterSpec{
		PostProcessor:      spec.Get("postprocessor").String(),
		GenerateDir:        clGenerateDir,
		GenerateShortNames: spec.Get("generate_short_names").Bool(),
		PruneParams:        spec.Get("prune_params").Bool(),
		ClusterDir:         clusterDir,
		Name:               clusterName,
	}, nil
}

type ComponentJsonnet struct {
	// The default namespace to deploy the component to (optional)
	Namespace string `json:"namespace"`
	// A unique name for the component (optional)
	ReleaseName string `json:"release_name"`
	// Component version number (optional)
	Version string `json:"version"`
	// Component-specific configuration for kr8 (required)
	Kr8Spec ComponentSpec `json:"kr8_spec"`
}

// kr8_spec object in cluster config
type ComponentSpec struct {
	// If true, includes the parameters of the current cluster when generating this component
	Kr8_allparams bool `json:"enable_kr8_allparams"`
	// If true, includes the parameters of all other clusters when generating this component
	Kr8_allclusters bool `json:"enable_kr8_allclusters"`
	// If false, all non-generated files present in the output directory will be removed
	DisableOutputDirClean bool `json:"disable_output_clean"`
	// A list of filenames to include as jsonnet vm external vars
	ExtFiles ExtFileVar `json:"extfiles"`
	// Additional jsonnet libs to the jsonnet vm, path component scoped
	JPaths []string `json:"jpaths"`
	// A list of filenames to include and output as files
	Includes []interface{} `json:"includes"`
}

func CreateComponentSpec(spec gjson.Result) (ComponentSpec, error) {
	specM := spec.Map()
	// spec is missing?
	if len(specM) == 0 {
		log.Fatal().Msg("Component has no `kr8_spec` object")
	}

	componentSpec := ComponentSpec{
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
			var include IncludeFileEntryStruct
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

// file to load as a string into the jsonnet vm
// name to reference the variable in jsonnet code through std.extvar()
// value of the variable, loaded from a file or provided directly
type ExtFileVar map[string]string

// A struct describing an included file
type IncludeFileSpec interface {
	string
	IncludeFileEntryStruct
}

type IncludeFileEntryStruct struct {
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
