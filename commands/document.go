package commands

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/jackofallops/tfpdk/helpers"
	"github.com/mitchellh/cli"
)

type DocumentCommand struct {
	Ui cli.Ui
}

var _ cli.Command = DocumentCommand{}

type DocumentData struct {
	Name                  string
	SnakeName             string
	ProviderName          string
	ProviderCanonicalName string
	ServicePackage        string
	DocType               string
	Resource              helpers.Resource
	IDExample             string
	Examples              []string
	ResourceData          string
}

func (d *DocumentData) ParseArgs(args []string) (errors []error) {
	docSet := flag.NewFlagSet("document", flag.ExitOnError)
	docSet.StringVar(&d.Name, "name", "", "The name of the resource")
	docSet.StringVar(&d.ServicePackage, "servicepackage", "", "The name of the Service Package the resource or data source belongs to")
	docSet.StringVar(&d.DocType, "type", "", "The type of item to document, one of `resource` or `datasource`")
	docSet.StringVar(&d.IDExample, "id", "", "An example of the ID this resource has when created (only valid for `-type=resource`)")
	err := docSet.Parse(args)
	if err != nil {
		errors = append(errors, err)
		return errors
	}

	if d.Name == "" {
		errors = append(errors, fmt.Errorf("required option `-name` missing\n"))
	}

	if strings.EqualFold(d.DocType, "resource") && d.IDExample == "" {
		errors = append(errors, fmt.Errorf("`-id` required when `-type=resource\n`"))
	}

	if config.ProviderCanonicalName != "" {
		d.ProviderCanonicalName = config.ProviderCanonicalName
	} else {
		d.ProviderCanonicalName = d.Name
	}

	return errors
}

func (c DocumentCommand) Run(args []string) int {
	data := &DocumentData{}

	if len(args) == 0 {
		fmt.Print(c.Help())
		return 1
	}

	if err := data.ParseArgs(args); err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}
		return 1
	}

	if err := data.generate(); err != nil {
		c.Ui.Error(fmt.Sprintf("Failed to generate: %+v", err))
		return 1
	}

	return 0
}

func (c DocumentCommand) Help() string {
	return "TODO"
}

func (c DocumentCommand) Synopsis() string {
	return "generates documentation from a resource"
}

func (d *DocumentData) generate() error {
	if d.ProviderName == "" {
		providerName, err := helpers.ProviderName()
		if err != nil {
			return err
		}
		d.ProviderName = *providerName
	}
	if d.SnakeName == "" {
		d.SnakeName = strcase.ToSnake(fmt.Sprintf("%s_%s", d.ProviderName, d.Name))
	}

	resource, err := helpers.GetSchema(config.SchemaAPIURL, d.DocType, d.SnakeName)
	if err != nil {
		return err
	}
	d.Resource = *resource

	tpl := template.Must(template.New("document.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/document.gotpl"))

	outputPath := fmt.Sprintf("website/%s", config.DocsPath) // TODO - AzureRM may drop the `website` in 3.0
	if d.DocType == "datasource" {
		outputPath = fmt.Sprintf("%s/%s", outputPath, config.DataSourceDocsDirname)
	} else {
		outputPath = fmt.Sprintf("%s/%s", outputPath, config.ResourceDocsDirname)
	}

	providerPrefix := fmt.Sprintf("%s_", d.ProviderName)
	outputPath = fmt.Sprintf("%s/%s.md", outputPath, strings.TrimPrefix(d.SnakeName, providerPrefix))

	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("[DEBUG] failed opening output resource file for writing: %+v", err.Error())
		return err
	}

	defer f.Close()

	err = tpl.Execute(f, d)
	if err != nil {
		return err
	}
	return nil
}

func buildAPIURL(base, resourceType string) string {
	return ""
}
