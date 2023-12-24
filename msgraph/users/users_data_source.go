package users

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
							Description: "Sets the age group of the user. Allowed values: null, Minor, NotAdult and Adult. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).",
							Computed:    true,
						},
						"assigned_licenses": schema.ListNestedAttribute{
							Description: "The licenses that are assigned to the user, including inherited (group-based) licenses. This property doesn't differentiate directly assigned and inherited licenses. Use the licenseAssignmentStates property to identify the directly assigned and inherited licenses.  Not nullable. Returned only on $select. Supports $filter (eq, not, /$count eq 0, /$count ne 0).",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"disabled_plans": schema.ListAttribute{
										Description: "A collection of the unique identifiers for plans that have been disabled.",
										Computed:    true,
										ElementType: types.StringType,
									},
									"sku_id": schema.StringAttribute{
										Description: "The unique identifier for the SKU.",
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
							Description: "The telephone numbers for the user. NOTE: Although this is a string collection, only one number can be set for this property. Read-only for users synced from on-premises directory. Returned by default. Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"city": schema.StringAttribute{
							Description: "The city where the user is located. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"company_name": schema.StringAttribute{
							Description: "The name of the company that the user is associated with. This property can be useful for describing the company that an external user comes from. The maximum length is 64 characters.Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"consent_provided_for_minor": schema.StringAttribute{
							Description: "Sets whether consent has been obtained for minors. Allowed values: null, Granted, Denied and NotRequired. Refer to the legal age group property definitions for further information. Returned only on $select. Supports $filter (eq, ne, not, and in).",
							Computed:    true,
						},
						"country": schema.StringAttribute{
							Description: "The country or region where the user is located; for example, US or UK. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"created_date_time": schema.StringAttribute{
							Description: "The date and time the user was created, in ISO 8601 format and in UTC time. The value cannot be modified and is automatically populated when the entity is created. Nullable. For on-premises users, the value represents when they were first created in Microsoft Entra ID. Property is null for some users created before June 2018 and on-premises users that were synced to Microsoft Entra ID before June 2018. Read-only. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
							Computed:    true,
						},
						"creation_type": schema.StringAttribute{
							Description: "Indicates whether the user account was created through one of the following methods:  As a regular school or work account (null). As an external account (Invitation). As a local account for an Azure Active Directory B2C tenant (LocalAccount). Through self-service sign-up by an internal user using email verification (EmailVerified). Through self-service sign-up by an external user signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.Returned only on $select. Supports $filter (eq, ne, not, in).",
							Computed:    true,
						},
						"custom_security_attributes": schema.SingleNestedAttribute{
							Description: "An open complex type that holds the value of a custom security attribute that is assigned to a directory object. Nullable. Returned only on $select. Supports $filter (eq, ne, not, startsWith). Filter value is case sensitive.",
							Computed:    true,
							Attributes:  map[string]schema.Attribute{},
						},
						"department": schema.StringAttribute{
							Description: "The name for the department in which the user works. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, and eq on null values).",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name. This property is required when a user is created and it cannot be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values), $orderby, and $search.",
							Computed:    true,
						},
						"employee_hire_date": schema.StringAttribute{
							Description: "The date and time when the user was hired or will start work in case of a future hire. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
							Computed:    true,
						},
						"employee_id": schema.StringAttribute{
							Description: "The employee identifier assigned to the user by the organization. The maximum length is 16 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"employee_leave_date_time": schema.StringAttribute{
							Description: "The date and time when the user left or will leave the organization. To read this property, the calling app must be assigned the User-LifeCycleInfo.Read.All permission. To write this property, the calling app must be assigned the User.Read.All and User-LifeCycleInfo.ReadWrite.All permissions. To read this property in delegated scenarios, the admin needs one of the following Microsoft Entra roles: Lifecycle Workflows Administrator, Global Reader, or Global Administrator. To write this property in delegated scenarios, the admin needs the Global Administrator role. Supports $filter (eq, ne, not , ge, le, in). For more information, see Configure the employeeLeaveDateTime property for a user.",
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
							Description: "For an external user invited to the tenant using the invitation API, this property represents the invited user's invitation status. For invited users, the state can be PendingAcceptance or Accepted, or null for all other users. Returned only on $select. Supports $filter (eq, ne, not , in).",
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
							Description: "Represents the identities that can be used to sign in to this user account. An identity can be provided by Microsoft (also known as a local account), by organizations, or by social identity providers such as Facebook, Google, and Microsoft, and tied to a user account. May contain multiple items with the same signInType value. Returned only on $select. Supports $filter (eq) including on null values, only where the signInType is not userPrincipalName.",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"issuer": schema.StringAttribute{
										Description: "Specifies the issuer of the identity, for example facebook.com.For local accounts (where signInType isn't federated), this property is the local B2C tenant default domain name, for example contoso.onmicrosoft.com.For guests from other Microsoft Entra organization, this is the domain of the federated organization, for example contoso.com.Supports $filter. 512 character limit.",
										Computed:    true,
									},
									"issuer_assigned_id": schema.StringAttribute{
										Description: "Specifies the unique identifier assigned to the user by the issuer. The combination of issuer and issuerAssignedId must be unique within the organization. Represents the sign-in name for the user, when signInType is set to emailAddress or userName (also known as local accounts).When signInType is set to: emailAddress, (or a custom string that starts with emailAddress like emailAddress1) issuerAssignedId must be a valid email addressuserName, issuerAssignedId must begin with alphabetical character or number, and can only contain alphanumeric characters and the following symbols: - or Supports $filter. 64 character limit.",
										Computed:    true,
									},
									"sign_in_type": schema.StringAttribute{
										Description: "Specifies the user sign-in types in your directory, such as emailAddress, userName, federated, or userPrincipalName. federated represents a unique identifier for a user from an issuer, that can be in any format chosen by the issuer. Setting or updating a userPrincipalName identity will update the value of the userPrincipalName property on the user object. The validations performed on the userPrincipalName property on the user object, for example, verified domains and acceptable characters, will be performed when setting or updating a userPrincipalName identity. Other validation is enforced on issuerAssignedId when the sign-in type is set to emailAddress or userName. This property can also be set to any custom string.",
										Computed:    true,
									},
								},
							},
						},
						"im_addresses": schema.ListAttribute{
							Description: "The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user. Read-only. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"is_resource_account": schema.BoolAttribute{
							Description: "Do not use – reserved for future use.",
							Computed:    true,
						},
						"job_title": schema.StringAttribute{
							Description: "The user's job title. Maximum length is 128 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"last_password_change_date_time": schema.StringAttribute{
							Description: "The time when this Microsoft Entra user last changed their password or when their password was created, whichever date the latest action was performed. The date and time information uses ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.",
							Computed:    true,
						},
						"legal_age_group_classification": schema.StringAttribute{
							Description: "Used by enterprise applications to determine the legal age group of the user. This property is read-only and calculated based on ageGroup and consentProvidedForMinor properties. Allowed values: null, MinorWithOutParentalConsent, MinorWithParentalConsent, MinorNoParentalConsentRequired, NotAdult and Adult. Refer to the legal age group property definitions for further information. Returned only on $select.",
							Computed:    true,
						},
						"license_assignment_states": schema.ListNestedAttribute{
							Description: "State of license assignments for this user. Also indicates licenses that are directly assigned or the user has inherited through group memberships. Read-only. Returned only on $select.",
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
							Description: "The SMTP address for the user, for example, jeff@contoso.onmicrosoft.com. Changes to this property will also update the user's proxyAddresses collection to include the value as an SMTP address. This property can't contain accent characters.  NOTE: We don't recommend updating this property for Azure AD B2C user profiles. Use the otherMails property instead. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith, and eq on null values).",
							Computed:    true,
						},
						"mail_nickname": schema.StringAttribute{
							Description: "The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"mobile_phone": schema.StringAttribute{
							Description: "The primary cellular telephone number for the user. Read-only for users synced from on-premises directory. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values) and $search.",
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
							Description: "Contains extensionAttributes1-15 for the user. These extension attributes are also known as Exchange custom attributes 1-15. For an onPremisesSyncEnabled user, the source of authority for this set of properties is the on-premises and is read-only. For a cloud-only user (where onPremisesSyncEnabled is false), these properties can be set during creation or update of a user object.  For a cloud-only user previously synced from on-premises Active Directory, these properties are read-only in Microsoft Graph but can be fully managed through the Exchange Admin Center or the Exchange Online V2 module in PowerShell. Returned only on $select. Supports $filter (eq, ne, not, in).",
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
							Description: "This property is used to associate an on-premises Active Directory user account to their Microsoft Entra user object. This property must be specified when creating a new user account in the Graph if you're using a federated domain for the user's userPrincipalName (UPN) property. NOTE: The $ and _ characters can't be used when specifying this property. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in)..",
							Computed:    true,
						},
						"on_premises_last_sync_date_time": schema.StringAttribute{
							Description: "Indicates the last time at which the object was synced with the on-premises directory; for example: 2013-02-16T03:04:54Z. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).",
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
							Description: "Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. Returned only on $select.  Supports $filter (eq including on null values).",
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
							Description: "A list of additional email addresses for the user; for example: ['bob@contoso.com', 'Robert@fabrikam.com']. NOTE: This property can't contain accent characters. Returned only on $select. Supports $filter (eq, not, ge, le, in, startsWith, endsWith, /$count eq 0, /$count ne 0).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"password_policies": schema.StringAttribute{
							Description: "Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two may be specified together; for example: DisablePasswordExpiration, DisableStrongPassword. Returned only on $select. For more information on the default password policies, see Microsoft Entra password policies. Supports $filter (ne, not, and eq on null values).",
							Computed:    true,
						},
						"password_profile": schema.SingleNestedAttribute{
							Description: "Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values).",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"force_change_password_next_sign_in": schema.BoolAttribute{
									Description: "true if the user must change her password on the next login; otherwise false.",
									Computed:    true,
								},
								"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
									Description: "If true, at next sign-in, the user must perform a multi-factor authentication (MFA) before being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that the user is required to first perform a multi-factor authentication before password change. After a password change, this property will be automatically reset to false. If not set, default is false.",
									Computed:    true,
								},
								"password": schema.StringAttribute{
									Description: "The password for the user. This property is required when a user is created. It can be updated, but the user will be required to change the password on the next login. The password must satisfy minimum requirements as specified by the user's passwordPolicies property. By default, a strong password is required.",
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
							Description: "For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. Changes to the mail property will also update this collection to include the value as an SMTP address. For more information, see mail and proxyAddresses properties. The proxy address prefixed with SMTP (capitalized) is the primary proxy address while those prefixed with smtp are the secondary proxy addresses. For Azure AD B2C accounts, this property has a limit of 10 unique addresses. Read-only in Microsoft Graph; you can update this property only through the Microsoft 365 admin center. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith, endsWith, /$count eq 0, /$count ne 0).",
							Computed:    true,
							ElementType: types.StringType,
						},
						"security_identifier": schema.StringAttribute{
							Description: "Security identifier (SID) of the user, used in Windows scenarios. Read-only. Returned by default. Supports $select and $filter (eq, not, ge, le, startsWith).",
							Computed:    true,
						},
						"service_provisioning_errors": schema.ListNestedAttribute{
							Description: "",
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"created_date_time": schema.StringAttribute{
										Description: "",
										Computed:    true,
									},
									"is_resolved": schema.BoolAttribute{
										Description: "",
										Computed:    true,
									},
									"service_instance": schema.StringAttribute{
										Description: "",
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
							Description: "Get the last signed-in date and request ID of the sign-in for a given user. Read-only.Returned only on $select. Supports $filter (eq, ne, not, ge, le) but not with any other filterable properties. Note: Details for this property require a Microsoft Entra ID P1 or P2 license and the AuditLog.Read.All permission.This property is not returned for a user who has never signed in or last signed in before April 2020.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"last_non_interactive_sign_in_date_time": schema.StringAttribute{
									Description: "The last non-interactive sign-in date for a specific user. You can use this field to calculate the last time a client attempted to sign into the directory on behalf of a user. Because some users may use clients to access tenant resources rather than signing into your tenant directly, you can use the non-interactive sign-in date to along with lastSignInDateTime to identify inactive users. The timestamp represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is: '2014-01-01T00:00:00Z'. Microsoft Entra ID maintains non-interactive sign-ins going back to May 2020. For more information about using the value of this property, see Manage inactive user accounts in Microsoft Entra ID.",
									Computed:    true,
								},
								"last_non_interactive_sign_in_request_id": schema.StringAttribute{
									Description: "Request identifier of the last non-interactive sign-in performed by this user.",
									Computed:    true,
								},
								"last_sign_in_date_time": schema.StringAttribute{
									Description: "The last interactive sign-in date and time for a specific user. You can use this field to calculate the last time a user attempted to sign into the directory with an interactive authentication method. This field can be used to build reports, such as inactive users. The timestamp represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is: '2014-01-01T00:00:00Z'. Microsoft Entra ID maintains interactive sign-ins going back to April 2020. For more information about using the value of this property, see Manage inactive user accounts in Microsoft Entra ID.",
									Computed:    true,
								},
								"last_sign_in_request_id": schema.StringAttribute{
									Description: "Request identifier of the last interactive sign-in performed by this user.",
									Computed:    true,
								},
							},
						},
						"sign_in_sessions_valid_from_date_time": schema.StringAttribute{
							Description: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).  If this happens, the application needs to acquire a new refresh token by making a request to the authorize endpoint. Read-only. Use revokeSignInSessions to reset. Returned only on $select.",
							Computed:    true,
						},
						"state": schema.StringAttribute{
							Description: "The state or province in the user's address. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"street_address": schema.StringAttribute{
							Description: "The street address of the user's place of business. Maximum length is 1024 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"surname": schema.StringAttribute{
							Description: "The user's surname (family name or last name). Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"usage_location": schema.StringAttribute{
							Description: "A two letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirement to check for availability of services in countries.  Examples include: US, JP, and GB. Not nullable. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
							Computed:    true,
						},
						"user_principal_name": schema.StringAttribute{
							Description: "The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this should map to the user's email name. The general format is alias@domain, where domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization.NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith) and $orderby.",
							Computed:    true,
						},
						"user_type": schema.StringAttribute{
							Description: "A string value that can be used to classify user types in your directory, such as Member and Guest. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). NOTE: For more information about the permissions for member and guest users, see What are the default user permissions in Microsoft Entra ID?",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

type usersDataSourceModel struct {
	Value []usersValueDataSourceModel `tfsdk:"value"`
}

type usersValueDataSourceModel struct {
	Id                              types.String                                       `tfsdk:"id"`
	DeletedDateTime                 types.String                                       `tfsdk:"deleted_date_time"`
	AccountEnabled                  types.Bool                                         `tfsdk:"account_enabled"`
	AgeGroup                        types.String                                       `tfsdk:"age_group"`
	AssignedLicenses                []usersAssignedLicensesDataSourceModel             `tfsdk:"assigned_licenses"`
	AssignedPlans                   []usersAssignedPlansDataSourceModel                `tfsdk:"assigned_plans"`
	AuthorizationInfo               *usersAuthorizationInfoDataSourceModel             `tfsdk:"authorization_info"`
	BusinessPhones                  []types.String                                     `tfsdk:"business_phones"`
	City                            types.String                                       `tfsdk:"city"`
	CompanyName                     types.String                                       `tfsdk:"company_name"`
	ConsentProvidedForMinor         types.String                                       `tfsdk:"consent_provided_for_minor"`
	Country                         types.String                                       `tfsdk:"country"`
	CreatedDateTime                 types.String                                       `tfsdk:"created_date_time"`
	CreationType                    types.String                                       `tfsdk:"creation_type"`
	CustomSecurityAttributes        *usersCustomSecurityAttributesDataSourceModel      `tfsdk:"custom_security_attributes"`
	Department                      types.String                                       `tfsdk:"department"`
	DisplayName                     types.String                                       `tfsdk:"display_name"`
	EmployeeHireDate                types.String                                       `tfsdk:"employee_hire_date"`
	EmployeeId                      types.String                                       `tfsdk:"employee_id"`
	EmployeeLeaveDateTime           types.String                                       `tfsdk:"employee_leave_date_time"`
	EmployeeOrgData                 *usersEmployeeOrgDataDataSourceModel               `tfsdk:"employee_org_data"`
	EmployeeType                    types.String                                       `tfsdk:"employee_type"`
	ExternalUserState               types.String                                       `tfsdk:"external_user_state"`
	ExternalUserStateChangeDateTime types.String                                       `tfsdk:"external_user_state_change_date_time"`
	FaxNumber                       types.String                                       `tfsdk:"fax_number"`
	GivenName                       types.String                                       `tfsdk:"given_name"`
	Identities                      []usersIdentitiesDataSourceModel                   `tfsdk:"identities"`
	ImAddresses                     []types.String                                     `tfsdk:"im_addresses"`
	IsResourceAccount               types.Bool                                         `tfsdk:"is_resource_account"`
	JobTitle                        types.String                                       `tfsdk:"job_title"`
	LastPasswordChangeDateTime      types.String                                       `tfsdk:"last_password_change_date_time"`
	LegalAgeGroupClassification     types.String                                       `tfsdk:"legal_age_group_classification"`
	LicenseAssignmentStates         []usersLicenseAssignmentStatesDataSourceModel      `tfsdk:"license_assignment_states"`
	Mail                            types.String                                       `tfsdk:"mail"`
	MailNickname                    types.String                                       `tfsdk:"mail_nickname"`
	MobilePhone                     types.String                                       `tfsdk:"mobile_phone"`
	OfficeLocation                  types.String                                       `tfsdk:"office_location"`
	OnPremisesDistinguishedName     types.String                                       `tfsdk:"on_premises_distinguished_name"`
	OnPremisesDomainName            types.String                                       `tfsdk:"on_premises_domain_name"`
	OnPremisesExtensionAttributes   *usersOnPremisesExtensionAttributesDataSourceModel `tfsdk:"on_premises_extension_attributes"`
	OnPremisesImmutableId           types.String                                       `tfsdk:"on_premises_immutable_id"`
	OnPremisesLastSyncDateTime      types.String                                       `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesProvisioningErrors    []usersOnPremisesProvisioningErrorsDataSourceModel `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName        types.String                                       `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier    types.String                                       `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled           types.Bool                                         `tfsdk:"on_premises_sync_enabled"`
	OnPremisesUserPrincipalName     types.String                                       `tfsdk:"on_premises_user_principal_name"`
	OtherMails                      []types.String                                     `tfsdk:"other_mails"`
	PasswordPolicies                types.String                                       `tfsdk:"password_policies"`
	PasswordProfile                 *usersPasswordProfileDataSourceModel               `tfsdk:"password_profile"`
	PostalCode                      types.String                                       `tfsdk:"postal_code"`
	PreferredDataLocation           types.String                                       `tfsdk:"preferred_data_location"`
	PreferredLanguage               types.String                                       `tfsdk:"preferred_language"`
	ProvisionedPlans                []usersProvisionedPlansDataSourceModel             `tfsdk:"provisioned_plans"`
	ProxyAddresses                  []types.String                                     `tfsdk:"proxy_addresses"`
	SecurityIdentifier              types.String                                       `tfsdk:"security_identifier"`
	ServiceProvisioningErrors       []usersServiceProvisioningErrorsDataSourceModel    `tfsdk:"service_provisioning_errors"`
	ShowInAddressList               types.Bool                                         `tfsdk:"show_in_address_list"`
	SignInActivity                  *usersSignInActivityDataSourceModel                `tfsdk:"sign_in_activity"`
	SignInSessionsValidFromDateTime types.String                                       `tfsdk:"sign_in_sessions_valid_from_date_time"`
	State                           types.String                                       `tfsdk:"state"`
	StreetAddress                   types.String                                       `tfsdk:"street_address"`
	Surname                         types.String                                       `tfsdk:"surname"`
	UsageLocation                   types.String                                       `tfsdk:"usage_location"`
	UserPrincipalName               types.String                                       `tfsdk:"user_principal_name"`
	UserType                        types.String                                       `tfsdk:"user_type"`
}

type usersAssignedLicensesDataSourceModel struct {
	DisabledPlans []types.String `tfsdk:"disabled_plans"`
	SkuId         types.String   `tfsdk:"sku_id"`
}

type usersAssignedPlansDataSourceModel struct {
	AssignedDateTime types.String `tfsdk:"assigned_date_time"`
	CapabilityStatus types.String `tfsdk:"capability_status"`
	Service          types.String `tfsdk:"service"`
	ServicePlanId    types.String `tfsdk:"service_plan_id"`
}

type usersAuthorizationInfoDataSourceModel struct {
	CertificateUserIds []types.String `tfsdk:"certificate_user_ids"`
}

type usersCustomSecurityAttributesDataSourceModel struct {
}

type usersEmployeeOrgDataDataSourceModel struct {
	CostCenter types.String `tfsdk:"cost_center"`
	Division   types.String `tfsdk:"division"`
}

type usersIdentitiesDataSourceModel struct {
	Issuer           types.String `tfsdk:"issuer"`
	IssuerAssignedId types.String `tfsdk:"issuer_assigned_id"`
	SignInType       types.String `tfsdk:"sign_in_type"`
}

type usersLicenseAssignmentStatesDataSourceModel struct {
	AssignedByGroup     types.String   `tfsdk:"assigned_by_group"`
	DisabledPlans       []types.String `tfsdk:"disabled_plans"`
	Error               types.String   `tfsdk:"error"`
	LastUpdatedDateTime types.String   `tfsdk:"last_updated_date_time"`
	SkuId               types.String   `tfsdk:"sku_id"`
	State               types.String   `tfsdk:"state"`
}

type usersOnPremisesExtensionAttributesDataSourceModel struct {
	ExtensionAttribute1  types.String `tfsdk:"extension_attribute_1"`
	ExtensionAttribute10 types.String `tfsdk:"extension_attribute_10"`
	ExtensionAttribute11 types.String `tfsdk:"extension_attribute_11"`
	ExtensionAttribute12 types.String `tfsdk:"extension_attribute_12"`
	ExtensionAttribute13 types.String `tfsdk:"extension_attribute_13"`
	ExtensionAttribute14 types.String `tfsdk:"extension_attribute_14"`
	ExtensionAttribute15 types.String `tfsdk:"extension_attribute_15"`
	ExtensionAttribute2  types.String `tfsdk:"extension_attribute_2"`
	ExtensionAttribute3  types.String `tfsdk:"extension_attribute_3"`
	ExtensionAttribute4  types.String `tfsdk:"extension_attribute_4"`
	ExtensionAttribute5  types.String `tfsdk:"extension_attribute_5"`
	ExtensionAttribute6  types.String `tfsdk:"extension_attribute_6"`
	ExtensionAttribute7  types.String `tfsdk:"extension_attribute_7"`
	ExtensionAttribute8  types.String `tfsdk:"extension_attribute_8"`
	ExtensionAttribute9  types.String `tfsdk:"extension_attribute_9"`
}

type usersOnPremisesProvisioningErrorsDataSourceModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

type usersPasswordProfileDataSourceModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
	Password                             types.String `tfsdk:"password"`
}

