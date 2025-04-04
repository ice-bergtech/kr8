package types

import (
	"encoding/json"
	"fmt"

	kompose "github.com/kubernetes/kompose/pkg/app"
	"github.com/kubernetes/kompose/pkg/kobject"

	"path/filepath"

	log "github.com/sirupsen/logrus"

	"github.com/kubernetes/kompose/pkg/loader"
	"github.com/kubernetes/kompose/pkg/transformer"
	"github.com/kubernetes/kompose/pkg/transformer/kubernetes"
	"github.com/kubernetes/kompose/pkg/transformer/openshift"
)

// A struct describing a compose file that will be processed by kompose to produce kubernetes manifests.
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
	// Specify a file name or directory to save objects to (if path does not exist, a file will be created)
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
	// Specify whether to generate network policies or not
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
	// Specify registry for pushing image, which will override registry from image name
	ImagePushRegistry string
}

// Initialie Kompose options with sensible defaults
func Create(inputFiles []string, outDir string, cmp Kr8ComponentJsonnet) *KomposeConvertOptions {
	return &KomposeConvertOptions{
		CreateChart: false,
		Controller:  "deployment",
		Replicas:    1,
		Namespace:   cmp.Namespace,

		ImagePush:    false,
		GenerateJSON: true,
		GenerateYaml: false,

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
	var resultOpts kobject.ConvertOptions

	resultOpts.ToStdout = k.GenerateToStdout

	resultOpts.InputFiles = k.InputFiles
	resultOpts.Profiles = []string{}
	resultOpts.Namespace = k.Namespace
	resultOpts.Provider = k.Provider

	// https://github.com/kubernetes/kompose/blob/v1.35.0/pkg/app/app.go#L166
	if k.Provider == "kubernetes" {
		if k.Controller == "deployment" {
			resultOpts.CreateD = true
			resultOpts.IsDeploymentFlag = true
		} else if k.Controller == "daemonSet" {
			resultOpts.CreateDS = true
			resultOpts.IsDaemonSetFlag = true
		} else if k.Controller == "replicationController" {
			resultOpts.CreateRC = true
			resultOpts.IsReplicationControllerFlag = true
		} else {
			resultOpts.CreateD = true
		}
	} else if k.Provider == "openshift" {
		resultOpts.CreateDeploymentConfig = k.OSCreateDeploymentConfig
		resultOpts.BuildRepo = k.OSBuildRepo
		resultOpts.BuildBranch = k.OSBuildBranch
	}

	resultOpts.PushImage = k.ImagePush
	resultOpts.PushImageRegistry = k.ImagePushRegistry

	resultOpts.CreateChart = k.CreateChart
	resultOpts.GenerateYaml = k.GenerateYaml
	resultOpts.GenerateJSON = k.GenerateJSON
	resultOpts.StoreManifest = k.StoreManifest
	resultOpts.EmptyVols = k.EmptyVols
	resultOpts.Volumes = k.Volumes
	resultOpts.PVCRequestSize = k.PVCRequestSize

	resultOpts.InsecureRepository = k.OSInsecureRepository
	resultOpts.Replicas = k.Replicas
	resultOpts.InputFiles = k.InputFiles
	resultOpts.OutFile = k.OutFile
	resultOpts.Provider = k.Provider
	resultOpts.Namespace = k.Namespace
	resultOpts.Controller = k.Controller

	resultOpts.BuildCommand = k.ImageBuildCommand
	resultOpts.PushCommand = k.ImagePushCommand

	resultOpts.Server = k.Server

	resultOpts.YAMLIndent = k.GenerateYAMLIndent

	resultOpts.WithKomposeAnnotation = k.WithKomposeAnnotation

	resultOpts.MultipleContainerMode = k.MultipleContainerMode
	resultOpts.ServiceGroupMode = k.ServiceGroupMode
	resultOpts.ServiceGroupName = k.ServiceGroupName

	resultOpts.SecretsAsFiles = k.SecretsAsFiles
	resultOpts.GenerateNetworkPolicies = k.NetworkPolicies

	return &resultOpts
}

// Validates a set of options for converting a Kubernetes manifest to a Docker Compose file.
func (k KomposeConvertOptions) Validate() error {
	if k.OutFile == "" {
		return fmt.Errorf("OutFile must be set")
	}
	if len(k.InputFiles) == 0 {
		return fmt.Errorf("InputFiles must be set")
	}
	// Makes sure the input files are present and are named in a compose-file way
	return kompose.ValidateComposeFile(k.GenKomposePkgOpts())
}

// Converts a Docker Compose file described by k into a set of kubernetes manifests.
func (k KomposeConvertOptions) Convert() (interface{}, error) {
	return Convert(*k.GenKomposePkgOpts())
}

// Convert transforms docker compose or dab file to k8s objects
func Convert(opt kobject.ConvertOptions) (interface{}, error) {
	// loader parses input from file into komposeObject.
	l, err := loader.GetLoader("compose")
	if err != nil {
		log.Fatal(err)
	}

	komposeObject := kobject.KomposeObject{
		ServiceConfigs: make(map[string]kobject.ServiceConfig),
	}
	komposeObject, err = l.LoadFile(opt.InputFiles, opt.Profiles)
	if err != nil {
		log.Fatalf(err.Error())
	}

	komposeObject.Namespace = opt.Namespace

	// Get the directory of the compose file
	workDir, err := transformer.GetComposeFileDir(opt.InputFiles)
	if err != nil {
		log.Fatalf("Unable to get compose file directory: %s", err)
	}

	// convert env_file from absolute to relative path
	for _, service := range komposeObject.ServiceConfigs {
		if len(service.EnvFile) <= 0 {
			continue
		}
		for i, envFile := range service.EnvFile {
			if !filepath.IsAbs(envFile) {
				continue
			}

			relPath, err := filepath.Rel(workDir, envFile)
			if err != nil {
				log.Fatalf(err.Error())
			}

			service.EnvFile[i] = filepath.ToSlash(relPath)
		}
	}

	// Get a transformer that maps komposeObject to provider's primitives
	t := getTransformer(opt)

	// Do the transformation
	objects, err := t.Transform(komposeObject, opt)

	if err != nil {
		log.Fatalf(err.Error())
	}

	list := []interface{}{}
	for _, obj := range objects {
		jsonBytes, err := json.Marshal(obj)
		list = append(list, string(jsonBytes))
		if err != nil {
			log.Fatalf("Failed to marshal object to JSON: %v", err)
		}
		//fmt.Println(string(jsonBytes))
	}

	result := make(map[string]interface{})
	result["objects"] = list
	return result, nil

	// // Print output
	// err = kubernetes.PrintList(objects, opt)
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// return objects, err
}

// Convenience method to return the appropriate Transformer based on
// what provider we are using.
func getTransformer(opt kobject.ConvertOptions) transformer.Transformer {
	var t transformer.Transformer
	if opt.Provider == "kubernetes" {
		// Create/Init new Kubernetes object with CLI opts
		t = &kubernetes.Kubernetes{Opt: opt}
	} else {
		// Create/Init new OpenShift object that is initialized with a newly
		// created Kubernetes object. Openshift inherits from Kubernetes
		t = &openshift.OpenShift{Kubernetes: kubernetes.Kubernetes{Opt: opt}}
	}
	return t
}
