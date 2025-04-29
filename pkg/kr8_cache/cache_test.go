package kr8_cache_test

import (
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
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotErr := kr8_cache.LoadClusterCache(testCase.cacheFile)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("LoadClusterCache() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("LoadClusterCache() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
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
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache, err := kr8_cache.LoadClusterCache(testCase.cacheFile)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			gotErr := cache.WriteCache(testCase.outFile)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("WriteCache() failed: %v", gotErr)
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
			// TODO: update the condition below to compare got with tt.want.
			if true {
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
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache, err := kr8_cache.LoadClusterCache(testCase.cacheFile)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := cache.CheckClusterComponentCache(
				testCase.config,
				testCase.componentName,
				testCase.componentPath,
				testCase.files,
				testCase.logger,
			)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CheckClusterComponentCache() = %v, want %v", got, testCase.want)
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
		config        string
		componentPath string
		listFiles     []string
		want          *kr8_cache.ComponentCache
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotErr := kr8_cache.CreateComponentCache(testCase.config, testCase.componentPath, testCase.listFiles)
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
		cconfig        string
		ccomponentPath string
		clistFiles     []string
		// Named input parameters for target function.
		config        string
		componentName string
		componentPath string
		files         []string
		logger        zerolog.Logger
		want          bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cache, err := kr8_cache.CreateComponentCache(testCase.cconfig, testCase.ccomponentPath, testCase.clistFiles)
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}
			got := cache.CheckComponentCache(
				testCase.config,
				testCase.componentName,
				testCase.componentPath,
				testCase.files,
				testCase.logger,
			)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CheckComponentCache() = %v, want %v", got, testCase.want)
			}
		})
	}
}
