package users

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &userDataSource{}
	_ datasource.DataSourceWithConfigure = &userDataSource{}
)

// NewUserDataSource is a helper function to simplify the provider implementation.
func NewUserDataSource() datasource.DataSource {
	return &userDataSource{}
}

// UserDataSource is the data source implementation.
type userDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *userDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Configure adds the provider configured client to the data source.
func (d *userDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *userDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"about_me": schema.StringAttribute{
				MarkdownDescription: "A freeform text entry field for the user to describe themselves. Returned only on `$select`.",
				Computed:            true,
			},
			"account_enabled": schema.BoolAttribute{
				MarkdownDescription: "`true` if the account is enabled; otherwise, `false`. This property is required when a user is created. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, and `in`).",
				Computed:            true,
			},
			"age_group": schema.StringAttribute{
				MarkdownDescription: "Sets the age group of the user. Allowed values: `null`, `Minor`, `NotAdult` and `Adult`. Refer to the [legal age group property definitions](#legal-age-group-property-definitions) for further information. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, and `in`).",
				Computed:            true,
			},
			"assigned_licenses": schema.SingleNestedAttribute{
				MarkdownDescription: "The licenses that are assigned to the user, including inherited (group-based) licenses. This property doesn't differentiate directly-assigned and inherited licenses. Use the **licenseAssignmentStates** property to identify the directly-assigned and inherited licenses.  Not nullable. Returned only on `$select`. Supports `$filter` (`eq`, `not`, `/$count eq 0`, `/$count ne 0`).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"disabled_plans": schema.ListAttribute{
						MarkdownDescription: "A collection of the unique identifiers for plans that have been disabled.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"sku_id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier for the SKU.",
						Computed:            true,
					},
				},
			},
			"assigned_plans": schema.SingleNestedAttribute{
				MarkdownDescription: "The plans that are assigned to the user. Read-only. Not nullable. <br><br>Returned only on `$select`. Supports `$filter` (`eq` and `not`).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"assigned_date_time": schema.StringAttribute{
						MarkdownDescription: "The date and time at which the plan was assigned. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`.",
						Computed:            true,
					},
					"capability_status": schema.StringAttribute{
						MarkdownDescription: "Condition of the capability assignment. The possible values are `Enabled`, `Warning`, `Suspended`, `Deleted`, `LockedOut`. See [a detailed description](#capabilitystatus-values) of each value.",
						Computed:            true,
					},
					"service": schema.StringAttribute{
						MarkdownDescription: "The name of the service; for example, `exchange`.",
						Computed:            true,
					},
					"service_plan_id": schema.StringAttribute{
						MarkdownDescription: "A GUID that identifies the service plan. For a complete list of GUIDs and their equivalent friendly service names, see [Product names and service plan identifiers for licensing](/azure/active-directory/enterprise-users/licensing-service-plan-reference).",
						Computed:            true,
					},
				},
			},
			"birthday": schema.StringAttribute{
				MarkdownDescription: "The birthday of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. <br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"business_phones": schema.ListAttribute{
				MarkdownDescription: "The telephone numbers for the user. NOTE: Although this is a string collection, only one number can be set for this property. Read-only for users synced from on-premises directory. <br><br>Returned by default. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"city": schema.StringAttribute{
				MarkdownDescription: "The city in which the user is located. Maximum length is 128 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"company_name": schema.StringAttribute{
				MarkdownDescription: "The company name which the user is associated. This property can be useful for describing the company that an external user comes from. The maximum length is 64 characters.<br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"consent_provided_for_minor": schema.StringAttribute{
				MarkdownDescription: "Sets whether consent has been obtained for minors. Allowed values: `null`, `Granted`, `Denied` and `NotRequired`. Refer to the [legal age group property definitions](#legal-age-group-property-definitions) for further information. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, and `in`).",
				Computed:            true,
			},
			"country": schema.StringAttribute{
				MarkdownDescription: "The country/region in which the user is located; for example, `US` or `UK`. Maximum length is 128 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"created_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the user was created, in ISO 8601 format and in UTC time. The value cannot be modified and is automatically populated when the entity is created. Nullable. For on-premises users, the value represents when they were first created in Azure AD. Property is `null` for some users created before June 2018 and on-premises users that were synced to Azure AD before June 2018. Read-only. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`).",
				Computed:            true,
			},
			"creation_type": schema.StringAttribute{
				MarkdownDescription: "Indicates whether the user account was created through one of the following methods: <br/> <ul><li>As a regular school or work account (`null`). <li>As an external account (`Invitation`). <li>As a local account for an Azure Active Directory B2C tenant (`LocalAccount`). <li>Through self-service sign-up by an internal user using email verification (`EmailVerified`). <li>Through self-service sign-up by an external user signing up through a link that is part of a user flow (`SelfServiceSignUp`).</ul> <br>Read-only.<br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
				Computed:            true,
			},
			"deleted_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time the user was deleted. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`).",
				Computed:            true,
			},
			"department": schema.StringAttribute{
				MarkdownDescription: "The name for the department in which the user works. Maximum length is 64 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`, and `eq` on `null` values).",
				Computed:            true,
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial and last name. This property is required when a user is created and it cannot be cleared during updates. Maximum length is 256 characters. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values), `$orderBy`, and `$search`.",
				Computed:            true,
			},
			"employee_hire_date": schema.StringAttribute{
				MarkdownDescription: "The date and time when the user was hired or will start work in case of a future hire. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`).",
				Computed:            true,
			},
			"employee_leave_date_time": schema.StringAttribute{
				MarkdownDescription: "The date and time when the user left or will leave the organization. <br><br>To read this property, the calling app must be assigned the *User-LifeCycleInfo.Read.All* permission. To write this property, the calling app must be assigned the *User.Read.All* and *User-LifeCycleInfo.ReadWrite.All* permissions. To read this property in delegated scenarios, the admin needs one of the following Azure AD roles: *Lifecycle Workflows Administrator*, *Global Reader*, or *Global Administrator*. To write this property in delegated scenarios, the admin needs the *Global Administrator* role. <br><br>Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`). <br><br>For more information, see [Configure the employeeLeaveDateTime property for a user](/graph/tutorial-lifecycle-workflows-set-employeeleavedatetime).",
				Computed:            true,
			},
			"employee_id": schema.StringAttribute{
				MarkdownDescription: "The employee identifier assigned to the user by the organization. The maximum length is 16 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"employee_org_data": schema.SingleNestedAttribute{
				MarkdownDescription: "Represents organization data (e.g. division and costCenter) associated with a user. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"division": schema.StringAttribute{
						MarkdownDescription: "The name of the division in which the user works. <br><br>Returned only on `$select`. Supports `$filter`.",
						Computed:            true,
					},
					"cost_center": schema.StringAttribute{
						MarkdownDescription: "The cost center associated with the user. <br><br>Returned only on `$select`. Supports `$filter`.",
						Computed:            true,
					},
				},
			},
			"employee_type": schema.StringAttribute{
				MarkdownDescription: "Captures enterprise worker type. For example, `Employee`, `Contractor`, `Consultant`, or `Vendor`. Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`, `startsWith`).",
				Computed:            true,
			},
			"external_user_state": schema.StringAttribute{
				MarkdownDescription: "For an external user invited to the tenant using the [invitation API](../api/invitation-post.md), this property represents the invited user's invitation status. For invited users, the state can be `PendingAcceptance` or `Accepted`, or `null` for all other users. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `in`).",
				Computed:            true,
			},
			"external_user_state_change_date_time": schema.StringAttribute{
				MarkdownDescription: "Shows the timestamp for the latest change to the **externalUserState** property. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `in`).",
				Computed:            true,
			},
			"fax_number": schema.StringAttribute{
				MarkdownDescription: "The fax number of the user. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"given_name": schema.StringAttribute{
				MarkdownDescription: "The given name (first name) of the user. Maximum length is 64 characters. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"hire_date": schema.StringAttribute{
				MarkdownDescription: "The hire date of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. <br><br>Returned only on `$select`. <br> **Note:** This property is specific to SharePoint Online. We recommend using the native **employeeHireDate** property to set and update hire date values using Microsoft Graph APIs.",
				Computed:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for the user. Should be treated as an opaque identifier. Inherited from [directoryObject](directoryobject.md). Key. Not nullable. Read-only. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
				Optional:            true,
				Computed:            true,
			},
			"identities": schema.SingleNestedAttribute{
				MarkdownDescription: "Represents the identities that can be used to sign in to this user account. An identity can be provided by Microsoft (also known as a local account), by organizations, or by social identity providers such as Facebook, Google, and Microsoft, and tied to a user account. May contain multiple items with the same **signInType** value. <br><br>Returned only on `$select`. Supports `$filter` (`eq`) including on `null` values, only where the **signInType** is not `userPrincipalName`.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"issuer": schema.StringAttribute{
						MarkdownDescription: "Specifies the issuer of the identity, for example `facebook.com`.<br>For local accounts (where **signInType** is not `federated`), this property is the local B2C tenant default domain name, for example `contoso.onmicrosoft.com`.<br>For external users from other Azure AD organization, this will be the domain of the federated organization, for example `contoso.com`.<br><br>Supports `$filter`. 512 character limit.",
						Computed:            true,
					},
					"issuer_assigned_id": schema.StringAttribute{
						MarkdownDescription: "Specifies the unique identifier assigned to the user by the issuer. The combination of **issuer** and **issuerAssignedId** must be unique within the organization. Represents the sign-in name for the user, when **signInType** is set to `emailAddress` or `userName` (also known as local accounts).<br>When **signInType** is set to: <ul><li>`emailAddress`, (or a custom string that starts with `emailAddress` like `emailAddress1`) **issuerAssignedId** must be a valid email address</li><li>`userName`, **issuerAssignedId** must begin with alphabetical character or number, and can only contain alphanumeric characters and the following symbols: - or _</li></ul>Supports `$filter`. 64 character limit.",
						Computed:            true,
					},
					"sign_in_type": schema.StringAttribute{
						MarkdownDescription: "Specifies the user sign-in types in your directory, such as `emailAddress`, `userName`, `federated`, or `userPrincipalName`. `federated` represents a unique identifier for a user from an issuer, that can be in any format chosen by the issuer. Setting or updating a `userPrincipalName` identity will update the value of the **userPrincipalName** property on the user object. The validations performed on the `userPrincipalName` property on the user object, for example, verified domains and acceptable characters, will be performed when setting or updating a `userPrincipalName` identity. Additional validation is enforced on **issuerAssignedId** when the sign-in type is set to `emailAddress` or `userName`. This property can also be set to any custom string.",
						Computed:            true,
					},
				},
			},
			"im_addresses": schema.ListAttribute{
				MarkdownDescription: "The instant message voice over IP (VOIP) session initiation protocol (SIP) addresses for the user. Read-only. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"interests": schema.ListAttribute{
				MarkdownDescription: "A list for the user to describe their interests. <br><br>Returned only on `$select`.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"is_resource_account": schema.BoolAttribute{
				MarkdownDescription: "Do not use – reserved for future use.",
				Computed:            true,
			},
			"job_title": schema.StringAttribute{
				MarkdownDescription: "The user's job title. Maximum length is 128 characters. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not` , `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"last_password_change_date_time": schema.StringAttribute{
				MarkdownDescription: "The time when this Azure AD user last changed their password or when their password was created, whichever date the latest action was performed. The date and time information uses ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. <br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"legal_age_group_classification": schema.StringAttribute{
				MarkdownDescription: "Used by enterprise applications to determine the legal age group of the user. This property is read-only and calculated based on **ageGroup** and **consentProvidedForMinor** properties. Allowed values: `null`, `MinorWithOutParentalConsent`, `MinorWithParentalConsent`, `MinorNoParentalConsentRequired`, `NotAdult` and `Adult`. Refer to the [legal age group property definitions](#legal-age-group-property-definitions) for further information. <br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"license_assignment_states": schema.SingleNestedAttribute{
				MarkdownDescription: "State of license assignments for this user. Also indicates licenses that are directly-assigned and those that the user has inherited through group memberships. Read-only. <br><br>Returned only on `$select`.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"assigned_by_group": schema.StringAttribute{
						MarkdownDescription: "Indicates whether the license is directly-assigned or inherited from a group. If directly-assigned, this field is `null`; if inherited through a group membership, this field contains the ID of the group. Read-Only.",
						Computed:            true,
					},
					"disabled_plans": schema.ListAttribute{
						MarkdownDescription: "The service plans that are disabled in this assignment. Read-Only.",
						Computed:            true,
						ElementType:         types.StringType,
					},
					"error": schema.StringAttribute{
						MarkdownDescription: "License assignment failure error. If the license is assigned successfully, this field will be Null. Read-Only. The possible values are `CountViolation`, `MutuallyExclusiveViolation`, `DependencyViolation`, `ProhibitedInUsageLocationViolation`, `UniquenessViolation`, and `Other`. For more information on how to identify and resolve license assignment errors see [here](/azure/active-directory/users-groups-roles/licensing-groups-resolve-problems).",
						Computed:            true,
					},
					"last_updated_date_time": schema.StringAttribute{
						MarkdownDescription: "The timestamp when the state of the license assignment was last updated.",
						Computed:            true,
					},
					"sku_id": schema.StringAttribute{
						MarkdownDescription: "The unique identifier for the SKU. Read-Only.",
						Computed:            true,
					},
					"state": schema.StringAttribute{
						MarkdownDescription: "Indicate the current state of this assignment. Read-Only. The possible values are `Active`, `ActiveWithError`, `Disabled`, and `Error`.",
						Computed:            true,
					},
				},
			},
			"mail": schema.StringAttribute{
				MarkdownDescription: "The SMTP address for the user, for example, `jeff@contoso.onmicrosoft.com`. Changes to this property will also update the user's **proxyAddresses** collection to include the value as an SMTP address. This property cannot contain accent characters. <br/> **NOTE:** We do not recommend updating this property for Azure AD B2C user profiles. Use the **otherMails** property instead. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, `endsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"mail_nickname": schema.StringAttribute{
				MarkdownDescription: "The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"mobile_phone": schema.StringAttribute{
				MarkdownDescription: "The primary cellular telephone number for the user. Read-only for users synced from on-premises directory. Maximum length is 64 characters. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"my_site": schema.StringAttribute{
				MarkdownDescription: "The URL for the user's personal site. <br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"office_location": schema.StringAttribute{
				MarkdownDescription: "The office location in the user's place of business. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"on_premises_distinguished_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises Active Directory `distinguished name` or `DN`. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect. Read-only. <br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"on_premises_domain_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises `domainFQDN`, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect. Read-only. <br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"on_premises_extension_attributes": schema.SingleNestedAttribute{
				MarkdownDescription: "Contains extensionAttributes1-15 for the user. These extension attributes are also known as Exchange custom attributes 1-15. <br><li>For an **onPremisesSyncEnabled** user, the source of authority for this set of properties is the on-premises and is read-only. </li><li>For a cloud-only user (where **onPremisesSyncEnabled** is `false`), these properties can be set during creation or update of a user object.  </li><li>For a cloud-only user previously synced from on-premises Active Directory, these properties are read-only in Microsoft Graph but can be fully managed through the Exchange Admin Center or the Exchange Online V2 module in PowerShell.</li><br> Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `in`).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"extension_attribute_1": schema.StringAttribute{
						MarkdownDescription: "First customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_2": schema.StringAttribute{
						MarkdownDescription: "Second customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_3": schema.StringAttribute{
						MarkdownDescription: "Third customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_4": schema.StringAttribute{
						MarkdownDescription: "Fourth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_5": schema.StringAttribute{
						MarkdownDescription: "Fifth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_6": schema.StringAttribute{
						MarkdownDescription: "Sixth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_7": schema.StringAttribute{
						MarkdownDescription: "Seventh customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_8": schema.StringAttribute{
						MarkdownDescription: "Eighth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_9": schema.StringAttribute{
						MarkdownDescription: "Ninth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_10": schema.StringAttribute{
						MarkdownDescription: "Tenth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_11": schema.StringAttribute{
						MarkdownDescription: "Eleventh customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_12": schema.StringAttribute{
						MarkdownDescription: "Twelfth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_13": schema.StringAttribute{
						MarkdownDescription: "Thirteenth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_14": schema.StringAttribute{
						MarkdownDescription: "Fourteenth customizable extension attribute.",
						Optional:            true,
					},
					"extension_attribute_15": schema.StringAttribute{
						MarkdownDescription: "Fifteenth customizable extension attribute.",
						Optional:            true,
					},
				},
			},
			"on_premises_immutable_id": schema.StringAttribute{
				MarkdownDescription: "This property is used to associate an on-premises Active Directory user account to their Azure AD user object. This property must be specified when creating a new user account in the Graph if you are using a federated domain for the user's **userPrincipalName** (UPN) property. **NOTE:** The **$** and **\\_** characters cannot be used when specifying this property. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`)..",
				Computed:            true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				MarkdownDescription: "Indicates the last time at which the object was synced with the on-premises directory; for example: `2013-02-16T03:04:54Z`. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is `2014-01-01T00:00:00Z`. Read-only. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`).",
				Computed:            true,
			},
			"on_premises_provisioning_errors": schema.SingleNestedAttribute{
				MarkdownDescription: "Errors when using Microsoft synchronization product during provisioning. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `not`, `ge`, `le`).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"category": schema.StringAttribute{
						MarkdownDescription: "Category of the provisioning error. Note: Currently, there is only one possible value. Possible value: *PropertyConflict* - indicates a property value is not unique. Other objects contain the same value for the property.",
						Optional:            true,
					},
					"occurred_date_time": schema.StringAttribute{
						MarkdownDescription: "The date and time at which the error occurred.",
						Optional:            true,
					},
					"property_causing_error": schema.StringAttribute{
						MarkdownDescription: "Name of the directory property causing the error. Current possible values: *UserPrincipalName* or *ProxyAddress*",
						Optional:            true,
					},
					"value": schema.StringAttribute{
						MarkdownDescription: "Value of the property causing the error.",
						Optional:            true,
					},
				},
			},
			"on_premises_sam_account_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises `samAccountName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect. Read-only. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`).",
				Computed:            true,
			},
			"on_premises_security_identifier": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. <br><br>Returned only on `$select`.  Supports `$filter` (`eq` including on `null` values).",
				Computed:            true,
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				MarkdownDescription: "`true` if this user object is currently being synced from an on-premises Active Directory (AD); otherwise the user isn't being synced and can be managed in Azure Active Directory (Azure AD). Read-only. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `in`, and `eq` on `null` values).",
				Computed:            true,
			},
			"on_premises_user_principal_name": schema.StringAttribute{
				MarkdownDescription: "Contains the on-premises `userPrincipalName` synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Azure Active Directory via Azure AD Connect. Read-only. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`).",
				Computed:            true,
			},
			"other_mails": schema.ListAttribute{
				MarkdownDescription: "A list of additional email addresses for the user; for example: `[\"bob@contoso.com\", \"Robert@fabrikam.com\"]`. <br>NOTE: This property cannot contain accent characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `not`, `ge`, `le`, `in`, `startsWith`, `endsWith`, `/$count eq 0`, `/$count ne 0`).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"password_policies": schema.StringAttribute{
				MarkdownDescription: "Specifies password policies for the user. This value is an enumeration with one possible value being `DisableStrongPassword`, which allows weaker passwords than the default policy to be specified. `DisablePasswordExpiration` can also be specified. The two may be specified together; for example: `DisablePasswordExpiration, DisableStrongPassword`. <br><br>Returned only on `$select`. For more information on the default password policies, see [Azure AD pasword policies](/azure/active-directory/authentication/concept-sspr-policy#password-policies-that-only-apply-to-cloud-user-accounts). Supports `$filter` (`ne`, `not`, and `eq` on `null` values).",
				Computed:            true,
			},
			"password_profile": schema.SingleNestedAttribute{
				MarkdownDescription: "Specifies the password profile for the user. The profile contains the user’s password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the **passwordPolicies** property. By default, a strong password is required. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `in`, and `eq` on `null` values).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"force_change_password_next_sign_in": schema.BoolAttribute{
						MarkdownDescription: "`true` if the user must change her password on the next login; otherwise `false`.",
						Optional:            true,
					},
					"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
						MarkdownDescription: "If `true`, at next sign-in, the user must perform a multi-factor authentication (MFA) before being forced to change their password. The behavior is identical to **forceChangePasswordNextSignIn** except that the user is required to first perform a multi-factor authentication before password change. After a password change, this property will be automatically reset to `false`. If not set, default is `false`.",
						Optional:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "The password for the user. This property is required when a user is created. It can be updated, but the user will be required to change the password on the next login. The password must satisfy minimum requirements as specified by the user’s **passwordPolicies** property. By default, a strong password is required.",
						Optional:            true,
					},
				},
			},
			"past_projects": schema.ListAttribute{
				MarkdownDescription: "A list for the user to enumerate their past projects. <br><br>Returned only on `$select`.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"postal_code": schema.StringAttribute{
				MarkdownDescription: "The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code. Maximum length is 40 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"preferred_data_location": schema.StringAttribute{
				MarkdownDescription: "The preferred data location for the user. For more information, see [OneDrive Online Multi-Geo](/sharepoint/dev/solution-guidance/multigeo-introduction).",
				Computed:            true,
			},
			"preferred_language": schema.StringAttribute{
				MarkdownDescription: "The preferred language for the user. Should follow ISO 639-1 Code; for example `en-US`. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values)",
				Computed:            true,
			},
			"preferred_name": schema.StringAttribute{
				MarkdownDescription: "The preferred name for the user. **Not Supported. This attribute returns an empty string.**<br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"provisioned_plans": schema.SingleNestedAttribute{
				MarkdownDescription: "The plans that are provisioned for the user. Read-only. Not nullable. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `not`, `ge`, `le`).",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"capability_status": schema.StringAttribute{
						MarkdownDescription: "For example, “Enabled”.",
						Optional:            true,
					},
					"provisioning_status": schema.StringAttribute{
						MarkdownDescription: "For example, “Success”.",
						Optional:            true,
					},
					"service": schema.StringAttribute{
						MarkdownDescription: "The name of the service; for example, “AccessControlS2S”",
						Optional:            true,
					},
				},
			},
			"proxy_addresses": schema.ListAttribute{
				MarkdownDescription: "For example: `[\"SMTP: bob@contoso.com\", \"smtp: bob@sales.contoso.com\"]`. Changes to the **mail** property will also update this collection to include the value as an SMTP address. For more information, see [mail and proxyAddresses properties](#mail-and-proxyaddresses-properties). The proxy address prefixed with `SMTP` (capitalized) is the primary proxy address while those prefixed with `smtp` are the secondary proxy addresses. For Azure AD B2C accounts, this property has a limit of ten unique addresses. Read-only in Microsoft Graph; you can update this property only through the [Microsoft 365 admin center](/exchange/recipients-in-exchange-online/manage-user-mailboxes/add-or-remove-email-addresses). Not nullable. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`, `endsWith`, `/$count eq 0`, `/$count ne 0`).",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"responsibilities": schema.ListAttribute{
				MarkdownDescription: "A list for the user to enumerate their responsibilities. <br><br>Returned only on `$select`.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"schools": schema.ListAttribute{
				MarkdownDescription: "A list for the user to enumerate the schools they have attended. <br><br>Returned only on `$select`.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"security_identifier": schema.StringAttribute{
				MarkdownDescription: "Security identifier (SID) of the user, used in Windows scenarios. <br><br>Read-only. Returned by default. <br>Supports `$select` and `$filter` (`eq`, `not`, `ge`, `le`, `startsWith`).",
				Computed:            true,
			},
			"show_in_address_list": schema.BoolAttribute{
				MarkdownDescription: "**Do not use in Microsoft Graph. Manage this property through the Microsoft 365 admin center instead.** Represents whether the user should be included in the Outlook global address list. See [Known issue](/graph/known-issues#showinaddresslist-property-is-out-of-sync-with-microsoft-exchange).",
				Computed:            true,
			},
			"sign_in_sessions_valid_from_date_time": schema.StringAttribute{
				MarkdownDescription: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications will get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).  If this happens, the application will need to acquire a new refresh token by making a request to the authorize endpoint. Read-only. Use [revokeSignInSessions](../api/user-revokesigninsessions.md) to reset. <br><br>Returned only on `$select`.",
				Computed:            true,
			},
			"skills": schema.ListAttribute{
				MarkdownDescription: "A list for the user to enumerate their skills. <br><br>Returned only on `$select`.",
				Computed:            true,
				ElementType:         types.StringType,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state or province in the user's address. Maximum length is 128 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"street_address": schema.StringAttribute{
				MarkdownDescription: "The street address of the user's place of business. Maximum length is 1024 characters. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"surname": schema.StringAttribute{
				MarkdownDescription: "The user's surname (family name or last name). Maximum length is 64 characters. <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"usage_location": schema.StringAttribute{
				MarkdownDescription: "A two letter country code (ISO standard 3166). Required for users that will be assigned licenses due to legal requirement to check for availability of services in countries.  Examples include: `US`, `JP`, and `GB`. Not nullable. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, and `eq` on `null` values).",
				Computed:            true,
			},
			"user_principal_name": schema.StringAttribute{
				MarkdownDescription: "The user principal name (UPN) of the user. The UPN is an Internet-style login name for the user based on the Internet standard RFC 822. By convention, this should map to the user's email name. The general format is alias@domain, where domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the **verifiedDomains** property of [organization](organization.md).<br>NOTE: This property cannot contain accent characters. Only the following characters are allowed `A - Z`, `a - z`, `0 - 9`, ` ' . - _ ! # ^ ~`. For the complete list of allowed characters, see [username policies](/azure/active-directory/authentication/concept-sspr-policy#userprincipalname-policies-that-apply-to-all-user-accounts). <br><br>Returned by default. Supports `$filter` (`eq`, `ne`, `not`, `ge`, `le`, `in`, `startsWith`, `endsWith`) and `$orderBy`.",
				Optional:            true,
				Computed:            true,
			},
			"user_type": schema.StringAttribute{
				MarkdownDescription: "A string value that can be used to classify user types in your directory, such as `Member` and `Guest`. <br><br>Returned only on `$select`. Supports `$filter` (`eq`, `ne`, `not`, `in`, and `eq` on `null` values). **NOTE:** For more information about the permissions for member and guest users, see [What are the default user permissions in Azure Active Directory?](/azure/active-directory/fundamentals/users-default-permissions?context=graph/context#member-and-guest-users)",
				Computed:            true,
			},
		},
	}
}

