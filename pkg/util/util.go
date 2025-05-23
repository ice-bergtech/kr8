// Package util contains various utility functions for directories and files.
// It includes functions for
// formatting JSON,
// writing to files,
// directory management,
// and go control-flow helpers
package util

import (
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// Filter returns a new slice containing only the elements that satisfy the predicate function.
// From https://gobyexample.com/collection-functions
func Filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}

	return vsf
}

// Fill with string to include and exclude, using kr8's special parsing.
type PathFilterOptions struct {
	// Comma-separated list of include filters
	// Filters can include:
	//
	// regex filters using the "~" operator. For example, "name~^myRegex$"
	// equality matches using the "=" operator. For example, "name=myValue"
	// substring matches using the "=" operator. For example, "name=myValue"
	//
	// If no operator is provided, it is treated as a substring match against the "name" field.
	Includes string
	// Comma-separated list of exclude filters.
	// Filters can include:
	//
	// regex filters using the "~" operator. For example, "name~^myRegex$"
	// equality matches using the "=" operator. For example, "name=myValue"
	// substring matches using the "=" operator. For example, "name=myValue"
	//
	// If no operator is provided, it is treated as a substring match against the "name" field.
	Excludes string
	// Comma separated cluster names.
	// Filters keys on exact match.
	Clusters string
	// Comma separated component names.
	Components string
}

// Checks if a input object matches a filter string.
// The filter string can be an equality match or a regex match.
func CheckObjectMatch(input gjson.Result, filterString string) bool {
	// equality match
	args := strings.SplitN(filterString, "=", 2)
	if len(args) == 2 {
		return input.Get(args[0]).String() == args[1]
	}
	// regex match
	args = strings.SplitN(filterString, "~", 2)
	if len(args) == 2 {
		matched, _ := regexp.MatchString(args[1], input.Get(args[0]).String())
		// Found a match, return
		return matched
	}

	// default to substring match of "name" field if no match type specified
	return strings.Contains(input.Get("name").String(), filterString)
}

// Given a map of string, filter them based on the provided options.
// The map value is parsed as a gjson result and then checked against the provided options.
func FilterItems(input map[string]string, pFilter PathFilterOptions) []string {
	if pFilter.Includes == "" && pFilter.Excludes == "" {
		// Exit hatch
		return []string{}
	}
	var clusterList []string
	for c := range input {
		gjResult := gjson.Parse(input[c])
		// filter on cluster parameters, passed in gjson path notation with either
		// "=" for equality or "~" for regex match
		include := false
		for _, b := range strings.Split(pFilter.Includes, ",") {
			include = include || CheckObjectMatch(gjResult, b)
		}
		if !include {
			continue
		}
		// filter on cluster parameters, passed in gjson path notation with either
		// "=" for equality or "~" for regex match
		var exclude bool
		exclude = false
		for _, b := range strings.Split(pFilter.Excludes, ",") {
			exclude = exclude || CheckObjectMatch(gjResult, b)
		}
		if exclude {
			continue
		}
	}

	return clusterList
}

// Using the allClusterParams variable and command flags to create a list of clusters to generate.
// Clusters can be filtered with "=" for equality or "~" for regex match.
func CalculateClusterIncludesExcludes(input map[string]string, filters PathFilterOptions) []string {
	// Defer to using clusters if set
	if filters.Clusters != "" {
		var clusterList []string
		// all clusters
		for _, key := range strings.Split(filters.Clusters, ",") {
			val, ok := input[key]
			if ok {
				clusterList = append(clusterList, val)
			}
		}

		return clusterList
	}

	return FilterItems(input, filters)
}

// Calculate the sha256 hash and returns the base64 encoded result.
func HashFile(path string) (string, error) {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return "", err
	}
	defer file.Close()

	hashBox := sha256.New()
	if _, err := io.Copy(hashBox, file); err != nil {
		return "", err
	}

	return base64.RawStdEncoding.EncodeToString(hashBox.Sum(nil)), nil
}
