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
	Schema                   []schemaInput
}

type schemaInput struct {
	NameUpperCamel      string
	NameLowerCamel      string
	NameSnake           string
	TypeModel           string
	TypeSchema          string
	MarkdownDescription string
	Required            bool
	Optional            bool
	Computed            bool
	ElementType         string
	NestedObject        bool
}

type csvSchema struct {
	Name        string `csv:"Property"`
	Type        string `csv:"Type"`
	Computed    bool   `csv:"Computed"`
	Optional    bool   `csv:"Optional"`
	Required    bool   `csv:"Required"`
	Description string `csv:"Description"`
}

func main() {

	// Get template
	templateDataSource := template.New("dataSource")
	templateFile, err := os.ReadFile("template/templates/data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	templateDataSource, err = templateDataSource.Parse(string(templateFile))

	// Get inputs
	packageName := os.Args[1]
	dataSourceName := os.Args[2]

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
	var schema []schemaInput
	for _, row := range rawCsv {
		schemaRow := new(schemaInput)
		schemaRow.NameUpperCamel = strcase.ToCamel(row.Name)
		schemaRow.NameLowerCamel = strcase.ToLowerCamel(row.Name)
		schemaRow.NameSnake = strcase.ToSnake(row.Name)

		// Convert types from MS Graph docs to Go and terraform types
		switch {
		case row.Type == "String":
			schemaRow.TypeSchema = "String"
			schemaRow.TypeModel = "types.String"
		case row.Type == "String collection":
			schemaRow.TypeSchema = "List"
			schemaRow.TypeModel = "[]types.String"
			schemaRow.ElementType = "types.StringType"
		case row.Type == "Boolean":
			schemaRow.TypeSchema = "Bool"
			schemaRow.TypeModel = "types.Bool"
		case row.Type == "DateTimeOffset":
			schemaRow.TypeSchema = "String"
			schemaRow.TypeModel = "types.String"
		case strings.HasSuffix(row.Type, "collection"):
			schemaRow.TypeSchema = "ListNested"
			schemaRow.TypeModel = "[]" + dataSourceName + "DataSource" + strcase.ToCamel(row.Type)
			schemaRow.NestedObject = true
		default:
			schemaRow.TypeSchema = "SingleNested"
			schemaRow.TypeModel = "*" + dataSourceName + "DataSource" + strcase.ToCamel(row.Type)
		}

		schemaRow.Computed = row.Computed
		schemaRow.Optional = row.Optional
		schemaRow.Required = row.Required
		schemaRow.MarkdownDescription = row.Description
		schema = append(schema, *schemaRow)
	}

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
