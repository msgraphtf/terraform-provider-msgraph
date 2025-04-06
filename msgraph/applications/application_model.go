package applications

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationModel struct {
	Id                                types.String `tfsdk:"id"`
	DeletedDateTime                   types.String `tfsdk:"deleted_date_time"`
	AddIns                            types.List   `tfsdk:"add_ins"`
	Api                               types.Object `tfsdk:"api"`
	AppId                             types.String `tfsdk:"app_id"`
	AppRoles                          types.List   `tfsdk:"app_roles"`
	ApplicationTemplateId             types.String `tfsdk:"application_template_id"`
	Certification                     types.Object `tfsdk:"certification"`
	CreatedDateTime                   types.String `tfsdk:"created_date_time"`
	DefaultRedirectUri                types.String `tfsdk:"default_redirect_uri"`
	Description                       types.String `tfsdk:"description"`
	DisabledByMicrosoftStatus         types.String `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                       types.String `tfsdk:"display_name"`
	GroupMembershipClaims             types.String `tfsdk:"group_membership_claims"`
	IdentifierUris                    types.List   `tfsdk:"identifier_uris"`
	Info                              types.Object `tfsdk:"info"`
	IsDeviceOnlyAuthSupported         types.Bool   `tfsdk:"is_device_only_auth_supported"`
	IsFallbackPublicClient            types.Bool   `tfsdk:"is_fallback_public_client"`
	KeyCredentials                    types.List   `tfsdk:"key_credentials"`
	Logo                              types.String `tfsdk:"logo"`
	NativeAuthenticationApisEnabled   types.String `tfsdk:"native_authentication_apis_enabled"`
	Notes                             types.String `tfsdk:"notes"`
	Oauth2RequirePostResponse         types.Bool   `tfsdk:"oauth_2_require_post_response"`
	OptionalClaims                    types.Object `tfsdk:"optional_claims"`
	ParentalControlSettings           types.Object `tfsdk:"parental_control_settings"`
	PasswordCredentials               types.List   `tfsdk:"password_credentials"`
	PublicClient                      types.Object `tfsdk:"public_client"`
	PublisherDomain                   types.String `tfsdk:"publisher_domain"`
	RequestSignatureVerification      types.Object `tfsdk:"request_signature_verification"`
	RequiredResourceAccess            types.List   `tfsdk:"required_resource_access"`
	SamlMetadataUrl                   types.String `tfsdk:"saml_metadata_url"`
	ServiceManagementReference        types.String `tfsdk:"service_management_reference"`
	ServicePrincipalLockConfiguration types.Object `tfsdk:"service_principal_lock_configuration"`
	SignInAudience                    types.String `tfsdk:"sign_in_audience"`
	Spa                               types.Object `tfsdk:"spa"`
	Tags                              types.List   `tfsdk:"tags"`
	TokenEncryptionKeyId              types.String `tfsdk:"token_encryption_key_id"`
	UniqueName                        types.String `tfsdk:"unique_name"`
	VerifiedPublisher                 types.Object `tfsdk:"verified_publisher"`
	Web                               types.Object `tfsdk:"web"`
}

func (m applicationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                                   types.StringType,
		"deleted_date_time":                    types.StringType,
		"add_ins":                              types.ListType{ElemType: types.ObjectType{AttrTypes: applicationAddInModel{}.AttributeTypes()}},
		"api":                                  types.ObjectType{AttrTypes: applicationApiApplicationModel{}.AttributeTypes()},
		"app_id":                               types.StringType,
		"app_roles":                            types.ListType{ElemType: types.ObjectType{AttrTypes: applicationAppRoleModel{}.AttributeTypes()}},
		"application_template_id":              types.StringType,
		"certification":                        types.ObjectType{AttrTypes: applicationCertificationModel{}.AttributeTypes()},
		"created_date_time":                    types.StringType,
		"default_redirect_uri":                 types.StringType,
		"description":                          types.StringType,
		"disabled_by_microsoft_status":         types.StringType,
		"display_name":                         types.StringType,
		"group_membership_claims":              types.StringType,
		"identifier_uris":                      types.ListType{ElemType: types.StringType},
		"info":                                 types.ObjectType{AttrTypes: applicationInformationalUrlModel{}.AttributeTypes()},
		"is_device_only_auth_supported":        types.BoolType,
		"is_fallback_public_client":            types.BoolType,
		"key_credentials":                      types.ListType{ElemType: types.ObjectType{AttrTypes: applicationKeyCredentialModel{}.AttributeTypes()}},
		"logo":                                 types.StringType,
		"native_authentication_apis_enabled":   types.StringType,
		"notes":                                types.StringType,
		"oauth_2_require_post_response":        types.BoolType,
		"optional_claims":                      types.ObjectType{AttrTypes: applicationOptionalClaimsModel{}.AttributeTypes()},
		"parental_control_settings":            types.ObjectType{AttrTypes: applicationParentalControlSettingsModel{}.AttributeTypes()},
		"password_credentials":                 types.ListType{ElemType: types.ObjectType{AttrTypes: applicationPasswordCredentialModel{}.AttributeTypes()}},
		"public_client":                        types.ObjectType{AttrTypes: applicationPublicClientApplicationModel{}.AttributeTypes()},
		"publisher_domain":                     types.StringType,
		"request_signature_verification":       types.ObjectType{AttrTypes: applicationRequestSignatureVerificationModel{}.AttributeTypes()},
		"required_resource_access":             types.ListType{ElemType: types.ObjectType{AttrTypes: applicationRequiredResourceAccessModel{}.AttributeTypes()}},
		"saml_metadata_url":                    types.StringType,
		"service_management_reference":         types.StringType,
		"service_principal_lock_configuration": types.ObjectType{AttrTypes: applicationServicePrincipalLockConfigurationModel{}.AttributeTypes()},
		"sign_in_audience":                     types.StringType,
		"spa":                                  types.ObjectType{AttrTypes: applicationSpaApplicationModel{}.AttributeTypes()},
		"tags":                                 types.ListType{ElemType: types.StringType},
		"token_encryption_key_id":              types.StringType,
		"unique_name":                          types.StringType,
		"verified_publisher":                   types.ObjectType{AttrTypes: applicationVerifiedPublisherModel{}.AttributeTypes()},
		"web":                                  types.ObjectType{AttrTypes: applicationWebApplicationModel{}.AttributeTypes()},
	}
}

type applicationAddInModel struct {
	Id         types.String `tfsdk:"id"`
	Properties types.List   `tfsdk:"properties"`
	Type       types.String `tfsdk:"type"`
}

func (m applicationAddInModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":         types.StringType,
		"properties": types.ListType{ElemType: types.ObjectType{AttrTypes: applicationKeyValueModel{}.AttributeTypes()}},
		"type":       types.StringType,
	}
}

type applicationKeyValueModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

func (m applicationKeyValueModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"key":   types.StringType,
		"value": types.StringType,
	}
}

type applicationApiApplicationModel struct {
	AcceptMappedClaims          types.Bool  `tfsdk:"accept_mapped_claims"`
	KnownClientApplications     types.List  `tfsdk:"known_client_applications"`
	Oauth2PermissionScopes      types.List  `tfsdk:"oauth_2_permission_scopes"`
	PreAuthorizedApplications   types.List  `tfsdk:"pre_authorized_applications"`
	RequestedAccessTokenVersion types.Int64 `tfsdk:"requested_access_token_version"`
}

func (m applicationApiApplicationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"accept_mapped_claims":           types.BoolType,
		"known_client_applications":      types.ListType{ElemType: types.StringType},
		"oauth_2_permission_scopes":      types.ListType{ElemType: types.ObjectType{AttrTypes: applicationPermissionScopeModel{}.AttributeTypes()}},
		"pre_authorized_applications":    types.ListType{ElemType: types.ObjectType{AttrTypes: applicationPreAuthorizedApplicationModel{}.AttributeTypes()}},
		"requested_access_token_version": types.Int64Type,
	}
}

type applicationPermissionScopeModel struct {
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

func (m applicationPermissionScopeModel) AttributeTypes() map[string]attr.Type {
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

type applicationPreAuthorizedApplicationModel struct {
	AppId                  types.String `tfsdk:"app_id"`
	DelegatedPermissionIds types.List   `tfsdk:"delegated_permission_ids"`
}

func (m applicationPreAuthorizedApplicationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"app_id":                   types.StringType,
		"delegated_permission_ids": types.ListType{ElemType: types.StringType},
	}
}

type applicationAppRoleModel struct {
	AllowedMemberTypes types.List   `tfsdk:"allowed_member_types"`
	Description        types.String `tfsdk:"description"`
	DisplayName        types.String `tfsdk:"display_name"`
	Id                 types.String `tfsdk:"id"`
	IsEnabled          types.Bool   `tfsdk:"is_enabled"`
	Origin             types.String `tfsdk:"origin"`
	Value              types.String `tfsdk:"value"`
}

func (m applicationAppRoleModel) AttributeTypes() map[string]attr.Type {
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

type applicationCertificationModel struct {
	CertificationDetailsUrl         types.String `tfsdk:"certification_details_url"`
	CertificationExpirationDateTime types.String `tfsdk:"certification_expiration_date_time"`
	IsCertifiedByMicrosoft          types.Bool   `tfsdk:"is_certified_by_microsoft"`
	IsPublisherAttested             types.Bool   `tfsdk:"is_publisher_attested"`
	LastCertificationDateTime       types.String `tfsdk:"last_certification_date_time"`
}

func (m applicationCertificationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"certification_details_url":          types.StringType,
		"certification_expiration_date_time": types.StringType,
		"is_certified_by_microsoft":          types.BoolType,
		"is_publisher_attested":              types.BoolType,
		"last_certification_date_time":       types.StringType,
	}
}

type applicationInformationalUrlModel struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

func (m applicationInformationalUrlModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"logo_url":              types.StringType,
		"marketing_url":         types.StringType,
		"privacy_statement_url": types.StringType,
		"support_url":           types.StringType,
		"terms_of_service_url":  types.StringType,
	}
}

type applicationKeyCredentialModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

func (m applicationKeyCredentialModel) AttributeTypes() map[string]attr.Type {
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

type applicationOptionalClaimsModel struct {
	AccessToken types.List `tfsdk:"access_token"`
	IdToken     types.List `tfsdk:"id_token"`
	Saml2Token  types.List `tfsdk:"saml_2_token"`
}

func (m applicationOptionalClaimsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"access_token": types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOptionalClaimModel{}.AttributeTypes()}},
		"id_token":     types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOptionalClaimModel{}.AttributeTypes()}},
		"saml_2_token": types.ListType{ElemType: types.ObjectType{AttrTypes: applicationOptionalClaimModel{}.AttributeTypes()}},
	}
}

type applicationOptionalClaimModel struct {
	AdditionalProperties types.List   `tfsdk:"additional_properties"`
	Essential            types.Bool   `tfsdk:"essential"`
	Name                 types.String `tfsdk:"name"`
	Source               types.String `tfsdk:"source"`
}

func (m applicationOptionalClaimModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"additional_properties": types.ListType{ElemType: types.StringType},
		"essential":             types.BoolType,
		"name":                  types.StringType,
		"source":                types.StringType,
	}
}

type applicationParentalControlSettingsModel struct {
	CountriesBlockedForMinors types.List   `tfsdk:"countries_blocked_for_minors"`
	LegalAgeGroupRule         types.String `tfsdk:"legal_age_group_rule"`
}

func (m applicationParentalControlSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"countries_blocked_for_minors": types.ListType{ElemType: types.StringType},
		"legal_age_group_rule":         types.StringType,
	}
}

type applicationPasswordCredentialModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

func (m applicationPasswordCredentialModel) AttributeTypes() map[string]attr.Type {
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

type applicationPublicClientApplicationModel struct {
	RedirectUris types.List `tfsdk:"redirect_uris"`
}

func (m applicationPublicClientApplicationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"redirect_uris": types.ListType{ElemType: types.StringType},
	}
}

type applicationRequestSignatureVerificationModel struct {
	AllowedWeakAlgorithms   types.String `tfsdk:"allowed_weak_algorithms"`
	IsSignedRequestRequired types.Bool   `tfsdk:"is_signed_request_required"`
}

func (m applicationRequestSignatureVerificationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_weak_algorithms":    types.StringType,
		"is_signed_request_required": types.BoolType,
	}
}

type applicationRequiredResourceAccessModel struct {
	ResourceAccess types.List   `tfsdk:"resource_access"`
	ResourceAppId  types.String `tfsdk:"resource_app_id"`
}

func (m applicationRequiredResourceAccessModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"resource_access": types.ListType{ElemType: types.ObjectType{AttrTypes: applicationResourceAccessModel{}.AttributeTypes()}},
		"resource_app_id": types.StringType,
	}
}

type applicationResourceAccessModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

func (m applicationResourceAccessModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"type": types.StringType,
	}
}

type applicationServicePrincipalLockConfigurationModel struct {
	AllProperties              types.Bool `tfsdk:"all_properties"`
	CredentialsWithUsageSign   types.Bool `tfsdk:"credentials_with_usage_sign"`
	CredentialsWithUsageVerify types.Bool `tfsdk:"credentials_with_usage_verify"`
	IsEnabled                  types.Bool `tfsdk:"is_enabled"`
	TokenEncryptionKeyId       types.Bool `tfsdk:"token_encryption_key_id"`
}

func (m applicationServicePrincipalLockConfigurationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"all_properties":                types.BoolType,
		"credentials_with_usage_sign":   types.BoolType,
		"credentials_with_usage_verify": types.BoolType,
		"is_enabled":                    types.BoolType,
		"token_encryption_key_id":       types.BoolType,
	}
}

type applicationSpaApplicationModel struct {
	RedirectUris types.List `tfsdk:"redirect_uris"`
}

func (m applicationSpaApplicationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"redirect_uris": types.ListType{ElemType: types.StringType},
	}
}

type applicationVerifiedPublisherModel struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}

func (m applicationVerifiedPublisherModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"added_date_time":       types.StringType,
		"display_name":          types.StringType,
		"verified_publisher_id": types.StringType,
	}
}

type applicationWebApplicationModel struct {
	HomePageUrl           types.String `tfsdk:"home_page_url"`
	ImplicitGrantSettings types.Object `tfsdk:"implicit_grant_settings"`
	LogoutUrl             types.String `tfsdk:"logout_url"`
	RedirectUriSettings   types.List   `tfsdk:"redirect_uri_settings"`
	RedirectUris          types.List   `tfsdk:"redirect_uris"`
}

func (m applicationWebApplicationModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"home_page_url":           types.StringType,
		"implicit_grant_settings": types.ObjectType{AttrTypes: applicationImplicitGrantSettingsModel{}.AttributeTypes()},
		"logout_url":              types.StringType,
		"redirect_uri_settings":   types.ListType{ElemType: types.ObjectType{AttrTypes: applicationRedirectUriSettingsModel{}.AttributeTypes()}},
		"redirect_uris":           types.ListType{ElemType: types.StringType},
	}
}

type applicationImplicitGrantSettingsModel struct {
	EnableAccessTokenIssuance types.Bool `tfsdk:"enable_access_token_issuance"`
	EnableIdTokenIssuance     types.Bool `tfsdk:"enable_id_token_issuance"`
}

func (m applicationImplicitGrantSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enable_access_token_issuance": types.BoolType,
		"enable_id_token_issuance":     types.BoolType,
	}
}

type applicationRedirectUriSettingsModel struct {
	Index types.Int64  `tfsdk:"index"`
	Uri   types.String `tfsdk:"uri"`
}

func (m applicationRedirectUriSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"index": types.Int64Type,
		"uri":   types.StringType,
	}
}
