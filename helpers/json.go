package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func ParseProviderJSON(input Provider, resourceName string, docType DocType) (resource *ResourceSchema, err error) {
	for _, p := range input.ProviderSchemas {
		switch docType {
		case DocTypeResource:
			if schema, ok := p.ResourceSchemas[resourceName]; ok {
				resource = &schema
			} else {
				return nil, fmt.Errorf("%s %s not found", docType, resourceName)
			}
		case DocTypeDataSource:
			if schema, ok := p.DataSourceSchemas[resourceName]; ok {
				resource = &schema
			} else {
				return nil, fmt.Errorf("%s %s not found", docType, resourceName)
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
	_ = json.Unmarshal(byteValue, &result)

	return result
}
