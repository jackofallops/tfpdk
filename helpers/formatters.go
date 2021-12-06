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

func SchemaItemFormatter(input interface{}, name string) string {
	data := Schema{}
	switch v := input.(type) {
	case Schema:
		data = v
	case map[string]interface{}:
		data = flattenMapToSchema(v)
	}
	// Block detection
	var optionalOrRequired, desc, forceNew string
	if data.Required {
		optionalOrRequired = "(Required) "
	} else if data.Optional {
		optionalOrRequired = "(Optional) "
	}
	if strings.EqualFold(data.Type, "typelist") || strings.EqualFold(data.Type, "typeset") {
		return fmt.Sprintf("* `%s` %s- %s block as detailed below.", name, optionalOrRequired, PrefixedLabelString(name))
	}

	if data.Description != "" {
		desc = strings.TrimSpace(data.Description) // TODO auto-fix line ends for MD linting?
	} else {
		desc = "// TODO - Add missing `Description` to schema for this property and regenerate this file."
	}

	if data.ForceNew {
		forceNew = "Changing this forces a new resource to be created."
		return fmt.Sprintf("* `%s` %s- %s %s", name, optionalOrRequired, desc, forceNew)
	}

	return fmt.Sprintf("* `%s` %s- %s", name, optionalOrRequired, desc)
}

func SchemaItemFormatterSpecial(input Schema, name string) string {
	switch name {
	case "name", "resource_group_name", "location":
		var optionalOrRequired, desc, forceNew string
		if input.Required {
			optionalOrRequired = "(Required)"
		} else if input.Optional {
			optionalOrRequired = "(Optional)"
		}
		if input.Description != "" {
			desc = strings.TrimSpace(input.Description) // TODO auto-fix line ends for MD linting?
		} else {
			desc = "// TODO - Add missing `Description` to schema for this property and regenerate this file."
		}
		if input.ForceNew {
			forceNew = "Changing this forces a new resource to be created."
			return fmt.Sprintf("* `%s` %s - %s %s", name, optionalOrRequired, desc, forceNew)
		}

		return fmt.Sprintf("* `%s` %s - %s", name, optionalOrRequired, desc)
	}
	return ""
}

func flattenMapToSchema(input map[string]interface{}) Schema {
	output := Schema{}
	if t, ok := input["type"]; ok {
		output.Type = t.(string)
	}

	if t, ok := input["config_mode"]; ok {
		output.ConfigMode = t.(string)
	}

	if t, ok := input["optional"]; ok {
		output.Optional = t.(bool)
	}

	if t, ok := input["required"]; ok {
		output.Required = t.(bool)
	}

	if t, ok := input["default"]; ok {
		output.Default = t
	}

	if t, ok := input["description"]; ok {
		output.Description = t.(string)
	}

	if t, ok := input["computed"]; ok {
		output.Computed = t.(bool)
	}

	if t, ok := input["force_new"]; ok {
		output.ForceNew = t.(bool)
	}

	if t, ok := input["elem"]; ok {
		output.Elem = t
	}

	if t, ok := input["max_items"]; ok {
		output.MaxItems = int(t.(float64))
	}

	if t, ok := input["min_items"]; ok {
		output.MinItems = int(t.(float64))
	}

	return output
}
