package users

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type usersModel struct {
	Value []usersValueModel `tfsdk:"value"`
}

type usersValueModel struct {
	Id                              types.String                             `tfsdk:"id"`
	DeletedDateTime                 types.String                             `tfsdk:"deleted_date_time"`
	AccountEnabled                  types.Bool                               `tfsdk:"account_enabled"`
	AgeGroup                        types.String                             `tfsdk:"age_group"`
	AssignedLicenses                []usersAssignedLicensesModel             `tfsdk:"assigned_licenses"`
	AssignedPlans                   []usersAssignedPlansModel                `tfsdk:"assigned_plans"`
	AuthorizationInfo               *usersAuthorizationInfoModel             `tfsdk:"authorization_info"`
	BusinessPhones                  []types.String                           `tfsdk:"business_phones"`
	City                            types.String                             `tfsdk:"city"`
	CompanyName                     types.String                             `tfsdk:"company_name"`
	ConsentProvidedForMinor         types.String                             `tfsdk:"consent_provided_for_minor"`
	Country                         types.String                             `tfsdk:"country"`
	CreatedDateTime                 types.String                             `tfsdk:"created_date_time"`
	CreationType                    types.String                             `tfsdk:"creation_type"`
	CustomSecurityAttributes        *usersCustomSecurityAttributesModel      `tfsdk:"custom_security_attributes"`
	Department                      types.String                             `tfsdk:"department"`
	DisplayName                     types.String                             `tfsdk:"display_name"`
	EmployeeHireDate                types.String                             `tfsdk:"employee_hire_date"`
	EmployeeId                      types.String                             `tfsdk:"employee_id"`
	EmployeeLeaveDateTime           types.String                             `tfsdk:"employee_leave_date_time"`
	EmployeeOrgData                 *usersEmployeeOrgDataModel               `tfsdk:"employee_org_data"`
	EmployeeType                    types.String                             `tfsdk:"employee_type"`
	ExternalUserState               types.String                             `tfsdk:"external_user_state"`
	ExternalUserStateChangeDateTime types.String                             `tfsdk:"external_user_state_change_date_time"`
	FaxNumber                       types.String                             `tfsdk:"fax_number"`
	GivenName                       types.String                             `tfsdk:"given_name"`
	Identities                      []usersIdentitiesModel                   `tfsdk:"identities"`
	ImAddresses                     []types.String                           `tfsdk:"im_addresses"`
	IsResourceAccount               types.Bool                               `tfsdk:"is_resource_account"`
	JobTitle                        types.String                             `tfsdk:"job_title"`
	LastPasswordChangeDateTime      types.String                             `tfsdk:"last_password_change_date_time"`
	LegalAgeGroupClassification     types.String                             `tfsdk:"legal_age_group_classification"`
	LicenseAssignmentStates         []usersLicenseAssignmentStatesModel      `tfsdk:"license_assignment_states"`
	Mail                            types.String                             `tfsdk:"mail"`
	MailNickname                    types.String                             `tfsdk:"mail_nickname"`
	MobilePhone                     types.String                             `tfsdk:"mobile_phone"`
	OfficeLocation                  types.String                             `tfsdk:"office_location"`
	OnPremisesDistinguishedName     types.String                             `tfsdk:"on_premises_distinguished_name"`
	OnPremisesDomainName            types.String                             `tfsdk:"on_premises_domain_name"`
	OnPremisesExtensionAttributes   *usersOnPremisesExtensionAttributesModel `tfsdk:"on_premises_extension_attributes"`
	OnPremisesImmutableId           types.String                             `tfsdk:"on_premises_immutable_id"`
	OnPremisesLastSyncDateTime      types.String                             `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesProvisioningErrors    []usersOnPremisesProvisioningErrorsModel `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName        types.String                             `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier    types.String                             `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled           types.Bool                               `tfsdk:"on_premises_sync_enabled"`
	OnPremisesUserPrincipalName     types.String                             `tfsdk:"on_premises_user_principal_name"`
	OtherMails                      []types.String                           `tfsdk:"other_mails"`
	PasswordPolicies                types.String                             `tfsdk:"password_policies"`
	PasswordProfile                 *usersPasswordProfileModel               `tfsdk:"password_profile"`
	PostalCode                      types.String                             `tfsdk:"postal_code"`
	PreferredDataLocation           types.String                             `tfsdk:"preferred_data_location"`
	PreferredLanguage               types.String                             `tfsdk:"preferred_language"`
	ProvisionedPlans                []usersProvisionedPlansModel             `tfsdk:"provisioned_plans"`
	ProxyAddresses                  []types.String                           `tfsdk:"proxy_addresses"`
	SecurityIdentifier              types.String                             `tfsdk:"security_identifier"`
	ServiceProvisioningErrors       []usersServiceProvisioningErrorsModel    `tfsdk:"service_provisioning_errors"`
	ShowInAddressList               types.Bool                               `tfsdk:"show_in_address_list"`
	SignInActivity                  *usersSignInActivityModel                `tfsdk:"sign_in_activity"`
	SignInSessionsValidFromDateTime types.String                             `tfsdk:"sign_in_sessions_valid_from_date_time"`
	State                           types.String                             `tfsdk:"state"`
	StreetAddress                   types.String                             `tfsdk:"street_address"`
	Surname                         types.String                             `tfsdk:"surname"`
	UsageLocation                   types.String                             `tfsdk:"usage_location"`
	UserPrincipalName               types.String                             `tfsdk:"user_principal_name"`
	UserType                        types.String                             `tfsdk:"user_type"`
}

