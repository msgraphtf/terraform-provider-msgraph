package main

import (
	//"encoding/csv"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/gocarina/gocsv"
	"github.com/iancoleman/strcase"
)

type templateInput struct {
	PackageName              string
	DataSourceName           string
	DataSourceNameUpperCamel string
	DataSourceNameLowerCamel string
	DataSourceNameSnake      string
	Schema                   []schemaAttribute
}

// Used by templates defined inside of data_source_template.go
type schemaAttribute struct {
	NameUpperCamel      string
	NameLowerCamel      string
	NameSnake           string
	TypeModel           string
	AttributeType       string
	MarkdownDescription string
	Required            bool
	Optional            bool
	Computed            bool
	ElementType         string
	Attributes          []schemaAttribute
	NestedObject        []schemaAttribute
}

type csvSchema struct {
	Name        string `csv:"Property"`
	Type        string `csv:"Type"`
	Computed    bool   `csv:"Computed"`
	Optional    bool   `csv:"Optional"`
	Required    bool   `csv:"Required"`
	Description string `csv:"Description"`
}

var dataSourceName string
var packageName string

func openCsv(path string) []*csvSchema {

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		r.LazyQuotes = true
		return r // Allows use pipe as delimiter
	})

	f, err := os.Open(path)
	if err != nil {
		fmt.Print(err)
	}
	defer f.Close()
	csv := []*csvSchema{}
	gocsv.UnmarshalFile(f, &csv)

	return csv

}

func generateSchema(schema *[]schemaAttribute, csv []*csvSchema) {
	for _, row := range csv {

		// Create new schemaAttribute for array
		nextAttribute := new(schemaAttribute)

		nextAttribute.NameUpperCamel = strcase.ToCamel(row.Name)
		nextAttribute.NameLowerCamel = strcase.ToLowerCamel(row.Name)
		nextAttribute.NameSnake = strcase.ToSnake(row.Name)

		// Convert types from MS Graph docs to Go and terraform types
		switch {
		case row.Type == "String":
			nextAttribute.AttributeType = "String"
			nextAttribute.TypeModel = "types.String"
		case row.Type == "String collection":
			nextAttribute.AttributeType = "List"
			nextAttribute.TypeModel = "[]types.String"
			nextAttribute.ElementType = "types.StringType"
		case row.Type == "Boolean":
			nextAttribute.AttributeType = "Bool"
			nextAttribute.TypeModel = "types.Bool"
		case row.Type == "DateTimeOffset":
			nextAttribute.AttributeType = "String"
			nextAttribute.TypeModel = "types.String"
		case strings.HasSuffix(row.Type, "collection"):
			nextAttribute.AttributeType = "ListNested"
			nextAttribute.TypeModel = "[]" + dataSourceName + "DataSource" + strcase.ToCamel(row.Type)

			nestedCsv := openCsv("template/input/" + packageName + "/" + nextAttribute.NameSnake + ".csv")
			var nestedAttributes []schemaAttribute
			generateSchema(&nestedAttributes, nestedCsv)

			nextAttribute.NestedObject = nestedAttributes
		default:
			nextAttribute.AttributeType = "SingleNested"
			nextAttribute.TypeModel = "*" + dataSourceName + "DataSource" + strcase.ToCamel(row.Type)

			nestedCsv := openCsv("template/input/" + packageName + "/" + nextAttribute.NameSnake + ".csv")
			var nestedAttributes []schemaAttribute
			generateSchema(&nestedAttributes, nestedCsv)

			nextAttribute.Attributes = nestedAttributes
		}

		nextAttribute.Computed = row.Computed
		nextAttribute.Optional = row.Optional
		nextAttribute.Required = row.Required
		nextAttribute.MarkdownDescription = row.Description

		*schema = append(*schema, *nextAttribute)
	}
}

func main() {

	// Get template
	templateDataSource := template.New("dataSource")
	templateFile, err := os.ReadFile("template/data_source_template.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource, err = templateDataSource.Parse(string(templateFile))

	// Get inputs
	packageName = os.Args[1]
	dataSourceName = os.Args[2]

	// Open top level CSV
	csv := openCsv("template/input/" + packageName + "/" + dataSourceName + ".csv")

	// Generate schema values from CSV columns
	var schema []schemaAttribute
	generateSchema(&schema, csv)

	// Set input values to top level template
	inputValues := templateInput{
		PackageName:              packageName,
		DataSourceName:           dataSourceName,
		DataSourceNameUpperCamel: strcase.ToCamel(dataSourceName),
		DataSourceNameLowerCamel: strcase.ToLowerCamel(dataSourceName),
		DataSourceNameSnake:      strcase.ToSnake(dataSourceName),
		Schema:                   schema,
	}

	os.MkdirAll("template/out/", os.ModePerm)
	outfile, err := os.Create("template/out/" + dataSourceName + "_data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource.Execute(outfile, inputValues)

}
