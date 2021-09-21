package helpers

import "encoding/json"

type DocType string

const DocTypeDataSource DocType = "datasource"
const DocTypeResource DocType = "resource"

type Provider struct {
	Version         string                    `json:"format_version"`
	ProviderSchemas map[string]ProviderSchema `json:"provider_schemas,omitempty"`
}

type ProviderSchema struct {
	Provider          map[string]json.RawMessage `json:"provider"`
	ResourceSchemas   map[string]ResourceSchema  `json:"resource_schemas,omitempty"`
	DataSourceSchemas map[string]ResourceSchema  `json:"data_source_schemas,omitempty"`
}

type ResourceSchema struct {
	Version int           `json:"version,omitempty"`
	Block   ResourceBlock `json:"block,omitempty"`
}

type ResourceBlock struct {
	Attributes      map[string]ResourceProperty  `json:"attributes,omitempty"`
	BlockTypes      map[string]ResourceBlockType `json:"block_types,omitempty"`
	DescriptionKind string                       `json:"description_kind,omitempty"`
}

type ResourceBlockType struct {
	NestingMode string        `json:"nesting_mode,omitempty"`
	Block       ResourceBlock `json:"block,omitempty"`
}

type ResourceProperty struct {
	Block       json.RawMessage `json:"block,omitempty"`
	Description string          `json:"description,omitempty"`
	Computed    bool            `json:"computed,omitempty"`
	Optional    bool            `json:"optional,omitempty"`
	Required    bool            `json:"required,omitempty"`
	Type        interface{}     `json:"type,omitempty"`    // This can be string, array, or object, or a mix...
	Default     interface{}     `json:"default,omitempty"` // This can be multiple primitive types
	ForceNew    bool            `json:"force_new,omitempty"`
	Deprecated  bool            `json:"deprecated,omitempty"`
	Kind        string          `json:"description_kind,omitempty"`
}

type BlockType struct {
}
