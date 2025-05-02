# kr8\_native\_funcs

```go
import "github.com/ice-bergtech/kr8/pkg/kr8_native_funcs"
```

Package kr8\_native\_funcs provides native functions that jsonnet code can reference. Functions include processing docker\-compose files, helm charts, templating, URL and IP address parsing.

## Index

- [func NativeHelmTemplate\(\) \*jsonnet.NativeFunction](<#NativeHelmTemplate>)
- [func NativeHelp\(allFuncs \[\]\*jsonnet.NativeFunction\) \*jsonnet.NativeFunction](<#NativeHelp>)
- [func NativeKompose\(\) \*jsonnet.NativeFunction](<#NativeKompose>)
- [func NativeNetAddressARPA\(\) \*jsonnet.NativeFunction](<#NativeNetAddressARPA>)
- [func NativeNetAddressBinary\(\) \*jsonnet.NativeFunction](<#NativeNetAddressBinary>)
- [func NativeNetAddressCalcSubnetsV4\(\) \*jsonnet.NativeFunction](<#NativeNetAddressCalcSubnetsV4>)
- [func NativeNetAddressCalcSubnetsV6\(\) \*jsonnet.NativeFunction](<#NativeNetAddressCalcSubnetsV6>)
- [func NativeNetAddressCompare\(\) \*jsonnet.NativeFunction](<#NativeNetAddressCompare>)
- [func NativeNetAddressDec\(\) \*jsonnet.NativeFunction](<#NativeNetAddressDec>)
- [func NativeNetAddressDecBy\(\) \*jsonnet.NativeFunction](<#NativeNetAddressDecBy>)
- [func NativeNetAddressDelta\(\) \*jsonnet.NativeFunction](<#NativeNetAddressDelta>)
- [func NativeNetAddressHex\(\) \*jsonnet.NativeFunction](<#NativeNetAddressHex>)
- [func NativeNetAddressInc\(\) \*jsonnet.NativeFunction](<#NativeNetAddressInc>)
- [func NativeNetAddressIncBy\(\) \*jsonnet.NativeFunction](<#NativeNetAddressIncBy>)
- [func NativeNetAddressNetsBetween\(\) \*jsonnet.NativeFunction](<#NativeNetAddressNetsBetween>)
- [func NativeNetAddressSort\(\) \*jsonnet.NativeFunction](<#NativeNetAddressSort>)
- [func NativeNetIPInfo\(\) \*jsonnet.NativeFunction](<#NativeNetIPInfo>)
- [func NativeNetUrl\(\) \*jsonnet.NativeFunction](<#NativeNetUrl>)
- [func NativeRegexEscape\(\) \*jsonnet.NativeFunction](<#NativeRegexEscape>)
- [func NativeRegexMatch\(\) \*jsonnet.NativeFunction](<#NativeRegexMatch>)
- [func NativeRegexSubst\(\) \*jsonnet.NativeFunction](<#NativeRegexSubst>)
- [func NativeSprigTemplate\(\) \*jsonnet.NativeFunction](<#NativeSprigTemplate>)
- [func RegisterNativeFuncs\(jvm \*jsonnet.VM\)](<#RegisterNativeFuncs>)
- [type IPV4](<#IPV4>)
  - [func IPV4Info\(rawIP string\) \(\*IPV4, error\)](<#IPV4Info>)
- [type IPV6](<#IPV6>)
  - [func IPV6Info\(rawIP string\) \(\*IPV6, error\)](<#IPV6Info>)
- [type KomposeHook](<#KomposeHook>)
  - [func \(\*KomposeHook\) Fire\(entry \*kLogger.Entry\) error](<#KomposeHook.Fire>)
  - [func \(\*KomposeHook\) Levels\(\) \[\]kLogger.Level](<#KomposeHook.Levels>)
- [type KomposeParams](<#KomposeParams>)
  - [func ParseKomposeParams\(args \[\]interface\{\}\) \(\*KomposeParams, error\)](<#ParseKomposeParams>)
  - [func \(params \*KomposeParams\) ExtractParameters\(\)](<#KomposeParams.ExtractParameters>)
- [type NativeFuncURL](<#NativeFuncURL>)


<a name="NativeHelmTemplate"></a>
## func [NativeHelmTemplate](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs.go#L96>)

```go
func NativeHelmTemplate() *jsonnet.NativeFunction
```

Allows executing helm template to process a helm chart and make available to kr8 configuration.

Source: https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23

<a name="NativeHelp"></a>
## func [NativeHelp](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs.go#L63>)

```go
func NativeHelp(allFuncs []*jsonnet.NativeFunction) *jsonnet.NativeFunction
```



<a name="NativeKompose"></a>
## func [NativeKompose](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_kompose.go#L34>)

```go
func NativeKompose() *jsonnet.NativeFunction
```

Allows converting a docker\-compose string into kubernetes resources using kompose. Files in the directory must be in the format \`\[docker\-\]compose.y\[a\]ml\`. RootDir is usually \`std.thisFile\(\)\`.

Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go

Inputs: \`rootDir\`, \`listFiles\`, \`namespace\`.

<a name="NativeNetAddressARPA"></a>
## func [NativeNetAddressARPA](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L404>)

```go
func NativeNetAddressARPA() *jsonnet.NativeFunction
```

Convert address to addr.APRA DNS name.

Inputs: "rawIP".

<a name="NativeNetAddressBinary"></a>
## func [NativeNetAddressBinary](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L448>)

```go
func NativeNetAddressBinary() *jsonnet.NativeFunction
```

Return binary string representation of address. This is the default stringer format for v6 net.IP.

Inputs: "rawIP".

<a name="NativeNetAddressCalcSubnetsV4"></a>
## func [NativeNetAddressCalcSubnetsV4](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L515>)

```go
func NativeNetAddressCalcSubnetsV4() *jsonnet.NativeFunction
```

Return a list of networks of a given masklen that can be extracted from an IPv4 CIDR. The mask provided must be a larger\-integer than the current mask. If set to 0 Subnet carves the network in half.

Inputs: "ip4Net", "maskLen".

<a name="NativeNetAddressCalcSubnetsV6"></a>
## func [NativeNetAddressCalcSubnetsV6](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L548>)

```go
func NativeNetAddressCalcSubnetsV6() *jsonnet.NativeFunction
```

Return a list of networks of a given masklen that can be extracted from an IPv6 CIDR. The mask provided must be a larger\-integer than the current mask. If set to 0 Subnet carves the network in half. Hostmask must be provided if desired.

Inputs: "ip6Net", "netMaskLen", "hostMaskLen".

<a name="NativeNetAddressCompare"></a>
## func [NativeNetAddressCompare](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L201>)

```go
func NativeNetAddressCompare() *jsonnet.NativeFunction
```

Compare two addresses.

0 if a==b, \-1 if a\<b, 1 if a\>b.

<a name="NativeNetAddressDec"></a>
## func [NativeNetAddressDec](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L353>)

```go
func NativeNetAddressDec() *jsonnet.NativeFunction
```

PreviousIP returns a net.IP decremented by one from the input address. If you underflow the IP space the zero address is returned.

Inputs: "rawIP".

<a name="NativeNetAddressDecBy"></a>
## func [NativeNetAddressDecBy](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L375>)

```go
func NativeNetAddressDecBy() *jsonnet.NativeFunction
```

Returns a net.IP that is lower than the supplied net.IP by the supplied integer value. If you underflow the IP space the zero address is returned.

Inputs: "rawIP", "count".

<a name="NativeNetAddressDelta"></a>
## func [NativeNetAddressDelta](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L233>)

```go
func NativeNetAddressDelta() *jsonnet.NativeFunction
```

Gets the delta of two addresses. Takes two net.IP's as input and returns the difference between them up to the limit of uint32.

Inputs: "rawIP, "otherIP".

<a name="NativeNetAddressHex"></a>
## func [NativeNetAddressHex](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L426>)

```go
func NativeNetAddressHex() *jsonnet.NativeFunction
```

Return hex representation of address. This is the default stringer format for v6 net.IP.

Inputs: "rawIP".

<a name="NativeNetAddressInc"></a>
## func [NativeNetAddressInc](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L301>)

```go
func NativeNetAddressInc() *jsonnet.NativeFunction
```

NextIP returns a net.IP incremented by one from the input address. If you overflow the IP space the all\-ones address is returned.

Inputs: "rawIP".

<a name="NativeNetAddressIncBy"></a>
## func [NativeNetAddressIncBy](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L323>)

```go
func NativeNetAddressIncBy() *jsonnet.NativeFunction
```

Returns a net.IP that is greater than the supplied net.IP by the supplied integer value. If you overflow the IP space the all\-ones address is returned.

Inputs: "rawIP", "count".

<a name="NativeNetAddressNetsBetween"></a>
## func [NativeNetAddressNetsBetween](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L470>)

```go
func NativeNetAddressNetsBetween() *jsonnet.NativeFunction
```

Returns a slice of netblocks spanning the range between the two networks, inclusively. Returns single\-address netblocks if required.

Inputs: "ipNet", "otherIPNet".

<a name="NativeNetAddressSort"></a>
## func [NativeNetAddressSort](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L264>)

```go
func NativeNetAddressSort() *jsonnet.NativeFunction
```

Sort list of ip addresses.

Inputs: "listIPs".

<a name="NativeNetIPInfo"></a>
## func [NativeNetIPInfo](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L175>)

```go
func NativeNetIPInfo() *jsonnet.NativeFunction
```

net.IP tools. https://github.com/c-robinson/iplib .

Inputs: "rawIP".

<a name="NativeNetUrl"></a>
## func [NativeNetUrl](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L49>)

```go
func NativeNetUrl() *jsonnet.NativeFunction
```

Decode URL information from a string. Based on https://github.com/lintnet/go-jsonnet-native-functions/blob/main/pkg/net/url/url.go .

Inputs: "rawURL".

<a name="NativeRegexEscape"></a>
## func [NativeRegexEscape](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_regex.go#L34>)

```go
func NativeRegexEscape() *jsonnet.NativeFunction
```

Escapes a string for use in regex.

Inputs: "str".

<a name="NativeRegexMatch"></a>
## func [NativeRegexMatch](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_regex.go#L46>)

```go
func NativeRegexMatch() *jsonnet.NativeFunction
```

Matches a string against a regex pattern.

Inputs: "regex", "string".

<a name="NativeRegexSubst"></a>
## func [NativeRegexSubst](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_regex.go#L58>)

```go
func NativeRegexSubst() *jsonnet.NativeFunction
```

Substitutes a regex pattern in a string with another string.

Inputs: "regex", "src", "repl".

<a name="NativeSprigTemplate"></a>
## func [NativeSprigTemplate](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs.go#L104>)

```go
func NativeSprigTemplate() *jsonnet.NativeFunction
```

Uses sprig to process passed in config data and template. Sprig template guide: https://masterminds.github.io/sprig/ .

Inputs: "config" "templateStr".

<a name="RegisterNativeFuncs"></a>
## func [RegisterNativeFuncs](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs.go#L23>)

```go
func RegisterNativeFuncs(jvm *jsonnet.VM)
```

Registers additional native functions in the jsonnet VM. These functions are used to extend the functionality of jsonnet. Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html

<a name="IPV4"></a>
## type [IPV4](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L98-L106>)

An IPv4 address or subnet range.

```go
type IPV4 struct {
    IP           string
    Mask         int
    CIDR         string
    Count        uint32
    FirstAddress string
    LastAddress  string
    Broadcast    string
}
```

<a name="IPV4Info"></a>
### func [IPV4Info](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L146>)

```go
func IPV4Info(rawIP string) (*IPV4, error)
```

Parses an IPv4 object from a string.

<a name="IPV6"></a>
## type [IPV6](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L109-L117>)

An IPv6 address or subnet range.

```go
type IPV6 struct {
    IP           string
    NetMask      string
    HostMask     string
    CIDR         string
    Count        uint128.Uint128
    FirstAddress string
    LastAddress  string
}
```

<a name="IPV6Info"></a>
### func [IPV6Info](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L120>)

```go
func IPV6Info(rawIP string) (*IPV6, error)
```

Parses an IPv6 object from a string.

<a name="KomposeHook"></a>
## type [KomposeHook](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_kompose.go#L115>)



```go
type KomposeHook struct{}
```

<a name="KomposeHook.Fire"></a>
### func \(\*KomposeHook\) [Fire](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_kompose.go#L126>)

```go
func (*KomposeHook) Fire(entry *kLogger.Entry) error
```



<a name="KomposeHook.Levels"></a>
### func \(\*KomposeHook\) [Levels](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_kompose.go#L117>)

```go
func (*KomposeHook) Levels() []kLogger.Level
```



<a name="KomposeParams"></a>
## type [KomposeParams](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_kompose.go#L16-L23>)



```go
type KomposeParams struct {
    // Root directory of the compose files
    RootDir string `json:"rootDir"`
    // The list of compose files to convert.
    ComposeFiles []string `json:"composeFileList"`
    // Namespace to assign to resources.
    Namespace string `json:"namespace"`
}
```

<a name="ParseKomposeParams"></a>
### func [ParseKomposeParams](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_kompose.go#L66>)

```go
func ParseKomposeParams(args []interface{}) (*KomposeParams, error)
```



<a name="KomposeParams.ExtractParameters"></a>
### func \(\*KomposeParams\) [ExtractParameters](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_kompose.go#L25>)

```go
func (params *KomposeParams) ExtractParameters()
```



<a name="NativeFuncURL"></a>
## type [NativeFuncURL](<https://github.com:icebergtech/kr8/blob/main/pkg/kr8_native_funcs/native_funcs_net.go#L19-L43>)

Contains the url information.

```go
type NativeFuncURL struct {
    Scheme string `json:"scheme"`
    // encoded opaque data
    Opaque string `json:"opaque"`
    // username information
    Username string `json:"username"`
    // Whether the password field is set
    PasswordSet bool `json:"passwordSet"`
    // password information
    Password string `json:"password"`
    // host or host:port (see Hostname and Port methods)
    Host string `json:"host"`
    // path (relative paths may omit leading slash)
    Path string `json:"path"`
    // encoded path hint (see EscapedPath method)
    RawPath string `json:"pathRaw"`
    // query values
    Query map[string]interface{} `json:"query"`
    // encoded query values, without '?'
    RawQuery string `json:"queryRaw"`
    // fragment for references, without '#'
    Fragment string `json:"fragment"`
    // encoded fragment hint (see EscapedFragment method)
    RawFragment string `json:"fragmentRaw"`
}
```