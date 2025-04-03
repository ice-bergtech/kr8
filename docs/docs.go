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
	buildPkg, err := build.ImportDir("../cmd", build.ImportComment)
	if err != nil {
		log.Fatal(err)
	}

	log := logger.New(logger.DebugLevel)
	pkg, err := lang.NewPackageFromBuild(log, buildPkg)
	output, err := out.Package(pkg)
	os.WriteFile("kr8-cmd.md", []byte(output), 0644)
}

func main() {
	CobraDocs()
	GoMarkDoc()
}
