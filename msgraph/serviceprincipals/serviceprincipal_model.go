package serviceprincipals

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type servicePrincipalModel struct {
	Id                                     types.String                                                  `tfsdk:"id"`
	DeletedDateTime                        types.String                                                  `tfsdk:"deleted_date_time"`
	AccountEnabled                         types.Bool                                                    `tfsdk:"account_enabled"`
	AddIns                                 []servicePrincipalAddInsModel                                 `tfsdk:"add_ins"`
	AlternativeNames                       []types.String                                                `tfsdk:"alternative_names"`
	AppDescription                         types.String                                                  `tfsdk:"app_description"`
	AppDisplayName                         types.String                                                  `tfsdk:"app_display_name"`
	AppId                                  types.String                                                  `tfsdk:"app_id"`
	AppOwnerOrganizationId                 types.String                                                  `tfsdk:"app_owner_organization_id"`
	AppRoleAssignmentRequired              types.Bool                                                    `tfsdk:"app_role_assignment_required"`
	AppRoles                               []servicePrincipalAppRolesModel                               `tfsdk:"app_roles"`
	ApplicationTemplateId                  types.String                                                  `tfsdk:"application_template_id"`
	CustomSecurityAttributes               *servicePrincipalCustomSecurityAttributesModel                `tfsdk:"custom_security_attributes"`
	Description                            types.String                                                  `tfsdk:"description"`
	DisabledByMicrosoftStatus              types.String                                                  `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                            types.String                                                  `tfsdk:"display_name"`
	Homepage                               types.String                                                  `tfsdk:"homepage"`
	Info                                   *servicePrincipalInfoModel                                    `tfsdk:"info"`
	KeyCredentials                         []servicePrincipalKeyCredentialsModel                         `tfsdk:"key_credentials"`
	LoginUrl                               types.String                                                  `tfsdk:"login_url"`
	LogoutUrl                              types.String                                                  `tfsdk:"logout_url"`
	Notes                                  types.String                                                  `tfsdk:"notes"`
	NotificationEmailAddresses             []types.String                                                `tfsdk:"notification_email_addresses"`
	Oauth2PermissionScopes                 []servicePrincipalOauth2PermissionScopesModel                 `tfsdk:"oauth_2_permission_scopes"`
	PasswordCredentials                    []servicePrincipalPasswordCredentialsModel                    `tfsdk:"password_credentials"`
	PreferredSingleSignOnMode              types.String                                                  `tfsdk:"preferred_single_sign_on_mode"`
	PreferredTokenSigningKeyThumbprint     types.String                                                  `tfsdk:"preferred_token_signing_key_thumbprint"`
	ReplyUrls                              []types.String                                                `tfsdk:"reply_urls"`
	ResourceSpecificApplicationPermissions []servicePrincipalResourceSpecificApplicationPermissionsModel `tfsdk:"resource_specific_application_permissions"`
	SamlSingleSignOnSettings               *servicePrincipalSamlSingleSignOnSettingsModel                `tfsdk:"saml_single_sign_on_settings"`
	ServicePrincipalNames                  []types.String                                                `tfsdk:"service_principal_names"`
	ServicePrincipalType                   types.String                                                  `tfsdk:"service_principal_type"`
	SignInAudience                         types.String                                                  `tfsdk:"sign_in_audience"`
	Tags                                   []types.String                                                `tfsdk:"tags"`
	TokenEncryptionKeyId                   types.String                                                  `tfsdk:"token_encryption_key_id"`
	VerifiedPublisher                      *servicePrincipalVerifiedPublisherModel                       `tfsdk:"verified_publisher"`
}

type servicePrincipalAddInsModel struct {
	Id         types.String                      `tfsdk:"id"`
	Properties []servicePrincipalPropertiesModel `tfsdk:"properties"`
	Type       types.String                      `tfsdk:"type"`
}

type servicePrincipalPropertiesModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type servicePrincipalAppRolesModel struct {
	AllowedMemberTypes []types.String `tfsdk:"allowed_member_types"`
	Description        types.String   `tfsdk:"description"`
	DisplayName        types.String   `tfsdk:"display_name"`
	Id                 types.String   `tfsdk:"id"`
	IsEnabled          types.Bool     `tfsdk:"is_enabled"`
	Origin             types.String   `tfsdk:"origin"`
	Value              types.String   `tfsdk:"value"`
}

type servicePrincipalCustomSecurityAttributesModel struct {
}

type servicePrincipalInfoModel struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

type servicePrincipalKeyCredentialsModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

type servicePrincipalOauth2PermissionScopesModel struct {
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

type servicePrincipalPasswordCredentialsModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

type servicePrincipalResourceSpecificApplicationPermissionsModel struct {
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
	IsEnabled   types.Bool   `tfsdk:"is_enabled"`
	Value       types.String `tfsdk:"value"`
}

type servicePrincipalSamlSingleSignOnSettingsModel struct {
	RelayState types.String `tfsdk:"relay_state"`
}

type servicePrincipalVerifiedPublisherModel struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}
