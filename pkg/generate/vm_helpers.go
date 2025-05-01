package generate

import (
	"os"
	"path/filepath"
	"strings"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/rs/zerolog"

	jnetvm "github.com/ice-bergtech/kr8/pkg/jnetvm"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

// This function sets up component-specific external code in the JVM.
// It makes the component config available to the jvm under the `kr8` extVar.
func SetupJvmForComponent(
	jvm *jsonnet.VM,
	config string,
	kr8Spec kr8_types.Kr8ClusterSpec,
	componentName string,
) {
	// check if we should prune params
	if kr8Spec.PruneParams {
		jvm.ExtCode("kr8", "std.prune("+config+"."+componentName+")")
	} else {
		jvm.ExtCode("kr8", config+"."+componentName)
	}
}

// This function sets up the JVM for a given component.
// It registers native functions, sets up post-processing, and prunes parameters as required.
// It's faster to create this VM for each component, rather than re-use.
// Default postprocessor just copies input to output.
func SetupBaseComponentJvm(
	vmconfig types.VMConfig,
	config string,
	kr8Spec kr8_types.Kr8ClusterSpec,
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

	return jvm, nil
}

// Loads Jsonnet Library paths from component spec.
// jPathResults always includes base lib folder, `filepath.Join(baseDir, "lib")`.
func loadLibPathsIntoVM(
	compSpec kr8_types.Kr8ComponentSpec,
	compPath string,
	baseDir string,
	jvm *jsonnet.VM,
	logger zerolog.Logger,
) {
	logger.Debug().Str("component", compPath).
		Msg("loadLibPathsIntoVM Loading JPaths into VM for component")

	jPathResults := []string{filepath.Join(baseDir, "lib")}
	for _, jPath := range compSpec.JPaths {
		jPathResults = append(jPathResults, filepath.Join(baseDir, compPath, jPath))
	}

	logger.Debug().Str("component", compPath).
		Msgf("loadLibPathsIntoVM JPaths: %v", jPathResults)

	jvm.Importer(&jsonnet.FileImporter{
		JPaths: jPathResults,
	})
}

// Load external files referenced by a component spec into jvm extVars.
func loadExtFilesIntoVM(
	compSpec kr8_types.Kr8ComponentSpec,
	compPath string,
	kr8Opts types.Kr8Opts,
	jvm *jsonnet.VM,
	logger zerolog.Logger,
) error {
	logger.Debug().Str("component path", compPath).Msgf("loadExtFilesIntoVars Loading extFiles")

	for key, val := range compSpec.ExtFiles {
		filePath := filepath.Join(kr8Opts.BaseDir, compPath, val)

		logger.Debug().Str("component path", compPath).
			Msg("loadExtFilesIntoVars ExtFileVar: " + key + "=" + val + "\n Path: " + filePath)

		if kr8Opts.BaseDir != "./" && !strings.HasPrefix(filePath, kr8Opts.BaseDir) {
			if err := util.ErrorIfCheck("Invalid file path: "+filePath, os.ErrNotExist); err != nil {
				return err
			}
		}
		extFile, err := os.ReadFile(filepath.Clean(filePath))
		if err := util.ErrorIfCheck("Error importing ExtFileVar", err); err != nil {
			return err
		}
		jvm.ExtVar(key, string(extFile))
	}

	return nil
}
