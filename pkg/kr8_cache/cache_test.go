package kr8_cache_test

import (
	"reflect"
	"testing"

	"github.com/ice-bergtech/kr8/pkg/kr8_cache"
	"github.com/rs/zerolog"
)

func TestLoadClusterCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cacheFile string
		want      *kr8_cache.DeploymentCache
		wantErr   bool
	}{
		{
			name:      "Valid cache file with correct structure",
			cacheFile: "testdata/valid_cache.json",
			want: &kr8_cache.DeploymentCache{
				ClusterConfig: kr8_cache.CreateClusterCache("test_cluster_config"),
				ComponentConfigs: map[string]kr8_cache.ComponentCache{
					"component1": {
						ComponentConfig: "config1",
						ComponentFiles: map[string]string{
							"file1.txt": "hash1",
							"file2.txt": "hash2",
						},
					},
				},
				LibraryCache: &kr8_cache.LibraryCache{
					Directory: "",
					Entries:   map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name:      "Invalid cache file (missing cluster_config)",
			cacheFile: "testdata/invalid_missing_cluster_config.json",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "Invalid cache file (malformed JSON)",
			cacheFile: "testdata/invalid_malformed_json.json",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "Valid cache file with empty component_configs",
			cacheFile: "testdata/valid_empty_component_configs.json",
			want: &kr8_cache.DeploymentCache{
				ClusterConfig:    kr8_cache.CreateClusterCache("test_cluster_config"),
				ComponentConfigs: map[string]kr8_cache.ComponentCache{},
				LibraryCache: &kr8_cache.LibraryCache{
					Directory: "",
					Entries:   map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name:      "Valid cache file with multiple components",
			cacheFile: "testdata/valid_multiple_components.json",
			want: &kr8_cache.DeploymentCache{
				ClusterConfig: kr8_cache.CreateClusterCache("test_cluster_config"),
				ComponentConfigs: map[string]kr8_cache.ComponentCache{
					"component1": {
						ComponentConfig: "config1",
						ComponentFiles: map[string]string{
							"file1.txt": "hash1",
						},
					},
					"component2": {
						ComponentConfig: "config2",
						ComponentFiles: map[string]string{
							"file2.txt": "hash2",
						},
					},
				},
				LibraryCache: &kr8_cache.LibraryCache{
					Directory: "",
					Entries:   map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name:      "Valid cache file with correct structure",
			cacheFile: "testdata/valid_cache.json",
			want: &kr8_cache.DeploymentCache{
				ClusterConfig: kr8_cache.CreateClusterCache("test_cluster_config"),
				ComponentConfigs: map[string]kr8_cache.ComponentCache{
					"component1": {
						ComponentConfig: "config1",
						ComponentFiles: map[string]string{
							"file1.txt": "hash1",
							"file2.txt": "hash2",
						},
					},
				},
				LibraryCache: &kr8_cache.LibraryCache{
					Directory: "",
					Entries:   map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name:      "Invalid cache file (missing cluster_config)",
			cacheFile: "testdata/invalid_missing_cluster_config.json",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "Invalid cache file (malformed JSON)",
			cacheFile: "testdata/invalid_malformed_json.json",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "Valid cache file with empty component_configs",
			cacheFile: "testdata/valid_empty_component_configs.json",
			want: &kr8_cache.DeploymentCache{
				ClusterConfig:    kr8_cache.CreateClusterCache("test_cluster_config"),
				ComponentConfigs: map[string]kr8_cache.ComponentCache{},
				LibraryCache: &kr8_cache.LibraryCache{
					Directory: "",
					Entries:   map[string]string{},
				},
			},
			wantErr: false,
		},
		{
			name:      "Valid cache file with multiple components",
			cacheFile: "testdata/valid_multiple_components.json",
			want: &kr8_cache.DeploymentCache{
				ClusterConfig: kr8_cache.CreateClusterCache("test_cluster_config"),
				ComponentConfigs: map[string]kr8_cache.ComponentCache{
					"component1": {
						ComponentConfig: "config1",
						ComponentFiles: map[string]string{
							"file1.txt": "hash1",
						},
					},
					"component2": {
						ComponentConfig: "config2",
						ComponentFiles: map[string]string{
							"file2.txt": "hash2",
						},
					},
				},
				LibraryCache: &kr8_cache.LibraryCache{
					Directory: "",
					Entries:   map[string]string{},
				},
			},
			wantErr: false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotErr := kr8_cache.LoadClusterCache(testCase.cacheFile)
			// _ = testCase.want.WriteCache(testCase.cacheFile, false)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("LoadClusterCache() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("LoadClusterCache() succeeded unexpectedly")
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("LoadClusterCache() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestDeploymentCache_WriteCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cacheFile string
		// Named input parameters for target function.
		outFile string
		wantErr bool
	}{
		{
			name:      "Valid cache file with correct structure",
			cacheFile: "testdata/valid_cache.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   false,
		},
		{
			name:      "Invalid cache file (missing cluster_config)",
			cacheFile: "testdata/invalid_missing_cluster_config.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   true,
		},
		{
			name:      "Invalid cache file (malformed JSON)",
			cacheFile: "testdata/invalid_malformed_json.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   true,
		},
		{
			name:      "Valid cache file with empty component_configs",
			cacheFile: "testdata/valid_empty_component_configs.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   false,
		},
		{
			name:      "Valid cache file with multiple components",
			cacheFile: "testdata/valid_multiple_components.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   false,
		},
		{
			name:      "Valid cache file with correct structure",
			cacheFile: "testdata/valid_cache.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   false,
		},
		{
			name:      "Invalid cache file (missing cluster_config)",
			cacheFile: "testdata/invalid_missing_cluster_config.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   true,
		},
		{
			name:      "Invalid cache file (malformed JSON)",
			cacheFile: "testdata/invalid_malformed_json.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   true,
		},
		{
			name:      "Valid cache file with empty component_configs",
			cacheFile: "testdata/valid_empty_component_configs.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   false,
		},
		{
			name:      "Valid cache file with multiple components",
			cacheFile: "testdata/valid_multiple_components.json",
			outFile:   "testdata/valid_cache_out.json",
			wantErr:   false,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache, err := kr8_cache.LoadClusterCache(testCase.cacheFile)
			if err != nil {
				if !testCase.wantErr {
					t.Fatalf("could not construct receiver type: %v", err)
				} else {
					return
				}
			}
			gotErr2 := cache.WriteCache(testCase.outFile, false)
			if gotErr2 != nil {
				if !testCase.wantErr {
					t.Errorf("WriteCache() failed: %v", gotErr2)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("WriteCache() succeeded unexpectedly")
			}
		})
	}
}

func TestDeploymentCache_CheckClusterCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cacheFile string
		// Named input parameters for target function.
		config string
		logger zerolog.Logger
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache, err := kr8_cache.LoadClusterCache(testCase.cacheFile)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := cache.CheckClusterCache(testCase.config, testCase.logger)
			if got != testCase.want {
				t.Errorf("CheckClusterCache() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestDeploymentCache_CheckClusterComponentCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cacheFile string
		// Named input parameters for target function.
		config        string
		componentName string
		componentPath string
		files         []string
		logger        zerolog.Logger
		want          bool
		want2         *kr8_cache.ComponentCache
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache, err := kr8_cache.LoadClusterCache(testCase.cacheFile)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, got2, gotErr := cache.CheckClusterComponentCache(
				testCase.config,
				testCase.componentName,
				testCase.componentPath,
				testCase.files,
				testCase.logger,
			)
			// TODO: update the condition below to compare got with tt.want.
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("ProcessComponentFinalizer() failed: %v", gotErr)
				}

				return
			}
			if true {
				t.Errorf("CheckClusterComponentCache() = %v, want %v", got, testCase.want)
			}
			if true {
				t.Errorf("CheckClusterComponentCache() = %v, want %v", got2, testCase.want2)
			}
		})
	}
}

func TestCreateClusterCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		config string
		want   *kr8_cache.ClusterCache
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := kr8_cache.CreateClusterCache(tt.config)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CreateClusterCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClusterCache_CheckClusterCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cconfig string
		// Named input parameters for target function.
		config string
		logger zerolog.Logger
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache := kr8_cache.CreateClusterCache(testCase.cconfig)
			got := cache.CheckClusterCache(testCase.config, testCase.logger)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CheckClusterCache() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestCreateComponentCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		config    string
		listFiles []string
		want      *kr8_cache.ComponentCache
		wantErr   bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotErr := kr8_cache.CreateComponentCache(testCase.config, testCase.listFiles)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("CreateComponentCache() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("CreateComponentCache() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CreateComponentCache() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestComponentCache_CheckComponentCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for receiver constructor.
		cconfig    string
		clistFiles []string
		// Named input parameters for target function.
		config        string
		componentName string
		componentPath string
		baseDir       string
		files         []string
		logger        zerolog.Logger
		want          bool
		want2         *kr8_cache.ComponentCache
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache, err := kr8_cache.CreateComponentCache(testCase.cconfig, testCase.clistFiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got, got2 := cache.CheckComponentCache(
				testCase.config,
				testCase.componentName,
				testCase.baseDir,
				testCase.files,
				testCase.logger,
			)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CheckComponentCache() = %v, want %v", got, testCase.want)
			}
			if true {
				t.Errorf("CheckComponentCache() = %v, want %v", got2, testCase.want2)
			}
		})
	}
}
