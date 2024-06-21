package users

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type userModel struct {
	Id                              types.String `tfsdk:"id"`
	DeletedDateTime                 types.String `tfsdk:"deleted_date_time"`
	AboutMe                         types.String `tfsdk:"about_me"`
	AccountEnabled                  types.Bool   `tfsdk:"account_enabled"`
	AgeGroup                        types.String `tfsdk:"age_group"`
	AssignedLicenses                types.List   `tfsdk:"assigned_licenses"`
	AssignedPlans                   types.List   `tfsdk:"assigned_plans"`
	AuthorizationInfo               types.Object `tfsdk:"authorization_info"`
	Birthday                        types.String `tfsdk:"birthday"`
	BusinessPhones                  types.List   `tfsdk:"business_phones"`
	City                            types.String `tfsdk:"city"`
	CompanyName                     types.String `tfsdk:"company_name"`
	ConsentProvidedForMinor         types.String `tfsdk:"consent_provided_for_minor"`
	Country                         types.String `tfsdk:"country"`
	CreatedDateTime                 types.String `tfsdk:"created_date_time"`
	CreationType                    types.String `tfsdk:"creation_type"`
	Department                      types.String `tfsdk:"department"`
	DisplayName                     types.String `tfsdk:"display_name"`
	EmployeeHireDate                types.String `tfsdk:"employee_hire_date"`
	EmployeeId                      types.String `tfsdk:"employee_id"`
	EmployeeLeaveDateTime           types.String `tfsdk:"employee_leave_date_time"`
	EmployeeOrgData                 types.Object `tfsdk:"employee_org_data"`
	EmployeeType                    types.String `tfsdk:"employee_type"`
	ExternalUserState               types.String `tfsdk:"external_user_state"`
	ExternalUserStateChangeDateTime types.String `tfsdk:"external_user_state_change_date_time"`
	FaxNumber                       types.String `tfsdk:"fax_number"`
	GivenName                       types.String `tfsdk:"given_name"`
	HireDate                        types.String `tfsdk:"hire_date"`
	Identities                      types.List   `tfsdk:"identities"`
	ImAddresses                     types.List   `tfsdk:"im_addresses"`
	Interests                       types.List   `tfsdk:"interests"`
	IsResourceAccount               types.Bool   `tfsdk:"is_resource_account"`
	JobTitle                        types.String `tfsdk:"job_title"`
	LastPasswordChangeDateTime      types.String `tfsdk:"last_password_change_date_time"`
	LegalAgeGroupClassification     types.String `tfsdk:"legal_age_group_classification"`
	LicenseAssignmentStates         types.List   `tfsdk:"license_assignment_states"`
	Mail                            types.String `tfsdk:"mail"`
	MailNickname                    types.String `tfsdk:"mail_nickname"`
	MobilePhone                     types.String `tfsdk:"mobile_phone"`
	MySite                          types.String `tfsdk:"my_site"`
	OfficeLocation                  types.String `tfsdk:"office_location"`
	OnPremisesDistinguishedName     types.String `tfsdk:"on_premises_distinguished_name"`
	OnPremisesDomainName            types.String `tfsdk:"on_premises_domain_name"`
	OnPremisesExtensionAttributes   types.Object `tfsdk:"on_premises_extension_attributes"`
	OnPremisesImmutableId           types.String `tfsdk:"on_premises_immutable_id"`
	OnPremisesLastSyncDateTime      types.String `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesProvisioningErrors    types.List   `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName        types.String `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier    types.String `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled           types.Bool   `tfsdk:"on_premises_sync_enabled"`
	OnPremisesUserPrincipalName     types.String `tfsdk:"on_premises_user_principal_name"`
	OtherMails                      types.List   `tfsdk:"other_mails"`
	PasswordPolicies                types.String `tfsdk:"password_policies"`
	PasswordProfile                 types.Object `tfsdk:"password_profile"`
	PastProjects                    types.List   `tfsdk:"past_projects"`
	PostalCode                      types.String `tfsdk:"postal_code"`
	PreferredDataLocation           types.String `tfsdk:"preferred_data_location"`
	PreferredLanguage               types.String `tfsdk:"preferred_language"`
	PreferredName                   types.String `tfsdk:"preferred_name"`
	ProvisionedPlans                types.List   `tfsdk:"provisioned_plans"`
	ProxyAddresses                  types.List   `tfsdk:"proxy_addresses"`
	Responsibilities                types.List   `tfsdk:"responsibilities"`
	Schools                         types.List   `tfsdk:"schools"`
	SecurityIdentifier              types.String `tfsdk:"security_identifier"`
	ServiceProvisioningErrors       types.List   `tfsdk:"service_provisioning_errors"`
	ShowInAddressList               types.Bool   `tfsdk:"show_in_address_list"`
	SignInActivity                  types.Object `tfsdk:"sign_in_activity"`
	SignInSessionsValidFromDateTime types.String `tfsdk:"sign_in_sessions_valid_from_date_time"`
	Skills                          types.List   `tfsdk:"skills"`
	State                           types.String `tfsdk:"state"`
	StreetAddress                   types.String `tfsdk:"street_address"`
	Surname                         types.String `tfsdk:"surname"`
	UsageLocation                   types.String `tfsdk:"usage_location"`
	UserPrincipalName               types.String `tfsdk:"user_principal_name"`
	UserType                        types.String `tfsdk:"user_type"`
}

