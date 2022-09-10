package app

import (
	"github.com/rezkit/cli/commands"
	"github.com/urfave/cli"
)

func GetApp() *cli.App {
	app := cli.NewApp()
	app.Name = "RezKit"
	app.Description = "RezKit CLI Tool"
	app.Authors = []cli.Author{
		{
			Name:  "RezKit",
			Email: "support@rezkit.app",
		},
	}

	app.Before = checkInstallation

	app.Commands = []cli.Command{
		{
			Name:     "login",
			Usage:    "login",
			Category: "config",
			Action:   commands.Login,
		},
	}

	return app
}

// Checks that the CLI app is properly installed.
func checkInstallation(cli *cli.Context) error {
	// Check to see if we're logged in. If not, tell the user we have to log in.

	// Don't check if we're trying to log in.
	if cli.Command.Name == "login" {
		return nil
	}

	return nil
}
