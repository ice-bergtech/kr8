package main

//go:generate go run ./docs.go

import (
	"log"

	cmd "github.com/ice-bergtech/kr8/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.RootCmd, "./cmd")
	if err != nil {
		log.Fatal(err)
	}
}
