package cli

import (
	"os"

	"github.com/urfave/cli/v2"
)

func Start() {
	app := &cli.App{
		Name:  "gexarch",
		Usage: "Generate on-demand scaffolding following the ports & adapters architecture",
		Commands: []*cli.Command{
			typeCommand(),
		},
	}

	app.Run(os.Args)
}
