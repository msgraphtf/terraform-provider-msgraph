package applications

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/applications"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"terraform-provider-msgraph/planmodifiers/boolplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/listplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/objectplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/stringplanmodifiers"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &applicationResource{}
	_ resource.ResourceWithConfigure = &applicationResource{}
)

// NewApplicationResource is a helper function to simplify the provider implementation.
func NewApplicationResource() resource.Resource {
	return &applicationResource{}
}

// applicationResource is the resource implementation.
type applicationResource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (d *applicationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

// Configure adds the provider configured client to the resource.
func (d *applicationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the resource.
func (d *applicationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"deleted_date_time": schema.StringAttribute{
				Description: "Date and time when this object was deleted. Always null when the object hasn't been deleted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"add_ins": schema.ListNestedAttribute{
				Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams can set the addIns property for its 'FileHandler' functionality. This lets services like Microsoft 365 call the application in the context of a document the user is working on.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier for the addIn object.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"properties": schema.ListNestedAttribute{
							Description: "The collection of key-value pairs that define parameters that the consuming service can use or call. You must specify this property when performing a POST or a PATCH operation on the addIns collection. Required.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.List{
								listplanmodifiers.UseStateForUnconfigured(),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"key": schema.StringAttribute{
										Description: "Key for the key-value pair.",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifiers.UseStateForUnconfigured(),
										},
									},
									"value": schema.StringAttribute{
										Description: "Value for the key-value pair.",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifiers.UseStateForUnconfigured(),
										},
									},
								},
							},
						},
						"type": schema.StringAttribute{
							Description: "The unique name for the functionality exposed by the app.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"api": schema.SingleNestedAttribute{
				Description: "Specifies settings for an application that implements a web API.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"accept_mapped_claims": schema.BoolAttribute{
						Description: "When true, allows an application to use claims mapping without specifying a custom signing key.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"known_client_applications": schema.ListAttribute{
						Description: "Used for bundling consent if you have a solution that contains two parts: a client app and a custom web API app. If you set the appID of the client app to this value, the user only consents once to the client app. Microsoft Entra ID knows that consenting to the client means implicitly consenting to the web API and automatically provisions service principals for both APIs at the same time. Both the client and the web API app must be registered in the same tenant.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						ElementType: types.StringType,
					},
					"oauth_2_permission_scopes": schema.ListNestedAttribute{
						Description: "The definition of the delegated permissions exposed by the web API represented by this application registration. These delegated permissions may be requested by a client application, and may be granted by users or administrators during consent. Delegated permissions are sometimes referred to as OAuth 2.0 scopes.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"admin_consent_description": schema.StringAttribute{
									Description: "A description of the delegated permissions, intended to be read by an administrator granting the permission on behalf of all users. This text appears in tenant-wide admin consent experiences.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"admin_consent_display_name": schema.StringAttribute{
									Description: "The permission's title, intended to be read by an administrator granting the permission on behalf of all users.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"id": schema.StringAttribute{
									Description: "Unique delegated permission identifier inside the collection of delegated permissions defined for a resource application.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"is_enabled": schema.BoolAttribute{
									Description: "When you create or update a permission, this property must be set to true (which is the default). To delete a permission, this property must first be set to false.  At that point, in a subsequent call, the permission may be removed.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"origin": schema.StringAttribute{
									Description: "",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"type": schema.StringAttribute{
									Description: "The possible values are: User and Admin. Specifies whether this delegated permission should be considered safe for non-admin users to consent to on behalf of themselves, or whether an administrator consent should always be required. While Microsoft Graph defines the default consent requirement for each permission, the tenant administrator may override the behavior in their organization (by allowing, restricting, or limiting user consent to this delegated permission). For more information, see Configure how users consent to applications.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"user_consent_description": schema.StringAttribute{
									Description: "A description of the delegated permissions, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"user_consent_display_name": schema.StringAttribute{
									Description: "A title for the permission, intended to be read by a user granting the permission on their own behalf. This text appears in consent experiences where the user is consenting only on behalf of themselves.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"value": schema.StringAttribute{
									Description: "Specifies the value to include in the scp (scope) claim in access tokens. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with ..",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
							},
						},
					},
					"pre_authorized_applications": schema.ListNestedAttribute{
						Description: "Lists the client applications that are preauthorized with the specified delegated permissions to access this application's APIs. Users aren't required to consent to any preauthorized application (for the permissions specified). However, any other permissions not listed in preAuthorizedApplications (requested through incremental consent for example) will require user consent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"app_id": schema.StringAttribute{
									Description: "The unique identifier for the application.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"delegated_permission_ids": schema.ListAttribute{
									Description: "The unique identifier for the oauth2PermissionScopes the application requires.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.List{
										listplanmodifiers.UseStateForUnconfigured(),
									},
									ElementType: types.StringType,
								},
							},
						},
					},
				},
			},
			"app_id": schema.StringAttribute{
				Description: "The unique identifier for the application that is assigned to an application by Microsoft Entra ID. Not nullable. Read-only. Alternate key. Supports $filter (eq).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"app_roles": schema.ListNestedAttribute{
				Description: "The collection of roles defined for the application. With app role assignments, these roles can be assigned to users, groups, or service principals associated with other applications. Not nullable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"allowed_member_types": schema.ListAttribute{
							Description: "Specifies whether this app role can be assigned to users and groups (by setting to ['User']), to other application's (by setting to ['Application'], or both (by setting to ['User', 'Application']). App roles supporting assignment to other applications' service principals are also known as application permissions. The 'Application' value is only supported for app roles defined on application entities.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.List{
								listplanmodifiers.UseStateForUnconfigured(),
							},
							ElementType: types.StringType,
						},
						"description": schema.StringAttribute{
							Description: "The description for the app role. This is displayed when the app role is being assigned and, if the app role functions as an application permission, during  consent experiences.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"display_name": schema.StringAttribute{
							Description: "Display name for the permission that appears in the app role assignment and consent experiences.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"id": schema.StringAttribute{
							Description: "Unique role identifier inside the appRoles collection. When creating a new app role, a new GUID identifier must be provided.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"is_enabled": schema.BoolAttribute{
							Description: "When creating or updating an app role, this must be set to true (which is the default). To delete a role, this must first be set to false.  At that point, in a subsequent call, this role may be removed.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"origin": schema.StringAttribute{
							Description: "Specifies if the app role is defined on the application object or on the servicePrincipal entity. Must not be included in any POST or PATCH requests. Read-only.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"value": schema.StringAttribute{
							Description: "Specifies the value to include in the roles claim in ID tokens and access tokens authenticating an assigned user or service principal. Must not exceed 120 characters in length. Allowed characters are : ! # $ % & ' ( ) * + , - . / : ;  =  ? @ [ ] ^ + _  {  } ~, and characters in the ranges 0-9, A-Z and a-z. Any other character, including the space character, aren't allowed. May not begin with ..",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"application_template_id": schema.StringAttribute{
				Description: "Unique identifier of the applicationTemplate. Supports $filter (eq, not, ne). Read-only. null if the app wasn't created from an application template.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"certification": schema.SingleNestedAttribute{
				Description: "Specifies the certification status of the application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"certification_details_url": schema.StringAttribute{
						Description: "URL that shows certification details for the application.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"certification_expiration_date_time": schema.StringAttribute{
						Description: "The timestamp when the current certification for the application expires.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"is_certified_by_microsoft": schema.BoolAttribute{
						Description: "Indicates whether the application is certified by Microsoft.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"is_publisher_attested": schema.BoolAttribute{
						Description: "Indicates whether the application has been self-attested by the application developer or the publisher.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"last_certification_date_time": schema.StringAttribute{
						Description: "The timestamp when the certification for the application was most recently added or updated.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"created_date_time": schema.StringAttribute{
				Description: "The date and time the application was registered. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.  Supports $filter (eq, ne, not, ge, le, in, and eq on null values) and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"default_redirect_uri": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"description": schema.StringAttribute{
				Description: "Free text field to provide a description of the application object to end users. The maximum allowed size is 1,024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"disabled_by_microsoft_status": schema.StringAttribute{
				Description: "Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name for the application. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"group_membership_claims": schema.StringAttribute{
				Description: "Configures the groups claim issued in a user or OAuth 2.0 access token that the application expects. To set this attribute, use one of the following valid string values: None, SecurityGroup (for security groups and Microsoft Entra roles), All (this gets all of the security groups, distribution groups, and Microsoft Entra directory roles that the signed-in user is a member of).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"identifier_uris": schema.ListAttribute{
				Description: "Also known as App ID URI, this value is set when an application is used as a resource app. The identifierUris acts as the prefix for the scopes you reference in your API's code, and it must be globally unique. You can use the default value provided, which is in the form api://<appId>, or specify a more readable URI like https://contoso.com/api. For more information on valid identifierUris patterns and best practices, see Microsoft Entra application registration security best practices. Not nullable. Supports $filter (eq, ne, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"info": schema.SingleNestedAttribute{
				Description: "Basic profile information of the application such as  app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"logo_url": schema.StringAttribute{
						Description: "CDN URL to the application's logo, Read-only.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"marketing_url": schema.StringAttribute{
						Description: "Link to the application's marketing page. For example, https://www.contoso.com/app/marketing",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"privacy_statement_url": schema.StringAttribute{
						Description: "Link to the application's privacy statement. For example, https://www.contoso.com/app/privacy",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"support_url": schema.StringAttribute{
						Description: "Link to the application's support page. For example, https://www.contoso.com/app/support",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"terms_of_service_url": schema.StringAttribute{
						Description: "Link to the application's terms of service statement. For example, https://www.contoso.com/app/termsofservice",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"is_device_only_auth_supported": schema.BoolAttribute{
				Description: "Specifies whether this application supports device authentication without a user. The default is false.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_fallback_public_client": schema.BoolAttribute{
				Description: "Specifies the fallback application type as public client, such as an installed application running on a mobile device. The default value is false, which means the fallback application type is confidential client such as a web app. There are certain scenarios where Microsoft Entra ID can't determine the client application type. For example, the ROPC flow where it's configured without specifying a redirect URI. In those cases, Microsoft Entra ID interprets the application type based on the value of this property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"key_credentials": schema.ListNestedAttribute{
				Description: "The collection of key credentials associated with the application. Not nullable. Supports $filter (eq, not, ge, le).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							Description: "A 40-character binary type that can be used to identify the credential. Optional. When not provided in the payload, defaults to the thumbprint of the certificate.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"display_name": schema.StringAttribute{
							Description: "The friendly name for the key, with a maximum length of 90 characters. Longer values are accepted but shortened. Optional.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"end_date_time": schema.StringAttribute{
							Description: "The date and time at which the credential expires. The DateTimeOffset type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"key": schema.StringAttribute{
							Description: "The certificate's raw data in byte array converted to Base64 string. Returned only on $select for a single object, that is, GET applications/{applicationId}?$select=keyCredentials or GET servicePrincipals/{servicePrincipalId}?$select=keyCredentials; otherwise, it's always null.  From a .cer certificate, you can read the key using the Convert.ToBase64String() method. For more information, see Get the certificate key.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"key_id": schema.StringAttribute{
							Description: "The unique identifier (GUID) for the key.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"start_date_time": schema.StringAttribute{
							Description: "The date and time at which the credential becomes valid.The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"type": schema.StringAttribute{
							Description: "The type of key credential; for example, Symmetric, AsymmetricX509Cert.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"usage": schema.StringAttribute{
							Description: "A string that describes the purpose for which the key can be used; for example, Verify.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"logo": schema.StringAttribute{
				Description: "The main logo for the application. Not nullable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"native_authentication_apis_enabled": schema.StringAttribute{
				Description: "Specifies whether the Native Authentication APIs are enabled for the application. The possible values are: none and all. Default is none. For more information, see Native Authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"notes": schema.StringAttribute{
				Description: "Notes relevant for the management of the application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"oauth_2_require_post_response": schema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"optional_claims": schema.SingleNestedAttribute{
				Description: "Application developers can configure optional claims in their Microsoft Entra applications to specify the claims that are sent to their application by the Microsoft security token service. For more information, see How to: Provide optional claims to your app.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"access_token": schema.ListNestedAttribute{
						Description: "The optional claims returned in the JWT access token.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"additional_properties": schema.ListAttribute{
									Description: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.List{
										listplanmodifiers.UseStateForUnconfigured(),
									},
									ElementType: types.StringType,
								},
								"essential": schema.BoolAttribute{
									Description: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"name": schema.StringAttribute{
									Description: "The name of the optional claim.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"source": schema.StringAttribute{
									Description: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
							},
						},
					},
					"id_token": schema.ListNestedAttribute{
						Description: "The optional claims returned in the JWT ID token.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"additional_properties": schema.ListAttribute{
									Description: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.List{
										listplanmodifiers.UseStateForUnconfigured(),
									},
									ElementType: types.StringType,
								},
								"essential": schema.BoolAttribute{
									Description: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"name": schema.StringAttribute{
									Description: "The name of the optional claim.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"source": schema.StringAttribute{
									Description: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
							},
						},
					},
					"saml_2_token": schema.ListNestedAttribute{
						Description: "The optional claims returned in the SAML token.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"additional_properties": schema.ListAttribute{
									Description: "Additional properties of the claim. If a property exists in this collection, it modifies the behavior of the optional claim specified in the name property.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.List{
										listplanmodifiers.UseStateForUnconfigured(),
									},
									ElementType: types.StringType,
								},
								"essential": schema.BoolAttribute{
									Description: "If the value is true, the claim specified by the client is necessary to ensure a smooth authorization experience for the specific task requested by the end user. The default value is false.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.Bool{
										boolplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"name": schema.StringAttribute{
									Description: "The name of the optional claim.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
								"source": schema.StringAttribute{
									Description: "The source (directory object) of the claim. There are predefined claims and user-defined claims from extension properties. If the source value is null, the claim is a predefined optional claim. If the source value is user, the value in the name property is the extension property from the user object.",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
							},
						},
					},
				},
			},
			"parental_control_settings": schema.SingleNestedAttribute{
				Description: "Specifies parental control settings for an application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"countries_blocked_for_minors": schema.ListAttribute{
						Description: "Specifies the two-letter ISO country codes. Access to the application will be blocked for minors from the countries specified in this list.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						ElementType: types.StringType,
					},
					"legal_age_group_rule": schema.StringAttribute{
						Description: "Specifies the legal age group rule that applies to users of the app. Can be set to one of the following values: ValueDescriptionAllowDefault. Enforces the legal minimum. This means parental consent is required for minors in the European Union and Korea.RequireConsentForPrivacyServicesEnforces the user to specify date of birth to comply with COPPA rules. RequireConsentForMinorsRequires parental consent for ages below 18, regardless of country minor rules.RequireConsentForKidsRequires parental consent for ages below 14, regardless of country minor rules.BlockMinorsBlocks minors from using the app.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"password_credentials": schema.ListNestedAttribute{
				Description: "The collection of password credentials associated with the application. Not nullable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"custom_key_identifier": schema.StringAttribute{
							Description: "Do not use.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"display_name": schema.StringAttribute{
							Description: "Friendly name for the password. Optional.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"end_date_time": schema.StringAttribute{
							Description: "The date and time at which the password expires represented using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"hint": schema.StringAttribute{
							Description: "Contains the first three characters of the password. Read-only.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"key_id": schema.StringAttribute{
							Description: "The unique identifier for the password.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"secret_text": schema.StringAttribute{
							Description: "Read-only; Contains the strong passwords generated by Microsoft Entra ID that are 16-64 characters in length. The generated password value is only returned during the initial POST request to addPassword. There is no way to retrieve this password in the future.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"start_date_time": schema.StringAttribute{
							Description: "The date and time at which the password becomes valid. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Optional.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"public_client": schema.SingleNestedAttribute{
				Description: "Specifies settings for installed clients such as desktop or mobile devices.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"redirect_uris": schema.ListAttribute{
						Description: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent. For iOS and macOS apps, specify the value following the syntax msauth.{BUNDLEID}://auth, replacing '{BUNDLEID}'. For example, if the bundle ID is com.microsoft.identitysample.MSALiOS, the URI is msauth.com.microsoft.identitysample.MSALiOS://auth.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						ElementType: types.StringType,
					},
				},
			},
			"publisher_domain": schema.StringAttribute{
				Description: "The verified publisher domain for the application. Read-only. For more information, see How to: Configure an application's publisher domain. Supports $filter (eq, ne, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"request_signature_verification": schema.SingleNestedAttribute{
				Description: "Specifies whether this application requires Microsoft Entra ID to verify the signed authentication requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"allowed_weak_algorithms": schema.StringAttribute{
						Description: "Specifies which weak algorithms are allowed.  The possible values are: rsaSha1, unknownFutureValue.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"is_signed_request_required": schema.BoolAttribute{
						Description: "Specifies whether signed authentication requests for this application should be required.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"required_resource_access": schema.ListNestedAttribute{
				Description: "Specifies the resources that the application needs to access. This property also specifies the set of delegated permissions and application roles that it needs for each of those resources. This configuration of access to the required resources drives the consent experience. No more than 50 resource services (APIs) can be configured. Beginning mid-October 2021, the total number of required permissions must not exceed 400. For more information, see Limits on requested permissions per app. Not nullable. Supports $filter (eq, not, ge, le).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"resource_access": schema.ListNestedAttribute{
							Description: "The list of OAuth2.0 permission scopes and app roles that the application requires from the specified resource.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.List{
								listplanmodifiers.UseStateForUnconfigured(),
							},
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Description: "The unique identifier of an app role or delegated permission exposed by the resource application. For delegated permissions, this should match the id property of one of the delegated permissions in the oauth2PermissionScopes collection of the resource application's service principal. For app roles (application permissions), this should match the id property of an app role in the appRoles collection of the resource application's service principal.",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifiers.UseStateForUnconfigured(),
										},
									},
									"type": schema.StringAttribute{
										Description: "Specifies whether the id property references a delegated permission or an app role (application permission). The possible values are: Scope (for delegated permissions) or Role (for app roles).",
										Optional:    true,
										Computed:    true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifiers.UseStateForUnconfigured(),
										},
									},
								},
							},
						},
						"resource_app_id": schema.StringAttribute{
							Description: "The unique identifier for the resource that the application requires access to. This should be equal to the appId declared on the target resource application.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"saml_metadata_url": schema.StringAttribute{
				Description: "The URL where the service exposes SAML metadata for federation. This property is valid only for single-tenant applications. Nullable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"service_management_reference": schema.StringAttribute{
				Description: "References application or service contact information from a Service or Asset Management database. Nullable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"service_principal_lock_configuration": schema.SingleNestedAttribute{
				Description: "Specifies whether sensitive properties of a multitenant application should be locked for editing after the application is provisioned in a tenant. Nullable. null by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"all_properties": schema.BoolAttribute{
						Description: "Enables locking all sensitive properties. The sensitive properties are keyCredentials, passwordCredentials, and tokenEncryptionKeyId.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"credentials_with_usage_sign": schema.BoolAttribute{
						Description: "Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Sign.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"credentials_with_usage_verify": schema.BoolAttribute{
						Description: "Locks the keyCredentials and passwordCredentials properties for modification where credential usage type is Verify. This locks OAuth service principals.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"is_enabled": schema.BoolAttribute{
						Description: "Enables or disables service principal lock configuration. To allow the sensitive properties to be updated, update this property to false to disable the lock on the service principal.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"token_encryption_key_id": schema.BoolAttribute{
						Description: "Locks the tokenEncryptionKeyId property for modification on the service principal.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"sign_in_audience": schema.StringAttribute{
				Description: "Specifies the Microsoft accounts that are supported for the current application. The possible values are: AzureADMyOrg (default), AzureADMultipleOrgs, AzureADandPersonalMicrosoftAccount, and PersonalMicrosoftAccount. See more in the table. The value of this object also limits the number of permissions an app can request. For more information, see Limits on requested permissions per app. The value for this property has implications on other app object properties. As a result, if you change this property, you might need to change other properties first. For more information, see Validation differences for signInAudience.Supports $filter (eq, ne, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"spa": schema.SingleNestedAttribute{
				Description: "Specifies settings for a single-page application, including sign out URLs and redirect URIs for authorization codes and access tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"redirect_uris": schema.ListAttribute{
						Description: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						ElementType: types.StringType,
					},
				},
			},
			"tags": schema.ListAttribute{
				Description: "Custom strings that can be used to categorize and identify the application. Not nullable. Strings added here will also appear in the tags property of any associated service principals.Supports $filter (eq, not, ge, le, startsWith) and $search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"token_encryption_key_id": schema.StringAttribute{
				Description: "Specifies the keyId of a public key from the keyCredentials collection. When configured, Microsoft Entra ID encrypts all the tokens it emits by using the key this property points to. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"unique_name": schema.StringAttribute{
				Description: "The unique identifier that can be assigned to an application and used as an alternate key. Immutable. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"verified_publisher": schema.SingleNestedAttribute{
				Description: "Specifies the verified publisher of the application. For more information about how publisher verification helps support application security, trustworthiness, and compliance, see Publisher verification.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"added_date_time": schema.StringAttribute{
						Description: "The timestamp when the verified publisher was first added or most recently updated.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"display_name": schema.StringAttribute{
						Description: "The verified publisher name from the app publisher's Partner Center account.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"verified_publisher_id": schema.StringAttribute{
						Description: "The ID of the verified publisher from the app publisher's Partner Center account.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"web": schema.SingleNestedAttribute{
				Description: "Specifies settings for a web application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"home_page_url": schema.StringAttribute{
						Description: "Home page or landing page of the application.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"implicit_grant_settings": schema.SingleNestedAttribute{
						Description: "Specifies whether this web application can request tokens using the OAuth 2.0 implicit flow.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifiers.UseStateForUnconfigured(),
						},
						Attributes: map[string]schema.Attribute{
							"enable_access_token_issuance": schema.BoolAttribute{
								Description: "Specifies whether this web application can request an access token using the OAuth 2.0 implicit flow.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifiers.UseStateForUnconfigured(),
								},
							},
							"enable_id_token_issuance": schema.BoolAttribute{
								Description: "Specifies whether this web application can request an ID token using the OAuth 2.0 implicit flow.",
								Optional:    true,
								Computed:    true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifiers.UseStateForUnconfigured(),
								},
							},
						},
					},
					"logout_url": schema.StringAttribute{
						Description: "Specifies the URL that is used by Microsoft's authorization service to log out a user using front-channel, back-channel or SAML logout protocols.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"redirect_uri_settings": schema.ListNestedAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"uri": schema.StringAttribute{
									Description: "",
									Optional:    true,
									Computed:    true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifiers.UseStateForUnconfigured(),
									},
								},
							},
						},
					},
					"redirect_uris": schema.ListAttribute{
						Description: "Specifies the URLs where user tokens are sent for sign-in, or the redirect URIs where OAuth 2.0 authorization codes and access tokens are sent.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						ElementType: types.StringType,
					},
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *applicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanApplication applicationModel
	diags := req.Plan.Get(ctx, &tfPlanApplication)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	requestBodyApplication := models.NewApplication()
	// START Id | CreateStringAttribute
	if !tfPlanApplication.Id.IsUnknown() {
		tfPlanId := tfPlanApplication.Id.ValueString()
		requestBodyApplication.SetId(&tfPlanId)
	} else {
		tfPlanApplication.Id = types.StringNull()
	}
	// END Id | CreateStringAttribute

	// START DeletedDateTime | CreateStringTimeAttribute
	if !tfPlanApplication.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlanApplication.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyApplication.SetDeletedDateTime(&t)
	} else {
		tfPlanApplication.DeletedDateTime = types.StringNull()
	}
	// END DeletedDateTime | CreateStringTimeAttribute

	// START AddIns | CreateArrayObjectAttribute
	if len(tfPlanApplication.AddIns.Elements()) > 0 {
		var requestBodyAddIns []models.AddInable
		for _, i := range tfPlanApplication.AddIns.Elements() {
			requestBodyAddIn := models.NewAddIn()
			tfPlanAddIn := applicationAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAddIn)

			// START Id | CreateStringUuidAttribute
			if !tfPlanAddIn.Id.IsUnknown() {
				tfPlanId := tfPlanAddIn.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAddIn.SetId(&u)
			} else {
				tfPlanAddIn.Id = types.StringNull()
			}
			// END Id | CreateStringUuidAttribute

			// START Properties | CreateArrayObjectAttribute
			if len(tfPlanAddIn.Properties.Elements()) > 0 {
				var requestBodyProperties []models.KeyValueable
				for _, i := range tfPlanAddIn.Properties.Elements() {
					requestBodyKeyValue := models.NewKeyValue()
					tfPlanKeyValue := applicationKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyValue)

					// START Key | CreateStringAttribute
					if !tfPlanKeyValue.Key.IsUnknown() {
						tfPlanKey := tfPlanKeyValue.Key.ValueString()
						requestBodyKeyValue.SetKey(&tfPlanKey)
					} else {
						tfPlanKeyValue.Key = types.StringNull()
					}
					// END Key | CreateStringAttribute

					// START Value | CreateStringAttribute
					if !tfPlanKeyValue.Value.IsUnknown() {
						tfPlanValue := tfPlanKeyValue.Value.ValueString()
						requestBodyKeyValue.SetValue(&tfPlanValue)
					} else {
						tfPlanKeyValue.Value = types.StringNull()
					}
					// END Value | CreateStringAttribute

				}
				requestBodyAddIn.SetProperties(requestBodyProperties)
			} else {
				tfPlanAddIn.Properties = types.ListNull(tfPlanAddIn.Properties.ElementType(ctx))
			}
			// END Properties | CreateArrayObjectAttribute

			// START Type | CreateStringAttribute
			if !tfPlanAddIn.Type.IsUnknown() {
				tfPlanType := tfPlanAddIn.Type.ValueString()
				requestBodyAddIn.SetTypeEscaped(&tfPlanType)
			} else {
				tfPlanAddIn.Type = types.StringNull()
			}
			// END Type | CreateStringAttribute

		}
		requestBodyApplication.SetAddIns(requestBodyAddIns)
	} else {
		tfPlanApplication.AddIns = types.ListNull(tfPlanApplication.AddIns.ElementType(ctx))
	}
	// END AddIns | CreateArrayObjectAttribute

	// START Api | CreateObjectAttribute
	if !tfPlanApplication.Api.IsUnknown() {
		requestBodyApiApplication := models.NewApiApplication()
		tfPlanApiApplication := applicationApiApplicationModel{}
		tfPlanApplication.Api.As(ctx, &tfPlanApiApplication, basetypes.ObjectAsOptions{})

		// START AcceptMappedClaims | CreateBoolAttribute
		if !tfPlanApiApplication.AcceptMappedClaims.IsUnknown() {
			tfPlanAcceptMappedClaims := tfPlanApiApplication.AcceptMappedClaims.ValueBool()
			requestBodyApiApplication.SetAcceptMappedClaims(&tfPlanAcceptMappedClaims)
		} else {
			tfPlanApiApplication.AcceptMappedClaims = types.BoolNull()
		}
		// END AcceptMappedClaims | CreateBoolAttribute

		// START KnownClientApplications | CreateArrayUuidAttribute
		if len(tfPlanApiApplication.KnownClientApplications.Elements()) > 0 {
			var uuidArrayKnownClientApplications []uuid.UUID
			for _, i := range tfPlanApiApplication.KnownClientApplications.Elements() {
				u, _ := uuid.Parse(i.String())
				uuidArrayKnownClientApplications = append(uuidArrayKnownClientApplications, u)
			}
			requestBodyApiApplication.SetKnownClientApplications(uuidArrayKnownClientApplications)
		} else {
			tfPlanApiApplication.KnownClientApplications = types.ListNull(types.StringType)
		}

		// END KnownClientApplications | CreateArrayUuidAttribute

		// START Oauth2PermissionScopes | CreateArrayObjectAttribute
		if len(tfPlanApiApplication.Oauth2PermissionScopes.Elements()) > 0 {
			var requestBodyOauth2PermissionScopes []models.PermissionScopeable
			for _, i := range tfPlanApiApplication.Oauth2PermissionScopes.Elements() {
				requestBodyPermissionScope := models.NewPermissionScope()
				tfPlanPermissionScope := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPermissionScope)

				// START AdminConsentDescription | CreateStringAttribute
				if !tfPlanPermissionScope.AdminConsentDescription.IsUnknown() {
					tfPlanAdminConsentDescription := tfPlanPermissionScope.AdminConsentDescription.ValueString()
					requestBodyPermissionScope.SetAdminConsentDescription(&tfPlanAdminConsentDescription)
				} else {
					tfPlanPermissionScope.AdminConsentDescription = types.StringNull()
				}
				// END AdminConsentDescription | CreateStringAttribute

				// START AdminConsentDisplayName | CreateStringAttribute
				if !tfPlanPermissionScope.AdminConsentDisplayName.IsUnknown() {
					tfPlanAdminConsentDisplayName := tfPlanPermissionScope.AdminConsentDisplayName.ValueString()
					requestBodyPermissionScope.SetAdminConsentDisplayName(&tfPlanAdminConsentDisplayName)
				} else {
					tfPlanPermissionScope.AdminConsentDisplayName = types.StringNull()
				}
				// END AdminConsentDisplayName | CreateStringAttribute

				// START Id | CreateStringUuidAttribute
				if !tfPlanPermissionScope.Id.IsUnknown() {
					tfPlanId := tfPlanPermissionScope.Id.ValueString()
					u, _ := uuid.Parse(tfPlanId)
					requestBodyPermissionScope.SetId(&u)
				} else {
					tfPlanPermissionScope.Id = types.StringNull()
				}
				// END Id | CreateStringUuidAttribute

				// START IsEnabled | CreateBoolAttribute
				if !tfPlanPermissionScope.IsEnabled.IsUnknown() {
					tfPlanIsEnabled := tfPlanPermissionScope.IsEnabled.ValueBool()
					requestBodyPermissionScope.SetIsEnabled(&tfPlanIsEnabled)
				} else {
					tfPlanPermissionScope.IsEnabled = types.BoolNull()
				}
				// END IsEnabled | CreateBoolAttribute

				// START Origin | CreateStringAttribute
				if !tfPlanPermissionScope.Origin.IsUnknown() {
					tfPlanOrigin := tfPlanPermissionScope.Origin.ValueString()
					requestBodyPermissionScope.SetOrigin(&tfPlanOrigin)
				} else {
					tfPlanPermissionScope.Origin = types.StringNull()
				}
				// END Origin | CreateStringAttribute

				// START Type | CreateStringAttribute
				if !tfPlanPermissionScope.Type.IsUnknown() {
					tfPlanType := tfPlanPermissionScope.Type.ValueString()
					requestBodyPermissionScope.SetTypeEscaped(&tfPlanType)
				} else {
					tfPlanPermissionScope.Type = types.StringNull()
				}
				// END Type | CreateStringAttribute

				// START UserConsentDescription | CreateStringAttribute
				if !tfPlanPermissionScope.UserConsentDescription.IsUnknown() {
					tfPlanUserConsentDescription := tfPlanPermissionScope.UserConsentDescription.ValueString()
					requestBodyPermissionScope.SetUserConsentDescription(&tfPlanUserConsentDescription)
				} else {
					tfPlanPermissionScope.UserConsentDescription = types.StringNull()
				}
				// END UserConsentDescription | CreateStringAttribute

				// START UserConsentDisplayName | CreateStringAttribute
				if !tfPlanPermissionScope.UserConsentDisplayName.IsUnknown() {
					tfPlanUserConsentDisplayName := tfPlanPermissionScope.UserConsentDisplayName.ValueString()
					requestBodyPermissionScope.SetUserConsentDisplayName(&tfPlanUserConsentDisplayName)
				} else {
					tfPlanPermissionScope.UserConsentDisplayName = types.StringNull()
				}
				// END UserConsentDisplayName | CreateStringAttribute

				// START Value | CreateStringAttribute
				if !tfPlanPermissionScope.Value.IsUnknown() {
					tfPlanValue := tfPlanPermissionScope.Value.ValueString()
					requestBodyPermissionScope.SetValue(&tfPlanValue)
				} else {
					tfPlanPermissionScope.Value = types.StringNull()
				}
				// END Value | CreateStringAttribute

			}
			requestBodyApiApplication.SetOauth2PermissionScopes(requestBodyOauth2PermissionScopes)
		} else {
			tfPlanApiApplication.Oauth2PermissionScopes = types.ListNull(tfPlanApiApplication.Oauth2PermissionScopes.ElementType(ctx))
		}
		// END Oauth2PermissionScopes | CreateArrayObjectAttribute

		// START PreAuthorizedApplications | CreateArrayObjectAttribute
		if len(tfPlanApiApplication.PreAuthorizedApplications.Elements()) > 0 {
			var requestBodyPreAuthorizedApplications []models.PreAuthorizedApplicationable
			for _, i := range tfPlanApiApplication.PreAuthorizedApplications.Elements() {
				requestBodyPreAuthorizedApplication := models.NewPreAuthorizedApplication()
				tfPlanPreAuthorizedApplication := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPreAuthorizedApplication)

				// START AppId | CreateStringAttribute
				if !tfPlanPreAuthorizedApplication.AppId.IsUnknown() {
					tfPlanAppId := tfPlanPreAuthorizedApplication.AppId.ValueString()
					requestBodyPreAuthorizedApplication.SetAppId(&tfPlanAppId)
				} else {
					tfPlanPreAuthorizedApplication.AppId = types.StringNull()
				}
				// END AppId | CreateStringAttribute

				// START DelegatedPermissionIds | CreateArrayStringAttribute
				if len(tfPlanPreAuthorizedApplication.DelegatedPermissionIds.Elements()) > 0 {
					var stringArrayDelegatedPermissionIds []string
					for _, i := range tfPlanPreAuthorizedApplication.DelegatedPermissionIds.Elements() {
						stringArrayDelegatedPermissionIds = append(stringArrayDelegatedPermissionIds, i.String())
					}
					requestBodyPreAuthorizedApplication.SetDelegatedPermissionIds(stringArrayDelegatedPermissionIds)
				} else {
					tfPlanPreAuthorizedApplication.DelegatedPermissionIds = types.ListNull(types.StringType)
				}
				// END DelegatedPermissionIds | CreateArrayStringAttribute

			}
			requestBodyApiApplication.SetPreAuthorizedApplications(requestBodyPreAuthorizedApplications)
		} else {
			tfPlanApiApplication.PreAuthorizedApplications = types.ListNull(tfPlanApiApplication.PreAuthorizedApplications.ElementType(ctx))
		}
		// END PreAuthorizedApplications | CreateArrayObjectAttribute

		// START RequestedAccessTokenVersion | UNKNOWN
		// END RequestedAccessTokenVersion | UNKNOWN

		requestBodyApplication.SetApi(requestBodyApiApplication)
		tfPlanApplication.Api, _ = types.ObjectValueFrom(ctx, tfPlanApiApplication.AttributeTypes(), requestBodyApiApplication)
	} else {
		tfPlanApplication.Api = types.ObjectNull(tfPlanApplication.Api.AttributeTypes(ctx))
	}
	// END Api | CreateObjectAttribute

	// START AppId | CreateStringAttribute
	if !tfPlanApplication.AppId.IsUnknown() {
		tfPlanAppId := tfPlanApplication.AppId.ValueString()
		requestBodyApplication.SetAppId(&tfPlanAppId)
	} else {
		tfPlanApplication.AppId = types.StringNull()
	}
	// END AppId | CreateStringAttribute

	// START AppRoles | CreateArrayObjectAttribute
	if len(tfPlanApplication.AppRoles.Elements()) > 0 {
		var requestBodyAppRoles []models.AppRoleable
		for _, i := range tfPlanApplication.AppRoles.Elements() {
			requestBodyAppRole := models.NewAppRole()
			tfPlanAppRole := applicationAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAppRole)

			// START AllowedMemberTypes | CreateArrayStringAttribute
			if len(tfPlanAppRole.AllowedMemberTypes.Elements()) > 0 {
				var stringArrayAllowedMemberTypes []string
				for _, i := range tfPlanAppRole.AllowedMemberTypes.Elements() {
					stringArrayAllowedMemberTypes = append(stringArrayAllowedMemberTypes, i.String())
				}
				requestBodyAppRole.SetAllowedMemberTypes(stringArrayAllowedMemberTypes)
			} else {
				tfPlanAppRole.AllowedMemberTypes = types.ListNull(types.StringType)
			}
			// END AllowedMemberTypes | CreateArrayStringAttribute

			// START Description | CreateStringAttribute
			if !tfPlanAppRole.Description.IsUnknown() {
				tfPlanDescription := tfPlanAppRole.Description.ValueString()
				requestBodyAppRole.SetDescription(&tfPlanDescription)
			} else {
				tfPlanAppRole.Description = types.StringNull()
			}
			// END Description | CreateStringAttribute

			// START DisplayName | CreateStringAttribute
			if !tfPlanAppRole.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanAppRole.DisplayName.ValueString()
				requestBodyAppRole.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanAppRole.DisplayName = types.StringNull()
			}
			// END DisplayName | CreateStringAttribute

			// START Id | CreateStringUuidAttribute
			if !tfPlanAppRole.Id.IsUnknown() {
				tfPlanId := tfPlanAppRole.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAppRole.SetId(&u)
			} else {
				tfPlanAppRole.Id = types.StringNull()
			}
			// END Id | CreateStringUuidAttribute

			// START IsEnabled | CreateBoolAttribute
			if !tfPlanAppRole.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := tfPlanAppRole.IsEnabled.ValueBool()
				requestBodyAppRole.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				tfPlanAppRole.IsEnabled = types.BoolNull()
			}
			// END IsEnabled | CreateBoolAttribute

			// START Origin | CreateStringAttribute
			if !tfPlanAppRole.Origin.IsUnknown() {
				tfPlanOrigin := tfPlanAppRole.Origin.ValueString()
				requestBodyAppRole.SetOrigin(&tfPlanOrigin)
			} else {
				tfPlanAppRole.Origin = types.StringNull()
			}
			// END Origin | CreateStringAttribute

			// START Value | CreateStringAttribute
			if !tfPlanAppRole.Value.IsUnknown() {
				tfPlanValue := tfPlanAppRole.Value.ValueString()
				requestBodyAppRole.SetValue(&tfPlanValue)
			} else {
				tfPlanAppRole.Value = types.StringNull()
			}
			// END Value | CreateStringAttribute

		}
		requestBodyApplication.SetAppRoles(requestBodyAppRoles)
	} else {
		tfPlanApplication.AppRoles = types.ListNull(tfPlanApplication.AppRoles.ElementType(ctx))
	}
	// END AppRoles | CreateArrayObjectAttribute

	// START ApplicationTemplateId | CreateStringAttribute
	if !tfPlanApplication.ApplicationTemplateId.IsUnknown() {
		tfPlanApplicationTemplateId := tfPlanApplication.ApplicationTemplateId.ValueString()
		requestBodyApplication.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	} else {
		tfPlanApplication.ApplicationTemplateId = types.StringNull()
	}
	// END ApplicationTemplateId | CreateStringAttribute

	// START Certification | CreateObjectAttribute
	if !tfPlanApplication.Certification.IsUnknown() {
		requestBodyCertification := models.NewCertification()
		tfPlanCertification := applicationCertificationModel{}
		tfPlanApplication.Certification.As(ctx, &tfPlanCertification, basetypes.ObjectAsOptions{})

		// START CertificationDetailsUrl | CreateStringAttribute
		if !tfPlanCertification.CertificationDetailsUrl.IsUnknown() {
			tfPlanCertificationDetailsUrl := tfPlanCertification.CertificationDetailsUrl.ValueString()
			requestBodyCertification.SetCertificationDetailsUrl(&tfPlanCertificationDetailsUrl)
		} else {
			tfPlanCertification.CertificationDetailsUrl = types.StringNull()
		}
		// END CertificationDetailsUrl | CreateStringAttribute

		// START CertificationExpirationDateTime | CreateStringTimeAttribute
		if !tfPlanCertification.CertificationExpirationDateTime.IsUnknown() {
			tfPlanCertificationExpirationDateTime := tfPlanCertification.CertificationExpirationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanCertificationExpirationDateTime)
			requestBodyCertification.SetCertificationExpirationDateTime(&t)
		} else {
			tfPlanCertification.CertificationExpirationDateTime = types.StringNull()
		}
		// END CertificationExpirationDateTime | CreateStringTimeAttribute

		// START IsCertifiedByMicrosoft | CreateBoolAttribute
		if !tfPlanCertification.IsCertifiedByMicrosoft.IsUnknown() {
			tfPlanIsCertifiedByMicrosoft := tfPlanCertification.IsCertifiedByMicrosoft.ValueBool()
			requestBodyCertification.SetIsCertifiedByMicrosoft(&tfPlanIsCertifiedByMicrosoft)
		} else {
			tfPlanCertification.IsCertifiedByMicrosoft = types.BoolNull()
		}
		// END IsCertifiedByMicrosoft | CreateBoolAttribute

		// START IsPublisherAttested | CreateBoolAttribute
		if !tfPlanCertification.IsPublisherAttested.IsUnknown() {
			tfPlanIsPublisherAttested := tfPlanCertification.IsPublisherAttested.ValueBool()
			requestBodyCertification.SetIsPublisherAttested(&tfPlanIsPublisherAttested)
		} else {
			tfPlanCertification.IsPublisherAttested = types.BoolNull()
		}
		// END IsPublisherAttested | CreateBoolAttribute

		// START LastCertificationDateTime | CreateStringTimeAttribute
		if !tfPlanCertification.LastCertificationDateTime.IsUnknown() {
			tfPlanLastCertificationDateTime := tfPlanCertification.LastCertificationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastCertificationDateTime)
			requestBodyCertification.SetLastCertificationDateTime(&t)
		} else {
			tfPlanCertification.LastCertificationDateTime = types.StringNull()
		}
		// END LastCertificationDateTime | CreateStringTimeAttribute

		requestBodyApplication.SetCertification(requestBodyCertification)
		tfPlanApplication.Certification, _ = types.ObjectValueFrom(ctx, tfPlanCertification.AttributeTypes(), requestBodyCertification)
	} else {
		tfPlanApplication.Certification = types.ObjectNull(tfPlanApplication.Certification.AttributeTypes(ctx))
	}
	// END Certification | CreateObjectAttribute

	// START CreatedDateTime | CreateStringTimeAttribute
	if !tfPlanApplication.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlanApplication.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyApplication.SetCreatedDateTime(&t)
	} else {
		tfPlanApplication.CreatedDateTime = types.StringNull()
	}
	// END CreatedDateTime | CreateStringTimeAttribute

	// START DefaultRedirectUri | CreateStringAttribute
	if !tfPlanApplication.DefaultRedirectUri.IsUnknown() {
		tfPlanDefaultRedirectUri := tfPlanApplication.DefaultRedirectUri.ValueString()
		requestBodyApplication.SetDefaultRedirectUri(&tfPlanDefaultRedirectUri)
	} else {
		tfPlanApplication.DefaultRedirectUri = types.StringNull()
	}
	// END DefaultRedirectUri | CreateStringAttribute

	// START Description | CreateStringAttribute
	if !tfPlanApplication.Description.IsUnknown() {
		tfPlanDescription := tfPlanApplication.Description.ValueString()
		requestBodyApplication.SetDescription(&tfPlanDescription)
	} else {
		tfPlanApplication.Description = types.StringNull()
	}
	// END Description | CreateStringAttribute

	// START DisabledByMicrosoftStatus | CreateStringAttribute
	if !tfPlanApplication.DisabledByMicrosoftStatus.IsUnknown() {
		tfPlanDisabledByMicrosoftStatus := tfPlanApplication.DisabledByMicrosoftStatus.ValueString()
		requestBodyApplication.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	} else {
		tfPlanApplication.DisabledByMicrosoftStatus = types.StringNull()
	}
	// END DisabledByMicrosoftStatus | CreateStringAttribute

	// START DisplayName | CreateStringAttribute
	if !tfPlanApplication.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanApplication.DisplayName.ValueString()
		requestBodyApplication.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanApplication.DisplayName = types.StringNull()
	}
	// END DisplayName | CreateStringAttribute

	// START GroupMembershipClaims | CreateStringAttribute
	if !tfPlanApplication.GroupMembershipClaims.IsUnknown() {
		tfPlanGroupMembershipClaims := tfPlanApplication.GroupMembershipClaims.ValueString()
		requestBodyApplication.SetGroupMembershipClaims(&tfPlanGroupMembershipClaims)
	} else {
		tfPlanApplication.GroupMembershipClaims = types.StringNull()
	}
	// END GroupMembershipClaims | CreateStringAttribute

	// START IdentifierUris | CreateArrayStringAttribute
	if len(tfPlanApplication.IdentifierUris.Elements()) > 0 {
		var stringArrayIdentifierUris []string
		for _, i := range tfPlanApplication.IdentifierUris.Elements() {
			stringArrayIdentifierUris = append(stringArrayIdentifierUris, i.String())
		}
		requestBodyApplication.SetIdentifierUris(stringArrayIdentifierUris)
	} else {
		tfPlanApplication.IdentifierUris = types.ListNull(types.StringType)
	}
	// END IdentifierUris | CreateArrayStringAttribute

	// START Info | CreateObjectAttribute
	if !tfPlanApplication.Info.IsUnknown() {
		requestBodyInformationalUrl := models.NewInformationalUrl()
		tfPlanInformationalUrl := applicationInformationalUrlModel{}
		tfPlanApplication.Info.As(ctx, &tfPlanInformationalUrl, basetypes.ObjectAsOptions{})

		// START LogoUrl | CreateStringAttribute
		if !tfPlanInformationalUrl.LogoUrl.IsUnknown() {
			tfPlanLogoUrl := tfPlanInformationalUrl.LogoUrl.ValueString()
			requestBodyInformationalUrl.SetLogoUrl(&tfPlanLogoUrl)
		} else {
			tfPlanInformationalUrl.LogoUrl = types.StringNull()
		}
		// END LogoUrl | CreateStringAttribute

		// START MarketingUrl | CreateStringAttribute
		if !tfPlanInformationalUrl.MarketingUrl.IsUnknown() {
			tfPlanMarketingUrl := tfPlanInformationalUrl.MarketingUrl.ValueString()
			requestBodyInformationalUrl.SetMarketingUrl(&tfPlanMarketingUrl)
		} else {
			tfPlanInformationalUrl.MarketingUrl = types.StringNull()
		}
		// END MarketingUrl | CreateStringAttribute

		// START PrivacyStatementUrl | CreateStringAttribute
		if !tfPlanInformationalUrl.PrivacyStatementUrl.IsUnknown() {
			tfPlanPrivacyStatementUrl := tfPlanInformationalUrl.PrivacyStatementUrl.ValueString()
			requestBodyInformationalUrl.SetPrivacyStatementUrl(&tfPlanPrivacyStatementUrl)
		} else {
			tfPlanInformationalUrl.PrivacyStatementUrl = types.StringNull()
		}
		// END PrivacyStatementUrl | CreateStringAttribute

		// START SupportUrl | CreateStringAttribute
		if !tfPlanInformationalUrl.SupportUrl.IsUnknown() {
			tfPlanSupportUrl := tfPlanInformationalUrl.SupportUrl.ValueString()
			requestBodyInformationalUrl.SetSupportUrl(&tfPlanSupportUrl)
		} else {
			tfPlanInformationalUrl.SupportUrl = types.StringNull()
		}
		// END SupportUrl | CreateStringAttribute

		// START TermsOfServiceUrl | CreateStringAttribute
		if !tfPlanInformationalUrl.TermsOfServiceUrl.IsUnknown() {
			tfPlanTermsOfServiceUrl := tfPlanInformationalUrl.TermsOfServiceUrl.ValueString()
			requestBodyInformationalUrl.SetTermsOfServiceUrl(&tfPlanTermsOfServiceUrl)
		} else {
			tfPlanInformationalUrl.TermsOfServiceUrl = types.StringNull()
		}
		// END TermsOfServiceUrl | CreateStringAttribute

		requestBodyApplication.SetInfo(requestBodyInformationalUrl)
		tfPlanApplication.Info, _ = types.ObjectValueFrom(ctx, tfPlanInformationalUrl.AttributeTypes(), requestBodyInformationalUrl)
	} else {
		tfPlanApplication.Info = types.ObjectNull(tfPlanApplication.Info.AttributeTypes(ctx))
	}
	// END Info | CreateObjectAttribute

	// START IsDeviceOnlyAuthSupported | CreateBoolAttribute
	if !tfPlanApplication.IsDeviceOnlyAuthSupported.IsUnknown() {
		tfPlanIsDeviceOnlyAuthSupported := tfPlanApplication.IsDeviceOnlyAuthSupported.ValueBool()
		requestBodyApplication.SetIsDeviceOnlyAuthSupported(&tfPlanIsDeviceOnlyAuthSupported)
	} else {
		tfPlanApplication.IsDeviceOnlyAuthSupported = types.BoolNull()
	}
	// END IsDeviceOnlyAuthSupported | CreateBoolAttribute

	// START IsFallbackPublicClient | CreateBoolAttribute
	if !tfPlanApplication.IsFallbackPublicClient.IsUnknown() {
		tfPlanIsFallbackPublicClient := tfPlanApplication.IsFallbackPublicClient.ValueBool()
		requestBodyApplication.SetIsFallbackPublicClient(&tfPlanIsFallbackPublicClient)
	} else {
		tfPlanApplication.IsFallbackPublicClient = types.BoolNull()
	}
	// END IsFallbackPublicClient | CreateBoolAttribute

	// START KeyCredentials | CreateArrayObjectAttribute
	if len(tfPlanApplication.KeyCredentials.Elements()) > 0 {
		var requestBodyKeyCredentials []models.KeyCredentialable
		for _, i := range tfPlanApplication.KeyCredentials.Elements() {
			requestBodyKeyCredential := models.NewKeyCredential()
			tfPlanKeyCredential := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyCredential)

			// START CustomKeyIdentifier | CreateStringBase64UrlAttribute
			if !tfPlanKeyCredential.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := tfPlanKeyCredential.CustomKeyIdentifier.ValueString()
				requestBodyKeyCredential.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				tfPlanKeyCredential.CustomKeyIdentifier = types.StringNull()
			}
			// END CustomKeyIdentifier | CreateStringBase64UrlAttribute

			// START DisplayName | CreateStringAttribute
			if !tfPlanKeyCredential.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanKeyCredential.DisplayName.ValueString()
				requestBodyKeyCredential.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanKeyCredential.DisplayName = types.StringNull()
			}
			// END DisplayName | CreateStringAttribute

			// START EndDateTime | CreateStringTimeAttribute
			if !tfPlanKeyCredential.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := tfPlanKeyCredential.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyKeyCredential.SetEndDateTime(&t)
			} else {
				tfPlanKeyCredential.EndDateTime = types.StringNull()
			}
			// END EndDateTime | CreateStringTimeAttribute

			// START Key | CreateStringBase64UrlAttribute
			if !tfPlanKeyCredential.Key.IsUnknown() {
				tfPlanKey := tfPlanKeyCredential.Key.ValueString()
				requestBodyKeyCredential.SetKey([]byte(tfPlanKey))
			} else {
				tfPlanKeyCredential.Key = types.StringNull()
			}
			// END Key | CreateStringBase64UrlAttribute

			// START KeyId | CreateStringUuidAttribute
			if !tfPlanKeyCredential.KeyId.IsUnknown() {
				tfPlanKeyId := tfPlanKeyCredential.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyKeyCredential.SetKeyId(&u)
			} else {
				tfPlanKeyCredential.KeyId = types.StringNull()
			}
			// END KeyId | CreateStringUuidAttribute

			// START StartDateTime | CreateStringTimeAttribute
			if !tfPlanKeyCredential.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := tfPlanKeyCredential.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyKeyCredential.SetStartDateTime(&t)
			} else {
				tfPlanKeyCredential.StartDateTime = types.StringNull()
			}
			// END StartDateTime | CreateStringTimeAttribute

			// START Type | CreateStringAttribute
			if !tfPlanKeyCredential.Type.IsUnknown() {
				tfPlanType := tfPlanKeyCredential.Type.ValueString()
				requestBodyKeyCredential.SetTypeEscaped(&tfPlanType)
			} else {
				tfPlanKeyCredential.Type = types.StringNull()
			}
			// END Type | CreateStringAttribute

			// START Usage | CreateStringAttribute
			if !tfPlanKeyCredential.Usage.IsUnknown() {
				tfPlanUsage := tfPlanKeyCredential.Usage.ValueString()
				requestBodyKeyCredential.SetUsage(&tfPlanUsage)
			} else {
				tfPlanKeyCredential.Usage = types.StringNull()
			}
			// END Usage | CreateStringAttribute

		}
		requestBodyApplication.SetKeyCredentials(requestBodyKeyCredentials)
	} else {
		tfPlanApplication.KeyCredentials = types.ListNull(tfPlanApplication.KeyCredentials.ElementType(ctx))
	}
	// END KeyCredentials | CreateArrayObjectAttribute

	// START Logo | CreateStringBase64UrlAttribute
	if !tfPlanApplication.Logo.IsUnknown() {
		tfPlanLogo := tfPlanApplication.Logo.ValueString()
		requestBodyApplication.SetLogo([]byte(tfPlanLogo))
	} else {
		tfPlanApplication.Logo = types.StringNull()
	}
	// END Logo | CreateStringBase64UrlAttribute

	// START NativeAuthenticationApisEnabled | CreateStringEnumAttribute
	if !tfPlanApplication.NativeAuthenticationApisEnabled.IsUnknown() {
		tfPlanNativeAuthenticationApisEnabled := tfPlanApplication.NativeAuthenticationApisEnabled.ValueString()
		parsedNativeAuthenticationApisEnabled, _ := models.ParseNativeAuthenticationApisEnabled(tfPlanNativeAuthenticationApisEnabled)
		assertedNativeAuthenticationApisEnabled := parsedNativeAuthenticationApisEnabled.(models.NativeAuthenticationApisEnabled)
		requestBodyApplication.SetNativeAuthenticationApisEnabled(&assertedNativeAuthenticationApisEnabled)
	} else {
		tfPlanApplication.NativeAuthenticationApisEnabled = types.StringNull()
	}
	// END NativeAuthenticationApisEnabled | CreateStringEnumAttribute

	// START Notes | CreateStringAttribute
	if !tfPlanApplication.Notes.IsUnknown() {
		tfPlanNotes := tfPlanApplication.Notes.ValueString()
		requestBodyApplication.SetNotes(&tfPlanNotes)
	} else {
		tfPlanApplication.Notes = types.StringNull()
	}
	// END Notes | CreateStringAttribute

	// START Oauth2RequirePostResponse | CreateBoolAttribute
	if !tfPlanApplication.Oauth2RequirePostResponse.IsUnknown() {
		tfPlanOauth2RequirePostResponse := tfPlanApplication.Oauth2RequirePostResponse.ValueBool()
		requestBodyApplication.SetOauth2RequirePostResponse(&tfPlanOauth2RequirePostResponse)
	} else {
		tfPlanApplication.Oauth2RequirePostResponse = types.BoolNull()
	}
	// END Oauth2RequirePostResponse | CreateBoolAttribute

	// START OptionalClaims | CreateObjectAttribute
	if !tfPlanApplication.OptionalClaims.IsUnknown() {
		requestBodyOptionalClaims := models.NewOptionalClaims()
		tfPlanOptionalClaims := applicationOptionalClaimsModel{}
		tfPlanApplication.OptionalClaims.As(ctx, &tfPlanOptionalClaims, basetypes.ObjectAsOptions{})

		// START AccessToken | CreateArrayObjectAttribute
		if len(tfPlanOptionalClaims.AccessToken.Elements()) > 0 {
			var requestBodyAccessToken []models.OptionalClaimable
			for _, i := range tfPlanOptionalClaims.AccessToken.Elements() {
				requestBodyOptionalClaim := models.NewOptionalClaim()
				tfPlanOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOptionalClaim)

				// START AdditionalProperties | CreateArrayStringAttribute
				if len(tfPlanOptionalClaim.AdditionalProperties.Elements()) > 0 {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanOptionalClaim.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyOptionalClaim.SetAdditionalProperties(stringArrayAdditionalProperties)
				} else {
					tfPlanOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				// END AdditionalProperties | CreateArrayStringAttribute

				// START Essential | CreateBoolAttribute
				if !tfPlanOptionalClaim.Essential.IsUnknown() {
					tfPlanEssential := tfPlanOptionalClaim.Essential.ValueBool()
					requestBodyOptionalClaim.SetEssential(&tfPlanEssential)
				} else {
					tfPlanOptionalClaim.Essential = types.BoolNull()
				}
				// END Essential | CreateBoolAttribute

				// START Name | CreateStringAttribute
				if !tfPlanOptionalClaim.Name.IsUnknown() {
					tfPlanName := tfPlanOptionalClaim.Name.ValueString()
					requestBodyOptionalClaim.SetName(&tfPlanName)
				} else {
					tfPlanOptionalClaim.Name = types.StringNull()
				}
				// END Name | CreateStringAttribute

				// START Source | CreateStringAttribute
				if !tfPlanOptionalClaim.Source.IsUnknown() {
					tfPlanSource := tfPlanOptionalClaim.Source.ValueString()
					requestBodyOptionalClaim.SetSource(&tfPlanSource)
				} else {
					tfPlanOptionalClaim.Source = types.StringNull()
				}
				// END Source | CreateStringAttribute

			}
			requestBodyOptionalClaims.SetAccessToken(requestBodyAccessToken)
		} else {
			tfPlanOptionalClaims.AccessToken = types.ListNull(tfPlanOptionalClaims.AccessToken.ElementType(ctx))
		}
		// END AccessToken | CreateArrayObjectAttribute

		// START IdToken | CreateArrayObjectAttribute
		if len(tfPlanOptionalClaims.IdToken.Elements()) > 0 {
			var requestBodyIdToken []models.OptionalClaimable
			for _, i := range tfPlanOptionalClaims.IdToken.Elements() {
				requestBodyOptionalClaim := models.NewOptionalClaim()
				tfPlanOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOptionalClaim)

				// START AdditionalProperties | CreateArrayStringAttribute
				if len(tfPlanOptionalClaim.AdditionalProperties.Elements()) > 0 {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanOptionalClaim.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyOptionalClaim.SetAdditionalProperties(stringArrayAdditionalProperties)
				} else {
					tfPlanOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				// END AdditionalProperties | CreateArrayStringAttribute

				// START Essential | CreateBoolAttribute
				if !tfPlanOptionalClaim.Essential.IsUnknown() {
					tfPlanEssential := tfPlanOptionalClaim.Essential.ValueBool()
					requestBodyOptionalClaim.SetEssential(&tfPlanEssential)
				} else {
					tfPlanOptionalClaim.Essential = types.BoolNull()
				}
				// END Essential | CreateBoolAttribute

				// START Name | CreateStringAttribute
				if !tfPlanOptionalClaim.Name.IsUnknown() {
					tfPlanName := tfPlanOptionalClaim.Name.ValueString()
					requestBodyOptionalClaim.SetName(&tfPlanName)
				} else {
					tfPlanOptionalClaim.Name = types.StringNull()
				}
				// END Name | CreateStringAttribute

				// START Source | CreateStringAttribute
				if !tfPlanOptionalClaim.Source.IsUnknown() {
					tfPlanSource := tfPlanOptionalClaim.Source.ValueString()
					requestBodyOptionalClaim.SetSource(&tfPlanSource)
				} else {
					tfPlanOptionalClaim.Source = types.StringNull()
				}
				// END Source | CreateStringAttribute

			}
			requestBodyOptionalClaims.SetIdToken(requestBodyIdToken)
		} else {
			tfPlanOptionalClaims.IdToken = types.ListNull(tfPlanOptionalClaims.IdToken.ElementType(ctx))
		}
		// END IdToken | CreateArrayObjectAttribute

		// START Saml2Token | CreateArrayObjectAttribute
		if len(tfPlanOptionalClaims.Saml2Token.Elements()) > 0 {
			var requestBodySaml2Token []models.OptionalClaimable
			for _, i := range tfPlanOptionalClaims.Saml2Token.Elements() {
				requestBodyOptionalClaim := models.NewOptionalClaim()
				tfPlanOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOptionalClaim)

				// START AdditionalProperties | CreateArrayStringAttribute
				if len(tfPlanOptionalClaim.AdditionalProperties.Elements()) > 0 {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanOptionalClaim.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyOptionalClaim.SetAdditionalProperties(stringArrayAdditionalProperties)
				} else {
					tfPlanOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				// END AdditionalProperties | CreateArrayStringAttribute

				// START Essential | CreateBoolAttribute
				if !tfPlanOptionalClaim.Essential.IsUnknown() {
					tfPlanEssential := tfPlanOptionalClaim.Essential.ValueBool()
					requestBodyOptionalClaim.SetEssential(&tfPlanEssential)
				} else {
					tfPlanOptionalClaim.Essential = types.BoolNull()
				}
				// END Essential | CreateBoolAttribute

				// START Name | CreateStringAttribute
				if !tfPlanOptionalClaim.Name.IsUnknown() {
					tfPlanName := tfPlanOptionalClaim.Name.ValueString()
					requestBodyOptionalClaim.SetName(&tfPlanName)
				} else {
					tfPlanOptionalClaim.Name = types.StringNull()
				}
				// END Name | CreateStringAttribute

				// START Source | CreateStringAttribute
				if !tfPlanOptionalClaim.Source.IsUnknown() {
					tfPlanSource := tfPlanOptionalClaim.Source.ValueString()
					requestBodyOptionalClaim.SetSource(&tfPlanSource)
				} else {
					tfPlanOptionalClaim.Source = types.StringNull()
				}
				// END Source | CreateStringAttribute

			}
			requestBodyOptionalClaims.SetSaml2Token(requestBodySaml2Token)
		} else {
			tfPlanOptionalClaims.Saml2Token = types.ListNull(tfPlanOptionalClaims.Saml2Token.ElementType(ctx))
		}
		// END Saml2Token | CreateArrayObjectAttribute

		requestBodyApplication.SetOptionalClaims(requestBodyOptionalClaims)
		tfPlanApplication.OptionalClaims, _ = types.ObjectValueFrom(ctx, tfPlanOptionalClaims.AttributeTypes(), requestBodyOptionalClaims)
	} else {
		tfPlanApplication.OptionalClaims = types.ObjectNull(tfPlanApplication.OptionalClaims.AttributeTypes(ctx))
	}
	// END OptionalClaims | CreateObjectAttribute

	// START ParentalControlSettings | CreateObjectAttribute
	if !tfPlanApplication.ParentalControlSettings.IsUnknown() {
		requestBodyParentalControlSettings := models.NewParentalControlSettings()
		tfPlanParentalControlSettings := applicationParentalControlSettingsModel{}
		tfPlanApplication.ParentalControlSettings.As(ctx, &tfPlanParentalControlSettings, basetypes.ObjectAsOptions{})

		// START CountriesBlockedForMinors | CreateArrayStringAttribute
		if len(tfPlanParentalControlSettings.CountriesBlockedForMinors.Elements()) > 0 {
			var stringArrayCountriesBlockedForMinors []string
			for _, i := range tfPlanParentalControlSettings.CountriesBlockedForMinors.Elements() {
				stringArrayCountriesBlockedForMinors = append(stringArrayCountriesBlockedForMinors, i.String())
			}
			requestBodyParentalControlSettings.SetCountriesBlockedForMinors(stringArrayCountriesBlockedForMinors)
		} else {
			tfPlanParentalControlSettings.CountriesBlockedForMinors = types.ListNull(types.StringType)
		}
		// END CountriesBlockedForMinors | CreateArrayStringAttribute

		// START LegalAgeGroupRule | CreateStringAttribute
		if !tfPlanParentalControlSettings.LegalAgeGroupRule.IsUnknown() {
			tfPlanLegalAgeGroupRule := tfPlanParentalControlSettings.LegalAgeGroupRule.ValueString()
			requestBodyParentalControlSettings.SetLegalAgeGroupRule(&tfPlanLegalAgeGroupRule)
		} else {
			tfPlanParentalControlSettings.LegalAgeGroupRule = types.StringNull()
		}
		// END LegalAgeGroupRule | CreateStringAttribute

		requestBodyApplication.SetParentalControlSettings(requestBodyParentalControlSettings)
		tfPlanApplication.ParentalControlSettings, _ = types.ObjectValueFrom(ctx, tfPlanParentalControlSettings.AttributeTypes(), requestBodyParentalControlSettings)
	} else {
		tfPlanApplication.ParentalControlSettings = types.ObjectNull(tfPlanApplication.ParentalControlSettings.AttributeTypes(ctx))
	}
	// END ParentalControlSettings | CreateObjectAttribute

	// START PasswordCredentials | CreateArrayObjectAttribute
	if len(tfPlanApplication.PasswordCredentials.Elements()) > 0 {
		var requestBodyPasswordCredentials []models.PasswordCredentialable
		for _, i := range tfPlanApplication.PasswordCredentials.Elements() {
			requestBodyPasswordCredential := models.NewPasswordCredential()
			tfPlanPasswordCredential := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPasswordCredential)

			// START CustomKeyIdentifier | CreateStringBase64UrlAttribute
			if !tfPlanPasswordCredential.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := tfPlanPasswordCredential.CustomKeyIdentifier.ValueString()
				requestBodyPasswordCredential.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				tfPlanPasswordCredential.CustomKeyIdentifier = types.StringNull()
			}
			// END CustomKeyIdentifier | CreateStringBase64UrlAttribute

			// START DisplayName | CreateStringAttribute
			if !tfPlanPasswordCredential.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanPasswordCredential.DisplayName.ValueString()
				requestBodyPasswordCredential.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanPasswordCredential.DisplayName = types.StringNull()
			}
			// END DisplayName | CreateStringAttribute

			// START EndDateTime | CreateStringTimeAttribute
			if !tfPlanPasswordCredential.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := tfPlanPasswordCredential.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyPasswordCredential.SetEndDateTime(&t)
			} else {
				tfPlanPasswordCredential.EndDateTime = types.StringNull()
			}
			// END EndDateTime | CreateStringTimeAttribute

			// START Hint | CreateStringAttribute
			if !tfPlanPasswordCredential.Hint.IsUnknown() {
				tfPlanHint := tfPlanPasswordCredential.Hint.ValueString()
				requestBodyPasswordCredential.SetHint(&tfPlanHint)
			} else {
				tfPlanPasswordCredential.Hint = types.StringNull()
			}
			// END Hint | CreateStringAttribute

			// START KeyId | CreateStringUuidAttribute
			if !tfPlanPasswordCredential.KeyId.IsUnknown() {
				tfPlanKeyId := tfPlanPasswordCredential.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyPasswordCredential.SetKeyId(&u)
			} else {
				tfPlanPasswordCredential.KeyId = types.StringNull()
			}
			// END KeyId | CreateStringUuidAttribute

			// START SecretText | CreateStringAttribute
			if !tfPlanPasswordCredential.SecretText.IsUnknown() {
				tfPlanSecretText := tfPlanPasswordCredential.SecretText.ValueString()
				requestBodyPasswordCredential.SetSecretText(&tfPlanSecretText)
			} else {
				tfPlanPasswordCredential.SecretText = types.StringNull()
			}
			// END SecretText | CreateStringAttribute

			// START StartDateTime | CreateStringTimeAttribute
			if !tfPlanPasswordCredential.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := tfPlanPasswordCredential.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyPasswordCredential.SetStartDateTime(&t)
			} else {
				tfPlanPasswordCredential.StartDateTime = types.StringNull()
			}
			// END StartDateTime | CreateStringTimeAttribute

		}
		requestBodyApplication.SetPasswordCredentials(requestBodyPasswordCredentials)
	} else {
		tfPlanApplication.PasswordCredentials = types.ListNull(tfPlanApplication.PasswordCredentials.ElementType(ctx))
	}
	// END PasswordCredentials | CreateArrayObjectAttribute

	// START PublicClient | CreateObjectAttribute
	if !tfPlanApplication.PublicClient.IsUnknown() {
		requestBodyPublicClientApplication := models.NewPublicClientApplication()
		tfPlanPublicClientApplication := applicationPublicClientApplicationModel{}
		tfPlanApplication.PublicClient.As(ctx, &tfPlanPublicClientApplication, basetypes.ObjectAsOptions{})

		// START RedirectUris | CreateArrayStringAttribute
		if len(tfPlanPublicClientApplication.RedirectUris.Elements()) > 0 {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanPublicClientApplication.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodyPublicClientApplication.SetRedirectUris(stringArrayRedirectUris)
		} else {
			tfPlanPublicClientApplication.RedirectUris = types.ListNull(types.StringType)
		}
		// END RedirectUris | CreateArrayStringAttribute

		requestBodyApplication.SetPublicClient(requestBodyPublicClientApplication)
		tfPlanApplication.PublicClient, _ = types.ObjectValueFrom(ctx, tfPlanPublicClientApplication.AttributeTypes(), requestBodyPublicClientApplication)
	} else {
		tfPlanApplication.PublicClient = types.ObjectNull(tfPlanApplication.PublicClient.AttributeTypes(ctx))
	}
	// END PublicClient | CreateObjectAttribute

	// START PublisherDomain | CreateStringAttribute
	if !tfPlanApplication.PublisherDomain.IsUnknown() {
		tfPlanPublisherDomain := tfPlanApplication.PublisherDomain.ValueString()
		requestBodyApplication.SetPublisherDomain(&tfPlanPublisherDomain)
	} else {
		tfPlanApplication.PublisherDomain = types.StringNull()
	}
	// END PublisherDomain | CreateStringAttribute

	// START RequestSignatureVerification | CreateObjectAttribute
	if !tfPlanApplication.RequestSignatureVerification.IsUnknown() {
		requestBodyRequestSignatureVerification := models.NewRequestSignatureVerification()
		tfPlanRequestSignatureVerification := applicationRequestSignatureVerificationModel{}
		tfPlanApplication.RequestSignatureVerification.As(ctx, &tfPlanRequestSignatureVerification, basetypes.ObjectAsOptions{})

		// START AllowedWeakAlgorithms | CreateStringEnumAttribute
		if !tfPlanRequestSignatureVerification.AllowedWeakAlgorithms.IsUnknown() {
			tfPlanAllowedWeakAlgorithms := tfPlanRequestSignatureVerification.AllowedWeakAlgorithms.ValueString()
			parsedAllowedWeakAlgorithms, _ := models.ParseWeakAlgorithms(tfPlanAllowedWeakAlgorithms)
			assertedAllowedWeakAlgorithms := parsedAllowedWeakAlgorithms.(models.WeakAlgorithms)
			requestBodyRequestSignatureVerification.SetAllowedWeakAlgorithms(&assertedAllowedWeakAlgorithms)
		} else {
			tfPlanRequestSignatureVerification.AllowedWeakAlgorithms = types.StringNull()
		}
		// END AllowedWeakAlgorithms | CreateStringEnumAttribute

		// START IsSignedRequestRequired | CreateBoolAttribute
		if !tfPlanRequestSignatureVerification.IsSignedRequestRequired.IsUnknown() {
			tfPlanIsSignedRequestRequired := tfPlanRequestSignatureVerification.IsSignedRequestRequired.ValueBool()
			requestBodyRequestSignatureVerification.SetIsSignedRequestRequired(&tfPlanIsSignedRequestRequired)
		} else {
			tfPlanRequestSignatureVerification.IsSignedRequestRequired = types.BoolNull()
		}
		// END IsSignedRequestRequired | CreateBoolAttribute

		requestBodyApplication.SetRequestSignatureVerification(requestBodyRequestSignatureVerification)
		tfPlanApplication.RequestSignatureVerification, _ = types.ObjectValueFrom(ctx, tfPlanRequestSignatureVerification.AttributeTypes(), requestBodyRequestSignatureVerification)
	} else {
		tfPlanApplication.RequestSignatureVerification = types.ObjectNull(tfPlanApplication.RequestSignatureVerification.AttributeTypes(ctx))
	}
	// END RequestSignatureVerification | CreateObjectAttribute

	// START RequiredResourceAccess | CreateArrayObjectAttribute
	if len(tfPlanApplication.RequiredResourceAccess.Elements()) > 0 {
		var requestBodyRequiredResourceAccess []models.RequiredResourceAccessable
		for _, i := range tfPlanApplication.RequiredResourceAccess.Elements() {
			requestBodyRequiredResourceAccess := models.NewRequiredResourceAccess()
			tfPlanRequiredResourceAccess := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanRequiredResourceAccess)

			// START ResourceAccess | CreateArrayObjectAttribute
			if len(tfPlanRequiredResourceAccess.ResourceAccess.Elements()) > 0 {
				var requestBodyResourceAccess []models.ResourceAccessable
				for _, i := range tfPlanRequiredResourceAccess.ResourceAccess.Elements() {
					requestBodyResourceAccess := models.NewResourceAccess()
					tfPlanResourceAccess := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanResourceAccess)

					// START Id | CreateStringUuidAttribute
					if !tfPlanResourceAccess.Id.IsUnknown() {
						tfPlanId := tfPlanResourceAccess.Id.ValueString()
						u, _ := uuid.Parse(tfPlanId)
						requestBodyResourceAccess.SetId(&u)
					} else {
						tfPlanResourceAccess.Id = types.StringNull()
					}
					// END Id | CreateStringUuidAttribute

					// START Type | CreateStringAttribute
					if !tfPlanResourceAccess.Type.IsUnknown() {
						tfPlanType := tfPlanResourceAccess.Type.ValueString()
						requestBodyResourceAccess.SetTypeEscaped(&tfPlanType)
					} else {
						tfPlanResourceAccess.Type = types.StringNull()
					}
					// END Type | CreateStringAttribute

				}
				requestBodyRequiredResourceAccess.SetResourceAccess(requestBodyResourceAccess)
			} else {
				tfPlanRequiredResourceAccess.ResourceAccess = types.ListNull(tfPlanRequiredResourceAccess.ResourceAccess.ElementType(ctx))
			}
			// END ResourceAccess | CreateArrayObjectAttribute

			// START ResourceAppId | CreateStringAttribute
			if !tfPlanRequiredResourceAccess.ResourceAppId.IsUnknown() {
				tfPlanResourceAppId := tfPlanRequiredResourceAccess.ResourceAppId.ValueString()
				requestBodyRequiredResourceAccess.SetResourceAppId(&tfPlanResourceAppId)
			} else {
				tfPlanRequiredResourceAccess.ResourceAppId = types.StringNull()
			}
			// END ResourceAppId | CreateStringAttribute

		}
		requestBodyApplication.SetRequiredResourceAccess(requestBodyRequiredResourceAccess)
	} else {
		tfPlanApplication.RequiredResourceAccess = types.ListNull(tfPlanApplication.RequiredResourceAccess.ElementType(ctx))
	}
	// END RequiredResourceAccess | CreateArrayObjectAttribute

	// START SamlMetadataUrl | CreateStringAttribute
	if !tfPlanApplication.SamlMetadataUrl.IsUnknown() {
		tfPlanSamlMetadataUrl := tfPlanApplication.SamlMetadataUrl.ValueString()
		requestBodyApplication.SetSamlMetadataUrl(&tfPlanSamlMetadataUrl)
	} else {
		tfPlanApplication.SamlMetadataUrl = types.StringNull()
	}
	// END SamlMetadataUrl | CreateStringAttribute

	// START ServiceManagementReference | CreateStringAttribute
	if !tfPlanApplication.ServiceManagementReference.IsUnknown() {
		tfPlanServiceManagementReference := tfPlanApplication.ServiceManagementReference.ValueString()
		requestBodyApplication.SetServiceManagementReference(&tfPlanServiceManagementReference)
	} else {
		tfPlanApplication.ServiceManagementReference = types.StringNull()
	}
	// END ServiceManagementReference | CreateStringAttribute

	// START ServicePrincipalLockConfiguration | CreateObjectAttribute
	if !tfPlanApplication.ServicePrincipalLockConfiguration.IsUnknown() {
		requestBodyServicePrincipalLockConfiguration := models.NewServicePrincipalLockConfiguration()
		tfPlanServicePrincipalLockConfiguration := applicationServicePrincipalLockConfigurationModel{}
		tfPlanApplication.ServicePrincipalLockConfiguration.As(ctx, &tfPlanServicePrincipalLockConfiguration, basetypes.ObjectAsOptions{})

		// START AllProperties | CreateBoolAttribute
		if !tfPlanServicePrincipalLockConfiguration.AllProperties.IsUnknown() {
			tfPlanAllProperties := tfPlanServicePrincipalLockConfiguration.AllProperties.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetAllProperties(&tfPlanAllProperties)
		} else {
			tfPlanServicePrincipalLockConfiguration.AllProperties = types.BoolNull()
		}
		// END AllProperties | CreateBoolAttribute

		// START CredentialsWithUsageSign | CreateBoolAttribute
		if !tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageSign.IsUnknown() {
			tfPlanCredentialsWithUsageSign := tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageSign.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetCredentialsWithUsageSign(&tfPlanCredentialsWithUsageSign)
		} else {
			tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolNull()
		}
		// END CredentialsWithUsageSign | CreateBoolAttribute

		// START CredentialsWithUsageVerify | CreateBoolAttribute
		if !tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageVerify.IsUnknown() {
			tfPlanCredentialsWithUsageVerify := tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageVerify.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetCredentialsWithUsageVerify(&tfPlanCredentialsWithUsageVerify)
		} else {
			tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolNull()
		}
		// END CredentialsWithUsageVerify | CreateBoolAttribute

		// START IsEnabled | CreateBoolAttribute
		if !tfPlanServicePrincipalLockConfiguration.IsEnabled.IsUnknown() {
			tfPlanIsEnabled := tfPlanServicePrincipalLockConfiguration.IsEnabled.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetIsEnabled(&tfPlanIsEnabled)
		} else {
			tfPlanServicePrincipalLockConfiguration.IsEnabled = types.BoolNull()
		}
		// END IsEnabled | CreateBoolAttribute

		// START TokenEncryptionKeyId | CreateBoolAttribute
		if !tfPlanServicePrincipalLockConfiguration.TokenEncryptionKeyId.IsUnknown() {
			tfPlanTokenEncryptionKeyId := tfPlanServicePrincipalLockConfiguration.TokenEncryptionKeyId.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetTokenEncryptionKeyId(&tfPlanTokenEncryptionKeyId)
		} else {
			tfPlanServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolNull()
		}
		// END TokenEncryptionKeyId | CreateBoolAttribute

		requestBodyApplication.SetServicePrincipalLockConfiguration(requestBodyServicePrincipalLockConfiguration)
		tfPlanApplication.ServicePrincipalLockConfiguration, _ = types.ObjectValueFrom(ctx, tfPlanServicePrincipalLockConfiguration.AttributeTypes(), requestBodyServicePrincipalLockConfiguration)
	} else {
		tfPlanApplication.ServicePrincipalLockConfiguration = types.ObjectNull(tfPlanApplication.ServicePrincipalLockConfiguration.AttributeTypes(ctx))
	}
	// END ServicePrincipalLockConfiguration | CreateObjectAttribute

	// START SignInAudience | CreateStringAttribute
	if !tfPlanApplication.SignInAudience.IsUnknown() {
		tfPlanSignInAudience := tfPlanApplication.SignInAudience.ValueString()
		requestBodyApplication.SetSignInAudience(&tfPlanSignInAudience)
	} else {
		tfPlanApplication.SignInAudience = types.StringNull()
	}
	// END SignInAudience | CreateStringAttribute

	// START Spa | CreateObjectAttribute
	if !tfPlanApplication.Spa.IsUnknown() {
		requestBodySpaApplication := models.NewSpaApplication()
		tfPlanSpaApplication := applicationSpaApplicationModel{}
		tfPlanApplication.Spa.As(ctx, &tfPlanSpaApplication, basetypes.ObjectAsOptions{})

		// START RedirectUris | CreateArrayStringAttribute
		if len(tfPlanSpaApplication.RedirectUris.Elements()) > 0 {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanSpaApplication.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodySpaApplication.SetRedirectUris(stringArrayRedirectUris)
		} else {
			tfPlanSpaApplication.RedirectUris = types.ListNull(types.StringType)
		}
		// END RedirectUris | CreateArrayStringAttribute

		requestBodyApplication.SetSpa(requestBodySpaApplication)
		tfPlanApplication.Spa, _ = types.ObjectValueFrom(ctx, tfPlanSpaApplication.AttributeTypes(), requestBodySpaApplication)
	} else {
		tfPlanApplication.Spa = types.ObjectNull(tfPlanApplication.Spa.AttributeTypes(ctx))
	}
	// END Spa | CreateObjectAttribute

	// START Tags | CreateArrayStringAttribute
	if len(tfPlanApplication.Tags.Elements()) > 0 {
		var stringArrayTags []string
		for _, i := range tfPlanApplication.Tags.Elements() {
			stringArrayTags = append(stringArrayTags, i.String())
		}
		requestBodyApplication.SetTags(stringArrayTags)
	} else {
		tfPlanApplication.Tags = types.ListNull(types.StringType)
	}
	// END Tags | CreateArrayStringAttribute

	// START TokenEncryptionKeyId | CreateStringUuidAttribute
	if !tfPlanApplication.TokenEncryptionKeyId.IsUnknown() {
		tfPlanTokenEncryptionKeyId := tfPlanApplication.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		requestBodyApplication.SetTokenEncryptionKeyId(&u)
	} else {
		tfPlanApplication.TokenEncryptionKeyId = types.StringNull()
	}
	// END TokenEncryptionKeyId | CreateStringUuidAttribute

	// START UniqueName | CreateStringAttribute
	if !tfPlanApplication.UniqueName.IsUnknown() {
		tfPlanUniqueName := tfPlanApplication.UniqueName.ValueString()
		requestBodyApplication.SetUniqueName(&tfPlanUniqueName)
	} else {
		tfPlanApplication.UniqueName = types.StringNull()
	}
	// END UniqueName | CreateStringAttribute

	// START VerifiedPublisher | CreateObjectAttribute
	if !tfPlanApplication.VerifiedPublisher.IsUnknown() {
		requestBodyVerifiedPublisher := models.NewVerifiedPublisher()
		tfPlanVerifiedPublisher := applicationVerifiedPublisherModel{}
		tfPlanApplication.VerifiedPublisher.As(ctx, &tfPlanVerifiedPublisher, basetypes.ObjectAsOptions{})

		// START AddedDateTime | CreateStringTimeAttribute
		if !tfPlanVerifiedPublisher.AddedDateTime.IsUnknown() {
			tfPlanAddedDateTime := tfPlanVerifiedPublisher.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanAddedDateTime)
			requestBodyVerifiedPublisher.SetAddedDateTime(&t)
		} else {
			tfPlanVerifiedPublisher.AddedDateTime = types.StringNull()
		}
		// END AddedDateTime | CreateStringTimeAttribute

		// START DisplayName | CreateStringAttribute
		if !tfPlanVerifiedPublisher.DisplayName.IsUnknown() {
			tfPlanDisplayName := tfPlanVerifiedPublisher.DisplayName.ValueString()
			requestBodyVerifiedPublisher.SetDisplayName(&tfPlanDisplayName)
		} else {
			tfPlanVerifiedPublisher.DisplayName = types.StringNull()
		}
		// END DisplayName | CreateStringAttribute

		// START VerifiedPublisherId | CreateStringAttribute
		if !tfPlanVerifiedPublisher.VerifiedPublisherId.IsUnknown() {
			tfPlanVerifiedPublisherId := tfPlanVerifiedPublisher.VerifiedPublisherId.ValueString()
			requestBodyVerifiedPublisher.SetVerifiedPublisherId(&tfPlanVerifiedPublisherId)
		} else {
			tfPlanVerifiedPublisher.VerifiedPublisherId = types.StringNull()
		}
		// END VerifiedPublisherId | CreateStringAttribute

		requestBodyApplication.SetVerifiedPublisher(requestBodyVerifiedPublisher)
		tfPlanApplication.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfPlanVerifiedPublisher.AttributeTypes(), requestBodyVerifiedPublisher)
	} else {
		tfPlanApplication.VerifiedPublisher = types.ObjectNull(tfPlanApplication.VerifiedPublisher.AttributeTypes(ctx))
	}
	// END VerifiedPublisher | CreateObjectAttribute

	// START Web | CreateObjectAttribute
	if !tfPlanApplication.Web.IsUnknown() {
		requestBodyWebApplication := models.NewWebApplication()
		tfPlanWebApplication := applicationWebApplicationModel{}
		tfPlanApplication.Web.As(ctx, &tfPlanWebApplication, basetypes.ObjectAsOptions{})

		// START HomePageUrl | CreateStringAttribute
		if !tfPlanWebApplication.HomePageUrl.IsUnknown() {
			tfPlanHomePageUrl := tfPlanWebApplication.HomePageUrl.ValueString()
			requestBodyWebApplication.SetHomePageUrl(&tfPlanHomePageUrl)
		} else {
			tfPlanWebApplication.HomePageUrl = types.StringNull()
		}
		// END HomePageUrl | CreateStringAttribute

		// START ImplicitGrantSettings | CreateObjectAttribute
		if !tfPlanWebApplication.ImplicitGrantSettings.IsUnknown() {
			requestBodyImplicitGrantSettings := models.NewImplicitGrantSettings()
			tfPlanImplicitGrantSettings := applicationImplicitGrantSettingsModel{}
			tfPlanWebApplication.ImplicitGrantSettings.As(ctx, &tfPlanImplicitGrantSettings, basetypes.ObjectAsOptions{})

			// START EnableAccessTokenIssuance | CreateBoolAttribute
			if !tfPlanImplicitGrantSettings.EnableAccessTokenIssuance.IsUnknown() {
				tfPlanEnableAccessTokenIssuance := tfPlanImplicitGrantSettings.EnableAccessTokenIssuance.ValueBool()
				requestBodyImplicitGrantSettings.SetEnableAccessTokenIssuance(&tfPlanEnableAccessTokenIssuance)
			} else {
				tfPlanImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolNull()
			}
			// END EnableAccessTokenIssuance | CreateBoolAttribute

			// START EnableIdTokenIssuance | CreateBoolAttribute
			if !tfPlanImplicitGrantSettings.EnableIdTokenIssuance.IsUnknown() {
				tfPlanEnableIdTokenIssuance := tfPlanImplicitGrantSettings.EnableIdTokenIssuance.ValueBool()
				requestBodyImplicitGrantSettings.SetEnableIdTokenIssuance(&tfPlanEnableIdTokenIssuance)
			} else {
				tfPlanImplicitGrantSettings.EnableIdTokenIssuance = types.BoolNull()
			}
			// END EnableIdTokenIssuance | CreateBoolAttribute

			requestBodyWebApplication.SetImplicitGrantSettings(requestBodyImplicitGrantSettings)
			tfPlanWebApplication.ImplicitGrantSettings, _ = types.ObjectValueFrom(ctx, tfPlanImplicitGrantSettings.AttributeTypes(), requestBodyImplicitGrantSettings)
		} else {
			tfPlanWebApplication.ImplicitGrantSettings = types.ObjectNull(tfPlanWebApplication.ImplicitGrantSettings.AttributeTypes(ctx))
		}
		// END ImplicitGrantSettings | CreateObjectAttribute

		// START LogoutUrl | CreateStringAttribute
		if !tfPlanWebApplication.LogoutUrl.IsUnknown() {
			tfPlanLogoutUrl := tfPlanWebApplication.LogoutUrl.ValueString()
			requestBodyWebApplication.SetLogoutUrl(&tfPlanLogoutUrl)
		} else {
			tfPlanWebApplication.LogoutUrl = types.StringNull()
		}
		// END LogoutUrl | CreateStringAttribute

		// START RedirectUriSettings | CreateArrayObjectAttribute
		if len(tfPlanWebApplication.RedirectUriSettings.Elements()) > 0 {
			var requestBodyRedirectUriSettings []models.RedirectUriSettingsable
			for _, i := range tfPlanWebApplication.RedirectUriSettings.Elements() {
				requestBodyRedirectUriSettings := models.NewRedirectUriSettings()
				tfPlanRedirectUriSettings := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanRedirectUriSettings)

				// START Index | UNKNOWN
				// END Index | UNKNOWN

				// START Uri | CreateStringAttribute
				if !tfPlanRedirectUriSettings.Uri.IsUnknown() {
					tfPlanUri := tfPlanRedirectUriSettings.Uri.ValueString()
					requestBodyRedirectUriSettings.SetUri(&tfPlanUri)
				} else {
					tfPlanRedirectUriSettings.Uri = types.StringNull()
				}
				// END Uri | CreateStringAttribute

			}
			requestBodyWebApplication.SetRedirectUriSettings(requestBodyRedirectUriSettings)
		} else {
			tfPlanWebApplication.RedirectUriSettings = types.ListNull(tfPlanWebApplication.RedirectUriSettings.ElementType(ctx))
		}
		// END RedirectUriSettings | CreateArrayObjectAttribute

		// START RedirectUris | CreateArrayStringAttribute
		if len(tfPlanWebApplication.RedirectUris.Elements()) > 0 {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanWebApplication.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodyWebApplication.SetRedirectUris(stringArrayRedirectUris)
		} else {
			tfPlanWebApplication.RedirectUris = types.ListNull(types.StringType)
		}
		// END RedirectUris | CreateArrayStringAttribute

		requestBodyApplication.SetWeb(requestBodyWebApplication)
		tfPlanApplication.Web, _ = types.ObjectValueFrom(ctx, tfPlanWebApplication.AttributeTypes(), requestBodyWebApplication)
	} else {
		tfPlanApplication.Web = types.ObjectNull(tfPlanApplication.Web.AttributeTypes(ctx))
	}
	// END Web | CreateObjectAttribute

	// Create new application
	result, err := r.client.Applications().Post(context.Background(), requestBodyApplication, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlanApplication.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlanApplication)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *applicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state applicationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	var result models.Applicationable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.Applications().ByApplicationId(state.Id.ValueString()).Get(context.Background(), &qparams)
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

	if result.GetId() != nil {
		state.Id = types.StringValue(*result.GetId())
	} else {
		state.Id = types.StringNull()
	}
	if result.GetDeletedDateTime() != nil {
		state.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	} else {
		state.DeletedDateTime = types.StringNull()
	}
	if len(result.GetAddIns()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAddIns() {
			addIns := new(applicationAddInModel)

			if v.GetId() != nil {
				addIns.Id = types.StringValue(v.GetId().String())
			} else {
				addIns.Id = types.StringNull()
			}
			if len(v.GetProperties()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetProperties() {
					properties := new(applicationKeyValueModel)

					if v.GetKey() != nil {
						properties.Key = types.StringValue(*v.GetKey())
					} else {
						properties.Key = types.StringNull()
					}
					if v.GetValue() != nil {
						properties.Value = types.StringValue(*v.GetValue())
					} else {
						properties.Value = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, properties.AttributeTypes(), properties)
					objectValues = append(objectValues, objectValue)
				}
				addIns.Properties, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetTypeEscaped() != nil {
				addIns.Type = types.StringValue(*v.GetTypeEscaped())
			} else {
				addIns.Type = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, addIns.AttributeTypes(), addIns)
			objectValues = append(objectValues, objectValue)
		}
		state.AddIns, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetApi() != nil {
		api := new(applicationApiApplicationModel)

		if result.GetApi().GetAcceptMappedClaims() != nil {
			api.AcceptMappedClaims = types.BoolValue(*result.GetApi().GetAcceptMappedClaims())
		} else {
			api.AcceptMappedClaims = types.BoolNull()
		}
		if len(result.GetApi().GetKnownClientApplications()) > 0 {
			var knownClientApplications []attr.Value
			for _, v := range result.GetApi().GetKnownClientApplications() {
				knownClientApplications = append(knownClientApplications, types.StringValue(v.String()))
			}
			listValue, _ := types.ListValue(types.StringType, knownClientApplications)
			api.KnownClientApplications = listValue
		} else {
			api.KnownClientApplications = types.ListNull(types.StringType)
		}
		if len(result.GetApi().GetOauth2PermissionScopes()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetApi().GetOauth2PermissionScopes() {
				oauth2PermissionScopes := new(applicationPermissionScopeModel)

				if v.GetAdminConsentDescription() != nil {
					oauth2PermissionScopes.AdminConsentDescription = types.StringValue(*v.GetAdminConsentDescription())
				} else {
					oauth2PermissionScopes.AdminConsentDescription = types.StringNull()
				}
				if v.GetAdminConsentDisplayName() != nil {
					oauth2PermissionScopes.AdminConsentDisplayName = types.StringValue(*v.GetAdminConsentDisplayName())
				} else {
					oauth2PermissionScopes.AdminConsentDisplayName = types.StringNull()
				}
				if v.GetId() != nil {
					oauth2PermissionScopes.Id = types.StringValue(v.GetId().String())
				} else {
					oauth2PermissionScopes.Id = types.StringNull()
				}
				if v.GetIsEnabled() != nil {
					oauth2PermissionScopes.IsEnabled = types.BoolValue(*v.GetIsEnabled())
				} else {
					oauth2PermissionScopes.IsEnabled = types.BoolNull()
				}
				if v.GetOrigin() != nil {
					oauth2PermissionScopes.Origin = types.StringValue(*v.GetOrigin())
				} else {
					oauth2PermissionScopes.Origin = types.StringNull()
				}
				if v.GetTypeEscaped() != nil {
					oauth2PermissionScopes.Type = types.StringValue(*v.GetTypeEscaped())
				} else {
					oauth2PermissionScopes.Type = types.StringNull()
				}
				if v.GetUserConsentDescription() != nil {
					oauth2PermissionScopes.UserConsentDescription = types.StringValue(*v.GetUserConsentDescription())
				} else {
					oauth2PermissionScopes.UserConsentDescription = types.StringNull()
				}
				if v.GetUserConsentDisplayName() != nil {
					oauth2PermissionScopes.UserConsentDisplayName = types.StringValue(*v.GetUserConsentDisplayName())
				} else {
					oauth2PermissionScopes.UserConsentDisplayName = types.StringNull()
				}
				if v.GetValue() != nil {
					oauth2PermissionScopes.Value = types.StringValue(*v.GetValue())
				} else {
					oauth2PermissionScopes.Value = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, oauth2PermissionScopes.AttributeTypes(), oauth2PermissionScopes)
				objectValues = append(objectValues, objectValue)
			}
			api.Oauth2PermissionScopes, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetApi().GetPreAuthorizedApplications()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetApi().GetPreAuthorizedApplications() {
				preAuthorizedApplications := new(applicationPreAuthorizedApplicationModel)

				if v.GetAppId() != nil {
					preAuthorizedApplications.AppId = types.StringValue(*v.GetAppId())
				} else {
					preAuthorizedApplications.AppId = types.StringNull()
				}
				if len(v.GetDelegatedPermissionIds()) > 0 {
					var delegatedPermissionIds []attr.Value
					for _, v := range v.GetDelegatedPermissionIds() {
						delegatedPermissionIds = append(delegatedPermissionIds, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, delegatedPermissionIds)
					preAuthorizedApplications.DelegatedPermissionIds = listValue
				} else {
					preAuthorizedApplications.DelegatedPermissionIds = types.ListNull(types.StringType)
				}
				objectValue, _ := types.ObjectValueFrom(ctx, preAuthorizedApplications.AttributeTypes(), preAuthorizedApplications)
				objectValues = append(objectValues, objectValue)
			}
			api.PreAuthorizedApplications, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}

		objectValue, _ := types.ObjectValueFrom(ctx, api.AttributeTypes(), api)
		state.Api = objectValue
	}
	if result.GetAppId() != nil {
		state.AppId = types.StringValue(*result.GetAppId())
	} else {
		state.AppId = types.StringNull()
	}
	if len(result.GetAppRoles()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAppRoles() {
			appRoles := new(applicationAppRoleModel)

			if len(v.GetAllowedMemberTypes()) > 0 {
				var allowedMemberTypes []attr.Value
				for _, v := range v.GetAllowedMemberTypes() {
					allowedMemberTypes = append(allowedMemberTypes, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, allowedMemberTypes)
				appRoles.AllowedMemberTypes = listValue
			} else {
				appRoles.AllowedMemberTypes = types.ListNull(types.StringType)
			}
			if v.GetDescription() != nil {
				appRoles.Description = types.StringValue(*v.GetDescription())
			} else {
				appRoles.Description = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				appRoles.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				appRoles.DisplayName = types.StringNull()
			}
			if v.GetId() != nil {
				appRoles.Id = types.StringValue(v.GetId().String())
			} else {
				appRoles.Id = types.StringNull()
			}
			if v.GetIsEnabled() != nil {
				appRoles.IsEnabled = types.BoolValue(*v.GetIsEnabled())
			} else {
				appRoles.IsEnabled = types.BoolNull()
			}
			if v.GetOrigin() != nil {
				appRoles.Origin = types.StringValue(*v.GetOrigin())
			} else {
				appRoles.Origin = types.StringNull()
			}
			if v.GetValue() != nil {
				appRoles.Value = types.StringValue(*v.GetValue())
			} else {
				appRoles.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, appRoles.AttributeTypes(), appRoles)
			objectValues = append(objectValues, objectValue)
		}
		state.AppRoles, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetApplicationTemplateId() != nil {
		state.ApplicationTemplateId = types.StringValue(*result.GetApplicationTemplateId())
	} else {
		state.ApplicationTemplateId = types.StringNull()
	}
	if result.GetCertification() != nil {
		certification := new(applicationCertificationModel)

		if result.GetCertification().GetCertificationDetailsUrl() != nil {
			certification.CertificationDetailsUrl = types.StringValue(*result.GetCertification().GetCertificationDetailsUrl())
		} else {
			certification.CertificationDetailsUrl = types.StringNull()
		}
		if result.GetCertification().GetCertificationExpirationDateTime() != nil {
			certification.CertificationExpirationDateTime = types.StringValue(result.GetCertification().GetCertificationExpirationDateTime().String())
		} else {
			certification.CertificationExpirationDateTime = types.StringNull()
		}
		if result.GetCertification().GetIsCertifiedByMicrosoft() != nil {
			certification.IsCertifiedByMicrosoft = types.BoolValue(*result.GetCertification().GetIsCertifiedByMicrosoft())
		} else {
			certification.IsCertifiedByMicrosoft = types.BoolNull()
		}
		if result.GetCertification().GetIsPublisherAttested() != nil {
			certification.IsPublisherAttested = types.BoolValue(*result.GetCertification().GetIsPublisherAttested())
		} else {
			certification.IsPublisherAttested = types.BoolNull()
		}
		if result.GetCertification().GetLastCertificationDateTime() != nil {
			certification.LastCertificationDateTime = types.StringValue(result.GetCertification().GetLastCertificationDateTime().String())
		} else {
			certification.LastCertificationDateTime = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, certification.AttributeTypes(), certification)
		state.Certification = objectValue
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		state.CreatedDateTime = types.StringNull()
	}
	if result.GetDefaultRedirectUri() != nil {
		state.DefaultRedirectUri = types.StringValue(*result.GetDefaultRedirectUri())
	} else {
		state.DefaultRedirectUri = types.StringNull()
	}
	if result.GetDescription() != nil {
		state.Description = types.StringValue(*result.GetDescription())
	} else {
		state.Description = types.StringNull()
	}
	if result.GetDisabledByMicrosoftStatus() != nil {
		state.DisabledByMicrosoftStatus = types.StringValue(*result.GetDisabledByMicrosoftStatus())
	} else {
		state.DisabledByMicrosoftStatus = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		state.DisplayName = types.StringNull()
	}
	if result.GetGroupMembershipClaims() != nil {
		state.GroupMembershipClaims = types.StringValue(*result.GetGroupMembershipClaims())
	} else {
		state.GroupMembershipClaims = types.StringNull()
	}
	if len(result.GetIdentifierUris()) > 0 {
		var identifierUris []attr.Value
		for _, v := range result.GetIdentifierUris() {
			identifierUris = append(identifierUris, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, identifierUris)
		state.IdentifierUris = listValue
	} else {
		state.IdentifierUris = types.ListNull(types.StringType)
	}
	if result.GetInfo() != nil {
		info := new(applicationInformationalUrlModel)

		if result.GetInfo().GetLogoUrl() != nil {
			info.LogoUrl = types.StringValue(*result.GetInfo().GetLogoUrl())
		} else {
			info.LogoUrl = types.StringNull()
		}
		if result.GetInfo().GetMarketingUrl() != nil {
			info.MarketingUrl = types.StringValue(*result.GetInfo().GetMarketingUrl())
		} else {
			info.MarketingUrl = types.StringNull()
		}
		if result.GetInfo().GetPrivacyStatementUrl() != nil {
			info.PrivacyStatementUrl = types.StringValue(*result.GetInfo().GetPrivacyStatementUrl())
		} else {
			info.PrivacyStatementUrl = types.StringNull()
		}
		if result.GetInfo().GetSupportUrl() != nil {
			info.SupportUrl = types.StringValue(*result.GetInfo().GetSupportUrl())
		} else {
			info.SupportUrl = types.StringNull()
		}
		if result.GetInfo().GetTermsOfServiceUrl() != nil {
			info.TermsOfServiceUrl = types.StringValue(*result.GetInfo().GetTermsOfServiceUrl())
		} else {
			info.TermsOfServiceUrl = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, info.AttributeTypes(), info)
		state.Info = objectValue
	}
	if result.GetIsDeviceOnlyAuthSupported() != nil {
		state.IsDeviceOnlyAuthSupported = types.BoolValue(*result.GetIsDeviceOnlyAuthSupported())
	} else {
		state.IsDeviceOnlyAuthSupported = types.BoolNull()
	}
	if result.GetIsFallbackPublicClient() != nil {
		state.IsFallbackPublicClient = types.BoolValue(*result.GetIsFallbackPublicClient())
	} else {
		state.IsFallbackPublicClient = types.BoolNull()
	}
	if len(result.GetKeyCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetKeyCredentials() {
			keyCredentials := new(applicationKeyCredentialModel)

			if v.GetCustomKeyIdentifier() != nil {
				keyCredentials.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
			} else {
				keyCredentials.CustomKeyIdentifier = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				keyCredentials.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				keyCredentials.DisplayName = types.StringNull()
			}
			if v.GetEndDateTime() != nil {
				keyCredentials.EndDateTime = types.StringValue(v.GetEndDateTime().String())
			} else {
				keyCredentials.EndDateTime = types.StringNull()
			}
			if v.GetKey() != nil {
				keyCredentials.Key = types.StringValue(string(v.GetKey()[:]))
			} else {
				keyCredentials.Key = types.StringNull()
			}
			if v.GetKeyId() != nil {
				keyCredentials.KeyId = types.StringValue(v.GetKeyId().String())
			} else {
				keyCredentials.KeyId = types.StringNull()
			}
			if v.GetStartDateTime() != nil {
				keyCredentials.StartDateTime = types.StringValue(v.GetStartDateTime().String())
			} else {
				keyCredentials.StartDateTime = types.StringNull()
			}
			if v.GetTypeEscaped() != nil {
				keyCredentials.Type = types.StringValue(*v.GetTypeEscaped())
			} else {
				keyCredentials.Type = types.StringNull()
			}
			if v.GetUsage() != nil {
				keyCredentials.Usage = types.StringValue(*v.GetUsage())
			} else {
				keyCredentials.Usage = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, keyCredentials.AttributeTypes(), keyCredentials)
			objectValues = append(objectValues, objectValue)
		}
		state.KeyCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetLogo() != nil {
		state.Logo = types.StringValue(string(result.GetLogo()[:]))
	} else {
		state.Logo = types.StringNull()
	}
	if result.GetNativeAuthenticationApisEnabled() != nil {
		state.NativeAuthenticationApisEnabled = types.StringValue(result.GetNativeAuthenticationApisEnabled().String())
	} else {
		state.NativeAuthenticationApisEnabled = types.StringNull()
	}
	if result.GetNotes() != nil {
		state.Notes = types.StringValue(*result.GetNotes())
	} else {
		state.Notes = types.StringNull()
	}
	if result.GetOauth2RequirePostResponse() != nil {
		state.Oauth2RequirePostResponse = types.BoolValue(*result.GetOauth2RequirePostResponse())
	} else {
		state.Oauth2RequirePostResponse = types.BoolNull()
	}
	if result.GetOptionalClaims() != nil {
		optionalClaims := new(applicationOptionalClaimsModel)

		if len(result.GetOptionalClaims().GetAccessToken()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetOptionalClaims().GetAccessToken() {
				accessToken := new(applicationOptionalClaimModel)

				if len(v.GetAdditionalProperties()) > 0 {
					var additionalProperties []attr.Value
					for _, v := range v.GetAdditionalProperties() {
						additionalProperties = append(additionalProperties, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, additionalProperties)
					accessToken.AdditionalProperties = listValue
				} else {
					accessToken.AdditionalProperties = types.ListNull(types.StringType)
				}
				if v.GetEssential() != nil {
					accessToken.Essential = types.BoolValue(*v.GetEssential())
				} else {
					accessToken.Essential = types.BoolNull()
				}
				if v.GetName() != nil {
					accessToken.Name = types.StringValue(*v.GetName())
				} else {
					accessToken.Name = types.StringNull()
				}
				if v.GetSource() != nil {
					accessToken.Source = types.StringValue(*v.GetSource())
				} else {
					accessToken.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, accessToken.AttributeTypes(), accessToken)
				objectValues = append(objectValues, objectValue)
			}
			optionalClaims.AccessToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetOptionalClaims().GetIdToken()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetOptionalClaims().GetIdToken() {
				idToken := new(applicationOptionalClaimModel)

				if len(v.GetAdditionalProperties()) > 0 {
					var additionalProperties []attr.Value
					for _, v := range v.GetAdditionalProperties() {
						additionalProperties = append(additionalProperties, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, additionalProperties)
					idToken.AdditionalProperties = listValue
				} else {
					idToken.AdditionalProperties = types.ListNull(types.StringType)
				}
				if v.GetEssential() != nil {
					idToken.Essential = types.BoolValue(*v.GetEssential())
				} else {
					idToken.Essential = types.BoolNull()
				}
				if v.GetName() != nil {
					idToken.Name = types.StringValue(*v.GetName())
				} else {
					idToken.Name = types.StringNull()
				}
				if v.GetSource() != nil {
					idToken.Source = types.StringValue(*v.GetSource())
				} else {
					idToken.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, idToken.AttributeTypes(), idToken)
				objectValues = append(objectValues, objectValue)
			}
			optionalClaims.IdToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetOptionalClaims().GetSaml2Token()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetOptionalClaims().GetSaml2Token() {
				saml2Token := new(applicationOptionalClaimModel)

				if len(v.GetAdditionalProperties()) > 0 {
					var additionalProperties []attr.Value
					for _, v := range v.GetAdditionalProperties() {
						additionalProperties = append(additionalProperties, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, additionalProperties)
					saml2Token.AdditionalProperties = listValue
				} else {
					saml2Token.AdditionalProperties = types.ListNull(types.StringType)
				}
				if v.GetEssential() != nil {
					saml2Token.Essential = types.BoolValue(*v.GetEssential())
				} else {
					saml2Token.Essential = types.BoolNull()
				}
				if v.GetName() != nil {
					saml2Token.Name = types.StringValue(*v.GetName())
				} else {
					saml2Token.Name = types.StringNull()
				}
				if v.GetSource() != nil {
					saml2Token.Source = types.StringValue(*v.GetSource())
				} else {
					saml2Token.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, saml2Token.AttributeTypes(), saml2Token)
				objectValues = append(objectValues, objectValue)
			}
			optionalClaims.Saml2Token, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}

		objectValue, _ := types.ObjectValueFrom(ctx, optionalClaims.AttributeTypes(), optionalClaims)
		state.OptionalClaims = objectValue
	}
	if result.GetParentalControlSettings() != nil {
		parentalControlSettings := new(applicationParentalControlSettingsModel)

		if len(result.GetParentalControlSettings().GetCountriesBlockedForMinors()) > 0 {
			var countriesBlockedForMinors []attr.Value
			for _, v := range result.GetParentalControlSettings().GetCountriesBlockedForMinors() {
				countriesBlockedForMinors = append(countriesBlockedForMinors, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, countriesBlockedForMinors)
			parentalControlSettings.CountriesBlockedForMinors = listValue
		} else {
			parentalControlSettings.CountriesBlockedForMinors = types.ListNull(types.StringType)
		}
		if result.GetParentalControlSettings().GetLegalAgeGroupRule() != nil {
			parentalControlSettings.LegalAgeGroupRule = types.StringValue(*result.GetParentalControlSettings().GetLegalAgeGroupRule())
		} else {
			parentalControlSettings.LegalAgeGroupRule = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, parentalControlSettings.AttributeTypes(), parentalControlSettings)
		state.ParentalControlSettings = objectValue
	}
	if len(result.GetPasswordCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetPasswordCredentials() {
			passwordCredentials := new(applicationPasswordCredentialModel)

			if v.GetCustomKeyIdentifier() != nil {
				passwordCredentials.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
			} else {
				passwordCredentials.CustomKeyIdentifier = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				passwordCredentials.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				passwordCredentials.DisplayName = types.StringNull()
			}
			if v.GetEndDateTime() != nil {
				passwordCredentials.EndDateTime = types.StringValue(v.GetEndDateTime().String())
			} else {
				passwordCredentials.EndDateTime = types.StringNull()
			}
			if v.GetHint() != nil {
				passwordCredentials.Hint = types.StringValue(*v.GetHint())
			} else {
				passwordCredentials.Hint = types.StringNull()
			}
			if v.GetKeyId() != nil {
				passwordCredentials.KeyId = types.StringValue(v.GetKeyId().String())
			} else {
				passwordCredentials.KeyId = types.StringNull()
			}
			if v.GetSecretText() != nil {
				passwordCredentials.SecretText = types.StringValue(*v.GetSecretText())
			} else {
				passwordCredentials.SecretText = types.StringNull()
			}
			if v.GetStartDateTime() != nil {
				passwordCredentials.StartDateTime = types.StringValue(v.GetStartDateTime().String())
			} else {
				passwordCredentials.StartDateTime = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, passwordCredentials.AttributeTypes(), passwordCredentials)
			objectValues = append(objectValues, objectValue)
		}
		state.PasswordCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetPublicClient() != nil {
		publicClient := new(applicationPublicClientApplicationModel)

		if len(result.GetPublicClient().GetRedirectUris()) > 0 {
			var redirectUris []attr.Value
			for _, v := range result.GetPublicClient().GetRedirectUris() {
				redirectUris = append(redirectUris, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, redirectUris)
			publicClient.RedirectUris = listValue
		} else {
			publicClient.RedirectUris = types.ListNull(types.StringType)
		}

		objectValue, _ := types.ObjectValueFrom(ctx, publicClient.AttributeTypes(), publicClient)
		state.PublicClient = objectValue
	}
	if result.GetPublisherDomain() != nil {
		state.PublisherDomain = types.StringValue(*result.GetPublisherDomain())
	} else {
		state.PublisherDomain = types.StringNull()
	}
	if result.GetRequestSignatureVerification() != nil {
		requestSignatureVerification := new(applicationRequestSignatureVerificationModel)

		if result.GetRequestSignatureVerification().GetAllowedWeakAlgorithms() != nil {
			requestSignatureVerification.AllowedWeakAlgorithms = types.StringValue(result.GetRequestSignatureVerification().GetAllowedWeakAlgorithms().String())
		} else {
			requestSignatureVerification.AllowedWeakAlgorithms = types.StringNull()
		}
		if result.GetRequestSignatureVerification().GetIsSignedRequestRequired() != nil {
			requestSignatureVerification.IsSignedRequestRequired = types.BoolValue(*result.GetRequestSignatureVerification().GetIsSignedRequestRequired())
		} else {
			requestSignatureVerification.IsSignedRequestRequired = types.BoolNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, requestSignatureVerification.AttributeTypes(), requestSignatureVerification)
		state.RequestSignatureVerification = objectValue
	}
	if len(result.GetRequiredResourceAccess()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetRequiredResourceAccess() {
			requiredResourceAccess := new(applicationRequiredResourceAccessModel)

			if len(v.GetResourceAccess()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetResourceAccess() {
					resourceAccess := new(applicationResourceAccessModel)

					if v.GetId() != nil {
						resourceAccess.Id = types.StringValue(v.GetId().String())
					} else {
						resourceAccess.Id = types.StringNull()
					}
					if v.GetTypeEscaped() != nil {
						resourceAccess.Type = types.StringValue(*v.GetTypeEscaped())
					} else {
						resourceAccess.Type = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, resourceAccess.AttributeTypes(), resourceAccess)
					objectValues = append(objectValues, objectValue)
				}
				requiredResourceAccess.ResourceAccess, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetResourceAppId() != nil {
				requiredResourceAccess.ResourceAppId = types.StringValue(*v.GetResourceAppId())
			} else {
				requiredResourceAccess.ResourceAppId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, requiredResourceAccess.AttributeTypes(), requiredResourceAccess)
			objectValues = append(objectValues, objectValue)
		}
		state.RequiredResourceAccess, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetSamlMetadataUrl() != nil {
		state.SamlMetadataUrl = types.StringValue(*result.GetSamlMetadataUrl())
	} else {
		state.SamlMetadataUrl = types.StringNull()
	}
	if result.GetServiceManagementReference() != nil {
		state.ServiceManagementReference = types.StringValue(*result.GetServiceManagementReference())
	} else {
		state.ServiceManagementReference = types.StringNull()
	}
	if result.GetServicePrincipalLockConfiguration() != nil {
		servicePrincipalLockConfiguration := new(applicationServicePrincipalLockConfigurationModel)

		if result.GetServicePrincipalLockConfiguration().GetAllProperties() != nil {
			servicePrincipalLockConfiguration.AllProperties = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetAllProperties())
		} else {
			servicePrincipalLockConfiguration.AllProperties = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign() != nil {
			servicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign())
		} else {
			servicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify() != nil {
			servicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify())
		} else {
			servicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetIsEnabled() != nil {
			servicePrincipalLockConfiguration.IsEnabled = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetIsEnabled())
		} else {
			servicePrincipalLockConfiguration.IsEnabled = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId() != nil {
			servicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId())
		} else {
			servicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, servicePrincipalLockConfiguration.AttributeTypes(), servicePrincipalLockConfiguration)
		state.ServicePrincipalLockConfiguration = objectValue
	}
	if result.GetSignInAudience() != nil {
		state.SignInAudience = types.StringValue(*result.GetSignInAudience())
	} else {
		state.SignInAudience = types.StringNull()
	}
	if result.GetSpa() != nil {
		spa := new(applicationSpaApplicationModel)

		if len(result.GetSpa().GetRedirectUris()) > 0 {
			var redirectUris []attr.Value
			for _, v := range result.GetSpa().GetRedirectUris() {
				redirectUris = append(redirectUris, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, redirectUris)
			spa.RedirectUris = listValue
		} else {
			spa.RedirectUris = types.ListNull(types.StringType)
		}

		objectValue, _ := types.ObjectValueFrom(ctx, spa.AttributeTypes(), spa)
		state.Spa = objectValue
	}
	if len(result.GetTags()) > 0 {
		var tags []attr.Value
		for _, v := range result.GetTags() {
			tags = append(tags, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, tags)
		state.Tags = listValue
	} else {
		state.Tags = types.ListNull(types.StringType)
	}
	if result.GetTokenEncryptionKeyId() != nil {
		state.TokenEncryptionKeyId = types.StringValue(result.GetTokenEncryptionKeyId().String())
	} else {
		state.TokenEncryptionKeyId = types.StringNull()
	}
	if result.GetUniqueName() != nil {
		state.UniqueName = types.StringValue(*result.GetUniqueName())
	} else {
		state.UniqueName = types.StringNull()
	}
	if result.GetVerifiedPublisher() != nil {
		verifiedPublisher := new(applicationVerifiedPublisherModel)

		if result.GetVerifiedPublisher().GetAddedDateTime() != nil {
			verifiedPublisher.AddedDateTime = types.StringValue(result.GetVerifiedPublisher().GetAddedDateTime().String())
		} else {
			verifiedPublisher.AddedDateTime = types.StringNull()
		}
		if result.GetVerifiedPublisher().GetDisplayName() != nil {
			verifiedPublisher.DisplayName = types.StringValue(*result.GetVerifiedPublisher().GetDisplayName())
		} else {
			verifiedPublisher.DisplayName = types.StringNull()
		}
		if result.GetVerifiedPublisher().GetVerifiedPublisherId() != nil {
			verifiedPublisher.VerifiedPublisherId = types.StringValue(*result.GetVerifiedPublisher().GetVerifiedPublisherId())
		} else {
			verifiedPublisher.VerifiedPublisherId = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, verifiedPublisher.AttributeTypes(), verifiedPublisher)
		state.VerifiedPublisher = objectValue
	}
	if result.GetWeb() != nil {
		web := new(applicationWebApplicationModel)

		if result.GetWeb().GetHomePageUrl() != nil {
			web.HomePageUrl = types.StringValue(*result.GetWeb().GetHomePageUrl())
		} else {
			web.HomePageUrl = types.StringNull()
		}
		if result.GetWeb().GetImplicitGrantSettings() != nil {
			implicitGrantSettings := new(applicationImplicitGrantSettingsModel)

			if result.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance() != nil {
				implicitGrantSettings.EnableAccessTokenIssuance = types.BoolValue(*result.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance())
			} else {
				implicitGrantSettings.EnableAccessTokenIssuance = types.BoolNull()
			}
			if result.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance() != nil {
				implicitGrantSettings.EnableIdTokenIssuance = types.BoolValue(*result.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance())
			} else {
				implicitGrantSettings.EnableIdTokenIssuance = types.BoolNull()
			}

			objectValue, _ := types.ObjectValueFrom(ctx, implicitGrantSettings.AttributeTypes(), implicitGrantSettings)
			web.ImplicitGrantSettings = objectValue
		}
		if result.GetWeb().GetLogoutUrl() != nil {
			web.LogoutUrl = types.StringValue(*result.GetWeb().GetLogoutUrl())
		} else {
			web.LogoutUrl = types.StringNull()
		}
		if len(result.GetWeb().GetRedirectUriSettings()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetWeb().GetRedirectUriSettings() {
				redirectUriSettings := new(applicationRedirectUriSettingsModel)

				if v.GetUri() != nil {
					redirectUriSettings.Uri = types.StringValue(*v.GetUri())
				} else {
					redirectUriSettings.Uri = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, redirectUriSettings.AttributeTypes(), redirectUriSettings)
				objectValues = append(objectValues, objectValue)
			}
			web.RedirectUriSettings, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetWeb().GetRedirectUris()) > 0 {
			var redirectUris []attr.Value
			for _, v := range result.GetWeb().GetRedirectUris() {
				redirectUris = append(redirectUris, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, redirectUris)
			web.RedirectUris = listValue
		} else {
			web.RedirectUris = types.ListNull(types.StringType)
		}

		objectValue, _ := types.ObjectValueFrom(ctx, web.AttributeTypes(), web)
		state.Web = objectValue
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlan applicationModel
	diags := req.Plan.Get(ctx, &tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfState applicationModel
	diags = req.State.Get(ctx, &tfState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewApplication()

	if !tfPlan.Id.Equal(tfState.Id) {
		tfPlanId := tfPlan.Id.ValueString()
		requestBody.SetId(&tfPlanId)
	}

	if !tfPlan.DeletedDateTime.Equal(tfState.DeletedDateTime) {
		tfPlanDeletedDateTime := tfPlan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !tfPlan.AddIns.Equal(tfState.AddIns) {
		var tfPlanAddIns []models.AddInable
		for k, i := range tfPlan.AddIns.Elements() {
			requestBodyAddIns := models.NewAddIn()
			tfPlanrequestBodyAddIns := applicationAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyAddIns)
			requestBodyAddInsState := applicationAddInModel{}
			types.ListValueFrom(ctx, tfState.AddIns.Elements()[k].Type(ctx), &tfPlanrequestBodyAddIns)

			if !tfPlanrequestBodyAddIns.Id.Equal(requestBodyAddInsState.Id) {
				tfPlanId := tfPlanrequestBodyAddIns.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAddIns.SetId(&u)
			}

			if !tfPlanrequestBodyAddIns.Properties.Equal(requestBodyAddInsState.Properties) {
				var tfPlanProperties []models.KeyValueable
				for k, i := range tfPlanrequestBodyAddIns.Properties.Elements() {
					requestBodyProperties := models.NewKeyValue()
					tfPlanrequestBodyProperties := applicationKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyProperties)
					requestBodyPropertiesState := applicationKeyValueModel{}
					types.ListValueFrom(ctx, requestBodyAddInsState.Properties.Elements()[k].Type(ctx), &tfPlanrequestBodyProperties)

					if !tfPlanrequestBodyProperties.Key.Equal(requestBodyPropertiesState.Key) {
						tfPlanKey := tfPlanrequestBodyProperties.Key.ValueString()
						requestBodyProperties.SetKey(&tfPlanKey)
					}

					if !tfPlanrequestBodyProperties.Value.Equal(requestBodyPropertiesState.Value) {
						tfPlanValue := tfPlanrequestBodyProperties.Value.ValueString()
						requestBodyProperties.SetValue(&tfPlanValue)
					}
				}
				requestBodyAddIns.SetProperties(tfPlanProperties)
			}

			if !tfPlanrequestBodyAddIns.Type.Equal(requestBodyAddInsState.Type) {
				tfPlanType := tfPlanrequestBodyAddIns.Type.ValueString()
				requestBodyAddIns.SetTypeEscaped(&tfPlanType)
			}
		}
		requestBody.SetAddIns(tfPlanAddIns)
	}

	if !tfPlan.Api.Equal(tfState.Api) {
		requestBodyApi := models.NewApiApplication()
		tfPlanrequestBodyApi := applicationApiApplicationModel{}
		tfPlan.Api.As(ctx, &tfPlanrequestBodyApi, basetypes.ObjectAsOptions{})
		requestBodyApiState := applicationApiApplicationModel{}
		tfState.Api.As(ctx, &requestBodyApiState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyApi.AcceptMappedClaims.Equal(requestBodyApiState.AcceptMappedClaims) {
			tfPlanAcceptMappedClaims := tfPlanrequestBodyApi.AcceptMappedClaims.ValueBool()
			requestBodyApi.SetAcceptMappedClaims(&tfPlanAcceptMappedClaims)
		}

		if !tfPlanrequestBodyApi.KnownClientApplications.Equal(requestBodyApiState.KnownClientApplications) {
			var KnownClientApplications []uuid.UUID
			for _, i := range tfPlanrequestBodyApi.KnownClientApplications.Elements() {
				u, _ := uuid.Parse(i.String())
				KnownClientApplications = append(KnownClientApplications, u)
			}
			requestBodyApi.SetKnownClientApplications(KnownClientApplications)
		}

		if !tfPlanrequestBodyApi.Oauth2PermissionScopes.Equal(requestBodyApiState.Oauth2PermissionScopes) {
			var tfPlanOauth2PermissionScopes []models.PermissionScopeable
			for k, i := range tfPlanrequestBodyApi.Oauth2PermissionScopes.Elements() {
				requestBodyOauth2PermissionScopes := models.NewPermissionScope()
				tfPlanrequestBodyOauth2PermissionScopes := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyOauth2PermissionScopes)
				requestBodyOauth2PermissionScopesState := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, requestBodyApiState.Oauth2PermissionScopes.Elements()[k].Type(ctx), &tfPlanrequestBodyOauth2PermissionScopes)

				if !tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDescription.Equal(requestBodyOauth2PermissionScopesState.AdminConsentDescription) {
					tfPlanAdminConsentDescription := tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDescription.ValueString()
					requestBodyOauth2PermissionScopes.SetAdminConsentDescription(&tfPlanAdminConsentDescription)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDisplayName.Equal(requestBodyOauth2PermissionScopesState.AdminConsentDisplayName) {
					tfPlanAdminConsentDisplayName := tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDisplayName.ValueString()
					requestBodyOauth2PermissionScopes.SetAdminConsentDisplayName(&tfPlanAdminConsentDisplayName)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.Id.Equal(requestBodyOauth2PermissionScopesState.Id) {
					tfPlanId := tfPlanrequestBodyOauth2PermissionScopes.Id.ValueString()
					u, _ := uuid.Parse(tfPlanId)
					requestBodyOauth2PermissionScopes.SetId(&u)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.IsEnabled.Equal(requestBodyOauth2PermissionScopesState.IsEnabled) {
					tfPlanIsEnabled := tfPlanrequestBodyOauth2PermissionScopes.IsEnabled.ValueBool()
					requestBodyOauth2PermissionScopes.SetIsEnabled(&tfPlanIsEnabled)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.Origin.Equal(requestBodyOauth2PermissionScopesState.Origin) {
					tfPlanOrigin := tfPlanrequestBodyOauth2PermissionScopes.Origin.ValueString()
					requestBodyOauth2PermissionScopes.SetOrigin(&tfPlanOrigin)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.Type.Equal(requestBodyOauth2PermissionScopesState.Type) {
					tfPlanType := tfPlanrequestBodyOauth2PermissionScopes.Type.ValueString()
					requestBodyOauth2PermissionScopes.SetTypeEscaped(&tfPlanType)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.UserConsentDescription.Equal(requestBodyOauth2PermissionScopesState.UserConsentDescription) {
					tfPlanUserConsentDescription := tfPlanrequestBodyOauth2PermissionScopes.UserConsentDescription.ValueString()
					requestBodyOauth2PermissionScopes.SetUserConsentDescription(&tfPlanUserConsentDescription)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.UserConsentDisplayName.Equal(requestBodyOauth2PermissionScopesState.UserConsentDisplayName) {
					tfPlanUserConsentDisplayName := tfPlanrequestBodyOauth2PermissionScopes.UserConsentDisplayName.ValueString()
					requestBodyOauth2PermissionScopes.SetUserConsentDisplayName(&tfPlanUserConsentDisplayName)
				}

				if !tfPlanrequestBodyOauth2PermissionScopes.Value.Equal(requestBodyOauth2PermissionScopesState.Value) {
					tfPlanValue := tfPlanrequestBodyOauth2PermissionScopes.Value.ValueString()
					requestBodyOauth2PermissionScopes.SetValue(&tfPlanValue)
				}
			}
			requestBodyApi.SetOauth2PermissionScopes(tfPlanOauth2PermissionScopes)
		}

		if !tfPlanrequestBodyApi.PreAuthorizedApplications.Equal(requestBodyApiState.PreAuthorizedApplications) {
			var tfPlanPreAuthorizedApplications []models.PreAuthorizedApplicationable
			for k, i := range tfPlanrequestBodyApi.PreAuthorizedApplications.Elements() {
				requestBodyPreAuthorizedApplications := models.NewPreAuthorizedApplication()
				tfPlanrequestBodyPreAuthorizedApplications := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyPreAuthorizedApplications)
				requestBodyPreAuthorizedApplicationsState := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, requestBodyApiState.PreAuthorizedApplications.Elements()[k].Type(ctx), &tfPlanrequestBodyPreAuthorizedApplications)

				if !tfPlanrequestBodyPreAuthorizedApplications.AppId.Equal(requestBodyPreAuthorizedApplicationsState.AppId) {
					tfPlanAppId := tfPlanrequestBodyPreAuthorizedApplications.AppId.ValueString()
					requestBodyPreAuthorizedApplications.SetAppId(&tfPlanAppId)
				}

				if !tfPlanrequestBodyPreAuthorizedApplications.DelegatedPermissionIds.Equal(requestBodyPreAuthorizedApplicationsState.DelegatedPermissionIds) {
					var stringArrayDelegatedPermissionIds []string
					for _, i := range tfPlanrequestBodyPreAuthorizedApplications.DelegatedPermissionIds.Elements() {
						stringArrayDelegatedPermissionIds = append(stringArrayDelegatedPermissionIds, i.String())
					}
					requestBodyPreAuthorizedApplications.SetDelegatedPermissionIds(stringArrayDelegatedPermissionIds)
				}
			}
			requestBodyApi.SetPreAuthorizedApplications(tfPlanPreAuthorizedApplications)
		}
		requestBody.SetApi(requestBodyApi)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyApi.AttributeTypes(), tfPlanrequestBodyApi)
		tfPlan.Api = objectValue
	}

	if !tfPlan.AppId.Equal(tfState.AppId) {
		tfPlanAppId := tfPlan.AppId.ValueString()
		requestBody.SetAppId(&tfPlanAppId)
	}

	if !tfPlan.AppRoles.Equal(tfState.AppRoles) {
		var tfPlanAppRoles []models.AppRoleable
		for k, i := range tfPlan.AppRoles.Elements() {
			requestBodyAppRoles := models.NewAppRole()
			tfPlanrequestBodyAppRoles := applicationAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyAppRoles)
			requestBodyAppRolesState := applicationAppRoleModel{}
			types.ListValueFrom(ctx, tfState.AppRoles.Elements()[k].Type(ctx), &tfPlanrequestBodyAppRoles)

			if !tfPlanrequestBodyAppRoles.AllowedMemberTypes.Equal(requestBodyAppRolesState.AllowedMemberTypes) {
				var stringArrayAllowedMemberTypes []string
				for _, i := range tfPlanrequestBodyAppRoles.AllowedMemberTypes.Elements() {
					stringArrayAllowedMemberTypes = append(stringArrayAllowedMemberTypes, i.String())
				}
				requestBodyAppRoles.SetAllowedMemberTypes(stringArrayAllowedMemberTypes)
			}

			if !tfPlanrequestBodyAppRoles.Description.Equal(requestBodyAppRolesState.Description) {
				tfPlanDescription := tfPlanrequestBodyAppRoles.Description.ValueString()
				requestBodyAppRoles.SetDescription(&tfPlanDescription)
			}

			if !tfPlanrequestBodyAppRoles.DisplayName.Equal(requestBodyAppRolesState.DisplayName) {
				tfPlanDisplayName := tfPlanrequestBodyAppRoles.DisplayName.ValueString()
				requestBodyAppRoles.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanrequestBodyAppRoles.Id.Equal(requestBodyAppRolesState.Id) {
				tfPlanId := tfPlanrequestBodyAppRoles.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAppRoles.SetId(&u)
			}

			if !tfPlanrequestBodyAppRoles.IsEnabled.Equal(requestBodyAppRolesState.IsEnabled) {
				tfPlanIsEnabled := tfPlanrequestBodyAppRoles.IsEnabled.ValueBool()
				requestBodyAppRoles.SetIsEnabled(&tfPlanIsEnabled)
			}

			if !tfPlanrequestBodyAppRoles.Origin.Equal(requestBodyAppRolesState.Origin) {
				tfPlanOrigin := tfPlanrequestBodyAppRoles.Origin.ValueString()
				requestBodyAppRoles.SetOrigin(&tfPlanOrigin)
			}

			if !tfPlanrequestBodyAppRoles.Value.Equal(requestBodyAppRolesState.Value) {
				tfPlanValue := tfPlanrequestBodyAppRoles.Value.ValueString()
				requestBodyAppRoles.SetValue(&tfPlanValue)
			}
		}
		requestBody.SetAppRoles(tfPlanAppRoles)
	}

	if !tfPlan.ApplicationTemplateId.Equal(tfState.ApplicationTemplateId) {
		tfPlanApplicationTemplateId := tfPlan.ApplicationTemplateId.ValueString()
		requestBody.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	}

	if !tfPlan.Certification.Equal(tfState.Certification) {
		requestBodyCertification := models.NewCertification()
		tfPlanrequestBodyCertification := applicationCertificationModel{}
		tfPlan.Certification.As(ctx, &tfPlanrequestBodyCertification, basetypes.ObjectAsOptions{})
		requestBodyCertificationState := applicationCertificationModel{}
		tfState.Certification.As(ctx, &requestBodyCertificationState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyCertification.CertificationDetailsUrl.Equal(requestBodyCertificationState.CertificationDetailsUrl) {
			tfPlanCertificationDetailsUrl := tfPlanrequestBodyCertification.CertificationDetailsUrl.ValueString()
			requestBodyCertification.SetCertificationDetailsUrl(&tfPlanCertificationDetailsUrl)
		}

		if !tfPlanrequestBodyCertification.CertificationExpirationDateTime.Equal(requestBodyCertificationState.CertificationExpirationDateTime) {
			tfPlanCertificationExpirationDateTime := tfPlanrequestBodyCertification.CertificationExpirationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanCertificationExpirationDateTime)
			requestBodyCertification.SetCertificationExpirationDateTime(&t)
		}

		if !tfPlanrequestBodyCertification.IsCertifiedByMicrosoft.Equal(requestBodyCertificationState.IsCertifiedByMicrosoft) {
			tfPlanIsCertifiedByMicrosoft := tfPlanrequestBodyCertification.IsCertifiedByMicrosoft.ValueBool()
			requestBodyCertification.SetIsCertifiedByMicrosoft(&tfPlanIsCertifiedByMicrosoft)
		}

		if !tfPlanrequestBodyCertification.IsPublisherAttested.Equal(requestBodyCertificationState.IsPublisherAttested) {
			tfPlanIsPublisherAttested := tfPlanrequestBodyCertification.IsPublisherAttested.ValueBool()
			requestBodyCertification.SetIsPublisherAttested(&tfPlanIsPublisherAttested)
		}

		if !tfPlanrequestBodyCertification.LastCertificationDateTime.Equal(requestBodyCertificationState.LastCertificationDateTime) {
			tfPlanLastCertificationDateTime := tfPlanrequestBodyCertification.LastCertificationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastCertificationDateTime)
			requestBodyCertification.SetLastCertificationDateTime(&t)
		}
		requestBody.SetCertification(requestBodyCertification)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyCertification.AttributeTypes(), tfPlanrequestBodyCertification)
		tfPlan.Certification = objectValue
	}

	if !tfPlan.CreatedDateTime.Equal(tfState.CreatedDateTime) {
		tfPlanCreatedDateTime := tfPlan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	}

	if !tfPlan.DefaultRedirectUri.Equal(tfState.DefaultRedirectUri) {
		tfPlanDefaultRedirectUri := tfPlan.DefaultRedirectUri.ValueString()
		requestBody.SetDefaultRedirectUri(&tfPlanDefaultRedirectUri)
	}

	if !tfPlan.Description.Equal(tfState.Description) {
		tfPlanDescription := tfPlan.Description.ValueString()
		requestBody.SetDescription(&tfPlanDescription)
	}

	if !tfPlan.DisabledByMicrosoftStatus.Equal(tfState.DisabledByMicrosoftStatus) {
		tfPlanDisabledByMicrosoftStatus := tfPlan.DisabledByMicrosoftStatus.ValueString()
		requestBody.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	}

	if !tfPlan.DisplayName.Equal(tfState.DisplayName) {
		tfPlanDisplayName := tfPlan.DisplayName.ValueString()
		requestBody.SetDisplayName(&tfPlanDisplayName)
	}

	if !tfPlan.GroupMembershipClaims.Equal(tfState.GroupMembershipClaims) {
		tfPlanGroupMembershipClaims := tfPlan.GroupMembershipClaims.ValueString()
		requestBody.SetGroupMembershipClaims(&tfPlanGroupMembershipClaims)
	}

	if !tfPlan.IdentifierUris.Equal(tfState.IdentifierUris) {
		var stringArrayIdentifierUris []string
		for _, i := range tfPlan.IdentifierUris.Elements() {
			stringArrayIdentifierUris = append(stringArrayIdentifierUris, i.String())
		}
		requestBody.SetIdentifierUris(stringArrayIdentifierUris)
	}

	if !tfPlan.Info.Equal(tfState.Info) {
		requestBodyInfo := models.NewInformationalUrl()
		tfPlanrequestBodyInfo := applicationInformationalUrlModel{}
		tfPlan.Info.As(ctx, &tfPlanrequestBodyInfo, basetypes.ObjectAsOptions{})
		requestBodyInfoState := applicationInformationalUrlModel{}
		tfState.Info.As(ctx, &requestBodyInfoState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyInfo.LogoUrl.Equal(requestBodyInfoState.LogoUrl) {
			tfPlanLogoUrl := tfPlanrequestBodyInfo.LogoUrl.ValueString()
			requestBodyInfo.SetLogoUrl(&tfPlanLogoUrl)
		}

		if !tfPlanrequestBodyInfo.MarketingUrl.Equal(requestBodyInfoState.MarketingUrl) {
			tfPlanMarketingUrl := tfPlanrequestBodyInfo.MarketingUrl.ValueString()
			requestBodyInfo.SetMarketingUrl(&tfPlanMarketingUrl)
		}

		if !tfPlanrequestBodyInfo.PrivacyStatementUrl.Equal(requestBodyInfoState.PrivacyStatementUrl) {
			tfPlanPrivacyStatementUrl := tfPlanrequestBodyInfo.PrivacyStatementUrl.ValueString()
			requestBodyInfo.SetPrivacyStatementUrl(&tfPlanPrivacyStatementUrl)
		}

		if !tfPlanrequestBodyInfo.SupportUrl.Equal(requestBodyInfoState.SupportUrl) {
			tfPlanSupportUrl := tfPlanrequestBodyInfo.SupportUrl.ValueString()
			requestBodyInfo.SetSupportUrl(&tfPlanSupportUrl)
		}

		if !tfPlanrequestBodyInfo.TermsOfServiceUrl.Equal(requestBodyInfoState.TermsOfServiceUrl) {
			tfPlanTermsOfServiceUrl := tfPlanrequestBodyInfo.TermsOfServiceUrl.ValueString()
			requestBodyInfo.SetTermsOfServiceUrl(&tfPlanTermsOfServiceUrl)
		}
		requestBody.SetInfo(requestBodyInfo)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyInfo.AttributeTypes(), tfPlanrequestBodyInfo)
		tfPlan.Info = objectValue
	}

	if !tfPlan.IsDeviceOnlyAuthSupported.Equal(tfState.IsDeviceOnlyAuthSupported) {
		tfPlanIsDeviceOnlyAuthSupported := tfPlan.IsDeviceOnlyAuthSupported.ValueBool()
		requestBody.SetIsDeviceOnlyAuthSupported(&tfPlanIsDeviceOnlyAuthSupported)
	}

	if !tfPlan.IsFallbackPublicClient.Equal(tfState.IsFallbackPublicClient) {
		tfPlanIsFallbackPublicClient := tfPlan.IsFallbackPublicClient.ValueBool()
		requestBody.SetIsFallbackPublicClient(&tfPlanIsFallbackPublicClient)
	}

	if !tfPlan.KeyCredentials.Equal(tfState.KeyCredentials) {
		var tfPlanKeyCredentials []models.KeyCredentialable
		for k, i := range tfPlan.KeyCredentials.Elements() {
			requestBodyKeyCredentials := models.NewKeyCredential()
			tfPlanrequestBodyKeyCredentials := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyKeyCredentials)
			requestBodyKeyCredentialsState := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, tfState.KeyCredentials.Elements()[k].Type(ctx), &tfPlanrequestBodyKeyCredentials)

			if !tfPlanrequestBodyKeyCredentials.CustomKeyIdentifier.Equal(requestBodyKeyCredentialsState.CustomKeyIdentifier) {
				tfPlanCustomKeyIdentifier := tfPlanrequestBodyKeyCredentials.CustomKeyIdentifier.ValueString()
				requestBodyKeyCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			}

			if !tfPlanrequestBodyKeyCredentials.DisplayName.Equal(requestBodyKeyCredentialsState.DisplayName) {
				tfPlanDisplayName := tfPlanrequestBodyKeyCredentials.DisplayName.ValueString()
				requestBodyKeyCredentials.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanrequestBodyKeyCredentials.EndDateTime.Equal(requestBodyKeyCredentialsState.EndDateTime) {
				tfPlanEndDateTime := tfPlanrequestBodyKeyCredentials.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyKeyCredentials.SetEndDateTime(&t)
			}

			if !tfPlanrequestBodyKeyCredentials.Key.Equal(requestBodyKeyCredentialsState.Key) {
				tfPlanKey := tfPlanrequestBodyKeyCredentials.Key.ValueString()
				requestBodyKeyCredentials.SetKey([]byte(tfPlanKey))
			}

			if !tfPlanrequestBodyKeyCredentials.KeyId.Equal(requestBodyKeyCredentialsState.KeyId) {
				tfPlanKeyId := tfPlanrequestBodyKeyCredentials.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyKeyCredentials.SetKeyId(&u)
			}

			if !tfPlanrequestBodyKeyCredentials.StartDateTime.Equal(requestBodyKeyCredentialsState.StartDateTime) {
				tfPlanStartDateTime := tfPlanrequestBodyKeyCredentials.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyKeyCredentials.SetStartDateTime(&t)
			}

			if !tfPlanrequestBodyKeyCredentials.Type.Equal(requestBodyKeyCredentialsState.Type) {
				tfPlanType := tfPlanrequestBodyKeyCredentials.Type.ValueString()
				requestBodyKeyCredentials.SetTypeEscaped(&tfPlanType)
			}

			if !tfPlanrequestBodyKeyCredentials.Usage.Equal(requestBodyKeyCredentialsState.Usage) {
				tfPlanUsage := tfPlanrequestBodyKeyCredentials.Usage.ValueString()
				requestBodyKeyCredentials.SetUsage(&tfPlanUsage)
			}
		}
		requestBody.SetKeyCredentials(tfPlanKeyCredentials)
	}

	if !tfPlan.Logo.Equal(tfState.Logo) {
		tfPlanLogo := tfPlan.Logo.ValueString()
		requestBody.SetLogo([]byte(tfPlanLogo))
	}

	if !tfPlan.NativeAuthenticationApisEnabled.Equal(tfState.NativeAuthenticationApisEnabled) {
		tfPlanNativeAuthenticationApisEnabled := tfPlan.NativeAuthenticationApisEnabled.ValueString()
		parsedNativeAuthenticationApisEnabled, _ := models.ParseNativeAuthenticationApisEnabled(tfPlanNativeAuthenticationApisEnabled)
		assertedNativeAuthenticationApisEnabled := parsedNativeAuthenticationApisEnabled.(models.NativeAuthenticationApisEnabled)
		requestBody.SetNativeAuthenticationApisEnabled(&assertedNativeAuthenticationApisEnabled)
	}

	if !tfPlan.Notes.Equal(tfState.Notes) {
		tfPlanNotes := tfPlan.Notes.ValueString()
		requestBody.SetNotes(&tfPlanNotes)
	}

	if !tfPlan.Oauth2RequirePostResponse.Equal(tfState.Oauth2RequirePostResponse) {
		tfPlanOauth2RequirePostResponse := tfPlan.Oauth2RequirePostResponse.ValueBool()
		requestBody.SetOauth2RequirePostResponse(&tfPlanOauth2RequirePostResponse)
	}

	if !tfPlan.OptionalClaims.Equal(tfState.OptionalClaims) {
		requestBodyOptionalClaims := models.NewOptionalClaims()
		tfPlanrequestBodyOptionalClaims := applicationOptionalClaimsModel{}
		tfPlan.OptionalClaims.As(ctx, &tfPlanrequestBodyOptionalClaims, basetypes.ObjectAsOptions{})
		requestBodyOptionalClaimsState := applicationOptionalClaimsModel{}
		tfState.OptionalClaims.As(ctx, &requestBodyOptionalClaimsState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyOptionalClaims.AccessToken.Equal(requestBodyOptionalClaimsState.AccessToken) {
			var tfPlanAccessToken []models.OptionalClaimable
			for k, i := range tfPlanrequestBodyOptionalClaims.AccessToken.Elements() {
				requestBodyAccessToken := models.NewOptionalClaim()
				tfPlanrequestBodyAccessToken := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyAccessToken)
				requestBodyAccessTokenState := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, requestBodyOptionalClaimsState.AccessToken.Elements()[k].Type(ctx), &tfPlanrequestBodyAccessToken)

				if !tfPlanrequestBodyAccessToken.AdditionalProperties.Equal(requestBodyAccessTokenState.AdditionalProperties) {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanrequestBodyAccessToken.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyAccessToken.SetAdditionalProperties(stringArrayAdditionalProperties)
				}

				if !tfPlanrequestBodyAccessToken.Essential.Equal(requestBodyAccessTokenState.Essential) {
					tfPlanEssential := tfPlanrequestBodyAccessToken.Essential.ValueBool()
					requestBodyAccessToken.SetEssential(&tfPlanEssential)
				}

				if !tfPlanrequestBodyAccessToken.Name.Equal(requestBodyAccessTokenState.Name) {
					tfPlanName := tfPlanrequestBodyAccessToken.Name.ValueString()
					requestBodyAccessToken.SetName(&tfPlanName)
				}

				if !tfPlanrequestBodyAccessToken.Source.Equal(requestBodyAccessTokenState.Source) {
					tfPlanSource := tfPlanrequestBodyAccessToken.Source.ValueString()
					requestBodyAccessToken.SetSource(&tfPlanSource)
				}
			}
			requestBodyOptionalClaims.SetAccessToken(tfPlanAccessToken)
		}

		if !tfPlanrequestBodyOptionalClaims.IdToken.Equal(requestBodyOptionalClaimsState.IdToken) {
			var tfPlanIdToken []models.OptionalClaimable
			for k, i := range tfPlanrequestBodyOptionalClaims.IdToken.Elements() {
				requestBodyIdToken := models.NewOptionalClaim()
				tfPlanrequestBodyIdToken := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyIdToken)
				requestBodyIdTokenState := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, requestBodyOptionalClaimsState.IdToken.Elements()[k].Type(ctx), &tfPlanrequestBodyIdToken)

				if !tfPlanrequestBodyIdToken.AdditionalProperties.Equal(requestBodyIdTokenState.AdditionalProperties) {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanrequestBodyIdToken.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyIdToken.SetAdditionalProperties(stringArrayAdditionalProperties)
				}

				if !tfPlanrequestBodyIdToken.Essential.Equal(requestBodyIdTokenState.Essential) {
					tfPlanEssential := tfPlanrequestBodyIdToken.Essential.ValueBool()
					requestBodyIdToken.SetEssential(&tfPlanEssential)
				}

				if !tfPlanrequestBodyIdToken.Name.Equal(requestBodyIdTokenState.Name) {
					tfPlanName := tfPlanrequestBodyIdToken.Name.ValueString()
					requestBodyIdToken.SetName(&tfPlanName)
				}

				if !tfPlanrequestBodyIdToken.Source.Equal(requestBodyIdTokenState.Source) {
					tfPlanSource := tfPlanrequestBodyIdToken.Source.ValueString()
					requestBodyIdToken.SetSource(&tfPlanSource)
				}
			}
			requestBodyOptionalClaims.SetIdToken(tfPlanIdToken)
		}

		if !tfPlanrequestBodyOptionalClaims.Saml2Token.Equal(requestBodyOptionalClaimsState.Saml2Token) {
			var tfPlanSaml2Token []models.OptionalClaimable
			for k, i := range tfPlanrequestBodyOptionalClaims.Saml2Token.Elements() {
				requestBodySaml2Token := models.NewOptionalClaim()
				tfPlanrequestBodySaml2Token := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodySaml2Token)
				requestBodySaml2TokenState := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, requestBodyOptionalClaimsState.Saml2Token.Elements()[k].Type(ctx), &tfPlanrequestBodySaml2Token)

				if !tfPlanrequestBodySaml2Token.AdditionalProperties.Equal(requestBodySaml2TokenState.AdditionalProperties) {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanrequestBodySaml2Token.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodySaml2Token.SetAdditionalProperties(stringArrayAdditionalProperties)
				}

				if !tfPlanrequestBodySaml2Token.Essential.Equal(requestBodySaml2TokenState.Essential) {
					tfPlanEssential := tfPlanrequestBodySaml2Token.Essential.ValueBool()
					requestBodySaml2Token.SetEssential(&tfPlanEssential)
				}

				if !tfPlanrequestBodySaml2Token.Name.Equal(requestBodySaml2TokenState.Name) {
					tfPlanName := tfPlanrequestBodySaml2Token.Name.ValueString()
					requestBodySaml2Token.SetName(&tfPlanName)
				}

				if !tfPlanrequestBodySaml2Token.Source.Equal(requestBodySaml2TokenState.Source) {
					tfPlanSource := tfPlanrequestBodySaml2Token.Source.ValueString()
					requestBodySaml2Token.SetSource(&tfPlanSource)
				}
			}
			requestBodyOptionalClaims.SetSaml2Token(tfPlanSaml2Token)
		}
		requestBody.SetOptionalClaims(requestBodyOptionalClaims)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyOptionalClaims.AttributeTypes(), tfPlanrequestBodyOptionalClaims)
		tfPlan.OptionalClaims = objectValue
	}

	if !tfPlan.ParentalControlSettings.Equal(tfState.ParentalControlSettings) {
		requestBodyParentalControlSettings := models.NewParentalControlSettings()
		tfPlanrequestBodyParentalControlSettings := applicationParentalControlSettingsModel{}
		tfPlan.ParentalControlSettings.As(ctx, &tfPlanrequestBodyParentalControlSettings, basetypes.ObjectAsOptions{})
		requestBodyParentalControlSettingsState := applicationParentalControlSettingsModel{}
		tfState.ParentalControlSettings.As(ctx, &requestBodyParentalControlSettingsState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyParentalControlSettings.CountriesBlockedForMinors.Equal(requestBodyParentalControlSettingsState.CountriesBlockedForMinors) {
			var stringArrayCountriesBlockedForMinors []string
			for _, i := range tfPlanrequestBodyParentalControlSettings.CountriesBlockedForMinors.Elements() {
				stringArrayCountriesBlockedForMinors = append(stringArrayCountriesBlockedForMinors, i.String())
			}
			requestBodyParentalControlSettings.SetCountriesBlockedForMinors(stringArrayCountriesBlockedForMinors)
		}

		if !tfPlanrequestBodyParentalControlSettings.LegalAgeGroupRule.Equal(requestBodyParentalControlSettingsState.LegalAgeGroupRule) {
			tfPlanLegalAgeGroupRule := tfPlanrequestBodyParentalControlSettings.LegalAgeGroupRule.ValueString()
			requestBodyParentalControlSettings.SetLegalAgeGroupRule(&tfPlanLegalAgeGroupRule)
		}
		requestBody.SetParentalControlSettings(requestBodyParentalControlSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyParentalControlSettings.AttributeTypes(), tfPlanrequestBodyParentalControlSettings)
		tfPlan.ParentalControlSettings = objectValue
	}

	if !tfPlan.PasswordCredentials.Equal(tfState.PasswordCredentials) {
		var tfPlanPasswordCredentials []models.PasswordCredentialable
		for k, i := range tfPlan.PasswordCredentials.Elements() {
			requestBodyPasswordCredentials := models.NewPasswordCredential()
			tfPlanrequestBodyPasswordCredentials := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyPasswordCredentials)
			requestBodyPasswordCredentialsState := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, tfState.PasswordCredentials.Elements()[k].Type(ctx), &tfPlanrequestBodyPasswordCredentials)

			if !tfPlanrequestBodyPasswordCredentials.CustomKeyIdentifier.Equal(requestBodyPasswordCredentialsState.CustomKeyIdentifier) {
				tfPlanCustomKeyIdentifier := tfPlanrequestBodyPasswordCredentials.CustomKeyIdentifier.ValueString()
				requestBodyPasswordCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			}

			if !tfPlanrequestBodyPasswordCredentials.DisplayName.Equal(requestBodyPasswordCredentialsState.DisplayName) {
				tfPlanDisplayName := tfPlanrequestBodyPasswordCredentials.DisplayName.ValueString()
				requestBodyPasswordCredentials.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanrequestBodyPasswordCredentials.EndDateTime.Equal(requestBodyPasswordCredentialsState.EndDateTime) {
				tfPlanEndDateTime := tfPlanrequestBodyPasswordCredentials.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyPasswordCredentials.SetEndDateTime(&t)
			}

			if !tfPlanrequestBodyPasswordCredentials.Hint.Equal(requestBodyPasswordCredentialsState.Hint) {
				tfPlanHint := tfPlanrequestBodyPasswordCredentials.Hint.ValueString()
				requestBodyPasswordCredentials.SetHint(&tfPlanHint)
			}

			if !tfPlanrequestBodyPasswordCredentials.KeyId.Equal(requestBodyPasswordCredentialsState.KeyId) {
				tfPlanKeyId := tfPlanrequestBodyPasswordCredentials.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyPasswordCredentials.SetKeyId(&u)
			}

			if !tfPlanrequestBodyPasswordCredentials.SecretText.Equal(requestBodyPasswordCredentialsState.SecretText) {
				tfPlanSecretText := tfPlanrequestBodyPasswordCredentials.SecretText.ValueString()
				requestBodyPasswordCredentials.SetSecretText(&tfPlanSecretText)
			}

			if !tfPlanrequestBodyPasswordCredentials.StartDateTime.Equal(requestBodyPasswordCredentialsState.StartDateTime) {
				tfPlanStartDateTime := tfPlanrequestBodyPasswordCredentials.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyPasswordCredentials.SetStartDateTime(&t)
			}
		}
		requestBody.SetPasswordCredentials(tfPlanPasswordCredentials)
	}

	if !tfPlan.PublicClient.Equal(tfState.PublicClient) {
		requestBodyPublicClient := models.NewPublicClientApplication()
		tfPlanrequestBodyPublicClient := applicationPublicClientApplicationModel{}
		tfPlan.PublicClient.As(ctx, &tfPlanrequestBodyPublicClient, basetypes.ObjectAsOptions{})
		requestBodyPublicClientState := applicationPublicClientApplicationModel{}
		tfState.PublicClient.As(ctx, &requestBodyPublicClientState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyPublicClient.RedirectUris.Equal(requestBodyPublicClientState.RedirectUris) {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanrequestBodyPublicClient.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodyPublicClient.SetRedirectUris(stringArrayRedirectUris)
		}
		requestBody.SetPublicClient(requestBodyPublicClient)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyPublicClient.AttributeTypes(), tfPlanrequestBodyPublicClient)
		tfPlan.PublicClient = objectValue
	}

	if !tfPlan.PublisherDomain.Equal(tfState.PublisherDomain) {
		tfPlanPublisherDomain := tfPlan.PublisherDomain.ValueString()
		requestBody.SetPublisherDomain(&tfPlanPublisherDomain)
	}

	if !tfPlan.RequestSignatureVerification.Equal(tfState.RequestSignatureVerification) {
		requestBodyRequestSignatureVerification := models.NewRequestSignatureVerification()
		tfPlanrequestBodyRequestSignatureVerification := applicationRequestSignatureVerificationModel{}
		tfPlan.RequestSignatureVerification.As(ctx, &tfPlanrequestBodyRequestSignatureVerification, basetypes.ObjectAsOptions{})
		requestBodyRequestSignatureVerificationState := applicationRequestSignatureVerificationModel{}
		tfState.RequestSignatureVerification.As(ctx, &requestBodyRequestSignatureVerificationState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyRequestSignatureVerification.AllowedWeakAlgorithms.Equal(requestBodyRequestSignatureVerificationState.AllowedWeakAlgorithms) {
			tfPlanAllowedWeakAlgorithms := tfPlanrequestBodyRequestSignatureVerification.AllowedWeakAlgorithms.ValueString()
			parsedAllowedWeakAlgorithms, _ := models.ParseWeakAlgorithms(tfPlanAllowedWeakAlgorithms)
			assertedAllowedWeakAlgorithms := parsedAllowedWeakAlgorithms.(models.WeakAlgorithms)
			requestBodyRequestSignatureVerification.SetAllowedWeakAlgorithms(&assertedAllowedWeakAlgorithms)
		}

		if !tfPlanrequestBodyRequestSignatureVerification.IsSignedRequestRequired.Equal(requestBodyRequestSignatureVerificationState.IsSignedRequestRequired) {
			tfPlanIsSignedRequestRequired := tfPlanrequestBodyRequestSignatureVerification.IsSignedRequestRequired.ValueBool()
			requestBodyRequestSignatureVerification.SetIsSignedRequestRequired(&tfPlanIsSignedRequestRequired)
		}
		requestBody.SetRequestSignatureVerification(requestBodyRequestSignatureVerification)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyRequestSignatureVerification.AttributeTypes(), tfPlanrequestBodyRequestSignatureVerification)
		tfPlan.RequestSignatureVerification = objectValue
	}

	if !tfPlan.RequiredResourceAccess.Equal(tfState.RequiredResourceAccess) {
		var tfPlanRequiredResourceAccess []models.RequiredResourceAccessable
		for k, i := range tfPlan.RequiredResourceAccess.Elements() {
			requestBodyRequiredResourceAccess := models.NewRequiredResourceAccess()
			tfPlanrequestBodyRequiredResourceAccess := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyRequiredResourceAccess)
			requestBodyRequiredResourceAccessState := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, tfState.RequiredResourceAccess.Elements()[k].Type(ctx), &tfPlanrequestBodyRequiredResourceAccess)

			if !tfPlanrequestBodyRequiredResourceAccess.ResourceAccess.Equal(requestBodyRequiredResourceAccessState.ResourceAccess) {
				var tfPlanResourceAccess []models.ResourceAccessable
				for k, i := range tfPlanrequestBodyRequiredResourceAccess.ResourceAccess.Elements() {
					requestBodyResourceAccess := models.NewResourceAccess()
					tfPlanrequestBodyResourceAccess := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyResourceAccess)
					requestBodyResourceAccessState := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, requestBodyRequiredResourceAccessState.ResourceAccess.Elements()[k].Type(ctx), &tfPlanrequestBodyResourceAccess)

					if !tfPlanrequestBodyResourceAccess.Id.Equal(requestBodyResourceAccessState.Id) {
						tfPlanId := tfPlanrequestBodyResourceAccess.Id.ValueString()
						u, _ := uuid.Parse(tfPlanId)
						requestBodyResourceAccess.SetId(&u)
					}

					if !tfPlanrequestBodyResourceAccess.Type.Equal(requestBodyResourceAccessState.Type) {
						tfPlanType := tfPlanrequestBodyResourceAccess.Type.ValueString()
						requestBodyResourceAccess.SetTypeEscaped(&tfPlanType)
					}
				}
				requestBodyRequiredResourceAccess.SetResourceAccess(tfPlanResourceAccess)
			}

			if !tfPlanrequestBodyRequiredResourceAccess.ResourceAppId.Equal(requestBodyRequiredResourceAccessState.ResourceAppId) {
				tfPlanResourceAppId := tfPlanrequestBodyRequiredResourceAccess.ResourceAppId.ValueString()
				requestBodyRequiredResourceAccess.SetResourceAppId(&tfPlanResourceAppId)
			}
		}
		requestBody.SetRequiredResourceAccess(tfPlanRequiredResourceAccess)
	}

	if !tfPlan.SamlMetadataUrl.Equal(tfState.SamlMetadataUrl) {
		tfPlanSamlMetadataUrl := tfPlan.SamlMetadataUrl.ValueString()
		requestBody.SetSamlMetadataUrl(&tfPlanSamlMetadataUrl)
	}

	if !tfPlan.ServiceManagementReference.Equal(tfState.ServiceManagementReference) {
		tfPlanServiceManagementReference := tfPlan.ServiceManagementReference.ValueString()
		requestBody.SetServiceManagementReference(&tfPlanServiceManagementReference)
	}

	if !tfPlan.ServicePrincipalLockConfiguration.Equal(tfState.ServicePrincipalLockConfiguration) {
		requestBodyServicePrincipalLockConfiguration := models.NewServicePrincipalLockConfiguration()
		tfPlanrequestBodyServicePrincipalLockConfiguration := applicationServicePrincipalLockConfigurationModel{}
		tfPlan.ServicePrincipalLockConfiguration.As(ctx, &tfPlanrequestBodyServicePrincipalLockConfiguration, basetypes.ObjectAsOptions{})
		requestBodyServicePrincipalLockConfigurationState := applicationServicePrincipalLockConfigurationModel{}
		tfState.ServicePrincipalLockConfiguration.As(ctx, &requestBodyServicePrincipalLockConfigurationState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyServicePrincipalLockConfiguration.AllProperties.Equal(requestBodyServicePrincipalLockConfigurationState.AllProperties) {
			tfPlanAllProperties := tfPlanrequestBodyServicePrincipalLockConfiguration.AllProperties.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetAllProperties(&tfPlanAllProperties)
		}

		if !tfPlanrequestBodyServicePrincipalLockConfiguration.CredentialsWithUsageSign.Equal(requestBodyServicePrincipalLockConfigurationState.CredentialsWithUsageSign) {
			tfPlanCredentialsWithUsageSign := tfPlanrequestBodyServicePrincipalLockConfiguration.CredentialsWithUsageSign.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetCredentialsWithUsageSign(&tfPlanCredentialsWithUsageSign)
		}

		if !tfPlanrequestBodyServicePrincipalLockConfiguration.CredentialsWithUsageVerify.Equal(requestBodyServicePrincipalLockConfigurationState.CredentialsWithUsageVerify) {
			tfPlanCredentialsWithUsageVerify := tfPlanrequestBodyServicePrincipalLockConfiguration.CredentialsWithUsageVerify.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetCredentialsWithUsageVerify(&tfPlanCredentialsWithUsageVerify)
		}

		if !tfPlanrequestBodyServicePrincipalLockConfiguration.IsEnabled.Equal(requestBodyServicePrincipalLockConfigurationState.IsEnabled) {
			tfPlanIsEnabled := tfPlanrequestBodyServicePrincipalLockConfiguration.IsEnabled.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetIsEnabled(&tfPlanIsEnabled)
		}

		if !tfPlanrequestBodyServicePrincipalLockConfiguration.TokenEncryptionKeyId.Equal(requestBodyServicePrincipalLockConfigurationState.TokenEncryptionKeyId) {
			tfPlanTokenEncryptionKeyId := tfPlanrequestBodyServicePrincipalLockConfiguration.TokenEncryptionKeyId.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetTokenEncryptionKeyId(&tfPlanTokenEncryptionKeyId)
		}
		requestBody.SetServicePrincipalLockConfiguration(requestBodyServicePrincipalLockConfiguration)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyServicePrincipalLockConfiguration.AttributeTypes(), tfPlanrequestBodyServicePrincipalLockConfiguration)
		tfPlan.ServicePrincipalLockConfiguration = objectValue
	}

	if !tfPlan.SignInAudience.Equal(tfState.SignInAudience) {
		tfPlanSignInAudience := tfPlan.SignInAudience.ValueString()
		requestBody.SetSignInAudience(&tfPlanSignInAudience)
	}

	if !tfPlan.Spa.Equal(tfState.Spa) {
		requestBodySpa := models.NewSpaApplication()
		tfPlanrequestBodySpa := applicationSpaApplicationModel{}
		tfPlan.Spa.As(ctx, &tfPlanrequestBodySpa, basetypes.ObjectAsOptions{})
		requestBodySpaState := applicationSpaApplicationModel{}
		tfState.Spa.As(ctx, &requestBodySpaState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodySpa.RedirectUris.Equal(requestBodySpaState.RedirectUris) {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanrequestBodySpa.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodySpa.SetRedirectUris(stringArrayRedirectUris)
		}
		requestBody.SetSpa(requestBodySpa)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodySpa.AttributeTypes(), tfPlanrequestBodySpa)
		tfPlan.Spa = objectValue
	}

	if !tfPlan.Tags.Equal(tfState.Tags) {
		var stringArrayTags []string
		for _, i := range tfPlan.Tags.Elements() {
			stringArrayTags = append(stringArrayTags, i.String())
		}
		requestBody.SetTags(stringArrayTags)
	}

	if !tfPlan.TokenEncryptionKeyId.Equal(tfState.TokenEncryptionKeyId) {
		tfPlanTokenEncryptionKeyId := tfPlan.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		requestBody.SetTokenEncryptionKeyId(&u)
	}

	if !tfPlan.UniqueName.Equal(tfState.UniqueName) {
		tfPlanUniqueName := tfPlan.UniqueName.ValueString()
		requestBody.SetUniqueName(&tfPlanUniqueName)
	}

	if !tfPlan.VerifiedPublisher.Equal(tfState.VerifiedPublisher) {
		requestBodyVerifiedPublisher := models.NewVerifiedPublisher()
		tfPlanrequestBodyVerifiedPublisher := applicationVerifiedPublisherModel{}
		tfPlan.VerifiedPublisher.As(ctx, &tfPlanrequestBodyVerifiedPublisher, basetypes.ObjectAsOptions{})
		requestBodyVerifiedPublisherState := applicationVerifiedPublisherModel{}
		tfState.VerifiedPublisher.As(ctx, &requestBodyVerifiedPublisherState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyVerifiedPublisher.AddedDateTime.Equal(requestBodyVerifiedPublisherState.AddedDateTime) {
			tfPlanAddedDateTime := tfPlanrequestBodyVerifiedPublisher.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanAddedDateTime)
			requestBodyVerifiedPublisher.SetAddedDateTime(&t)
		}

		if !tfPlanrequestBodyVerifiedPublisher.DisplayName.Equal(requestBodyVerifiedPublisherState.DisplayName) {
			tfPlanDisplayName := tfPlanrequestBodyVerifiedPublisher.DisplayName.ValueString()
			requestBodyVerifiedPublisher.SetDisplayName(&tfPlanDisplayName)
		}

		if !tfPlanrequestBodyVerifiedPublisher.VerifiedPublisherId.Equal(requestBodyVerifiedPublisherState.VerifiedPublisherId) {
			tfPlanVerifiedPublisherId := tfPlanrequestBodyVerifiedPublisher.VerifiedPublisherId.ValueString()
			requestBodyVerifiedPublisher.SetVerifiedPublisherId(&tfPlanVerifiedPublisherId)
		}
		requestBody.SetVerifiedPublisher(requestBodyVerifiedPublisher)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyVerifiedPublisher.AttributeTypes(), tfPlanrequestBodyVerifiedPublisher)
		tfPlan.VerifiedPublisher = objectValue
	}

	if !tfPlan.Web.Equal(tfState.Web) {
		requestBodyWeb := models.NewWebApplication()
		tfPlanrequestBodyWeb := applicationWebApplicationModel{}
		tfPlan.Web.As(ctx, &tfPlanrequestBodyWeb, basetypes.ObjectAsOptions{})
		requestBodyWebState := applicationWebApplicationModel{}
		tfState.Web.As(ctx, &requestBodyWebState, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyWeb.HomePageUrl.Equal(requestBodyWebState.HomePageUrl) {
			tfPlanHomePageUrl := tfPlanrequestBodyWeb.HomePageUrl.ValueString()
			requestBodyWeb.SetHomePageUrl(&tfPlanHomePageUrl)
		}

		if !tfPlanrequestBodyWeb.ImplicitGrantSettings.Equal(requestBodyWebState.ImplicitGrantSettings) {
			requestBodyImplicitGrantSettings := models.NewImplicitGrantSettings()
			tfPlanrequestBodyImplicitGrantSettings := applicationImplicitGrantSettingsModel{}
			tfPlanrequestBodyWeb.ImplicitGrantSettings.As(ctx, &tfPlanrequestBodyImplicitGrantSettings, basetypes.ObjectAsOptions{})
			requestBodyImplicitGrantSettingsState := applicationImplicitGrantSettingsModel{}
			requestBodyWebState.ImplicitGrantSettings.As(ctx, &requestBodyImplicitGrantSettingsState, basetypes.ObjectAsOptions{})

			if !tfPlanrequestBodyImplicitGrantSettings.EnableAccessTokenIssuance.Equal(requestBodyImplicitGrantSettingsState.EnableAccessTokenIssuance) {
				tfPlanEnableAccessTokenIssuance := tfPlanrequestBodyImplicitGrantSettings.EnableAccessTokenIssuance.ValueBool()
				requestBodyImplicitGrantSettings.SetEnableAccessTokenIssuance(&tfPlanEnableAccessTokenIssuance)
			}

			if !tfPlanrequestBodyImplicitGrantSettings.EnableIdTokenIssuance.Equal(requestBodyImplicitGrantSettingsState.EnableIdTokenIssuance) {
				tfPlanEnableIdTokenIssuance := tfPlanrequestBodyImplicitGrantSettings.EnableIdTokenIssuance.ValueBool()
				requestBodyImplicitGrantSettings.SetEnableIdTokenIssuance(&tfPlanEnableIdTokenIssuance)
			}
			requestBodyWeb.SetImplicitGrantSettings(requestBodyImplicitGrantSettings)
			objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyImplicitGrantSettings.AttributeTypes(), tfPlanrequestBodyImplicitGrantSettings)
			tfPlanrequestBodyWeb.ImplicitGrantSettings = objectValue
		}

		if !tfPlanrequestBodyWeb.LogoutUrl.Equal(requestBodyWebState.LogoutUrl) {
			tfPlanLogoutUrl := tfPlanrequestBodyWeb.LogoutUrl.ValueString()
			requestBodyWeb.SetLogoutUrl(&tfPlanLogoutUrl)
		}

		if !tfPlanrequestBodyWeb.RedirectUriSettings.Equal(requestBodyWebState.RedirectUriSettings) {
			var tfPlanRedirectUriSettings []models.RedirectUriSettingsable
			for k, i := range tfPlanrequestBodyWeb.RedirectUriSettings.Elements() {
				requestBodyRedirectUriSettings := models.NewRedirectUriSettings()
				tfPlanrequestBodyRedirectUriSettings := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyRedirectUriSettings)
				requestBodyRedirectUriSettingsState := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, requestBodyWebState.RedirectUriSettings.Elements()[k].Type(ctx), &tfPlanrequestBodyRedirectUriSettings)

				if !tfPlanrequestBodyRedirectUriSettings.Uri.Equal(requestBodyRedirectUriSettingsState.Uri) {
					tfPlanUri := tfPlanrequestBodyRedirectUriSettings.Uri.ValueString()
					requestBodyRedirectUriSettings.SetUri(&tfPlanUri)
				}
			}
			requestBodyWeb.SetRedirectUriSettings(tfPlanRedirectUriSettings)
		}

		if !tfPlanrequestBodyWeb.RedirectUris.Equal(requestBodyWebState.RedirectUris) {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanrequestBodyWeb.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodyWeb.SetRedirectUris(stringArrayRedirectUris)
		}
		requestBody.SetWeb(requestBodyWeb)
		objectValue, _ := types.ObjectValueFrom(ctx, tfPlanrequestBodyWeb.AttributeTypes(), tfPlanrequestBodyWeb)
		tfPlan.Web = objectValue
	}

	// Update application
	_, err := r.client.Applications().ByApplicationId(tfState.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating application",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfState applicationModel
	diags := req.State.Get(ctx, &tfState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete application
	err := r.client.Applications().ByApplicationId(tfState.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting application",
			err.Error(),
		)
		return
	}

}
