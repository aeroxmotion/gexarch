package cli

import (
	"strings"

	"github.com/aeroxmotion/gexarch/config"
	"github.com/aeroxmotion/gexarch/processor"
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
	targetType := strings.Title(ctx.Args().Get(0))

	processor := processor.NewTemplateProcessor(config.GetProcessorConfigByType(targetType))
	processor.Process()

	return nil
}
