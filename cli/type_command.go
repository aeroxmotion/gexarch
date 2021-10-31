package cli

import (
	"errors"
	"strings"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/processor"
	"github.com/iancoleman/strcase"
	"github.com/urfave/cli/v2"
)

func typeCommand() *cli.Command {
	return &cli.Command{
		Name:    "type",
		Aliases: []string{"t"},
		Usage:   "Generate scaffold by `type`",
		Action:  typeCommandAction,
	}
}

func typeCommandAction(ctx *cli.Context) error {
	targetType := strings.TrimSpace(strcase.ToCamel(ctx.Args().Get(0)))

	if targetType == "" {
		return errors.New("missing `type` name")
	}

	processor := processor.NewTemplateProcessor(config.GetProcessorConfigByType(targetType))
	processor.ProcessByType()

	return nil
}
