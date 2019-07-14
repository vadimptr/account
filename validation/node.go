package validation

type Node struct {
	Type               []string      `json:"type,omitempty"`
	AdditionalProperty *bool         `json:"additionalProperties,omitempty"`
	Required           []string      `json:"required,omitempty"`
	Properties         Properties    `json:"properties,omitempty"`
	MaxLength          *int          `json:"maxLength,omitempty"`
	MinLength          *int          `json:"minLength,omitempty"`
	Minimum            *int          `json:"minimum,omitempty"`
	Maximum            *int          `json:"maximum,omitempty"`
	Default            *int          `json:"default,omitempty"`
	Pattern            *string       `json:"pattern,omitempty"`
	OneOf              []Node        `json:"oneOf,omitempty"`
	Not                *Node         `json:"not,omitempty"`
	Format             *string       `json:"format,omitempty"`
	Enum               []interface{} `json:"enum,omitempty"`
	Items              *Node         `json:"items,omitempty"`
}

type Properties map[string]Node

var StringType = []string{"string"}
var MaybeStringType = []string{"null", "string"}
var IntegerType = []string{"integer"}
var MaybeIntegerType = []string{"null", "integer"}
var ObjectType = []string{"object"}
var MaybeObjectType = []string{"null", "object"}
var ArrayType = []string{"array"}
var BooleanType = []string{"boolean"}
var MaybeBooleanType = []string{"null", "boolean"}
var ValueTrue = true
var ValueFalse = false
var Value0 = 0
var Value1 = 1
var Value6 = 6
var Value20 = 20
var Value100 = 100
var Value256 = 256
var Value1024 = 1024
