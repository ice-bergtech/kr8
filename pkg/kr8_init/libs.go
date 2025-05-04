package kr8_init

import (
	"os"

	util "github.com/ice-bergtech/kr8/pkg/util"
)

// Downloads a starter kr8+ jsonnet lib from github.
// If fetch is true, downloads the repo in the /lib directory.
// If false, prints the git commands to run.
// Repo: https://github.com/ice-bergtech/kr8-libsonnet .
// return util.FetchRepoUrl("https://github.com/ice-bergtech/kr8-libsonnet", dstDir+"/kr8-lib", !fetch).
func GenerateLib(fetch bool, dstDir string) error {
	if err := util.ErrorIfCheck("error creating lib directory", os.MkdirAll(dstDir, 0750)); err != nil {
		return err
	}

	return util.FetchRepoUrl("https://github.com/kube-libsonnet/kube-libsonnet.git", dstDir+"/kube-lib", !fetch)
}
