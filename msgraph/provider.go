package msgraph

import (
	"context"
	"encoding/base64"
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

// Ensure the implementation satisfies various provider interfaces.
var (
	_ provider.Provider = &MsGraphProvider{}
)

// MsGraphProvider defines the provider implementation
type MsGraphProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// msgraphProviderModel describes the provider data model.
type msgraphProviderModel struct {
	TenantID                  types.String `tfsdk:"tenant_id"`
	ClientID                  types.String `tfsdk:"client_id"`
	ClientSecret              types.String `tfsdk:"client_secret"`
	ClientCertificate         types.String `tfsdk:"client_certificate"`
	ClientCertificatePath     types.String `tfsdk:"client_certificate_path"`
	ClientCertificatePassword types.String `tfsdk:"client_certificate_password"`
}

func (p *MsGraphProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "msgraph"
	resp.Version = p.version
}

func (p *MsGraphProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
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
			"client_certificate": schema.StringAttribute{
				Description: "Service Principal client certificate",
				Optional:    true,
			},
			"client_certificate_path": schema.StringAttribute{
				Description: "Service Principal client certificate path",
				Optional:    true,
			},
			"client_certificate_password": schema.StringAttribute{
				Description: "Service Principal client certificate password",
				Optional:    true,
			},
		},
	}
}

func (p *MsGraphProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
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
	client_certificate := os.Getenv("MSGRAPH_CLIENT_CERTIFICATE")
	client_certificate_path := os.Getenv("MSGRAPH_CLIENT_CERTIFICATE_PATH")
	client_certificate_password := os.Getenv("MSGRAPH_CLIENT_CERTIFICATE_PASSWORD")

	if provider_config.TenantID.ValueString() != "" {
		tenant_id = provider_config.TenantID.ValueString()
	}
	if provider_config.ClientID.ValueString() != "" {
		client_id = provider_config.ClientID.ValueString()
	}
	if provider_config.ClientSecret.ValueString() != "" {
		client_secret = provider_config.ClientSecret.ValueString()
	}
	if provider_config.ClientCertificate.ValueString() != "" {
		client_certificate = provider_config.ClientCertificate.ValueString()
	}
	if provider_config.ClientCertificatePath.ValueString() != "" {
		client_certificate_path = provider_config.ClientCertificatePath.ValueString()
	}
	if provider_config.ClientCertificatePassword.ValueString() != "" {
		client_certificate_password = provider_config.ClientCertificatePassword.ValueString()
	}

	var cred azcore.TokenCredential
	var err error

	if tenant_id != "" && client_id != "" && client_secret != "" {
		cred, err = azidentity.NewClientSecretCredential(tenant_id, client_id, client_secret, nil)
	} else if tenant_id != "" && client_id != "" && client_certificate != "" && client_certificate_password != "" {
		decoded_client_certificate, err := base64.StdEncoding.DecodeString(client_certificate)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error decoding certificate",
				err.Error(),
			)
		}
		certificate, private_key, err := azidentity.ParseCertificates(decoded_client_certificate, []byte(client_certificate_password))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing certificate",
				err.Error(),
			)
		}
		cred, err = azidentity.NewClientCertificateCredential(tenant_id, client_id, certificate, private_key, nil)
	} else if tenant_id != "" && client_id != "" && client_certificate_path != "" && client_certificate_password != "" {
		certificate_file, err := os.ReadFile(client_certificate_path)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error opening certificate file",
				err.Error(),
			)
		}
		certificate, private_key, err := azidentity.ParseCertificates(certificate_file, []byte(client_certificate_password))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing certificate",
				err.Error(),
			)
		}
		cred, err = azidentity.NewClientCertificateCredential(tenant_id, client_id, certificate, private_key, nil)
	} else {
		cred, err = azidentity.NewAzureCLICredential(nil)
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting credential",
			err.Error(),
		)
	}

	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, []string{"https://graph.microsoft.com/.default"})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting client",
			err.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		fmt.Printf("Error")
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *MsGraphProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
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

func (p *MsGraphProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// Provider specific implementation
		applications.NewApplicationResource,
		devices.NewDeviceResource,
		groups.NewGroupResource,
		serviceprincipals.NewServicePrincipalResource,
		teams.NewTeamResource,
		users.NewUserResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &MsGraphProvider{
			version: version,
		}
	}
}
