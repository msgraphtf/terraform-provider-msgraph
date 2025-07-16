package msgraph

import (
	"context"
	"fmt"
	"os"

	"terraform-provider-msgraph/msgraph/applications"
	"terraform-provider-msgraph/msgraph/devices"
	"terraform-provider-msgraph/msgraph/groups"
	"terraform-provider-msgraph/msgraph/serviceprincipals"
	"terraform-provider-msgraph/msgraph/sites"
	"terraform-provider-msgraph/msgraph/teams"
	"terraform-provider-msgraph/msgraph/users"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"tenant_id": schema.StringAttribute{
				Description: "Azure AD Tenant ID.",
				Optional:    true,
			},
			"client_id": schema.StringAttribute{
				Description: "Service Principal client ID",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "Service Principal client secret",
				Optional:    true,
			},
		},
	}
}

// msgraphProviderModel maps provider schema data to a Go type.
type msgraphProviderModel struct {
	TenantID     types.String `tfsdk:"tenant_id"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

// Configure satisfies the provider.Provider interface for msGraphProvider.
func (p *msGraphProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Provider specific implementation.

	var provider_config msgraphProviderModel
	diags := req.Config.Get(ctx, &provider_config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant_id := os.Getenv("MSGRAPH_TENANT_ID")
	client_id := os.Getenv("MSGRAPH_CLIENT_ID")
	client_secret := os.Getenv("MSGRAPH_CLIENT_SECRET")

	if provider_config.TenantID.ValueString() != "" {
		tenant_id = provider_config.TenantID.ValueString()
	}
	if provider_config.ClientID.ValueString() != "" {
		client_id = provider_config.ClientID.ValueString()
	}
	if provider_config.ClientSecret.ValueString() != "" {
		client_secret = provider_config.ClientSecret.ValueString()
	}

	var cred azcore.TokenCredential
	var err error

	if tenant_id != "" && client_id != "" && client_secret != "" {
		cred, err = azidentity.NewClientSecretCredential(tenant_id, client_id, client_secret, nil)
	} else {
		cred, err = azidentity.NewAzureCLICredential(nil)
	}
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

}

// DataSources satisfies the provider.Provider interface for msGraphProvider.
func (p *msGraphProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// Provider specific implementation
		applications.NewApplicationDataSource,
		applications.NewApplicationsDataSource,
		devices.NewDeviceDataSource,
		devices.NewDevicesDataSource,
		groups.NewGroupDataSource,
		groups.NewGroupsDataSource,
		serviceprincipals.NewServicePrincipalDataSource,
		serviceprincipals.NewServicePrincipalsDataSource,
		sites.NewSiteDataSource,
		sites.NewSitesDataSource,
		teams.NewTeamDataSource,
		users.NewUserDataSource,
		users.NewUsersDataSource,
	}
}

// Resources satisfies the provider.Provider interface for msGraphProvider.
func (p *msGraphProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Provider specific implementation
		users.NewUserResource,
		groups.NewGroupResource,
	}
}
