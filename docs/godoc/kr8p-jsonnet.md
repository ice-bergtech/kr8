# jnetvm

```go
import "github.com/ice-bergtech/kr8p/pkg/jnetvm"
```

Package jvm contains the jsonnet rendering logic.

## Index

- [func JsonnetRender\(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig\) error](<#JsonnetRender>)
- [func JsonnetRenderClusterParams\(vmconfig types.VMConfig, clusterName string, componentNames \[\]string, clusterParams string, prune bool\) \(string, error\)](<#JsonnetRenderClusterParams>)
- [func JsonnetRenderClusterParamsOnly\(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool\) \(string, error\)](<#JsonnetRenderClusterParamsOnly>)
- [func JsonnetRenderFiles\(vmConfig types.VMConfig, files \[\]string, param string, prune bool, prepend string, source string\) \(string, error\)](<#JsonnetRenderFiles>)
- [func JsonnetVM\(vmconfig types.VMConfig\) \(\*jsonnet.VM, error\)](<#JsonnetVM>)
- [func MergeComponentDefaults\(componentMap map\[string\]types.Kr8ClusterComponentRef, componentNames \[\]string, vmconfig types.VMConfig\) \(string, error\)](<#MergeComponentDefaults>)
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
- [type NativeFuncURL](<#NativeFuncURL>)


<a name="JsonnetRender"></a>
## func [JsonnetRender](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L122>)

```go
func JsonnetRender(cmdFlagsJsonnet types.CmdJsonnetOptions, filename string, vmConfig types.VMConfig) error
```

Renders a jsonnet file with the specified options.

<a name="JsonnetRenderClusterParams"></a>
## func [JsonnetRenderClusterParams](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L195-L201>)

```go
func JsonnetRenderClusterParams(vmconfig types.VMConfig, clusterName string, componentNames []string, clusterParams string, prune bool) (string, error)
```

Render cluster params, merged with one or more component's parameters. Empty componentName list renders all component parameters.

<a name="JsonnetRenderClusterParamsOnly"></a>
## func [JsonnetRenderClusterParamsOnly](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L172-L177>)

```go
func JsonnetRenderClusterParamsOnly(vmconfig types.VMConfig, clusterName string, clusterParams string, prune bool) (string, error)
```

Only render cluster params \(\_cluster\), without components.

<a name="JsonnetRenderFiles"></a>
## func [JsonnetRenderFiles](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L73-L80>)

```go
func JsonnetRenderFiles(vmConfig types.VMConfig, files []string, param string, prune bool, prepend string, source string) (string, error)
```

Takes a list of jsonnet files and imports each one. Formats the string for jsonnet using "\+".

<a name="JsonnetVM"></a>
## func [JsonnetVM](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L41>)

```go
func JsonnetVM(vmconfig types.VMConfig) (*jsonnet.VM, error)
```

Create a Jsonnet VM to run commands in.

<a name="MergeComponentDefaults"></a>
## func [MergeComponentDefaults](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/jsonnet.go#L240-L244>)

```go
func MergeComponentDefaults(componentMap map[string]types.Kr8ClusterComponentRef, componentNames []string, vmconfig types.VMConfig) (string, error)
```



<a name="NativeHelmTemplate"></a>
## func [NativeHelmTemplate](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs.go#L97>)

```go
func NativeHelmTemplate() *jsonnet.NativeFunction
```

Allows executing helm template to process a helm chart and make available to kr8p configuration.

Source: https://github.com/grafana/tanka/blob/v0.27.1/pkg/helm/template.go#L23

<a name="NativeHelp"></a>
## func [NativeHelp](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs.go#L64>)

```go
func NativeHelp(allFuncs []*jsonnet.NativeFunction) *jsonnet.NativeFunction
```



<a name="NativeKompose"></a>
## func [NativeKompose](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs.go#L144>)

```go
func NativeKompose() *jsonnet.NativeFunction
```

Allows converting a docker\-compose file string into kubernetes resources using kompose. Files in the directory must be in the format \`\[docker\-\]compose.ym\[a\]l\`.

Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go

Inputs: \`inFile\`, \`outPath\`, \`opts\`.

<a name="NativeNetAddressARPA"></a>
## func [NativeNetAddressARPA](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L399>)

```go
func NativeNetAddressARPA() *jsonnet.NativeFunction
```

Convert address to addr.APRA DNS name.

Inputs: "rawIP".

<a name="NativeNetAddressBinary"></a>
## func [NativeNetAddressBinary](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L443>)

```go
func NativeNetAddressBinary() *jsonnet.NativeFunction
```

Return binary string representation of address. This is the default stringer format for v6 net.IP.

Inputs: "rawIP".

<a name="NativeNetAddressCalcSubnetsV4"></a>
## func [NativeNetAddressCalcSubnetsV4](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L510>)

```go
func NativeNetAddressCalcSubnetsV4() *jsonnet.NativeFunction
```

Return a list of networks of a given masklen that can be extracted from an IPv4 CIDR. The mask provided must be a larger\-integer than the current mask. If set to 0 Subnet will carve the network in half.

Inputs: "ip4Net", "maskLen".

<a name="NativeNetAddressCalcSubnetsV6"></a>
## func [NativeNetAddressCalcSubnetsV6](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L543>)

```go
func NativeNetAddressCalcSubnetsV6() *jsonnet.NativeFunction
```

Return a list of networks of a given masklen that can be extracted from an IPv6 CIDR. The mask provided must be a larger\-integer than the current mask. If set to 0 Subnet will carve the network in half. Hostmask must be provided if desired.

Inputs: "ip6Net", "netMaskLen", "hostMaskLen".

<a name="NativeNetAddressCompare"></a>
## func [NativeNetAddressCompare](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L196>)

```go
func NativeNetAddressCompare() *jsonnet.NativeFunction
```

Compare two addresses.

0 if a==b, \-1 if a\<b, 1 if a\>b.

<a name="NativeNetAddressDec"></a>
## func [NativeNetAddressDec](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L348>)

```go
func NativeNetAddressDec() *jsonnet.NativeFunction
```

PreviousIP returns a net.IP decremented by one from the input address. If you underflow the IP space it will return the zero address.

Inputs: "rawIP".

<a name="NativeNetAddressDecBy"></a>
## func [NativeNetAddressDecBy](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L370>)

```go
func NativeNetAddressDecBy() *jsonnet.NativeFunction
```

Returns a net.IP that is lower than the supplied net.IP by the supplied integer value. If you underflow the IP space it will return the zero address.

Inputs: "rawIP", "count".

<a name="NativeNetAddressDelta"></a>
## func [NativeNetAddressDelta](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L228>)

```go
func NativeNetAddressDelta() *jsonnet.NativeFunction
```

Gets the delta of two addresses. Takes two net.IP's as input and returns the difference between them up to the limit of uint32.

Inputs: "rawIP, "otherIP".

<a name="NativeNetAddressHex"></a>
## func [NativeNetAddressHex](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L421>)

```go
func NativeNetAddressHex() *jsonnet.NativeFunction
```

Return hex representation of address. This is the default stringer format for v6 net.IP.

Inputs: "rawIP".

<a name="NativeNetAddressInc"></a>
## func [NativeNetAddressInc](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L296>)

```go
func NativeNetAddressInc() *jsonnet.NativeFunction
```

NextIP returns a net.IP incremented by one from the input address. If you overflow the IP space it will return the all\-ones address.

Inputs: "rawIP".

<a name="NativeNetAddressIncBy"></a>
## func [NativeNetAddressIncBy](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L318>)

```go
func NativeNetAddressIncBy() *jsonnet.NativeFunction
```

Returns a net.IP that is greater than the supplied net.IP by the supplied integer value. If you overflow the IP space it will return the all\-ones address.

Inputs: "rawIP", "count".

<a name="NativeNetAddressNetsBetween"></a>
## func [NativeNetAddressNetsBetween](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L465>)

```go
func NativeNetAddressNetsBetween() *jsonnet.NativeFunction
```

Returns a slice of netblocks spanning the range between the two networks, inclusively. Returns single\-address netblocks if required.

Inputs: "ipNet", "otherIPNet".

<a name="NativeNetAddressSort"></a>
## func [NativeNetAddressSort](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L259>)

```go
func NativeNetAddressSort() *jsonnet.NativeFunction
```

Sort list of ip addresses.

Inputs: "listIPs".

<a name="NativeNetIPInfo"></a>
## func [NativeNetIPInfo](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L170>)

```go
func NativeNetIPInfo() *jsonnet.NativeFunction
```

net.IP tools. https://github.com/c-robinson/iplib .

Inputs: "rawIP".

<a name="NativeNetUrl"></a>
## func [NativeNetUrl](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L48>)

```go
func NativeNetUrl() *jsonnet.NativeFunction
```

Decode URL information from a string. Based on https://github.com/lintnet/go-jsonnet-native-functions/blob/main/pkg/net/url/url.go .

Inputs: "rawURL".

<a name="NativeRegexEscape"></a>
## func [NativeRegexEscape](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_regex.go#L34>)

```go
func NativeRegexEscape() *jsonnet.NativeFunction
```

Escapes a string for use in regex.

Inputs: "str".

<a name="NativeRegexMatch"></a>
## func [NativeRegexMatch](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_regex.go#L46>)

```go
func NativeRegexMatch() *jsonnet.NativeFunction
```

Matches a string against a regex pattern.

Inputs: "regex", "string".

<a name="NativeRegexSubst"></a>
## func [NativeRegexSubst](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_regex.go#L58>)

```go
func NativeRegexSubst() *jsonnet.NativeFunction
```

Substitutes a regex pattern in a string with another string.

Inputs: "regex", "src", "repl".

<a name="NativeSprigTemplate"></a>
## func [NativeSprigTemplate](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs.go#L105>)

```go
func NativeSprigTemplate() *jsonnet.NativeFunction
```

Uses sprig to process passed in config data and template. Sprig template guide: https://masterminds.github.io/sprig/ .

Inputs: "config" "templateStr".

<a name="RegisterNativeFuncs"></a>
## func [RegisterNativeFuncs](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs.go#L24>)

```go
func RegisterNativeFuncs(jvm *jsonnet.VM)
```

Registers additional native functions in the jsonnet VM. These functions are used to extend the functionality of jsonnet. Adds on to functions part of the jsonnet stdlib: https://jsonnet.org/ref/stdlib.html

<a name="IPV4"></a>
## type [IPV4](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L96-L104>)



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
### func [IPV4Info](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L141>)

```go
func IPV4Info(rawIP string) (*IPV4, error)
```



<a name="IPV6"></a>
## type [IPV6](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L106-L114>)



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
### func [IPV6Info](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L116>)

```go
func IPV6Info(rawIP string) (*IPV6, error)
```



<a name="NativeFuncURL"></a>
## type [NativeFuncURL](<https://github.com/ice-bergtech/kr8/blob/main/pkg/jnetvm/native_funcs_net.go#L18-L42>)

Contains the url information.

```go
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
```