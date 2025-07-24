package {{.PackageName}}

import (
	"context"

	{{- if .ReadResponse.IfAttrImportNeeded }}
	"github.com/hashicorp/terraform-plugin-framework/attr"
	{{- end}}
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types"
	{{- if .ReadResponse.IfBasetypesImportNeeded }}
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	{{- end}}

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	{{- if .ReadQuery.MultipleGetMethodParameters }}
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	{{- end}}
	"github.com/microsoftgraph/msgraph-sdk-go/{{.PackageName}}"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ datasource.DataSource = &{{.BlockName.LowerCamel}}DataSource{}
    _ datasource.DataSourceWithConfigure = &{{.BlockName.LowerCamel}}DataSource{}
)

// New{{.BlockName.UpperCamel}}DataSource is a helper function to simplify the provider implementation.
func New{{.BlockName.UpperCamel}}DataSource() datasource.DataSource {
    return &{{.BlockName.LowerCamel}}DataSource{}
}

// {{.BlockName.LowerCamel}}DataSource is the data source implementation.
type {{.BlockName.LowerCamel}}DataSource struct{
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *{{.BlockName.LowerCamel}}DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_{{.BlockName.Snake}}"
}

// Configure adds the provider configured client to the data source.
func (d *{{.BlockName.LowerCamel}}DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *{{.BlockName.LowerCamel}}DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
		Description: "{{- .SchemaDescription }}",
		Attributes: map[string]schema.Attribute{
			{{- template "schema_template.go" .SchemaDataSource}}
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *{{.BlockName.LowerCamel}}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfState{{.BlockName.UpperCamel}} {{.BlockName.LowerCamel}}Model
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfState{{.BlockName.UpperCamel}})...)
	if resp.Diagnostics.HasError() {
		return
	}

	{{ template "read_query_template.go" .ReadQuery}}

	{{ template "read_response_template.go" .ReadResponse}}


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfState{{.BlockName.UpperCamel}})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


}
