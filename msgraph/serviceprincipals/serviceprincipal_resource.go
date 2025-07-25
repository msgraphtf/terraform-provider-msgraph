package serviceprincipals

import (
	"context"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"

	"terraform-provider-msgraph/planmodifiers/boolplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/listplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/objectplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/stringplanmodifiers"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &servicePrincipalResource{}
	_ resource.ResourceWithConfigure = &servicePrincipalResource{}
)

// NewServicePrincipalResource is a helper function to simplify the provider implementation.
func NewServicePrincipalResource() resource.Resource {
	return &servicePrincipalResource{}
}

// servicePrincipalResource is the resource implementation.
type servicePrincipalResource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (d *servicePrincipalResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_principal"
}

// Configure adds the provider configured client to the resource.
func (d *servicePrincipalResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the resource.
func (d *servicePrincipalResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "",
		Attributes: map[string]schema.Attribute{
			"account_enabled": schema.BoolAttribute{
				Description: "true if the service principal account is enabled; otherwise, false. If set to false, then no users are able to sign in to this app, even if they're assigned to it. Supports $filter (eq, ne, not, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"add_ins": schema.ListNestedAttribute{
				Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams may set the addIns property for its 'FileHandler' functionality. This lets services like Microsoft 365 call the application in the context of a document the user is working on.",
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
			"alternative_names": schema.ListAttribute{
				Description: "Used to retrieve service principals by subscription, identify resource group and full resource IDs for managed identities. Supports $filter (eq, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"app_description": schema.StringAttribute{
				Description: "The description exposed by the associated application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"app_display_name": schema.StringAttribute{
				Description: "The display name exposed by the associated application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"app_id": schema.StringAttribute{
				Description: "The unique identifier for the associated application (its appId property). Alternate key. Supports $filter (eq, ne, not, in, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"app_owner_organization_id": schema.StringAttribute{
				Description: "Contains the tenant ID where the application is registered. This is applicable only to service principals backed by applications. Supports $filter (eq, ne, NOT, ge, le).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"app_role_assignment_required": schema.BoolAttribute{
				Description: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports $filter (eq, ne, NOT).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"app_roles": schema.ListNestedAttribute{
				Description: "The roles exposed by the application that's linked to this service principal. For more information, see the appRoles property definition on the application entity. Not nullable.",
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
				Description: "Unique identifier of the applicationTemplate. Supports $filter (eq, not, ne). Read-only. null if the service principal wasn't created from an application template.",
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
			"description": schema.StringAttribute{
				Description: "Free text field to provide an internal end-user facing description of the service principal. End-user portals such MyApps displays the application description in this field. The maximum allowed size is 1,024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
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
				Description: "The display name for the service principal. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"homepage": schema.StringAttribute{
				Description: "Home page or landing page of the application.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"info": schema.SingleNestedAttribute{
				Description: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).",
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
			"key_credentials": schema.ListNestedAttribute{
				Description: "The collection of key credentials associated with the service principal. Not nullable. Supports $filter (eq, not, ge, le).",
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
			"login_url": schema.StringAttribute{
				Description: "Specifies the URL where the service provider redirects the user to Microsoft Entra ID to authenticate. Microsoft Entra ID uses the URL to launch the application from Microsoft 365 or the Microsoft Entra My Apps. When blank, Microsoft Entra ID performs IdP-initiated sign-on for applications configured with SAML-based single sign-on. The user launches the application from Microsoft 365, the Microsoft Entra My Apps, or the Microsoft Entra SSO URL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"logout_url": schema.StringAttribute{
				Description: "Specifies the URL that the Microsoft's authorization service uses to sign out a user using OpenID Connect front-channel, back-channel, or SAML sign out protocols.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"notes": schema.StringAttribute{
				Description: "Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1,024 characters.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"notification_email_addresses": schema.ListAttribute{
				Description: "Specifies the list of email addresses where Microsoft Entra ID sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Microsoft Entra Gallery applications.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"oauth_2_permission_scopes": schema.ListNestedAttribute{
				Description: "The delegated permissions exposed by the application. For more information, see the oauth2PermissionScopes property on the application entity's api property. Not nullable.",
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
			"preferred_single_sign_on_mode": schema.StringAttribute{
				Description: "Specifies the single sign-on mode configured for this application. Microsoft Entra ID uses the preferred single sign-on mode to launch the application from Microsoft 365 or the My Apps portal. The supported values are password, saml, notSupported, and oidc. Note: This field might be null for older SAML apps and for OIDC applications where it isn't set automatically.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"preferred_token_signing_key_thumbprint": schema.StringAttribute{
				Description: "This property can be used on SAML applications (apps that have preferredSingleSignOnMode set to saml) to control which certificate is used to sign the SAML responses. For applications that aren't SAML, don't write or otherwise rely on this property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"reply_urls": schema.ListAttribute{
				Description: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application. Not nullable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"resource_specific_application_permissions": schema.ListNestedAttribute{
				Description: "The resource-specific application permissions exposed by this application. Currently, resource-specific permissions are only supported for Teams apps accessing to specific chats and teams using Microsoft Graph. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"description": schema.StringAttribute{
							Description: "Describes the level of access that the resource-specific permission represents.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"display_name": schema.StringAttribute{
							Description: "The display name for the resource-specific permission.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"id": schema.StringAttribute{
							Description: "The unique identifier for the resource-specific application permission.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"is_enabled": schema.BoolAttribute{
							Description: "Indicates whether the permission is enabled.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"value": schema.StringAttribute{
							Description: "The value of the permission.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"saml_single_sign_on_settings": schema.SingleNestedAttribute{
				Description: "The collection for settings related to saml single sign-on.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"relay_state": schema.StringAttribute{
						Description: "The relative URI the service provider would redirect to after completion of the single sign-on flow.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"service_principal_names": schema.ListAttribute{
				Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Microsoft Entra ID. For example,Client apps can specify a resource URI that is based on the values of this property to acquire an access token, which is the URI returned in the 'aud' claim.The any operator is required for filter expressions on multi-valued properties. Not nullable.  Supports $filter (eq, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"service_principal_type": schema.StringAttribute{
				Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Microsoft Entra ID internally. The servicePrincipalType property can be set to three different values: Application - A service principal that represents an application or service. The appId property identifies the associated app registration, and matches the appId of an application, possibly from a different tenant. If the associated app registration is missing, tokens aren't issued for the service principal.ManagedIdentity - A service principal that represents a managed identity. Service principals representing managed identities can be granted access and permissions, but can't be updated or modified directly.Legacy - A service principal that represents an app created before app registrations, or through legacy experiences. A legacy service principal can have credentials, service principal names, reply URLs, and other properties that are editable by an authorized user, but doesn't have an associated app registration. The appId value doesn't associate the service principal with an app registration. The service principal can only be used in the tenant where it was created.SocialIdp - For internal use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"sign_in_audience": schema.StringAttribute{
				Description: "Specifies the Microsoft accounts that are supported for the current application. Read-only. Supported values are:AzureADMyOrg: Users with a Microsoft work or school account in my organization's Microsoft Entra tenant (single-tenant).AzureADMultipleOrgs: Users with a Microsoft work or school account in any organization's Microsoft Entra tenant (multitenant).AzureADandPersonalMicrosoftAccount: Users with a personal Microsoft account, or a work or school account in any organization's Microsoft Entra tenant.PersonalMicrosoftAccount: Users with a personal Microsoft account only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"tags": schema.ListAttribute{
				Description: "Custom strings that can be used to categorize and identify the service principal. Not nullable. The value is the union of strings set here and on the associated application entity's tags property.Supports $filter (eq, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"token_encryption_key_id": schema.StringAttribute{
				Description: "Specifies the keyId of a public key from the keyCredentials collection. When configured, Microsoft Entra ID issues tokens for this application encrypted using the key specified by this property. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"verified_publisher": schema.SingleNestedAttribute{
				Description: "Specifies the verified publisher of the application that's linked to this service principal.",
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
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *servicePrincipalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanServicePrincipal servicePrincipalModel
	diags := req.Plan.Get(ctx, &tfPlanServicePrincipal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	requestBodyServicePrincipal := models.NewServicePrincipal()
	if !tfPlanServicePrincipal.AccountEnabled.IsUnknown() {
		tfPlanAccountEnabled := tfPlanServicePrincipal.AccountEnabled.ValueBool()
		requestBodyServicePrincipal.SetAccountEnabled(&tfPlanAccountEnabled)
	} else {
		tfPlanServicePrincipal.AccountEnabled = types.BoolNull()
	}

	if len(tfPlanServicePrincipal.AddIns.Elements()) > 0 {
		var requestBodyAddIns []models.AddInable
		for _, i := range tfPlanServicePrincipal.AddIns.Elements() {
			requestBodyAddIn := models.NewAddIn()
			tfPlanAddIn := servicePrincipalAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAddIn)

			if !tfPlanAddIn.Id.IsUnknown() {
				tfPlanId := tfPlanAddIn.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAddIn.SetId(&u)
			} else {
				tfPlanAddIn.Id = types.StringNull()
			}

			if len(tfPlanAddIn.Properties.Elements()) > 0 {
				var requestBodyProperties []models.KeyValueable
				for _, i := range tfPlanAddIn.Properties.Elements() {
					requestBodyKeyValue := models.NewKeyValue()
					tfPlanKeyValue := servicePrincipalKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyValue)

					if !tfPlanKeyValue.Key.IsUnknown() {
						tfPlanKey := tfPlanKeyValue.Key.ValueString()
						requestBodyKeyValue.SetKey(&tfPlanKey)
					} else {
						tfPlanKeyValue.Key = types.StringNull()
					}

					if !tfPlanKeyValue.Value.IsUnknown() {
						tfPlanValue := tfPlanKeyValue.Value.ValueString()
						requestBodyKeyValue.SetValue(&tfPlanValue)
					} else {
						tfPlanKeyValue.Value = types.StringNull()
					}

				}
				requestBodyAddIn.SetProperties(requestBodyProperties)
			} else {
				tfPlanAddIn.Properties = types.ListNull(tfPlanAddIn.Properties.ElementType(ctx))
			}

			if !tfPlanAddIn.Type.IsUnknown() {
				tfPlanType := tfPlanAddIn.Type.ValueString()
				requestBodyAddIn.SetTypeEscaped(&tfPlanType)
			} else {
				tfPlanAddIn.Type = types.StringNull()
			}

		}
		requestBodyServicePrincipal.SetAddIns(requestBodyAddIns)
	} else {
		tfPlanServicePrincipal.AddIns = types.ListNull(tfPlanServicePrincipal.AddIns.ElementType(ctx))
	}

	if len(tfPlanServicePrincipal.AlternativeNames.Elements()) > 0 {
		var stringArrayAlternativeNames []string
		for _, i := range tfPlanServicePrincipal.AlternativeNames.Elements() {
			stringArrayAlternativeNames = append(stringArrayAlternativeNames, i.String())
		}
		requestBodyServicePrincipal.SetAlternativeNames(stringArrayAlternativeNames)
	} else {
		tfPlanServicePrincipal.AlternativeNames = types.ListNull(types.StringType)
	}

	if !tfPlanServicePrincipal.AppDescription.IsUnknown() {
		tfPlanAppDescription := tfPlanServicePrincipal.AppDescription.ValueString()
		requestBodyServicePrincipal.SetAppDescription(&tfPlanAppDescription)
	} else {
		tfPlanServicePrincipal.AppDescription = types.StringNull()
	}

	if !tfPlanServicePrincipal.AppDisplayName.IsUnknown() {
		tfPlanAppDisplayName := tfPlanServicePrincipal.AppDisplayName.ValueString()
		requestBodyServicePrincipal.SetAppDisplayName(&tfPlanAppDisplayName)
	} else {
		tfPlanServicePrincipal.AppDisplayName = types.StringNull()
	}

	if !tfPlanServicePrincipal.AppId.IsUnknown() {
		tfPlanAppId := tfPlanServicePrincipal.AppId.ValueString()
		requestBodyServicePrincipal.SetAppId(&tfPlanAppId)
	} else {
		tfPlanServicePrincipal.AppId = types.StringNull()
	}

	if !tfPlanServicePrincipal.AppOwnerOrganizationId.IsUnknown() {
		tfPlanAppOwnerOrganizationId := tfPlanServicePrincipal.AppOwnerOrganizationId.ValueString()
		u, _ := uuid.Parse(tfPlanAppOwnerOrganizationId)
		requestBodyServicePrincipal.SetAppOwnerOrganizationId(&u)
	} else {
		tfPlanServicePrincipal.AppOwnerOrganizationId = types.StringNull()
	}

	if !tfPlanServicePrincipal.AppRoleAssignmentRequired.IsUnknown() {
		tfPlanAppRoleAssignmentRequired := tfPlanServicePrincipal.AppRoleAssignmentRequired.ValueBool()
		requestBodyServicePrincipal.SetAppRoleAssignmentRequired(&tfPlanAppRoleAssignmentRequired)
	} else {
		tfPlanServicePrincipal.AppRoleAssignmentRequired = types.BoolNull()
	}

	if len(tfPlanServicePrincipal.AppRoles.Elements()) > 0 {
		var requestBodyAppRoles []models.AppRoleable
		for _, i := range tfPlanServicePrincipal.AppRoles.Elements() {
			requestBodyAppRole := models.NewAppRole()
			tfPlanAppRole := servicePrincipalAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAppRole)

			if len(tfPlanAppRole.AllowedMemberTypes.Elements()) > 0 {
				var stringArrayAllowedMemberTypes []string
				for _, i := range tfPlanAppRole.AllowedMemberTypes.Elements() {
					stringArrayAllowedMemberTypes = append(stringArrayAllowedMemberTypes, i.String())
				}
				requestBodyAppRole.SetAllowedMemberTypes(stringArrayAllowedMemberTypes)
			} else {
				tfPlanAppRole.AllowedMemberTypes = types.ListNull(types.StringType)
			}

			if !tfPlanAppRole.Description.IsUnknown() {
				tfPlanDescription := tfPlanAppRole.Description.ValueString()
				requestBodyAppRole.SetDescription(&tfPlanDescription)
			} else {
				tfPlanAppRole.Description = types.StringNull()
			}

			if !tfPlanAppRole.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanAppRole.DisplayName.ValueString()
				requestBodyAppRole.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanAppRole.DisplayName = types.StringNull()
			}

			if !tfPlanAppRole.Id.IsUnknown() {
				tfPlanId := tfPlanAppRole.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAppRole.SetId(&u)
			} else {
				tfPlanAppRole.Id = types.StringNull()
			}

			if !tfPlanAppRole.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := tfPlanAppRole.IsEnabled.ValueBool()
				requestBodyAppRole.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				tfPlanAppRole.IsEnabled = types.BoolNull()
			}

			if !tfPlanAppRole.Origin.IsUnknown() {
				tfPlanOrigin := tfPlanAppRole.Origin.ValueString()
				requestBodyAppRole.SetOrigin(&tfPlanOrigin)
			} else {
				tfPlanAppRole.Origin = types.StringNull()
			}

			if !tfPlanAppRole.Value.IsUnknown() {
				tfPlanValue := tfPlanAppRole.Value.ValueString()
				requestBodyAppRole.SetValue(&tfPlanValue)
			} else {
				tfPlanAppRole.Value = types.StringNull()
			}

		}
		requestBodyServicePrincipal.SetAppRoles(requestBodyAppRoles)
	} else {
		tfPlanServicePrincipal.AppRoles = types.ListNull(tfPlanServicePrincipal.AppRoles.ElementType(ctx))
	}

	if !tfPlanServicePrincipal.ApplicationTemplateId.IsUnknown() {
		tfPlanApplicationTemplateId := tfPlanServicePrincipal.ApplicationTemplateId.ValueString()
		requestBodyServicePrincipal.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	} else {
		tfPlanServicePrincipal.ApplicationTemplateId = types.StringNull()
	}

	if !tfPlanServicePrincipal.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlanServicePrincipal.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyServicePrincipal.SetDeletedDateTime(&t)
	} else {
		tfPlanServicePrincipal.DeletedDateTime = types.StringNull()
	}

	if !tfPlanServicePrincipal.Description.IsUnknown() {
		tfPlanDescription := tfPlanServicePrincipal.Description.ValueString()
		requestBodyServicePrincipal.SetDescription(&tfPlanDescription)
	} else {
		tfPlanServicePrincipal.Description = types.StringNull()
	}

	if !tfPlanServicePrincipal.DisabledByMicrosoftStatus.IsUnknown() {
		tfPlanDisabledByMicrosoftStatus := tfPlanServicePrincipal.DisabledByMicrosoftStatus.ValueString()
		requestBodyServicePrincipal.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	} else {
		tfPlanServicePrincipal.DisabledByMicrosoftStatus = types.StringNull()
	}

	if !tfPlanServicePrincipal.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanServicePrincipal.DisplayName.ValueString()
		requestBodyServicePrincipal.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanServicePrincipal.DisplayName = types.StringNull()
	}

	if !tfPlanServicePrincipal.Homepage.IsUnknown() {
		tfPlanHomepage := tfPlanServicePrincipal.Homepage.ValueString()
		requestBodyServicePrincipal.SetHomepage(&tfPlanHomepage)
	} else {
		tfPlanServicePrincipal.Homepage = types.StringNull()
	}

	if !tfPlanServicePrincipal.Id.IsUnknown() {
		tfPlanId := tfPlanServicePrincipal.Id.ValueString()
		requestBodyServicePrincipal.SetId(&tfPlanId)
	} else {
		tfPlanServicePrincipal.Id = types.StringNull()
	}

	if !tfPlanServicePrincipal.Info.IsUnknown() {
		requestBodyInformationalUrl := models.NewInformationalUrl()
		tfPlanInformationalUrl := servicePrincipalInformationalUrlModel{}
		tfPlanServicePrincipal.Info.As(ctx, &tfPlanInformationalUrl, basetypes.ObjectAsOptions{})

		if !tfPlanInformationalUrl.LogoUrl.IsUnknown() {
			tfPlanLogoUrl := tfPlanInformationalUrl.LogoUrl.ValueString()
			requestBodyInformationalUrl.SetLogoUrl(&tfPlanLogoUrl)
		} else {
			tfPlanInformationalUrl.LogoUrl = types.StringNull()
		}

		if !tfPlanInformationalUrl.MarketingUrl.IsUnknown() {
			tfPlanMarketingUrl := tfPlanInformationalUrl.MarketingUrl.ValueString()
			requestBodyInformationalUrl.SetMarketingUrl(&tfPlanMarketingUrl)
		} else {
			tfPlanInformationalUrl.MarketingUrl = types.StringNull()
		}

		if !tfPlanInformationalUrl.PrivacyStatementUrl.IsUnknown() {
			tfPlanPrivacyStatementUrl := tfPlanInformationalUrl.PrivacyStatementUrl.ValueString()
			requestBodyInformationalUrl.SetPrivacyStatementUrl(&tfPlanPrivacyStatementUrl)
		} else {
			tfPlanInformationalUrl.PrivacyStatementUrl = types.StringNull()
		}

		if !tfPlanInformationalUrl.SupportUrl.IsUnknown() {
			tfPlanSupportUrl := tfPlanInformationalUrl.SupportUrl.ValueString()
			requestBodyInformationalUrl.SetSupportUrl(&tfPlanSupportUrl)
		} else {
			tfPlanInformationalUrl.SupportUrl = types.StringNull()
		}

		if !tfPlanInformationalUrl.TermsOfServiceUrl.IsUnknown() {
			tfPlanTermsOfServiceUrl := tfPlanInformationalUrl.TermsOfServiceUrl.ValueString()
			requestBodyInformationalUrl.SetTermsOfServiceUrl(&tfPlanTermsOfServiceUrl)
		} else {
			tfPlanInformationalUrl.TermsOfServiceUrl = types.StringNull()
		}

		requestBodyServicePrincipal.SetInfo(requestBodyInformationalUrl)
		tfPlanServicePrincipal.Info, _ = types.ObjectValueFrom(ctx, tfPlanInformationalUrl.AttributeTypes(), tfPlanInformationalUrl)
	} else {
		tfPlanServicePrincipal.Info = types.ObjectNull(tfPlanServicePrincipal.Info.AttributeTypes(ctx))
	}

	if len(tfPlanServicePrincipal.KeyCredentials.Elements()) > 0 {
		var requestBodyKeyCredentials []models.KeyCredentialable
		for _, i := range tfPlanServicePrincipal.KeyCredentials.Elements() {
			requestBodyKeyCredential := models.NewKeyCredential()
			tfPlanKeyCredential := servicePrincipalKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyCredential)

			if !tfPlanKeyCredential.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := tfPlanKeyCredential.CustomKeyIdentifier.ValueString()
				requestBodyKeyCredential.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				tfPlanKeyCredential.CustomKeyIdentifier = types.StringNull()
			}

			if !tfPlanKeyCredential.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanKeyCredential.DisplayName.ValueString()
				requestBodyKeyCredential.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanKeyCredential.DisplayName = types.StringNull()
			}

			if !tfPlanKeyCredential.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := tfPlanKeyCredential.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyKeyCredential.SetEndDateTime(&t)
			} else {
				tfPlanKeyCredential.EndDateTime = types.StringNull()
			}

			if !tfPlanKeyCredential.Key.IsUnknown() {
				tfPlanKey := tfPlanKeyCredential.Key.ValueString()
				requestBodyKeyCredential.SetKey([]byte(tfPlanKey))
			} else {
				tfPlanKeyCredential.Key = types.StringNull()
			}

			if !tfPlanKeyCredential.KeyId.IsUnknown() {
				tfPlanKeyId := tfPlanKeyCredential.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyKeyCredential.SetKeyId(&u)
			} else {
				tfPlanKeyCredential.KeyId = types.StringNull()
			}

			if !tfPlanKeyCredential.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := tfPlanKeyCredential.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyKeyCredential.SetStartDateTime(&t)
			} else {
				tfPlanKeyCredential.StartDateTime = types.StringNull()
			}

			if !tfPlanKeyCredential.Type.IsUnknown() {
				tfPlanType := tfPlanKeyCredential.Type.ValueString()
				requestBodyKeyCredential.SetTypeEscaped(&tfPlanType)
			} else {
				tfPlanKeyCredential.Type = types.StringNull()
			}

			if !tfPlanKeyCredential.Usage.IsUnknown() {
				tfPlanUsage := tfPlanKeyCredential.Usage.ValueString()
				requestBodyKeyCredential.SetUsage(&tfPlanUsage)
			} else {
				tfPlanKeyCredential.Usage = types.StringNull()
			}

		}
		requestBodyServicePrincipal.SetKeyCredentials(requestBodyKeyCredentials)
	} else {
		tfPlanServicePrincipal.KeyCredentials = types.ListNull(tfPlanServicePrincipal.KeyCredentials.ElementType(ctx))
	}

	if !tfPlanServicePrincipal.LoginUrl.IsUnknown() {
		tfPlanLoginUrl := tfPlanServicePrincipal.LoginUrl.ValueString()
		requestBodyServicePrincipal.SetLoginUrl(&tfPlanLoginUrl)
	} else {
		tfPlanServicePrincipal.LoginUrl = types.StringNull()
	}

	if !tfPlanServicePrincipal.LogoutUrl.IsUnknown() {
		tfPlanLogoutUrl := tfPlanServicePrincipal.LogoutUrl.ValueString()
		requestBodyServicePrincipal.SetLogoutUrl(&tfPlanLogoutUrl)
	} else {
		tfPlanServicePrincipal.LogoutUrl = types.StringNull()
	}

	if !tfPlanServicePrincipal.Notes.IsUnknown() {
		tfPlanNotes := tfPlanServicePrincipal.Notes.ValueString()
		requestBodyServicePrincipal.SetNotes(&tfPlanNotes)
	} else {
		tfPlanServicePrincipal.Notes = types.StringNull()
	}

	if len(tfPlanServicePrincipal.NotificationEmailAddresses.Elements()) > 0 {
		var stringArrayNotificationEmailAddresses []string
		for _, i := range tfPlanServicePrincipal.NotificationEmailAddresses.Elements() {
			stringArrayNotificationEmailAddresses = append(stringArrayNotificationEmailAddresses, i.String())
		}
		requestBodyServicePrincipal.SetNotificationEmailAddresses(stringArrayNotificationEmailAddresses)
	} else {
		tfPlanServicePrincipal.NotificationEmailAddresses = types.ListNull(types.StringType)
	}

	if len(tfPlanServicePrincipal.Oauth2PermissionScopes.Elements()) > 0 {
		var requestBodyOauth2PermissionScopes []models.PermissionScopeable
		for _, i := range tfPlanServicePrincipal.Oauth2PermissionScopes.Elements() {
			requestBodyPermissionScope := models.NewPermissionScope()
			tfPlanPermissionScope := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPermissionScope)

			if !tfPlanPermissionScope.AdminConsentDescription.IsUnknown() {
				tfPlanAdminConsentDescription := tfPlanPermissionScope.AdminConsentDescription.ValueString()
				requestBodyPermissionScope.SetAdminConsentDescription(&tfPlanAdminConsentDescription)
			} else {
				tfPlanPermissionScope.AdminConsentDescription = types.StringNull()
			}

			if !tfPlanPermissionScope.AdminConsentDisplayName.IsUnknown() {
				tfPlanAdminConsentDisplayName := tfPlanPermissionScope.AdminConsentDisplayName.ValueString()
				requestBodyPermissionScope.SetAdminConsentDisplayName(&tfPlanAdminConsentDisplayName)
			} else {
				tfPlanPermissionScope.AdminConsentDisplayName = types.StringNull()
			}

			if !tfPlanPermissionScope.Id.IsUnknown() {
				tfPlanId := tfPlanPermissionScope.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyPermissionScope.SetId(&u)
			} else {
				tfPlanPermissionScope.Id = types.StringNull()
			}

			if !tfPlanPermissionScope.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := tfPlanPermissionScope.IsEnabled.ValueBool()
				requestBodyPermissionScope.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				tfPlanPermissionScope.IsEnabled = types.BoolNull()
			}

			if !tfPlanPermissionScope.Origin.IsUnknown() {
				tfPlanOrigin := tfPlanPermissionScope.Origin.ValueString()
				requestBodyPermissionScope.SetOrigin(&tfPlanOrigin)
			} else {
				tfPlanPermissionScope.Origin = types.StringNull()
			}

			if !tfPlanPermissionScope.Type.IsUnknown() {
				tfPlanType := tfPlanPermissionScope.Type.ValueString()
				requestBodyPermissionScope.SetTypeEscaped(&tfPlanType)
			} else {
				tfPlanPermissionScope.Type = types.StringNull()
			}

			if !tfPlanPermissionScope.UserConsentDescription.IsUnknown() {
				tfPlanUserConsentDescription := tfPlanPermissionScope.UserConsentDescription.ValueString()
				requestBodyPermissionScope.SetUserConsentDescription(&tfPlanUserConsentDescription)
			} else {
				tfPlanPermissionScope.UserConsentDescription = types.StringNull()
			}

			if !tfPlanPermissionScope.UserConsentDisplayName.IsUnknown() {
				tfPlanUserConsentDisplayName := tfPlanPermissionScope.UserConsentDisplayName.ValueString()
				requestBodyPermissionScope.SetUserConsentDisplayName(&tfPlanUserConsentDisplayName)
			} else {
				tfPlanPermissionScope.UserConsentDisplayName = types.StringNull()
			}

			if !tfPlanPermissionScope.Value.IsUnknown() {
				tfPlanValue := tfPlanPermissionScope.Value.ValueString()
				requestBodyPermissionScope.SetValue(&tfPlanValue)
			} else {
				tfPlanPermissionScope.Value = types.StringNull()
			}

		}
		requestBodyServicePrincipal.SetOauth2PermissionScopes(requestBodyOauth2PermissionScopes)
	} else {
		tfPlanServicePrincipal.Oauth2PermissionScopes = types.ListNull(tfPlanServicePrincipal.Oauth2PermissionScopes.ElementType(ctx))
	}

	if len(tfPlanServicePrincipal.PasswordCredentials.Elements()) > 0 {
		var requestBodyPasswordCredentials []models.PasswordCredentialable
		for _, i := range tfPlanServicePrincipal.PasswordCredentials.Elements() {
			requestBodyPasswordCredential := models.NewPasswordCredential()
			tfPlanPasswordCredential := servicePrincipalPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPasswordCredential)

			if !tfPlanPasswordCredential.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := tfPlanPasswordCredential.CustomKeyIdentifier.ValueString()
				requestBodyPasswordCredential.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				tfPlanPasswordCredential.CustomKeyIdentifier = types.StringNull()
			}

			if !tfPlanPasswordCredential.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanPasswordCredential.DisplayName.ValueString()
				requestBodyPasswordCredential.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanPasswordCredential.DisplayName = types.StringNull()
			}

			if !tfPlanPasswordCredential.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := tfPlanPasswordCredential.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyPasswordCredential.SetEndDateTime(&t)
			} else {
				tfPlanPasswordCredential.EndDateTime = types.StringNull()
			}

			if !tfPlanPasswordCredential.Hint.IsUnknown() {
				tfPlanHint := tfPlanPasswordCredential.Hint.ValueString()
				requestBodyPasswordCredential.SetHint(&tfPlanHint)
			} else {
				tfPlanPasswordCredential.Hint = types.StringNull()
			}

			if !tfPlanPasswordCredential.KeyId.IsUnknown() {
				tfPlanKeyId := tfPlanPasswordCredential.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyPasswordCredential.SetKeyId(&u)
			} else {
				tfPlanPasswordCredential.KeyId = types.StringNull()
			}

			if !tfPlanPasswordCredential.SecretText.IsUnknown() {
				tfPlanSecretText := tfPlanPasswordCredential.SecretText.ValueString()
				requestBodyPasswordCredential.SetSecretText(&tfPlanSecretText)
			} else {
				tfPlanPasswordCredential.SecretText = types.StringNull()
			}

			if !tfPlanPasswordCredential.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := tfPlanPasswordCredential.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyPasswordCredential.SetStartDateTime(&t)
			} else {
				tfPlanPasswordCredential.StartDateTime = types.StringNull()
			}

		}
		requestBodyServicePrincipal.SetPasswordCredentials(requestBodyPasswordCredentials)
	} else {
		tfPlanServicePrincipal.PasswordCredentials = types.ListNull(tfPlanServicePrincipal.PasswordCredentials.ElementType(ctx))
	}

	if !tfPlanServicePrincipal.PreferredSingleSignOnMode.IsUnknown() {
		tfPlanPreferredSingleSignOnMode := tfPlanServicePrincipal.PreferredSingleSignOnMode.ValueString()
		requestBodyServicePrincipal.SetPreferredSingleSignOnMode(&tfPlanPreferredSingleSignOnMode)
	} else {
		tfPlanServicePrincipal.PreferredSingleSignOnMode = types.StringNull()
	}

	if !tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint.IsUnknown() {
		tfPlanPreferredTokenSigningKeyThumbprint := tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint.ValueString()
		requestBodyServicePrincipal.SetPreferredTokenSigningKeyThumbprint(&tfPlanPreferredTokenSigningKeyThumbprint)
	} else {
		tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint = types.StringNull()
	}

	if len(tfPlanServicePrincipal.ReplyUrls.Elements()) > 0 {
		var stringArrayReplyUrls []string
		for _, i := range tfPlanServicePrincipal.ReplyUrls.Elements() {
			stringArrayReplyUrls = append(stringArrayReplyUrls, i.String())
		}
		requestBodyServicePrincipal.SetReplyUrls(stringArrayReplyUrls)
	} else {
		tfPlanServicePrincipal.ReplyUrls = types.ListNull(types.StringType)
	}

	if len(tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.Elements()) > 0 {
		var requestBodyResourceSpecificApplicationPermissions []models.ResourceSpecificPermissionable
		for _, i := range tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.Elements() {
			requestBodyResourceSpecificPermission := models.NewResourceSpecificPermission()
			tfPlanResourceSpecificPermission := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanResourceSpecificPermission)

			if !tfPlanResourceSpecificPermission.Description.IsUnknown() {
				tfPlanDescription := tfPlanResourceSpecificPermission.Description.ValueString()
				requestBodyResourceSpecificPermission.SetDescription(&tfPlanDescription)
			} else {
				tfPlanResourceSpecificPermission.Description = types.StringNull()
			}

			if !tfPlanResourceSpecificPermission.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanResourceSpecificPermission.DisplayName.ValueString()
				requestBodyResourceSpecificPermission.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanResourceSpecificPermission.DisplayName = types.StringNull()
			}

			if !tfPlanResourceSpecificPermission.Id.IsUnknown() {
				tfPlanId := tfPlanResourceSpecificPermission.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyResourceSpecificPermission.SetId(&u)
			} else {
				tfPlanResourceSpecificPermission.Id = types.StringNull()
			}

			if !tfPlanResourceSpecificPermission.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := tfPlanResourceSpecificPermission.IsEnabled.ValueBool()
				requestBodyResourceSpecificPermission.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				tfPlanResourceSpecificPermission.IsEnabled = types.BoolNull()
			}

			if !tfPlanResourceSpecificPermission.Value.IsUnknown() {
				tfPlanValue := tfPlanResourceSpecificPermission.Value.ValueString()
				requestBodyResourceSpecificPermission.SetValue(&tfPlanValue)
			} else {
				tfPlanResourceSpecificPermission.Value = types.StringNull()
			}

		}
		requestBodyServicePrincipal.SetResourceSpecificApplicationPermissions(requestBodyResourceSpecificApplicationPermissions)
	} else {
		tfPlanServicePrincipal.ResourceSpecificApplicationPermissions = types.ListNull(tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.ElementType(ctx))
	}

	if !tfPlanServicePrincipal.SamlSingleSignOnSettings.IsUnknown() {
		requestBodySamlSingleSignOnSettings := models.NewSamlSingleSignOnSettings()
		tfPlanSamlSingleSignOnSettings := servicePrincipalSamlSingleSignOnSettingsModel{}
		tfPlanServicePrincipal.SamlSingleSignOnSettings.As(ctx, &tfPlanSamlSingleSignOnSettings, basetypes.ObjectAsOptions{})

		if !tfPlanSamlSingleSignOnSettings.RelayState.IsUnknown() {
			tfPlanRelayState := tfPlanSamlSingleSignOnSettings.RelayState.ValueString()
			requestBodySamlSingleSignOnSettings.SetRelayState(&tfPlanRelayState)
		} else {
			tfPlanSamlSingleSignOnSettings.RelayState = types.StringNull()
		}

		requestBodyServicePrincipal.SetSamlSingleSignOnSettings(requestBodySamlSingleSignOnSettings)
		tfPlanServicePrincipal.SamlSingleSignOnSettings, _ = types.ObjectValueFrom(ctx, tfPlanSamlSingleSignOnSettings.AttributeTypes(), tfPlanSamlSingleSignOnSettings)
	} else {
		tfPlanServicePrincipal.SamlSingleSignOnSettings = types.ObjectNull(tfPlanServicePrincipal.SamlSingleSignOnSettings.AttributeTypes(ctx))
	}

	if len(tfPlanServicePrincipal.ServicePrincipalNames.Elements()) > 0 {
		var stringArrayServicePrincipalNames []string
		for _, i := range tfPlanServicePrincipal.ServicePrincipalNames.Elements() {
			stringArrayServicePrincipalNames = append(stringArrayServicePrincipalNames, i.String())
		}
		requestBodyServicePrincipal.SetServicePrincipalNames(stringArrayServicePrincipalNames)
	} else {
		tfPlanServicePrincipal.ServicePrincipalNames = types.ListNull(types.StringType)
	}

	if !tfPlanServicePrincipal.ServicePrincipalType.IsUnknown() {
		tfPlanServicePrincipalType := tfPlanServicePrincipal.ServicePrincipalType.ValueString()
		requestBodyServicePrincipal.SetServicePrincipalType(&tfPlanServicePrincipalType)
	} else {
		tfPlanServicePrincipal.ServicePrincipalType = types.StringNull()
	}

	if !tfPlanServicePrincipal.SignInAudience.IsUnknown() {
		tfPlanSignInAudience := tfPlanServicePrincipal.SignInAudience.ValueString()
		requestBodyServicePrincipal.SetSignInAudience(&tfPlanSignInAudience)
	} else {
		tfPlanServicePrincipal.SignInAudience = types.StringNull()
	}

	if len(tfPlanServicePrincipal.Tags.Elements()) > 0 {
		var stringArrayTags []string
		for _, i := range tfPlanServicePrincipal.Tags.Elements() {
			stringArrayTags = append(stringArrayTags, i.String())
		}
		requestBodyServicePrincipal.SetTags(stringArrayTags)
	} else {
		tfPlanServicePrincipal.Tags = types.ListNull(types.StringType)
	}

	if !tfPlanServicePrincipal.TokenEncryptionKeyId.IsUnknown() {
		tfPlanTokenEncryptionKeyId := tfPlanServicePrincipal.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		requestBodyServicePrincipal.SetTokenEncryptionKeyId(&u)
	} else {
		tfPlanServicePrincipal.TokenEncryptionKeyId = types.StringNull()
	}

	if !tfPlanServicePrincipal.VerifiedPublisher.IsUnknown() {
		requestBodyVerifiedPublisher := models.NewVerifiedPublisher()
		tfPlanVerifiedPublisher := servicePrincipalVerifiedPublisherModel{}
		tfPlanServicePrincipal.VerifiedPublisher.As(ctx, &tfPlanVerifiedPublisher, basetypes.ObjectAsOptions{})

		if !tfPlanVerifiedPublisher.AddedDateTime.IsUnknown() {
			tfPlanAddedDateTime := tfPlanVerifiedPublisher.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanAddedDateTime)
			requestBodyVerifiedPublisher.SetAddedDateTime(&t)
		} else {
			tfPlanVerifiedPublisher.AddedDateTime = types.StringNull()
		}

		if !tfPlanVerifiedPublisher.DisplayName.IsUnknown() {
			tfPlanDisplayName := tfPlanVerifiedPublisher.DisplayName.ValueString()
			requestBodyVerifiedPublisher.SetDisplayName(&tfPlanDisplayName)
		} else {
			tfPlanVerifiedPublisher.DisplayName = types.StringNull()
		}

		if !tfPlanVerifiedPublisher.VerifiedPublisherId.IsUnknown() {
			tfPlanVerifiedPublisherId := tfPlanVerifiedPublisher.VerifiedPublisherId.ValueString()
			requestBodyVerifiedPublisher.SetVerifiedPublisherId(&tfPlanVerifiedPublisherId)
		} else {
			tfPlanVerifiedPublisher.VerifiedPublisherId = types.StringNull()
		}

		requestBodyServicePrincipal.SetVerifiedPublisher(requestBodyVerifiedPublisher)
		tfPlanServicePrincipal.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfPlanVerifiedPublisher.AttributeTypes(), tfPlanVerifiedPublisher)
	} else {
		tfPlanServicePrincipal.VerifiedPublisher = types.ObjectNull(tfPlanServicePrincipal.VerifiedPublisher.AttributeTypes(ctx))
	}

	// Create new ServicePrincipal
	result, err := r.client.ServicePrincipals().Post(context.Background(), requestBodyServicePrincipal, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating ServicePrincipal",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlanServicePrincipal.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlanServicePrincipal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *servicePrincipalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tfStateServicePrincipal servicePrincipalModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfStateServicePrincipal)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := serviceprincipals.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipals.ServicePrincipalItemRequestBuilderGetQueryParameters{
			Select: []string{
				"accountEnabled",
				"addIns",
				"alternativeNames",
				"appDescription",
				"appDisplayName",
				"appId",
				"appOwnerOrganizationId",
				"appRoleAssignmentRequired",
				"appRoles",
				"applicationTemplateId",
				"deletedDateTime",
				"description",
				"disabledByMicrosoftStatus",
				"displayName",
				"homepage",
				"id",
				"info",
				"keyCredentials",
				"loginUrl",
				"logoutUrl",
				"notes",
				"notificationEmailAddresses",
				"oauth2PermissionScopes",
				"passwordCredentials",
				"preferredSingleSignOnMode",
				"preferredTokenSigningKeyThumbprint",
				"replyUrls",
				"resourceSpecificApplicationPermissions",
				"samlSingleSignOnSettings",
				"servicePrincipalNames",
				"servicePrincipalType",
				"signInAudience",
				"tags",
				"tokenEncryptionKeyId",
				"verifiedPublisher",
			},
		},
	}

	var responseServicePrincipal models.ServicePrincipalable
	var err error

	if !tfStateServicePrincipal.Id.IsNull() {
		responseServicePrincipal, err = d.client.ServicePrincipals().ByServicePrincipalId(tfStateServicePrincipal.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting ServicePrincipal",
			err.Error(),
		)
		return
	}

	if responseServicePrincipal.GetAccountEnabled() != nil {
		tfStateServicePrincipal.AccountEnabled = types.BoolValue(*responseServicePrincipal.GetAccountEnabled())
	} else {
		tfStateServicePrincipal.AccountEnabled = types.BoolNull()
	}
	if len(responseServicePrincipal.GetAddIns()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseAddIn := range responseServicePrincipal.GetAddIns() {
			tfStateAddIn := servicePrincipalAddInModel{}

			if responseAddIn.GetId() != nil {
				tfStateAddIn.Id = types.StringValue(responseAddIn.GetId().String())
			} else {
				tfStateAddIn.Id = types.StringNull()
			}
			if len(responseAddIn.GetProperties()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseKeyValue := range responseAddIn.GetProperties() {
					tfStateKeyValue := servicePrincipalKeyValueModel{}

					if responseKeyValue.GetKey() != nil {
						tfStateKeyValue.Key = types.StringValue(*responseKeyValue.GetKey())
					} else {
						tfStateKeyValue.Key = types.StringNull()
					}
					if responseKeyValue.GetValue() != nil {
						tfStateKeyValue.Value = types.StringValue(*responseKeyValue.GetValue())
					} else {
						tfStateKeyValue.Value = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateKeyValue.AttributeTypes(), tfStateKeyValue)
					objectValues = append(objectValues, objectValue)
				}
				tfStateAddIn.Properties, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if responseAddIn.GetTypeEscaped() != nil {
				tfStateAddIn.Type = types.StringValue(*responseAddIn.GetTypeEscaped())
			} else {
				tfStateAddIn.Type = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAddIn.AttributeTypes(), tfStateAddIn)
			objectValues = append(objectValues, objectValue)
		}
		tfStateServicePrincipal.AddIns, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(responseServicePrincipal.GetAlternativeNames()) > 0 {
		var valueArrayAlternativeNames []attr.Value
		for _, responseAlternativeNames := range responseServicePrincipal.GetAlternativeNames() {
			valueArrayAlternativeNames = append(valueArrayAlternativeNames, types.StringValue(responseAlternativeNames))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayAlternativeNames)
		tfStateServicePrincipal.AlternativeNames = listValue
	} else {
		tfStateServicePrincipal.AlternativeNames = types.ListNull(types.StringType)
	}
	if responseServicePrincipal.GetAppDescription() != nil {
		tfStateServicePrincipal.AppDescription = types.StringValue(*responseServicePrincipal.GetAppDescription())
	} else {
		tfStateServicePrincipal.AppDescription = types.StringNull()
	}
	if responseServicePrincipal.GetAppDisplayName() != nil {
		tfStateServicePrincipal.AppDisplayName = types.StringValue(*responseServicePrincipal.GetAppDisplayName())
	} else {
		tfStateServicePrincipal.AppDisplayName = types.StringNull()
	}
	if responseServicePrincipal.GetAppId() != nil {
		tfStateServicePrincipal.AppId = types.StringValue(*responseServicePrincipal.GetAppId())
	} else {
		tfStateServicePrincipal.AppId = types.StringNull()
	}
	if responseServicePrincipal.GetAppOwnerOrganizationId() != nil {
		tfStateServicePrincipal.AppOwnerOrganizationId = types.StringValue(responseServicePrincipal.GetAppOwnerOrganizationId().String())
	} else {
		tfStateServicePrincipal.AppOwnerOrganizationId = types.StringNull()
	}
	if responseServicePrincipal.GetAppRoleAssignmentRequired() != nil {
		tfStateServicePrincipal.AppRoleAssignmentRequired = types.BoolValue(*responseServicePrincipal.GetAppRoleAssignmentRequired())
	} else {
		tfStateServicePrincipal.AppRoleAssignmentRequired = types.BoolNull()
	}
	if len(responseServicePrincipal.GetAppRoles()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseAppRole := range responseServicePrincipal.GetAppRoles() {
			tfStateAppRole := servicePrincipalAppRoleModel{}

			if len(responseAppRole.GetAllowedMemberTypes()) > 0 {
				var valueArrayAllowedMemberTypes []attr.Value
				for _, responseAllowedMemberTypes := range responseAppRole.GetAllowedMemberTypes() {
					valueArrayAllowedMemberTypes = append(valueArrayAllowedMemberTypes, types.StringValue(responseAllowedMemberTypes))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayAllowedMemberTypes)
				tfStateAppRole.AllowedMemberTypes = listValue
			} else {
				tfStateAppRole.AllowedMemberTypes = types.ListNull(types.StringType)
			}
			if responseAppRole.GetDescription() != nil {
				tfStateAppRole.Description = types.StringValue(*responseAppRole.GetDescription())
			} else {
				tfStateAppRole.Description = types.StringNull()
			}
			if responseAppRole.GetDisplayName() != nil {
				tfStateAppRole.DisplayName = types.StringValue(*responseAppRole.GetDisplayName())
			} else {
				tfStateAppRole.DisplayName = types.StringNull()
			}
			if responseAppRole.GetId() != nil {
				tfStateAppRole.Id = types.StringValue(responseAppRole.GetId().String())
			} else {
				tfStateAppRole.Id = types.StringNull()
			}
			if responseAppRole.GetIsEnabled() != nil {
				tfStateAppRole.IsEnabled = types.BoolValue(*responseAppRole.GetIsEnabled())
			} else {
				tfStateAppRole.IsEnabled = types.BoolNull()
			}
			if responseAppRole.GetOrigin() != nil {
				tfStateAppRole.Origin = types.StringValue(*responseAppRole.GetOrigin())
			} else {
				tfStateAppRole.Origin = types.StringNull()
			}
			if responseAppRole.GetValue() != nil {
				tfStateAppRole.Value = types.StringValue(*responseAppRole.GetValue())
			} else {
				tfStateAppRole.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAppRole.AttributeTypes(), tfStateAppRole)
			objectValues = append(objectValues, objectValue)
		}
		tfStateServicePrincipal.AppRoles, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if responseServicePrincipal.GetApplicationTemplateId() != nil {
		tfStateServicePrincipal.ApplicationTemplateId = types.StringValue(*responseServicePrincipal.GetApplicationTemplateId())
	} else {
		tfStateServicePrincipal.ApplicationTemplateId = types.StringNull()
	}
	if responseServicePrincipal.GetDeletedDateTime() != nil {
		tfStateServicePrincipal.DeletedDateTime = types.StringValue(responseServicePrincipal.GetDeletedDateTime().String())
	} else {
		tfStateServicePrincipal.DeletedDateTime = types.StringNull()
	}
	if responseServicePrincipal.GetDescription() != nil {
		tfStateServicePrincipal.Description = types.StringValue(*responseServicePrincipal.GetDescription())
	} else {
		tfStateServicePrincipal.Description = types.StringNull()
	}
	if responseServicePrincipal.GetDisabledByMicrosoftStatus() != nil {
		tfStateServicePrincipal.DisabledByMicrosoftStatus = types.StringValue(*responseServicePrincipal.GetDisabledByMicrosoftStatus())
	} else {
		tfStateServicePrincipal.DisabledByMicrosoftStatus = types.StringNull()
	}
	if responseServicePrincipal.GetDisplayName() != nil {
		tfStateServicePrincipal.DisplayName = types.StringValue(*responseServicePrincipal.GetDisplayName())
	} else {
		tfStateServicePrincipal.DisplayName = types.StringNull()
	}
	if responseServicePrincipal.GetHomepage() != nil {
		tfStateServicePrincipal.Homepage = types.StringValue(*responseServicePrincipal.GetHomepage())
	} else {
		tfStateServicePrincipal.Homepage = types.StringNull()
	}
	if responseServicePrincipal.GetId() != nil {
		tfStateServicePrincipal.Id = types.StringValue(*responseServicePrincipal.GetId())
	} else {
		tfStateServicePrincipal.Id = types.StringNull()
	}
	if responseServicePrincipal.GetInfo() != nil {
		tfStateInformationalUrl := servicePrincipalInformationalUrlModel{}
		responseInformationalUrl := responseServicePrincipal.GetInfo()

		if responseInformationalUrl.GetLogoUrl() != nil {
			tfStateInformationalUrl.LogoUrl = types.StringValue(*responseInformationalUrl.GetLogoUrl())
		} else {
			tfStateInformationalUrl.LogoUrl = types.StringNull()
		}
		if responseInformationalUrl.GetMarketingUrl() != nil {
			tfStateInformationalUrl.MarketingUrl = types.StringValue(*responseInformationalUrl.GetMarketingUrl())
		} else {
			tfStateInformationalUrl.MarketingUrl = types.StringNull()
		}
		if responseInformationalUrl.GetPrivacyStatementUrl() != nil {
			tfStateInformationalUrl.PrivacyStatementUrl = types.StringValue(*responseInformationalUrl.GetPrivacyStatementUrl())
		} else {
			tfStateInformationalUrl.PrivacyStatementUrl = types.StringNull()
		}
		if responseInformationalUrl.GetSupportUrl() != nil {
			tfStateInformationalUrl.SupportUrl = types.StringValue(*responseInformationalUrl.GetSupportUrl())
		} else {
			tfStateInformationalUrl.SupportUrl = types.StringNull()
		}
		if responseInformationalUrl.GetTermsOfServiceUrl() != nil {
			tfStateInformationalUrl.TermsOfServiceUrl = types.StringValue(*responseInformationalUrl.GetTermsOfServiceUrl())
		} else {
			tfStateInformationalUrl.TermsOfServiceUrl = types.StringNull()
		}

		tfStateServicePrincipal.Info, _ = types.ObjectValueFrom(ctx, tfStateInformationalUrl.AttributeTypes(), tfStateInformationalUrl)
	}
	if len(responseServicePrincipal.GetKeyCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseKeyCredential := range responseServicePrincipal.GetKeyCredentials() {
			tfStateKeyCredential := servicePrincipalKeyCredentialModel{}

			if responseKeyCredential.GetCustomKeyIdentifier() != nil {
				tfStateKeyCredential.CustomKeyIdentifier = types.StringValue(string(responseKeyCredential.GetCustomKeyIdentifier()[:]))
			} else {
				tfStateKeyCredential.CustomKeyIdentifier = types.StringNull()
			}
			if responseKeyCredential.GetDisplayName() != nil {
				tfStateKeyCredential.DisplayName = types.StringValue(*responseKeyCredential.GetDisplayName())
			} else {
				tfStateKeyCredential.DisplayName = types.StringNull()
			}
			if responseKeyCredential.GetEndDateTime() != nil {
				tfStateKeyCredential.EndDateTime = types.StringValue(responseKeyCredential.GetEndDateTime().String())
			} else {
				tfStateKeyCredential.EndDateTime = types.StringNull()
			}
			if responseKeyCredential.GetKey() != nil {
				tfStateKeyCredential.Key = types.StringValue(string(responseKeyCredential.GetKey()[:]))
			} else {
				tfStateKeyCredential.Key = types.StringNull()
			}
			if responseKeyCredential.GetKeyId() != nil {
				tfStateKeyCredential.KeyId = types.StringValue(responseKeyCredential.GetKeyId().String())
			} else {
				tfStateKeyCredential.KeyId = types.StringNull()
			}
			if responseKeyCredential.GetStartDateTime() != nil {
				tfStateKeyCredential.StartDateTime = types.StringValue(responseKeyCredential.GetStartDateTime().String())
			} else {
				tfStateKeyCredential.StartDateTime = types.StringNull()
			}
			if responseKeyCredential.GetTypeEscaped() != nil {
				tfStateKeyCredential.Type = types.StringValue(*responseKeyCredential.GetTypeEscaped())
			} else {
				tfStateKeyCredential.Type = types.StringNull()
			}
			if responseKeyCredential.GetUsage() != nil {
				tfStateKeyCredential.Usage = types.StringValue(*responseKeyCredential.GetUsage())
			} else {
				tfStateKeyCredential.Usage = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateKeyCredential.AttributeTypes(), tfStateKeyCredential)
			objectValues = append(objectValues, objectValue)
		}
		tfStateServicePrincipal.KeyCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if responseServicePrincipal.GetLoginUrl() != nil {
		tfStateServicePrincipal.LoginUrl = types.StringValue(*responseServicePrincipal.GetLoginUrl())
	} else {
		tfStateServicePrincipal.LoginUrl = types.StringNull()
	}
	if responseServicePrincipal.GetLogoutUrl() != nil {
		tfStateServicePrincipal.LogoutUrl = types.StringValue(*responseServicePrincipal.GetLogoutUrl())
	} else {
		tfStateServicePrincipal.LogoutUrl = types.StringNull()
	}
	if responseServicePrincipal.GetNotes() != nil {
		tfStateServicePrincipal.Notes = types.StringValue(*responseServicePrincipal.GetNotes())
	} else {
		tfStateServicePrincipal.Notes = types.StringNull()
	}
	if len(responseServicePrincipal.GetNotificationEmailAddresses()) > 0 {
		var valueArrayNotificationEmailAddresses []attr.Value
		for _, responseNotificationEmailAddresses := range responseServicePrincipal.GetNotificationEmailAddresses() {
			valueArrayNotificationEmailAddresses = append(valueArrayNotificationEmailAddresses, types.StringValue(responseNotificationEmailAddresses))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayNotificationEmailAddresses)
		tfStateServicePrincipal.NotificationEmailAddresses = listValue
	} else {
		tfStateServicePrincipal.NotificationEmailAddresses = types.ListNull(types.StringType)
	}
	if len(responseServicePrincipal.GetOauth2PermissionScopes()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responsePermissionScope := range responseServicePrincipal.GetOauth2PermissionScopes() {
			tfStatePermissionScope := servicePrincipalPermissionScopeModel{}

			if responsePermissionScope.GetAdminConsentDescription() != nil {
				tfStatePermissionScope.AdminConsentDescription = types.StringValue(*responsePermissionScope.GetAdminConsentDescription())
			} else {
				tfStatePermissionScope.AdminConsentDescription = types.StringNull()
			}
			if responsePermissionScope.GetAdminConsentDisplayName() != nil {
				tfStatePermissionScope.AdminConsentDisplayName = types.StringValue(*responsePermissionScope.GetAdminConsentDisplayName())
			} else {
				tfStatePermissionScope.AdminConsentDisplayName = types.StringNull()
			}
			if responsePermissionScope.GetId() != nil {
				tfStatePermissionScope.Id = types.StringValue(responsePermissionScope.GetId().String())
			} else {
				tfStatePermissionScope.Id = types.StringNull()
			}
			if responsePermissionScope.GetIsEnabled() != nil {
				tfStatePermissionScope.IsEnabled = types.BoolValue(*responsePermissionScope.GetIsEnabled())
			} else {
				tfStatePermissionScope.IsEnabled = types.BoolNull()
			}
			if responsePermissionScope.GetOrigin() != nil {
				tfStatePermissionScope.Origin = types.StringValue(*responsePermissionScope.GetOrigin())
			} else {
				tfStatePermissionScope.Origin = types.StringNull()
			}
			if responsePermissionScope.GetTypeEscaped() != nil {
				tfStatePermissionScope.Type = types.StringValue(*responsePermissionScope.GetTypeEscaped())
			} else {
				tfStatePermissionScope.Type = types.StringNull()
			}
			if responsePermissionScope.GetUserConsentDescription() != nil {
				tfStatePermissionScope.UserConsentDescription = types.StringValue(*responsePermissionScope.GetUserConsentDescription())
			} else {
				tfStatePermissionScope.UserConsentDescription = types.StringNull()
			}
			if responsePermissionScope.GetUserConsentDisplayName() != nil {
				tfStatePermissionScope.UserConsentDisplayName = types.StringValue(*responsePermissionScope.GetUserConsentDisplayName())
			} else {
				tfStatePermissionScope.UserConsentDisplayName = types.StringNull()
			}
			if responsePermissionScope.GetValue() != nil {
				tfStatePermissionScope.Value = types.StringValue(*responsePermissionScope.GetValue())
			} else {
				tfStatePermissionScope.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStatePermissionScope.AttributeTypes(), tfStatePermissionScope)
			objectValues = append(objectValues, objectValue)
		}
		tfStateServicePrincipal.Oauth2PermissionScopes, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(responseServicePrincipal.GetPasswordCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responsePasswordCredential := range responseServicePrincipal.GetPasswordCredentials() {
			tfStatePasswordCredential := servicePrincipalPasswordCredentialModel{}

			if responsePasswordCredential.GetCustomKeyIdentifier() != nil {
				tfStatePasswordCredential.CustomKeyIdentifier = types.StringValue(string(responsePasswordCredential.GetCustomKeyIdentifier()[:]))
			} else {
				tfStatePasswordCredential.CustomKeyIdentifier = types.StringNull()
			}
			if responsePasswordCredential.GetDisplayName() != nil {
				tfStatePasswordCredential.DisplayName = types.StringValue(*responsePasswordCredential.GetDisplayName())
			} else {
				tfStatePasswordCredential.DisplayName = types.StringNull()
			}
			if responsePasswordCredential.GetEndDateTime() != nil {
				tfStatePasswordCredential.EndDateTime = types.StringValue(responsePasswordCredential.GetEndDateTime().String())
			} else {
				tfStatePasswordCredential.EndDateTime = types.StringNull()
			}
			if responsePasswordCredential.GetHint() != nil {
				tfStatePasswordCredential.Hint = types.StringValue(*responsePasswordCredential.GetHint())
			} else {
				tfStatePasswordCredential.Hint = types.StringNull()
			}
			if responsePasswordCredential.GetKeyId() != nil {
				tfStatePasswordCredential.KeyId = types.StringValue(responsePasswordCredential.GetKeyId().String())
			} else {
				tfStatePasswordCredential.KeyId = types.StringNull()
			}
			if responsePasswordCredential.GetSecretText() != nil {
				tfStatePasswordCredential.SecretText = types.StringValue(*responsePasswordCredential.GetSecretText())
			} else {
				tfStatePasswordCredential.SecretText = types.StringNull()
			}
			if responsePasswordCredential.GetStartDateTime() != nil {
				tfStatePasswordCredential.StartDateTime = types.StringValue(responsePasswordCredential.GetStartDateTime().String())
			} else {
				tfStatePasswordCredential.StartDateTime = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStatePasswordCredential.AttributeTypes(), tfStatePasswordCredential)
			objectValues = append(objectValues, objectValue)
		}
		tfStateServicePrincipal.PasswordCredentials, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if responseServicePrincipal.GetPreferredSingleSignOnMode() != nil {
		tfStateServicePrincipal.PreferredSingleSignOnMode = types.StringValue(*responseServicePrincipal.GetPreferredSingleSignOnMode())
	} else {
		tfStateServicePrincipal.PreferredSingleSignOnMode = types.StringNull()
	}
	if responseServicePrincipal.GetPreferredTokenSigningKeyThumbprint() != nil {
		tfStateServicePrincipal.PreferredTokenSigningKeyThumbprint = types.StringValue(*responseServicePrincipal.GetPreferredTokenSigningKeyThumbprint())
	} else {
		tfStateServicePrincipal.PreferredTokenSigningKeyThumbprint = types.StringNull()
	}
	if len(responseServicePrincipal.GetReplyUrls()) > 0 {
		var valueArrayReplyUrls []attr.Value
		for _, responseReplyUrls := range responseServicePrincipal.GetReplyUrls() {
			valueArrayReplyUrls = append(valueArrayReplyUrls, types.StringValue(responseReplyUrls))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayReplyUrls)
		tfStateServicePrincipal.ReplyUrls = listValue
	} else {
		tfStateServicePrincipal.ReplyUrls = types.ListNull(types.StringType)
	}
	if len(responseServicePrincipal.GetResourceSpecificApplicationPermissions()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseResourceSpecificPermission := range responseServicePrincipal.GetResourceSpecificApplicationPermissions() {
			tfStateResourceSpecificPermission := servicePrincipalResourceSpecificPermissionModel{}

			if responseResourceSpecificPermission.GetDescription() != nil {
				tfStateResourceSpecificPermission.Description = types.StringValue(*responseResourceSpecificPermission.GetDescription())
			} else {
				tfStateResourceSpecificPermission.Description = types.StringNull()
			}
			if responseResourceSpecificPermission.GetDisplayName() != nil {
				tfStateResourceSpecificPermission.DisplayName = types.StringValue(*responseResourceSpecificPermission.GetDisplayName())
			} else {
				tfStateResourceSpecificPermission.DisplayName = types.StringNull()
			}
			if responseResourceSpecificPermission.GetId() != nil {
				tfStateResourceSpecificPermission.Id = types.StringValue(responseResourceSpecificPermission.GetId().String())
			} else {
				tfStateResourceSpecificPermission.Id = types.StringNull()
			}
			if responseResourceSpecificPermission.GetIsEnabled() != nil {
				tfStateResourceSpecificPermission.IsEnabled = types.BoolValue(*responseResourceSpecificPermission.GetIsEnabled())
			} else {
				tfStateResourceSpecificPermission.IsEnabled = types.BoolNull()
			}
			if responseResourceSpecificPermission.GetValue() != nil {
				tfStateResourceSpecificPermission.Value = types.StringValue(*responseResourceSpecificPermission.GetValue())
			} else {
				tfStateResourceSpecificPermission.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateResourceSpecificPermission.AttributeTypes(), tfStateResourceSpecificPermission)
			objectValues = append(objectValues, objectValue)
		}
		tfStateServicePrincipal.ResourceSpecificApplicationPermissions, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if responseServicePrincipal.GetSamlSingleSignOnSettings() != nil {
		tfStateSamlSingleSignOnSettings := servicePrincipalSamlSingleSignOnSettingsModel{}
		responseSamlSingleSignOnSettings := responseServicePrincipal.GetSamlSingleSignOnSettings()

		if responseSamlSingleSignOnSettings.GetRelayState() != nil {
			tfStateSamlSingleSignOnSettings.RelayState = types.StringValue(*responseSamlSingleSignOnSettings.GetRelayState())
		} else {
			tfStateSamlSingleSignOnSettings.RelayState = types.StringNull()
		}

		tfStateServicePrincipal.SamlSingleSignOnSettings, _ = types.ObjectValueFrom(ctx, tfStateSamlSingleSignOnSettings.AttributeTypes(), tfStateSamlSingleSignOnSettings)
	}
	if len(responseServicePrincipal.GetServicePrincipalNames()) > 0 {
		var valueArrayServicePrincipalNames []attr.Value
		for _, responseServicePrincipalNames := range responseServicePrincipal.GetServicePrincipalNames() {
			valueArrayServicePrincipalNames = append(valueArrayServicePrincipalNames, types.StringValue(responseServicePrincipalNames))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayServicePrincipalNames)
		tfStateServicePrincipal.ServicePrincipalNames = listValue
	} else {
		tfStateServicePrincipal.ServicePrincipalNames = types.ListNull(types.StringType)
	}
	if responseServicePrincipal.GetServicePrincipalType() != nil {
		tfStateServicePrincipal.ServicePrincipalType = types.StringValue(*responseServicePrincipal.GetServicePrincipalType())
	} else {
		tfStateServicePrincipal.ServicePrincipalType = types.StringNull()
	}
	if responseServicePrincipal.GetSignInAudience() != nil {
		tfStateServicePrincipal.SignInAudience = types.StringValue(*responseServicePrincipal.GetSignInAudience())
	} else {
		tfStateServicePrincipal.SignInAudience = types.StringNull()
	}
	if len(responseServicePrincipal.GetTags()) > 0 {
		var valueArrayTags []attr.Value
		for _, responseTags := range responseServicePrincipal.GetTags() {
			valueArrayTags = append(valueArrayTags, types.StringValue(responseTags))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayTags)
		tfStateServicePrincipal.Tags = listValue
	} else {
		tfStateServicePrincipal.Tags = types.ListNull(types.StringType)
	}
	if responseServicePrincipal.GetTokenEncryptionKeyId() != nil {
		tfStateServicePrincipal.TokenEncryptionKeyId = types.StringValue(responseServicePrincipal.GetTokenEncryptionKeyId().String())
	} else {
		tfStateServicePrincipal.TokenEncryptionKeyId = types.StringNull()
	}
	if responseServicePrincipal.GetVerifiedPublisher() != nil {
		tfStateVerifiedPublisher := servicePrincipalVerifiedPublisherModel{}
		responseVerifiedPublisher := responseServicePrincipal.GetVerifiedPublisher()

		if responseVerifiedPublisher.GetAddedDateTime() != nil {
			tfStateVerifiedPublisher.AddedDateTime = types.StringValue(responseVerifiedPublisher.GetAddedDateTime().String())
		} else {
			tfStateVerifiedPublisher.AddedDateTime = types.StringNull()
		}
		if responseVerifiedPublisher.GetDisplayName() != nil {
			tfStateVerifiedPublisher.DisplayName = types.StringValue(*responseVerifiedPublisher.GetDisplayName())
		} else {
			tfStateVerifiedPublisher.DisplayName = types.StringNull()
		}
		if responseVerifiedPublisher.GetVerifiedPublisherId() != nil {
			tfStateVerifiedPublisher.VerifiedPublisherId = types.StringValue(*responseVerifiedPublisher.GetVerifiedPublisherId())
		} else {
			tfStateVerifiedPublisher.VerifiedPublisherId = types.StringNull()
		}

		tfStateServicePrincipal.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfStateVerifiedPublisher.AttributeTypes(), tfStateVerifiedPublisher)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateServicePrincipal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *servicePrincipalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanServicePrincipal servicePrincipalModel
	diags := req.Plan.Get(ctx, &tfPlanServicePrincipal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfStateServicePrincipal servicePrincipalModel
	diags = req.State.Get(ctx, &tfStateServicePrincipal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBodyServicePrincipal := models.NewServicePrincipal()

	if !tfPlanServicePrincipal.AccountEnabled.Equal(tfStateServicePrincipal.AccountEnabled) {
		tfPlanAccountEnabled := tfPlanServicePrincipal.AccountEnabled.ValueBool()
		requestBodyServicePrincipal.SetAccountEnabled(&tfPlanAccountEnabled)
	}

	if !tfPlanServicePrincipal.AddIns.Equal(tfStateServicePrincipal.AddIns) {
		var tfPlanAddIns []models.AddInable
		for k, i := range tfPlanServicePrincipal.AddIns.Elements() {
			requestBodyAddIn := models.NewAddIn()
			tfPlanAddIn := servicePrincipalAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAddIn)
			tfStateAddIn := servicePrincipalAddInModel{}
			types.ListValueFrom(ctx, tfStateServicePrincipal.AddIns.Elements()[k].Type(ctx), &tfPlanAddIn)

			if !tfPlanAddIn.Id.Equal(tfStateAddIn.Id) {
				tfPlanId := tfPlanAddIn.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAddIn.SetId(&u)
			}

			if !tfPlanAddIn.Properties.Equal(tfStateAddIn.Properties) {
				var tfPlanProperties []models.KeyValueable
				for k, i := range tfPlanAddIn.Properties.Elements() {
					requestBodyKeyValue := models.NewKeyValue()
					tfPlanKeyValue := servicePrincipalKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyValue)
					tfStateKeyValue := servicePrincipalKeyValueModel{}
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
		requestBodyServicePrincipal.SetAddIns(tfPlanAddIns)
	}

	if !tfPlanServicePrincipal.AlternativeNames.Equal(tfStateServicePrincipal.AlternativeNames) {
		var stringArrayAlternativeNames []string
		for _, i := range tfPlanServicePrincipal.AlternativeNames.Elements() {
			stringArrayAlternativeNames = append(stringArrayAlternativeNames, i.String())
		}
		requestBodyServicePrincipal.SetAlternativeNames(stringArrayAlternativeNames)
	}

	if !tfPlanServicePrincipal.AppDescription.Equal(tfStateServicePrincipal.AppDescription) {
		tfPlanAppDescription := tfPlanServicePrincipal.AppDescription.ValueString()
		requestBodyServicePrincipal.SetAppDescription(&tfPlanAppDescription)
	}

	if !tfPlanServicePrincipal.AppDisplayName.Equal(tfStateServicePrincipal.AppDisplayName) {
		tfPlanAppDisplayName := tfPlanServicePrincipal.AppDisplayName.ValueString()
		requestBodyServicePrincipal.SetAppDisplayName(&tfPlanAppDisplayName)
	}

	if !tfPlanServicePrincipal.AppId.Equal(tfStateServicePrincipal.AppId) {
		tfPlanAppId := tfPlanServicePrincipal.AppId.ValueString()
		requestBodyServicePrincipal.SetAppId(&tfPlanAppId)
	}

	if !tfPlanServicePrincipal.AppOwnerOrganizationId.Equal(tfStateServicePrincipal.AppOwnerOrganizationId) {
		tfPlanAppOwnerOrganizationId := tfPlanServicePrincipal.AppOwnerOrganizationId.ValueString()
		u, _ := uuid.Parse(tfPlanAppOwnerOrganizationId)
		requestBodyServicePrincipal.SetAppOwnerOrganizationId(&u)
	}

	if !tfPlanServicePrincipal.AppRoleAssignmentRequired.Equal(tfStateServicePrincipal.AppRoleAssignmentRequired) {
		tfPlanAppRoleAssignmentRequired := tfPlanServicePrincipal.AppRoleAssignmentRequired.ValueBool()
		requestBodyServicePrincipal.SetAppRoleAssignmentRequired(&tfPlanAppRoleAssignmentRequired)
	}

	if !tfPlanServicePrincipal.AppRoles.Equal(tfStateServicePrincipal.AppRoles) {
		var tfPlanAppRoles []models.AppRoleable
		for k, i := range tfPlanServicePrincipal.AppRoles.Elements() {
			requestBodyAppRole := models.NewAppRole()
			tfPlanAppRole := servicePrincipalAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAppRole)
			tfStateAppRole := servicePrincipalAppRoleModel{}
			types.ListValueFrom(ctx, tfStateServicePrincipal.AppRoles.Elements()[k].Type(ctx), &tfPlanAppRole)

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
		requestBodyServicePrincipal.SetAppRoles(tfPlanAppRoles)
	}

	if !tfPlanServicePrincipal.ApplicationTemplateId.Equal(tfStateServicePrincipal.ApplicationTemplateId) {
		tfPlanApplicationTemplateId := tfPlanServicePrincipal.ApplicationTemplateId.ValueString()
		requestBodyServicePrincipal.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	}

	if !tfPlanServicePrincipal.DeletedDateTime.Equal(tfStateServicePrincipal.DeletedDateTime) {
		tfPlanDeletedDateTime := tfPlanServicePrincipal.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyServicePrincipal.SetDeletedDateTime(&t)
	}

	if !tfPlanServicePrincipal.Description.Equal(tfStateServicePrincipal.Description) {
		tfPlanDescription := tfPlanServicePrincipal.Description.ValueString()
		requestBodyServicePrincipal.SetDescription(&tfPlanDescription)
	}

	if !tfPlanServicePrincipal.DisabledByMicrosoftStatus.Equal(tfStateServicePrincipal.DisabledByMicrosoftStatus) {
		tfPlanDisabledByMicrosoftStatus := tfPlanServicePrincipal.DisabledByMicrosoftStatus.ValueString()
		requestBodyServicePrincipal.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	}

	if !tfPlanServicePrincipal.DisplayName.Equal(tfStateServicePrincipal.DisplayName) {
		tfPlanDisplayName := tfPlanServicePrincipal.DisplayName.ValueString()
		requestBodyServicePrincipal.SetDisplayName(&tfPlanDisplayName)
	}

	if !tfPlanServicePrincipal.Homepage.Equal(tfStateServicePrincipal.Homepage) {
		tfPlanHomepage := tfPlanServicePrincipal.Homepage.ValueString()
		requestBodyServicePrincipal.SetHomepage(&tfPlanHomepage)
	}

	if !tfPlanServicePrincipal.Id.Equal(tfStateServicePrincipal.Id) {
		tfPlanId := tfPlanServicePrincipal.Id.ValueString()
		requestBodyServicePrincipal.SetId(&tfPlanId)
	}

	if !tfPlanServicePrincipal.Info.Equal(tfStateServicePrincipal.Info) {
		requestBodyInformationalUrl := models.NewInformationalUrl()
		tfPlanInformationalUrl := servicePrincipalInformationalUrlModel{}
		tfPlanServicePrincipal.Info.As(ctx, &tfPlanInformationalUrl, basetypes.ObjectAsOptions{})
		tfStateInformationalUrl := servicePrincipalInformationalUrlModel{}
		tfStateServicePrincipal.Info.As(ctx, &tfStateInformationalUrl, basetypes.ObjectAsOptions{})

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
		requestBodyServicePrincipal.SetInfo(requestBodyInformationalUrl)
		tfPlanServicePrincipal.Info, _ = types.ObjectValueFrom(ctx, tfPlanInformationalUrl.AttributeTypes(), tfPlanInformationalUrl)
	}

	if !tfPlanServicePrincipal.KeyCredentials.Equal(tfStateServicePrincipal.KeyCredentials) {
		var tfPlanKeyCredentials []models.KeyCredentialable
		for k, i := range tfPlanServicePrincipal.KeyCredentials.Elements() {
			requestBodyKeyCredential := models.NewKeyCredential()
			tfPlanKeyCredential := servicePrincipalKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanKeyCredential)
			tfStateKeyCredential := servicePrincipalKeyCredentialModel{}
			types.ListValueFrom(ctx, tfStateServicePrincipal.KeyCredentials.Elements()[k].Type(ctx), &tfPlanKeyCredential)

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
		requestBodyServicePrincipal.SetKeyCredentials(tfPlanKeyCredentials)
	}

	if !tfPlanServicePrincipal.LoginUrl.Equal(tfStateServicePrincipal.LoginUrl) {
		tfPlanLoginUrl := tfPlanServicePrincipal.LoginUrl.ValueString()
		requestBodyServicePrincipal.SetLoginUrl(&tfPlanLoginUrl)
	}

	if !tfPlanServicePrincipal.LogoutUrl.Equal(tfStateServicePrincipal.LogoutUrl) {
		tfPlanLogoutUrl := tfPlanServicePrincipal.LogoutUrl.ValueString()
		requestBodyServicePrincipal.SetLogoutUrl(&tfPlanLogoutUrl)
	}

	if !tfPlanServicePrincipal.Notes.Equal(tfStateServicePrincipal.Notes) {
		tfPlanNotes := tfPlanServicePrincipal.Notes.ValueString()
		requestBodyServicePrincipal.SetNotes(&tfPlanNotes)
	}

	if !tfPlanServicePrincipal.NotificationEmailAddresses.Equal(tfStateServicePrincipal.NotificationEmailAddresses) {
		var stringArrayNotificationEmailAddresses []string
		for _, i := range tfPlanServicePrincipal.NotificationEmailAddresses.Elements() {
			stringArrayNotificationEmailAddresses = append(stringArrayNotificationEmailAddresses, i.String())
		}
		requestBodyServicePrincipal.SetNotificationEmailAddresses(stringArrayNotificationEmailAddresses)
	}

	if !tfPlanServicePrincipal.Oauth2PermissionScopes.Equal(tfStateServicePrincipal.Oauth2PermissionScopes) {
		var tfPlanOauth2PermissionScopes []models.PermissionScopeable
		for k, i := range tfPlanServicePrincipal.Oauth2PermissionScopes.Elements() {
			requestBodyPermissionScope := models.NewPermissionScope()
			tfPlanPermissionScope := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPermissionScope)
			tfStatePermissionScope := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, tfStateServicePrincipal.Oauth2PermissionScopes.Elements()[k].Type(ctx), &tfPlanPermissionScope)

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
		requestBodyServicePrincipal.SetOauth2PermissionScopes(tfPlanOauth2PermissionScopes)
	}

	if !tfPlanServicePrincipal.PasswordCredentials.Equal(tfStateServicePrincipal.PasswordCredentials) {
		var tfPlanPasswordCredentials []models.PasswordCredentialable
		for k, i := range tfPlanServicePrincipal.PasswordCredentials.Elements() {
			requestBodyPasswordCredential := models.NewPasswordCredential()
			tfPlanPasswordCredential := servicePrincipalPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanPasswordCredential)
			tfStatePasswordCredential := servicePrincipalPasswordCredentialModel{}
			types.ListValueFrom(ctx, tfStateServicePrincipal.PasswordCredentials.Elements()[k].Type(ctx), &tfPlanPasswordCredential)

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
		requestBodyServicePrincipal.SetPasswordCredentials(tfPlanPasswordCredentials)
	}

	if !tfPlanServicePrincipal.PreferredSingleSignOnMode.Equal(tfStateServicePrincipal.PreferredSingleSignOnMode) {
		tfPlanPreferredSingleSignOnMode := tfPlanServicePrincipal.PreferredSingleSignOnMode.ValueString()
		requestBodyServicePrincipal.SetPreferredSingleSignOnMode(&tfPlanPreferredSingleSignOnMode)
	}

	if !tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint.Equal(tfStateServicePrincipal.PreferredTokenSigningKeyThumbprint) {
		tfPlanPreferredTokenSigningKeyThumbprint := tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint.ValueString()
		requestBodyServicePrincipal.SetPreferredTokenSigningKeyThumbprint(&tfPlanPreferredTokenSigningKeyThumbprint)
	}

	if !tfPlanServicePrincipal.ReplyUrls.Equal(tfStateServicePrincipal.ReplyUrls) {
		var stringArrayReplyUrls []string
		for _, i := range tfPlanServicePrincipal.ReplyUrls.Elements() {
			stringArrayReplyUrls = append(stringArrayReplyUrls, i.String())
		}
		requestBodyServicePrincipal.SetReplyUrls(stringArrayReplyUrls)
	}

	if !tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.Equal(tfStateServicePrincipal.ResourceSpecificApplicationPermissions) {
		var tfPlanResourceSpecificApplicationPermissions []models.ResourceSpecificPermissionable
		for k, i := range tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.Elements() {
			requestBodyResourceSpecificPermission := models.NewResourceSpecificPermission()
			tfPlanResourceSpecificPermission := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanResourceSpecificPermission)
			tfStateResourceSpecificPermission := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, tfStateServicePrincipal.ResourceSpecificApplicationPermissions.Elements()[k].Type(ctx), &tfPlanResourceSpecificPermission)

			if !tfPlanResourceSpecificPermission.Description.Equal(tfStateResourceSpecificPermission.Description) {
				tfPlanDescription := tfPlanResourceSpecificPermission.Description.ValueString()
				requestBodyResourceSpecificPermission.SetDescription(&tfPlanDescription)
			}

			if !tfPlanResourceSpecificPermission.DisplayName.Equal(tfStateResourceSpecificPermission.DisplayName) {
				tfPlanDisplayName := tfPlanResourceSpecificPermission.DisplayName.ValueString()
				requestBodyResourceSpecificPermission.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanResourceSpecificPermission.Id.Equal(tfStateResourceSpecificPermission.Id) {
				tfPlanId := tfPlanResourceSpecificPermission.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyResourceSpecificPermission.SetId(&u)
			}

			if !tfPlanResourceSpecificPermission.IsEnabled.Equal(tfStateResourceSpecificPermission.IsEnabled) {
				tfPlanIsEnabled := tfPlanResourceSpecificPermission.IsEnabled.ValueBool()
				requestBodyResourceSpecificPermission.SetIsEnabled(&tfPlanIsEnabled)
			}

			if !tfPlanResourceSpecificPermission.Value.Equal(tfStateResourceSpecificPermission.Value) {
				tfPlanValue := tfPlanResourceSpecificPermission.Value.ValueString()
				requestBodyResourceSpecificPermission.SetValue(&tfPlanValue)
			}
		}
		requestBodyServicePrincipal.SetResourceSpecificApplicationPermissions(tfPlanResourceSpecificApplicationPermissions)
	}

	if !tfPlanServicePrincipal.SamlSingleSignOnSettings.Equal(tfStateServicePrincipal.SamlSingleSignOnSettings) {
		requestBodySamlSingleSignOnSettings := models.NewSamlSingleSignOnSettings()
		tfPlanSamlSingleSignOnSettings := servicePrincipalSamlSingleSignOnSettingsModel{}
		tfPlanServicePrincipal.SamlSingleSignOnSettings.As(ctx, &tfPlanSamlSingleSignOnSettings, basetypes.ObjectAsOptions{})
		tfStateSamlSingleSignOnSettings := servicePrincipalSamlSingleSignOnSettingsModel{}
		tfStateServicePrincipal.SamlSingleSignOnSettings.As(ctx, &tfStateSamlSingleSignOnSettings, basetypes.ObjectAsOptions{})

		if !tfPlanSamlSingleSignOnSettings.RelayState.Equal(tfStateSamlSingleSignOnSettings.RelayState) {
			tfPlanRelayState := tfPlanSamlSingleSignOnSettings.RelayState.ValueString()
			requestBodySamlSingleSignOnSettings.SetRelayState(&tfPlanRelayState)
		}
		requestBodyServicePrincipal.SetSamlSingleSignOnSettings(requestBodySamlSingleSignOnSettings)
		tfPlanServicePrincipal.SamlSingleSignOnSettings, _ = types.ObjectValueFrom(ctx, tfPlanSamlSingleSignOnSettings.AttributeTypes(), tfPlanSamlSingleSignOnSettings)
	}

	if !tfPlanServicePrincipal.ServicePrincipalNames.Equal(tfStateServicePrincipal.ServicePrincipalNames) {
		var stringArrayServicePrincipalNames []string
		for _, i := range tfPlanServicePrincipal.ServicePrincipalNames.Elements() {
			stringArrayServicePrincipalNames = append(stringArrayServicePrincipalNames, i.String())
		}
		requestBodyServicePrincipal.SetServicePrincipalNames(stringArrayServicePrincipalNames)
	}

	if !tfPlanServicePrincipal.ServicePrincipalType.Equal(tfStateServicePrincipal.ServicePrincipalType) {
		tfPlanServicePrincipalType := tfPlanServicePrincipal.ServicePrincipalType.ValueString()
		requestBodyServicePrincipal.SetServicePrincipalType(&tfPlanServicePrincipalType)
	}

	if !tfPlanServicePrincipal.SignInAudience.Equal(tfStateServicePrincipal.SignInAudience) {
		tfPlanSignInAudience := tfPlanServicePrincipal.SignInAudience.ValueString()
		requestBodyServicePrincipal.SetSignInAudience(&tfPlanSignInAudience)
	}

	if !tfPlanServicePrincipal.Tags.Equal(tfStateServicePrincipal.Tags) {
		var stringArrayTags []string
		for _, i := range tfPlanServicePrincipal.Tags.Elements() {
			stringArrayTags = append(stringArrayTags, i.String())
		}
		requestBodyServicePrincipal.SetTags(stringArrayTags)
	}

	if !tfPlanServicePrincipal.TokenEncryptionKeyId.Equal(tfStateServicePrincipal.TokenEncryptionKeyId) {
		tfPlanTokenEncryptionKeyId := tfPlanServicePrincipal.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		requestBodyServicePrincipal.SetTokenEncryptionKeyId(&u)
	}

	if !tfPlanServicePrincipal.VerifiedPublisher.Equal(tfStateServicePrincipal.VerifiedPublisher) {
		requestBodyVerifiedPublisher := models.NewVerifiedPublisher()
		tfPlanVerifiedPublisher := servicePrincipalVerifiedPublisherModel{}
		tfPlanServicePrincipal.VerifiedPublisher.As(ctx, &tfPlanVerifiedPublisher, basetypes.ObjectAsOptions{})
		tfStateVerifiedPublisher := servicePrincipalVerifiedPublisherModel{}
		tfStateServicePrincipal.VerifiedPublisher.As(ctx, &tfStateVerifiedPublisher, basetypes.ObjectAsOptions{})

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
		requestBodyServicePrincipal.SetVerifiedPublisher(requestBodyVerifiedPublisher)
		tfPlanServicePrincipal.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfPlanVerifiedPublisher.AttributeTypes(), tfPlanVerifiedPublisher)
	}

	// Update servicePrincipal
	_, err := r.client.ServicePrincipals().ByServicePrincipalId(tfStateServicePrincipal.Id.ValueString()).Patch(context.Background(), requestBodyServicePrincipal, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating service_principal",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlanServicePrincipal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *servicePrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfStateServicePrincipal servicePrincipalModel
	diags := req.State.Get(ctx, &tfStateServicePrincipal)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete servicePrincipal
	err := r.client.ServicePrincipals().ByServicePrincipalId(tfStateServicePrincipal.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting service_principal",
			err.Error(),
		)
		return
	}

}

func (r *servicePrincipalResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
