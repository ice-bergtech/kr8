package jnetvm

import (
	"fmt"
	"net"
	"net/url"
	"sort"
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
// Based on https://github.com/lintnet/go-jsonnet-native-functions/blob/main/pkg/net/url/url.go .
//
// Inputs: "rawURL".
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

			urlClean, err := url.Parse(rawURL)
			if err != nil {
				return nil, err
			}

			q := urlClean.Query()
			query := make(map[string]any, len(q))
			for k, v := range q {
				a := make([]any, len(v))
				for i, b := range v {
					a[i] = b
				}
				query[k] = a
			}

			pass, passSet := urlClean.User.Password()

			return NativeFuncURL{
				Scheme:      urlClean.Scheme,
				Opaque:      urlClean.Opaque,
				Username:    urlClean.User.Username(),
				Password:    pass,
				PasswordSet: passSet,
				Host:        urlClean.Host,
				Path:        urlClean.Path,
				RawPath:     urlClean.RawPath,
				Query:       query,
				RawQuery:    urlClean.RawQuery,
				Fragment:    urlClean.Fragment,
				RawFragment: urlClean.RawFragment,
			}, nil
		},
	}
}

type IPV4 struct {
	IP           string
	Mask         int
	CIDR         string
	Count        uint32
	FirstAddress string
	LastAddress  string
	Broadcast    string
}

type IPV6 struct {
	IP           string
	NetMask      string
	HostMask     string
	CIDR         string
	Count        uint128.Uint128
	FirstAddress string
	LastAddress  string
}

func IPV6Info(rawIP string) (*IPV6, error) {
	ipa := net.ParseIP(rawIP)
	mask := 128
	if strings.Contains(rawIP, "/") {
		parts := strings.Split(rawIP, "/")
		var err error
		mask, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
	}
	ipNet := iplib.NewNet6(ipa, mask, 0)

	// ipv6 address
	return &IPV6{
		IP:           ipNet.IP().String(),
		NetMask:      ipNet.Mask().String(),
		HostMask:     ipNet.Hostmask.String(),
		CIDR:         ipNet.String(),
		Count:        ipNet.Count(),
		FirstAddress: ipNet.FirstAddress().String(),
		LastAddress:  ipNet.LastAddress().String(),
	}, nil
}

func IPV4Info(rawIP string) (*IPV4, error) {
	ipa := net.ParseIP(rawIP)
	mask := 32
	if strings.Contains(rawIP, "/") {
		parts := strings.Split(rawIP, "/")
		var err error
		mask, err = strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
	}
	ipNet := iplib.NewNet4(ipa, mask)

	// ipv4 address
	return &IPV4{
		IP:           ipNet.IP().String(),
		Mask:         mask,
		CIDR:         ipNet.String(),
		Count:        ipNet.Count(),
		FirstAddress: ipNet.FirstAddress().String(),
		LastAddress:  ipNet.LastAddress().String(),
		Broadcast:    ipNet.BroadcastAddress().String(),
	}, nil
}

// net.IP tools.
// https://github.com/c-robinson/iplib .
//
// Inputs: "rawIP".
func NativeNetIPInfo() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPInfo",
		Params: []jsonnetAst.Identifier{"rawIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, ok := args[0].(string)
			if !ok {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			ipa := net.ParseIP(rawIP)
			if ipa.To4() == nil {
				return IPV6Info(rawIP)
			} else {
				return IPV4Info(rawIP)
			}
		},
	}
}

// Compare two addresses.
//
// 0 if a==b, -1 if a<b, 1 if a>b.
func NativeNetAddressCompare() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPCompare",
		Params: []jsonnetAst.Identifier{"rawIP", "otherIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, pOk := args[0].(string)
			if !pOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}
			otherIP, pOk := args[1].(string)
			if !pOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'otherIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			ipa := net.ParseIP(rawIP)
			ipb := net.ParseIP(otherIP)

			return iplib.CompareIPs(ipa, ipb), nil
		},
	}
}

// Gets the delta of two addresses.
// Takes two net.IP's as input and returns the difference between them up to the limit of uint32.
//
// Inputs: "rawIP, "otherIP".
func NativeNetAddressDelta() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPDelta",
		Params: []jsonnetAst.Identifier{"rawIP", "otherIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, pOk := args[0].(string)
			if !pOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}
			otherIP, pOk := args[1].(string)
			if !pOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'otherIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			ipa := net.ParseIP(rawIP)
			ipb := net.ParseIP(otherIP)

			return iplib.DeltaIP(ipa, ipb), nil
		},
	}
}

// Sort list of ip addresses.
//
// Inputs: "listIPs".
func NativeNetAddressSort() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPDelta",
		Params: []jsonnetAst.Identifier{"listIPs"},
		Func: func(args []interface{}) (interface{}, error) {
			listIPs, ok := args[0].([]string)
			if !ok {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'listIPs' must be of '[]string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			// Marshal items into a net.IP object
			iplist := []net.IP{}
			for _, ip := range listIPs {
				iplist = append(iplist, net.ParseIP(ip))
			}

			// Perform the Sort
			sort.Sort(iplib.ByIP(iplist))

			// Unmarshal into string list
			result := make([]string, len(iplist))
			for i, ipo := range iplist {
				result[i] = ipo.String()
			}

			return result, nil
		},
	}
}

// NextIP returns a net.IP incremented by one from the input address.
// If you overflow the IP space the all-ones address is returned.
//
// Inputs: "rawIP".
func NativeNetAddressInc() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPInc",
		Params: []jsonnetAst.Identifier{"rawIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, ok := args[0].(string)
			if !ok {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			return iplib.NextIP(net.ParseIP(rawIP)), nil
		},
	}
}

