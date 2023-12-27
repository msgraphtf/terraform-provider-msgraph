package {{.PackageName}}

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	{{- if gt .ReadQueryGetMethodParametersCount 0 }}
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	{{- end}}
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
