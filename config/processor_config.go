package config

import (
	"fmt"
	"os"
	"path"

	"github.com/aeroxmotion/gexarch/util"
	"golang.org/x/mod/modfile"
)

type ProcessorConfig struct {
	*CliConfig
	ModulePath     string
	TypeName       string
	UseCaseName    string
	RepositoryName string
	EntityName     string
}

func GetProcessorConfigByType(Type string) *ProcessorConfig {
	modFileBytes, err := os.ReadFile("go.mod")
	util.PanicIfError(err)

	parsedModFile, err := modfile.Parse("go.mod", modFileBytes, func(_, version string) (string, error) {
		return version, nil
	})
	util.PanicIfError(err)

	cliConfig := GetCliConfig()

	return &ProcessorConfig{
		CliConfig:      cliConfig,
		ModulePath:     path.Join(parsedModFile.Module.Mod.Path, cliConfig.TargetPath),
		TypeName:       Type,
		UseCaseName:    fmt.Sprintf("%sFinder", Type),
		RepositoryName: fmt.Sprintf("%sRepository", Type),
		EntityName:     Type,
	}
}
