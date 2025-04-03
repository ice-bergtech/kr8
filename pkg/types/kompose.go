package types

import (
	"fmt"

	kompose "github.com/kubernetes/kompose/pkg/app"
	"github.com/kubernetes/kompose/pkg/kobject"
)

// A struct describing a compose file that will be processed by kompose to produce kubernetes manifests
// Based on https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
type KomposeConvertOptions struct {
	// Kubernetes: Create a Helm chart for converted objects
	CreateChart bool
	// Kubernetes: Set the output controller ("deployment"|"daemonSet"|"replicationController")
	Controller string

	// Print converted objects to stdout
	ToStdout bool

	// Set the type of build ("local"|"build-config"(OpenShift only)|"none")
	Build string
	// Openshift: Specify source repository for buildconfig (default remote origin)
	BuildRepo string
	// Openshift: Specify repository branch to use for buildconfig (default master)
	BuildBranch string

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
	Replicas int
	// List of compose file filenames
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

func Create(inputFiles []string, outDir string, cmp Kr8ComponentJsonnet) *KomposeConvertOptions {
	return &KomposeConvertOptions{
		CreateChart: false,
		Controller:  "deployment",
		CreateD:     true,
		Replicas:    1,
		Namespace:   cmp.Namespace,

		PushImage:    false,
		GenerateJSON: true,
		GenerateYaml: false,

		EmptyVols:               false,
		PVCRequestSize:          "100m",
		SecretsAsFiles:          true,
		GenerateNetworkPolicies: true,

		Provider:   "kubernetes",
		Build:      "local",
		InputFiles: inputFiles,
		OutFile:    outDir,
		YAMLIndent: 2,
		ToStdout:   false,
	}
}

func (k KomposeConvertOptions) genKomposePkgOpts() *kobject.ConvertOptions {
	var resultOpts kobject.ConvertOptions

	resultOpts.InputFiles = k.InputFiles
	resultOpts.Profiles = ""
	resultOpts.Namespace = k.Namespace
	resultOpts.Provider = k.Provider

	// https://github.com/kubernetes/kompose/blob/v1.35.0/pkg/app/app.go#L166
	if k.Provider == "kubernetes" {
		if k.Controller == "deployment" {
			resultOpts.CreateD = true
		} else if k.Controller == "daemonSet" {
			resultOpts.CreateDS = true
		} else if k.Controller == "replicationController" {
			resultOpts.CreateRC = true
		} else {
			resultOpts.CreateD = true
		}
	} else if k.Provider == "openshift" {
	}

	//

	return &resultOpts
}

func (k KomposeConvertOptions) Validate() error {
	if k.OutFile == "" {
		return fmt.Errorf("OutFile must be set")
	}
	if len(k.InputFiles) == 0 {
		return fmt.Errorf("InputFiles must be set")
	}
	// Makes sure the input files are present and are named in a compose-file way
	return kompose.ValidateComposeFile(k.genKomposePkgOpts())
}

func (k KomposeConvertOptions) Convert() (interface{}, error) {
	return kompose.Convert(*k.genKomposePkgOpts())
}