// Returns a net.IP that is greater than the supplied net.IP by the supplied integer value.
// If you overflow the IP space the all-ones address is returned.
//
// Inputs: "rawIP", "count".
func NativeNetAddressIncBy() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPIncBy",
		Params: []jsonnetAst.Identifier{"rawIP", "count"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, pOk := args[0].(string)
			if !pOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			count, pOk := args[1].(uint32)
			if !pOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'count' must be of 'uint32' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			return iplib.IncrementIPBy(net.ParseIP(rawIP), count), nil
		},
	}
}

// PreviousIP returns a net.IP decremented by one from the input address.
// If you underflow the IP space the zero address is returned.
//
// Inputs: "rawIP".
func NativeNetAddressDec() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPDec",
		Params: []jsonnetAst.Identifier{"rawIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, ok := args[0].(string)
			if !ok {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			return iplib.PreviousIP(net.ParseIP(rawIP)), nil
		},
	}
}

// Returns a net.IP that is lower than the supplied net.IP by the supplied integer value.
// If you underflow the IP space the zero address is returned.
//
// Inputs: "rawIP", "count".
func NativeNetAddressDecBy() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPDecBy",
		Params: []jsonnetAst.Identifier{"rawIP", "count"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, varOk := args[0].(string)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			count, varOk := args[1].(uint32)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'count' must be of 'uint32' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			return iplib.DecrementIPBy(net.ParseIP(rawIP), count), nil
		},
	}
}

// Convert address to addr.APRA DNS name.
//
// Inputs: "rawIP".
func NativeNetAddressARPA() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPARPA",
		Params: []jsonnetAst.Identifier{"rawIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, varOk := args[0].(string)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			return iplib.IPToARPA(net.ParseIP(rawIP)), nil
		},
	}
}

// Return hex representation of address.
// This is the default stringer format for v6 net.IP.
//
// Inputs: "rawIP".
func NativeNetAddressHex() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPHex",
		Params: []jsonnetAst.Identifier{"rawIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, ok := args[0].(string)
			if !ok {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			return iplib.IPToHexString(net.ParseIP(rawIP)), nil
		},
	}
}

// Return binary string representation of address.
// This is the default stringer format for v6 net.IP.
//
// Inputs: "rawIP".
func NativeNetAddressBinary() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPBinary",
		Params: []jsonnetAst.Identifier{"rawIP"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, ok := args[0].(string)
			if !ok {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'rawIP' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			return iplib.IPToBinaryString(net.ParseIP(rawIP)), nil
		},
	}
}

// Returns a slice of netblocks spanning the range between the two networks, inclusively.
// Returns single-address netblocks if required.
//
// Inputs: "ipNet", "otherIPNet".
func NativeNetAddressNetsBetween() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPNetsBetween",
		Params: []jsonnetAst.Identifier{"ipNet", "otherIPNet"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, varOk := args[0].(string)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'ipNet' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			otherIP, varOk := args[1].(string)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'otherIPNet' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			netsBetween, err := iplib.AllNetsBetween(net.ParseIP(rawIP), net.ParseIP(otherIP))
			if err != nil {
				return nil, err
			}

			// Perform the Sort
			sort.Sort(iplib.ByNet(netsBetween))

			// Unmarshal into string list
			result := make([]string, len(netsBetween))
			for i, ipo := range netsBetween {
				result[i] = ipo.String()
			}

			return result, nil
		},
	}
}

// Return a list of networks of a given masklen that can be extracted from an IPv4 CIDR.
// The mask provided must be a larger-integer than the current mask.
// If set to 0 Subnet carves the network in half.
//
// Inputs: "ip4Net", "maskLen".
func NativeNetAddressCalcSubnetsV4() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPCalcSubnetsV4",
		Params: []jsonnetAst.Identifier{"ip4Net", "maskLen"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, varOk := args[0].(string)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'ip4Net' must be of 'string' in CIDR notation type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			maskResult, varOk := args[1].(int)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'maskLen' must be of 'int' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			// ipv4 address
			return iplib.Net4FromStr(rawIP).Subnet(maskResult)
		},
	}
}

// Return a list of networks of a given masklen that can be extracted from an IPv6 CIDR.
// The mask provided must be a larger-integer than the current mask.
// If set to 0 Subnet carves the network in half.
// Hostmask must be provided if desired.
//
// Inputs: "ip6Net", "netMaskLen", "hostMaskLen".
func NativeNetAddressCalcSubnetsV6() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "netIPCalcSubnetsV6",
		Params: []jsonnetAst.Identifier{"ip6Net", "netMaskLen", "hostMaskLen"},
		Func: func(args []interface{}) (interface{}, error) {
			rawIP, varOk := args[0].(string)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "first argument 'ip6Net' must be of 'string' in CIDR notation type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			netMask, varOk := args[1].(int)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "second argument 'netMaskLen' must be of 'int' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			hostMask, varOk := args[2].(int)
			if !varOk {
				return nil, jsonnet.RuntimeError{
					Msg:        "third argument 'hostMaskLen' must be of 'int' type, got " + fmt.Sprintf("%T", args[0]),
					StackTrace: nil,
				}
			}

			// ipv6 address
			return iplib.Net6FromStr(rawIP).Subnet(netMask, hostMask)
		},
	}
}

// TODO(): expand ipv6 address

// TODO(): map an ipv4 address to an ipv6 address space

// TODO(): retrieve wildcard mask

// TODO(): enumerate all or part of a netblock to []net.IP

// TODO(): decrement or increment addresses within the boundaries of the netblock

// TODO(): return the supernet of a netblock

// TODO(): return next- or previous-adjacent netblocks
