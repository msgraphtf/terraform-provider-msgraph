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
	CustomSecurityAttributes        types.Object `tfsdk:"custom_security_attributes"`
	Department                      types.String `tfsdk:"department"`
	DeviceEnrollmentLimit           types.Int64  `tfsdk:"device_enrollment_limit"`
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
	IsManagementRestricted          types.Bool   `tfsdk:"is_management_restricted"`
	IsResourceAccount               types.Bool   `tfsdk:"is_resource_account"`
	JobTitle                        types.String `tfsdk:"job_title"`
	LastPasswordChangeDateTime      types.String `tfsdk:"last_password_change_date_time"`
	LegalAgeGroupClassification     types.String `tfsdk:"legal_age_group_classification"`
	LicenseAssignmentStates         types.List   `tfsdk:"license_assignment_states"`
	Mail                            types.String `tfsdk:"mail"`
	MailNickname                    types.String `tfsdk:"mail_nickname"`
	MailboxSettings                 types.Object `tfsdk:"mailbox_settings"`
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
	Print                           types.Object `tfsdk:"print"`
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
	userAssignedLicenses := userAssignedLicenseModel{}
	userAssignedPlans := userAssignedPlanModel{}
	userAuthorizationInfo := userAuthorizationInfoModel{}
	userCustomSecurityAttributes := userCustomSecurityAttributeValueModel{}
	userEmployeeOrgData := userEmployeeOrgDataModel{}
	userIdentities := userObjectIdentityModel{}
	userLicenseAssignmentStates := userLicenseAssignmentStateModel{}
	userMailboxSettings := userMailboxSettingsModel{}
	userOnPremisesExtensionAttributes := userOnPremisesExtensionAttributesModel{}
	userOnPremisesProvisioningErrors := userOnPremisesProvisioningErrorModel{}
	userPasswordProfile := userPasswordProfileModel{}
	userPrint := userUserPrintModel{}
	userProvisionedPlans := userProvisionedPlanModel{}
	userServiceProvisioningErrors := userServiceProvisioningErrorModel{}
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
		"custom_security_attributes":            types.ObjectType{AttrTypes: userCustomSecurityAttributes.AttributeTypes()},
		"department":                            types.StringType,
		"device_enrollment_limit":               types.Int64Type,
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
		"is_management_restricted":              types.BoolType,
		"is_resource_account":                   types.BoolType,
		"job_title":                             types.StringType,
		"last_password_change_date_time":        types.StringType,
		"legal_age_group_classification":        types.StringType,
		"license_assignment_states":             types.ListType{ElemType: types.ObjectType{AttrTypes: userLicenseAssignmentStates.AttributeTypes()}},
		"mail":                                  types.StringType,
		"mail_nickname":                         types.StringType,
		"mailbox_settings":                      types.ObjectType{AttrTypes: userMailboxSettings.AttributeTypes()},
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
		"print":                                 types.ObjectType{AttrTypes: userPrint.AttributeTypes()},
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

type userAssignedLicenseModel struct {
	DisabledPlans types.List   `tfsdk:"disabled_plans"`
	SkuId         types.String `tfsdk:"sku_id"`
}

func (m userAssignedLicenseModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"disabled_plans": types.ListType{ElemType: types.StringType},
		"sku_id":         types.StringType,
	}
}

type userAssignedPlanModel struct {
	AssignedDateTime types.String `tfsdk:"assigned_date_time"`
	CapabilityStatus types.String `tfsdk:"capability_status"`
	Service          types.String `tfsdk:"service"`
	ServicePlanId    types.String `tfsdk:"service_plan_id"`
}

func (m userAssignedPlanModel) AttributeTypes() map[string]attr.Type {
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

type userCustomSecurityAttributeValueModel struct {
}

func (m userCustomSecurityAttributeValueModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{}
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

type userObjectIdentityModel struct {
	Issuer           types.String `tfsdk:"issuer"`
	IssuerAssignedId types.String `tfsdk:"issuer_assigned_id"`
	SignInType       types.String `tfsdk:"sign_in_type"`
}

func (m userObjectIdentityModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"issuer":             types.StringType,
		"issuer_assigned_id": types.StringType,
		"sign_in_type":       types.StringType,
	}
}

type userLicenseAssignmentStateModel struct {
	AssignedByGroup     types.String `tfsdk:"assigned_by_group"`
	DisabledPlans       types.List   `tfsdk:"disabled_plans"`
	Error               types.String `tfsdk:"error"`
	LastUpdatedDateTime types.String `tfsdk:"last_updated_date_time"`
	SkuId               types.String `tfsdk:"sku_id"`
	State               types.String `tfsdk:"state"`
}

func (m userLicenseAssignmentStateModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"assigned_by_group":      types.StringType,
		"disabled_plans":         types.ListType{ElemType: types.StringType},
		"error":                  types.StringType,
		"last_updated_date_time": types.StringType,
		"sku_id":                 types.StringType,
		"state":                  types.StringType,
	}
}

