package main

import (
	//"encoding/csv"
	"fmt"
	"io"
	"os"
	"text/template"
	"encoding/csv"

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
	NameUpperCamel string
	NameLowerCamel string
	NameSnake      string
	TypeSchema     string
	TypeModel      string
	Computed       bool
	Optional       bool
	Required       bool
}

type csvSchema struct {
	Name string `csv:"Property"`
	Type string `csv:"Type"`
	Computed bool `csv:"Computed"`
	Optional bool `csv:"Optional"`
	Required bool `csv:"Required"`
}

func main() {

	// Get template
	t1 := template.New("dataSource")
	templateFile, err := os.ReadFile("template/templates/data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	t1, err = t1.Parse(string(templateFile))

	// Get inputs
	packageName    := os.Args[1]
	dataSourceName := os.Args[2]
	f, err := os.Open("template/input/"+dataSourceName+".csv")
    defer f.Close()
	rawCsv := []*csvSchema{}
	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
            r := csv.NewReader(in)
            r.Comma = '|'
			r.LazyQuotes = true
            return r // Allows use pipe as delimiter
        })
	gocsv.UnmarshalFile(f, &rawCsv)
	fmt.Print(rawCsv)

	// Generate schema values from CSV columns
	var schema []schemaInput
	for _, row := range rawCsv {
		schemaRow := new(schemaInput)
		schemaRow.NameUpperCamel = strcase.ToCamel(row.Name)
		schemaRow.NameLowerCamel = strcase.ToLowerCamel(row.Name)
		schemaRow.NameSnake      = strcase.ToSnake(row.Name)

		// Convert types from MS Graph docs to Go and terraform types
		switch row.Type {
			case "String":
				schemaRow.TypeSchema = "String"
				schemaRow.TypeModel  = "String"
			case "Boolean":
				schemaRow.TypeSchema = "Bool"
				schemaRow.TypeModel  = "Bool"
			case "DateTimeOffset":
				schemaRow.TypeSchema = "String"
				schemaRow.TypeModel  = "String"
			default:
				schemaRow.TypeSchema = "FIXME"
				schemaRow.TypeModel  = "FIXME"
		}

		schemaRow.Computed = row.Computed
		schemaRow.Optional = row.Optional
		schemaRow.Required = row.Required
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
	outfile, err := os.Create("template/out/"+dataSourceName+"_data_source.go")
	if err != nil {
		fmt.Print(err)
	}
	t1.Execute(outfile, inputValues)

}
