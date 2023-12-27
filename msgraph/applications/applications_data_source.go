package applications

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &applicationsDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationsDataSource{}
)

// NewApplicationsDataSource is a helper function to simplify the provider implementation.
func NewApplicationsDataSource() datasource.DataSource {
	return &applicationsDataSource{}
}

// applicationsDataSource is the data source implementation.
type applicationsDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *applicationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_applications"
}

// Configure adds the provider configured client to the data source.
func (d *applicationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *applicationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"value": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier for an entity. Read-only.",
							Computed:    true,
						},
						"deleted_date_time": schema.StringAttribute{
							Description: "Date and time when this object was deleted. Always null when the object hasn't been deleted.",
							Computed:    true,
						},
						"add_ins": schema.ListNestedAttribute{
							Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams may set the addIns property for its 'FileHandler' functionality. This will let services like Office 365 call the application in the context of a document the user is working on.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"properties": schema.ListNestedAttribute{
										Description: "",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"key": schema.StringAttribute{
													Description: "Key for the key-value pair.",
													Computed:    true,
												},
												"value": schema.StringAttribute{
													Description: "Value for the key-value pair.",
													Computed:    true,
												},
											},
										},
									},
									"type": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"api": schema.SingleNestedAttribute{
							Description: "Specifies settings for an application that implements a web API.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"accept_mapped_claims": schema.BoolAttribute{
									Description: "When true, allows an application to use claims mapping without specifying a custom signing key.",
									Computed:    true,
								},
								"known_client_applications": schema.ListAttribute{
									Description: "Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app. If you set the appID of the client app to this value, the user only consents once to the client app. Microsoft Entra ID knows that consenting to the client means implicitly consenting to the web API and automatically provisions service principals for both APIs at the same time. Both the client and the web API app must be registered in the same tenant.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"oauth_2_permission_scopes": schema.ListNestedAttribute{
									Description: "The definition of the delegated permissions exposed by the web API represented by this application registration. These delegated permissions may be requested by a client application, and may be granted by users or administrators during consent. Delegated permissions are sometimes referred to as OAuth 2.0 scopes.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"admin_consent_description": schema.StringAttribute{
												Description: "A description of the delegated permissions, intended to be read by an administrator granting the permission on behalf of all users. This text appears in tenant-wide admin consent experiences.",
												Computed:    true,
											},
											"admin_consent_display_name": schema.StringAttribute{
												Description: "The permission's title, intended to be read by an administrator granting the permission on behalf of all users.",
												Computed:    true,
											},
											"id": schema.StringAttribute{
												Description: "Unique delegated permission identifier inside the collection of delegated permissions defined for a resource application.",
												Computed:    true,
											},
											"is_enabled": schema.BoolAttribute{
												Description: "When you create or update a permission, this property must be set to true (which is the default). To delete a permission, this property must first be set to false.  At that point, in a subsequent call, the permission may be removed.",
												Computed:    true,
											},
											"origin": schema.StringAttribute{
												Description: "",
												Computed:    true,
											},
											"type": schema.StringAttribute{
												Description: "The possible values are: User and Admin. Specifies whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator consent should always be required. While Microsoft Graph defines the default consent requirement for each permission, the tenant administrator may override the behavior in their organization (by allowing, restricting, or limiting user consent to this delegated permission). For more information, see Configure how users consent to applications.",
												Computed:    true,
											},
											"user_consent_description": schema.StringAttribute{
												Description: "A description of the delegated permissions, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
												Computed:    true,
											},
											"user_consent_display_name": schema.StringAttribute{
												Description: "A title for the permission, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
												Computed:    true,
											},
											"value": schema.StringAttribute{
												Description: "Specifies the value to include in the scp (scope) claim in access tokens. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with ..",
												Computed:    true,
											},
										},
									},
								},
								"pre_authorized_applications": schema.ListNestedAttribute{
									Description: "Lists the client applications that are preauthorized with the specified delegated permissions to access this application's APIs. Users aren't required to consent to any preauthorized application (for the permissions specified). However, any other permissions not listed in preAuthorizedApplications (requested through incremental consent for example) will require user consent.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"app_id": schema.StringAttribute{
												Description: "The unique identifier for the application.",
												Computed:    true,
											},
											"delegated_permission_ids": schema.ListAttribute{
												Description: "The unique identifier for the oauth2PermissionScopes the application requires.",
												Computed:    true,
												ElementType: types.StringType,
											},
										},
									},
								},
								"requested_access_token_version": schema.Int64Attribute{
									Description: "Specifies the access token version expected by this resource. This changes the version and format of the JWT produced independent of the endpoint or client used to request the access token.  The endpoint used, v1.0 or v2.0, is chosen by the client and only impacts the version of id_tokens. Resources need to explicitly configure requestedAccessTokenVersion to indicate the supported access token format.  Possible values for requestedAccessTokenVersion are 1, 2, or null. If the value is null, this defaults to 1, which corresponds to the v1.0 endpoint.  If signInAudience on the application is configured as AzureADandPersonalMicrosoftAccount or PersonalMicrosoftAccount, the value for this property must be 2.",
									Computed:    true,
								},
							},
						},
						"app_id": schema.StringAttribute{
							Description: "The unique identifier for the application that is assigned to an application by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports $filter (eq).",
							Computed:    true,
						},
						"app_roles": schema.ListNestedAttribute{
							Description: "The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"allowed_member_types": schema.ListAttribute{
										Description: "Specifies whether this app role can be assigned to users and groups (by setting to ['User']), to other application's (by setting to ['Application'], or both (by setting to ['User', 'Application']). App roles supporting assignment to other applications' service principals are also known as application permissions. The 'Application' value is only supported for app roles defined on application entities.",
										Computed:    true,
										ElementType: types.StringType,
									},
									"description": schema.StringAttribute{
										Description: "The description for the app role. This is displayed when the app role is being assigned and, if the app role functions as an application permission, during  consent experiences.",
										Computed:    true,
									},
									"display_name": schema.StringAttribute{
										Description: "Display name for the permission that appears in the app role assignment and consent experiences.",
										Computed:    true,
									},
									"id": schema.StringAttribute{
										Description: "Unique role identifier inside the appRoles collection. When creating a new app role, a new GUID identifier must be provided.",
										Computed:    true,
									},
									"is_enabled": schema.BoolAttribute{
										Description: "When creating or updating an app role, this must be set to true (which is the default). To delete a role, this must first be set to false.  At that point, in a subsequent call, this role may be removed.",
										Computed:    true,
									},
									"origin": schema.StringAttribute{
										Description: "Specifies if the app role is defined on the application object or on the servicePrincipal entity. Must not be included in any POST or PATCH requests. Read-only.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Specifies the value to include in the roles claim in ID tokens and access tokens authenticating an assigned user or service principal. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with ..",
										Computed:    true,
									},
								},
							},
						},
						"application_template_id": schema.StringAttribute{
							Description: "Unique identifier of the applicationTemplate. Supports $filter (eq, not, ne).",
							Computed:    true,
						},
						"certification": schema.SingleNestedAttribute{
							Description: "Specifies the certification status of the application.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"certification_details_url": schema.StringAttribute{
									Description: "URL that shows certification details for the application.",
									Computed:    true,
								},
								"certification_expiration_date_time": schema.StringAttribute{
									Description: "The timestamp when the current certification for the application expires.",
									Computed:    true,
								},
								"is_certified_by_microsoft": schema.BoolAttribute{
									Description: "Indicates whether the application is certified by Microsoft.",
									Computed:    true,
								},
								"is_publisher_attested": schema.BoolAttribute{
									Description: "Indicates whether the application has been self-attested by the application developer or the publisher.",
									Computed:    true,
								},
								"last_certification_date_time": schema.StringAttribute{
									Description: "The timestamp when the certification for the application was most recently added or updated.",
									Computed:    true,
								},
							},
						},
						"created_date_time": schema.StringAttribute{
							Description: "The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.  Supports $filter (eq, ne, not, ge, le, in, and eq on null values) and $orderby.",
							Computed:    true,
						},
						"default_redirect_uri": schema.StringAttribute{
							Description: "",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Free text field to provide a description of the application object to end users. The maximum allowed size is 1024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
							Computed:    true,
						},
						"disabled_by_microsoft_status": schema.StringAttribute{
							Description: "Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The display name for the application. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
							Computed:    true,
						},
						"group_membership_claims": schema.StringAttribute{
							Description: "Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following valid string values: None, SecurityGroup (for security groups and Microsoft Entra roles), All (this gets all of the security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).",
							Computed:    true,
						},
						"identifier_uris": schema.ListAttribute{
							Description: "Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you'll reference in your API's code, and it must be globally unique. You can use the default value provided, which is in the form api://<application-client-id>, or specify a more readable URI like https://contoso.com/api. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports $filter (eq, ne, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"info": schema.SingleNestedAttribute{
							Description: "Basic profile information of the application such as  app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"logo_url": schema.StringAttribute{
									Description: "CDN URL to the application's logo, Read-only.",
									Computed:    true,
								},
								"marketing_url": schema.StringAttribute{
									Description: "Link to the application's marketing page. For example, https://www.contoso.com/app/marketing",
									Computed:    true,
								},
								"privacy_statement_url": schema.StringAttribute{
									Description: "Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy",
									Computed:    true,
								},
								"support_url": schema.StringAttribute{
									Description: "Link to the application's support page. For example, https://www.contoso.com/app/support",
									Computed:    true,
								},
								"terms_of_service_url": schema.StringAttribute{
									Description: "Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice",
									Computed:    true,
								},
							},
						},
						"is_device_only_auth_supported": schema.BoolAttribute{
							Description: "Specifies whether this application supports device authentication without a user. The default is false.",
							Computed:    true,
						},
						"is_fallback_public_client": schema.BoolAttribute{
							Description: "Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID cannot determine the client application type. For example, the ROPC flow where it is configured without specifying a redirect URI. In those cases Microsoft Entra ID interprets the application type based on the value of this property.",
							Computed:    true,
						},
						"key_credentials": schema.ListNestedAttribute{
							Description: "The collection of key credentials associated with the application. Not nullable. Supports $filter (eq, not, ge, le).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"custom_key_identifier": schema.StringAttribute{
										Description: "A 40-character binary type that can be used to identify the credential. Optional. When not provided in the payload, defaults to the thumbprint of the certificate.",
										Computed:    true,
									},
									"display_name": schema.StringAttribute{
										Description: "Friendly name for the key. Optional.",
										Computed:    true,
									},
									"end_date_time": schema.StringAttribute{
										Description: "The date and time at which the credential expires. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
										Computed:    true,
									},
									"key": schema.StringAttribute{
										Description: "The certificate's raw data in byte array converted to Base64 string. Returned only on $select for a single object, that is, GET applications/{applicationId}?$select=keyCredentials or GET servicePrincipals/{servicePrincipalId}?$select=keyCredentials; otherwise, it is always null.  From a .cer certificate, you can read the key using the Convert.ToBase64String() method. For more information, see Get the certificate key.",
										Computed:    true,
									},
									"key_id": schema.StringAttribute{
										Description: "The unique identifier (GUID) for the key.",
										Computed:    true,
									},
									"start_date_time": schema.StringAttribute{
										Description: "The date and time at which the credential becomes valid.The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
										Computed:    true,
									},
									"type": schema.StringAttribute{
										Description: "The type of key credential; for example, Symmetric, AsymmetricX509Cert.",
										Computed:    true,
									},
									"usage": schema.StringAttribute{
										Description: "A string that describes the purpose for which the key can be used; for example, Verify.",
										Computed:    true,
									},
								},
							},
						},
						"logo": schema.StringAttribute{
							Description: "The main logo for the application. Not nullable.",
							Computed:    true,
						},
						"notes": schema.StringAttribute{
							Description: "Notes relevant for the management of the application.",
							Computed:    true,
						},
						"oauth_2_require_post_response": schema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"optional_claims": schema.SingleNestedAttribute{
							Description: "Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"access_token": schema.ListNestedAttribute{
									Description: "The optional claims returned in the JWT access token.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"additional_properties": schema.ListAttribute{
												Description: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"essential": schema.BoolAttribute{
												Description: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "The name of the optional claim.",
												Computed:    true,
											},
											"source": schema.StringAttribute{
												Description: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
												Computed:    true,
											},
										},
									},
								},
								"id_token": schema.ListNestedAttribute{
									Description: "The optional claims returned in the JWT ID token.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"additional_properties": schema.ListAttribute{
												Description: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"essential": schema.BoolAttribute{
												Description: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "The name of the optional claim.",
												Computed:    true,
											},
											"source": schema.StringAttribute{
												Description: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
												Computed:    true,
											},
										},
									},
								},
								"saml_2_token": schema.ListNestedAttribute{
									Description: "The optional claims returned in the SAML token.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"additional_properties": schema.ListAttribute{
												Description: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
												Computed:    true,
												ElementType: types.StringType,
											},
											"essential": schema.BoolAttribute{
												Description: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
												Computed:    true,
											},
											"name": schema.StringAttribute{
												Description: "The name of the optional claim.",
												Computed:    true,
											},
											"source": schema.StringAttribute{
												Description: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"parental_control_settings": schema.SingleNestedAttribute{
							Description: "Specifies parental control settings for an application.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"countries_blocked_for_minors": schema.ListAttribute{
									Description: "Specifies the two-letter ISO country codes. Access to the application will be blocked for minors from the countries specified in this list.",
									Computed:    true,
									ElementType: types.StringType,
								},
								"legal_age_group_rule": schema.StringAttribute{
									Description: "Specifies the legal age group rule that applies to users of the app. Can be set to one of the following values: ValueDescriptionAllowDefault. Enforces the legal minimum. This means parental consent is required for minors in the European Union and Korea.RequireConsentForPrivacyServicesEnforces the user to specify date of birth to comply with COPPA rules. RequireConsentForMinorsRequires parental consent for ages below 18, regardless of country minor rules.RequireConsentForKidsRequires parental consent for ages below 14, regardless of country minor rules.BlockMinorsBlocks minors from using the app.",
									Computed:    true,
								},
							},
						},
						"password_credentials": schema.ListNestedAttribute{
							Description: "The collection of password credentials associated with the application. Not nullable.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"custom_key_identifier": schema.StringAttribute{
										Description: "Do not use.",
										Computed:    true,
									},
									"display_name": schema.StringAttribute{
										Description: "Friendly name for the password. Optional.",
										Computed:    true,
									},
									"end_date_time": schema.StringAttribute{
										Description: "The date and time at which the password expires represented using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.",
										Computed:    true,
									},
									"hint": schema.StringAttribute{
										Description: "Contains the first three characters of the password. Read-only.",
										Computed:    true,
									},
									"key_id": schema.StringAttribute{
										Description: "The unique identifier for the password.",
										Computed:    true,
									},
									"secret_text": schema.StringAttribute{
										Description: "Read-only; Contains the strong passwords generated by Microsoft Entra ID that are 16-64 characters in length. The generated password value is only returned during the initial POST request to addPassword. There is no way to retrieve this password in the future.",
										Computed:    true,
									},
									"start_date_time": schema.StringAttribute{
										Description: "The date and time at which the password becomes valid. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.",
										Computed:    true,
									},
								},
							},
						},
						"public_client": schema.SingleNestedAttribute{
							Description: "Specifies settings for installed clients such as desktop or mobile devices.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"redirect_uris": schema.ListAttribute{
									Description: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"publisher_domain": schema.StringAttribute{
							Description: "The verified publisher domain for the application. Read-only. For more information, see How to: Configure an application's publisher domain. Supports $filter (eq, ne, ge, le, startsWith).",
							Computed:    true,
						},
						"request_signature_verification": schema.SingleNestedAttribute{
							Description: "Specifies whether this application requires Microsoft Entra ID to verify the signed authentication requests.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"allowed_weak_algorithms": schema.StringAttribute{
									Description: "Specifies which weak algorithms are allowed.  The possible values are: rsaSha1, unknownFutureValue.",
									Computed:    true,
								},
								"is_signed_request_required": schema.BoolAttribute{
									Description: "Specifies whether signed authentication requests for this application should be required.",
									Computed:    true,
								},
							},
						},
						"required_resource_access": schema.ListNestedAttribute{
							Description: "Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports $filter (eq, not, ge, le).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"resource_access": schema.ListNestedAttribute{
										Description: "The list of OAuth2.0 permission scopes and app roles that the application requires from the specified resource.",
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Description: "The unique identifier of an app role or delegated permission exposed by the resource application. For delegated permissions, this should match the id property of one of the delegated permissions in the oauth2PermissionScopes collection of the resource application's service principal. For app roles (application permissions), this should match the id property of an app role in the appRoles collection of the resource application's service principal.",
													Computed:    true,
												},
												"type": schema.StringAttribute{
													Description: "Specifies whether the id property references a delegated permission or an app role (application permission). The possible values are: Scope (for delegated permissions) or Role (for app roles).",
													Computed:    true,
												},
											},
										},
									},
									"resource_app_id": schema.StringAttribute{
										Description: "The unique identifier for the resource that the application requires access to. This should be equal to the appId declared on the target resource application.",
										Computed:    true,
									},
								},
							},
						},
						"saml_metadata_url": schema.StringAttribute{
							Description: "The URL where the service exposes SAML metadata for federation. This property is valid only for single-tenant applications. Nullable.",
							Computed:    true,
						},
						"service_management_reference": schema.StringAttribute{
							Description: "References application or service contact information from a Service or Asset Management database. Nullable.",
							Computed:    true,
						},
						"service_principal_lock_configuration": schema.SingleNestedAttribute{
							Description: "Specifies whether sensitive properties of a multi-tenant application should be locked for editing after the application is provisioned in a tenant. Nullable. null by default.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"all_properties": schema.BoolAttribute{
									Description: "Enables locking all sensitive properties. The sensitive properties are keyCredentials, passwordCredentials, and tokenEncryptionKeyId.",
									Computed:    true,
								},
								"credentials_with_usage_sign": schema.BoolAttribute{
									Description: "Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Sign.",
									Computed:    true,
								},
								"credentials_with_usage_verify": schema.BoolAttribute{
									Description: "Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Verify. This locks OAuth service principals.",
									Computed:    true,
								},
								"is_enabled": schema.BoolAttribute{
									Description: "Enables or disables service principal lock configuration. To allow the sensitive properties to be updated, update this property to false to disable the lock on the service principal.",
									Computed:    true,
								},
								"token_encryption_key_id": schema.BoolAttribute{
									Description: "Locks the tokenEncryptionKeyId property for modification on the service principal.",
									Computed:    true,
								},
							},
						},
						"sign_in_audience": schema.StringAttribute{
							Description: "Specifies the Microsoft accounts that are supported for the current application. The possible values are: AzureADMyOrg, AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount (default), and PersonalMicrosoftAccount. See more in the table. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you may need to change other properties first. For more information, see Validation differences for signInAudience.Supports $filter (eq, ne, not).",
							Computed:    true,
						},
						"spa": schema.SingleNestedAttribute{
							Description: "Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"redirect_uris": schema.ListAttribute{
									Description: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"tags": schema.ListAttribute{
							Description: "Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here will also appear in the tags property of any associated service principals.Supports $filter (eq, not, ge, le, startsWith) and $search.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"token_encryption_key_id": schema.StringAttribute{
							Description: "Specifies the keyId of a public key from the keyCredentials collection. When configured, Microsoft Entra ID encrypts all the tokens it emits by using the key this property points to. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.",
							Computed:    true,
						},
						"verified_publisher": schema.SingleNestedAttribute{
							Description: "Specifies the verified publisher of the application. For more information about how publisher verification helps support application security, trustworthiness, and compliance, see Publisher verification.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"added_date_time": schema.StringAttribute{
									Description: "The timestamp when the verified publisher was first added or most recently updated.",
									Computed:    true,
								},
								"display_name": schema.StringAttribute{
									Description: "The verified publisher name from the app publisher's Partner Center account.",
									Computed:    true,
								},
								"verified_publisher_id": schema.StringAttribute{
									Description: "The ID of the verified publisher from the app publisher's Partner Center account.",
									Computed:    true,
								},
							},
						},
						"web": schema.SingleNestedAttribute{
							Description: "Specifies settings for a web application.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"home_page_url": schema.StringAttribute{
									Description: "Home page or landing page of the application.",
									Computed:    true,
								},
								"implicit_grant_settings": schema.SingleNestedAttribute{
									Description: "Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"enable_access_token_issuance": schema.BoolAttribute{
											Description: "Specifies whether this web application can request an access token using the OAuth 2.0 implicit flow.",
											Computed:    true,
										},
										"enable_id_token_issuance": schema.BoolAttribute{
											Description: "Specifies whether this web application can request an ID token using the OAuth 2.0 implicit flow.",
											Computed:    true,
										},
									},
								},
								"logout_url": schema.StringAttribute{
									Description: "Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.",
									Computed:    true,
								},
								"redirect_uri_settings": schema.ListNestedAttribute{
									Description: "",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"index": schema.Int64Attribute{
												Description: "",
												Computed:    true,
											},
											"uri": schema.StringAttribute{
												Description: "",
												Computed:    true,
											},
										},
									},
								},
								"redirect_uris": schema.ListAttribute{
									Description: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *applicationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state applicationsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := applications.ApplicationsRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationsRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"deletedDateTime",
				"addIns",
				"api",
				"appId",
				"applicationTemplateId",
				"appRoles",
				"certification",
				"createdDateTime",
				"defaultRedirectUri",
				"description",
				"disabledByMicrosoftStatus",
				"displayName",
				"groupMembershipClaims",
				"identifierUris",
				"info",
				"isDeviceOnlyAuthSupported",
				"isFallbackPublicClient",
				"keyCredentials",
				"logo",
				"notes",
				"oauth2RequirePostResponse",
				"optionalClaims",
				"parentalControlSettings",
				"passwordCredentials",
				"publicClient",
				"publisherDomain",
				"requestSignatureVerification",
				"requiredResourceAccess",
				"samlMetadataUrl",
				"serviceManagementReference",
				"servicePrincipalLockConfiguration",
				"signInAudience",
				"spa",
				"tags",
				"tokenEncryptionKeyId",
				"verifiedPublisher",
				"web",
				"appManagementPolicies",
				"createdOnBehalfOf",
				"extensionProperties",
				"federatedIdentityCredentials",
				"homeRealmDiscoveryPolicies",
				"owners",
				"tokenIssuancePolicies",
				"tokenLifetimePolicies",
				"synchronization",
			},
		},
	}

	result, err := d.client.Applications().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting applications",
			err.Error(),
		)
		return
	}

	for _, v := range result.GetValue() {
		value := new(applicationsValueDataSourceModel)

		if v.GetId() != nil {
			value.Id = types.StringValue(*v.GetId())
		}
		if v.GetDeletedDateTime() != nil {
			value.DeletedDateTime = types.StringValue(v.GetDeletedDateTime().String())
		}
		for _, v := range v.GetAddIns() {
			addIns := new(applicationsAddInsDataSourceModel)

			if v.GetId() != nil {
				addIns.Id = types.StringValue(v.GetId().String())
			}
			for _, v := range v.GetProperties() {
				properties := new(applicationsPropertiesDataSourceModel)

				if v.GetKey() != nil {
					properties.Key = types.StringValue(*v.GetKey())
				}
				if v.GetValue() != nil {
					properties.Value = types.StringValue(*v.GetValue())
				}
				addIns.Properties = append(addIns.Properties, *properties)
			}
			if v.GetTypeEscaped() != nil {
				addIns.Type = types.StringValue(*v.GetTypeEscaped())
			}
			value.AddIns = append(value.AddIns, *addIns)
		}
		if v.GetApi() != nil {
			value.Api = new(applicationsApiDataSourceModel)

			if v.GetApi().GetAcceptMappedClaims() != nil {
				value.Api.AcceptMappedClaims = types.BoolValue(*v.GetApi().GetAcceptMappedClaims())
			}
			for _, v := range v.GetApi().GetKnownClientApplications() {
				value.Api.KnownClientApplications = append(value.Api.KnownClientApplications, types.StringValue(v.String()))
			}
			for _, v := range v.GetApi().GetOauth2PermissionScopes() {
				oauth2PermissionScopes := new(applicationsOauth2PermissionScopesDataSourceModel)

				if v.GetAdminConsentDescription() != nil {
					oauth2PermissionScopes.AdminConsentDescription = types.StringValue(*v.GetAdminConsentDescription())
				}
				if v.GetAdminConsentDisplayName() != nil {
					oauth2PermissionScopes.AdminConsentDisplayName = types.StringValue(*v.GetAdminConsentDisplayName())
				}
				if v.GetId() != nil {
					oauth2PermissionScopes.Id = types.StringValue(v.GetId().String())
				}
				if v.GetIsEnabled() != nil {
					oauth2PermissionScopes.IsEnabled = types.BoolValue(*v.GetIsEnabled())
				}
				if v.GetOrigin() != nil {
					oauth2PermissionScopes.Origin = types.StringValue(*v.GetOrigin())
				}
				if v.GetTypeEscaped() != nil {
					oauth2PermissionScopes.Type = types.StringValue(*v.GetTypeEscaped())
				}
				if v.GetUserConsentDescription() != nil {
					oauth2PermissionScopes.UserConsentDescription = types.StringValue(*v.GetUserConsentDescription())
				}
				if v.GetUserConsentDisplayName() != nil {
					oauth2PermissionScopes.UserConsentDisplayName = types.StringValue(*v.GetUserConsentDisplayName())
				}
				if v.GetValue() != nil {
					oauth2PermissionScopes.Value = types.StringValue(*v.GetValue())
				}
				value.Api.Oauth2PermissionScopes = append(value.Api.Oauth2PermissionScopes, *oauth2PermissionScopes)
			}
			for _, v := range v.GetApi().GetPreAuthorizedApplications() {
				preAuthorizedApplications := new(applicationsPreAuthorizedApplicationsDataSourceModel)

				if v.GetAppId() != nil {
					preAuthorizedApplications.AppId = types.StringValue(*v.GetAppId())
				}
				for _, v := range v.GetDelegatedPermissionIds() {
					preAuthorizedApplications.DelegatedPermissionIds = append(preAuthorizedApplications.DelegatedPermissionIds, types.StringValue(v))
				}
				value.Api.PreAuthorizedApplications = append(value.Api.PreAuthorizedApplications, *preAuthorizedApplications)
			}
			if v.GetApi().GetRequestedAccessTokenVersion() != nil {
				value.Api.RequestedAccessTokenVersion = types.Int64Value(int64(*v.GetApi().GetRequestedAccessTokenVersion()))
			}
		}
		if v.GetAppId() != nil {
			value.AppId = types.StringValue(*v.GetAppId())
		}
		for _, v := range v.GetAppRoles() {
			appRoles := new(applicationsAppRolesDataSourceModel)

			for _, v := range v.GetAllowedMemberTypes() {
				appRoles.AllowedMemberTypes = append(appRoles.AllowedMemberTypes, types.StringValue(v))
			}
			if v.GetDescription() != nil {
				appRoles.Description = types.StringValue(*v.GetDescription())
			}
			if v.GetDisplayName() != nil {
				appRoles.DisplayName = types.StringValue(*v.GetDisplayName())
			}
			if v.GetId() != nil {
				appRoles.Id = types.StringValue(v.GetId().String())
			}
			if v.GetIsEnabled() != nil {
				appRoles.IsEnabled = types.BoolValue(*v.GetIsEnabled())
			}
			if v.GetOrigin() != nil {
				appRoles.Origin = types.StringValue(*v.GetOrigin())
			}
			if v.GetValue() != nil {
				appRoles.Value = types.StringValue(*v.GetValue())
			}
			value.AppRoles = append(value.AppRoles, *appRoles)
		}
		if v.GetApplicationTemplateId() != nil {
			value.ApplicationTemplateId = types.StringValue(*v.GetApplicationTemplateId())
		}
		if v.GetCertification() != nil {
			value.Certification = new(applicationsCertificationDataSourceModel)

			if v.GetCertification().GetCertificationDetailsUrl() != nil {
				value.Certification.CertificationDetailsUrl = types.StringValue(*v.GetCertification().GetCertificationDetailsUrl())
			}
			if v.GetCertification().GetCertificationExpirationDateTime() != nil {
				value.Certification.CertificationExpirationDateTime = types.StringValue(v.GetCertification().GetCertificationExpirationDateTime().String())
			}
			if v.GetCertification().GetIsCertifiedByMicrosoft() != nil {
				value.Certification.IsCertifiedByMicrosoft = types.BoolValue(*v.GetCertification().GetIsCertifiedByMicrosoft())
			}
			if v.GetCertification().GetIsPublisherAttested() != nil {
				value.Certification.IsPublisherAttested = types.BoolValue(*v.GetCertification().GetIsPublisherAttested())
			}
			if v.GetCertification().GetLastCertificationDateTime() != nil {
				value.Certification.LastCertificationDateTime = types.StringValue(v.GetCertification().GetLastCertificationDateTime().String())
			}
		}
		if v.GetCreatedDateTime() != nil {
			value.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
		}
		if v.GetDefaultRedirectUri() != nil {
			value.DefaultRedirectUri = types.StringValue(*v.GetDefaultRedirectUri())
		}
		if v.GetDescription() != nil {
			value.Description = types.StringValue(*v.GetDescription())
		}
		if v.GetDisabledByMicrosoftStatus() != nil {
			value.DisabledByMicrosoftStatus = types.StringValue(*v.GetDisabledByMicrosoftStatus())
		}
		if v.GetDisplayName() != nil {
			value.DisplayName = types.StringValue(*v.GetDisplayName())
		}
		if v.GetGroupMembershipClaims() != nil {
			value.GroupMembershipClaims = types.StringValue(*v.GetGroupMembershipClaims())
		}
		for _, v := range v.GetIdentifierUris() {
			value.IdentifierUris = append(value.IdentifierUris, types.StringValue(v))
		}
		if v.GetInfo() != nil {
			value.Info = new(applicationsInfoDataSourceModel)

			if v.GetInfo().GetLogoUrl() != nil {
				value.Info.LogoUrl = types.StringValue(*v.GetInfo().GetLogoUrl())
			}
			if v.GetInfo().GetMarketingUrl() != nil {
				value.Info.MarketingUrl = types.StringValue(*v.GetInfo().GetMarketingUrl())
			}
			if v.GetInfo().GetPrivacyStatementUrl() != nil {
				value.Info.PrivacyStatementUrl = types.StringValue(*v.GetInfo().GetPrivacyStatementUrl())
			}
			if v.GetInfo().GetSupportUrl() != nil {
				value.Info.SupportUrl = types.StringValue(*v.GetInfo().GetSupportUrl())
			}
			if v.GetInfo().GetTermsOfServiceUrl() != nil {
				value.Info.TermsOfServiceUrl = types.StringValue(*v.GetInfo().GetTermsOfServiceUrl())
			}
		}
		if v.GetIsDeviceOnlyAuthSupported() != nil {
			value.IsDeviceOnlyAuthSupported = types.BoolValue(*v.GetIsDeviceOnlyAuthSupported())
		}
		if v.GetIsFallbackPublicClient() != nil {
			value.IsFallbackPublicClient = types.BoolValue(*v.GetIsFallbackPublicClient())
		}
		for _, v := range v.GetKeyCredentials() {
			keyCredentials := new(applicationsKeyCredentialsDataSourceModel)

			if v.GetCustomKeyIdentifier() != nil {
				keyCredentials.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
			}
			if v.GetDisplayName() != nil {
				keyCredentials.DisplayName = types.StringValue(*v.GetDisplayName())
			}
			if v.GetEndDateTime() != nil {
				keyCredentials.EndDateTime = types.StringValue(v.GetEndDateTime().String())
			}
			if v.GetKey() != nil {
				keyCredentials.Key = types.StringValue(string(v.GetKey()[:]))
			}
			if v.GetKeyId() != nil {
				keyCredentials.KeyId = types.StringValue(v.GetKeyId().String())
			}
			if v.GetStartDateTime() != nil {
				keyCredentials.StartDateTime = types.StringValue(v.GetStartDateTime().String())
			}
			if v.GetTypeEscaped() != nil {
				keyCredentials.Type = types.StringValue(*v.GetTypeEscaped())
			}
			if v.GetUsage() != nil {
				keyCredentials.Usage = types.StringValue(*v.GetUsage())
			}
			value.KeyCredentials = append(value.KeyCredentials, *keyCredentials)
		}
		if v.GetLogo() != nil {
			value.Logo = types.StringValue(string(v.GetLogo()[:]))
		}
		if v.GetNotes() != nil {
			value.Notes = types.StringValue(*v.GetNotes())
		}
		if v.GetOauth2RequirePostResponse() != nil {
			value.Oauth2RequirePostResponse = types.BoolValue(*v.GetOauth2RequirePostResponse())
		}
		if v.GetOptionalClaims() != nil {
			value.OptionalClaims = new(applicationsOptionalClaimsDataSourceModel)

			for _, v := range v.GetOptionalClaims().GetAccessToken() {
				accessToken := new(applicationsAccessTokenDataSourceModel)

				for _, v := range v.GetAdditionalProperties() {
					accessToken.AdditionalProperties = append(accessToken.AdditionalProperties, types.StringValue(v))
				}
				if v.GetEssential() != nil {
					accessToken.Essential = types.BoolValue(*v.GetEssential())
				}
				if v.GetName() != nil {
					accessToken.Name = types.StringValue(*v.GetName())
				}
				if v.GetSource() != nil {
					accessToken.Source = types.StringValue(*v.GetSource())
				}
				value.OptionalClaims.AccessToken = append(value.OptionalClaims.AccessToken, *accessToken)
			}
			for _, v := range v.GetOptionalClaims().GetIdToken() {
				idToken := new(applicationsIdTokenDataSourceModel)

				for _, v := range v.GetAdditionalProperties() {
					idToken.AdditionalProperties = append(idToken.AdditionalProperties, types.StringValue(v))
				}
				if v.GetEssential() != nil {
					idToken.Essential = types.BoolValue(*v.GetEssential())
				}
				if v.GetName() != nil {
					idToken.Name = types.StringValue(*v.GetName())
				}
				if v.GetSource() != nil {
					idToken.Source = types.StringValue(*v.GetSource())
				}
				value.OptionalClaims.IdToken = append(value.OptionalClaims.IdToken, *idToken)
			}
			for _, v := range v.GetOptionalClaims().GetSaml2Token() {
				saml2Token := new(applicationsSaml2TokenDataSourceModel)

				for _, v := range v.GetAdditionalProperties() {
					saml2Token.AdditionalProperties = append(saml2Token.AdditionalProperties, types.StringValue(v))
				}
				if v.GetEssential() != nil {
					saml2Token.Essential = types.BoolValue(*v.GetEssential())
				}
				if v.GetName() != nil {
					saml2Token.Name = types.StringValue(*v.GetName())
				}
				if v.GetSource() != nil {
					saml2Token.Source = types.StringValue(*v.GetSource())
				}
				value.OptionalClaims.Saml2Token = append(value.OptionalClaims.Saml2Token, *saml2Token)
			}
		}
		if v.GetParentalControlSettings() != nil {
			value.ParentalControlSettings = new(applicationsParentalControlSettingsDataSourceModel)

			for _, v := range v.GetParentalControlSettings().GetCountriesBlockedForMinors() {
				value.ParentalControlSettings.CountriesBlockedForMinors = append(value.ParentalControlSettings.CountriesBlockedForMinors, types.StringValue(v))
			}
			if v.GetParentalControlSettings().GetLegalAgeGroupRule() != nil {
				value.ParentalControlSettings.LegalAgeGroupRule = types.StringValue(*v.GetParentalControlSettings().GetLegalAgeGroupRule())
			}
		}
		for _, v := range v.GetPasswordCredentials() {
			passwordCredentials := new(applicationsPasswordCredentialsDataSourceModel)

			if v.GetCustomKeyIdentifier() != nil {
				passwordCredentials.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
			}
			if v.GetDisplayName() != nil {
				passwordCredentials.DisplayName = types.StringValue(*v.GetDisplayName())
			}
			if v.GetEndDateTime() != nil {
				passwordCredentials.EndDateTime = types.StringValue(v.GetEndDateTime().String())
			}
			if v.GetHint() != nil {
				passwordCredentials.Hint = types.StringValue(*v.GetHint())
			}
			if v.GetKeyId() != nil {
				passwordCredentials.KeyId = types.StringValue(v.GetKeyId().String())
			}
			if v.GetSecretText() != nil {
				passwordCredentials.SecretText = types.StringValue(*v.GetSecretText())
			}
			if v.GetStartDateTime() != nil {
				passwordCredentials.StartDateTime = types.StringValue(v.GetStartDateTime().String())
			}
			value.PasswordCredentials = append(value.PasswordCredentials, *passwordCredentials)
		}
		if v.GetPublicClient() != nil {
			value.PublicClient = new(applicationsPublicClientDataSourceModel)

			for _, v := range v.GetPublicClient().GetRedirectUris() {
				value.PublicClient.RedirectUris = append(value.PublicClient.RedirectUris, types.StringValue(v))
			}
		}
		if v.GetPublisherDomain() != nil {
			value.PublisherDomain = types.StringValue(*v.GetPublisherDomain())
		}
		if v.GetRequestSignatureVerification() != nil {
			value.RequestSignatureVerification = new(applicationsRequestSignatureVerificationDataSourceModel)

			if v.GetRequestSignatureVerification().GetAllowedWeakAlgorithms() != nil {
				value.RequestSignatureVerification.AllowedWeakAlgorithms = types.StringValue(v.GetRequestSignatureVerification().GetAllowedWeakAlgorithms().String())
			}
			if v.GetRequestSignatureVerification().GetIsSignedRequestRequired() != nil {
				value.RequestSignatureVerification.IsSignedRequestRequired = types.BoolValue(*v.GetRequestSignatureVerification().GetIsSignedRequestRequired())
			}
		}
		for _, v := range v.GetRequiredResourceAccess() {
			requiredResourceAccess := new(applicationsRequiredResourceAccessDataSourceModel)

			for _, v := range v.GetResourceAccess() {
				resourceAccess := new(applicationsResourceAccessDataSourceModel)

				if v.GetId() != nil {
					resourceAccess.Id = types.StringValue(v.GetId().String())
				}
				if v.GetTypeEscaped() != nil {
					resourceAccess.Type = types.StringValue(*v.GetTypeEscaped())
				}
				requiredResourceAccess.ResourceAccess = append(requiredResourceAccess.ResourceAccess, *resourceAccess)
			}
			if v.GetResourceAppId() != nil {
				requiredResourceAccess.ResourceAppId = types.StringValue(*v.GetResourceAppId())
			}
			value.RequiredResourceAccess = append(value.RequiredResourceAccess, *requiredResourceAccess)
		}
		if v.GetSamlMetadataUrl() != nil {
			value.SamlMetadataUrl = types.StringValue(*v.GetSamlMetadataUrl())
		}
		if v.GetServiceManagementReference() != nil {
			value.ServiceManagementReference = types.StringValue(*v.GetServiceManagementReference())
		}
		if v.GetServicePrincipalLockConfiguration() != nil {
			value.ServicePrincipalLockConfiguration = new(applicationsServicePrincipalLockConfigurationDataSourceModel)

			if v.GetServicePrincipalLockConfiguration().GetAllProperties() != nil {
				value.ServicePrincipalLockConfiguration.AllProperties = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetAllProperties())
			}
			if v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign() != nil {
				value.ServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign())
			}
			if v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify() != nil {
				value.ServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify())
			}
			if v.GetServicePrincipalLockConfiguration().GetIsEnabled() != nil {
				value.ServicePrincipalLockConfiguration.IsEnabled = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetIsEnabled())
			}
			if v.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId() != nil {
				value.ServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId())
			}
		}
		if v.GetSignInAudience() != nil {
			value.SignInAudience = types.StringValue(*v.GetSignInAudience())
		}
		if v.GetSpa() != nil {
			value.Spa = new(applicationsSpaDataSourceModel)

			for _, v := range v.GetSpa().GetRedirectUris() {
				value.Spa.RedirectUris = append(value.Spa.RedirectUris, types.StringValue(v))
			}
		}
		for _, v := range v.GetTags() {
			value.Tags = append(value.Tags, types.StringValue(v))
		}
		if v.GetTokenEncryptionKeyId() != nil {
			value.TokenEncryptionKeyId = types.StringValue(v.GetTokenEncryptionKeyId().String())
		}
		if v.GetVerifiedPublisher() != nil {
			value.VerifiedPublisher = new(applicationsVerifiedPublisherDataSourceModel)

			if v.GetVerifiedPublisher().GetAddedDateTime() != nil {
				value.VerifiedPublisher.AddedDateTime = types.StringValue(v.GetVerifiedPublisher().GetAddedDateTime().String())
			}
			if v.GetVerifiedPublisher().GetDisplayName() != nil {
				value.VerifiedPublisher.DisplayName = types.StringValue(*v.GetVerifiedPublisher().GetDisplayName())
			}
			if v.GetVerifiedPublisher().GetVerifiedPublisherId() != nil {
				value.VerifiedPublisher.VerifiedPublisherId = types.StringValue(*v.GetVerifiedPublisher().GetVerifiedPublisherId())
			}
		}
		if v.GetWeb() != nil {
			value.Web = new(applicationsWebDataSourceModel)

			if v.GetWeb().GetHomePageUrl() != nil {
				value.Web.HomePageUrl = types.StringValue(*v.GetWeb().GetHomePageUrl())
			}
			if v.GetWeb().GetImplicitGrantSettings() != nil {
				value.Web.ImplicitGrantSettings = new(applicationsImplicitGrantSettingsDataSourceModel)

				if v.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance() != nil {
					value.Web.ImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolValue(*v.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance())
				}
				if v.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance() != nil {
					value.Web.ImplicitGrantSettings.EnableIdTokenIssuance = types.BoolValue(*v.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance())
				}
			}
			if v.GetWeb().GetLogoutUrl() != nil {
				value.Web.LogoutUrl = types.StringValue(*v.GetWeb().GetLogoutUrl())
			}
			for _, v := range v.GetWeb().GetRedirectUriSettings() {
				redirectUriSettings := new(applicationsRedirectUriSettingsDataSourceModel)

				if v.GetIndex() != nil {
					redirectUriSettings.Index = types.Int64Value(int64(*v.GetIndex()))
				}
				if v.GetUri() != nil {
					redirectUriSettings.Uri = types.StringValue(*v.GetUri())
				}
				value.Web.RedirectUriSettings = append(value.Web.RedirectUriSettings, *redirectUriSettings)
			}
			for _, v := range v.GetWeb().GetRedirectUris() {
				value.Web.RedirectUris = append(value.Web.RedirectUris, types.StringValue(v))
			}
		}
		state.Value = append(state.Value, *value)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
