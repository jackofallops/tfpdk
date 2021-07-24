package commands

type InitialiseCommand struct{}

func (i InitialiseCommand) Run(args []string) int {
	// TODO - Clone the scaffold project to path from `args`?
	return 0
}

func (InitialiseCommand) Help() string {
	return ""
}

func (InitialiseCommand) Synopsis() string {
	return ""
}
