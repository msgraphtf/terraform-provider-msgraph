package users

import (
	"context"
	//"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	//"github.com/microsoftgraph/msgraph-sdk-go/models"
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
			"about_me": schema.StringAttribute{
				Computed: true,
			},
			"account_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"age_group": schema.StringAttribute{
				Computed: true,
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"disabled_plans": schema.ListAttribute{
							Computed: true,
							ElementType: types.StringType,
						},
						"skus": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"assigned_plans": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assigned_date_time": schema.StringAttribute{
							Computed: true,
						},
						"capability_status": schema.StringAttribute{
							Computed: true,
						},
						"service": schema.StringAttribute{
							Computed: true,
						},
						"service_plan_id": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"birthday": schema.StringAttribute{
				Computed: true,
			},
			"business_phones": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"mail_nickname": schema.StringAttribute{
				Computed: true,
				//TODO: Optional: true,
			},
			"password_profile": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"force_change_password_next_sign_in": schema.BoolAttribute{
						Computed: true,
					},
					"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
						Computed: true,
					},
					"password": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"user_principal_name": schema.StringAttribute{
				//TODO: Optional: true,
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
	AboutMe           types.String   `tfsdk:"about_me"`
	AccountEnabled    types.Bool     `tfsdk:"account_enabled"`
	AgeGroup          types.String   `tfsdk:"age_group"`
	AssignedLicenses  []userDataSourceAssignedLicenseModel `tfsdk:"assigned_licenses"`
	AssignedPlans     []userDataSourceAssignedPlanModel `tfsdk:"assigned_plans"`
	Birthday          types.String   `tfsdk:"birthday"`
	BusinessPhones    []types.String `tfsdk:"business_phones"`
	DisplayName       types.String   `tfsdk:"display_name"`
	Id                types.String   `tfsdk:"id"`
	MailNickname      types.String   `tfsdk:"mail_nickname"`
	PasswordProfile   *userDataSourcePasswordProfileModel `tfsdk:"password_profile"`
	UserPrincipalName types.String   `tfsdk:"user_principal_name"`
}

type userDataSourceAssignedLicenseModel struct {
	DisabledPlans []types.String `tfsdk:"disabled_plans"`
	Sku           types.String `tfsdk:"skus"`
}

type userDataSourceAssignedPlanModel struct {
	AssignedDateTime types.String `tfsdk:"assigned_date_time"`
	CapabilityStatus types.String `tfsdk:"capability_status"`
	Service          types.String `tfsdk:"service"`
	ServicePlanID    types.String `tfsdk:"service_plan_id"`
}

type userDataSourcePasswordProfileModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
	Password                             types.String `tfsdk:"password"`
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
			Select: UserProperties[:],
		},
	}

	// TODO: Allow using UserPrincipalName or MailNickname as ID to search for
	result, err := d.client.UsersById(state.Id.ValueString()).Get(context.Background(), &qparams)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting user",
			err.Error(),
		)
		return
	}


	// Map response to model
	if result.GetAboutMe()  != nil {state.AboutMe = types.StringValue(*result.GetAboutMe())}
	state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	if result.GetAgeGroup() != nil {state.AboutMe = types.StringValue(*result.GetAgeGroup())}

	// Map assigned licenses
	for _, license := range result.GetAssignedLicenses(){
		assignedLicensesState := userDataSourceAssignedLicenseModel{
			Sku: types.StringValue(license.GetSkuId().String()),
		}
		for _, disabledLicense := range license.GetDisabledPlans() {
			assignedLicensesState.DisabledPlans = append(assignedLicensesState.DisabledPlans, types.StringValue(disabledLicense.String()))
		}
		state.AssignedLicenses = append(state.AssignedLicenses, assignedLicensesState)
	}

	// Map assigned plans
	for _, plan := range result.GetAssignedPlans(){
		assignedPlansState := userDataSourceAssignedPlanModel{
			AssignedDateTime: types.StringValue(plan.GetAssignedDateTime().String()),
			CapabilityStatus: types.StringValue(*plan.GetCapabilityStatus()),
			Service:          types.StringValue(*plan.GetService()),
			ServicePlanID:    types.StringValue(plan.GetServicePlanId().String()),
		}
		state.AssignedPlans = append(state.AssignedPlans, assignedPlansState)
	}

	if result.GetBirthday() != nil {state.Birthday = types.StringValue(result.GetBirthday().String())}

	for _, businessPhone := range result.GetBusinessPhones() {
		state.BusinessPhones = append(state.BusinessPhones, types.StringValue(businessPhone))
	}

	state.DisplayName       = types.StringValue(*result.GetDisplayName())
	state.MailNickname      = types.StringValue(*result.GetMailNickname())
	state.Id                = types.StringValue(*result.GetId())

	// Map password profile
	passwordProfile := new(userDataSourcePasswordProfileModel)
	passwordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
	passwordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
	state.PasswordProfile = passwordProfile

	state.UserPrincipalName = types.StringValue(*result.GetUserPrincipalName())


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
