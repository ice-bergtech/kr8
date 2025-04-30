package kr8_types

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/ice-bergtech/kr8/pkg/types"
	"github.com/ice-bergtech/kr8/pkg/util"
	"github.com/tidwall/gjson"
)

func TestCreateClusterSpec(t *testing.T) {
	tests := []struct {
		name               string
		clusterName        string
		spec               gjson.Result
		kr8Opts            types.Kr8Opts
		genDirOverride     string
		wantKr8ClusterSpec Kr8ClusterSpec
		wantErr            bool
	}{
		{
			name:        "default values",
			clusterName: "test-cluster",
			spec: gjson.Parse(`{
					"postprocessor": "",
					"generate_dir": ""
			}`),
			kr8Opts: types.Kr8Opts{
				BaseDir:      "/path/to/kr8",
				ComponentDir: "",
				ClusterDir:   "",
			},
			genDirOverride: "",
			wantKr8ClusterSpec: Kr8ClusterSpec{
				PostProcessor:      "",
				GenerateDir:        "/path/to/kr8/generated",
				GenerateShortNames: false,
				PruneParams:        false,
				ClusterOutputDir:   filepath.Join("/path/to/kr8/generated/", "test-cluster"),
				Name:               "test-cluster",
				EnableCache:        false,
				CompressCache:      true,
			},
			wantErr: false,
		},
		{
			name:        "custom values",
			clusterName: "test-cluster",
			spec: gjson.Parse(`{
					"postprocessor": "function(input) input",
					"generate_dir": "/path/to/custom/dir"
			}`),
			kr8Opts: types.Kr8Opts{
				BaseDir:      "/path/to/kr8",
				ComponentDir: "",
				ClusterDir:   "",
			},
			genDirOverride: "",
			wantKr8ClusterSpec: Kr8ClusterSpec{
				PostProcessor:      "function(input) input",
				GenerateDir:        "/path/to/custom/dir",
				GenerateShortNames: false,
				PruneParams:        false,
				ClusterOutputDir:   "/path/to/custom/dir/test-cluster",
				Name:               "test-cluster",
				EnableCache:        false,
				CompressCache:      true,
			},
			wantErr: false,
		},
		{
			name:        "custom values relative generate dir",
			clusterName: "test-cluster",
			spec: gjson.Parse(`{
					"postprocessor": "function(input) input",
					"generate_dir": "rel/custom/dir"
			}`),
			kr8Opts: types.Kr8Opts{
				BaseDir:      "/path/to/kr8",
				ComponentDir: "",
				ClusterDir:   "",
			},
			genDirOverride: "",
			wantKr8ClusterSpec: Kr8ClusterSpec{
				PostProcessor:      "function(input) input",
				GenerateDir:        "/path/to/kr8/rel/custom/dir",
				GenerateShortNames: false,
				PruneParams:        false,
				ClusterOutputDir:   filepath.Join("/path/to/kr8/rel/custom/dir", "test-cluster"),
				Name:               "test-cluster",
				EnableCache:        false,
				CompressCache:      true,
			},
			wantErr: false,
		},
		{
			name:        "genDirOverride set relative path",
			clusterName: "test-cluster",
			spec: gjson.Parse(`{
					"postprocessor": "",
					"generate_dir": "/path/to/custom/dir"
			}`),
			kr8Opts: types.Kr8Opts{
				BaseDir:      "/path/to/kr8",
				ComponentDir: "",
				ClusterDir:   "",
			},
			genDirOverride: "alt/gen/dir",
			wantKr8ClusterSpec: Kr8ClusterSpec{
				PostProcessor:      "",
				GenerateDir:        "/path/to/kr8/alt/gen/dir",
				GenerateShortNames: false,
				PruneParams:        false,
				ClusterOutputDir:   filepath.Join("/path/to/kr8/alt/gen/dir", "test-cluster"),
				Name:               "test-cluster",
				EnableCache:        false,
				CompressCache:      true,
			},
			wantErr: false,
		},
		{
			name:        "genDirOverride set absolute path",
			clusterName: "test-cluster",
			spec: gjson.Parse(`{
					"postprocessor": "",
					"generate_dir": "/path/to/custom/dir"
			}`),
			kr8Opts: types.Kr8Opts{
				BaseDir:      "/path/to/kr8",
				ComponentDir: "",
				ClusterDir:   "",
			},
			genDirOverride: "/absolute/path/to/gen/dir",
			wantKr8ClusterSpec: Kr8ClusterSpec{
				PostProcessor:      "",
				GenerateDir:        "/absolute/path/to/gen/dir",
				GenerateShortNames: false,
				PruneParams:        false,
				ClusterOutputDir:   filepath.Join("/absolute/path/to/gen/dir", "test-cluster"),
				Name:               "test-cluster",
				EnableCache:        false,
				CompressCache:      true,
			},
			wantErr: false,
		},
	}

	for _, testEntry := range tests {
		t.Run(testEntry.name, func(t *testing.T) {
			gotKr8ClusterSpec, err := CreateClusterSpec(
				testEntry.clusterName,
				testEntry.spec,
				testEntry.kr8Opts,
				testEntry.genDirOverride,
				util.SetupLogger(true),
			)
			if (err != nil) != testEntry.wantErr {
				t.Errorf("CreateClusterSpec() `%v` error = \n%v\n-wantErr-\n%v", testEntry.name, err, testEntry.wantErr)

				return
			}
			if !reflect.DeepEqual(gotKr8ClusterSpec, testEntry.wantKr8ClusterSpec) {
				t.Errorf("CreateClusterSpec() `%v` got Kr8ClusterSpec = \n%v\n-want-\n%v",
					testEntry.name,
					gotKr8ClusterSpec,
					testEntry.wantKr8ClusterSpec,
				)
			}
		})
	}
}

