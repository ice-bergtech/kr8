// Generates docs for kr8+ code and commands.
//
//go:generate go run ./docs.go
package main

import (
	"errors"
	"go/build"
	"os"
	"strings"

	cmd "github.com/ice-bergtech/kr8/cmd"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	"github.com/invopop/jsonschema"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra/doc"
)

func GenerateCobraDocs() error {
	err := os.Mkdir("cmd", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}
	err = doc.GenMarkdownTree(cmd.RootCmd, "./cmd")
	if err != nil {
		return err
	}

	return nil
}

func GenerateReflector(pkg string, path string) (*jsonschema.Reflector, error) {
	reflector := new(jsonschema.Reflector)
	err := reflector.AddGoComments(pkg, path)

	return reflector, err
}

func WriteSchema(filename string, schema *jsonschema.Schema) error {
	outFile, err := schema.MarshalJSON()
	if err != nil {
		return err
	}

	return os.WriteFile(filename, outFile, 0600)
}

// Generates jsonschema files for kr8+ resources.
//
//nolint:exhaustruct
func GenerateKr8Schemas() error {
	err := os.Mkdir("schemas", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	reflector, err := GenerateReflector(
		"github.com/ice-bergtech/kr8/pkg/kr8_types",
		"../pkg/kr8_types",
	)
	if err != nil {
		return err
	}

	reflector.AllowAdditionalProperties = true
	schema := reflector.Reflect(&kr8_types.Kr8ClusterJsonnet{})

	if err := WriteSchema("./schemas/kr8-plus-cluster-schema.json", schema); err != nil {
		return err
	}

	schema = reflector.Reflect(&kr8_types.Kr8ComponentJsonnet{})

	return WriteSchema("./schemas/kr8-plus-component-schema.json", schema)
}

// Copies various repo files into documentation directory.
// Copies Readme, patching link paths.
// Copies taskfile, for reference.
func CopyRepoFiles() error {
	destinationFile := "./README-repo.md"
	iFile, err := os.ReadFile("../README.md")
	if err != nil {
		return err
	}

	fixed := strings.ReplaceAll(string(iFile), "docs/", "")
	err = os.WriteFile(destinationFile, []byte(fixed), 0600)
	if err != nil {
		return err
	}

	// copy over referenced task file
	destinationFile = "./Taskfile.yml"
	iFile, err = os.ReadFile("../Taskfile.yml")
	if err != nil {
		return err
	}
	err = os.WriteFile(destinationFile, iFile, 0600)
	if err != nil {
		return err
	}

	return nil
}

func GoMarkDoc() error {
	docRenderer, err := gomarkdoc.NewRenderer()
	if err != nil {
		return err
	}

	repo := lang.Repo{
		Remote:        "https://github.com:icebergtech/kr8",
		DefaultBranch: "main",
		PathFromRoot:  "",
	}

	err = os.Mkdir("godoc", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	docFiles := map[string]string{
		"../cmd":                  "kr8-cmd.md",
		"../pkg/jnetvm":           "kr8-jsonnet.md",
		"../pkg/kr8_cache":        "kr8-cache.md",
		"../pkg/kr8_types":        "kr8-types.md",
		"../pkg/types":            "types.md",
		"../pkg/util":             "kr8-util.md",
		"../pkg/kr8_init":         "kr8-init.md",
		"../pkg/generate":         "kr8-generate.md",
		"../pkg/kr8_native_funcs": "kr8-native-functions.md",
	}

	for pkgPath, pkgDoc := range docFiles {
		buildPkg, err := build.ImportDir(pkgPath, build.ImportComment)
		if err != nil {
			return err
		}

		logger := logger.New(logger.DebugLevel)
		pkg, err := lang.NewPackageFromBuild(logger, buildPkg, lang.PackageWithRepositoryOverrides(&repo))
		if err != nil {
			return err
		}
		output, err := docRenderer.Package(pkg)
		if err != nil {
			return err
		}
		err = os.WriteFile("godoc/"+pkgDoc, []byte(output), 0600)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	if err := GenerateCobraDocs(); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := GoMarkDoc(); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := CopyRepoFiles(); err != nil {
		log.Fatal().Err(err).Send()
	}
	if err := GenerateKr8Schemas(); err != nil {
		log.Fatal().Err(err).Send()
	}
}
