package transform

import (
	"os"
	"gopkg.in/yaml.v3"
	"strings"
	"terraform-provider-msgraph/generate/extract"

	"github.com/iancoleman/strcase"
)

func upperFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
}

func pathFieldName(s string) (string, string) {
	s = strings.TrimLeft(s, "{")
	s = strings.TrimRight(s, "}")
	pLeft, pRight, _ := strings.Cut(s, "-")
	return pLeft, pRight
}

type strWithCases struct {
	String string
}

func (s strWithCases) LowerCamel() string {
	return strcase.ToLowerCamel(s.String)
}

func (s strWithCases) UpperCamel() string {
	return strcase.ToCamel(s.String)
}

func (s strWithCases) Snake() string {
	return strcase.ToSnake(s.String)
}

func (s strWithCases) UpperFirst() string {
	return strings.ToUpper(s.String[0:1]) + s.String[1:]
}

// Represents a method used to perform a query using msgraph-sdk-go
type queryMethod struct {
	MethodName string
	Parameter  string
}

type TemplateInput struct {
	OpenAPIPath   extract.OpenAPIPathObject
}

func (ti TemplateInput) PackageName() string {
	return strings.ToLower(strings.Split(ti.OpenAPIPath.Path, "/")[1])
}

func (ti TemplateInput) BlockName() strWithCases {

	pathFields := strings.Split(ti.OpenAPIPath.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array

	// Generate name of the terraform block
	blockName := ""
	if len(pathFields) > 1 {
		for _, p := range pathFields[1:] {
			if strings.HasPrefix(p, "{") {
				pLeft, _ := pathFieldName(p)
				blockName += pLeft
			} else {
				blockName += p
			}
		}
	} else {
		blockName = pathFields[0]
	}

	return strWithCases{String: blockName}
}

func (ti TemplateInput) Model() model {
	return model{Template: &ti}
}

func (ti TemplateInput) SchemaDescription() string {
	return strings.Replace(ti.OpenAPIPath.Get().Response().Description(), "\n", " ", -1)
}

func (ti TemplateInput) SchemaDataSource() schema {
	return schema{Template: &ti, BehaviourMode: "DataSource"}
}

func (ti TemplateInput) SchemaResource() schema {
	return schema{Template: &ti, BehaviourMode: "Resource"}
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
type templateAugment struct {
	ExcludedProperties       []string            `yaml:"excludedProperties"`
	AltReadMethods           []map[string]string `yaml:"altReadMethods"`
	DataSourceExtraOptionals []string            `yaml:"dataSourceExtraOptionals"`
	ResourceExtraComputed    []string            `yaml:"resourceExtraComputed"`
}

func (ti TemplateInput) Augment() templateAugment {

	pathFields := strings.Split(ti.OpenAPIPath.Path, "/")[1:] // Paths start with a '/', so we need to get rid of the first empty entry in the array
	packageName := strings.ToLower(pathFields[0])

	// Open augment file if available
	var err error = nil
	augment := templateAugment{}
	augmentFile, err := os.ReadFile("generate/augment/" + packageName + "/" + ti.BlockName().LowerCamel() + ".yaml")
	if err == nil {
		yaml.Unmarshal(augmentFile, &augment)
	}

	return augment
}

