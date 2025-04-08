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

	result, err := d.client.Users().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting users",
			err.Error(),
		)
		return
	}

	if len(result.GetValue()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetValue() {
			tfStateValue := usersUserModel{}

			if v.GetId() != nil {
				tfStateValue.Id = types.StringValue(*v.GetId())
			} else {
				tfStateValue.Id = types.StringNull()
			}
			if v.GetDeletedDateTime() != nil {
				tfStateValue.DeletedDateTime = types.StringValue(v.GetDeletedDateTime().String())
			} else {
				tfStateValue.DeletedDateTime = types.StringNull()
			}
			if v.GetAccountEnabled() != nil {
				tfStateValue.AccountEnabled = types.BoolValue(*v.GetAccountEnabled())
			} else {
				tfStateValue.AccountEnabled = types.BoolNull()
			}
			if v.GetAgeGroup() != nil {
				tfStateValue.AgeGroup = types.StringValue(*v.GetAgeGroup())
			} else {
				tfStateValue.AgeGroup = types.StringNull()
			}
			if len(v.GetAssignedLicenses()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetAssignedLicenses() {
					tfStateAssignedLicenses := usersAssignedLicenseModel{}

					if len(v.GetDisabledPlans()) > 0 {
						var valueArrayDisabledPlans []attr.Value
						for _, v := range v.GetDisabledPlans() {
							valueArrayDisabledPlans = append(valueArrayDisabledPlans, types.StringValue(v.String()))
						}
						tfStateAssignedLicenses.DisabledPlans, _ = types.ListValue(types.StringType, valueArrayDisabledPlans)
					} else {
						tfStateAssignedLicenses.DisabledPlans = types.ListNull(types.StringType)
					}
					if v.GetSkuId() != nil {
						tfStateAssignedLicenses.SkuId = types.StringValue(v.GetSkuId().String())
					} else {
						tfStateAssignedLicenses.SkuId = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedLicenses.AttributeTypes(), tfStateAssignedLicenses)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.AssignedLicenses, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if len(v.GetAssignedPlans()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetAssignedPlans() {
					tfStateAssignedPlans := usersAssignedPlanModel{}

					if v.GetAssignedDateTime() != nil {
						tfStateAssignedPlans.AssignedDateTime = types.StringValue(v.GetAssignedDateTime().String())
					} else {
						tfStateAssignedPlans.AssignedDateTime = types.StringNull()
					}
					if v.GetCapabilityStatus() != nil {
						tfStateAssignedPlans.CapabilityStatus = types.StringValue(*v.GetCapabilityStatus())
					} else {
						tfStateAssignedPlans.CapabilityStatus = types.StringNull()
					}
					if v.GetService() != nil {
						tfStateAssignedPlans.Service = types.StringValue(*v.GetService())
					} else {
						tfStateAssignedPlans.Service = types.StringNull()
					}
					if v.GetServicePlanId() != nil {
						tfStateAssignedPlans.ServicePlanId = types.StringValue(v.GetServicePlanId().String())
					} else {
						tfStateAssignedPlans.ServicePlanId = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedPlans.AttributeTypes(), tfStateAssignedPlans)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.AssignedPlans, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetAuthorizationInfo() != nil {
				tfStateAuthorizationInfo := usersAuthorizationInfoModel{}

				if len(v.GetAuthorizationInfo().GetCertificateUserIds()) > 0 {
					var valueArrayCertificateUserIds []attr.Value
					for _, v := range v.GetAuthorizationInfo().GetCertificateUserIds() {
						valueArrayCertificateUserIds = append(valueArrayCertificateUserIds, types.StringValue(v))
					}
					listValue, _ := types.ListValue(types.StringType, valueArrayCertificateUserIds)
					tfStateAuthorizationInfo.CertificateUserIds = listValue
				} else {
					tfStateAuthorizationInfo.CertificateUserIds = types.ListNull(types.StringType)
				}

				tfStateValue.AuthorizationInfo, _ = types.ObjectValueFrom(ctx, tfStateAuthorizationInfo.AttributeTypes(), tfStateAuthorizationInfo)
			}
			if len(v.GetBusinessPhones()) > 0 {
				var valueArrayBusinessPhones []attr.Value
				for _, v := range v.GetBusinessPhones() {
					valueArrayBusinessPhones = append(valueArrayBusinessPhones, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayBusinessPhones)
				tfStateValue.BusinessPhones = listValue
			} else {
				tfStateValue.BusinessPhones = types.ListNull(types.StringType)
			}
			if v.GetCity() != nil {
				tfStateValue.City = types.StringValue(*v.GetCity())
			} else {
				tfStateValue.City = types.StringNull()
			}
			if v.GetCompanyName() != nil {
				tfStateValue.CompanyName = types.StringValue(*v.GetCompanyName())
			} else {
				tfStateValue.CompanyName = types.StringNull()
			}
			if v.GetConsentProvidedForMinor() != nil {
				tfStateValue.ConsentProvidedForMinor = types.StringValue(*v.GetConsentProvidedForMinor())
			} else {
				tfStateValue.ConsentProvidedForMinor = types.StringNull()
			}
			if v.GetCountry() != nil {
				tfStateValue.Country = types.StringValue(*v.GetCountry())
			} else {
				tfStateValue.Country = types.StringNull()
			}
			if v.GetCreatedDateTime() != nil {
				tfStateValue.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			} else {
				tfStateValue.CreatedDateTime = types.StringNull()
			}
			if v.GetCreationType() != nil {
				tfStateValue.CreationType = types.StringValue(*v.GetCreationType())
			} else {
				tfStateValue.CreationType = types.StringNull()
			}
			if v.GetCustomSecurityAttributes() != nil {
				tfStateCustomSecurityAttributes := usersCustomSecurityAttributeValueModel{}

				tfStateValue.CustomSecurityAttributes, _ = types.ObjectValueFrom(ctx, tfStateCustomSecurityAttributes.AttributeTypes(), tfStateCustomSecurityAttributes)
			}
			if v.GetDepartment() != nil {
				tfStateValue.Department = types.StringValue(*v.GetDepartment())
			} else {
				tfStateValue.Department = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				tfStateValue.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				tfStateValue.DisplayName = types.StringNull()
			}
			if v.GetEmployeeHireDate() != nil {
				tfStateValue.EmployeeHireDate = types.StringValue(v.GetEmployeeHireDate().String())
			} else {
				tfStateValue.EmployeeHireDate = types.StringNull()
			}
			if v.GetEmployeeId() != nil {
				tfStateValue.EmployeeId = types.StringValue(*v.GetEmployeeId())
			} else {
				tfStateValue.EmployeeId = types.StringNull()
			}
			if v.GetEmployeeLeaveDateTime() != nil {
				tfStateValue.EmployeeLeaveDateTime = types.StringValue(v.GetEmployeeLeaveDateTime().String())
			} else {
				tfStateValue.EmployeeLeaveDateTime = types.StringNull()
			}
			if v.GetEmployeeOrgData() != nil {
				tfStateEmployeeOrgData := usersEmployeeOrgDataModel{}

				if v.GetEmployeeOrgData().GetCostCenter() != nil {
					tfStateEmployeeOrgData.CostCenter = types.StringValue(*v.GetEmployeeOrgData().GetCostCenter())
				} else {
					tfStateEmployeeOrgData.CostCenter = types.StringNull()
				}
				if v.GetEmployeeOrgData().GetDivision() != nil {
					tfStateEmployeeOrgData.Division = types.StringValue(*v.GetEmployeeOrgData().GetDivision())
				} else {
					tfStateEmployeeOrgData.Division = types.StringNull()
				}

				tfStateValue.EmployeeOrgData, _ = types.ObjectValueFrom(ctx, tfStateEmployeeOrgData.AttributeTypes(), tfStateEmployeeOrgData)
			}
			if v.GetEmployeeType() != nil {
				tfStateValue.EmployeeType = types.StringValue(*v.GetEmployeeType())
			} else {
				tfStateValue.EmployeeType = types.StringNull()
			}
			if v.GetExternalUserState() != nil {
				tfStateValue.ExternalUserState = types.StringValue(*v.GetExternalUserState())
			} else {
				tfStateValue.ExternalUserState = types.StringNull()
			}
			if v.GetExternalUserStateChangeDateTime() != nil {
				tfStateValue.ExternalUserStateChangeDateTime = types.StringValue(v.GetExternalUserStateChangeDateTime().String())
			} else {
				tfStateValue.ExternalUserStateChangeDateTime = types.StringNull()
			}
			if v.GetFaxNumber() != nil {
				tfStateValue.FaxNumber = types.StringValue(*v.GetFaxNumber())
			} else {
				tfStateValue.FaxNumber = types.StringNull()
			}
			if v.GetGivenName() != nil {
				tfStateValue.GivenName = types.StringValue(*v.GetGivenName())
			} else {
				tfStateValue.GivenName = types.StringNull()
			}
			if len(v.GetIdentities()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetIdentities() {
					tfStateIdentities := usersObjectIdentityModel{}

					if v.GetIssuer() != nil {
						tfStateIdentities.Issuer = types.StringValue(*v.GetIssuer())
					} else {
						tfStateIdentities.Issuer = types.StringNull()
					}
					if v.GetIssuerAssignedId() != nil {
						tfStateIdentities.IssuerAssignedId = types.StringValue(*v.GetIssuerAssignedId())
					} else {
						tfStateIdentities.IssuerAssignedId = types.StringNull()
					}
					if v.GetSignInType() != nil {
						tfStateIdentities.SignInType = types.StringValue(*v.GetSignInType())
					} else {
						tfStateIdentities.SignInType = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateIdentities.AttributeTypes(), tfStateIdentities)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.Identities, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if len(v.GetImAddresses()) > 0 {
				var valueArrayImAddresses []attr.Value
				for _, v := range v.GetImAddresses() {
					valueArrayImAddresses = append(valueArrayImAddresses, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayImAddresses)
				tfStateValue.ImAddresses = listValue
			} else {
				tfStateValue.ImAddresses = types.ListNull(types.StringType)
			}
			if v.GetIsManagementRestricted() != nil {
				tfStateValue.IsManagementRestricted = types.BoolValue(*v.GetIsManagementRestricted())
			} else {
				tfStateValue.IsManagementRestricted = types.BoolNull()
			}
			if v.GetIsResourceAccount() != nil {
				tfStateValue.IsResourceAccount = types.BoolValue(*v.GetIsResourceAccount())
			} else {
				tfStateValue.IsResourceAccount = types.BoolNull()
			}
			if v.GetJobTitle() != nil {
				tfStateValue.JobTitle = types.StringValue(*v.GetJobTitle())
			} else {
				tfStateValue.JobTitle = types.StringNull()
			}
			if v.GetLastPasswordChangeDateTime() != nil {
				tfStateValue.LastPasswordChangeDateTime = types.StringValue(v.GetLastPasswordChangeDateTime().String())
			} else {
				tfStateValue.LastPasswordChangeDateTime = types.StringNull()
			}
			if v.GetLegalAgeGroupClassification() != nil {
				tfStateValue.LegalAgeGroupClassification = types.StringValue(*v.GetLegalAgeGroupClassification())
			} else {
				tfStateValue.LegalAgeGroupClassification = types.StringNull()
			}
			if len(v.GetLicenseAssignmentStates()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetLicenseAssignmentStates() {
					tfStateLicenseAssignmentStates := usersLicenseAssignmentStateModel{}

					if v.GetAssignedByGroup() != nil {
						tfStateLicenseAssignmentStates.AssignedByGroup = types.StringValue(*v.GetAssignedByGroup())
					} else {
						tfStateLicenseAssignmentStates.AssignedByGroup = types.StringNull()
					}
					if len(v.GetDisabledPlans()) > 0 {
						var valueArrayDisabledPlans []attr.Value
						for _, v := range v.GetDisabledPlans() {
							valueArrayDisabledPlans = append(valueArrayDisabledPlans, types.StringValue(v.String()))
						}
						tfStateLicenseAssignmentStates.DisabledPlans, _ = types.ListValue(types.StringType, valueArrayDisabledPlans)
					} else {
						tfStateLicenseAssignmentStates.DisabledPlans = types.ListNull(types.StringType)
					}
					if v.GetError() != nil {
						tfStateLicenseAssignmentStates.Error = types.StringValue(*v.GetError())
					} else {
						tfStateLicenseAssignmentStates.Error = types.StringNull()
					}
					if v.GetLastUpdatedDateTime() != nil {
						tfStateLicenseAssignmentStates.LastUpdatedDateTime = types.StringValue(v.GetLastUpdatedDateTime().String())
					} else {
						tfStateLicenseAssignmentStates.LastUpdatedDateTime = types.StringNull()
					}
					if v.GetSkuId() != nil {
						tfStateLicenseAssignmentStates.SkuId = types.StringValue(v.GetSkuId().String())
					} else {
						tfStateLicenseAssignmentStates.SkuId = types.StringNull()
					}
					if v.GetState() != nil {
						tfStateLicenseAssignmentStates.State = types.StringValue(*v.GetState())
					} else {
						tfStateLicenseAssignmentStates.State = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateLicenseAssignmentStates.AttributeTypes(), tfStateLicenseAssignmentStates)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.LicenseAssignmentStates, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetMail() != nil {
				tfStateValue.Mail = types.StringValue(*v.GetMail())
			} else {
				tfStateValue.Mail = types.StringNull()
			}
			if v.GetMailNickname() != nil {
				tfStateValue.MailNickname = types.StringValue(*v.GetMailNickname())
			} else {
				tfStateValue.MailNickname = types.StringNull()
			}
			if v.GetMobilePhone() != nil {
				tfStateValue.MobilePhone = types.StringValue(*v.GetMobilePhone())
			} else {
				tfStateValue.MobilePhone = types.StringNull()
			}
			if v.GetOfficeLocation() != nil {
				tfStateValue.OfficeLocation = types.StringValue(*v.GetOfficeLocation())
			} else {
				tfStateValue.OfficeLocation = types.StringNull()
			}
			if v.GetOnPremisesDistinguishedName() != nil {
				tfStateValue.OnPremisesDistinguishedName = types.StringValue(*v.GetOnPremisesDistinguishedName())
			} else {
				tfStateValue.OnPremisesDistinguishedName = types.StringNull()
			}
			if v.GetOnPremisesDomainName() != nil {
				tfStateValue.OnPremisesDomainName = types.StringValue(*v.GetOnPremisesDomainName())
			} else {
				tfStateValue.OnPremisesDomainName = types.StringNull()
			}
			if v.GetOnPremisesExtensionAttributes() != nil {
				tfStateOnPremisesExtensionAttributes := usersOnPremisesExtensionAttributesModel{}

				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute1() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute1())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute10() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute10())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute11() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute11())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute12() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute12())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute13() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute13())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute14() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute14())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute15() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute15())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute2() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute2())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute3() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute3())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute4() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute4())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute5() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute5())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute6() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute6())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute7() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute7())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute8() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute8())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringNull()
				}
				if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute9() != nil {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute9())
				} else {
					tfStateOnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringNull()
				}

				tfStateValue.OnPremisesExtensionAttributes, _ = types.ObjectValueFrom(ctx, tfStateOnPremisesExtensionAttributes.AttributeTypes(), tfStateOnPremisesExtensionAttributes)
			}
			if v.GetOnPremisesImmutableId() != nil {
				tfStateValue.OnPremisesImmutableId = types.StringValue(*v.GetOnPremisesImmutableId())
			} else {
				tfStateValue.OnPremisesImmutableId = types.StringNull()
			}
			if v.GetOnPremisesLastSyncDateTime() != nil {
				tfStateValue.OnPremisesLastSyncDateTime = types.StringValue(v.GetOnPremisesLastSyncDateTime().String())
			} else {
				tfStateValue.OnPremisesLastSyncDateTime = types.StringNull()
			}
			if len(v.GetOnPremisesProvisioningErrors()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetOnPremisesProvisioningErrors() {
					tfStateOnPremisesProvisioningErrors := usersOnPremisesProvisioningErrorModel{}

					if v.GetCategory() != nil {
						tfStateOnPremisesProvisioningErrors.Category = types.StringValue(*v.GetCategory())
					} else {
						tfStateOnPremisesProvisioningErrors.Category = types.StringNull()
					}
					if v.GetOccurredDateTime() != nil {
						tfStateOnPremisesProvisioningErrors.OccurredDateTime = types.StringValue(v.GetOccurredDateTime().String())
					} else {
						tfStateOnPremisesProvisioningErrors.OccurredDateTime = types.StringNull()
					}
					if v.GetPropertyCausingError() != nil {
						tfStateOnPremisesProvisioningErrors.PropertyCausingError = types.StringValue(*v.GetPropertyCausingError())
					} else {
						tfStateOnPremisesProvisioningErrors.PropertyCausingError = types.StringNull()
					}
					if v.GetValue() != nil {
						tfStateOnPremisesProvisioningErrors.Value = types.StringValue(*v.GetValue())
					} else {
						tfStateOnPremisesProvisioningErrors.Value = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateOnPremisesProvisioningErrors.AttributeTypes(), tfStateOnPremisesProvisioningErrors)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.OnPremisesProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetOnPremisesSamAccountName() != nil {
				tfStateValue.OnPremisesSamAccountName = types.StringValue(*v.GetOnPremisesSamAccountName())
			} else {
				tfStateValue.OnPremisesSamAccountName = types.StringNull()
			}
			if v.GetOnPremisesSecurityIdentifier() != nil {
				tfStateValue.OnPremisesSecurityIdentifier = types.StringValue(*v.GetOnPremisesSecurityIdentifier())
			} else {
				tfStateValue.OnPremisesSecurityIdentifier = types.StringNull()
			}
			if v.GetOnPremisesSyncEnabled() != nil {
				tfStateValue.OnPremisesSyncEnabled = types.BoolValue(*v.GetOnPremisesSyncEnabled())
			} else {
				tfStateValue.OnPremisesSyncEnabled = types.BoolNull()
			}
			if v.GetOnPremisesUserPrincipalName() != nil {
				tfStateValue.OnPremisesUserPrincipalName = types.StringValue(*v.GetOnPremisesUserPrincipalName())
			} else {
				tfStateValue.OnPremisesUserPrincipalName = types.StringNull()
			}
			if len(v.GetOtherMails()) > 0 {
				var valueArrayOtherMails []attr.Value
				for _, v := range v.GetOtherMails() {
					valueArrayOtherMails = append(valueArrayOtherMails, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayOtherMails)
				tfStateValue.OtherMails = listValue
			} else {
				tfStateValue.OtherMails = types.ListNull(types.StringType)
			}
			if v.GetPasswordPolicies() != nil {
				tfStateValue.PasswordPolicies = types.StringValue(*v.GetPasswordPolicies())
			} else {
				tfStateValue.PasswordPolicies = types.StringNull()
			}
			if v.GetPasswordProfile() != nil {
				tfStatePasswordProfile := usersPasswordProfileModel{}

				if v.GetPasswordProfile().GetForceChangePasswordNextSignIn() != nil {
					tfStatePasswordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*v.GetPasswordProfile().GetForceChangePasswordNextSignIn())
				} else {
					tfStatePasswordProfile.ForceChangePasswordNextSignIn = types.BoolNull()
				}
				if v.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa() != nil {
					tfStatePasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*v.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
				} else {
					tfStatePasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolNull()
				}
				if v.GetPasswordProfile().GetPassword() != nil {
					tfStatePasswordProfile.Password = types.StringValue(*v.GetPasswordProfile().GetPassword())
				} else {
					tfStatePasswordProfile.Password = types.StringNull()
				}

				tfStateValue.PasswordProfile, _ = types.ObjectValueFrom(ctx, tfStatePasswordProfile.AttributeTypes(), tfStatePasswordProfile)
			}
			if v.GetPostalCode() != nil {
				tfStateValue.PostalCode = types.StringValue(*v.GetPostalCode())
			} else {
				tfStateValue.PostalCode = types.StringNull()
			}
			if v.GetPreferredDataLocation() != nil {
				tfStateValue.PreferredDataLocation = types.StringValue(*v.GetPreferredDataLocation())
			} else {
				tfStateValue.PreferredDataLocation = types.StringNull()
			}
			if v.GetPreferredLanguage() != nil {
				tfStateValue.PreferredLanguage = types.StringValue(*v.GetPreferredLanguage())
			} else {
				tfStateValue.PreferredLanguage = types.StringNull()
			}
			if len(v.GetProvisionedPlans()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetProvisionedPlans() {
					tfStateProvisionedPlans := usersProvisionedPlanModel{}

					if v.GetCapabilityStatus() != nil {
						tfStateProvisionedPlans.CapabilityStatus = types.StringValue(*v.GetCapabilityStatus())
					} else {
						tfStateProvisionedPlans.CapabilityStatus = types.StringNull()
					}
					if v.GetProvisioningStatus() != nil {
						tfStateProvisionedPlans.ProvisioningStatus = types.StringValue(*v.GetProvisioningStatus())
					} else {
						tfStateProvisionedPlans.ProvisioningStatus = types.StringNull()
					}
					if v.GetService() != nil {
						tfStateProvisionedPlans.Service = types.StringValue(*v.GetService())
					} else {
						tfStateProvisionedPlans.Service = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateProvisionedPlans.AttributeTypes(), tfStateProvisionedPlans)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.ProvisionedPlans, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if len(v.GetProxyAddresses()) > 0 {
				var valueArrayProxyAddresses []attr.Value
				for _, v := range v.GetProxyAddresses() {
					valueArrayProxyAddresses = append(valueArrayProxyAddresses, types.StringValue(v))
				}
				listValue, _ := types.ListValue(types.StringType, valueArrayProxyAddresses)
				tfStateValue.ProxyAddresses = listValue
			} else {
				tfStateValue.ProxyAddresses = types.ListNull(types.StringType)
			}
			if v.GetSecurityIdentifier() != nil {
				tfStateValue.SecurityIdentifier = types.StringValue(*v.GetSecurityIdentifier())
			} else {
				tfStateValue.SecurityIdentifier = types.StringNull()
			}
			if len(v.GetServiceProvisioningErrors()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range v.GetServiceProvisioningErrors() {
					tfStateServiceProvisioningErrors := usersServiceProvisioningErrorModel{}

					if v.GetCreatedDateTime() != nil {
						tfStateServiceProvisioningErrors.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
					} else {
						tfStateServiceProvisioningErrors.CreatedDateTime = types.StringNull()
					}
					if v.GetIsResolved() != nil {
						tfStateServiceProvisioningErrors.IsResolved = types.BoolValue(*v.GetIsResolved())
					} else {
						tfStateServiceProvisioningErrors.IsResolved = types.BoolNull()
					}
					if v.GetServiceInstance() != nil {
						tfStateServiceProvisioningErrors.ServiceInstance = types.StringValue(*v.GetServiceInstance())
					} else {
						tfStateServiceProvisioningErrors.ServiceInstance = types.StringNull()
					}
					objectValue, _ := types.ObjectValueFrom(ctx, tfStateServiceProvisioningErrors.AttributeTypes(), tfStateServiceProvisioningErrors)
					objectValues = append(objectValues, objectValue)
				}
				tfStateValue.ServiceProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
			}
			if v.GetShowInAddressList() != nil {
				tfStateValue.ShowInAddressList = types.BoolValue(*v.GetShowInAddressList())
			} else {
				tfStateValue.ShowInAddressList = types.BoolNull()
			}
			if v.GetSignInActivity() != nil {
				tfStateSignInActivity := usersSignInActivityModel{}

				if v.GetSignInActivity().GetLastNonInteractiveSignInDateTime() != nil {
					tfStateSignInActivity.LastNonInteractiveSignInDateTime = types.StringValue(v.GetSignInActivity().GetLastNonInteractiveSignInDateTime().String())
				} else {
					tfStateSignInActivity.LastNonInteractiveSignInDateTime = types.StringNull()
				}
				if v.GetSignInActivity().GetLastNonInteractiveSignInRequestId() != nil {
					tfStateSignInActivity.LastNonInteractiveSignInRequestId = types.StringValue(*v.GetSignInActivity().GetLastNonInteractiveSignInRequestId())
				} else {
					tfStateSignInActivity.LastNonInteractiveSignInRequestId = types.StringNull()
				}
				if v.GetSignInActivity().GetLastSignInDateTime() != nil {
					tfStateSignInActivity.LastSignInDateTime = types.StringValue(v.GetSignInActivity().GetLastSignInDateTime().String())
				} else {
					tfStateSignInActivity.LastSignInDateTime = types.StringNull()
				}
				if v.GetSignInActivity().GetLastSignInRequestId() != nil {
					tfStateSignInActivity.LastSignInRequestId = types.StringValue(*v.GetSignInActivity().GetLastSignInRequestId())
				} else {
					tfStateSignInActivity.LastSignInRequestId = types.StringNull()
				}
				if v.GetSignInActivity().GetLastSuccessfulSignInDateTime() != nil {
					tfStateSignInActivity.LastSuccessfulSignInDateTime = types.StringValue(v.GetSignInActivity().GetLastSuccessfulSignInDateTime().String())
				} else {
					tfStateSignInActivity.LastSuccessfulSignInDateTime = types.StringNull()
				}
				if v.GetSignInActivity().GetLastSuccessfulSignInRequestId() != nil {
					tfStateSignInActivity.LastSuccessfulSignInRequestId = types.StringValue(*v.GetSignInActivity().GetLastSuccessfulSignInRequestId())
				} else {
					tfStateSignInActivity.LastSuccessfulSignInRequestId = types.StringNull()
				}

				tfStateValue.SignInActivity, _ = types.ObjectValueFrom(ctx, tfStateSignInActivity.AttributeTypes(), tfStateSignInActivity)
			}
			if v.GetSignInSessionsValidFromDateTime() != nil {
				tfStateValue.SignInSessionsValidFromDateTime = types.StringValue(v.GetSignInSessionsValidFromDateTime().String())
			} else {
				tfStateValue.SignInSessionsValidFromDateTime = types.StringNull()
			}
			if v.GetState() != nil {
				tfStateValue.State = types.StringValue(*v.GetState())
			} else {
				tfStateValue.State = types.StringNull()
			}
			if v.GetStreetAddress() != nil {
				tfStateValue.StreetAddress = types.StringValue(*v.GetStreetAddress())
			} else {
				tfStateValue.StreetAddress = types.StringNull()
			}
			if v.GetSurname() != nil {
				tfStateValue.Surname = types.StringValue(*v.GetSurname())
			} else {
				tfStateValue.Surname = types.StringNull()
			}
			if v.GetUsageLocation() != nil {
				tfStateValue.UsageLocation = types.StringValue(*v.GetUsageLocation())
			} else {
				tfStateValue.UsageLocation = types.StringNull()
			}
			if v.GetUserPrincipalName() != nil {
				tfStateValue.UserPrincipalName = types.StringValue(*v.GetUserPrincipalName())
			} else {
				tfStateValue.UserPrincipalName = types.StringNull()
			}
			if v.GetUserType() != nil {
				tfStateValue.UserType = types.StringValue(*v.GetUserType())
			} else {
				tfStateValue.UserType = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateValue.AttributeTypes(), tfStateValue)
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
