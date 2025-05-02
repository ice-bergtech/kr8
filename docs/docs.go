// Generates docs for kr8+ code and commands.
//
//go:generate go run ./docs.go
package main

import (
	"errors"
	"go/build"
	"log"
	"os"
	"strings"

	cmd "github.com/ice-bergtech/kr8/cmd"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
	"github.com/spf13/cobra/doc"
)

func CobraDocs() error {
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

func CopyReadme() error {
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
	if err := CobraDocs(); err != nil {
		log.Fatal(err)
	}
	if err := GoMarkDoc(); err != nil {
		log.Fatal(err)
	}
	if err := CopyReadme(); err != nil {
		log.Fatal(err)
	}
}
