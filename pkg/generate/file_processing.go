package generate

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	goyaml "github.com/ghodss/yaml"
	jsonnet "github.com/google/go-jsonnet"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"

	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
)

func processIncludesFile(
	jvm *jsonnet.VM,
	config string,
	kr8Spec kr8_types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	componentName string,
	componentPath string,
	componentOutputDir string,
	incInfo kr8_types.Kr8ComponentSpecIncludeObject,
	outputFileMap map[string]bool,
	logger zerolog.Logger,
) error {
	// ensure this directory exists
	outputDir := componentOutputDir
	if incInfo.DestDir != "" {
		outputDir = filepath.Join(componentOutputDir, incInfo.DestDir)
		logger.Debug().Msg("includes destdir override: " + outputDir)
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0750)
		if err := util.ErrorIfCheck("error creating alternate directory", err); err != nil {
			return err
		}
	}
	inputFile := filepath.Join(kr8Opts.BaseDir, componentPath, incInfo.File)
	outputFile := filepath.Join(outputDir, filepath.Base(incInfo.DestName+"."+incInfo.DestExt))
	// remember output filename for purging files
	outputFileMap[filepath.Base(incInfo.DestName+"."+incInfo.DestExt)] = true

	outStr, err := ProcessFile(inputFile, outputFile, kr8Spec, componentName, config, incInfo, jvm, logger)
	if err := util.ErrorIfCheck("error processing file", err); err != nil {
		return err
	}

	logger.Debug().Str("cluster", kr8Spec.Name).Str("component", componentName).Msg("Checking if file needs updating...")

	// only write file if it does not exist, or the generated contents does not match what is on disk
	updateNeeded, err := CheckIfUpdateNeeded(outputFile, outStr)
	if err != nil {
		return util.ErrorIfCheck("Error checking if file needs updating", err)
	}
	if updateNeeded {
		return os.WriteFile(outputFile, []byte(outStr), 0600)
	}

	return nil
}

// Process an includes file.
// Based on the extension, the file is processed differently.
//   - .jsonnet: Imported and processed using jsonnet VM.
//   - .yml, .yaml: Imported and processed through native function ParseYaml.
//   - .tpl, .tmpl: Processed using component config and Sprig templating.
func ProcessFile(
	inputFile string,
	outputFile string,
	kr8Spec kr8_types.Kr8ClusterSpec,
	componentName string,
	config string,
	incInfo kr8_types.Kr8ComponentSpecIncludeObject,
	jvm *jsonnet.VM,
	logger zerolog.Logger,
) (string, error) {
	logger.Debug().Str("cluster", kr8Spec.Name).
		Str("component", componentName).
		Msg("Process file: " + inputFile + " -> " + outputFile)

	file_extension := filepath.Ext(incInfo.File)

	var input string
	var outStr string
	var err error
	switch file_extension {
	case ".jsonnet":
		// file is processed as an ExtCode input, so that we can postprocess it
		// in the snippet
		input = "( import '" + inputFile + "')"
		outStr, err = processJsonnet(jvm, input, incInfo.File)
	case ".yml":
	case ".yaml":
		input = "std.native('parseYaml')(importstr '" + inputFile + "')"
		outStr, err = processJsonnet(jvm, input, incInfo.File)
	case ".tmpl":
	case ".tpl":
		// Pass component config as data for the template
		if len(incInfo.Config) > 0 {
			outStr, err = processTemplate(inputFile, gjson.Parse(incInfo.Config))
		} else {
			outStr, err = processTemplate(inputFile, gjson.Get(config, componentName))
		}
	default:
		outStr, err = "", os.ErrInvalid
	}
	if err != nil {
		logger.Error().
			Err(err).
			Msg(outStr)
	}

	return outStr, err
}

func processJsonnet(jvm *jsonnet.VM, input string, snippetFilename string) (string, error) {
	jvm.ExtCode("input", input)
	jsonStr, err := jvm.EvaluateAnonymousSnippet(snippetFilename, "std.extVar('process')(std.extVar('input'))")
	if err != nil {
		return "Error evaluating jsonnet snippet", err
	}

	// create output file contents in a string first, as a yaml stream
	var listObjOut []interface{}
	var outStr string
	if err := util.ErrorIfCheck("Error unmarshalling jsonnet output to go slice",
		json.Unmarshal([]byte(jsonStr), &listObjOut),
	); err != nil {
		return "", err
	}
	for _, jObj := range listObjOut {
		buf, err := goyaml.Marshal(jObj)
		if err := util.ErrorIfCheck("Error marshalling jsonnet object to yaml", err); err != nil {
			return "", err
		}
		outStr += string(buf)
		// Place yml new document at end of each object
		outStr += "\n---\n"
	}

	return outStr, nil
}

func processTemplate(filename string, data gjson.Result) (string, error) {
	var tInput []byte
	var tmpl *template.Template
	var buffer bytes.Buffer
	var err error

	tInput, err = os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return "Error loading template", err
	}
	tmpl, err = template.New(filepath.Base(filename)).Funcs(sprig.FuncMap()).Parse(string(tInput))
	if err != nil {
		return "Error parsing template", err
	}
	if err = tmpl.Execute(&buffer, data.Value().(map[string]interface{})); err != nil {
		return "Error executing templating", err
	}

	return buffer.String(), nil
}
