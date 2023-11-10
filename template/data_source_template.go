package {{.PackageName}}

import (
    "context"

	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/{{.PackageName}}"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ datasource.DataSource = &{{.DataSourceNameLowerCamel}}DataSource{}
    _ datasource.DataSourceWithConfigure = &{{.DataSourceNameLowerCamel}}DataSource{}
)

// New{{.DataSourceNameUpperCamel}}DataSource is a helper function to simplify the provider implementation.
func New{{.DataSourceNameUpperCamel}}DataSource() datasource.DataSource {
    return &{{.DataSourceNameLowerCamel}}DataSource{}
}

// {{.DataSourceNameUpperCamel}}DataSource is the data source implementation.
type {{.DataSourceName}}DataSource struct{
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *{{.DataSourceNameLowerCamel}}DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_{{.DataSourceAttributeName}}"
}

// Configure adds the provider configured client to the data source.
func (d *{{.DataSourceNameLowerCamel}}DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *{{.DataSourceNameLowerCamel}}DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{

			{{- /* Define templates for different Attribute types */}}
			{{- define "StringAttribute" }}
			"{{.AttributeName}}": schema.StringAttribute{
				Description: "{{.Description}}",
				{{- if .Required}}
				Required: true,
				{{- end}}
				{{- if .Optional}}
				Optional: true,
				{{- end}}
				{{- if .Computed}}
				Computed: true,
				{{- end}}
			},
			{{- end }}

			{{- define "Int64Attribute" }}
			"{{.AttributeName}}": schema.Int64Attribute{
				Description: "{{.Description}}",
				{{- if .Required}}
				Required: true,
				{{- end}}
				{{- if .Optional}}
				Optional: true,
				{{- end}}
				{{- if .Computed}}
				Computed: true,
				{{- end}}
			},
			{{- end }}

			{{- define "BoolAttribute" }}
			"{{.AttributeName}}": schema.BoolAttribute{
				Description: "{{.Description}}",
				{{- if .Required}}
				Required: true,
				{{- end}}
				{{- if .Optional}}
				Optional: true,
				{{- end}}
				{{- if .Computed}}
				Computed: true,
				{{- end}}
			},
			{{- end }}

			{{- define "ListAttribute" }}
			"{{.AttributeName}}": schema.ListAttribute{
				Description: "{{.Description}}",
				{{- if .Required}}
				Required: true,
				{{- end}}
				{{- if .Optional}}
				Optional: true,
				{{- end}}
				{{- if .Computed}}
				Computed: true,
				{{- end}}
				ElementType: {{.ElementType}},
			},
			{{- end }}

			{{- define "SingleNestedAttribute" }}
			"{{.AttributeName}}": schema.SingleNestedAttribute{
				Description: "{{.Description}}",
				{{- if .Required}}
				Required: true,
				{{- end}}
				{{- if .Optional}}
				Optional: true,
				{{- end}}
				{{- if .Computed}}
				Computed: true,
				{{- end}}
				Attributes: map[string]schema.Attribute{
				{{- template "generate_schema" .Attributes}}
				},
			},
			{{- end }}

			{{- define "ListNestedAttribute" }}
			"{{.AttributeName}}": schema.ListNestedAttribute{
				Description: "{{.Description}}",
				{{- if .Required}}
				Required: true,
				{{- end}}
				{{- if .Optional}}
				Optional: true,
				{{- end}}
				{{- if .Computed}}
				Computed: true,
				{{- end}}
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						{{- template "generate_schema" .NestedObject}}
					},
				},
			},
			{{- end }}

			{{- /* Generate our Attributes from our defined templates above */}}
			{{- block "generate_schema" .Schema}}
			{{- range .}}
			{{- if eq .AttributeType "String" }}
			{{- template "StringAttribute" .}}
			{{- else if eq .AttributeType "Integer" }}
			{{- template "Int64Attribute" .}}
			{{- else if eq .AttributeType "Bool" }}
			{{- template "BoolAttribute" .}}
			{{- else if eq .AttributeType "List" }}
			{{- template "ListAttribute" .}}
			{{- else if eq .AttributeType "SingleNested" }}
			{{- template "SingleNestedAttribute" .}}
			{{- else if eq .AttributeType "ListNested" }}
			{{- template "ListNestedAttribute" .}}
			{{- end }}
			{{- end}}
			{{- end}}
		},
	}
}

