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
    _ datasource.DataSource = &{{.DataSourceName.LowerCamel}}DataSource{}
    _ datasource.DataSourceWithConfigure = &{{.DataSourceName.LowerCamel}}DataSource{}
)

// New{{.DataSourceName.UpperCamel}}DataSource is a helper function to simplify the provider implementation.
func New{{.DataSourceName.UpperCamel}}DataSource() datasource.DataSource {
    return &{{.DataSourceName.LowerCamel}}DataSource{}
}

// {{.DataSourceName.LowerCamel}}DataSource is the data source implementation.
type {{.DataSourceName.LowerCamel}}DataSource struct{
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *{{.DataSourceName.LowerCamel}}DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_{{.DataSourceName.Snake}}"
}

// Configure adds the provider configured client to the data source.
func (d *{{.DataSourceName.LowerCamel}}DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *{{.DataSourceName.LowerCamel}}DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			{{- if eq .AttributeType "StringAttribute" }}
			{{- template "StringAttribute" .}}
			{{- else if eq .AttributeType "Int64Attribute" }}
			{{- template "Int64Attribute" .}}
			{{- else if eq .AttributeType "BoolAttribute" }}
			{{- template "BoolAttribute" .}}
			{{- else if eq .AttributeType "ListAttribute" }}
			{{- template "ListAttribute" .}}
			{{- else if eq .AttributeType "SingleNestedAttribute" }}
			{{- template "SingleNestedAttribute" .}}
			{{- else if eq .AttributeType "ListNestedAttribute" }}
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
func (d *{{.DataSourceName.LowerCamel}}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state {{.DataSourceName.LowerCamel}}DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	{{.PreRead}}

	{{- /* Define templates for mapping each response type to state */}}
	{{- define "ReadStringAttribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.StringValue(*{{.GetMethod}})}
	{{- end}}

	{{- define "ReadStringFormattedAttribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.StringValue({{.GetMethod}}.String())}
	{{- end}}

	{{- define "ReadInt64Attribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.Int64Value(int64(*{{.GetMethod}}))}
	{{- end}}

	{{- define "ReadBoolAttribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.BoolValue(*{{.GetMethod}})}
	{{- end}}

	{{- define "ReadSingleNestedAttribute" }}
	if {{.GetMethod}} != nil {
		{{.StateVarName}} = new({{.ModelName}})
		{{template "generate_read" .NestedRead}}
	}
	{{- end}}

	{{- define "ReadListStringAttribute" }}
	for _, value := range {{.GetMethod}} {
		{{.StateVarName}} = append({{.StateVarName}}, types.StringValue(value))
	}
	{{- end}}

	{{- define "ReadListStringFormattedAttribute" }}
	for _, value := range {{.GetMethod}} {
		{{.StateVarName}} = append({{.StateVarName}}, types.StringValue(value.String()))
	}
	{{- end}}

	{{- define "ReadListNestedAttribute" }}
	for _, value := range {{.GetMethod}} {
		{{.ModelVarName}} := new({{.ModelName}})
			{{template "generate_read" .NestedRead}}
		{{.StateVarName}} = append({{.StateVarName}}, *{{.ModelVarName}})
	}
	{{- end}}


	{{/* Generate statements to map response to state */}}
	{{- block "generate_read" .Read}}
	{{- range .}}
	{{- if eq .AttributeType "ReadStringAttribute"}}
	{{- template "ReadStringAttribute" .}}
	{{- else if eq .AttributeType "ReadStringFormattedAttribute"}}
	{{- template "ReadStringFormattedAttribute" .}}
	{{- else if eq .AttributeType "ReadInt64Attribute"}}
	{{- template "ReadInt64Attribute" .}}
	{{- else if eq .AttributeType "ReadBoolAttribute"}}
	{{- template "ReadBoolAttribute" .}}
	{{- else if eq .AttributeType "ReadListStringAttribute"}}
	{{- template "ReadListStringAttribute" .}}
	{{- else if eq .AttributeType "ReadListStringFormattedAttribute"}}
	{{- template "ReadListStringFormattedAttribute" .}}
	{{- else if eq .AttributeType "ReadSingleNestedAttribute"}}
	{{- template "ReadSingleNestedAttribute" .}}
	{{- else if eq .AttributeType "ReadListNestedAttribute"}}
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
