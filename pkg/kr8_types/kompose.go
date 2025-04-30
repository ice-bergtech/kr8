package kr8_types

import (
	"encoding/json"

	"github.com/ice-bergtech/kr8/pkg/types"

	kompose "github.com/kubernetes/kompose/pkg/app"
	"github.com/kubernetes/kompose/pkg/kobject"
	"github.com/kubernetes/kompose/pkg/loader"
	"github.com/kubernetes/kompose/pkg/transformer"
	"github.com/kubernetes/kompose/pkg/transformer/kubernetes"
	"github.com/kubernetes/kompose/pkg/transformer/openshift"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	sjson "k8s.io/apimachinery/pkg/runtime/serializer/json"
)

// A struct describing a compose file to be processed by kompose to produce kubernetes manifests.
//
// Based on https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
type KomposeConvertOptions struct {
	// Kubernetes: Set the output controller ("deployment"|"daemonSet"|"replicationController")
	Controller string

	// The kubecfg (?) profile to use, can use multiple profiles
	Profiles []string

	// List of compose file filenames.
	// Filenames should be in the format `[docker-]compose.ym[a]l`
	InputFiles []string
	// Specify a file name or directory to save objects to.
	// if path does not exist, a file is created)
	OutFile string
	// Generate a Helm chart for converted objects
	CreateChart bool
	// Add kompose annotations to generated resource
	WithKomposeAnnotation bool
	// Generate resource files into YAML format
	GenerateYaml bool
	// Spaces length to indent generated yaml files
	GenerateYAMLIndent int
	// Generate resource files into JSON format
	GenerateJSON bool
	// Print converted objects to stdout
	GenerateToStdout bool

	// Set the type of build ("local"|"build-config"(OpenShift only)|"none")
	Build string

	// Specify the namespace of the generated resources`)
	Namespace string
	// Specify the number of replicas in the generated resource spec
	Replicas int
	// Convert docker-compose secrets into files instead of symlinked directories
	SecretsAsFiles bool
	// Use Empty Volumes. Do not generate PVCs
	EmptyVols bool
	// Volumes to be generated ("persistentVolumeClaim"|"emptyDir"|"hostPath" | "configMap")
	Volumes string
	// Specify the size of pvc storage requests in the generated resource spec
	PVCRequestSize string
	// Determine whether to generate network policies
	NetworkPolicies bool

	// Create multiple containers grouped by 'kompose.service.group' label
	MultipleContainerMode bool
	// Group multiple service to create single workload by `label`(`kompose.service.group`) or `volume`(shared volumes)
	ServiceGroupMode string
	// Using with --service-group-mode=volume to specific a final service name for the group
	ServiceGroupName string

	// ??
	Provider string
	// ??
	StoreManifest bool
	// ??
	Server string

	// OpenShift: ??
	OSCreateDeploymentConfig bool
	// Openshift: Specify source repository for buildconfig (default remote origin)
	OSBuildRepo string
	// Openshift: Use an insecure Docker repository for OpenShift ImageStream
	OSInsecureRepository bool
	// Openshift: Specify repository branch to use for buildconfig (default master)
	OSBuildBranch string

	// Whether to push built docker image to remote registry.
	ImagePush bool
	// Command used to build to image.  Used with PushCommand
	ImageBuildCommand string
	// Command used to push image
	ImagePushCommand string
	// Specify registry for pushing image, which overrides the registry derived from image name
	ImagePushRegistry string
}

// Initialie Kompose options with sensible defaults.
func Create(inputFiles []string, outDir string, cmp Kr8ComponentJsonnet) *KomposeConvertOptions {
	return &KomposeConvertOptions{
		CreateChart: false,
		Controller:  "deployment",
		Replicas:    1,
		Namespace:   cmp.Namespace,

		ImagePush:    false,
		GenerateJSON: false,
		GenerateYaml: false,

		Volumes:         "persistentVolumeClaim",
		EmptyVols:       false,
		PVCRequestSize:  "100m",
		SecretsAsFiles:  true,
		NetworkPolicies: true,

		Provider:           "kubernetes",
		Build:              "local",
		InputFiles:         inputFiles,
		OutFile:            outDir,
		GenerateYAMLIndent: 2,
		GenerateToStdout:   false,

		Profiles:              []string{},
		WithKomposeAnnotation: true,

		StoreManifest: true,

		MultipleContainerMode: false,
		ServiceGroupMode:      "",
		ServiceGroupName:      "",

		Server:                   "",
		OSCreateDeploymentConfig: false,
		OSBuildRepo:              "",
		OSInsecureRepository:     false,
		OSBuildBranch:            "",
		ImageBuildCommand:        "",
		ImagePushCommand:         "",
		ImagePushRegistry:        "",
	}
}

