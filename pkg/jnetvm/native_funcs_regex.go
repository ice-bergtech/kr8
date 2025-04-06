/*

Much of this code is based on the kubecfg project: https://github.com/ksonnet/kubecfg
Native funcs: https://github.com/kubecfg/kubecfg/blob/main/utils/nativefuncs.go

Copyright 2018 ksonnet

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package jnetvm

import (
	"regexp"

	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
)

// Escapes a string for use in regex.
//
// Inputs: "str".
func NativeRegexEscape() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "regexEscapeString",
		Params: []jsonnetAst.Identifier{"str"},
		Func: func(args []interface{}) (interface{}, error) {
			return regexp.QuoteMeta(args[0].(string)), nil
		}}
}

// Matches a string against a regex pattern.
//
// Inputs: "regex", "string".
func NativeRegexMatch() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "regexMatch",
		Params: []jsonnetAst.Identifier{"regex", "string"},
		Func: func(args []interface{}) (interface{}, error) {
			return regexp.MatchString(args[0].(string), args[1].(string))
		}}
}

// Substitutes a regex pattern in a string with another string.
//
// Inputs: "regex", "src", "repl".
func NativeRegexSubst() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "regexSubst",
		Params: []jsonnetAst.Identifier{"regex", "src", "repl"},
		Func: func(args []interface{}) (interface{}, error) {
			regex := args[0].(string)
			src := args[1].(string)
			repl := args[2].(string)

			r, err := regexp.Compile(regex)
			if err != nil {
				return "", err
			}

			return r.ReplaceAllString(src, repl), nil
		}}
}
