package serviceprincipals

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/serviceprincipals"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &servicePrincipalsDataSource{}
	_ datasource.DataSourceWithConfigure = &servicePrincipalsDataSource{}
)

// NewServicePrincipalsDataSource is a helper function to simplify the provider implementation.
func NewServicePrincipalsDataSource() datasource.DataSource {
	return &servicePrincipalsDataSource{}
}

// servicePrincipalsDataSource is the data source implementation.
type servicePrincipalsDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *servicePrincipalsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_principals"
}

// Configure adds the provider configured client to the data source.
func (d *servicePrincipalsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *servicePrincipalsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"value": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"account_enabled": schema.BoolAttribute{
							Description: "true if the service principal account is enabled; otherwise, false. If set to false, then no users are able to sign in to this app, even if they're assigned to it. Supports $filter (eq, ne, not, in).",
							Computed:    true,
						},
						"add_ins": schema.ListNestedAttribute{
							Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams may set the addIns property for its 'FileHandler' functionality. This lets services like Microsoft 365 call the application in the context of a document the user is working on.",
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
						"alternative_names": schema.ListAttribute{
							Description: "Used to retrieve service principals by subscription, identify resource group and full resource IDs for managed identities. Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"app_description": schema.StringAttribute{
							Description: "The description exposed by the associated application.",
							Computed:    true,
						},
						"app_display_name": schema.StringAttribute{
							Description: "The display name exposed by the associated application.",
							Computed:    true,
						},
						"app_id": schema.StringAttribute{
							Description: "The unique identifier for the associated application (its appId property). Alternate key. Supports $filter (eq, ne, not, in, startsWith).",
							Computed:    true,
						},
						"app_owner_organization_id": schema.StringAttribute{
							Description: "Contains the tenant ID where the application is registered. This is applicable only to service principals backed by applications. Supports $filter (eq, ne, NOT, ge, le).",
							Computed:    true,
						},
						"app_role_assignment_required": schema.BoolAttribute{
							Description: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports $filter (eq, ne, NOT).",
							Computed:    true,
						},
						"app_roles": schema.ListNestedAttribute{
							Description: "The roles exposed by the application that's linked to this service principal. For more information, see the appRoles property definition on the application entity. Not nullable.",
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
							Description: "Unique identifier of the applicationTemplate. Supports $filter (eq, not, ne). Read-only. null if the service principal wasn't created from an application template.",
							Computed:    true,
						},
						"deleted_date_time": schema.StringAttribute{
							Description: "Date and time when this object was deleted. Always null when the object hasn't been deleted.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Free text field to provide an internal end-user facing description of the service principal. End-user portals such MyApps displays the application description in this field. The maximum allowed size is 1,024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
							Computed:    true,
						},
						"disabled_by_microsoft_status": schema.StringAttribute{
							Description: "Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The display name for the service principal. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
							Computed:    true,
						},
						"homepage": schema.StringAttribute{
							Description: "Home page or landing page of the application.",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "The unique identifier for an entity. Read-only.",
							Computed:    true,
						},
						"info": schema.SingleNestedAttribute{
							Description: "Basic profile information of the acquired application such as app's marketing, support, terms of service and privacy statement URLs. The terms of service and privacy statement are surfaced to users through the user consent experience. For more info, see How to: Add Terms of service and privacy statement for registered Microsoft Entra apps. Supports $filter (eq, ne, not, ge, le, and eq on null values).",
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
						"key_credentials": schema.ListNestedAttribute{
							Description: "The collection of key credentials associated with the service principal. Not nullable. Supports $filter (eq, not, ge, le).",
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
						"login_url": schema.StringAttribute{
							Description: "Specifies the URL where the service provider redirects the user to Microsoft Entra ID to authenticate. Microsoft Entra ID uses the URL to launch the application from Microsoft 365 or the Microsoft Entra My Apps. When blank, Microsoft Entra ID performs IdP-initiated sign-on for applications configured with SAML-based single sign-on. The user launches the application from Microsoft 365, the Microsoft Entra My Apps, or the Microsoft Entra SSO URL.",
							Computed:    true,
						},
						"logout_url": schema.StringAttribute{
							Description: "Specifies the URL that the Microsoft's authorization service uses to sign out a user using OpenID Connect front-channel, back-channel, or SAML sign out protocols.",
							Computed:    true,
						},
						"notes": schema.StringAttribute{
							Description: "Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1,024 characters.",
							Computed:    true,
						},
						"notification_email_addresses": schema.ListAttribute{
							Description: "Specifies the list of email addresses where Microsoft Entra ID sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Microsoft Entra Gallery applications.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"oauth_2_permission_scopes": schema.ListNestedAttribute{
							Description: "The delegated permissions exposed by the application. For more information, see the oauth2PermissionScopes property on the application entity's api property. Not nullable.",
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
						"preferred_single_sign_on_mode": schema.StringAttribute{
							Description: "Specifies the single sign-on mode configured for this application. Microsoft Entra ID uses the preferred single sign-on mode to launch the application from Microsoft 365 or the My Apps portal. The supported values are password, saml, notSupported, and oidc. Note: This field might be null for older SAML apps and for OIDC applications where it isn't set automatically.",
							Computed:    true,
						},
						"preferred_token_signing_key_thumbprint": schema.StringAttribute{
							Description: "This property can be used on SAML applications (apps that have preferredSingleSignOnMode set to saml) to control which certificate is used to sign the SAML responses. For applications that aren't SAML, don't write or otherwise rely on this property.",
							Computed:    true,
						},
						"reply_urls": schema.ListAttribute{
							Description: "The URLs that user tokens are sent to for sign in with the associated application, or the redirect URIs that OAuth 2.0 authorization codes and access tokens are sent to for the associated application. Not nullable.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"resource_specific_application_permissions": schema.ListNestedAttribute{
							Description: "The resource-specific application permissions exposed by this application. Currently, resource-specific permissions are only supported for Teams apps accessing to specific chats and teams using Microsoft Graph. Read-only.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"description": schema.StringAttribute{
										Description: "Describes the level of access that the resource-specific permission represents.",
										Computed:    true,
									},
									"display_name": schema.StringAttribute{
										Description: "The display name for the resource-specific permission.",
										Computed:    true,
									},
									"id": schema.StringAttribute{
										Description: "The unique identifier for the resource-specific application permission.",
										Computed:    true,
									},
									"is_enabled": schema.BoolAttribute{
										Description: "Indicates whether the permission is enabled.",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "The value of the permission.",
										Computed:    true,
									},
								},
							},
						},
						"saml_single_sign_on_settings": schema.SingleNestedAttribute{
							Description: "The collection for settings related to saml single sign-on.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"relay_state": schema.StringAttribute{
									Description: "The relative URI the service provider would redirect to after completion of the single sign-on flow.",
									Computed:    true,
								},
							},
						},
						"service_principal_names": schema.ListAttribute{
							Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Microsoft Entra ID. For example,Client apps can specify a resource URI that is based on the values of this property to acquire an access token, which is the URI returned in the 'aud' claim.The any operator is required for filter expressions on multi-valued properties. Not nullable.  Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"service_principal_type": schema.StringAttribute{
							Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Microsoft Entra ID internally. The servicePrincipalType property can be set to three different values: Application - A service principal that represents an application or service. The appId property identifies the associated app registration, and matches the appId of an application, possibly from a different tenant. If the associated app registration is missing, tokens aren't issued for the service principal.ManagedIdentity - A service principal that represents a managed identity. Service principals representing managed identities can be granted access and permissions, but can't be updated or modified directly.Legacy - A service principal that represents an app created before app registrations, or through legacy experiences. A legacy service principal can have credentials, service principal names, reply URLs, and other properties that are editable by an authorized user, but doesn't have an associated app registration. The appId value doesn't associate the service principal with an app registration. The service principal can only be used in the tenant where it was created.SocialIdp - For internal use.",
							Computed:    true,
						},
						"sign_in_audience": schema.StringAttribute{
							Description: "Specifies the Microsoft accounts that are supported for the current application. Read-only. Supported values are:AzureADMyOrg: Users with a Microsoft work or school account in my organization's Microsoft Entra tenant (single-tenant).AzureADMultipleOrgs: Users with a Microsoft work or school account in any organization's Microsoft Entra tenant (multitenant).AzureADandPersonalMicrosoftAccount: Users with a personal Microsoft account, or a work or school account in any organization's Microsoft Entra tenant.PersonalMicrosoftAccount: Users with a personal Microsoft account only.",
							Computed:    true,
						},
						"tags": schema.ListAttribute{
							Description: "Custom strings that can be used to categorize and identify the service principal. Not nullable. The value is the union of strings set here and on the associated application entity's tags property.Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"token_encryption_key_id": schema.StringAttribute{
							Description: "Specifies the keyId of a public key from the keyCredentials collection. When configured, Microsoft Entra ID issues tokens for this application encrypted using the key specified by this property. The application code that receives the encrypted token must use the matching private key to decrypt the token before it can be used for the signed-in user.",
							Computed:    true,
						},
						"verified_publisher": schema.SingleNestedAttribute{
							Description: "Specifies the verified publisher of the application that's linked to this service principal.",
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
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *servicePrincipalsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfStateServicePrincipals servicePrincipalsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateServicePrincipals)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
			Select: []string{
				"value",
			},
		},
	}

	responseServicePrincipals, err := d.client.ServicePrincipals().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting ServicePrincipals",
			err.Error(),
		)
		return
	}

	if len(responseServicePrincipals.GetValue()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseServicePrincipal := range responseServicePrincipals.GetValue() {
			tfStateServicePrincipal := servicePrincipalsServicePrincipalModel{}

			if responseServicePrincipal.GetAccountEnabled() != nil {
				tfStateServicePrincipal.AccountEnabled = types.BoolValue(*responseServicePrincipal.GetAccountEnabled())
			} else {
				tfStateServicePrincipal.AccountEnabled = types.BoolNull()
			}
			if len(responseServicePrincipal.GetAddIns()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseAddIn := range responseServicePrincipal.GetAddIns() {
					tfStateAddIn := servicePrincipalsAddInModel{}

					if responseAddIn.GetId() != nil {
						tfStateAddIn.Id = types.StringValue(responseAddIn.GetId().String())
					} else {
						tfStateAddIn.Id = types.StringNull()
					}
					if len(responseAddIn.GetProperties()) > 0 {
						objectValues := []basetypes.ObjectValue{}
						for _, responseKeyValue := range responseAddIn.GetProperties() {
							tfStateKeyValue := servicePrincipalsKeyValueModel{}

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
					tfStateAppRole := servicePrincipalsAppRoleModel{}

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
				tfStateInformationalUrl := servicePrincipalsInformationalUrlModel{}
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
					tfStateKeyCredential := servicePrincipalsKeyCredentialModel{}

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
					tfStatePermissionScope := servicePrincipalsPermissionScopeModel{}

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
					tfStatePasswordCredential := servicePrincipalsPasswordCredentialModel{}

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
					tfStateResourceSpecificPermission := servicePrincipalsResourceSpecificPermissionModel{}

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
				tfStateSamlSingleSignOnSettings := servicePrincipalsSamlSingleSignOnSettingsModel{}
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
				tfStateVerifiedPublisher := servicePrincipalsVerifiedPublisherModel{}
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
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateServicePrincipal.AttributeTypes(), tfStateServicePrincipal)
			objectValues = append(objectValues, objectValue)
		}
		tfStateServicePrincipals.Value, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateServicePrincipals)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
