package commands

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/jackofallops/tfpdk/helpers"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/cli"
)

type ResourceCommand struct {
	Ui cli.Ui
}

var _ cli.Command = ResourceCommand{}

var config = &helpers.Configuration{}

type ResourceData struct {
	HasUpdate      bool
	Name           string
	ProviderName   string
	ServicePackage string
	Typed          bool
}

func (d *ResourceData) ParseArgs(args []string) (errors []error) {
	resourceSet := flag.NewFlagSet("resource", flag.ExitOnError)
	resourceSet.StringVar(&d.Name, "name", "", "(Required) the name of the new resource, can be in the form resource_name, ResourceName, or resource-name")
	resourceSet.StringVar(&d.ServicePackage, "servicepackage", "", "(Optional) place the resource under the named service package")
	resourceSet.BoolVar(&d.HasUpdate, "has-update", true, "(Optional) Use if the new resource supports updating in place. Note if not used all schema properties must be 'ForceNew: true'")
	resourceSet.BoolVar(&d.Typed, "typed", config.TypedSDK, "(Optional) Generate a resource for use with the Typed Resource SDK")
	err := resourceSet.Parse(args)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	if d.Name == "" {
		errors = append(errors, fmt.Errorf("required option `-name` missing\n"))
	}

	return errors
}

func (c ResourceCommand) Run(args []string) int {
	data := &ResourceData{}
	if err := data.ParseArgs(args); err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}
		return 1
	}

	if err := data.generate(); err != nil {
		c.Ui.Error(fmt.Sprintf("generating resource %s: %+v", data.Name, err))
		return 1
	}

	return 0
}

func (d ResourceData) generate() error {
	providerName, err := helpers.ProviderName()
	if err != nil {
		return err
	}
	d.ProviderName = *providerName

	tpl := template.Must(template.New("resource.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/resource.gotpl"))

	outputPath := ""
	if d.ServicePackage != "" {
		path := fmt.Sprintf("%s/%s", config.ServicePackagesPath, d.ServicePackage)
		_, err := os.Stat(path)
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		}
		outputPath = fmt.Sprintf("%s/%s/%s_resource.go", config.ServicePackagesPath, strings.ToLower(strcase.ToCamel(d.ServicePackage)), strcase.ToSnake(d.Name))
	} else {
		outputPath = fmt.Sprintf("internal/%s_resource.go", strcase.ToSnake(d.Name))
	}

	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("a resource with this name already exists and will not be overwritten. Please remove this file if you wish to regenerate")
	}

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed opening output resource file for writing: %+v", err.Error())
	}

	err = tpl.Execute(f, d)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("failed writing to file: %+v", err.Error())
	}

	// Make sure the generated template complies with the local Go's view of how it should be fmt'd
	if err := helpers.GoFmt(outputPath); err != nil {
		return err
	}

	return nil
}

func (c ResourceCommand) Help() string {
	return `
Usage: tfpdk resource [options]

Generates a scaffolded resource, optionally under a service package for this provider

Example: 
$ tfpdk resource -name=MyNewResource -servicepackage=SomeExistingService -typed -has-update
`
}

func (c ResourceCommand) Synopsis() string {
	return "creates boiler-plate resources."
}
