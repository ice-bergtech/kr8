package util

import (
	"fmt"
	"os"
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
	Includes string
	Excludes string
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

func FilterItems(input map[string]string, pf PathFilterOptions) ([]string, error) {
	if pf.Includes == "" && pf.Excludes == "" {
		return []string{}, fmt.Errorf("no filter conditions provided")
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
	return clusterList, nil
}

// util.FatalErrorCheck is a helper function that logs an error and exits the program if the error is not nil.
// Saves 3 lines per use and centralizes fatal errors for rewriting
func FatalErrorCheck(err error, message string) {
	if err != nil {
		log.Fatal().Err(err).Msg(message)
	}
}

// colorize function from zerolog console.go file to replicate their coloring functionality.
// https://github.com/rs/zerolog/blob/a21d6107dcda23e36bc5cfd00ce8fdbe8f3ddc23/console.go#L389
func Colorize(s interface{}, c int, disabled bool) string {
	e := os.Getenv("NO_COLOR")
	if e != "" || c == 0 {
		disabled = true
	}

	if disabled {
		return fmt.Sprintf("%s", s)
	}
	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}
