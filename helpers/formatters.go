package helpers

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

func TerraformResourceName(provider, resourceName string) string {
	fmtStr := "%s_%s"
	return fmt.Sprintf(fmtStr, strings.ToLower(provider), strcase.ToSnake(resourceName))
}

func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case DocType:
		return string(v)
	default:
		return ""
	}
}

func PrefixedDescriptionString(input string) string {
	prefix := "a"
	first := input[0:1]
	vowel, _ := regexp.Match(first, []byte(`aeiouAEIOU`))

	if vowel {
		prefix = "an"
	}
	return fmt.Sprintf("%s %s", prefix, strings.Title(strcase.ToDelimited(input, ' ')))
}

func ToDelimTitle(input string) string {
	return strings.Title(strcase.ToDelimited(input, ' '))
}

func PrefixedLabelString(input string) string {
	prefix := "a"
	first := input[0:1]
	vowel, _ := regexp.Match(first, []byte(`aeiouAEIOU`))

	if vowel {
		prefix = "an"
	}

	return fmt.Sprintf("%s `%s`", prefix, input)
}
