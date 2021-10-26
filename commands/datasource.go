package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/jackofallops/tfpdk/helpers"

	"github.com/iancoleman/strcase"
	"github.com/mitchellh/cli"
)

type DataSourceCommand struct {
	Ui cli.Ui
}

var _ cli.Command = DataSourceCommand{}

type DataSourceData struct {
	Name           string
	ProviderName   string
	ServicePackage string
	Typed          bool
	Config         *helpers.Configuration
}

func (d *DataSourceData) ParseArgs(args []string) (errors []error) {
	dsSet := flag.NewFlagSet("datasource", flag.ExitOnError)
	dsSet.StringVar(&d.Name, "name", "", "(Required) the name of the new Resource, can be in the form resource_name, ResourceName, or resource-name")
	dsSet.StringVar(&d.ServicePackage, "service-package", "", "(Optional) place the Data Source under the named service package")
	dsSet.BoolVar(&d.Typed, "typed", config.TypedSDK, "(Optional) Generate a Data Source for use with the Typed Resource SDK")
	err := dsSet.Parse(args)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	if d.Name == "" {
		errors = append(errors, fmt.Errorf("required option `-name` missing\n"))
	}

	return errors
}

func (c DataSourceCommand) Run(args []string) int {
	data := &DataSourceData{
		Config: config,
	}

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

func (d DataSourceData) generate() error {
	if d.ProviderName == "" {
		providerName, err := helpers.ProviderName()
		if err != nil {
			return err
		}
		d.ProviderName = *providerName
	}
	tpl := template.Must(template.New("datasource.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/datasource.gotpl"))

	outputPath := ""
	if d.ServicePackage != "" {
		outputPath = fmt.Sprintf("%s/%s/%s_data_source.go", config.ServicePackagesPath, strings.ToLower(strcase.ToCamel(d.ServicePackage)), strcase.ToSnake(d.Name))
	} else {
		outputPath = fmt.Sprintf("internal/%s_data_source.go", strcase.ToSnake(d.Name))
	}

	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("a data source with this name already exists and will not be overwritten. Please remove this file if you wish to regenerate it")
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

	if err := helpers.GoFmt(outputPath); err != nil {
		return err
	}

	_ = helpers.UpdateRegistration(config.ServicePackagesPath+"/"+d.ServicePackage, d.Name, helpers.TypeDatasource, helpers.Register, d.Typed)

	return nil
}

func (c DataSourceCommand) Help() string {
	return `
Usage: tfpdk datasource [options]

Example:
$ tfpdk datasource -name=MyNewResource -service-package=MyExistingService -typed -use-resource-model

`
}

func (c DataSourceCommand) Synopsis() string {
	return "creates boiler-plate Data Sources."
}