func (m userModel) AttributeTypes() map[string]attr.Type {
	userAssignedLicenses := userAssignedLicensesModel{}
	userAssignedPlans := userAssignedPlansModel{}
	userAuthorizationInfo := userAuthorizationInfoModel{}
	userEmployeeOrgData := userEmployeeOrgDataModel{}
	userIdentities := userIdentitiesModel{}
	userLicenseAssignmentStates := userLicenseAssignmentStatesModel{}
	userOnPremisesExtensionAttributes := userOnPremisesExtensionAttributesModel{}
	userOnPremisesProvisioningErrors := userOnPremisesProvisioningErrorsModel{}
	userPasswordProfile := userPasswordProfileModel{}
	userProvisionedPlans := userProvisionedPlansModel{}
	userServiceProvisioningErrors := userServiceProvisioningErrorsModel{}
	userSignInActivity := userSignInActivityModel{}
	return map[string]attr.Type{
		"id":                                    types.StringType,
		"deleted_date_time":                     types.StringType,
		"about_me":                              types.StringType,
		"account_enabled":                       types.BoolType,
		"age_group":                             types.StringType,
		"assigned_licenses":                     types.ListType{ElemType: types.ObjectType{AttrTypes: userAssignedLicenses.AttributeTypes()}},
		"assigned_plans":                        types.ListType{ElemType: types.ObjectType{AttrTypes: userAssignedPlans.AttributeTypes()}},
		"authorization_info":                    types.ObjectType{AttrTypes: userAuthorizationInfo.AttributeTypes()},
		"birthday":                              types.StringType,
		"business_phones":                       types.ListType{ElemType: types.StringType},
		"city":                                  types.StringType,
		"company_name":                          types.StringType,
		"consent_provided_for_minor":            types.StringType,
		"country":                               types.StringType,
		"created_date_time":                     types.StringType,
		"creation_type":                         types.StringType,
		"department":                            types.StringType,
		"display_name":                          types.StringType,
		"employee_hire_date":                    types.StringType,
		"employee_id":                           types.StringType,
		"employee_leave_date_time":              types.StringType,
		"employee_org_data":                     types.ObjectType{AttrTypes: userEmployeeOrgData.AttributeTypes()},
		"employee_type":                         types.StringType,
		"external_user_state":                   types.StringType,
		"external_user_state_change_date_time":  types.StringType,
		"fax_number":                            types.StringType,
		"given_name":                            types.StringType,
		"hire_date":                             types.StringType,
		"identities":                            types.ListType{ElemType: types.ObjectType{AttrTypes: userIdentities.AttributeTypes()}},
		"im_addresses":                          types.ListType{ElemType: types.StringType},
		"interests":                             types.ListType{ElemType: types.StringType},
		"is_resource_account":                   types.BoolType,
		"job_title":                             types.StringType,
		"last_password_change_date_time":        types.StringType,
		"legal_age_group_classification":        types.StringType,
		"license_assignment_states":             types.ListType{ElemType: types.ObjectType{AttrTypes: userLicenseAssignmentStates.AttributeTypes()}},
		"mail":                                  types.StringType,
		"mail_nickname":                         types.StringType,
		"mobile_phone":                          types.StringType,
		"my_site":                               types.StringType,
		"office_location":                       types.StringType,
		"on_premises_distinguished_name":        types.StringType,
		"on_premises_domain_name":               types.StringType,
		"on_premises_extension_attributes":      types.ObjectType{AttrTypes: userOnPremisesExtensionAttributes.AttributeTypes()},
		"on_premises_immutable_id":              types.StringType,
		"on_premises_last_sync_date_time":       types.StringType,
		"on_premises_provisioning_errors":       types.ListType{ElemType: types.ObjectType{AttrTypes: userOnPremisesProvisioningErrors.AttributeTypes()}},
		"on_premises_sam_account_name":          types.StringType,
		"on_premises_security_identifier":       types.StringType,
		"on_premises_sync_enabled":              types.BoolType,
		"on_premises_user_principal_name":       types.StringType,
		"other_mails":                           types.ListType{ElemType: types.StringType},
		"password_policies":                     types.StringType,
		"password_profile":                      types.ObjectType{AttrTypes: userPasswordProfile.AttributeTypes()},
		"past_projects":                         types.ListType{ElemType: types.StringType},
		"postal_code":                           types.StringType,
		"preferred_data_location":               types.StringType,
		"preferred_language":                    types.StringType,
		"preferred_name":                        types.StringType,
		"provisioned_plans":                     types.ListType{ElemType: types.ObjectType{AttrTypes: userProvisionedPlans.AttributeTypes()}},
		"proxy_addresses":                       types.ListType{ElemType: types.StringType},
		"responsibilities":                      types.ListType{ElemType: types.StringType},
		"schools":                               types.ListType{ElemType: types.StringType},
		"security_identifier":                   types.StringType,
		"service_provisioning_errors":           types.ListType{ElemType: types.ObjectType{AttrTypes: userServiceProvisioningErrors.AttributeTypes()}},
		"show_in_address_list":                  types.BoolType,
		"sign_in_activity":                      types.ObjectType{AttrTypes: userSignInActivity.AttributeTypes()},
		"sign_in_sessions_valid_from_date_time": types.StringType,
		"skills":                                types.ListType{ElemType: types.StringType},
		"state":                                 types.StringType,
		"street_address":                        types.StringType,
		"surname":                               types.StringType,
		"usage_location":                        types.StringType,
		"user_principal_name":                   types.StringType,
		"user_type":                             types.StringType,
	}
}

