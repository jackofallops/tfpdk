package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// Deprecated: ParseProviderJSON parsed the Terraform core view of the schema into a model that could be passed to the
// documentation template. This lacked the required detail to generate provider resource docs without user effort.
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

// Deprecated: OpenProviderJSON opens an exported Terraform Provider JSON file for parsing.
// This has been deprecated in favour of creating a local HTTP server in the provider, which will hopefully
// be moved upstream to the PluginSDK in the future.
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

// Deprecated: GetTerraformSchemaJSON calls the main Terraform binary to
// export the known provider Schemas and marshals the data in to Go structs.
// This has been deprecated in favour of creating a local HTTP server in the provider, which will hopefully
// be moved upstream to the PluginSDK in the future.
func GetTerraformSchemaJSON() Provider {
	f, err := CallTerraform("providers", "schema", "-json")

	if err != nil {
		fmt.Printf("could read Provider JSON from Terraform: %+v", err)
	}

	var result Provider
	_ = json.Unmarshal(f, &result)

	return result
}

func GetSchema(apiPath string, docType string, name string) (*Resource, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%ss/%s", apiPath, docType, name))
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
