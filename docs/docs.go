package main

//go:generate go run ./docs.go

import (
	"errors"
	"go/build"
	"log"
	"os"
	"strings"

	cmd "github.com/ice-bergtech/kr8p/cmd"
	"github.com/princjef/gomarkdoc"
	"github.com/princjef/gomarkdoc/lang"
	"github.com/princjef/gomarkdoc/logger"
	"github.com/spf13/cobra/doc"
)

func CobraDocs() {
	err := os.Mkdir("cmd", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}
	err = doc.GenMarkdownTree(cmd.RootCmd, "./cmd")
	if err != nil {
		log.Fatal(err)
	}
}

func CopyReadme() {
	destinationFile := "./README-repo.md"
	iFile, err := os.ReadFile("../README.md")
	if err != nil {
		log.Fatal(err)
	}

	fixed := strings.ReplaceAll(string(iFile), "docs/", "")
	err = os.WriteFile(destinationFile, []byte(fixed), 0600)
	if err != nil {
		log.Fatal(err)
	}

	// copy over referenced task file
	destinationFile = "./Taskfile.yml"
	iFile, err = os.ReadFile("../Taskfile.yml")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(destinationFile, iFile, 0600)
	if err != nil {
		log.Fatal(err)
	}
}

func GoMarkDoc() {
	out, err := gomarkdoc.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	err = os.Mkdir("godoc", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		log.Fatal(err)
	}

	docfiles := map[string]string{
		"../cmd":          "kr8p-cmd.md",
		"../pkg/jnetvm":   "kr8p-jsonnet.md",
		"../pkg/types":    "kr8p-types.md",
		"../pkg/util":     "kr8p-util.md",
		"../pkg/kr8_init": "kr8p-init.md",
		"../pkg/generate": "kr8p-generate.md",
	}

	for pkgPath, pkgDoc := range docfiles {
		buildPkg, err := build.ImportDir(pkgPath, build.ImportComment)
		if err != nil {
			log.Fatal(err)
		}

		logger := logger.New(logger.DebugLevel)
		pkg, err := lang.NewPackageFromBuild(logger, buildPkg)
		if err != nil {
			log.Fatal(err)
		}
		output, err := out.Package(pkg)
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile("godoc/"+pkgDoc, []byte(output), 0600)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	CobraDocs()
	GoMarkDoc()
	CopyReadme()
}
