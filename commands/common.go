package commands

import (
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/jackofallops/tfpdk/helpers"
)

var TplFuncMap = template.FuncMap{
	"ToLower":             strings.ToLower,
	"ToTitle":             strings.Title,
	"ToCamel":             strcase.ToCamel,
	"ToSnake":             strcase.ToSnake,
	"TfName":              helpers.TerraformResourceName,
	"ToString":            helpers.ToString,
	"ToDelim":             strcase.ToDelimited,
	"ToDescriptionString": helpers.ToDescriptionString,
}

func init() {
	config = helpers.LoadConfig()
}
