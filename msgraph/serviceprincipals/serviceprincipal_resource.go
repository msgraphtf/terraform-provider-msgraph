package serviceprincipals

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
			"custom_security_attributes": schema.SingleNestedAttribute{
				Description: "An open complex type that holds the value of a custom security attribute that is assigned to a directory object. Nullable. Returned only on $select. Supports $filter (eq, ne, not, startsWith). Filter value is case sensitive. To read this property, the calling app must be assigned the CustomSecAttributeAssignment.Read.All permission. To write this property, the calling app must be assigned the CustomSecAttributeAssignment.ReadWrite.All permissions. To read or write this property in delegated scenarios, the admin must be assigned the Attribute Assignment Administrator role.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{},
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
	// START Id | CreateStringAttribute
	if !tfPlanServicePrincipal.Id.IsUnknown() {
		tfPlanId := tfPlanServicePrincipal.Id.ValueString()
		requestBodyServicePrincipal.SetId(&tfPlanId)
	} else {
		tfPlanServicePrincipal.Id = types.StringNull()
	}
	// END Id | CreateStringAttribute

	// START DeletedDateTime | CreateStringTimeAttribute
	if !tfPlanServicePrincipal.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlanServicePrincipal.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyServicePrincipal.SetDeletedDateTime(&t)
	} else {
		tfPlanServicePrincipal.DeletedDateTime = types.StringNull()
	}
	// END DeletedDateTime | CreateStringTimeAttribute

	// START AccountEnabled | CreateBoolAttribute
	if !tfPlanServicePrincipal.AccountEnabled.IsUnknown() {
		tfPlanAccountEnabled := tfPlanServicePrincipal.AccountEnabled.ValueBool()
		requestBodyServicePrincipal.SetAccountEnabled(&tfPlanAccountEnabled)
	} else {
		tfPlanServicePrincipal.AccountEnabled = types.BoolNull()
	}
	// END AccountEnabled | CreateBoolAttribute

	// START AddIns | CreateArrayObjectAttribute
	if len(tfPlanServicePrincipal.AddIns.Elements()) > 0 {
		var requestBodyAddIns []models.AddInable
		for _, i := range tfPlanServicePrincipal.AddIns.Elements() {
			requestBodyAddIn := models.NewAddIn()
			tfPlanAddIn := servicePrincipalAddInModel{}
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
					tfPlanKeyValue := servicePrincipalKeyValueModel{}
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
		requestBodyServicePrincipal.SetAddIns(requestBodyAddIns)
	} else {
		tfPlanServicePrincipal.AddIns = types.ListNull(tfPlanServicePrincipal.AddIns.ElementType(ctx))
	}
	// END AddIns | CreateArrayObjectAttribute

	// START AlternativeNames | CreateArrayStringAttribute
	if len(tfPlanServicePrincipal.AlternativeNames.Elements()) > 0 {
		var stringArrayAlternativeNames []string
		for _, i := range tfPlanServicePrincipal.AlternativeNames.Elements() {
			stringArrayAlternativeNames = append(stringArrayAlternativeNames, i.String())
		}
		requestBodyServicePrincipal.SetAlternativeNames(stringArrayAlternativeNames)
	} else {
		tfPlanServicePrincipal.AlternativeNames = types.ListNull(types.StringType)
	}
	// END AlternativeNames | CreateArrayStringAttribute

	// START AppDescription | CreateStringAttribute
	if !tfPlanServicePrincipal.AppDescription.IsUnknown() {
		tfPlanAppDescription := tfPlanServicePrincipal.AppDescription.ValueString()
		requestBodyServicePrincipal.SetAppDescription(&tfPlanAppDescription)
	} else {
		tfPlanServicePrincipal.AppDescription = types.StringNull()
	}
	// END AppDescription | CreateStringAttribute

	// START AppDisplayName | CreateStringAttribute
	if !tfPlanServicePrincipal.AppDisplayName.IsUnknown() {
		tfPlanAppDisplayName := tfPlanServicePrincipal.AppDisplayName.ValueString()
		requestBodyServicePrincipal.SetAppDisplayName(&tfPlanAppDisplayName)
	} else {
		tfPlanServicePrincipal.AppDisplayName = types.StringNull()
	}
	// END AppDisplayName | CreateStringAttribute

	// START AppId | CreateStringAttribute
	if !tfPlanServicePrincipal.AppId.IsUnknown() {
		tfPlanAppId := tfPlanServicePrincipal.AppId.ValueString()
		requestBodyServicePrincipal.SetAppId(&tfPlanAppId)
	} else {
		tfPlanServicePrincipal.AppId = types.StringNull()
	}
	// END AppId | CreateStringAttribute

	// START AppOwnerOrganizationId | CreateStringUuidAttribute
	if !tfPlanServicePrincipal.AppOwnerOrganizationId.IsUnknown() {
		tfPlanAppOwnerOrganizationId := tfPlanServicePrincipal.AppOwnerOrganizationId.ValueString()
		u, _ := uuid.Parse(tfPlanAppOwnerOrganizationId)
		requestBodyServicePrincipal.SetAppOwnerOrganizationId(&u)
	} else {
		tfPlanServicePrincipal.AppOwnerOrganizationId = types.StringNull()
	}
	// END AppOwnerOrganizationId | CreateStringUuidAttribute

	// START AppRoleAssignmentRequired | CreateBoolAttribute
	if !tfPlanServicePrincipal.AppRoleAssignmentRequired.IsUnknown() {
		tfPlanAppRoleAssignmentRequired := tfPlanServicePrincipal.AppRoleAssignmentRequired.ValueBool()
		requestBodyServicePrincipal.SetAppRoleAssignmentRequired(&tfPlanAppRoleAssignmentRequired)
	} else {
		tfPlanServicePrincipal.AppRoleAssignmentRequired = types.BoolNull()
	}
	// END AppRoleAssignmentRequired | CreateBoolAttribute

	// START AppRoles | CreateArrayObjectAttribute
	if len(tfPlanServicePrincipal.AppRoles.Elements()) > 0 {
		var requestBodyAppRoles []models.AppRoleable
		for _, i := range tfPlanServicePrincipal.AppRoles.Elements() {
			requestBodyAppRole := models.NewAppRole()
			tfPlanAppRole := servicePrincipalAppRoleModel{}
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
		requestBodyServicePrincipal.SetAppRoles(requestBodyAppRoles)
	} else {
		tfPlanServicePrincipal.AppRoles = types.ListNull(tfPlanServicePrincipal.AppRoles.ElementType(ctx))
	}
	// END AppRoles | CreateArrayObjectAttribute

	// START ApplicationTemplateId | CreateStringAttribute
	if !tfPlanServicePrincipal.ApplicationTemplateId.IsUnknown() {
		tfPlanApplicationTemplateId := tfPlanServicePrincipal.ApplicationTemplateId.ValueString()
		requestBodyServicePrincipal.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	} else {
		tfPlanServicePrincipal.ApplicationTemplateId = types.StringNull()
	}
	// END ApplicationTemplateId | CreateStringAttribute

	// START CustomSecurityAttributes | CreateObjectAttribute
	if !tfPlanServicePrincipal.CustomSecurityAttributes.IsUnknown() {
		requestBodyCustomSecurityAttributeValue := models.NewCustomSecurityAttributeValue()
		tfPlanCustomSecurityAttributeValue := servicePrincipalCustomSecurityAttributeValueModel{}
		tfPlanServicePrincipal.CustomSecurityAttributes.As(ctx, &tfPlanCustomSecurityAttributeValue, basetypes.ObjectAsOptions{})

		requestBodyServicePrincipal.SetCustomSecurityAttributes(requestBodyCustomSecurityAttributeValue)
		tfPlanServicePrincipal.CustomSecurityAttributes, _ = types.ObjectValueFrom(ctx, tfPlanCustomSecurityAttributeValue.AttributeTypes(), requestBodyCustomSecurityAttributeValue)
	} else {
		tfPlanServicePrincipal.CustomSecurityAttributes = types.ObjectNull(tfPlanServicePrincipal.CustomSecurityAttributes.AttributeTypes(ctx))
	}
	// END CustomSecurityAttributes | CreateObjectAttribute

	// START Description | CreateStringAttribute
	if !tfPlanServicePrincipal.Description.IsUnknown() {
		tfPlanDescription := tfPlanServicePrincipal.Description.ValueString()
		requestBodyServicePrincipal.SetDescription(&tfPlanDescription)
	} else {
		tfPlanServicePrincipal.Description = types.StringNull()
	}
	// END Description | CreateStringAttribute

	// START DisabledByMicrosoftStatus | CreateStringAttribute
	if !tfPlanServicePrincipal.DisabledByMicrosoftStatus.IsUnknown() {
		tfPlanDisabledByMicrosoftStatus := tfPlanServicePrincipal.DisabledByMicrosoftStatus.ValueString()
		requestBodyServicePrincipal.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	} else {
		tfPlanServicePrincipal.DisabledByMicrosoftStatus = types.StringNull()
	}
	// END DisabledByMicrosoftStatus | CreateStringAttribute

	// START DisplayName | CreateStringAttribute
	if !tfPlanServicePrincipal.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanServicePrincipal.DisplayName.ValueString()
		requestBodyServicePrincipal.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanServicePrincipal.DisplayName = types.StringNull()
	}
	// END DisplayName | CreateStringAttribute

	// START Homepage | CreateStringAttribute
	if !tfPlanServicePrincipal.Homepage.IsUnknown() {
		tfPlanHomepage := tfPlanServicePrincipal.Homepage.ValueString()
		requestBodyServicePrincipal.SetHomepage(&tfPlanHomepage)
	} else {
		tfPlanServicePrincipal.Homepage = types.StringNull()
	}
	// END Homepage | CreateStringAttribute

	// START Info | CreateObjectAttribute
	if !tfPlanServicePrincipal.Info.IsUnknown() {
		requestBodyInformationalUrl := models.NewInformationalUrl()
		tfPlanInformationalUrl := servicePrincipalInformationalUrlModel{}
		tfPlanServicePrincipal.Info.As(ctx, &tfPlanInformationalUrl, basetypes.ObjectAsOptions{})

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

		requestBodyServicePrincipal.SetInfo(requestBodyInformationalUrl)
		tfPlanServicePrincipal.Info, _ = types.ObjectValueFrom(ctx, tfPlanInformationalUrl.AttributeTypes(), requestBodyInformationalUrl)
	} else {
		tfPlanServicePrincipal.Info = types.ObjectNull(tfPlanServicePrincipal.Info.AttributeTypes(ctx))
	}
	// END Info | CreateObjectAttribute

	// START KeyCredentials | CreateArrayObjectAttribute
	if len(tfPlanServicePrincipal.KeyCredentials.Elements()) > 0 {
		var requestBodyKeyCredentials []models.KeyCredentialable
		for _, i := range tfPlanServicePrincipal.KeyCredentials.Elements() {
			requestBodyKeyCredential := models.NewKeyCredential()
			tfPlanKeyCredential := servicePrincipalKeyCredentialModel{}
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
		requestBodyServicePrincipal.SetKeyCredentials(requestBodyKeyCredentials)
	} else {
		tfPlanServicePrincipal.KeyCredentials = types.ListNull(tfPlanServicePrincipal.KeyCredentials.ElementType(ctx))
	}
	// END KeyCredentials | CreateArrayObjectAttribute

	// START LoginUrl | CreateStringAttribute
	if !tfPlanServicePrincipal.LoginUrl.IsUnknown() {
		tfPlanLoginUrl := tfPlanServicePrincipal.LoginUrl.ValueString()
		requestBodyServicePrincipal.SetLoginUrl(&tfPlanLoginUrl)
	} else {
		tfPlanServicePrincipal.LoginUrl = types.StringNull()
	}
	// END LoginUrl | CreateStringAttribute

	// START LogoutUrl | CreateStringAttribute
	if !tfPlanServicePrincipal.LogoutUrl.IsUnknown() {
		tfPlanLogoutUrl := tfPlanServicePrincipal.LogoutUrl.ValueString()
		requestBodyServicePrincipal.SetLogoutUrl(&tfPlanLogoutUrl)
	} else {
		tfPlanServicePrincipal.LogoutUrl = types.StringNull()
	}
	// END LogoutUrl | CreateStringAttribute

	// START Notes | CreateStringAttribute
	if !tfPlanServicePrincipal.Notes.IsUnknown() {
		tfPlanNotes := tfPlanServicePrincipal.Notes.ValueString()
		requestBodyServicePrincipal.SetNotes(&tfPlanNotes)
	} else {
		tfPlanServicePrincipal.Notes = types.StringNull()
	}
	// END Notes | CreateStringAttribute

	// START NotificationEmailAddresses | CreateArrayStringAttribute
	if len(tfPlanServicePrincipal.NotificationEmailAddresses.Elements()) > 0 {
		var stringArrayNotificationEmailAddresses []string
		for _, i := range tfPlanServicePrincipal.NotificationEmailAddresses.Elements() {
			stringArrayNotificationEmailAddresses = append(stringArrayNotificationEmailAddresses, i.String())
		}
		requestBodyServicePrincipal.SetNotificationEmailAddresses(stringArrayNotificationEmailAddresses)
	} else {
		tfPlanServicePrincipal.NotificationEmailAddresses = types.ListNull(types.StringType)
	}
	// END NotificationEmailAddresses | CreateArrayStringAttribute

	// START Oauth2PermissionScopes | CreateArrayObjectAttribute
	if len(tfPlanServicePrincipal.Oauth2PermissionScopes.Elements()) > 0 {
		var requestBodyOauth2PermissionScopes []models.PermissionScopeable
		for _, i := range tfPlanServicePrincipal.Oauth2PermissionScopes.Elements() {
			requestBodyPermissionScope := models.NewPermissionScope()
			tfPlanPermissionScope := servicePrincipalPermissionScopeModel{}
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
		requestBodyServicePrincipal.SetOauth2PermissionScopes(requestBodyOauth2PermissionScopes)
	} else {
		tfPlanServicePrincipal.Oauth2PermissionScopes = types.ListNull(tfPlanServicePrincipal.Oauth2PermissionScopes.ElementType(ctx))
	}
	// END Oauth2PermissionScopes | CreateArrayObjectAttribute

	// START PasswordCredentials | CreateArrayObjectAttribute
	if len(tfPlanServicePrincipal.PasswordCredentials.Elements()) > 0 {
		var requestBodyPasswordCredentials []models.PasswordCredentialable
		for _, i := range tfPlanServicePrincipal.PasswordCredentials.Elements() {
			requestBodyPasswordCredential := models.NewPasswordCredential()
			tfPlanPasswordCredential := servicePrincipalPasswordCredentialModel{}
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
		requestBodyServicePrincipal.SetPasswordCredentials(requestBodyPasswordCredentials)
	} else {
		tfPlanServicePrincipal.PasswordCredentials = types.ListNull(tfPlanServicePrincipal.PasswordCredentials.ElementType(ctx))
	}
	// END PasswordCredentials | CreateArrayObjectAttribute

	// START PreferredSingleSignOnMode | CreateStringAttribute
	if !tfPlanServicePrincipal.PreferredSingleSignOnMode.IsUnknown() {
		tfPlanPreferredSingleSignOnMode := tfPlanServicePrincipal.PreferredSingleSignOnMode.ValueString()
		requestBodyServicePrincipal.SetPreferredSingleSignOnMode(&tfPlanPreferredSingleSignOnMode)
	} else {
		tfPlanServicePrincipal.PreferredSingleSignOnMode = types.StringNull()
	}
	// END PreferredSingleSignOnMode | CreateStringAttribute

	// START PreferredTokenSigningKeyThumbprint | CreateStringAttribute
	if !tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint.IsUnknown() {
		tfPlanPreferredTokenSigningKeyThumbprint := tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint.ValueString()
		requestBodyServicePrincipal.SetPreferredTokenSigningKeyThumbprint(&tfPlanPreferredTokenSigningKeyThumbprint)
	} else {
		tfPlanServicePrincipal.PreferredTokenSigningKeyThumbprint = types.StringNull()
	}
	// END PreferredTokenSigningKeyThumbprint | CreateStringAttribute

	// START ReplyUrls | CreateArrayStringAttribute
	if len(tfPlanServicePrincipal.ReplyUrls.Elements()) > 0 {
		var stringArrayReplyUrls []string
		for _, i := range tfPlanServicePrincipal.ReplyUrls.Elements() {
			stringArrayReplyUrls = append(stringArrayReplyUrls, i.String())
		}
		requestBodyServicePrincipal.SetReplyUrls(stringArrayReplyUrls)
	} else {
		tfPlanServicePrincipal.ReplyUrls = types.ListNull(types.StringType)
	}
	// END ReplyUrls | CreateArrayStringAttribute

	// START ResourceSpecificApplicationPermissions | CreateArrayObjectAttribute
	if len(tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.Elements()) > 0 {
		var requestBodyResourceSpecificApplicationPermissions []models.ResourceSpecificPermissionable
		for _, i := range tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.Elements() {
			requestBodyResourceSpecificPermission := models.NewResourceSpecificPermission()
			tfPlanResourceSpecificPermission := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanResourceSpecificPermission)

			// START Description | CreateStringAttribute
			if !tfPlanResourceSpecificPermission.Description.IsUnknown() {
				tfPlanDescription := tfPlanResourceSpecificPermission.Description.ValueString()
				requestBodyResourceSpecificPermission.SetDescription(&tfPlanDescription)
			} else {
				tfPlanResourceSpecificPermission.Description = types.StringNull()
			}
			// END Description | CreateStringAttribute

			// START DisplayName | CreateStringAttribute
			if !tfPlanResourceSpecificPermission.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanResourceSpecificPermission.DisplayName.ValueString()
				requestBodyResourceSpecificPermission.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanResourceSpecificPermission.DisplayName = types.StringNull()
			}
			// END DisplayName | CreateStringAttribute

			// START Id | CreateStringUuidAttribute
			if !tfPlanResourceSpecificPermission.Id.IsUnknown() {
				tfPlanId := tfPlanResourceSpecificPermission.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyResourceSpecificPermission.SetId(&u)
			} else {
				tfPlanResourceSpecificPermission.Id = types.StringNull()
			}
			// END Id | CreateStringUuidAttribute

			// START IsEnabled | CreateBoolAttribute
			if !tfPlanResourceSpecificPermission.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := tfPlanResourceSpecificPermission.IsEnabled.ValueBool()
				requestBodyResourceSpecificPermission.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				tfPlanResourceSpecificPermission.IsEnabled = types.BoolNull()
			}
			// END IsEnabled | CreateBoolAttribute

			// START Value | CreateStringAttribute
			if !tfPlanResourceSpecificPermission.Value.IsUnknown() {
				tfPlanValue := tfPlanResourceSpecificPermission.Value.ValueString()
				requestBodyResourceSpecificPermission.SetValue(&tfPlanValue)
			} else {
				tfPlanResourceSpecificPermission.Value = types.StringNull()
			}
			// END Value | CreateStringAttribute

		}
		requestBodyServicePrincipal.SetResourceSpecificApplicationPermissions(requestBodyResourceSpecificApplicationPermissions)
	} else {
		tfPlanServicePrincipal.ResourceSpecificApplicationPermissions = types.ListNull(tfPlanServicePrincipal.ResourceSpecificApplicationPermissions.ElementType(ctx))
	}
	// END ResourceSpecificApplicationPermissions | CreateArrayObjectAttribute

	// START SamlSingleSignOnSettings | CreateObjectAttribute
	if !tfPlanServicePrincipal.SamlSingleSignOnSettings.IsUnknown() {
		requestBodySamlSingleSignOnSettings := models.NewSamlSingleSignOnSettings()
		tfPlanSamlSingleSignOnSettings := servicePrincipalSamlSingleSignOnSettingsModel{}
		tfPlanServicePrincipal.SamlSingleSignOnSettings.As(ctx, &tfPlanSamlSingleSignOnSettings, basetypes.ObjectAsOptions{})

		// START RelayState | CreateStringAttribute
		if !tfPlanSamlSingleSignOnSettings.RelayState.IsUnknown() {
			tfPlanRelayState := tfPlanSamlSingleSignOnSettings.RelayState.ValueString()
			requestBodySamlSingleSignOnSettings.SetRelayState(&tfPlanRelayState)
		} else {
			tfPlanSamlSingleSignOnSettings.RelayState = types.StringNull()
		}
		// END RelayState | CreateStringAttribute

		requestBodyServicePrincipal.SetSamlSingleSignOnSettings(requestBodySamlSingleSignOnSettings)
		tfPlanServicePrincipal.SamlSingleSignOnSettings, _ = types.ObjectValueFrom(ctx, tfPlanSamlSingleSignOnSettings.AttributeTypes(), requestBodySamlSingleSignOnSettings)
	} else {
		tfPlanServicePrincipal.SamlSingleSignOnSettings = types.ObjectNull(tfPlanServicePrincipal.SamlSingleSignOnSettings.AttributeTypes(ctx))
	}
	// END SamlSingleSignOnSettings | CreateObjectAttribute

	// START ServicePrincipalNames | CreateArrayStringAttribute
	if len(tfPlanServicePrincipal.ServicePrincipalNames.Elements()) > 0 {
		var stringArrayServicePrincipalNames []string
		for _, i := range tfPlanServicePrincipal.ServicePrincipalNames.Elements() {
			stringArrayServicePrincipalNames = append(stringArrayServicePrincipalNames, i.String())
		}
		requestBodyServicePrincipal.SetServicePrincipalNames(stringArrayServicePrincipalNames)
	} else {
		tfPlanServicePrincipal.ServicePrincipalNames = types.ListNull(types.StringType)
	}
	// END ServicePrincipalNames | CreateArrayStringAttribute

	// START ServicePrincipalType | CreateStringAttribute
	if !tfPlanServicePrincipal.ServicePrincipalType.IsUnknown() {
		tfPlanServicePrincipalType := tfPlanServicePrincipal.ServicePrincipalType.ValueString()
		requestBodyServicePrincipal.SetServicePrincipalType(&tfPlanServicePrincipalType)
	} else {
		tfPlanServicePrincipal.ServicePrincipalType = types.StringNull()
	}
	// END ServicePrincipalType | CreateStringAttribute

	// START SignInAudience | CreateStringAttribute
	if !tfPlanServicePrincipal.SignInAudience.IsUnknown() {
		tfPlanSignInAudience := tfPlanServicePrincipal.SignInAudience.ValueString()
		requestBodyServicePrincipal.SetSignInAudience(&tfPlanSignInAudience)
	} else {
		tfPlanServicePrincipal.SignInAudience = types.StringNull()
	}
	// END SignInAudience | CreateStringAttribute

	// START Tags | CreateArrayStringAttribute
	if len(tfPlanServicePrincipal.Tags.Elements()) > 0 {
		var stringArrayTags []string
		for _, i := range tfPlanServicePrincipal.Tags.Elements() {
			stringArrayTags = append(stringArrayTags, i.String())
		}
		requestBodyServicePrincipal.SetTags(stringArrayTags)
	} else {
		tfPlanServicePrincipal.Tags = types.ListNull(types.StringType)
	}
	// END Tags | CreateArrayStringAttribute

	// START TokenEncryptionKeyId | CreateStringUuidAttribute
	if !tfPlanServicePrincipal.TokenEncryptionKeyId.IsUnknown() {
		tfPlanTokenEncryptionKeyId := tfPlanServicePrincipal.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		requestBodyServicePrincipal.SetTokenEncryptionKeyId(&u)
	} else {
		tfPlanServicePrincipal.TokenEncryptionKeyId = types.StringNull()
	}
	// END TokenEncryptionKeyId | CreateStringUuidAttribute

	// START VerifiedPublisher | CreateObjectAttribute
	if !tfPlanServicePrincipal.VerifiedPublisher.IsUnknown() {
		requestBodyVerifiedPublisher := models.NewVerifiedPublisher()
		tfPlanVerifiedPublisher := servicePrincipalVerifiedPublisherModel{}
		tfPlanServicePrincipal.VerifiedPublisher.As(ctx, &tfPlanVerifiedPublisher, basetypes.ObjectAsOptions{})

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

		requestBodyServicePrincipal.SetVerifiedPublisher(requestBodyVerifiedPublisher)
		tfPlanServicePrincipal.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfPlanVerifiedPublisher.AttributeTypes(), requestBodyVerifiedPublisher)
	} else {
		tfPlanServicePrincipal.VerifiedPublisher = types.ObjectNull(tfPlanServicePrincipal.VerifiedPublisher.AttributeTypes(ctx))
	}
	// END VerifiedPublisher | CreateObjectAttribute

	// Create new servicePrincipal
	result, err := r.client.ServicePrincipals().Post(context.Background(), requestBodyServicePrincipal, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service_principal",
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
	var state servicePrincipalModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := serviceprincipals.ServicePrincipalItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipals.ServicePrincipalItemRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"deletedDateTime",
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
				"customSecurityAttributes",
				"description",
				"disabledByMicrosoftStatus",
				"displayName",
				"homepage",
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

	var result models.ServicePrincipalable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.ServicePrincipals().ByServicePrincipalId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting service_principal",
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
	if result.GetAccountEnabled() != nil {
		state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	} else {
		state.AccountEnabled = types.BoolNull()
	}
	if len(result.GetAddIns()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAddIns() {
			addIns := new(servicePrincipalAddInModel)

			if v.GetId() != nil {
				addIns.Id = types.StringValue(v.GetId().String())
			} else {
				addIns.Id = types.StringNull()
			}
			if len(v.GetProperties()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetProperties() {
					properties := new(servicePrincipalKeyValueModel)

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
	if len(result.GetAlternativeNames()) > 0 {
		var alternativeNames []attr.Value
		for _, v := range result.GetAlternativeNames() {
			alternativeNames = append(alternativeNames, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, alternativeNames)
		state.AlternativeNames = listValue
	} else {
		state.AlternativeNames = types.ListNull(types.StringType)
	}
	if result.GetAppDescription() != nil {
		state.AppDescription = types.StringValue(*result.GetAppDescription())
	} else {
		state.AppDescription = types.StringNull()
	}
	if result.GetAppDisplayName() != nil {
		state.AppDisplayName = types.StringValue(*result.GetAppDisplayName())
	} else {
		state.AppDisplayName = types.StringNull()
	}
	if result.GetAppId() != nil {
		state.AppId = types.StringValue(*result.GetAppId())
	} else {
		state.AppId = types.StringNull()
	}
	if result.GetAppOwnerOrganizationId() != nil {
		state.AppOwnerOrganizationId = types.StringValue(result.GetAppOwnerOrganizationId().String())
	} else {
		state.AppOwnerOrganizationId = types.StringNull()
	}
	if result.GetAppRoleAssignmentRequired() != nil {
		state.AppRoleAssignmentRequired = types.BoolValue(*result.GetAppRoleAssignmentRequired())
	} else {
		state.AppRoleAssignmentRequired = types.BoolNull()
	}
	if len(result.GetAppRoles()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAppRoles() {
			appRoles := new(servicePrincipalAppRoleModel)

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
	if result.GetCustomSecurityAttributes() != nil {
		customSecurityAttributes := new(servicePrincipalCustomSecurityAttributeValueModel)

		objectValue, _ := types.ObjectValueFrom(ctx, customSecurityAttributes.AttributeTypes(), customSecurityAttributes)
		state.CustomSecurityAttributes = objectValue
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
	if result.GetHomepage() != nil {
		state.Homepage = types.StringValue(*result.GetHomepage())
	} else {
		state.Homepage = types.StringNull()
	}
	if result.GetInfo() != nil {
		info := new(servicePrincipalInformationalUrlModel)

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
	if len(result.GetKeyCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetKeyCredentials() {
			keyCredentials := new(servicePrincipalKeyCredentialModel)

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
	if result.GetLoginUrl() != nil {
		state.LoginUrl = types.StringValue(*result.GetLoginUrl())
	} else {
		state.LoginUrl = types.StringNull()
	}
	if result.GetLogoutUrl() != nil {
		state.LogoutUrl = types.StringValue(*result.GetLogoutUrl())
	} else {
		state.LogoutUrl = types.StringNull()
	}
	if result.GetNotes() != nil {
		state.Notes = types.StringValue(*result.GetNotes())
	} else {
		state.Notes = types.StringNull()
	}
	if len(result.GetNotificationEmailAddresses()) > 0 {
		var notificationEmailAddresses []attr.Value
		for _, v := range result.GetNotificationEmailAddresses() {
			notificationEmailAddresses = append(notificationEmailAddresses, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, notificationEmailAddresses)
		state.NotificationEmailAddresses = listValue
	} else {
		state.NotificationEmailAddresses = types.ListNull(types.StringType)
	}
	if len(result.GetOauth2PermissionScopes()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetOauth2PermissionScopes() {
			oauth2PermissionScopes := new(servicePrincipalPermissionScopeModel)

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
		state.Oauth2PermissionScopes, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(result.GetPasswordCredentials()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetPasswordCredentials() {
			passwordCredentials := new(servicePrincipalPasswordCredentialModel)

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
	if result.GetPreferredSingleSignOnMode() != nil {
		state.PreferredSingleSignOnMode = types.StringValue(*result.GetPreferredSingleSignOnMode())
	} else {
		state.PreferredSingleSignOnMode = types.StringNull()
	}
	if result.GetPreferredTokenSigningKeyThumbprint() != nil {
		state.PreferredTokenSigningKeyThumbprint = types.StringValue(*result.GetPreferredTokenSigningKeyThumbprint())
	} else {
		state.PreferredTokenSigningKeyThumbprint = types.StringNull()
	}
	if len(result.GetReplyUrls()) > 0 {
		var replyUrls []attr.Value
		for _, v := range result.GetReplyUrls() {
			replyUrls = append(replyUrls, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, replyUrls)
		state.ReplyUrls = listValue
	} else {
		state.ReplyUrls = types.ListNull(types.StringType)
	}
	if len(result.GetResourceSpecificApplicationPermissions()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetResourceSpecificApplicationPermissions() {
			resourceSpecificApplicationPermissions := new(servicePrincipalResourceSpecificPermissionModel)

			if v.GetDescription() != nil {
				resourceSpecificApplicationPermissions.Description = types.StringValue(*v.GetDescription())
			} else {
				resourceSpecificApplicationPermissions.Description = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				resourceSpecificApplicationPermissions.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				resourceSpecificApplicationPermissions.DisplayName = types.StringNull()
			}
			if v.GetId() != nil {
				resourceSpecificApplicationPermissions.Id = types.StringValue(v.GetId().String())
			} else {
				resourceSpecificApplicationPermissions.Id = types.StringNull()
			}
			if v.GetIsEnabled() != nil {
				resourceSpecificApplicationPermissions.IsEnabled = types.BoolValue(*v.GetIsEnabled())
			} else {
				resourceSpecificApplicationPermissions.IsEnabled = types.BoolNull()
			}
			if v.GetValue() != nil {
				resourceSpecificApplicationPermissions.Value = types.StringValue(*v.GetValue())
			} else {
				resourceSpecificApplicationPermissions.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, resourceSpecificApplicationPermissions.AttributeTypes(), resourceSpecificApplicationPermissions)
			objectValues = append(objectValues, objectValue)
		}
		state.ResourceSpecificApplicationPermissions, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetSamlSingleSignOnSettings() != nil {
		samlSingleSignOnSettings := new(servicePrincipalSamlSingleSignOnSettingsModel)

		if result.GetSamlSingleSignOnSettings().GetRelayState() != nil {
			samlSingleSignOnSettings.RelayState = types.StringValue(*result.GetSamlSingleSignOnSettings().GetRelayState())
		} else {
			samlSingleSignOnSettings.RelayState = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, samlSingleSignOnSettings.AttributeTypes(), samlSingleSignOnSettings)
		state.SamlSingleSignOnSettings = objectValue
	}
	if len(result.GetServicePrincipalNames()) > 0 {
		var servicePrincipalNames []attr.Value
		for _, v := range result.GetServicePrincipalNames() {
			servicePrincipalNames = append(servicePrincipalNames, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, servicePrincipalNames)
		state.ServicePrincipalNames = listValue
	} else {
		state.ServicePrincipalNames = types.ListNull(types.StringType)
	}
	if result.GetServicePrincipalType() != nil {
		state.ServicePrincipalType = types.StringValue(*result.GetServicePrincipalType())
	} else {
		state.ServicePrincipalType = types.StringNull()
	}
	if result.GetSignInAudience() != nil {
		state.SignInAudience = types.StringValue(*result.GetSignInAudience())
	} else {
		state.SignInAudience = types.StringNull()
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
	if result.GetVerifiedPublisher() != nil {
		verifiedPublisher := new(servicePrincipalVerifiedPublisherModel)

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

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *servicePrincipalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlan servicePrincipalModel
	diags := req.Plan.Get(ctx, &tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfState servicePrincipalModel
	diags = req.State.Get(ctx, &tfState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewServicePrincipal()

	if !tfPlan.Id.Equal(tfState.Id) {
		tfPlanId := tfPlan.Id.ValueString()
		requestBody.SetId(&tfPlanId)
	}

	if !tfPlan.DeletedDateTime.Equal(tfState.DeletedDateTime) {
		tfPlanDeletedDateTime := tfPlan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !tfPlan.AccountEnabled.Equal(tfState.AccountEnabled) {
		tfPlanAccountEnabled := tfPlan.AccountEnabled.ValueBool()
		requestBody.SetAccountEnabled(&tfPlanAccountEnabled)
	}

	if !tfPlan.AddIns.Equal(tfState.AddIns) {
		var tfPlanAddIns []models.AddInable
		for k, i := range tfPlan.AddIns.Elements() {
			requestBodyAddIns := models.NewAddIn()
			tfPlanrequestBodyAddIns := servicePrincipalAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyAddIns)
			tfStaterequestBodyAddIns := servicePrincipalAddInModel{}
			types.ListValueFrom(ctx, tfState.AddIns.Elements()[k].Type(ctx), &tfPlanrequestBodyAddIns)

			if !tfPlanrequestBodyAddIns.Id.Equal(tfStaterequestBodyAddIns.Id) {
				tfPlanId := tfPlanrequestBodyAddIns.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAddIns.SetId(&u)
			}

			if !tfPlanrequestBodyAddIns.Properties.Equal(tfStaterequestBodyAddIns.Properties) {
				var tfPlanProperties []models.KeyValueable
				for k, i := range tfPlanrequestBodyAddIns.Properties.Elements() {
					requestBodyProperties := models.NewKeyValue()
					tfPlanrequestBodyProperties := servicePrincipalKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyProperties)
					tfStaterequestBodyProperties := servicePrincipalKeyValueModel{}
					types.ListValueFrom(ctx, tfStaterequestBodyAddIns.Properties.Elements()[k].Type(ctx), &tfPlanrequestBodyProperties)

					if !tfPlanrequestBodyProperties.Key.Equal(tfStaterequestBodyProperties.Key) {
						tfPlanKey := tfPlanrequestBodyProperties.Key.ValueString()
						requestBodyProperties.SetKey(&tfPlanKey)
					}

					if !tfPlanrequestBodyProperties.Value.Equal(tfStaterequestBodyProperties.Value) {
						tfPlanValue := tfPlanrequestBodyProperties.Value.ValueString()
						requestBodyProperties.SetValue(&tfPlanValue)
					}
				}
				requestBodyAddIns.SetProperties(tfPlanProperties)
			}

			if !tfPlanrequestBodyAddIns.Type.Equal(tfStaterequestBodyAddIns.Type) {
				tfPlanType := tfPlanrequestBodyAddIns.Type.ValueString()
				requestBodyAddIns.SetTypeEscaped(&tfPlanType)
			}
		}
		requestBody.SetAddIns(tfPlanAddIns)
	}

	if !tfPlan.AlternativeNames.Equal(tfState.AlternativeNames) {
		var stringArrayAlternativeNames []string
		for _, i := range tfPlan.AlternativeNames.Elements() {
			stringArrayAlternativeNames = append(stringArrayAlternativeNames, i.String())
		}
		requestBody.SetAlternativeNames(stringArrayAlternativeNames)
	}

	if !tfPlan.AppDescription.Equal(tfState.AppDescription) {
		tfPlanAppDescription := tfPlan.AppDescription.ValueString()
		requestBody.SetAppDescription(&tfPlanAppDescription)
	}

	if !tfPlan.AppDisplayName.Equal(tfState.AppDisplayName) {
		tfPlanAppDisplayName := tfPlan.AppDisplayName.ValueString()
		requestBody.SetAppDisplayName(&tfPlanAppDisplayName)
	}

	if !tfPlan.AppId.Equal(tfState.AppId) {
		tfPlanAppId := tfPlan.AppId.ValueString()
		requestBody.SetAppId(&tfPlanAppId)
	}

	if !tfPlan.AppOwnerOrganizationId.Equal(tfState.AppOwnerOrganizationId) {
		tfPlanAppOwnerOrganizationId := tfPlan.AppOwnerOrganizationId.ValueString()
		u, _ := uuid.Parse(tfPlanAppOwnerOrganizationId)
		requestBody.SetAppOwnerOrganizationId(&u)
	}

	if !tfPlan.AppRoleAssignmentRequired.Equal(tfState.AppRoleAssignmentRequired) {
		tfPlanAppRoleAssignmentRequired := tfPlan.AppRoleAssignmentRequired.ValueBool()
		requestBody.SetAppRoleAssignmentRequired(&tfPlanAppRoleAssignmentRequired)
	}

	if !tfPlan.AppRoles.Equal(tfState.AppRoles) {
		var tfPlanAppRoles []models.AppRoleable
		for k, i := range tfPlan.AppRoles.Elements() {
			requestBodyAppRoles := models.NewAppRole()
			tfPlanrequestBodyAppRoles := servicePrincipalAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyAppRoles)
			tfStaterequestBodyAppRoles := servicePrincipalAppRoleModel{}
			types.ListValueFrom(ctx, tfState.AppRoles.Elements()[k].Type(ctx), &tfPlanrequestBodyAppRoles)

			if !tfPlanrequestBodyAppRoles.AllowedMemberTypes.Equal(tfStaterequestBodyAppRoles.AllowedMemberTypes) {
				var stringArrayAllowedMemberTypes []string
				for _, i := range tfPlanrequestBodyAppRoles.AllowedMemberTypes.Elements() {
					stringArrayAllowedMemberTypes = append(stringArrayAllowedMemberTypes, i.String())
				}
				requestBodyAppRoles.SetAllowedMemberTypes(stringArrayAllowedMemberTypes)
			}

			if !tfPlanrequestBodyAppRoles.Description.Equal(tfStaterequestBodyAppRoles.Description) {
				tfPlanDescription := tfPlanrequestBodyAppRoles.Description.ValueString()
				requestBodyAppRoles.SetDescription(&tfPlanDescription)
			}

			if !tfPlanrequestBodyAppRoles.DisplayName.Equal(tfStaterequestBodyAppRoles.DisplayName) {
				tfPlanDisplayName := tfPlanrequestBodyAppRoles.DisplayName.ValueString()
				requestBodyAppRoles.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanrequestBodyAppRoles.Id.Equal(tfStaterequestBodyAppRoles.Id) {
				tfPlanId := tfPlanrequestBodyAppRoles.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyAppRoles.SetId(&u)
			}

			if !tfPlanrequestBodyAppRoles.IsEnabled.Equal(tfStaterequestBodyAppRoles.IsEnabled) {
				tfPlanIsEnabled := tfPlanrequestBodyAppRoles.IsEnabled.ValueBool()
				requestBodyAppRoles.SetIsEnabled(&tfPlanIsEnabled)
			}

			if !tfPlanrequestBodyAppRoles.Origin.Equal(tfStaterequestBodyAppRoles.Origin) {
				tfPlanOrigin := tfPlanrequestBodyAppRoles.Origin.ValueString()
				requestBodyAppRoles.SetOrigin(&tfPlanOrigin)
			}

			if !tfPlanrequestBodyAppRoles.Value.Equal(tfStaterequestBodyAppRoles.Value) {
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

	if !tfPlan.CustomSecurityAttributes.Equal(tfState.CustomSecurityAttributes) {
		requestBodyCustomSecurityAttributes := models.NewCustomSecurityAttributeValue()
		tfPlanrequestBodyCustomSecurityAttributes := servicePrincipalCustomSecurityAttributeValueModel{}
		tfPlan.CustomSecurityAttributes.As(ctx, &tfPlanrequestBodyCustomSecurityAttributes, basetypes.ObjectAsOptions{})
		tfStaterequestBodyCustomSecurityAttributes := servicePrincipalCustomSecurityAttributeValueModel{}
		tfState.CustomSecurityAttributes.As(ctx, &tfStaterequestBodyCustomSecurityAttributes, basetypes.ObjectAsOptions{})

		requestBody.SetCustomSecurityAttributes(requestBodyCustomSecurityAttributes)
		tfPlan.CustomSecurityAttributes, _ = types.ObjectValueFrom(ctx, tfPlanrequestBodyCustomSecurityAttributes.AttributeTypes(), tfPlanrequestBodyCustomSecurityAttributes)
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

	if !tfPlan.Homepage.Equal(tfState.Homepage) {
		tfPlanHomepage := tfPlan.Homepage.ValueString()
		requestBody.SetHomepage(&tfPlanHomepage)
	}

	if !tfPlan.Info.Equal(tfState.Info) {
		requestBodyInfo := models.NewInformationalUrl()
		tfPlanrequestBodyInfo := servicePrincipalInformationalUrlModel{}
		tfPlan.Info.As(ctx, &tfPlanrequestBodyInfo, basetypes.ObjectAsOptions{})
		tfStaterequestBodyInfo := servicePrincipalInformationalUrlModel{}
		tfState.Info.As(ctx, &tfStaterequestBodyInfo, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyInfo.LogoUrl.Equal(tfStaterequestBodyInfo.LogoUrl) {
			tfPlanLogoUrl := tfPlanrequestBodyInfo.LogoUrl.ValueString()
			requestBodyInfo.SetLogoUrl(&tfPlanLogoUrl)
		}

		if !tfPlanrequestBodyInfo.MarketingUrl.Equal(tfStaterequestBodyInfo.MarketingUrl) {
			tfPlanMarketingUrl := tfPlanrequestBodyInfo.MarketingUrl.ValueString()
			requestBodyInfo.SetMarketingUrl(&tfPlanMarketingUrl)
		}

		if !tfPlanrequestBodyInfo.PrivacyStatementUrl.Equal(tfStaterequestBodyInfo.PrivacyStatementUrl) {
			tfPlanPrivacyStatementUrl := tfPlanrequestBodyInfo.PrivacyStatementUrl.ValueString()
			requestBodyInfo.SetPrivacyStatementUrl(&tfPlanPrivacyStatementUrl)
		}

		if !tfPlanrequestBodyInfo.SupportUrl.Equal(tfStaterequestBodyInfo.SupportUrl) {
			tfPlanSupportUrl := tfPlanrequestBodyInfo.SupportUrl.ValueString()
			requestBodyInfo.SetSupportUrl(&tfPlanSupportUrl)
		}

		if !tfPlanrequestBodyInfo.TermsOfServiceUrl.Equal(tfStaterequestBodyInfo.TermsOfServiceUrl) {
			tfPlanTermsOfServiceUrl := tfPlanrequestBodyInfo.TermsOfServiceUrl.ValueString()
			requestBodyInfo.SetTermsOfServiceUrl(&tfPlanTermsOfServiceUrl)
		}
		requestBody.SetInfo(requestBodyInfo)
		tfPlan.Info, _ = types.ObjectValueFrom(ctx, tfPlanrequestBodyInfo.AttributeTypes(), tfPlanrequestBodyInfo)
	}

	if !tfPlan.KeyCredentials.Equal(tfState.KeyCredentials) {
		var tfPlanKeyCredentials []models.KeyCredentialable
		for k, i := range tfPlan.KeyCredentials.Elements() {
			requestBodyKeyCredentials := models.NewKeyCredential()
			tfPlanrequestBodyKeyCredentials := servicePrincipalKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyKeyCredentials)
			tfStaterequestBodyKeyCredentials := servicePrincipalKeyCredentialModel{}
			types.ListValueFrom(ctx, tfState.KeyCredentials.Elements()[k].Type(ctx), &tfPlanrequestBodyKeyCredentials)

			if !tfPlanrequestBodyKeyCredentials.CustomKeyIdentifier.Equal(tfStaterequestBodyKeyCredentials.CustomKeyIdentifier) {
				tfPlanCustomKeyIdentifier := tfPlanrequestBodyKeyCredentials.CustomKeyIdentifier.ValueString()
				requestBodyKeyCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			}

			if !tfPlanrequestBodyKeyCredentials.DisplayName.Equal(tfStaterequestBodyKeyCredentials.DisplayName) {
				tfPlanDisplayName := tfPlanrequestBodyKeyCredentials.DisplayName.ValueString()
				requestBodyKeyCredentials.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanrequestBodyKeyCredentials.EndDateTime.Equal(tfStaterequestBodyKeyCredentials.EndDateTime) {
				tfPlanEndDateTime := tfPlanrequestBodyKeyCredentials.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyKeyCredentials.SetEndDateTime(&t)
			}

			if !tfPlanrequestBodyKeyCredentials.Key.Equal(tfStaterequestBodyKeyCredentials.Key) {
				tfPlanKey := tfPlanrequestBodyKeyCredentials.Key.ValueString()
				requestBodyKeyCredentials.SetKey([]byte(tfPlanKey))
			}

			if !tfPlanrequestBodyKeyCredentials.KeyId.Equal(tfStaterequestBodyKeyCredentials.KeyId) {
				tfPlanKeyId := tfPlanrequestBodyKeyCredentials.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyKeyCredentials.SetKeyId(&u)
			}

			if !tfPlanrequestBodyKeyCredentials.StartDateTime.Equal(tfStaterequestBodyKeyCredentials.StartDateTime) {
				tfPlanStartDateTime := tfPlanrequestBodyKeyCredentials.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyKeyCredentials.SetStartDateTime(&t)
			}

			if !tfPlanrequestBodyKeyCredentials.Type.Equal(tfStaterequestBodyKeyCredentials.Type) {
				tfPlanType := tfPlanrequestBodyKeyCredentials.Type.ValueString()
				requestBodyKeyCredentials.SetTypeEscaped(&tfPlanType)
			}

			if !tfPlanrequestBodyKeyCredentials.Usage.Equal(tfStaterequestBodyKeyCredentials.Usage) {
				tfPlanUsage := tfPlanrequestBodyKeyCredentials.Usage.ValueString()
				requestBodyKeyCredentials.SetUsage(&tfPlanUsage)
			}
		}
		requestBody.SetKeyCredentials(tfPlanKeyCredentials)
	}

	if !tfPlan.LoginUrl.Equal(tfState.LoginUrl) {
		tfPlanLoginUrl := tfPlan.LoginUrl.ValueString()
		requestBody.SetLoginUrl(&tfPlanLoginUrl)
	}

	if !tfPlan.LogoutUrl.Equal(tfState.LogoutUrl) {
		tfPlanLogoutUrl := tfPlan.LogoutUrl.ValueString()
		requestBody.SetLogoutUrl(&tfPlanLogoutUrl)
	}

	if !tfPlan.Notes.Equal(tfState.Notes) {
		tfPlanNotes := tfPlan.Notes.ValueString()
		requestBody.SetNotes(&tfPlanNotes)
	}

	if !tfPlan.NotificationEmailAddresses.Equal(tfState.NotificationEmailAddresses) {
		var stringArrayNotificationEmailAddresses []string
		for _, i := range tfPlan.NotificationEmailAddresses.Elements() {
			stringArrayNotificationEmailAddresses = append(stringArrayNotificationEmailAddresses, i.String())
		}
		requestBody.SetNotificationEmailAddresses(stringArrayNotificationEmailAddresses)
	}

	if !tfPlan.Oauth2PermissionScopes.Equal(tfState.Oauth2PermissionScopes) {
		var tfPlanOauth2PermissionScopes []models.PermissionScopeable
		for k, i := range tfPlan.Oauth2PermissionScopes.Elements() {
			requestBodyOauth2PermissionScopes := models.NewPermissionScope()
			tfPlanrequestBodyOauth2PermissionScopes := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyOauth2PermissionScopes)
			tfStaterequestBodyOauth2PermissionScopes := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, tfState.Oauth2PermissionScopes.Elements()[k].Type(ctx), &tfPlanrequestBodyOauth2PermissionScopes)

			if !tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDescription.Equal(tfStaterequestBodyOauth2PermissionScopes.AdminConsentDescription) {
				tfPlanAdminConsentDescription := tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDescription.ValueString()
				requestBodyOauth2PermissionScopes.SetAdminConsentDescription(&tfPlanAdminConsentDescription)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDisplayName.Equal(tfStaterequestBodyOauth2PermissionScopes.AdminConsentDisplayName) {
				tfPlanAdminConsentDisplayName := tfPlanrequestBodyOauth2PermissionScopes.AdminConsentDisplayName.ValueString()
				requestBodyOauth2PermissionScopes.SetAdminConsentDisplayName(&tfPlanAdminConsentDisplayName)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.Id.Equal(tfStaterequestBodyOauth2PermissionScopes.Id) {
				tfPlanId := tfPlanrequestBodyOauth2PermissionScopes.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyOauth2PermissionScopes.SetId(&u)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.IsEnabled.Equal(tfStaterequestBodyOauth2PermissionScopes.IsEnabled) {
				tfPlanIsEnabled := tfPlanrequestBodyOauth2PermissionScopes.IsEnabled.ValueBool()
				requestBodyOauth2PermissionScopes.SetIsEnabled(&tfPlanIsEnabled)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.Origin.Equal(tfStaterequestBodyOauth2PermissionScopes.Origin) {
				tfPlanOrigin := tfPlanrequestBodyOauth2PermissionScopes.Origin.ValueString()
				requestBodyOauth2PermissionScopes.SetOrigin(&tfPlanOrigin)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.Type.Equal(tfStaterequestBodyOauth2PermissionScopes.Type) {
				tfPlanType := tfPlanrequestBodyOauth2PermissionScopes.Type.ValueString()
				requestBodyOauth2PermissionScopes.SetTypeEscaped(&tfPlanType)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.UserConsentDescription.Equal(tfStaterequestBodyOauth2PermissionScopes.UserConsentDescription) {
				tfPlanUserConsentDescription := tfPlanrequestBodyOauth2PermissionScopes.UserConsentDescription.ValueString()
				requestBodyOauth2PermissionScopes.SetUserConsentDescription(&tfPlanUserConsentDescription)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.UserConsentDisplayName.Equal(tfStaterequestBodyOauth2PermissionScopes.UserConsentDisplayName) {
				tfPlanUserConsentDisplayName := tfPlanrequestBodyOauth2PermissionScopes.UserConsentDisplayName.ValueString()
				requestBodyOauth2PermissionScopes.SetUserConsentDisplayName(&tfPlanUserConsentDisplayName)
			}

			if !tfPlanrequestBodyOauth2PermissionScopes.Value.Equal(tfStaterequestBodyOauth2PermissionScopes.Value) {
				tfPlanValue := tfPlanrequestBodyOauth2PermissionScopes.Value.ValueString()
				requestBodyOauth2PermissionScopes.SetValue(&tfPlanValue)
			}
		}
		requestBody.SetOauth2PermissionScopes(tfPlanOauth2PermissionScopes)
	}

	if !tfPlan.PasswordCredentials.Equal(tfState.PasswordCredentials) {
		var tfPlanPasswordCredentials []models.PasswordCredentialable
		for k, i := range tfPlan.PasswordCredentials.Elements() {
			requestBodyPasswordCredentials := models.NewPasswordCredential()
			tfPlanrequestBodyPasswordCredentials := servicePrincipalPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyPasswordCredentials)
			tfStaterequestBodyPasswordCredentials := servicePrincipalPasswordCredentialModel{}
			types.ListValueFrom(ctx, tfState.PasswordCredentials.Elements()[k].Type(ctx), &tfPlanrequestBodyPasswordCredentials)

			if !tfPlanrequestBodyPasswordCredentials.CustomKeyIdentifier.Equal(tfStaterequestBodyPasswordCredentials.CustomKeyIdentifier) {
				tfPlanCustomKeyIdentifier := tfPlanrequestBodyPasswordCredentials.CustomKeyIdentifier.ValueString()
				requestBodyPasswordCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			}

			if !tfPlanrequestBodyPasswordCredentials.DisplayName.Equal(tfStaterequestBodyPasswordCredentials.DisplayName) {
				tfPlanDisplayName := tfPlanrequestBodyPasswordCredentials.DisplayName.ValueString()
				requestBodyPasswordCredentials.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanrequestBodyPasswordCredentials.EndDateTime.Equal(tfStaterequestBodyPasswordCredentials.EndDateTime) {
				tfPlanEndDateTime := tfPlanrequestBodyPasswordCredentials.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				requestBodyPasswordCredentials.SetEndDateTime(&t)
			}

			if !tfPlanrequestBodyPasswordCredentials.Hint.Equal(tfStaterequestBodyPasswordCredentials.Hint) {
				tfPlanHint := tfPlanrequestBodyPasswordCredentials.Hint.ValueString()
				requestBodyPasswordCredentials.SetHint(&tfPlanHint)
			}

			if !tfPlanrequestBodyPasswordCredentials.KeyId.Equal(tfStaterequestBodyPasswordCredentials.KeyId) {
				tfPlanKeyId := tfPlanrequestBodyPasswordCredentials.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				requestBodyPasswordCredentials.SetKeyId(&u)
			}

			if !tfPlanrequestBodyPasswordCredentials.SecretText.Equal(tfStaterequestBodyPasswordCredentials.SecretText) {
				tfPlanSecretText := tfPlanrequestBodyPasswordCredentials.SecretText.ValueString()
				requestBodyPasswordCredentials.SetSecretText(&tfPlanSecretText)
			}

			if !tfPlanrequestBodyPasswordCredentials.StartDateTime.Equal(tfStaterequestBodyPasswordCredentials.StartDateTime) {
				tfPlanStartDateTime := tfPlanrequestBodyPasswordCredentials.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				requestBodyPasswordCredentials.SetStartDateTime(&t)
			}
		}
		requestBody.SetPasswordCredentials(tfPlanPasswordCredentials)
	}

	if !tfPlan.PreferredSingleSignOnMode.Equal(tfState.PreferredSingleSignOnMode) {
		tfPlanPreferredSingleSignOnMode := tfPlan.PreferredSingleSignOnMode.ValueString()
		requestBody.SetPreferredSingleSignOnMode(&tfPlanPreferredSingleSignOnMode)
	}

	if !tfPlan.PreferredTokenSigningKeyThumbprint.Equal(tfState.PreferredTokenSigningKeyThumbprint) {
		tfPlanPreferredTokenSigningKeyThumbprint := tfPlan.PreferredTokenSigningKeyThumbprint.ValueString()
		requestBody.SetPreferredTokenSigningKeyThumbprint(&tfPlanPreferredTokenSigningKeyThumbprint)
	}

	if !tfPlan.ReplyUrls.Equal(tfState.ReplyUrls) {
		var stringArrayReplyUrls []string
		for _, i := range tfPlan.ReplyUrls.Elements() {
			stringArrayReplyUrls = append(stringArrayReplyUrls, i.String())
		}
		requestBody.SetReplyUrls(stringArrayReplyUrls)
	}

	if !tfPlan.ResourceSpecificApplicationPermissions.Equal(tfState.ResourceSpecificApplicationPermissions) {
		var tfPlanResourceSpecificApplicationPermissions []models.ResourceSpecificPermissionable
		for k, i := range tfPlan.ResourceSpecificApplicationPermissions.Elements() {
			requestBodyResourceSpecificApplicationPermissions := models.NewResourceSpecificPermission()
			tfPlanrequestBodyResourceSpecificApplicationPermissions := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanrequestBodyResourceSpecificApplicationPermissions)
			tfStaterequestBodyResourceSpecificApplicationPermissions := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, tfState.ResourceSpecificApplicationPermissions.Elements()[k].Type(ctx), &tfPlanrequestBodyResourceSpecificApplicationPermissions)

			if !tfPlanrequestBodyResourceSpecificApplicationPermissions.Description.Equal(tfStaterequestBodyResourceSpecificApplicationPermissions.Description) {
				tfPlanDescription := tfPlanrequestBodyResourceSpecificApplicationPermissions.Description.ValueString()
				requestBodyResourceSpecificApplicationPermissions.SetDescription(&tfPlanDescription)
			}

			if !tfPlanrequestBodyResourceSpecificApplicationPermissions.DisplayName.Equal(tfStaterequestBodyResourceSpecificApplicationPermissions.DisplayName) {
				tfPlanDisplayName := tfPlanrequestBodyResourceSpecificApplicationPermissions.DisplayName.ValueString()
				requestBodyResourceSpecificApplicationPermissions.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanrequestBodyResourceSpecificApplicationPermissions.Id.Equal(tfStaterequestBodyResourceSpecificApplicationPermissions.Id) {
				tfPlanId := tfPlanrequestBodyResourceSpecificApplicationPermissions.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				requestBodyResourceSpecificApplicationPermissions.SetId(&u)
			}

			if !tfPlanrequestBodyResourceSpecificApplicationPermissions.IsEnabled.Equal(tfStaterequestBodyResourceSpecificApplicationPermissions.IsEnabled) {
				tfPlanIsEnabled := tfPlanrequestBodyResourceSpecificApplicationPermissions.IsEnabled.ValueBool()
				requestBodyResourceSpecificApplicationPermissions.SetIsEnabled(&tfPlanIsEnabled)
			}

			if !tfPlanrequestBodyResourceSpecificApplicationPermissions.Value.Equal(tfStaterequestBodyResourceSpecificApplicationPermissions.Value) {
				tfPlanValue := tfPlanrequestBodyResourceSpecificApplicationPermissions.Value.ValueString()
				requestBodyResourceSpecificApplicationPermissions.SetValue(&tfPlanValue)
			}
		}
		requestBody.SetResourceSpecificApplicationPermissions(tfPlanResourceSpecificApplicationPermissions)
	}

	if !tfPlan.SamlSingleSignOnSettings.Equal(tfState.SamlSingleSignOnSettings) {
		requestBodySamlSingleSignOnSettings := models.NewSamlSingleSignOnSettings()
		tfPlanrequestBodySamlSingleSignOnSettings := servicePrincipalSamlSingleSignOnSettingsModel{}
		tfPlan.SamlSingleSignOnSettings.As(ctx, &tfPlanrequestBodySamlSingleSignOnSettings, basetypes.ObjectAsOptions{})
		tfStaterequestBodySamlSingleSignOnSettings := servicePrincipalSamlSingleSignOnSettingsModel{}
		tfState.SamlSingleSignOnSettings.As(ctx, &tfStaterequestBodySamlSingleSignOnSettings, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodySamlSingleSignOnSettings.RelayState.Equal(tfStaterequestBodySamlSingleSignOnSettings.RelayState) {
			tfPlanRelayState := tfPlanrequestBodySamlSingleSignOnSettings.RelayState.ValueString()
			requestBodySamlSingleSignOnSettings.SetRelayState(&tfPlanRelayState)
		}
		requestBody.SetSamlSingleSignOnSettings(requestBodySamlSingleSignOnSettings)
		tfPlan.SamlSingleSignOnSettings, _ = types.ObjectValueFrom(ctx, tfPlanrequestBodySamlSingleSignOnSettings.AttributeTypes(), tfPlanrequestBodySamlSingleSignOnSettings)
	}

	if !tfPlan.ServicePrincipalNames.Equal(tfState.ServicePrincipalNames) {
		var stringArrayServicePrincipalNames []string
		for _, i := range tfPlan.ServicePrincipalNames.Elements() {
			stringArrayServicePrincipalNames = append(stringArrayServicePrincipalNames, i.String())
		}
		requestBody.SetServicePrincipalNames(stringArrayServicePrincipalNames)
	}

	if !tfPlan.ServicePrincipalType.Equal(tfState.ServicePrincipalType) {
		tfPlanServicePrincipalType := tfPlan.ServicePrincipalType.ValueString()
		requestBody.SetServicePrincipalType(&tfPlanServicePrincipalType)
	}

	if !tfPlan.SignInAudience.Equal(tfState.SignInAudience) {
		tfPlanSignInAudience := tfPlan.SignInAudience.ValueString()
		requestBody.SetSignInAudience(&tfPlanSignInAudience)
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

	if !tfPlan.VerifiedPublisher.Equal(tfState.VerifiedPublisher) {
		requestBodyVerifiedPublisher := models.NewVerifiedPublisher()
		tfPlanrequestBodyVerifiedPublisher := servicePrincipalVerifiedPublisherModel{}
		tfPlan.VerifiedPublisher.As(ctx, &tfPlanrequestBodyVerifiedPublisher, basetypes.ObjectAsOptions{})
		tfStaterequestBodyVerifiedPublisher := servicePrincipalVerifiedPublisherModel{}
		tfState.VerifiedPublisher.As(ctx, &tfStaterequestBodyVerifiedPublisher, basetypes.ObjectAsOptions{})

		if !tfPlanrequestBodyVerifiedPublisher.AddedDateTime.Equal(tfStaterequestBodyVerifiedPublisher.AddedDateTime) {
			tfPlanAddedDateTime := tfPlanrequestBodyVerifiedPublisher.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanAddedDateTime)
			requestBodyVerifiedPublisher.SetAddedDateTime(&t)
		}

		if !tfPlanrequestBodyVerifiedPublisher.DisplayName.Equal(tfStaterequestBodyVerifiedPublisher.DisplayName) {
			tfPlanDisplayName := tfPlanrequestBodyVerifiedPublisher.DisplayName.ValueString()
			requestBodyVerifiedPublisher.SetDisplayName(&tfPlanDisplayName)
		}

		if !tfPlanrequestBodyVerifiedPublisher.VerifiedPublisherId.Equal(tfStaterequestBodyVerifiedPublisher.VerifiedPublisherId) {
			tfPlanVerifiedPublisherId := tfPlanrequestBodyVerifiedPublisher.VerifiedPublisherId.ValueString()
			requestBodyVerifiedPublisher.SetVerifiedPublisherId(&tfPlanVerifiedPublisherId)
		}
		requestBody.SetVerifiedPublisher(requestBodyVerifiedPublisher)
		tfPlan.VerifiedPublisher, _ = types.ObjectValueFrom(ctx, tfPlanrequestBodyVerifiedPublisher.AttributeTypes(), tfPlanrequestBodyVerifiedPublisher)
	}

	// Update servicePrincipal
	_, err := r.client.ServicePrincipals().ByServicePrincipalId(tfState.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating service_principal",
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
func (r *servicePrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfState servicePrincipalModel
	diags := req.State.Get(ctx, &tfState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete servicePrincipal
	err := r.client.ServicePrincipals().ByServicePrincipalId(tfState.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting service_principal",
			err.Error(),
		)
		return
	}

}