type userDataSourceModel struct {
	AboutMe                         types.String                                      `tfsdk:"about_me"`
	AccountEnabled                  types.Bool                                        `tfsdk:"account_enabled"`
	AgeGroup                        types.String                                      `tfsdk:"age_group"`
	AssignedLicenses                []userAssignedLicensesDataSourceModel             `tfsdk:"assigned_licenses"`
	AssignedPlans                   []userAssignedPlansDataSourceModel                `tfsdk:"assigned_plans"`
	Birthday                        types.String                                      `tfsdk:"birthday"`
	BusinessPhones                  []types.String                                    `tfsdk:"business_phones"`
	City                            types.String                                      `tfsdk:"city"`
	CompanyName                     types.String                                      `tfsdk:"company_name"`
	ConsentProvidedForMinor         types.String                                      `tfsdk:"consent_provided_for_minor"`
	Country                         types.String                                      `tfsdk:"country"`
	CreatedDateTime                 types.String                                      `tfsdk:"created_date_time"`
	CreationType                    types.String                                      `tfsdk:"creation_type"`
	DeletedDateTime                 types.String                                      `tfsdk:"deleted_date_time"`
	Department                      types.String                                      `tfsdk:"department"`
	DisplayName                     types.String                                      `tfsdk:"display_name"`
	EmployeeHireDate                types.String                                      `tfsdk:"employee_hire_date"`
	EmployeeLeaveDateTime           types.String                                      `tfsdk:"employee_leave_date_time"`
	EmployeeId                      types.String                                      `tfsdk:"employee_id"`
	EmployeeOrgData                 *userEmployeeOrgDataDataSourceModel               `tfsdk:"employee_org_data"`
	EmployeeType                    types.String                                      `tfsdk:"employee_type"`
	ExternalUserState               types.String                                      `tfsdk:"external_user_state"`
	ExternalUserStateChangeDateTime types.String                                      `tfsdk:"external_user_state_change_date_time"`
	FaxNumber                       types.String                                      `tfsdk:"fax_number"`
	GivenName                       types.String                                      `tfsdk:"given_name"`
	HireDate                        types.String                                      `tfsdk:"hire_date"`
	Id                              types.String                                      `tfsdk:"id"`
	Identities                      []userIdentitiesDataSourceModel                   `tfsdk:"identities"`
	ImAddresses                     []types.String                                    `tfsdk:"im_addresses"`
	Interests                       []types.String                                    `tfsdk:"interests"`
	IsResourceAccount               types.Bool                                        `tfsdk:"is_resource_account"`
	JobTitle                        types.String                                      `tfsdk:"job_title"`
	LastPasswordChangeDateTime      types.String                                      `tfsdk:"last_password_change_date_time"`
	LegalAgeGroupClassification     types.String                                      `tfsdk:"legal_age_group_classification"`
	LicenseAssignmentStates         []userLicenseAssignmentStatesDataSourceModel      `tfsdk:"license_assignment_states"`
	Mail                            types.String                                      `tfsdk:"mail"`
	MailNickname                    types.String                                      `tfsdk:"mail_nickname"`
	MobilePhone                     types.String                                      `tfsdk:"mobile_phone"`
	MySite                          types.String                                      `tfsdk:"my_site"`
	OfficeLocation                  types.String                                      `tfsdk:"office_location"`
	OnPremisesDistinguishedName     types.String                                      `tfsdk:"on_premises_distinguished_name"`
	OnPremisesDomainName            types.String                                      `tfsdk:"on_premises_domain_name"`
	OnPremisesExtensionAttributes   *userOnPremisesExtensionAttributesDataSourceModel `tfsdk:"on_premises_extension_attributes"`
	OnPremisesImmutableId           types.String                                      `tfsdk:"on_premises_immutable_id"`
	OnPremisesLastSyncDateTime      types.String                                      `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesProvisioningErrors    []userOnPremisesProvisioningErrorsDataSourceModel `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName        types.String                                      `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier    types.String                                      `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled           types.Bool                                        `tfsdk:"on_premises_sync_enabled"`
	OnPremisesUserPrincipalName     types.String                                      `tfsdk:"on_premises_user_principal_name"`
	OtherMails                      []types.String                                    `tfsdk:"other_mails"`
	PasswordPolicies                types.String                                      `tfsdk:"password_policies"`
	PasswordProfile                 *userPasswordProfileDataSourceModel               `tfsdk:"password_profile"`
	PastProjects                    []types.String                                    `tfsdk:"past_projects"`
	PostalCode                      types.String                                      `tfsdk:"postal_code"`
	PreferredDataLocation           types.String                                      `tfsdk:"preferred_data_location"`
	PreferredLanguage               types.String                                      `tfsdk:"preferred_language"`
	PreferredName                   types.String                                      `tfsdk:"preferred_name"`
	ProvisionedPlans                []userProvisionedPlansDataSourceModel             `tfsdk:"provisioned_plans"`
	ProxyAddresses                  []types.String                                    `tfsdk:"proxy_addresses"`
	Responsibilities                []types.String                                    `tfsdk:"responsibilities"`
	Schools                         []types.String                                    `tfsdk:"schools"`
	SecurityIdentifier              types.String                                      `tfsdk:"security_identifier"`
	ShowInAddressList               types.Bool                                        `tfsdk:"show_in_address_list"`
	SignInSessionsValidFromDateTime types.String                                      `tfsdk:"sign_in_sessions_valid_from_date_time"`
	Skills                          []types.String                                    `tfsdk:"skills"`
	State                           types.String                                      `tfsdk:"state"`
	StreetAddress                   types.String                                      `tfsdk:"street_address"`
	Surname                         types.String                                      `tfsdk:"surname"`
	UsageLocation                   types.String                                      `tfsdk:"usage_location"`
	UserPrincipalName               types.String                                      `tfsdk:"user_principal_name"`
	UserType                        types.String                                      `tfsdk:"user_type"`
}

