package applications

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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
							Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams can set the addIns property for its 'FileHandler' functionality. This lets services like Microsoft 365 call the application in the context of a document the user is working on.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The unique identifier for the addIn object.",
										Computed:    true,
									},
									"properties": schema.ListNestedAttribute{
										Description: "The collection of key-value pairs that define parameters that the consuming service can use or call. You must specify this property when performing a POST or a PATCH operation on the addIns collection. Required.",
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
										Description: "The unique name for the functionality exposed by the app.",
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
							Description: "Unique identifier of the applicationTemplate. Supports $filter (eq, not, ne). Read-only. null if the app wasn't created from an application template.",
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
							Description: "Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
							Computed:    true,
						},
						"disabled_by_microsoft_status": schema.StringAttribute{
							Description: "Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).",
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
							Description: "Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique. You can use the default value provided, which is in the form api://<appId>, or specify a more readable URI like https://contoso.com/api. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports $filter (eq, ne, ge, le, startsWith).",
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
							Description: "Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where it's configured without specifying a redirect URI. In those cases, Microsoft Entra ID interprets the application type based on the value of this property.",
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
										Description: "The friendly name for the key, with a maximum length of 90 characters. Longer values are accepted but shortened. Optional.",
										Computed:    true,
									},
									"end_date_time": schema.StringAttribute{
										Description: "The date and time at which the credential expires. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
										Computed:    true,
									},
									"key": schema.StringAttribute{
										Description: "The certificate's raw data in byte array converted to Base64 string. Returned only on $select for a single object, that is, GET applications/{applicationId}?$select=keyCredentials or GET servicePrincipals/{servicePrincipalId}?$select=keyCredentials; otherwise, it's always null.  From a .cer certificate, you can read the key using the Convert.ToBase64String() method. For more information, see Get the certificate key.",
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
						"native_authentication_apis_enabled": schema.StringAttribute{
							Description: "Specifies whether the Native Authentication APIs are enabled for the application. The possible values are: none and all. Default is none. For more information, see Native Authentication.",
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
									Description: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent. For iOS and macOS apps, specify the value following the syntax msauth.{BUNDLEID}://auth, replacing '{BUNDLEID}'. For example, if the bundle ID is com.microsoft.identitysample.MSALiOS, the URI is msauth.com.microsoft.identitysample.MSALiOS://auth.",
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
							Description: "Specifies whether sensitive properties of a multitenant application should be locked for editing after the application is provisioned in a tenant. Nullable. null by default.",
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
							Description: "Specifies the Microsoft accounts that are supported for the current application. The possible values are: AzureADMyOrg (default), AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, and PersonalMicrosoftAccount. See more in the table. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you might need to change other properties first. For more information, see Validation differences for signInAudience.Supports $filter (eq, ne, not).",
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
						"unique_name": schema.StringAttribute{
							Description: "The unique identifier that can be assigned to an application and used as an alternate key. Immutable. Read-only.",
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
	var tfStateApplications applicationsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateApplications)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := applications.ApplicationsRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationsRequestBuilderGetQueryParameters{
			Select: []string{
				"value",
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

	if len(result.GetValue()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetValue() {
			tfStateValue := applicationsApplicationModel{}

			if v.GetId() != nil {
				tfStateValue.Id = types.StringValue(*v.GetId())
			} else {
				tfStateValue.Id = types.StringNull()
			}
			if v.GetDeletedDateTime() != nil {
				tfStateValue.DeletedDateTime = types.StringValue(v.GetDeletedDateTime().String())
			} else {
				tfStateValue.DeletedDateTime = types.StringNull()
			}
			if len(v.GetAddIns()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetAddIns() {
					tfStateAddIns := applicationsAddInModel{}

					if v.GetId() != nil {
						tfStateAddIns.Id = types.StringValue(v.GetId().String())
					} else {
						tfStateAddIns.Id = types.StringNull()
					}
					if len(v.GetProperties()) > 0 {
						objectValues := []basetypes.ObjectValue{}
						for _, v := range v.GetProperties() {
							tfStateProperties := applicationsKeyValueModel{}

							if v.GetKey() != nil {
								tfStateProperties.Key = types.StringValue(*v.GetKey())
							} else {
								tfStateProperties.Key = types.StringNull()
							}
							if v.GetValue() != nil {
								tfStateProperties.Value = types.StringValue(*v.GetValue())
							} else {
								tfStateProperties.Value = types.StringNull()
							}
							objectValue, _ := types.ObjectValueFrom(ctx, tfStateProperties.AttributeTypes(), tfStateProperties)
							objectValues = append(objectValues, objectValue)
						}
						tfStateAddIns.Properties, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
					}
					if v.GetTypeEscaped() != nil {
						tfStateAddIns.Type = types.StringValue(*v.GetTypeEscaped())
					} else {
						tfStateAddIns.Type = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateAddIns.AttributeTypes(), tfStateAddIns)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.AddIns, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetApi() != nil {
				tfStateApi := applicationsApiApplicationModel{}

				if v.GetApi().GetAcceptMappedClaims() != nil {
					tfStateApi.AcceptMappedClaims = types.BoolValue(*v.GetApi().GetAcceptMappedClaims())
				} else {
					tfStateApi.AcceptMappedClaims = types.BoolNull()
				}
				if len(v.GetApi().GetKnownClientApplications()) > 0 {
					var valueArrayKnownClientApplications []attr.Value
					for _, v := range v.GetApi().GetKnownClientApplications() {
						valueArrayKnownClientApplications = append(valueArrayKnownClientApplications, types.StringValue(v.String()))
					}
					tfStateApi.KnownClientApplications, _ = types.ListValue(types.StringType, valueArrayKnownClientApplications)
				} else {
					tfStateApi.KnownClientApplications = types.ListNull(types.StringType)
				}
				if len(v.GetApi().GetOauth2PermissionScopes()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, v := range v.GetApi().GetOauth2PermissionScopes() {
						tfStateOauth2PermissionScopes := applicationsPermissionScopeModel{}

						if v.GetAdminConsentDescription() != nil {
							tfStateOauth2PermissionScopes.AdminConsentDescription = types.StringValue(*v.GetAdminConsentDescription())
						} else {
							tfStateOauth2PermissionScopes.AdminConsentDescription = types.StringNull()
						}
						if v.GetAdminConsentDisplayName() != nil {
							tfStateOauth2PermissionScopes.AdminConsentDisplayName = types.StringValue(*v.GetAdminConsentDisplayName())
						} else {
							tfStateOauth2PermissionScopes.AdminConsentDisplayName = types.StringNull()
						}
						if v.GetId() != nil {
							tfStateOauth2PermissionScopes.Id = types.StringValue(v.GetId().String())
						} else {
							tfStateOauth2PermissionScopes.Id = types.StringNull()
						}
						if v.GetIsEnabled() != nil {
							tfStateOauth2PermissionScopes.IsEnabled = types.BoolValue(*v.GetIsEnabled())
						} else {
							tfStateOauth2PermissionScopes.IsEnabled = types.BoolNull()
						}
						if v.GetOrigin() != nil {
							tfStateOauth2PermissionScopes.Origin = types.StringValue(*v.GetOrigin())
						} else {
							tfStateOauth2PermissionScopes.Origin = types.StringNull()
						}
						if v.GetTypeEscaped() != nil {
							tfStateOauth2PermissionScopes.Type = types.StringValue(*v.GetTypeEscaped())
						} else {
							tfStateOauth2PermissionScopes.Type = types.StringNull()
						}
						if v.GetUserConsentDescription() != nil {
							tfStateOauth2PermissionScopes.UserConsentDescription = types.StringValue(*v.GetUserConsentDescription())
						} else {
							tfStateOauth2PermissionScopes.UserConsentDescription = types.StringNull()
						}
						if v.GetUserConsentDisplayName() != nil {
							tfStateOauth2PermissionScopes.UserConsentDisplayName = types.StringValue(*v.GetUserConsentDisplayName())
						} else {
							tfStateOauth2PermissionScopes.UserConsentDisplayName = types.StringNull()
						}
						if v.GetValue() != nil {
							tfStateOauth2PermissionScopes.Value = types.StringValue(*v.GetValue())
						} else {
							tfStateOauth2PermissionScopes.Value = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStateOauth2PermissionScopes.AttributeTypes(), tfStateOauth2PermissionScopes)
						objectValues = append(objectValues, objectValue)
					}
					tfStateApi.Oauth2PermissionScopes, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}
				if len(v.GetApi().GetPreAuthorizedApplications()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, v := range v.GetApi().GetPreAuthorizedApplications() {
						tfStatePreAuthorizedApplications := applicationsPreAuthorizedApplicationModel{}

						if v.GetAppId() != nil {
							tfStatePreAuthorizedApplications.AppId = types.StringValue(*v.GetAppId())
						} else {
							tfStatePreAuthorizedApplications.AppId = types.StringNull()
						}
						if len(v.GetDelegatedPermissionIds()) > 0 {
							var valueArrayDelegatedPermissionIds []attr.Value
							for _, v := range v.GetDelegatedPermissionIds() {
								valueArrayDelegatedPermissionIds = append(valueArrayDelegatedPermissionIds, types.StringValue(v))
							}
							listValue, _ := types.ListValue(types.StringType, valueArrayDelegatedPermissionIds)
							tfStatePreAuthorizedApplications.DelegatedPermissionIds = listValue
						} else {
							tfStatePreAuthorizedApplications.DelegatedPermissionIds = types.ListNull(types.StringType)
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStatePreAuthorizedApplications.AttributeTypes(), tfStatePreAuthorizedApplications)
						objectValues = append(objectValues, objectValue)
					}
					tfStateApi.PreAuthorizedApplications, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}

				tfStateValue.Api, _ = types.ObjectValueFrom(ctx, tfStateApi.AttributeTypes(), tfStateApi)
			}
			if v.GetAppId() != nil {
				tfStateValue.AppId = types.StringValue(*v.GetAppId())
			} else {
				tfStateValue.AppId = types.StringNull()
			}
			if len(v.GetAppRoles()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetAppRoles() {
					tfStateAppRoles := applicationsAppRoleModel{}

					if len(v.GetAllowedMemberTypes()) > 0 {
						var valueArrayAllowedMemberTypes []attr.Value
						for _, v := range v.GetAllowedMemberTypes() {
							valueArrayAllowedMemberTypes = append(valueArrayAllowedMemberTypes, types.StringValue(v))
						}
						listValue, _ := types.ListValue(types.StringType, valueArrayAllowedMemberTypes)
						tfStateAppRoles.AllowedMemberTypes = listValue
					} else {
						tfStateAppRoles.AllowedMemberTypes = types.ListNull(types.StringType)
					}
					if v.GetDescription() != nil {
						tfStateAppRoles.Description = types.StringValue(*v.GetDescription())
					} else {
						tfStateAppRoles.Description = types.StringNull()
					}
					if v.GetDisplayName() != nil {
						tfStateAppRoles.DisplayName = types.StringValue(*v.GetDisplayName())
					} else {
						tfStateAppRoles.DisplayName = types.StringNull()
					}
					if v.GetId() != nil {
						tfStateAppRoles.Id = types.StringValue(v.GetId().String())
					} else {
						tfStateAppRoles.Id = types.StringNull()
					}
					if v.GetIsEnabled() != nil {
						tfStateAppRoles.IsEnabled = types.BoolValue(*v.GetIsEnabled())
					} else {
						tfStateAppRoles.IsEnabled = types.BoolNull()
					}
					if v.GetOrigin() != nil {
						tfStateAppRoles.Origin = types.StringValue(*v.GetOrigin())
					} else {
						tfStateAppRoles.Origin = types.StringNull()
					}
					if v.GetValue() != nil {
						tfStateAppRoles.Value = types.StringValue(*v.GetValue())
					} else {
						tfStateAppRoles.Value = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateAppRoles.AttributeTypes(), tfStateAppRoles)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.AppRoles, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetApplicationTemplateId() != nil {
				tfStateValue.ApplicationTemplateId = types.StringValue(*v.GetApplicationTemplateId())
			} else {
				tfStateValue.ApplicationTemplateId = types.StringNull()
			}
			if v.GetCertification() != nil {
				tfStateCertification := applicationsCertificationModel{}

				if v.GetCertification().GetCertificationDetailsUrl() != nil {
					tfStateCertification.CertificationDetailsUrl = types.StringValue(*v.GetCertification().GetCertificationDetailsUrl())
				} else {
					tfStateCertification.CertificationDetailsUrl = types.StringNull()
				}
				if v.GetCertification().GetCertificationExpirationDateTime() != nil {
					tfStateCertification.CertificationExpirationDateTime = types.StringValue(v.GetCertification().GetCertificationExpirationDateTime().String())
				} else {
					tfStateCertification.CertificationExpirationDateTime = types.StringNull()
				}
				if v.GetCertification().GetIsCertifiedByMicrosoft() != nil {
					tfStateCertification.IsCertifiedByMicrosoft = types.BoolValue(*v.GetCertification().GetIsCertifiedByMicrosoft())
				} else {
					tfStateCertification.IsCertifiedByMicrosoft = types.BoolNull()
				}
				if v.GetCertification().GetIsPublisherAttested() != nil {
					tfStateCertification.IsPublisherAttested = types.BoolValue(*v.GetCertification().GetIsPublisherAttested())
				} else {
					tfStateCertification.IsPublisherAttested = types.BoolNull()
				}
				if v.GetCertification().GetLastCertificationDateTime() != nil {
					tfStateCertification.LastCertificationDateTime = types.StringValue(v.GetCertification().GetLastCertificationDateTime().String())
				} else {
					tfStateCertification.LastCertificationDateTime = types.StringNull()
				}

				tfStateValue.Certification, _ = types.ObjectValueFrom(ctx, tfStateCertification.AttributeTypes(), tfStateCertification)
			}
			if v.GetCreatedDateTime() != nil {
				tfStateValue.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			} else {
				tfStateValue.CreatedDateTime = types.StringNull()
			}
			if v.GetDefaultRedirectUri() != nil {
				tfStateValue.DefaultRedirectUri = types.StringValue(*v.GetDefaultRedirectUri())
			} else {
				tfStateValue.DefaultRedirectUri = types.StringNull()
			}
			if v.GetDescription() != nil {
				tfStateValue.Description = types.StringValue(*v.GetDescription())
			} else {
				tfStateValue.Description = types.StringNull()
			}
			if v.GetDisabledByMicrosoftStatus() != nil {
				tfStateValue.DisabledByMicrosoftStatus = types.StringValue(*v.GetDisabledByMicrosoftStatus())
			} else {
				tfStateValue.DisabledByMicrosoftStatus = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				tfStateValue.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				tfStateValue.DisplayName = types.StringNull()
			}
			if v.GetGroupMembershipClaims() != nil {
				tfStateValue.GroupMembershipClaims = types.StringValue(*v.GetGroupMembershipClaims())
			} else {
				tfStateValue.GroupMembershipClaims = types.StringNull()
			}
			if len(v.GetIdentifierUris()) > 0 {
				var valueArrayIdentifierUris []attr.Value
				for _, v := range v.GetIdentifierUris() {
					valueArrayIdentifierUris = append(valueArrayIdentifierUris, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayIdentifierUris)
				tfStateValue.IdentifierUris = listValue
			} else {
				tfStateValue.IdentifierUris = types.ListNull(types.StringType)
			}
			if v.GetInfo() != nil {
				tfStateInfo := applicationsInformationalUrlModel{}

				if v.GetInfo().GetLogoUrl() != nil {
					tfStateInfo.LogoUrl = types.StringValue(*v.GetInfo().GetLogoUrl())
				} else {
					tfStateInfo.LogoUrl = types.StringNull()
				}
				if v.GetInfo().GetMarketingUrl() != nil {
					tfStateInfo.MarketingUrl = types.StringValue(*v.GetInfo().GetMarketingUrl())
				} else {
					tfStateInfo.MarketingUrl = types.StringNull()
				}
				if v.GetInfo().GetPrivacyStatementUrl() != nil {
					tfStateInfo.PrivacyStatementUrl = types.StringValue(*v.GetInfo().GetPrivacyStatementUrl())
				} else {
					tfStateInfo.PrivacyStatementUrl = types.StringNull()
				}
				if v.GetInfo().GetSupportUrl() != nil {
					tfStateInfo.SupportUrl = types.StringValue(*v.GetInfo().GetSupportUrl())
				} else {
					tfStateInfo.SupportUrl = types.StringNull()
				}
				if v.GetInfo().GetTermsOfServiceUrl() != nil {
					tfStateInfo.TermsOfServiceUrl = types.StringValue(*v.GetInfo().GetTermsOfServiceUrl())
				} else {
					tfStateInfo.TermsOfServiceUrl = types.StringNull()
				}

				tfStateValue.Info, _ = types.ObjectValueFrom(ctx, tfStateInfo.AttributeTypes(), tfStateInfo)
			}
			if v.GetIsDeviceOnlyAuthSupported() != nil {
				tfStateValue.IsDeviceOnlyAuthSupported = types.BoolValue(*v.GetIsDeviceOnlyAuthSupported())
			} else {
				tfStateValue.IsDeviceOnlyAuthSupported = types.BoolNull()
			}
			if v.GetIsFallbackPublicClient() != nil {
				tfStateValue.IsFallbackPublicClient = types.BoolValue(*v.GetIsFallbackPublicClient())
			} else {
				tfStateValue.IsFallbackPublicClient = types.BoolNull()
			}
			if len(v.GetKeyCredentials()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetKeyCredentials() {
					tfStateKeyCredentials := applicationsKeyCredentialModel{}

					if v.GetCustomKeyIdentifier() != nil {
						tfStateKeyCredentials.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
					} else {
						tfStateKeyCredentials.CustomKeyIdentifier = types.StringNull()
					}
					if v.GetDisplayName() != nil {
						tfStateKeyCredentials.DisplayName = types.StringValue(*v.GetDisplayName())
					} else {
						tfStateKeyCredentials.DisplayName = types.StringNull()
					}
					if v.GetEndDateTime() != nil {
						tfStateKeyCredentials.EndDateTime = types.StringValue(v.GetEndDateTime().String())
					} else {
						tfStateKeyCredentials.EndDateTime = types.StringNull()
					}
					if v.GetKey() != nil {
						tfStateKeyCredentials.Key = types.StringValue(string(v.GetKey()[:]))
					} else {
						tfStateKeyCredentials.Key = types.StringNull()
					}
					if v.GetKeyId() != nil {
						tfStateKeyCredentials.KeyId = types.StringValue(v.GetKeyId().String())
					} else {
						tfStateKeyCredentials.KeyId = types.StringNull()
					}
					if v.GetStartDateTime() != nil {
						tfStateKeyCredentials.StartDateTime = types.StringValue(v.GetStartDateTime().String())
					} else {
						tfStateKeyCredentials.StartDateTime = types.StringNull()
					}
					if v.GetTypeEscaped() != nil {
						tfStateKeyCredentials.Type = types.StringValue(*v.GetTypeEscaped())
					} else {
						tfStateKeyCredentials.Type = types.StringNull()
					}
					if v.GetUsage() != nil {
						tfStateKeyCredentials.Usage = types.StringValue(*v.GetUsage())
					} else {
						tfStateKeyCredentials.Usage = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateKeyCredentials.AttributeTypes(), tfStateKeyCredentials)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.KeyCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetLogo() != nil {
				tfStateValue.Logo = types.StringValue(string(v.GetLogo()[:]))
			} else {
				tfStateValue.Logo = types.StringNull()
			}
			if v.GetNativeAuthenticationApisEnabled() != nil {
				tfStateValue.NativeAuthenticationApisEnabled = types.StringValue(v.GetNativeAuthenticationApisEnabled().String())
			} else {
				tfStateValue.NativeAuthenticationApisEnabled = types.StringNull()
			}
			if v.GetNotes() != nil {
				tfStateValue.Notes = types.StringValue(*v.GetNotes())
			} else {
				tfStateValue.Notes = types.StringNull()
			}
			if v.GetOauth2RequirePostResponse() != nil {
				tfStateValue.Oauth2RequirePostResponse = types.BoolValue(*v.GetOauth2RequirePostResponse())
			} else {
				tfStateValue.Oauth2RequirePostResponse = types.BoolNull()
			}
			if v.GetOptionalClaims() != nil {
				tfStateOptionalClaims := applicationsOptionalClaimsModel{}

				if len(v.GetOptionalClaims().GetAccessToken()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, v := range v.GetOptionalClaims().GetAccessToken() {
						tfStateAccessToken := applicationsOptionalClaimModel{}

						if len(v.GetAdditionalProperties()) > 0 {
							var valueArrayAdditionalProperties []attr.Value
							for _, v := range v.GetAdditionalProperties() {
								valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(v))
							}
							listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
							tfStateAccessToken.AdditionalProperties = listValue
						} else {
							tfStateAccessToken.AdditionalProperties = types.ListNull(types.StringType)
						}
						if v.GetEssential() != nil {
							tfStateAccessToken.Essential = types.BoolValue(*v.GetEssential())
						} else {
							tfStateAccessToken.Essential = types.BoolNull()
						}
						if v.GetName() != nil {
							tfStateAccessToken.Name = types.StringValue(*v.GetName())
						} else {
							tfStateAccessToken.Name = types.StringNull()
						}
						if v.GetSource() != nil {
							tfStateAccessToken.Source = types.StringValue(*v.GetSource())
						} else {
							tfStateAccessToken.Source = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStateAccessToken.AttributeTypes(), tfStateAccessToken)
						objectValues = append(objectValues, objectValue)
					}
					tfStateOptionalClaims.AccessToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}
				if len(v.GetOptionalClaims().GetIdToken()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, v := range v.GetOptionalClaims().GetIdToken() {
						tfStateIdToken := applicationsOptionalClaimModel{}

						if len(v.GetAdditionalProperties()) > 0 {
							var valueArrayAdditionalProperties []attr.Value
							for _, v := range v.GetAdditionalProperties() {
								valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(v))
							}
							listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
							tfStateIdToken.AdditionalProperties = listValue
						} else {
							tfStateIdToken.AdditionalProperties = types.ListNull(types.StringType)
						}
						if v.GetEssential() != nil {
							tfStateIdToken.Essential = types.BoolValue(*v.GetEssential())
						} else {
							tfStateIdToken.Essential = types.BoolNull()
						}
						if v.GetName() != nil {
							tfStateIdToken.Name = types.StringValue(*v.GetName())
						} else {
							tfStateIdToken.Name = types.StringNull()
						}
						if v.GetSource() != nil {
							tfStateIdToken.Source = types.StringValue(*v.GetSource())
						} else {
							tfStateIdToken.Source = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStateIdToken.AttributeTypes(), tfStateIdToken)
						objectValues = append(objectValues, objectValue)
					}
					tfStateOptionalClaims.IdToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}
				if len(v.GetOptionalClaims().GetSaml2Token()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, v := range v.GetOptionalClaims().GetSaml2Token() {
						tfStateSaml2Token := applicationsOptionalClaimModel{}

						if len(v.GetAdditionalProperties()) > 0 {
							var valueArrayAdditionalProperties []attr.Value
							for _, v := range v.GetAdditionalProperties() {
								valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(v))
							}
							listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
							tfStateSaml2Token.AdditionalProperties = listValue
						} else {
							tfStateSaml2Token.AdditionalProperties = types.ListNull(types.StringType)
						}
						if v.GetEssential() != nil {
							tfStateSaml2Token.Essential = types.BoolValue(*v.GetEssential())
						} else {
							tfStateSaml2Token.Essential = types.BoolNull()
						}
						if v.GetName() != nil {
							tfStateSaml2Token.Name = types.StringValue(*v.GetName())
						} else {
							tfStateSaml2Token.Name = types.StringNull()
						}
						if v.GetSource() != nil {
							tfStateSaml2Token.Source = types.StringValue(*v.GetSource())
						} else {
							tfStateSaml2Token.Source = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStateSaml2Token.AttributeTypes(), tfStateSaml2Token)
						objectValues = append(objectValues, objectValue)
					}
					tfStateOptionalClaims.Saml2Token, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}

				tfStateValue.OptionalClaims, _ = types.ObjectValueFrom(ctx, tfStateOptionalClaims.AttributeTypes(), tfStateOptionalClaims)
			}
			if v.GetParentalControlSettings() != nil {
				tfStateParentalControlSettings := applicationsParentalControlSettingsModel{}

				if len(v.GetParentalControlSettings().GetCountriesBlockedForMinors()) > 0 {
					var valueArrayCountriesBlockedForMinors []attr.Value
					for _, v := range v.GetParentalControlSettings().GetCountriesBlockedForMinors() {
						valueArrayCountriesBlockedForMinors = append(valueArrayCountriesBlockedForMinors, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayCountriesBlockedForMinors)
					tfStateParentalControlSettings.CountriesBlockedForMinors = listValue
				} else {
					tfStateParentalControlSettings.CountriesBlockedForMinors = types.ListNull(types.StringType)
				}
				if v.GetParentalControlSettings().GetLegalAgeGroupRule() != nil {
					tfStateParentalControlSettings.LegalAgeGroupRule = types.StringValue(*v.GetParentalControlSettings().GetLegalAgeGroupRule())
				} else {
					tfStateParentalControlSettings.LegalAgeGroupRule = types.StringNull()
				}

				tfStateValue.ParentalControlSettings, _ = types.ObjectValueFrom(ctx, tfStateParentalControlSettings.AttributeTypes(), tfStateParentalControlSettings)
			}
			if len(v.GetPasswordCredentials()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetPasswordCredentials() {
					tfStatePasswordCredentials := applicationsPasswordCredentialModel{}

					if v.GetCustomKeyIdentifier() != nil {
						tfStatePasswordCredentials.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
					} else {
						tfStatePasswordCredentials.CustomKeyIdentifier = types.StringNull()
					}
					if v.GetDisplayName() != nil {
						tfStatePasswordCredentials.DisplayName = types.StringValue(*v.GetDisplayName())
					} else {
						tfStatePasswordCredentials.DisplayName = types.StringNull()
					}
					if v.GetEndDateTime() != nil {
						tfStatePasswordCredentials.EndDateTime = types.StringValue(v.GetEndDateTime().String())
					} else {
						tfStatePasswordCredentials.EndDateTime = types.StringNull()
					}
					if v.GetHint() != nil {
						tfStatePasswordCredentials.Hint = types.StringValue(*v.GetHint())
					} else {
						tfStatePasswordCredentials.Hint = types.StringNull()
					}
					if v.GetKeyId() != nil {
						tfStatePasswordCredentials.KeyId = types.StringValue(v.GetKeyId().String())
					} else {
						tfStatePasswordCredentials.KeyId = types.StringNull()
					}
					if v.GetSecretText() != nil {
						tfStatePasswordCredentials.SecretText = types.StringValue(*v.GetSecretText())
					} else {
						tfStatePasswordCredentials.SecretText = types.StringNull()
					}
					if v.GetStartDateTime() != nil {
						tfStatePasswordCredentials.StartDateTime = types.StringValue(v.GetStartDateTime().String())
					} else {
						tfStatePasswordCredentials.StartDateTime = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStatePasswordCredentials.AttributeTypes(), tfStatePasswordCredentials)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.PasswordCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetPublicClient() != nil {
				tfStatePublicClient := applicationsPublicClientApplicationModel{}

				if len(v.GetPublicClient().GetRedirectUris()) > 0 {
					var valueArrayRedirectUris []attr.Value
					for _, v := range v.GetPublicClient().GetRedirectUris() {
						valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
					tfStatePublicClient.RedirectUris = listValue
				} else {
					tfStatePublicClient.RedirectUris = types.ListNull(types.StringType)
				}

				tfStateValue.PublicClient, _ = types.ObjectValueFrom(ctx, tfStatePublicClient.AttributeTypes(), tfStatePublicClient)
			}
			if v.GetPublisherDomain() != nil {
				tfStateValue.PublisherDomain = types.StringValue(*v.GetPublisherDomain())
			} else {
				tfStateValue.PublisherDomain = types.StringNull()
			}
			if v.GetRequestSignatureVerification() != nil {
				tfStateRequestSignatureVerification := applicationsRequestSignatureVerificationModel{}

				if v.GetRequestSignatureVerification().GetAllowedWeakAlgorithms() != nil {
					tfStateRequestSignatureVerification.AllowedWeakAlgorithms = types.StringValue(v.GetRequestSignatureVerification().GetAllowedWeakAlgorithms().String())
				} else {
					tfStateRequestSignatureVerification.AllowedWeakAlgorithms = types.StringNull()
				}
				if v.GetRequestSignatureVerification().GetIsSignedRequestRequired() != nil {
					tfStateRequestSignatureVerification.IsSignedRequestRequired = types.BoolValue(*v.GetRequestSignatureVerification().GetIsSignedRequestRequired())
				} else {
					tfStateRequestSignatureVerification.IsSignedRequestRequired = types.BoolNull()
				}

				tfStateValue.RequestSignatureVerification, _ = types.ObjectValueFrom(ctx, tfStateRequestSignatureVerification.AttributeTypes(), tfStateRequestSignatureVerification)
			}
			if len(v.GetRequiredResourceAccess()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetRequiredResourceAccess() {
					tfStateRequiredResourceAccess := applicationsRequiredResourceAccessModel{}

					if len(v.GetResourceAccess()) > 0 {
						objectValues := []basetypes.ObjectValue{}
						for _, v := range v.GetResourceAccess() {
							tfStateResourceAccess := applicationsResourceAccessModel{}

							if v.GetId() != nil {
								tfStateResourceAccess.Id = types.StringValue(v.GetId().String())
							} else {
								tfStateResourceAccess.Id = types.StringNull()
							}
							if v.GetTypeEscaped() != nil {
								tfStateResourceAccess.Type = types.StringValue(*v.GetTypeEscaped())
							} else {
								tfStateResourceAccess.Type = types.StringNull()
							}
							objectValue, _ := types.ObjectValueFrom(ctx, tfStateResourceAccess.AttributeTypes(), tfStateResourceAccess)
							objectValues = append(objectValues, objectValue)
						}
						tfStateRequiredResourceAccess.ResourceAccess, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
					}
					if v.GetResourceAppId() != nil {
						tfStateRequiredResourceAccess.ResourceAppId = types.StringValue(*v.GetResourceAppId())
					} else {
						tfStateRequiredResourceAccess.ResourceAppId = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateRequiredResourceAccess.AttributeTypes(), tfStateRequiredResourceAccess)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.RequiredResourceAccess, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetSamlMetadataUrl() != nil {
				tfStateValue.SamlMetadataUrl = types.StringValue(*v.GetSamlMetadataUrl())
			} else {
				tfStateValue.SamlMetadataUrl = types.StringNull()
			}
			if v.GetServiceManagementReference() != nil {
				tfStateValue.ServiceManagementReference = types.StringValue(*v.GetServiceManagementReference())
			} else {
				tfStateValue.ServiceManagementReference = types.StringNull()
			}
			if v.GetServicePrincipalLockConfiguration() != nil {
				tfStateServicePrincipalLockConfiguration := applicationsServicePrincipalLockConfigurationModel{}

				if v.GetServicePrincipalLockConfiguration().GetAllProperties() != nil {
					tfStateServicePrincipalLockConfiguration.AllProperties = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetAllProperties())
				} else {
					tfStateServicePrincipalLockConfiguration.AllProperties = types.BoolNull()
				}
				if v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign() != nil {
					tfStateServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign())
				} else {
					tfStateServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolNull()
				}
				if v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify() != nil {
					tfStateServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify())
				} else {
					tfStateServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolNull()
				}
				if v.GetServicePrincipalLockConfiguration().GetIsEnabled() != nil {
					tfStateServicePrincipalLockConfiguration.IsEnabled = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetIsEnabled())
				} else {
					tfStateServicePrincipalLockConfiguration.IsEnabled = types.BoolNull()
				}
				if v.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId() != nil {
					tfStateServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolValue(*v.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId())
				} else {
					tfStateServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolNull()
				}

				tfStateValue.ServicePrincipalLockConfiguration, _ = types.ObjectValueFrom(ctx, tfStateServicePrincipalLockConfiguration.AttributeTypes(), tfStateServicePrincipalLockConfiguration)
			}
			if v.GetSignInAudience() != nil {
				tfStateValue.SignInAudience = types.StringValue(*v.GetSignInAudience())
			} else {
				tfStateValue.SignInAudience = types.StringNull()
			}
			if v.GetSpa() != nil {
				tfStateSpa := applicationsSpaApplicationModel{}

				if len(v.GetSpa().GetRedirectUris()) > 0 {
					var valueArrayRedirectUris []attr.Value
					for _, v := range v.GetSpa().GetRedirectUris() {
						valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
					tfStateSpa.RedirectUris = listValue
				} else {
					tfStateSpa.RedirectUris = types.ListNull(types.StringType)
				}

				tfStateValue.Spa, _ = types.ObjectValueFrom(ctx, tfStateSpa.AttributeTypes(), tfStateSpa)
			}
			if len(v.GetTags()) > 0 {
				var valueArrayTags []attr.Value
				for _, v := range v.GetTags() {
					valueArrayTags = append(valueArrayTags, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayTags)
				tfStateValue.Tags = listValue
			} else {
				tfStateValue.Tags = types.ListNull(types.StringType)
			}
			if v.GetTokenEncryptionKeyId() != nil {
				tfStateValue.TokenEncryptionKeyId = types.StringValue(v.GetTokenEncryptionKeyId().String())
			} else {
				tfStateValue.TokenEncryptionKeyId = types.StringNull()
			}
			if v.GetUniqueName() != nil {
				tfStateValue.UniqueName = types.StringValue(*v.GetUniqueName())
			} else {
				tfStateValue.UniqueName = types.StringNull()
			}
			if v.GetVerifiedPublisher() != nil {
				tfStateVerifiedPublisher := applicationsVerifiedPublisherModel{}

				if v.GetVerifiedPublisher().GetAddedDateTime() != nil {
					tfStateVerifiedPublisher.AddedDateTime = types.StringValue(v.GetVerifiedPublisher().GetAddedDateTime().String())
				} else {
					tfStateVerifiedPublisher.AddedDateTime = types.StringNull()
				}
				if v.GetVerifiedPublisher().GetDisplayName() != nil {
					tfStateVerifiedPublisher.DisplayName = types.StringValue(*v.GetVerifiedPublisher().GetDisplayName())
				} else {
					tfStateVerifiedPublisher.DisplayName = types.StringNull()
				}
				if v.GetVerifiedPublisher().GetVerifiedPublisherId() != nil {
					tfStateVerifiedPublisher.VerifiedPublisherId = types.StringValue(*v.GetVerifiedPublisher().GetVerifiedPublisherId())
				} else {
					tfStateVerifiedPublisher.VerifiedPublisherId = types.StringNull()
				}

				tfStateValue.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfStateVerifiedPublisher.AttributeTypes(), tfStateVerifiedPublisher)
			}
			if v.GetWeb() != nil {
				tfStateWeb := applicationsWebApplicationModel{}

				if v.GetWeb().GetHomePageUrl() != nil {
					tfStateWeb.HomePageUrl = types.StringValue(*v.GetWeb().GetHomePageUrl())
				} else {
					tfStateWeb.HomePageUrl = types.StringNull()
				}
				if v.GetWeb().GetImplicitGrantSettings() != nil {
					tfStateImplicitGrantSettings := applicationsImplicitGrantSettingsModel{}

					if v.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance() != nil {
						tfStateImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolValue(*v.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance())
					} else {
						tfStateImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolNull()
					}
					if v.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance() != nil {
						tfStateImplicitGrantSettings.EnableIdTokenIssuance = types.BoolValue(*v.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance())
					} else {
						tfStateImplicitGrantSettings.EnableIdTokenIssuance = types.BoolNull()
					}

					tfStateWeb.ImplicitGrantSettings, _ = types.ObjectValueFrom(ctx, tfStateImplicitGrantSettings.AttributeTypes(), tfStateImplicitGrantSettings)
				}
				if v.GetWeb().GetLogoutUrl() != nil {
					tfStateWeb.LogoutUrl = types.StringValue(*v.GetWeb().GetLogoutUrl())
				} else {
					tfStateWeb.LogoutUrl = types.StringNull()
				}
				if len(v.GetWeb().GetRedirectUriSettings()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, v := range v.GetWeb().GetRedirectUriSettings() {
						tfStateRedirectUriSettings := applicationsRedirectUriSettingsModel{}

						if v.GetUri() != nil {
							tfStateRedirectUriSettings.Uri = types.StringValue(*v.GetUri())
						} else {
							tfStateRedirectUriSettings.Uri = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStateRedirectUriSettings.AttributeTypes(), tfStateRedirectUriSettings)
						objectValues = append(objectValues, objectValue)
					}
					tfStateWeb.RedirectUriSettings, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}
				if len(v.GetWeb().GetRedirectUris()) > 0 {
					var valueArrayRedirectUris []attr.Value
					for _, v := range v.GetWeb().GetRedirectUris() {
						valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
					tfStateWeb.RedirectUris = listValue
				} else {
					tfStateWeb.RedirectUris = types.ListNull(types.StringType)
				}

				tfStateValue.Web, _ = types.ObjectValueFrom(ctx, tfStateWeb.AttributeTypes(), tfStateWeb)
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateValue.AttributeTypes(), tfStateValue)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplications.Value, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateApplications)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
