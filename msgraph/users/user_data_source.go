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
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Computed:    true,
				Optional:    true,
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
			"birthday": schema.StringAttribute{
				Description: "The birthday of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.",
				Computed:    true,
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
			"device_enrollment_limit": schema.Int64Attribute{
				Description: "The limit on the maximum number of devices that the user is permitted to enroll. Allowed values are 5 or 1000.",
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
			"hire_date": schema.StringAttribute{
				Description: "The hire date of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.  Note: This property is specific to SharePoint Online. We recommend using the native employeeHireDate property to set and update hire date values using Microsoft Graph APIs.",
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
			"mailbox_settings": schema.SingleNestedAttribute{
				Description: "Settings for the primary mailbox of the signed-in user. You can get or update settings for sending automatic replies to incoming messages, locale and time zone. Returned only on $select.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"archive_folder": schema.StringAttribute{
						Description: "Folder ID of an archive folder for the user.",
						Computed:    true,
					},
					"automatic_replies_setting": schema.SingleNestedAttribute{
						Description: "Configuration settings to automatically notify the sender of an incoming email with a message from the signed-in user.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"external_audience": schema.StringAttribute{
								Description: "The set of audience external to the signed-in user's organization who will receive the ExternalReplyMessage, if Status is AlwaysEnabled or Scheduled. The possible values are: none, contactsOnly, all.",
								Computed:    true,
							},
							"external_reply_message": schema.StringAttribute{
								Description: "The automatic reply to send to the specified external audience, if Status is AlwaysEnabled or Scheduled.",
								Computed:    true,
							},
							"internal_reply_message": schema.StringAttribute{
								Description: "The automatic reply to send to the audience internal to the signed-in user's organization, if Status is AlwaysEnabled or Scheduled.",
								Computed:    true,
							},
							"scheduled_end_date_time": schema.SingleNestedAttribute{
								Description: "The date and time that automatic replies are set to end, if Status is set to Scheduled.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"date_time": schema.StringAttribute{
										Description: "A single point of time in a combined date and time representation ({date}T{time}; for example, 2017-08-29T04:00:00.0000000).",
										Computed:    true,
									},
									"time_zone": schema.StringAttribute{
										Description: "Represents a time zone, for example, 'Pacific Standard Time'. See below for more possible values.",
										Computed:    true,
									},
								},
							},
							"scheduled_start_date_time": schema.SingleNestedAttribute{
								Description: "The date and time that automatic replies are set to begin, if Status is set to Scheduled.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"date_time": schema.StringAttribute{
										Description: "A single point of time in a combined date and time representation ({date}T{time}; for example, 2017-08-29T04:00:00.0000000).",
										Computed:    true,
									},
									"time_zone": schema.StringAttribute{
										Description: "Represents a time zone, for example, 'Pacific Standard Time'. See below for more possible values.",
										Computed:    true,
									},
								},
							},
							"status": schema.StringAttribute{
								Description: "Configurations status for automatic replies. The possible values are: disabled, alwaysEnabled, scheduled.",
								Computed:    true,
							},
						},
					},
					"date_format": schema.StringAttribute{
						Description: "The date format for the user's mailbox.",
						Computed:    true,
					},
					"delegate_meeting_message_delivery_options": schema.StringAttribute{
						Description: "If the user has a calendar delegate, this specifies whether the delegate, mailbox owner, or both receive meeting messages and meeting responses. Possible values are: sendToDelegateAndInformationToPrincipal, sendToDelegateAndPrincipal, sendToDelegateOnly.",
						Computed:    true,
					},
					"language": schema.SingleNestedAttribute{
						Description: "The locale information for the user, including the preferred language and country/region.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"display_name": schema.StringAttribute{
								Description: "A name representing the user's locale in natural language, for example, 'English (United States)'.",
								Computed:    true,
							},
							"locale": schema.StringAttribute{
								Description: "A locale representation for the user, which includes the user's preferred language and country/region. For example, 'en-us'. The language component follows 2-letter codes as defined in ISO 639-1, and the country component follows 2-letter codes as defined in ISO 3166-1 alpha-2.",
								Computed:    true,
							},
						},
					},
					"time_format": schema.StringAttribute{
						Description: "The time format for the user's mailbox.",
						Computed:    true,
					},
					"time_zone": schema.StringAttribute{
						Description: "The default time zone for the user's mailbox.",
						Computed:    true,
					},
					"user_purpose": schema.StringAttribute{
						Description: "The purpose of the mailbox. Differentiates a mailbox for a single user from a shared mailbox and equipment mailbox in Exchange Online. Possible values are: user, linked, shared, room, equipment, others, unknownFutureValue. Read-only.",
						Computed:    true,
					},
					"working_hours": schema.SingleNestedAttribute{
						Description: "The days of the week and hours in a specific time zone that the user works.",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"days_of_week": schema.ListAttribute{
								Description: "The days of the week on which the user works.",
								Computed:    true,
								ElementType: types.StringType,
							},
							"end_time": schema.StringAttribute{
								Description: "The time of the day that the user stops working.",
								Computed:    true,
							},
							"start_time": schema.StringAttribute{
								Description: "The time of the day that the user starts working.",
								Computed:    true,
							},
							"time_zone": schema.SingleNestedAttribute{
								Description: "The time zone to which the working hours apply.",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Description: "The name of a time zone. It can be a standard time zone name such as 'Hawaii-Aleutian Standard Time', or 'Customized Time Zone' for a custom time zone.",
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
			"mobile_phone": schema.StringAttribute{
				Description: "The primary cellular telephone number for the user. Read-only for users synced from on-premises directory. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values) and $search.",
				Computed:    true,
			},
			"my_site": schema.StringAttribute{
				Description: "The URL for the user's personal site. Returned only on $select.",
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
				Description: "The preferred language for the user. The preferred language format is based on RFC 4646. The name is a combination of an ISO 639 two-letter lowercase culture code associated with the language, and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: 'en-US', or 'es-ES'. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values)",
				Computed:    true,
			},
			"preferred_name": schema.StringAttribute{
				Description: "The preferred name for the user. Not Supported. This attribute returns an empty string.Returned only on $select.",
				Computed:    true,
			},
			"print": schema.SingleNestedAttribute{
				Description: "",
				Computed:    true,
				Attributes:  map[string]schema.Attribute{},
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
	}
}

type userDataSourceModel struct {
	Id                              types.String                                      `tfsdk:"id"`
	DeletedDateTime                 types.String                                      `tfsdk:"deleted_date_time"`
	AboutMe                         types.String                                      `tfsdk:"about_me"`
	AccountEnabled                  types.Bool                                        `tfsdk:"account_enabled"`
	AgeGroup                        types.String                                      `tfsdk:"age_group"`
	AssignedLicenses                []userAssignedLicensesDataSourceModel             `tfsdk:"assigned_licenses"`
	AssignedPlans                   []userAssignedPlansDataSourceModel                `tfsdk:"assigned_plans"`
	AuthorizationInfo               *userAuthorizationInfoDataSourceModel             `tfsdk:"authorization_info"`
	Birthday                        types.String                                      `tfsdk:"birthday"`
	BusinessPhones                  []types.String                                    `tfsdk:"business_phones"`
	City                            types.String                                      `tfsdk:"city"`
	CompanyName                     types.String                                      `tfsdk:"company_name"`
	ConsentProvidedForMinor         types.String                                      `tfsdk:"consent_provided_for_minor"`
	Country                         types.String                                      `tfsdk:"country"`
	CreatedDateTime                 types.String                                      `tfsdk:"created_date_time"`
	CreationType                    types.String                                      `tfsdk:"creation_type"`
	CustomSecurityAttributes        *userCustomSecurityAttributesDataSourceModel      `tfsdk:"custom_security_attributes"`
	Department                      types.String                                      `tfsdk:"department"`
	DeviceEnrollmentLimit           types.Int64                                       `tfsdk:"device_enrollment_limit"`
	DisplayName                     types.String                                      `tfsdk:"display_name"`
	EmployeeHireDate                types.String                                      `tfsdk:"employee_hire_date"`
	EmployeeId                      types.String                                      `tfsdk:"employee_id"`
	EmployeeLeaveDateTime           types.String                                      `tfsdk:"employee_leave_date_time"`
	EmployeeOrgData                 *userEmployeeOrgDataDataSourceModel               `tfsdk:"employee_org_data"`
	EmployeeType                    types.String                                      `tfsdk:"employee_type"`
	ExternalUserState               types.String                                      `tfsdk:"external_user_state"`
	ExternalUserStateChangeDateTime types.String                                      `tfsdk:"external_user_state_change_date_time"`
	FaxNumber                       types.String                                      `tfsdk:"fax_number"`
	GivenName                       types.String                                      `tfsdk:"given_name"`
	HireDate                        types.String                                      `tfsdk:"hire_date"`
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
	MailboxSettings                 *userMailboxSettingsDataSourceModel               `tfsdk:"mailbox_settings"`
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
	Print                           *userPrintDataSourceModel                         `tfsdk:"print"`
	ProvisionedPlans                []userProvisionedPlansDataSourceModel             `tfsdk:"provisioned_plans"`
	ProxyAddresses                  []types.String                                    `tfsdk:"proxy_addresses"`
	Responsibilities                []types.String                                    `tfsdk:"responsibilities"`
	Schools                         []types.String                                    `tfsdk:"schools"`
	SecurityIdentifier              types.String                                      `tfsdk:"security_identifier"`
	ServiceProvisioningErrors       []userServiceProvisioningErrorsDataSourceModel    `tfsdk:"service_provisioning_errors"`
	ShowInAddressList               types.Bool                                        `tfsdk:"show_in_address_list"`
	SignInActivity                  *userSignInActivityDataSourceModel                `tfsdk:"sign_in_activity"`
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

type userAuthorizationInfoDataSourceModel struct {
	CertificateUserIds []types.String `tfsdk:"certificate_user_ids"`
}

type userCustomSecurityAttributesDataSourceModel struct {
}

type userEmployeeOrgDataDataSourceModel struct {
	CostCenter types.String `tfsdk:"cost_center"`
	Division   types.String `tfsdk:"division"`
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

type userMailboxSettingsDataSourceModel struct {
	ArchiveFolder                         types.String                                `tfsdk:"archive_folder"`
	AutomaticRepliesSetting               *userAutomaticRepliesSettingDataSourceModel `tfsdk:"automatic_replies_setting"`
	DateFormat                            types.String                                `tfsdk:"date_format"`
	DelegateMeetingMessageDeliveryOptions types.String                                `tfsdk:"delegate_meeting_message_delivery_options"`
	Language                              *userLanguageDataSourceModel                `tfsdk:"language"`
	TimeFormat                            types.String                                `tfsdk:"time_format"`
	TimeZone                              types.String                                `tfsdk:"time_zone"`
	UserPurpose                           types.String                                `tfsdk:"user_purpose"`
	WorkingHours                          *userWorkingHoursDataSourceModel            `tfsdk:"working_hours"`
}

type userAutomaticRepliesSettingDataSourceModel struct {
	ExternalAudience       types.String                               `tfsdk:"external_audience"`
	ExternalReplyMessage   types.String                               `tfsdk:"external_reply_message"`
	InternalReplyMessage   types.String                               `tfsdk:"internal_reply_message"`
	ScheduledEndDateTime   *userScheduledEndDateTimeDataSourceModel   `tfsdk:"scheduled_end_date_time"`
	ScheduledStartDateTime *userScheduledStartDateTimeDataSourceModel `tfsdk:"scheduled_start_date_time"`
	Status                 types.String                               `tfsdk:"status"`
}

type userScheduledEndDateTimeDataSourceModel struct {
	DateTime types.String `tfsdk:"date_time"`
	TimeZone types.String `tfsdk:"time_zone"`
}

type userScheduledStartDateTimeDataSourceModel struct {
	DateTime types.String `tfsdk:"date_time"`
	TimeZone types.String `tfsdk:"time_zone"`
}

type userLanguageDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Locale      types.String `tfsdk:"locale"`
}

type userWorkingHoursDataSourceModel struct {
	DaysOfWeek []types.String               `tfsdk:"days_of_week"`
	EndTime    types.String                 `tfsdk:"end_time"`
	StartTime  types.String                 `tfsdk:"start_time"`
	TimeZone   *userTimeZoneDataSourceModel `tfsdk:"time_zone"`
}

type userTimeZoneDataSourceModel struct {
	Name types.String `tfsdk:"name"`
}

type userOnPremisesExtensionAttributesDataSourceModel struct {
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

type userPrintDataSourceModel struct {
}

type userProvisionedPlansDataSourceModel struct {
	CapabilityStatus   types.String `tfsdk:"capability_status"`
	ProvisioningStatus types.String `tfsdk:"provisioning_status"`
	Service            types.String `tfsdk:"service"`
}

type userServiceProvisioningErrorsDataSourceModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

type userSignInActivityDataSourceModel struct {
	LastNonInteractiveSignInDateTime  types.String `tfsdk:"last_non_interactive_sign_in_date_time"`
	LastNonInteractiveSignInRequestId types.String `tfsdk:"last_non_interactive_sign_in_request_id"`
	LastSignInDateTime                types.String `tfsdk:"last_sign_in_date_time"`
	LastSignInRequestId               types.String `tfsdk:"last_sign_in_request_id"`
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
		result, err = d.client.Users().ByUserId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else if !state.UserPrincipalName.IsNull() {
		result, err = d.client.Users().ByUserId(state.UserPrincipalName.ValueString()).Get(context.Background(), &qparams)
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
	authorizationInfo := new(userAuthorizationInfoDataSourceModel)
	if result.GetAuthorizationInfo() != nil {

		for _, value := range result.GetAuthorizationInfo().GetCertificateUserIds() {
			state.AuthorizationInfo.CertificateUserIds = append(state.AuthorizationInfo.CertificateUserIds, types.StringValue(value))
		}
	}
	state.AuthorizationInfo = authorizationInfo
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
	customSecurityAttributes := new(userCustomSecurityAttributesDataSourceModel)
	if result.GetCustomSecurityAttributes() != nil {

	}
	state.CustomSecurityAttributes = customSecurityAttributes
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
	employeeOrgData := new(userEmployeeOrgDataDataSourceModel)
	if result.GetEmployeeOrgData() != nil {

		if result.GetEmployeeOrgData().GetCostCenter() != nil {
			state.EmployeeOrgData.CostCenter = types.StringValue(*result.GetEmployeeOrgData().GetCostCenter())
		}
		if result.GetEmployeeOrgData().GetDivision() != nil {
			state.EmployeeOrgData.Division = types.StringValue(*result.GetEmployeeOrgData().GetDivision())
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
	mailboxSettings := new(userMailboxSettingsDataSourceModel)
	if result.GetMailboxSettings() != nil {

		if result.GetMailboxSettings().GetArchiveFolder() != nil {
			state.MailboxSettings.ArchiveFolder = types.StringValue(*result.GetMailboxSettings().GetArchiveFolder())
		}
		automaticRepliesSetting := new(userAutomaticRepliesSettingDataSourceModel)
		if result.GetMailboxSettings().GetAutomaticRepliesSetting() != nil {

			if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetExternalAudience() != nil {
				state.MailboxSettings.AutomaticRepliesSetting.ExternalAudience = types.StringValue(result.GetMailboxSettings().GetAutomaticRepliesSetting().GetExternalAudience().String())
			}
			if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetExternalReplyMessage() != nil {
				state.MailboxSettings.AutomaticRepliesSetting.ExternalReplyMessage = types.StringValue(*result.GetMailboxSettings().GetAutomaticRepliesSetting().GetExternalReplyMessage())
			}
			if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetInternalReplyMessage() != nil {
				state.MailboxSettings.AutomaticRepliesSetting.InternalReplyMessage = types.StringValue(*result.GetMailboxSettings().GetAutomaticRepliesSetting().GetInternalReplyMessage())
			}
			scheduledEndDateTime := new(userScheduledEndDateTimeDataSourceModel)
			if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledEndDateTime() != nil {

				if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledEndDateTime().GetDateTime() != nil {
					state.MailboxSettings.AutomaticRepliesSetting.ScheduledEndDateTime.DateTime = types.StringValue(*result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledEndDateTime().GetDateTime())
				}
				if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledEndDateTime().GetTimeZone() != nil {
					state.MailboxSettings.AutomaticRepliesSetting.ScheduledEndDateTime.TimeZone = types.StringValue(*result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledEndDateTime().GetTimeZone())
				}
			}
			state.MailboxSettings.AutomaticRepliesSetting.ScheduledEndDateTime = scheduledEndDateTime
			scheduledStartDateTime := new(userScheduledStartDateTimeDataSourceModel)
			if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledStartDateTime() != nil {

				if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledStartDateTime().GetDateTime() != nil {
					state.MailboxSettings.AutomaticRepliesSetting.ScheduledStartDateTime.DateTime = types.StringValue(*result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledStartDateTime().GetDateTime())
				}
				if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledStartDateTime().GetTimeZone() != nil {
					state.MailboxSettings.AutomaticRepliesSetting.ScheduledStartDateTime.TimeZone = types.StringValue(*result.GetMailboxSettings().GetAutomaticRepliesSetting().GetScheduledStartDateTime().GetTimeZone())
				}
			}
			state.MailboxSettings.AutomaticRepliesSetting.ScheduledStartDateTime = scheduledStartDateTime
			if result.GetMailboxSettings().GetAutomaticRepliesSetting().GetStatus() != nil {
				state.MailboxSettings.AutomaticRepliesSetting.Status = types.StringValue(result.GetMailboxSettings().GetAutomaticRepliesSetting().GetStatus().String())
			}
		}
		state.MailboxSettings.AutomaticRepliesSetting = automaticRepliesSetting
		if result.GetMailboxSettings().GetDateFormat() != nil {
			state.MailboxSettings.DateFormat = types.StringValue(*result.GetMailboxSettings().GetDateFormat())
		}
		if result.GetMailboxSettings().GetDelegateMeetingMessageDeliveryOptions() != nil {
			state.MailboxSettings.DelegateMeetingMessageDeliveryOptions = types.StringValue(result.GetMailboxSettings().GetDelegateMeetingMessageDeliveryOptions().String())
		}
		language := new(userLanguageDataSourceModel)
		if result.GetMailboxSettings().GetLanguage() != nil {

			if result.GetMailboxSettings().GetLanguage().GetDisplayName() != nil {
				state.MailboxSettings.Language.DisplayName = types.StringValue(*result.GetMailboxSettings().GetLanguage().GetDisplayName())
			}
			if result.GetMailboxSettings().GetLanguage().GetLocale() != nil {
				state.MailboxSettings.Language.Locale = types.StringValue(*result.GetMailboxSettings().GetLanguage().GetLocale())
			}
		}
		state.MailboxSettings.Language = language
		if result.GetMailboxSettings().GetTimeFormat() != nil {
			state.MailboxSettings.TimeFormat = types.StringValue(*result.GetMailboxSettings().GetTimeFormat())
		}
		if result.GetMailboxSettings().GetTimeZone() != nil {
			state.MailboxSettings.TimeZone = types.StringValue(*result.GetMailboxSettings().GetTimeZone())
		}
		if result.GetMailboxSettings().GetUserPurpose() != nil {
			state.MailboxSettings.UserPurpose = types.StringValue(result.GetMailboxSettings().GetUserPurpose().String())
		}
		workingHours := new(userWorkingHoursDataSourceModel)
		if result.GetMailboxSettings().GetWorkingHours() != nil {

			for _, value := range result.GetMailboxSettings().GetWorkingHours().GetDaysOfWeek() {
				state.MailboxSettings.WorkingHours.DaysOfWeek = append(state.MailboxSettings.WorkingHours.DaysOfWeek, types.StringValue(value.String()))
			}
			if result.GetMailboxSettings().GetWorkingHours().GetEndTime() != nil {
				state.MailboxSettings.WorkingHours.EndTime = types.StringValue(result.GetMailboxSettings().GetWorkingHours().GetEndTime().String())
			}
			if result.GetMailboxSettings().GetWorkingHours().GetStartTime() != nil {
				state.MailboxSettings.WorkingHours.StartTime = types.StringValue(result.GetMailboxSettings().GetWorkingHours().GetStartTime().String())
			}
			timeZone := new(userTimeZoneDataSourceModel)
			if result.GetMailboxSettings().GetWorkingHours().GetTimeZone() != nil {

				if result.GetMailboxSettings().GetWorkingHours().GetTimeZone().GetName() != nil {
					state.MailboxSettings.WorkingHours.TimeZone.Name = types.StringValue(*result.GetMailboxSettings().GetWorkingHours().GetTimeZone().GetName())
				}
			}
			state.MailboxSettings.WorkingHours.TimeZone = timeZone
		}
		state.MailboxSettings.WorkingHours = workingHours
	}
	state.MailboxSettings = mailboxSettings
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
	state.PasswordProfile = new(userPasswordProfileDataSourceModel)
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
	print := new(userPrintDataSourceModel)
	if result.GetPrint() != nil {

	}
	state.Print = print
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
	for _, value := range result.GetServiceProvisioningErrors() {
		serviceProvisioningErrors := new(userServiceProvisioningErrorsDataSourceModel)

		if value.GetCreatedDateTime() != nil {
			serviceProvisioningErrors.CreatedDateTime = types.StringValue(value.GetCreatedDateTime().String())
		}
		if value.GetIsResolved() != nil {
			serviceProvisioningErrors.IsResolved = types.BoolValue(*value.GetIsResolved())
		}
		if value.GetServiceInstance() != nil {
			serviceProvisioningErrors.ServiceInstance = types.StringValue(*value.GetServiceInstance())
		}
		state.ServiceProvisioningErrors = append(state.ServiceProvisioningErrors, *serviceProvisioningErrors)
	}
	if result.GetShowInAddressList() != nil {
		state.ShowInAddressList = types.BoolValue(*result.GetShowInAddressList())
	}
	signInActivity := new(userSignInActivityDataSourceModel)
	if result.GetSignInActivity() != nil {

		if result.GetSignInActivity().GetLastNonInteractiveSignInDateTime() != nil {
			state.SignInActivity.LastNonInteractiveSignInDateTime = types.StringValue(result.GetSignInActivity().GetLastNonInteractiveSignInDateTime().String())
		}
		if result.GetSignInActivity().GetLastNonInteractiveSignInRequestId() != nil {
			state.SignInActivity.LastNonInteractiveSignInRequestId = types.StringValue(*result.GetSignInActivity().GetLastNonInteractiveSignInRequestId())
		}
		if result.GetSignInActivity().GetLastSignInDateTime() != nil {
			state.SignInActivity.LastSignInDateTime = types.StringValue(result.GetSignInActivity().GetLastSignInDateTime().String())
		}
		if result.GetSignInActivity().GetLastSignInRequestId() != nil {
			state.SignInActivity.LastSignInRequestId = types.StringValue(*result.GetSignInActivity().GetLastSignInRequestId())
		}
	}
	state.SignInActivity = signInActivity
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