func TestExtractExtFiles(t *testing.T) {
	tests := []struct {
		name string
		spec gjson.Result
		want map[string]string
	}{
		{
			name: "empty spec",
			spec: gjson.Parse(`{}`),
			want: make(map[string]string),
		},
		{
			name: "single extfile",
			spec: gjson.Parse(`{
				"extfiles": {
					"var1": "/path/to/file1"
				}
			}`),
			want: map[string]string{
				"var1": "/path/to/file1",
			},
		},
		{
			name: "multiple extfiles",
			spec: gjson.Parse(`{
				"extfiles": {
					"var1": "/path/to/file1",
					"var2": "/path/to/file2"
				}
			}`),
			want: map[string]string{
				"var1": "/path/to/file1",
				"var2": "/path/to/file2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractExtFiles(tt.spec)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractExtFiles() `%v` got = \n%v\n-want-\n%v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExtractJpaths(t *testing.T) {
	tests := []struct {
		name string
		spec gjson.Result
		want []string
	}{
		{
			name: "empty spec",
			spec: gjson.Parse(`{}`),
			want: make([]string, 0),
		},
		{
			name: "single jpath",
			spec: gjson.Parse(`{
				"jpaths": ["/path/to/jpath"]
			}`),
			want: []string{"/path/to/jpath"},
		},
		{
			name: "multiple jpaths",
			spec: gjson.Parse(`{
				"jpaths": ["/path/to/jpath1", "/path/to/jpath2"]
			}`),
			want: []string{"/path/to/jpath1", "/path/to/jpath2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractJpaths(tt.spec)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractJpaths() `%v` got = \n%v\n-want-\n%v", tt.name, got, tt.want)
			}
		})
	}
}

func TestExtractIncludes(t *testing.T) {
	tests := []struct {
		name string
		spec gjson.Result
		want Kr8ComponentSpecIncludes
	}{
		{
			name: "empty spec",
			spec: gjson.Parse(`{}`),
			want: make(Kr8ComponentSpecIncludes, 0),
		},
		{
			name: "single include file",
			spec: gjson.Parse(`{
				"includes": ["/path/to/file1"]
			}`),
			want: Kr8ComponentSpecIncludes{
				Kr8ComponentSpecIncludeObject{
					File:     "/path/to/file1",
					DestName: "/path/to/file1",
					DestExt:  "yaml",
					DestDir:  "",
					Config:   "",
				},
			},
		},
		{
			name: "multiple include files",
			spec: gjson.Parse(`{
				"includes": ["/path/to/file1", "/path/to/file2"]
			}`),
			want: Kr8ComponentSpecIncludes{
				Kr8ComponentSpecIncludeObject{
					File:     "/path/to/file1",
					DestName: "/path/to/file1",
					DestExt:  "yaml",
					DestDir:  "",
					Config:   "",
				},
				Kr8ComponentSpecIncludeObject{
					File:     "/path/to/file2",
					DestName: "/path/to/file2",
					DestExt:  "yaml",
					DestDir:  "",
					Config:   "",
				},
			},
		},
		{
			name: "single include object",
			spec: gjson.Parse(`{
				"includes": [
					{
						"file": "/path/to/file1"
					}
				]
			}`),
			want: Kr8ComponentSpecIncludes{
				Kr8ComponentSpecIncludeObject{
					File:     "/path/to/file1",
					DestDir:  "",
					Config:   "",
					DestExt:  "",
					DestName: "",
				},
			},
		},
		{
			name: "multiple include objects",
			spec: gjson.Parse(`{
				"includes": [
					{
						"file": "/path/to/file1"
					},
					{
						"file": "/path/to/file2"
					}
				]
			}`),
			want: Kr8ComponentSpecIncludes{
				Kr8ComponentSpecIncludeObject{
					File:     "/path/to/file1",
					DestDir:  "",
					Config:   "",
					DestExt:  "",
					DestName: "",
				},
				Kr8ComponentSpecIncludeObject{
					File:     "/path/to/file2",
					DestDir:  "",
					Config:   "",
					DestExt:  "",
					DestName: "",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := ExtractIncludes(tt.spec)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractIncludes() `%v` got = \n%v\n-want-\n%v", tt.name, got, tt.want)
			}
		})
	}
}

func TestCreateComponentSpec(t *testing.T) {
	tests := []struct {
		name string
		spec gjson.Result
		want Kr8ComponentSpec
		err  bool
	}{
		{
			name: "empty spec",
			spec: gjson.Parse(`{}`),
			//nolint:exhaustruct
			want: Kr8ComponentSpec{},
			err:  true,
		},
		{
			name: "single boolean option",
			spec: gjson.Parse(`{
				"enable_kr8_allparams": true
			}`),
			//nolint:exhaustruct
			want: Kr8ComponentSpec{
				Kr8_allParams: true,
			},
			err: false,
		},
		{
			name: "multiple boolean options",
			spec: gjson.Parse(`{
				"enable_kr8_allparams": true,
				"enable_kr8_allclusters": false
			}`),
			//nolint:exhaustruct
			want: Kr8ComponentSpec{
				Kr8_allParams:   true,
				Kr8_allClusters: false,
			},
			err: false,
		},
		{
			name: "string option",
			spec: gjson.Parse(`{
				"disable_output_clean": false
			}`),
			//nolint:exhaustruct
			want: Kr8ComponentSpec{
				DisableOutputDirClean: false,
			},
			err: false,
		},
		{
			name: "extfiles option",
			spec: gjson.Parse(`{
				"extfiles": {
					"var1": "/path/to/file1"
				}
			}`),
			//nolint:exhaustruct
			want: Kr8ComponentSpec{
				ExtFiles: map[string]string{
					"var1": "/path/to/file1",
				},
			},
			err: false,
		},
		{
			name: "jpaths option",
			spec: gjson.Parse(`{
				"jpaths": ["/path/to/jpath"]
			}`),
			//nolint:exhaustruct
			want: Kr8ComponentSpec{
				JPaths: []string{"/path/to/jpath"},
			},
			err: false,
		},
		{
			name: "includes option with string values",
			spec: gjson.Parse(`{
				"includes": ["/path/to/file1", "/path/to/file2"]
			}`),
			//nolint:exhaustruct
			want: Kr8ComponentSpec{
				Includes: Kr8ComponentSpecIncludes{
					Kr8ComponentSpecIncludeObject{
						File:     "/path/to/file1",
						DestDir:  "",
						Config:   "",
						DestExt:  "",
						DestName: ""},
					Kr8ComponentSpecIncludeObject{
						File:     "/path/to/file2",
						DestDir:  "",
						Config:   "",
						DestExt:  "",
						DestName: ""},
				},
			},
			err: false,
		},
	}

	logger := util.SetupLogger(true)
	for _, testEntry := range tests {
		t.Run(testEntry.name, func(t *testing.T) {
			_, err := CreateComponentSpec(testEntry.spec, logger)
			if err != nil && !testEntry.err {
				t.Errorf("CreateComponentSpec() error = %v", err)

				return
			}
			// TODO: fix includes interfaces
			// if !reflect.DeepEqual(got, testEntry.want) {
			// 	t.Errorf("CreateComponentSpec() `%v` got = \n%v\n-want-\n%v", testEntry.name, got, testEntry.want)
			// }
		})
	}
}
