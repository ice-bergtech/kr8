package generate

import (
	"os"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	jnetvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/rs/zerolog/log"
)

// This function sets up the JVM for a given component.
// It registers native functions, sets up post-processing, and prunes parameters as required.
// It's faster to create this VM for each component, rather than re-use.
// Default postprocessor just copies input to output.
func SetupJvmForComponent(
	vmconfig types.VMConfig,
	config string,
	kr8Spec types.Kr8ClusterSpec,
	componentName string,
) (*jsonnet.VM, error) {
	jvm, err := jnetvm.JsonnetVM(vmconfig)
	if err != nil {
		return nil, err
	}
	jnetvm.RegisterNativeFuncs(jvm)
	jvm.ExtCode("kr8_cluster", "std.prune("+config+"._cluster)")

	if kr8Spec.PostProcessor != "" {
		jvm.ExtCode("process", kr8Spec.PostProcessor)
	} else {
		// Default PostProcessor passes input to output
		jvm.ExtCode("process", "function(input) input")
	}

	// check if we should prune params
	if kr8Spec.PruneParams {
		jvm.ExtCode("kr8", "std.prune("+config+"."+componentName+")")
	} else {
		jvm.ExtCode("kr8", config+"."+componentName)
	}

	return jvm, nil
}

// jPathResults always includes base lib.
// Adds jpaths from spec if set.
func loadJPathsIntoVM(compSpec types.Kr8ComponentSpec, compPath string, baseDir string, jvm *jsonnet.VM) {
	jPathResults := []string{filepath.Join(baseDir, "lib")}
	for _, jPath := range compSpec.JPaths {
		jPathResults = append(jPathResults, filepath.Join(baseDir, compPath, jPath))
	}
	jvm.Importer(&jsonnet.FileImporter{
		JPaths: jPathResults,
	})
}

func loadExtFilesIntoVars(
	compSpec types.Kr8ComponentSpec,
	compPath string,
	kr8Spec types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	componentName string,
	jvm *jsonnet.VM,
) {
	for key, val := range compSpec.ExtFiles {
		log.Debug().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Msg("Extfile: " + key + "=" + val)
		filePath := filepath.Join(kr8Opts.BaseDir, compPath, val)
		if kr8Opts.BaseDir != "./" && !strings.HasPrefix(filePath, kr8Opts.BaseDir) {
			util.FatalErrorCheck("Invalid file path: "+filePath, os.ErrNotExist)
		}
		extFile, err := os.ReadFile(filepath.Clean(filePath))
		util.FatalErrorCheck("Error importing extfiles item", err)
		jvm.ExtVar(key, string(extFile))
	}
}
