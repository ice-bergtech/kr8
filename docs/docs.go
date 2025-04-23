//go:build ignore

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
	docRenderrer, err := gomarkdoc.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}

	repo := lang.Repo{
		Remote:        "https://github.com:icebergtech/kr8",
		DefaultBranch: "main",
		PathFromRoot:  "",
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
		pkg, err := lang.NewPackageFromBuild(logger, buildPkg, lang.PackageWithRepositoryOverrides(&repo))
		if err != nil {
			log.Fatal(err)
		}
		output, err := docRenderrer.Package(pkg)
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
