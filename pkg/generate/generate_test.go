package generate_test

import (
	"reflect"
	"testing"

	jsonnet "github.com/google/go-jsonnet"
	"github.com/ice-bergtech/kr8/pkg/generate"
	"github.com/ice-bergtech/kr8/pkg/kr8_cache"
	"github.com/ice-bergtech/kr8/pkg/kr8_types"
	types "github.com/ice-bergtech/kr8/pkg/types"
	util "github.com/ice-bergtech/kr8/pkg/util"
	"github.com/panjf2000/ants/v2"
	"github.com/rs/zerolog"
	"github.com/tidwall/gjson"
)

func TestGetClusterParams(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clusterDir string
		vmConfig   types.VMConfig
		logger     zerolog.Logger
		want       map[string]string
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotErr := generate.GetClusterParams(testCase.clusterDir, testCase.vmConfig, testCase.logger)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("GetClusterParams() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("GetClusterParams() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GetClusterParams() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestCalculateClusterComponentList(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clusterComponents         map[string]gjson.Result
		filters                   util.PathFilterOptions
		existingClusterComponents []string
		want                      []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generate.CalculateClusterComponentList(tt.clusterComponents, tt.filters, tt.existingClusterComponents)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CalculateClusterComponentList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenProcessComponent(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		vmConfig      types.VMConfig
		componentName string
		kr8Spec       kr8_types.Kr8ClusterSpec
		kr8Opts       types.Kr8Opts
		config        string
		allConfig     *generate.SafeString
		filters       util.PathFilterOptions
		paramsFile    string
		cache         *kr8_cache.DeploymentCache
		logger        zerolog.Logger
		want          bool
		want2         *kr8_cache.ComponentCache
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, got2, gotErr := generate.GenProcessComponent(
				testCase.vmConfig,
				testCase.componentName,
				testCase.kr8Spec,
				testCase.kr8Opts,
				testCase.config,
				testCase.allConfig,
				testCase.filters,
				testCase.paramsFile,
				testCase.cache,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("GenProcessComponent() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("GenProcessComponent() succeeded unexpectedly")
			}
			if got != testCase.want {
				t.Errorf("GenProcessComponent() = %v, want %v", got, testCase.want)
			}
			if !reflect.DeepEqual(got2, testCase.want2) {
				t.Errorf("GenProcessComponent() = %v, want %v", got2, testCase.want2)
			}
		})
	}
}

func TestProcessComponentFinalizer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		kr8Opts            types.Kr8Opts
		config             string
		compPath           string
		compSpec           kr8_types.Kr8ComponentSpec
		componentOutputDir string
		outputFileMap      map[string]bool
		logger             zerolog.Logger
		want               *kr8_cache.ComponentCache
		wantErr            bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			gotErr := generate.ProcessComponentFinalizer(
				testCase.compSpec,
				testCase.componentOutputDir,
				testCase.outputFileMap,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("ProcessComponentFinalizer() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("ProcessComponentFinalizer() succeeded unexpectedly")
			}
		})
	}
}

func TestCheckComponentCache(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		cache         *kr8_cache.DeploymentCache
		compSpec      kr8_types.Kr8ComponentSpec
		config        string
		componentName string
		baseDir       string
		logger        zerolog.Logger
		want          bool
		want2         *kr8_cache.ComponentCache
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, got2, gotErr := generate.CheckComponentCache(
				testCase.cache,
				testCase.compSpec,
				testCase.config,
				testCase.componentName,
				testCase.baseDir,
				testCase.logger,
			)
			// TODO: update the condition below to compare got with tt.want.
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("ProcessComponentFinalizer() failed: %v", gotErr)
				}

				return
			}
			if got != testCase.want {
				t.Errorf("CheckComponentCache() = %v, want %v", got, testCase.want)
			}
			if !reflect.DeepEqual(got2, testCase.want2) {
				t.Errorf("CheckComponentCache() = %v, want %v", got2, testCase.want2)
			}
		})
	}
}

