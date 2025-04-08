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
	var tfStateApplication applicationModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfStateApplication)...)
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

	if !tfStateApplication.Id.IsNull() {
		result, err = d.client.Applications().ByApplicationId(tfStateApplication.Id.ValueString()).Get(context.Background(), &qparams)
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
		tfStateApplication.Id = types.StringValue(*result.GetId())
	} else {
		tfStateApplication.Id = types.StringNull()
	}
	if result.GetDeletedDateTime() != nil {
		tfStateApplication.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	} else {
		tfStateApplication.DeletedDateTime = types.StringNull()
	}
	if len(result.GetAddIns()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAddIns() {
			tfStateAddIn := applicationAddInModel{}

			if v.GetId() != nil {
				tfStateAddIn.Id = types.StringValue(v.GetId().String())
			} else {
				tfStateAddIn.Id = types.StringNull()
			}
			if len(v.GetProperties()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetProperties() {
					tfStateKeyValue := applicationKeyValueModel{}

					if v.GetKey() != nil {
						tfStateKeyValue.Key = types.StringValue(*v.GetKey())
					} else {
						tfStateKeyValue.Key = types.StringNull()
					}
					if v.GetValue() != nil {
						tfStateKeyValue.Value = types.StringValue(*v.GetValue())
					} else {
						tfStateKeyValue.Value = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateKeyValue.AttributeTypes(), tfStateKeyValue)
					objectValues = append(objectValues, objectValue)
				}
				tfStateAddIn.Properties, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetTypeEscaped() != nil {
				tfStateAddIn.Type = types.StringValue(*v.GetTypeEscaped())
			} else {
				tfStateAddIn.Type = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAddIn.AttributeTypes(), tfStateAddIn)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.AddIns, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetApi() != nil {
		tfStateApiApplication := applicationApiApplicationModel{}

		if result.GetApi().GetAcceptMappedClaims() != nil {
			tfStateApiApplication.AcceptMappedClaims = types.BoolValue(*result.GetApi().GetAcceptMappedClaims())
		} else {
			tfStateApiApplication.AcceptMappedClaims = types.BoolNull()
		}
		if len(result.GetApi().GetKnownClientApplications()) > 0 {
			var valueArrayKnownClientApplications []attr.Value
			for _, resultKnownClientApplications := range result.GetApi().GetKnownClientApplications() {
				valueArrayKnownClientApplications = append(valueArrayKnownClientApplications, types.StringValue(resultKnownClientApplications.String()))
			}
			tfStateApiApplication.KnownClientApplications, _ = types.ListValue(types.StringType, valueArrayKnownClientApplications)
		} else {
			tfStateApiApplication.KnownClientApplications = types.ListNull(types.StringType)
		}
		if len(result.GetApi().GetOauth2PermissionScopes()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetApi().GetOauth2PermissionScopes() {
				tfStatePermissionScope := applicationPermissionScopeModel{}

				if v.GetAdminConsentDescription() != nil {
					tfStatePermissionScope.AdminConsentDescription = types.StringValue(*v.GetAdminConsentDescription())
				} else {
					tfStatePermissionScope.AdminConsentDescription = types.StringNull()
				}
				if v.GetAdminConsentDisplayName() != nil {
					tfStatePermissionScope.AdminConsentDisplayName = types.StringValue(*v.GetAdminConsentDisplayName())
				} else {
					tfStatePermissionScope.AdminConsentDisplayName = types.StringNull()
				}
				if v.GetId() != nil {
					tfStatePermissionScope.Id = types.StringValue(v.GetId().String())
				} else {
					tfStatePermissionScope.Id = types.StringNull()
				}
				if v.GetIsEnabled() != nil {
					tfStatePermissionScope.IsEnabled = types.BoolValue(*v.GetIsEnabled())
				} else {
					tfStatePermissionScope.IsEnabled = types.BoolNull()
				}
				if v.GetOrigin() != nil {
					tfStatePermissionScope.Origin = types.StringValue(*v.GetOrigin())
				} else {
					tfStatePermissionScope.Origin = types.StringNull()
				}
				if v.GetTypeEscaped() != nil {
					tfStatePermissionScope.Type = types.StringValue(*v.GetTypeEscaped())
				} else {
					tfStatePermissionScope.Type = types.StringNull()
				}
				if v.GetUserConsentDescription() != nil {
					tfStatePermissionScope.UserConsentDescription = types.StringValue(*v.GetUserConsentDescription())
				} else {
					tfStatePermissionScope.UserConsentDescription = types.StringNull()
				}
				if v.GetUserConsentDisplayName() != nil {
					tfStatePermissionScope.UserConsentDisplayName = types.StringValue(*v.GetUserConsentDisplayName())
				} else {
					tfStatePermissionScope.UserConsentDisplayName = types.StringNull()
				}
				if v.GetValue() != nil {
					tfStatePermissionScope.Value = types.StringValue(*v.GetValue())
				} else {
					tfStatePermissionScope.Value = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStatePermissionScope.AttributeTypes(), tfStatePermissionScope)
				objectValues = append(objectValues, objectValue)
			}
			tfStateApiApplication.Oauth2PermissionScopes, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetApi().GetPreAuthorizedApplications()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetApi().GetPreAuthorizedApplications() {
				tfStatePreAuthorizedApplication := applicationPreAuthorizedApplicationModel{}

				if v.GetAppId() != nil {
					tfStatePreAuthorizedApplication.AppId = types.StringValue(*v.GetAppId())
				} else {
					tfStatePreAuthorizedApplication.AppId = types.StringNull()
				}
				if len(v.GetDelegatedPermissionIds()) > 0 {
					var valueArrayDelegatedPermissionIds []attr.Value
					for _, v := range v.GetDelegatedPermissionIds() {
						valueArrayDelegatedPermissionIds = append(valueArrayDelegatedPermissionIds, types.StringValue(v))
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
	if result.GetAppId() != nil {
		tfStateApplication.AppId = types.StringValue(*result.GetAppId())
	} else {
		tfStateApplication.AppId = types.StringNull()
	}
	if len(result.GetAppRoles()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAppRoles() {
			tfStateAppRole := applicationAppRoleModel{}

			if len(v.GetAllowedMemberTypes()) > 0 {
				var valueArrayAllowedMemberTypes []attr.Value
				for _, v := range v.GetAllowedMemberTypes() {
					valueArrayAllowedMemberTypes = append(valueArrayAllowedMemberTypes, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayAllowedMemberTypes)
				tfStateAppRole.AllowedMemberTypes = listValue
			} else {
				tfStateAppRole.AllowedMemberTypes = types.ListNull(types.StringType)
			}
			if v.GetDescription() != nil {
				tfStateAppRole.Description = types.StringValue(*v.GetDescription())
			} else {
				tfStateAppRole.Description = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				tfStateAppRole.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				tfStateAppRole.DisplayName = types.StringNull()
			}
			if v.GetId() != nil {
				tfStateAppRole.Id = types.StringValue(v.GetId().String())
			} else {
				tfStateAppRole.Id = types.StringNull()
			}
			if v.GetIsEnabled() != nil {
				tfStateAppRole.IsEnabled = types.BoolValue(*v.GetIsEnabled())
			} else {
				tfStateAppRole.IsEnabled = types.BoolNull()
			}
			if v.GetOrigin() != nil {
				tfStateAppRole.Origin = types.StringValue(*v.GetOrigin())
			} else {
				tfStateAppRole.Origin = types.StringNull()
			}
			if v.GetValue() != nil {
				tfStateAppRole.Value = types.StringValue(*v.GetValue())
			} else {
				tfStateAppRole.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAppRole.AttributeTypes(), tfStateAppRole)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.AppRoles, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetApplicationTemplateId() != nil {
		tfStateApplication.ApplicationTemplateId = types.StringValue(*result.GetApplicationTemplateId())
	} else {
		tfStateApplication.ApplicationTemplateId = types.StringNull()
	}
	if result.GetCertification() != nil {
		tfStateCertification := applicationCertificationModel{}

		if result.GetCertification().GetCertificationDetailsUrl() != nil {
			tfStateCertification.CertificationDetailsUrl = types.StringValue(*result.GetCertification().GetCertificationDetailsUrl())
		} else {
			tfStateCertification.CertificationDetailsUrl = types.StringNull()
		}
		if result.GetCertification().GetCertificationExpirationDateTime() != nil {
			tfStateCertification.CertificationExpirationDateTime = types.StringValue(result.GetCertification().GetCertificationExpirationDateTime().String())
		} else {
			tfStateCertification.CertificationExpirationDateTime = types.StringNull()
		}
		if result.GetCertification().GetIsCertifiedByMicrosoft() != nil {
			tfStateCertification.IsCertifiedByMicrosoft = types.BoolValue(*result.GetCertification().GetIsCertifiedByMicrosoft())
		} else {
			tfStateCertification.IsCertifiedByMicrosoft = types.BoolNull()
		}
		if result.GetCertification().GetIsPublisherAttested() != nil {
			tfStateCertification.IsPublisherAttested = types.BoolValue(*result.GetCertification().GetIsPublisherAttested())
		} else {
			tfStateCertification.IsPublisherAttested = types.BoolNull()
		}
		if result.GetCertification().GetLastCertificationDateTime() != nil {
			tfStateCertification.LastCertificationDateTime = types.StringValue(result.GetCertification().GetLastCertificationDateTime().String())
		} else {
			tfStateCertification.LastCertificationDateTime = types.StringNull()
		}

		tfStateApplication.Certification, _ = types.ObjectValueFrom(ctx, tfStateCertification.AttributeTypes(), tfStateCertification)
	}
	if result.GetCreatedDateTime() != nil {
		tfStateApplication.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		tfStateApplication.CreatedDateTime = types.StringNull()
	}
	if result.GetDefaultRedirectUri() != nil {
		tfStateApplication.DefaultRedirectUri = types.StringValue(*result.GetDefaultRedirectUri())
	} else {
		tfStateApplication.DefaultRedirectUri = types.StringNull()
	}
	if result.GetDescription() != nil {
		tfStateApplication.Description = types.StringValue(*result.GetDescription())
	} else {
		tfStateApplication.Description = types.StringNull()
	}
	if result.GetDisabledByMicrosoftStatus() != nil {
		tfStateApplication.DisabledByMicrosoftStatus = types.StringValue(*result.GetDisabledByMicrosoftStatus())
	} else {
		tfStateApplication.DisabledByMicrosoftStatus = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		tfStateApplication.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		tfStateApplication.DisplayName = types.StringNull()
	}
	if result.GetGroupMembershipClaims() != nil {
		tfStateApplication.GroupMembershipClaims = types.StringValue(*result.GetGroupMembershipClaims())
	} else {
		tfStateApplication.GroupMembershipClaims = types.StringNull()
	}
	if len(result.GetIdentifierUris()) > 0 {
		var valueArrayIdentifierUris []attr.Value
		for _, v := range result.GetIdentifierUris() {
			valueArrayIdentifierUris = append(valueArrayIdentifierUris, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayIdentifierUris)
		tfStateApplication.IdentifierUris = listValue
	} else {
		tfStateApplication.IdentifierUris = types.ListNull(types.StringType)
	}
	if result.GetInfo() != nil {
		tfStateInformationalUrl := applicationInformationalUrlModel{}

		if result.GetInfo().GetLogoUrl() != nil {
			tfStateInformationalUrl.LogoUrl = types.StringValue(*result.GetInfo().GetLogoUrl())
		} else {
			tfStateInformationalUrl.LogoUrl = types.StringNull()
		}
		if result.GetInfo().GetMarketingUrl() != nil {
			tfStateInformationalUrl.MarketingUrl = types.StringValue(*result.GetInfo().GetMarketingUrl())
		} else {
			tfStateInformationalUrl.MarketingUrl = types.StringNull()
		}
		if result.GetInfo().GetPrivacyStatementUrl() != nil {
			tfStateInformationalUrl.PrivacyStatementUrl = types.StringValue(*result.GetInfo().GetPrivacyStatementUrl())
		} else {
			tfStateInformationalUrl.PrivacyStatementUrl = types.StringNull()
		}
		if result.GetInfo().GetSupportUrl() != nil {
			tfStateInformationalUrl.SupportUrl = types.StringValue(*result.GetInfo().GetSupportUrl())
		} else {
			tfStateInformationalUrl.SupportUrl = types.StringNull()
		}
		if result.GetInfo().GetTermsOfServiceUrl() != nil {
			tfStateInformationalUrl.TermsOfServiceUrl = types.StringValue(*result.GetInfo().GetTermsOfServiceUrl())
		} else {
			tfStateInformationalUrl.TermsOfServiceUrl = types.StringNull()
		}

		tfStateApplication.Info, _ = types.ObjectValueFrom(ctx, tfStateInformationalUrl.AttributeTypes(), tfStateInformationalUrl)
	}
	if result.GetIsDeviceOnlyAuthSupported() != nil {
		tfStateApplication.IsDeviceOnlyAuthSupported = types.BoolValue(*result.GetIsDeviceOnlyAuthSupported())
	} else {
		tfStateApplication.IsDeviceOnlyAuthSupported = types.BoolNull()
	}
	if result.GetIsFallbackPublicClient() != nil {
		tfStateApplication.IsFallbackPublicClient = types.BoolValue(*result.GetIsFallbackPublicClient())
	} else {
		tfStateApplication.IsFallbackPublicClient = types.BoolNull()
	}
	if len(result.GetKeyCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetKeyCredentials() {
			tfStateKeyCredential := applicationKeyCredentialModel{}

			if v.GetCustomKeyIdentifier() != nil {
				tfStateKeyCredential.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
			} else {
				tfStateKeyCredential.CustomKeyIdentifier = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				tfStateKeyCredential.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				tfStateKeyCredential.DisplayName = types.StringNull()
			}
			if v.GetEndDateTime() != nil {
				tfStateKeyCredential.EndDateTime = types.StringValue(v.GetEndDateTime().String())
			} else {
				tfStateKeyCredential.EndDateTime = types.StringNull()
			}
			if v.GetKey() != nil {
				tfStateKeyCredential.Key = types.StringValue(string(v.GetKey()[:]))
			} else {
				tfStateKeyCredential.Key = types.StringNull()
			}
			if v.GetKeyId() != nil {
				tfStateKeyCredential.KeyId = types.StringValue(v.GetKeyId().String())
			} else {
				tfStateKeyCredential.KeyId = types.StringNull()
			}
			if v.GetStartDateTime() != nil {
				tfStateKeyCredential.StartDateTime = types.StringValue(v.GetStartDateTime().String())
			} else {
				tfStateKeyCredential.StartDateTime = types.StringNull()
			}
			if v.GetTypeEscaped() != nil {
				tfStateKeyCredential.Type = types.StringValue(*v.GetTypeEscaped())
			} else {
				tfStateKeyCredential.Type = types.StringNull()
			}
			if v.GetUsage() != nil {
				tfStateKeyCredential.Usage = types.StringValue(*v.GetUsage())
			} else {
				tfStateKeyCredential.Usage = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateKeyCredential.AttributeTypes(), tfStateKeyCredential)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.KeyCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetLogo() != nil {
		tfStateApplication.Logo = types.StringValue(string(result.GetLogo()[:]))
	} else {
		tfStateApplication.Logo = types.StringNull()
	}
	if result.GetNativeAuthenticationApisEnabled() != nil {
		tfStateApplication.NativeAuthenticationApisEnabled = types.StringValue(result.GetNativeAuthenticationApisEnabled().String())
	} else {
		tfStateApplication.NativeAuthenticationApisEnabled = types.StringNull()
	}
	if result.GetNotes() != nil {
		tfStateApplication.Notes = types.StringValue(*result.GetNotes())
	} else {
		tfStateApplication.Notes = types.StringNull()
	}
	if result.GetOauth2RequirePostResponse() != nil {
		tfStateApplication.Oauth2RequirePostResponse = types.BoolValue(*result.GetOauth2RequirePostResponse())
	} else {
		tfStateApplication.Oauth2RequirePostResponse = types.BoolNull()
	}
	if result.GetOptionalClaims() != nil {
		tfStateOptionalClaims := applicationOptionalClaimsModel{}

		if len(result.GetOptionalClaims().GetAccessToken()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetOptionalClaims().GetAccessToken() {
				tfStateOptionalClaim := applicationOptionalClaimModel{}

				if len(v.GetAdditionalProperties()) > 0 {
					var valueArrayAdditionalProperties []attr.Value
					for _, v := range v.GetAdditionalProperties() {
						valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
					tfStateOptionalClaim.AdditionalProperties = listValue
				} else {
					tfStateOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				if v.GetEssential() != nil {
					tfStateOptionalClaim.Essential = types.BoolValue(*v.GetEssential())
				} else {
					tfStateOptionalClaim.Essential = types.BoolNull()
				}
				if v.GetName() != nil {
					tfStateOptionalClaim.Name = types.StringValue(*v.GetName())
				} else {
					tfStateOptionalClaim.Name = types.StringNull()
				}
				if v.GetSource() != nil {
					tfStateOptionalClaim.Source = types.StringValue(*v.GetSource())
				} else {
					tfStateOptionalClaim.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStateOptionalClaim.AttributeTypes(), tfStateOptionalClaim)
				objectValues = append(objectValues, objectValue)
			}
			tfStateOptionalClaims.AccessToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetOptionalClaims().GetIdToken()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetOptionalClaims().GetIdToken() {
				tfStateOptionalClaim := applicationOptionalClaimModel{}

				if len(v.GetAdditionalProperties()) > 0 {
					var valueArrayAdditionalProperties []attr.Value
					for _, v := range v.GetAdditionalProperties() {
						valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
					tfStateOptionalClaim.AdditionalProperties = listValue
				} else {
					tfStateOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				if v.GetEssential() != nil {
					tfStateOptionalClaim.Essential = types.BoolValue(*v.GetEssential())
				} else {
					tfStateOptionalClaim.Essential = types.BoolNull()
				}
				if v.GetName() != nil {
					tfStateOptionalClaim.Name = types.StringValue(*v.GetName())
				} else {
					tfStateOptionalClaim.Name = types.StringNull()
				}
				if v.GetSource() != nil {
					tfStateOptionalClaim.Source = types.StringValue(*v.GetSource())
				} else {
					tfStateOptionalClaim.Source = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStateOptionalClaim.AttributeTypes(), tfStateOptionalClaim)
				objectValues = append(objectValues, objectValue)
			}
			tfStateOptionalClaims.IdToken, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetOptionalClaims().GetSaml2Token()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetOptionalClaims().GetSaml2Token() {
				tfStateOptionalClaim := applicationOptionalClaimModel{}

				if len(v.GetAdditionalProperties()) > 0 {
					var valueArrayAdditionalProperties []attr.Value
					for _, v := range v.GetAdditionalProperties() {
						valueArrayAdditionalProperties = append(valueArrayAdditionalProperties, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayAdditionalProperties)
					tfStateOptionalClaim.AdditionalProperties = listValue
				} else {
					tfStateOptionalClaim.AdditionalProperties = types.ListNull(types.StringType)
				}
				if v.GetEssential() != nil {
					tfStateOptionalClaim.Essential = types.BoolValue(*v.GetEssential())
				} else {
					tfStateOptionalClaim.Essential = types.BoolNull()
				}
				if v.GetName() != nil {
					tfStateOptionalClaim.Name = types.StringValue(*v.GetName())
				} else {
					tfStateOptionalClaim.Name = types.StringNull()
				}
				if v.GetSource() != nil {
					tfStateOptionalClaim.Source = types.StringValue(*v.GetSource())
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
	if result.GetParentalControlSettings() != nil {
		tfStateParentalControlSettings := applicationParentalControlSettingsModel{}

		if len(result.GetParentalControlSettings().GetCountriesBlockedForMinors()) > 0 {
			var valueArrayCountriesBlockedForMinors []attr.Value
			for _, v := range result.GetParentalControlSettings().GetCountriesBlockedForMinors() {
				valueArrayCountriesBlockedForMinors = append(valueArrayCountriesBlockedForMinors, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, valueArrayCountriesBlockedForMinors)
			tfStateParentalControlSettings.CountriesBlockedForMinors = listValue
		} else {
			tfStateParentalControlSettings.CountriesBlockedForMinors = types.ListNull(types.StringType)
		}
		if result.GetParentalControlSettings().GetLegalAgeGroupRule() != nil {
			tfStateParentalControlSettings.LegalAgeGroupRule = types.StringValue(*result.GetParentalControlSettings().GetLegalAgeGroupRule())
		} else {
			tfStateParentalControlSettings.LegalAgeGroupRule = types.StringNull()
		}

		tfStateApplication.ParentalControlSettings, _ = types.ObjectValueFrom(ctx, tfStateParentalControlSettings.AttributeTypes(), tfStateParentalControlSettings)
	}
	if len(result.GetPasswordCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetPasswordCredentials() {
			tfStatePasswordCredential := applicationPasswordCredentialModel{}

			if v.GetCustomKeyIdentifier() != nil {
				tfStatePasswordCredential.CustomKeyIdentifier = types.StringValue(string(v.GetCustomKeyIdentifier()[:]))
			} else {
				tfStatePasswordCredential.CustomKeyIdentifier = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				tfStatePasswordCredential.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				tfStatePasswordCredential.DisplayName = types.StringNull()
			}
			if v.GetEndDateTime() != nil {
				tfStatePasswordCredential.EndDateTime = types.StringValue(v.GetEndDateTime().String())
			} else {
				tfStatePasswordCredential.EndDateTime = types.StringNull()
			}
			if v.GetHint() != nil {
				tfStatePasswordCredential.Hint = types.StringValue(*v.GetHint())
			} else {
				tfStatePasswordCredential.Hint = types.StringNull()
			}
			if v.GetKeyId() != nil {
				tfStatePasswordCredential.KeyId = types.StringValue(v.GetKeyId().String())
			} else {
				tfStatePasswordCredential.KeyId = types.StringNull()
			}
			if v.GetSecretText() != nil {
				tfStatePasswordCredential.SecretText = types.StringValue(*v.GetSecretText())
			} else {
				tfStatePasswordCredential.SecretText = types.StringNull()
			}
			if v.GetStartDateTime() != nil {
				tfStatePasswordCredential.StartDateTime = types.StringValue(v.GetStartDateTime().String())
			} else {
				tfStatePasswordCredential.StartDateTime = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStatePasswordCredential.AttributeTypes(), tfStatePasswordCredential)
			objectValues = append(objectValues, objectValue)
		}
		tfStateApplication.PasswordCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetPublicClient() != nil {
		tfStatePublicClientApplication := applicationPublicClientApplicationModel{}

		if len(result.GetPublicClient().GetRedirectUris()) > 0 {
			var valueArrayRedirectUris []attr.Value
			for _, v := range result.GetPublicClient().GetRedirectUris() {
				valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
			tfStatePublicClientApplication.RedirectUris = listValue
		} else {
			tfStatePublicClientApplication.RedirectUris = types.ListNull(types.StringType)
		}

		tfStateApplication.PublicClient, _ = types.ObjectValueFrom(ctx, tfStatePublicClientApplication.AttributeTypes(), tfStatePublicClientApplication)
	}
	if result.GetPublisherDomain() != nil {
		tfStateApplication.PublisherDomain = types.StringValue(*result.GetPublisherDomain())
	} else {
		tfStateApplication.PublisherDomain = types.StringNull()
	}
	if result.GetRequestSignatureVerification() != nil {
		tfStateRequestSignatureVerification := applicationRequestSignatureVerificationModel{}

		if result.GetRequestSignatureVerification().GetAllowedWeakAlgorithms() != nil {
			tfStateRequestSignatureVerification.AllowedWeakAlgorithms = types.StringValue(result.GetRequestSignatureVerification().GetAllowedWeakAlgorithms().String())
		} else {
			tfStateRequestSignatureVerification.AllowedWeakAlgorithms = types.StringNull()
		}
		if result.GetRequestSignatureVerification().GetIsSignedRequestRequired() != nil {
			tfStateRequestSignatureVerification.IsSignedRequestRequired = types.BoolValue(*result.GetRequestSignatureVerification().GetIsSignedRequestRequired())
		} else {
			tfStateRequestSignatureVerification.IsSignedRequestRequired = types.BoolNull()
		}

		tfStateApplication.RequestSignatureVerification, _ = types.ObjectValueFrom(ctx, tfStateRequestSignatureVerification.AttributeTypes(), tfStateRequestSignatureVerification)
	}
	if len(result.GetRequiredResourceAccess()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetRequiredResourceAccess() {
			tfStateRequiredResourceAccess := applicationRequiredResourceAccessModel{}

			if len(v.GetResourceAccess()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetResourceAccess() {
					tfStateResourceAccess := applicationResourceAccessModel{}

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
		tfStateApplication.RequiredResourceAccess, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetSamlMetadataUrl() != nil {
		tfStateApplication.SamlMetadataUrl = types.StringValue(*result.GetSamlMetadataUrl())
	} else {
		tfStateApplication.SamlMetadataUrl = types.StringNull()
	}
	if result.GetServiceManagementReference() != nil {
		tfStateApplication.ServiceManagementReference = types.StringValue(*result.GetServiceManagementReference())
	} else {
		tfStateApplication.ServiceManagementReference = types.StringNull()
	}
	if result.GetServicePrincipalLockConfiguration() != nil {
		tfStateServicePrincipalLockConfiguration := applicationServicePrincipalLockConfigurationModel{}

		if result.GetServicePrincipalLockConfiguration().GetAllProperties() != nil {
			tfStateServicePrincipalLockConfiguration.AllProperties = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetAllProperties())
		} else {
			tfStateServicePrincipalLockConfiguration.AllProperties = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign() != nil {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageSign())
		} else {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify() != nil {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetCredentialsWithUsageVerify())
		} else {
			tfStateServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetIsEnabled() != nil {
			tfStateServicePrincipalLockConfiguration.IsEnabled = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetIsEnabled())
		} else {
			tfStateServicePrincipalLockConfiguration.IsEnabled = types.BoolNull()
		}
		if result.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId() != nil {
			tfStateServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolValue(*result.GetServicePrincipalLockConfiguration().GetTokenEncryptionKeyId())
		} else {
			tfStateServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolNull()
		}

		tfStateApplication.ServicePrincipalLockConfiguration, _ = types.ObjectValueFrom(ctx, tfStateServicePrincipalLockConfiguration.AttributeTypes(), tfStateServicePrincipalLockConfiguration)
	}
	if result.GetSignInAudience() != nil {
		tfStateApplication.SignInAudience = types.StringValue(*result.GetSignInAudience())
	} else {
		tfStateApplication.SignInAudience = types.StringNull()
	}
	if result.GetSpa() != nil {
		tfStateSpaApplication := applicationSpaApplicationModel{}

		if len(result.GetSpa().GetRedirectUris()) > 0 {
			var valueArrayRedirectUris []attr.Value
			for _, v := range result.GetSpa().GetRedirectUris() {
				valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, valueArrayRedirectUris)
			tfStateSpaApplication.RedirectUris = listValue
		} else {
			tfStateSpaApplication.RedirectUris = types.ListNull(types.StringType)
		}

		tfStateApplication.Spa, _ = types.ObjectValueFrom(ctx, tfStateSpaApplication.AttributeTypes(), tfStateSpaApplication)
	}
	if len(result.GetTags()) > 0 {
		var valueArrayTags []attr.Value
		for _, v := range result.GetTags() {
			valueArrayTags = append(valueArrayTags, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayTags)
		tfStateApplication.Tags = listValue
	} else {
		tfStateApplication.Tags = types.ListNull(types.StringType)
	}
	if result.GetTokenEncryptionKeyId() != nil {
		tfStateApplication.TokenEncryptionKeyId = types.StringValue(result.GetTokenEncryptionKeyId().String())
	} else {
		tfStateApplication.TokenEncryptionKeyId = types.StringNull()
	}
	if result.GetUniqueName() != nil {
		tfStateApplication.UniqueName = types.StringValue(*result.GetUniqueName())
	} else {
		tfStateApplication.UniqueName = types.StringNull()
	}
	if result.GetVerifiedPublisher() != nil {
		tfStateVerifiedPublisher := applicationVerifiedPublisherModel{}

		if result.GetVerifiedPublisher().GetAddedDateTime() != nil {
			tfStateVerifiedPublisher.AddedDateTime = types.StringValue(result.GetVerifiedPublisher().GetAddedDateTime().String())
		} else {
			tfStateVerifiedPublisher.AddedDateTime = types.StringNull()
		}
		if result.GetVerifiedPublisher().GetDisplayName() != nil {
			tfStateVerifiedPublisher.DisplayName = types.StringValue(*result.GetVerifiedPublisher().GetDisplayName())
		} else {
			tfStateVerifiedPublisher.DisplayName = types.StringNull()
		}
		if result.GetVerifiedPublisher().GetVerifiedPublisherId() != nil {
			tfStateVerifiedPublisher.VerifiedPublisherId = types.StringValue(*result.GetVerifiedPublisher().GetVerifiedPublisherId())
		} else {
			tfStateVerifiedPublisher.VerifiedPublisherId = types.StringNull()
		}

		tfStateApplication.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfStateVerifiedPublisher.AttributeTypes(), tfStateVerifiedPublisher)
	}
	if result.GetWeb() != nil {
		tfStateWebApplication := applicationWebApplicationModel{}

		if result.GetWeb().GetHomePageUrl() != nil {
			tfStateWebApplication.HomePageUrl = types.StringValue(*result.GetWeb().GetHomePageUrl())
		} else {
			tfStateWebApplication.HomePageUrl = types.StringNull()
		}
		if result.GetWeb().GetImplicitGrantSettings() != nil {
			tfStateImplicitGrantSettings := applicationImplicitGrantSettingsModel{}

			if result.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance() != nil {
				tfStateImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolValue(*result.GetWeb().GetImplicitGrantSettings().GetEnableAccessTokenIssuance())
			} else {
				tfStateImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolNull()
			}
			if result.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance() != nil {
				tfStateImplicitGrantSettings.EnableIdTokenIssuance = types.BoolValue(*result.GetWeb().GetImplicitGrantSettings().GetEnableIdTokenIssuance())
			} else {
				tfStateImplicitGrantSettings.EnableIdTokenIssuance = types.BoolNull()
			}

			tfStateWebApplication.ImplicitGrantSettings, _ = types.ObjectValueFrom(ctx, tfStateImplicitGrantSettings.AttributeTypes(), tfStateImplicitGrantSettings)
		}
		if result.GetWeb().GetLogoutUrl() != nil {
			tfStateWebApplication.LogoutUrl = types.StringValue(*result.GetWeb().GetLogoutUrl())
		} else {
			tfStateWebApplication.LogoutUrl = types.StringNull()
		}
		if len(result.GetWeb().GetRedirectUriSettings()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetWeb().GetRedirectUriSettings() {
				tfStateRedirectUriSettings := applicationRedirectUriSettingsModel{}

				if v.GetUri() != nil {
					tfStateRedirectUriSettings.Uri = types.StringValue(*v.GetUri())
				} else {
					tfStateRedirectUriSettings.Uri = types.StringNull()
				}
				objectValue, _ := types.ObjectValueFrom(ctx, tfStateRedirectUriSettings.AttributeTypes(), tfStateRedirectUriSettings)
				objectValues = append(objectValues, objectValue)
			}
			tfStateWebApplication.RedirectUriSettings, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
		}
		if len(result.GetWeb().GetRedirectUris()) > 0 {
			var valueArrayRedirectUris []attr.Value
			for _, v := range result.GetWeb().GetRedirectUris() {
				valueArrayRedirectUris = append(valueArrayRedirectUris, types.StringValue(v))
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

// Update updates the resource and sets the updated Terraform state on success.
func (r *applicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanApplication applicationModel
	diags := req.Plan.Get(ctx, &tfPlanApplication)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfStateApplication applicationModel
	diags = req.State.Get(ctx, &tfStateApplication)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBodyApplication := models.NewApplication()

	if !tfPlanApplication.Id.Equal(tfStateApplication.Id) {
		tfPlanId := tfPlanApplication.Id.ValueString()
		requestBodyApplication.SetId(&tfPlanId)
	}

	if !tfPlanApplication.DeletedDateTime.Equal(tfStateApplication.DeletedDateTime) {
		tfPlanDeletedDateTime := tfPlanApplication.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyApplication.SetDeletedDateTime(&t)
	}

	if !tfPlanApplication.AddIns.Equal(tfStateApplication.AddIns) {
		var tfPlanAddIns []models.AddInable
		for k, i := range tfPlanApplication.AddIns.Elements() {
			requestBodyAddIn := models.NewAddIn()
			tfPlanAddIn := applicationAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAddIn)
			tfStateAddIn := applicationAddInModel{}
			types.ListValueFrom(ctx, tfStateApplication.AddIns.Elements()[k].Type(ctx), &tfPlanAddIn)

			if !tfPlanAddIn.Id.Equal(tfStateAddIn.Id) {
				tfPlanId := tfPlanAddIn.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAddIn.SetId(&u)
			}

			if !tfPlanAddIn.Properties.Equal(tfStateAddIn.Properties) {
				var tfPlanProperties []models.KeyValueable
				for k, i := range tfPlanAddIn.Properties.Elements() {
					requestBodyKeyValue := models.NewKeyValue()
					tfPlanKeyValue := applicationKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyValue)
					tfStateKeyValue := applicationKeyValueModel{}
					types.ListValueFrom(ctx, tfStateAddIn.Properties.Elements()[k].Type(ctx), &tfPlanKeyValue)

					if !tfPlanKeyValue.Key.Equal(tfStateKeyValue.Key) {
						tfPlanKey := tfPlanKeyValue.Key.ValueString()
						requestBodyKeyValue.SetKey(&tfPlanKey)
					}

					if !tfPlanKeyValue.Value.Equal(tfStateKeyValue.Value) {
						tfPlanValue := tfPlanKeyValue.Value.ValueString()
						requestBodyKeyValue.SetValue(&tfPlanValue)
					}
				}
				requestBodyAddIn.SetProperties(tfPlanProperties)
			}

			if !tfPlanAddIn.Type.Equal(tfStateAddIn.Type) {
				tfPlanType := tfPlanAddIn.Type.ValueString()
				requestBodyAddIn.SetTypeEscaped(&tfPlanType)
			}
		}
		requestBodyApplication.SetAddIns(tfPlanAddIns)
	}

	if !tfPlanApplication.Api.Equal(tfStateApplication.Api) {
		requestBodyApiApplication := models.NewApiApplication()
		tfPlanApiApplication := applicationApiApplicationModel{}
		tfPlanApplication.Api.As(ctx, &tfPlanApiApplication, basetypes.ObjectAsOptions{})
		tfStateApiApplication := applicationApiApplicationModel{}
		tfStateApplication.Api.As(ctx, &tfStateApiApplication, basetypes.ObjectAsOptions{})

		if !tfPlanApiApplication.AcceptMappedClaims.Equal(tfStateApiApplication.AcceptMappedClaims) {
			tfPlanAcceptMappedClaims := tfPlanApiApplication.AcceptMappedClaims.ValueBool()
			requestBodyApiApplication.SetAcceptMappedClaims(&tfPlanAcceptMappedClaims)
		}

		if !tfPlanApiApplication.KnownClientApplications.Equal(tfStateApiApplication.KnownClientApplications) {
			var KnownClientApplications []uuid.UUID
			for _, i := range tfPlanApiApplication.KnownClientApplications.Elements() {
				u, _ := uuid.Parse(i.String())
				KnownClientApplications = append(KnownClientApplications, u)
			}
			requestBodyApiApplication.SetKnownClientApplications(KnownClientApplications)
		}

		if !tfPlanApiApplication.Oauth2PermissionScopes.Equal(tfStateApiApplication.Oauth2PermissionScopes) {
			var tfPlanOauth2PermissionScopes []models.PermissionScopeable
			for k, i := range tfPlanApiApplication.Oauth2PermissionScopes.Elements() {
				requestBodyPermissionScope := models.NewPermissionScope()
				tfPlanPermissionScope := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPermissionScope)
				tfStatePermissionScope := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, tfStateApiApplication.Oauth2PermissionScopes.Elements()[k].Type(ctx), &tfPlanPermissionScope)

				if !tfPlanPermissionScope.AdminConsentDescription.Equal(tfStatePermissionScope.AdminConsentDescription) {
					tfPlanAdminConsentDescription := tfPlanPermissionScope.AdminConsentDescription.ValueString()
					requestBodyPermissionScope.SetAdminConsentDescription(&tfPlanAdminConsentDescription)
				}

				if !tfPlanPermissionScope.AdminConsentDisplayName.Equal(tfStatePermissionScope.AdminConsentDisplayName) {
					tfPlanAdminConsentDisplayName := tfPlanPermissionScope.AdminConsentDisplayName.ValueString()
					requestBodyPermissionScope.SetAdminConsentDisplayName(&tfPlanAdminConsentDisplayName)
				}

				if !tfPlanPermissionScope.Id.Equal(tfStatePermissionScope.Id) {
					tfPlanId := tfPlanPermissionScope.Id.ValueString()
					u, _ := uuid.Parse(tfPlanId)
					requestBodyPermissionScope.SetId(&u)
				}

				if !tfPlanPermissionScope.IsEnabled.Equal(tfStatePermissionScope.IsEnabled) {
					tfPlanIsEnabled := tfPlanPermissionScope.IsEnabled.ValueBool()
					requestBodyPermissionScope.SetIsEnabled(&tfPlanIsEnabled)
				}

				if !tfPlanPermissionScope.Origin.Equal(tfStatePermissionScope.Origin) {
					tfPlanOrigin := tfPlanPermissionScope.Origin.ValueString()
					requestBodyPermissionScope.SetOrigin(&tfPlanOrigin)
				}

				if !tfPlanPermissionScope.Type.Equal(tfStatePermissionScope.Type) {
					tfPlanType := tfPlanPermissionScope.Type.ValueString()
					requestBodyPermissionScope.SetTypeEscaped(&tfPlanType)
				}

				if !tfPlanPermissionScope.UserConsentDescription.Equal(tfStatePermissionScope.UserConsentDescription) {
					tfPlanUserConsentDescription := tfPlanPermissionScope.UserConsentDescription.ValueString()
					requestBodyPermissionScope.SetUserConsentDescription(&tfPlanUserConsentDescription)
				}

				if !tfPlanPermissionScope.UserConsentDisplayName.Equal(tfStatePermissionScope.UserConsentDisplayName) {
					tfPlanUserConsentDisplayName := tfPlanPermissionScope.UserConsentDisplayName.ValueString()
					requestBodyPermissionScope.SetUserConsentDisplayName(&tfPlanUserConsentDisplayName)
				}

				if !tfPlanPermissionScope.Value.Equal(tfStatePermissionScope.Value) {
					tfPlanValue := tfPlanPermissionScope.Value.ValueString()
					requestBodyPermissionScope.SetValue(&tfPlanValue)
				}
			}
			requestBodyApiApplication.SetOauth2PermissionScopes(tfPlanOauth2PermissionScopes)
		}

		if !tfPlanApiApplication.PreAuthorizedApplications.Equal(tfStateApiApplication.PreAuthorizedApplications) {
			var tfPlanPreAuthorizedApplications []models.PreAuthorizedApplicationable
			for k, i := range tfPlanApiApplication.PreAuthorizedApplications.Elements() {
				requestBodyPreAuthorizedApplication := models.NewPreAuthorizedApplication()
				tfPlanPreAuthorizedApplication := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPreAuthorizedApplication)
				tfStatePreAuthorizedApplication := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, tfStateApiApplication.PreAuthorizedApplications.Elements()[k].Type(ctx), &tfPlanPreAuthorizedApplication)

				if !tfPlanPreAuthorizedApplication.AppId.Equal(tfStatePreAuthorizedApplication.AppId) {
					tfPlanAppId := tfPlanPreAuthorizedApplication.AppId.ValueString()
					requestBodyPreAuthorizedApplication.SetAppId(&tfPlanAppId)
				}

				if !tfPlanPreAuthorizedApplication.DelegatedPermissionIds.Equal(tfStatePreAuthorizedApplication.DelegatedPermissionIds) {
					var stringArrayDelegatedPermissionIds []string
					for _, i := range tfPlanPreAuthorizedApplication.DelegatedPermissionIds.Elements() {
						stringArrayDelegatedPermissionIds = append(stringArrayDelegatedPermissionIds, i.String())
					}
					requestBodyPreAuthorizedApplication.SetDelegatedPermissionIds(stringArrayDelegatedPermissionIds)
				}
			}
			requestBodyApiApplication.SetPreAuthorizedApplications(tfPlanPreAuthorizedApplications)
		}
		requestBodyApplication.SetApi(requestBodyApiApplication)
		tfPlanApplication.Api, _ = types.ObjectValueFrom(ctx, tfPlanApiApplication.AttributeTypes(), tfPlanApiApplication)
	}

	if !tfPlanApplication.AppId.Equal(tfStateApplication.AppId) {
		tfPlanAppId := tfPlanApplication.AppId.ValueString()
		requestBodyApplication.SetAppId(&tfPlanAppId)
	}

	if !tfPlanApplication.AppRoles.Equal(tfStateApplication.AppRoles) {
		var tfPlanAppRoles []models.AppRoleable
		for k, i := range tfPlanApplication.AppRoles.Elements() {
			requestBodyAppRole := models.NewAppRole()
			tfPlanAppRole := applicationAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAppRole)
			tfStateAppRole := applicationAppRoleModel{}
			types.ListValueFrom(ctx, tfStateApplication.AppRoles.Elements()[k].Type(ctx), &tfPlanAppRole)

			if !tfPlanAppRole.AllowedMemberTypes.Equal(tfStateAppRole.AllowedMemberTypes) {
				var stringArrayAllowedMemberTypes []string
				for _, i := range tfPlanAppRole.AllowedMemberTypes.Elements() {
					stringArrayAllowedMemberTypes = append(stringArrayAllowedMemberTypes, i.String())
				}
				requestBodyAppRole.SetAllowedMemberTypes(stringArrayAllowedMemberTypes)
			}

			if !tfPlanAppRole.Description.Equal(tfStateAppRole.Description) {
				tfPlanDescription := tfPlanAppRole.Description.ValueString()
				requestBodyAppRole.SetDescription(&tfPlanDescription)
			}

			if !tfPlanAppRole.DisplayName.Equal(tfStateAppRole.DisplayName) {
				tfPlanDisplayName := tfPlanAppRole.DisplayName.ValueString()
				requestBodyAppRole.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanAppRole.Id.Equal(tfStateAppRole.Id) {
				tfPlanId := tfPlanAppRole.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAppRole.SetId(&u)
			}

			if !tfPlanAppRole.IsEnabled.Equal(tfStateAppRole.IsEnabled) {
				tfPlanIsEnabled := tfPlanAppRole.IsEnabled.ValueBool()
				requestBodyAppRole.SetIsEnabled(&tfPlanIsEnabled)
			}

			if !tfPlanAppRole.Origin.Equal(tfStateAppRole.Origin) {
				tfPlanOrigin := tfPlanAppRole.Origin.ValueString()
				requestBodyAppRole.SetOrigin(&tfPlanOrigin)
			}

			if !tfPlanAppRole.Value.Equal(tfStateAppRole.Value) {
				tfPlanValue := tfPlanAppRole.Value.ValueString()
				requestBodyAppRole.SetValue(&tfPlanValue)
			}
		}
		requestBodyApplication.SetAppRoles(tfPlanAppRoles)
	}

	if !tfPlanApplication.ApplicationTemplateId.Equal(tfStateApplication.ApplicationTemplateId) {
		tfPlanApplicationTemplateId := tfPlanApplication.ApplicationTemplateId.ValueString()
		requestBodyApplication.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	}

	if !tfPlanApplication.Certification.Equal(tfStateApplication.Certification) {
		requestBodyCertification := models.NewCertification()
		tfPlanCertification := applicationCertificationModel{}
		tfPlanApplication.Certification.As(ctx, &tfPlanCertification, basetypes.ObjectAsOptions{})
		tfStateCertification := applicationCertificationModel{}
		tfStateApplication.Certification.As(ctx, &tfStateCertification, basetypes.ObjectAsOptions{})

		if !tfPlanCertification.CertificationDetailsUrl.Equal(tfStateCertification.CertificationDetailsUrl) {
			tfPlanCertificationDetailsUrl := tfPlanCertification.CertificationDetailsUrl.ValueString()
			requestBodyCertification.SetCertificationDetailsUrl(&tfPlanCertificationDetailsUrl)
		}

		if !tfPlanCertification.CertificationExpirationDateTime.Equal(tfStateCertification.CertificationExpirationDateTime) {
			tfPlanCertificationExpirationDateTime := tfPlanCertification.CertificationExpirationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanCertificationExpirationDateTime)
			requestBodyCertification.SetCertificationExpirationDateTime(&t)
		}

		if !tfPlanCertification.IsCertifiedByMicrosoft.Equal(tfStateCertification.IsCertifiedByMicrosoft) {
			tfPlanIsCertifiedByMicrosoft := tfPlanCertification.IsCertifiedByMicrosoft.ValueBool()
			requestBodyCertification.SetIsCertifiedByMicrosoft(&tfPlanIsCertifiedByMicrosoft)
		}

		if !tfPlanCertification.IsPublisherAttested.Equal(tfStateCertification.IsPublisherAttested) {
			tfPlanIsPublisherAttested := tfPlanCertification.IsPublisherAttested.ValueBool()
			requestBodyCertification.SetIsPublisherAttested(&tfPlanIsPublisherAttested)
		}

		if !tfPlanCertification.LastCertificationDateTime.Equal(tfStateCertification.LastCertificationDateTime) {
			tfPlanLastCertificationDateTime := tfPlanCertification.LastCertificationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastCertificationDateTime)
			requestBodyCertification.SetLastCertificationDateTime(&t)
		}
		requestBodyApplication.SetCertification(requestBodyCertification)
		tfPlanApplication.Certification, _ = types.ObjectValueFrom(ctx, tfPlanCertification.AttributeTypes(), tfPlanCertification)
	}

	if !tfPlanApplication.CreatedDateTime.Equal(tfStateApplication.CreatedDateTime) {
		tfPlanCreatedDateTime := tfPlanApplication.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyApplication.SetCreatedDateTime(&t)
	}

	if !tfPlanApplication.DefaultRedirectUri.Equal(tfStateApplication.DefaultRedirectUri) {
		tfPlanDefaultRedirectUri := tfPlanApplication.DefaultRedirectUri.ValueString()
		requestBodyApplication.SetDefaultRedirectUri(&tfPlanDefaultRedirectUri)
	}

	if !tfPlanApplication.Description.Equal(tfStateApplication.Description) {
		tfPlanDescription := tfPlanApplication.Description.ValueString()
		requestBodyApplication.SetDescription(&tfPlanDescription)
	}

	if !tfPlanApplication.DisabledByMicrosoftStatus.Equal(tfStateApplication.DisabledByMicrosoftStatus) {
		tfPlanDisabledByMicrosoftStatus := tfPlanApplication.DisabledByMicrosoftStatus.ValueString()
		requestBodyApplication.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	}

	if !tfPlanApplication.DisplayName.Equal(tfStateApplication.DisplayName) {
		tfPlanDisplayName := tfPlanApplication.DisplayName.ValueString()
		requestBodyApplication.SetDisplayName(&tfPlanDisplayName)
	}

	if !tfPlanApplication.GroupMembershipClaims.Equal(tfStateApplication.GroupMembershipClaims) {
		tfPlanGroupMembershipClaims := tfPlanApplication.GroupMembershipClaims.ValueString()
		requestBodyApplication.SetGroupMembershipClaims(&tfPlanGroupMembershipClaims)
	}

	if !tfPlanApplication.IdentifierUris.Equal(tfStateApplication.IdentifierUris) {
		var stringArrayIdentifierUris []string
		for _, i := range tfPlanApplication.IdentifierUris.Elements() {
			stringArrayIdentifierUris = append(stringArrayIdentifierUris, i.String())
		}
		requestBodyApplication.SetIdentifierUris(stringArrayIdentifierUris)
	}

	if !tfPlanApplication.Info.Equal(tfStateApplication.Info) {
		requestBodyInformationalUrl := models.NewInformationalUrl()
		tfPlanInformationalUrl := applicationInformationalUrlModel{}
		tfPlanApplication.Info.As(ctx, &tfPlanInformationalUrl, basetypes.ObjectAsOptions{})
		tfStateInformationalUrl := applicationInformationalUrlModel{}
		tfStateApplication.Info.As(ctx, &tfStateInformationalUrl, basetypes.ObjectAsOptions{})

		if !tfPlanInformationalUrl.LogoUrl.Equal(tfStateInformationalUrl.LogoUrl) {
			tfPlanLogoUrl := tfPlanInformationalUrl.LogoUrl.ValueString()
			requestBodyInformationalUrl.SetLogoUrl(&tfPlanLogoUrl)
		}

		if !tfPlanInformationalUrl.MarketingUrl.Equal(tfStateInformationalUrl.MarketingUrl) {
			tfPlanMarketingUrl := tfPlanInformationalUrl.MarketingUrl.ValueString()
			requestBodyInformationalUrl.SetMarketingUrl(&tfPlanMarketingUrl)
		}

		if !tfPlanInformationalUrl.PrivacyStatementUrl.Equal(tfStateInformationalUrl.PrivacyStatementUrl) {
			tfPlanPrivacyStatementUrl := tfPlanInformationalUrl.PrivacyStatementUrl.ValueString()
			requestBodyInformationalUrl.SetPrivacyStatementUrl(&tfPlanPrivacyStatementUrl)
		}

		if !tfPlanInformationalUrl.SupportUrl.Equal(tfStateInformationalUrl.SupportUrl) {
			tfPlanSupportUrl := tfPlanInformationalUrl.SupportUrl.ValueString()
			requestBodyInformationalUrl.SetSupportUrl(&tfPlanSupportUrl)
		}

		if !tfPlanInformationalUrl.TermsOfServiceUrl.Equal(tfStateInformationalUrl.TermsOfServiceUrl) {
			tfPlanTermsOfServiceUrl := tfPlanInformationalUrl.TermsOfServiceUrl.ValueString()
			requestBodyInformationalUrl.SetTermsOfServiceUrl(&tfPlanTermsOfServiceUrl)
		}
		requestBodyApplication.SetInfo(requestBodyInformationalUrl)
		tfPlanApplication.Info, _ = types.ObjectValueFrom(ctx, tfPlanInformationalUrl.AttributeTypes(), tfPlanInformationalUrl)
	}

	if !tfPlanApplication.IsDeviceOnlyAuthSupported.Equal(tfStateApplication.IsDeviceOnlyAuthSupported) {
		tfPlanIsDeviceOnlyAuthSupported := tfPlanApplication.IsDeviceOnlyAuthSupported.ValueBool()
		requestBodyApplication.SetIsDeviceOnlyAuthSupported(&tfPlanIsDeviceOnlyAuthSupported)
	}

	if !tfPlanApplication.IsFallbackPublicClient.Equal(tfStateApplication.IsFallbackPublicClient) {
		tfPlanIsFallbackPublicClient := tfPlanApplication.IsFallbackPublicClient.ValueBool()
		requestBodyApplication.SetIsFallbackPublicClient(&tfPlanIsFallbackPublicClient)
	}

	if !tfPlanApplication.KeyCredentials.Equal(tfStateApplication.KeyCredentials) {
		var tfPlanKeyCredentials []models.KeyCredentialable
		for k, i := range tfPlanApplication.KeyCredentials.Elements() {
			requestBodyKeyCredential := models.NewKeyCredential()
			tfPlanKeyCredential := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyCredential)
			tfStateKeyCredential := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, tfStateApplication.KeyCredentials.Elements()[k].Type(ctx), &tfPlanKeyCredential)

			if !tfPlanKeyCredential.CustomKeyIdentifier.Equal(tfStateKeyCredential.CustomKeyIdentifier) {
				tfPlanCustomKeyIdentifier := tfPlanKeyCredential.CustomKeyIdentifier.ValueString()
				requestBodyKeyCredential.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			}

			if !tfPlanKeyCredential.DisplayName.Equal(tfStateKeyCredential.DisplayName) {
				tfPlanDisplayName := tfPlanKeyCredential.DisplayName.ValueString()
				requestBodyKeyCredential.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanKeyCredential.EndDateTime.Equal(tfStateKeyCredential.EndDateTime) {
				tfPlanEndDateTime := tfPlanKeyCredential.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyKeyCredential.SetEndDateTime(&t)
			}

			if !tfPlanKeyCredential.Key.Equal(tfStateKeyCredential.Key) {
				tfPlanKey := tfPlanKeyCredential.Key.ValueString()
				requestBodyKeyCredential.SetKey([]byte(tfPlanKey))
			}

			if !tfPlanKeyCredential.KeyId.Equal(tfStateKeyCredential.KeyId) {
				tfPlanKeyId := tfPlanKeyCredential.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyKeyCredential.SetKeyId(&u)
			}

			if !tfPlanKeyCredential.StartDateTime.Equal(tfStateKeyCredential.StartDateTime) {
				tfPlanStartDateTime := tfPlanKeyCredential.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyKeyCredential.SetStartDateTime(&t)
			}

			if !tfPlanKeyCredential.Type.Equal(tfStateKeyCredential.Type) {
				tfPlanType := tfPlanKeyCredential.Type.ValueString()
				requestBodyKeyCredential.SetTypeEscaped(&tfPlanType)
			}

			if !tfPlanKeyCredential.Usage.Equal(tfStateKeyCredential.Usage) {
				tfPlanUsage := tfPlanKeyCredential.Usage.ValueString()
				requestBodyKeyCredential.SetUsage(&tfPlanUsage)
			}
		}
		requestBodyApplication.SetKeyCredentials(tfPlanKeyCredentials)
	}

	if !tfPlanApplication.Logo.Equal(tfStateApplication.Logo) {
		tfPlanLogo := tfPlanApplication.Logo.ValueString()
		requestBodyApplication.SetLogo([]byte(tfPlanLogo))
	}

	if !tfPlanApplication.NativeAuthenticationApisEnabled.Equal(tfStateApplication.NativeAuthenticationApisEnabled) {
		tfPlanNativeAuthenticationApisEnabled := tfPlanApplication.NativeAuthenticationApisEnabled.ValueString()
		parsedNativeAuthenticationApisEnabled, _ := models.ParseNativeAuthenticationApisEnabled(tfPlanNativeAuthenticationApisEnabled)
		assertedNativeAuthenticationApisEnabled := parsedNativeAuthenticationApisEnabled.(models.NativeAuthenticationApisEnabled)
		requestBodyApplication.SetNativeAuthenticationApisEnabled(&assertedNativeAuthenticationApisEnabled)
	}

	if !tfPlanApplication.Notes.Equal(tfStateApplication.Notes) {
		tfPlanNotes := tfPlanApplication.Notes.ValueString()
		requestBodyApplication.SetNotes(&tfPlanNotes)
	}

	if !tfPlanApplication.Oauth2RequirePostResponse.Equal(tfStateApplication.Oauth2RequirePostResponse) {
		tfPlanOauth2RequirePostResponse := tfPlanApplication.Oauth2RequirePostResponse.ValueBool()
		requestBodyApplication.SetOauth2RequirePostResponse(&tfPlanOauth2RequirePostResponse)
	}

	if !tfPlanApplication.OptionalClaims.Equal(tfStateApplication.OptionalClaims) {
		requestBodyOptionalClaims := models.NewOptionalClaims()
		tfPlanOptionalClaims := applicationOptionalClaimsModel{}
		tfPlanApplication.OptionalClaims.As(ctx, &tfPlanOptionalClaims, basetypes.ObjectAsOptions{})
		tfStateOptionalClaims := applicationOptionalClaimsModel{}
		tfStateApplication.OptionalClaims.As(ctx, &tfStateOptionalClaims, basetypes.ObjectAsOptions{})

		if !tfPlanOptionalClaims.AccessToken.Equal(tfStateOptionalClaims.AccessToken) {
			var tfPlanAccessToken []models.OptionalClaimable
			for k, i := range tfPlanOptionalClaims.AccessToken.Elements() {
				requestBodyOptionalClaim := models.NewOptionalClaim()
				tfPlanOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOptionalClaim)
				tfStateOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, tfStateOptionalClaims.AccessToken.Elements()[k].Type(ctx), &tfPlanOptionalClaim)

				if !tfPlanOptionalClaim.AdditionalProperties.Equal(tfStateOptionalClaim.AdditionalProperties) {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanOptionalClaim.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyOptionalClaim.SetAdditionalProperties(stringArrayAdditionalProperties)
				}

				if !tfPlanOptionalClaim.Essential.Equal(tfStateOptionalClaim.Essential) {
					tfPlanEssential := tfPlanOptionalClaim.Essential.ValueBool()
					requestBodyOptionalClaim.SetEssential(&tfPlanEssential)
				}

				if !tfPlanOptionalClaim.Name.Equal(tfStateOptionalClaim.Name) {
					tfPlanName := tfPlanOptionalClaim.Name.ValueString()
					requestBodyOptionalClaim.SetName(&tfPlanName)
				}

				if !tfPlanOptionalClaim.Source.Equal(tfStateOptionalClaim.Source) {
					tfPlanSource := tfPlanOptionalClaim.Source.ValueString()
					requestBodyOptionalClaim.SetSource(&tfPlanSource)
				}
			}
			requestBodyOptionalClaims.SetAccessToken(tfPlanAccessToken)
		}

		if !tfPlanOptionalClaims.IdToken.Equal(tfStateOptionalClaims.IdToken) {
			var tfPlanIdToken []models.OptionalClaimable
			for k, i := range tfPlanOptionalClaims.IdToken.Elements() {
				requestBodyOptionalClaim := models.NewOptionalClaim()
				tfPlanOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOptionalClaim)
				tfStateOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, tfStateOptionalClaims.IdToken.Elements()[k].Type(ctx), &tfPlanOptionalClaim)

				if !tfPlanOptionalClaim.AdditionalProperties.Equal(tfStateOptionalClaim.AdditionalProperties) {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanOptionalClaim.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyOptionalClaim.SetAdditionalProperties(stringArrayAdditionalProperties)
				}

				if !tfPlanOptionalClaim.Essential.Equal(tfStateOptionalClaim.Essential) {
					tfPlanEssential := tfPlanOptionalClaim.Essential.ValueBool()
					requestBodyOptionalClaim.SetEssential(&tfPlanEssential)
				}

				if !tfPlanOptionalClaim.Name.Equal(tfStateOptionalClaim.Name) {
					tfPlanName := tfPlanOptionalClaim.Name.ValueString()
					requestBodyOptionalClaim.SetName(&tfPlanName)
				}

				if !tfPlanOptionalClaim.Source.Equal(tfStateOptionalClaim.Source) {
					tfPlanSource := tfPlanOptionalClaim.Source.ValueString()
					requestBodyOptionalClaim.SetSource(&tfPlanSource)
				}
			}
			requestBodyOptionalClaims.SetIdToken(tfPlanIdToken)
		}

		if !tfPlanOptionalClaims.Saml2Token.Equal(tfStateOptionalClaims.Saml2Token) {
			var tfPlanSaml2Token []models.OptionalClaimable
			for k, i := range tfPlanOptionalClaims.Saml2Token.Elements() {
				requestBodyOptionalClaim := models.NewOptionalClaim()
				tfPlanOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOptionalClaim)
				tfStateOptionalClaim := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, tfStateOptionalClaims.Saml2Token.Elements()[k].Type(ctx), &tfPlanOptionalClaim)

				if !tfPlanOptionalClaim.AdditionalProperties.Equal(tfStateOptionalClaim.AdditionalProperties) {
					var stringArrayAdditionalProperties []string
					for _, i := range tfPlanOptionalClaim.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					requestBodyOptionalClaim.SetAdditionalProperties(stringArrayAdditionalProperties)
				}

				if !tfPlanOptionalClaim.Essential.Equal(tfStateOptionalClaim.Essential) {
					tfPlanEssential := tfPlanOptionalClaim.Essential.ValueBool()
					requestBodyOptionalClaim.SetEssential(&tfPlanEssential)
				}

				if !tfPlanOptionalClaim.Name.Equal(tfStateOptionalClaim.Name) {
					tfPlanName := tfPlanOptionalClaim.Name.ValueString()
					requestBodyOptionalClaim.SetName(&tfPlanName)
				}

				if !tfPlanOptionalClaim.Source.Equal(tfStateOptionalClaim.Source) {
					tfPlanSource := tfPlanOptionalClaim.Source.ValueString()
					requestBodyOptionalClaim.SetSource(&tfPlanSource)
				}
			}
			requestBodyOptionalClaims.SetSaml2Token(tfPlanSaml2Token)
		}
		requestBodyApplication.SetOptionalClaims(requestBodyOptionalClaims)
		tfPlanApplication.OptionalClaims, _ = types.ObjectValueFrom(ctx, tfPlanOptionalClaims.AttributeTypes(), tfPlanOptionalClaims)
	}

	if !tfPlanApplication.ParentalControlSettings.Equal(tfStateApplication.ParentalControlSettings) {
		requestBodyParentalControlSettings := models.NewParentalControlSettings()
		tfPlanParentalControlSettings := applicationParentalControlSettingsModel{}
		tfPlanApplication.ParentalControlSettings.As(ctx, &tfPlanParentalControlSettings, basetypes.ObjectAsOptions{})
		tfStateParentalControlSettings := applicationParentalControlSettingsModel{}
		tfStateApplication.ParentalControlSettings.As(ctx, &tfStateParentalControlSettings, basetypes.ObjectAsOptions{})

		if !tfPlanParentalControlSettings.CountriesBlockedForMinors.Equal(tfStateParentalControlSettings.CountriesBlockedForMinors) {
			var stringArrayCountriesBlockedForMinors []string
			for _, i := range tfPlanParentalControlSettings.CountriesBlockedForMinors.Elements() {
				stringArrayCountriesBlockedForMinors = append(stringArrayCountriesBlockedForMinors, i.String())
			}
			requestBodyParentalControlSettings.SetCountriesBlockedForMinors(stringArrayCountriesBlockedForMinors)
		}

		if !tfPlanParentalControlSettings.LegalAgeGroupRule.Equal(tfStateParentalControlSettings.LegalAgeGroupRule) {
			tfPlanLegalAgeGroupRule := tfPlanParentalControlSettings.LegalAgeGroupRule.ValueString()
			requestBodyParentalControlSettings.SetLegalAgeGroupRule(&tfPlanLegalAgeGroupRule)
		}
		requestBodyApplication.SetParentalControlSettings(requestBodyParentalControlSettings)
		tfPlanApplication.ParentalControlSettings, _ = types.ObjectValueFrom(ctx, tfPlanParentalControlSettings.AttributeTypes(), tfPlanParentalControlSettings)
	}

	if !tfPlanApplication.PasswordCredentials.Equal(tfStateApplication.PasswordCredentials) {
		var tfPlanPasswordCredentials []models.PasswordCredentialable
		for k, i := range tfPlanApplication.PasswordCredentials.Elements() {
			requestBodyPasswordCredential := models.NewPasswordCredential()
			tfPlanPasswordCredential := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPasswordCredential)
			tfStatePasswordCredential := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, tfStateApplication.PasswordCredentials.Elements()[k].Type(ctx), &tfPlanPasswordCredential)

			if !tfPlanPasswordCredential.CustomKeyIdentifier.Equal(tfStatePasswordCredential.CustomKeyIdentifier) {
				tfPlanCustomKeyIdentifier := tfPlanPasswordCredential.CustomKeyIdentifier.ValueString()
				requestBodyPasswordCredential.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			}

			if !tfPlanPasswordCredential.DisplayName.Equal(tfStatePasswordCredential.DisplayName) {
				tfPlanDisplayName := tfPlanPasswordCredential.DisplayName.ValueString()
				requestBodyPasswordCredential.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanPasswordCredential.EndDateTime.Equal(tfStatePasswordCredential.EndDateTime) {
				tfPlanEndDateTime := tfPlanPasswordCredential.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyPasswordCredential.SetEndDateTime(&t)
			}

			if !tfPlanPasswordCredential.Hint.Equal(tfStatePasswordCredential.Hint) {
				tfPlanHint := tfPlanPasswordCredential.Hint.ValueString()
				requestBodyPasswordCredential.SetHint(&tfPlanHint)
			}

			if !tfPlanPasswordCredential.KeyId.Equal(tfStatePasswordCredential.KeyId) {
				tfPlanKeyId := tfPlanPasswordCredential.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyPasswordCredential.SetKeyId(&u)
			}

			if !tfPlanPasswordCredential.SecretText.Equal(tfStatePasswordCredential.SecretText) {
				tfPlanSecretText := tfPlanPasswordCredential.SecretText.ValueString()
				requestBodyPasswordCredential.SetSecretText(&tfPlanSecretText)
			}

			if !tfPlanPasswordCredential.StartDateTime.Equal(tfStatePasswordCredential.StartDateTime) {
				tfPlanStartDateTime := tfPlanPasswordCredential.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyPasswordCredential.SetStartDateTime(&t)
			}
		}
		requestBodyApplication.SetPasswordCredentials(tfPlanPasswordCredentials)
	}

	if !tfPlanApplication.PublicClient.Equal(tfStateApplication.PublicClient) {
		requestBodyPublicClientApplication := models.NewPublicClientApplication()
		tfPlanPublicClientApplication := applicationPublicClientApplicationModel{}
		tfPlanApplication.PublicClient.As(ctx, &tfPlanPublicClientApplication, basetypes.ObjectAsOptions{})
		tfStatePublicClientApplication := applicationPublicClientApplicationModel{}
		tfStateApplication.PublicClient.As(ctx, &tfStatePublicClientApplication, basetypes.ObjectAsOptions{})

		if !tfPlanPublicClientApplication.RedirectUris.Equal(tfStatePublicClientApplication.RedirectUris) {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanPublicClientApplication.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodyPublicClientApplication.SetRedirectUris(stringArrayRedirectUris)
		}
		requestBodyApplication.SetPublicClient(requestBodyPublicClientApplication)
		tfPlanApplication.PublicClient, _ = types.ObjectValueFrom(ctx, tfPlanPublicClientApplication.AttributeTypes(), tfPlanPublicClientApplication)
	}

	if !tfPlanApplication.PublisherDomain.Equal(tfStateApplication.PublisherDomain) {
		tfPlanPublisherDomain := tfPlanApplication.PublisherDomain.ValueString()
		requestBodyApplication.SetPublisherDomain(&tfPlanPublisherDomain)
	}

	if !tfPlanApplication.RequestSignatureVerification.Equal(tfStateApplication.RequestSignatureVerification) {
		requestBodyRequestSignatureVerification := models.NewRequestSignatureVerification()
		tfPlanRequestSignatureVerification := applicationRequestSignatureVerificationModel{}
		tfPlanApplication.RequestSignatureVerification.As(ctx, &tfPlanRequestSignatureVerification, basetypes.ObjectAsOptions{})
		tfStateRequestSignatureVerification := applicationRequestSignatureVerificationModel{}
		tfStateApplication.RequestSignatureVerification.As(ctx, &tfStateRequestSignatureVerification, basetypes.ObjectAsOptions{})

		if !tfPlanRequestSignatureVerification.AllowedWeakAlgorithms.Equal(tfStateRequestSignatureVerification.AllowedWeakAlgorithms) {
			tfPlanAllowedWeakAlgorithms := tfPlanRequestSignatureVerification.AllowedWeakAlgorithms.ValueString()
			parsedAllowedWeakAlgorithms, _ := models.ParseWeakAlgorithms(tfPlanAllowedWeakAlgorithms)
			assertedAllowedWeakAlgorithms := parsedAllowedWeakAlgorithms.(models.WeakAlgorithms)
			requestBodyRequestSignatureVerification.SetAllowedWeakAlgorithms(&assertedAllowedWeakAlgorithms)
		}

		if !tfPlanRequestSignatureVerification.IsSignedRequestRequired.Equal(tfStateRequestSignatureVerification.IsSignedRequestRequired) {
			tfPlanIsSignedRequestRequired := tfPlanRequestSignatureVerification.IsSignedRequestRequired.ValueBool()
			requestBodyRequestSignatureVerification.SetIsSignedRequestRequired(&tfPlanIsSignedRequestRequired)
		}
		requestBodyApplication.SetRequestSignatureVerification(requestBodyRequestSignatureVerification)
		tfPlanApplication.RequestSignatureVerification, _ = types.ObjectValueFrom(ctx, tfPlanRequestSignatureVerification.AttributeTypes(), tfPlanRequestSignatureVerification)
	}

	if !tfPlanApplication.RequiredResourceAccess.Equal(tfStateApplication.RequiredResourceAccess) {
		var tfPlanRequiredResourceAccess []models.RequiredResourceAccessable
		for k, i := range tfPlanApplication.RequiredResourceAccess.Elements() {
			requestBodyRequiredResourceAccess := models.NewRequiredResourceAccess()
			tfPlanRequiredResourceAccess := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanRequiredResourceAccess)
			tfStateRequiredResourceAccess := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, tfStateApplication.RequiredResourceAccess.Elements()[k].Type(ctx), &tfPlanRequiredResourceAccess)

			if !tfPlanRequiredResourceAccess.ResourceAccess.Equal(tfStateRequiredResourceAccess.ResourceAccess) {
				var tfPlanResourceAccess []models.ResourceAccessable
				for k, i := range tfPlanRequiredResourceAccess.ResourceAccess.Elements() {
					requestBodyResourceAccess := models.NewResourceAccess()
					tfPlanResourceAccess := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanResourceAccess)
					tfStateResourceAccess := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, tfStateRequiredResourceAccess.ResourceAccess.Elements()[k].Type(ctx), &tfPlanResourceAccess)

					if !tfPlanResourceAccess.Id.Equal(tfStateResourceAccess.Id) {
						tfPlanId := tfPlanResourceAccess.Id.ValueString()
						u, _ := uuid.Parse(tfPlanId)
						requestBodyResourceAccess.SetId(&u)
					}

					if !tfPlanResourceAccess.Type.Equal(tfStateResourceAccess.Type) {
						tfPlanType := tfPlanResourceAccess.Type.ValueString()
						requestBodyResourceAccess.SetTypeEscaped(&tfPlanType)
					}
				}
				requestBodyRequiredResourceAccess.SetResourceAccess(tfPlanResourceAccess)
			}

			if !tfPlanRequiredResourceAccess.ResourceAppId.Equal(tfStateRequiredResourceAccess.ResourceAppId) {
				tfPlanResourceAppId := tfPlanRequiredResourceAccess.ResourceAppId.ValueString()
				requestBodyRequiredResourceAccess.SetResourceAppId(&tfPlanResourceAppId)
			}
		}
		requestBodyApplication.SetRequiredResourceAccess(tfPlanRequiredResourceAccess)
	}

	if !tfPlanApplication.SamlMetadataUrl.Equal(tfStateApplication.SamlMetadataUrl) {
		tfPlanSamlMetadataUrl := tfPlanApplication.SamlMetadataUrl.ValueString()
		requestBodyApplication.SetSamlMetadataUrl(&tfPlanSamlMetadataUrl)
	}

	if !tfPlanApplication.ServiceManagementReference.Equal(tfStateApplication.ServiceManagementReference) {
		tfPlanServiceManagementReference := tfPlanApplication.ServiceManagementReference.ValueString()
		requestBodyApplication.SetServiceManagementReference(&tfPlanServiceManagementReference)
	}

	if !tfPlanApplication.ServicePrincipalLockConfiguration.Equal(tfStateApplication.ServicePrincipalLockConfiguration) {
		requestBodyServicePrincipalLockConfiguration := models.NewServicePrincipalLockConfiguration()
		tfPlanServicePrincipalLockConfiguration := applicationServicePrincipalLockConfigurationModel{}
		tfPlanApplication.ServicePrincipalLockConfiguration.As(ctx, &tfPlanServicePrincipalLockConfiguration, basetypes.ObjectAsOptions{})
		tfStateServicePrincipalLockConfiguration := applicationServicePrincipalLockConfigurationModel{}
		tfStateApplication.ServicePrincipalLockConfiguration.As(ctx, &tfStateServicePrincipalLockConfiguration, basetypes.ObjectAsOptions{})

		if !tfPlanServicePrincipalLockConfiguration.AllProperties.Equal(tfStateServicePrincipalLockConfiguration.AllProperties) {
			tfPlanAllProperties := tfPlanServicePrincipalLockConfiguration.AllProperties.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetAllProperties(&tfPlanAllProperties)
		}

		if !tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageSign.Equal(tfStateServicePrincipalLockConfiguration.CredentialsWithUsageSign) {
			tfPlanCredentialsWithUsageSign := tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageSign.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetCredentialsWithUsageSign(&tfPlanCredentialsWithUsageSign)
		}

		if !tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageVerify.Equal(tfStateServicePrincipalLockConfiguration.CredentialsWithUsageVerify) {
			tfPlanCredentialsWithUsageVerify := tfPlanServicePrincipalLockConfiguration.CredentialsWithUsageVerify.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetCredentialsWithUsageVerify(&tfPlanCredentialsWithUsageVerify)
		}

		if !tfPlanServicePrincipalLockConfiguration.IsEnabled.Equal(tfStateServicePrincipalLockConfiguration.IsEnabled) {
			tfPlanIsEnabled := tfPlanServicePrincipalLockConfiguration.IsEnabled.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetIsEnabled(&tfPlanIsEnabled)
		}

		if !tfPlanServicePrincipalLockConfiguration.TokenEncryptionKeyId.Equal(tfStateServicePrincipalLockConfiguration.TokenEncryptionKeyId) {
			tfPlanTokenEncryptionKeyId := tfPlanServicePrincipalLockConfiguration.TokenEncryptionKeyId.ValueBool()
			requestBodyServicePrincipalLockConfiguration.SetTokenEncryptionKeyId(&tfPlanTokenEncryptionKeyId)
		}
		requestBodyApplication.SetServicePrincipalLockConfiguration(requestBodyServicePrincipalLockConfiguration)
		tfPlanApplication.ServicePrincipalLockConfiguration, _ = types.ObjectValueFrom(ctx, tfPlanServicePrincipalLockConfiguration.AttributeTypes(), tfPlanServicePrincipalLockConfiguration)
	}

	if !tfPlanApplication.SignInAudience.Equal(tfStateApplication.SignInAudience) {
		tfPlanSignInAudience := tfPlanApplication.SignInAudience.ValueString()
		requestBodyApplication.SetSignInAudience(&tfPlanSignInAudience)
	}

	if !tfPlanApplication.Spa.Equal(tfStateApplication.Spa) {
		requestBodySpaApplication := models.NewSpaApplication()
		tfPlanSpaApplication := applicationSpaApplicationModel{}
		tfPlanApplication.Spa.As(ctx, &tfPlanSpaApplication, basetypes.ObjectAsOptions{})
		tfStateSpaApplication := applicationSpaApplicationModel{}
		tfStateApplication.Spa.As(ctx, &tfStateSpaApplication, basetypes.ObjectAsOptions{})

		if !tfPlanSpaApplication.RedirectUris.Equal(tfStateSpaApplication.RedirectUris) {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanSpaApplication.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodySpaApplication.SetRedirectUris(stringArrayRedirectUris)
		}
		requestBodyApplication.SetSpa(requestBodySpaApplication)
		tfPlanApplication.Spa, _ = types.ObjectValueFrom(ctx, tfPlanSpaApplication.AttributeTypes(), tfPlanSpaApplication)
	}

	if !tfPlanApplication.Tags.Equal(tfStateApplication.Tags) {
		var stringArrayTags []string
		for _, i := range tfPlanApplication.Tags.Elements() {
			stringArrayTags = append(stringArrayTags, i.String())
		}
		requestBodyApplication.SetTags(stringArrayTags)
	}

	if !tfPlanApplication.TokenEncryptionKeyId.Equal(tfStateApplication.TokenEncryptionKeyId) {
		tfPlanTokenEncryptionKeyId := tfPlanApplication.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		requestBodyApplication.SetTokenEncryptionKeyId(&u)
	}

	if !tfPlanApplication.UniqueName.Equal(tfStateApplication.UniqueName) {
		tfPlanUniqueName := tfPlanApplication.UniqueName.ValueString()
		requestBodyApplication.SetUniqueName(&tfPlanUniqueName)
	}

	if !tfPlanApplication.VerifiedPublisher.Equal(tfStateApplication.VerifiedPublisher) {
		requestBodyVerifiedPublisher := models.NewVerifiedPublisher()
		tfPlanVerifiedPublisher := applicationVerifiedPublisherModel{}
		tfPlanApplication.VerifiedPublisher.As(ctx, &tfPlanVerifiedPublisher, basetypes.ObjectAsOptions{})
		tfStateVerifiedPublisher := applicationVerifiedPublisherModel{}
		tfStateApplication.VerifiedPublisher.As(ctx, &tfStateVerifiedPublisher, basetypes.ObjectAsOptions{})

		if !tfPlanVerifiedPublisher.AddedDateTime.Equal(tfStateVerifiedPublisher.AddedDateTime) {
			tfPlanAddedDateTime := tfPlanVerifiedPublisher.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanAddedDateTime)
			requestBodyVerifiedPublisher.SetAddedDateTime(&t)
		}

		if !tfPlanVerifiedPublisher.DisplayName.Equal(tfStateVerifiedPublisher.DisplayName) {
			tfPlanDisplayName := tfPlanVerifiedPublisher.DisplayName.ValueString()
			requestBodyVerifiedPublisher.SetDisplayName(&tfPlanDisplayName)
		}

		if !tfPlanVerifiedPublisher.VerifiedPublisherId.Equal(tfStateVerifiedPublisher.VerifiedPublisherId) {
			tfPlanVerifiedPublisherId := tfPlanVerifiedPublisher.VerifiedPublisherId.ValueString()
			requestBodyVerifiedPublisher.SetVerifiedPublisherId(&tfPlanVerifiedPublisherId)
		}
		requestBodyApplication.SetVerifiedPublisher(requestBodyVerifiedPublisher)
		tfPlanApplication.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfPlanVerifiedPublisher.AttributeTypes(), tfPlanVerifiedPublisher)
	}

	if !tfPlanApplication.Web.Equal(tfStateApplication.Web) {
		requestBodyWebApplication := models.NewWebApplication()
		tfPlanWebApplication := applicationWebApplicationModel{}
		tfPlanApplication.Web.As(ctx, &tfPlanWebApplication, basetypes.ObjectAsOptions{})
		tfStateWebApplication := applicationWebApplicationModel{}
		tfStateApplication.Web.As(ctx, &tfStateWebApplication, basetypes.ObjectAsOptions{})

		if !tfPlanWebApplication.HomePageUrl.Equal(tfStateWebApplication.HomePageUrl) {
			tfPlanHomePageUrl := tfPlanWebApplication.HomePageUrl.ValueString()
			requestBodyWebApplication.SetHomePageUrl(&tfPlanHomePageUrl)
		}

		if !tfPlanWebApplication.ImplicitGrantSettings.Equal(tfStateWebApplication.ImplicitGrantSettings) {
			requestBodyImplicitGrantSettings := models.NewImplicitGrantSettings()
			tfPlanImplicitGrantSettings := applicationImplicitGrantSettingsModel{}
			tfPlanWebApplication.ImplicitGrantSettings.As(ctx, &tfPlanImplicitGrantSettings, basetypes.ObjectAsOptions{})
			tfStateImplicitGrantSettings := applicationImplicitGrantSettingsModel{}
			tfStateWebApplication.ImplicitGrantSettings.As(ctx, &tfStateImplicitGrantSettings, basetypes.ObjectAsOptions{})

			if !tfPlanImplicitGrantSettings.EnableAccessTokenIssuance.Equal(tfStateImplicitGrantSettings.EnableAccessTokenIssuance) {
				tfPlanEnableAccessTokenIssuance := tfPlanImplicitGrantSettings.EnableAccessTokenIssuance.ValueBool()
				requestBodyImplicitGrantSettings.SetEnableAccessTokenIssuance(&tfPlanEnableAccessTokenIssuance)
			}

			if !tfPlanImplicitGrantSettings.EnableIdTokenIssuance.Equal(tfStateImplicitGrantSettings.EnableIdTokenIssuance) {
				tfPlanEnableIdTokenIssuance := tfPlanImplicitGrantSettings.EnableIdTokenIssuance.ValueBool()
				requestBodyImplicitGrantSettings.SetEnableIdTokenIssuance(&tfPlanEnableIdTokenIssuance)
			}
			requestBodyWebApplication.SetImplicitGrantSettings(requestBodyImplicitGrantSettings)
			tfPlanWebApplication.ImplicitGrantSettings, _ = types.ObjectValueFrom(ctx, tfPlanImplicitGrantSettings.AttributeTypes(), tfPlanImplicitGrantSettings)
		}

		if !tfPlanWebApplication.LogoutUrl.Equal(tfStateWebApplication.LogoutUrl) {
			tfPlanLogoutUrl := tfPlanWebApplication.LogoutUrl.ValueString()
			requestBodyWebApplication.SetLogoutUrl(&tfPlanLogoutUrl)
		}

		if !tfPlanWebApplication.RedirectUriSettings.Equal(tfStateWebApplication.RedirectUriSettings) {
			var tfPlanRedirectUriSettings []models.RedirectUriSettingsable
			for k, i := range tfPlanWebApplication.RedirectUriSettings.Elements() {
				requestBodyRedirectUriSettings := models.NewRedirectUriSettings()
				tfPlanRedirectUriSettings := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfPlanRedirectUriSettings)
				tfStateRedirectUriSettings := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, tfStateWebApplication.RedirectUriSettings.Elements()[k].Type(ctx), &tfPlanRedirectUriSettings)

				if !tfPlanRedirectUriSettings.Uri.Equal(tfStateRedirectUriSettings.Uri) {
					tfPlanUri := tfPlanRedirectUriSettings.Uri.ValueString()
					requestBodyRedirectUriSettings.SetUri(&tfPlanUri)
				}
			}
			requestBodyWebApplication.SetRedirectUriSettings(tfPlanRedirectUriSettings)
		}

		if !tfPlanWebApplication.RedirectUris.Equal(tfStateWebApplication.RedirectUris) {
			var stringArrayRedirectUris []string
			for _, i := range tfPlanWebApplication.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			requestBodyWebApplication.SetRedirectUris(stringArrayRedirectUris)
		}
		requestBodyApplication.SetWeb(requestBodyWebApplication)
		tfPlanApplication.Web, _ = types.ObjectValueFrom(ctx, tfPlanWebApplication.AttributeTypes(), tfPlanWebApplication)
	}

	// Update application
	_, err := r.client.Applications().ByApplicationId(tfStateApplication.Id.ValueString()).Patch(context.Background(), requestBodyApplication, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating application",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlanApplication)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfStateApplication applicationModel
	diags := req.State.Get(ctx, &tfStateApplication)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete application
	err := r.client.Applications().ByApplicationId(tfStateApplication.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting application",
			err.Error(),
		)
		return
	}

}
