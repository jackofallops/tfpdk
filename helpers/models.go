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

type Schema struct {
	Type        string      `json:"type,omitempty"`
	ConfigMode  string      `json:"config_mode,omitempty"`
	Optional    bool        `json:"optional,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Description string      `json:"description,omitempty"`
	Computed    bool        `json:"computed,omitempty"`
	ForceNew    bool        `json:"force_new,omitempty"`
	Elem        interface{} `json:"elem,omitempty"`
	MaxItems    int         `json:"max_items,omitempty"`
	MinItems    int         `json:"min_items,omitempty"`
}

type Resource struct {
	Schema   map[string]Schema `json:"schema"`
	Timeouts *Timeouts         `json:"timeouts,omitempty"`
}

type Timeouts struct {
	Create int `json:"create,omitempty"`
	Read   int `json:"read,omitempty"`
	Delete int `json:"delete,omitempty"`
	Update int `json:"update,omitempty"`
}
