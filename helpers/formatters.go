package helpers

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
)

func TerraformResourceName(provider, resourceName string) string {
	fmtStr := "%s_%s"
	return fmt.Sprintf(fmtStr, strings.ToLower(provider), strcase.ToSnake(resourceName))
}
