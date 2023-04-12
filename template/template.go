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
	PreRead                  string
	Read                     []attributeRead
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
	Fields    []attributeModelField
}

type attributeModelField struct {
	FieldName     string
	FieldType     string
	AttributeName string
}

type attributeRead struct {
	GetMethod string
	StateAttributeName string
	ModelVarName string
	ModelName string
	AttributeType string
	DataSourceName string
	NestedRead []attributeRead
	ParentRead *attributeRead
	ResultVarName string
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
		case row.Type == "String" || row.Type == "Guid":
			nextAttributeSchema.AttributeType = "String"
		case row.Type == "String collection" || row.Type == "Guid collection":
			nextAttributeSchema.AttributeType = "List"
			nextAttributeSchema.ElementType = "types.StringType"
		case row.Type == "Boolean":
			nextAttributeSchema.AttributeType = "Bool"
		case row.Type == "DateTimeOffset":
			nextAttributeSchema.AttributeType = "String"
		case row.Type == "String collection":
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
	var nestedModels []attributeModel

	for _, row := range csv {

		nextModelField := new(attributeModelField)
		nextModelField.FieldName = strcase.ToCamel(row.Name)
		nextModelField.AttributeName = strcase.ToSnake(row.Name)

		switch {
		case row.Type == "String" || row.Type == "Guid":
			nextModelField.FieldType = "types.String"
		case row.Type == "Boolean":
			nextModelField.FieldType = "types.Bool"
		case row.Type == "DateTimeOffset":
			nextModelField.FieldType = "types.String"
		case row.Type == "String collection" || row.Type == "Guid collection":
			nextModelField.FieldType = "[]types.String"
		case strings.HasSuffix(row.Type, "collection"):
			nextModelField.FieldType = "[]" + dataSourceName + strcase.ToCamel(row.Name) + "DataSourceModel"

			nestedCsv := openCsv("template/input/" + packageName + "/" + strcase.ToSnake(row.Name) + ".csv")
			generateModel(dataSourceName+strcase.ToCamel(row.Name)+"DataSourceModel", &nestedModels, nestedCsv)

		default:
			nextModelField.FieldType = "*" + dataSourceName + strcase.ToCamel(row.Name) + "DataSourceModel"

			nestedCsv := openCsv("template/input/" + packageName + "/" + strcase.ToSnake(row.Name) + ".csv")
			generateModel(dataSourceName+strcase.ToCamel(row.Name)+"DataSourceModel", &nestedModels, nestedCsv)

		}

		newModel.Fields = append(newModel.Fields, *nextModelField)

	}

	*model = append(*model, newModel)
	if len(nestedModels) != 0 {
		*model = append(*model, nestedModels...)
	}

}

func generateRead(read *[]attributeRead, csv []*csvSchema, parent *attributeRead) {

	for _, row := range csv {

		nextAttributeRead := attributeRead{
			ModelVarName: strcase.ToLowerCamel(row.Name),
			DataSourceName: dataSourceName,
			ResultVarName: "result",
		}
		if parent != nil && parent.AttributeType == "SingleNested" {
			nextAttributeRead.ParentRead = parent
			nextAttributeRead.GetMethod = parent.GetMethod+".Get"+strcase.ToCamel(row.Name)+"()"
			nextAttributeRead.StateAttributeName = parent.StateAttributeName+"."+strcase.ToCamel(row.Name)
		} else if parent != nil && parent.AttributeType == "ListNested" {
			nextAttributeRead.ParentRead = parent
			nextAttributeRead.GetMethod = "Get"+strcase.ToCamel(row.Name)+"()"
			nextAttributeRead.StateAttributeName = parent.ModelVarName+"."+strcase.ToCamel(row.Name)
			nextAttributeRead.ResultVarName = "value"
		} else {
			nextAttributeRead.GetMethod = "Get"+strcase.ToCamel(row.Name)+"()"
			nextAttributeRead.StateAttributeName = "state."+strcase.ToCamel(row.Name)
			nextAttributeRead.ModelName = dataSourceName+strcase.ToCamel(row.Name)+"DataSourceModel"
		}

		switch {
		case row.Type == "String":
			nextAttributeRead.AttributeType = "String"
		case row.Type == "Guid":
			nextAttributeRead.AttributeType = "DateTimeOffset"
		case row.Type == "Boolean":
			nextAttributeRead.AttributeType = "Boolean"
		case row.Type == "String collection":
			nextAttributeRead.AttributeType = "StringCollection"
		case row.Type == "Guid collection":
			nextAttributeRead.AttributeType = "GuidCollection"
		case row.Type == "DateTimeOffset":
			nextAttributeRead.AttributeType = "DateTimeOffset"
		case strings.HasSuffix(row.Type, "collection"):
			nextAttributeRead.AttributeType = "ListNested"

			nestedCsv := openCsv("template/input/" + packageName + "/" + strcase.ToSnake(row.Name) + ".csv")
			var nestedRead []attributeRead
			generateRead(&nestedRead, nestedCsv, &nextAttributeRead)

			nextAttributeRead.NestedRead = nestedRead
		default:
			nextAttributeRead.AttributeType = "SingleNested"

			nestedCsv := openCsv("template/input/" + packageName + "/" + strcase.ToSnake(row.Name) + ".csv")
			var nestedRead []attributeRead
			generateRead(&nestedRead, nestedCsv, &nextAttributeRead)

			nextAttributeRead.NestedRead = nestedRead
		}

		*read = append(*read, nextAttributeRead)
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

	// Generate schema values from CSV
	var schema []attributeSchema
	generateSchema(&schema, csv)

	// Generate model values from CSV
	var model []attributeModel
	generateModel(strcase.ToLowerCamel(dataSourceName)+"DataSourceModel", &model, csv)

	// Generate schema values from CSV
	var read []attributeRead
	generateRead(&read, csv, nil)
	preRead, err := os.ReadFile("template/input/"+packageName+"/pre_read.go")

	// Set input values to top level template
	inputValues := templateInput{
		PackageName:              packageName,
		DataSourceName:           dataSourceName,
		DataSourceNameUpperCamel: strcase.ToCamel(dataSourceName),
		DataSourceNameLowerCamel: strcase.ToLowerCamel(dataSourceName),
		DataSourceAttributeName:  strcase.ToSnake(dataSourceName),
		Schema:                   schema,
		Model:                    model,
		PreRead:                  string(preRead),
		Read:                     read,
	}

	os.MkdirAll("template/out/", os.ModePerm)
	outfile, err := os.Create("template/out/" + dataSourceName + "_data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource.Execute(outfile, inputValues)

}
