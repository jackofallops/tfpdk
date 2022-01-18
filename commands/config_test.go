package commands

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

/*
`tfpdk config` currently does not take cli options, but it may later, so leaving the test capable of checking them

*/

func TestConfigTemplateOutput(t *testing.T) {
	cases := []struct {
		ServicePackagePath    string
		ProviderGithubOrg     string
		DocsPath              string
		ResourcesDocsDirname  string
		DataSourceDocsDirname string
		Typed                 bool
		Expected              string
	}{
		{
			// No User
			Expected: `service_packages_path           = "internal/services"
provider_github_org             = "hashicorp"
docs_path                       = "docs"
resource_docs_directory_name    = "r"
data_source_docs_directory_name = "d"
use_typed_sdk                   = false
`,
		},
	}

	for _, tc := range cases {
		lc := localConfig(*config)
		err := lc.generate()
		if err != nil {
			t.Fatalf("generating output config file: %+v", err)
		}
		actual, err := os.ReadFile(".tfpdk.hcl")
		if err != nil {
			t.Fatalf("reading config output file: %+v", err)
		}
		if !reflect.DeepEqual(string(actual), tc.Expected) {
			t.Fatalf("generated config does not match expected")
		}

		// Tidy up
		err = os.Remove(".tfpdk.hcl")
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			t.Errorf("removing test .tfpdk.hcl file")
		}
	}
}