func TestGetComponentFiles(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		compSpec kr8_types.Kr8ComponentSpec
		want     []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generate.GetComponentFiles(tt.compSpec)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetComponentFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetupComponentVM(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		vmConfig      types.VMConfig
		config        string
		kr8Spec       kr8_types.Kr8ClusterSpec
		componentName string
		compSpec      kr8_types.Kr8ComponentSpec
		allConfig     *generate.SafeString
		filters       util.PathFilterOptions
		paramsFile    string
		kr8Opts       types.Kr8Opts
		logger        zerolog.Logger
		want          *jsonnet.VM
		want2         string
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, got2, gotErr := generate.SetupComponentVM(
				testCase.vmConfig,
				testCase.config,
				testCase.kr8Spec,
				testCase.componentName,
				testCase.compSpec,
				testCase.allConfig,
				testCase.filters,
				testCase.paramsFile,
				testCase.kr8Opts,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("SetupComponentVM() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("SetupComponentVM() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("SetupComponentVM() = %v, want %v", got, testCase.want)
			}
			if got2 != testCase.want2 {
				t.Errorf("SetupComponentVM() = %v, want %v", got2, testCase.want2)
			}
		})
	}
}

func TestGetAllClusterParams(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clusterDir string
		vmConfig   types.VMConfig
		jvm        *jsonnet.VM
		logger     zerolog.Logger
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			gotErr := generate.GetAllClusterParams(testCase.clusterDir, testCase.vmConfig, testCase.jvm, testCase.logger)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("GetAllClusterParams() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("GetAllClusterParams() succeeded unexpectedly")
			}
		})
	}
}

func TestGetClusterComponentParamsThreadsafe(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		allConfig  *generate.SafeString
		config     string
		vmConfig   types.VMConfig
		kr8Spec    kr8_types.Kr8ClusterSpec
		filters    util.PathFilterOptions
		paramsFile string
		jvm        *jsonnet.VM
		logger     zerolog.Logger
		wantErr    bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			gotErr := generate.GetClusterComponentParamsThreadsafe(
				testCase.allConfig,
				testCase.config,
				testCase.vmConfig,
				testCase.kr8Spec,
				testCase.filters,
				testCase.paramsFile,
				testCase.jvm,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("GetClusterComponentParamsThreadsafe() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("GetClusterComponentParamsThreadsafe() succeeded unexpectedly")
			}
		})
	}
}

func TestGenerateIncludesFiles(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		includesFiles      []kr8_types.Kr8ComponentSpecIncludeObject
		kr8Spec            kr8_types.Kr8ClusterSpec
		kr8Opts            types.Kr8Opts
		config             string
		componentName      string
		compPath           string
		componentOutputDir string
		jvm                *jsonnet.VM
		logger             zerolog.Logger
		want               map[string]bool
		wantErr            bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotErr := generate.GenerateIncludesFiles(
				testCase.includesFiles,
				testCase.kr8Spec,
				testCase.kr8Opts,
				testCase.config,
				testCase.componentName,
				testCase.compPath,
				testCase.componentOutputDir,
				testCase.jvm,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("GenerateIncludesFiles() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("GenerateIncludesFiles() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GenerateIncludesFiles() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestGenProcessCluster(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clusterName         string
		clusterdir          string
		baseDir             string
		generateDirOverride string
		kr8Opts             types.Kr8Opts
		clusterParamsFile   string
		filters             util.PathFilterOptions
		vmConfig            types.VMConfig
		pool                *ants.Pool
		logger              zerolog.Logger
		wantErr             bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			gotErr := generate.GenProcessCluster(
				testCase.clusterName,
				testCase.clusterdir,
				testCase.baseDir,
				testCase.generateDirOverride,
				testCase.kr8Opts,
				testCase.clusterParamsFile,
				testCase.filters,
				testCase.vmConfig,
				testCase.pool,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("GenProcessCluster() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("GenProcessCluster() succeeded unexpectedly")
			}
		})
	}
}

func TestGatherClusterConfig(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clusterName         string
		clusterDir          string
		kr8Opts             types.Kr8Opts
		vmConfig            types.VMConfig
		generateDirOverride string
		filters             util.PathFilterOptions
		clusterParamsFile   string
		logger              zerolog.Logger
		want                *kr8_types.Kr8ClusterSpec
		want2               []string
		want3               string
		wantErr             bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, got2, got3, gotErr := generate.GatherClusterConfig(
				testCase.clusterName,
				testCase.clusterDir,
				testCase.kr8Opts,
				testCase.vmConfig,
				testCase.generateDirOverride,
				testCase.filters,
				testCase.clusterParamsFile,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("GatherClusterConfig() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("GatherClusterConfig() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GatherClusterConfig() = %v, want %v", got, testCase.want)
			}
			if true {
				t.Errorf("GatherClusterConfig() = %v, want %v", got2, testCase.want2)
			}
			if true {
				t.Errorf("GatherClusterConfig() = %v, want %v", got3, testCase.want3)
			}
		})
	}
}