type userMailboxSettingsModel struct {
	ArchiveFolder                         types.String `tfsdk:"archive_folder"`
	AutomaticRepliesSetting               types.Object `tfsdk:"automatic_replies_setting"`
	DateFormat                            types.String `tfsdk:"date_format"`
	DelegateMeetingMessageDeliveryOptions types.String `tfsdk:"delegate_meeting_message_delivery_options"`
	Language                              types.Object `tfsdk:"language"`
	TimeFormat                            types.String `tfsdk:"time_format"`
	TimeZone                              types.String `tfsdk:"time_zone"`
	UserPurpose                           types.String `tfsdk:"user_purpose"`
	WorkingHours                          types.Object `tfsdk:"working_hours"`
}

func (m userMailboxSettingsModel) AttributeTypes() map[string]attr.Type {
	userAutomaticRepliesSetting := userAutomaticRepliesSettingModel{}
	userLanguage := userLocaleInfoModel{}
	userWorkingHours := userWorkingHoursModel{}
	return map[string]attr.Type{
		"archive_folder":                            types.StringType,
		"automatic_replies_setting":                 types.ObjectType{AttrTypes: userAutomaticRepliesSetting.AttributeTypes()},
		"date_format":                               types.StringType,
		"delegate_meeting_message_delivery_options": types.StringType,
		"language":                                  types.ObjectType{AttrTypes: userLanguage.AttributeTypes()},
		"time_format":                               types.StringType,
		"time_zone":                                 types.StringType,
		"user_purpose":                              types.StringType,
		"working_hours":                             types.ObjectType{AttrTypes: userWorkingHours.AttributeTypes()},
	}
}

type userAutomaticRepliesSettingModel struct {
	ExternalAudience       types.String `tfsdk:"external_audience"`
	ExternalReplyMessage   types.String `tfsdk:"external_reply_message"`
	InternalReplyMessage   types.String `tfsdk:"internal_reply_message"`
	ScheduledEndDateTime   types.Object `tfsdk:"scheduled_end_date_time"`
	ScheduledStartDateTime types.Object `tfsdk:"scheduled_start_date_time"`
	Status                 types.String `tfsdk:"status"`
}

func (m userAutomaticRepliesSettingModel) AttributeTypes() map[string]attr.Type {
	userScheduledEndDateTime := userDateTimeTimeZoneModel{}
	userScheduledStartDateTime := userDateTimeTimeZoneModel{}
	return map[string]attr.Type{
		"external_audience":         types.StringType,
		"external_reply_message":    types.StringType,
		"internal_reply_message":    types.StringType,
		"scheduled_end_date_time":   types.ObjectType{AttrTypes: userScheduledEndDateTime.AttributeTypes()},
		"scheduled_start_date_time": types.ObjectType{AttrTypes: userScheduledStartDateTime.AttributeTypes()},
		"status":                    types.StringType,
	}
}

type userDateTimeTimeZoneModel struct {
	DateTime types.String `tfsdk:"date_time"`
	TimeZone types.String `tfsdk:"time_zone"`
}

func (m userDateTimeTimeZoneModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"date_time": types.StringType,
		"time_zone": types.StringType,
	}
}

type userLocaleInfoModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Locale      types.String `tfsdk:"locale"`
}

func (m userLocaleInfoModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"display_name": types.StringType,
		"locale":       types.StringType,
	}
}

type userWorkingHoursModel struct {
	DaysOfWeek types.List   `tfsdk:"days_of_week"`
	EndTime    types.String `tfsdk:"end_time"`
	StartTime  types.String `tfsdk:"start_time"`
	TimeZone   types.Object `tfsdk:"time_zone"`
}

func (m userWorkingHoursModel) AttributeTypes() map[string]attr.Type {
	userTimeZone := userTimeZoneBaseModel{}
	return map[string]attr.Type{
		"days_of_week": types.ListType{ElemType: types.StringType},
		"end_time":     types.StringType,
		"start_time":   types.StringType,
		"time_zone":    types.ObjectType{AttrTypes: userTimeZone.AttributeTypes()},
	}
}

type userTimeZoneBaseModel struct {
	Name types.String `tfsdk:"name"`
}

func (m userTimeZoneBaseModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
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

type userOnPremisesProvisioningErrorModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

func (m userOnPremisesProvisioningErrorModel) AttributeTypes() map[string]attr.Type {
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

type userUserPrintModel struct {
}

func (m userUserPrintModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{}
}

type userProvisionedPlanModel struct {
	CapabilityStatus   types.String `tfsdk:"capability_status"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	Service            types.String `tfsdk:"service"`
}

func (m userProvisionedPlanModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"capability_status":   types.StringType,
		"provisioning_status": types.StringType,
		"service":             types.StringType,
	}
}

type userServiceProvisioningErrorModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

func (m userServiceProvisioningErrorModel) AttributeTypes() map[string]attr.Type {
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
