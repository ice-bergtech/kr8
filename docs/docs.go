package main

//go:generate go run ./docs.go

import (
	"errors"
	"go/build"
	"log"
	"os"

	cmd "github.com/ice-bergtech/kr8/cmd"
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
	err = os.WriteFile(destinationFile, iFile, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func CopyLicenses() {
	destinationFile := "./README-repo.md"
	iFile, err := os.ReadFile("../README.md")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(destinationFile, iFile, 0644)
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
		"../cmd":          "kr8-cmd.md",
		"../pkg/jnetvm":   "kr8-jsonnet.md",
		"../pkg/types":    "kr8-types.md",
		"../pkg/util":     "kr8-util.md",
		"../pkg/kr8_init": "kr8-init.md",
		"../pkg/generate": "kr8-generate.md",
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
