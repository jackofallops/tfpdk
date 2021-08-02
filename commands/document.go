package commands

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/jackofallops/tfpdk/helpers"
	"github.com/mitchellh/cli"
)

type DocumentCommand struct{}

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

func (c DocumentCommand) Run(args []string) int {
	data := DocumentData{}

	if len(args) == 0 {
		fmt.Print(c.Help())
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
			fmt.Printf(c.Help())
			return 0

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
				fmt.Println("argument `servicepackage` requires a value, eg `-servicepackage=somePackageName`")
				return 1
			}
		case "type":
			if len(arg) == 2 {
				data.DocType = arg[1]
			} else {
				fmt.Println("argument `type` requires a value, eg `-type=resource`")
				return 1
			}
		case "id":
			if len(arg) == 2 {
				data.IDExample = arg[1]
			} else {
				fmt.Println("argument `id` requires a quoted value, eg `-id=\"/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/group1\"`")
				return 1
			}

		default:
			fmt.Printf("unrecognised option %q", arg[0])
			return 1
		}
	}

	if data.Name == "" {
		fmt.Printf("Error: missing required argument `-name`")
		return 1
	}

	if err := data.generate(); err != nil {
		fmt.Printf("Failed to generate: %+v", err)
		return 1
	}

	return 0
}

func (c DocumentCommand) Help() string {
	return "Not implemented yet..."
}

func (c DocumentCommand) Synopsis() string {
	return "generates documentation from a resource"
}

func (d DocumentData) generate() error {
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
	if err != nil {
		fmt.Printf("[ERROR] reading %s %s from provider JSON", d.DocType, d.Name)
		return err
	}

	// terraform marks the id as optional in the output JSON, this makes the template super messy, so we'll flip the bool here
	tmpID := schema.Block.Attributes["id"]
	tmpID.Optional = false
	schema.Block.Attributes["id"] = tmpID

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
