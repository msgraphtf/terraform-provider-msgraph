package users

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

type userDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"account_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"mail_nickname": schema.StringAttribute{
				Computed: true,
				Optional: true,
			},
			"password_profile": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"force_change_password_next_sign_in": schema.BoolAttribute{
						Computed: true,
					},
					"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
						Computed: true,
					},
				},
			},
			"user_principal_name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

type userDataSourceModel struct {
	AccountEnabled    types.Bool               `tfsdk:"account_enabled"`
	DisplayName       types.String             `tfsdk:"display_name"`
	MailNickname      types.String             `tfsdk:"mail_nickname"`
	PasswordProfile   *userDataSourcePasswordProfileModel `tfsdk:"password_profile"`
	UserPrincipalName types.String             `tfsdk:"user_principal_name"`
	Id                types.String             `tfsdk:"id"`
}

type userDataSourcePasswordProfileModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
}

type userModel struct {
	DisplayName types.String `tfsdk:"display_name"`
}

// Read refreshes the Terraform state with the latest data.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: []string{"accountEnabled, displayName, mailNickname, passwordProfile, userPrincipalName, Id"},
		},
	}

	result, err := d.client.UsersById(state.Id.ValueString()).Get(context.Background(), &qparams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting user",
			err.Error(),
		)
		return
	}

	tflog.Info(ctx, "PASSWORD: "+strconv.FormatBool(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn()))

	// Map response to model
	state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	state.DisplayName = types.StringValue(*result.GetDisplayName())
	state.MailNickname = types.StringValue(*result.GetMailNickname())

	passwordProfile := new(userDataSourcePasswordProfileModel)
	passwordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
	passwordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
	state.PasswordProfile = passwordProfile

	state.UserPrincipalName = types.StringValue(*result.GetUserPrincipalName())
	state.Id = types.StringValue(*result.GetId())


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