type userAssignedLicensesDataSourceModel struct {
	DisabledPlans []types.String `tfsdk:"disabled_plans"`
	SkuId         types.String   `tfsdk:"sku_id"`
}

type userAssignedPlansDataSourceModel struct {
	AssignedDateTime types.String `tfsdk:"assigned_date_time"`
	CapabilityStatus types.String `tfsdk:"capability_status"`
	Service          types.String `tfsdk:"service"`
	ServicePlanId    types.String `tfsdk:"service_plan_id"`
}

type userEmployeeOrgDataDataSourceModel struct {
	Division   types.String `tfsdk:"division"`
	CostCenter types.String `tfsdk:"cost_center"`
}

type userIdentitiesDataSourceModel struct {
	Issuer           types.String `tfsdk:"issuer"`
	IssuerAssignedId types.String `tfsdk:"issuer_assigned_id"`
	SignInType       types.String `tfsdk:"sign_in_type"`
}

type userLicenseAssignmentStatesDataSourceModel struct {
	AssignedByGroup     types.String   `tfsdk:"assigned_by_group"`
	DisabledPlans       []types.String `tfsdk:"disabled_plans"`
	Error               types.String   `tfsdk:"error"`
	LastUpdatedDateTime types.String   `tfsdk:"last_updated_date_time"`
	SkuId               types.String   `tfsdk:"sku_id"`
	State               types.String   `tfsdk:"state"`
}

