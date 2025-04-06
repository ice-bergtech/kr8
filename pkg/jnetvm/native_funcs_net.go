package jnetvm

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/c-robinson/iplib/v2"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"lukechampine.com/uint128"
)

// Contains the url information.
type NativeFuncURL struct {
	Scheme string
	// encoded opaque data
	Opaque string
	// username information
	Username string
	// Whether the password field is set
	PasswordSet bool
	// password information
	Password string
	// host or host:port (see Hostname and Port methods)
	Host string
	// path (relative paths may omit leading slash)
	Path string
	// encoded path hint (see EscapedPath method)
	RawPath string
	// query values
	Query map[string]interface{}
	// encoded query values, without '?'
	RawQuery string
	// fragment for references, without '#'
	Fragment string
	// encoded fragment hint (see EscapedFragment method)
	RawFragment string
}

// Decode URL information from a string.
// Based on https://github.com/lintnet/go-jsonnet-native-functions/blob/main/pkg/net/url/url.go
func NativeNetUrl() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "url",
		Params: []jsonnetAst.Identifier{"rawURL"},
		Func: func(args []interface{}) (interface{}, error) {
			rawURL, ok := args[0].(string)
			if !ok {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawURL' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			u, err := url.Parse(rawURL)
			if err != nil {
				return nil, err
			}

			q := u.Query()
			query := make(map[string]any, len(q))
			for k, v := range q {
				a := make([]any, len(v))
				for i, b := range v {
					a[i] = b
				}
				query[k] = a
			}

			pass, passSet := u.User.Password()

			return NativeFuncURL{
				Scheme:      u.Scheme,
				Opaque:      u.Opaque,
				Username:    u.User.Username(),
				Password:    pass,
				PasswordSet: passSet,
				Host:        u.Host,
				Path:        u.Path,
				RawPath:     u.RawPath,
				Query:       query,
				RawQuery:    u.RawQuery,
				Fragment:    u.Fragment,
				RawFragment: u.RawFragment,
			}, nil
		},
	}
}

