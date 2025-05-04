// The package kr8 is an opinionated Kubernetes cluster configuration management tool.
// It is designed to simplify and standardize the process of managing Kubernetes clusters.
package main

import "github.com/ice-bergtech/kr8/cmd"

var version = "v0.2.0"

func main() {
	cmd.Execute(version)
}
