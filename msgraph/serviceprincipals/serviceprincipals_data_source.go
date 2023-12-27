package serviceprincipals

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
						"id": schema.StringAttribute{
							Description: "The unique identifier for an entity. Read-only.",
							Computed:    true,
						},
						"deleted_date_time": schema.StringAttribute{
							Description: "Date and time when this object was deleted. Always null when the object hasn't been deleted.",
							Computed:    true,
						},
						"account_enabled": schema.BoolAttribute{
							Description: "true if the service principal account is enabled; otherwise, false. If set to false, then no users will be able to sign in to this app, even if they are assigned to it. Supports $filter (eq, ne, not, in).",
							Computed:    true,
						},
						"add_ins": schema.ListNestedAttribute{
							Description: "Defines custom behavior that a consuming service can use to call an app in specific contexts. For example, applications that can render file streams may set the addIns property for its 'FileHandler' functionality. This will let services like Microsoft 365 call the application in the context of a document the user is working on.",
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
						"alternative_names": schema.ListAttribute{
							Description: "Used to retrieve service principals by subscription, identify resource group and full resource ids for managed identities. Supports $filter (eq, not, ge, le, startsWith).",
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
							Description: "Contains the tenant id where the application is registered. This is applicable only to service principals backed by applications. Supports $filter (eq, ne, NOT, ge, le).",
							Computed:    true,
						},
						"app_role_assignment_required": schema.BoolAttribute{
							Description: "Specifies whether users or other service principals need to be granted an app role assignment for this service principal before users can sign in or apps can get tokens. The default value is false. Not nullable. Supports $filter (eq, ne, NOT).",
							Computed:    true,
						},
						"app_roles": schema.ListNestedAttribute{
							Description: "The roles exposed by the application which this service principal represents. For more information see the appRoles property definition on the application entity. Not nullable.",
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
							Description: "Unique identifier of the applicationTemplate that the servicePrincipal was created from. Read-only. Supports $filter (eq, ne, NOT, startsWith).",
							Computed:    true,
						},
						"custom_security_attributes": schema.SingleNestedAttribute{
							Description: "An open complex type that holds the value of a custom security attribute that is assigned to a directory object. Nullable. Returned only on $select. Supports $filter (eq, ne, not, startsWith). Filter value is case sensitive.",
							Computed:    true,
							Attributes:  map[string]schema.Attribute{},
						},
						"description": schema.StringAttribute{
							Description: "Free text field to provide an internal end-user facing description of the service principal. End-user portals such MyApps will display the application description in this field. The maximum allowed size is 1024 characters. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
							Computed:    true,
						},
						"disabled_by_microsoft_status": schema.StringAttribute{
							Description: "Specifies whether Microsoft has disabled the registered application. Possible values are: null (default value), NotDisabled, and DisabledDueToViolationOfServicesAgreement (reasons may include suspicious, abusive, or malicious activity, or a violation of the Microsoft Services Agreement).  Supports $filter (eq, ne, not).",
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
						"login_url": schema.StringAttribute{
							Description: "Specifies the URL where the service provider redirects the user to Microsoft Entra ID to authenticate. Microsoft Entra ID uses the URL to launch the application from Microsoft 365 or the Microsoft Entra My Apps. When blank, Microsoft Entra ID performs IdP-initiated sign-on for applications configured with SAML-based single sign-on. The user launches the application from Microsoft 365, the Microsoft Entra My Apps, or the Microsoft Entra SSO URL.",
							Computed:    true,
						},
						"logout_url": schema.StringAttribute{
							Description: "Specifies the URL that will be used by Microsoft's authorization service to logout an user using OpenId Connect front-channel, back-channel or SAML logout protocols.",
							Computed:    true,
						},
						"notes": schema.StringAttribute{
							Description: "Free text field to capture information about the service principal, typically used for operational purposes. Maximum allowed size is 1024 characters.",
							Computed:    true,
						},
						"notification_email_addresses": schema.ListAttribute{
							Description: "Specifies the list of email addresses where Microsoft Entra ID sends a notification when the active certificate is near the expiration date. This is only for the certificates used to sign the SAML token issued for Microsoft Entra Gallery applications.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"oauth_2_permission_scopes": schema.ListNestedAttribute{
							Description: "The delegated permissions exposed by the application. For more information see the oauth2PermissionScopes property on the application entity's api property. Not nullable.",
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
							Description: "Specifies the single sign-on mode configured for this application. Microsoft Entra ID uses the preferred single sign-on mode to launch the application from Microsoft 365 or the My Apps portal. The supported values are password, saml, notSupported, and oidc.",
							Computed:    true,
						},
						"preferred_token_signing_key_thumbprint": schema.StringAttribute{
							Description: "This property can be used on SAML applications (apps that have preferredSingleSignOnMode set to saml) to control which certificate is used to sign the SAML responses. For applications that are not SAML, do not write or otherwise rely on this property.",
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
							Description: "Contains the list of identifiersUris, copied over from the associated application. Additional values can be added to hybrid applications. These values can be used to identify the permissions exposed by this app within Microsoft Entra ID. For example,Client apps can specify a resource URI which is based on the values of this property to acquire an access token, which is the URI returned in the 'aud' claim.The any operator is required for filter expressions on multi-valued properties. Not nullable.  Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"service_principal_type": schema.StringAttribute{
							Description: "Identifies whether the service principal represents an application, a managed identity, or a legacy application. This is set by Microsoft Entra ID internally. The servicePrincipalType property can be set to three different values: Application - A service principal that represents an application or service. The appId property identifies the associated app registration, and matches the appId of an application, possibly from a different tenant. If the associated app registration is missing, tokens are not issued for the service principal.ManagedIdentity - A service principal that represents a managed identity. Service principals representing managed identities can be granted access and permissions, but cannot be updated or modified directly.Legacy - A service principal that represents an app created before app registrations, or through legacy experiences. Legacy service principal can have credentials, service principal names, reply URLs, and other properties which are editable by an authorized user, but does not have an associated app registration. The appId value does not associate the service principal with an app registration. The service principal can only be used in the tenant where it was created.SocialIdp - For internal use.",
							Computed:    true,
						},
						"sign_in_audience": schema.StringAttribute{
							Description: "Specifies the Microsoft accounts that are supported for the current application. Read-only. Supported values are:AzureADMyOrg: Users with a Microsoft work or school account in my organization's Microsoft Entra tenant (single-tenant).AzureADMultipleOrgs: Users with a Microsoft work or school account in any organization's Microsoft Entra tenant (multi-tenant).AzureADandPersonalMicrosoftAccount: Users with a personal Microsoft account, or a work or school account in any organization's Microsoft Entra tenant.PersonalMicrosoftAccount: Users with a personal Microsoft account only.",
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
							Description: "Specifies the verified publisher of the application which this service principal represents.",
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
	var state servicePrincipalsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := serviceprincipals.ServicePrincipalsRequestBuilderGetRequestConfiguration{
		QueryParameters: &serviceprincipals.ServicePrincipalsRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"deletedDateTime",
				"accountEnabled",
				"addIns",
				"alternativeNames",
				"appDescription",
				"appDisplayName",
				"appId",
				"applicationTemplateId",
				"appOwnerOrganizationId",
				"appRoleAssignmentRequired",
				"appRoles",
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
				"appManagementPolicies",
				"appRoleAssignedTo",
				"appRoleAssignments",
				"claimsMappingPolicies",
				"createdObjects",
				"delegatedPermissionClassifications",
				"endpoints",
				"federatedIdentityCredentials",
				"homeRealmDiscoveryPolicies",
				"memberOf",
				"oauth2PermissionGrants",
				"ownedObjects",
				"owners",
				"remoteDesktopSecurityConfiguration",
				"tokenIssuancePolicies",
				"tokenLifetimePolicies",
				"transitiveMemberOf",
				"synchronization",
			},
		},
	}

	result, err := d.client.ServicePrincipals().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting service_principals",
			err.Error(),
		)
		return
	}

	for _, v := range result.GetValue() {
		value := new(servicePrincipalsValueModel)

		if v.GetId() != nil {
			value.Id = types.StringValue(*v.GetId())
		}
		if v.GetDeletedDateTime() != nil {
			value.DeletedDateTime = types.StringValue(v.GetDeletedDateTime().String())
		}
		if v.GetAccountEnabled() != nil {
			value.AccountEnabled = types.BoolValue(*v.GetAccountEnabled())
		}
		for _, v := range v.GetAddIns() {
			addIns := new(servicePrincipalsAddInsModel)

			if v.GetId() != nil {
				addIns.Id = types.StringValue(v.GetId().String())
			}
			for _, v := range v.GetProperties() {
				properties := new(servicePrincipalsPropertiesModel)

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
		for _, v := range v.GetAlternativeNames() {
			value.AlternativeNames = append(value.AlternativeNames, types.StringValue(v))
		}
		if v.GetAppDescription() != nil {
			value.AppDescription = types.StringValue(*v.GetAppDescription())
		}
		if v.GetAppDisplayName() != nil {
			value.AppDisplayName = types.StringValue(*v.GetAppDisplayName())
		}
		if v.GetAppId() != nil {
			value.AppId = types.StringValue(*v.GetAppId())
		}
		if v.GetAppOwnerOrganizationId() != nil {
			value.AppOwnerOrganizationId = types.StringValue(v.GetAppOwnerOrganizationId().String())
		}
		if v.GetAppRoleAssignmentRequired() != nil {
			value.AppRoleAssignmentRequired = types.BoolValue(*v.GetAppRoleAssignmentRequired())
		}
		for _, v := range v.GetAppRoles() {
			appRoles := new(servicePrincipalsAppRolesModel)

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
		if v.GetCustomSecurityAttributes() != nil {
			value.CustomSecurityAttributes = new(servicePrincipalsCustomSecurityAttributesModel)

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
		if v.GetHomepage() != nil {
			value.Homepage = types.StringValue(*v.GetHomepage())
		}
		if v.GetInfo() != nil {
			value.Info = new(servicePrincipalsInfoModel)

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
		for _, v := range v.GetKeyCredentials() {
			keyCredentials := new(servicePrincipalsKeyCredentialsModel)

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
		if v.GetLoginUrl() != nil {
			value.LoginUrl = types.StringValue(*v.GetLoginUrl())
		}
		if v.GetLogoutUrl() != nil {
			value.LogoutUrl = types.StringValue(*v.GetLogoutUrl())
		}
		if v.GetNotes() != nil {
			value.Notes = types.StringValue(*v.GetNotes())
		}
		for _, v := range v.GetNotificationEmailAddresses() {
			value.NotificationEmailAddresses = append(value.NotificationEmailAddresses, types.StringValue(v))
		}
		for _, v := range v.GetOauth2PermissionScopes() {
			oauth2PermissionScopes := new(servicePrincipalsOauth2PermissionScopesModel)

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
			value.Oauth2PermissionScopes = append(value.Oauth2PermissionScopes, *oauth2PermissionScopes)
		}
		for _, v := range v.GetPasswordCredentials() {
			passwordCredentials := new(servicePrincipalsPasswordCredentialsModel)

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
		if v.GetPreferredSingleSignOnMode() != nil {
			value.PreferredSingleSignOnMode = types.StringValue(*v.GetPreferredSingleSignOnMode())
		}
		if v.GetPreferredTokenSigningKeyThumbprint() != nil {
			value.PreferredTokenSigningKeyThumbprint = types.StringValue(*v.GetPreferredTokenSigningKeyThumbprint())
		}
		for _, v := range v.GetReplyUrls() {
			value.ReplyUrls = append(value.ReplyUrls, types.StringValue(v))
		}
		for _, v := range v.GetResourceSpecificApplicationPermissions() {
			resourceSpecificApplicationPermissions := new(servicePrincipalsResourceSpecificApplicationPermissionsModel)

			if v.GetDescription() != nil {
				resourceSpecificApplicationPermissions.Description = types.StringValue(*v.GetDescription())
			}
			if v.GetDisplayName() != nil {
				resourceSpecificApplicationPermissions.DisplayName = types.StringValue(*v.GetDisplayName())
			}
			if v.GetId() != nil {
				resourceSpecificApplicationPermissions.Id = types.StringValue(v.GetId().String())
			}
			if v.GetIsEnabled() != nil {
				resourceSpecificApplicationPermissions.IsEnabled = types.BoolValue(*v.GetIsEnabled())
			}
			if v.GetValue() != nil {
				resourceSpecificApplicationPermissions.Value = types.StringValue(*v.GetValue())
			}
			value.ResourceSpecificApplicationPermissions = append(value.ResourceSpecificApplicationPermissions, *resourceSpecificApplicationPermissions)
		}
		if v.GetSamlSingleSignOnSettings() != nil {
			value.SamlSingleSignOnSettings = new(servicePrincipalsSamlSingleSignOnSettingsModel)

			if v.GetSamlSingleSignOnSettings().GetRelayState() != nil {
				value.SamlSingleSignOnSettings.RelayState = types.StringValue(*v.GetSamlSingleSignOnSettings().GetRelayState())
			}
		}
		for _, v := range v.GetServicePrincipalNames() {
			value.ServicePrincipalNames = append(value.ServicePrincipalNames, types.StringValue(v))
		}
		if v.GetServicePrincipalType() != nil {
			value.ServicePrincipalType = types.StringValue(*v.GetServicePrincipalType())
		}
		if v.GetSignInAudience() != nil {
			value.SignInAudience = types.StringValue(*v.GetSignInAudience())
		}
		for _, v := range v.GetTags() {
			value.Tags = append(value.Tags, types.StringValue(v))
		}
		if v.GetTokenEncryptionKeyId() != nil {
			value.TokenEncryptionKeyId = types.StringValue(v.GetTokenEncryptionKeyId().String())
		}
		if v.GetVerifiedPublisher() != nil {
			value.VerifiedPublisher = new(servicePrincipalsVerifiedPublisherModel)

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
		state.Value = append(state.Value, *value)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
