package commands

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
)

type InitialiseCommand struct {
	Ui cli.Ui
}

func (i InitialiseCommand) Run(args []string) int {
	// TODO - Clone the scaffold project to path from `args`?
	if len(args) == 0 {
		fmt.Print(i.Help())
		return 1
	}

	for _, v := range args {
		arg := strings.Split(v, "=")
		if len(arg) > 2 {
			fmt.Printf("malformed argument %q", arg)
			return 1
		}

	}

	return 0
}

func (InitialiseCommand) Help() string {
	return "Not yet implemented - sorry!"
}

func (InitialiseCommand) Synopsis() string {
	return "Not yet implemented - sorry!"
}
