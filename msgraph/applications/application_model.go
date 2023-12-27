package applications

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationModel struct {
	Id                                types.String                                       `tfsdk:"id"`
	DeletedDateTime                   types.String                                       `tfsdk:"deleted_date_time"`
	AddIns                            []applicationAddInsModel                           `tfsdk:"add_ins"`
	Api                               *applicationApiModel                               `tfsdk:"api"`
	AppId                             types.String                                       `tfsdk:"app_id"`
	AppRoles                          []applicationAppRolesModel                         `tfsdk:"app_roles"`
	ApplicationTemplateId             types.String                                       `tfsdk:"application_template_id"`
	Certification                     *applicationCertificationModel                     `tfsdk:"certification"`
	CreatedDateTime                   types.String                                       `tfsdk:"created_date_time"`
	DefaultRedirectUri                types.String                                       `tfsdk:"default_redirect_uri"`
	Description                       types.String                                       `tfsdk:"description"`
	DisabledByMicrosoftStatus         types.String                                       `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                       types.String                                       `tfsdk:"display_name"`
	GroupMembershipClaims             types.String                                       `tfsdk:"group_membership_claims"`
	IdentifierUris                    []types.String                                     `tfsdk:"identifier_uris"`
	Info                              *applicationInfoModel                              `tfsdk:"info"`
	IsDeviceOnlyAuthSupported         types.Bool                                         `tfsdk:"is_device_only_auth_supported"`
	IsFallbackPublicClient            types.Bool                                         `tfsdk:"is_fallback_public_client"`
	KeyCredentials                    []applicationKeyCredentialsModel                   `tfsdk:"key_credentials"`
	Logo                              types.String                                       `tfsdk:"logo"`
	Notes                             types.String                                       `tfsdk:"notes"`
	Oauth2RequirePostResponse         types.Bool                                         `tfsdk:"oauth_2_require_post_response"`
	OptionalClaims                    *applicationOptionalClaimsModel                    `tfsdk:"optional_claims"`
	ParentalControlSettings           *applicationParentalControlSettingsModel           `tfsdk:"parental_control_settings"`
	PasswordCredentials               []applicationPasswordCredentialsModel              `tfsdk:"password_credentials"`
	PublicClient                      *applicationPublicClientModel                      `tfsdk:"public_client"`
	PublisherDomain                   types.String                                       `tfsdk:"publisher_domain"`
	RequestSignatureVerification      *applicationRequestSignatureVerificationModel      `tfsdk:"request_signature_verification"`
	RequiredResourceAccess            []applicationRequiredResourceAccessModel           `tfsdk:"required_resource_access"`
	SamlMetadataUrl                   types.String                                       `tfsdk:"saml_metadata_url"`
	ServiceManagementReference        types.String                                       `tfsdk:"service_management_reference"`
	ServicePrincipalLockConfiguration *applicationServicePrincipalLockConfigurationModel `tfsdk:"service_principal_lock_configuration"`
	SignInAudience                    types.String                                       `tfsdk:"sign_in_audience"`
	Spa                               *applicationSpaModel                               `tfsdk:"spa"`
	Tags                              []types.String                                     `tfsdk:"tags"`
	TokenEncryptionKeyId              types.String                                       `tfsdk:"token_encryption_key_id"`
	VerifiedPublisher                 *applicationVerifiedPublisherModel                 `tfsdk:"verified_publisher"`
	Web                               *applicationWebModel                               `tfsdk:"web"`
}

type applicationAddInsModel struct {
	Id         types.String                 `tfsdk:"id"`
	Properties []applicationPropertiesModel `tfsdk:"properties"`
	Type       types.String                 `tfsdk:"type"`
}

type applicationPropertiesModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type applicationApiModel struct {
	AcceptMappedClaims          types.Bool                                  `tfsdk:"accept_mapped_claims"`
	KnownClientApplications     []types.String                              `tfsdk:"known_client_applications"`
	Oauth2PermissionScopes      []applicationOauth2PermissionScopesModel    `tfsdk:"oauth_2_permission_scopes"`
	PreAuthorizedApplications   []applicationPreAuthorizedApplicationsModel `tfsdk:"pre_authorized_applications"`
	RequestedAccessTokenVersion types.Int64                                 `tfsdk:"requested_access_token_version"`
}

type applicationOauth2PermissionScopesModel struct {
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

type applicationPreAuthorizedApplicationsModel struct {
	AppId                  types.String   `tfsdk:"app_id"`
	DelegatedPermissionIds []types.String `tfsdk:"delegated_permission_ids"`
}

type applicationAppRolesModel struct {
	AllowedMemberTypes []types.String `tfsdk:"allowed_member_types"`
	Description        types.String   `tfsdk:"description"`
	DisplayName        types.String   `tfsdk:"display_name"`
	Id                 types.String   `tfsdk:"id"`
	IsEnabled          types.Bool     `tfsdk:"is_enabled"`
	Origin             types.String   `tfsdk:"origin"`
	Value              types.String   `tfsdk:"value"`
}

type applicationCertificationModel struct {
	CertificationDetailsUrl         types.String `tfsdk:"certification_details_url"`
	CertificationExpirationDateTime types.String `tfsdk:"certification_expiration_date_time"`
	IsCertifiedByMicrosoft          types.Bool   `tfsdk:"is_certified_by_microsoft"`
	IsPublisherAttested             types.Bool   `tfsdk:"is_publisher_attested"`
	LastCertificationDateTime       types.String `tfsdk:"last_certification_date_time"`
}

type applicationInfoModel struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

type applicationKeyCredentialsModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

type applicationOptionalClaimsModel struct {
	AccessToken []applicationAccessTokenModel `tfsdk:"access_token"`
	IdToken     []applicationIdTokenModel     `tfsdk:"id_token"`
	Saml2Token  []applicationSaml2TokenModel  `tfsdk:"saml_2_token"`
}

type applicationAccessTokenModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationIdTokenModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationSaml2TokenModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationParentalControlSettingsModel struct {
	CountriesBlockedForMinors []types.String `tfsdk:"countries_blocked_for_minors"`
	LegalAgeGroupRule         types.String   `tfsdk:"legal_age_group_rule"`
}

type applicationPasswordCredentialsModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

type applicationPublicClientModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationRequestSignatureVerificationModel struct {
	AllowedWeakAlgorithms   types.String `tfsdk:"allowed_weak_algorithms"`
	IsSignedRequestRequired types.Bool   `tfsdk:"is_signed_request_required"`
}

type applicationRequiredResourceAccessModel struct {
	ResourceAccess []applicationResourceAccessModel `tfsdk:"resource_access"`
	ResourceAppId  types.String                     `tfsdk:"resource_app_id"`
}

type applicationResourceAccessModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type applicationServicePrincipalLockConfigurationModel struct {
	AllProperties              types.Bool `tfsdk:"all_properties"`
	CredentialsWithUsageSign   types.Bool `tfsdk:"credentials_with_usage_sign"`
	CredentialsWithUsageVerify types.Bool `tfsdk:"credentials_with_usage_verify"`
	IsEnabled                  types.Bool `tfsdk:"is_enabled"`
	TokenEncryptionKeyId       types.Bool `tfsdk:"token_encryption_key_id"`
}

type applicationSpaModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationVerifiedPublisherModel struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}

type applicationWebModel struct {
	HomePageUrl           types.String                           `tfsdk:"home_page_url"`
	ImplicitGrantSettings *applicationImplicitGrantSettingsModel `tfsdk:"implicit_grant_settings"`
	LogoutUrl             types.String                           `tfsdk:"logout_url"`
	RedirectUriSettings   []applicationRedirectUriSettingsModel  `tfsdk:"redirect_uri_settings"`
	RedirectUris          []types.String                         `tfsdk:"redirect_uris"`
}

type applicationImplicitGrantSettingsModel struct {
	EnableAccessTokenIssuance types.Bool `tfsdk:"enable_access_token_issuance"`
	EnableIdTokenIssuance     types.Bool `tfsdk:"enable_id_token_issuance"`
}

type applicationRedirectUriSettingsModel struct {
	Index types.Int64  `tfsdk:"index"`
	Uri   types.String `tfsdk:"uri"`
}
