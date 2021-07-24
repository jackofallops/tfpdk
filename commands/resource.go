package commands

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/cli"
)

type ResourceCommand struct {
}

var _ cli.Command = ResourceCommand{}

type TypedResourceData struct {
	Name           string
	ServicePackage string
	HasUpdate      bool
	ProviderName   string
	Typed          bool
}

func (r ResourceCommand) Run(args []string) int {
	data := TypedResourceData{
		ProviderName: "azurerm", // Defaulting for now, will need to eval from os.path later
	}

	if len(args) == 0 {
		fmt.Print(r.Help())
		return 1
	}

	for _, v := range args {
		arg := strings.Split(v, "=")
		if len(arg) > 2 {
			fmt.Printf("malformed argument %q", arg)
			return 1
		}

		switch strings.ToLower(strings.TrimLeft(arg[0], "-")) {
		case "help":
			fmt.Printf(r.Help())
			return 0
		case "name":
			if len(arg) == 2 {
				data.Name = arg[1]
			} else {
				fmt.Println("argument `name` requires a value, eg `-name=some_resource_name`")
			}

		case "has-update":
			data.HasUpdate = true

		case "servicepackage":
			data.ServicePackage = arg[1]

		case "typed":
			data.Typed = true

		default:
			fmt.Printf("unrecognised option %q", arg[0])
			return 1
		}
	}

	if data.Name == "" {
		fmt.Printf("Error: missing required argument `-name`")
		return 1
	}

	tpl := template.Must(template.New("resource.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/resource.gotpl"))

	// output file - Note we assume that the PATH to the file already exists, having been through init and, optionally, service package creation
	outputPath := ""
	if data.ServicePackage != "" {
		outputPath = fmt.Sprintf("%s/internal/services/%s/%s_resource.go", strings.ToLower(data.ProviderName), strings.ToLower(strcase.ToCamel(data.ServicePackage)), strcase.ToSnake(data.Name))
	} else {
		outputPath = fmt.Sprintf("%s/internal/%s_resource.go", data.ProviderName, strcase.ToSnake(data.Name))
	}

	//if _, err := os.Stat(outputPath); err == nil {
	//	fmt.Printf("Error: A resource with this name already exists and will not be overwritten. Please remove this file if you wish to regenerate.")
	//	return 1
	//}

	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error: failed opening output resource file for writing: %+v", err.Error())
		return 1
	}

	err = tpl.Execute(f, data)
	if err != nil {
		fmt.Println(err.Error())
		return 1
	}
	if err := f.Close(); err != nil {
		fmt.Printf("Error: Writing to file: %+v", err.Error())
	}

	return 0
}

func (r ResourceCommand) Help() string {
	return fmt.Sprintf(`
Usage: tfpdk resource [options]

Generates a scaffolded resource, optionally under a service package for this provider

Options:

-name=string			(Required) the name of the new resource, can be in the form resource_name, ResourceName, or resource-name

-servicepackage=string		(Optional) place the resource under the named service package

-has-update			(Optional) Use if the new resouce supports updating in place. Note if not used all schema properties must be 'ForceNew: true'

-typed				(Optional) Generate a resource for use with the Typed Resource SDK

`)
}

func (r ResourceCommand) Synopsis() string {
	return "creates boiler-plate resources."
}