// Generates a ConvertOptions struct that kompose expects from our commented KomposeConvertOptions
//
// References:
//
// https://pkg.go.dev/github.com/kubernetes/kompose@v1.35.0/pkg/kobject#ConvertOptions
//
// https://github.com/kubernetes/kompose/blob/v1.35.0/pkg/app/app.go#L166
func (k KomposeConvertOptions) GenKomposePkgOpts() *kobject.ConvertOptions {
	isKube := k.Provider == "kubernetes"

	// Initialize te result options with sensible defaults
	resultOpts := kobject.ConvertOptions{
		ToStdout:   k.GenerateToStdout,
		InputFiles: k.InputFiles,
		Profiles:   []string{},
		Namespace:  k.Namespace,
		Provider:   k.Provider,

		PushImage:         k.ImagePush,
		PushImageRegistry: k.ImagePushRegistry,

		CreateChart:    k.CreateChart,
		GenerateYaml:   k.GenerateYaml,
		GenerateJSON:   k.GenerateJSON,
		StoreManifest:  k.StoreManifest,
		EmptyVols:      k.EmptyVols,
		Volumes:        k.Volumes,
		PVCRequestSize: k.PVCRequestSize,

		InsecureRepository: k.OSInsecureRepository,
		Replicas:           k.Replicas,
		OutFile:            k.OutFile,
		Controller:         k.Controller,

		BuildCommand: k.ImageBuildCommand,
		PushCommand:  k.ImagePushCommand,

		Server: k.Server,

		YAMLIndent: k.GenerateYAMLIndent,

		WithKomposeAnnotation: k.WithKomposeAnnotation,

		MultipleContainerMode: k.MultipleContainerMode,
		ServiceGroupMode:      k.ServiceGroupMode,
		ServiceGroupName:      k.ServiceGroupName,

		SecretsAsFiles:          k.SecretsAsFiles,
		GenerateNetworkPolicies: k.NetworkPolicies,
		// https://github.com/kubernetes/kompose/blob/v1.35.0/pkg/app/app.go#L166

		// deployment
		CreateD:          k.Controller == "deployment" && isKube,
		IsDeploymentFlag: k.Controller == "deployment" && isKube,
		// daemon sets
		CreateDS:        k.Controller == "daemonSet" && isKube,
		IsDaemonSetFlag: k.Controller == "daemonSet" && isKube,
		// replication controller
		CreateRC:                    k.Controller == "replicationController" && isKube,
		IsReplicationControllerFlag: k.Controller == "replicationController" && isKube,

		// Openshift specific params
		CreateDeploymentConfig: k.OSCreateDeploymentConfig && !isKube,
		BuildRepo:              k.OSBuildRepo,
		BuildBranch:            k.OSBuildBranch,

		// TODO ??
		Build:                  "",
		IsReplicaSetFlag:       false,
		IsDeploymentConfigFlag: false,
		IsNamespaceFlag:        false,
	}

	return &resultOpts
}

// Validates a set of options for converting a Kubernetes manifest to a Docker Compose file.
func (k KomposeConvertOptions) Validate() error {
	if k.OutFile == "" {
		return types.Kr8Error{Message: "OutFile must be set", Value: ""}
	}
	if len(k.InputFiles) == 0 {
		return types.Kr8Error{Message: "InputFiles must be set", Value: 0}
	}
	// Makes sure the input files are present and are named in a compose-file way
	return kompose.ValidateComposeFile(k.GenKomposePkgOpts())
}

// Converts a Docker Compose file described by k into a set of kubernetes manifests.
func (k KomposeConvertOptions) Convert() (interface{}, error) {
	return convertComposeToK8s(*k.GenKomposePkgOpts())
}

// Convenience method to return the appropriate Transformer based on
// what provider we are using.
func getTransformer(opt kobject.ConvertOptions) *transformer.Transformer {
	var tFormer transformer.Transformer
	if opt.Provider == "kubernetes" {
		// Create/Init new Kubernetes object with CLI opts
		tFormer = &kubernetes.Kubernetes{Opt: opt}
	} else {
		// Create/Init new OpenShift object that is initialized with a newly
		// created Kubernetes object. Openshift inherits from Kubernetes
		tFormer = &openshift.OpenShift{Kubernetes: kubernetes.Kubernetes{Opt: opt}}
	}

	return &tFormer
}

// Convert transforms docker compose or dab file to k8s objects
//
// Based on https://github.com/kubernetes/kompose/blob/main/pkg/app/app.go#L209
func convertComposeToK8s(opt kobject.ConvertOptions) ([]interface{}, error) {
	loader, err := loader.GetLoader("compose")
	if err != nil {
		return nil, err
	}

	// Load the docker-compose file
	objects, err := loader.LoadFile(opt.InputFiles, []string{})
	if err != nil {
		return nil, err
	}

	// Transform the loaded objects into Kubernetes objects
	tFormer := *getTransformer(opt)
	k8sObjects, err := tFormer.Transform(objects, opt)
	if err != nil {
		return nil, err
	}

	// Create a Scheme
	scheme := runtime.NewScheme()
	_ = v1.AddToScheme(scheme)

	// Create a JSON serializer
	jsonSerializer := sjson.NewSerializerWithOptions(
		sjson.SimpleMetaFactory{}, scheme, scheme,
		sjson.SerializerOptions{Yaml: false, Pretty: false, Strict: false},
	)
	// Convert the Kubernetes objects to a format that Jsonnet can use
	result := make([]interface{}, len(k8sObjects))
	for idx, obj := range k8sObjects {
		jsonObj, err := runtime.Encode(jsonSerializer, obj)
		if err != nil {
			return nil, err
		}

		var mapObj map[string]interface{}
		if err := json.Unmarshal(jsonObj, &mapObj); err != nil {
			return nil, err
		}

		result[idx] = mapObj
	}

	return result, nil
}