type userOnPremisesExtensionAttributesDataSourceModel struct {
	ExtensionAttribute1  types.String `tfsdk:"extension_attribute_1"`
	ExtensionAttribute2  types.String `tfsdk:"extension_attribute_2"`
	ExtensionAttribute3  types.String `tfsdk:"extension_attribute_3"`
	ExtensionAttribute4  types.String `tfsdk:"extension_attribute_4"`
	ExtensionAttribute5  types.String `tfsdk:"extension_attribute_5"`
	ExtensionAttribute6  types.String `tfsdk:"extension_attribute_6"`
	ExtensionAttribute7  types.String `tfsdk:"extension_attribute_7"`
	ExtensionAttribute8  types.String `tfsdk:"extension_attribute_8"`
	ExtensionAttribute9  types.String `tfsdk:"extension_attribute_9"`
	ExtensionAttribute10 types.String `tfsdk:"extension_attribute_10"`
	ExtensionAttribute11 types.String `tfsdk:"extension_attribute_11"`
	ExtensionAttribute12 types.String `tfsdk:"extension_attribute_12"`
	ExtensionAttribute13 types.String `tfsdk:"extension_attribute_13"`
	ExtensionAttribute14 types.String `tfsdk:"extension_attribute_14"`
	ExtensionAttribute15 types.String `tfsdk:"extension_attribute_15"`
}

