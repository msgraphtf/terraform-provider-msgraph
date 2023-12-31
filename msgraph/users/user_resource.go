package users

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &userResource{}
	_ resource.ResourceWithConfigure = &userResource{}
)

// NewUserResource is a helper function to simplify the provider implementation.
func NewUserResource() resource.Resource {
	return &userResource{}
}

// userResource is the resource implementation.
type userResource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (d *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

// Configure adds the provider configured client to the resource.
func (d *userResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the resource.
func (d *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
			},
			"deleted_date_time": schema.StringAttribute{
				Description: "Date and time when this object was deleted. Always null when the object hasn't been deleted.",
				Computed:    true,
			},
			"about_me": schema.StringAttribute{
				Description: "A freeform text entry field for the user to describe themselves. Returned only on $select.",
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
				Description: "The licenses that are assigned to the user, including inherited (group-based) licenses. This property doesn't differentiate between directly assigned and inherited licenses. Use the licenseAssignmentStates property to identify the directly assigned and inherited licenses.  Not nullable. Returned only on $select. Supports $filter (eq, not, /$count eq 0, /$count ne 0).",
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
			"birthday": schema.StringAttribute{
				Description: "The birthday of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.",
				Computed:    true,
			},
			"business_phones": schema.ListAttribute{
				Description: "The telephone numbers for the user. NOTE: Although it is a string collection, only one number can be set for this property. Read-only for users synced from the on-premises directory. Returned by default. Supports $filter (eq, not, ge, le, startsWith).",
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
				Description: "Sets whether consent was obtained for minors. Allowed values: null, Granted, Denied and NotRequired. Refer to the legal age group property definitions for further information. Returned only on $select. Supports $filter (eq, ne, not, and in).",
				Computed:    true,
			},
			"country": schema.StringAttribute{
				Description: "The country or region where the user is located; for example, US or UK. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Computed:    true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "The date and time the user was created, in ISO 8601 format and UTC. The value cannot be modified and is automatically populated when the entity is created. Nullable. For on-premises users, the value represents when they were first created in Microsoft Entra ID. Property is null for some users created before June 2018 and on-premises users that were synced to Microsoft Entra ID before June 2018. Read-only. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
				Computed:    true,
			},
			"creation_type": schema.StringAttribute{
				Description: "Indicates whether the user account was created through one of the following methods:  As a regular school or work account (null). As an external account (Invitation). As a local account for an Azure Active Directory B2C tenant (LocalAccount). Through self-service sign-up by an internal user using email verification (EmailVerified). Through self-service sign-up by an external user signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.Returned only on $select. Supports $filter (eq, ne, not, in).",
				Computed:    true,
			},
			"department": schema.StringAttribute{
				Description: "The name of the department in which the user works. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, and eq on null values).",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The name displayed in the address book for the user. This is usually the combination of the user's first name, middle initial, and last name. This property is required when a user is created and it cannot be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values), $orderby, and $search.",
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
			"hire_date": schema.StringAttribute{
				Description: "The hire date of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.  Note: This property is specific to SharePoint Online. We recommend using the native employeeHireDate property to set and update hire date values using Microsoft Graph APIs.",
				Computed:    true,
			},
			"identities": schema.ListNestedAttribute{
				Description: "Represents the identities that can be used to sign in to this user account. Microsoft (also known as a local account), organizations, or social identity providers such as Facebook, Google, and Microsoft can provide identity and tie it to a user account. It may contain multiple items with the same signInType value. Returned only on $select. Supports $filter (eq) including on null values, only where the signInType is not userPrincipalName.",
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
				Description: "The instant message voice-over IP (VOIP) session initiation protocol (SIP) addresses for the user. Read-only. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"interests": schema.ListAttribute{
				Description: "A list for the user to describe their interests. Returned only on $select.",
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
				Description: "The time when this Microsoft Entra user last changed their password or when their password was created, whichever date the latest action was performed. The date and time information uses ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.",
				Computed:    true,
			},
			"legal_age_group_classification": schema.StringAttribute{
				Description: "Used by enterprise applications to determine the legal age group of the user. This property is read-only and calculated based on ageGroup and consentProvidedForMinor properties. Allowed values: null, MinorWithOutParentalConsent, MinorWithParentalConsent, MinorNoParentalConsentRequired, NotAdult, and Adult. Refer to the legal age group property definitions for further information. Returned only on $select.",
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
				Description: "The SMTP address for the user, for example, jeff@contoso.onmicrosoft.com. Changes to this property update the user's proxyAddresses collection to include the value as an SMTP address. This property can't contain accent characters.  NOTE: We don't recommend updating this property for Azure AD B2C user profiles. Use the otherMails property instead. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith, and eq on null values).",
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
			"my_site": schema.StringAttribute{
				Description: "The URL for the user's site. Returned only on $select.",
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
				Description: "Contains extensionAttributes1-15 for the user. These extension attributes are also known as Exchange custom attributes 1-15. For an onPremisesSyncEnabled user, the source of authority for this set of properties is the on-premises and is read-only. For a cloud-only user (where onPremisesSyncEnabled is false), these properties can be set during the creation or update of a user object.  For a cloud-only user previously synced from on-premises Active Directory, these properties are read-only in Microsoft Graph but can be fully managed through the Exchange Admin Center or the Exchange Online V2 module in PowerShell. Returned only on $select. Supports $filter (eq, ne, not, in).",
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
			"past_projects": schema.ListAttribute{
				Description: "A list for the user to enumerate their past projects. Returned only on $select.",
				Computed:    true,
				ElementType: types.StringType,
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
				Description: "The preferred language for the user. The preferred language format is based on RFC 4646. The name is a combination of an ISO 639 two-letter lowercase culture code associated with the language and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: 'en-US', or 'es-ES'. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values)",
				Computed:    true,
			},
			"preferred_name": schema.StringAttribute{
				Description: "The preferred name for the user. Not Supported. This attribute returns an empty string.Returned only on $select.",
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
			"responsibilities": schema.ListAttribute{
				Description: "A list for the user to enumerate their responsibilities. Returned only on $select.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"schools": schema.ListAttribute{
				Description: "A list for the user to enumerate the schools they have attended. Returned only on $select.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"security_identifier": schema.StringAttribute{
				Description: "Security identifier (SID) of the user, used in Windows scenarios. Read-only. Returned by default. Supports $select and $filter (eq, not, ge, le, startsWith).",
				Computed:    true,
			},
			"service_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors published by a federated service describing a non-transient, service-specific error regarding the properties or link from a user object .  Supports $filter (eq, not, for isResolved and serviceInstance).",
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
				Description: "Any refresh tokens or sessions tokens (session cookies) issued before this time are invalid, and applications get an error when using an invalid refresh or sessions token to acquire a delegated access token (to access APIs such as Microsoft Graph).  If this happens, the application needs to acquire a new refresh token by requesting the authorized endpoint. Read-only. Use revokeSignInSessions to reset. Returned only on $select.",
				Computed:    true,
			},
			"skills": schema.ListAttribute{
				Description: "A list for the user to enumerate their skills. Returned only on $select.",
				Computed:    true,
				ElementType: types.StringType,
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
				Description: "A two-letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirements to check for availability of services in countries.  Examples include: US, JP, and GB. Not nullable. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Computed:    true,
			},
			"user_principal_name": schema.StringAttribute{
				Description: "The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this should map to the user's email name. The general format is alias@domain, where the domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization.NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith) and $orderby.",
				Optional:    true,
				Computed:    true,
			},
			"user_type": schema.StringAttribute{
				Description: "A string value that can be used to classify user types in your directory, such as Member and Guest. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). NOTE: For more information about the permissions for member and guest users, see What are the default user permissions in Microsoft Entra ID?",
				Computed:    true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan userModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var t time.Time
	var u uuid.UUID

	// Generate API request body from Plan
	requestBody := models.NewUser()

	if !plan.Id.IsUnknown() {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	} else {
	}

	if !plan.AccountEnabled.IsUnknown() {
		planAccountEnabled := plan.AccountEnabled.ValueBool()
		requestBody.SetAccountEnabled(&planAccountEnabled)
	}

	if !plan.AgeGroup.IsUnknown() {
		planAgeGroup := plan.AgeGroup.ValueString()
		requestBody.SetAgeGroup(&planAgeGroup)
	} else {
	}

	var planAssignedLicenses []models.AssignedLicenseable
	for _, i := range plan.AssignedLicenses.Elements() {
		assignedLicense := models.NewAssignedLicense()
		assignedLicenseModel := userAssignedLicensesModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &assignedLicenseModel)

		if len(assignedLicenseModel.DisabledPlans.Elements()) > 0 {
			var DisabledPlans []uuid.UUID
			for _, i := range assignedLicenseModel.DisabledPlans.Elements() {
				u, _ = uuid.Parse(i.String())
				DisabledPlans = append(DisabledPlans, u)
			}
			assignedLicense.SetDisabledPlans(DisabledPlans)
		}

		if !assignedLicenseModel.SkuId.IsUnknown() {
			planSkuId := assignedLicenseModel.SkuId.ValueString()
			u, _ = uuid.Parse(planSkuId)
			assignedLicense.SetSkuId(&u)
		}
	}
	requestBody.SetAssignedLicenses(planAssignedLicenses)

	var planAssignedPlans []models.AssignedPlanable
	for _, i := range plan.AssignedPlans.Elements() {
		assignedPlan := models.NewAssignedPlan()
		assignedPlanModel := userAssignedPlansModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &assignedPlanModel)

		if !assignedPlanModel.AssignedDateTime.IsUnknown() {
			planAssignedDateTime := assignedPlanModel.AssignedDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planAssignedDateTime)
			assignedPlan.SetAssignedDateTime(&t)
		}

		if !assignedPlanModel.CapabilityStatus.IsUnknown() {
			planCapabilityStatus := assignedPlanModel.CapabilityStatus.ValueString()
			assignedPlan.SetCapabilityStatus(&planCapabilityStatus)
		} else {
		}

		if !assignedPlanModel.Service.IsUnknown() {
			planService := assignedPlanModel.Service.ValueString()
			assignedPlan.SetService(&planService)
		} else {
		}

		if !assignedPlanModel.ServicePlanId.IsUnknown() {
			planServicePlanId := assignedPlanModel.ServicePlanId.ValueString()
			u, _ = uuid.Parse(planServicePlanId)
			assignedPlan.SetServicePlanId(&u)
		}
	}
	requestBody.SetAssignedPlans(planAssignedPlans)

	authorizationInfo := models.NewAuthorizationInfo()
	authorizationInfoModel := userAuthorizationInfoModel{}
	plan.AuthorizationInfo.As(ctx, &authorizationInfoModel, basetypes.ObjectAsOptions{})

	if len(authorizationInfoModel.CertificateUserIds.Elements()) > 0 {
		var certificateUserIds []string
		for _, i := range authorizationInfoModel.CertificateUserIds.Elements() {
			certificateUserIds = append(certificateUserIds, i.String())
		}
		authorizationInfo.SetCertificateUserIds(certificateUserIds)
	}
	requestBody.SetAuthorizationInfo(authorizationInfo)

	if len(plan.BusinessPhones.Elements()) > 0 {
		var businessPhones []string
		for _, i := range plan.BusinessPhones.Elements() {
			businessPhones = append(businessPhones, i.String())
		}
		requestBody.SetBusinessPhones(businessPhones)
	}

	if !plan.City.IsUnknown() {
		planCity := plan.City.ValueString()
		requestBody.SetCity(&planCity)
	} else {
	}

	if !plan.CompanyName.IsUnknown() {
		planCompanyName := plan.CompanyName.ValueString()
		requestBody.SetCompanyName(&planCompanyName)
	} else {
	}

	if !plan.ConsentProvidedForMinor.IsUnknown() {
		planConsentProvidedForMinor := plan.ConsentProvidedForMinor.ValueString()
		requestBody.SetConsentProvidedForMinor(&planConsentProvidedForMinor)
	} else {
	}

	if !plan.Country.IsUnknown() {
		planCountry := plan.Country.ValueString()
		requestBody.SetCountry(&planCountry)
	} else {
	}

	if !plan.CreatedDateTime.IsUnknown() {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	}

	if !plan.CreationType.IsUnknown() {
		planCreationType := plan.CreationType.ValueString()
		requestBody.SetCreationType(&planCreationType)
	} else {
	}

	if !plan.Department.IsUnknown() {
		planDepartment := plan.Department.ValueString()
		requestBody.SetDepartment(&planDepartment)
	} else {
	}

	if !plan.DisplayName.IsUnknown() {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	} else {
	}

	if !plan.EmployeeHireDate.IsUnknown() {
		planEmployeeHireDate := plan.EmployeeHireDate.ValueString()
		t, _ = time.Parse(time.RFC3339, planEmployeeHireDate)
		requestBody.SetEmployeeHireDate(&t)
	}

	if !plan.EmployeeId.IsUnknown() {
		planEmployeeId := plan.EmployeeId.ValueString()
		requestBody.SetEmployeeId(&planEmployeeId)
	} else {
	}

	if !plan.EmployeeLeaveDateTime.IsUnknown() {
		planEmployeeLeaveDateTime := plan.EmployeeLeaveDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planEmployeeLeaveDateTime)
		requestBody.SetEmployeeLeaveDateTime(&t)
	}

	employeeOrgData := models.NewEmployeeOrgData()
	employeeOrgDataModel := userEmployeeOrgDataModel{}
	plan.EmployeeOrgData.As(ctx, &employeeOrgDataModel, basetypes.ObjectAsOptions{})

	if !employeeOrgDataModel.CostCenter.IsUnknown() {
		planCostCenter := employeeOrgDataModel.CostCenter.ValueString()
		employeeOrgData.SetCostCenter(&planCostCenter)
	} else {
	}

	if !employeeOrgDataModel.Division.IsUnknown() {
		planDivision := employeeOrgDataModel.Division.ValueString()
		employeeOrgData.SetDivision(&planDivision)
	} else {
	}
	requestBody.SetEmployeeOrgData(employeeOrgData)

	if !plan.EmployeeType.IsUnknown() {
		planEmployeeType := plan.EmployeeType.ValueString()
		requestBody.SetEmployeeType(&planEmployeeType)
	} else {
	}

	if !plan.ExternalUserState.IsUnknown() {
		planExternalUserState := plan.ExternalUserState.ValueString()
		requestBody.SetExternalUserState(&planExternalUserState)
	} else {
	}

	if !plan.ExternalUserStateChangeDateTime.IsUnknown() {
		planExternalUserStateChangeDateTime := plan.ExternalUserStateChangeDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planExternalUserStateChangeDateTime)
		requestBody.SetExternalUserStateChangeDateTime(&t)
	}

	if !plan.FaxNumber.IsUnknown() {
		planFaxNumber := plan.FaxNumber.ValueString()
		requestBody.SetFaxNumber(&planFaxNumber)
	} else {
	}

	if !plan.GivenName.IsUnknown() {
		planGivenName := plan.GivenName.ValueString()
		requestBody.SetGivenName(&planGivenName)
	} else {
	}

	var planIdentities []models.ObjectIdentityable
	for _, i := range plan.Identities.Elements() {
		objectIdentity := models.NewObjectIdentity()
		objectIdentityModel := userIdentitiesModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &objectIdentityModel)

		if !objectIdentityModel.Issuer.IsUnknown() {
			planIssuer := objectIdentityModel.Issuer.ValueString()
			objectIdentity.SetIssuer(&planIssuer)
		} else {
		}

		if !objectIdentityModel.IssuerAssignedId.IsUnknown() {
			planIssuerAssignedId := objectIdentityModel.IssuerAssignedId.ValueString()
			objectIdentity.SetIssuerAssignedId(&planIssuerAssignedId)
		} else {
		}

		if !objectIdentityModel.SignInType.IsUnknown() {
			planSignInType := objectIdentityModel.SignInType.ValueString()
			objectIdentity.SetSignInType(&planSignInType)
		} else {
		}
	}
	requestBody.SetIdentities(planIdentities)

	if len(plan.ImAddresses.Elements()) > 0 {
		var imAddresses []string
		for _, i := range plan.ImAddresses.Elements() {
			imAddresses = append(imAddresses, i.String())
		}
		requestBody.SetImAddresses(imAddresses)
	}

	if len(plan.Interests.Elements()) > 0 {
		var interests []string
		for _, i := range plan.Interests.Elements() {
			interests = append(interests, i.String())
		}
		requestBody.SetInterests(interests)
	}

	if !plan.IsResourceAccount.IsUnknown() {
		planIsResourceAccount := plan.IsResourceAccount.ValueBool()
		requestBody.SetIsResourceAccount(&planIsResourceAccount)
	}

	if !plan.JobTitle.IsUnknown() {
		planJobTitle := plan.JobTitle.ValueString()
		requestBody.SetJobTitle(&planJobTitle)
	} else {
	}

	if !plan.LastPasswordChangeDateTime.IsUnknown() {
		planLastPasswordChangeDateTime := plan.LastPasswordChangeDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planLastPasswordChangeDateTime)
		requestBody.SetLastPasswordChangeDateTime(&t)
	}

	if !plan.LegalAgeGroupClassification.IsUnknown() {
		planLegalAgeGroupClassification := plan.LegalAgeGroupClassification.ValueString()
		requestBody.SetLegalAgeGroupClassification(&planLegalAgeGroupClassification)
	} else {
	}

	var planLicenseAssignmentStates []models.LicenseAssignmentStateable
	for _, i := range plan.LicenseAssignmentStates.Elements() {
		licenseAssignmentState := models.NewLicenseAssignmentState()
		licenseAssignmentStateModel := userLicenseAssignmentStatesModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &licenseAssignmentStateModel)

		if !licenseAssignmentStateModel.AssignedByGroup.IsUnknown() {
			planAssignedByGroup := licenseAssignmentStateModel.AssignedByGroup.ValueString()
			licenseAssignmentState.SetAssignedByGroup(&planAssignedByGroup)
		} else {
		}

		if len(licenseAssignmentStateModel.DisabledPlans.Elements()) > 0 {
			var DisabledPlans []uuid.UUID
			for _, i := range licenseAssignmentStateModel.DisabledPlans.Elements() {
				u, _ = uuid.Parse(i.String())
				DisabledPlans = append(DisabledPlans, u)
			}
			licenseAssignmentState.SetDisabledPlans(DisabledPlans)
		}

		if !licenseAssignmentStateModel.Error.IsUnknown() {
			planError := licenseAssignmentStateModel.Error.ValueString()
			licenseAssignmentState.SetError(&planError)
		} else {
		}

		if !licenseAssignmentStateModel.LastUpdatedDateTime.IsUnknown() {
			planLastUpdatedDateTime := licenseAssignmentStateModel.LastUpdatedDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planLastUpdatedDateTime)
			licenseAssignmentState.SetLastUpdatedDateTime(&t)
		}

		if !licenseAssignmentStateModel.SkuId.IsUnknown() {
			planSkuId := licenseAssignmentStateModel.SkuId.ValueString()
			u, _ = uuid.Parse(planSkuId)
			licenseAssignmentState.SetSkuId(&u)
		}

		if !licenseAssignmentStateModel.State.IsUnknown() {
			planState := licenseAssignmentStateModel.State.ValueString()
			licenseAssignmentState.SetState(&planState)
		} else {
		}
	}
	requestBody.SetLicenseAssignmentStates(planLicenseAssignmentStates)

	if !plan.Mail.IsUnknown() {
		planMail := plan.Mail.ValueString()
		requestBody.SetMail(&planMail)
	} else {
	}

	if !plan.MailNickname.IsUnknown() {
		planMailNickname := plan.MailNickname.ValueString()
		requestBody.SetMailNickname(&planMailNickname)
	} else {
	}

	if !plan.MobilePhone.IsUnknown() {
		planMobilePhone := plan.MobilePhone.ValueString()
		requestBody.SetMobilePhone(&planMobilePhone)
	} else {
	}

	if !plan.OfficeLocation.IsUnknown() {
		planOfficeLocation := plan.OfficeLocation.ValueString()
		requestBody.SetOfficeLocation(&planOfficeLocation)
	} else {
	}

	if !plan.OnPremisesDistinguishedName.IsUnknown() {
		planOnPremisesDistinguishedName := plan.OnPremisesDistinguishedName.ValueString()
		requestBody.SetOnPremisesDistinguishedName(&planOnPremisesDistinguishedName)
	} else {
	}

	if !plan.OnPremisesDomainName.IsUnknown() {
		planOnPremisesDomainName := plan.OnPremisesDomainName.ValueString()
		requestBody.SetOnPremisesDomainName(&planOnPremisesDomainName)
	} else {
	}

	onPremisesExtensionAttributes := models.NewOnPremisesExtensionAttributes()
	onPremisesExtensionAttributesModel := userOnPremisesExtensionAttributesModel{}
	plan.OnPremisesExtensionAttributes.As(ctx, &onPremisesExtensionAttributesModel, basetypes.ObjectAsOptions{})

	if !onPremisesExtensionAttributesModel.ExtensionAttribute1.IsUnknown() {
		planExtensionAttribute1 := onPremisesExtensionAttributesModel.ExtensionAttribute1.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute1(&planExtensionAttribute1)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute10.IsUnknown() {
		planExtensionAttribute10 := onPremisesExtensionAttributesModel.ExtensionAttribute10.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute10(&planExtensionAttribute10)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute11.IsUnknown() {
		planExtensionAttribute11 := onPremisesExtensionAttributesModel.ExtensionAttribute11.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute11(&planExtensionAttribute11)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute12.IsUnknown() {
		planExtensionAttribute12 := onPremisesExtensionAttributesModel.ExtensionAttribute12.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute12(&planExtensionAttribute12)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute13.IsUnknown() {
		planExtensionAttribute13 := onPremisesExtensionAttributesModel.ExtensionAttribute13.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute13(&planExtensionAttribute13)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute14.IsUnknown() {
		planExtensionAttribute14 := onPremisesExtensionAttributesModel.ExtensionAttribute14.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute14(&planExtensionAttribute14)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute15.IsUnknown() {
		planExtensionAttribute15 := onPremisesExtensionAttributesModel.ExtensionAttribute15.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute15(&planExtensionAttribute15)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute2.IsUnknown() {
		planExtensionAttribute2 := onPremisesExtensionAttributesModel.ExtensionAttribute2.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute2(&planExtensionAttribute2)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute3.IsUnknown() {
		planExtensionAttribute3 := onPremisesExtensionAttributesModel.ExtensionAttribute3.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute3(&planExtensionAttribute3)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute4.IsUnknown() {
		planExtensionAttribute4 := onPremisesExtensionAttributesModel.ExtensionAttribute4.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute4(&planExtensionAttribute4)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute5.IsUnknown() {
		planExtensionAttribute5 := onPremisesExtensionAttributesModel.ExtensionAttribute5.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute5(&planExtensionAttribute5)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute6.IsUnknown() {
		planExtensionAttribute6 := onPremisesExtensionAttributesModel.ExtensionAttribute6.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute6(&planExtensionAttribute6)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute7.IsUnknown() {
		planExtensionAttribute7 := onPremisesExtensionAttributesModel.ExtensionAttribute7.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute7(&planExtensionAttribute7)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute8.IsUnknown() {
		planExtensionAttribute8 := onPremisesExtensionAttributesModel.ExtensionAttribute8.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute8(&planExtensionAttribute8)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute9.IsUnknown() {
		planExtensionAttribute9 := onPremisesExtensionAttributesModel.ExtensionAttribute9.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute9(&planExtensionAttribute9)
	} else {
	}
	requestBody.SetOnPremisesExtensionAttributes(onPremisesExtensionAttributes)

	if !plan.OnPremisesImmutableId.IsUnknown() {
		planOnPremisesImmutableId := plan.OnPremisesImmutableId.ValueString()
		requestBody.SetOnPremisesImmutableId(&planOnPremisesImmutableId)
	} else {
	}

	if !plan.OnPremisesLastSyncDateTime.IsUnknown() {
		planOnPremisesLastSyncDateTime := plan.OnPremisesLastSyncDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planOnPremisesLastSyncDateTime)
		requestBody.SetOnPremisesLastSyncDateTime(&t)
	}

	var planOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
	for _, i := range plan.OnPremisesProvisioningErrors.Elements() {
		onPremisesProvisioningError := models.NewOnPremisesProvisioningError()
		onPremisesProvisioningErrorModel := userOnPremisesProvisioningErrorsModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &onPremisesProvisioningErrorModel)

		if !onPremisesProvisioningErrorModel.Category.IsUnknown() {
			planCategory := onPremisesProvisioningErrorModel.Category.ValueString()
			onPremisesProvisioningError.SetCategory(&planCategory)
		} else {
		}

		if !onPremisesProvisioningErrorModel.OccurredDateTime.IsUnknown() {
			planOccurredDateTime := onPremisesProvisioningErrorModel.OccurredDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planOccurredDateTime)
			onPremisesProvisioningError.SetOccurredDateTime(&t)
		}

		if !onPremisesProvisioningErrorModel.PropertyCausingError.IsUnknown() {
			planPropertyCausingError := onPremisesProvisioningErrorModel.PropertyCausingError.ValueString()
			onPremisesProvisioningError.SetPropertyCausingError(&planPropertyCausingError)
		} else {
		}

		if !onPremisesProvisioningErrorModel.Value.IsUnknown() {
			planValue := onPremisesProvisioningErrorModel.Value.ValueString()
			onPremisesProvisioningError.SetValue(&planValue)
		} else {
		}
	}
	requestBody.SetOnPremisesProvisioningErrors(planOnPremisesProvisioningErrors)

	if !plan.OnPremisesSamAccountName.IsUnknown() {
		planOnPremisesSamAccountName := plan.OnPremisesSamAccountName.ValueString()
		requestBody.SetOnPremisesSamAccountName(&planOnPremisesSamAccountName)
	} else {
	}

	if !plan.OnPremisesSecurityIdentifier.IsUnknown() {
		planOnPremisesSecurityIdentifier := plan.OnPremisesSecurityIdentifier.ValueString()
		requestBody.SetOnPremisesSecurityIdentifier(&planOnPremisesSecurityIdentifier)
	} else {
	}

	if !plan.OnPremisesSyncEnabled.IsUnknown() {
		planOnPremisesSyncEnabled := plan.OnPremisesSyncEnabled.ValueBool()
		requestBody.SetOnPremisesSyncEnabled(&planOnPremisesSyncEnabled)
	}

	if !plan.OnPremisesUserPrincipalName.IsUnknown() {
		planOnPremisesUserPrincipalName := plan.OnPremisesUserPrincipalName.ValueString()
		requestBody.SetOnPremisesUserPrincipalName(&planOnPremisesUserPrincipalName)
	} else {
	}

	if len(plan.OtherMails.Elements()) > 0 {
		var otherMails []string
		for _, i := range plan.OtherMails.Elements() {
			otherMails = append(otherMails, i.String())
		}
		requestBody.SetOtherMails(otherMails)
	}

	if !plan.PasswordPolicies.IsUnknown() {
		planPasswordPolicies := plan.PasswordPolicies.ValueString()
		requestBody.SetPasswordPolicies(&planPasswordPolicies)
	} else {
	}

	passwordProfile := models.NewPasswordProfile()
	passwordProfileModel := userPasswordProfileModel{}
	plan.PasswordProfile.As(ctx, &passwordProfileModel, basetypes.ObjectAsOptions{})

	if !passwordProfileModel.ForceChangePasswordNextSignIn.IsUnknown() {
		planForceChangePasswordNextSignIn := passwordProfileModel.ForceChangePasswordNextSignIn.ValueBool()
		passwordProfile.SetForceChangePasswordNextSignIn(&planForceChangePasswordNextSignIn)
	}

	if !passwordProfileModel.ForceChangePasswordNextSignInWithMfa.IsUnknown() {
		planForceChangePasswordNextSignInWithMfa := passwordProfileModel.ForceChangePasswordNextSignInWithMfa.ValueBool()
		passwordProfile.SetForceChangePasswordNextSignInWithMfa(&planForceChangePasswordNextSignInWithMfa)
	}

	if !passwordProfileModel.Password.IsUnknown() {
		planPassword := passwordProfileModel.Password.ValueString()
		passwordProfile.SetPassword(&planPassword)
	} else {
	}
	requestBody.SetPasswordProfile(passwordProfile)

	if len(plan.PastProjects.Elements()) > 0 {
		var pastProjects []string
		for _, i := range plan.PastProjects.Elements() {
			pastProjects = append(pastProjects, i.String())
		}
		requestBody.SetPastProjects(pastProjects)
	}

	if !plan.PostalCode.IsUnknown() {
		planPostalCode := plan.PostalCode.ValueString()
		requestBody.SetPostalCode(&planPostalCode)
	} else {
	}

	if !plan.PreferredDataLocation.IsUnknown() {
		planPreferredDataLocation := plan.PreferredDataLocation.ValueString()
		requestBody.SetPreferredDataLocation(&planPreferredDataLocation)
	} else {
	}

	if !plan.PreferredLanguage.IsUnknown() {
		planPreferredLanguage := plan.PreferredLanguage.ValueString()
		requestBody.SetPreferredLanguage(&planPreferredLanguage)
	} else {
	}

	var planProvisionedPlans []models.ProvisionedPlanable
	for _, i := range plan.ProvisionedPlans.Elements() {
		provisionedPlan := models.NewProvisionedPlan()
		provisionedPlanModel := userProvisionedPlansModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &provisionedPlanModel)

		if !provisionedPlanModel.CapabilityStatus.IsUnknown() {
			planCapabilityStatus := provisionedPlanModel.CapabilityStatus.ValueString()
			provisionedPlan.SetCapabilityStatus(&planCapabilityStatus)
		} else {
		}

		if !provisionedPlanModel.ProvisioningStatus.IsUnknown() {
			planProvisioningStatus := provisionedPlanModel.ProvisioningStatus.ValueString()
			provisionedPlan.SetProvisioningStatus(&planProvisioningStatus)
		} else {
		}

		if !provisionedPlanModel.Service.IsUnknown() {
			planService := provisionedPlanModel.Service.ValueString()
			provisionedPlan.SetService(&planService)
		} else {
		}
	}
	requestBody.SetProvisionedPlans(planProvisionedPlans)

	if len(plan.ProxyAddresses.Elements()) > 0 {
		var proxyAddresses []string
		for _, i := range plan.ProxyAddresses.Elements() {
			proxyAddresses = append(proxyAddresses, i.String())
		}
		requestBody.SetProxyAddresses(proxyAddresses)
	}

	if len(plan.Responsibilities.Elements()) > 0 {
		var responsibilities []string
		for _, i := range plan.Responsibilities.Elements() {
			responsibilities = append(responsibilities, i.String())
		}
		requestBody.SetResponsibilities(responsibilities)
	}

	if len(plan.Schools.Elements()) > 0 {
		var schools []string
		for _, i := range plan.Schools.Elements() {
			schools = append(schools, i.String())
		}
		requestBody.SetSchools(schools)
	}

	if !plan.SecurityIdentifier.IsUnknown() {
		planSecurityIdentifier := plan.SecurityIdentifier.ValueString()
		requestBody.SetSecurityIdentifier(&planSecurityIdentifier)
	} else {
	}

	var planServiceProvisioningErrors []models.ServiceProvisioningErrorable
	for _, i := range plan.ServiceProvisioningErrors.Elements() {
		serviceProvisioningError := models.NewServiceProvisioningError()
		serviceProvisioningErrorModel := userServiceProvisioningErrorsModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &serviceProvisioningErrorModel)

		if !serviceProvisioningErrorModel.CreatedDateTime.IsUnknown() {
			planCreatedDateTime := serviceProvisioningErrorModel.CreatedDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
			serviceProvisioningError.SetCreatedDateTime(&t)
		}

		if !serviceProvisioningErrorModel.IsResolved.IsUnknown() {
			planIsResolved := serviceProvisioningErrorModel.IsResolved.ValueBool()
			serviceProvisioningError.SetIsResolved(&planIsResolved)
		}

		if !serviceProvisioningErrorModel.ServiceInstance.IsUnknown() {
			planServiceInstance := serviceProvisioningErrorModel.ServiceInstance.ValueString()
			serviceProvisioningError.SetServiceInstance(&planServiceInstance)
		} else {
		}
	}
	requestBody.SetServiceProvisioningErrors(planServiceProvisioningErrors)

	if !plan.ShowInAddressList.IsUnknown() {
		planShowInAddressList := plan.ShowInAddressList.ValueBool()
		requestBody.SetShowInAddressList(&planShowInAddressList)
	}

	if !plan.SignInSessionsValidFromDateTime.IsUnknown() {
		planSignInSessionsValidFromDateTime := plan.SignInSessionsValidFromDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planSignInSessionsValidFromDateTime)
		requestBody.SetSignInSessionsValidFromDateTime(&t)
	}

	if len(plan.Skills.Elements()) > 0 {
		var skills []string
		for _, i := range plan.Skills.Elements() {
			skills = append(skills, i.String())
		}
		requestBody.SetSkills(skills)
	}

	if !plan.State.IsUnknown() {
		planState := plan.State.ValueString()
		requestBody.SetState(&planState)
	} else {
	}

	if !plan.StreetAddress.IsUnknown() {
		planStreetAddress := plan.StreetAddress.ValueString()
		requestBody.SetStreetAddress(&planStreetAddress)
	} else {
	}

	if !plan.Surname.IsUnknown() {
		planSurname := plan.Surname.ValueString()
		requestBody.SetSurname(&planSurname)
	} else {
	}

	if !plan.UsageLocation.IsUnknown() {
		planUsageLocation := plan.UsageLocation.ValueString()
		requestBody.SetUsageLocation(&planUsageLocation)
	} else {
	}

	if !plan.UserPrincipalName.IsUnknown() {
		planUserPrincipalName := plan.UserPrincipalName.ValueString()
		requestBody.SetUserPrincipalName(&planUserPrincipalName)
	} else {
	}

	if !plan.UserType.IsUnknown() {
		planUserType := plan.UserType.ValueString()
		requestBody.SetUserType(&planUserType)
	} else {
	}

	// Create new user
	result, err := r.client.Users().Post(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	plan.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state userModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
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
				"aboutMe",
				"birthday",
				"hireDate",
				"interests",
				"mySite",
				"pastProjects",
				"preferredName",
				"responsibilities",
				"schools",
				"skills",
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
				"permissionGrants",
				"teamwork",
				"todo",
				"employeeExperience",
			},
		},
	}

	var result models.Userable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.Users().ByUserId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else if !state.UserPrincipalName.IsNull() {
		result, err = d.client.Users().ByUserId(state.UserPrincipalName.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"`id` or `user_principal_name` must be supplied.",
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

	if result.GetId() != nil {
		state.Id = types.StringValue(*result.GetId())
	}
	if result.GetDeletedDateTime() != nil {
		state.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
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
	if len(result.GetAssignedLicenses()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedLicenses() {
			assignedLicenses := new(userAssignedLicensesModel)

			if len(v.GetDisabledPlans()) > 0 {
				var disabledPlans []attr.Value
				for _, v := range v.GetDisabledPlans() {
					disabledPlans = append(disabledPlans, types.StringValue(v.String()))
				}
				listValue, _ := types.ListValue(types.StringType, disabledPlans)
				assignedLicenses.DisabledPlans = listValue
			} else {
				assignedLicenses.DisabledPlans = types.ListNull(types.StringType)
			}
			if v.GetSkuId() != nil {
				assignedLicenses.SkuId = types.StringValue(v.GetSkuId().String())
			}
			objectValue, _ := types.ObjectValueFrom(ctx, assignedLicenses.AttributeTypes(), assignedLicenses)
			objectValues = append(objectValues, objectValue)
		}
		state.AssignedLicenses, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(result.GetAssignedPlans()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedPlans() {
			assignedPlans := new(userAssignedPlansModel)

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
			objectValue, _ := types.ObjectValueFrom(ctx, assignedPlans.AttributeTypes(), assignedPlans)
			objectValues = append(objectValues, objectValue)
		}
		state.AssignedPlans, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetAuthorizationInfo() != nil {
		authorizationInfo := new(userAuthorizationInfoModel)

		if len(result.GetAuthorizationInfo().GetCertificateUserIds()) > 0 {
			var certificateUserIds []attr.Value
			for _, v := range result.GetAuthorizationInfo().GetCertificateUserIds() {
				certificateUserIds = append(certificateUserIds, types.StringValue(v))
			}
			listValue, _ := types.ListValue(types.StringType, certificateUserIds)
			authorizationInfo.CertificateUserIds = listValue
		} else {
			authorizationInfo.CertificateUserIds = types.ListNull(types.StringType)
		}

		objectValue, _ := types.ObjectValueFrom(ctx, authorizationInfo.AttributeTypes(), authorizationInfo)
		state.AuthorizationInfo = objectValue
	}
	if result.GetBirthday() != nil {
		state.Birthday = types.StringValue(result.GetBirthday().String())
	}
	if len(result.GetBusinessPhones()) > 0 {
		var businessPhones []attr.Value
		for _, v := range result.GetBusinessPhones() {
			businessPhones = append(businessPhones, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, businessPhones)
		state.BusinessPhones = listValue
	} else {
		state.BusinessPhones = types.ListNull(types.StringType)
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
	if result.GetDepartment() != nil {
		state.Department = types.StringValue(*result.GetDepartment())
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	}
	if result.GetEmployeeHireDate() != nil {
		state.EmployeeHireDate = types.StringValue(result.GetEmployeeHireDate().String())
	}
	if result.GetEmployeeId() != nil {
		state.EmployeeId = types.StringValue(*result.GetEmployeeId())
	}
	if result.GetEmployeeLeaveDateTime() != nil {
		state.EmployeeLeaveDateTime = types.StringValue(result.GetEmployeeLeaveDateTime().String())
	}
	if result.GetEmployeeOrgData() != nil {
		employeeOrgData := new(userEmployeeOrgDataModel)

		if result.GetEmployeeOrgData().GetCostCenter() != nil {
			employeeOrgData.CostCenter = types.StringValue(*result.GetEmployeeOrgData().GetCostCenter())
		}
		if result.GetEmployeeOrgData().GetDivision() != nil {
			employeeOrgData.Division = types.StringValue(*result.GetEmployeeOrgData().GetDivision())
		}

		objectValue, _ := types.ObjectValueFrom(ctx, employeeOrgData.AttributeTypes(), employeeOrgData)
		state.EmployeeOrgData = objectValue
	}
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
	if len(result.GetIdentities()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetIdentities() {
			identities := new(userIdentitiesModel)

			if v.GetIssuer() != nil {
				identities.Issuer = types.StringValue(*v.GetIssuer())
			}
			if v.GetIssuerAssignedId() != nil {
				identities.IssuerAssignedId = types.StringValue(*v.GetIssuerAssignedId())
			}
			if v.GetSignInType() != nil {
				identities.SignInType = types.StringValue(*v.GetSignInType())
			}
			objectValue, _ := types.ObjectValueFrom(ctx, identities.AttributeTypes(), identities)
			objectValues = append(objectValues, objectValue)
		}
		state.Identities, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(result.GetImAddresses()) > 0 {
		var imAddresses []attr.Value
		for _, v := range result.GetImAddresses() {
			imAddresses = append(imAddresses, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, imAddresses)
		state.ImAddresses = listValue
	} else {
		state.ImAddresses = types.ListNull(types.StringType)
	}
	if len(result.GetInterests()) > 0 {
		var interests []attr.Value
		for _, v := range result.GetInterests() {
			interests = append(interests, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, interests)
		state.Interests = listValue
	} else {
		state.Interests = types.ListNull(types.StringType)
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
	if len(result.GetLicenseAssignmentStates()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetLicenseAssignmentStates() {
			licenseAssignmentStates := new(userLicenseAssignmentStatesModel)

			if v.GetAssignedByGroup() != nil {
				licenseAssignmentStates.AssignedByGroup = types.StringValue(*v.GetAssignedByGroup())
			}
			if len(v.GetDisabledPlans()) > 0 {
				var disabledPlans []attr.Value
				for _, v := range v.GetDisabledPlans() {
					disabledPlans = append(disabledPlans, types.StringValue(v.String()))
				}
				listValue, _ := types.ListValue(types.StringType, disabledPlans)
				licenseAssignmentStates.DisabledPlans = listValue
			} else {
				licenseAssignmentStates.DisabledPlans = types.ListNull(types.StringType)
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
			objectValue, _ := types.ObjectValueFrom(ctx, licenseAssignmentStates.AttributeTypes(), licenseAssignmentStates)
			objectValues = append(objectValues, objectValue)
		}
		state.LicenseAssignmentStates, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
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
	if result.GetOnPremisesExtensionAttributes() != nil {
		onPremisesExtensionAttributes := new(userOnPremisesExtensionAttributesModel)

		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute1 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute10 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute11 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute12 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute13 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute14 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute15 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute2 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute3 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute4 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute5 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute6 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute7 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute8 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8())
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute9 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9())
		}

		objectValue, _ := types.ObjectValueFrom(ctx, onPremisesExtensionAttributes.AttributeTypes(), onPremisesExtensionAttributes)
		state.OnPremisesExtensionAttributes = objectValue
	}
	if result.GetOnPremisesImmutableId() != nil {
		state.OnPremisesImmutableId = types.StringValue(*result.GetOnPremisesImmutableId())
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	}
	if len(result.GetOnPremisesProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetOnPremisesProvisioningErrors() {
			onPremisesProvisioningErrors := new(userOnPremisesProvisioningErrorsModel)

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
			objectValue, _ := types.ObjectValueFrom(ctx, onPremisesProvisioningErrors.AttributeTypes(), onPremisesProvisioningErrors)
			objectValues = append(objectValues, objectValue)
		}
		state.OnPremisesProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
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
	if len(result.GetOtherMails()) > 0 {
		var otherMails []attr.Value
		for _, v := range result.GetOtherMails() {
			otherMails = append(otherMails, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, otherMails)
		state.OtherMails = listValue
	} else {
		state.OtherMails = types.ListNull(types.StringType)
	}
	if result.GetPasswordPolicies() != nil {
		state.PasswordPolicies = types.StringValue(*result.GetPasswordPolicies())
	}
	if result.GetPasswordProfile() != nil {
		passwordProfile := new(userPasswordProfileModel)

		if result.GetPasswordProfile().GetForceChangePasswordNextSignIn() != nil {
			passwordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
		}
		if result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa() != nil {
			passwordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
		}
		if result.GetPasswordProfile().GetPassword() != nil {
			passwordProfile.Password = types.StringValue(*result.GetPasswordProfile().GetPassword())
		}

		objectValue, _ := types.ObjectValueFrom(ctx, passwordProfile.AttributeTypes(), passwordProfile)
		state.PasswordProfile = objectValue
	}
	if len(result.GetPastProjects()) > 0 {
		var pastProjects []attr.Value
		for _, v := range result.GetPastProjects() {
			pastProjects = append(pastProjects, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, pastProjects)
		state.PastProjects = listValue
	} else {
		state.PastProjects = types.ListNull(types.StringType)
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
	if len(result.GetProvisionedPlans()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetProvisionedPlans() {
			provisionedPlans := new(userProvisionedPlansModel)

			if v.GetCapabilityStatus() != nil {
				provisionedPlans.CapabilityStatus = types.StringValue(*v.GetCapabilityStatus())
			}
			if v.GetProvisioningStatus() != nil {
				provisionedPlans.ProvisioningStatus = types.StringValue(*v.GetProvisioningStatus())
			}
			if v.GetService() != nil {
				provisionedPlans.Service = types.StringValue(*v.GetService())
			}
			objectValue, _ := types.ObjectValueFrom(ctx, provisionedPlans.AttributeTypes(), provisionedPlans)
			objectValues = append(objectValues, objectValue)
		}
		state.ProvisionedPlans, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(result.GetProxyAddresses()) > 0 {
		var proxyAddresses []attr.Value
		for _, v := range result.GetProxyAddresses() {
			proxyAddresses = append(proxyAddresses, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, proxyAddresses)
		state.ProxyAddresses = listValue
	} else {
		state.ProxyAddresses = types.ListNull(types.StringType)
	}
	if len(result.GetResponsibilities()) > 0 {
		var responsibilities []attr.Value
		for _, v := range result.GetResponsibilities() {
			responsibilities = append(responsibilities, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, responsibilities)
		state.Responsibilities = listValue
	} else {
		state.Responsibilities = types.ListNull(types.StringType)
	}
	if len(result.GetSchools()) > 0 {
		var schools []attr.Value
		for _, v := range result.GetSchools() {
			schools = append(schools, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, schools)
		state.Schools = listValue
	} else {
		state.Schools = types.ListNull(types.StringType)
	}
	if result.GetSecurityIdentifier() != nil {
		state.SecurityIdentifier = types.StringValue(*result.GetSecurityIdentifier())
	}
	if len(result.GetServiceProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetServiceProvisioningErrors() {
			serviceProvisioningErrors := new(userServiceProvisioningErrorsModel)

			if v.GetCreatedDateTime() != nil {
				serviceProvisioningErrors.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			}
			if v.GetIsResolved() != nil {
				serviceProvisioningErrors.IsResolved = types.BoolValue(*v.GetIsResolved())
			}
			if v.GetServiceInstance() != nil {
				serviceProvisioningErrors.ServiceInstance = types.StringValue(*v.GetServiceInstance())
			}
			objectValue, _ := types.ObjectValueFrom(ctx, serviceProvisioningErrors.AttributeTypes(), serviceProvisioningErrors)
			objectValues = append(objectValues, objectValue)
		}
		state.ServiceProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetShowInAddressList() != nil {
		state.ShowInAddressList = types.BoolValue(*result.GetShowInAddressList())
	}
	if result.GetSignInActivity() != nil {
		signInActivity := new(userSignInActivityModel)

		if result.GetSignInActivity().GetLastNonInteractiveSignInDateTime() != nil {
			signInActivity.LastNonInteractiveSignInDateTime = types.StringValue(result.GetSignInActivity().GetLastNonInteractiveSignInDateTime().String())
		}
		if result.GetSignInActivity().GetLastNonInteractiveSignInRequestId() != nil {
			signInActivity.LastNonInteractiveSignInRequestId = types.StringValue(*result.GetSignInActivity().GetLastNonInteractiveSignInRequestId())
		}
		if result.GetSignInActivity().GetLastSignInDateTime() != nil {
			signInActivity.LastSignInDateTime = types.StringValue(result.GetSignInActivity().GetLastSignInDateTime().String())
		}
		if result.GetSignInActivity().GetLastSignInRequestId() != nil {
			signInActivity.LastSignInRequestId = types.StringValue(*result.GetSignInActivity().GetLastSignInRequestId())
		}

		objectValue, _ := types.ObjectValueFrom(ctx, signInActivity.AttributeTypes(), signInActivity)
		state.SignInActivity = objectValue
	}
	if result.GetSignInSessionsValidFromDateTime() != nil {
		state.SignInSessionsValidFromDateTime = types.StringValue(result.GetSignInSessionsValidFromDateTime().String())
	}
	if len(result.GetSkills()) > 0 {
		var skills []attr.Value
		for _, v := range result.GetSkills() {
			skills = append(skills, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, skills)
		state.Skills = listValue
	} else {
		state.Skills = types.ListNull(types.StringType)
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

// Update updates the resource and sets the updated Terraform state on success.
func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan userModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state userModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewUser()
	var t time.Time
	var u uuid.UUID

	if !plan.Id.IsUnknown() {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	} else {
	}

	if !plan.AccountEnabled.IsUnknown() {
		planAccountEnabled := plan.AccountEnabled.ValueBool()
		requestBody.SetAccountEnabled(&planAccountEnabled)
	}

	if !plan.AgeGroup.IsUnknown() {
		planAgeGroup := plan.AgeGroup.ValueString()
		requestBody.SetAgeGroup(&planAgeGroup)
	} else {
	}

	var planAssignedLicenses []models.AssignedLicenseable
	for _, i := range plan.AssignedLicenses.Elements() {
		assignedLicense := models.NewAssignedLicense()
		assignedLicenseModel := userAssignedLicensesModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &assignedLicenseModel)

		if len(assignedLicenseModel.DisabledPlans.Elements()) > 0 {
			var DisabledPlans []uuid.UUID
			for _, i := range assignedLicenseModel.DisabledPlans.Elements() {
				u, _ = uuid.Parse(i.String())
				DisabledPlans = append(DisabledPlans, u)
			}
			assignedLicense.SetDisabledPlans(DisabledPlans)
		}

		if !assignedLicenseModel.SkuId.IsUnknown() {
			planSkuId := assignedLicenseModel.SkuId.ValueString()
			u, _ = uuid.Parse(planSkuId)
			assignedLicense.SetSkuId(&u)
		}
	}
	requestBody.SetAssignedLicenses(planAssignedLicenses)

	var planAssignedPlans []models.AssignedPlanable
	for _, i := range plan.AssignedPlans.Elements() {
		assignedPlan := models.NewAssignedPlan()
		assignedPlanModel := userAssignedPlansModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &assignedPlanModel)

		if !assignedPlanModel.AssignedDateTime.IsUnknown() {
			planAssignedDateTime := assignedPlanModel.AssignedDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planAssignedDateTime)
			assignedPlan.SetAssignedDateTime(&t)
		}

		if !assignedPlanModel.CapabilityStatus.IsUnknown() {
			planCapabilityStatus := assignedPlanModel.CapabilityStatus.ValueString()
			assignedPlan.SetCapabilityStatus(&planCapabilityStatus)
		} else {
		}

		if !assignedPlanModel.Service.IsUnknown() {
			planService := assignedPlanModel.Service.ValueString()
			assignedPlan.SetService(&planService)
		} else {
		}

		if !assignedPlanModel.ServicePlanId.IsUnknown() {
			planServicePlanId := assignedPlanModel.ServicePlanId.ValueString()
			u, _ = uuid.Parse(planServicePlanId)
			assignedPlan.SetServicePlanId(&u)
		}
	}
	requestBody.SetAssignedPlans(planAssignedPlans)

	authorizationInfo := models.NewAuthorizationInfo()
	authorizationInfoModel := userAuthorizationInfoModel{}
	plan.AuthorizationInfo.As(ctx, &authorizationInfoModel, basetypes.ObjectAsOptions{})

	if len(authorizationInfoModel.CertificateUserIds.Elements()) > 0 {
		var certificateUserIds []string
		for _, i := range authorizationInfoModel.CertificateUserIds.Elements() {
			certificateUserIds = append(certificateUserIds, i.String())
		}
		authorizationInfo.SetCertificateUserIds(certificateUserIds)
	}
	requestBody.SetAuthorizationInfo(authorizationInfo)

	if len(plan.BusinessPhones.Elements()) > 0 {
		var businessPhones []string
		for _, i := range plan.BusinessPhones.Elements() {
			businessPhones = append(businessPhones, i.String())
		}
		requestBody.SetBusinessPhones(businessPhones)
	}

	if !plan.City.IsUnknown() {
		planCity := plan.City.ValueString()
		requestBody.SetCity(&planCity)
	} else {
	}

	if !plan.CompanyName.IsUnknown() {
		planCompanyName := plan.CompanyName.ValueString()
		requestBody.SetCompanyName(&planCompanyName)
	} else {
	}

	if !plan.ConsentProvidedForMinor.IsUnknown() {
		planConsentProvidedForMinor := plan.ConsentProvidedForMinor.ValueString()
		requestBody.SetConsentProvidedForMinor(&planConsentProvidedForMinor)
	} else {
	}

	if !plan.Country.IsUnknown() {
		planCountry := plan.Country.ValueString()
		requestBody.SetCountry(&planCountry)
	} else {
	}

	if !plan.CreatedDateTime.IsUnknown() {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	}

	if !plan.CreationType.IsUnknown() {
		planCreationType := plan.CreationType.ValueString()
		requestBody.SetCreationType(&planCreationType)
	} else {
	}

	if !plan.Department.IsUnknown() {
		planDepartment := plan.Department.ValueString()
		requestBody.SetDepartment(&planDepartment)
	} else {
	}

	if !plan.DisplayName.IsUnknown() {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	} else {
	}

	if !plan.EmployeeHireDate.IsUnknown() {
		planEmployeeHireDate := plan.EmployeeHireDate.ValueString()
		t, _ = time.Parse(time.RFC3339, planEmployeeHireDate)
		requestBody.SetEmployeeHireDate(&t)
	}

	if !plan.EmployeeId.IsUnknown() {
		planEmployeeId := plan.EmployeeId.ValueString()
		requestBody.SetEmployeeId(&planEmployeeId)
	} else {
	}

	if !plan.EmployeeLeaveDateTime.IsUnknown() {
		planEmployeeLeaveDateTime := plan.EmployeeLeaveDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planEmployeeLeaveDateTime)
		requestBody.SetEmployeeLeaveDateTime(&t)
	}

	employeeOrgData := models.NewEmployeeOrgData()
	employeeOrgDataModel := userEmployeeOrgDataModel{}
	plan.EmployeeOrgData.As(ctx, &employeeOrgDataModel, basetypes.ObjectAsOptions{})

	if !employeeOrgDataModel.CostCenter.IsUnknown() {
		planCostCenter := employeeOrgDataModel.CostCenter.ValueString()
		employeeOrgData.SetCostCenter(&planCostCenter)
	} else {
	}

	if !employeeOrgDataModel.Division.IsUnknown() {
		planDivision := employeeOrgDataModel.Division.ValueString()
		employeeOrgData.SetDivision(&planDivision)
	} else {
	}
	requestBody.SetEmployeeOrgData(employeeOrgData)

	if !plan.EmployeeType.IsUnknown() {
		planEmployeeType := plan.EmployeeType.ValueString()
		requestBody.SetEmployeeType(&planEmployeeType)
	} else {
	}

	if !plan.ExternalUserState.IsUnknown() {
		planExternalUserState := plan.ExternalUserState.ValueString()
		requestBody.SetExternalUserState(&planExternalUserState)
	} else {
	}

	if !plan.ExternalUserStateChangeDateTime.IsUnknown() {
		planExternalUserStateChangeDateTime := plan.ExternalUserStateChangeDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planExternalUserStateChangeDateTime)
		requestBody.SetExternalUserStateChangeDateTime(&t)
	}

	if !plan.FaxNumber.IsUnknown() {
		planFaxNumber := plan.FaxNumber.ValueString()
		requestBody.SetFaxNumber(&planFaxNumber)
	} else {
	}

	if !plan.GivenName.IsUnknown() {
		planGivenName := plan.GivenName.ValueString()
		requestBody.SetGivenName(&planGivenName)
	} else {
	}

	var planIdentities []models.ObjectIdentityable
	for _, i := range plan.Identities.Elements() {
		objectIdentity := models.NewObjectIdentity()
		objectIdentityModel := userIdentitiesModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &objectIdentityModel)

		if !objectIdentityModel.Issuer.IsUnknown() {
			planIssuer := objectIdentityModel.Issuer.ValueString()
			objectIdentity.SetIssuer(&planIssuer)
		} else {
		}

		if !objectIdentityModel.IssuerAssignedId.IsUnknown() {
			planIssuerAssignedId := objectIdentityModel.IssuerAssignedId.ValueString()
			objectIdentity.SetIssuerAssignedId(&planIssuerAssignedId)
		} else {
		}

		if !objectIdentityModel.SignInType.IsUnknown() {
			planSignInType := objectIdentityModel.SignInType.ValueString()
			objectIdentity.SetSignInType(&planSignInType)
		} else {
		}
	}
	requestBody.SetIdentities(planIdentities)

	if len(plan.ImAddresses.Elements()) > 0 {
		var imAddresses []string
		for _, i := range plan.ImAddresses.Elements() {
			imAddresses = append(imAddresses, i.String())
		}
		requestBody.SetImAddresses(imAddresses)
	}

	if len(plan.Interests.Elements()) > 0 {
		var interests []string
		for _, i := range plan.Interests.Elements() {
			interests = append(interests, i.String())
		}
		requestBody.SetInterests(interests)
	}

	if !plan.IsResourceAccount.IsUnknown() {
		planIsResourceAccount := plan.IsResourceAccount.ValueBool()
		requestBody.SetIsResourceAccount(&planIsResourceAccount)
	}

	if !plan.JobTitle.IsUnknown() {
		planJobTitle := plan.JobTitle.ValueString()
		requestBody.SetJobTitle(&planJobTitle)
	} else {
	}

	if !plan.LastPasswordChangeDateTime.IsUnknown() {
		planLastPasswordChangeDateTime := plan.LastPasswordChangeDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planLastPasswordChangeDateTime)
		requestBody.SetLastPasswordChangeDateTime(&t)
	}

	if !plan.LegalAgeGroupClassification.IsUnknown() {
		planLegalAgeGroupClassification := plan.LegalAgeGroupClassification.ValueString()
		requestBody.SetLegalAgeGroupClassification(&planLegalAgeGroupClassification)
	} else {
	}

	var planLicenseAssignmentStates []models.LicenseAssignmentStateable
	for _, i := range plan.LicenseAssignmentStates.Elements() {
		licenseAssignmentState := models.NewLicenseAssignmentState()
		licenseAssignmentStateModel := userLicenseAssignmentStatesModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &licenseAssignmentStateModel)

		if !licenseAssignmentStateModel.AssignedByGroup.IsUnknown() {
			planAssignedByGroup := licenseAssignmentStateModel.AssignedByGroup.ValueString()
			licenseAssignmentState.SetAssignedByGroup(&planAssignedByGroup)
		} else {
		}

		if len(licenseAssignmentStateModel.DisabledPlans.Elements()) > 0 {
			var DisabledPlans []uuid.UUID
			for _, i := range licenseAssignmentStateModel.DisabledPlans.Elements() {
				u, _ = uuid.Parse(i.String())
				DisabledPlans = append(DisabledPlans, u)
			}
			licenseAssignmentState.SetDisabledPlans(DisabledPlans)
		}

		if !licenseAssignmentStateModel.Error.IsUnknown() {
			planError := licenseAssignmentStateModel.Error.ValueString()
			licenseAssignmentState.SetError(&planError)
		} else {
		}

		if !licenseAssignmentStateModel.LastUpdatedDateTime.IsUnknown() {
			planLastUpdatedDateTime := licenseAssignmentStateModel.LastUpdatedDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planLastUpdatedDateTime)
			licenseAssignmentState.SetLastUpdatedDateTime(&t)
		}

		if !licenseAssignmentStateModel.SkuId.IsUnknown() {
			planSkuId := licenseAssignmentStateModel.SkuId.ValueString()
			u, _ = uuid.Parse(planSkuId)
			licenseAssignmentState.SetSkuId(&u)
		}

		if !licenseAssignmentStateModel.State.IsUnknown() {
			planState := licenseAssignmentStateModel.State.ValueString()
			licenseAssignmentState.SetState(&planState)
		} else {
		}
	}
	requestBody.SetLicenseAssignmentStates(planLicenseAssignmentStates)

	if !plan.Mail.IsUnknown() {
		planMail := plan.Mail.ValueString()
		requestBody.SetMail(&planMail)
	} else {
	}

	if !plan.MailNickname.IsUnknown() {
		planMailNickname := plan.MailNickname.ValueString()
		requestBody.SetMailNickname(&planMailNickname)
	} else {
	}

	if !plan.MobilePhone.IsUnknown() {
		planMobilePhone := plan.MobilePhone.ValueString()
		requestBody.SetMobilePhone(&planMobilePhone)
	} else {
	}

	if !plan.OfficeLocation.IsUnknown() {
		planOfficeLocation := plan.OfficeLocation.ValueString()
		requestBody.SetOfficeLocation(&planOfficeLocation)
	} else {
	}

	if !plan.OnPremisesDistinguishedName.IsUnknown() {
		planOnPremisesDistinguishedName := plan.OnPremisesDistinguishedName.ValueString()
		requestBody.SetOnPremisesDistinguishedName(&planOnPremisesDistinguishedName)
	} else {
	}

	if !plan.OnPremisesDomainName.IsUnknown() {
		planOnPremisesDomainName := plan.OnPremisesDomainName.ValueString()
		requestBody.SetOnPremisesDomainName(&planOnPremisesDomainName)
	} else {
	}

	onPremisesExtensionAttributes := models.NewOnPremisesExtensionAttributes()
	onPremisesExtensionAttributesModel := userOnPremisesExtensionAttributesModel{}
	plan.OnPremisesExtensionAttributes.As(ctx, &onPremisesExtensionAttributesModel, basetypes.ObjectAsOptions{})

	if !onPremisesExtensionAttributesModel.ExtensionAttribute1.IsUnknown() {
		planExtensionAttribute1 := onPremisesExtensionAttributesModel.ExtensionAttribute1.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute1(&planExtensionAttribute1)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute10.IsUnknown() {
		planExtensionAttribute10 := onPremisesExtensionAttributesModel.ExtensionAttribute10.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute10(&planExtensionAttribute10)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute11.IsUnknown() {
		planExtensionAttribute11 := onPremisesExtensionAttributesModel.ExtensionAttribute11.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute11(&planExtensionAttribute11)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute12.IsUnknown() {
		planExtensionAttribute12 := onPremisesExtensionAttributesModel.ExtensionAttribute12.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute12(&planExtensionAttribute12)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute13.IsUnknown() {
		planExtensionAttribute13 := onPremisesExtensionAttributesModel.ExtensionAttribute13.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute13(&planExtensionAttribute13)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute14.IsUnknown() {
		planExtensionAttribute14 := onPremisesExtensionAttributesModel.ExtensionAttribute14.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute14(&planExtensionAttribute14)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute15.IsUnknown() {
		planExtensionAttribute15 := onPremisesExtensionAttributesModel.ExtensionAttribute15.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute15(&planExtensionAttribute15)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute2.IsUnknown() {
		planExtensionAttribute2 := onPremisesExtensionAttributesModel.ExtensionAttribute2.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute2(&planExtensionAttribute2)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute3.IsUnknown() {
		planExtensionAttribute3 := onPremisesExtensionAttributesModel.ExtensionAttribute3.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute3(&planExtensionAttribute3)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute4.IsUnknown() {
		planExtensionAttribute4 := onPremisesExtensionAttributesModel.ExtensionAttribute4.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute4(&planExtensionAttribute4)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute5.IsUnknown() {
		planExtensionAttribute5 := onPremisesExtensionAttributesModel.ExtensionAttribute5.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute5(&planExtensionAttribute5)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute6.IsUnknown() {
		planExtensionAttribute6 := onPremisesExtensionAttributesModel.ExtensionAttribute6.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute6(&planExtensionAttribute6)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute7.IsUnknown() {
		planExtensionAttribute7 := onPremisesExtensionAttributesModel.ExtensionAttribute7.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute7(&planExtensionAttribute7)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute8.IsUnknown() {
		planExtensionAttribute8 := onPremisesExtensionAttributesModel.ExtensionAttribute8.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute8(&planExtensionAttribute8)
	} else {
	}

	if !onPremisesExtensionAttributesModel.ExtensionAttribute9.IsUnknown() {
		planExtensionAttribute9 := onPremisesExtensionAttributesModel.ExtensionAttribute9.ValueString()
		onPremisesExtensionAttributes.SetExtensionAttribute9(&planExtensionAttribute9)
	} else {
	}
	requestBody.SetOnPremisesExtensionAttributes(onPremisesExtensionAttributes)

	if !plan.OnPremisesImmutableId.IsUnknown() {
		planOnPremisesImmutableId := plan.OnPremisesImmutableId.ValueString()
		requestBody.SetOnPremisesImmutableId(&planOnPremisesImmutableId)
	} else {
	}

	if !plan.OnPremisesLastSyncDateTime.IsUnknown() {
		planOnPremisesLastSyncDateTime := plan.OnPremisesLastSyncDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planOnPremisesLastSyncDateTime)
		requestBody.SetOnPremisesLastSyncDateTime(&t)
	}

	var planOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
	for _, i := range plan.OnPremisesProvisioningErrors.Elements() {
		onPremisesProvisioningError := models.NewOnPremisesProvisioningError()
		onPremisesProvisioningErrorModel := userOnPremisesProvisioningErrorsModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &onPremisesProvisioningErrorModel)

		if !onPremisesProvisioningErrorModel.Category.IsUnknown() {
			planCategory := onPremisesProvisioningErrorModel.Category.ValueString()
			onPremisesProvisioningError.SetCategory(&planCategory)
		} else {
		}

		if !onPremisesProvisioningErrorModel.OccurredDateTime.IsUnknown() {
			planOccurredDateTime := onPremisesProvisioningErrorModel.OccurredDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planOccurredDateTime)
			onPremisesProvisioningError.SetOccurredDateTime(&t)
		}

		if !onPremisesProvisioningErrorModel.PropertyCausingError.IsUnknown() {
			planPropertyCausingError := onPremisesProvisioningErrorModel.PropertyCausingError.ValueString()
			onPremisesProvisioningError.SetPropertyCausingError(&planPropertyCausingError)
		} else {
		}

		if !onPremisesProvisioningErrorModel.Value.IsUnknown() {
			planValue := onPremisesProvisioningErrorModel.Value.ValueString()
			onPremisesProvisioningError.SetValue(&planValue)
		} else {
		}
	}
	requestBody.SetOnPremisesProvisioningErrors(planOnPremisesProvisioningErrors)

	if !plan.OnPremisesSamAccountName.IsUnknown() {
		planOnPremisesSamAccountName := plan.OnPremisesSamAccountName.ValueString()
		requestBody.SetOnPremisesSamAccountName(&planOnPremisesSamAccountName)
	} else {
	}

	if !plan.OnPremisesSecurityIdentifier.IsUnknown() {
		planOnPremisesSecurityIdentifier := plan.OnPremisesSecurityIdentifier.ValueString()
		requestBody.SetOnPremisesSecurityIdentifier(&planOnPremisesSecurityIdentifier)
	} else {
	}

	if !plan.OnPremisesSyncEnabled.IsUnknown() {
		planOnPremisesSyncEnabled := plan.OnPremisesSyncEnabled.ValueBool()
		requestBody.SetOnPremisesSyncEnabled(&planOnPremisesSyncEnabled)
	}

	if !plan.OnPremisesUserPrincipalName.IsUnknown() {
		planOnPremisesUserPrincipalName := plan.OnPremisesUserPrincipalName.ValueString()
		requestBody.SetOnPremisesUserPrincipalName(&planOnPremisesUserPrincipalName)
	} else {
	}

	if len(plan.OtherMails.Elements()) > 0 {
		var otherMails []string
		for _, i := range plan.OtherMails.Elements() {
			otherMails = append(otherMails, i.String())
		}
		requestBody.SetOtherMails(otherMails)
	}

	if !plan.PasswordPolicies.IsUnknown() {
		planPasswordPolicies := plan.PasswordPolicies.ValueString()
		requestBody.SetPasswordPolicies(&planPasswordPolicies)
	} else {
	}

	passwordProfile := models.NewPasswordProfile()
	passwordProfileModel := userPasswordProfileModel{}
	plan.PasswordProfile.As(ctx, &passwordProfileModel, basetypes.ObjectAsOptions{})

	if !passwordProfileModel.ForceChangePasswordNextSignIn.IsUnknown() {
		planForceChangePasswordNextSignIn := passwordProfileModel.ForceChangePasswordNextSignIn.ValueBool()
		passwordProfile.SetForceChangePasswordNextSignIn(&planForceChangePasswordNextSignIn)
	}

	if !passwordProfileModel.ForceChangePasswordNextSignInWithMfa.IsUnknown() {
		planForceChangePasswordNextSignInWithMfa := passwordProfileModel.ForceChangePasswordNextSignInWithMfa.ValueBool()
		passwordProfile.SetForceChangePasswordNextSignInWithMfa(&planForceChangePasswordNextSignInWithMfa)
	}

	if !passwordProfileModel.Password.IsUnknown() {
		planPassword := passwordProfileModel.Password.ValueString()
		passwordProfile.SetPassword(&planPassword)
	} else {
	}
	requestBody.SetPasswordProfile(passwordProfile)

	if len(plan.PastProjects.Elements()) > 0 {
		var pastProjects []string
		for _, i := range plan.PastProjects.Elements() {
			pastProjects = append(pastProjects, i.String())
		}
		requestBody.SetPastProjects(pastProjects)
	}

	if !plan.PostalCode.IsUnknown() {
		planPostalCode := plan.PostalCode.ValueString()
		requestBody.SetPostalCode(&planPostalCode)
	} else {
	}

	if !plan.PreferredDataLocation.IsUnknown() {
		planPreferredDataLocation := plan.PreferredDataLocation.ValueString()
		requestBody.SetPreferredDataLocation(&planPreferredDataLocation)
	} else {
	}

	if !plan.PreferredLanguage.IsUnknown() {
		planPreferredLanguage := plan.PreferredLanguage.ValueString()
		requestBody.SetPreferredLanguage(&planPreferredLanguage)
	} else {
	}

	var planProvisionedPlans []models.ProvisionedPlanable
	for _, i := range plan.ProvisionedPlans.Elements() {
		provisionedPlan := models.NewProvisionedPlan()
		provisionedPlanModel := userProvisionedPlansModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &provisionedPlanModel)

		if !provisionedPlanModel.CapabilityStatus.IsUnknown() {
			planCapabilityStatus := provisionedPlanModel.CapabilityStatus.ValueString()
			provisionedPlan.SetCapabilityStatus(&planCapabilityStatus)
		} else {
		}

		if !provisionedPlanModel.ProvisioningStatus.IsUnknown() {
			planProvisioningStatus := provisionedPlanModel.ProvisioningStatus.ValueString()
			provisionedPlan.SetProvisioningStatus(&planProvisioningStatus)
		} else {
		}

		if !provisionedPlanModel.Service.IsUnknown() {
			planService := provisionedPlanModel.Service.ValueString()
			provisionedPlan.SetService(&planService)
		} else {
		}
	}
	requestBody.SetProvisionedPlans(planProvisionedPlans)

	if len(plan.ProxyAddresses.Elements()) > 0 {
		var proxyAddresses []string
		for _, i := range plan.ProxyAddresses.Elements() {
			proxyAddresses = append(proxyAddresses, i.String())
		}
		requestBody.SetProxyAddresses(proxyAddresses)
	}

	if len(plan.Responsibilities.Elements()) > 0 {
		var responsibilities []string
		for _, i := range plan.Responsibilities.Elements() {
			responsibilities = append(responsibilities, i.String())
		}
		requestBody.SetResponsibilities(responsibilities)
	}

	if len(plan.Schools.Elements()) > 0 {
		var schools []string
		for _, i := range plan.Schools.Elements() {
			schools = append(schools, i.String())
		}
		requestBody.SetSchools(schools)
	}

	if !plan.SecurityIdentifier.IsUnknown() {
		planSecurityIdentifier := plan.SecurityIdentifier.ValueString()
		requestBody.SetSecurityIdentifier(&planSecurityIdentifier)
	} else {
	}

	var planServiceProvisioningErrors []models.ServiceProvisioningErrorable
	for _, i := range plan.ServiceProvisioningErrors.Elements() {
		serviceProvisioningError := models.NewServiceProvisioningError()
		serviceProvisioningErrorModel := userServiceProvisioningErrorsModel{}
		types.ListValueFrom(ctx, i.Type(ctx), &serviceProvisioningErrorModel)

		if !serviceProvisioningErrorModel.CreatedDateTime.IsUnknown() {
			planCreatedDateTime := serviceProvisioningErrorModel.CreatedDateTime.ValueString()
			t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
			serviceProvisioningError.SetCreatedDateTime(&t)
		}

		if !serviceProvisioningErrorModel.IsResolved.IsUnknown() {
			planIsResolved := serviceProvisioningErrorModel.IsResolved.ValueBool()
			serviceProvisioningError.SetIsResolved(&planIsResolved)
		}

		if !serviceProvisioningErrorModel.ServiceInstance.IsUnknown() {
			planServiceInstance := serviceProvisioningErrorModel.ServiceInstance.ValueString()
			serviceProvisioningError.SetServiceInstance(&planServiceInstance)
		} else {
		}
	}
	requestBody.SetServiceProvisioningErrors(planServiceProvisioningErrors)

	if !plan.ShowInAddressList.IsUnknown() {
		planShowInAddressList := plan.ShowInAddressList.ValueBool()
		requestBody.SetShowInAddressList(&planShowInAddressList)
	}

	if !plan.SignInSessionsValidFromDateTime.IsUnknown() {
		planSignInSessionsValidFromDateTime := plan.SignInSessionsValidFromDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planSignInSessionsValidFromDateTime)
		requestBody.SetSignInSessionsValidFromDateTime(&t)
	}

	if len(plan.Skills.Elements()) > 0 {
		var skills []string
		for _, i := range plan.Skills.Elements() {
			skills = append(skills, i.String())
		}
		requestBody.SetSkills(skills)
	}

	if !plan.State.IsUnknown() {
		planState := plan.State.ValueString()
		requestBody.SetState(&planState)
	} else {
	}

	if !plan.StreetAddress.IsUnknown() {
		planStreetAddress := plan.StreetAddress.ValueString()
		requestBody.SetStreetAddress(&planStreetAddress)
	} else {
	}

	if !plan.Surname.IsUnknown() {
		planSurname := plan.Surname.ValueString()
		requestBody.SetSurname(&planSurname)
	} else {
	}

	if !plan.UsageLocation.IsUnknown() {
		planUsageLocation := plan.UsageLocation.ValueString()
		requestBody.SetUsageLocation(&planUsageLocation)
	} else {
	}

	if !plan.UserPrincipalName.IsUnknown() {
		planUserPrincipalName := plan.UserPrincipalName.ValueString()
		requestBody.SetUserPrincipalName(&planUserPrincipalName)
	} else {
	}

	if !plan.UserType.IsUnknown() {
		planUserType := plan.UserType.ValueString()
		requestBody.SetUserType(&planUserType)
	} else {
	}

	// Update user
	_, err := r.client.Users().ByUserId(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
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
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state userModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete user
	err := r.client.Users().ByUserId(state.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			err.Error(),
		)
		return
	}

}
