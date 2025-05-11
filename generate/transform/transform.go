package transform

import (
	"strings"
	"terraform-provider-msgraph/generate/openapi"

	"github.com/iancoleman/strcase"
)

func upperFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func PathFieldName(s string) (string, string) {
	s = strings.TrimLeft(s, "{")
	s = strings.TrimRight(s, "}")
	pLeft, pRight, _ := strings.Cut(s, "-")
	return pLeft, pRight
}

type StrWithCases struct {
	String string
}

func (s StrWithCases) LowerCamel() string {
	return strcase.ToLowerCamel(s.String)
}

func (s StrWithCases) UpperCamel() string {
	return strcase.ToCamel(s.String)
}

func (s StrWithCases) Snake() string {
	return strcase.ToSnake(s.String)
}

func (s StrWithCases) UpperFirst() string {
	return strings.ToUpper(s.String[0:1]) + s.String[1:]
}

// Represents a method used to perform a query using msgraph-sdk-go
type queryMethod struct {
	MethodName string
	Parameter  string
}

type TemplateInput struct {
	PackageName   string
	OpenAPIPath   openapi.OpenAPIPathObject
	Schema        TerraformSchema
	Model         Model
	Augment       TemplateAugment
}

func (ti TemplateInput) BlockName() StrWithCases {

	pathFields := strings.Split(ti.OpenAPIPath.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array

	// Generate name of the terraform block
	blockName := ""
	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := PathFieldName(p)
				blockName += pLeft
			} else {
				blockName += p
			}
		}
	} else {
		blockName = pathFields[0]
	}

	return StrWithCases{String: blockName}
}

func (ti TemplateInput) ReadQuery() readQuery {
	return readQuery{Template: &ti}
}

func (ti TemplateInput) ReadResponse() readResponse {
	return readResponse{Template: &ti}
}

func (ti TemplateInput) CreateRequest() createRequest {
	return createRequest{Template: &ti}
}

func (ti TemplateInput) UpdateRequest() updateRequest {
	return updateRequest{Template: &ti}
}

// Represents an 'augment' YAML file, used to describe manual changes from the MS Graph OpenAPI spec
type TemplateAugment struct {
	ExcludedProperties       []string            `yaml:"excludedProperties"`
	AltReadMethods           []map[string]string `yaml:"altReadMethods"`
	DataSourceExtraOptionals []string            `yaml:"dataSourceExtraOptionals"`
	ResourceExtraComputed    []string            `yaml:"resourceExtraComputed"`
}
