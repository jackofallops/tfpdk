package main

import (
	"log"
	"os"

	"github.com/jackofallops/tfpdk/commands"
	"github.com/mitchellh/cli"
)

var Commands map[string]cli.CommandFactory

func main() {
	args := os.Args[1:]
	Commands = map[string]cli.CommandFactory{
		"init": func() (cli.Command, error) {
			return &commands.InitialiseCommand{}, nil
		},
		"resource": func() (cli.Command, error) {
			return &commands.ResourceCommand{}, nil
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

	os.Exit(exitStatus)
}
