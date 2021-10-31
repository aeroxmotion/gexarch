package config

import (
	"fmt"
	"path"

	"github.com/aeroxmotion/gexarch/util"
	"github.com/iancoleman/strcase"
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
	modFile := util.ParseModfile()
	cliConfig := GetCliConfig()

	return &ProcessorConfig{
		CliConfig:      cliConfig,
		ModulePath:     path.Join(modFile.Module.Mod.Path, cliConfig.TargetPath),
		TypeName:       Type,
		UseCaseName:    fmt.Sprintf("%sFinder", Type),
		RepositoryName: fmt.Sprintf("%sRepository", Type),
		EntityName:     Type,
	}
}

func (config *ProcessorConfig) ToSnakeValues() *ProcessorConfig {
	return &ProcessorConfig{
		CliConfig:      config.CliConfig,
		ModulePath:     config.ModulePath,
		TypeName:       strcase.ToSnake(config.TypeName),
		EntityName:     strcase.ToSnake(config.EntityName),
		RepositoryName: strcase.ToSnake(config.RepositoryName),
		UseCaseName:    strcase.ToSnake(config.UseCaseName),
	}
}
