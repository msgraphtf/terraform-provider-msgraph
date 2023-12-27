package applications

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type applicationsModel struct {
	Value []applicationsValueModel `tfsdk:"value"`
}

type applicationsValueModel struct {
	Id                                types.String                                        `tfsdk:"id"`
	DeletedDateTime                   types.String                                        `tfsdk:"deleted_date_time"`
	AddIns                            []applicationsAddInsModel                           `tfsdk:"add_ins"`
	Api                               *applicationsApiModel                               `tfsdk:"api"`
	AppId                             types.String                                        `tfsdk:"app_id"`
	AppRoles                          []applicationsAppRolesModel                         `tfsdk:"app_roles"`
	ApplicationTemplateId             types.String                                        `tfsdk:"application_template_id"`
	Certification                     *applicationsCertificationModel                     `tfsdk:"certification"`
	CreatedDateTime                   types.String                                        `tfsdk:"created_date_time"`
	DefaultRedirectUri                types.String                                        `tfsdk:"default_redirect_uri"`
	Description                       types.String                                        `tfsdk:"description"`
	DisabledByMicrosoftStatus         types.String                                        `tfsdk:"disabled_by_microsoft_status"`
	DisplayName                       types.String                                        `tfsdk:"display_name"`
	GroupMembershipClaims             types.String                                        `tfsdk:"group_membership_claims"`
	IdentifierUris                    []types.String                                      `tfsdk:"identifier_uris"`
	Info                              *applicationsInfoModel                              `tfsdk:"info"`
	IsDeviceOnlyAuthSupported         types.Bool                                          `tfsdk:"is_device_only_auth_supported"`
	IsFallbackPublicClient            types.Bool                                          `tfsdk:"is_fallback_public_client"`
	KeyCredentials                    []applicationsKeyCredentialsModel                   `tfsdk:"key_credentials"`
	Logo                              types.String                                        `tfsdk:"logo"`
	Notes                             types.String                                        `tfsdk:"notes"`
	Oauth2RequirePostResponse         types.Bool                                          `tfsdk:"oauth_2_require_post_response"`
	OptionalClaims                    *applicationsOptionalClaimsModel                    `tfsdk:"optional_claims"`
	ParentalControlSettings           *applicationsParentalControlSettingsModel           `tfsdk:"parental_control_settings"`
	PasswordCredentials               []applicationsPasswordCredentialsModel              `tfsdk:"password_credentials"`
	PublicClient                      *applicationsPublicClientModel                      `tfsdk:"public_client"`
	PublisherDomain                   types.String                                        `tfsdk:"publisher_domain"`
	RequestSignatureVerification      *applicationsRequestSignatureVerificationModel      `tfsdk:"request_signature_verification"`
	RequiredResourceAccess            []applicationsRequiredResourceAccessModel           `tfsdk:"required_resource_access"`
	SamlMetadataUrl                   types.String                                        `tfsdk:"saml_metadata_url"`
	ServiceManagementReference        types.String                                        `tfsdk:"service_management_reference"`
	ServicePrincipalLockConfiguration *applicationsServicePrincipalLockConfigurationModel `tfsdk:"service_principal_lock_configuration"`
	SignInAudience                    types.String                                        `tfsdk:"sign_in_audience"`
	Spa                               *applicationsSpaModel                               `tfsdk:"spa"`
	Tags                              []types.String                                      `tfsdk:"tags"`
	TokenEncryptionKeyId              types.String                                        `tfsdk:"token_encryption_key_id"`
	VerifiedPublisher                 *applicationsVerifiedPublisherModel                 `tfsdk:"verified_publisher"`
	Web                               *applicationsWebModel                               `tfsdk:"web"`
}

type applicationsAddInsModel struct {
	Id         types.String                  `tfsdk:"id"`
	Properties []applicationsPropertiesModel `tfsdk:"properties"`
	Type       types.String                  `tfsdk:"type"`
}

type applicationsPropertiesModel struct {
	Key   types.String `tfsdk:"key"`
	Value types.String `tfsdk:"value"`
}

type applicationsApiModel struct {
	AcceptMappedClaims          types.Bool                                   `tfsdk:"accept_mapped_claims"`
	KnownClientApplications     []types.String                               `tfsdk:"known_client_applications"`
	Oauth2PermissionScopes      []applicationsOauth2PermissionScopesModel    `tfsdk:"oauth_2_permission_scopes"`
	PreAuthorizedApplications   []applicationsPreAuthorizedApplicationsModel `tfsdk:"pre_authorized_applications"`
	RequestedAccessTokenVersion types.Int64                                  `tfsdk:"requested_access_token_version"`
}

