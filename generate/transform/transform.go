package transform

import (
	"github.com/iancoleman/strcase"
	"strings"
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

// Represents an 'augment' YAML file, used to describe manual changes from the MS Graph OpenAPI spec
type TemplateAugment struct {
	ExcludedProperties       []string            `yaml:"excludedProperties"`
	AltReadMethods           []map[string]string `yaml:"altReadMethods"`
	DataSourceExtraOptionals []string            `yaml:"dataSourceExtraOptionals"`
	ResourceExtraComputed    []string            `yaml:"resourceExtraComputed"`
}
