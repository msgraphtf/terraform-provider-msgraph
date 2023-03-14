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

type schemaAttribute struct {
	NameUpperCamel      string
	NameLowerCamel      string
	NameSnake           string
	TypeModel           string
	AttributeType          string
	MarkdownDescription string
	Required            bool
	Optional            bool
	Computed            bool
	ElementType         string
	NestedObject        *schemaAttribute
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
			//nextAttribute.NestedObject = true
		default:
			nextAttribute.AttributeType = "SingleNested"
			nextAttribute.TypeModel = "*" + dataSourceName + "DataSource" + strcase.ToCamel(row.Type)
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

	// Configure gocsv
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.Comma = '|'
		r.LazyQuotes = true
		return r // Allows use pipe as delimiter
	})

	f, err := os.Open("template/input/" + packageName + "/" + dataSourceName + ".csv")
	defer f.Close()
	rawCsv := []*csvSchema{}
	gocsv.UnmarshalFile(f, &rawCsv)

	// Generate schema values from CSV columns
	var schema []schemaAttribute
	generateSchema(&schema, rawCsv)

	// Set inputs on struct
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
