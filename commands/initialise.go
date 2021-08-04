package commands

import (
	"flag"
	"fmt"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/mitchellh/cli"
)

type InitialiseCommand struct {
	Ui cli.Ui
}

type InitData struct {
	ProviderName    string
	ServicePackages bool
	TypedSDK        bool
}

func (d *InitData) parseArgs(args []string) (errors []error) {
	initSet := flag.NewFlagSet("init", flag.ExitOnError)
	initSet.StringVar(&d.ProviderName, "name", "", "(Required) name for the new provider")
	initSet.BoolVar(&d.ServicePackages, "service-packages", false, "(Optional) This provider groups resources and data sources into service packages")
	initSet.BoolVar(&d.TypedSDK, "use-typed-sdk", false, "(Optional) This provider uses the Typed SDK for resources and Data Sources")

	if err := initSet.Parse(args); err != nil {
		errors = append(errors, err)
		return errors
	}

	if d.ProviderName == "" {
		errors = append(errors, fmt.Errorf("required option `-name` missing\n"))
		return errors
	}
	return errors
}

func (i InitialiseCommand) Run(args []string) int {
	data := &InitData{}

	if err := data.parseArgs(args); err != nil {
		for _, e := range err {
			i.Ui.Error(e.Error())
		}
		return 1
	}

	if err := data.exec(); err != nil {
		i.Ui.Error(err.Error())
		return 1
	}

	return 0
}

func (d InitData) exec() error {
	cloneOpts := &git.CloneOptions{
		URL:   "https://github.com/hashicorp/terraform-provider-scaffolding.git",
		Depth: 1,
		Tags:  git.NoTags,
	}
	providerName := fmt.Sprintf("terraform-provider-%s", strings.ToLower(strings.TrimPrefix(d.ProviderName, "terraform-provider-")))

	// Clone scaffold repo
	r, err := git.PlainClone(providerName, false, cloneOpts)
	if err != nil {
		return fmt.Errorf("cloning: %+v", err)
	}

	r.DeleteRemote("origin")

	return nil
}

func (InitialiseCommand) Help() string {
	return `
Usage: tfpdk init -name=[providername] [-typed] [-service-packages]

Example:
$ tfpdk init -name=myprovider -typed -service-packages
`
}

func (InitialiseCommand) Synopsis() string {
	return "initialises a new provider based on the scaffold project"
}
