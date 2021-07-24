package commands

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/cli"
)

type DataSourceCommand struct{}

var _ cli.Command = DataSourceCommand{}

type TypedDataSourceData struct {
	Name             string
	ProviderName     string
	ServicePackage   string
	Typed            bool
	UseResourceModel bool
}

func (d DataSourceCommand) Run(args []string) int {
	data := TypedDataSourceData{
		ProviderName:     "azurerm",
		Typed:            false,
		UseResourceModel: false,
	}

	if len(args) == 0 {
		fmt.Print(d.Help())
		return 1
	}

	for _, v := range args {
		arg := strings.Split(v, "=")
		if len(arg) > 2 {
			fmt.Printf("malformed argument %q", arg)
			return 1
		}

		switch strings.ToLower(strings.TrimLeft(arg[0], "-")) {
		case "name":
			if len(arg) == 2 {
				data.Name = arg[1]
			} else {
				fmt.Println("argument `name` requires a value, eg `-name=some_resource_name`")
				return 1
			}

		case "servicepackage":
			if len(arg) == 2 {
				data.ServicePackage = arg[1]
			} else {
				fmt.Println("argument `servicepackage` requires a value, eg `-servicepackage=some_service_name`")
				return 1
			}

		case "typed":
			data.Typed = true

		case "useresourcemodel":
			data.UseResourceModel = true

		default:
			fmt.Printf("unrecognised option %q", arg[0])
			return 1
		}
	}

	if data.Name == "" {
		fmt.Printf("Error: missing required argument `-name`")
		return 1
	}

	tpl := template.Must(template.New("datasource.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/datasource.gotpl"))

	outputPath := ""
	if data.ServicePackage != "" {
		outputPath = fmt.Sprintf("%s/internal/services/%s/%s_data_source.go", strings.ToLower(data.ProviderName), strings.ToLower(strcase.ToCamel(data.ServicePackage)), strcase.ToSnake(data.Name))
	} else {
		outputPath = fmt.Sprintf("%s/internal/%s_data_source.go", data.ProviderName, strcase.ToSnake(data.Name))
	}

	//if _, err := os.Stat(outputPath); err == nil {
	//	fmt.Printf("Error: A data source with this name already exists and will not be overwritten. Please remove this file if you wish to regenerate.")
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

func (d DataSourceCommand) Help() string {
	return fmt.Sprintf(`
Usage: tfpdk datasource [options]

Generates a scaffolded datasource, optionally under a service package for this provider

Options:

-name=string			(Required) the name of the new Resource, can be in the form resource_name, ResourceName, or resource-name

-servicepackage=string		(Optional) place the Data Source under the named service package

-typed				(Optional) Generate a Data Source for use with the Typed Resource SDK

-useresourcemodel		(Optional) Use the related resouce model for this. Use this if there is an Existing Resource for this Data Source that already contains a suitable model for representing the Schema.

Example:
 tfpdk datasource -name=MyNewResource -servicepackage=MyExistingService -typed -useresourcemodel

`)
}

func (d DataSourceCommand) Synopsis() string {
	return ""
}