type usersProvisionedPlansDataSourceModel struct {
	CapabilityStatus   types.String `tfsdk:"capability_status"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	Service            types.String `tfsdk:"service"`
}

type usersServiceProvisioningErrorsDataSourceModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

type usersSignInActivityDataSourceModel struct {
	LastNonInteractiveSignInDateTime  types.String `tfsdk:"last_non_interactive_sign_in_date_time"`
	LastNonInteractiveSignInRequestId types.String `tfsdk:"last_non_interactive_sign_in_request_id"`
	LastSignInDateTime                types.String `tfsdk:"last_sign_in_date_time"`
	LastSignInRequestId               types.String `tfsdk:"last_sign_in_request_id"`
}

// Read refreshes the Terraform state with the latest data.
func (d *usersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state usersDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"deletedDateTime",
				"signInActivity",
				"accountEnabled",
				"ageGroup",
				"assignedLicenses",
				"assignedPlans",
				"authorizationInfo",
				"businessPhones",
				"city",
				"companyName",
				"consentProvidedForMinor",
				"country",
				"createdDateTime",
				"creationType",
				"customSecurityAttributes",
				"department",
				"displayName",
				"employeeHireDate",
				"employeeId",
				"employeeLeaveDateTime",
				"employeeOrgData",
				"employeeType",
				"externalUserState",
				"externalUserStateChangeDateTime",
				"faxNumber",
				"givenName",
				"identities",
				"imAddresses",
				"isResourceAccount",
				"jobTitle",
				"lastPasswordChangeDateTime",
				"legalAgeGroupClassification",
				"licenseAssignmentStates",
				"mail",
				"mailNickname",
				"mobilePhone",
				"officeLocation",
				"onPremisesDistinguishedName",
				"onPremisesDomainName",
				"onPremisesExtensionAttributes",
				"onPremisesImmutableId",
				"onPremisesLastSyncDateTime",
				"onPremisesProvisioningErrors",
				"onPremisesSamAccountName",
				"onPremisesSecurityIdentifier",
				"onPremisesSyncEnabled",
				"onPremisesUserPrincipalName",
				"otherMails",
				"passwordPolicies",
				"passwordProfile",
				"postalCode",
				"preferredDataLocation",
				"preferredLanguage",
				"provisionedPlans",
				"proxyAddresses",
				"securityIdentifier",
				"serviceProvisioningErrors",
				"showInAddressList",
				"signInSessionsValidFromDateTime",
				"state",
				"streetAddress",
				"surname",
				"usageLocation",
				"userPrincipalName",
				"userType",
				"appRoleAssignments",
				"createdObjects",
				"directReports",
				"licenseDetails",
				"manager",
				"memberOf",
				"oauth2PermissionGrants",
				"ownedDevices",
				"ownedObjects",
				"registeredDevices",
				"scopedRoleMemberOf",
				"transitiveMemberOf",
				"calendar",
				"calendarGroups",
				"calendars",
				"calendarView",
				"contactFolders",
				"contacts",
				"events",
				"inferenceClassification",
				"mailFolders",
				"messages",
				"outlook",
				"people",
				"drive",
				"drives",
				"followedSites",
				"extensions",
				"agreementAcceptances",
				"managedDevices",
				"managedAppRegistrations",
				"deviceManagementTroubleshootingEvents",
				"planner",
				"insights",
				"settings",
				"onenote",
				"photo",
				"photos",
				"activities",
				"onlineMeetings",
				"presence",
				"authentication",
				"chats",
				"joinedTeams",
				"teamwork",
				"todo",
				"employeeExperience",
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

	for _, v := range result.GetValue() {
		value := new(usersValueDataSourceModel)

		if v.GetId() != nil {
			value.Id = types.StringValue(*v.GetId())
		}
		if v.GetDeletedDateTime() != nil {
			value.DeletedDateTime = types.StringValue(v.GetDeletedDateTime().String())
		}
		if v.GetAccountEnabled() != nil {
			value.AccountEnabled = types.BoolValue(*v.GetAccountEnabled())
		}
		if v.GetAgeGroup() != nil {
			value.AgeGroup = types.StringValue(*v.GetAgeGroup())
		}
		for _, v := range v.GetAssignedLicenses() {
			assignedLicenses := new(usersAssignedLicensesDataSourceModel)

			for _, v := range v.GetDisabledPlans() {
				assignedLicenses.DisabledPlans = append(assignedLicenses.DisabledPlans, types.StringValue(v.String()))
			}
			if v.GetSkuId() != nil {
				assignedLicenses.SkuId = types.StringValue(v.GetSkuId().String())
			}
			value.AssignedLicenses = append(value.AssignedLicenses, *assignedLicenses)
		}
		for _, v := range v.GetAssignedPlans() {
			assignedPlans := new(usersAssignedPlansDataSourceModel)

			if v.GetAssignedDateTime() != nil {
				assignedPlans.AssignedDateTime = types.StringValue(v.GetAssignedDateTime().String())
			}
			if v.GetCapabilityStatus() != nil {
				assignedPlans.CapabilityStatus = types.StringValue(*v.GetCapabilityStatus())
			}
			if v.GetService() != nil {
				assignedPlans.Service = types.StringValue(*v.GetService())
			}
			if v.GetServicePlanId() != nil {
				assignedPlans.ServicePlanId = types.StringValue(v.GetServicePlanId().String())
			}
			value.AssignedPlans = append(value.AssignedPlans, *assignedPlans)
		}
		if v.GetAuthorizationInfo() != nil {
			value.AuthorizationInfo = new(usersAuthorizationInfoDataSourceModel)

			for _, v := range v.GetAuthorizationInfo().GetCertificateUserIds() {
				value.AuthorizationInfo.CertificateUserIds = append(value.AuthorizationInfo.CertificateUserIds, types.StringValue(v))
			}
		}
		for _, v := range v.GetBusinessPhones() {
			value.BusinessPhones = append(value.BusinessPhones, types.StringValue(v))
		}
		if v.GetCity() != nil {
			value.City = types.StringValue(*v.GetCity())
		}
		if v.GetCompanyName() != nil {
			value.CompanyName = types.StringValue(*v.GetCompanyName())
		}
		if v.GetConsentProvidedForMinor() != nil {
			value.ConsentProvidedForMinor = types.StringValue(*v.GetConsentProvidedForMinor())
		}
		if v.GetCountry() != nil {
			value.Country = types.StringValue(*v.GetCountry())
		}
		if v.GetCreatedDateTime() != nil {
			value.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
		}
		if v.GetCreationType() != nil {
			value.CreationType = types.StringValue(*v.GetCreationType())
		}
		if v.GetCustomSecurityAttributes() != nil {
			value.CustomSecurityAttributes = new(usersCustomSecurityAttributesDataSourceModel)

		}
		if v.GetDepartment() != nil {
			value.Department = types.StringValue(*v.GetDepartment())
		}
		if v.GetDisplayName() != nil {
			value.DisplayName = types.StringValue(*v.GetDisplayName())
		}
		if v.GetEmployeeHireDate() != nil {
			value.EmployeeHireDate = types.StringValue(v.GetEmployeeHireDate().String())
		}
		if v.GetEmployeeId() != nil {
			value.EmployeeId = types.StringValue(*v.GetEmployeeId())
		}
		if v.GetEmployeeLeaveDateTime() != nil {
			value.EmployeeLeaveDateTime = types.StringValue(v.GetEmployeeLeaveDateTime().String())
		}
		if v.GetEmployeeOrgData() != nil {
			value.EmployeeOrgData = new(usersEmployeeOrgDataDataSourceModel)

			if v.GetEmployeeOrgData().GetCostCenter() != nil {
				value.EmployeeOrgData.CostCenter = types.StringValue(*v.GetEmployeeOrgData().GetCostCenter())
			}
			if v.GetEmployeeOrgData().GetDivision() != nil {
				value.EmployeeOrgData.Division = types.StringValue(*v.GetEmployeeOrgData().GetDivision())
			}
		}
		if v.GetEmployeeType() != nil {
			value.EmployeeType = types.StringValue(*v.GetEmployeeType())
		}
		if v.GetExternalUserState() != nil {
			value.ExternalUserState = types.StringValue(*v.GetExternalUserState())
		}
		if v.GetExternalUserStateChangeDateTime() != nil {
			value.ExternalUserStateChangeDateTime = types.StringValue(v.GetExternalUserStateChangeDateTime().String())
		}
		if v.GetFaxNumber() != nil {
			value.FaxNumber = types.StringValue(*v.GetFaxNumber())
		}
		if v.GetGivenName() != nil {
			value.GivenName = types.StringValue(*v.GetGivenName())
		}
		for _, v := range v.GetIdentities() {
			identities := new(usersIdentitiesDataSourceModel)

			if v.GetIssuer() != nil {
				identities.Issuer = types.StringValue(*v.GetIssuer())
			}
			if v.GetIssuerAssignedId() != nil {
				identities.IssuerAssignedId = types.StringValue(*v.GetIssuerAssignedId())
			}
			if v.GetSignInType() != nil {
				identities.SignInType = types.StringValue(*v.GetSignInType())
			}
			value.Identities = append(value.Identities, *identities)
		}
		for _, v := range v.GetImAddresses() {
			value.ImAddresses = append(value.ImAddresses, types.StringValue(v))
		}
		if v.GetIsResourceAccount() != nil {
			value.IsResourceAccount = types.BoolValue(*v.GetIsResourceAccount())
		}
		if v.GetJobTitle() != nil {
			value.JobTitle = types.StringValue(*v.GetJobTitle())
		}
		if v.GetLastPasswordChangeDateTime() != nil {
			value.LastPasswordChangeDateTime = types.StringValue(v.GetLastPasswordChangeDateTime().String())
		}
		if v.GetLegalAgeGroupClassification() != nil {
			value.LegalAgeGroupClassification = types.StringValue(*v.GetLegalAgeGroupClassification())
		}
		for _, v := range v.GetLicenseAssignmentStates() {
			licenseAssignmentStates := new(usersLicenseAssignmentStatesDataSourceModel)

			if v.GetAssignedByGroup() != nil {
				licenseAssignmentStates.AssignedByGroup = types.StringValue(*v.GetAssignedByGroup())
			}
			for _, v := range v.GetDisabledPlans() {
				licenseAssignmentStates.DisabledPlans = append(licenseAssignmentStates.DisabledPlans, types.StringValue(v.String()))
			}
			if v.GetError() != nil {
				licenseAssignmentStates.Error = types.StringValue(*v.GetError())
			}
			if v.GetLastUpdatedDateTime() != nil {
				licenseAssignmentStates.LastUpdatedDateTime = types.StringValue(v.GetLastUpdatedDateTime().String())
			}
			if v.GetSkuId() != nil {
				licenseAssignmentStates.SkuId = types.StringValue(v.GetSkuId().String())
			}
			if v.GetState() != nil {
				licenseAssignmentStates.State = types.StringValue(*v.GetState())
			}
			value.LicenseAssignmentStates = append(value.LicenseAssignmentStates, *licenseAssignmentStates)
		}
		if v.GetMail() != nil {
			value.Mail = types.StringValue(*v.GetMail())
		}
		if v.GetMailNickname() != nil {
			value.MailNickname = types.StringValue(*v.GetMailNickname())
		}
		if v.GetMobilePhone() != nil {
			value.MobilePhone = types.StringValue(*v.GetMobilePhone())
		}
		if v.GetOfficeLocation() != nil {
			value.OfficeLocation = types.StringValue(*v.GetOfficeLocation())
		}
		if v.GetOnPremisesDistinguishedName() != nil {
			value.OnPremisesDistinguishedName = types.StringValue(*v.GetOnPremisesDistinguishedName())
		}
		if v.GetOnPremisesDomainName() != nil {
			value.OnPremisesDomainName = types.StringValue(*v.GetOnPremisesDomainName())
		}
		if v.GetOnPremisesExtensionAttributes() != nil {
			value.OnPremisesExtensionAttributes = new(usersOnPremisesExtensionAttributesDataSourceModel)

			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute1() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute1())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute10() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute10())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute11() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute11())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute12() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute12())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute13() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute13())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute14() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute14())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute15() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute15())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute2() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute2())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute3() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute3())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute4() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute4())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute5() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute5())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute6() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute6())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute7() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute7())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute8() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute8())
			}
			if v.GetOnPremisesExtensionAttributes().GetExtensionAttribute9() != nil {
				value.OnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringValue(*v.GetOnPremisesExtensionAttributes().GetExtensionAttribute9())
			}
		}
		if v.GetOnPremisesImmutableId() != nil {
			value.OnPremisesImmutableId = types.StringValue(*v.GetOnPremisesImmutableId())
		}
		if v.GetOnPremisesLastSyncDateTime() != nil {
			value.OnPremisesLastSyncDateTime = types.StringValue(v.GetOnPremisesLastSyncDateTime().String())
		}
		for _, v := range v.GetOnPremisesProvisioningErrors() {
			onPremisesProvisioningErrors := new(usersOnPremisesProvisioningErrorsDataSourceModel)

			if v.GetCategory() != nil {
				onPremisesProvisioningErrors.Category = types.StringValue(*v.GetCategory())
			}
			if v.GetOccurredDateTime() != nil {
				onPremisesProvisioningErrors.OccurredDateTime = types.StringValue(v.GetOccurredDateTime().String())
			}
			if v.GetPropertyCausingError() != nil {
				onPremisesProvisioningErrors.PropertyCausingError = types.StringValue(*v.GetPropertyCausingError())
			}
			if v.GetValue() != nil {
				onPremisesProvisioningErrors.Value = types.StringValue(*v.GetValue())
			}
			value.OnPremisesProvisioningErrors = append(value.OnPremisesProvisioningErrors, *onPremisesProvisioningErrors)
		}
		if v.GetOnPremisesSamAccountName() != nil {
			value.OnPremisesSamAccountName = types.StringValue(*v.GetOnPremisesSamAccountName())
		}
		if v.GetOnPremisesSecurityIdentifier() != nil {
			value.OnPremisesSecurityIdentifier = types.StringValue(*v.GetOnPremisesSecurityIdentifier())
		}
		if v.GetOnPremisesSyncEnabled() != nil {
			value.OnPremisesSyncEnabled = types.BoolValue(*v.GetOnPremisesSyncEnabled())
		}
		if v.GetOnPremisesUserPrincipalName() != nil {
			value.OnPremisesUserPrincipalName = types.StringValue(*v.GetOnPremisesUserPrincipalName())
		}
		for _, v := range v.GetOtherMails() {
			value.OtherMails = append(value.OtherMails, types.StringValue(v))
		}
		if v.GetPasswordPolicies() != nil {
			value.PasswordPolicies = types.StringValue(*v.GetPasswordPolicies())
		}
		if v.GetPasswordProfile() != nil {
			value.PasswordProfile = new(usersPasswordProfileDataSourceModel)

			if v.GetPasswordProfile().GetForceChangePasswordNextSignIn() != nil {
				value.PasswordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*v.GetPasswordProfile().GetForceChangePasswordNextSignIn())
			}
			if v.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa() != nil {
				value.PasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*v.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
			}
			if v.GetPasswordProfile().GetPassword() != nil {
				value.PasswordProfile.Password = types.StringValue(*v.GetPasswordProfile().GetPassword())
			}
		}
		if v.GetPostalCode() != nil {
			value.PostalCode = types.StringValue(*v.GetPostalCode())
		}
		if v.GetPreferredDataLocation() != nil {
			value.PreferredDataLocation = types.StringValue(*v.GetPreferredDataLocation())
		}
		if v.GetPreferredLanguage() != nil {
			value.PreferredLanguage = types.StringValue(*v.GetPreferredLanguage())
		}
		for _, v := range v.GetProvisionedPlans() {
			provisionedPlans := new(usersProvisionedPlansDataSourceModel)

			if v.GetCapabilityStatus() != nil {
				provisionedPlans.CapabilityStatus = types.StringValue(*v.GetCapabilityStatus())
			}
			if v.GetProvisioningStatus() != nil {
				provisionedPlans.ProvisioningStatus = types.StringValue(*v.GetProvisioningStatus())
			}
			if v.GetService() != nil {
				provisionedPlans.Service = types.StringValue(*v.GetService())
			}
			value.ProvisionedPlans = append(value.ProvisionedPlans, *provisionedPlans)
		}
		for _, v := range v.GetProxyAddresses() {
			value.ProxyAddresses = append(value.ProxyAddresses, types.StringValue(v))
		}
		if v.GetSecurityIdentifier() != nil {
			value.SecurityIdentifier = types.StringValue(*v.GetSecurityIdentifier())
		}
		for _, v := range v.GetServiceProvisioningErrors() {
			serviceProvisioningErrors := new(usersServiceProvisioningErrorsDataSourceModel)

			if v.GetCreatedDateTime() != nil {
				serviceProvisioningErrors.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			}
			if v.GetIsResolved() != nil {
				serviceProvisioningErrors.IsResolved = types.BoolValue(*v.GetIsResolved())
			}
			if v.GetServiceInstance() != nil {
				serviceProvisioningErrors.ServiceInstance = types.StringValue(*v.GetServiceInstance())
			}
			value.ServiceProvisioningErrors = append(value.ServiceProvisioningErrors, *serviceProvisioningErrors)
		}
		if v.GetShowInAddressList() != nil {
			value.ShowInAddressList = types.BoolValue(*v.GetShowInAddressList())
		}
		if v.GetSignInActivity() != nil {
			value.SignInActivity = new(usersSignInActivityDataSourceModel)

			if v.GetSignInActivity().GetLastNonInteractiveSignInDateTime() != nil {
				value.SignInActivity.LastNonInteractiveSignInDateTime = types.StringValue(v.GetSignInActivity().GetLastNonInteractiveSignInDateTime().String())
			}
			if v.GetSignInActivity().GetLastNonInteractiveSignInRequestId() != nil {
				value.SignInActivity.LastNonInteractiveSignInRequestId = types.StringValue(*v.GetSignInActivity().GetLastNonInteractiveSignInRequestId())
			}
			if v.GetSignInActivity().GetLastSignInDateTime() != nil {
				value.SignInActivity.LastSignInDateTime = types.StringValue(v.GetSignInActivity().GetLastSignInDateTime().String())
			}
			if v.GetSignInActivity().GetLastSignInRequestId() != nil {
				value.SignInActivity.LastSignInRequestId = types.StringValue(*v.GetSignInActivity().GetLastSignInRequestId())
			}
		}
		if v.GetSignInSessionsValidFromDateTime() != nil {
			value.SignInSessionsValidFromDateTime = types.StringValue(v.GetSignInSessionsValidFromDateTime().String())
		}
		if v.GetState() != nil {
			value.State = types.StringValue(*v.GetState())
		}
		if v.GetStreetAddress() != nil {
			value.StreetAddress = types.StringValue(*v.GetStreetAddress())
		}
		if v.GetSurname() != nil {
			value.Surname = types.StringValue(*v.GetSurname())
		}
		if v.GetUsageLocation() != nil {
			value.UsageLocation = types.StringValue(*v.GetUsageLocation())
		}
		if v.GetUserPrincipalName() != nil {
			value.UserPrincipalName = types.StringValue(*v.GetUserPrincipalName())
		}
		if v.GetUserType() != nil {
			value.UserType = types.StringValue(*v.GetUserType())
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
