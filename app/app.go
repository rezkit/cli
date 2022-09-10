package app

import (
	"errors"
	"fmt"
	"os"

	"github.com/rezkit/cli/commands"
	"github.com/rezkit/cli/internal/config"
	"github.com/spf13/viper"
	"github.com/urfave/cli"
)

var (
	ErrInvalidFormat = errors.New("Invalid output format")
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

	app.Before = func(ctx *cli.Context) error {
		if err := loadConfig(ctx); err != nil {
			return err
		}

		if err := checkParams(ctx); err != nil {
			return err
		}

		return checkInstallation(ctx)
	}

	app.After = func(ctx *cli.Context) error {
		return config.GetConfig().WriteConfig()
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:     "format",
			Value:    "table",
			Required: false,
			Usage:    "Set the output format: table|json|yaml",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:     "login",
			Usage:    "login",
			Category: "config",
			Action:   commands.Login,
		},

		{
			Name:      "organizations",
			ShortName: "org",
			Subcommands: []cli.Command{
				{
					Name:      "list",
					ShortName: "ls",
					Action:    commands.ListOrganizations,
					Usage:     "List available organizations",
				},
			},
		},
	}

	return app
}

func loadConfig(ctx *cli.Context) error {
	if err := config.GetConfig().ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			config.GetConfig().WriteConfig()
			return nil
		default:
			return err
		}
	}

	return nil
}

// Checks that the CLI app is properly installed.
func checkInstallation(cli *cli.Context) error {

	// Don't check if we're trying to log in.
	if cli.Command.Name == "login" {
		return nil
	}

	if token := config.GetConfig().GetString("authentication.access_token"); token == "" {
		fmt.Fprintln(os.Stderr, "Please run `login` first to log in.")
		return errors.New("Please run `login` first")
	}

	return nil
}

func checkParams(ctx *cli.Context) error {
	format := ctx.GlobalString("format")

	if format == "table" || format == "yaml" || format == "json" {
		return nil
	}

	return ErrInvalidFormat
}
