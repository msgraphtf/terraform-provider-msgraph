package {{.PackageName}}

import (
    "context"

	"github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
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
    resp.TypeName = req.ProviderTypeName + "_{{.DataSourceNameSnake}}"
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
			{{- range .Schema}}
			"{{.NameSnake}}": schema.{{.TypeSchema}}Attribute{
				{{- if .Required}}
				Required: true,
				{{- else if .Optional}}
				Optional: true,
				{{- else if .Computed}}
				Computed: true,
				{{- end}}
			},{{end}}
		},
	}
}

type {{.DataSourceNameLowerCamel}}DataSourceModel struct {
	{{- range .Schema}}
	{{.NameUpperCamel}} types.{{.TypeModel}}
	{{- end}}
}

// Read refreshes the Terraform state with the latest data.
func (d *{{.DataSourceNameLowerCamel}}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state {{.DataSourceNameLowerCamel}}DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

}
