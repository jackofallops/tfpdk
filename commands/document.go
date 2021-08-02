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
	Name           string
	SnakeName      string
	ProviderName   string
	ServicePackage string
	DocType        string
	Schema         helpers.ResourceSchema
	IDExample      string
	Examples       []string
	ResourceData   string
}

func (d *DocumentData) ParseArgs(args []string) (errors []error) {
	docSet := flag.NewFlagSet("document", flag.ExitOnError)
	docSet.StringVar(&d.Name, "name", "", "The name of the resource")
	docSet.StringVar(&d.ServicePackage, "servicepackage", "", "The name of the Service Package the resource or data source belongs to")
	docSet.StringVar(&d.DocType, "type", "", "The type of item to document, one of `resource` or `datasource`")
	docSet.StringVar(&d.IDExample, "id", "", "An example of the ID this resource has when created (only valid for `-type=resource`)")
	docSet.Parse(args)

	if d.Name == "" {
		errors = append(errors, fmt.Errorf("required option `-name` missing\n"))
	}

	if strings.EqualFold(d.DocType, "resource") && d.IDExample == "" {
		errors = append(errors, fmt.Errorf("`-id` required when `-type=resource\n`"))
	}

	return errors
}

func (c DocumentCommand) Run(args []string) int {
	data := &DocumentData{}

	if len(args) == 0 {
		fmt.Print(c.Help())
		return 1
	}

	var err []error
	err = data.ParseArgs(args)
	if err != nil {
		for _, e := range err {
			c.Ui.Error(e.Error())
		}
		return 1
	}

	if err := data.generate(); err != nil {
		fmt.Printf("Failed to generate: %+v", err)
		return 1
	}

	return 0
}

func (c DocumentCommand) Help() string {
	return "partially there..."
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

	schema, err := helpers.ParseProviderJSON(helpers.OpenProviderJSON("/tmp/azurerm-provider-out.json"), d.SnakeName, helpers.DocType(d.DocType))
	if err != nil || schema == nil {
		fmt.Printf("[ERROR] reading %s %s from provider JSON", d.DocType, d.Name)
		return err
	}

	// terraform marks the id as optional in the output JSON, this makes the template super messy, so we'll flip the bool here
	if tmpID, ok := schema.Block.Attributes["id"]; ok {
		tmpID.Optional = false
		schema.Block.Attributes["id"] = tmpID
	}

	d.Schema = *schema

	tpl := template.Must(template.New("document.gotpl").Funcs(TplFuncMap).ParseFS(Templatedir, "templates/document.gotpl"))

	outputPath := ""
	if d.ServicePackage != "" {
		outputPath = fmt.Sprintf("%s/internal/services/%s/docs/%s.md", strings.ToLower(d.ProviderName), strings.ToLower(strcase.ToCamel(d.ServicePackage)), d.SnakeName)
	} else {
		outputPath = fmt.Sprintf("website/docs/%s.md", d.SnakeName)
	}

	f, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("[DEBUG] failed opening output resource file for writing: %+v", err.Error())
		return err
	}

	defer f.Close()

	//fmt.Printf("[STEBUG] d: %+v", d)
	err = tpl.Execute(f, d)
	if err != nil {
		return err
	}
	return nil
}
