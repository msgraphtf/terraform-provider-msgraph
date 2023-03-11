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
			//"mail": schema.StringAttribute{ // This property can only be used to access the authenticated users mailbox, not other users
			//	Computed: true,
			//},
			"mail_nickname": schema.StringAttribute{
				Computed: true,
				//TODO: Optional: true,
			},
			"mobile_phone": schema.StringAttribute{
				Computed: true,
			},
			"my_site": schema.StringAttribute{
				Computed: true,
			},
			"office_location": schema.StringAttribute{
				Computed: true,
			},
			"on_premises_distinguished_name": schema.StringAttribute{
				Computed: true,
			},
			"on_premises_domain_name": schema.StringAttribute{
				Computed: true,
			},
			"on_premises_extension_attributes": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"extension_attribute_1": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_2": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_3": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_4": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_5": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_6": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_7": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_8": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_9": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_10": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_11": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_12": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_13": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_14": schema.StringAttribute{
						Computed: true,
					},
					"extension_attribute_15": schema.StringAttribute{
						Computed: true,
					},
				},
			},
			"on_premises_immutable_id": schema.StringAttribute{
				Computed: true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Computed: true,
			},
			"on_premises_provisioning_errors": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category": schema.StringAttribute{
							Computed: true,
						},
						"occured_date_time": schema.StringAttribute{
							Computed: true,
						},
						"property_causing_error": schema.StringAttribute{
							Computed: true,
						},
						"value": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
			"on_premises_sam_account_name": schema.StringAttribute{
				Computed: true,
			},
			"on_premises_security_identifier": schema.StringAttribute{
				Computed: true,
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Computed: true,
			},
			"on_premises_user_principal_name": schema.StringAttribute{
				Computed: true,
			},
			"other_mails": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
			},
			"password_policies": schema.StringAttribute{
				Computed: true,
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
			"past_projects": schema.ListAttribute{
				Computed: true,
				ElementType: types.StringType,
			},
			"postal_code": schema.StringAttribute{
				Computed: true,
			},
			"preferred_data_location": schema.StringAttribute{
				Computed: true,
			},
			"preferred_language": schema.StringAttribute{
				Computed: true,
			},
			"preferred_name": schema.StringAttribute{
				Computed: true,
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
	//Mail                            types.String                         `tfsdk:"mail"`
	MailNickname                    types.String                         `tfsdk:"mail_nickname"`
	MobilePhone                     types.String                         `tfsdk:"mobile_phone"`
	OfficeLocation                  types.String                         `tfsdk:"office_location"`
	MySite                          types.String                         `tfsdk:"my_site"`
	OnPremisesDistinguishedName     types.String                         `tfsdk:"on_premises_distinguished_name"`
	OnPremisesDomainName            types.String                         `tfsdk:"on_premises_domain_name"`
	OnPremisesExtensionAttributes   *userDataSourceOnPremisesExtensionAttributesModel `tfsdk:"on_premises_extension_attributes"`
	OnPremisesImmutableId           types.String                         `tfsdk:"on_premises_immutable_id"`
	OnPremisesLastSyncDateTime      types.String                         `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesProvisioningErrors    []userDataSourceOnPremisesProvisioningErrorModel `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName        types.String                         `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier    types.String                         `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled           types.Bool                           `tfsdk:"on_premises_sync_enabled"`
	OnPremisesUserPrincipalName     types.String                         `tfsdk:"on_premises_user_principal_name"`
	OtherMails                      []types.String                       `tfsdk:"other_mails"`
	PasswordPolicies                types.String                         `tfsdk:"password_policies"`
	PasswordProfile                 *userDataSourcePasswordProfileModel  `tfsdk:"password_profile"`
	PastProjects                    []types.String                       `tfsdk:"past_projects"`
	PostalCode                      types.String                         `tfsdk:"postal_code"`
	PreferredDataLocation           types.String                         `tfsdk:"preferred_data_location"`
	PreferredLanguage               types.String                         `tfsdk:"preferred_language"`
	PreferredName                   types.String                         `tfsdk:"preferred_name"`
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

type userDataSourceOnPremisesExtensionAttributesModel struct {
	ExtensionAttribute1  types.String `tfsdk:"extension_attribute_1"`
	ExtensionAttribute2  types.String `tfsdk:"extension_attribute_2"`
	ExtensionAttribute3  types.String `tfsdk:"extension_attribute_3"`
	ExtensionAttribute4  types.String `tfsdk:"extension_attribute_4"`
	ExtensionAttribute5  types.String `tfsdk:"extension_attribute_5"`
	ExtensionAttribute6  types.String `tfsdk:"extension_attribute_6"`
	ExtensionAttribute7  types.String `tfsdk:"extension_attribute_7"`
	ExtensionAttribute8  types.String `tfsdk:"extension_attribute_8"`
	ExtensionAttribute9  types.String `tfsdk:"extension_attribute_9"`
	ExtensionAttribute10 types.String `tfsdk:"extension_attribute_10"`
	ExtensionAttribute11 types.String `tfsdk:"extension_attribute_11"`
	ExtensionAttribute12 types.String `tfsdk:"extension_attribute_12"`
	ExtensionAttribute13 types.String `tfsdk:"extension_attribute_13"`
	ExtensionAttribute14 types.String `tfsdk:"extension_attribute_14"`
	ExtensionAttribute15 types.String `tfsdk:"extension_attribute_15"`
}

type userDataSourceOnPremisesProvisioningErrorModel struct {
	Category             types.String `tfsdk:"category"`
	OccuredDateTime      types.String `tfsdk:"occured_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
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

	//if result.GetMail() != nil {state.Mail = types.StringValue(*result.GetMail())}
	state.MailNickname = types.StringValue(*result.GetMailNickname())
	if result.GetMobilePhone()                 != nil {state.MobilePhone = types.StringValue(*result.GetMobilePhone())}
	if result.GetMySite()                      != nil {state.MySite      = types.StringValue(*result.GetMySite())}
	if result.GetOfficeLocation()              != nil {state.OfficeLocation = types.StringValue(*result.GetOfficeLocation())}
	if result.GetOnPremisesDistinguishedName() != nil {state.OnPremisesDistinguishedName = types.StringValue(*result.GetOnPremisesDistinguishedName())}
	if result.GetOnPremisesDomainName()        != nil {state.OnPremisesDomainName = types.StringValue(*result.GetOnPremisesDomainName())}

	onPremisesExtensionAttributes := new(userDataSourceOnPremisesExtensionAttributesModel)
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1()  != nil {onPremisesExtensionAttributes.ExtensionAttribute1  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2()  != nil {onPremisesExtensionAttributes.ExtensionAttribute2  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3()  != nil {onPremisesExtensionAttributes.ExtensionAttribute3  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4()  != nil {onPremisesExtensionAttributes.ExtensionAttribute4  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5()  != nil {onPremisesExtensionAttributes.ExtensionAttribute5  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6()  != nil {onPremisesExtensionAttributes.ExtensionAttribute6  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7()  != nil {onPremisesExtensionAttributes.ExtensionAttribute7  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8()  != nil {onPremisesExtensionAttributes.ExtensionAttribute8  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9()  != nil {onPremisesExtensionAttributes.ExtensionAttribute9  = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10() != nil {onPremisesExtensionAttributes.ExtensionAttribute10 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11() != nil {onPremisesExtensionAttributes.ExtensionAttribute11 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12() != nil {onPremisesExtensionAttributes.ExtensionAttribute12 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13() != nil {onPremisesExtensionAttributes.ExtensionAttribute13 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14() != nil {onPremisesExtensionAttributes.ExtensionAttribute14 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14())}
	if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15() != nil {onPremisesExtensionAttributes.ExtensionAttribute15 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15())}
	state.OnPremisesExtensionAttributes = onPremisesExtensionAttributes

	if result.GetOnPremisesImmutableId() != nil {state.OnPremisesImmutableId = types.StringValue(*result.GetOnPremisesImmutableId())}
	if result.GetOnPremisesLastSyncDateTime() != nil {state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())}

	for _, onPremisesProvisioningError := range result.GetOnPremisesProvisioningErrors() {
			onPremisesProvisioningErrorState := userDataSourceOnPremisesProvisioningErrorModel{
				Category:              types.StringValue(*onPremisesProvisioningError.GetCategory()),
				OccuredDateTime:        types.StringValue(onPremisesProvisioningError.GetOccurredDateTime().String()),
				PropertyCausingError:   types.StringValue(*onPremisesProvisioningError.GetPropertyCausingError()),
				Value:                  types.StringValue(*onPremisesProvisioningError.GetValue()),
			}
		state.OnPremisesProvisioningErrors = append(state.OnPremisesProvisioningErrors, onPremisesProvisioningErrorState)
	}

	if result.GetOnPremisesSamAccountName()     != nil {state.OnPremisesSamAccountName     = types.StringValue(*result.GetOnPremisesSamAccountName())}
	if result.GetOnPremisesSecurityIdentifier() != nil {state.OnPremisesSecurityIdentifier = types.StringValue(*result.GetOnPremisesSecurityIdentifier())}
	if result.GetOnPremisesSyncEnabled()        != nil {state.OnPremisesSyncEnabled        = types.BoolValue(*result.GetOnPremisesSyncEnabled())}
	if result.GetOnPremisesUserPrincipalName()  != nil {state.OnPremisesUserPrincipalName  = types.StringValue(*result.GetOnPremisesUserPrincipalName())}

	for _, other_mail := range result.GetOtherMails() {
		state.OtherMails = append(state.OtherMails, types.StringValue(other_mail))
	}

	if result.GetPasswordPolicies() != nil {state.PasswordPolicies = types.StringValue(*result.GetPasswordPolicies())}

	passwordProfile := new(userDataSourcePasswordProfileModel)
	passwordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
	passwordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
	state.PasswordProfile = passwordProfile

	for _, past_project := range result.GetPastProjects() {
		state.PastProjects = append(state.PastProjects, types.StringValue(past_project))
	}

	if result.GetPostalCode()            != nil {state.PostalCode            = types.StringValue(*result.GetPostalCode())}
	if result.GetPreferredDataLocation() != nil {state.PreferredDataLocation = types.StringValue(*result.GetPreferredDataLocation())}
	if result.GetPreferredLanguage()     != nil {state.PreferredLanguage     = types.StringValue(*result.GetPreferredLanguage())}
	if result.GetPreferredName()         != nil {state.PreferredName         = types.StringValue(*result.GetPreferredName())}

	state.UserPrincipalName = types.StringValue(*result.GetUserPrincipalName())


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