type userAssignedLicensesModel struct {
	DisabledPlans types.List   `tfsdk:"disabled_plans"`
	SkuId         types.String `tfsdk:"sku_id"`
}

func (m userAssignedLicensesModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"disabled_plans": types.ListType{ElemType: types.StringType},
		"sku_id":         types.StringType,
	}
}

type userAssignedPlansModel struct {
	AssignedDateTime types.String `tfsdk:"assigned_date_time"`
	CapabilityStatus types.String `tfsdk:"capability_status"`
	Service          types.String `tfsdk:"service"`
	ServicePlanId    types.String `tfsdk:"service_plan_id"`
}

func (m userAssignedPlansModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"assigned_date_time": types.StringType,
		"capability_status":  types.StringType,
		"service":            types.StringType,
		"service_plan_id":    types.StringType,
	}
}

type userAuthorizationInfoModel struct {
	CertificateUserIds types.List `tfsdk:"certificate_user_ids"`
}

func (m userAuthorizationInfoModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"certificate_user_ids": types.ListType{ElemType: types.StringType},
	}
}

type userEmployeeOrgDataModel struct {
	CostCenter types.String `tfsdk:"cost_center"`
	Division   types.String `tfsdk:"division"`
}

func (m userEmployeeOrgDataModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cost_center": types.StringType,
		"division":    types.StringType,
	}
}

