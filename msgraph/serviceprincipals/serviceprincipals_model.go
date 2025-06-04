package serviceprincipals

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type servicePrincipalsModel struct {
	Value types.List `tfsdk:"value"`
}

func (m servicePrincipalsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsServicePrincipalModel{}.AttributeTypes()}},
	}
}

type servicePrincipalsServicePrincipalModel struct {
	AccountEnabled                         types.Bool   `tfsdk:"account_enabled"`
	AddIns                                 types.List   `tfsdk:"add_ins"`
	AlternativeNames                       types.List   `tfsdk:"alternative_names"`
	AppDescription                         types.String `tfsdk:"app_description"`
	AppDisplayName                         types.String `tfsdk:"app_display_name"`
	AppId                                  types.String `tfsdk:"app_id"`
	AppOwnerOrganizationId                 types.String `tfsdk:"app_owner_organization_id"`
	AppRoleAssignmentRequired              types.Bool   `tfsdk:"app_role_assignment_required"`
	AppRoles                               types.List   `tfsdk:"app_roles"`
	ApplicationTemplateId                  types.String `tfsdk:"application_template_id"`
	DeletedDateTime                        types.String `tfsdk:"deleted_date_time"`
	Description                            types.String `tfsdk:"description"`
	DisabledByMicrosoftStatus              types.String `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                            types.String `tfsdk:"display_name"`
	Homepage                               types.String `tfsdk:"homepage"`
	Id                                     types.String `tfsdk:"id"`
	Info                                   types.Object `tfsdk:"info"`
	KeyCredentials                         types.List   `tfsdk:"key_credentials"`
	LoginUrl                               types.String `tfsdk:"login_url"`
	LogoutUrl                              types.String `tfsdk:"logout_url"`
	Notes                                  types.String `tfsdk:"notes"`
	NotificationEmailAddresses             types.List   `tfsdk:"notification_email_addresses"`
	Oauth2PermissionScopes                 types.List   `tfsdk:"oauth_2_permission_scopes"`
	PasswordCredentials                    types.List   `tfsdk:"password_credentials"`
	PreferredSingleSignOnMode              types.String `tfsdk:"preferred_single_sign_on_mode"`
	PreferredTokenSigningKeyThumbprint     types.String `tfsdk:"preferred_token_signing_key_thumbprint"`
	ReplyUrls                              types.List   `tfsdk:"reply_urls"`
	ResourceSpecificApplicationPermissions types.List   `tfsdk:"resource_specific_application_permissions"`
	SamlSingleSignOnSettings               types.Object `tfsdk:"saml_single_sign_on_settings"`
	ServicePrincipalNames                  types.List   `tfsdk:"service_principal_names"`
	ServicePrincipalType                   types.String `tfsdk:"service_principal_type"`
	SignInAudience                         types.String `tfsdk:"sign_in_audience"`
	Tags                                   types.List   `tfsdk:"tags"`
	TokenEncryptionKeyId                   types.String `tfsdk:"token_encryption_key_id"`
	VerifiedPublisher                      types.Object `tfsdk:"verified_publisher"`
}

func (m servicePrincipalsServicePrincipalModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_enabled":                        types.BoolType,
		"add_ins":                                types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsAddInModel{}.AttributeTypes()}},
		"alternative_names":                      types.ListType{ElemType: types.StringType},
		"app_description":                        types.StringType,
		"app_display_name":                       types.StringType,
		"app_id":                                 types.StringType,
		"app_owner_organization_id":              types.StringType,
		"app_role_assignment_required":           types.BoolType,
		"app_roles":                              types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsAppRoleModel{}.AttributeTypes()}},
		"application_template_id":                types.StringType,
		"deleted_date_time":                      types.StringType,
		"description":                            types.StringType,
		"disabled_by_microsoft_status":           types.StringType,
		"display_name":                           types.StringType,
		"homepage":                               types.StringType,
		"id":                                     types.StringType,
		"info":                                   types.ObjectType{AttrTypes: servicePrincipalsInformationalUrlModel{}.AttributeTypes()},
		"key_credentials":                        types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsKeyCredentialModel{}.AttributeTypes()}},
		"login_url":                              types.StringType,
		"logout_url":                             types.StringType,
		"notes":                                  types.StringType,
		"notification_email_addresses":           types.ListType{ElemType: types.StringType},
		"oauth_2_permission_scopes":              types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsPermissionScopeModel{}.AttributeTypes()}},
		"password_credentials":                   types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsPasswordCredentialModel{}.AttributeTypes()}},
		"preferred_single_sign_on_mode":          types.StringType,
		"preferred_token_signing_key_thumbprint": types.StringType,
		"reply_urls":                             types.ListType{ElemType: types.StringType},
		"resource_specific_application_permissions": types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsResourceSpecificPermissionModel{}.AttributeTypes()}},
		"saml_single_sign_on_settings":              types.ObjectType{AttrTypes: servicePrincipalsSamlSingleSignOnSettingsModel{}.AttributeTypes()},
		"service_principal_names":                   types.ListType{ElemType: types.StringType},
		"service_principal_type":                    types.StringType,
		"sign_in_audience":                          types.StringType,
		"tags":                                      types.ListType{ElemType: types.StringType},
		"token_encryption_key_id":                   types.StringType,
		"verified_publisher":                        types.ObjectType{AttrTypes: servicePrincipalsVerifiedPublisherModel{}.AttributeTypes()},
	}
}

type servicePrincipalsAddInModel struct {
	Id         types.String `tfsdk:"id"`
	Properties types.List   `tfsdk:"properties"`
	Type       types.String `tfsdk:"type"`
}

func (m servicePrincipalsAddInModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"properties": types.ListType{ElemType: types.ObjectType{AttrTypes: servicePrincipalsKeyValueModel{}.AttributeTypes()}},
		"type":       types.StringType,
	}
}

type servicePrincipalsAppRoleModel struct {
	AllowedMemberTypes types.List   `tfsdk:"allowed_member_types"`
	Description        types.String `tfsdk:"description"`
	DisplayName        types.String `tfsdk:"display_name"`
	Id                 types.String `tfsdk:"id"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Origin             types.String `tfsdk:"origin"`
	Value              types.String `tfsdk:"value"`
}

func (m servicePrincipalsAppRoleModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_member_types": types.ListType{ElemType: types.StringType},
		"description":          types.StringType,
		"display_name":         types.StringType,
		"id":                   types.StringType,
		"is_enabled":           types.BoolType,
		"origin":               types.StringType,
		"value":                types.StringType,
	}
}

type servicePrincipalsInformationalUrlModel struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

func (m servicePrincipalsInformationalUrlModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"logo_url":              types.StringType,
		"marketing_url":         types.StringType,
		"privacy_statement_url": types.StringType,
		"support_url":           types.StringType,
		"terms_of_service_url":  types.StringType,
	}
}

type servicePrincipalsKeyCredentialModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

func (m servicePrincipalsKeyCredentialModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"custom_key_identifier": types.StringType,
		"display_name":          types.StringType,
		"end_date_time":         types.StringType,
		"key":                   types.StringType,
		"key_id":                types.StringType,
		"start_date_time":       types.StringType,
		"type":                  types.StringType,
		"usage":                 types.StringType,
	}
}

type servicePrincipalsPermissionScopeModel struct {
	AdminConsentDescription types.String `tfsdk:"admin_consent_description"`
	AdminConsentDisplayName types.String `tfsdk:"admin_consent_display_name"`
	Id                      types.String `tfsdk:"id"`
	IsEnabled               types.Bool   `tfsdk:"is_enabled"`
	Origin                  types.String `tfsdk:"origin"`
	Type                    types.String `tfsdk:"type"`
	UserConsentDescription  types.String `tfsdk:"user_consent_description"`
	UserConsentDisplayName  types.String `tfsdk:"user_consent_display_name"`
	Value                   types.String `tfsdk:"value"`
}

func (m servicePrincipalsPermissionScopeModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"admin_consent_description":  types.StringType,
		"admin_consent_display_name": types.StringType,
		"id":                         types.StringType,
		"is_enabled":                 types.BoolType,
		"origin":                     types.StringType,
		"type":                       types.StringType,
		"user_consent_description":   types.StringType,
		"user_consent_display_name":  types.StringType,
		"value":                      types.StringType,
	}
}

type servicePrincipalsPasswordCredentialModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

func (m servicePrincipalsPasswordCredentialModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"custom_key_identifier": types.StringType,
		"display_name":          types.StringType,
		"end_date_time":         types.StringType,
		"hint":                  types.StringType,
		"key_id":                types.StringType,
		"secret_text":           types.StringType,
		"start_date_time":       types.StringType,
	}
}

type servicePrincipalsResourceSpecificPermissionModel struct {
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
	IsEnabled   types.Bool   `tfsdk:"is_enabled"`
	Value       types.String `tfsdk:"value"`
}

func (m servicePrincipalsResourceSpecificPermissionModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description":  types.StringType,
		"display_name": types.StringType,
		"id":           types.StringType,
		"is_enabled":   types.BoolType,
		"value":        types.StringType,
	}
}

type servicePrincipalsSamlSingleSignOnSettingsModel struct {
	RelayState types.String `tfsdk:"relay_state"`
}

func (m servicePrincipalsSamlSingleSignOnSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"relay_state": types.StringType,
	}
}

type servicePrincipalsVerifiedPublisherModel struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}

func (m servicePrincipalsVerifiedPublisherModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"added_date_time":       types.StringType,
		"display_name":          types.StringType,
		"verified_publisher_id": types.StringType,
	}
}

type servicePrincipalsKeyValueModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (m servicePrincipalsKeyValueModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}
