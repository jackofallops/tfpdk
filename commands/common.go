package commands

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/jackofallops/tfpdk/helpers"
)

var config = &helpers.Configuration{}

var TplFuncMap = template.FuncMap{
	"ToLower":                   strings.ToLower,
	"ToTitle":                   strings.Title,
	"ToCamel":                   strcase.ToCamel,
	"ToSnake":                   strcase.ToSnake,
	"TfName":                    helpers.TerraformResourceName,
	"ToString":                  helpers.ToString,
	"ToDelim":                   strcase.ToDelimited,
	"ToDelimTitle":              helpers.ToDelimTitle,
	"PrefixedDescriptionString": helpers.PrefixedDescriptionString,
}

func init() {
	config = helpers.LoadConfig()
}