type userOnPremisesProvisioningErrorsDataSourceModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

type userPasswordProfileDataSourceModel struct {
	ForceChangePasswordNextSignIn        types.Bool   `tfsdk:"force_change_password_next_sign_in"`
	ForceChangePasswordNextSignInWithMfa types.Bool   `tfsdk:"force_change_password_next_sign_in_with_mfa"`
	Password                             types.String `tfsdk:"password"`
}

type userProvisionedPlansDataSourceModel struct {
	CapabilityStatus   types.String `tfsdk:"capability_status"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	Service            types.String `tfsdk:"service"`
}

// Read refreshes the Terraform state with the latest data.
func (d *userDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state userDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: UserProperties[:],
		},
	}

	var result models.Userable
	var err error
	if !state.Id.IsNull() {
		result, err = d.client.UsersById(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else if !state.UserPrincipalName.IsNull() {
		result, err = d.client.UsersById(state.UserPrincipalName.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"Either `id` or `user_principal_name` must be supplied.",
		)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting user",
			err.Error(),
		)
		return
	}

	if result.GetAboutMe() != nil {
		state.AboutMe = types.StringValue(*result.GetAboutMe())
	}
	if result.GetAccountEnabled() != nil {
		state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	}
	if result.GetAgeGroup() != nil {
		state.AgeGroup = types.StringValue(*result.GetAgeGroup())
	}
	for _, value := range result.GetAssignedLicenses() {
		assignedLicenses := new(userAssignedLicensesDataSourceModel)

		for _, value := range value.GetDisabledPlans() {
			assignedLicenses.DisabledPlans = append(assignedLicenses.DisabledPlans, types.StringValue(value.String()))
		}
		if value.GetSkuId() != nil {
			assignedLicenses.SkuId = types.StringValue(value.GetSkuId().String())
		}
		state.AssignedLicenses = append(state.AssignedLicenses, *assignedLicenses)
	}
	for _, value := range result.GetAssignedPlans() {
		assignedPlans := new(userAssignedPlansDataSourceModel)

		if value.GetAssignedDateTime() != nil {
			assignedPlans.AssignedDateTime = types.StringValue(value.GetAssignedDateTime().String())
		}
		if value.GetCapabilityStatus() != nil {
			assignedPlans.CapabilityStatus = types.StringValue(*value.GetCapabilityStatus())
		}
		if value.GetService() != nil {
			assignedPlans.Service = types.StringValue(*value.GetService())
		}
		if value.GetServicePlanId() != nil {
			assignedPlans.ServicePlanId = types.StringValue(value.GetServicePlanId().String())
		}
		state.AssignedPlans = append(state.AssignedPlans, *assignedPlans)
	}
	if result.GetBirthday() != nil {
		state.Birthday = types.StringValue(result.GetBirthday().String())
	}
	for _, value := range result.GetBusinessPhones() {
		state.BusinessPhones = append(state.BusinessPhones, types.StringValue(value))
	}
	if result.GetCity() != nil {
		state.City = types.StringValue(*result.GetCity())
	}
	if result.GetCompanyName() != nil {
		state.CompanyName = types.StringValue(*result.GetCompanyName())
	}
	if result.GetConsentProvidedForMinor() != nil {
		state.ConsentProvidedForMinor = types.StringValue(*result.GetConsentProvidedForMinor())
	}
	if result.GetCountry() != nil {
		state.Country = types.StringValue(*result.GetCountry())
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	}
	if result.GetCreationType() != nil {
		state.CreationType = types.StringValue(*result.GetCreationType())
	}
	if result.GetDeletedDateTime() != nil {
		state.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	}
	if result.GetDepartment() != nil {
		state.Department = types.StringValue(*result.GetDepartment())
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	}
	if result.GetEmployeeHireDate() != nil {
		state.EmployeeHireDate = types.StringValue(result.GetEmployeeHireDate().String())
	}
	if result.GetEmployeeLeaveDateTime() != nil {
		state.EmployeeLeaveDateTime = types.StringValue(result.GetEmployeeLeaveDateTime().String())
	}
	if result.GetEmployeeId() != nil {
		state.EmployeeId = types.StringValue(*result.GetEmployeeId())
	}
	employeeOrgData := new(userEmployeeOrgDataDataSourceModel)
	if result.GetEmployeeOrgData() != nil {

		if result.GetEmployeeOrgData().GetDivision() != nil {
			state.EmployeeOrgData.Division = types.StringValue(*result.GetEmployeeOrgData().GetDivision())
		}
		if result.GetEmployeeOrgData().GetCostCenter() != nil {
			state.EmployeeOrgData.CostCenter = types.StringValue(*result.GetEmployeeOrgData().GetCostCenter())
		}
	}
	state.EmployeeOrgData = employeeOrgData
	if result.GetEmployeeType() != nil {
		state.EmployeeType = types.StringValue(*result.GetEmployeeType())
	}
	if result.GetExternalUserState() != nil {
		state.ExternalUserState = types.StringValue(*result.GetExternalUserState())
	}
	if result.GetExternalUserStateChangeDateTime() != nil {
		state.ExternalUserStateChangeDateTime = types.StringValue(result.GetExternalUserStateChangeDateTime().String())
	}
	if result.GetFaxNumber() != nil {
		state.FaxNumber = types.StringValue(*result.GetFaxNumber())
	}
	if result.GetGivenName() != nil {
		state.GivenName = types.StringValue(*result.GetGivenName())
	}
	if result.GetHireDate() != nil {
		state.HireDate = types.StringValue(result.GetHireDate().String())
	}
	if result.GetId() != nil {
		state.Id = types.StringValue(*result.GetId())
	}
	for _, value := range result.GetIdentities() {
		identities := new(userIdentitiesDataSourceModel)

		if value.GetIssuer() != nil {
			identities.Issuer = types.StringValue(*value.GetIssuer())
		}
		if value.GetIssuerAssignedId() != nil {
			identities.IssuerAssignedId = types.StringValue(*value.GetIssuerAssignedId())
		}
		if value.GetSignInType() != nil {
			identities.SignInType = types.StringValue(*value.GetSignInType())
		}
		state.Identities = append(state.Identities, *identities)
	}
	for _, value := range result.GetImAddresses() {
		state.ImAddresses = append(state.ImAddresses, types.StringValue(value))
	}
	for _, value := range result.GetInterests() {
		state.Interests = append(state.Interests, types.StringValue(value))
	}
	if result.GetIsResourceAccount() != nil {
		state.IsResourceAccount = types.BoolValue(*result.GetIsResourceAccount())
	}
	if result.GetJobTitle() != nil {
		state.JobTitle = types.StringValue(*result.GetJobTitle())
	}
	if result.GetLastPasswordChangeDateTime() != nil {
		state.LastPasswordChangeDateTime = types.StringValue(result.GetLastPasswordChangeDateTime().String())
	}
	if result.GetLegalAgeGroupClassification() != nil {
		state.LegalAgeGroupClassification = types.StringValue(*result.GetLegalAgeGroupClassification())
	}
	for _, value := range result.GetLicenseAssignmentStates() {
		licenseAssignmentStates := new(userLicenseAssignmentStatesDataSourceModel)

		if value.GetAssignedByGroup() != nil {
			licenseAssignmentStates.AssignedByGroup = types.StringValue(*value.GetAssignedByGroup())
		}
		for _, value := range value.GetDisabledPlans() {
			licenseAssignmentStates.DisabledPlans = append(licenseAssignmentStates.DisabledPlans, types.StringValue(value.String()))
		}
		if value.GetError() != nil {
			licenseAssignmentStates.Error = types.StringValue(*value.GetError())
		}
		if value.GetLastUpdatedDateTime() != nil {
			licenseAssignmentStates.LastUpdatedDateTime = types.StringValue(value.GetLastUpdatedDateTime().String())
		}
		if value.GetSkuId() != nil {
			licenseAssignmentStates.SkuId = types.StringValue(value.GetSkuId().String())
		}
		if value.GetState() != nil {
			licenseAssignmentStates.State = types.StringValue(*value.GetState())
		}
		state.LicenseAssignmentStates = append(state.LicenseAssignmentStates, *licenseAssignmentStates)
	}
	if result.GetMail() != nil {
		state.Mail = types.StringValue(*result.GetMail())
	}
	if result.GetMailNickname() != nil {
		state.MailNickname = types.StringValue(*result.GetMailNickname())
	}
	if result.GetMobilePhone() != nil {
		state.MobilePhone = types.StringValue(*result.GetMobilePhone())
	}
	if result.GetMySite() != nil {
		state.MySite = types.StringValue(*result.GetMySite())
	}
	if result.GetOfficeLocation() != nil {
		state.OfficeLocation = types.StringValue(*result.GetOfficeLocation())
	}
	if result.GetOnPremisesDistinguishedName() != nil {
		state.OnPremisesDistinguishedName = types.StringValue(*result.GetOnPremisesDistinguishedName())
	}
	if result.GetOnPremisesDomainName() != nil {
		state.OnPremisesDomainName = types.StringValue(*result.GetOnPremisesDomainName())
	}
	onPremisesExtensionAttributes := new(userOnPremisesExtensionAttributesDataSourceModel)
	if result.GetOnPremisesExtensionAttributes() != nil {

		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15() != nil {
			state.OnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15())
		}
	}
	state.OnPremisesExtensionAttributes = onPremisesExtensionAttributes
	if result.GetOnPremisesImmutableId() != nil {
		state.OnPremisesImmutableId = types.StringValue(*result.GetOnPremisesImmutableId())
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	}
	for _, value := range result.GetOnPremisesProvisioningErrors() {
		onPremisesProvisioningErrors := new(userOnPremisesProvisioningErrorsDataSourceModel)

		if value.GetCategory() != nil {
			onPremisesProvisioningErrors.Category = types.StringValue(*value.GetCategory())
		}
		if value.GetOccurredDateTime() != nil {
			onPremisesProvisioningErrors.OccurredDateTime = types.StringValue(value.GetOccurredDateTime().String())
		}
		if value.GetPropertyCausingError() != nil {
			onPremisesProvisioningErrors.PropertyCausingError = types.StringValue(*value.GetPropertyCausingError())
		}
		if value.GetValue() != nil {
			onPremisesProvisioningErrors.Value = types.StringValue(*value.GetValue())
		}
		state.OnPremisesProvisioningErrors = append(state.OnPremisesProvisioningErrors, *onPremisesProvisioningErrors)
	}
	if result.GetOnPremisesSamAccountName() != nil {
		state.OnPremisesSamAccountName = types.StringValue(*result.GetOnPremisesSamAccountName())
	}
	if result.GetOnPremisesSecurityIdentifier() != nil {
		state.OnPremisesSecurityIdentifier = types.StringValue(*result.GetOnPremisesSecurityIdentifier())
	}
	if result.GetOnPremisesSyncEnabled() != nil {
		state.OnPremisesSyncEnabled = types.BoolValue(*result.GetOnPremisesSyncEnabled())
	}
	if result.GetOnPremisesUserPrincipalName() != nil {
		state.OnPremisesUserPrincipalName = types.StringValue(*result.GetOnPremisesUserPrincipalName())
	}
	for _, value := range result.GetOtherMails() {
		state.OtherMails = append(state.OtherMails, types.StringValue(value))
	}
	if result.GetPasswordPolicies() != nil {
		state.PasswordPolicies = types.StringValue(*result.GetPasswordPolicies())
	}
	passwordProfile := new(userPasswordProfileDataSourceModel)
	if result.GetPasswordProfile() != nil {

		if result.GetPasswordProfile().GetForceChangePasswordNextSignIn() != nil {
			state.PasswordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
		}
		if result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa() != nil {
			state.PasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
		}
		if result.GetPasswordProfile().GetPassword() != nil {
			state.PasswordProfile.Password = types.StringValue(*result.GetPasswordProfile().GetPassword())
		}
	}
	state.PasswordProfile = passwordProfile
	for _, value := range result.GetPastProjects() {
		state.PastProjects = append(state.PastProjects, types.StringValue(value))
	}
	if result.GetPostalCode() != nil {
		state.PostalCode = types.StringValue(*result.GetPostalCode())
	}
	if result.GetPreferredDataLocation() != nil {
		state.PreferredDataLocation = types.StringValue(*result.GetPreferredDataLocation())
	}
	if result.GetPreferredLanguage() != nil {
		state.PreferredLanguage = types.StringValue(*result.GetPreferredLanguage())
	}
	if result.GetPreferredName() != nil {
		state.PreferredName = types.StringValue(*result.GetPreferredName())
	}
	for _, value := range result.GetProvisionedPlans() {
		provisionedPlans := new(userProvisionedPlansDataSourceModel)

		if value.GetCapabilityStatus() != nil {
			provisionedPlans.CapabilityStatus = types.StringValue(*value.GetCapabilityStatus())
		}
		if value.GetProvisioningStatus() != nil {
			provisionedPlans.ProvisioningStatus = types.StringValue(*value.GetProvisioningStatus())
		}
		if value.GetService() != nil {
			provisionedPlans.Service = types.StringValue(*value.GetService())
		}
		state.ProvisionedPlans = append(state.ProvisionedPlans, *provisionedPlans)
	}
	for _, value := range result.GetProxyAddresses() {
		state.ProxyAddresses = append(state.ProxyAddresses, types.StringValue(value))
	}
	for _, value := range result.GetResponsibilities() {
		state.Responsibilities = append(state.Responsibilities, types.StringValue(value))
	}
	for _, value := range result.GetSchools() {
		state.Schools = append(state.Schools, types.StringValue(value))
	}
	if result.GetSecurityIdentifier() != nil {
		state.SecurityIdentifier = types.StringValue(*result.GetSecurityIdentifier())
	}
	if result.GetShowInAddressList() != nil {
		state.ShowInAddressList = types.BoolValue(*result.GetShowInAddressList())
	}
	if result.GetSignInSessionsValidFromDateTime() != nil {
		state.SignInSessionsValidFromDateTime = types.StringValue(result.GetSignInSessionsValidFromDateTime().String())
	}
	for _, value := range result.GetSkills() {
		state.Skills = append(state.Skills, types.StringValue(value))
	}
	if result.GetState() != nil {
		state.State = types.StringValue(*result.GetState())
	}
	if result.GetStreetAddress() != nil {
		state.StreetAddress = types.StringValue(*result.GetStreetAddress())
	}
	if result.GetSurname() != nil {
		state.Surname = types.StringValue(*result.GetSurname())
	}
	if result.GetUsageLocation() != nil {
		state.UsageLocation = types.StringValue(*result.GetUsageLocation())
	}
	if result.GetUserPrincipalName() != nil {
		state.UserPrincipalName = types.StringValue(*result.GetUserPrincipalName())
	}
	if result.GetUserType() != nil {
		state.UserType = types.StringValue(*result.GetUserType())
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
