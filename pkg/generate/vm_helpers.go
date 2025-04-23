package generate

import (
	"os"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/rs/zerolog/log"

	jnetvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// This function sets up the JVM for a given component.
// It registers native functions, sets up post-processing, and prunes parameters as required.
// It's faster to create this VM for each component, rather than re-use.
// Default postprocessor just copies input to output.
func SetupJvmForComponent(
	vmconfig types.VMConfig,
	config string,
	kr8Spec kr8_types.Kr8ClusterSpec,
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
func loadJPathsIntoVM(compSpec kr8_types.Kr8ComponentSpec, compPath string, baseDir string, jvm *jsonnet.VM) {
	log.Debug().
		Str("component path", compPath).
		Msg("Loading JPaths into VM for component")
	jPathResults := []string{filepath.Join(baseDir, "lib")}
	for _, jPath := range compSpec.JPaths {
		jPathResults = append(jPathResults, filepath.Join(baseDir, compPath, jPath))
	}
	log.Debug().Str("component path", compPath).Msgf("JPaths: %v", jPathResults)
	jvm.Importer(&jsonnet.FileImporter{
		JPaths: jPathResults,
	})
}

func loadExtFilesIntoVars(
	compSpec kr8_types.Kr8ComponentSpec,
	compPath string,
	kr8Spec kr8_types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	componentName string,
	jvm *jsonnet.VM,
) error {
	log.Debug().Str("component path", compPath).Msgf("Loading extFiles")
	for key, val := range compSpec.ExtFiles {
		filePath := filepath.Join(kr8Opts.BaseDir, compPath, val)
		log.Debug().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Msg("Extfile: " + key + "=" + val + "\n Path: " + filePath)
		if kr8Opts.BaseDir != "./" && !strings.HasPrefix(filePath, kr8Opts.BaseDir) {
			if err := util.GenErrorIfCheck("Invalid file path: "+filePath, os.ErrNotExist); err != nil {
				return err
			}
		}
		extFile, err := os.ReadFile(filepath.Clean(filePath))
		if err := util.GenErrorIfCheck("Error importing extfiles item", err); err != nil {
			return err
		}
		jvm.ExtVar(key, string(extFile))
	}

	return nil
}