type userIdentitiesModel struct {
	Issuer           types.String `tfsdk:"issuer"`
	IssuerAssignedId types.String `tfsdk:"issuer_assigned_id"`
	SignInType       types.String `tfsdk:"sign_in_type"`
}

func (m userIdentitiesModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"issuer":             types.StringType,
		"issuer_assigned_id": types.StringType,
		"sign_in_type":       types.StringType,
	}
}

type userLicenseAssignmentStatesModel struct {
	AssignedByGroup     types.String `tfsdk:"assigned_by_group"`
	DisabledPlans       types.List   `tfsdk:"disabled_plans"`
	Error               types.String `tfsdk:"error"`
	LastUpdatedDateTime types.String `tfsdk:"last_updated_date_time"`
	SkuId               types.String `tfsdk:"sku_id"`
	State               types.String `tfsdk:"state"`
}

func (m userLicenseAssignmentStatesModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"assigned_by_group":      types.StringType,
		"disabled_plans":         types.ListType{ElemType: types.StringType},
		"error":                  types.StringType,
		"last_updated_date_time": types.StringType,
		"sku_id":                 types.StringType,
		"state":                  types.StringType,
	}
}

type userOnPremisesExtensionAttributesModel struct {
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

func (m userOnPremisesExtensionAttributesModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"extension_attribute_1":  types.StringType,
		"extension_attribute_10": types.StringType,
		"extension_attribute_11": types.StringType,
		"extension_attribute_12": types.StringType,
		"extension_attribute_13": types.StringType,
		"extension_attribute_14": types.StringType,
		"extension_attribute_15": types.StringType,
		"extension_attribute_2":  types.StringType,
		"extension_attribute_3":  types.StringType,
		"extension_attribute_4":  types.StringType,
		"extension_attribute_5":  types.StringType,
		"extension_attribute_6":  types.StringType,
		"extension_attribute_7":  types.StringType,
		"extension_attribute_8":  types.StringType,
		"extension_attribute_9":  types.StringType,
	}
}

type userOnPremisesProvisioningErrorsModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

func (m userOnPremisesProvisioningErrorsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"category":               types.StringType,
		"occurred_date_time":     types.StringType,
		"property_causing_error": types.StringType,
		"value":                  types.StringType,
	}
}

type userPasswordProfileModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
	Password                             types.String `tfsdk:"password"`
}

func (m userPasswordProfileModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"force_change_password_next_sign_in":          types.BoolType,
		"force_change_password_next_sign_in_with_mfa": types.BoolType,
		"password": types.StringType,
	}
}

type userProvisionedPlansModel struct {
	CapabilityStatus   types.String `tfsdk:"capability_status"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	Service            types.String `tfsdk:"service"`
}

func (m userProvisionedPlansModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"capability_status":   types.StringType,
		"provisioning_status": types.StringType,
		"service":             types.StringType,
	}
}

type userServiceProvisioningErrorsModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

func (m userServiceProvisioningErrorsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_date_time": types.StringType,
		"is_resolved":       types.BoolType,
		"service_instance":  types.StringType,
	}
}

type userSignInActivityModel struct {
	LastNonInteractiveSignInDateTime  types.String `tfsdk:"last_non_interactive_sign_in_date_time"`
	LastNonInteractiveSignInRequestId types.String `tfsdk:"last_non_interactive_sign_in_request_id"`
	LastSignInDateTime                types.String `tfsdk:"last_sign_in_date_time"`
	LastSignInRequestId               types.String `tfsdk:"last_sign_in_request_id"`
	LastSuccessfulSignInDateTime      types.String `tfsdk:"last_successful_sign_in_date_time"`
	LastSuccessfulSignInRequestId     types.String `tfsdk:"last_successful_sign_in_request_id"`
}

func (m userSignInActivityModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"last_non_interactive_sign_in_date_time":  types.StringType,
		"last_non_interactive_sign_in_request_id": types.StringType,
		"last_sign_in_date_time":                  types.StringType,
		"last_sign_in_request_id":                 types.StringType,
		"last_successful_sign_in_date_time":       types.StringType,
		"last_successful_sign_in_request_id":      types.StringType,
	}
}
