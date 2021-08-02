package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ParseProviderJSON(input Provider, resourceName string, Type DocType) (resource *ResourceSchema, err error) {
	for _, p := range input.ProviderSchemas {
		switch Type {
		case DocTypeResource:
			if schema, ok := p.ResourceSchemas[resourceName]; ok {
				resource = &schema
			} else {
				return nil, fmt.Errorf("%s %s not found", Type, resourceName)
			}
		case DocTypeDataSource:
			if schema, ok := p.DataSourceSchemas[resourceName]; ok {
				resource = &schema
			} else {
				return nil, fmt.Errorf("%s %s not found", Type, resourceName)
			}
		}
	}
	return resource, nil
}

func OpenProviderJSON(filename string) Provider {
	f, err := os.Open(filename)

	if err != nil {
		fmt.Printf("could not open Provider JSON file: %+v", err)
	}
	defer f.Close()
	byteValue, _ := ioutil.ReadAll(f)

	var result Provider
	err = json.Unmarshal(byteValue, &result)

	return result
}
