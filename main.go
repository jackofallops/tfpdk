package main

import (
	"log"
	"os"

	"github.com/jackofallops/tfpdk/commands"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func main() {
	os.Exit(realMain(os.Args[1:]))
}

func realMain(args []string) int {
	var ui cli.Ui = &cli.ColoredUi{
		ErrorColor: cli.UiColorRed,
		WarnColor:  cli.UiColorYellow,
		InfoColor:  cli.UiColorNone,

		Ui: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
	}

	Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &commands.InitialiseCommand{
				Ui: ui,
			}, nil
		},
		"resource": func() (cli.Command, error) {
			return &commands.ResourceCommand{
				Ui: ui,
			}, nil
		},
		"datasource": func() (cli.Command, error) {
			return &commands.DataSourceCommand{
				Ui: ui,
			}, nil
		},
		"document": func() (cli.Command, error) {
			return &commands.DocumentCommand{
				Ui: ui,
			}, nil
		},
		"config": func() (cli.Command, error) {
			return &commands.GenConfigCommand{
				Ui: ui,
			}, nil
		},
		"servicepackage": func() (cli.Command, error) {
			return &commands.ServicePackageCommand{
				Ui: ui,
			}, nil
		},
	}

	tfpdk := cli.CLI{
		Args:     args,
		Commands: Commands,
		Name:     "tfpdk",
		Version:  "0.1.0",
	}

	exitStatus, err := tfpdk.Run()
	if err != nil {
		log.Println(err)
	}
	return exitStatus
}
