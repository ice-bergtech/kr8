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
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	types "github.com/ice-bergtech/kr8p/pkg/types"
	util "github.com/ice-bergtech/kr8p/pkg/util"
)

func processIncludesFile(
	jvm *jsonnet.VM,
	config string,
	kr8Spec types.Kr8ClusterSpec,
	kr8Opts types.Kr8Opts,
	componentName string,
	componentPath string,
	componentOutputDir string,
	incInfo types.Kr8ComponentSpecIncludeObject,
	outputFileMap map[string]bool,
) error {
	// ensure this directory exists
	outputDir := componentOutputDir
	if incInfo.DestDir != "" {
		outputDir = filepath.Join(kr8Spec.ClusterDir, incInfo.DestDir)
	}
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0750)
		if err := util.GenErrorIfCheck("Error creating alternate directory", err); err != nil {
			return err
		}
	}
	outputFile := filepath.Clean(filepath.Join(outputDir, incInfo.DestName+"."+incInfo.DestExt))
	inputFile := filepath.Clean(filepath.Join(kr8Opts.BaseDir, componentPath, incInfo.File))

	// remember output filename for purging files
	outputFileMap[incInfo.DestName+"."+incInfo.DestExt] = true

	outStr, err := ProcessFile(inputFile, outputFile, kr8Spec, componentName, config, incInfo, jvm)
	if err := util.GenErrorIfCheck("Error processing file", err); err != nil {
		return err
	}

	log.Debug().Str("cluster", kr8Spec.Name).Str("component", componentName).Msg("Checking if file needs updating...")

	// only write file if it does not exist, or the generated contents does not match what is on disk
	updateNeeded, err := CheckIfUpdateNeeded(outputFile, outStr)
	if err != nil {
		return util.GenErrorIfCheck("Error checking if file needs updating", err)
	}
	if updateNeeded {
		file, err := os.Create(outputFile)
		if err := util.GenErrorIfCheck("Error creating file", err); err != nil {
			return err
		}
		_, err = file.WriteString(outStr)
		if err := util.GenErrorIfCheck("Error writing to file", err); err != nil {
			return err
		}

		return util.GenErrorIfCheck("Error closing file", file.Close())
	}

	return nil
}

// Process an includes file.
// Based on the extension, it will process it differently.
//
// .jsonnet: Imported and processed using jsonnet VM.
//
// .yml, .yaml: Imported and processed through native function ParseYaml.
//
// .tpl, .tmpl: Processed using component config and Sprig templating.
func ProcessFile(
	inputFile string,
	outputFile string,
	kr8Spec types.Kr8ClusterSpec,
	componentName string,
	config string,
	incInfo types.Kr8ComponentSpecIncludeObject,
	jvm *jsonnet.VM,
) (string, error) {
	log.Debug().Str("cluster", kr8Spec.Name).
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
		outStr, err = processTemplate(inputFile, gjson.Get(config, componentName).Map())
	default:
		outStr, err = "", os.ErrInvalid
	}
	if err != nil {
		log.Error().Str("cluster", kr8Spec.Name).
			Str("component", componentName).
			Str("file", incInfo.File).
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
	if err := util.GenErrorIfCheck("Error unmarshalling jsonnet output to go slice",
		json.Unmarshal([]byte(jsonStr), &listObjOut),
	); err != nil {
		return "", err
	}
	for _, jObj := range listObjOut {
		buf, err := goyaml.Marshal(jObj)
		if err := util.GenErrorIfCheck("Error marshalling jsonnet object to yaml", err); err != nil {
			return "", err
		}
		outStr += string(buf)
		// Place yml new document at end of each object
		outStr += "\n---\n"
	}

	return outStr, nil
}

func processTemplate(filename string, data map[string]gjson.Result) (string, error) {
	var tInput []byte
	var tmpl *template.Template
	var buffer bytes.Buffer
	var err error

	tInput, err = os.ReadFile(filepath.Clean(filename))
	if err != nil {
		return "Error loading template", err
	}
	tmpl, err = template.New("file").Funcs(sprig.FuncMap()).Parse(string(tInput))
	if err != nil {
		return "Error parsing template", err
	}
	if err = tmpl.Execute(&buffer, data); err != nil {
		return "Error executing templating", err
	}

	return buffer.String(), nil
}
