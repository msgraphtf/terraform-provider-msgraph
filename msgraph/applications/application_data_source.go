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
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &applicationDataSource{}
	_ datasource.DataSourceWithConfigure = &applicationDataSource{}
)

// NewApplicationDataSource is a helper function to simplify the provider implementation.
func NewApplicationDataSource() datasource.DataSource {
	return &applicationDataSource{}
}

// applicationDataSource is the data source implementation.
type applicationDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *applicationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Configure adds the provider configured client to the data source.
func (d *applicationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *applicationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
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
							Optional:    true,
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
									Optional:    true,
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
							Optional:    true,
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
										Optional:    true,
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
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *applicationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfStateApplication applicationModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateApplication)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := applications.ApplicationItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &applications.ApplicationItemRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"deletedDateTime",
				"addIns",
				"api",
				"appId",
				"appRoles",
				"applicationTemplateId",
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
				"nativeAuthenticationApisEnabled",
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
				"uniqueName",
				"verifiedPublisher",
				"web",
			},
		},
	}

	var response models.Applicationable
	var err error

	if !tfStateApplication.Id.IsNull() {
		response, err = d.client.Applications().ByApplicationId(tfStateApplication.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting application",
			err.Error(),
		)
		return
	}

	if response.GetId() != nil {
		tfStateApplication.Id = types.StringValue(*response.GetId())
	} else {
		tfStateApplication.Id = types.StringNull()
	}
	if response.GetDeletedDateTime() != nil {
		tfStateApplication.DeletedDateTime = types.StringValue(response.GetDeletedDateTime().String())
	} else {
		tfStateApplication.DeletedDateTime = types.StringNull()
	}
	if len(response.GetAddIns()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseAddIns := range response.GetAddIns() {
			tfStateAddIn := applicationAddInModel{}

			if responseAddIns.GetId() != nil {
				tfStateAddIn.Id = types.StringValue(responseAddIns.GetId().String())
			} else {
				tfStateAddIn.Id = types.StringNull()
			}
			if len(responseAddIns.GetProperties()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseProperties := range responseAddIns.GetProperties() {
					tfStateKeyValue := applicationKeyValueModel{}

					if responseProperties.GetKey() != nil {
						tfStateKeyValue.Key = types.StringValue(*responseProperties.GetKey())
					} else {
						tfStateKeyValue.Key = types.StringNull()
					}
					if responseProperties.GetValue() != nil {
						tfStateKeyValue.Value = types.StringValue(*responseProperties.GetValue())
					} else {
						tfStateKeyValue.Value = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateKeyValue.AttributeTypes(), tfStateKeyValue)
					objectValues = append(objectValues, objectValue)
				}
				tfStateAddIn.Properties, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if responseAddIns.GetTypeEscaped() != nil {
				tfStateAddIn.Type = types.StringValue(*responseAddIns.GetTypeEscaped())
			} else {
				tfStateAddIn.Type = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAddIn.AttributeTypes(), tfStateAddIn)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.AddIns, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if response.GetApi() != nil {
		tfStateApiApplication := applicationApiApplicationModel{}

		if response.GetApi().GetAcceptMappedClaims() != nil {
			tfStateApiApplication.AcceptMappedClaims = types.BoolValue(*response.GetApi().GetAcceptMappedClaims())
		} else {
			tfStateApiApplication.AcceptMappedClaims = types.BoolNull()
		}
		if len(response.GetApi().GetKnownClientApplications()) > 0 {
			var valueArrayKnownClientApplications []attr.Value
			for _, responseKnownClientApplications := range response.GetApi().GetKnownClientApplications() {
				valueArrayKnownClientApplications = append(valueArrayKnownClientApplications, types.StringValue(responseKnownClientApplications.String()))
			}
			tfStateApiApplication.KnownClientApplications, _ = types.ListValue(types.StringType, valueArrayKnownClientApplications)
		} else {
			tfStateApiApplication.KnownClientApplications = types.ListNull(types.StringType)
		}
		if len(response.GetApi().GetOauth2PermissionScopes()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, responseOauth2PermissionScopes := range response.GetApi().GetOauth2PermissionScopes() {
				tfStatePermissionScope := applicationPermissionScopeModel{}

				if responseOauth2PermissionScopes.GetAdminConsentDescription() != nil {
					tfStatePermissionScope.AdminConsentDescription = types.StringValue(*responseOauth2PermissionScopes.GetAdminConsentDescription())
				} else {
					tfStatePermissionScope.AdminConsentDescription = types.StringNull()
				}
				if responseOauth2PermissionScopes.GetAdminConsentDisplayName() != nil {
					tfStatePermissionScope.AdminConsentDisplayName = types.StringValue(*responseOauth2PermissionScopes.GetAdminConsentDisplayName())
				} else {
					tfStatePermissionScope.AdminConsentDisplayName = types.StringNull()
				}
				if responseOauth2PermissionScopes.GetId() != nil {
					tfStatePermissionScope.Id = types.StringValue(responseOauth2PermissionScopes.GetId().String())
				} else {
					tfStatePermissionScope.Id = types.StringNull()
				}
				if responseOauth2PermissionScopes.GetIsEnabled() != nil {
					tfStatePermissionScope.IsEnabled = types.BoolValue(*responseOauth2PermissionScopes.GetIsEnabled())
				} else {
					tfStatePermissionScope.IsEnabled = types.BoolNull()
				}
				if responseOauth2PermissionScopes.GetOrigin() != nil {
					tfStatePermissionScope.Origin = types.StringValue(*responseOauth2PermissionScopes.GetOrigin())
				} else {
					tfStatePermissionScope.Origin = types.StringNull()
				}
				if responseOauth2PermissionScopes.GetTypeEscaped() != nil {
					tfStatePermissionScope.Type = types.StringValue(*responseOauth2PermissionScopes.GetTypeEscaped())
				} else {
					tfStatePermissionScope.Type = types.StringNull()
				}
				if responseOauth2PermissionScopes.GetUserConsentDescription() != nil {
					tfStatePermissionScope.UserConsentDescription = types.StringValue(*responseOauth2PermissionScopes.GetUserConsentDescription())
				} else {
					tfStatePermissionScope.UserConsentDescription = types.StringNull()
				}
				if responseOauth2PermissionScopes.GetUserConsentDisplayName() != nil {
					tfStatePermissionScope.UserConsentDisplayName = types.StringValue(*responseOauth2PermissionScopes.GetUserConsentDisplayName())
				} else {
					tfStatePermissionScope.UserConsentDisplayName = types.StringNull()
				}
				if responseOauth2PermissionScopes.GetValue() != nil {
					tfStatePermissionScope.Value = types.StringValue(*responseOauth2PermissionScopes.GetValue())
				} else {
					tfStatePermissionScope.Value = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStatePermissionScope.AttributeTypes(), tfStatePermissionScope)
				objectValues = append(objectValues, objectValue)
			}
			tfStateApiApplication.Oauth2PermissionScopes, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(response.GetApi().GetPreAuthorizedApplications()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, responsePreAuthorizedApplications := range response.GetApi().GetPreAuthorizedApplications() {
				tfStatePreAuthorizedApplication := applicationPreAuthorizedApplicationModel{}

				if responsePreAuthorizedApplications.GetAppId() != nil {
					tfStatePreAuthorizedApplication.AppId = types.StringValue(*responsePreAuthorizedApplications.GetAppId())
				} else {
					tfStatePreAuthorizedApplication.AppId = types.StringNull()
				}
				if len(responsePreAuthorizedApplications.GetDelegatedPermissionIds()) > 0 {
					var valueArrayDelegatedPermissionIds []attr.Value
					for _, responseDelegatedPermissionIds := range responsePreAuthorizedApplications.GetDelegatedPermissionIds() {
						valueArrayDelegatedPermissionIds = append(valueArrayDelegatedPermissionIds, types.StringValue(responseDelegatedPermissionIds))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayDelegatedPermissionIds)
					tfStatePreAuthorizedApplication.DelegatedPermissionIds = listValue
				} else {
					tfStatePreAuthorizedApplication.DelegatedPermissionIds = types.ListNull(types.StringType)
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStatePreAuthorizedApplication.AttributeTypes(), tfStatePreAuthorizedApplication)
				objectValues = append(objectValues, objectValue)
			}
			tfStateApiApplication.PreAuthorizedApplications, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}

		tfStateApplication.Api, _ = types.ObjectValueFrom(ctx, tfStateApiApplication.AttributeTypes(), tfStateApiApplication)
	}
	if response.GetAppId() != nil {
		tfStateApplication.AppId = types.StringValue(*response.GetAppId())
	} else {
		tfStateApplication.AppId = types.StringNull()
	}
	if len(response.GetAppRoles()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseAppRoles := range response.GetAppRoles() {
			tfStateAppRole := applicationAppRoleModel{}

			if len(responseAppRoles.GetAllowedMemberTypes()) > 0 {
				var valueArrayAllowedMemberTypes []attr.Value
				for _, responseAllowedMemberTypes := range responseAppRoles.GetAllowedMemberTypes() {
					valueArrayAllowedMemberTypes = append(valueArrayAllowedMemberTypes, types.StringValue(responseAllowedMemberTypes))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayAllowedMemberTypes)
				tfStateAppRole.AllowedMemberTypes = listValue
			} else {
				tfStateAppRole.AllowedMemberTypes = types.ListNull(types.StringType)
			}
			if responseAppRoles.GetDescription() != nil {
				tfStateAppRole.Description = types.StringValue(*responseAppRoles.GetDescription())
			} else {
				tfStateAppRole.Description = types.StringNull()
			}
			if responseAppRoles.GetDisplayName() != nil {
				tfStateAppRole.DisplayName = types.StringValue(*responseAppRoles.GetDisplayName())
			} else {
				tfStateAppRole.DisplayName = types.StringNull()
			}
			if responseAppRoles.GetId() != nil {
				tfStateAppRole.Id = types.StringValue(responseAppRoles.GetId().String())
			} else {
				tfStateAppRole.Id = types.StringNull()
			}
			if responseAppRoles.GetIsEnabled() != nil {
				tfStateAppRole.IsEnabled = types.BoolValue(*responseAppRoles.GetIsEnabled())
			} else {
				tfStateAppRole.IsEnabled = types.BoolNull()
			}
			if responseAppRoles.GetOrigin() != nil {
				tfStateAppRole.Origin = types.StringValue(*responseAppRoles.GetOrigin())
			} else {
				tfStateAppRole.Origin = types.StringNull()
			}
			if responseAppRoles.GetValue() != nil {
				tfStateAppRole.Value = types.StringValue(*responseAppRoles.GetValue())
			} else {
				tfStateAppRole.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAppRole.AttributeTypes(), tfStateAppRole)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.AppRoles, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if response.GetApplicationTemplateId() != nil {
		tfStateApplication.ApplicationTemplateId = types.StringValue(*response.GetApplicationTemplateId())
	} else {
		tfStateApplication.ApplicationTemplateId = types.StringNull()
	}
	if response.GetCertification() != nil {
		tfStateCertification := applicationCertificationModel{}

		if response.GetCertification().GetCertificationDetailsUrl() != nil {
			tfStateCertification.CertificationDetailsUrl = types.StringValue(*response.GetCertification().GetCertificationDetailsUrl())
		} else {
			tfStateCertification.CertificationDetailsUrl = types.StringNull()
		}
		if response.GetCertification().GetCertificationExpirationDateTime() != nil {
			tfStateCertification.CertificationExpirationDateTime = types.StringValue(response.GetCertification().GetCertificationExpirationDateTime().String())
		} else {
			tfStateCertification.CertificationExpirationDateTime = types.StringNull()
		}
		if response.GetCertification().GetIsCertifiedByMicrosoft() != nil {
			tfStateCertification.IsCertifiedByMicrosoft = types.BoolValue(*response.GetCertification().GetIsCertifiedByMicrosoft())
		} else {
			tfStateCertification.IsCertifiedByMicrosoft = types.BoolNull()
		}
		if response.GetCertification().GetIsPublisherAttested() != nil {
			tfStateCertification.IsPublisherAttested = types.BoolValue(*response.GetCertification().GetIsPublisherAttested())
		} else {
			tfStateCertification.IsPublisherAttested = types.BoolNull()
		}
		if response.GetCertification().GetLastCertificationDateTime() != nil {
			tfStateCertification.LastCertificationDateTime = types.StringValue(response.GetCertification().GetLastCertificationDateTime().String())
		} else {
			tfStateCertification.LastCertificationDateTime = types.StringNull()
		}

		tfStateApplication.Certification, _ = types.ObjectValueFrom(ctx, tfStateCertification.AttributeTypes(), tfStateCertification)
	}
	if response.GetCreatedDateTime() != nil {
		tfStateApplication.CreatedDateTime = types.StringValue(response.GetCreatedDateTime().String())
	} else {
		tfStateApplication.CreatedDateTime = types.StringNull()
	}
	if response.GetDefaultRedirectUri() != nil {
		tfStateApplication.DefaultRedirectUri = types.StringValue(*response.GetDefaultRedirectUri())
	} else {
		tfStateApplication.DefaultRedirectUri = types.StringNull()
	}
	if response.GetDescription() != nil {
		tfStateApplication.Description = types.StringValue(*response.GetDescription())
	} else {
		tfStateApplication.Description = types.StringNull()
	}
	if response.GetDisabledByMicrosoftStatus() != nil {
		tfStateApplication.DisabledByMicrosoftStatus = types.StringValue(*response.GetDisabledByMicrosoftStatus())
	} else {
		tfStateApplication.DisabledByMicrosoftStatus = types.StringNull()
	}
	if response.GetDisplayName() != nil {
		tfStateApplication.DisplayName = types.StringValue(*response.GetDisplayName())
	} else {
		tfStateApplication.DisplayName = types.StringNull()
	}
	if response.GetGroupMembershipClaims() != nil {
		tfStateApplication.GroupMembershipClaims = types.StringValue(*response.GetGroupMembershipClaims())
	} else {
		tfStateApplication.GroupMembershipClaims = types.StringNull()
	}
	if len(response.GetIdentifierUris()) > 0 {
		var valueArrayIdentifierUris []attr.Value
		for _, responseIdentifierUris := range response.GetIdentifierUris() {
			valueArrayIdentifierUris = append(valueArrayIdentifierUris, types.StringValue(responseIdentifierUris))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayIdentifierUris)
		tfStateApplication.IdentifierUris = listValue
	} else {
		tfStateApplication.IdentifierUris = types.ListNull(types.StringType)
	}
	if response.GetInfo() != nil {
		tfStateInformationalUrl := applicationInformationalUrlModel{}

		if response.GetInfo().GetLogoUrl() != nil {
			tfStateInformationalUrl.LogoUrl = types.StringValue(*response.GetInfo().GetLogoUrl())
		} else {
			tfStateInformationalUrl.LogoUrl = types.StringNull()
		}
		if response.GetInfo().GetMarketingUrl() != nil {
			tfStateInformationalUrl.MarketingUrl = types.StringValue(*response.GetInfo().GetMarketingUrl())
		} else {
			tfStateInformationalUrl.MarketingUrl = types.StringNull()
		}
		if response.GetInfo().GetPrivacyStatementUrl() != nil {
			tfStateInformationalUrl.PrivacyStatementUrl = types.StringValue(*response.GetInfo().GetPrivacyStatementUrl())
		} else {
			tfStateInformationalUrl.PrivacyStatementUrl = types.StringNull()
		}
		if response.GetInfo().GetSupportUrl() != nil {
			tfStateInformationalUrl.SupportUrl = types.StringValue(*response.GetInfo().GetSupportUrl())
		} else {
			tfStateInformationalUrl.SupportUrl = types.StringNull()
		}
		if response.GetInfo().GetTermsOfServiceUrl() != nil {
			tfStateInformationalUrl.TermsOfServiceUrl = types.StringValue(*response.GetInfo().GetTermsOfServiceUrl())
		} else {
			tfStateInformationalUrl.TermsOfServiceUrl = types.StringNull()
		}

		tfStateApplication.Info, _ = types.ObjectValueFrom(ctx, tfStateInformationalUrl.AttributeTypes(), tfStateInformationalUrl)
	}
	if response.GetIsDeviceOnlyAuthSupported() != nil {
		tfStateApplication.IsDeviceOnlyAuthSupported = types.BoolValue(*response.GetIsDeviceOnlyAuthSupported())
	} else {
		tfStateApplication.IsDeviceOnlyAuthSupported = types.BoolNull()
	}
	if response.GetIsFallbackPublicClient() != nil {
		tfStateApplication.IsFallbackPublicClient = types.BoolValue(*response.GetIsFallbackPublicClient())
	} else {
		tfStateApplication.IsFallbackPublicClient = types.BoolNull()
	}
	if len(response.GetKeyCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseKeyCredentials := range response.GetKeyCredentials() {
			tfStateKeyCredential := applicationKeyCredentialModel{}

			if responseKeyCredentials.GetCustomKeyIdentifier() != nil {
				tfStateKeyCredential.CustomKeyIdentifier = types.StringValue(string(responseKeyCredentials.GetCustomKeyIdentifier()[:]))
			} else {
				tfStateKeyCredential.CustomKeyIdentifier = types.StringNull()
			}
			if responseKeyCredentials.GetDisplayName() != nil {
				tfStateKeyCredential.DisplayName = types.StringValue(*responseKeyCredentials.GetDisplayName())
			} else {
				tfStateKeyCredential.DisplayName = types.StringNull()
			}
			if responseKeyCredentials.GetEndDateTime() != nil {
				tfStateKeyCredential.EndDateTime = types.StringValue(responseKeyCredentials.GetEndDateTime().String())
			} else {
				tfStateKeyCredential.EndDateTime = types.StringNull()
			}
			if responseKeyCredentials.GetKey() != nil {
				tfStateKeyCredential.Key = types.StringValue(string(responseKeyCredentials.GetKey()[:]))
			} else {
				tfStateKeyCredential.Key = types.StringNull()
			}
			if responseKeyCredentials.GetKeyId() != nil {
				tfStateKeyCredential.KeyId = types.StringValue(responseKeyCredentials.GetKeyId().String())
			} else {
				tfStateKeyCredential.KeyId = types.StringNull()
			}
			if responseKeyCredentials.GetStartDateTime() != nil {
				tfStateKeyCredential.StartDateTime = types.StringValue(responseKeyCredentials.GetStartDateTime().String())
			} else {
				tfStateKeyCredential.StartDateTime = types.StringNull()
			}
			if responseKeyCredentials.GetTypeEscaped() != nil {
				tfStateKeyCredential.Type = types.StringValue(*responseKeyCredentials.GetTypeEscaped())
			} else {
				tfStateKeyCredential.Type = types.StringNull()
			}
			if responseKeyCredentials.GetUsage() != nil {
				tfStateKeyCredential.Usage = types.StringValue(*responseKeyCredentials.GetUsage())
			} else {
				tfStateKeyCredential.Usage = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateKeyCredential.AttributeTypes(), tfStateKeyCredential)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.KeyCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if response.GetLogo() != nil {
		tfStateApplication.Logo = types.StringValue(string(response.GetLogo()[:]))
	} else {
		tfStateApplication.Logo = types.StringNull()
	}
	if response.GetNativeAuthenticationApisEnabled() != nil {
		tfStateApplication.NativeAuthenticationApisEnabled = types.StringValue(response.GetNativeAuthenticationApisEnabled().String())
	} else {
		tfStateApplication.NativeAuthenticationApisEnabled = types.StringNull()
	}
	if response.GetNotes() != nil {
		tfStateApplication.Notes = types.StringValue(*response.GetNotes())
	} else {
		tfStateApplication.Notes = types.StringNull()
	}
	if response.GetOauth2RequirePostResponse() != nil {
		tfStateApplication.Oauth2RequirePostResponse = types.BoolValue(*response.GetOauth2RequirePostResponse())
	} else {
		tfStateApplication.Oauth2RequirePostResponse = types.BoolNull()
	}
	if response.GetOptionalClaims() != nil {
		tfStateOptionalClaims := applicationOptionalClaimsModel{}

		if len(response.GetOptionalClaims().GetAccessToken()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, responseAccessToken := range response.GetOptionalClaims().GetAccessToken() {
				tfStateOptionalClaim := applicationOptionalClaimModel{}

				if len(responseAccessToken.GetAdditionalProperties()) > 0 {
					var valueArrayAdditionalProperties []attr.Value
					for _, responseAdditionalProperties := range responseAccessToken.GetAdditionalProperties() {
						valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(responseAdditionalProperties))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
					tfStateOptionalClaim.AdditionalProperties = listValue
				} else {
					tfStateOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				if responseAccessToken.GetEssential() != nil {
					tfStateOptionalClaim.Essential = types.BoolValue(*responseAccessToken.GetEssential())
				} else {
					tfStateOptionalClaim.Essential = types.BoolNull()
				}
				if responseAccessToken.GetName() != nil {
					tfStateOptionalClaim.Name = types.StringValue(*responseAccessToken.GetName())
				} else {
					tfStateOptionalClaim.Name = types.StringNull()
				}
				if responseAccessToken.GetSource() != nil {
					tfStateOptionalClaim.Source = types.StringValue(*responseAccessToken.GetSource())
				} else {
					tfStateOptionalClaim.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStateOptionalClaim.AttributeTypes(), tfStateOptionalClaim)
				objectValues = append(objectValues, objectValue)
			}
			tfStateOptionalClaims.AccessToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(response.GetOptionalClaims().GetIdToken()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, responseIdToken := range response.GetOptionalClaims().GetIdToken() {
				tfStateOptionalClaim := applicationOptionalClaimModel{}

				if len(responseIdToken.GetAdditionalProperties()) > 0 {
					var valueArrayAdditionalProperties []attr.Value
					for _, responseAdditionalProperties := range responseIdToken.GetAdditionalProperties() {
						valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(responseAdditionalProperties))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
					tfStateOptionalClaim.AdditionalProperties = listValue
				} else {
					tfStateOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				if responseIdToken.GetEssential() != nil {
					tfStateOptionalClaim.Essential = types.BoolValue(*responseIdToken.GetEssential())
				} else {
					tfStateOptionalClaim.Essential = types.BoolNull()
				}
				if responseIdToken.GetName() != nil {
					tfStateOptionalClaim.Name = types.StringValue(*responseIdToken.GetName())
				} else {
					tfStateOptionalClaim.Name = types.StringNull()
				}
				if responseIdToken.GetSource() != nil {
					tfStateOptionalClaim.Source = types.StringValue(*responseIdToken.GetSource())
				} else {
					tfStateOptionalClaim.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStateOptionalClaim.AttributeTypes(), tfStateOptionalClaim)
				objectValues = append(objectValues, objectValue)
			}
			tfStateOptionalClaims.IdToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(response.GetOptionalClaims().GetSaml2Token()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, responseSaml2Token := range response.GetOptionalClaims().GetSaml2Token() {
				tfStateOptionalClaim := applicationOptionalClaimModel{}

				if len(responseSaml2Token.GetAdditionalProperties()) > 0 {
					var valueArrayAdditionalProperties []attr.Value
					for _, responseAdditionalProperties := range responseSaml2Token.GetAdditionalProperties() {
						valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(responseAdditionalProperties))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
					tfStateOptionalClaim.AdditionalProperties = listValue
				} else {
					tfStateOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				if responseSaml2Token.GetEssential() != nil {
					tfStateOptionalClaim.Essential = types.BoolValue(*responseSaml2Token.GetEssential())
				} else {
					tfStateOptionalClaim.Essential = types.BoolNull()
				}
				if responseSaml2Token.GetName() != nil {
					tfStateOptionalClaim.Name = types.StringValue(*responseSaml2Token.GetName())
				} else {
					tfStateOptionalClaim.Name = types.StringNull()
				}
				if responseSaml2Token.GetSource() != nil {
					tfStateOptionalClaim.Source = types.StringValue(*responseSaml2Token.GetSource())
				} else {
					tfStateOptionalClaim.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStateOptionalClaim.AttributeTypes(), tfStateOptionalClaim)
				objectValues = append(objectValues, objectValue)
			}
			tfStateOptionalClaims.Saml2Token, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}

		tfStateApplication.OptionalClaims, _ = types.ObjectValueFrom(ctx, tfStateOptionalClaims.AttributeTypes(), tfStateOptionalClaims)
	}
	if response.GetParentalControlSettings() != nil {
		tfStateParentalControlSettings := applicationParentalControlSettingsModel{}

		if len(response.GetParentalControlSettings().GetCountriesBlockedForMinors()) > 0 {
			var valueArrayCountriesBlockedForMinors []attr.Value
			for _, responseCountriesBlockedForMinors := range response.GetParentalControlSettings().GetCountriesBlockedForMinors() {
				valueArrayCountriesBlockedForMinors = append(valueArrayCountriesBlockedForMinors, types.StringValue(responseCountriesBlockedForMinors))
			}
			listValue, _ := types.ListValue(types.StringType, valueArrayCountriesBlockedForMinors)
			tfStateParentalControlSettings.CountriesBlockedForMinors = listValue
		} else {
			tfStateParentalControlSettings.CountriesBlockedForMinors = types.ListNull(types.StringType)
		}
		if response.GetParentalControlSettings().GetLegalAgeGroupRule() != nil {
			tfStateParentalControlSettings.LegalAgeGroupRule = types.StringValue(*response.GetParentalControlSettings().GetLegalAgeGroupRule())
		} else {
			tfStateParentalControlSettings.LegalAgeGroupRule = types.StringNull()
		}

		tfStateApplication.ParentalControlSettings, _ = types.ObjectValueFrom(ctx, tfStateParentalControlSettings.AttributeTypes(), tfStateParentalControlSettings)
	}
	if len(response.GetPasswordCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responsePasswordCredentials := range response.GetPasswordCredentials() {
			tfStatePasswordCredential := applicationPasswordCredentialModel{}

			if responsePasswordCredentials.GetCustomKeyIdentifier() != nil {
				tfStatePasswordCredential.CustomKeyIdentifier = types.StringValue(string(responsePasswordCredentials.GetCustomKeyIdentifier()[:]))
			} else {
				tfStatePasswordCredential.CustomKeyIdentifier = types.StringNull()
			}
			if responsePasswordCredentials.GetDisplayName() != nil {
				tfStatePasswordCredential.DisplayName = types.StringValue(*responsePasswordCredentials.GetDisplayName())
			} else {
				tfStatePasswordCredential.DisplayName = types.StringNull()
			}
			if responsePasswordCredentials.GetEndDateTime() != nil {
				tfStatePasswordCredential.EndDateTime = types.StringValue(responsePasswordCredentials.GetEndDateTime().String())
			} else {
				tfStatePasswordCredential.EndDateTime = types.StringNull()
			}
			if responsePasswordCredentials.GetHint() != nil {
				tfStatePasswordCredential.Hint = types.StringValue(*responsePasswordCredentials.GetHint())
			} else {
				tfStatePasswordCredential.Hint = types.StringNull()
			}
			if responsePasswordCredentials.GetKeyId() != nil {
				tfStatePasswordCredential.KeyId = types.StringValue(responsePasswordCredentials.GetKeyId().String())
			} else {
				tfStatePasswordCredential.KeyId = types.StringNull()
			}
			if responsePasswordCredentials.GetSecretText() != nil {
				tfStatePasswordCredential.SecretText = types.StringValue(*responsePasswordCredentials.GetSecretText())
			} else {
				tfStatePasswordCredential.SecretText = types.StringNull()
			}
			if responsePasswordCredentials.GetStartDateTime() != nil {
				tfStatePasswordCredential.StartDateTime = types.StringValue(responsePasswordCredentials.GetStartDateTime().String())
			} else {
				tfStatePasswordCredential.StartDateTime = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStatePasswordCredential.AttributeTypes(), tfStatePasswordCredential)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.PasswordCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if response.GetPublicClient() != nil {
		tfStatePublicClientApplication := applicationPublicClientApplicationModel{}

		if len(response.GetPublicClient().GetRedirectUris()) > 0 {
			var valueArrayRedirectUris []attr.Value
			for _, responseRedirectUris := range response.GetPublicClient().GetRedirectUris() {
				valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(responseRedirectUris))
			}
			listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
			tfStatePublicClientApplication.RedirectUris = listValue
		} else {
			tfStatePublicClientApplication.RedirectUris = types.ListNull(types.StringType)
		}

		tfStateApplication.PublicClient, _ = types.ObjectValueFrom(ctx, tfStatePublicClientApplication.AttributeTypes(), tfStatePublicClientApplication)
	}
	if response.GetPublisherDomain() != nil {
		tfStateApplication.PublisherDomain = types.StringValue(*response.GetPublisherDomain())
	} else {
		tfStateApplication.PublisherDomain = types.StringNull()
	}
	if response.GetRequestSignatureVerification() != nil {
		tfStateRequestSignatureVerification := applicationRequestSignatureVerificationModel{}

		if response.GetRequestSignatureVerification().GetAllowedWeakAlgorithms() != nil {
			tfStateRequestSignatureVerification.AllowedWeakAlgorithms = types.StringValue(response.GetRequestSignatureVerification().GetAllowedWeakAlgorithms().String())
		} else {
			tfStateRequestSignatureVerification.AllowedWeakAlgorithms = types.StringNull()
		}
		if response.GetRequestSignatureVerification().GetIsSignedRequestRequired() != nil {
			tfStateRequestSignatureVerification.IsSignedRequestRequired = types.BoolValue(*response.GetRequestSignatureVerification().GetIsSignedRequestRequired())
		} else {
			tfStateRequestSignatureVerification.IsSignedRequestRequired = types.BoolNull()
		}

		tfStateApplication.RequestSignatureVerification, _ = types.ObjectValueFrom(ctx, tfStateRequestSignatureVerification.AttributeTypes(), tfStateRequestSignatureVerification)
	}
	if len(response.GetRequiredResourceAccess()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseRequiredResourceAccess := range response.GetRequiredResourceAccess() {
			tfStateRequiredResourceAccess := applicationRequiredResourceAccessModel{}

			if len(responseRequiredResourceAccess.GetResourceAccess()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseResourceAccess := range responseRequiredResourceAccess.GetResourceAccess() {
					tfStateResourceAccess := applicationResourceAccessModel{}

					if responseResourceAccess.GetId() != nil {
						tfStateResourceAccess.Id = types.StringValue(responseResourceAccess.GetId().String())
					} else {
						tfStateResourceAccess.Id = types.StringNull()
					}
					if responseResourceAccess.GetTypeEscaped() != nil {
						tfStateResourceAccess.Type = types.StringValue(*responseResourceAccess.GetTypeEscaped())
					} else {
						tfStateResourceAccess.Type = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateResourceAccess.AttributeTypes(), tfStateResourceAccess)
					objectValues = append(objectValues, objectValue)
				}
				tfStateRequiredResourceAccess.ResourceAccess, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if responseRequiredResourceAccess.GetResourceAppId() != nil {
				tfStateRequiredResourceAccess.ResourceAppId = types.StringValue(*responseRequiredResourceAccess.GetResourceAppId())
			} else {
				tfStateRequiredResourceAccess.ResourceAppId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateRequiredResourceAccess.AttributeTypes(), tfStateRequiredResourceAccess)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.RequiredResourceAccess, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if response.GetSamlMetadataUrl() != nil {
		tfStateApplication.SamlMetadataUrl = types.StringValue(*response.GetSamlMetadataUrl())
	} else {
		tfStateApplication.SamlMetadataUrl = types.StringNull()
	}
	if response.GetServiceManagementReference() != nil {
		tfStateApplication.ServiceManagementReference = types.StringValue(*response.GetServiceManagementReference())
	} else {
		tfStateApplication.ServiceManagementReference = types.StringNull()
	}
	if response.GetServicePrincipalLockConfiguration() != nil {
		tfStateServicePrincipalLockConfiguration := applicationServicePrincipalLockConfigurationModel{}

		if response.GetServicePrincipalLockConfiguration().GetAllProperties() != nil {
			tfStateServicePrincipalLockConfiguration.AllProperties = types.BoolValue(*response.GetServicePrincipalLockConfiguration().GetAllProperties())
		} else {
			tfStateServicePrincipalLockConfiguration.AllProperties = types.BoolNull()
		}
		if response.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign() != nil {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolValue(*response.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign())
		} else {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolNull()
		}
		if response.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify() != nil {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolValue(*response.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify())
		} else {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolNull()
		}
		if response.GetServicePrincipalLockConfiguration().GetIsEnabled() != nil {
			tfStateServicePrincipalLockConfiguration.IsEnabled = types.BoolValue(*response.GetServicePrincipalLockConfiguration().GetIsEnabled())
		} else {
			tfStateServicePrincipalLockConfiguration.IsEnabled = types.BoolNull()
		}
		if response.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId() != nil {
			tfStateServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolValue(*response.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId())
		} else {
			tfStateServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolNull()
		}

		tfStateApplication.ServicePrincipalLockConfiguration, _ = types.ObjectValueFrom(ctx, tfStateServicePrincipalLockConfiguration.AttributeTypes(), tfStateServicePrincipalLockConfiguration)
	}
	if response.GetSignInAudience() != nil {
		tfStateApplication.SignInAudience = types.StringValue(*response.GetSignInAudience())
	} else {
		tfStateApplication.SignInAudience = types.StringNull()
	}
	if response.GetSpa() != nil {
		tfStateSpaApplication := applicationSpaApplicationModel{}

		if len(response.GetSpa().GetRedirectUris()) > 0 {
			var valueArrayRedirectUris []attr.Value
			for _, responseRedirectUris := range response.GetSpa().GetRedirectUris() {
				valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(responseRedirectUris))
			}
			listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
			tfStateSpaApplication.RedirectUris = listValue
		} else {
			tfStateSpaApplication.RedirectUris = types.ListNull(types.StringType)
		}

		tfStateApplication.Spa, _ = types.ObjectValueFrom(ctx, tfStateSpaApplication.AttributeTypes(), tfStateSpaApplication)
	}
	if len(response.GetTags()) > 0 {
		var valueArrayTags []attr.Value
		for _, responseTags := range response.GetTags() {
			valueArrayTags = append(valueArrayTags, types.StringValue(responseTags))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayTags)
		tfStateApplication.Tags = listValue
	} else {
		tfStateApplication.Tags = types.ListNull(types.StringType)
	}
	if response.GetTokenEncryptionKeyId() != nil {
		tfStateApplication.TokenEncryptionKeyId = types.StringValue(response.GetTokenEncryptionKeyId().String())
	} else {
		tfStateApplication.TokenEncryptionKeyId = types.StringNull()
	}
	if response.GetUniqueName() != nil {
		tfStateApplication.UniqueName = types.StringValue(*response.GetUniqueName())
	} else {
		tfStateApplication.UniqueName = types.StringNull()
	}
	if response.GetVerifiedPublisher() != nil {
		tfStateVerifiedPublisher := applicationVerifiedPublisherModel{}

		if response.GetVerifiedPublisher().GetAddedDateTime() != nil {
			tfStateVerifiedPublisher.AddedDateTime = types.StringValue(response.GetVerifiedPublisher().GetAddedDateTime().String())
		} else {
			tfStateVerifiedPublisher.AddedDateTime = types.StringNull()
		}
		if response.GetVerifiedPublisher().GetDisplayName() != nil {
			tfStateVerifiedPublisher.DisplayName = types.StringValue(*response.GetVerifiedPublisher().GetDisplayName())
		} else {
			tfStateVerifiedPublisher.DisplayName = types.StringNull()
		}
		if response.GetVerifiedPublisher().GetVerifiedPublisherId() != nil {
			tfStateVerifiedPublisher.VerifiedPublisherId = types.StringValue(*response.GetVerifiedPublisher().GetVerifiedPublisherId())
		} else {
			tfStateVerifiedPublisher.VerifiedPublisherId = types.StringNull()
		}

		tfStateApplication.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfStateVerifiedPublisher.AttributeTypes(), tfStateVerifiedPublisher)
	}
	if response.GetWeb() != nil {
		tfStateWebApplication := applicationWebApplicationModel{}

		if response.GetWeb().GetHomePageUrl() != nil {
			tfStateWebApplication.HomePageUrl = types.StringValue(*response.GetWeb().GetHomePageUrl())
		} else {
			tfStateWebApplication.HomePageUrl = types.StringNull()
		}
		if response.GetWeb().GetImplicitGrantSettings() != nil {
			tfStateImplicitGrantSettings := applicationImplicitGrantSettingsModel{}

			if response.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance() != nil {
				tfStateImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolValue(*response.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance())
			} else {
				tfStateImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolNull()
			}
			if response.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance() != nil {
				tfStateImplicitGrantSettings.EnableIdTokenIssuance = types.BoolValue(*response.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance())
			} else {
				tfStateImplicitGrantSettings.EnableIdTokenIssuance = types.BoolNull()
			}

			tfStateWebApplication.ImplicitGrantSettings, _ = types.ObjectValueFrom(ctx, tfStateImplicitGrantSettings.AttributeTypes(), tfStateImplicitGrantSettings)
		}
		if response.GetWeb().GetLogoutUrl() != nil {
			tfStateWebApplication.LogoutUrl = types.StringValue(*response.GetWeb().GetLogoutUrl())
		} else {
			tfStateWebApplication.LogoutUrl = types.StringNull()
		}
		if len(response.GetWeb().GetRedirectUriSettings()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, responseRedirectUriSettings := range response.GetWeb().GetRedirectUriSettings() {
				tfStateRedirectUriSettings := applicationRedirectUriSettingsModel{}

				if responseRedirectUriSettings.GetUri() != nil {
					tfStateRedirectUriSettings.Uri = types.StringValue(*responseRedirectUriSettings.GetUri())
				} else {
					tfStateRedirectUriSettings.Uri = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStateRedirectUriSettings.AttributeTypes(), tfStateRedirectUriSettings)
				objectValues = append(objectValues, objectValue)
			}
			tfStateWebApplication.RedirectUriSettings, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(response.GetWeb().GetRedirectUris()) > 0 {
			var valueArrayRedirectUris []attr.Value
			for _, responseRedirectUris := range response.GetWeb().GetRedirectUris() {
				valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(responseRedirectUris))
			}
			listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
			tfStateWebApplication.RedirectUris = listValue
		} else {
			tfStateWebApplication.RedirectUris = types.ListNull(types.StringType)
		}

		tfStateApplication.Web, _ = types.ObjectValueFrom(ctx, tfStateWebApplication.AttributeTypes(), tfStateWebApplication)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateApplication)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
