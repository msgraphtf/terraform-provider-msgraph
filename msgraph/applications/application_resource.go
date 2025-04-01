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
	var tfPlan applicationModel
	diags := req.Plan.Get(ctx, &tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	sdkModelApplication := models.NewApplication()

	if !tfPlan.Id.IsUnknown() {
		tfPlanId := tfPlan.Id.ValueString()
		sdkModelApplication.SetId(&tfPlanId)
	} else {
		tfPlan.Id = types.StringNull()
	}

	if !tfPlan.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		sdkModelApplication.SetDeletedDateTime(&t)
	} else {
		tfPlan.DeletedDateTime = types.StringNull()
	}

	if len(tfPlan.AddIns.Elements()) > 0 {
		var tfPlanAddIns []models.AddInable
		for _, i := range tfPlan.AddIns.Elements() {
			sdkModelAddIns := models.NewAddIn()
			tfModelAddIns := applicationAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelAddIns)

			if !tfModelAddIns.Id.IsUnknown() {
				tfPlanId := tfModelAddIns.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				sdkModelAddIns.SetId(&u)
			} else {
				tfModelAddIns.Id = types.StringNull()
			}

			if len(tfModelAddIns.Properties.Elements()) > 0 {
				var tfPlanProperties []models.KeyValueable
				for _, i := range tfModelAddIns.Properties.Elements() {
					sdkModelProperties := models.NewKeyValue()
					tfModelProperties := applicationKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfModelProperties)

					if !tfModelProperties.Key.IsUnknown() {
						tfPlanKey := tfModelProperties.Key.ValueString()
						sdkModelProperties.SetKey(&tfPlanKey)
					} else {
						tfModelProperties.Key = types.StringNull()
					}

					if !tfModelProperties.Value.IsUnknown() {
						tfPlanValue := tfModelProperties.Value.ValueString()
						sdkModelProperties.SetValue(&tfPlanValue)
					} else {
						tfModelProperties.Value = types.StringNull()
					}
				}
				sdkModelAddIns.SetProperties(tfPlanProperties)
			} else {
				tfModelAddIns.Properties = types.ListNull(tfModelAddIns.Properties.ElementType(ctx))
			}

			if !tfModelAddIns.Type.IsUnknown() {
				tfPlanType := tfModelAddIns.Type.ValueString()
				sdkModelAddIns.SetTypeEscaped(&tfPlanType)
			} else {
				tfModelAddIns.Type = types.StringNull()
			}
		}
		sdkModelApplication.SetAddIns(tfPlanAddIns)
	} else {
		tfPlan.AddIns = types.ListNull(tfPlan.AddIns.ElementType(ctx))
	}

	if !tfPlan.Api.IsUnknown() {
		sdkModelApi := models.NewApiApplication()
		tfModelApi := applicationApiApplicationModel{}
		tfPlan.Api.As(ctx, &tfModelApi, basetypes.ObjectAsOptions{})

		if !tfModelApi.AcceptMappedClaims.IsUnknown() {
			tfPlanAcceptMappedClaims := tfModelApi.AcceptMappedClaims.ValueBool()
			sdkModelApi.SetAcceptMappedClaims(&tfPlanAcceptMappedClaims)
		} else {
			tfModelApi.AcceptMappedClaims = types.BoolNull()
		}

		if len(tfModelApi.KnownClientApplications.Elements()) > 0 {
			var uuidArrayKnownClientApplications []uuid.UUID
			for _, i := range tfModelApi.KnownClientApplications.Elements() {
				u, _ := uuid.Parse(i.String())
				uuidArrayKnownClientApplications = append(uuidArrayKnownClientApplications, u)
			}
			sdkModelApi.SetKnownClientApplications(uuidArrayKnownClientApplications)
		} else {
			tfModelApi.KnownClientApplications = types.ListNull(types.StringType)
		}

		if len(tfModelApi.Oauth2PermissionScopes.Elements()) > 0 {
			var tfPlanOauth2PermissionScopes []models.PermissionScopeable
			for _, i := range tfModelApi.Oauth2PermissionScopes.Elements() {
				sdkModelOauth2PermissionScopes := models.NewPermissionScope()
				tfModelOauth2PermissionScopes := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfModelOauth2PermissionScopes)

				if !tfModelOauth2PermissionScopes.AdminConsentDescription.IsUnknown() {
					tfPlanAdminConsentDescription := tfModelOauth2PermissionScopes.AdminConsentDescription.ValueString()
					sdkModelOauth2PermissionScopes.SetAdminConsentDescription(&tfPlanAdminConsentDescription)
				} else {
					tfModelOauth2PermissionScopes.AdminConsentDescription = types.StringNull()
				}

				if !tfModelOauth2PermissionScopes.AdminConsentDisplayName.IsUnknown() {
					tfPlanAdminConsentDisplayName := tfModelOauth2PermissionScopes.AdminConsentDisplayName.ValueString()
					sdkModelOauth2PermissionScopes.SetAdminConsentDisplayName(&tfPlanAdminConsentDisplayName)
				} else {
					tfModelOauth2PermissionScopes.AdminConsentDisplayName = types.StringNull()
				}

				if !tfModelOauth2PermissionScopes.Id.IsUnknown() {
					tfPlanId := tfModelOauth2PermissionScopes.Id.ValueString()
					u, _ := uuid.Parse(tfPlanId)
					sdkModelOauth2PermissionScopes.SetId(&u)
				} else {
					tfModelOauth2PermissionScopes.Id = types.StringNull()
				}

				if !tfModelOauth2PermissionScopes.IsEnabled.IsUnknown() {
					tfPlanIsEnabled := tfModelOauth2PermissionScopes.IsEnabled.ValueBool()
					sdkModelOauth2PermissionScopes.SetIsEnabled(&tfPlanIsEnabled)
				} else {
					tfModelOauth2PermissionScopes.IsEnabled = types.BoolNull()
				}

				if !tfModelOauth2PermissionScopes.Origin.IsUnknown() {
					tfPlanOrigin := tfModelOauth2PermissionScopes.Origin.ValueString()
					sdkModelOauth2PermissionScopes.SetOrigin(&tfPlanOrigin)
				} else {
					tfModelOauth2PermissionScopes.Origin = types.StringNull()
				}

				if !tfModelOauth2PermissionScopes.Type.IsUnknown() {
					tfPlanType := tfModelOauth2PermissionScopes.Type.ValueString()
					sdkModelOauth2PermissionScopes.SetTypeEscaped(&tfPlanType)
				} else {
					tfModelOauth2PermissionScopes.Type = types.StringNull()
				}

				if !tfModelOauth2PermissionScopes.UserConsentDescription.IsUnknown() {
					tfPlanUserConsentDescription := tfModelOauth2PermissionScopes.UserConsentDescription.ValueString()
					sdkModelOauth2PermissionScopes.SetUserConsentDescription(&tfPlanUserConsentDescription)
				} else {
					tfModelOauth2PermissionScopes.UserConsentDescription = types.StringNull()
				}

				if !tfModelOauth2PermissionScopes.UserConsentDisplayName.IsUnknown() {
					tfPlanUserConsentDisplayName := tfModelOauth2PermissionScopes.UserConsentDisplayName.ValueString()
					sdkModelOauth2PermissionScopes.SetUserConsentDisplayName(&tfPlanUserConsentDisplayName)
				} else {
					tfModelOauth2PermissionScopes.UserConsentDisplayName = types.StringNull()
				}

				if !tfModelOauth2PermissionScopes.Value.IsUnknown() {
					tfPlanValue := tfModelOauth2PermissionScopes.Value.ValueString()
					sdkModelOauth2PermissionScopes.SetValue(&tfPlanValue)
				} else {
					tfModelOauth2PermissionScopes.Value = types.StringNull()
				}
			}
			sdkModelApi.SetOauth2PermissionScopes(tfPlanOauth2PermissionScopes)
		} else {
			tfModelApi.Oauth2PermissionScopes = types.ListNull(tfModelApi.Oauth2PermissionScopes.ElementType(ctx))
		}

		if len(tfModelApi.PreAuthorizedApplications.Elements()) > 0 {
			var tfPlanPreAuthorizedApplications []models.PreAuthorizedApplicationable
			for _, i := range tfModelApi.PreAuthorizedApplications.Elements() {
				sdkModelPreAuthorizedApplications := models.NewPreAuthorizedApplication()
				tfModelPreAuthorizedApplications := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfModelPreAuthorizedApplications)

				if !tfModelPreAuthorizedApplications.AppId.IsUnknown() {
					tfPlanAppId := tfModelPreAuthorizedApplications.AppId.ValueString()
					sdkModelPreAuthorizedApplications.SetAppId(&tfPlanAppId)
				} else {
					tfModelPreAuthorizedApplications.AppId = types.StringNull()
				}

				if len(tfModelPreAuthorizedApplications.DelegatedPermissionIds.Elements()) > 0 {
					var stringArrayDelegatedPermissionIds []string
					for _, i := range tfModelPreAuthorizedApplications.DelegatedPermissionIds.Elements() {
						stringArrayDelegatedPermissionIds = append(stringArrayDelegatedPermissionIds, i.String())
					}
					sdkModelPreAuthorizedApplications.SetDelegatedPermissionIds(stringArrayDelegatedPermissionIds)
				} else {
					tfModelPreAuthorizedApplications.DelegatedPermissionIds = types.ListNull(types.StringType)
				}
			}
			sdkModelApi.SetPreAuthorizedApplications(tfPlanPreAuthorizedApplications)
		} else {
			tfModelApi.PreAuthorizedApplications = types.ListNull(tfModelApi.PreAuthorizedApplications.ElementType(ctx))
		}
		sdkModelApplication.SetApi(sdkModelApi)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelApi.AttributeTypes(), sdkModelApi)
		tfPlan.Api = objectValue
	} else {
		tfPlan.Api = types.ObjectNull(tfPlan.Api.AttributeTypes(ctx))
	}

	if !tfPlan.AppId.IsUnknown() {
		tfPlanAppId := tfPlan.AppId.ValueString()
		sdkModelApplication.SetAppId(&tfPlanAppId)
	} else {
		tfPlan.AppId = types.StringNull()
	}

	if len(tfPlan.AppRoles.Elements()) > 0 {
		var tfPlanAppRoles []models.AppRoleable
		for _, i := range tfPlan.AppRoles.Elements() {
			sdkModelAppRoles := models.NewAppRole()
			tfModelAppRoles := applicationAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelAppRoles)

			if len(tfModelAppRoles.AllowedMemberTypes.Elements()) > 0 {
				var stringArrayAllowedMemberTypes []string
				for _, i := range tfModelAppRoles.AllowedMemberTypes.Elements() {
					stringArrayAllowedMemberTypes = append(stringArrayAllowedMemberTypes, i.String())
				}
				sdkModelAppRoles.SetAllowedMemberTypes(stringArrayAllowedMemberTypes)
			} else {
				tfModelAppRoles.AllowedMemberTypes = types.ListNull(types.StringType)
			}

			if !tfModelAppRoles.Description.IsUnknown() {
				tfPlanDescription := tfModelAppRoles.Description.ValueString()
				sdkModelAppRoles.SetDescription(&tfPlanDescription)
			} else {
				tfModelAppRoles.Description = types.StringNull()
			}

			if !tfModelAppRoles.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfModelAppRoles.DisplayName.ValueString()
				sdkModelAppRoles.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfModelAppRoles.DisplayName = types.StringNull()
			}

			if !tfModelAppRoles.Id.IsUnknown() {
				tfPlanId := tfModelAppRoles.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				sdkModelAppRoles.SetId(&u)
			} else {
				tfModelAppRoles.Id = types.StringNull()
			}

			if !tfModelAppRoles.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := tfModelAppRoles.IsEnabled.ValueBool()
				sdkModelAppRoles.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				tfModelAppRoles.IsEnabled = types.BoolNull()
			}

			if !tfModelAppRoles.Origin.IsUnknown() {
				tfPlanOrigin := tfModelAppRoles.Origin.ValueString()
				sdkModelAppRoles.SetOrigin(&tfPlanOrigin)
			} else {
				tfModelAppRoles.Origin = types.StringNull()
			}

			if !tfModelAppRoles.Value.IsUnknown() {
				tfPlanValue := tfModelAppRoles.Value.ValueString()
				sdkModelAppRoles.SetValue(&tfPlanValue)
			} else {
				tfModelAppRoles.Value = types.StringNull()
			}
		}
		sdkModelApplication.SetAppRoles(tfPlanAppRoles)
	} else {
		tfPlan.AppRoles = types.ListNull(tfPlan.AppRoles.ElementType(ctx))
	}

	if !tfPlan.ApplicationTemplateId.IsUnknown() {
		tfPlanApplicationTemplateId := tfPlan.ApplicationTemplateId.ValueString()
		sdkModelApplication.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	} else {
		tfPlan.ApplicationTemplateId = types.StringNull()
	}

	if !tfPlan.Certification.IsUnknown() {
		sdkModelCertification := models.NewCertification()
		tfModelCertification := applicationCertificationModel{}
		tfPlan.Certification.As(ctx, &tfModelCertification, basetypes.ObjectAsOptions{})

		if !tfModelCertification.CertificationDetailsUrl.IsUnknown() {
			tfPlanCertificationDetailsUrl := tfModelCertification.CertificationDetailsUrl.ValueString()
			sdkModelCertification.SetCertificationDetailsUrl(&tfPlanCertificationDetailsUrl)
		} else {
			tfModelCertification.CertificationDetailsUrl = types.StringNull()
		}

		if !tfModelCertification.CertificationExpirationDateTime.IsUnknown() {
			tfPlanCertificationExpirationDateTime := tfModelCertification.CertificationExpirationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanCertificationExpirationDateTime)
			sdkModelCertification.SetCertificationExpirationDateTime(&t)
		} else {
			tfModelCertification.CertificationExpirationDateTime = types.StringNull()
		}

		if !tfModelCertification.IsCertifiedByMicrosoft.IsUnknown() {
			tfPlanIsCertifiedByMicrosoft := tfModelCertification.IsCertifiedByMicrosoft.ValueBool()
			sdkModelCertification.SetIsCertifiedByMicrosoft(&tfPlanIsCertifiedByMicrosoft)
		} else {
			tfModelCertification.IsCertifiedByMicrosoft = types.BoolNull()
		}

		if !tfModelCertification.IsPublisherAttested.IsUnknown() {
			tfPlanIsPublisherAttested := tfModelCertification.IsPublisherAttested.ValueBool()
			sdkModelCertification.SetIsPublisherAttested(&tfPlanIsPublisherAttested)
		} else {
			tfModelCertification.IsPublisherAttested = types.BoolNull()
		}

		if !tfModelCertification.LastCertificationDateTime.IsUnknown() {
			tfPlanLastCertificationDateTime := tfModelCertification.LastCertificationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastCertificationDateTime)
			sdkModelCertification.SetLastCertificationDateTime(&t)
		} else {
			tfModelCertification.LastCertificationDateTime = types.StringNull()
		}
		sdkModelApplication.SetCertification(sdkModelCertification)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelCertification.AttributeTypes(), sdkModelCertification)
		tfPlan.Certification = objectValue
	} else {
		tfPlan.Certification = types.ObjectNull(tfPlan.Certification.AttributeTypes(ctx))
	}

	if !tfPlan.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		sdkModelApplication.SetCreatedDateTime(&t)
	} else {
		tfPlan.CreatedDateTime = types.StringNull()
	}

	if !tfPlan.DefaultRedirectUri.IsUnknown() {
		tfPlanDefaultRedirectUri := tfPlan.DefaultRedirectUri.ValueString()
		sdkModelApplication.SetDefaultRedirectUri(&tfPlanDefaultRedirectUri)
	} else {
		tfPlan.DefaultRedirectUri = types.StringNull()
	}

	if !tfPlan.Description.IsUnknown() {
		tfPlanDescription := tfPlan.Description.ValueString()
		sdkModelApplication.SetDescription(&tfPlanDescription)
	} else {
		tfPlan.Description = types.StringNull()
	}

	if !tfPlan.DisabledByMicrosoftStatus.IsUnknown() {
		tfPlanDisabledByMicrosoftStatus := tfPlan.DisabledByMicrosoftStatus.ValueString()
		sdkModelApplication.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	} else {
		tfPlan.DisabledByMicrosoftStatus = types.StringNull()
	}

	if !tfPlan.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlan.DisplayName.ValueString()
		sdkModelApplication.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlan.DisplayName = types.StringNull()
	}

	if !tfPlan.GroupMembershipClaims.IsUnknown() {
		tfPlanGroupMembershipClaims := tfPlan.GroupMembershipClaims.ValueString()
		sdkModelApplication.SetGroupMembershipClaims(&tfPlanGroupMembershipClaims)
	} else {
		tfPlan.GroupMembershipClaims = types.StringNull()
	}

	if len(tfPlan.IdentifierUris.Elements()) > 0 {
		var stringArrayIdentifierUris []string
		for _, i := range tfPlan.IdentifierUris.Elements() {
			stringArrayIdentifierUris = append(stringArrayIdentifierUris, i.String())
		}
		sdkModelApplication.SetIdentifierUris(stringArrayIdentifierUris)
	} else {
		tfPlan.IdentifierUris = types.ListNull(types.StringType)
	}

	if !tfPlan.Info.IsUnknown() {
		sdkModelInfo := models.NewInformationalUrl()
		tfModelInfo := applicationInformationalUrlModel{}
		tfPlan.Info.As(ctx, &tfModelInfo, basetypes.ObjectAsOptions{})

		if !tfModelInfo.LogoUrl.IsUnknown() {
			tfPlanLogoUrl := tfModelInfo.LogoUrl.ValueString()
			sdkModelInfo.SetLogoUrl(&tfPlanLogoUrl)
		} else {
			tfModelInfo.LogoUrl = types.StringNull()
		}

		if !tfModelInfo.MarketingUrl.IsUnknown() {
			tfPlanMarketingUrl := tfModelInfo.MarketingUrl.ValueString()
			sdkModelInfo.SetMarketingUrl(&tfPlanMarketingUrl)
		} else {
			tfModelInfo.MarketingUrl = types.StringNull()
		}

		if !tfModelInfo.PrivacyStatementUrl.IsUnknown() {
			tfPlanPrivacyStatementUrl := tfModelInfo.PrivacyStatementUrl.ValueString()
			sdkModelInfo.SetPrivacyStatementUrl(&tfPlanPrivacyStatementUrl)
		} else {
			tfModelInfo.PrivacyStatementUrl = types.StringNull()
		}

		if !tfModelInfo.SupportUrl.IsUnknown() {
			tfPlanSupportUrl := tfModelInfo.SupportUrl.ValueString()
			sdkModelInfo.SetSupportUrl(&tfPlanSupportUrl)
		} else {
			tfModelInfo.SupportUrl = types.StringNull()
		}

		if !tfModelInfo.TermsOfServiceUrl.IsUnknown() {
			tfPlanTermsOfServiceUrl := tfModelInfo.TermsOfServiceUrl.ValueString()
			sdkModelInfo.SetTermsOfServiceUrl(&tfPlanTermsOfServiceUrl)
		} else {
			tfModelInfo.TermsOfServiceUrl = types.StringNull()
		}
		sdkModelApplication.SetInfo(sdkModelInfo)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelInfo.AttributeTypes(), sdkModelInfo)
		tfPlan.Info = objectValue
	} else {
		tfPlan.Info = types.ObjectNull(tfPlan.Info.AttributeTypes(ctx))
	}

	if !tfPlan.IsDeviceOnlyAuthSupported.IsUnknown() {
		tfPlanIsDeviceOnlyAuthSupported := tfPlan.IsDeviceOnlyAuthSupported.ValueBool()
		sdkModelApplication.SetIsDeviceOnlyAuthSupported(&tfPlanIsDeviceOnlyAuthSupported)
	} else {
		tfPlan.IsDeviceOnlyAuthSupported = types.BoolNull()
	}

	if !tfPlan.IsFallbackPublicClient.IsUnknown() {
		tfPlanIsFallbackPublicClient := tfPlan.IsFallbackPublicClient.ValueBool()
		sdkModelApplication.SetIsFallbackPublicClient(&tfPlanIsFallbackPublicClient)
	} else {
		tfPlan.IsFallbackPublicClient = types.BoolNull()
	}

	if len(tfPlan.KeyCredentials.Elements()) > 0 {
		var tfPlanKeyCredentials []models.KeyCredentialable
		for _, i := range tfPlan.KeyCredentials.Elements() {
			sdkModelKeyCredentials := models.NewKeyCredential()
			tfModelKeyCredentials := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelKeyCredentials)

			if !tfModelKeyCredentials.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := tfModelKeyCredentials.CustomKeyIdentifier.ValueString()
				sdkModelKeyCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				tfModelKeyCredentials.CustomKeyIdentifier = types.StringNull()
			}

			if !tfModelKeyCredentials.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfModelKeyCredentials.DisplayName.ValueString()
				sdkModelKeyCredentials.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfModelKeyCredentials.DisplayName = types.StringNull()
			}

			if !tfModelKeyCredentials.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := tfModelKeyCredentials.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				sdkModelKeyCredentials.SetEndDateTime(&t)
			} else {
				tfModelKeyCredentials.EndDateTime = types.StringNull()
			}

			if !tfModelKeyCredentials.Key.IsUnknown() {
				tfPlanKey := tfModelKeyCredentials.Key.ValueString()
				sdkModelKeyCredentials.SetKey([]byte(tfPlanKey))
			} else {
				tfModelKeyCredentials.Key = types.StringNull()
			}

			if !tfModelKeyCredentials.KeyId.IsUnknown() {
				tfPlanKeyId := tfModelKeyCredentials.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				sdkModelKeyCredentials.SetKeyId(&u)
			} else {
				tfModelKeyCredentials.KeyId = types.StringNull()
			}

			if !tfModelKeyCredentials.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := tfModelKeyCredentials.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				sdkModelKeyCredentials.SetStartDateTime(&t)
			} else {
				tfModelKeyCredentials.StartDateTime = types.StringNull()
			}

			if !tfModelKeyCredentials.Type.IsUnknown() {
				tfPlanType := tfModelKeyCredentials.Type.ValueString()
				sdkModelKeyCredentials.SetTypeEscaped(&tfPlanType)
			} else {
				tfModelKeyCredentials.Type = types.StringNull()
			}

			if !tfModelKeyCredentials.Usage.IsUnknown() {
				tfPlanUsage := tfModelKeyCredentials.Usage.ValueString()
				sdkModelKeyCredentials.SetUsage(&tfPlanUsage)
			} else {
				tfModelKeyCredentials.Usage = types.StringNull()
			}
		}
		sdkModelApplication.SetKeyCredentials(tfPlanKeyCredentials)
	} else {
		tfPlan.KeyCredentials = types.ListNull(tfPlan.KeyCredentials.ElementType(ctx))
	}

	if !tfPlan.Logo.IsUnknown() {
		tfPlanLogo := tfPlan.Logo.ValueString()
		sdkModelApplication.SetLogo([]byte(tfPlanLogo))
	} else {
		tfPlan.Logo = types.StringNull()
	}

	if !tfPlan.NativeAuthenticationApisEnabled.IsUnknown() {
		tfPlanNativeAuthenticationApisEnabled := tfPlan.NativeAuthenticationApisEnabled.ValueString()
		parsedNativeAuthenticationApisEnabled, _ := models.ParseNativeAuthenticationApisEnabled(tfPlanNativeAuthenticationApisEnabled)
		assertedNativeAuthenticationApisEnabled := parsedNativeAuthenticationApisEnabled.(models.NativeAuthenticationApisEnabled)
		sdkModelApplication.SetNativeAuthenticationApisEnabled(&assertedNativeAuthenticationApisEnabled)
	} else {
		tfPlan.NativeAuthenticationApisEnabled = types.StringNull()
	}

	if !tfPlan.Notes.IsUnknown() {
		tfPlanNotes := tfPlan.Notes.ValueString()
		sdkModelApplication.SetNotes(&tfPlanNotes)
	} else {
		tfPlan.Notes = types.StringNull()
	}

	if !tfPlan.Oauth2RequirePostResponse.IsUnknown() {
		tfPlanOauth2RequirePostResponse := tfPlan.Oauth2RequirePostResponse.ValueBool()
		sdkModelApplication.SetOauth2RequirePostResponse(&tfPlanOauth2RequirePostResponse)
	} else {
		tfPlan.Oauth2RequirePostResponse = types.BoolNull()
	}

	if !tfPlan.OptionalClaims.IsUnknown() {
		sdkModelOptionalClaims := models.NewOptionalClaims()
		tfModelOptionalClaims := applicationOptionalClaimsModel{}
		tfPlan.OptionalClaims.As(ctx, &tfModelOptionalClaims, basetypes.ObjectAsOptions{})

		if len(tfModelOptionalClaims.AccessToken.Elements()) > 0 {
			var tfPlanAccessToken []models.OptionalClaimable
			for _, i := range tfModelOptionalClaims.AccessToken.Elements() {
				sdkModelAccessToken := models.NewOptionalClaim()
				tfModelAccessToken := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfModelAccessToken)

				if len(tfModelAccessToken.AdditionalProperties.Elements()) > 0 {
					var stringArrayAdditionalProperties []string
					for _, i := range tfModelAccessToken.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					sdkModelAccessToken.SetAdditionalProperties(stringArrayAdditionalProperties)
				} else {
					tfModelAccessToken.AdditionalProperties = types.ListNull(types.StringType)
				}

				if !tfModelAccessToken.Essential.IsUnknown() {
					tfPlanEssential := tfModelAccessToken.Essential.ValueBool()
					sdkModelAccessToken.SetEssential(&tfPlanEssential)
				} else {
					tfModelAccessToken.Essential = types.BoolNull()
				}

				if !tfModelAccessToken.Name.IsUnknown() {
					tfPlanName := tfModelAccessToken.Name.ValueString()
					sdkModelAccessToken.SetName(&tfPlanName)
				} else {
					tfModelAccessToken.Name = types.StringNull()
				}

				if !tfModelAccessToken.Source.IsUnknown() {
					tfPlanSource := tfModelAccessToken.Source.ValueString()
					sdkModelAccessToken.SetSource(&tfPlanSource)
				} else {
					tfModelAccessToken.Source = types.StringNull()
				}
			}
			sdkModelOptionalClaims.SetAccessToken(tfPlanAccessToken)
		} else {
			tfModelOptionalClaims.AccessToken = types.ListNull(tfModelOptionalClaims.AccessToken.ElementType(ctx))
		}

		if len(tfModelOptionalClaims.IdToken.Elements()) > 0 {
			var tfPlanIdToken []models.OptionalClaimable
			for _, i := range tfModelOptionalClaims.IdToken.Elements() {
				sdkModelIdToken := models.NewOptionalClaim()
				tfModelIdToken := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfModelIdToken)

				if len(tfModelIdToken.AdditionalProperties.Elements()) > 0 {
					var stringArrayAdditionalProperties []string
					for _, i := range tfModelIdToken.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					sdkModelIdToken.SetAdditionalProperties(stringArrayAdditionalProperties)
				} else {
					tfModelIdToken.AdditionalProperties = types.ListNull(types.StringType)
				}

				if !tfModelIdToken.Essential.IsUnknown() {
					tfPlanEssential := tfModelIdToken.Essential.ValueBool()
					sdkModelIdToken.SetEssential(&tfPlanEssential)
				} else {
					tfModelIdToken.Essential = types.BoolNull()
				}

				if !tfModelIdToken.Name.IsUnknown() {
					tfPlanName := tfModelIdToken.Name.ValueString()
					sdkModelIdToken.SetName(&tfPlanName)
				} else {
					tfModelIdToken.Name = types.StringNull()
				}

				if !tfModelIdToken.Source.IsUnknown() {
					tfPlanSource := tfModelIdToken.Source.ValueString()
					sdkModelIdToken.SetSource(&tfPlanSource)
				} else {
					tfModelIdToken.Source = types.StringNull()
				}
			}
			sdkModelOptionalClaims.SetIdToken(tfPlanIdToken)
		} else {
			tfModelOptionalClaims.IdToken = types.ListNull(tfModelOptionalClaims.IdToken.ElementType(ctx))
		}

		if len(tfModelOptionalClaims.Saml2Token.Elements()) > 0 {
			var tfPlanSaml2Token []models.OptionalClaimable
			for _, i := range tfModelOptionalClaims.Saml2Token.Elements() {
				sdkModelSaml2Token := models.NewOptionalClaim()
				tfModelSaml2Token := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfModelSaml2Token)

				if len(tfModelSaml2Token.AdditionalProperties.Elements()) > 0 {
					var stringArrayAdditionalProperties []string
					for _, i := range tfModelSaml2Token.AdditionalProperties.Elements() {
						stringArrayAdditionalProperties = append(stringArrayAdditionalProperties, i.String())
					}
					sdkModelSaml2Token.SetAdditionalProperties(stringArrayAdditionalProperties)
				} else {
					tfModelSaml2Token.AdditionalProperties = types.ListNull(types.StringType)
				}

				if !tfModelSaml2Token.Essential.IsUnknown() {
					tfPlanEssential := tfModelSaml2Token.Essential.ValueBool()
					sdkModelSaml2Token.SetEssential(&tfPlanEssential)
				} else {
					tfModelSaml2Token.Essential = types.BoolNull()
				}

				if !tfModelSaml2Token.Name.IsUnknown() {
					tfPlanName := tfModelSaml2Token.Name.ValueString()
					sdkModelSaml2Token.SetName(&tfPlanName)
				} else {
					tfModelSaml2Token.Name = types.StringNull()
				}

				if !tfModelSaml2Token.Source.IsUnknown() {
					tfPlanSource := tfModelSaml2Token.Source.ValueString()
					sdkModelSaml2Token.SetSource(&tfPlanSource)
				} else {
					tfModelSaml2Token.Source = types.StringNull()
				}
			}
			sdkModelOptionalClaims.SetSaml2Token(tfPlanSaml2Token)
		} else {
			tfModelOptionalClaims.Saml2Token = types.ListNull(tfModelOptionalClaims.Saml2Token.ElementType(ctx))
		}
		sdkModelApplication.SetOptionalClaims(sdkModelOptionalClaims)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelOptionalClaims.AttributeTypes(), sdkModelOptionalClaims)
		tfPlan.OptionalClaims = objectValue
	} else {
		tfPlan.OptionalClaims = types.ObjectNull(tfPlan.OptionalClaims.AttributeTypes(ctx))
	}

	if !tfPlan.ParentalControlSettings.IsUnknown() {
		sdkModelParentalControlSettings := models.NewParentalControlSettings()
		tfModelParentalControlSettings := applicationParentalControlSettingsModel{}
		tfPlan.ParentalControlSettings.As(ctx, &tfModelParentalControlSettings, basetypes.ObjectAsOptions{})

		if len(tfModelParentalControlSettings.CountriesBlockedForMinors.Elements()) > 0 {
			var stringArrayCountriesBlockedForMinors []string
			for _, i := range tfModelParentalControlSettings.CountriesBlockedForMinors.Elements() {
				stringArrayCountriesBlockedForMinors = append(stringArrayCountriesBlockedForMinors, i.String())
			}
			sdkModelParentalControlSettings.SetCountriesBlockedForMinors(stringArrayCountriesBlockedForMinors)
		} else {
			tfModelParentalControlSettings.CountriesBlockedForMinors = types.ListNull(types.StringType)
		}

		if !tfModelParentalControlSettings.LegalAgeGroupRule.IsUnknown() {
			tfPlanLegalAgeGroupRule := tfModelParentalControlSettings.LegalAgeGroupRule.ValueString()
			sdkModelParentalControlSettings.SetLegalAgeGroupRule(&tfPlanLegalAgeGroupRule)
		} else {
			tfModelParentalControlSettings.LegalAgeGroupRule = types.StringNull()
		}
		sdkModelApplication.SetParentalControlSettings(sdkModelParentalControlSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelParentalControlSettings.AttributeTypes(), sdkModelParentalControlSettings)
		tfPlan.ParentalControlSettings = objectValue
	} else {
		tfPlan.ParentalControlSettings = types.ObjectNull(tfPlan.ParentalControlSettings.AttributeTypes(ctx))
	}

	if len(tfPlan.PasswordCredentials.Elements()) > 0 {
		var tfPlanPasswordCredentials []models.PasswordCredentialable
		for _, i := range tfPlan.PasswordCredentials.Elements() {
			sdkModelPasswordCredentials := models.NewPasswordCredential()
			tfModelPasswordCredentials := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelPasswordCredentials)

			if !tfModelPasswordCredentials.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := tfModelPasswordCredentials.CustomKeyIdentifier.ValueString()
				sdkModelPasswordCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				tfModelPasswordCredentials.CustomKeyIdentifier = types.StringNull()
			}

			if !tfModelPasswordCredentials.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfModelPasswordCredentials.DisplayName.ValueString()
				sdkModelPasswordCredentials.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfModelPasswordCredentials.DisplayName = types.StringNull()
			}

			if !tfModelPasswordCredentials.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := tfModelPasswordCredentials.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				sdkModelPasswordCredentials.SetEndDateTime(&t)
			} else {
				tfModelPasswordCredentials.EndDateTime = types.StringNull()
			}

			if !tfModelPasswordCredentials.Hint.IsUnknown() {
				tfPlanHint := tfModelPasswordCredentials.Hint.ValueString()
				sdkModelPasswordCredentials.SetHint(&tfPlanHint)
			} else {
				tfModelPasswordCredentials.Hint = types.StringNull()
			}

			if !tfModelPasswordCredentials.KeyId.IsUnknown() {
				tfPlanKeyId := tfModelPasswordCredentials.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				sdkModelPasswordCredentials.SetKeyId(&u)
			} else {
				tfModelPasswordCredentials.KeyId = types.StringNull()
			}

			if !tfModelPasswordCredentials.SecretText.IsUnknown() {
				tfPlanSecretText := tfModelPasswordCredentials.SecretText.ValueString()
				sdkModelPasswordCredentials.SetSecretText(&tfPlanSecretText)
			} else {
				tfModelPasswordCredentials.SecretText = types.StringNull()
			}

			if !tfModelPasswordCredentials.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := tfModelPasswordCredentials.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				sdkModelPasswordCredentials.SetStartDateTime(&t)
			} else {
				tfModelPasswordCredentials.StartDateTime = types.StringNull()
			}
		}
		sdkModelApplication.SetPasswordCredentials(tfPlanPasswordCredentials)
	} else {
		tfPlan.PasswordCredentials = types.ListNull(tfPlan.PasswordCredentials.ElementType(ctx))
	}

	if !tfPlan.PublicClient.IsUnknown() {
		sdkModelPublicClient := models.NewPublicClientApplication()
		tfModelPublicClient := applicationPublicClientApplicationModel{}
		tfPlan.PublicClient.As(ctx, &tfModelPublicClient, basetypes.ObjectAsOptions{})

		if len(tfModelPublicClient.RedirectUris.Elements()) > 0 {
			var stringArrayRedirectUris []string
			for _, i := range tfModelPublicClient.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			sdkModelPublicClient.SetRedirectUris(stringArrayRedirectUris)
		} else {
			tfModelPublicClient.RedirectUris = types.ListNull(types.StringType)
		}
		sdkModelApplication.SetPublicClient(sdkModelPublicClient)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelPublicClient.AttributeTypes(), sdkModelPublicClient)
		tfPlan.PublicClient = objectValue
	} else {
		tfPlan.PublicClient = types.ObjectNull(tfPlan.PublicClient.AttributeTypes(ctx))
	}

	if !tfPlan.PublisherDomain.IsUnknown() {
		tfPlanPublisherDomain := tfPlan.PublisherDomain.ValueString()
		sdkModelApplication.SetPublisherDomain(&tfPlanPublisherDomain)
	} else {
		tfPlan.PublisherDomain = types.StringNull()
	}

	if !tfPlan.RequestSignatureVerification.IsUnknown() {
		sdkModelRequestSignatureVerification := models.NewRequestSignatureVerification()
		tfModelRequestSignatureVerification := applicationRequestSignatureVerificationModel{}
		tfPlan.RequestSignatureVerification.As(ctx, &tfModelRequestSignatureVerification, basetypes.ObjectAsOptions{})

		if !tfModelRequestSignatureVerification.AllowedWeakAlgorithms.IsUnknown() {
			tfPlanAllowedWeakAlgorithms := tfModelRequestSignatureVerification.AllowedWeakAlgorithms.ValueString()
			parsedAllowedWeakAlgorithms, _ := models.ParseWeakAlgorithms(tfPlanAllowedWeakAlgorithms)
			assertedAllowedWeakAlgorithms := parsedAllowedWeakAlgorithms.(models.WeakAlgorithms)
			sdkModelRequestSignatureVerification.SetAllowedWeakAlgorithms(&assertedAllowedWeakAlgorithms)
		} else {
			tfModelRequestSignatureVerification.AllowedWeakAlgorithms = types.StringNull()
		}

		if !tfModelRequestSignatureVerification.IsSignedRequestRequired.IsUnknown() {
			tfPlanIsSignedRequestRequired := tfModelRequestSignatureVerification.IsSignedRequestRequired.ValueBool()
			sdkModelRequestSignatureVerification.SetIsSignedRequestRequired(&tfPlanIsSignedRequestRequired)
		} else {
			tfModelRequestSignatureVerification.IsSignedRequestRequired = types.BoolNull()
		}
		sdkModelApplication.SetRequestSignatureVerification(sdkModelRequestSignatureVerification)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelRequestSignatureVerification.AttributeTypes(), sdkModelRequestSignatureVerification)
		tfPlan.RequestSignatureVerification = objectValue
	} else {
		tfPlan.RequestSignatureVerification = types.ObjectNull(tfPlan.RequestSignatureVerification.AttributeTypes(ctx))
	}

	if len(tfPlan.RequiredResourceAccess.Elements()) > 0 {
		var tfPlanRequiredResourceAccess []models.RequiredResourceAccessable
		for _, i := range tfPlan.RequiredResourceAccess.Elements() {
			sdkModelRequiredResourceAccess := models.NewRequiredResourceAccess()
			tfModelRequiredResourceAccess := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelRequiredResourceAccess)

			if len(tfModelRequiredResourceAccess.ResourceAccess.Elements()) > 0 {
				var tfPlanResourceAccess []models.ResourceAccessable
				for _, i := range tfModelRequiredResourceAccess.ResourceAccess.Elements() {
					sdkModelResourceAccess := models.NewResourceAccess()
					tfModelResourceAccess := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfModelResourceAccess)

					if !tfModelResourceAccess.Id.IsUnknown() {
						tfPlanId := tfModelResourceAccess.Id.ValueString()
						u, _ := uuid.Parse(tfPlanId)
						sdkModelResourceAccess.SetId(&u)
					} else {
						tfModelResourceAccess.Id = types.StringNull()
					}

					if !tfModelResourceAccess.Type.IsUnknown() {
						tfPlanType := tfModelResourceAccess.Type.ValueString()
						sdkModelResourceAccess.SetTypeEscaped(&tfPlanType)
					} else {
						tfModelResourceAccess.Type = types.StringNull()
					}
				}
				sdkModelRequiredResourceAccess.SetResourceAccess(tfPlanResourceAccess)
			} else {
				tfModelRequiredResourceAccess.ResourceAccess = types.ListNull(tfModelRequiredResourceAccess.ResourceAccess.ElementType(ctx))
			}

			if !tfModelRequiredResourceAccess.ResourceAppId.IsUnknown() {
				tfPlanResourceAppId := tfModelRequiredResourceAccess.ResourceAppId.ValueString()
				sdkModelRequiredResourceAccess.SetResourceAppId(&tfPlanResourceAppId)
			} else {
				tfModelRequiredResourceAccess.ResourceAppId = types.StringNull()
			}
		}
		sdkModelApplication.SetRequiredResourceAccess(tfPlanRequiredResourceAccess)
	} else {
		tfPlan.RequiredResourceAccess = types.ListNull(tfPlan.RequiredResourceAccess.ElementType(ctx))
	}

	if !tfPlan.SamlMetadataUrl.IsUnknown() {
		tfPlanSamlMetadataUrl := tfPlan.SamlMetadataUrl.ValueString()
		sdkModelApplication.SetSamlMetadataUrl(&tfPlanSamlMetadataUrl)
	} else {
		tfPlan.SamlMetadataUrl = types.StringNull()
	}

	if !tfPlan.ServiceManagementReference.IsUnknown() {
		tfPlanServiceManagementReference := tfPlan.ServiceManagementReference.ValueString()
		sdkModelApplication.SetServiceManagementReference(&tfPlanServiceManagementReference)
	} else {
		tfPlan.ServiceManagementReference = types.StringNull()
	}

	if !tfPlan.ServicePrincipalLockConfiguration.IsUnknown() {
		sdkModelServicePrincipalLockConfiguration := models.NewServicePrincipalLockConfiguration()
		tfModelServicePrincipalLockConfiguration := applicationServicePrincipalLockConfigurationModel{}
		tfPlan.ServicePrincipalLockConfiguration.As(ctx, &tfModelServicePrincipalLockConfiguration, basetypes.ObjectAsOptions{})

		if !tfModelServicePrincipalLockConfiguration.AllProperties.IsUnknown() {
			tfPlanAllProperties := tfModelServicePrincipalLockConfiguration.AllProperties.ValueBool()
			sdkModelServicePrincipalLockConfiguration.SetAllProperties(&tfPlanAllProperties)
		} else {
			tfModelServicePrincipalLockConfiguration.AllProperties = types.BoolNull()
		}

		if !tfModelServicePrincipalLockConfiguration.CredentialsWithUsageSign.IsUnknown() {
			tfPlanCredentialsWithUsageSign := tfModelServicePrincipalLockConfiguration.CredentialsWithUsageSign.ValueBool()
			sdkModelServicePrincipalLockConfiguration.SetCredentialsWithUsageSign(&tfPlanCredentialsWithUsageSign)
		} else {
			tfModelServicePrincipalLockConfiguration.CredentialsWithUsageSign = types.BoolNull()
		}

		if !tfModelServicePrincipalLockConfiguration.CredentialsWithUsageVerify.IsUnknown() {
			tfPlanCredentialsWithUsageVerify := tfModelServicePrincipalLockConfiguration.CredentialsWithUsageVerify.ValueBool()
			sdkModelServicePrincipalLockConfiguration.SetCredentialsWithUsageVerify(&tfPlanCredentialsWithUsageVerify)
		} else {
			tfModelServicePrincipalLockConfiguration.CredentialsWithUsageVerify = types.BoolNull()
		}

		if !tfModelServicePrincipalLockConfiguration.IsEnabled.IsUnknown() {
			tfPlanIsEnabled := tfModelServicePrincipalLockConfiguration.IsEnabled.ValueBool()
			sdkModelServicePrincipalLockConfiguration.SetIsEnabled(&tfPlanIsEnabled)
		} else {
			tfModelServicePrincipalLockConfiguration.IsEnabled = types.BoolNull()
		}

		if !tfModelServicePrincipalLockConfiguration.TokenEncryptionKeyId.IsUnknown() {
			tfPlanTokenEncryptionKeyId := tfModelServicePrincipalLockConfiguration.TokenEncryptionKeyId.ValueBool()
			sdkModelServicePrincipalLockConfiguration.SetTokenEncryptionKeyId(&tfPlanTokenEncryptionKeyId)
		} else {
			tfModelServicePrincipalLockConfiguration.TokenEncryptionKeyId = types.BoolNull()
		}
		sdkModelApplication.SetServicePrincipalLockConfiguration(sdkModelServicePrincipalLockConfiguration)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelServicePrincipalLockConfiguration.AttributeTypes(), sdkModelServicePrincipalLockConfiguration)
		tfPlan.ServicePrincipalLockConfiguration = objectValue
	} else {
		tfPlan.ServicePrincipalLockConfiguration = types.ObjectNull(tfPlan.ServicePrincipalLockConfiguration.AttributeTypes(ctx))
	}

	if !tfPlan.SignInAudience.IsUnknown() {
		tfPlanSignInAudience := tfPlan.SignInAudience.ValueString()
		sdkModelApplication.SetSignInAudience(&tfPlanSignInAudience)
	} else {
		tfPlan.SignInAudience = types.StringNull()
	}

	if !tfPlan.Spa.IsUnknown() {
		sdkModelSpa := models.NewSpaApplication()
		tfModelSpa := applicationSpaApplicationModel{}
		tfPlan.Spa.As(ctx, &tfModelSpa, basetypes.ObjectAsOptions{})

		if len(tfModelSpa.RedirectUris.Elements()) > 0 {
			var stringArrayRedirectUris []string
			for _, i := range tfModelSpa.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			sdkModelSpa.SetRedirectUris(stringArrayRedirectUris)
		} else {
			tfModelSpa.RedirectUris = types.ListNull(types.StringType)
		}
		sdkModelApplication.SetSpa(sdkModelSpa)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelSpa.AttributeTypes(), sdkModelSpa)
		tfPlan.Spa = objectValue
	} else {
		tfPlan.Spa = types.ObjectNull(tfPlan.Spa.AttributeTypes(ctx))
	}

	if len(tfPlan.Tags.Elements()) > 0 {
		var stringArrayTags []string
		for _, i := range tfPlan.Tags.Elements() {
			stringArrayTags = append(stringArrayTags, i.String())
		}
		sdkModelApplication.SetTags(stringArrayTags)
	} else {
		tfPlan.Tags = types.ListNull(types.StringType)
	}

	if !tfPlan.TokenEncryptionKeyId.IsUnknown() {
		tfPlanTokenEncryptionKeyId := tfPlan.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		sdkModelApplication.SetTokenEncryptionKeyId(&u)
	} else {
		tfPlan.TokenEncryptionKeyId = types.StringNull()
	}

	if !tfPlan.UniqueName.IsUnknown() {
		tfPlanUniqueName := tfPlan.UniqueName.ValueString()
		sdkModelApplication.SetUniqueName(&tfPlanUniqueName)
	} else {
		tfPlan.UniqueName = types.StringNull()
	}

	if !tfPlan.VerifiedPublisher.IsUnknown() {
		sdkModelVerifiedPublisher := models.NewVerifiedPublisher()
		tfModelVerifiedPublisher := applicationVerifiedPublisherModel{}
		tfPlan.VerifiedPublisher.As(ctx, &tfModelVerifiedPublisher, basetypes.ObjectAsOptions{})

		if !tfModelVerifiedPublisher.AddedDateTime.IsUnknown() {
			tfPlanAddedDateTime := tfModelVerifiedPublisher.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanAddedDateTime)
			sdkModelVerifiedPublisher.SetAddedDateTime(&t)
		} else {
			tfModelVerifiedPublisher.AddedDateTime = types.StringNull()
		}

		if !tfModelVerifiedPublisher.DisplayName.IsUnknown() {
			tfPlanDisplayName := tfModelVerifiedPublisher.DisplayName.ValueString()
			sdkModelVerifiedPublisher.SetDisplayName(&tfPlanDisplayName)
		} else {
			tfModelVerifiedPublisher.DisplayName = types.StringNull()
		}

		if !tfModelVerifiedPublisher.VerifiedPublisherId.IsUnknown() {
			tfPlanVerifiedPublisherId := tfModelVerifiedPublisher.VerifiedPublisherId.ValueString()
			sdkModelVerifiedPublisher.SetVerifiedPublisherId(&tfPlanVerifiedPublisherId)
		} else {
			tfModelVerifiedPublisher.VerifiedPublisherId = types.StringNull()
		}
		sdkModelApplication.SetVerifiedPublisher(sdkModelVerifiedPublisher)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelVerifiedPublisher.AttributeTypes(), sdkModelVerifiedPublisher)
		tfPlan.VerifiedPublisher = objectValue
	} else {
		tfPlan.VerifiedPublisher = types.ObjectNull(tfPlan.VerifiedPublisher.AttributeTypes(ctx))
	}

	if !tfPlan.Web.IsUnknown() {
		sdkModelWeb := models.NewWebApplication()
		tfModelWeb := applicationWebApplicationModel{}
		tfPlan.Web.As(ctx, &tfModelWeb, basetypes.ObjectAsOptions{})

		if !tfModelWeb.HomePageUrl.IsUnknown() {
			tfPlanHomePageUrl := tfModelWeb.HomePageUrl.ValueString()
			sdkModelWeb.SetHomePageUrl(&tfPlanHomePageUrl)
		} else {
			tfModelWeb.HomePageUrl = types.StringNull()
		}

		if !tfModelWeb.ImplicitGrantSettings.IsUnknown() {
			sdkModelImplicitGrantSettings := models.NewImplicitGrantSettings()
			tfModelImplicitGrantSettings := applicationImplicitGrantSettingsModel{}
			tfModelWeb.ImplicitGrantSettings.As(ctx, &tfModelImplicitGrantSettings, basetypes.ObjectAsOptions{})

			if !tfModelImplicitGrantSettings.EnableAccessTokenIssuance.IsUnknown() {
				tfPlanEnableAccessTokenIssuance := tfModelImplicitGrantSettings.EnableAccessTokenIssuance.ValueBool()
				sdkModelImplicitGrantSettings.SetEnableAccessTokenIssuance(&tfPlanEnableAccessTokenIssuance)
			} else {
				tfModelImplicitGrantSettings.EnableAccessTokenIssuance = types.BoolNull()
			}

			if !tfModelImplicitGrantSettings.EnableIdTokenIssuance.IsUnknown() {
				tfPlanEnableIdTokenIssuance := tfModelImplicitGrantSettings.EnableIdTokenIssuance.ValueBool()
				sdkModelImplicitGrantSettings.SetEnableIdTokenIssuance(&tfPlanEnableIdTokenIssuance)
			} else {
				tfModelImplicitGrantSettings.EnableIdTokenIssuance = types.BoolNull()
			}
			sdkModelWeb.SetImplicitGrantSettings(sdkModelImplicitGrantSettings)
			objectValue, _ := types.ObjectValueFrom(ctx, tfModelImplicitGrantSettings.AttributeTypes(), sdkModelImplicitGrantSettings)
			tfModelWeb.ImplicitGrantSettings = objectValue
		} else {
			tfModelWeb.ImplicitGrantSettings = types.ObjectNull(tfModelWeb.ImplicitGrantSettings.AttributeTypes(ctx))
		}

		if !tfModelWeb.LogoutUrl.IsUnknown() {
			tfPlanLogoutUrl := tfModelWeb.LogoutUrl.ValueString()
			sdkModelWeb.SetLogoutUrl(&tfPlanLogoutUrl)
		} else {
			tfModelWeb.LogoutUrl = types.StringNull()
		}

		if len(tfModelWeb.RedirectUriSettings.Elements()) > 0 {
			var tfPlanRedirectUriSettings []models.RedirectUriSettingsable
			for _, i := range tfModelWeb.RedirectUriSettings.Elements() {
				sdkModelRedirectUriSettings := models.NewRedirectUriSettings()
				tfModelRedirectUriSettings := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &tfModelRedirectUriSettings)

				if !tfModelRedirectUriSettings.Uri.IsUnknown() {
					tfPlanUri := tfModelRedirectUriSettings.Uri.ValueString()
					sdkModelRedirectUriSettings.SetUri(&tfPlanUri)
				} else {
					tfModelRedirectUriSettings.Uri = types.StringNull()
				}
			}
			sdkModelWeb.SetRedirectUriSettings(tfPlanRedirectUriSettings)
		} else {
			tfModelWeb.RedirectUriSettings = types.ListNull(tfModelWeb.RedirectUriSettings.ElementType(ctx))
		}

		if len(tfModelWeb.RedirectUris.Elements()) > 0 {
			var stringArrayRedirectUris []string
			for _, i := range tfModelWeb.RedirectUris.Elements() {
				stringArrayRedirectUris = append(stringArrayRedirectUris, i.String())
			}
			sdkModelWeb.SetRedirectUris(stringArrayRedirectUris)
		} else {
			tfModelWeb.RedirectUris = types.ListNull(types.StringType)
		}
		sdkModelApplication.SetWeb(sdkModelWeb)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelWeb.AttributeTypes(), sdkModelWeb)
		tfPlan.Web = objectValue
	} else {
		tfPlan.Web = types.ObjectNull(tfPlan.Web.AttributeTypes(ctx))
	}

	// Create new application
	result, err := r.client.Applications().Post(context.Background(), sdkModelApplication, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating application",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlan.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlan)
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
	// Retrieve values from plan
	var plan applicationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state applicationModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewApplication()

	if !plan.Id.Equal(state.Id) {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	}

	if !plan.DeletedDateTime.Equal(state.DeletedDateTime) {
		planDeletedDateTime := plan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !plan.AddIns.Equal(state.AddIns) {
		var planAddIns []models.AddInable
		for k, i := range plan.AddIns.Elements() {
			addIns := models.NewAddIn()
			addInsModel := applicationAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &addInsModel)
			addInsState := applicationAddInModel{}
			types.ListValueFrom(ctx, state.AddIns.Elements()[k].Type(ctx), &addInsModel)

			if !addInsModel.Id.Equal(addInsState.Id) {
				planId := addInsModel.Id.ValueString()
				u, _ := uuid.Parse(planId)
				addIns.SetId(&u)
			}

			if !addInsModel.Properties.Equal(addInsState.Properties) {
				var planProperties []models.KeyValueable
				for k, i := range addInsModel.Properties.Elements() {
					properties := models.NewKeyValue()
					propertiesModel := applicationKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &propertiesModel)
					propertiesState := applicationKeyValueModel{}
					types.ListValueFrom(ctx, addInsState.Properties.Elements()[k].Type(ctx), &propertiesModel)

					if !propertiesModel.Key.Equal(propertiesState.Key) {
						planKey := propertiesModel.Key.ValueString()
						properties.SetKey(&planKey)
					}

					if !propertiesModel.Value.Equal(propertiesState.Value) {
						planValue := propertiesModel.Value.ValueString()
						properties.SetValue(&planValue)
					}
				}
				addIns.SetProperties(planProperties)
			}

			if !addInsModel.Type.Equal(addInsState.Type) {
				planType := addInsModel.Type.ValueString()
				addIns.SetTypeEscaped(&planType)
			}
		}
		requestBody.SetAddIns(planAddIns)
	}

	if !plan.Api.Equal(state.Api) {
		api := models.NewApiApplication()
		apiModel := applicationApiApplicationModel{}
		plan.Api.As(ctx, &apiModel, basetypes.ObjectAsOptions{})
		apiState := applicationApiApplicationModel{}
		state.Api.As(ctx, &apiState, basetypes.ObjectAsOptions{})

		if !apiModel.AcceptMappedClaims.Equal(apiState.AcceptMappedClaims) {
			planAcceptMappedClaims := apiModel.AcceptMappedClaims.ValueBool()
			api.SetAcceptMappedClaims(&planAcceptMappedClaims)
		}

		if !apiModel.KnownClientApplications.Equal(apiState.KnownClientApplications) {
			var KnownClientApplications []uuid.UUID
			for _, i := range apiModel.KnownClientApplications.Elements() {
				u, _ := uuid.Parse(i.String())
				KnownClientApplications = append(KnownClientApplications, u)
			}
			api.SetKnownClientApplications(KnownClientApplications)
		}

		if !apiModel.Oauth2PermissionScopes.Equal(apiState.Oauth2PermissionScopes) {
			var planOauth2PermissionScopes []models.PermissionScopeable
			for k, i := range apiModel.Oauth2PermissionScopes.Elements() {
				oauth2PermissionScopes := models.NewPermissionScope()
				oauth2PermissionScopesModel := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &oauth2PermissionScopesModel)
				oauth2PermissionScopesState := applicationPermissionScopeModel{}
				types.ListValueFrom(ctx, apiState.Oauth2PermissionScopes.Elements()[k].Type(ctx), &oauth2PermissionScopesModel)

				if !oauth2PermissionScopesModel.AdminConsentDescription.Equal(oauth2PermissionScopesState.AdminConsentDescription) {
					planAdminConsentDescription := oauth2PermissionScopesModel.AdminConsentDescription.ValueString()
					oauth2PermissionScopes.SetAdminConsentDescription(&planAdminConsentDescription)
				}

				if !oauth2PermissionScopesModel.AdminConsentDisplayName.Equal(oauth2PermissionScopesState.AdminConsentDisplayName) {
					planAdminConsentDisplayName := oauth2PermissionScopesModel.AdminConsentDisplayName.ValueString()
					oauth2PermissionScopes.SetAdminConsentDisplayName(&planAdminConsentDisplayName)
				}

				if !oauth2PermissionScopesModel.Id.Equal(oauth2PermissionScopesState.Id) {
					planId := oauth2PermissionScopesModel.Id.ValueString()
					u, _ := uuid.Parse(planId)
					oauth2PermissionScopes.SetId(&u)
				}

				if !oauth2PermissionScopesModel.IsEnabled.Equal(oauth2PermissionScopesState.IsEnabled) {
					planIsEnabled := oauth2PermissionScopesModel.IsEnabled.ValueBool()
					oauth2PermissionScopes.SetIsEnabled(&planIsEnabled)
				}

				if !oauth2PermissionScopesModel.Origin.Equal(oauth2PermissionScopesState.Origin) {
					planOrigin := oauth2PermissionScopesModel.Origin.ValueString()
					oauth2PermissionScopes.SetOrigin(&planOrigin)
				}

				if !oauth2PermissionScopesModel.Type.Equal(oauth2PermissionScopesState.Type) {
					planType := oauth2PermissionScopesModel.Type.ValueString()
					oauth2PermissionScopes.SetTypeEscaped(&planType)
				}

				if !oauth2PermissionScopesModel.UserConsentDescription.Equal(oauth2PermissionScopesState.UserConsentDescription) {
					planUserConsentDescription := oauth2PermissionScopesModel.UserConsentDescription.ValueString()
					oauth2PermissionScopes.SetUserConsentDescription(&planUserConsentDescription)
				}

				if !oauth2PermissionScopesModel.UserConsentDisplayName.Equal(oauth2PermissionScopesState.UserConsentDisplayName) {
					planUserConsentDisplayName := oauth2PermissionScopesModel.UserConsentDisplayName.ValueString()
					oauth2PermissionScopes.SetUserConsentDisplayName(&planUserConsentDisplayName)
				}

				if !oauth2PermissionScopesModel.Value.Equal(oauth2PermissionScopesState.Value) {
					planValue := oauth2PermissionScopesModel.Value.ValueString()
					oauth2PermissionScopes.SetValue(&planValue)
				}
			}
			api.SetOauth2PermissionScopes(planOauth2PermissionScopes)
		}

		if !apiModel.PreAuthorizedApplications.Equal(apiState.PreAuthorizedApplications) {
			var planPreAuthorizedApplications []models.PreAuthorizedApplicationable
			for k, i := range apiModel.PreAuthorizedApplications.Elements() {
				preAuthorizedApplications := models.NewPreAuthorizedApplication()
				preAuthorizedApplicationsModel := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &preAuthorizedApplicationsModel)
				preAuthorizedApplicationsState := applicationPreAuthorizedApplicationModel{}
				types.ListValueFrom(ctx, apiState.PreAuthorizedApplications.Elements()[k].Type(ctx), &preAuthorizedApplicationsModel)

				if !preAuthorizedApplicationsModel.AppId.Equal(preAuthorizedApplicationsState.AppId) {
					planAppId := preAuthorizedApplicationsModel.AppId.ValueString()
					preAuthorizedApplications.SetAppId(&planAppId)
				}

				if !preAuthorizedApplicationsModel.DelegatedPermissionIds.Equal(preAuthorizedApplicationsState.DelegatedPermissionIds) {
					var delegatedPermissionIds []string
					for _, i := range preAuthorizedApplicationsModel.DelegatedPermissionIds.Elements() {
						delegatedPermissionIds = append(delegatedPermissionIds, i.String())
					}
					preAuthorizedApplications.SetDelegatedPermissionIds(delegatedPermissionIds)
				}
			}
			api.SetPreAuthorizedApplications(planPreAuthorizedApplications)
		}
		requestBody.SetApi(api)
		objectValue, _ := types.ObjectValueFrom(ctx, apiModel.AttributeTypes(), apiModel)
		plan.Api = objectValue
	}

	if !plan.AppId.Equal(state.AppId) {
		planAppId := plan.AppId.ValueString()
		requestBody.SetAppId(&planAppId)
	}

	if !plan.AppRoles.Equal(state.AppRoles) {
		var planAppRoles []models.AppRoleable
		for k, i := range plan.AppRoles.Elements() {
			appRoles := models.NewAppRole()
			appRolesModel := applicationAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &appRolesModel)
			appRolesState := applicationAppRoleModel{}
			types.ListValueFrom(ctx, state.AppRoles.Elements()[k].Type(ctx), &appRolesModel)

			if !appRolesModel.AllowedMemberTypes.Equal(appRolesState.AllowedMemberTypes) {
				var allowedMemberTypes []string
				for _, i := range appRolesModel.AllowedMemberTypes.Elements() {
					allowedMemberTypes = append(allowedMemberTypes, i.String())
				}
				appRoles.SetAllowedMemberTypes(allowedMemberTypes)
			}

			if !appRolesModel.Description.Equal(appRolesState.Description) {
				planDescription := appRolesModel.Description.ValueString()
				appRoles.SetDescription(&planDescription)
			}

			if !appRolesModel.DisplayName.Equal(appRolesState.DisplayName) {
				planDisplayName := appRolesModel.DisplayName.ValueString()
				appRoles.SetDisplayName(&planDisplayName)
			}

			if !appRolesModel.Id.Equal(appRolesState.Id) {
				planId := appRolesModel.Id.ValueString()
				u, _ := uuid.Parse(planId)
				appRoles.SetId(&u)
			}

			if !appRolesModel.IsEnabled.Equal(appRolesState.IsEnabled) {
				planIsEnabled := appRolesModel.IsEnabled.ValueBool()
				appRoles.SetIsEnabled(&planIsEnabled)
			}

			if !appRolesModel.Origin.Equal(appRolesState.Origin) {
				planOrigin := appRolesModel.Origin.ValueString()
				appRoles.SetOrigin(&planOrigin)
			}

			if !appRolesModel.Value.Equal(appRolesState.Value) {
				planValue := appRolesModel.Value.ValueString()
				appRoles.SetValue(&planValue)
			}
		}
		requestBody.SetAppRoles(planAppRoles)
	}

	if !plan.ApplicationTemplateId.Equal(state.ApplicationTemplateId) {
		planApplicationTemplateId := plan.ApplicationTemplateId.ValueString()
		requestBody.SetApplicationTemplateId(&planApplicationTemplateId)
	}

	if !plan.Certification.Equal(state.Certification) {
		certification := models.NewCertification()
		certificationModel := applicationCertificationModel{}
		plan.Certification.As(ctx, &certificationModel, basetypes.ObjectAsOptions{})
		certificationState := applicationCertificationModel{}
		state.Certification.As(ctx, &certificationState, basetypes.ObjectAsOptions{})

		if !certificationModel.CertificationDetailsUrl.Equal(certificationState.CertificationDetailsUrl) {
			planCertificationDetailsUrl := certificationModel.CertificationDetailsUrl.ValueString()
			certification.SetCertificationDetailsUrl(&planCertificationDetailsUrl)
		}

		if !certificationModel.CertificationExpirationDateTime.Equal(certificationState.CertificationExpirationDateTime) {
			planCertificationExpirationDateTime := certificationModel.CertificationExpirationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, planCertificationExpirationDateTime)
			certification.SetCertificationExpirationDateTime(&t)
		}

		if !certificationModel.IsCertifiedByMicrosoft.Equal(certificationState.IsCertifiedByMicrosoft) {
			planIsCertifiedByMicrosoft := certificationModel.IsCertifiedByMicrosoft.ValueBool()
			certification.SetIsCertifiedByMicrosoft(&planIsCertifiedByMicrosoft)
		}

		if !certificationModel.IsPublisherAttested.Equal(certificationState.IsPublisherAttested) {
			planIsPublisherAttested := certificationModel.IsPublisherAttested.ValueBool()
			certification.SetIsPublisherAttested(&planIsPublisherAttested)
		}

		if !certificationModel.LastCertificationDateTime.Equal(certificationState.LastCertificationDateTime) {
			planLastCertificationDateTime := certificationModel.LastCertificationDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, planLastCertificationDateTime)
			certification.SetLastCertificationDateTime(&t)
		}
		requestBody.SetCertification(certification)
		objectValue, _ := types.ObjectValueFrom(ctx, certificationModel.AttributeTypes(), certificationModel)
		plan.Certification = objectValue
	}

	if !plan.CreatedDateTime.Equal(state.CreatedDateTime) {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	}

	if !plan.DefaultRedirectUri.Equal(state.DefaultRedirectUri) {
		planDefaultRedirectUri := plan.DefaultRedirectUri.ValueString()
		requestBody.SetDefaultRedirectUri(&planDefaultRedirectUri)
	}

	if !plan.Description.Equal(state.Description) {
		planDescription := plan.Description.ValueString()
		requestBody.SetDescription(&planDescription)
	}

	if !plan.DisabledByMicrosoftStatus.Equal(state.DisabledByMicrosoftStatus) {
		planDisabledByMicrosoftStatus := plan.DisabledByMicrosoftStatus.ValueString()
		requestBody.SetDisabledByMicrosoftStatus(&planDisabledByMicrosoftStatus)
	}

	if !plan.DisplayName.Equal(state.DisplayName) {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	}

	if !plan.GroupMembershipClaims.Equal(state.GroupMembershipClaims) {
		planGroupMembershipClaims := plan.GroupMembershipClaims.ValueString()
		requestBody.SetGroupMembershipClaims(&planGroupMembershipClaims)
	}

	if !plan.IdentifierUris.Equal(state.IdentifierUris) {
		var identifierUris []string
		for _, i := range plan.IdentifierUris.Elements() {
			identifierUris = append(identifierUris, i.String())
		}
		requestBody.SetIdentifierUris(identifierUris)
	}

	if !plan.Info.Equal(state.Info) {
		info := models.NewInformationalUrl()
		infoModel := applicationInformationalUrlModel{}
		plan.Info.As(ctx, &infoModel, basetypes.ObjectAsOptions{})
		infoState := applicationInformationalUrlModel{}
		state.Info.As(ctx, &infoState, basetypes.ObjectAsOptions{})

		if !infoModel.LogoUrl.Equal(infoState.LogoUrl) {
			planLogoUrl := infoModel.LogoUrl.ValueString()
			info.SetLogoUrl(&planLogoUrl)
		}

		if !infoModel.MarketingUrl.Equal(infoState.MarketingUrl) {
			planMarketingUrl := infoModel.MarketingUrl.ValueString()
			info.SetMarketingUrl(&planMarketingUrl)
		}

		if !infoModel.PrivacyStatementUrl.Equal(infoState.PrivacyStatementUrl) {
			planPrivacyStatementUrl := infoModel.PrivacyStatementUrl.ValueString()
			info.SetPrivacyStatementUrl(&planPrivacyStatementUrl)
		}

		if !infoModel.SupportUrl.Equal(infoState.SupportUrl) {
			planSupportUrl := infoModel.SupportUrl.ValueString()
			info.SetSupportUrl(&planSupportUrl)
		}

		if !infoModel.TermsOfServiceUrl.Equal(infoState.TermsOfServiceUrl) {
			planTermsOfServiceUrl := infoModel.TermsOfServiceUrl.ValueString()
			info.SetTermsOfServiceUrl(&planTermsOfServiceUrl)
		}
		requestBody.SetInfo(info)
		objectValue, _ := types.ObjectValueFrom(ctx, infoModel.AttributeTypes(), infoModel)
		plan.Info = objectValue
	}

	if !plan.IsDeviceOnlyAuthSupported.Equal(state.IsDeviceOnlyAuthSupported) {
		planIsDeviceOnlyAuthSupported := plan.IsDeviceOnlyAuthSupported.ValueBool()
		requestBody.SetIsDeviceOnlyAuthSupported(&planIsDeviceOnlyAuthSupported)
	}

	if !plan.IsFallbackPublicClient.Equal(state.IsFallbackPublicClient) {
		planIsFallbackPublicClient := plan.IsFallbackPublicClient.ValueBool()
		requestBody.SetIsFallbackPublicClient(&planIsFallbackPublicClient)
	}

	if !plan.KeyCredentials.Equal(state.KeyCredentials) {
		var planKeyCredentials []models.KeyCredentialable
		for k, i := range plan.KeyCredentials.Elements() {
			keyCredentials := models.NewKeyCredential()
			keyCredentialsModel := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &keyCredentialsModel)
			keyCredentialsState := applicationKeyCredentialModel{}
			types.ListValueFrom(ctx, state.KeyCredentials.Elements()[k].Type(ctx), &keyCredentialsModel)

			if !keyCredentialsModel.CustomKeyIdentifier.Equal(keyCredentialsState.CustomKeyIdentifier) {
				planCustomKeyIdentifier := keyCredentialsModel.CustomKeyIdentifier.ValueString()
				keyCredentials.SetCustomKeyIdentifier([]byte(planCustomKeyIdentifier))
			}

			if !keyCredentialsModel.DisplayName.Equal(keyCredentialsState.DisplayName) {
				planDisplayName := keyCredentialsModel.DisplayName.ValueString()
				keyCredentials.SetDisplayName(&planDisplayName)
			}

			if !keyCredentialsModel.EndDateTime.Equal(keyCredentialsState.EndDateTime) {
				planEndDateTime := keyCredentialsModel.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planEndDateTime)
				keyCredentials.SetEndDateTime(&t)
			}

			if !keyCredentialsModel.Key.Equal(keyCredentialsState.Key) {
				planKey := keyCredentialsModel.Key.ValueString()
				keyCredentials.SetKey([]byte(planKey))
			}

			if !keyCredentialsModel.KeyId.Equal(keyCredentialsState.KeyId) {
				planKeyId := keyCredentialsModel.KeyId.ValueString()
				u, _ := uuid.Parse(planKeyId)
				keyCredentials.SetKeyId(&u)
			}

			if !keyCredentialsModel.StartDateTime.Equal(keyCredentialsState.StartDateTime) {
				planStartDateTime := keyCredentialsModel.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planStartDateTime)
				keyCredentials.SetStartDateTime(&t)
			}

			if !keyCredentialsModel.Type.Equal(keyCredentialsState.Type) {
				planType := keyCredentialsModel.Type.ValueString()
				keyCredentials.SetTypeEscaped(&planType)
			}

			if !keyCredentialsModel.Usage.Equal(keyCredentialsState.Usage) {
				planUsage := keyCredentialsModel.Usage.ValueString()
				keyCredentials.SetUsage(&planUsage)
			}
		}
		requestBody.SetKeyCredentials(planKeyCredentials)
	}

	if !plan.Logo.Equal(state.Logo) {
		planLogo := plan.Logo.ValueString()
		requestBody.SetLogo([]byte(planLogo))
	}

	if !plan.NativeAuthenticationApisEnabled.Equal(state.NativeAuthenticationApisEnabled) {
		planNativeAuthenticationApisEnabled := plan.NativeAuthenticationApisEnabled.ValueString()
		parsedNativeAuthenticationApisEnabled, _ := models.ParseNativeAuthenticationApisEnabled(planNativeAuthenticationApisEnabled)
		assertedNativeAuthenticationApisEnabled := parsedNativeAuthenticationApisEnabled.(models.NativeAuthenticationApisEnabled)
		requestBody.SetNativeAuthenticationApisEnabled(&assertedNativeAuthenticationApisEnabled)
	}

	if !plan.Notes.Equal(state.Notes) {
		planNotes := plan.Notes.ValueString()
		requestBody.SetNotes(&planNotes)
	}

	if !plan.Oauth2RequirePostResponse.Equal(state.Oauth2RequirePostResponse) {
		planOauth2RequirePostResponse := plan.Oauth2RequirePostResponse.ValueBool()
		requestBody.SetOauth2RequirePostResponse(&planOauth2RequirePostResponse)
	}

	if !plan.OptionalClaims.Equal(state.OptionalClaims) {
		optionalClaims := models.NewOptionalClaims()
		optionalClaimsModel := applicationOptionalClaimsModel{}
		plan.OptionalClaims.As(ctx, &optionalClaimsModel, basetypes.ObjectAsOptions{})
		optionalClaimsState := applicationOptionalClaimsModel{}
		state.OptionalClaims.As(ctx, &optionalClaimsState, basetypes.ObjectAsOptions{})

		if !optionalClaimsModel.AccessToken.Equal(optionalClaimsState.AccessToken) {
			var planAccessToken []models.OptionalClaimable
			for k, i := range optionalClaimsModel.AccessToken.Elements() {
				accessToken := models.NewOptionalClaim()
				accessTokenModel := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &accessTokenModel)
				accessTokenState := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, optionalClaimsState.AccessToken.Elements()[k].Type(ctx), &accessTokenModel)

				if !accessTokenModel.AdditionalProperties.Equal(accessTokenState.AdditionalProperties) {
					var additionalProperties []string
					for _, i := range accessTokenModel.AdditionalProperties.Elements() {
						additionalProperties = append(additionalProperties, i.String())
					}
					accessToken.SetAdditionalProperties(additionalProperties)
				}

				if !accessTokenModel.Essential.Equal(accessTokenState.Essential) {
					planEssential := accessTokenModel.Essential.ValueBool()
					accessToken.SetEssential(&planEssential)
				}

				if !accessTokenModel.Name.Equal(accessTokenState.Name) {
					planName := accessTokenModel.Name.ValueString()
					accessToken.SetName(&planName)
				}

				if !accessTokenModel.Source.Equal(accessTokenState.Source) {
					planSource := accessTokenModel.Source.ValueString()
					accessToken.SetSource(&planSource)
				}
			}
			optionalClaims.SetAccessToken(planAccessToken)
		}

		if !optionalClaimsModel.IdToken.Equal(optionalClaimsState.IdToken) {
			var planIdToken []models.OptionalClaimable
			for k, i := range optionalClaimsModel.IdToken.Elements() {
				idToken := models.NewOptionalClaim()
				idTokenModel := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &idTokenModel)
				idTokenState := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, optionalClaimsState.IdToken.Elements()[k].Type(ctx), &idTokenModel)

				if !idTokenModel.AdditionalProperties.Equal(idTokenState.AdditionalProperties) {
					var additionalProperties []string
					for _, i := range idTokenModel.AdditionalProperties.Elements() {
						additionalProperties = append(additionalProperties, i.String())
					}
					idToken.SetAdditionalProperties(additionalProperties)
				}

				if !idTokenModel.Essential.Equal(idTokenState.Essential) {
					planEssential := idTokenModel.Essential.ValueBool()
					idToken.SetEssential(&planEssential)
				}

				if !idTokenModel.Name.Equal(idTokenState.Name) {
					planName := idTokenModel.Name.ValueString()
					idToken.SetName(&planName)
				}

				if !idTokenModel.Source.Equal(idTokenState.Source) {
					planSource := idTokenModel.Source.ValueString()
					idToken.SetSource(&planSource)
				}
			}
			optionalClaims.SetIdToken(planIdToken)
		}

		if !optionalClaimsModel.Saml2Token.Equal(optionalClaimsState.Saml2Token) {
			var planSaml2Token []models.OptionalClaimable
			for k, i := range optionalClaimsModel.Saml2Token.Elements() {
				saml2Token := models.NewOptionalClaim()
				saml2TokenModel := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &saml2TokenModel)
				saml2TokenState := applicationOptionalClaimModel{}
				types.ListValueFrom(ctx, optionalClaimsState.Saml2Token.Elements()[k].Type(ctx), &saml2TokenModel)

				if !saml2TokenModel.AdditionalProperties.Equal(saml2TokenState.AdditionalProperties) {
					var additionalProperties []string
					for _, i := range saml2TokenModel.AdditionalProperties.Elements() {
						additionalProperties = append(additionalProperties, i.String())
					}
					saml2Token.SetAdditionalProperties(additionalProperties)
				}

				if !saml2TokenModel.Essential.Equal(saml2TokenState.Essential) {
					planEssential := saml2TokenModel.Essential.ValueBool()
					saml2Token.SetEssential(&planEssential)
				}

				if !saml2TokenModel.Name.Equal(saml2TokenState.Name) {
					planName := saml2TokenModel.Name.ValueString()
					saml2Token.SetName(&planName)
				}

				if !saml2TokenModel.Source.Equal(saml2TokenState.Source) {
					planSource := saml2TokenModel.Source.ValueString()
					saml2Token.SetSource(&planSource)
				}
			}
			optionalClaims.SetSaml2Token(planSaml2Token)
		}
		requestBody.SetOptionalClaims(optionalClaims)
		objectValue, _ := types.ObjectValueFrom(ctx, optionalClaimsModel.AttributeTypes(), optionalClaimsModel)
		plan.OptionalClaims = objectValue
	}

	if !plan.ParentalControlSettings.Equal(state.ParentalControlSettings) {
		parentalControlSettings := models.NewParentalControlSettings()
		parentalControlSettingsModel := applicationParentalControlSettingsModel{}
		plan.ParentalControlSettings.As(ctx, &parentalControlSettingsModel, basetypes.ObjectAsOptions{})
		parentalControlSettingsState := applicationParentalControlSettingsModel{}
		state.ParentalControlSettings.As(ctx, &parentalControlSettingsState, basetypes.ObjectAsOptions{})

		if !parentalControlSettingsModel.CountriesBlockedForMinors.Equal(parentalControlSettingsState.CountriesBlockedForMinors) {
			var countriesBlockedForMinors []string
			for _, i := range parentalControlSettingsModel.CountriesBlockedForMinors.Elements() {
				countriesBlockedForMinors = append(countriesBlockedForMinors, i.String())
			}
			parentalControlSettings.SetCountriesBlockedForMinors(countriesBlockedForMinors)
		}

		if !parentalControlSettingsModel.LegalAgeGroupRule.Equal(parentalControlSettingsState.LegalAgeGroupRule) {
			planLegalAgeGroupRule := parentalControlSettingsModel.LegalAgeGroupRule.ValueString()
			parentalControlSettings.SetLegalAgeGroupRule(&planLegalAgeGroupRule)
		}
		requestBody.SetParentalControlSettings(parentalControlSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, parentalControlSettingsModel.AttributeTypes(), parentalControlSettingsModel)
		plan.ParentalControlSettings = objectValue
	}

	if !plan.PasswordCredentials.Equal(state.PasswordCredentials) {
		var planPasswordCredentials []models.PasswordCredentialable
		for k, i := range plan.PasswordCredentials.Elements() {
			passwordCredentials := models.NewPasswordCredential()
			passwordCredentialsModel := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &passwordCredentialsModel)
			passwordCredentialsState := applicationPasswordCredentialModel{}
			types.ListValueFrom(ctx, state.PasswordCredentials.Elements()[k].Type(ctx), &passwordCredentialsModel)

			if !passwordCredentialsModel.CustomKeyIdentifier.Equal(passwordCredentialsState.CustomKeyIdentifier) {
				planCustomKeyIdentifier := passwordCredentialsModel.CustomKeyIdentifier.ValueString()
				passwordCredentials.SetCustomKeyIdentifier([]byte(planCustomKeyIdentifier))
			}

			if !passwordCredentialsModel.DisplayName.Equal(passwordCredentialsState.DisplayName) {
				planDisplayName := passwordCredentialsModel.DisplayName.ValueString()
				passwordCredentials.SetDisplayName(&planDisplayName)
			}

			if !passwordCredentialsModel.EndDateTime.Equal(passwordCredentialsState.EndDateTime) {
				planEndDateTime := passwordCredentialsModel.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planEndDateTime)
				passwordCredentials.SetEndDateTime(&t)
			}

			if !passwordCredentialsModel.Hint.Equal(passwordCredentialsState.Hint) {
				planHint := passwordCredentialsModel.Hint.ValueString()
				passwordCredentials.SetHint(&planHint)
			}

			if !passwordCredentialsModel.KeyId.Equal(passwordCredentialsState.KeyId) {
				planKeyId := passwordCredentialsModel.KeyId.ValueString()
				u, _ := uuid.Parse(planKeyId)
				passwordCredentials.SetKeyId(&u)
			}

			if !passwordCredentialsModel.SecretText.Equal(passwordCredentialsState.SecretText) {
				planSecretText := passwordCredentialsModel.SecretText.ValueString()
				passwordCredentials.SetSecretText(&planSecretText)
			}

			if !passwordCredentialsModel.StartDateTime.Equal(passwordCredentialsState.StartDateTime) {
				planStartDateTime := passwordCredentialsModel.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planStartDateTime)
				passwordCredentials.SetStartDateTime(&t)
			}
		}
		requestBody.SetPasswordCredentials(planPasswordCredentials)
	}

	if !plan.PublicClient.Equal(state.PublicClient) {
		publicClient := models.NewPublicClientApplication()
		publicClientModel := applicationPublicClientApplicationModel{}
		plan.PublicClient.As(ctx, &publicClientModel, basetypes.ObjectAsOptions{})
		publicClientState := applicationPublicClientApplicationModel{}
		state.PublicClient.As(ctx, &publicClientState, basetypes.ObjectAsOptions{})

		if !publicClientModel.RedirectUris.Equal(publicClientState.RedirectUris) {
			var redirectUris []string
			for _, i := range publicClientModel.RedirectUris.Elements() {
				redirectUris = append(redirectUris, i.String())
			}
			publicClient.SetRedirectUris(redirectUris)
		}
		requestBody.SetPublicClient(publicClient)
		objectValue, _ := types.ObjectValueFrom(ctx, publicClientModel.AttributeTypes(), publicClientModel)
		plan.PublicClient = objectValue
	}

	if !plan.PublisherDomain.Equal(state.PublisherDomain) {
		planPublisherDomain := plan.PublisherDomain.ValueString()
		requestBody.SetPublisherDomain(&planPublisherDomain)
	}

	if !plan.RequestSignatureVerification.Equal(state.RequestSignatureVerification) {
		requestSignatureVerification := models.NewRequestSignatureVerification()
		requestSignatureVerificationModel := applicationRequestSignatureVerificationModel{}
		plan.RequestSignatureVerification.As(ctx, &requestSignatureVerificationModel, basetypes.ObjectAsOptions{})
		requestSignatureVerificationState := applicationRequestSignatureVerificationModel{}
		state.RequestSignatureVerification.As(ctx, &requestSignatureVerificationState, basetypes.ObjectAsOptions{})

		if !requestSignatureVerificationModel.AllowedWeakAlgorithms.Equal(requestSignatureVerificationState.AllowedWeakAlgorithms) {
			planAllowedWeakAlgorithms := requestSignatureVerificationModel.AllowedWeakAlgorithms.ValueString()
			parsedAllowedWeakAlgorithms, _ := models.ParseWeakAlgorithms(planAllowedWeakAlgorithms)
			assertedAllowedWeakAlgorithms := parsedAllowedWeakAlgorithms.(models.WeakAlgorithms)
			requestSignatureVerification.SetAllowedWeakAlgorithms(&assertedAllowedWeakAlgorithms)
		}

		if !requestSignatureVerificationModel.IsSignedRequestRequired.Equal(requestSignatureVerificationState.IsSignedRequestRequired) {
			planIsSignedRequestRequired := requestSignatureVerificationModel.IsSignedRequestRequired.ValueBool()
			requestSignatureVerification.SetIsSignedRequestRequired(&planIsSignedRequestRequired)
		}
		requestBody.SetRequestSignatureVerification(requestSignatureVerification)
		objectValue, _ := types.ObjectValueFrom(ctx, requestSignatureVerificationModel.AttributeTypes(), requestSignatureVerificationModel)
		plan.RequestSignatureVerification = objectValue
	}

	if !plan.RequiredResourceAccess.Equal(state.RequiredResourceAccess) {
		var planRequiredResourceAccess []models.RequiredResourceAccessable
		for k, i := range plan.RequiredResourceAccess.Elements() {
			requiredResourceAccess := models.NewRequiredResourceAccess()
			requiredResourceAccessModel := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &requiredResourceAccessModel)
			requiredResourceAccessState := applicationRequiredResourceAccessModel{}
			types.ListValueFrom(ctx, state.RequiredResourceAccess.Elements()[k].Type(ctx), &requiredResourceAccessModel)

			if !requiredResourceAccessModel.ResourceAccess.Equal(requiredResourceAccessState.ResourceAccess) {
				var planResourceAccess []models.ResourceAccessable
				for k, i := range requiredResourceAccessModel.ResourceAccess.Elements() {
					resourceAccess := models.NewResourceAccess()
					resourceAccessModel := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &resourceAccessModel)
					resourceAccessState := applicationResourceAccessModel{}
					types.ListValueFrom(ctx, requiredResourceAccessState.ResourceAccess.Elements()[k].Type(ctx), &resourceAccessModel)

					if !resourceAccessModel.Id.Equal(resourceAccessState.Id) {
						planId := resourceAccessModel.Id.ValueString()
						u, _ := uuid.Parse(planId)
						resourceAccess.SetId(&u)
					}

					if !resourceAccessModel.Type.Equal(resourceAccessState.Type) {
						planType := resourceAccessModel.Type.ValueString()
						resourceAccess.SetTypeEscaped(&planType)
					}
				}
				requiredResourceAccess.SetResourceAccess(planResourceAccess)
			}

			if !requiredResourceAccessModel.ResourceAppId.Equal(requiredResourceAccessState.ResourceAppId) {
				planResourceAppId := requiredResourceAccessModel.ResourceAppId.ValueString()
				requiredResourceAccess.SetResourceAppId(&planResourceAppId)
			}
		}
		requestBody.SetRequiredResourceAccess(planRequiredResourceAccess)
	}

	if !plan.SamlMetadataUrl.Equal(state.SamlMetadataUrl) {
		planSamlMetadataUrl := plan.SamlMetadataUrl.ValueString()
		requestBody.SetSamlMetadataUrl(&planSamlMetadataUrl)
	}

	if !plan.ServiceManagementReference.Equal(state.ServiceManagementReference) {
		planServiceManagementReference := plan.ServiceManagementReference.ValueString()
		requestBody.SetServiceManagementReference(&planServiceManagementReference)
	}

	if !plan.ServicePrincipalLockConfiguration.Equal(state.ServicePrincipalLockConfiguration) {
		servicePrincipalLockConfiguration := models.NewServicePrincipalLockConfiguration()
		servicePrincipalLockConfigurationModel := applicationServicePrincipalLockConfigurationModel{}
		plan.ServicePrincipalLockConfiguration.As(ctx, &servicePrincipalLockConfigurationModel, basetypes.ObjectAsOptions{})
		servicePrincipalLockConfigurationState := applicationServicePrincipalLockConfigurationModel{}
		state.ServicePrincipalLockConfiguration.As(ctx, &servicePrincipalLockConfigurationState, basetypes.ObjectAsOptions{})

		if !servicePrincipalLockConfigurationModel.AllProperties.Equal(servicePrincipalLockConfigurationState.AllProperties) {
			planAllProperties := servicePrincipalLockConfigurationModel.AllProperties.ValueBool()
			servicePrincipalLockConfiguration.SetAllProperties(&planAllProperties)
		}

		if !servicePrincipalLockConfigurationModel.CredentialsWithUsageSign.Equal(servicePrincipalLockConfigurationState.CredentialsWithUsageSign) {
			planCredentialsWithUsageSign := servicePrincipalLockConfigurationModel.CredentialsWithUsageSign.ValueBool()
			servicePrincipalLockConfiguration.SetCredentialsWithUsageSign(&planCredentialsWithUsageSign)
		}

		if !servicePrincipalLockConfigurationModel.CredentialsWithUsageVerify.Equal(servicePrincipalLockConfigurationState.CredentialsWithUsageVerify) {
			planCredentialsWithUsageVerify := servicePrincipalLockConfigurationModel.CredentialsWithUsageVerify.ValueBool()
			servicePrincipalLockConfiguration.SetCredentialsWithUsageVerify(&planCredentialsWithUsageVerify)
		}

		if !servicePrincipalLockConfigurationModel.IsEnabled.Equal(servicePrincipalLockConfigurationState.IsEnabled) {
			planIsEnabled := servicePrincipalLockConfigurationModel.IsEnabled.ValueBool()
			servicePrincipalLockConfiguration.SetIsEnabled(&planIsEnabled)
		}

		if !servicePrincipalLockConfigurationModel.TokenEncryptionKeyId.Equal(servicePrincipalLockConfigurationState.TokenEncryptionKeyId) {
			planTokenEncryptionKeyId := servicePrincipalLockConfigurationModel.TokenEncryptionKeyId.ValueBool()
			servicePrincipalLockConfiguration.SetTokenEncryptionKeyId(&planTokenEncryptionKeyId)
		}
		requestBody.SetServicePrincipalLockConfiguration(servicePrincipalLockConfiguration)
		objectValue, _ := types.ObjectValueFrom(ctx, servicePrincipalLockConfigurationModel.AttributeTypes(), servicePrincipalLockConfigurationModel)
		plan.ServicePrincipalLockConfiguration = objectValue
	}

	if !plan.SignInAudience.Equal(state.SignInAudience) {
		planSignInAudience := plan.SignInAudience.ValueString()
		requestBody.SetSignInAudience(&planSignInAudience)
	}

	if !plan.Spa.Equal(state.Spa) {
		spa := models.NewSpaApplication()
		spaModel := applicationSpaApplicationModel{}
		plan.Spa.As(ctx, &spaModel, basetypes.ObjectAsOptions{})
		spaState := applicationSpaApplicationModel{}
		state.Spa.As(ctx, &spaState, basetypes.ObjectAsOptions{})

		if !spaModel.RedirectUris.Equal(spaState.RedirectUris) {
			var redirectUris []string
			for _, i := range spaModel.RedirectUris.Elements() {
				redirectUris = append(redirectUris, i.String())
			}
			spa.SetRedirectUris(redirectUris)
		}
		requestBody.SetSpa(spa)
		objectValue, _ := types.ObjectValueFrom(ctx, spaModel.AttributeTypes(), spaModel)
		plan.Spa = objectValue
	}

	if !plan.Tags.Equal(state.Tags) {
		var tags []string
		for _, i := range plan.Tags.Elements() {
			tags = append(tags, i.String())
		}
		requestBody.SetTags(tags)
	}

	if !plan.TokenEncryptionKeyId.Equal(state.TokenEncryptionKeyId) {
		planTokenEncryptionKeyId := plan.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(planTokenEncryptionKeyId)
		requestBody.SetTokenEncryptionKeyId(&u)
	}

	if !plan.UniqueName.Equal(state.UniqueName) {
		planUniqueName := plan.UniqueName.ValueString()
		requestBody.SetUniqueName(&planUniqueName)
	}

	if !plan.VerifiedPublisher.Equal(state.VerifiedPublisher) {
		verifiedPublisher := models.NewVerifiedPublisher()
		verifiedPublisherModel := applicationVerifiedPublisherModel{}
		plan.VerifiedPublisher.As(ctx, &verifiedPublisherModel, basetypes.ObjectAsOptions{})
		verifiedPublisherState := applicationVerifiedPublisherModel{}
		state.VerifiedPublisher.As(ctx, &verifiedPublisherState, basetypes.ObjectAsOptions{})

		if !verifiedPublisherModel.AddedDateTime.Equal(verifiedPublisherState.AddedDateTime) {
			planAddedDateTime := verifiedPublisherModel.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, planAddedDateTime)
			verifiedPublisher.SetAddedDateTime(&t)
		}

		if !verifiedPublisherModel.DisplayName.Equal(verifiedPublisherState.DisplayName) {
			planDisplayName := verifiedPublisherModel.DisplayName.ValueString()
			verifiedPublisher.SetDisplayName(&planDisplayName)
		}

		if !verifiedPublisherModel.VerifiedPublisherId.Equal(verifiedPublisherState.VerifiedPublisherId) {
			planVerifiedPublisherId := verifiedPublisherModel.VerifiedPublisherId.ValueString()
			verifiedPublisher.SetVerifiedPublisherId(&planVerifiedPublisherId)
		}
		requestBody.SetVerifiedPublisher(verifiedPublisher)
		objectValue, _ := types.ObjectValueFrom(ctx, verifiedPublisherModel.AttributeTypes(), verifiedPublisherModel)
		plan.VerifiedPublisher = objectValue
	}

	if !plan.Web.Equal(state.Web) {
		web := models.NewWebApplication()
		webModel := applicationWebApplicationModel{}
		plan.Web.As(ctx, &webModel, basetypes.ObjectAsOptions{})
		webState := applicationWebApplicationModel{}
		state.Web.As(ctx, &webState, basetypes.ObjectAsOptions{})

		if !webModel.HomePageUrl.Equal(webState.HomePageUrl) {
			planHomePageUrl := webModel.HomePageUrl.ValueString()
			web.SetHomePageUrl(&planHomePageUrl)
		}

		if !webModel.ImplicitGrantSettings.Equal(webState.ImplicitGrantSettings) {
			implicitGrantSettings := models.NewImplicitGrantSettings()
			implicitGrantSettingsModel := applicationImplicitGrantSettingsModel{}
			webModel.ImplicitGrantSettings.As(ctx, &implicitGrantSettingsModel, basetypes.ObjectAsOptions{})
			implicitGrantSettingsState := applicationImplicitGrantSettingsModel{}
			webState.ImplicitGrantSettings.As(ctx, &implicitGrantSettingsState, basetypes.ObjectAsOptions{})

			if !implicitGrantSettingsModel.EnableAccessTokenIssuance.Equal(implicitGrantSettingsState.EnableAccessTokenIssuance) {
				planEnableAccessTokenIssuance := implicitGrantSettingsModel.EnableAccessTokenIssuance.ValueBool()
				implicitGrantSettings.SetEnableAccessTokenIssuance(&planEnableAccessTokenIssuance)
			}

			if !implicitGrantSettingsModel.EnableIdTokenIssuance.Equal(implicitGrantSettingsState.EnableIdTokenIssuance) {
				planEnableIdTokenIssuance := implicitGrantSettingsModel.EnableIdTokenIssuance.ValueBool()
				implicitGrantSettings.SetEnableIdTokenIssuance(&planEnableIdTokenIssuance)
			}
			web.SetImplicitGrantSettings(implicitGrantSettings)
			objectValue, _ := types.ObjectValueFrom(ctx, implicitGrantSettingsModel.AttributeTypes(), implicitGrantSettingsModel)
			webModel.ImplicitGrantSettings = objectValue
		}

		if !webModel.LogoutUrl.Equal(webState.LogoutUrl) {
			planLogoutUrl := webModel.LogoutUrl.ValueString()
			web.SetLogoutUrl(&planLogoutUrl)
		}

		if !webModel.RedirectUriSettings.Equal(webState.RedirectUriSettings) {
			var planRedirectUriSettings []models.RedirectUriSettingsable
			for k, i := range webModel.RedirectUriSettings.Elements() {
				redirectUriSettings := models.NewRedirectUriSettings()
				redirectUriSettingsModel := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, i.Type(ctx), &redirectUriSettingsModel)
				redirectUriSettingsState := applicationRedirectUriSettingsModel{}
				types.ListValueFrom(ctx, webState.RedirectUriSettings.Elements()[k].Type(ctx), &redirectUriSettingsModel)

				if !redirectUriSettingsModel.Uri.Equal(redirectUriSettingsState.Uri) {
					planUri := redirectUriSettingsModel.Uri.ValueString()
					redirectUriSettings.SetUri(&planUri)
				}
			}
			web.SetRedirectUriSettings(planRedirectUriSettings)
		}

		if !webModel.RedirectUris.Equal(webState.RedirectUris) {
			var redirectUris []string
			for _, i := range webModel.RedirectUris.Elements() {
				redirectUris = append(redirectUris, i.String())
			}
			web.SetRedirectUris(redirectUris)
		}
		requestBody.SetWeb(web)
		objectValue, _ := types.ObjectValueFrom(ctx, webModel.AttributeTypes(), webModel)
		plan.Web = objectValue
	}

	// Update application
	_, err := r.client.Applications().ByApplicationId(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating application",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *applicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state applicationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete application
	err := r.client.Applications().ByApplicationId(state.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting application",
			err.Error(),
		)
		return
	}

}