type applicationsOauth2PermissionScopesModel struct {
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

type applicationsPreAuthorizedApplicationsModel struct {
	AppId                  types.String   `tfsdk:"app_id"`
	DelegatedPermissionIds []types.String `tfsdk:"delegated_permission_ids"`
}

type applicationsAppRolesModel struct {
	AllowedMemberTypes []types.String `tfsdk:"allowed_member_types"`
	Description        types.String   `tfsdk:"description"`
	DisplayName        types.String   `tfsdk:"display_name"`
	Id                 types.String   `tfsdk:"id"`
	IsEnabled          types.Bool     `tfsdk:"is_enabled"`
	Origin             types.String   `tfsdk:"origin"`
	Value              types.String   `tfsdk:"value"`
}

type applicationsCertificationModel struct {
	CertificationDetailsUrl         types.String `tfsdk:"certification_details_url"`
	CertificationExpirationDateTime types.String `tfsdk:"certification_expiration_date_time"`
	IsCertifiedByMicrosoft          types.Bool   `tfsdk:"is_certified_by_microsoft"`
	IsPublisherAttested             types.Bool   `tfsdk:"is_publisher_attested"`
	LastCertificationDateTime       types.String `tfsdk:"last_certification_date_time"`
}

type applicationsInfoModel struct {
	LogoUrl             types.String `tfsdk:"logo_url"`
	MarketingUrl        types.String `tfsdk:"marketing_url"`
	PrivacyStatementUrl types.String `tfsdk:"privacy_statement_url"`
	SupportUrl          types.String `tfsdk:"support_url"`
	TermsOfServiceUrl   types.String `tfsdk:"terms_of_service_url"`
}

type applicationsKeyCredentialsModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Key                 types.String `tfsdk:"key"`
	KeyId               types.String `tfsdk:"key_id"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
	Type                types.String `tfsdk:"type"`
	Usage               types.String `tfsdk:"usage"`
}

type applicationsOptionalClaimsModel struct {
	AccessToken []applicationsAccessTokenModel `tfsdk:"access_token"`
	IdToken     []applicationsIdTokenModel     `tfsdk:"id_token"`
	Saml2Token  []applicationsSaml2TokenModel  `tfsdk:"saml_2_token"`
}

type applicationsAccessTokenModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationsIdTokenModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationsSaml2TokenModel struct {
	AdditionalProperties []types.String `tfsdk:"additional_properties"`
	Essential            types.Bool     `tfsdk:"essential"`
	Name                 types.String   `tfsdk:"name"`
	Source               types.String   `tfsdk:"source"`
}

type applicationsParentalControlSettingsModel struct {
	CountriesBlockedForMinors []types.String `tfsdk:"countries_blocked_for_minors"`
	LegalAgeGroupRule         types.String   `tfsdk:"legal_age_group_rule"`
}

type applicationsPasswordCredentialsModel struct {
	CustomKeyIdentifier types.String `tfsdk:"custom_key_identifier"`
	DisplayName         types.String `tfsdk:"display_name"`
	EndDateTime         types.String `tfsdk:"end_date_time"`
	Hint                types.String `tfsdk:"hint"`
	KeyId               types.String `tfsdk:"key_id"`
	SecretText          types.String `tfsdk:"secret_text"`
	StartDateTime       types.String `tfsdk:"start_date_time"`
}

type applicationsPublicClientModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationsRequestSignatureVerificationModel struct {
	AllowedWeakAlgorithms   types.String `tfsdk:"allowed_weak_algorithms"`
	IsSignedRequestRequired types.Bool   `tfsdk:"is_signed_request_required"`
}

type applicationsRequiredResourceAccessModel struct {
	ResourceAccess []applicationsResourceAccessModel `tfsdk:"resource_access"`
	ResourceAppId  types.String                      `tfsdk:"resource_app_id"`
}

type applicationsResourceAccessModel struct {
	Id   types.String `tfsdk:"id"`
	Type types.String `tfsdk:"type"`
}

type applicationsServicePrincipalLockConfigurationModel struct {
	AllProperties              types.Bool `tfsdk:"all_properties"`
	CredentialsWithUsageSign   types.Bool `tfsdk:"credentials_with_usage_sign"`
	CredentialsWithUsageVerify types.Bool `tfsdk:"credentials_with_usage_verify"`
	IsEnabled                  types.Bool `tfsdk:"is_enabled"`
	TokenEncryptionKeyId       types.Bool `tfsdk:"token_encryption_key_id"`
}

type applicationsSpaModel struct {
	RedirectUris []types.String `tfsdk:"redirect_uris"`
}

type applicationsVerifiedPublisherModel struct {
	AddedDateTime       types.String `tfsdk:"added_date_time"`
	DisplayName         types.String `tfsdk:"display_name"`
	VerifiedPublisherId types.String `tfsdk:"verified_publisher_id"`
}

type applicationsWebModel struct {
	HomePageUrl           types.String                            `tfsdk:"home_page_url"`
	ImplicitGrantSettings *applicationsImplicitGrantSettingsModel `tfsdk:"implicit_grant_settings"`
	LogoutUrl             types.String                            `tfsdk:"logout_url"`
	RedirectUriSettings   []applicationsRedirectUriSettingsModel  `tfsdk:"redirect_uri_settings"`
	RedirectUris          []types.String                          `tfsdk:"redirect_uris"`
}

type applicationsImplicitGrantSettingsModel struct {
	EnableAccessTokenIssuance types.Bool `tfsdk:"enable_access_token_issuance"`
	EnableIdTokenIssuance     types.Bool `tfsdk:"enable_id_token_issuance"`
}

type applicationsRedirectUriSettingsModel struct {
	Index types.Int64  `tfsdk:"index"`
	Uri   types.String `tfsdk:"uri"`
}
