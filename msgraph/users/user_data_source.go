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
			"city": schema.StringAttribute{
				Computed: true,
			},
			"company_name": schema.StringAttribute{
				Computed: true,
			},
			"consent_provided_for_minor": schema.StringAttribute{
				Computed: true,
			},
			"country": schema.StringAttribute{
				Computed: true,
			},
			"creation_type": schema.StringAttribute{
				Computed: true,
			},
			"deleted_date_time": schema.StringAttribute{
				Computed: true,
			},
			"created_date_time": schema.StringAttribute{
				Computed: true,
			},
			"department": schema.StringAttribute{
				Computed: true,
			},
			"employee_hire_date": schema.StringAttribute{
				Computed: true,
			},
			"employee_id": schema.StringAttribute{
				Computed: true,
			},
			"employee_leave_date_time": schema.StringAttribute{
				Computed: true,
			},
			"employee_org_data": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"cost_center": schema.StringAttribute{
						Computed: true,
					},
					"division": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"employee_type": schema.StringAttribute{
				Computed: true,
			},
			"external_user_state": schema.StringAttribute{
				Computed: true,
			},
			"external_user_state_change_date_time": schema.StringAttribute{
				Computed: true,
			},
			"fax_number": schema.StringAttribute{
				Computed: true,
			},
			"given_name": schema.StringAttribute{
				Computed: true,
			},
			"hire_date": schema.StringAttribute{
				Computed: true,
			},
			"display_name": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"identities": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"issuer": schema.StringAttribute{
							Computed: true,
						},
						"issuer_assigned_id": schema.StringAttribute{
							Computed: true,
						},
						"sign_in_type": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"im_addresses": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
			},
			"interests": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
			},
			"is_resource_account": schema.BoolAttribute{
				Computed: true,
			},
			"job_title": schema.StringAttribute{
				Computed: true,
			},
			"last_password_change_date_time": schema.StringAttribute{
				Computed: true,
			},
			"legal_age_group_classification": schema.StringAttribute{
				Computed: true,
			},
			"license_assignment_states": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assigned_by_group": schema.StringAttribute{
							Computed: true,
						},
						"disabled_plans": schema.ListAttribute{
							Computed: true,
							ElementType: types.StringType,
						},
						"error": schema.StringAttribute{
							Computed: true,
						},
						"last_updated_date_time": schema.StringAttribute{
							Computed: true,
						},
						"sku_id": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"mail": schema.StringAttribute{
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
	AboutMe                         types.String                         `tfsdk:"about_me"`
	AccountEnabled                  types.Bool                           `tfsdk:"account_enabled"`
	AgeGroup                        types.String                         `tfsdk:"age_group"`
	AssignedLicenses                []userDataSourceAssignedLicenseModel `tfsdk:"assigned_licenses"`
	AssignedPlans                   []userDataSourceAssignedPlanModel    `tfsdk:"assigned_plans"`
	Birthday                        types.String                         `tfsdk:"birthday"`
	BusinessPhones                  []types.String                       `tfsdk:"business_phones"`
	City                            types.String                         `tfsdk:"city"`
	CompanyName                     types.String                         `tfsdk:"company_name"`
	ConsentProvidedForMinor         types.String                         `tfsdk:"consent_provided_for_minor"`
	Country                         types.String                         `tfsdk:"country"`
	CreatedDateTime                 types.String                         `tfsdk:"created_date_time"`
	CreationType                    types.String                         `tfsdk:"creation_type"`
	DeletedDateTime                 types.String                         `tfsdk:"deleted_date_time"`
	Department                      types.String                         `tfsdk:"department"`
	DisplayName                     types.String                         `tfsdk:"display_name"`
	EmployeeHireDate                types.String                         `tfsdk:"employee_hire_date"`
	EmployeeId                      types.String                         `tfsdk:"employee_id"`
	EmployeeLeaveDateTime           types.String                         `tfsdk:"employee_leave_date_time"`
	EmployeeOrgData                 *userDataSourceEmployeeOrgData       `tfsdk:"employee_org_data"`
	EmployeeType                    types.String                         `tfsdk:"employee_type"`
	ExternalUserState               types.String                         `tfsdk:"external_user_state"`
	ExternalUserStateChangeDateTime types.String                         `tfsdk:"external_user_state_change_date_time"`
	FaxNumber                       types.String                         `tfsdk:"fax_number"`
	GivenName                       types.String                         `tfsdk:"given_name"`
	HireDate                        types.String                         `tfsdk:"hire_date"`
	Id                              types.String                         `tfsdk:"id"`
	Identities                      []userDataSourceIdentities           `tfsdk:"identities"`
	ImAddresses                     []types.String                       `tfsdk:"im_addresses"`
	Interests                       []types.String                       `tfsdk:"interests"`
	IsResourceAccount               types.Bool                           `tfsdk:"is_resource_account"`
	JobTitle                        types.String                         `tfsdk:"job_title"`
	LegalAgeGroupClassification     types.String                         `tfsdk:"legal_age_group_classification"`
	LastPasswordChangeDateTime      types.String                         `tfsdk:"last_password_change_date_time"`
	LicenseAssignmentStates         []userDataSourceLicenseAssignmentStatesModel `tfsdk:"license_assignment_states"`
	Mail                            types.String                         `tfsdk:"mail"`
	MailNickname                    types.String                         `tfsdk:"mail_nickname"`
	PasswordProfile                 *userDataSourcePasswordProfileModel  `tfsdk:"password_profile"`
	UserPrincipalName               types.String                         `tfsdk:"user_principal_name"`
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

type userDataSourceEmployeeOrgData struct {
	CostCenter types.String `tfsdk:"cost_center"`
	Division   types.String `tfsdk:"division"`
}

type userDataSourceIdentities struct {
	Issuer           types.String `tfsdk:"issuer"`
	IssuerAssignedId types.String `tfsdk:"issuer_assigned_id"`
	SignInType       types.String `tfsdk:"sign_in_type"`
}

type userDataSourcePasswordProfileModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
	Password                             types.String `tfsdk:"password"`
}

