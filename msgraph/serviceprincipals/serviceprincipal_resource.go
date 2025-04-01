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
	var tfPlan servicePrincipalModel
	diags := req.Plan.Get(ctx, &tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Plan
	requestBody := models.NewServicePrincipal()

	if !tfPlan.Id.IsUnknown() {
		tfPlanId := tfPlan.Id.ValueString()
		requestBody.SetId(&tfPlanId)
	} else {
		tfPlan.Id = types.StringNull()
	}

	if !tfPlan.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	} else {
		tfPlan.DeletedDateTime = types.StringNull()
	}

	if !tfPlan.AccountEnabled.IsUnknown() {
		tfPlanAccountEnabled := tfPlan.AccountEnabled.ValueBool()
		requestBody.SetAccountEnabled(&tfPlanAccountEnabled)
	} else {
		tfPlan.AccountEnabled = types.BoolNull()
	}

	if len(tfPlan.AddIns.Elements()) > 0 {
		var tfPlanAddIns []models.AddInable
		for _, i := range tfPlan.AddIns.Elements() {
			addIns := models.NewAddIn()
			addInsModel := servicePrincipalAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &addInsModel)

			if !addInsModel.Id.IsUnknown() {
				tfPlanId := addInsModel.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				addIns.SetId(&u)
			} else {
				addInsModel.Id = types.StringNull()
			}

			if len(addInsModel.Properties.Elements()) > 0 {
				var tfPlanProperties []models.KeyValueable
				for _, i := range addInsModel.Properties.Elements() {
					properties := models.NewKeyValue()
					propertiesModel := servicePrincipalKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &propertiesModel)

					if !propertiesModel.Key.IsUnknown() {
						tfPlanKey := propertiesModel.Key.ValueString()
						properties.SetKey(&tfPlanKey)
					} else {
						propertiesModel.Key = types.StringNull()
					}

					if !propertiesModel.Value.IsUnknown() {
						tfPlanValue := propertiesModel.Value.ValueString()
						properties.SetValue(&tfPlanValue)
					} else {
						propertiesModel.Value = types.StringNull()
					}
				}
				addIns.SetProperties(tfPlanProperties)
			} else {
				addInsModel.Properties = types.ListNull(addInsModel.Properties.ElementType(ctx))
			}

			if !addInsModel.Type.IsUnknown() {
				tfPlanType := addInsModel.Type.ValueString()
				addIns.SetTypeEscaped(&tfPlanType)
			} else {
				addInsModel.Type = types.StringNull()
			}
		}
		requestBody.SetAddIns(tfPlanAddIns)
	} else {
		tfPlan.AddIns = types.ListNull(tfPlan.AddIns.ElementType(ctx))
	}

	if len(tfPlan.AlternativeNames.Elements()) > 0 {
		var alternativeNames []string
		for _, i := range tfPlan.AlternativeNames.Elements() {
			alternativeNames = append(alternativeNames, i.String())
		}
		requestBody.SetAlternativeNames(alternativeNames)
	} else {
		tfPlan.AlternativeNames = types.ListNull(types.StringType)
	}

	if !tfPlan.AppDescription.IsUnknown() {
		tfPlanAppDescription := tfPlan.AppDescription.ValueString()
		requestBody.SetAppDescription(&tfPlanAppDescription)
	} else {
		tfPlan.AppDescription = types.StringNull()
	}

	if !tfPlan.AppDisplayName.IsUnknown() {
		tfPlanAppDisplayName := tfPlan.AppDisplayName.ValueString()
		requestBody.SetAppDisplayName(&tfPlanAppDisplayName)
	} else {
		tfPlan.AppDisplayName = types.StringNull()
	}

	if !tfPlan.AppId.IsUnknown() {
		tfPlanAppId := tfPlan.AppId.ValueString()
		requestBody.SetAppId(&tfPlanAppId)
	} else {
		tfPlan.AppId = types.StringNull()
	}

	if !tfPlan.AppOwnerOrganizationId.IsUnknown() {
		tfPlanAppOwnerOrganizationId := tfPlan.AppOwnerOrganizationId.ValueString()
		u, _ := uuid.Parse(tfPlanAppOwnerOrganizationId)
		requestBody.SetAppOwnerOrganizationId(&u)
	} else {
		tfPlan.AppOwnerOrganizationId = types.StringNull()
	}

	if !tfPlan.AppRoleAssignmentRequired.IsUnknown() {
		tfPlanAppRoleAssignmentRequired := tfPlan.AppRoleAssignmentRequired.ValueBool()
		requestBody.SetAppRoleAssignmentRequired(&tfPlanAppRoleAssignmentRequired)
	} else {
		tfPlan.AppRoleAssignmentRequired = types.BoolNull()
	}

	if len(tfPlan.AppRoles.Elements()) > 0 {
		var tfPlanAppRoles []models.AppRoleable
		for _, i := range tfPlan.AppRoles.Elements() {
			appRoles := models.NewAppRole()
			appRolesModel := servicePrincipalAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &appRolesModel)

			if len(appRolesModel.AllowedMemberTypes.Elements()) > 0 {
				var allowedMemberTypes []string
				for _, i := range appRolesModel.AllowedMemberTypes.Elements() {
					allowedMemberTypes = append(allowedMemberTypes, i.String())
				}
				appRoles.SetAllowedMemberTypes(allowedMemberTypes)
			} else {
				appRolesModel.AllowedMemberTypes = types.ListNull(types.StringType)
			}

			if !appRolesModel.Description.IsUnknown() {
				tfPlanDescription := appRolesModel.Description.ValueString()
				appRoles.SetDescription(&tfPlanDescription)
			} else {
				appRolesModel.Description = types.StringNull()
			}

			if !appRolesModel.DisplayName.IsUnknown() {
				tfPlanDisplayName := appRolesModel.DisplayName.ValueString()
				appRoles.SetDisplayName(&tfPlanDisplayName)
			} else {
				appRolesModel.DisplayName = types.StringNull()
			}

			if !appRolesModel.Id.IsUnknown() {
				tfPlanId := appRolesModel.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				appRoles.SetId(&u)
			} else {
				appRolesModel.Id = types.StringNull()
			}

			if !appRolesModel.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := appRolesModel.IsEnabled.ValueBool()
				appRoles.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				appRolesModel.IsEnabled = types.BoolNull()
			}

			if !appRolesModel.Origin.IsUnknown() {
				tfPlanOrigin := appRolesModel.Origin.ValueString()
				appRoles.SetOrigin(&tfPlanOrigin)
			} else {
				appRolesModel.Origin = types.StringNull()
			}

			if !appRolesModel.Value.IsUnknown() {
				tfPlanValue := appRolesModel.Value.ValueString()
				appRoles.SetValue(&tfPlanValue)
			} else {
				appRolesModel.Value = types.StringNull()
			}
		}
		requestBody.SetAppRoles(tfPlanAppRoles)
	} else {
		tfPlan.AppRoles = types.ListNull(tfPlan.AppRoles.ElementType(ctx))
	}

	if !tfPlan.ApplicationTemplateId.IsUnknown() {
		tfPlanApplicationTemplateId := tfPlan.ApplicationTemplateId.ValueString()
		requestBody.SetApplicationTemplateId(&tfPlanApplicationTemplateId)
	} else {
		tfPlan.ApplicationTemplateId = types.StringNull()
	}

	if !tfPlan.CustomSecurityAttributes.IsUnknown() {
		customSecurityAttributes := models.NewCustomSecurityAttributeValue()
		customSecurityAttributesModel := servicePrincipalCustomSecurityAttributeValueModel{}
		tfPlan.CustomSecurityAttributes.As(ctx, &customSecurityAttributesModel, basetypes.ObjectAsOptions{})

		requestBody.SetCustomSecurityAttributes(customSecurityAttributes)
		objectValue, _ := types.ObjectValueFrom(ctx, customSecurityAttributesModel.AttributeTypes(), customSecurityAttributesModel)
		tfPlan.CustomSecurityAttributes = objectValue
	} else {
		tfPlan.CustomSecurityAttributes = types.ObjectNull(tfPlan.CustomSecurityAttributes.AttributeTypes(ctx))
	}

	if !tfPlan.Description.IsUnknown() {
		tfPlanDescription := tfPlan.Description.ValueString()
		requestBody.SetDescription(&tfPlanDescription)
	} else {
		tfPlan.Description = types.StringNull()
	}

	if !tfPlan.DisabledByMicrosoftStatus.IsUnknown() {
		tfPlanDisabledByMicrosoftStatus := tfPlan.DisabledByMicrosoftStatus.ValueString()
		requestBody.SetDisabledByMicrosoftStatus(&tfPlanDisabledByMicrosoftStatus)
	} else {
		tfPlan.DisabledByMicrosoftStatus = types.StringNull()
	}

	if !tfPlan.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlan.DisplayName.ValueString()
		requestBody.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlan.DisplayName = types.StringNull()
	}

	if !tfPlan.Homepage.IsUnknown() {
		tfPlanHomepage := tfPlan.Homepage.ValueString()
		requestBody.SetHomepage(&tfPlanHomepage)
	} else {
		tfPlan.Homepage = types.StringNull()
	}

	if !tfPlan.Info.IsUnknown() {
		info := models.NewInformationalUrl()
		infoModel := servicePrincipalInformationalUrlModel{}
		tfPlan.Info.As(ctx, &infoModel, basetypes.ObjectAsOptions{})

		if !infoModel.LogoUrl.IsUnknown() {
			tfPlanLogoUrl := infoModel.LogoUrl.ValueString()
			info.SetLogoUrl(&tfPlanLogoUrl)
		} else {
			infoModel.LogoUrl = types.StringNull()
		}

		if !infoModel.MarketingUrl.IsUnknown() {
			tfPlanMarketingUrl := infoModel.MarketingUrl.ValueString()
			info.SetMarketingUrl(&tfPlanMarketingUrl)
		} else {
			infoModel.MarketingUrl = types.StringNull()
		}

		if !infoModel.PrivacyStatementUrl.IsUnknown() {
			tfPlanPrivacyStatementUrl := infoModel.PrivacyStatementUrl.ValueString()
			info.SetPrivacyStatementUrl(&tfPlanPrivacyStatementUrl)
		} else {
			infoModel.PrivacyStatementUrl = types.StringNull()
		}

		if !infoModel.SupportUrl.IsUnknown() {
			tfPlanSupportUrl := infoModel.SupportUrl.ValueString()
			info.SetSupportUrl(&tfPlanSupportUrl)
		} else {
			infoModel.SupportUrl = types.StringNull()
		}

		if !infoModel.TermsOfServiceUrl.IsUnknown() {
			tfPlanTermsOfServiceUrl := infoModel.TermsOfServiceUrl.ValueString()
			info.SetTermsOfServiceUrl(&tfPlanTermsOfServiceUrl)
		} else {
			infoModel.TermsOfServiceUrl = types.StringNull()
		}
		requestBody.SetInfo(info)
		objectValue, _ := types.ObjectValueFrom(ctx, infoModel.AttributeTypes(), infoModel)
		tfPlan.Info = objectValue
	} else {
		tfPlan.Info = types.ObjectNull(tfPlan.Info.AttributeTypes(ctx))
	}

	if len(tfPlan.KeyCredentials.Elements()) > 0 {
		var tfPlanKeyCredentials []models.KeyCredentialable
		for _, i := range tfPlan.KeyCredentials.Elements() {
			keyCredentials := models.NewKeyCredential()
			keyCredentialsModel := servicePrincipalKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &keyCredentialsModel)

			if !keyCredentialsModel.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := keyCredentialsModel.CustomKeyIdentifier.ValueString()
				keyCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				keyCredentialsModel.CustomKeyIdentifier = types.StringNull()
			}

			if !keyCredentialsModel.DisplayName.IsUnknown() {
				tfPlanDisplayName := keyCredentialsModel.DisplayName.ValueString()
				keyCredentials.SetDisplayName(&tfPlanDisplayName)
			} else {
				keyCredentialsModel.DisplayName = types.StringNull()
			}

			if !keyCredentialsModel.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := keyCredentialsModel.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				keyCredentials.SetEndDateTime(&t)
			} else {
				keyCredentialsModel.EndDateTime = types.StringNull()
			}

			if !keyCredentialsModel.Key.IsUnknown() {
				tfPlanKey := keyCredentialsModel.Key.ValueString()
				keyCredentials.SetKey([]byte(tfPlanKey))
			} else {
				keyCredentialsModel.Key = types.StringNull()
			}

			if !keyCredentialsModel.KeyId.IsUnknown() {
				tfPlanKeyId := keyCredentialsModel.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				keyCredentials.SetKeyId(&u)
			} else {
				keyCredentialsModel.KeyId = types.StringNull()
			}

			if !keyCredentialsModel.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := keyCredentialsModel.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				keyCredentials.SetStartDateTime(&t)
			} else {
				keyCredentialsModel.StartDateTime = types.StringNull()
			}

			if !keyCredentialsModel.Type.IsUnknown() {
				tfPlanType := keyCredentialsModel.Type.ValueString()
				keyCredentials.SetTypeEscaped(&tfPlanType)
			} else {
				keyCredentialsModel.Type = types.StringNull()
			}

			if !keyCredentialsModel.Usage.IsUnknown() {
				tfPlanUsage := keyCredentialsModel.Usage.ValueString()
				keyCredentials.SetUsage(&tfPlanUsage)
			} else {
				keyCredentialsModel.Usage = types.StringNull()
			}
		}
		requestBody.SetKeyCredentials(tfPlanKeyCredentials)
	} else {
		tfPlan.KeyCredentials = types.ListNull(tfPlan.KeyCredentials.ElementType(ctx))
	}

	if !tfPlan.LoginUrl.IsUnknown() {
		tfPlanLoginUrl := tfPlan.LoginUrl.ValueString()
		requestBody.SetLoginUrl(&tfPlanLoginUrl)
	} else {
		tfPlan.LoginUrl = types.StringNull()
	}

	if !tfPlan.LogoutUrl.IsUnknown() {
		tfPlanLogoutUrl := tfPlan.LogoutUrl.ValueString()
		requestBody.SetLogoutUrl(&tfPlanLogoutUrl)
	} else {
		tfPlan.LogoutUrl = types.StringNull()
	}

	if !tfPlan.Notes.IsUnknown() {
		tfPlanNotes := tfPlan.Notes.ValueString()
		requestBody.SetNotes(&tfPlanNotes)
	} else {
		tfPlan.Notes = types.StringNull()
	}

	if len(tfPlan.NotificationEmailAddresses.Elements()) > 0 {
		var notificationEmailAddresses []string
		for _, i := range tfPlan.NotificationEmailAddresses.Elements() {
			notificationEmailAddresses = append(notificationEmailAddresses, i.String())
		}
		requestBody.SetNotificationEmailAddresses(notificationEmailAddresses)
	} else {
		tfPlan.NotificationEmailAddresses = types.ListNull(types.StringType)
	}

	if len(tfPlan.Oauth2PermissionScopes.Elements()) > 0 {
		var tfPlanOauth2PermissionScopes []models.PermissionScopeable
		for _, i := range tfPlan.Oauth2PermissionScopes.Elements() {
			oauth2PermissionScopes := models.NewPermissionScope()
			oauth2PermissionScopesModel := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &oauth2PermissionScopesModel)

			if !oauth2PermissionScopesModel.AdminConsentDescription.IsUnknown() {
				tfPlanAdminConsentDescription := oauth2PermissionScopesModel.AdminConsentDescription.ValueString()
				oauth2PermissionScopes.SetAdminConsentDescription(&tfPlanAdminConsentDescription)
			} else {
				oauth2PermissionScopesModel.AdminConsentDescription = types.StringNull()
			}

			if !oauth2PermissionScopesModel.AdminConsentDisplayName.IsUnknown() {
				tfPlanAdminConsentDisplayName := oauth2PermissionScopesModel.AdminConsentDisplayName.ValueString()
				oauth2PermissionScopes.SetAdminConsentDisplayName(&tfPlanAdminConsentDisplayName)
			} else {
				oauth2PermissionScopesModel.AdminConsentDisplayName = types.StringNull()
			}

			if !oauth2PermissionScopesModel.Id.IsUnknown() {
				tfPlanId := oauth2PermissionScopesModel.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				oauth2PermissionScopes.SetId(&u)
			} else {
				oauth2PermissionScopesModel.Id = types.StringNull()
			}

			if !oauth2PermissionScopesModel.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := oauth2PermissionScopesModel.IsEnabled.ValueBool()
				oauth2PermissionScopes.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				oauth2PermissionScopesModel.IsEnabled = types.BoolNull()
			}

			if !oauth2PermissionScopesModel.Origin.IsUnknown() {
				tfPlanOrigin := oauth2PermissionScopesModel.Origin.ValueString()
				oauth2PermissionScopes.SetOrigin(&tfPlanOrigin)
			} else {
				oauth2PermissionScopesModel.Origin = types.StringNull()
			}

			if !oauth2PermissionScopesModel.Type.IsUnknown() {
				tfPlanType := oauth2PermissionScopesModel.Type.ValueString()
				oauth2PermissionScopes.SetTypeEscaped(&tfPlanType)
			} else {
				oauth2PermissionScopesModel.Type = types.StringNull()
			}

			if !oauth2PermissionScopesModel.UserConsentDescription.IsUnknown() {
				tfPlanUserConsentDescription := oauth2PermissionScopesModel.UserConsentDescription.ValueString()
				oauth2PermissionScopes.SetUserConsentDescription(&tfPlanUserConsentDescription)
			} else {
				oauth2PermissionScopesModel.UserConsentDescription = types.StringNull()
			}

			if !oauth2PermissionScopesModel.UserConsentDisplayName.IsUnknown() {
				tfPlanUserConsentDisplayName := oauth2PermissionScopesModel.UserConsentDisplayName.ValueString()
				oauth2PermissionScopes.SetUserConsentDisplayName(&tfPlanUserConsentDisplayName)
			} else {
				oauth2PermissionScopesModel.UserConsentDisplayName = types.StringNull()
			}

			if !oauth2PermissionScopesModel.Value.IsUnknown() {
				tfPlanValue := oauth2PermissionScopesModel.Value.ValueString()
				oauth2PermissionScopes.SetValue(&tfPlanValue)
			} else {
				oauth2PermissionScopesModel.Value = types.StringNull()
			}
		}
		requestBody.SetOauth2PermissionScopes(tfPlanOauth2PermissionScopes)
	} else {
		tfPlan.Oauth2PermissionScopes = types.ListNull(tfPlan.Oauth2PermissionScopes.ElementType(ctx))
	}

	if len(tfPlan.PasswordCredentials.Elements()) > 0 {
		var tfPlanPasswordCredentials []models.PasswordCredentialable
		for _, i := range tfPlan.PasswordCredentials.Elements() {
			passwordCredentials := models.NewPasswordCredential()
			passwordCredentialsModel := servicePrincipalPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &passwordCredentialsModel)

			if !passwordCredentialsModel.CustomKeyIdentifier.IsUnknown() {
				tfPlanCustomKeyIdentifier := passwordCredentialsModel.CustomKeyIdentifier.ValueString()
				passwordCredentials.SetCustomKeyIdentifier([]byte(tfPlanCustomKeyIdentifier))
			} else {
				passwordCredentialsModel.CustomKeyIdentifier = types.StringNull()
			}

			if !passwordCredentialsModel.DisplayName.IsUnknown() {
				tfPlanDisplayName := passwordCredentialsModel.DisplayName.ValueString()
				passwordCredentials.SetDisplayName(&tfPlanDisplayName)
			} else {
				passwordCredentialsModel.DisplayName = types.StringNull()
			}

			if !passwordCredentialsModel.EndDateTime.IsUnknown() {
				tfPlanEndDateTime := passwordCredentialsModel.EndDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanEndDateTime)
				passwordCredentials.SetEndDateTime(&t)
			} else {
				passwordCredentialsModel.EndDateTime = types.StringNull()
			}

			if !passwordCredentialsModel.Hint.IsUnknown() {
				tfPlanHint := passwordCredentialsModel.Hint.ValueString()
				passwordCredentials.SetHint(&tfPlanHint)
			} else {
				passwordCredentialsModel.Hint = types.StringNull()
			}

			if !passwordCredentialsModel.KeyId.IsUnknown() {
				tfPlanKeyId := passwordCredentialsModel.KeyId.ValueString()
				u, _ := uuid.Parse(tfPlanKeyId)
				passwordCredentials.SetKeyId(&u)
			} else {
				passwordCredentialsModel.KeyId = types.StringNull()
			}

			if !passwordCredentialsModel.SecretText.IsUnknown() {
				tfPlanSecretText := passwordCredentialsModel.SecretText.ValueString()
				passwordCredentials.SetSecretText(&tfPlanSecretText)
			} else {
				passwordCredentialsModel.SecretText = types.StringNull()
			}

			if !passwordCredentialsModel.StartDateTime.IsUnknown() {
				tfPlanStartDateTime := passwordCredentialsModel.StartDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanStartDateTime)
				passwordCredentials.SetStartDateTime(&t)
			} else {
				passwordCredentialsModel.StartDateTime = types.StringNull()
			}
		}
		requestBody.SetPasswordCredentials(tfPlanPasswordCredentials)
	} else {
		tfPlan.PasswordCredentials = types.ListNull(tfPlan.PasswordCredentials.ElementType(ctx))
	}

	if !tfPlan.PreferredSingleSignOnMode.IsUnknown() {
		tfPlanPreferredSingleSignOnMode := tfPlan.PreferredSingleSignOnMode.ValueString()
		requestBody.SetPreferredSingleSignOnMode(&tfPlanPreferredSingleSignOnMode)
	} else {
		tfPlan.PreferredSingleSignOnMode = types.StringNull()
	}

	if !tfPlan.PreferredTokenSigningKeyThumbprint.IsUnknown() {
		tfPlanPreferredTokenSigningKeyThumbprint := tfPlan.PreferredTokenSigningKeyThumbprint.ValueString()
		requestBody.SetPreferredTokenSigningKeyThumbprint(&tfPlanPreferredTokenSigningKeyThumbprint)
	} else {
		tfPlan.PreferredTokenSigningKeyThumbprint = types.StringNull()
	}

	if len(tfPlan.ReplyUrls.Elements()) > 0 {
		var replyUrls []string
		for _, i := range tfPlan.ReplyUrls.Elements() {
			replyUrls = append(replyUrls, i.String())
		}
		requestBody.SetReplyUrls(replyUrls)
	} else {
		tfPlan.ReplyUrls = types.ListNull(types.StringType)
	}

	if len(tfPlan.ResourceSpecificApplicationPermissions.Elements()) > 0 {
		var tfPlanResourceSpecificApplicationPermissions []models.ResourceSpecificPermissionable
		for _, i := range tfPlan.ResourceSpecificApplicationPermissions.Elements() {
			resourceSpecificApplicationPermissions := models.NewResourceSpecificPermission()
			resourceSpecificApplicationPermissionsModel := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &resourceSpecificApplicationPermissionsModel)

			if !resourceSpecificApplicationPermissionsModel.Description.IsUnknown() {
				tfPlanDescription := resourceSpecificApplicationPermissionsModel.Description.ValueString()
				resourceSpecificApplicationPermissions.SetDescription(&tfPlanDescription)
			} else {
				resourceSpecificApplicationPermissionsModel.Description = types.StringNull()
			}

			if !resourceSpecificApplicationPermissionsModel.DisplayName.IsUnknown() {
				tfPlanDisplayName := resourceSpecificApplicationPermissionsModel.DisplayName.ValueString()
				resourceSpecificApplicationPermissions.SetDisplayName(&tfPlanDisplayName)
			} else {
				resourceSpecificApplicationPermissionsModel.DisplayName = types.StringNull()
			}

			if !resourceSpecificApplicationPermissionsModel.Id.IsUnknown() {
				tfPlanId := resourceSpecificApplicationPermissionsModel.Id.ValueString()
				u, _ := uuid.Parse(tfPlanId)
				resourceSpecificApplicationPermissions.SetId(&u)
			} else {
				resourceSpecificApplicationPermissionsModel.Id = types.StringNull()
			}

			if !resourceSpecificApplicationPermissionsModel.IsEnabled.IsUnknown() {
				tfPlanIsEnabled := resourceSpecificApplicationPermissionsModel.IsEnabled.ValueBool()
				resourceSpecificApplicationPermissions.SetIsEnabled(&tfPlanIsEnabled)
			} else {
				resourceSpecificApplicationPermissionsModel.IsEnabled = types.BoolNull()
			}

			if !resourceSpecificApplicationPermissionsModel.Value.IsUnknown() {
				tfPlanValue := resourceSpecificApplicationPermissionsModel.Value.ValueString()
				resourceSpecificApplicationPermissions.SetValue(&tfPlanValue)
			} else {
				resourceSpecificApplicationPermissionsModel.Value = types.StringNull()
			}
		}
		requestBody.SetResourceSpecificApplicationPermissions(tfPlanResourceSpecificApplicationPermissions)
	} else {
		tfPlan.ResourceSpecificApplicationPermissions = types.ListNull(tfPlan.ResourceSpecificApplicationPermissions.ElementType(ctx))
	}

	if !tfPlan.SamlSingleSignOnSettings.IsUnknown() {
		samlSingleSignOnSettings := models.NewSamlSingleSignOnSettings()
		samlSingleSignOnSettingsModel := servicePrincipalSamlSingleSignOnSettingsModel{}
		tfPlan.SamlSingleSignOnSettings.As(ctx, &samlSingleSignOnSettingsModel, basetypes.ObjectAsOptions{})

		if !samlSingleSignOnSettingsModel.RelayState.IsUnknown() {
			tfPlanRelayState := samlSingleSignOnSettingsModel.RelayState.ValueString()
			samlSingleSignOnSettings.SetRelayState(&tfPlanRelayState)
		} else {
			samlSingleSignOnSettingsModel.RelayState = types.StringNull()
		}
		requestBody.SetSamlSingleSignOnSettings(samlSingleSignOnSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, samlSingleSignOnSettingsModel.AttributeTypes(), samlSingleSignOnSettingsModel)
		tfPlan.SamlSingleSignOnSettings = objectValue
	} else {
		tfPlan.SamlSingleSignOnSettings = types.ObjectNull(tfPlan.SamlSingleSignOnSettings.AttributeTypes(ctx))
	}

	if len(tfPlan.ServicePrincipalNames.Elements()) > 0 {
		var servicePrincipalNames []string
		for _, i := range tfPlan.ServicePrincipalNames.Elements() {
			servicePrincipalNames = append(servicePrincipalNames, i.String())
		}
		requestBody.SetServicePrincipalNames(servicePrincipalNames)
	} else {
		tfPlan.ServicePrincipalNames = types.ListNull(types.StringType)
	}

	if !tfPlan.ServicePrincipalType.IsUnknown() {
		tfPlanServicePrincipalType := tfPlan.ServicePrincipalType.ValueString()
		requestBody.SetServicePrincipalType(&tfPlanServicePrincipalType)
	} else {
		tfPlan.ServicePrincipalType = types.StringNull()
	}

	if !tfPlan.SignInAudience.IsUnknown() {
		tfPlanSignInAudience := tfPlan.SignInAudience.ValueString()
		requestBody.SetSignInAudience(&tfPlanSignInAudience)
	} else {
		tfPlan.SignInAudience = types.StringNull()
	}

	if len(tfPlan.Tags.Elements()) > 0 {
		var tags []string
		for _, i := range tfPlan.Tags.Elements() {
			tags = append(tags, i.String())
		}
		requestBody.SetTags(tags)
	} else {
		tfPlan.Tags = types.ListNull(types.StringType)
	}

	if !tfPlan.TokenEncryptionKeyId.IsUnknown() {
		tfPlanTokenEncryptionKeyId := tfPlan.TokenEncryptionKeyId.ValueString()
		u, _ := uuid.Parse(tfPlanTokenEncryptionKeyId)
		requestBody.SetTokenEncryptionKeyId(&u)
	} else {
		tfPlan.TokenEncryptionKeyId = types.StringNull()
	}

	if !tfPlan.VerifiedPublisher.IsUnknown() {
		verifiedPublisher := models.NewVerifiedPublisher()
		verifiedPublisherModel := servicePrincipalVerifiedPublisherModel{}
		tfPlan.VerifiedPublisher.As(ctx, &verifiedPublisherModel, basetypes.ObjectAsOptions{})

		if !verifiedPublisherModel.AddedDateTime.IsUnknown() {
			tfPlanAddedDateTime := verifiedPublisherModel.AddedDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanAddedDateTime)
			verifiedPublisher.SetAddedDateTime(&t)
		} else {
			verifiedPublisherModel.AddedDateTime = types.StringNull()
		}

		if !verifiedPublisherModel.DisplayName.IsUnknown() {
			tfPlanDisplayName := verifiedPublisherModel.DisplayName.ValueString()
			verifiedPublisher.SetDisplayName(&tfPlanDisplayName)
		} else {
			verifiedPublisherModel.DisplayName = types.StringNull()
		}

		if !verifiedPublisherModel.VerifiedPublisherId.IsUnknown() {
			tfPlanVerifiedPublisherId := verifiedPublisherModel.VerifiedPublisherId.ValueString()
			verifiedPublisher.SetVerifiedPublisherId(&tfPlanVerifiedPublisherId)
		} else {
			verifiedPublisherModel.VerifiedPublisherId = types.StringNull()
		}
		requestBody.SetVerifiedPublisher(verifiedPublisher)
		objectValue, _ := types.ObjectValueFrom(ctx, verifiedPublisherModel.AttributeTypes(), verifiedPublisherModel)
		tfPlan.VerifiedPublisher = objectValue
	} else {
		tfPlan.VerifiedPublisher = types.ObjectNull(tfPlan.VerifiedPublisher.AttributeTypes(ctx))
	}

	// Create new servicePrincipal
	result, err := r.client.ServicePrincipals().Post(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating service_principal",
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
	// Retrieve values from plan
	var plan servicePrincipalModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state servicePrincipalModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewServicePrincipal()

	if !plan.Id.Equal(state.Id) {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	}

	if !plan.DeletedDateTime.Equal(state.DeletedDateTime) {
		planDeletedDateTime := plan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !plan.AccountEnabled.Equal(state.AccountEnabled) {
		planAccountEnabled := plan.AccountEnabled.ValueBool()
		requestBody.SetAccountEnabled(&planAccountEnabled)
	}

	if !plan.AddIns.Equal(state.AddIns) {
		var planAddIns []models.AddInable
		for k, i := range plan.AddIns.Elements() {
			addIns := models.NewAddIn()
			addInsModel := servicePrincipalAddInModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &addInsModel)
			addInsState := servicePrincipalAddInModel{}
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
					propertiesModel := servicePrincipalKeyValueModel{}
					types.ListValueFrom(ctx, i.Type(ctx), &propertiesModel)
					propertiesState := servicePrincipalKeyValueModel{}
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

	if !plan.AlternativeNames.Equal(state.AlternativeNames) {
		var alternativeNames []string
		for _, i := range plan.AlternativeNames.Elements() {
			alternativeNames = append(alternativeNames, i.String())
		}
		requestBody.SetAlternativeNames(alternativeNames)
	}

	if !plan.AppDescription.Equal(state.AppDescription) {
		planAppDescription := plan.AppDescription.ValueString()
		requestBody.SetAppDescription(&planAppDescription)
	}

	if !plan.AppDisplayName.Equal(state.AppDisplayName) {
		planAppDisplayName := plan.AppDisplayName.ValueString()
		requestBody.SetAppDisplayName(&planAppDisplayName)
	}

	if !plan.AppId.Equal(state.AppId) {
		planAppId := plan.AppId.ValueString()
		requestBody.SetAppId(&planAppId)
	}

	if !plan.AppOwnerOrganizationId.Equal(state.AppOwnerOrganizationId) {
		planAppOwnerOrganizationId := plan.AppOwnerOrganizationId.ValueString()
		u, _ := uuid.Parse(planAppOwnerOrganizationId)
		requestBody.SetAppOwnerOrganizationId(&u)
	}

	if !plan.AppRoleAssignmentRequired.Equal(state.AppRoleAssignmentRequired) {
		planAppRoleAssignmentRequired := plan.AppRoleAssignmentRequired.ValueBool()
		requestBody.SetAppRoleAssignmentRequired(&planAppRoleAssignmentRequired)
	}

	if !plan.AppRoles.Equal(state.AppRoles) {
		var planAppRoles []models.AppRoleable
		for k, i := range plan.AppRoles.Elements() {
			appRoles := models.NewAppRole()
			appRolesModel := servicePrincipalAppRoleModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &appRolesModel)
			appRolesState := servicePrincipalAppRoleModel{}
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

	if !plan.CustomSecurityAttributes.Equal(state.CustomSecurityAttributes) {
		customSecurityAttributes := models.NewCustomSecurityAttributeValue()
		customSecurityAttributesModel := servicePrincipalCustomSecurityAttributeValueModel{}
		plan.CustomSecurityAttributes.As(ctx, &customSecurityAttributesModel, basetypes.ObjectAsOptions{})
		customSecurityAttributesState := servicePrincipalCustomSecurityAttributeValueModel{}
		state.CustomSecurityAttributes.As(ctx, &customSecurityAttributesState, basetypes.ObjectAsOptions{})

		requestBody.SetCustomSecurityAttributes(customSecurityAttributes)
		objectValue, _ := types.ObjectValueFrom(ctx, customSecurityAttributesModel.AttributeTypes(), customSecurityAttributesModel)
		plan.CustomSecurityAttributes = objectValue
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

	if !plan.Homepage.Equal(state.Homepage) {
		planHomepage := plan.Homepage.ValueString()
		requestBody.SetHomepage(&planHomepage)
	}

	if !plan.Info.Equal(state.Info) {
		info := models.NewInformationalUrl()
		infoModel := servicePrincipalInformationalUrlModel{}
		plan.Info.As(ctx, &infoModel, basetypes.ObjectAsOptions{})
		infoState := servicePrincipalInformationalUrlModel{}
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

	if !plan.KeyCredentials.Equal(state.KeyCredentials) {
		var planKeyCredentials []models.KeyCredentialable
		for k, i := range plan.KeyCredentials.Elements() {
			keyCredentials := models.NewKeyCredential()
			keyCredentialsModel := servicePrincipalKeyCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &keyCredentialsModel)
			keyCredentialsState := servicePrincipalKeyCredentialModel{}
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

	if !plan.LoginUrl.Equal(state.LoginUrl) {
		planLoginUrl := plan.LoginUrl.ValueString()
		requestBody.SetLoginUrl(&planLoginUrl)
	}

	if !plan.LogoutUrl.Equal(state.LogoutUrl) {
		planLogoutUrl := plan.LogoutUrl.ValueString()
		requestBody.SetLogoutUrl(&planLogoutUrl)
	}

	if !plan.Notes.Equal(state.Notes) {
		planNotes := plan.Notes.ValueString()
		requestBody.SetNotes(&planNotes)
	}

	if !plan.NotificationEmailAddresses.Equal(state.NotificationEmailAddresses) {
		var notificationEmailAddresses []string
		for _, i := range plan.NotificationEmailAddresses.Elements() {
			notificationEmailAddresses = append(notificationEmailAddresses, i.String())
		}
		requestBody.SetNotificationEmailAddresses(notificationEmailAddresses)
	}

	if !plan.Oauth2PermissionScopes.Equal(state.Oauth2PermissionScopes) {
		var planOauth2PermissionScopes []models.PermissionScopeable
		for k, i := range plan.Oauth2PermissionScopes.Elements() {
			oauth2PermissionScopes := models.NewPermissionScope()
			oauth2PermissionScopesModel := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &oauth2PermissionScopesModel)
			oauth2PermissionScopesState := servicePrincipalPermissionScopeModel{}
			types.ListValueFrom(ctx, state.Oauth2PermissionScopes.Elements()[k].Type(ctx), &oauth2PermissionScopesModel)

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
		requestBody.SetOauth2PermissionScopes(planOauth2PermissionScopes)
	}

	if !plan.PasswordCredentials.Equal(state.PasswordCredentials) {
		var planPasswordCredentials []models.PasswordCredentialable
		for k, i := range plan.PasswordCredentials.Elements() {
			passwordCredentials := models.NewPasswordCredential()
			passwordCredentialsModel := servicePrincipalPasswordCredentialModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &passwordCredentialsModel)
			passwordCredentialsState := servicePrincipalPasswordCredentialModel{}
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

	if !plan.PreferredSingleSignOnMode.Equal(state.PreferredSingleSignOnMode) {
		planPreferredSingleSignOnMode := plan.PreferredSingleSignOnMode.ValueString()
		requestBody.SetPreferredSingleSignOnMode(&planPreferredSingleSignOnMode)
	}

	if !plan.PreferredTokenSigningKeyThumbprint.Equal(state.PreferredTokenSigningKeyThumbprint) {
		planPreferredTokenSigningKeyThumbprint := plan.PreferredTokenSigningKeyThumbprint.ValueString()
		requestBody.SetPreferredTokenSigningKeyThumbprint(&planPreferredTokenSigningKeyThumbprint)
	}

	if !plan.ReplyUrls.Equal(state.ReplyUrls) {
		var replyUrls []string
		for _, i := range plan.ReplyUrls.Elements() {
			replyUrls = append(replyUrls, i.String())
		}
		requestBody.SetReplyUrls(replyUrls)
	}

	if !plan.ResourceSpecificApplicationPermissions.Equal(state.ResourceSpecificApplicationPermissions) {
		var planResourceSpecificApplicationPermissions []models.ResourceSpecificPermissionable
		for k, i := range plan.ResourceSpecificApplicationPermissions.Elements() {
			resourceSpecificApplicationPermissions := models.NewResourceSpecificPermission()
			resourceSpecificApplicationPermissionsModel := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &resourceSpecificApplicationPermissionsModel)
			resourceSpecificApplicationPermissionsState := servicePrincipalResourceSpecificPermissionModel{}
			types.ListValueFrom(ctx, state.ResourceSpecificApplicationPermissions.Elements()[k].Type(ctx), &resourceSpecificApplicationPermissionsModel)

			if !resourceSpecificApplicationPermissionsModel.Description.Equal(resourceSpecificApplicationPermissionsState.Description) {
				planDescription := resourceSpecificApplicationPermissionsModel.Description.ValueString()
				resourceSpecificApplicationPermissions.SetDescription(&planDescription)
			}

			if !resourceSpecificApplicationPermissionsModel.DisplayName.Equal(resourceSpecificApplicationPermissionsState.DisplayName) {
				planDisplayName := resourceSpecificApplicationPermissionsModel.DisplayName.ValueString()
				resourceSpecificApplicationPermissions.SetDisplayName(&planDisplayName)
			}

			if !resourceSpecificApplicationPermissionsModel.Id.Equal(resourceSpecificApplicationPermissionsState.Id) {
				planId := resourceSpecificApplicationPermissionsModel.Id.ValueString()
				u, _ := uuid.Parse(planId)
				resourceSpecificApplicationPermissions.SetId(&u)
			}

			if !resourceSpecificApplicationPermissionsModel.IsEnabled.Equal(resourceSpecificApplicationPermissionsState.IsEnabled) {
				planIsEnabled := resourceSpecificApplicationPermissionsModel.IsEnabled.ValueBool()
				resourceSpecificApplicationPermissions.SetIsEnabled(&planIsEnabled)
			}

			if !resourceSpecificApplicationPermissionsModel.Value.Equal(resourceSpecificApplicationPermissionsState.Value) {
				planValue := resourceSpecificApplicationPermissionsModel.Value.ValueString()
				resourceSpecificApplicationPermissions.SetValue(&planValue)
			}
		}
		requestBody.SetResourceSpecificApplicationPermissions(planResourceSpecificApplicationPermissions)
	}

	if !plan.SamlSingleSignOnSettings.Equal(state.SamlSingleSignOnSettings) {
		samlSingleSignOnSettings := models.NewSamlSingleSignOnSettings()
		samlSingleSignOnSettingsModel := servicePrincipalSamlSingleSignOnSettingsModel{}
		plan.SamlSingleSignOnSettings.As(ctx, &samlSingleSignOnSettingsModel, basetypes.ObjectAsOptions{})
		samlSingleSignOnSettingsState := servicePrincipalSamlSingleSignOnSettingsModel{}
		state.SamlSingleSignOnSettings.As(ctx, &samlSingleSignOnSettingsState, basetypes.ObjectAsOptions{})

		if !samlSingleSignOnSettingsModel.RelayState.Equal(samlSingleSignOnSettingsState.RelayState) {
			planRelayState := samlSingleSignOnSettingsModel.RelayState.ValueString()
			samlSingleSignOnSettings.SetRelayState(&planRelayState)
		}
		requestBody.SetSamlSingleSignOnSettings(samlSingleSignOnSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, samlSingleSignOnSettingsModel.AttributeTypes(), samlSingleSignOnSettingsModel)
		plan.SamlSingleSignOnSettings = objectValue
	}

	if !plan.ServicePrincipalNames.Equal(state.ServicePrincipalNames) {
		var servicePrincipalNames []string
		for _, i := range plan.ServicePrincipalNames.Elements() {
			servicePrincipalNames = append(servicePrincipalNames, i.String())
		}
		requestBody.SetServicePrincipalNames(servicePrincipalNames)
	}

	if !plan.ServicePrincipalType.Equal(state.ServicePrincipalType) {
		planServicePrincipalType := plan.ServicePrincipalType.ValueString()
		requestBody.SetServicePrincipalType(&planServicePrincipalType)
	}

	if !plan.SignInAudience.Equal(state.SignInAudience) {
		planSignInAudience := plan.SignInAudience.ValueString()
		requestBody.SetSignInAudience(&planSignInAudience)
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

	if !plan.VerifiedPublisher.Equal(state.VerifiedPublisher) {
		verifiedPublisher := models.NewVerifiedPublisher()
		verifiedPublisherModel := servicePrincipalVerifiedPublisherModel{}
		plan.VerifiedPublisher.As(ctx, &verifiedPublisherModel, basetypes.ObjectAsOptions{})
		verifiedPublisherState := servicePrincipalVerifiedPublisherModel{}
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

	// Update servicePrincipal
	_, err := r.client.ServicePrincipals().ByServicePrincipalId(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating service_principal",
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
func (r *servicePrincipalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state servicePrincipalModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete servicePrincipal
	err := r.client.ServicePrincipals().ByServicePrincipalId(state.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting service_principal",
			err.Error(),
		)
		return
	}

}
