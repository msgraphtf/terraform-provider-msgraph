package users

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &usersDataSource{}
	_ datasource.DataSourceWithConfigure = &usersDataSource{}
)

// NewUsersDataSource is a helper function to simplify the provider implementation.
func NewUsersDataSource() datasource.DataSource {
	return &usersDataSource{}
}

// usersDataSource is the data source implementation.
type usersDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *usersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_users"
}

// Configure adds the provider configured client to the data source.
func (d *usersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *usersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
							Description: "true if the account is enabled; otherwise, false. This property is required when a user is created. Returned only on $select. Supports $filter (eq, ne, not, and in).",
							Computed:    true,
						},
						"age_group": schema.StringAttribute{
							Description: "Sets the age group of the user. Allowed values: null, Minor, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).",
							Computed:    true,
						},
						"assigned_licenses": schema.ListNestedAttribute{
							Description: "The licenses that are assigned to the user, including inherited (group-based) licenses. This property doesn't differentiate between directly assigned and inherited licenses. Use the licenseAssignmentStates property to identify the directly assigned and inherited licenses. Not nullable. Returned only on $select. Supports $filter (eq, not, /$count eq 0, /$count ne 0).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"disabled_plans": schema.ListAttribute{
										Description: "A collection of the unique identifiers for plans that have been disabled. IDs are available in servicePlans > servicePlanId in the tenant's subscribedSkus or serviceStatus > servicePlanId in the tenant's companySubscription.",
										Computed:    true,
										ElementType: types.StringType,
									},
									"sku_id": schema.StringAttribute{
										Description: "The unique identifier for the SKU. Corresponds to the skuId from subscribedSkus or companySubscription.",
										Computed:    true,
									},
								},
							},
						},
						"assigned_plans": schema.ListNestedAttribute{
							Description: "The plans that are assigned to the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq and not).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"assigned_date_time": schema.StringAttribute{
										Description: "The date and time at which the plan was assigned. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
										Computed:    true,
									},
									"capability_status": schema.StringAttribute{
										Description: "Condition of the capability assignment. The possible values are Enabled, Warning, Suspended, Deleted, LockedOut. See a detailed description of each value.",
										Computed:    true,
									},
									"service": schema.StringAttribute{
										Description: "The name of the service; for example, exchange.",
										Computed:    true,
									},
									"service_plan_id": schema.StringAttribute{
										Description: "A GUID that identifies the service plan. For a complete list of GUIDs and their equivalent friendly service names, see Product names and service plan identifiers for licensing.",
										Computed:    true,
									},
								},
							},
						},
						"authorization_info": schema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"certificate_user_ids": schema.ListAttribute{
									Description: "",
									Computed:    true,
									ElementType: types.StringType,
								},
							},
						},
						"business_phones": schema.ListAttribute{
							Description: "The telephone numbers for the user. NOTE: Although it's a string collection, only one number can be set for this property. Read-only for users synced from the on-premises directory. Returned by default. Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"city": schema.StringAttribute{
							Description: "The city where the user is located. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"company_name": schema.StringAttribute{
							Description: "The name of the company that the user is associated with. This property can be useful for describing the company that a guest comes from. The maximum length is 64 characters.Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"consent_provided_for_minor": schema.StringAttribute{
							Description: "Sets whether consent was obtained for minors. Allowed values: null, Granted, Denied, and NotRequired. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).",
							Computed:    true,
						},
						"country": schema.StringAttribute{
							Description: "The country/region where the user is located; for example, US or UK. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"created_date_time": schema.StringAttribute{
							Description: "The date and time the user was created, in ISO 8601 format and UTC. The value can't be modified and is automatically populated when the entity is created. Nullable. For on-premises users, the value represents when they were first created in Microsoft Entra ID. Property is null for some users created before June 2018 and on-premises users that were synced to Microsoft Entra ID before June 2018. Read-only. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
							Computed:    true,
						},
						"creation_type": schema.StringAttribute{
							Description: "Indicates whether the user account was created through one of the following methods:  As a regular school or work account (null). As an external account (Invitation). As a local account for an Azure Active Directory B2C tenant (LocalAccount). Through self-service sign-up by an internal user using email verification (EmailVerified). Through self-service sign-up by a guest signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.Returned only on $select. Supports $filter (eq, ne, not, in).",
							Computed:    true,
						},
						"custom_security_attributes": schema.SingleNestedAttribute{
							Description: "An open complex type that holds the value of a custom security attribute that is assigned to a directory object. Nullable. Returned only on $select. Supports $filter (eq, ne, not, startsWith). The filter value is case-sensitive. To read this property, the calling app must be assigned the CustomSecAttributeAssignment.Read.All permission. To write this property, the calling app must be assigned the CustomSecAttributeAssignment.ReadWrite.All permissions. To read or write this property in delegated scenarios, the admin must be assigned the Attribute Assignment Administrator role.",
							Computed:    true,
							Attributes:  map[string]schema.Attribute{},
						},
						"department": schema.StringAttribute{
							Description: "The name of the department in which the user works. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, and eq on null values).",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The name displayed in the address book for the user. This value is usually the combination of the user's first name, middle initial, and family name. This property is required when a user is created and it can't be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values), $orderby, and $search.",
							Computed:    true,
						},
						"employee_hire_date": schema.StringAttribute{
							Description: "The date and time when the user was hired or will start work in a future hire. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
							Computed:    true,
						},
						"employee_id": schema.StringAttribute{
							Description: "The employee identifier assigned to the user by the organization. The maximum length is 16 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"employee_leave_date_time": schema.StringAttribute{
							Description: "The date and time when the user left or will leave the organization. To read this property, the calling app must be assigned the User-LifeCycleInfo.Read.All permission. To write this property, the calling app must be assigned the User.Read.All and User-LifeCycleInfo.ReadWrite.All permissions. To read this property in delegated scenarios, the admin needs at least one of the following Microsoft Entra roles: Lifecycle Workflows Administrator (least privilege), Global Reader. To write this property in delegated scenarios, the admin needs the Global Administrator role. Supports $filter (eq, ne, not , ge, le, in). For more information, see Configure the employeeLeaveDateTime property for a user.",
							Computed:    true,
						},
						"employee_org_data": schema.SingleNestedAttribute{
							Description: "Represents organization data (for example, division and costCenter) associated with a user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"cost_center": schema.StringAttribute{
									Description: "The cost center associated with the user. Returned only on $select. Supports $filter.",
									Computed:    true,
								},
								"division": schema.StringAttribute{
									Description: "The name of the division in which the user works. Returned only on $select. Supports $filter.",
									Computed:    true,
								},
							},
						},
						"employee_type": schema.StringAttribute{
							Description: "Captures enterprise worker type. For example, Employee, Contractor, Consultant, or Vendor. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith).",
							Computed:    true,
						},
						"external_user_state": schema.StringAttribute{
							Description: "For a guest invited to the tenant using the invitation API, this property represents the invited user's invitation status. For invited users, the state can be PendingAcceptance or Accepted, or null for all other users. Returned only on $select. Supports $filter (eq, ne, not , in).",
							Computed:    true,
						},
						"external_user_state_change_date_time": schema.StringAttribute{
							Description: "Shows the timestamp for the latest change to the externalUserState property. Returned only on $select. Supports $filter (eq, ne, not , in).",
							Computed:    true,
						},
						"fax_number": schema.StringAttribute{
							Description: "The fax number of the user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"given_name": schema.StringAttribute{
							Description: "The given name (first name) of the user. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"identities": schema.ListNestedAttribute{
							Description: "Represents the identities that can be used to sign in to this user account. Microsoft (also known as a local account), organizations, or social identity providers such as Facebook, Google, and Microsoft can provide identity and tie it to a user account. It might contain multiple items with the same signInType value. Returned only on $select.  Supports $filter (eq) with limitations.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"issuer": schema.StringAttribute{
										Description: "Specifies the issuer of the identity, for example facebook.com. 512 character limit. For local accounts (where signInType isn't federated), this property is the local default domain name for the tenant, for example contoso.com.  For guests from other Microsoft Entra organizations, this is the domain of the federated organization, for example contoso.com. For more information about filtering behavior for this property, see Filtering on the identities property of a user.",
										Computed:    true,
									},
									"issuer_assigned_id": schema.StringAttribute{
										Description: "Specifies the unique identifier assigned to the user by the issuer. 64 character limit. The combination of issuer and issuerAssignedId must be unique within the organization. Represents the sign-in name for the user, when signInType is set to emailAddress or userName (also known as local accounts).When signInType is set to: emailAddress (or a custom string that starts with emailAddress like emailAddress1), issuerAssignedId must be a valid email addressuserName, issuerAssignedId must begin with an alphabetical character or number, and can only contain alphanumeric characters and the following symbols: - or _  For more information about filtering behavior for this property, see Filtering on the identities property of a user.",
										Computed:    true,
									},
									"sign_in_type": schema.StringAttribute{
										Description: "Specifies the user sign-in types in your directory, such as emailAddress, userName, federated, or userPrincipalName. federated represents a unique identifier for a user from an issuer that can be in any format chosen by the issuer. Setting or updating a userPrincipalName identity updates the value of the userPrincipalName property on the user object. The validations performed on the userPrincipalName property on the user object, for example, verified domains and acceptable characters, are performed when setting or updating a userPrincipalName identity. Extra validation is enforced on issuerAssignedId when the sign-in type is set to emailAddress or userName. This property can also be set to any custom string.  For more information about filtering behavior for this property, see Filtering on the identities property of a user.",
										Computed:    true,
									},
								},
							},
						},
						"im_addresses": schema.ListAttribute{
							Description: "The instant message voice-over IP (VOIP) session initiation protocol (SIP) addresses for the user. Read-only. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"is_management_restricted": schema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"is_resource_account": schema.BoolAttribute{
							Description: "Don't use – reserved for future use.",
							Computed:    true,
						},
						"job_title": schema.StringAttribute{
							Description: "The user's job title. Maximum length is 128 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"last_password_change_date_time": schema.StringAttribute{
							Description: "The time when this Microsoft Entra user last changed their password or when their password was created, whichever date the latest action was performed. The date and time information uses ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.",
							Computed:    true,
						},
						"legal_age_group_classification": schema.StringAttribute{
							Description: "Used by enterprise applications to determine the legal age group of the user. This property is read-only and calculated based on ageGroup and consentProvidedForMinor properties. Allowed values: null, MinorWithOutParentalConsent, MinorWithParentalConsent, MinorNoParentalConsentRequired, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select.",
							Computed:    true,
						},
						"license_assignment_states": schema.ListNestedAttribute{
							Description: "State of license assignments for this user. Also indicates licenses that are directly assigned or the user inherited through group memberships. Read-only. Returned only on $select.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"assigned_by_group": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"disabled_plans": schema.ListAttribute{
										Description: "",
										Computed:    true,
										ElementType: types.StringType,
									},
									"error": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"last_updated_date_time": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"sku_id": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"state": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
								},
							},
						},
						"mail": schema.StringAttribute{
							Description: "The SMTP address for the user, for example, jeff@contoso.com. Changes to this property update the user's proxyAddresses collection to include the value as an SMTP address. This property can't contain accent characters.  NOTE: We don't recommend updating this property for Azure AD B2C user profiles. Use the otherMails property instead. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith, and eq on null values).",
							Computed:    true,
						},
						"mail_nickname": schema.StringAttribute{
							Description: "The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"mobile_phone": schema.StringAttribute{
							Description: "The primary cellular telephone number for the user. Read-only for users synced from the on-premises directory. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values) and $search.",
							Computed:    true,
						},
						"office_location": schema.StringAttribute{
							Description: "The office location in the user's place of business. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"on_premises_distinguished_name": schema.StringAttribute{
							Description: "Contains the on-premises Active Directory distinguished name or DN. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.",
							Computed:    true,
						},
						"on_premises_domain_name": schema.StringAttribute{
							Description: "Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.",
							Computed:    true,
						},
						"on_premises_extension_attributes": schema.SingleNestedAttribute{
							Description: "Contains extensionAttributes1-15 for the user. These extension attributes are also known as Exchange custom attributes 1-15. Each attribute can store up to 1024 characters. For an onPremisesSyncEnabled user, the source of authority for this set of properties is the on-premises and is read-only. For a cloud-only user (where onPremisesSyncEnabled is false), these properties can be set during the creation or update of a user object.  For a cloud-only user previously synced from on-premises Active Directory, these properties are read-only in Microsoft Graph but can be fully managed through the Exchange Admin Center or the Exchange Online V2 module in PowerShell. Returned only on $select. Supports $filter (eq, ne, not, in).",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"extension_attribute_1": schema.StringAttribute{
									Description: "First customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_10": schema.StringAttribute{
									Description: "Tenth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_11": schema.StringAttribute{
									Description: "Eleventh customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_12": schema.StringAttribute{
									Description: "Twelfth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_13": schema.StringAttribute{
									Description: "Thirteenth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_14": schema.StringAttribute{
									Description: "Fourteenth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_15": schema.StringAttribute{
									Description: "Fifteenth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_2": schema.StringAttribute{
									Description: "Second customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_3": schema.StringAttribute{
									Description: "Third customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_4": schema.StringAttribute{
									Description: "Fourth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_5": schema.StringAttribute{
									Description: "Fifth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_6": schema.StringAttribute{
									Description: "Sixth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_7": schema.StringAttribute{
									Description: "Seventh customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_8": schema.StringAttribute{
									Description: "Eighth customizable extension attribute.",
									Computed:    true,
								},
								"extension_attribute_9": schema.StringAttribute{
									Description: "Ninth customizable extension attribute.",
									Computed:    true,
								},
							},
						},
						"on_premises_immutable_id": schema.StringAttribute{
							Description: "This property is used to associate an on-premises Active Directory user account to their Microsoft Entra user object. This property must be specified when creating a new user account in the Graph if you're using a federated domain for the user's userPrincipalName (UPN) property. NOTE: The $ and _ characters can't be used when specifying this property. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).",
							Computed:    true,
						},
						"on_premises_last_sync_date_time": schema.StringAttribute{
							Description: "Indicates the last time at which the object was synced with the on-premises directory; for example: 2013-02-16T03:04:54Z. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).",
							Computed:    true,
						},
						"on_premises_provisioning_errors": schema.ListNestedAttribute{
							Description: "Errors when using Microsoft synchronization product during provisioning. Returned only on $select. Supports $filter (eq, not, ge, le).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"category": schema.StringAttribute{
										Description: "Category of the provisioning error. Note: Currently, there is only one possible value. Possible value: PropertyConflict - indicates a property value is not unique. Other objects contain the same value for the property.",
										Computed:    true,
									},
									"occurred_date_time": schema.StringAttribute{
										Description: "The date and time at which the error occurred.",
										Computed:    true,
									},
									"property_causing_error": schema.StringAttribute{
										Description: "Name of the directory property causing the error. Current possible values: UserPrincipalName or ProxyAddress",
										Computed:    true,
									},
									"value": schema.StringAttribute{
										Description: "Value of the property causing the error.",
										Computed:    true,
									},
								},
							},
						},
						"on_premises_sam_account_name": schema.StringAttribute{
							Description: "Contains the on-premises samAccountName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).",
							Computed:    true,
						},
						"on_premises_security_identifier": schema.StringAttribute{
							Description: "Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. Returned only on $select. Supports $filter (eq including on null values).",
							Computed:    true,
						},
						"on_premises_sync_enabled": schema.BoolAttribute{
							Description: "true if this user object is currently being synced from an on-premises Active Directory (AD); otherwise the user isn't being synced and can be managed in Microsoft Entra ID. Read-only. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values).",
							Computed:    true,
						},
						"on_premises_user_principal_name": schema.StringAttribute{
							Description: "Contains the on-premises userPrincipalName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).",
							Computed:    true,
						},
						"other_mails": schema.ListAttribute{
							Description: "A list of other email addresses for the user; for example: ['bob@contoso.com', 'Robert@fabrikam.com']. NOTE: This property can't contain accent characters. Returned only on $select. Supports $filter (eq, not, ge, le, in, startsWith, endsWith, /$count eq 0, /$count ne 0).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"password_policies": schema.StringAttribute{
							Description: "Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two might be specified together; for example: DisablePasswordExpiration, DisableStrongPassword. Returned only on $select. For more information on the default password policies, see Microsoft Entra password policies. Supports $filter (ne, not, and eq on null values).",
							Computed:    true,
						},
						"password_profile": schema.SingleNestedAttribute{
							Description: "Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). To update this property:  In delegated access, the calling app must be assigned the Directory.AccessAsUser.All delegated permission on behalf of the signed-in user.  In application-only access, the calling app must be assigned the User.ReadWrite.All (least privilege) or Directory.ReadWrite.All (higher privilege) application permission and at least the User Administrator Microsoft Entra role.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"force_change_password_next_sign_in": schema.BoolAttribute{
									Description: "true if the user must change their password on the next sign-in; otherwise false.",
									Computed:    true,
								},
								"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
									Description: "If true, at next sign-in, the user must perform a multifactor authentication (MFA) before being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that the user is required to first perform a multifactor authentication before password change. After a password change, this property will be automatically reset to false. If not set, default is false.",
									Computed:    true,
								},
								"password": schema.StringAttribute{
									Description: "The password for the user. This property is required when a user is created. It can be updated, but the user will be required to change the password on the next sign-in. The password must satisfy minimum requirements as specified by the user's passwordPolicies property. By default, a strong password is required.",
									Computed:    true,
								},
							},
						},
						"postal_code": schema.StringAttribute{
							Description: "The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code. Maximum length is 40 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"preferred_data_location": schema.StringAttribute{
							Description: "The preferred data location for the user. For more information, see OneDrive Online Multi-Geo.",
							Computed:    true,
						},
						"preferred_language": schema.StringAttribute{
							Description: "The preferred language for the user. The preferred language format is based on RFC 4646. The name is a combination of an ISO 639 two-letter lowercase culture code associated with the language, and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: 'en-US', or 'es-ES'. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values)",
							Computed:    true,
						},
						"provisioned_plans": schema.ListNestedAttribute{
							Description: "The plans that are provisioned for the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"capability_status": schema.StringAttribute{
										Description: "For example, 'Enabled'.",
										Computed:    true,
									},
									"provisioning_status": schema.StringAttribute{
										Description: "For example, 'Success'.",
										Computed:    true,
									},
									"service": schema.StringAttribute{
										Description: "The name of the service; for example, 'AccessControlS2S'",
										Computed:    true,
									},
								},
							},
						},
						"proxy_addresses": schema.ListAttribute{
							Description: "For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. Changes to the mail property update this collection to include the value as an SMTP address. For more information, see mail and proxyAddresses properties. The proxy address prefixed with SMTP (capitalized) is the primary proxy address, while those addresses prefixed with smtp are the secondary proxy addresses. For Azure AD B2C accounts, this property has a limit of 10 unique addresses. Read-only in Microsoft Graph; you can update this property only through the Microsoft 365 admin center. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith, endsWith, /$count eq 0, /$count ne 0).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"security_identifier": schema.StringAttribute{
							Description: "Security identifier (SID) of the user, used in Windows scenarios. Read-only. Returned by default. Supports $select and $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
						},
						"service_provisioning_errors": schema.ListNestedAttribute{
							Description: "Errors published by a federated service describing a nontransient, service-specific error regarding the properties or link from a user object.  Supports $filter (eq, not, for isResolved and serviceInstance).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_date_time": schema.StringAttribute{
										Description: "The date and time at which the error occurred.",
										Computed:    true,
									},
									"is_resolved": schema.BoolAttribute{
										Description: "Indicates whether the error has been attended to.",
										Computed:    true,
									},
									"service_instance": schema.StringAttribute{
										Description: "Qualified service instance (for example, 'SharePoint/Dublin') that published the service error information.",
										Computed:    true,
									},
								},
							},
						},
						"show_in_address_list": schema.BoolAttribute{
							Description: "Do not use in Microsoft Graph. Manage this property through the Microsoft 365 admin center instead. Represents whether the user should be included in the Outlook global address list. See Known issue.",
							Computed:    true,
						},
						"sign_in_activity": schema.SingleNestedAttribute{
							Description: "Get the last signed-in date and request ID of the sign-in for a given user. Read-only.Returned only on $select. Supports $filter (eq, ne, not, ge, le) but not with any other filterable properties. Note: Details for this property require a Microsoft Entra ID P1 or P2 license and the AuditLog.Read.All permission.This property isn't returned for a user who never signed in or last signed in before April 2020.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"last_non_interactive_sign_in_date_time": schema.StringAttribute{
									Description: "The last non-interactive sign-in date for a specific user. You can use this field to calculate the last time a client attempted (either successfully or unsuccessfully) to sign in to the directory on behalf of a user. Because some users may use clients to access tenant resources rather than signing into your tenant directly, you can use the non-interactive sign-in date to along with lastSignInDateTime to identify inactive users. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Microsoft Entra ID maintains non-interactive sign-ins going back to May 2020. For more information about using the value of this property, see Manage inactive user accounts in Microsoft Entra ID.",
									Computed:    true,
								},
								"last_non_interactive_sign_in_request_id": schema.StringAttribute{
									Description: "Request identifier of the last non-interactive sign-in performed by this user.",
									Computed:    true,
								},
								"last_sign_in_date_time": schema.StringAttribute{
									Description: "The last interactive sign-in date and time for a specific user. You can use this field to calculate the last time a user attempted (either successfully or unsuccessfully) to sign in to the directory with an interactive authentication method. This field can be used to build reports, such as inactive users. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Microsoft Entra ID maintains interactive sign-ins going back to April 2020. For more information about using the value of this property, see Manage inactive user accounts in Microsoft Entra ID.",
									Computed:    true,
								},
								"last_sign_in_request_id": schema.StringAttribute{
									Description: "Request identifier of the last interactive sign-in performed by this user.",
									Computed:    true,
								},
								"last_successful_sign_in_date_time": schema.StringAttribute{
									Description: "The date and time of the user's most recent successful sign-in activity. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
									Computed:    true,
								},
								"last_successful_sign_in_request_id": schema.StringAttribute{
									Description: "The request ID of the last successful sign-in.",
									Computed:    true,
								},
							},
						},
						"sign_in_sessions_valid_from_date_time": schema.StringAttribute{
							Description: "Any refresh tokens or session tokens (session cookies) issued before this time are invalid. Applications get an error when using an invalid refresh or session token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application needs to acquire a new refresh token by requesting the authorized endpoint. Read-only. Use revokeSignInSessions to reset. Returned only on $select.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "The state or province in the user's address. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"street_address": schema.StringAttribute{
							Description: "The street address of the user's place of business. Maximum length is 1,024 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"surname": schema.StringAttribute{
							Description: "The user's surname (family name or last name). Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"usage_location": schema.StringAttribute{
							Description: "A two-letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirements to check for availability of services in countries. Examples include: US, JP, and GB. Not nullable. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"user_principal_name": schema.StringAttribute{
							Description: "The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this value should map to the user's email name. The general format is alias@domain, where the domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization.NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith) and $orderby.",
							Computed:    true,
						},
						"user_type": schema.StringAttribute{
							Description: "A string value that can be used to classify user types in your directory. The possible values are Member and Guest. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). NOTE: For more information about the permissions for members and guests, see What are the default user permissions in Microsoft Entra ID?",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfStateUsers usersModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateUsers)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
			Select: []string{
				"value",
			},
		},
	}

	response, err := d.client.Users().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting users",
			err.Error(),
		)
		return
	}

	if len(response.GetValue()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseValue := range response.GetValue() {
			tfStateUser := usersUserModel{}

			if responseValue.GetId() != nil {
				tfStateUser.Id = types.StringValue(*responseValue.GetId())
			} else {
				tfStateUser.Id = types.StringNull()
			}
			if responseValue.GetDeletedDateTime() != nil {
				tfStateUser.DeletedDateTime = types.StringValue(responseValue.GetDeletedDateTime().String())
			} else {
				tfStateUser.DeletedDateTime = types.StringNull()
			}
			if responseValue.GetAccountEnabled() != nil {
				tfStateUser.AccountEnabled = types.BoolValue(*responseValue.GetAccountEnabled())
			} else {
				tfStateUser.AccountEnabled = types.BoolNull()
			}
			if responseValue.GetAgeGroup() != nil {
				tfStateUser.AgeGroup = types.StringValue(*responseValue.GetAgeGroup())
			} else {
				tfStateUser.AgeGroup = types.StringNull()
			}
			if len(responseValue.GetAssignedLicenses()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseAssignedLicenses := range responseValue.GetAssignedLicenses() {
					tfStateAssignedLicense := usersAssignedLicenseModel{}

					if len(responseAssignedLicenses.GetDisabledPlans()) > 0 {
						var valueArrayDisabledPlans []attr.Value
						for _, responseDisabledPlans := range responseAssignedLicenses.GetDisabledPlans() {
							valueArrayDisabledPlans = append(valueArrayDisabledPlans, types.StringValue(responseDisabledPlans.String()))
						}
						tfStateAssignedLicense.DisabledPlans, _ = types.ListValue(types.StringType, valueArrayDisabledPlans)
					} else {
						tfStateAssignedLicense.DisabledPlans = types.ListNull(types.StringType)
					}
					if responseAssignedLicenses.GetSkuId() != nil {
						tfStateAssignedLicense.SkuId = types.StringValue(responseAssignedLicenses.GetSkuId().String())
					} else {
						tfStateAssignedLicense.SkuId = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedLicense.AttributeTypes(), tfStateAssignedLicense)
					objectValues = append(objectValues, objectValue)
				}
				tfStateUser.AssignedLicenses, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if len(responseValue.GetAssignedPlans()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseAssignedPlans := range responseValue.GetAssignedPlans() {
					tfStateAssignedPlan := usersAssignedPlanModel{}

					if responseAssignedPlans.GetAssignedDateTime() != nil {
						tfStateAssignedPlan.AssignedDateTime = types.StringValue(responseAssignedPlans.GetAssignedDateTime().String())
					} else {
						tfStateAssignedPlan.AssignedDateTime = types.StringNull()
					}
					if responseAssignedPlans.GetCapabilityStatus() != nil {
						tfStateAssignedPlan.CapabilityStatus = types.StringValue(*responseAssignedPlans.GetCapabilityStatus())
					} else {
						tfStateAssignedPlan.CapabilityStatus = types.StringNull()
					}
					if responseAssignedPlans.GetService() != nil {
						tfStateAssignedPlan.Service = types.StringValue(*responseAssignedPlans.GetService())
					} else {
						tfStateAssignedPlan.Service = types.StringNull()
					}
					if responseAssignedPlans.GetServicePlanId() != nil {
						tfStateAssignedPlan.ServicePlanId = types.StringValue(responseAssignedPlans.GetServicePlanId().String())
					} else {
						tfStateAssignedPlan.ServicePlanId = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedPlan.AttributeTypes(), tfStateAssignedPlan)
					objectValues = append(objectValues, objectValue)
				}
				tfStateUser.AssignedPlans, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if responseValue.GetAuthorizationInfo() != nil {
				tfStateAuthorizationInfo := usersAuthorizationInfoModel{}

				if len(responseValue.GetAuthorizationInfo().GetCertificateUserIds()) > 0 {
					var valueArrayCertificateUserIds []attr.Value
					for _, responseCertificateUserIds := range responseValue.GetAuthorizationInfo().GetCertificateUserIds() {
						valueArrayCertificateUserIds = append(valueArrayCertificateUserIds, types.StringValue(responseCertificateUserIds))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayCertificateUserIds)
					tfStateAuthorizationInfo.CertificateUserIds = listValue
				} else {
					tfStateAuthorizationInfo.CertificateUserIds = types.ListNull(types.StringType)
				}

				tfStateUser.AuthorizationInfo, _ = types.ObjectValueFrom(ctx, tfStateAuthorizationInfo.AttributeTypes(), tfStateAuthorizationInfo)
			}
			if len(responseValue.GetBusinessPhones()) > 0 {
				var valueArrayBusinessPhones []attr.Value
				for _, responseBusinessPhones := range responseValue.GetBusinessPhones() {
					valueArrayBusinessPhones = append(valueArrayBusinessPhones, types.StringValue(responseBusinessPhones))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayBusinessPhones)
				tfStateUser.BusinessPhones = listValue
			} else {
				tfStateUser.BusinessPhones = types.ListNull(types.StringType)
			}
			if responseValue.GetCity() != nil {
				tfStateUser.City = types.StringValue(*responseValue.GetCity())
			} else {
				tfStateUser.City = types.StringNull()
			}
			if responseValue.GetCompanyName() != nil {
				tfStateUser.CompanyName = types.StringValue(*responseValue.GetCompanyName())
			} else {
				tfStateUser.CompanyName = types.StringNull()
			}
			if responseValue.GetConsentProvidedForMinor() != nil {
				tfStateUser.ConsentProvidedForMinor = types.StringValue(*responseValue.GetConsentProvidedForMinor())
			} else {
				tfStateUser.ConsentProvidedForMinor = types.StringNull()
			}
			if responseValue.GetCountry() != nil {
				tfStateUser.Country = types.StringValue(*responseValue.GetCountry())
			} else {
				tfStateUser.Country = types.StringNull()
			}
			if responseValue.GetCreatedDateTime() != nil {
				tfStateUser.CreatedDateTime = types.StringValue(responseValue.GetCreatedDateTime().String())
			} else {
				tfStateUser.CreatedDateTime = types.StringNull()
			}
			if responseValue.GetCreationType() != nil {
				tfStateUser.CreationType = types.StringValue(*responseValue.GetCreationType())
			} else {
				tfStateUser.CreationType = types.StringNull()
			}
			if responseValue.GetCustomSecurityAttributes() != nil {
				tfStateCustomSecurityAttributeValue := usersCustomSecurityAttributeValueModel{}

				tfStateUser.CustomSecurityAttributes, _ = types.ObjectValueFrom(ctx, tfStateCustomSecurityAttributeValue.AttributeTypes(), tfStateCustomSecurityAttributeValue)
			}
			if responseValue.GetDepartment() != nil {
				tfStateUser.Department = types.StringValue(*responseValue.GetDepartment())
			} else {
				tfStateUser.Department = types.StringNull()
			}
			if responseValue.GetDisplayName() != nil {
				tfStateUser.DisplayName = types.StringValue(*responseValue.GetDisplayName())
			} else {
				tfStateUser.DisplayName = types.StringNull()
			}
			if responseValue.GetEmployeeHireDate() != nil {
				tfStateUser.EmployeeHireDate = types.StringValue(responseValue.GetEmployeeHireDate().String())
			} else {
				tfStateUser.EmployeeHireDate = types.StringNull()
			}
			if responseValue.GetEmployeeId() != nil {
				tfStateUser.EmployeeId = types.StringValue(*responseValue.GetEmployeeId())
			} else {
				tfStateUser.EmployeeId = types.StringNull()
			}
			if responseValue.GetEmployeeLeaveDateTime() != nil {
				tfStateUser.EmployeeLeaveDateTime = types.StringValue(responseValue.GetEmployeeLeaveDateTime().String())
			} else {
				tfStateUser.EmployeeLeaveDateTime = types.StringNull()
			}
			if responseValue.GetEmployeeOrgData() != nil {
				tfStateEmployeeOrgData := usersEmployeeOrgDataModel{}

				if responseValue.GetEmployeeOrgData().GetCostCenter() != nil {
					tfStateEmployeeOrgData.CostCenter = types.StringValue(*responseValue.GetEmployeeOrgData().GetCostCenter())
				} else {
					tfStateEmployeeOrgData.CostCenter = types.StringNull()
				}
				if responseValue.GetEmployeeOrgData().GetDivision() != nil {
					tfStateEmployeeOrgData.Division = types.StringValue(*responseValue.GetEmployeeOrgData().GetDivision())
				} else {
					tfStateEmployeeOrgData.Division = types.StringNull()
				}

				tfStateUser.EmployeeOrgData, _ = types.ObjectValueFrom(ctx, tfStateEmployeeOrgData.AttributeTypes(), tfStateEmployeeOrgData)
			}
			if responseValue.GetEmployeeType() != nil {
				tfStateUser.EmployeeType = types.StringValue(*responseValue.GetEmployeeType())
			} else {
				tfStateUser.EmployeeType = types.StringNull()
			}
			if responseValue.GetExternalUserState() != nil {
				tfStateUser.ExternalUserState = types.StringValue(*responseValue.GetExternalUserState())
			} else {
				tfStateUser.ExternalUserState = types.StringNull()
			}
			if responseValue.GetExternalUserStateChangeDateTime() != nil {
				tfStateUser.ExternalUserStateChangeDateTime = types.StringValue(responseValue.GetExternalUserStateChangeDateTime().String())
			} else {
				tfStateUser.ExternalUserStateChangeDateTime = types.StringNull()
			}
			if responseValue.GetFaxNumber() != nil {
				tfStateUser.FaxNumber = types.StringValue(*responseValue.GetFaxNumber())
			} else {
				tfStateUser.FaxNumber = types.StringNull()
			}
			if responseValue.GetGivenName() != nil {
				tfStateUser.GivenName = types.StringValue(*responseValue.GetGivenName())
			} else {
				tfStateUser.GivenName = types.StringNull()
			}
			if len(responseValue.GetIdentities()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseIdentities := range responseValue.GetIdentities() {
					tfStateObjectIdentity := usersObjectIdentityModel{}

					if responseIdentities.GetIssuer() != nil {
						tfStateObjectIdentity.Issuer = types.StringValue(*responseIdentities.GetIssuer())
					} else {
						tfStateObjectIdentity.Issuer = types.StringNull()
					}
					if responseIdentities.GetIssuerAssignedId() != nil {
						tfStateObjectIdentity.IssuerAssignedId = types.StringValue(*responseIdentities.GetIssuerAssignedId())
					} else {
						tfStateObjectIdentity.IssuerAssignedId = types.StringNull()
					}
					if responseIdentities.GetSignInType() != nil {
						tfStateObjectIdentity.SignInType = types.StringValue(*responseIdentities.GetSignInType())
					} else {
						tfStateObjectIdentity.SignInType = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateObjectIdentity.AttributeTypes(), tfStateObjectIdentity)
					objectValues = append(objectValues, objectValue)
				}
				tfStateUser.Identities, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if len(responseValue.GetImAddresses()) > 0 {
				var valueArrayImAddresses []attr.Value
				for _, responseImAddresses := range responseValue.GetImAddresses() {
					valueArrayImAddresses = append(valueArrayImAddresses, types.StringValue(responseImAddresses))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayImAddresses)
				tfStateUser.ImAddresses = listValue
			} else {
				tfStateUser.ImAddresses = types.ListNull(types.StringType)
			}
			if responseValue.GetIsManagementRestricted() != nil {
				tfStateUser.IsManagementRestricted = types.BoolValue(*responseValue.GetIsManagementRestricted())
			} else {
				tfStateUser.IsManagementRestricted = types.BoolNull()
			}
			if responseValue.GetIsResourceAccount() != nil {
				tfStateUser.IsResourceAccount = types.BoolValue(*responseValue.GetIsResourceAccount())
			} else {
				tfStateUser.IsResourceAccount = types.BoolNull()
			}
			if responseValue.GetJobTitle() != nil {
				tfStateUser.JobTitle = types.StringValue(*responseValue.GetJobTitle())
			} else {
				tfStateUser.JobTitle = types.StringNull()
			}
			if responseValue.GetLastPasswordChangeDateTime() != nil {
				tfStateUser.LastPasswordChangeDateTime = types.StringValue(responseValue.GetLastPasswordChangeDateTime().String())
			} else {
				tfStateUser.LastPasswordChangeDateTime = types.StringNull()
			}
			if responseValue.GetLegalAgeGroupClassification() != nil {
				tfStateUser.LegalAgeGroupClassification = types.StringValue(*responseValue.GetLegalAgeGroupClassification())
			} else {
				tfStateUser.LegalAgeGroupClassification = types.StringNull()
			}
			if len(responseValue.GetLicenseAssignmentStates()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseLicenseAssignmentStates := range responseValue.GetLicenseAssignmentStates() {
					tfStateLicenseAssignmentState := usersLicenseAssignmentStateModel{}

					if responseLicenseAssignmentStates.GetAssignedByGroup() != nil {
						tfStateLicenseAssignmentState.AssignedByGroup = types.StringValue(*responseLicenseAssignmentStates.GetAssignedByGroup())
					} else {
						tfStateLicenseAssignmentState.AssignedByGroup = types.StringNull()
					}
					if len(responseLicenseAssignmentStates.GetDisabledPlans()) > 0 {
						var valueArrayDisabledPlans []attr.Value
						for _, responseDisabledPlans := range responseLicenseAssignmentStates.GetDisabledPlans() {
							valueArrayDisabledPlans = append(valueArrayDisabledPlans, types.StringValue(responseDisabledPlans.String()))
						}
						tfStateLicenseAssignmentState.DisabledPlans, _ = types.ListValue(types.StringType, valueArrayDisabledPlans)
					} else {
						tfStateLicenseAssignmentState.DisabledPlans = types.ListNull(types.StringType)
					}
					if responseLicenseAssignmentStates.GetError() != nil {
						tfStateLicenseAssignmentState.Error = types.StringValue(*responseLicenseAssignmentStates.GetError())
					} else {
						tfStateLicenseAssignmentState.Error = types.StringNull()
					}
					if responseLicenseAssignmentStates.GetLastUpdatedDateTime() != nil {
						tfStateLicenseAssignmentState.LastUpdatedDateTime = types.StringValue(responseLicenseAssignmentStates.GetLastUpdatedDateTime().String())
					} else {
						tfStateLicenseAssignmentState.LastUpdatedDateTime = types.StringNull()
					}
					if responseLicenseAssignmentStates.GetSkuId() != nil {
						tfStateLicenseAssignmentState.SkuId = types.StringValue(responseLicenseAssignmentStates.GetSkuId().String())
					} else {
						tfStateLicenseAssignmentState.SkuId = types.StringNull()
					}
					if responseLicenseAssignmentStates.GetState() != nil {
						tfStateLicenseAssignmentState.State = types.StringValue(*responseLicenseAssignmentStates.GetState())
					} else {
						tfStateLicenseAssignmentState.State = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateLicenseAssignmentState.AttributeTypes(), tfStateLicenseAssignmentState)
					objectValues = append(objectValues, objectValue)
				}
				tfStateUser.LicenseAssignmentStates, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if responseValue.GetMail() != nil {
				tfStateUser.Mail = types.StringValue(*responseValue.GetMail())
			} else {
				tfStateUser.Mail = types.StringNull()
			}
			if responseValue.GetMailNickname() != nil {
				tfStateUser.MailNickname = types.StringValue(*responseValue.GetMailNickname())
			} else {
				tfStateUser.MailNickname = types.StringNull()
			}
			if responseValue.GetMobilePhone() != nil {
				tfStateUser.MobilePhone = types.StringValue(*responseValue.GetMobilePhone())
			} else {
				tfStateUser.MobilePhone = types.StringNull()
			}
			if responseValue.GetOfficeLocation() != nil {
				tfStateUser.OfficeLocation = types.StringValue(*responseValue.GetOfficeLocation())
			} else {
				tfStateUser.OfficeLocation = types.StringNull()
			}
			if responseValue.GetOnPremisesDistinguishedName() != nil {
				tfStateUser.OnPremisesDistinguishedName = types.StringValue(*responseValue.GetOnPremisesDistinguishedName())
			} else {
				tfStateUser.OnPremisesDistinguishedName = types.StringNull()
			}
			if responseValue.GetOnPremisesDomainName() != nil {
				tfStateUser.OnPremisesDomainName = types.StringValue(*responseValue.GetOnPremisesDomainName())
			} else {
				tfStateUser.OnPremisesDomainName = types.StringNull()
			}
			if responseValue.GetOnPremisesExtensionAttributes() != nil {
				tfStateOnPremisesExtensionAttributes := usersOnPremisesExtensionAttributesModel{}

				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute1() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute1())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute10() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute10())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute11() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute11())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute12() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute12())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute13() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute13())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute14() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute14())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute15() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute15())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute2() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute2())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute3() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute3())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute4() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute4())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute5() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute5())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute6() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute6())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute7() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute7())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute8() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute8())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringNull()
				}
				if responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute9() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringValue(*responseValue.GetOnPremisesExtensionAttributes().GetExtensionAttribute9())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringNull()
				}

				tfStateUser.OnPremisesExtensionAttributes, _ = types.ObjectValueFrom(ctx, tfStateOnPremisesExtensionAttributes.AttributeTypes(), tfStateOnPremisesExtensionAttributes)
			}
			if responseValue.GetOnPremisesImmutableId() != nil {
				tfStateUser.OnPremisesImmutableId = types.StringValue(*responseValue.GetOnPremisesImmutableId())
			} else {
				tfStateUser.OnPremisesImmutableId = types.StringNull()
			}
			if responseValue.GetOnPremisesLastSyncDateTime() != nil {
				tfStateUser.OnPremisesLastSyncDateTime = types.StringValue(responseValue.GetOnPremisesLastSyncDateTime().String())
			} else {
				tfStateUser.OnPremisesLastSyncDateTime = types.StringNull()
			}
			if len(responseValue.GetOnPremisesProvisioningErrors()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseOnPremisesProvisioningErrors := range responseValue.GetOnPremisesProvisioningErrors() {
					tfStateOnPremisesProvisioningError := usersOnPremisesProvisioningErrorModel{}

					if responseOnPremisesProvisioningErrors.GetCategory() != nil {
						tfStateOnPremisesProvisioningError.Category = types.StringValue(*responseOnPremisesProvisioningErrors.GetCategory())
					} else {
						tfStateOnPremisesProvisioningError.Category = types.StringNull()
					}
					if responseOnPremisesProvisioningErrors.GetOccurredDateTime() != nil {
						tfStateOnPremisesProvisioningError.OccurredDateTime = types.StringValue(responseOnPremisesProvisioningErrors.GetOccurredDateTime().String())
					} else {
						tfStateOnPremisesProvisioningError.OccurredDateTime = types.StringNull()
					}
					if responseOnPremisesProvisioningErrors.GetPropertyCausingError() != nil {
						tfStateOnPremisesProvisioningError.PropertyCausingError = types.StringValue(*responseOnPremisesProvisioningErrors.GetPropertyCausingError())
					} else {
						tfStateOnPremisesProvisioningError.PropertyCausingError = types.StringNull()
					}
					if responseOnPremisesProvisioningErrors.GetValue() != nil {
						tfStateOnPremisesProvisioningError.Value = types.StringValue(*responseOnPremisesProvisioningErrors.GetValue())
					} else {
						tfStateOnPremisesProvisioningError.Value = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateOnPremisesProvisioningError.AttributeTypes(), tfStateOnPremisesProvisioningError)
					objectValues = append(objectValues, objectValue)
				}
				tfStateUser.OnPremisesProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if responseValue.GetOnPremisesSamAccountName() != nil {
				tfStateUser.OnPremisesSamAccountName = types.StringValue(*responseValue.GetOnPremisesSamAccountName())
			} else {
				tfStateUser.OnPremisesSamAccountName = types.StringNull()
			}
			if responseValue.GetOnPremisesSecurityIdentifier() != nil {
				tfStateUser.OnPremisesSecurityIdentifier = types.StringValue(*responseValue.GetOnPremisesSecurityIdentifier())
			} else {
				tfStateUser.OnPremisesSecurityIdentifier = types.StringNull()
			}
			if responseValue.GetOnPremisesSyncEnabled() != nil {
				tfStateUser.OnPremisesSyncEnabled = types.BoolValue(*responseValue.GetOnPremisesSyncEnabled())
			} else {
				tfStateUser.OnPremisesSyncEnabled = types.BoolNull()
			}
			if responseValue.GetOnPremisesUserPrincipalName() != nil {
				tfStateUser.OnPremisesUserPrincipalName = types.StringValue(*responseValue.GetOnPremisesUserPrincipalName())
			} else {
				tfStateUser.OnPremisesUserPrincipalName = types.StringNull()
			}
			if len(responseValue.GetOtherMails()) > 0 {
				var valueArrayOtherMails []attr.Value
				for _, responseOtherMails := range responseValue.GetOtherMails() {
					valueArrayOtherMails = append(valueArrayOtherMails, types.StringValue(responseOtherMails))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayOtherMails)
				tfStateUser.OtherMails = listValue
			} else {
				tfStateUser.OtherMails = types.ListNull(types.StringType)
			}
			if responseValue.GetPasswordPolicies() != nil {
				tfStateUser.PasswordPolicies = types.StringValue(*responseValue.GetPasswordPolicies())
			} else {
				tfStateUser.PasswordPolicies = types.StringNull()
			}
			if responseValue.GetPasswordProfile() != nil {
				tfStatePasswordProfile := usersPasswordProfileModel{}

				if responseValue.GetPasswordProfile().GetForceChangePasswordNextSignIn() != nil {
					tfStatePasswordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*responseValue.GetPasswordProfile().GetForceChangePasswordNextSignIn())
				} else {
					tfStatePasswordProfile.ForceChangePasswordNextSignIn = types.BoolNull()
				}
				if responseValue.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa() != nil {
					tfStatePasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*responseValue.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
				} else {
					tfStatePasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolNull()
				}
				if responseValue.GetPasswordProfile().GetPassword() != nil {
					tfStatePasswordProfile.Password = types.StringValue(*responseValue.GetPasswordProfile().GetPassword())
				} else {
					tfStatePasswordProfile.Password = types.StringNull()
				}

				tfStateUser.PasswordProfile, _ = types.ObjectValueFrom(ctx, tfStatePasswordProfile.AttributeTypes(), tfStatePasswordProfile)
			}
			if responseValue.GetPostalCode() != nil {
				tfStateUser.PostalCode = types.StringValue(*responseValue.GetPostalCode())
			} else {
				tfStateUser.PostalCode = types.StringNull()
			}
			if responseValue.GetPreferredDataLocation() != nil {
				tfStateUser.PreferredDataLocation = types.StringValue(*responseValue.GetPreferredDataLocation())
			} else {
				tfStateUser.PreferredDataLocation = types.StringNull()
			}
			if responseValue.GetPreferredLanguage() != nil {
				tfStateUser.PreferredLanguage = types.StringValue(*responseValue.GetPreferredLanguage())
			} else {
				tfStateUser.PreferredLanguage = types.StringNull()
			}
			if len(responseValue.GetProvisionedPlans()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseProvisionedPlans := range responseValue.GetProvisionedPlans() {
					tfStateProvisionedPlan := usersProvisionedPlanModel{}

					if responseProvisionedPlans.GetCapabilityStatus() != nil {
						tfStateProvisionedPlan.CapabilityStatus = types.StringValue(*responseProvisionedPlans.GetCapabilityStatus())
					} else {
						tfStateProvisionedPlan.CapabilityStatus = types.StringNull()
					}
					if responseProvisionedPlans.GetProvisioningStatus() != nil {
						tfStateProvisionedPlan.ProvisioningStatus = types.StringValue(*responseProvisionedPlans.GetProvisioningStatus())
					} else {
						tfStateProvisionedPlan.ProvisioningStatus = types.StringNull()
					}
					if responseProvisionedPlans.GetService() != nil {
						tfStateProvisionedPlan.Service = types.StringValue(*responseProvisionedPlans.GetService())
					} else {
						tfStateProvisionedPlan.Service = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateProvisionedPlan.AttributeTypes(), tfStateProvisionedPlan)
					objectValues = append(objectValues, objectValue)
				}
				tfStateUser.ProvisionedPlans, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if len(responseValue.GetProxyAddresses()) > 0 {
				var valueArrayProxyAddresses []attr.Value
				for _, responseProxyAddresses := range responseValue.GetProxyAddresses() {
					valueArrayProxyAddresses = append(valueArrayProxyAddresses, types.StringValue(responseProxyAddresses))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayProxyAddresses)
				tfStateUser.ProxyAddresses = listValue
			} else {
				tfStateUser.ProxyAddresses = types.ListNull(types.StringType)
			}
			if responseValue.GetSecurityIdentifier() != nil {
				tfStateUser.SecurityIdentifier = types.StringValue(*responseValue.GetSecurityIdentifier())
			} else {
				tfStateUser.SecurityIdentifier = types.StringNull()
			}
			if len(responseValue.GetServiceProvisioningErrors()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseServiceProvisioningErrors := range responseValue.GetServiceProvisioningErrors() {
					tfStateServiceProvisioningError := usersServiceProvisioningErrorModel{}

					if responseServiceProvisioningErrors.GetCreatedDateTime() != nil {
						tfStateServiceProvisioningError.CreatedDateTime = types.StringValue(responseServiceProvisioningErrors.GetCreatedDateTime().String())
					} else {
						tfStateServiceProvisioningError.CreatedDateTime = types.StringNull()
					}
					if responseServiceProvisioningErrors.GetIsResolved() != nil {
						tfStateServiceProvisioningError.IsResolved = types.BoolValue(*responseServiceProvisioningErrors.GetIsResolved())
					} else {
						tfStateServiceProvisioningError.IsResolved = types.BoolNull()
					}
					if responseServiceProvisioningErrors.GetServiceInstance() != nil {
						tfStateServiceProvisioningError.ServiceInstance = types.StringValue(*responseServiceProvisioningErrors.GetServiceInstance())
					} else {
						tfStateServiceProvisioningError.ServiceInstance = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateServiceProvisioningError.AttributeTypes(), tfStateServiceProvisioningError)
					objectValues = append(objectValues, objectValue)
				}
				tfStateUser.ServiceProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if responseValue.GetShowInAddressList() != nil {
				tfStateUser.ShowInAddressList = types.BoolValue(*responseValue.GetShowInAddressList())
			} else {
				tfStateUser.ShowInAddressList = types.BoolNull()
			}
			if responseValue.GetSignInActivity() != nil {
				tfStateSignInActivity := usersSignInActivityModel{}

				if responseValue.GetSignInActivity().GetLastNonInteractiveSignInDateTime() != nil {
					tfStateSignInActivity.LastNonInteractiveSignInDateTime = types.StringValue(responseValue.GetSignInActivity().GetLastNonInteractiveSignInDateTime().String())
				} else {
					tfStateSignInActivity.LastNonInteractiveSignInDateTime = types.StringNull()
				}
				if responseValue.GetSignInActivity().GetLastNonInteractiveSignInRequestId() != nil {
					tfStateSignInActivity.LastNonInteractiveSignInRequestId = types.StringValue(*responseValue.GetSignInActivity().GetLastNonInteractiveSignInRequestId())
				} else {
					tfStateSignInActivity.LastNonInteractiveSignInRequestId = types.StringNull()
				}
				if responseValue.GetSignInActivity().GetLastSignInDateTime() != nil {
					tfStateSignInActivity.LastSignInDateTime = types.StringValue(responseValue.GetSignInActivity().GetLastSignInDateTime().String())
				} else {
					tfStateSignInActivity.LastSignInDateTime = types.StringNull()
				}
				if responseValue.GetSignInActivity().GetLastSignInRequestId() != nil {
					tfStateSignInActivity.LastSignInRequestId = types.StringValue(*responseValue.GetSignInActivity().GetLastSignInRequestId())
				} else {
					tfStateSignInActivity.LastSignInRequestId = types.StringNull()
				}
				if responseValue.GetSignInActivity().GetLastSuccessfulSignInDateTime() != nil {
					tfStateSignInActivity.LastSuccessfulSignInDateTime = types.StringValue(responseValue.GetSignInActivity().GetLastSuccessfulSignInDateTime().String())
				} else {
					tfStateSignInActivity.LastSuccessfulSignInDateTime = types.StringNull()
				}
				if responseValue.GetSignInActivity().GetLastSuccessfulSignInRequestId() != nil {
					tfStateSignInActivity.LastSuccessfulSignInRequestId = types.StringValue(*responseValue.GetSignInActivity().GetLastSuccessfulSignInRequestId())
				} else {
					tfStateSignInActivity.LastSuccessfulSignInRequestId = types.StringNull()
				}

				tfStateUser.SignInActivity, _ = types.ObjectValueFrom(ctx, tfStateSignInActivity.AttributeTypes(), tfStateSignInActivity)
			}
			if responseValue.GetSignInSessionsValidFromDateTime() != nil {
				tfStateUser.SignInSessionsValidFromDateTime = types.StringValue(responseValue.GetSignInSessionsValidFromDateTime().String())
			} else {
				tfStateUser.SignInSessionsValidFromDateTime = types.StringNull()
			}
			if responseValue.GetState() != nil {
				tfStateUser.State = types.StringValue(*responseValue.GetState())
			} else {
				tfStateUser.State = types.StringNull()
			}
			if responseValue.GetStreetAddress() != nil {
				tfStateUser.StreetAddress = types.StringValue(*responseValue.GetStreetAddress())
			} else {
				tfStateUser.StreetAddress = types.StringNull()
			}
			if responseValue.GetSurname() != nil {
				tfStateUser.Surname = types.StringValue(*responseValue.GetSurname())
			} else {
				tfStateUser.Surname = types.StringNull()
			}
			if responseValue.GetUsageLocation() != nil {
				tfStateUser.UsageLocation = types.StringValue(*responseValue.GetUsageLocation())
			} else {
				tfStateUser.UsageLocation = types.StringNull()
			}
			if responseValue.GetUserPrincipalName() != nil {
				tfStateUser.UserPrincipalName = types.StringValue(*responseValue.GetUserPrincipalName())
			} else {
				tfStateUser.UserPrincipalName = types.StringNull()
			}
			if responseValue.GetUserType() != nil {
				tfStateUser.UserType = types.StringValue(*responseValue.GetUserType())
			} else {
				tfStateUser.UserType = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateUser.AttributeTypes(), tfStateUser)
			objectValues = append(objectValues, objectValue)
		}
		tfStateUsers.Value, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateUsers)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
