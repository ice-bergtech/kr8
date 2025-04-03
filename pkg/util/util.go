package util

import (
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
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

// Fill with string to include and exclude, using kr8's special parsing
type PathFilterOptions struct {
	// Comma-separated list of include filters
	// Filters can include regex filters using the "~" operator. For example, "name~^myregex$"
	// Filters can include equality matches using the "=" operator. For example, "name=myvalue"
	// Filters can include substring matches using the "=" operator. For example, "name=myvalue"
	// If no operator is provided, it is treated as a substring match against the "name" field.
	Includes string
	// Comma-separated list of exclude filters
	// Filters can include regex filters using the "~" operator. For example, "name~^myregex$"
	// Filters can include equality matches using the "=" operator. For example, "name=myvalue"
	// Filters can include substring matches using the "=" operator. For example, "name=myvalue"
	// If no operator is provided, it is treated as a substring match against the "name" field.
	Excludes string
	// Comma separated cluster names
	// Filters keys on exact match
	Clusters string
}

// Checks if a input object matches a filter string.
// The filter string can be an equality match or a regex match.
func CheckObjectMatch(input gjson.Result, filterString string) bool {
	// equality match
	kv := strings.SplitN(filterString, "=", 2)
	if len(kv) == 2 {
		return input.Get(kv[0]).String() == kv[1]
	}
	// regex match
	kv = strings.SplitN(filterString, "~", 2)
	if len(kv) == 2 {
		matched, _ := regexp.MatchString(kv[1], input.Get(kv[0]).String())
		return matched
	}

	// default to substring match of "name" field if no match type specified
	return strings.Contains(input.Get("name").String(), filterString)
}

func FilterItems(input map[string]string, pf PathFilterOptions) []string {
	if pf.Includes == "" && pf.Excludes == "" {
		return []string{}
	}
	var clusterList []string
	for c := range input {
		if pf.Includes != "" || pf.Excludes != "" {
			gjResult := gjson.Parse(input[c])
			// filter on cluster parameters, passed in gjson path notation with either
			// "=" for equality or "~" for regex match
			include := false
			for _, b := range strings.Split(pf.Includes, ",") {
				include = include || CheckObjectMatch(gjResult, b)
			}
			if !include {
				continue
			}
			// filter on cluster parameters, passed in gjson path notation with either
			// "=" for equality or "~" for regex match
			var exclude bool
			exclude = false
			for _, b := range strings.Split(pf.Excludes, ",") {
				exclude = exclude || CheckObjectMatch(gjResult, b)
			}
			if exclude {
				continue
			}
		}
	}
	return clusterList
}

// util.FatalErrorCheck is a helper function that logs an error and exits the program if the error is not nil.
// Saves 3 lines per use and centralizes fatal errors for rewriting
func FatalErrorCheck(err error, message string) {
	if err != nil {
		log.Fatal().Err(err).Msg(message)
	}
}

// Using the allClusterParams variable and command flags to create a list of clusters to generate
// Clusters can be filtered with "=" for equality or "~" for regex match
func CalculateClusterIncludesExcludes(input map[string]string, filters PathFilterOptions) []string {
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
