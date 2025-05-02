package jnetvm

import (
	"fmt"
	"path/filepath"

	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	"github.com/rs/zerolog/log"
	kLogger "github.com/sirupsen/logrus"

	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	types "github.com/ice-bergtech/kr8/pkg/types"
)

type KomposeParams struct {
	// Root directory of the compose files
	RootDir string `json:"rootDir"`
	// The list of compose files to convert.
	ComposeFiles []string `json:"composeFileList"`
	// Namespace to assign to resources.
	Namespace string `json:"namespace"`
}

func (params *KomposeParams) ExtractParameters() {}

// Allows converting a docker-compose string into kubernetes resources using kompose.
// Files in the directory must be in the format `[docker-]compose.y[a]ml`.
//
// Source: https://github.com/kubernetes/kompose/blob/main/cmd/convert.go
//
// Inputs: `rootDir`, `listFiles`, `namespace`.
func NativeKompose() *jsonnet.NativeFunction {
	return &jsonnet.NativeFunction{
		Name:   "komposeFile",
		Params: jsonnetAst.Identifiers{"rootDir", "listFiles", "namespace"},
		Func: func(args []interface{}) (interface{}, error) {
			var argOk bool
			kParams := KomposeParams{}

			kParams.RootDir, argOk = args[0].(string)
			if !argOk {
				return nil, types.Kr8Error{
					Message: "first argument 'rootDir' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					Value:   args[0],
				}
			}

			kParams.ComposeFiles, argOk = args[1].([]string)
			if !argOk {
				return nil, types.Kr8Error{
					Message: "second argument 'listFiles' must be of '[]string' type, got " + fmt.Sprintf("%T", args[0]),
					Value:   args[1],
				}
			}

			kParams.Namespace, argOk = args[2].(string)
			if !argOk {
				return nil, types.Kr8Error{
					Message: "third argument 'namespace' must be of 'string' type, got " + fmt.Sprintf("%T", args[0]),
					Value:   args[2],
				}
			}

			inFiles := make([]string, len(kParams.ComposeFiles))
			for i, s := range kParams.ComposeFiles {
				inFiles[i] = filepath.Join(kParams.RootDir, s)
			}

			// set the logger that kompose uses is set to warn and above
			kLogger.SetLevel(kLogger.WarnLevel)
			kLogger.AddHook(&KomposeHook{})
			komposeConfig, err := kr8_types.CreateKomposeOpts(inFiles, kParams.Namespace)
			if err != nil {
				return "", err
			}

			return komposeConfig.Convert()
		},
	}
}

type KomposeHook struct{}

func (*KomposeHook) Levels() []kLogger.Level {
	return []kLogger.Level{
		kLogger.WarnLevel,
		kLogger.ErrorLevel,
		kLogger.FatalLevel,
		kLogger.PanicLevel,
	}
}

func (*KomposeHook) Fire(entry *kLogger.Entry) error {
	log.Warn().Str("nativeFunc", "kompose").Msg(entry.Message)
	return nil
}