type userDataSourceLicenseAssignmentStatesModel struct {
	AssignedByGroup     types.String   `tfsdk:"assigned_by_group"`
	DisabledPlans       []types.String `tfsdk:"disabled_plans"` // WARNING: Is this a issue being a duplicated attribute name? 
	Error               types.String   `tfsdk:"error"`
	LastUpdatedDateTime types.String   `tfsdk:"last_updated_date_time"`
	SkuId               types.String   `tfsdk:"sku_id"`
	State               types.String   `tfsdk:"state"`
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
	for _, license := range result.GetAssignedLicenses() {
		assignedLicensesState := userDataSourceAssignedLicenseModel{
			Sku: types.StringValue(license.GetSkuId().String()),
		}
		for _, disabledLicense := range license.GetDisabledPlans() {
			assignedLicensesState.DisabledPlans = append(assignedLicensesState.DisabledPlans, types.StringValue(disabledLicense.String()))
		}
		state.AssignedLicenses = append(state.AssignedLicenses, assignedLicensesState)
	}

	// Map assigned plans
	for _, plan := range result.GetAssignedPlans() {
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

	if result.GetCity()                    != nil {state.City                    = types.StringValue(*result.GetCity())}
	if result.GetCompanyName()             != nil {state.CompanyName             = types.StringValue(*result.GetCompanyName())}
	if result.GetConsentProvidedForMinor() != nil {state.ConsentProvidedForMinor = types.StringValue(*result.GetConsentProvidedForMinor())}
	state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	if result.GetCreationType()            != nil {state.CreationType            = types.StringValue(*result.GetCreationType())}
	if result.GetDeletedDateTime()         != nil {state.DeletedDateTime         = types.StringValue(result.GetDeletedDateTime().String())}
	if result.GetDepartment()              != nil {state.Department              = types.StringValue(*result.GetDepartment())}
	state.DisplayName       = types.StringValue(*result.GetDisplayName())
	if result.GetEmployeeHireDate()        != nil {state.EmployeeHireDate        = types.StringValue(result.GetEmployeeHireDate().String())}
	if result.GetEmployeeId()              != nil {state.EmployeeId              = types.StringValue(*result.GetEmployeeId())}
	if result.GetEmployeeLeaveDateTime()   != nil {state.EmployeeLeaveDateTime   = types.StringValue(result.GetEmployeeLeaveDateTime().String())}

	employeeOrgData := new(userDataSourceEmployeeOrgData)
	if result.GetEmployeeOrgData() != nil {
		if result.GetEmployeeOrgData().GetCostCenter() != nil {employeeOrgData.CostCenter = types.StringValue(*result.GetEmployeeOrgData().GetCostCenter())}
		if result.GetEmployeeOrgData().GetDivision()   != nil {employeeOrgData.Division   = types.StringValue(*result.GetEmployeeOrgData().GetDivision())}
	}
	state.EmployeeOrgData = employeeOrgData

	if result.GetEmployeeType()                    != nil {state.EmployeeType                    = types.StringValue(*result.GetEmployeeType())}
	if result.GetExternalUserState()               != nil {state.ExternalUserState               = types.StringValue(*result.GetExternalUserState())}
	if result.GetExternalUserStateChangeDateTime() != nil {state.ExternalUserStateChangeDateTime = types.StringValue(result.GetExternalUserStateChangeDateTime().String())}
	if result.GetFaxNumber()                       != nil {state.FaxNumber                       = types.StringValue(*result.GetFaxNumber())}
	if result.GetGivenName()                       != nil {state.GivenName                       = types.StringValue(*result.GetGivenName())}
	if result.GetHireDate()                        != nil {state.HireDate                        = types.StringValue(result.GetHireDate().String())}
	state.Id = types.StringValue(*result.GetId())

	for _, identity := range result.GetIdentities() {
		identitiesState := userDataSourceIdentities{
			Issuer:           types.StringValue(*identity.GetIssuer()),
			IssuerAssignedId: types.StringValue(*identity.GetIssuerAssignedId()),
			SignInType:       types.StringValue(*identity.GetSignInType()),
		}
		state.Identities = append(state.Identities, identitiesState)
	}

	for _, imAddress := range result.GetImAddresses() {
		state.ImAddresses = append(state.ImAddresses, types.StringValue(imAddress))
	}
	for _, interest := range result.GetInterests() {
		state.Interests = append(state.Interests, types.StringValue(interest))
	}

	if result.GetIsResourceAccount()           != nil {state.IsResourceAccount           = types.BoolValue(*result.GetIsResourceAccount())}
	if result.GetJobTitle()                    != nil {state.JobTitle                    = types.StringValue(*result.GetJobTitle())}
	if result.GetLastPasswordChangeDateTime()  != nil {state.LastPasswordChangeDateTime  = types.StringValue(result.GetLastPasswordChangeDateTime().String())}
	if result.GetLegalAgeGroupClassification() != nil {state.LegalAgeGroupClassification = types.StringValue(*result.GetLegalAgeGroupClassification())}

	for _, licenseAssignmentState := range result.GetLicenseAssignmentStates() {
		licenseAssignmentStateState := new(userDataSourceLicenseAssignmentStatesModel)
		if licenseAssignmentState.GetAssignedByGroup() != nil {
			licenseAssignmentStateState.AssignedByGroup = types.StringValue(*licenseAssignmentState.GetAssignedByGroup())
		}
		for _, disabledLicense := range licenseAssignmentState.GetDisabledPlans() {
			licenseAssignmentStateState.DisabledPlans = append(licenseAssignmentStateState.DisabledPlans, types.StringValue(disabledLicense.String()))
		}
		if licenseAssignmentState.GetError() != nil {
			licenseAssignmentStateState.Error = types.StringValue(*licenseAssignmentState.GetError())
		}
		if licenseAssignmentState.GetLastUpdatedDateTime() != nil {
			licenseAssignmentStateState.LastUpdatedDateTime = types.StringValue(licenseAssignmentState.GetLastUpdatedDateTime().String())
		}
		if licenseAssignmentState.GetSkuId() != nil {
			licenseAssignmentStateState.SkuId = types.StringValue(licenseAssignmentState.GetSkuId().String())
		}
		if licenseAssignmentState.GetState() != nil {
			licenseAssignmentStateState.State = types.StringValue(*licenseAssignmentState.GetState())
		}
		state.LicenseAssignmentStates = append(state.LicenseAssignmentStates, *licenseAssignmentStateState)
	}

	if result.GetMail() != nil {state.Mail = types.StringValue(*result.GetMail())}
	state.MailNickname      = types.StringValue(*result.GetMailNickname())

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
