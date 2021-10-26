package helpers

import (
	"errors"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
)

const (
	ConfigFileName = ".tfpdk.hcl"
)

type Configuration struct {
	ProviderName          string `hcl:"provider_name,optional"`
	ServicePackagesPath   string `hcl:"service_packages_path,optional"`
	DocsPath              string `hcl:"docs_path,optional"`
	ProviderGithubOrg     string `hcl:"provider_github_org,optional"`
	ResourceDocsDirname   string `hcl:"resource_docs_directory_name,optional"`
	DataSourceDocsDirname string `hcl:"data_source_docs_directory_name,optional"`
	TypedSDK              bool   `hcl:"use_typed_sdk,optional"`
}

// LoadConfig loads the configuration file if present to allow users to override various settings in the
// tool, such as path to services, docs and any SDK options.
func LoadConfig() *Configuration {
	config := Configuration{
		ServicePackagesPath:   "internal/services",
		ProviderGithubOrg:     "hashicorp",
		DocsPath:              "docs",
		ResourceDocsDirname:   "r",
		DataSourceDocsDirname: "d",
		TypedSDK:              false,
	}
	p, err := ProviderName()
	if err != nil || p == nil {
		log.Printf("failed to determine provider name %+v", err)
	}
	config.ProviderName = *p
	info, err := os.Stat(ConfigFileName)
	if !errors.Is(err, os.ErrNotExist) && !info.IsDir() {
		err := hclsimple.DecodeFile(ConfigFileName, nil, &config)
		if err != nil {
			log.Printf("Failed to load configuration: %s", err)
		}
		// fmt.Printf("Using configuration %+v", config)
	}

	return &config
}