func TestGenerateCacheInitializer(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		kr8Spec     *kr8_types.Kr8ClusterSpec
		enableCache bool
		logger      zerolog.Logger
		want        *kr8_cache.DeploymentCache
		want2       string
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, got2 := generate.LoadClusterCache(testCase.kr8Spec, testCase.logger)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("GenerateCacheInitializer() = %v, want %v", got, testCase.want)
			}
			if true {
				t.Errorf("GenerateCacheInitializer() = %v, want %v", got2, testCase.want2)
			}
		})
	}
}

func TestCleanupOldComponentDirs(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		existingComponents []string
		clusterComponents  map[string]gjson.Result
		kr8Spec            *kr8_types.Kr8ClusterSpec
		logger             zerolog.Logger
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			generate.CleanupOldComponentDirs(
				testCase.existingComponents,
				testCase.clusterComponents,
				testCase.kr8Spec,
				testCase.logger,
			)
		})
	}
}

func TestCompileClusterConfiguration(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		clusterName         string
		clusterDir          string
		kr8Opts             types.Kr8Opts
		vmConfig            types.VMConfig
		generateDirOverride string
		logger              zerolog.Logger
		want                *kr8_types.Kr8ClusterSpec
		want2               map[string]gjson.Result
		wantErr             bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, got2, gotErr := generate.CompileClusterConfiguration(
				testCase.clusterName,
				testCase.clusterDir,
				testCase.kr8Opts,
				testCase.vmConfig,
				testCase.generateDirOverride,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("CompileClusterConfiguration() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("CompileClusterConfiguration() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("CompileClusterConfiguration() = %v, want %v", got, testCase.want)
			}
			if true {
				t.Errorf("CompileClusterConfiguration() = %v, want %v", got2, testCase.want2)
			}
		})
	}
}

func TestRenderComponents(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		config            string
		vmConfig          types.VMConfig
		kr8Spec           kr8_types.Kr8ClusterSpec
		cache             *kr8_cache.DeploymentCache
		compList          []string
		clusterParamsFile string
		pool              *ants.Pool
		kr8Opts           types.Kr8Opts
		filters           util.PathFilterOptions
		logger            zerolog.Logger
		want              map[string]kr8_cache.ComponentCache
		wantErr           bool
	}{
		// TODO: Add test cases.
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			got, gotErr := generate.RenderComponents(
				testCase.config,
				testCase.vmConfig,
				testCase.kr8Spec,
				testCase.cache,
				testCase.compList,
				testCase.clusterParamsFile,
				testCase.pool,
				testCase.kr8Opts,
				testCase.filters,
				testCase.logger,
			)
			if gotErr != nil {
				if !testCase.wantErr {
					t.Errorf("RenderComponents() failed: %v", gotErr)
				}

				return
			}
			if testCase.wantErr {
				t.Fatal("RenderComponents() succeeded unexpectedly")
			}
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("RenderComponents() = %v, want %v", got, testCase.want)
			}
		})
	}
}
