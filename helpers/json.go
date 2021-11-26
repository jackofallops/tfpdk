package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ParseProviderJSON(input Provider, providerName string, resourceName string, docType DocType) (resource *ResourceSchema, err error) {
	for n, p := range input.ProviderSchemas {
		if n == fmt.Sprintf("registry.terraform.io/hashicorp/%s", providerName) {
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
	}
	return resource, nil
}

// OpenProviderJSON opens an exported Terraform Provider JSON file for parsing
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

// GetTerraformSchemaJSON calls the main Terraform binary to
// export the known provider Schemas and marshals the data in to Go structs
func GetTerraformSchemaJSON() Provider {
	f, err := CallTerraform("providers", "schema", "-json")

	if err != nil {
		fmt.Printf("could read Provider JSON from Terraform: %+v", err)
	}

	var result Provider
	_ = json.Unmarshal(f, &result)

	return result
}

func GetSchema(name string) (*Resource, error) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/schema-data/v1/resources/%s", name))
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &Resource{}
	err = json.Unmarshal(content, result)

	return result, err
}
