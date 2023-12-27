package applications

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationDataSourceModel struct {
	Id                                types.String                                                 `tfsdk:"id"`
	DeletedDateTime                   types.String                                                 `tfsdk:"deleted_date_time"`
	AddIns                            []applicationAddInsDataSourceModel                           `tfsdk:"add_ins"`
	Api                               *applicationApiDataSourceModel                               `tfsdk:"api"`
	AppId                             types.String                                                 `tfsdk:"app_id"`
	AppRoles                          []applicationAppRolesDataSourceModel                         `tfsdk:"app_roles"`
	ApplicationTemplateId             types.String                                                 `tfsdk:"application_template_id"`
	Certification                     *applicationCertificationDataSourceModel                     `tfsdk:"certification"`
	CreatedDateTime                   types.String                                                 `tfsdk:"created_date_time"`
	DefaultRedirectUri                types.String                                                 `tfsdk:"default_redirect_uri"`
	Description                       types.String                                                 `tfsdk:"description"`
	DisabledByMicrosoftStatus         types.String                                                 `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                       types.String                                                 `tfsdk:"display_name"`
	GroupMembershipClaims             types.String                                                 `tfsdk:"group_membership_claims"`
	IdentifierUris                    []types.String                                               `tfsdk:"identifier_uris"`
	Info                              *applicationInfoDataSourceModel                              `tfsdk:"info"`
	IsDeviceOnlyAuthSupported         types.Bool                                                   `tfsdk:"is_device_only_auth_supported"`
	IsFallbackPublicClient            types.Bool                                                   `tfsdk:"is_fallback_public_client"`
	KeyCredentials                    []applicationKeyCredentialsDataSourceModel                   `tfsdk:"key_credentials"`
	Logo                              types.String                                                 `tfsdk:"logo"`
	Notes                             types.String                                                 `tfsdk:"notes"`
	Oauth2RequirePostResponse         types.Bool                                                   `tfsdk:"oauth_2_require_post_response"`
	OptionalClaims                    *applicationOptionalClaimsDataSourceModel                    `tfsdk:"optional_claims"`
	ParentalControlSettings           *applicationParentalControlSettingsDataSourceModel           `tfsdk:"parental_control_settings"`
	PasswordCredentials               []applicationPasswordCredentialsDataSourceModel              `tfsdk:"password_credentials"`
	PublicClient                      *applicationPublicClientDataSourceModel                      `tfsdk:"public_client"`
	PublisherDomain                   types.String                                                 `tfsdk:"publisher_domain"`
	RequestSignatureVerification      *applicationRequestSignatureVerificationDataSourceModel      `tfsdk:"request_signature_verification"`
	RequiredResourceAccess            []applicationRequiredResourceAccessDataSourceModel           `tfsdk:"required_resource_access"`
	SamlMetadataUrl                   types.String                                                 `tfsdk:"saml_metadata_url"`
	ServiceManagementReference        types.String                                                 `tfsdk:"service_management_reference"`
	ServicePrincipalLockConfiguration *applicationServicePrincipalLockConfigurationDataSourceModel `tfsdk:"service_principal_lock_configuration"`
	SignInAudience                    types.String                                                 `tfsdk:"sign_in_audience"`
	Spa                               *applicationSpaDataSourceModel                               `tfsdk:"spa"`
	Tags                              []types.String                                               `tfsdk:"tags"`
	TokenEncryptionKeyId              types.String                                                 `tfsdk:"token_encryption_key_id"`
	VerifiedPublisher                 *applicationVerifiedPublisherDataSourceModel                 `tfsdk:"verified_publisher"`
	Web                               *applicationWebDataSourceModel                               `tfsdk:"web"`
}

type applicationAddInsDataSourceModel struct {
	Id         types.String                           `tfsdk:"id"`
	Properties []applicationPropertiesDataSourceModel `tfsdk:"properties"`
	Type       types.String                           `tfsdk:"type"`
}

type applicationPropertiesDataSourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type applicationApiDataSourceModel struct {
	AcceptMappedClaims          types.Bool                                            `tfsdk:"accept_mapped_claims"`
	KnownClientApplications     []types.String                                        `tfsdk:"known_client_applications"`
	Oauth2PermissionScopes      []applicationOauth2PermissionScopesDataSourceModel    `tfsdk:"oauth_2_permission_scopes"`
	PreAuthorizedApplications   []applicationPreAuthorizedApplicationsDataSourceModel `tfsdk:"pre_authorized_applications"`
	RequestedAccessTokenVersion types.Int64                                           `tfsdk:"requested_access_token_version"`
}

type applicationOauth2PermissionScopesDataSourceModel struct {
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

type applicationPreAuthorizedApplicationsDataSourceModel struct {
	AppId                  types.String   `tfsdk:"app_id"`
	DelegatedPermissionIds []types.String `tfsdk:"delegated_permission_ids"`
}

type applicationAppRolesDataSourceModel struct {
	AllowedMemberTypes []types.String `tfsdk:"allowed_member_types"`
	Description        types.String   `tfsdk:"description"`
	DisplayName        types.String   `tfsdk:"display_name"`
	Id                 types.String   `tfsdk:"id"`
	IsEnabled          types.Bool     `tfsdk:"is_enabled"`
	Origin             types.String   `tfsdk:"origin"`
	Value              types.String   `tfsdk:"value"`
}

type applicationCertificationDataSourceModel struct {
	CertificationDetailsUrl         types.String `tfsdk:"certification_details_url"`
	CertificationExpirationDateTime types.String `tfsdk:"certification_expiration_date_time"`
	IsCertifiedByMicrosoft          types.Bool   `tfsdk:"is_certified_by_microsoft"`
	IsPublisherAttested             types.Bool   `tfsdk:"is_publisher_attested"`
	LastCertificationDateTime       types.String `tfsdk:"last_certification_date_time"`
}

type applicationInfoDataSourceModel struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

type applicationKeyCredentialsDataSourceModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

type applicationOptionalClaimsDataSourceModel struct {
	AccessToken []applicationAccessTokenDataSourceModel `tfsdk:"access_token"`
	IdToken     []applicationIdTokenDataSourceModel     `tfsdk:"id_token"`
	Saml2Token  []applicationSaml2TokenDataSourceModel  `tfsdk:"saml_2_token"`
}

type applicationAccessTokenDataSourceModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationIdTokenDataSourceModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationSaml2TokenDataSourceModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationParentalControlSettingsDataSourceModel struct {
	CountriesBlockedForMinors []types.String `tfsdk:"countries_blocked_for_minors"`
	LegalAgeGroupRule         types.String   `tfsdk:"legal_age_group_rule"`
}

type applicationPasswordCredentialsDataSourceModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

type applicationPublicClientDataSourceModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationRequestSignatureVerificationDataSourceModel struct {
	AllowedWeakAlgorithms   types.String `tfsdk:"allowed_weak_algorithms"`
	IsSignedRequestRequired types.Bool   `tfsdk:"is_signed_request_required"`
}

type applicationRequiredResourceAccessDataSourceModel struct {
	ResourceAccess []applicationResourceAccessDataSourceModel `tfsdk:"resource_access"`
	ResourceAppId  types.String                               `tfsdk:"resource_app_id"`
}

type applicationResourceAccessDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type applicationServicePrincipalLockConfigurationDataSourceModel struct {
	AllProperties              types.Bool `tfsdk:"all_properties"`
	CredentialsWithUsageSign   types.Bool `tfsdk:"credentials_with_usage_sign"`
	CredentialsWithUsageVerify types.Bool `tfsdk:"credentials_with_usage_verify"`
	IsEnabled                  types.Bool `tfsdk:"is_enabled"`
	TokenEncryptionKeyId       types.Bool `tfsdk:"token_encryption_key_id"`
}

type applicationSpaDataSourceModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationVerifiedPublisherDataSourceModel struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}

type applicationWebDataSourceModel struct {
	HomePageUrl           types.String                                     `tfsdk:"home_page_url"`
	ImplicitGrantSettings *applicationImplicitGrantSettingsDataSourceModel `tfsdk:"implicit_grant_settings"`
	LogoutUrl             types.String                                     `tfsdk:"logout_url"`
	RedirectUriSettings   []applicationRedirectUriSettingsDataSourceModel  `tfsdk:"redirect_uri_settings"`
	RedirectUris          []types.String                                   `tfsdk:"redirect_uris"`
}

type applicationImplicitGrantSettingsDataSourceModel struct {
	EnableAccessTokenIssuance types.Bool `tfsdk:"enable_access_token_issuance"`
	EnableIdTokenIssuance     types.Bool `tfsdk:"enable_id_token_issuance"`
}

type applicationRedirectUriSettingsDataSourceModel struct {
	Index types.Int64  `tfsdk:"index"`
	Uri   types.String `tfsdk:"uri"`
}