type usersAssignedLicensesModel struct {
	DisabledPlans []types.String `tfsdk:"disabled_plans"`
	SkuId         types.String   `tfsdk:"sku_id"`
}

type usersAssignedPlansModel struct {
	AssignedDateTime types.String `tfsdk:"assigned_date_time"`
	CapabilityStatus types.String `tfsdk:"capability_status"`
	Service          types.String `tfsdk:"service"`
	ServicePlanId    types.String `tfsdk:"service_plan_id"`
}

type usersAuthorizationInfoModel struct {
	CertificateUserIds []types.String `tfsdk:"certificate_user_ids"`
}

type usersCustomSecurityAttributesModel struct {
}

type usersEmployeeOrgDataModel struct {
	CostCenter types.String `tfsdk:"cost_center"`
	Division   types.String `tfsdk:"division"`
}

type usersIdentitiesModel struct {
	Issuer           types.String `tfsdk:"issuer"`
	IssuerAssignedId types.String `tfsdk:"issuer_assigned_id"`
	SignInType       types.String `tfsdk:"sign_in_type"`
}

type usersLicenseAssignmentStatesModel struct {
	AssignedByGroup     types.String   `tfsdk:"assigned_by_group"`
	DisabledPlans       []types.String `tfsdk:"disabled_plans"`
	Error               types.String   `tfsdk:"error"`
	LastUpdatedDateTime types.String   `tfsdk:"last_updated_date_time"`
	SkuId               types.String   `tfsdk:"sku_id"`
	State               types.String   `tfsdk:"state"`
}

type usersOnPremisesExtensionAttributesModel struct {
	ExtensionAttribute1  types.String `tfsdk:"extension_attribute_1"`
	ExtensionAttribute10 types.String `tfsdk:"extension_attribute_10"`
	ExtensionAttribute11 types.String `tfsdk:"extension_attribute_11"`
	ExtensionAttribute12 types.String `tfsdk:"extension_attribute_12"`
	ExtensionAttribute13 types.String `tfsdk:"extension_attribute_13"`
	ExtensionAttribute14 types.String `tfsdk:"extension_attribute_14"`
	ExtensionAttribute15 types.String `tfsdk:"extension_attribute_15"`
	ExtensionAttribute2  types.String `tfsdk:"extension_attribute_2"`
	ExtensionAttribute3  types.String `tfsdk:"extension_attribute_3"`
	ExtensionAttribute4  types.String `tfsdk:"extension_attribute_4"`
	ExtensionAttribute5  types.String `tfsdk:"extension_attribute_5"`
	ExtensionAttribute6  types.String `tfsdk:"extension_attribute_6"`
	ExtensionAttribute7  types.String `tfsdk:"extension_attribute_7"`
	ExtensionAttribute8  types.String `tfsdk:"extension_attribute_8"`
	ExtensionAttribute9  types.String `tfsdk:"extension_attribute_9"`
}

type usersOnPremisesProvisioningErrorsModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

type usersPasswordProfileModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
	Password                             types.String `tfsdk:"password"`
}

type usersProvisionedPlansModel struct {
	CapabilityStatus   types.String `tfsdk:"capability_status"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	Service            types.String `tfsdk:"service"`
}

type usersServiceProvisioningErrorsModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

type usersSignInActivityModel struct {
	LastNonInteractiveSignInDateTime  types.String `tfsdk:"last_non_interactive_sign_in_date_time"`
	LastNonInteractiveSignInRequestId types.String `tfsdk:"last_non_interactive_sign_in_request_id"`
	LastSignInDateTime                types.String `tfsdk:"last_sign_in_date_time"`
	LastSignInRequestId               types.String `tfsdk:"last_sign_in_request_id"`
}
