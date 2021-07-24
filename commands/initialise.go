package commands

import "fmt"

type InitialiseCommand struct{}

func (i InitialiseCommand) Run(args []string) int {
	// TODO - Clone the scaffold project to path from `args`?
	fmt.Printf("Not yet implemented - sorry!")
	return 0
}

func (InitialiseCommand) Help() string {
	return "Not yet implemented - sorry!"
}

func (InitialiseCommand) Synopsis() string {
	return "Not yet implemented - sorry!"
}