{{/* Generate data models from provided data */}}
{{- range .Model}}
type {{.ModelName}} struct {
{{- range .Fields}}
{{.FieldName}} {{.FieldType}} `tfsdk:"{{.AttributeName}}"`
{{- end}}
}
{{end}}

// Read refreshes the Terraform state with the latest data.
func (d *{{.DataSourceNameLowerCamel}}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state {{.DataSourceNameLowerCamel}}DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	{{.PreRead}}

	{{- /* Define templates for mapping each response type to state */}}
	{{- define "ReadStringAttribute" }}
	if {{.ResultVarName}}.{{.GetMethod}}  != nil { {{- .StateAttributeName}} = types.StringValue(*{{.ResultVarName}}.{{.GetMethod}})}
	{{- end}}

	{{- define "ReadStringFormatted" }}
	if {{.ResultVarName}}.{{.GetMethod}}  != nil { {{- .StateAttributeName}} = types.StringValue({{.ResultVarName}}.{{.GetMethod}}.String())}
	{{- end}}

	{{- define "ReadInt64Attribute" }}
	if {{.ResultVarName}}.{{.GetMethod}}  != nil { {{- .StateAttributeName}} = types.Int64Value(int64(*{{.ResultVarName}}.{{.GetMethod}}))}
	{{- end}}

	{{- define "ReadBooleanAttribute" }}
	if {{.ResultVarName}}.{{.GetMethod}}  != nil { {{- .StateAttributeName}} = types.BoolValue(*{{.ResultVarName}}.{{.GetMethod}})}
	{{- end}}

	{{- define "ReadStringCollection" }}
	for _, value := range {{.ResultVarName}}.{{.GetMethod}} {
		{{.StateAttributeName}}= append({{.StateAttributeName}}, types.StringValue(value))
	}
	{{- end}}

	{{- define "ReadStringFormattedCollection" }}
	for _, value := range {{.ResultVarName}}.{{.GetMethod}} {
		{{.StateAttributeName}} = append({{.StateAttributeName}}, types.StringValue(value.String()))
	}
	{{- end}}

	{{- define "ReadSingleNestedAttribute" }}
	{{.ModelVarName}} := new({{.ModelName}})
	if {{.ResultVarName}}.{{.GetMethod}} != nil {
		{{template "generate_read" .NestedRead}}
	}
	{{.StateAttributeName}} = {{.ModelVarName}}
	{{- end}}

	{{- define "ReadListNestedAttribute" }}
	for _, value := range {{.ResultVarName}}.{{.GetMethod}} {
		{{.ModelVarName}} := new({{.ModelName}})
			{{template "generate_read" .NestedRead}}
		{{.StateAttributeName}} = append({{.StateAttributeName}}, *{{.ModelVarName}})
	}
	{{- end}}


	{{/* Generate statements to map response to state */}}
	{{- block "generate_read" .Read}}
	{{- range .}}
	{{- if eq .AttributeType "String"}}
	{{- template "ReadStringAttribute" .}}
	{{- else if eq .AttributeType "StringFormatted"}}
	{{- template "ReadStringFormatted" .}}
	{{- else if eq .AttributeType "Integer"}}
	{{- template "ReadInt64Attribute" .}}
	{{- else if eq .AttributeType "Boolean"}}
	{{- template "ReadBooleanAttribute" .}}
	{{- else if eq .AttributeType "StringCollection"}}
	{{- template "ReadStringCollection" .}}
	{{- else if eq .AttributeType "StringFormattedCollection"}}
	{{- template "ReadStringFormattedCollection" .}}
	{{- else if eq .AttributeType "SingleNested"}}
	{{- template "ReadSingleNestedAttribute" .}}
	{{- else if eq .AttributeType "ListNested"}}
	{{- template "ReadListNestedAttribute" .}}
	{{- end}}
	{{- end}}
	{{- end}}


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


}
