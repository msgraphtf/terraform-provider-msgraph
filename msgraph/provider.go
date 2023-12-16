package msgraph

import (
	"context"
	"fmt"

	"terraform-provider-msgraph/msgraph/users"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the provider.Provider interface.
var (
	_ provider.Provider = &msGraphProvider{}
)

func New() provider.Provider {
	return &msGraphProvider{}
}

type msGraphProvider struct{}

// Metadata satisfies the provider.Provider interface for msGraphProvider
func (p *msGraphProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "msgraph"
}

// Schema satisfies the provider.Provider interface for msGraphProvider.
func (p *msGraphProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"tennant_id": schema.StringAttribute{
				Description: "Azure AD Tenant ID.",
				Optional:    true,
			},
		},
	}
}

// msgraphProviderModel maps provider schema data to a Go type.
type msgraphProviderModel struct {
	TennantID types.String `tfsdk:"tennant_id"`
}

// Configure satisfies the provider.Provider interface for msGraphProvider.
func (p *msGraphProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Provider specific implementation.

	tflog.Info(ctx, "Configuring MS Graph client")

	var config msgraphProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting credential",
			"Error getting credential",
		)
	}

	if resp.Diagnostics.HasError() {
		fmt.Printf("Error")
		return
	}
	tflog.Info(ctx, "Creating MS Graph client")

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting client",
			"Error getting client",
		)
	}

	if resp.Diagnostics.HasError() {
		fmt.Printf("Error")
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured MS Graph client", map[string]any{"success": true})
}

// DataSources satisfies the provider.Provider interface for msGraphProvider.
func (p *msGraphProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Provider specific implementation
		users.NewUserDataSource,
		users.NewUsersDataSource,
	}
}

// Resources satisfies the provider.Provider interface for msGraphProvider.
func (p *msGraphProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Provider specific implementation
		users.NewUserResource,
	}
}
