package main

import (
	"os"
	"github.com/codegangsta/cli"
)


func main() {
	app := cli.NewApp()

	app.Name = "gnvm"
	app.Version = "0.1.0"
	app.Usage = "install and switch node.js versions in a snap"

	app.Commands = []cli.Command{
		{
			Name: "list",
			Aliases: []string{"ls" },
			Usage: "list all installed versions of node.js",
			Action: ListAction,
		},
		{
			Name: "list-remote",
			Aliases: []string{"ls-remote" },
			Usage: "list all existing versions of node.js",
			Action: ListRemoteAction,
		},
		{
			Name: "current",
			Usage: "show current version of node.js",
			Action: CurrentAction,
		},
		{
			Name: "load",
			Usage: "install particular version of node.js",
			Action: LoadAction,
		},
		{
			Name: "use",
			Usage: "start using particular version of node.js",
			Action: UseAction,
		},
		{
			Name: "run",
			Usage: "run the given version of node.js",
			Action: RunAction,
		},
		{
			Name: "uninstall",
			Usage: "uninstall particular version of node.js",
			Action: UninstallAction,
		},
	}

	app.Run(os.Args)
}
