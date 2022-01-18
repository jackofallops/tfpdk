package helpers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetSchema(apiPath string, docType string, name string) (*Resource, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%ss/%s", apiPath, docType, name))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("error, schema for item %q not found", name)
	}

	if contentType := resp.Header.Get("Content-Type"); !strings.EqualFold(contentType, "application/json; charset=UTF-8") {
		return nil, fmt.Errorf("error reading API, expected content type 'application/json; charset=UTF-8', got %q", contentType)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	result := &Resource{}
	err = json.Unmarshal(content, result)

	return result, err
}
