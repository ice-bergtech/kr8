package main

//go:generate go run ./docs.go

import (
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
	err := doc.GenMarkdownTree(cmd.RootCmd, "./cmd")
	if err != nil {
		log.Fatal(err)
	}
}

func GoMarkDoc() {
	out, err := gomarkdoc.NewRenderer()
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}

	os.Mkdir("godoc", 0755)

	docfiles := map[string]string{
		"../cmd":          "kr8-cmd.md",
		"../pkg/jvm":      "kr8-jsonnet.md",
		"../pkg/types":    "kr8-types.md",
		"../pkg/util":     "kr8-util.md",
		"../pkg/kr8_init": "kr8-init.md",
	}

	for k, v := range docfiles {
		buildPkg, err := build.ImportDir(k, build.ImportComment)
		if err != nil {
			log.Fatal(err)
		}

		log := logger.New(logger.DebugLevel)
		pkg, err := lang.NewPackageFromBuild(log, buildPkg)
		output, err := out.Package(pkg)
		os.WriteFile("godoc/"+v, []byte(output), 0644)
	}
}

func main() {
	CobraDocs()
	GoMarkDoc()
}
