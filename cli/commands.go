package cli

import "github.com/codegangsta/cli"

var commands = []cli.Command{
	{
		Name:   "deploy",
		Usage:  "despliega el servicio",
		Flags:  deployFlags(),
		Before: deployBefore,
		Action: deployCmd,
	},
}
