package main

import (
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
	DataSourceAttributeName  string
	Schema                   []attributeSchema
	Model                    []attributeModel
}

// Used by templates defined inside of data_source_template.go to generate the schema
type attributeSchema struct {
	AttributeName       string
	AttributeType       string
	MarkdownDescription string
	Required            bool
	Optional            bool
	Computed            bool
	ElementType         string
	Attributes          []attributeSchema
	NestedObject        []attributeSchema
}

// Used by templates defined inside of data_source_template.go to generate the data models
type attributeModel struct {
	ModelName string
	Fields    []modelField
}

type modelField struct {
	FieldName     string
	FieldType     string
	AttributeName string
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

func generateSchema(schema *[]attributeSchema, csv []*csvSchema) {
	for _, row := range csv {

		// Create new attribute schema and model for array
		nextAttributeSchema := new(attributeSchema)

		nextAttributeSchema.AttributeName = strcase.ToSnake(row.Name)

		// Convert types from MS Graph docs to Go and terraform types
		switch {
		case row.Type == "String":
			nextAttributeSchema.AttributeType = "String"
		case row.Type == "String collection":
			nextAttributeSchema.AttributeType = "List"
			nextAttributeSchema.ElementType = "types.StringType"
		case row.Type == "Boolean":
			nextAttributeSchema.AttributeType = "Bool"
		case row.Type == "DateTimeOffset":
			nextAttributeSchema.AttributeType = "String"
		case strings.HasSuffix(row.Type, "collection"):
			nextAttributeSchema.AttributeType = "ListNested"

			nestedCsv := openCsv("template/input/" + packageName + "/" + nextAttributeSchema.AttributeName + ".csv")
			var nestedAttributes []attributeSchema
			generateSchema(&nestedAttributes, nestedCsv)

			nextAttributeSchema.NestedObject = nestedAttributes
		default:
			nextAttributeSchema.AttributeType = "SingleNested"

			nestedCsv := openCsv("template/input/" + packageName + "/" + nextAttributeSchema.AttributeName + ".csv")
			var nestedAttributes []attributeSchema
			generateSchema(&nestedAttributes, nestedCsv)

			nextAttributeSchema.Attributes = nestedAttributes
		}

		nextAttributeSchema.Computed = row.Computed
		nextAttributeSchema.Optional = row.Optional
		nextAttributeSchema.Required = row.Required
		nextAttributeSchema.MarkdownDescription = row.Description

		*schema = append(*schema, *nextAttributeSchema)
	}
}

func generateModel(modelName string, model *[]attributeModel, csv []*csvSchema) {

	newModel := attributeModel{
		ModelName: modelName,
	}

	for _, row := range csv {

		nextModelField := new(modelField)
		nextModelField.FieldName = strcase.ToCamel(row.Name)
		nextModelField.AttributeName = strcase.ToSnake(row.Name)

		switch {
		case row.Type == "String":
			nextModelField.FieldType = "types.String"
		case row.Type == "Boolean":
			nextModelField.FieldType = "types.Bool"
		case row.Type == "DateTimeOffset":
			nextModelField.FieldType = "types.String"
		case row.Type == "String collection":
			nextModelField.FieldType = "[]types.String"
		case strings.HasSuffix(row.Type, "collection"):
			nestedModel := attributeModel{
				ModelName: "[]" + dataSourceName + strcase.ToCamel(row.Name) + "Model",
			}
			nextModelField.FieldType = nestedModel.ModelName
		default:
			nestedModel := attributeModel{
				ModelName: dataSourceName + strcase.ToCamel(row.Name) + "Model",
			}
			nextModelField.FieldType = nestedModel.ModelName
		}

		newModel.Fields = append(newModel.Fields, *nextModelField)

	}

	*model = append(*model, newModel)

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

	// Generate schema values from CSV
	var schema []attributeSchema
	generateSchema(&schema, csv)

	// Generate model values from CSV
	var model []attributeModel
	generateModel(strcase.ToLowerCamel(dataSourceName)+"DataSourceModel", &model, csv)

	// Set input values to top level template
	inputValues := templateInput{
		PackageName:              packageName,
		DataSourceName:           dataSourceName,
		DataSourceNameUpperCamel: strcase.ToCamel(dataSourceName),
		DataSourceNameLowerCamel: strcase.ToLowerCamel(dataSourceName),
		DataSourceAttributeName:  strcase.ToSnake(dataSourceName),
		Schema:                   schema,
		Model:                    model,
	}

	os.MkdirAll("template/out/", os.ModePerm)
	outfile, err := os.Create("template/out/" + dataSourceName + "_data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource.Execute(outfile, inputValues)

}
