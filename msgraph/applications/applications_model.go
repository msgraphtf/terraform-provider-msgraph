package applications

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationsDataSourceModel struct {
	Value []applicationsValueDataSourceModel `tfsdk:"value"`
}

type applicationsValueDataSourceModel struct {
	Id                                types.String                                                  `tfsdk:"id"`
	DeletedDateTime                   types.String                                                  `tfsdk:"deleted_date_time"`
	AddIns                            []applicationsAddInsDataSourceModel                           `tfsdk:"add_ins"`
	Api                               *applicationsApiDataSourceModel                               `tfsdk:"api"`
	AppId                             types.String                                                  `tfsdk:"app_id"`
	AppRoles                          []applicationsAppRolesDataSourceModel                         `tfsdk:"app_roles"`
	ApplicationTemplateId             types.String                                                  `tfsdk:"application_template_id"`
	Certification                     *applicationsCertificationDataSourceModel                     `tfsdk:"certification"`
	CreatedDateTime                   types.String                                                  `tfsdk:"created_date_time"`
	DefaultRedirectUri                types.String                                                  `tfsdk:"default_redirect_uri"`
	Description                       types.String                                                  `tfsdk:"description"`
	DisabledByMicrosoftStatus         types.String                                                  `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                       types.String                                                  `tfsdk:"display_name"`
	GroupMembershipClaims             types.String                                                  `tfsdk:"group_membership_claims"`
	IdentifierUris                    []types.String                                                `tfsdk:"identifier_uris"`
	Info                              *applicationsInfoDataSourceModel                              `tfsdk:"info"`
	IsDeviceOnlyAuthSupported         types.Bool                                                    `tfsdk:"is_device_only_auth_supported"`
	IsFallbackPublicClient            types.Bool                                                    `tfsdk:"is_fallback_public_client"`
	KeyCredentials                    []applicationsKeyCredentialsDataSourceModel                   `tfsdk:"key_credentials"`
	Logo                              types.String                                                  `tfsdk:"logo"`
	Notes                             types.String                                                  `tfsdk:"notes"`
	Oauth2RequirePostResponse         types.Bool                                                    `tfsdk:"oauth_2_require_post_response"`
	OptionalClaims                    *applicationsOptionalClaimsDataSourceModel                    `tfsdk:"optional_claims"`
	ParentalControlSettings           *applicationsParentalControlSettingsDataSourceModel           `tfsdk:"parental_control_settings"`
	PasswordCredentials               []applicationsPasswordCredentialsDataSourceModel              `tfsdk:"password_credentials"`
	PublicClient                      *applicationsPublicClientDataSourceModel                      `tfsdk:"public_client"`
	PublisherDomain                   types.String                                                  `tfsdk:"publisher_domain"`
	RequestSignatureVerification      *applicationsRequestSignatureVerificationDataSourceModel      `tfsdk:"request_signature_verification"`
	RequiredResourceAccess            []applicationsRequiredResourceAccessDataSourceModel           `tfsdk:"required_resource_access"`
	SamlMetadataUrl                   types.String                                                  `tfsdk:"saml_metadata_url"`
	ServiceManagementReference        types.String                                                  `tfsdk:"service_management_reference"`
	ServicePrincipalLockConfiguration *applicationsServicePrincipalLockConfigurationDataSourceModel `tfsdk:"service_principal_lock_configuration"`
	SignInAudience                    types.String                                                  `tfsdk:"sign_in_audience"`
	Spa                               *applicationsSpaDataSourceModel                               `tfsdk:"spa"`
	Tags                              []types.String                                                `tfsdk:"tags"`
	TokenEncryptionKeyId              types.String                                                  `tfsdk:"token_encryption_key_id"`
	VerifiedPublisher                 *applicationsVerifiedPublisherDataSourceModel                 `tfsdk:"verified_publisher"`
	Web                               *applicationsWebDataSourceModel                               `tfsdk:"web"`
}

type applicationsAddInsDataSourceModel struct {
	Id         types.String                            `tfsdk:"id"`
	Properties []applicationsPropertiesDataSourceModel `tfsdk:"properties"`
	Type       types.String                            `tfsdk:"type"`
}

type applicationsPropertiesDataSourceModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type applicationsApiDataSourceModel struct {
	AcceptMappedClaims          types.Bool                                             `tfsdk:"accept_mapped_claims"`
	KnownClientApplications     []types.String                                         `tfsdk:"known_client_applications"`
	Oauth2PermissionScopes      []applicationsOauth2PermissionScopesDataSourceModel    `tfsdk:"oauth_2_permission_scopes"`
	PreAuthorizedApplications   []applicationsPreAuthorizedApplicationsDataSourceModel `tfsdk:"pre_authorized_applications"`
	RequestedAccessTokenVersion types.Int64                                            `tfsdk:"requested_access_token_version"`
}

type applicationsOauth2PermissionScopesDataSourceModel struct {
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

type applicationsPreAuthorizedApplicationsDataSourceModel struct {
	AppId                  types.String   `tfsdk:"app_id"`
	DelegatedPermissionIds []types.String `tfsdk:"delegated_permission_ids"`
}

type applicationsAppRolesDataSourceModel struct {
	AllowedMemberTypes []types.String `tfsdk:"allowed_member_types"`
	Description        types.String   `tfsdk:"description"`
	DisplayName        types.String   `tfsdk:"display_name"`
	Id                 types.String   `tfsdk:"id"`
	IsEnabled          types.Bool     `tfsdk:"is_enabled"`
	Origin             types.String   `tfsdk:"origin"`
	Value              types.String   `tfsdk:"value"`
}

type applicationsCertificationDataSourceModel struct {
	CertificationDetailsUrl         types.String `tfsdk:"certification_details_url"`
	CertificationExpirationDateTime types.String `tfsdk:"certification_expiration_date_time"`
	IsCertifiedByMicrosoft          types.Bool   `tfsdk:"is_certified_by_microsoft"`
	IsPublisherAttested             types.Bool   `tfsdk:"is_publisher_attested"`
	LastCertificationDateTime       types.String `tfsdk:"last_certification_date_time"`
}

type applicationsInfoDataSourceModel struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

type applicationsKeyCredentialsDataSourceModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

type applicationsOptionalClaimsDataSourceModel struct {
	AccessToken []applicationsAccessTokenDataSourceModel `tfsdk:"access_token"`
	IdToken     []applicationsIdTokenDataSourceModel     `tfsdk:"id_token"`
	Saml2Token  []applicationsSaml2TokenDataSourceModel  `tfsdk:"saml_2_token"`
}

type applicationsAccessTokenDataSourceModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationsIdTokenDataSourceModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationsSaml2TokenDataSourceModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationsParentalControlSettingsDataSourceModel struct {
	CountriesBlockedForMinors []types.String `tfsdk:"countries_blocked_for_minors"`
	LegalAgeGroupRule         types.String   `tfsdk:"legal_age_group_rule"`
}

type applicationsPasswordCredentialsDataSourceModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

type applicationsPublicClientDataSourceModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationsRequestSignatureVerificationDataSourceModel struct {
	AllowedWeakAlgorithms   types.String `tfsdk:"allowed_weak_algorithms"`
	IsSignedRequestRequired types.Bool   `tfsdk:"is_signed_request_required"`
}

type applicationsRequiredResourceAccessDataSourceModel struct {
	ResourceAccess []applicationsResourceAccessDataSourceModel `tfsdk:"resource_access"`
	ResourceAppId  types.String                                `tfsdk:"resource_app_id"`
}

type applicationsResourceAccessDataSourceModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type applicationsServicePrincipalLockConfigurationDataSourceModel struct {
	AllProperties              types.Bool `tfsdk:"all_properties"`
	CredentialsWithUsageSign   types.Bool `tfsdk:"credentials_with_usage_sign"`
	CredentialsWithUsageVerify types.Bool `tfsdk:"credentials_with_usage_verify"`
	IsEnabled                  types.Bool `tfsdk:"is_enabled"`
	TokenEncryptionKeyId       types.Bool `tfsdk:"token_encryption_key_id"`
}

type applicationsSpaDataSourceModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationsVerifiedPublisherDataSourceModel struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}

type applicationsWebDataSourceModel struct {
	HomePageUrl           types.String                                      `tfsdk:"home_page_url"`
	ImplicitGrantSettings *applicationsImplicitGrantSettingsDataSourceModel `tfsdk:"implicit_grant_settings"`
	LogoutUrl             types.String                                      `tfsdk:"logout_url"`
	RedirectUriSettings   []applicationsRedirectUriSettingsDataSourceModel  `tfsdk:"redirect_uri_settings"`
	RedirectUris          []types.String                                    `tfsdk:"redirect_uris"`
}

type applicationsImplicitGrantSettingsDataSourceModel struct {
	EnableAccessTokenIssuance types.Bool `tfsdk:"enable_access_token_issuance"`
	EnableIdTokenIssuance     types.Bool `tfsdk:"enable_id_token_issuance"`
}

type applicationsRedirectUriSettingsDataSourceModel struct {
	Index types.Int64  `tfsdk:"index"`
	Uri   types.String `tfsdk:"uri"`
}
