package users

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
	"github.com/microsoftgraph/msgraph-sdk-go/users"

	"terraform-provider-msgraph/planmodifiers/boolplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/listplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/objectplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/stringplanmodifiers"
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
			"about_me": schema.StringAttribute{
				Description: "A freeform text entry field for the user to describe themselves. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"account_enabled": schema.BoolAttribute{
				Description: "true if the account is enabled; otherwise, false. This property is required when a user is created. Returned only on $select. Supports $filter (eq, ne, not, and in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"age_group": schema.StringAttribute{
				Description: "Sets the age group of the user. Allowed values: null, Minor, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Description: "The licenses that are assigned to the user, including inherited (group-based) licenses. This property doesn't differentiate between directly assigned and inherited licenses. Use the licenseAssignmentStates property to identify the directly assigned and inherited licenses. Not nullable. Returned only on $select. Supports $filter (eq, not, /$count eq 0, /$count ne 0).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"disabled_plans": schema.ListAttribute{
							Description: "A collection of the unique identifiers for plans that have been disabled. IDs are available in servicePlans > servicePlanId in the tenant's subscribedSkus or serviceStatus > servicePlanId in the tenant's companySubscription.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.List{
								listplanmodifiers.UseStateForUnconfigured(),
							},
							ElementType: types.StringType,
						},
						"sku_id": schema.StringAttribute{
							Description: "The unique identifier for the SKU. Corresponds to the skuId from subscribedSkus or companySubscription.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"assigned_plans": schema.ListNestedAttribute{
				Description: "The plans that are assigned to the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq and not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assigned_date_time": schema.StringAttribute{
							Description: "The date and time at which the plan was assigned. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"capability_status": schema.StringAttribute{
							Description: "Condition of the capability assignment. The possible values are Enabled, Warning, Suspended, Deleted, LockedOut. See a detailed description of each value.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"service": schema.StringAttribute{
							Description: "The name of the service; for example, exchange.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"service_plan_id": schema.StringAttribute{
							Description: "A GUID that identifies the service plan. For a complete list of GUIDs and their equivalent friendly service names, see Product names and service plan identifiers for licensing.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"authorization_info": schema.SingleNestedAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"certificate_user_ids": schema.ListAttribute{
						Description: "",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.List{
							listplanmodifiers.UseStateForUnconfigured(),
						},
						ElementType: types.StringType,
					},
				},
			},
			"birthday": schema.StringAttribute{
				Description: "The birthday of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"business_phones": schema.ListAttribute{
				Description: "The telephone numbers for the user. NOTE: Although it's a string collection, only one number can be set for this property. Read-only for users synced from the on-premises directory. Returned by default. Supports $filter (eq, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"city": schema.StringAttribute{
				Description: "The city where the user is located. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"company_name": schema.StringAttribute{
				Description: "The name of the company that the user is associated with. This property can be useful for describing the company that a guest comes from. The maximum length is 64 characters.Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"consent_provided_for_minor": schema.StringAttribute{
				Description: "Sets whether consent was obtained for minors. Allowed values: null, Granted, Denied, and NotRequired. For more information, see legal age group property definitions. Returned only on $select. Supports $filter (eq, ne, not, and in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"country": schema.StringAttribute{
				Description: "The country/region where the user is located; for example, US or UK. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"created_date_time": schema.StringAttribute{
				Description: "The date and time the user was created, in ISO 8601 format and UTC. The value can't be modified and is automatically populated when the entity is created. Nullable. For on-premises users, the value represents when they were first created in Microsoft Entra ID. Property is null for some users created before June 2018 and on-premises users that were synced to Microsoft Entra ID before June 2018. Read-only. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"creation_type": schema.StringAttribute{
				Description: "Indicates whether the user account was created through one of the following methods:  As a regular school or work account (null). As an external account (Invitation). As a local account for an Azure Active Directory B2C tenant (LocalAccount). Through self-service sign-up by an internal user using email verification (EmailVerified). Through self-service sign-up by a guest signing up through a link that is part of a user flow (SelfServiceSignUp). Read-only.Returned only on $select. Supports $filter (eq, ne, not, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"department": schema.StringAttribute{
				Description: "The name of the department in which the user works. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The name displayed in the address book for the user. This value is usually the combination of the user's first name, middle initial, and family name. This property is required when a user is created and it can't be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values), $orderby, and $search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"employee_hire_date": schema.StringAttribute{
				Description: "The date and time when the user was hired or will start work in a future hire. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"employee_id": schema.StringAttribute{
				Description: "The employee identifier assigned to the user by the organization. The maximum length is 16 characters. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"employee_leave_date_time": schema.StringAttribute{
				Description: "The date and time when the user left or will leave the organization. To read this property, the calling app must be assigned the User-LifeCycleInfo.Read.All permission. To write this property, the calling app must be assigned the User.Read.All and User-LifeCycleInfo.ReadWrite.All permissions. To read this property in delegated scenarios, the admin needs at least one of the following Microsoft Entra roles: Lifecycle Workflows Administrator (least privilege), Global Reader. To write this property in delegated scenarios, the admin needs the Global Administrator role. Supports $filter (eq, ne, not , ge, le, in). For more information, see Configure the employeeLeaveDateTime property for a user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"employee_org_data": schema.SingleNestedAttribute{
				Description: "Represents organization data (for example, division and costCenter) associated with a user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"cost_center": schema.StringAttribute{
						Description: "The cost center associated with the user. Returned only on $select. Supports $filter.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"division": schema.StringAttribute{
						Description: "The name of the division in which the user works. Returned only on $select. Supports $filter.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"employee_type": schema.StringAttribute{
				Description: "Captures enterprise worker type. For example, Employee, Contractor, Consultant, or Vendor. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"external_user_state": schema.StringAttribute{
				Description: "For a guest invited to the tenant using the invitation API, this property represents the invited user's invitation status. For invited users, the state can be PendingAcceptance or Accepted, or null for all other users. Returned only on $select. Supports $filter (eq, ne, not , in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"external_user_state_change_date_time": schema.StringAttribute{
				Description: "Shows the timestamp for the latest change to the externalUserState property. Returned only on $select. Supports $filter (eq, ne, not , in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"fax_number": schema.StringAttribute{
				Description: "The fax number of the user. Returned only on $select. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"given_name": schema.StringAttribute{
				Description: "The given name (first name) of the user. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"hire_date": schema.StringAttribute{
				Description: "The hire date of the user. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014, is 2014-01-01T00:00:00Z. Returned only on $select.  Note: This property is specific to SharePoint in Microsoft 365. We recommend using the native employeeHireDate property to set and update hire date values using Microsoft Graph APIs.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"identities": schema.ListNestedAttribute{
				Description: "Represents the identities that can be used to sign in to this user account. Microsoft (also known as a local account), organizations, or social identity providers such as Facebook, Google, and Microsoft can provide identity and tie it to a user account. It might contain multiple items with the same signInType value. Returned only on $select.  Supports $filter (eq) with limitations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"issuer": schema.StringAttribute{
							Description: "Specifies the issuer of the identity, for example facebook.com. 512 character limit. For local accounts (where signInType isn't federated), this property is the local default domain name for the tenant, for example contoso.com.  For guests from other Microsoft Entra organizations, this is the domain of the federated organization, for example contoso.com. For more information about filtering behavior for this property, see Filtering on the identities property of a user.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"issuer_assigned_id": schema.StringAttribute{
							Description: "Specifies the unique identifier assigned to the user by the issuer. 64 character limit. The combination of issuer and issuerAssignedId must be unique within the organization. Represents the sign-in name for the user, when signInType is set to emailAddress or userName (also known as local accounts).When signInType is set to: emailAddress (or a custom string that starts with emailAddress like emailAddress1), issuerAssignedId must be a valid email addressuserName, issuerAssignedId must begin with an alphabetical character or number, and can only contain alphanumeric characters and the following symbols: - or _  For more information about filtering behavior for this property, see Filtering on the identities property of a user.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"sign_in_type": schema.StringAttribute{
							Description: "Specifies the user sign-in types in your directory, such as emailAddress, userName, federated, or userPrincipalName. federated represents a unique identifier for a user from an issuer that can be in any format chosen by the issuer. Setting or updating a userPrincipalName identity updates the value of the userPrincipalName property on the user object. The validations performed on the userPrincipalName property on the user object, for example, verified domains and acceptable characters, are performed when setting or updating a userPrincipalName identity. Extra validation is enforced on issuerAssignedId when the sign-in type is set to emailAddress or userName. This property can also be set to any custom string.  For more information about filtering behavior for this property, see Filtering on the identities property of a user.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"im_addresses": schema.ListAttribute{
				Description: "The instant message voice-over IP (VOIP) session initiation protocol (SIP) addresses for the user. Read-only. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"interests": schema.ListAttribute{
				Description: "A list for the user to describe their interests. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"is_management_restricted": schema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_resource_account": schema.BoolAttribute{
				Description: "Don't use – reserved for future use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"job_title": schema.StringAttribute{
				Description: "The user's job title. Maximum length is 128 characters. Returned by default. Supports $filter (eq, ne, not , ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"last_password_change_date_time": schema.StringAttribute{
				Description: "The time when this Microsoft Entra user last changed their password or when their password was created, whichever date the latest action was performed. The date and time information uses ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"legal_age_group_classification": schema.StringAttribute{
				Description: "Used by enterprise applications to determine the legal age group of the user. This property is read-only and calculated based on ageGroup and consentProvidedForMinor properties. Allowed values: null, MinorWithOutParentalConsent, MinorWithParentalConsent, MinorNoParentalConsentRequired, NotAdult, and Adult. For more information, see legal age group property definitions. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"license_assignment_states": schema.ListNestedAttribute{
				Description: "State of license assignments for this user. Also indicates licenses that are directly assigned or the user inherited through group memberships. Read-only. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"assigned_by_group": schema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"disabled_plans": schema.ListAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.List{
								listplanmodifiers.UseStateForUnconfigured(),
							},
							ElementType: types.StringType,
						},
						"error": schema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"last_updated_date_time": schema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"sku_id": schema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"state": schema.StringAttribute{
							Description: "",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"mail": schema.StringAttribute{
				Description: "The SMTP address for the user, for example, jeff@contoso.com. Changes to this property update the user's proxyAddresses collection to include the value as an SMTP address. This property can't contain accent characters.  NOTE: We don't recommend updating this property for Azure AD B2C user profiles. Use the otherMails property instead. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"mail_nickname": schema.StringAttribute{
				Description: "The mail alias for the user. This property must be specified when a user is created. Maximum length is 64 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"mobile_phone": schema.StringAttribute{
				Description: "The primary cellular telephone number for the user. Read-only for users synced from the on-premises directory. Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values) and $search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"my_site": schema.StringAttribute{
				Description: "The URL for the user's site. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"office_location": schema.StringAttribute{
				Description: "The office location in the user's place of business. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_distinguished_name": schema.StringAttribute{
				Description: "Contains the on-premises Active Directory distinguished name or DN. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_domain_name": schema.StringAttribute{
				Description: "Contains the on-premises domainFQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_extension_attributes": schema.SingleNestedAttribute{
				Description: "Contains extensionAttributes1-15 for the user. These extension attributes are also known as Exchange custom attributes 1-15. Each attribute can store up to 1024 characters. For an onPremisesSyncEnabled user, the source of authority for this set of properties is the on-premises and is read-only. For a cloud-only user (where onPremisesSyncEnabled is false), these properties can be set during the creation or update of a user object.  For a cloud-only user previously synced from on-premises Active Directory, these properties are read-only in Microsoft Graph but can be fully managed through the Exchange Admin Center or the Exchange Online V2 module in PowerShell. Returned only on $select. Supports $filter (eq, ne, not, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"extension_attribute_1": schema.StringAttribute{
						Description: "First customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_10": schema.StringAttribute{
						Description: "Tenth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_11": schema.StringAttribute{
						Description: "Eleventh customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_12": schema.StringAttribute{
						Description: "Twelfth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_13": schema.StringAttribute{
						Description: "Thirteenth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_14": schema.StringAttribute{
						Description: "Fourteenth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_15": schema.StringAttribute{
						Description: "Fifteenth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_2": schema.StringAttribute{
						Description: "Second customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_3": schema.StringAttribute{
						Description: "Third customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_4": schema.StringAttribute{
						Description: "Fourth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_5": schema.StringAttribute{
						Description: "Fifth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_6": schema.StringAttribute{
						Description: "Sixth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_7": schema.StringAttribute{
						Description: "Seventh customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_8": schema.StringAttribute{
						Description: "Eighth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"extension_attribute_9": schema.StringAttribute{
						Description: "Ninth customizable extension attribute.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"on_premises_immutable_id": schema.StringAttribute{
				Description: "This property is used to associate an on-premises Active Directory user account to their Microsoft Entra user object. This property must be specified when creating a new user account in the Graph if you're using a federated domain for the user's userPrincipalName (UPN) property. NOTE: The $ and _ characters can't be used when specifying this property. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "Indicates the last time at which the object was synced with the on-premises directory; for example: 2013-02-16T03:04:54Z. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors when using Microsoft synchronization product during provisioning. Returned only on $select. Supports $filter (eq, not, ge, le).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category": schema.StringAttribute{
							Description: "Category of the provisioning error. Note: Currently, there is only one possible value. Possible value: PropertyConflict - indicates a property value is not unique. Other objects contain the same value for the property.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"occurred_date_time": schema.StringAttribute{
							Description: "The date and time at which the error occurred.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"property_causing_error": schema.StringAttribute{
							Description: "Name of the directory property causing the error. Current possible values: UserPrincipalName or ProxyAddress",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"value": schema.StringAttribute{
							Description: "Value of the property causing the error.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"on_premises_sam_account_name": schema.StringAttribute{
				Description: "Contains the on-premises samAccountName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_security_identifier": schema.StringAttribute{
				Description: "Contains the on-premises security identifier (SID) for the user that was synchronized from on-premises to the cloud. Read-only. Returned only on $select. Supports $filter (eq including on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Description: "true if this user object is currently being synced from an on-premises Active Directory (AD); otherwise the user isn't being synced and can be managed in Microsoft Entra ID. Read-only. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_user_principal_name": schema.StringAttribute{
				Description: "Contains the on-premises userPrincipalName synchronized from the on-premises directory. The property is only populated for customers who are synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect. Read-only. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"other_mails": schema.ListAttribute{
				Description: "A list of other email addresses for the user; for example: ['bob@contoso.com', 'Robert@fabrikam.com']. NOTE: This property can't contain accent characters. Returned only on $select. Supports $filter (eq, not, ge, le, in, startsWith, endsWith, /$count eq 0, /$count ne 0).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"password_policies": schema.StringAttribute{
				Description: "Specifies password policies for the user. This value is an enumeration with one possible value being DisableStrongPassword, which allows weaker passwords than the default policy to be specified. DisablePasswordExpiration can also be specified. The two might be specified together; for example: DisablePasswordExpiration, DisableStrongPassword. Returned only on $select. For more information on the default password policies, see Microsoft Entra password policies. Supports $filter (ne, not, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"password_profile": schema.SingleNestedAttribute{
				Description: "Specifies the password profile for the user. The profile contains the user's password. This property is required when a user is created. The password in the profile must satisfy minimum requirements as specified by the passwordPolicies property. By default, a strong password is required. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). To update this property:  In delegated access, the calling app must be assigned the Directory.AccessAsUser.All delegated permission on behalf of the signed-in user.  In application-only access, the calling app must be assigned the User.ReadWrite.All (least privilege) or Directory.ReadWrite.All (higher privilege) application permission and at least the User Administrator Microsoft Entra role.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"force_change_password_next_sign_in": schema.BoolAttribute{
						Description: "true if the user must change their password on the next sign-in; otherwise false.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"force_change_password_next_sign_in_with_mfa": schema.BoolAttribute{
						Description: "If true, at next sign-in, the user must perform a multifactor authentication (MFA) before being forced to change their password. The behavior is identical to forceChangePasswordNextSignIn except that the user is required to first perform a multifactor authentication before password change. After a password change, this property will be automatically reset to false. If not set, default is false.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"password": schema.StringAttribute{
						Description: "The password for the user. This property is required when a user is created. It can be updated, but the user will be required to change the password on the next sign-in. The password must satisfy minimum requirements as specified by the user's passwordPolicies property. By default, a strong password is required.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"past_projects": schema.ListAttribute{
				Description: "A list for the user to enumerate their past projects. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"postal_code": schema.StringAttribute{
				Description: "The postal code for the user's postal address. The postal code is specific to the user's country/region. In the United States of America, this attribute contains the ZIP code. Maximum length is 40 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"preferred_data_location": schema.StringAttribute{
				Description: "The preferred data location for the user. For more information, see OneDrive Online Multi-Geo.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"preferred_language": schema.StringAttribute{
				Description: "The preferred language for the user. The preferred language format is based on RFC 4646. The name is a combination of an ISO 639 two-letter lowercase culture code associated with the language, and an ISO 3166 two-letter uppercase subculture code associated with the country or region. Example: 'en-US', or 'es-ES'. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values)",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"preferred_name": schema.StringAttribute{
				Description: "The preferred name for the user. Not Supported. This attribute returns an empty string.Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"provisioned_plans": schema.ListNestedAttribute{
				Description: "The plans that are provisioned for the user. Read-only. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"capability_status": schema.StringAttribute{
							Description: "For example, 'Enabled'.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"provisioning_status": schema.StringAttribute{
							Description: "For example, 'Success'.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"service": schema.StringAttribute{
							Description: "The name of the service; for example, 'AccessControlS2S'",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"proxy_addresses": schema.ListAttribute{
				Description: "For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. Changes to the mail property update this collection to include the value as an SMTP address. For more information, see mail and proxyAddresses properties. The proxy address prefixed with SMTP (capitalized) is the primary proxy address, while those addresses prefixed with smtp are the secondary proxy addresses. For Azure AD B2C accounts, this property has a limit of 10 unique addresses. Read-only in Microsoft Graph; you can update this property only through the Microsoft 365 admin center. Not nullable. Returned only on $select. Supports $filter (eq, not, ge, le, startsWith, endsWith, /$count eq 0, /$count ne 0).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"responsibilities": schema.ListAttribute{
				Description: "A list for the user to enumerate their responsibilities. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"schools": schema.ListAttribute{
				Description: "A list for the user to enumerate the schools they attended. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"security_identifier": schema.StringAttribute{
				Description: "Security identifier (SID) of the user, used in Windows scenarios. Read-only. Returned by default. Supports $select and $filter (eq, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"service_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors published by a federated service describing a nontransient, service-specific error regarding the properties or link from a user object.  Supports $filter (eq, not, for isResolved and serviceInstance).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"created_date_time": schema.StringAttribute{
							Description: "The date and time at which the error occurred.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"is_resolved": schema.BoolAttribute{
							Description: "Indicates whether the error has been attended to.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"service_instance": schema.StringAttribute{
							Description: "Qualified service instance (for example, 'SharePoint/Dublin') that published the service error information.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"show_in_address_list": schema.BoolAttribute{
				Description: "Do not use in Microsoft Graph. Manage this property through the Microsoft 365 admin center instead. Represents whether the user should be included in the Outlook global address list. See Known issue.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"sign_in_activity": schema.SingleNestedAttribute{
				Description: "Get the last signed-in date and request ID of the sign-in for a given user. Read-only.Returned only on $select. Supports $filter (eq, ne, not, ge, le) but not with any other filterable properties. Note: Details for this property require a Microsoft Entra ID P1 or P2 license and the AuditLog.Read.All permission.This property isn't returned for a user who never signed in or last signed in before April 2020.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"last_non_interactive_sign_in_date_time": schema.StringAttribute{
						Description: "The last non-interactive sign-in date for a specific user. You can use this field to calculate the last time a client attempted (either successfully or unsuccessfully) to sign in to the directory on behalf of a user. Because some users may use clients to access tenant resources rather than signing into your tenant directly, you can use the non-interactive sign-in date to along with lastSignInDateTime to identify inactive users. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Microsoft Entra ID maintains non-interactive sign-ins going back to May 2020. For more information about using the value of this property, see Manage inactive user accounts in Microsoft Entra ID.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"last_non_interactive_sign_in_request_id": schema.StringAttribute{
						Description: "Request identifier of the last non-interactive sign-in performed by this user.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"last_sign_in_date_time": schema.StringAttribute{
						Description: "The last interactive sign-in date and time for a specific user. You can use this field to calculate the last time a user attempted (either successfully or unsuccessfully) to sign in to the directory with an interactive authentication method. This field can be used to build reports, such as inactive users. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Microsoft Entra ID maintains interactive sign-ins going back to April 2020. For more information about using the value of this property, see Manage inactive user accounts in Microsoft Entra ID.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"last_sign_in_request_id": schema.StringAttribute{
						Description: "Request identifier of the last interactive sign-in performed by this user.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"last_successful_sign_in_date_time": schema.StringAttribute{
						Description: "The date and time of the user's most recent successful sign-in activity. The timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"last_successful_sign_in_request_id": schema.StringAttribute{
						Description: "The request ID of the last successful sign-in.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"sign_in_sessions_valid_from_date_time": schema.StringAttribute{
				Description: "Any refresh tokens or session tokens (session cookies) issued before this time are invalid. Applications get an error when using an invalid refresh or session token to acquire a delegated access token (to access APIs such as Microsoft Graph). If this happens, the application needs to acquire a new refresh token by requesting the authorized endpoint. Read-only. Use revokeSignInSessions to reset. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"skills": schema.ListAttribute{
				Description: "A list for the user to enumerate their skills. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"state": schema.StringAttribute{
				Description: "The state or province in the user's address. Maximum length is 128 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"street_address": schema.StringAttribute{
				Description: "The street address of the user's place of business. Maximum length is 1,024 characters. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"surname": schema.StringAttribute{
				Description: "The user's surname (family name or last name). Maximum length is 64 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"usage_location": schema.StringAttribute{
				Description: "A two-letter country code (ISO standard 3166). Required for users that are assigned licenses due to legal requirements to check for availability of services in countries. Examples include: US, JP, and GB. Not nullable. Returned only on $select. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"user_principal_name": schema.StringAttribute{
				Description: "The user principal name (UPN) of the user. The UPN is an Internet-style sign-in name for the user based on the Internet standard RFC 822. By convention, this value should map to the user's email name. The general format is alias@domain, where the domain must be present in the tenant's collection of verified domains. This property is required when a user is created. The verified domains for the tenant can be accessed from the verifiedDomains property of organization.NOTE: This property can't contain accent characters. Only the following characters are allowed A - Z, a - z, 0 - 9, ' . - _ ! # ^ ~. For the complete list of allowed characters, see username policies. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, endsWith) and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"user_type": schema.StringAttribute{
				Description: "A string value that can be used to classify user types in your directory. The possible values are Member and Guest. Returned only on $select. Supports $filter (eq, ne, not, in, and eq on null values). NOTE: For more information about the permissions for members and guests, see What are the default user permissions in Microsoft Entra ID?",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from Terraform plan
	var tfPlan userModel
	diags := req.Plan.Get(ctx, &tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	sdkModelUser := models.NewUser()

	if !tfPlan.Id.IsUnknown() {
		tfPlanId := tfPlan.Id.ValueString()
		sdkModelUser.SetId(&tfPlanId)
	} else {
		tfPlan.Id = types.StringNull()
	}

	if !tfPlan.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		sdkModelUser.SetDeletedDateTime(&t)
	} else {
		tfPlan.DeletedDateTime = types.StringNull()
	}

	if !tfPlan.AboutMe.IsUnknown() {
		tfPlanAboutMe := tfPlan.AboutMe.ValueString()
		sdkModelUser.SetAboutMe(&tfPlanAboutMe)
	} else {
		tfPlan.AboutMe = types.StringNull()
	}

	if !tfPlan.AccountEnabled.IsUnknown() {
		tfPlanAccountEnabled := tfPlan.AccountEnabled.ValueBool()
		sdkModelUser.SetAccountEnabled(&tfPlanAccountEnabled)
	} else {
		tfPlan.AccountEnabled = types.BoolNull()
	}

	if !tfPlan.AgeGroup.IsUnknown() {
		tfPlanAgeGroup := tfPlan.AgeGroup.ValueString()
		sdkModelUser.SetAgeGroup(&tfPlanAgeGroup)
	} else {
		tfPlan.AgeGroup = types.StringNull()
	}

	if len(tfPlan.AssignedLicenses.Elements()) > 0 {
		var tfPlanAssignedLicenses []models.AssignedLicenseable
		for _, i := range tfPlan.AssignedLicenses.Elements() {
			sdkModelAssignedLicenses := models.NewAssignedLicense()
			tfModelAssignedLicenses := userAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelAssignedLicenses)

			if len(tfModelAssignedLicenses.DisabledPlans.Elements()) > 0 {
				var DisabledPlans []uuid.UUID
				for _, i := range tfModelAssignedLicenses.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				sdkModelAssignedLicenses.SetDisabledPlans(DisabledPlans)
			} else {
				tfModelAssignedLicenses.DisabledPlans = types.ListNull(types.StringType)
			}

			if !tfModelAssignedLicenses.SkuId.IsUnknown() {
				tfPlanSkuId := tfModelAssignedLicenses.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				sdkModelAssignedLicenses.SetSkuId(&u)
			} else {
				tfModelAssignedLicenses.SkuId = types.StringNull()
			}
		}
		sdkModelUser.SetAssignedLicenses(tfPlanAssignedLicenses)
	} else {
		tfPlan.AssignedLicenses = types.ListNull(tfPlan.AssignedLicenses.ElementType(ctx))
	}

	if len(tfPlan.AssignedPlans.Elements()) > 0 {
		var tfPlanAssignedPlans []models.AssignedPlanable
		for _, i := range tfPlan.AssignedPlans.Elements() {
			sdkModelAssignedPlans := models.NewAssignedPlan()
			tfModelAssignedPlans := userAssignedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelAssignedPlans)

			if !tfModelAssignedPlans.AssignedDateTime.IsUnknown() {
				tfPlanAssignedDateTime := tfModelAssignedPlans.AssignedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanAssignedDateTime)
				sdkModelAssignedPlans.SetAssignedDateTime(&t)
			} else {
				tfModelAssignedPlans.AssignedDateTime = types.StringNull()
			}

			if !tfModelAssignedPlans.CapabilityStatus.IsUnknown() {
				tfPlanCapabilityStatus := tfModelAssignedPlans.CapabilityStatus.ValueString()
				sdkModelAssignedPlans.SetCapabilityStatus(&tfPlanCapabilityStatus)
			} else {
				tfModelAssignedPlans.CapabilityStatus = types.StringNull()
			}

			if !tfModelAssignedPlans.Service.IsUnknown() {
				tfPlanService := tfModelAssignedPlans.Service.ValueString()
				sdkModelAssignedPlans.SetService(&tfPlanService)
			} else {
				tfModelAssignedPlans.Service = types.StringNull()
			}

			if !tfModelAssignedPlans.ServicePlanId.IsUnknown() {
				tfPlanServicePlanId := tfModelAssignedPlans.ServicePlanId.ValueString()
				u, _ := uuid.Parse(tfPlanServicePlanId)
				sdkModelAssignedPlans.SetServicePlanId(&u)
			} else {
				tfModelAssignedPlans.ServicePlanId = types.StringNull()
			}
		}
		sdkModelUser.SetAssignedPlans(tfPlanAssignedPlans)
	} else {
		tfPlan.AssignedPlans = types.ListNull(tfPlan.AssignedPlans.ElementType(ctx))
	}

	if !tfPlan.AuthorizationInfo.IsUnknown() {
		sdkModelAuthorizationInfo := models.NewAuthorizationInfo()
		tfModelAuthorizationInfo := userAuthorizationInfoModel{}
		tfPlan.AuthorizationInfo.As(ctx, &tfModelAuthorizationInfo, basetypes.ObjectAsOptions{})

		if len(tfModelAuthorizationInfo.CertificateUserIds.Elements()) > 0 {
			var certificateUserIds []string
			for _, i := range tfModelAuthorizationInfo.CertificateUserIds.Elements() {
				certificateUserIds = append(certificateUserIds, i.String())
			}
			sdkModelAuthorizationInfo.SetCertificateUserIds(certificateUserIds)
		} else {
			tfModelAuthorizationInfo.CertificateUserIds = types.ListNull(types.StringType)
		}
		sdkModelUser.SetAuthorizationInfo(sdkModelAuthorizationInfo)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelAuthorizationInfo.AttributeTypes(), sdkModelAuthorizationInfo)
		tfPlan.AuthorizationInfo = objectValue
	} else {
		tfPlan.AuthorizationInfo = types.ObjectNull(tfPlan.AuthorizationInfo.AttributeTypes(ctx))
	}

	if !tfPlan.Birthday.IsUnknown() {
		tfPlanBirthday := tfPlan.Birthday.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanBirthday)
		sdkModelUser.SetBirthday(&t)
	} else {
		tfPlan.Birthday = types.StringNull()
	}

	if len(tfPlan.BusinessPhones.Elements()) > 0 {
		var businessPhones []string
		for _, i := range tfPlan.BusinessPhones.Elements() {
			businessPhones = append(businessPhones, i.String())
		}
		sdkModelUser.SetBusinessPhones(businessPhones)
	} else {
		tfPlan.BusinessPhones = types.ListNull(types.StringType)
	}

	if !tfPlan.City.IsUnknown() {
		tfPlanCity := tfPlan.City.ValueString()
		sdkModelUser.SetCity(&tfPlanCity)
	} else {
		tfPlan.City = types.StringNull()
	}

	if !tfPlan.CompanyName.IsUnknown() {
		tfPlanCompanyName := tfPlan.CompanyName.ValueString()
		sdkModelUser.SetCompanyName(&tfPlanCompanyName)
	} else {
		tfPlan.CompanyName = types.StringNull()
	}

	if !tfPlan.ConsentProvidedForMinor.IsUnknown() {
		tfPlanConsentProvidedForMinor := tfPlan.ConsentProvidedForMinor.ValueString()
		sdkModelUser.SetConsentProvidedForMinor(&tfPlanConsentProvidedForMinor)
	} else {
		tfPlan.ConsentProvidedForMinor = types.StringNull()
	}

	if !tfPlan.Country.IsUnknown() {
		tfPlanCountry := tfPlan.Country.ValueString()
		sdkModelUser.SetCountry(&tfPlanCountry)
	} else {
		tfPlan.Country = types.StringNull()
	}

	if !tfPlan.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		sdkModelUser.SetCreatedDateTime(&t)
	} else {
		tfPlan.CreatedDateTime = types.StringNull()
	}

	if !tfPlan.CreationType.IsUnknown() {
		tfPlanCreationType := tfPlan.CreationType.ValueString()
		sdkModelUser.SetCreationType(&tfPlanCreationType)
	} else {
		tfPlan.CreationType = types.StringNull()
	}

	if !tfPlan.Department.IsUnknown() {
		tfPlanDepartment := tfPlan.Department.ValueString()
		sdkModelUser.SetDepartment(&tfPlanDepartment)
	} else {
		tfPlan.Department = types.StringNull()
	}

	if !tfPlan.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlan.DisplayName.ValueString()
		sdkModelUser.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlan.DisplayName = types.StringNull()
	}

	if !tfPlan.EmployeeHireDate.IsUnknown() {
		tfPlanEmployeeHireDate := tfPlan.EmployeeHireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanEmployeeHireDate)
		sdkModelUser.SetEmployeeHireDate(&t)
	} else {
		tfPlan.EmployeeHireDate = types.StringNull()
	}

	if !tfPlan.EmployeeId.IsUnknown() {
		tfPlanEmployeeId := tfPlan.EmployeeId.ValueString()
		sdkModelUser.SetEmployeeId(&tfPlanEmployeeId)
	} else {
		tfPlan.EmployeeId = types.StringNull()
	}

	if !tfPlan.EmployeeLeaveDateTime.IsUnknown() {
		tfPlanEmployeeLeaveDateTime := tfPlan.EmployeeLeaveDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanEmployeeLeaveDateTime)
		sdkModelUser.SetEmployeeLeaveDateTime(&t)
	} else {
		tfPlan.EmployeeLeaveDateTime = types.StringNull()
	}

	if !tfPlan.EmployeeOrgData.IsUnknown() {
		sdkModelEmployeeOrgData := models.NewEmployeeOrgData()
		tfModelEmployeeOrgData := userEmployeeOrgDataModel{}
		tfPlan.EmployeeOrgData.As(ctx, &tfModelEmployeeOrgData, basetypes.ObjectAsOptions{})

		if !tfModelEmployeeOrgData.CostCenter.IsUnknown() {
			tfPlanCostCenter := tfModelEmployeeOrgData.CostCenter.ValueString()
			sdkModelEmployeeOrgData.SetCostCenter(&tfPlanCostCenter)
		} else {
			tfModelEmployeeOrgData.CostCenter = types.StringNull()
		}

		if !tfModelEmployeeOrgData.Division.IsUnknown() {
			tfPlanDivision := tfModelEmployeeOrgData.Division.ValueString()
			sdkModelEmployeeOrgData.SetDivision(&tfPlanDivision)
		} else {
			tfModelEmployeeOrgData.Division = types.StringNull()
		}
		sdkModelUser.SetEmployeeOrgData(sdkModelEmployeeOrgData)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelEmployeeOrgData.AttributeTypes(), sdkModelEmployeeOrgData)
		tfPlan.EmployeeOrgData = objectValue
	} else {
		tfPlan.EmployeeOrgData = types.ObjectNull(tfPlan.EmployeeOrgData.AttributeTypes(ctx))
	}

	if !tfPlan.EmployeeType.IsUnknown() {
		tfPlanEmployeeType := tfPlan.EmployeeType.ValueString()
		sdkModelUser.SetEmployeeType(&tfPlanEmployeeType)
	} else {
		tfPlan.EmployeeType = types.StringNull()
	}

	if !tfPlan.ExternalUserState.IsUnknown() {
		tfPlanExternalUserState := tfPlan.ExternalUserState.ValueString()
		sdkModelUser.SetExternalUserState(&tfPlanExternalUserState)
	} else {
		tfPlan.ExternalUserState = types.StringNull()
	}

	if !tfPlan.ExternalUserStateChangeDateTime.IsUnknown() {
		tfPlanExternalUserStateChangeDateTime := tfPlan.ExternalUserStateChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanExternalUserStateChangeDateTime)
		sdkModelUser.SetExternalUserStateChangeDateTime(&t)
	} else {
		tfPlan.ExternalUserStateChangeDateTime = types.StringNull()
	}

	if !tfPlan.FaxNumber.IsUnknown() {
		tfPlanFaxNumber := tfPlan.FaxNumber.ValueString()
		sdkModelUser.SetFaxNumber(&tfPlanFaxNumber)
	} else {
		tfPlan.FaxNumber = types.StringNull()
	}

	if !tfPlan.GivenName.IsUnknown() {
		tfPlanGivenName := tfPlan.GivenName.ValueString()
		sdkModelUser.SetGivenName(&tfPlanGivenName)
	} else {
		tfPlan.GivenName = types.StringNull()
	}

	if !tfPlan.HireDate.IsUnknown() {
		tfPlanHireDate := tfPlan.HireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanHireDate)
		sdkModelUser.SetHireDate(&t)
	} else {
		tfPlan.HireDate = types.StringNull()
	}

	if len(tfPlan.Identities.Elements()) > 0 {
		var tfPlanIdentities []models.ObjectIdentityable
		for _, i := range tfPlan.Identities.Elements() {
			sdkModelIdentities := models.NewObjectIdentity()
			tfModelIdentities := userObjectIdentityModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelIdentities)

			if !tfModelIdentities.Issuer.IsUnknown() {
				tfPlanIssuer := tfModelIdentities.Issuer.ValueString()
				sdkModelIdentities.SetIssuer(&tfPlanIssuer)
			} else {
				tfModelIdentities.Issuer = types.StringNull()
			}

			if !tfModelIdentities.IssuerAssignedId.IsUnknown() {
				tfPlanIssuerAssignedId := tfModelIdentities.IssuerAssignedId.ValueString()
				sdkModelIdentities.SetIssuerAssignedId(&tfPlanIssuerAssignedId)
			} else {
				tfModelIdentities.IssuerAssignedId = types.StringNull()
			}

			if !tfModelIdentities.SignInType.IsUnknown() {
				tfPlanSignInType := tfModelIdentities.SignInType.ValueString()
				sdkModelIdentities.SetSignInType(&tfPlanSignInType)
			} else {
				tfModelIdentities.SignInType = types.StringNull()
			}
		}
		sdkModelUser.SetIdentities(tfPlanIdentities)
	} else {
		tfPlan.Identities = types.ListNull(tfPlan.Identities.ElementType(ctx))
	}

	if len(tfPlan.ImAddresses.Elements()) > 0 {
		var imAddresses []string
		for _, i := range tfPlan.ImAddresses.Elements() {
			imAddresses = append(imAddresses, i.String())
		}
		sdkModelUser.SetImAddresses(imAddresses)
	} else {
		tfPlan.ImAddresses = types.ListNull(types.StringType)
	}

	if len(tfPlan.Interests.Elements()) > 0 {
		var interests []string
		for _, i := range tfPlan.Interests.Elements() {
			interests = append(interests, i.String())
		}
		sdkModelUser.SetInterests(interests)
	} else {
		tfPlan.Interests = types.ListNull(types.StringType)
	}

	if !tfPlan.IsManagementRestricted.IsUnknown() {
		tfPlanIsManagementRestricted := tfPlan.IsManagementRestricted.ValueBool()
		sdkModelUser.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	} else {
		tfPlan.IsManagementRestricted = types.BoolNull()
	}

	if !tfPlan.IsResourceAccount.IsUnknown() {
		tfPlanIsResourceAccount := tfPlan.IsResourceAccount.ValueBool()
		sdkModelUser.SetIsResourceAccount(&tfPlanIsResourceAccount)
	} else {
		tfPlan.IsResourceAccount = types.BoolNull()
	}

	if !tfPlan.JobTitle.IsUnknown() {
		tfPlanJobTitle := tfPlan.JobTitle.ValueString()
		sdkModelUser.SetJobTitle(&tfPlanJobTitle)
	} else {
		tfPlan.JobTitle = types.StringNull()
	}

	if !tfPlan.LastPasswordChangeDateTime.IsUnknown() {
		tfPlanLastPasswordChangeDateTime := tfPlan.LastPasswordChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanLastPasswordChangeDateTime)
		sdkModelUser.SetLastPasswordChangeDateTime(&t)
	} else {
		tfPlan.LastPasswordChangeDateTime = types.StringNull()
	}

	if !tfPlan.LegalAgeGroupClassification.IsUnknown() {
		tfPlanLegalAgeGroupClassification := tfPlan.LegalAgeGroupClassification.ValueString()
		sdkModelUser.SetLegalAgeGroupClassification(&tfPlanLegalAgeGroupClassification)
	} else {
		tfPlan.LegalAgeGroupClassification = types.StringNull()
	}

	if len(tfPlan.LicenseAssignmentStates.Elements()) > 0 {
		var tfPlanLicenseAssignmentStates []models.LicenseAssignmentStateable
		for _, i := range tfPlan.LicenseAssignmentStates.Elements() {
			sdkModelLicenseAssignmentStates := models.NewLicenseAssignmentState()
			tfModelLicenseAssignmentStates := userLicenseAssignmentStateModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelLicenseAssignmentStates)

			if !tfModelLicenseAssignmentStates.AssignedByGroup.IsUnknown() {
				tfPlanAssignedByGroup := tfModelLicenseAssignmentStates.AssignedByGroup.ValueString()
				sdkModelLicenseAssignmentStates.SetAssignedByGroup(&tfPlanAssignedByGroup)
			} else {
				tfModelLicenseAssignmentStates.AssignedByGroup = types.StringNull()
			}

			if len(tfModelLicenseAssignmentStates.DisabledPlans.Elements()) > 0 {
				var DisabledPlans []uuid.UUID
				for _, i := range tfModelLicenseAssignmentStates.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				sdkModelLicenseAssignmentStates.SetDisabledPlans(DisabledPlans)
			} else {
				tfModelLicenseAssignmentStates.DisabledPlans = types.ListNull(types.StringType)
			}

			if !tfModelLicenseAssignmentStates.Error.IsUnknown() {
				tfPlanError := tfModelLicenseAssignmentStates.Error.ValueString()
				sdkModelLicenseAssignmentStates.SetError(&tfPlanError)
			} else {
				tfModelLicenseAssignmentStates.Error = types.StringNull()
			}

			if !tfModelLicenseAssignmentStates.LastUpdatedDateTime.IsUnknown() {
				tfPlanLastUpdatedDateTime := tfModelLicenseAssignmentStates.LastUpdatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanLastUpdatedDateTime)
				sdkModelLicenseAssignmentStates.SetLastUpdatedDateTime(&t)
			} else {
				tfModelLicenseAssignmentStates.LastUpdatedDateTime = types.StringNull()
			}

			if !tfModelLicenseAssignmentStates.SkuId.IsUnknown() {
				tfPlanSkuId := tfModelLicenseAssignmentStates.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				sdkModelLicenseAssignmentStates.SetSkuId(&u)
			} else {
				tfModelLicenseAssignmentStates.SkuId = types.StringNull()
			}

			if !tfModelLicenseAssignmentStates.State.IsUnknown() {
				tfPlanState := tfModelLicenseAssignmentStates.State.ValueString()
				sdkModelLicenseAssignmentStates.SetState(&tfPlanState)
			} else {
				tfModelLicenseAssignmentStates.State = types.StringNull()
			}
		}
		sdkModelUser.SetLicenseAssignmentStates(tfPlanLicenseAssignmentStates)
	} else {
		tfPlan.LicenseAssignmentStates = types.ListNull(tfPlan.LicenseAssignmentStates.ElementType(ctx))
	}

	if !tfPlan.Mail.IsUnknown() {
		tfPlanMail := tfPlan.Mail.ValueString()
		sdkModelUser.SetMail(&tfPlanMail)
	} else {
		tfPlan.Mail = types.StringNull()
	}

	if !tfPlan.MailNickname.IsUnknown() {
		tfPlanMailNickname := tfPlan.MailNickname.ValueString()
		sdkModelUser.SetMailNickname(&tfPlanMailNickname)
	} else {
		tfPlan.MailNickname = types.StringNull()
	}

	if !tfPlan.MobilePhone.IsUnknown() {
		tfPlanMobilePhone := tfPlan.MobilePhone.ValueString()
		sdkModelUser.SetMobilePhone(&tfPlanMobilePhone)
	} else {
		tfPlan.MobilePhone = types.StringNull()
	}

	if !tfPlan.MySite.IsUnknown() {
		tfPlanMySite := tfPlan.MySite.ValueString()
		sdkModelUser.SetMySite(&tfPlanMySite)
	} else {
		tfPlan.MySite = types.StringNull()
	}

	if !tfPlan.OfficeLocation.IsUnknown() {
		tfPlanOfficeLocation := tfPlan.OfficeLocation.ValueString()
		sdkModelUser.SetOfficeLocation(&tfPlanOfficeLocation)
	} else {
		tfPlan.OfficeLocation = types.StringNull()
	}

	if !tfPlan.OnPremisesDistinguishedName.IsUnknown() {
		tfPlanOnPremisesDistinguishedName := tfPlan.OnPremisesDistinguishedName.ValueString()
		sdkModelUser.SetOnPremisesDistinguishedName(&tfPlanOnPremisesDistinguishedName)
	} else {
		tfPlan.OnPremisesDistinguishedName = types.StringNull()
	}

	if !tfPlan.OnPremisesDomainName.IsUnknown() {
		tfPlanOnPremisesDomainName := tfPlan.OnPremisesDomainName.ValueString()
		sdkModelUser.SetOnPremisesDomainName(&tfPlanOnPremisesDomainName)
	} else {
		tfPlan.OnPremisesDomainName = types.StringNull()
	}

	if !tfPlan.OnPremisesExtensionAttributes.IsUnknown() {
		sdkModelOnPremisesExtensionAttributes := models.NewOnPremisesExtensionAttributes()
		tfModelOnPremisesExtensionAttributes := userOnPremisesExtensionAttributesModel{}
		tfPlan.OnPremisesExtensionAttributes.As(ctx, &tfModelOnPremisesExtensionAttributes, basetypes.ObjectAsOptions{})

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute1.IsUnknown() {
			tfPlanExtensionAttribute1 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute1.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute1(&tfPlanExtensionAttribute1)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute10.IsUnknown() {
			tfPlanExtensionAttribute10 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute10.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute10(&tfPlanExtensionAttribute10)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute11.IsUnknown() {
			tfPlanExtensionAttribute11 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute11.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute11(&tfPlanExtensionAttribute11)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute12.IsUnknown() {
			tfPlanExtensionAttribute12 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute12.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute12(&tfPlanExtensionAttribute12)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute13.IsUnknown() {
			tfPlanExtensionAttribute13 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute13.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute13(&tfPlanExtensionAttribute13)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute14.IsUnknown() {
			tfPlanExtensionAttribute14 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute14.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute14(&tfPlanExtensionAttribute14)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute15.IsUnknown() {
			tfPlanExtensionAttribute15 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute15.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute15(&tfPlanExtensionAttribute15)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute2.IsUnknown() {
			tfPlanExtensionAttribute2 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute2.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute2(&tfPlanExtensionAttribute2)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute3.IsUnknown() {
			tfPlanExtensionAttribute3 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute3.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute3(&tfPlanExtensionAttribute3)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute4.IsUnknown() {
			tfPlanExtensionAttribute4 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute4.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute4(&tfPlanExtensionAttribute4)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute5.IsUnknown() {
			tfPlanExtensionAttribute5 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute5.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute5(&tfPlanExtensionAttribute5)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute6.IsUnknown() {
			tfPlanExtensionAttribute6 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute6.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute6(&tfPlanExtensionAttribute6)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute7.IsUnknown() {
			tfPlanExtensionAttribute7 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute7.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute7(&tfPlanExtensionAttribute7)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute8.IsUnknown() {
			tfPlanExtensionAttribute8 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute8.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute8(&tfPlanExtensionAttribute8)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringNull()
		}

		if !tfModelOnPremisesExtensionAttributes.ExtensionAttribute9.IsUnknown() {
			tfPlanExtensionAttribute9 := tfModelOnPremisesExtensionAttributes.ExtensionAttribute9.ValueString()
			sdkModelOnPremisesExtensionAttributes.SetExtensionAttribute9(&tfPlanExtensionAttribute9)
		} else {
			tfModelOnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringNull()
		}
		sdkModelUser.SetOnPremisesExtensionAttributes(sdkModelOnPremisesExtensionAttributes)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelOnPremisesExtensionAttributes.AttributeTypes(), sdkModelOnPremisesExtensionAttributes)
		tfPlan.OnPremisesExtensionAttributes = objectValue
	} else {
		tfPlan.OnPremisesExtensionAttributes = types.ObjectNull(tfPlan.OnPremisesExtensionAttributes.AttributeTypes(ctx))
	}

	if !tfPlan.OnPremisesImmutableId.IsUnknown() {
		tfPlanOnPremisesImmutableId := tfPlan.OnPremisesImmutableId.ValueString()
		sdkModelUser.SetOnPremisesImmutableId(&tfPlanOnPremisesImmutableId)
	} else {
		tfPlan.OnPremisesImmutableId = types.StringNull()
	}

	if !tfPlan.OnPremisesLastSyncDateTime.IsUnknown() {
		tfPlanOnPremisesLastSyncDateTime := tfPlan.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		sdkModelUser.SetOnPremisesLastSyncDateTime(&t)
	} else {
		tfPlan.OnPremisesLastSyncDateTime = types.StringNull()
	}

	if len(tfPlan.OnPremisesProvisioningErrors.Elements()) > 0 {
		var tfPlanOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for _, i := range tfPlan.OnPremisesProvisioningErrors.Elements() {
			sdkModelOnPremisesProvisioningErrors := models.NewOnPremisesProvisioningError()
			tfModelOnPremisesProvisioningErrors := userOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelOnPremisesProvisioningErrors)

			if !tfModelOnPremisesProvisioningErrors.Category.IsUnknown() {
				tfPlanCategory := tfModelOnPremisesProvisioningErrors.Category.ValueString()
				sdkModelOnPremisesProvisioningErrors.SetCategory(&tfPlanCategory)
			} else {
				tfModelOnPremisesProvisioningErrors.Category = types.StringNull()
			}

			if !tfModelOnPremisesProvisioningErrors.OccurredDateTime.IsUnknown() {
				tfPlanOccurredDateTime := tfModelOnPremisesProvisioningErrors.OccurredDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanOccurredDateTime)
				sdkModelOnPremisesProvisioningErrors.SetOccurredDateTime(&t)
			} else {
				tfModelOnPremisesProvisioningErrors.OccurredDateTime = types.StringNull()
			}

			if !tfModelOnPremisesProvisioningErrors.PropertyCausingError.IsUnknown() {
				tfPlanPropertyCausingError := tfModelOnPremisesProvisioningErrors.PropertyCausingError.ValueString()
				sdkModelOnPremisesProvisioningErrors.SetPropertyCausingError(&tfPlanPropertyCausingError)
			} else {
				tfModelOnPremisesProvisioningErrors.PropertyCausingError = types.StringNull()
			}

			if !tfModelOnPremisesProvisioningErrors.Value.IsUnknown() {
				tfPlanValue := tfModelOnPremisesProvisioningErrors.Value.ValueString()
				sdkModelOnPremisesProvisioningErrors.SetValue(&tfPlanValue)
			} else {
				tfModelOnPremisesProvisioningErrors.Value = types.StringNull()
			}
		}
		sdkModelUser.SetOnPremisesProvisioningErrors(tfPlanOnPremisesProvisioningErrors)
	} else {
		tfPlan.OnPremisesProvisioningErrors = types.ListNull(tfPlan.OnPremisesProvisioningErrors.ElementType(ctx))
	}

	if !tfPlan.OnPremisesSamAccountName.IsUnknown() {
		tfPlanOnPremisesSamAccountName := tfPlan.OnPremisesSamAccountName.ValueString()
		sdkModelUser.SetOnPremisesSamAccountName(&tfPlanOnPremisesSamAccountName)
	} else {
		tfPlan.OnPremisesSamAccountName = types.StringNull()
	}

	if !tfPlan.OnPremisesSecurityIdentifier.IsUnknown() {
		tfPlanOnPremisesSecurityIdentifier := tfPlan.OnPremisesSecurityIdentifier.ValueString()
		sdkModelUser.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	} else {
		tfPlan.OnPremisesSecurityIdentifier = types.StringNull()
	}

	if !tfPlan.OnPremisesSyncEnabled.IsUnknown() {
		tfPlanOnPremisesSyncEnabled := tfPlan.OnPremisesSyncEnabled.ValueBool()
		sdkModelUser.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	} else {
		tfPlan.OnPremisesSyncEnabled = types.BoolNull()
	}

	if !tfPlan.OnPremisesUserPrincipalName.IsUnknown() {
		tfPlanOnPremisesUserPrincipalName := tfPlan.OnPremisesUserPrincipalName.ValueString()
		sdkModelUser.SetOnPremisesUserPrincipalName(&tfPlanOnPremisesUserPrincipalName)
	} else {
		tfPlan.OnPremisesUserPrincipalName = types.StringNull()
	}

	if len(tfPlan.OtherMails.Elements()) > 0 {
		var otherMails []string
		for _, i := range tfPlan.OtherMails.Elements() {
			otherMails = append(otherMails, i.String())
		}
		sdkModelUser.SetOtherMails(otherMails)
	} else {
		tfPlan.OtherMails = types.ListNull(types.StringType)
	}

	if !tfPlan.PasswordPolicies.IsUnknown() {
		tfPlanPasswordPolicies := tfPlan.PasswordPolicies.ValueString()
		sdkModelUser.SetPasswordPolicies(&tfPlanPasswordPolicies)
	} else {
		tfPlan.PasswordPolicies = types.StringNull()
	}

	if !tfPlan.PasswordProfile.IsUnknown() {
		sdkModelPasswordProfile := models.NewPasswordProfile()
		tfModelPasswordProfile := userPasswordProfileModel{}
		tfPlan.PasswordProfile.As(ctx, &tfModelPasswordProfile, basetypes.ObjectAsOptions{})

		if !tfModelPasswordProfile.ForceChangePasswordNextSignIn.IsUnknown() {
			tfPlanForceChangePasswordNextSignIn := tfModelPasswordProfile.ForceChangePasswordNextSignIn.ValueBool()
			sdkModelPasswordProfile.SetForceChangePasswordNextSignIn(&tfPlanForceChangePasswordNextSignIn)
		} else {
			tfModelPasswordProfile.ForceChangePasswordNextSignIn = types.BoolNull()
		}

		if !tfModelPasswordProfile.ForceChangePasswordNextSignInWithMfa.IsUnknown() {
			tfPlanForceChangePasswordNextSignInWithMfa := tfModelPasswordProfile.ForceChangePasswordNextSignInWithMfa.ValueBool()
			sdkModelPasswordProfile.SetForceChangePasswordNextSignInWithMfa(&tfPlanForceChangePasswordNextSignInWithMfa)
		} else {
			tfModelPasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolNull()
		}

		if !tfModelPasswordProfile.Password.IsUnknown() {
			tfPlanPassword := tfModelPasswordProfile.Password.ValueString()
			sdkModelPasswordProfile.SetPassword(&tfPlanPassword)
		} else {
			tfModelPasswordProfile.Password = types.StringNull()
		}
		sdkModelUser.SetPasswordProfile(sdkModelPasswordProfile)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelPasswordProfile.AttributeTypes(), sdkModelPasswordProfile)
		tfPlan.PasswordProfile = objectValue
	} else {
		tfPlan.PasswordProfile = types.ObjectNull(tfPlan.PasswordProfile.AttributeTypes(ctx))
	}

	if len(tfPlan.PastProjects.Elements()) > 0 {
		var pastProjects []string
		for _, i := range tfPlan.PastProjects.Elements() {
			pastProjects = append(pastProjects, i.String())
		}
		sdkModelUser.SetPastProjects(pastProjects)
	} else {
		tfPlan.PastProjects = types.ListNull(types.StringType)
	}

	if !tfPlan.PostalCode.IsUnknown() {
		tfPlanPostalCode := tfPlan.PostalCode.ValueString()
		sdkModelUser.SetPostalCode(&tfPlanPostalCode)
	} else {
		tfPlan.PostalCode = types.StringNull()
	}

	if !tfPlan.PreferredDataLocation.IsUnknown() {
		tfPlanPreferredDataLocation := tfPlan.PreferredDataLocation.ValueString()
		sdkModelUser.SetPreferredDataLocation(&tfPlanPreferredDataLocation)
	} else {
		tfPlan.PreferredDataLocation = types.StringNull()
	}

	if !tfPlan.PreferredLanguage.IsUnknown() {
		tfPlanPreferredLanguage := tfPlan.PreferredLanguage.ValueString()
		sdkModelUser.SetPreferredLanguage(&tfPlanPreferredLanguage)
	} else {
		tfPlan.PreferredLanguage = types.StringNull()
	}

	if !tfPlan.PreferredName.IsUnknown() {
		tfPlanPreferredName := tfPlan.PreferredName.ValueString()
		sdkModelUser.SetPreferredName(&tfPlanPreferredName)
	} else {
		tfPlan.PreferredName = types.StringNull()
	}

	if len(tfPlan.ProvisionedPlans.Elements()) > 0 {
		var tfPlanProvisionedPlans []models.ProvisionedPlanable
		for _, i := range tfPlan.ProvisionedPlans.Elements() {
			sdkModelProvisionedPlans := models.NewProvisionedPlan()
			tfModelProvisionedPlans := userProvisionedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelProvisionedPlans)

			if !tfModelProvisionedPlans.CapabilityStatus.IsUnknown() {
				tfPlanCapabilityStatus := tfModelProvisionedPlans.CapabilityStatus.ValueString()
				sdkModelProvisionedPlans.SetCapabilityStatus(&tfPlanCapabilityStatus)
			} else {
				tfModelProvisionedPlans.CapabilityStatus = types.StringNull()
			}

			if !tfModelProvisionedPlans.ProvisioningStatus.IsUnknown() {
				tfPlanProvisioningStatus := tfModelProvisionedPlans.ProvisioningStatus.ValueString()
				sdkModelProvisionedPlans.SetProvisioningStatus(&tfPlanProvisioningStatus)
			} else {
				tfModelProvisionedPlans.ProvisioningStatus = types.StringNull()
			}

			if !tfModelProvisionedPlans.Service.IsUnknown() {
				tfPlanService := tfModelProvisionedPlans.Service.ValueString()
				sdkModelProvisionedPlans.SetService(&tfPlanService)
			} else {
				tfModelProvisionedPlans.Service = types.StringNull()
			}
		}
		sdkModelUser.SetProvisionedPlans(tfPlanProvisionedPlans)
	} else {
		tfPlan.ProvisionedPlans = types.ListNull(tfPlan.ProvisionedPlans.ElementType(ctx))
	}

	if len(tfPlan.ProxyAddresses.Elements()) > 0 {
		var proxyAddresses []string
		for _, i := range tfPlan.ProxyAddresses.Elements() {
			proxyAddresses = append(proxyAddresses, i.String())
		}
		sdkModelUser.SetProxyAddresses(proxyAddresses)
	} else {
		tfPlan.ProxyAddresses = types.ListNull(types.StringType)
	}

	if len(tfPlan.Responsibilities.Elements()) > 0 {
		var responsibilities []string
		for _, i := range tfPlan.Responsibilities.Elements() {
			responsibilities = append(responsibilities, i.String())
		}
		sdkModelUser.SetResponsibilities(responsibilities)
	} else {
		tfPlan.Responsibilities = types.ListNull(types.StringType)
	}

	if len(tfPlan.Schools.Elements()) > 0 {
		var schools []string
		for _, i := range tfPlan.Schools.Elements() {
			schools = append(schools, i.String())
		}
		sdkModelUser.SetSchools(schools)
	} else {
		tfPlan.Schools = types.ListNull(types.StringType)
	}

	if !tfPlan.SecurityIdentifier.IsUnknown() {
		tfPlanSecurityIdentifier := tfPlan.SecurityIdentifier.ValueString()
		sdkModelUser.SetSecurityIdentifier(&tfPlanSecurityIdentifier)
	} else {
		tfPlan.SecurityIdentifier = types.StringNull()
	}

	if len(tfPlan.ServiceProvisioningErrors.Elements()) > 0 {
		var tfPlanServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for _, i := range tfPlan.ServiceProvisioningErrors.Elements() {
			sdkModelServiceProvisioningErrors := models.NewServiceProvisioningError()
			tfModelServiceProvisioningErrors := userServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfModelServiceProvisioningErrors)

			if !tfModelServiceProvisioningErrors.CreatedDateTime.IsUnknown() {
				tfPlanCreatedDateTime := tfModelServiceProvisioningErrors.CreatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
				sdkModelServiceProvisioningErrors.SetCreatedDateTime(&t)
			} else {
				tfModelServiceProvisioningErrors.CreatedDateTime = types.StringNull()
			}

			if !tfModelServiceProvisioningErrors.IsResolved.IsUnknown() {
				tfPlanIsResolved := tfModelServiceProvisioningErrors.IsResolved.ValueBool()
				sdkModelServiceProvisioningErrors.SetIsResolved(&tfPlanIsResolved)
			} else {
				tfModelServiceProvisioningErrors.IsResolved = types.BoolNull()
			}

			if !tfModelServiceProvisioningErrors.ServiceInstance.IsUnknown() {
				tfPlanServiceInstance := tfModelServiceProvisioningErrors.ServiceInstance.ValueString()
				sdkModelServiceProvisioningErrors.SetServiceInstance(&tfPlanServiceInstance)
			} else {
				tfModelServiceProvisioningErrors.ServiceInstance = types.StringNull()
			}
		}
		sdkModelUser.SetServiceProvisioningErrors(tfPlanServiceProvisioningErrors)
	} else {
		tfPlan.ServiceProvisioningErrors = types.ListNull(tfPlan.ServiceProvisioningErrors.ElementType(ctx))
	}

	if !tfPlan.ShowInAddressList.IsUnknown() {
		tfPlanShowInAddressList := tfPlan.ShowInAddressList.ValueBool()
		sdkModelUser.SetShowInAddressList(&tfPlanShowInAddressList)
	} else {
		tfPlan.ShowInAddressList = types.BoolNull()
	}

	if !tfPlan.SignInActivity.IsUnknown() {
		sdkModelSignInActivity := models.NewSignInActivity()
		tfModelSignInActivity := userSignInActivityModel{}
		tfPlan.SignInActivity.As(ctx, &tfModelSignInActivity, basetypes.ObjectAsOptions{})

		if !tfModelSignInActivity.LastNonInteractiveSignInDateTime.IsUnknown() {
			tfPlanLastNonInteractiveSignInDateTime := tfModelSignInActivity.LastNonInteractiveSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastNonInteractiveSignInDateTime)
			sdkModelSignInActivity.SetLastNonInteractiveSignInDateTime(&t)
		} else {
			tfModelSignInActivity.LastNonInteractiveSignInDateTime = types.StringNull()
		}

		if !tfModelSignInActivity.LastNonInteractiveSignInRequestId.IsUnknown() {
			tfPlanLastNonInteractiveSignInRequestId := tfModelSignInActivity.LastNonInteractiveSignInRequestId.ValueString()
			sdkModelSignInActivity.SetLastNonInteractiveSignInRequestId(&tfPlanLastNonInteractiveSignInRequestId)
		} else {
			tfModelSignInActivity.LastNonInteractiveSignInRequestId = types.StringNull()
		}

		if !tfModelSignInActivity.LastSignInDateTime.IsUnknown() {
			tfPlanLastSignInDateTime := tfModelSignInActivity.LastSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastSignInDateTime)
			sdkModelSignInActivity.SetLastSignInDateTime(&t)
		} else {
			tfModelSignInActivity.LastSignInDateTime = types.StringNull()
		}

		if !tfModelSignInActivity.LastSignInRequestId.IsUnknown() {
			tfPlanLastSignInRequestId := tfModelSignInActivity.LastSignInRequestId.ValueString()
			sdkModelSignInActivity.SetLastSignInRequestId(&tfPlanLastSignInRequestId)
		} else {
			tfModelSignInActivity.LastSignInRequestId = types.StringNull()
		}

		if !tfModelSignInActivity.LastSuccessfulSignInDateTime.IsUnknown() {
			tfPlanLastSuccessfulSignInDateTime := tfModelSignInActivity.LastSuccessfulSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastSuccessfulSignInDateTime)
			sdkModelSignInActivity.SetLastSuccessfulSignInDateTime(&t)
		} else {
			tfModelSignInActivity.LastSuccessfulSignInDateTime = types.StringNull()
		}

		if !tfModelSignInActivity.LastSuccessfulSignInRequestId.IsUnknown() {
			tfPlanLastSuccessfulSignInRequestId := tfModelSignInActivity.LastSuccessfulSignInRequestId.ValueString()
			sdkModelSignInActivity.SetLastSuccessfulSignInRequestId(&tfPlanLastSuccessfulSignInRequestId)
		} else {
			tfModelSignInActivity.LastSuccessfulSignInRequestId = types.StringNull()
		}
		sdkModelUser.SetSignInActivity(sdkModelSignInActivity)
		objectValue, _ := types.ObjectValueFrom(ctx, tfModelSignInActivity.AttributeTypes(), sdkModelSignInActivity)
		tfPlan.SignInActivity = objectValue
	} else {
		tfPlan.SignInActivity = types.ObjectNull(tfPlan.SignInActivity.AttributeTypes(ctx))
	}

	if !tfPlan.SignInSessionsValidFromDateTime.IsUnknown() {
		tfPlanSignInSessionsValidFromDateTime := tfPlan.SignInSessionsValidFromDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanSignInSessionsValidFromDateTime)
		sdkModelUser.SetSignInSessionsValidFromDateTime(&t)
	} else {
		tfPlan.SignInSessionsValidFromDateTime = types.StringNull()
	}

	if len(tfPlan.Skills.Elements()) > 0 {
		var skills []string
		for _, i := range tfPlan.Skills.Elements() {
			skills = append(skills, i.String())
		}
		sdkModelUser.SetSkills(skills)
	} else {
		tfPlan.Skills = types.ListNull(types.StringType)
	}

	if !tfPlan.State.IsUnknown() {
		tfPlanState := tfPlan.State.ValueString()
		sdkModelUser.SetState(&tfPlanState)
	} else {
		tfPlan.State = types.StringNull()
	}

	if !tfPlan.StreetAddress.IsUnknown() {
		tfPlanStreetAddress := tfPlan.StreetAddress.ValueString()
		sdkModelUser.SetStreetAddress(&tfPlanStreetAddress)
	} else {
		tfPlan.StreetAddress = types.StringNull()
	}

	if !tfPlan.Surname.IsUnknown() {
		tfPlanSurname := tfPlan.Surname.ValueString()
		sdkModelUser.SetSurname(&tfPlanSurname)
	} else {
		tfPlan.Surname = types.StringNull()
	}

	if !tfPlan.UsageLocation.IsUnknown() {
		tfPlanUsageLocation := tfPlan.UsageLocation.ValueString()
		sdkModelUser.SetUsageLocation(&tfPlanUsageLocation)
	} else {
		tfPlan.UsageLocation = types.StringNull()
	}

	if !tfPlan.UserPrincipalName.IsUnknown() {
		tfPlanUserPrincipalName := tfPlan.UserPrincipalName.ValueString()
		sdkModelUser.SetUserPrincipalName(&tfPlanUserPrincipalName)
	} else {
		tfPlan.UserPrincipalName = types.StringNull()
	}

	if !tfPlan.UserType.IsUnknown() {
		tfPlanUserType := tfPlan.UserType.ValueString()
		sdkModelUser.SetUserType(&tfPlanUserType)
	} else {
		tfPlan.UserType = types.StringNull()
	}

	// Create new user
	result, err := r.client.Users().Post(context.Background(), sdkModelUser, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
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
				"aboutMe",
				"accountEnabled",
				"ageGroup",
				"assignedLicenses",
				"assignedPlans",
				"authorizationInfo",
				"birthday",
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
				"hireDate",
				"identities",
				"imAddresses",
				"interests",
				"isManagementRestricted",
				"isResourceAccount",
				"jobTitle",
				"lastPasswordChangeDateTime",
				"legalAgeGroupClassification",
				"licenseAssignmentStates",
				"mail",
				"mailNickname",
				"mobilePhone",
				"mySite",
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
				"pastProjects",
				"postalCode",
				"preferredDataLocation",
				"preferredLanguage",
				"preferredName",
				"provisionedPlans",
				"proxyAddresses",
				"responsibilities",
				"schools",
				"securityIdentifier",
				"serviceProvisioningErrors",
				"showInAddressList",
				"signInActivity",
				"signInSessionsValidFromDateTime",
				"skills",
				"state",
				"streetAddress",
				"surname",
				"usageLocation",
				"userPrincipalName",
				"userType",
			},
		},
	}

	var result models.Userable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.Users().ByUserId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
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
	} else {
		state.Id = types.StringNull()
	}
	if result.GetDeletedDateTime() != nil {
		state.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	} else {
		state.DeletedDateTime = types.StringNull()
	}
	if result.GetAboutMe() != nil {
		state.AboutMe = types.StringValue(*result.GetAboutMe())
	} else {
		state.AboutMe = types.StringNull()
	}
	if result.GetAccountEnabled() != nil {
		state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	} else {
		state.AccountEnabled = types.BoolNull()
	}
	if result.GetAgeGroup() != nil {
		state.AgeGroup = types.StringValue(*result.GetAgeGroup())
	} else {
		state.AgeGroup = types.StringNull()
	}
	if len(result.GetAssignedLicenses()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedLicenses() {
			assignedLicenses := new(userAssignedLicenseModel)

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
			} else {
				assignedLicenses.SkuId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, assignedLicenses.AttributeTypes(), assignedLicenses)
			objectValues = append(objectValues, objectValue)
		}
		state.AssignedLicenses, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(result.GetAssignedPlans()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedPlans() {
			assignedPlans := new(userAssignedPlanModel)

			if v.GetAssignedDateTime() != nil {
				assignedPlans.AssignedDateTime = types.StringValue(v.GetAssignedDateTime().String())
			} else {
				assignedPlans.AssignedDateTime = types.StringNull()
			}
			if v.GetCapabilityStatus() != nil {
				assignedPlans.CapabilityStatus = types.StringValue(*v.GetCapabilityStatus())
			} else {
				assignedPlans.CapabilityStatus = types.StringNull()
			}
			if v.GetService() != nil {
				assignedPlans.Service = types.StringValue(*v.GetService())
			} else {
				assignedPlans.Service = types.StringNull()
			}
			if v.GetServicePlanId() != nil {
				assignedPlans.ServicePlanId = types.StringValue(v.GetServicePlanId().String())
			} else {
				assignedPlans.ServicePlanId = types.StringNull()
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
	} else {
		state.Birthday = types.StringNull()
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
	} else {
		state.City = types.StringNull()
	}
	if result.GetCompanyName() != nil {
		state.CompanyName = types.StringValue(*result.GetCompanyName())
	} else {
		state.CompanyName = types.StringNull()
	}
	if result.GetConsentProvidedForMinor() != nil {
		state.ConsentProvidedForMinor = types.StringValue(*result.GetConsentProvidedForMinor())
	} else {
		state.ConsentProvidedForMinor = types.StringNull()
	}
	if result.GetCountry() != nil {
		state.Country = types.StringValue(*result.GetCountry())
	} else {
		state.Country = types.StringNull()
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		state.CreatedDateTime = types.StringNull()
	}
	if result.GetCreationType() != nil {
		state.CreationType = types.StringValue(*result.GetCreationType())
	} else {
		state.CreationType = types.StringNull()
	}
	if result.GetDepartment() != nil {
		state.Department = types.StringValue(*result.GetDepartment())
	} else {
		state.Department = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		state.DisplayName = types.StringNull()
	}
	if result.GetEmployeeHireDate() != nil {
		state.EmployeeHireDate = types.StringValue(result.GetEmployeeHireDate().String())
	} else {
		state.EmployeeHireDate = types.StringNull()
	}
	if result.GetEmployeeId() != nil {
		state.EmployeeId = types.StringValue(*result.GetEmployeeId())
	} else {
		state.EmployeeId = types.StringNull()
	}
	if result.GetEmployeeLeaveDateTime() != nil {
		state.EmployeeLeaveDateTime = types.StringValue(result.GetEmployeeLeaveDateTime().String())
	} else {
		state.EmployeeLeaveDateTime = types.StringNull()
	}
	if result.GetEmployeeOrgData() != nil {
		employeeOrgData := new(userEmployeeOrgDataModel)

		if result.GetEmployeeOrgData().GetCostCenter() != nil {
			employeeOrgData.CostCenter = types.StringValue(*result.GetEmployeeOrgData().GetCostCenter())
		} else {
			employeeOrgData.CostCenter = types.StringNull()
		}
		if result.GetEmployeeOrgData().GetDivision() != nil {
			employeeOrgData.Division = types.StringValue(*result.GetEmployeeOrgData().GetDivision())
		} else {
			employeeOrgData.Division = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, employeeOrgData.AttributeTypes(), employeeOrgData)
		state.EmployeeOrgData = objectValue
	}
	if result.GetEmployeeType() != nil {
		state.EmployeeType = types.StringValue(*result.GetEmployeeType())
	} else {
		state.EmployeeType = types.StringNull()
	}
	if result.GetExternalUserState() != nil {
		state.ExternalUserState = types.StringValue(*result.GetExternalUserState())
	} else {
		state.ExternalUserState = types.StringNull()
	}
	if result.GetExternalUserStateChangeDateTime() != nil {
		state.ExternalUserStateChangeDateTime = types.StringValue(result.GetExternalUserStateChangeDateTime().String())
	} else {
		state.ExternalUserStateChangeDateTime = types.StringNull()
	}
	if result.GetFaxNumber() != nil {
		state.FaxNumber = types.StringValue(*result.GetFaxNumber())
	} else {
		state.FaxNumber = types.StringNull()
	}
	if result.GetGivenName() != nil {
		state.GivenName = types.StringValue(*result.GetGivenName())
	} else {
		state.GivenName = types.StringNull()
	}
	if result.GetHireDate() != nil {
		state.HireDate = types.StringValue(result.GetHireDate().String())
	} else {
		state.HireDate = types.StringNull()
	}
	if len(result.GetIdentities()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetIdentities() {
			identities := new(userObjectIdentityModel)

			if v.GetIssuer() != nil {
				identities.Issuer = types.StringValue(*v.GetIssuer())
			} else {
				identities.Issuer = types.StringNull()
			}
			if v.GetIssuerAssignedId() != nil {
				identities.IssuerAssignedId = types.StringValue(*v.GetIssuerAssignedId())
			} else {
				identities.IssuerAssignedId = types.StringNull()
			}
			if v.GetSignInType() != nil {
				identities.SignInType = types.StringValue(*v.GetSignInType())
			} else {
				identities.SignInType = types.StringNull()
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
	if result.GetIsManagementRestricted() != nil {
		state.IsManagementRestricted = types.BoolValue(*result.GetIsManagementRestricted())
	} else {
		state.IsManagementRestricted = types.BoolNull()
	}
	if result.GetIsResourceAccount() != nil {
		state.IsResourceAccount = types.BoolValue(*result.GetIsResourceAccount())
	} else {
		state.IsResourceAccount = types.BoolNull()
	}
	if result.GetJobTitle() != nil {
		state.JobTitle = types.StringValue(*result.GetJobTitle())
	} else {
		state.JobTitle = types.StringNull()
	}
	if result.GetLastPasswordChangeDateTime() != nil {
		state.LastPasswordChangeDateTime = types.StringValue(result.GetLastPasswordChangeDateTime().String())
	} else {
		state.LastPasswordChangeDateTime = types.StringNull()
	}
	if result.GetLegalAgeGroupClassification() != nil {
		state.LegalAgeGroupClassification = types.StringValue(*result.GetLegalAgeGroupClassification())
	} else {
		state.LegalAgeGroupClassification = types.StringNull()
	}
	if len(result.GetLicenseAssignmentStates()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetLicenseAssignmentStates() {
			licenseAssignmentStates := new(userLicenseAssignmentStateModel)

			if v.GetAssignedByGroup() != nil {
				licenseAssignmentStates.AssignedByGroup = types.StringValue(*v.GetAssignedByGroup())
			} else {
				licenseAssignmentStates.AssignedByGroup = types.StringNull()
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
			} else {
				licenseAssignmentStates.Error = types.StringNull()
			}
			if v.GetLastUpdatedDateTime() != nil {
				licenseAssignmentStates.LastUpdatedDateTime = types.StringValue(v.GetLastUpdatedDateTime().String())
			} else {
				licenseAssignmentStates.LastUpdatedDateTime = types.StringNull()
			}
			if v.GetSkuId() != nil {
				licenseAssignmentStates.SkuId = types.StringValue(v.GetSkuId().String())
			} else {
				licenseAssignmentStates.SkuId = types.StringNull()
			}
			if v.GetState() != nil {
				licenseAssignmentStates.State = types.StringValue(*v.GetState())
			} else {
				licenseAssignmentStates.State = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, licenseAssignmentStates.AttributeTypes(), licenseAssignmentStates)
			objectValues = append(objectValues, objectValue)
		}
		state.LicenseAssignmentStates, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetMail() != nil {
		state.Mail = types.StringValue(*result.GetMail())
	} else {
		state.Mail = types.StringNull()
	}
	if result.GetMailNickname() != nil {
		state.MailNickname = types.StringValue(*result.GetMailNickname())
	} else {
		state.MailNickname = types.StringNull()
	}
	if result.GetMobilePhone() != nil {
		state.MobilePhone = types.StringValue(*result.GetMobilePhone())
	} else {
		state.MobilePhone = types.StringNull()
	}
	if result.GetMySite() != nil {
		state.MySite = types.StringValue(*result.GetMySite())
	} else {
		state.MySite = types.StringNull()
	}
	if result.GetOfficeLocation() != nil {
		state.OfficeLocation = types.StringValue(*result.GetOfficeLocation())
	} else {
		state.OfficeLocation = types.StringNull()
	}
	if result.GetOnPremisesDistinguishedName() != nil {
		state.OnPremisesDistinguishedName = types.StringValue(*result.GetOnPremisesDistinguishedName())
	} else {
		state.OnPremisesDistinguishedName = types.StringNull()
	}
	if result.GetOnPremisesDomainName() != nil {
		state.OnPremisesDomainName = types.StringValue(*result.GetOnPremisesDomainName())
	} else {
		state.OnPremisesDomainName = types.StringNull()
	}
	if result.GetOnPremisesExtensionAttributes() != nil {
		onPremisesExtensionAttributes := new(userOnPremisesExtensionAttributesModel)

		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute1 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute1())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute1 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute10 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute10())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute10 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute11 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute11())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute11 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute12 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute12())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute12 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute13 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute13())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute13 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute14 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute14())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute14 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute15 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute15())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute15 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute2 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute2())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute2 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute3 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute3())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute3 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute4 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute4())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute4 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute5 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute5())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute5 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute6 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute6())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute6 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute7 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute7())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute7 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute8 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute8())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute8 = types.StringNull()
		}
		if result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9() != nil {
			onPremisesExtensionAttributes.ExtensionAttribute9 = types.StringValue(*result.GetOnPremisesExtensionAttributes().GetExtensionAttribute9())
		} else {
			onPremisesExtensionAttributes.ExtensionAttribute9 = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, onPremisesExtensionAttributes.AttributeTypes(), onPremisesExtensionAttributes)
		state.OnPremisesExtensionAttributes = objectValue
	}
	if result.GetOnPremisesImmutableId() != nil {
		state.OnPremisesImmutableId = types.StringValue(*result.GetOnPremisesImmutableId())
	} else {
		state.OnPremisesImmutableId = types.StringNull()
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	} else {
		state.OnPremisesLastSyncDateTime = types.StringNull()
	}
	if len(result.GetOnPremisesProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetOnPremisesProvisioningErrors() {
			onPremisesProvisioningErrors := new(userOnPremisesProvisioningErrorModel)

			if v.GetCategory() != nil {
				onPremisesProvisioningErrors.Category = types.StringValue(*v.GetCategory())
			} else {
				onPremisesProvisioningErrors.Category = types.StringNull()
			}
			if v.GetOccurredDateTime() != nil {
				onPremisesProvisioningErrors.OccurredDateTime = types.StringValue(v.GetOccurredDateTime().String())
			} else {
				onPremisesProvisioningErrors.OccurredDateTime = types.StringNull()
			}
			if v.GetPropertyCausingError() != nil {
				onPremisesProvisioningErrors.PropertyCausingError = types.StringValue(*v.GetPropertyCausingError())
			} else {
				onPremisesProvisioningErrors.PropertyCausingError = types.StringNull()
			}
			if v.GetValue() != nil {
				onPremisesProvisioningErrors.Value = types.StringValue(*v.GetValue())
			} else {
				onPremisesProvisioningErrors.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, onPremisesProvisioningErrors.AttributeTypes(), onPremisesProvisioningErrors)
			objectValues = append(objectValues, objectValue)
		}
		state.OnPremisesProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetOnPremisesSamAccountName() != nil {
		state.OnPremisesSamAccountName = types.StringValue(*result.GetOnPremisesSamAccountName())
	} else {
		state.OnPremisesSamAccountName = types.StringNull()
	}
	if result.GetOnPremisesSecurityIdentifier() != nil {
		state.OnPremisesSecurityIdentifier = types.StringValue(*result.GetOnPremisesSecurityIdentifier())
	} else {
		state.OnPremisesSecurityIdentifier = types.StringNull()
	}
	if result.GetOnPremisesSyncEnabled() != nil {
		state.OnPremisesSyncEnabled = types.BoolValue(*result.GetOnPremisesSyncEnabled())
	} else {
		state.OnPremisesSyncEnabled = types.BoolNull()
	}
	if result.GetOnPremisesUserPrincipalName() != nil {
		state.OnPremisesUserPrincipalName = types.StringValue(*result.GetOnPremisesUserPrincipalName())
	} else {
		state.OnPremisesUserPrincipalName = types.StringNull()
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
	} else {
		state.PasswordPolicies = types.StringNull()
	}
	if result.GetPasswordProfile() != nil {
		passwordProfile := new(userPasswordProfileModel)

		if result.GetPasswordProfile().GetForceChangePasswordNextSignIn() != nil {
			passwordProfile.ForceChangePasswordNextSignIn = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignIn())
		} else {
			passwordProfile.ForceChangePasswordNextSignIn = types.BoolNull()
		}
		if result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa() != nil {
			passwordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolValue(*result.GetPasswordProfile().GetForceChangePasswordNextSignInWithMfa())
		} else {
			passwordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolNull()
		}
		if result.GetPasswordProfile().GetPassword() != nil {
			passwordProfile.Password = types.StringValue(*result.GetPasswordProfile().GetPassword())
		} else {
			passwordProfile.Password = types.StringNull()
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
	} else {
		state.PostalCode = types.StringNull()
	}
	if result.GetPreferredDataLocation() != nil {
		state.PreferredDataLocation = types.StringValue(*result.GetPreferredDataLocation())
	} else {
		state.PreferredDataLocation = types.StringNull()
	}
	if result.GetPreferredLanguage() != nil {
		state.PreferredLanguage = types.StringValue(*result.GetPreferredLanguage())
	} else {
		state.PreferredLanguage = types.StringNull()
	}
	if result.GetPreferredName() != nil {
		state.PreferredName = types.StringValue(*result.GetPreferredName())
	} else {
		state.PreferredName = types.StringNull()
	}
	if len(result.GetProvisionedPlans()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetProvisionedPlans() {
			provisionedPlans := new(userProvisionedPlanModel)

			if v.GetCapabilityStatus() != nil {
				provisionedPlans.CapabilityStatus = types.StringValue(*v.GetCapabilityStatus())
			} else {
				provisionedPlans.CapabilityStatus = types.StringNull()
			}
			if v.GetProvisioningStatus() != nil {
				provisionedPlans.ProvisioningStatus = types.StringValue(*v.GetProvisioningStatus())
			} else {
				provisionedPlans.ProvisioningStatus = types.StringNull()
			}
			if v.GetService() != nil {
				provisionedPlans.Service = types.StringValue(*v.GetService())
			} else {
				provisionedPlans.Service = types.StringNull()
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
	} else {
		state.SecurityIdentifier = types.StringNull()
	}
	if len(result.GetServiceProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetServiceProvisioningErrors() {
			serviceProvisioningErrors := new(userServiceProvisioningErrorModel)

			if v.GetCreatedDateTime() != nil {
				serviceProvisioningErrors.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			} else {
				serviceProvisioningErrors.CreatedDateTime = types.StringNull()
			}
			if v.GetIsResolved() != nil {
				serviceProvisioningErrors.IsResolved = types.BoolValue(*v.GetIsResolved())
			} else {
				serviceProvisioningErrors.IsResolved = types.BoolNull()
			}
			if v.GetServiceInstance() != nil {
				serviceProvisioningErrors.ServiceInstance = types.StringValue(*v.GetServiceInstance())
			} else {
				serviceProvisioningErrors.ServiceInstance = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, serviceProvisioningErrors.AttributeTypes(), serviceProvisioningErrors)
			objectValues = append(objectValues, objectValue)
		}
		state.ServiceProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetShowInAddressList() != nil {
		state.ShowInAddressList = types.BoolValue(*result.GetShowInAddressList())
	} else {
		state.ShowInAddressList = types.BoolNull()
	}
	if result.GetSignInActivity() != nil {
		signInActivity := new(userSignInActivityModel)

		if result.GetSignInActivity().GetLastNonInteractiveSignInDateTime() != nil {
			signInActivity.LastNonInteractiveSignInDateTime = types.StringValue(result.GetSignInActivity().GetLastNonInteractiveSignInDateTime().String())
		} else {
			signInActivity.LastNonInteractiveSignInDateTime = types.StringNull()
		}
		if result.GetSignInActivity().GetLastNonInteractiveSignInRequestId() != nil {
			signInActivity.LastNonInteractiveSignInRequestId = types.StringValue(*result.GetSignInActivity().GetLastNonInteractiveSignInRequestId())
		} else {
			signInActivity.LastNonInteractiveSignInRequestId = types.StringNull()
		}
		if result.GetSignInActivity().GetLastSignInDateTime() != nil {
			signInActivity.LastSignInDateTime = types.StringValue(result.GetSignInActivity().GetLastSignInDateTime().String())
		} else {
			signInActivity.LastSignInDateTime = types.StringNull()
		}
		if result.GetSignInActivity().GetLastSignInRequestId() != nil {
			signInActivity.LastSignInRequestId = types.StringValue(*result.GetSignInActivity().GetLastSignInRequestId())
		} else {
			signInActivity.LastSignInRequestId = types.StringNull()
		}
		if result.GetSignInActivity().GetLastSuccessfulSignInDateTime() != nil {
			signInActivity.LastSuccessfulSignInDateTime = types.StringValue(result.GetSignInActivity().GetLastSuccessfulSignInDateTime().String())
		} else {
			signInActivity.LastSuccessfulSignInDateTime = types.StringNull()
		}
		if result.GetSignInActivity().GetLastSuccessfulSignInRequestId() != nil {
			signInActivity.LastSuccessfulSignInRequestId = types.StringValue(*result.GetSignInActivity().GetLastSuccessfulSignInRequestId())
		} else {
			signInActivity.LastSuccessfulSignInRequestId = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, signInActivity.AttributeTypes(), signInActivity)
		state.SignInActivity = objectValue
	}
	if result.GetSignInSessionsValidFromDateTime() != nil {
		state.SignInSessionsValidFromDateTime = types.StringValue(result.GetSignInSessionsValidFromDateTime().String())
	} else {
		state.SignInSessionsValidFromDateTime = types.StringNull()
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
	} else {
		state.State = types.StringNull()
	}
	if result.GetStreetAddress() != nil {
		state.StreetAddress = types.StringValue(*result.GetStreetAddress())
	} else {
		state.StreetAddress = types.StringNull()
	}
	if result.GetSurname() != nil {
		state.Surname = types.StringValue(*result.GetSurname())
	} else {
		state.Surname = types.StringNull()
	}
	if result.GetUsageLocation() != nil {
		state.UsageLocation = types.StringValue(*result.GetUsageLocation())
	} else {
		state.UsageLocation = types.StringNull()
	}
	if result.GetUserPrincipalName() != nil {
		state.UserPrincipalName = types.StringValue(*result.GetUserPrincipalName())
	} else {
		state.UserPrincipalName = types.StringNull()
	}
	if result.GetUserType() != nil {
		state.UserType = types.StringValue(*result.GetUserType())
	} else {
		state.UserType = types.StringNull()
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

	if !plan.Id.Equal(state.Id) {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	}

	if !plan.DeletedDateTime.Equal(state.DeletedDateTime) {
		planDeletedDateTime := plan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !plan.AboutMe.Equal(state.AboutMe) {
		planAboutMe := plan.AboutMe.ValueString()
		requestBody.SetAboutMe(&planAboutMe)
	}

	if !plan.AccountEnabled.Equal(state.AccountEnabled) {
		planAccountEnabled := plan.AccountEnabled.ValueBool()
		requestBody.SetAccountEnabled(&planAccountEnabled)
	}

	if !plan.AgeGroup.Equal(state.AgeGroup) {
		planAgeGroup := plan.AgeGroup.ValueString()
		requestBody.SetAgeGroup(&planAgeGroup)
	}

	if !plan.AssignedLicenses.Equal(state.AssignedLicenses) {
		var planAssignedLicenses []models.AssignedLicenseable
		for k, i := range plan.AssignedLicenses.Elements() {
			assignedLicenses := models.NewAssignedLicense()
			assignedLicensesModel := userAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedLicensesModel)
			assignedLicensesState := userAssignedLicenseModel{}
			types.ListValueFrom(ctx, state.AssignedLicenses.Elements()[k].Type(ctx), &assignedLicensesModel)

			if !assignedLicensesModel.DisabledPlans.Equal(assignedLicensesState.DisabledPlans) {
				var DisabledPlans []uuid.UUID
				for _, i := range assignedLicensesModel.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				assignedLicenses.SetDisabledPlans(DisabledPlans)
			}

			if !assignedLicensesModel.SkuId.Equal(assignedLicensesState.SkuId) {
				planSkuId := assignedLicensesModel.SkuId.ValueString()
				u, _ := uuid.Parse(planSkuId)
				assignedLicenses.SetSkuId(&u)
			}
		}
		requestBody.SetAssignedLicenses(planAssignedLicenses)
	}

	if !plan.AssignedPlans.Equal(state.AssignedPlans) {
		var planAssignedPlans []models.AssignedPlanable
		for k, i := range plan.AssignedPlans.Elements() {
			assignedPlans := models.NewAssignedPlan()
			assignedPlansModel := userAssignedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedPlansModel)
			assignedPlansState := userAssignedPlanModel{}
			types.ListValueFrom(ctx, state.AssignedPlans.Elements()[k].Type(ctx), &assignedPlansModel)

			if !assignedPlansModel.AssignedDateTime.Equal(assignedPlansState.AssignedDateTime) {
				planAssignedDateTime := assignedPlansModel.AssignedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planAssignedDateTime)
				assignedPlans.SetAssignedDateTime(&t)
			}

			if !assignedPlansModel.CapabilityStatus.Equal(assignedPlansState.CapabilityStatus) {
				planCapabilityStatus := assignedPlansModel.CapabilityStatus.ValueString()
				assignedPlans.SetCapabilityStatus(&planCapabilityStatus)
			}

			if !assignedPlansModel.Service.Equal(assignedPlansState.Service) {
				planService := assignedPlansModel.Service.ValueString()
				assignedPlans.SetService(&planService)
			}

			if !assignedPlansModel.ServicePlanId.Equal(assignedPlansState.ServicePlanId) {
				planServicePlanId := assignedPlansModel.ServicePlanId.ValueString()
				u, _ := uuid.Parse(planServicePlanId)
				assignedPlans.SetServicePlanId(&u)
			}
		}
		requestBody.SetAssignedPlans(planAssignedPlans)
	}

	if !plan.AuthorizationInfo.Equal(state.AuthorizationInfo) {
		authorizationInfo := models.NewAuthorizationInfo()
		authorizationInfoModel := userAuthorizationInfoModel{}
		plan.AuthorizationInfo.As(ctx, &authorizationInfoModel, basetypes.ObjectAsOptions{})
		authorizationInfoState := userAuthorizationInfoModel{}
		state.AuthorizationInfo.As(ctx, &authorizationInfoState, basetypes.ObjectAsOptions{})

		if !authorizationInfoModel.CertificateUserIds.Equal(authorizationInfoState.CertificateUserIds) {
			var certificateUserIds []string
			for _, i := range authorizationInfoModel.CertificateUserIds.Elements() {
				certificateUserIds = append(certificateUserIds, i.String())
			}
			authorizationInfo.SetCertificateUserIds(certificateUserIds)
		}
		requestBody.SetAuthorizationInfo(authorizationInfo)
		objectValue, _ := types.ObjectValueFrom(ctx, authorizationInfoModel.AttributeTypes(), authorizationInfoModel)
		plan.AuthorizationInfo = objectValue
	}

	if !plan.Birthday.Equal(state.Birthday) {
		planBirthday := plan.Birthday.ValueString()
		t, _ := time.Parse(time.RFC3339, planBirthday)
		requestBody.SetBirthday(&t)
	}

	if !plan.BusinessPhones.Equal(state.BusinessPhones) {
		var businessPhones []string
		for _, i := range plan.BusinessPhones.Elements() {
			businessPhones = append(businessPhones, i.String())
		}
		requestBody.SetBusinessPhones(businessPhones)
	}

	if !plan.City.Equal(state.City) {
		planCity := plan.City.ValueString()
		requestBody.SetCity(&planCity)
	}

	if !plan.CompanyName.Equal(state.CompanyName) {
		planCompanyName := plan.CompanyName.ValueString()
		requestBody.SetCompanyName(&planCompanyName)
	}

	if !plan.ConsentProvidedForMinor.Equal(state.ConsentProvidedForMinor) {
		planConsentProvidedForMinor := plan.ConsentProvidedForMinor.ValueString()
		requestBody.SetConsentProvidedForMinor(&planConsentProvidedForMinor)
	}

	if !plan.Country.Equal(state.Country) {
		planCountry := plan.Country.ValueString()
		requestBody.SetCountry(&planCountry)
	}

	if !plan.CreatedDateTime.Equal(state.CreatedDateTime) {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	}

	if !plan.CreationType.Equal(state.CreationType) {
		planCreationType := plan.CreationType.ValueString()
		requestBody.SetCreationType(&planCreationType)
	}

	if !plan.Department.Equal(state.Department) {
		planDepartment := plan.Department.ValueString()
		requestBody.SetDepartment(&planDepartment)
	}

	if !plan.DisplayName.Equal(state.DisplayName) {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	}

	if !plan.EmployeeHireDate.Equal(state.EmployeeHireDate) {
		planEmployeeHireDate := plan.EmployeeHireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, planEmployeeHireDate)
		requestBody.SetEmployeeHireDate(&t)
	}

	if !plan.EmployeeId.Equal(state.EmployeeId) {
		planEmployeeId := plan.EmployeeId.ValueString()
		requestBody.SetEmployeeId(&planEmployeeId)
	}

	if !plan.EmployeeLeaveDateTime.Equal(state.EmployeeLeaveDateTime) {
		planEmployeeLeaveDateTime := plan.EmployeeLeaveDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planEmployeeLeaveDateTime)
		requestBody.SetEmployeeLeaveDateTime(&t)
	}

	if !plan.EmployeeOrgData.Equal(state.EmployeeOrgData) {
		employeeOrgData := models.NewEmployeeOrgData()
		employeeOrgDataModel := userEmployeeOrgDataModel{}
		plan.EmployeeOrgData.As(ctx, &employeeOrgDataModel, basetypes.ObjectAsOptions{})
		employeeOrgDataState := userEmployeeOrgDataModel{}
		state.EmployeeOrgData.As(ctx, &employeeOrgDataState, basetypes.ObjectAsOptions{})

		if !employeeOrgDataModel.CostCenter.Equal(employeeOrgDataState.CostCenter) {
			planCostCenter := employeeOrgDataModel.CostCenter.ValueString()
			employeeOrgData.SetCostCenter(&planCostCenter)
		}

		if !employeeOrgDataModel.Division.Equal(employeeOrgDataState.Division) {
			planDivision := employeeOrgDataModel.Division.ValueString()
			employeeOrgData.SetDivision(&planDivision)
		}
		requestBody.SetEmployeeOrgData(employeeOrgData)
		objectValue, _ := types.ObjectValueFrom(ctx, employeeOrgDataModel.AttributeTypes(), employeeOrgDataModel)
		plan.EmployeeOrgData = objectValue
	}

	if !plan.EmployeeType.Equal(state.EmployeeType) {
		planEmployeeType := plan.EmployeeType.ValueString()
		requestBody.SetEmployeeType(&planEmployeeType)
	}

	if !plan.ExternalUserState.Equal(state.ExternalUserState) {
		planExternalUserState := plan.ExternalUserState.ValueString()
		requestBody.SetExternalUserState(&planExternalUserState)
	}

	if !plan.ExternalUserStateChangeDateTime.Equal(state.ExternalUserStateChangeDateTime) {
		planExternalUserStateChangeDateTime := plan.ExternalUserStateChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planExternalUserStateChangeDateTime)
		requestBody.SetExternalUserStateChangeDateTime(&t)
	}

	if !plan.FaxNumber.Equal(state.FaxNumber) {
		planFaxNumber := plan.FaxNumber.ValueString()
		requestBody.SetFaxNumber(&planFaxNumber)
	}

	if !plan.GivenName.Equal(state.GivenName) {
		planGivenName := plan.GivenName.ValueString()
		requestBody.SetGivenName(&planGivenName)
	}

	if !plan.HireDate.Equal(state.HireDate) {
		planHireDate := plan.HireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, planHireDate)
		requestBody.SetHireDate(&t)
	}

	if !plan.Identities.Equal(state.Identities) {
		var planIdentities []models.ObjectIdentityable
		for k, i := range plan.Identities.Elements() {
			identities := models.NewObjectIdentity()
			identitiesModel := userObjectIdentityModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &identitiesModel)
			identitiesState := userObjectIdentityModel{}
			types.ListValueFrom(ctx, state.Identities.Elements()[k].Type(ctx), &identitiesModel)

			if !identitiesModel.Issuer.Equal(identitiesState.Issuer) {
				planIssuer := identitiesModel.Issuer.ValueString()
				identities.SetIssuer(&planIssuer)
			}

			if !identitiesModel.IssuerAssignedId.Equal(identitiesState.IssuerAssignedId) {
				planIssuerAssignedId := identitiesModel.IssuerAssignedId.ValueString()
				identities.SetIssuerAssignedId(&planIssuerAssignedId)
			}

			if !identitiesModel.SignInType.Equal(identitiesState.SignInType) {
				planSignInType := identitiesModel.SignInType.ValueString()
				identities.SetSignInType(&planSignInType)
			}
		}
		requestBody.SetIdentities(planIdentities)
	}

	if !plan.ImAddresses.Equal(state.ImAddresses) {
		var imAddresses []string
		for _, i := range plan.ImAddresses.Elements() {
			imAddresses = append(imAddresses, i.String())
		}
		requestBody.SetImAddresses(imAddresses)
	}

	if !plan.Interests.Equal(state.Interests) {
		var interests []string
		for _, i := range plan.Interests.Elements() {
			interests = append(interests, i.String())
		}
		requestBody.SetInterests(interests)
	}

	if !plan.IsManagementRestricted.Equal(state.IsManagementRestricted) {
		planIsManagementRestricted := plan.IsManagementRestricted.ValueBool()
		requestBody.SetIsManagementRestricted(&planIsManagementRestricted)
	}

	if !plan.IsResourceAccount.Equal(state.IsResourceAccount) {
		planIsResourceAccount := plan.IsResourceAccount.ValueBool()
		requestBody.SetIsResourceAccount(&planIsResourceAccount)
	}

	if !plan.JobTitle.Equal(state.JobTitle) {
		planJobTitle := plan.JobTitle.ValueString()
		requestBody.SetJobTitle(&planJobTitle)
	}

	if !plan.LastPasswordChangeDateTime.Equal(state.LastPasswordChangeDateTime) {
		planLastPasswordChangeDateTime := plan.LastPasswordChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planLastPasswordChangeDateTime)
		requestBody.SetLastPasswordChangeDateTime(&t)
	}

	if !plan.LegalAgeGroupClassification.Equal(state.LegalAgeGroupClassification) {
		planLegalAgeGroupClassification := plan.LegalAgeGroupClassification.ValueString()
		requestBody.SetLegalAgeGroupClassification(&planLegalAgeGroupClassification)
	}

	if !plan.LicenseAssignmentStates.Equal(state.LicenseAssignmentStates) {
		var planLicenseAssignmentStates []models.LicenseAssignmentStateable
		for k, i := range plan.LicenseAssignmentStates.Elements() {
			licenseAssignmentStates := models.NewLicenseAssignmentState()
			licenseAssignmentStatesModel := userLicenseAssignmentStateModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &licenseAssignmentStatesModel)
			licenseAssignmentStatesState := userLicenseAssignmentStateModel{}
			types.ListValueFrom(ctx, state.LicenseAssignmentStates.Elements()[k].Type(ctx), &licenseAssignmentStatesModel)

			if !licenseAssignmentStatesModel.AssignedByGroup.Equal(licenseAssignmentStatesState.AssignedByGroup) {
				planAssignedByGroup := licenseAssignmentStatesModel.AssignedByGroup.ValueString()
				licenseAssignmentStates.SetAssignedByGroup(&planAssignedByGroup)
			}

			if !licenseAssignmentStatesModel.DisabledPlans.Equal(licenseAssignmentStatesState.DisabledPlans) {
				var DisabledPlans []uuid.UUID
				for _, i := range licenseAssignmentStatesModel.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				licenseAssignmentStates.SetDisabledPlans(DisabledPlans)
			}

			if !licenseAssignmentStatesModel.Error.Equal(licenseAssignmentStatesState.Error) {
				planError := licenseAssignmentStatesModel.Error.ValueString()
				licenseAssignmentStates.SetError(&planError)
			}

			if !licenseAssignmentStatesModel.LastUpdatedDateTime.Equal(licenseAssignmentStatesState.LastUpdatedDateTime) {
				planLastUpdatedDateTime := licenseAssignmentStatesModel.LastUpdatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planLastUpdatedDateTime)
				licenseAssignmentStates.SetLastUpdatedDateTime(&t)
			}

			if !licenseAssignmentStatesModel.SkuId.Equal(licenseAssignmentStatesState.SkuId) {
				planSkuId := licenseAssignmentStatesModel.SkuId.ValueString()
				u, _ := uuid.Parse(planSkuId)
				licenseAssignmentStates.SetSkuId(&u)
			}

			if !licenseAssignmentStatesModel.State.Equal(licenseAssignmentStatesState.State) {
				planState := licenseAssignmentStatesModel.State.ValueString()
				licenseAssignmentStates.SetState(&planState)
			}
		}
		requestBody.SetLicenseAssignmentStates(planLicenseAssignmentStates)
	}

	if !plan.Mail.Equal(state.Mail) {
		planMail := plan.Mail.ValueString()
		requestBody.SetMail(&planMail)
	}

	if !plan.MailNickname.Equal(state.MailNickname) {
		planMailNickname := plan.MailNickname.ValueString()
		requestBody.SetMailNickname(&planMailNickname)
	}

	if !plan.MobilePhone.Equal(state.MobilePhone) {
		planMobilePhone := plan.MobilePhone.ValueString()
		requestBody.SetMobilePhone(&planMobilePhone)
	}

	if !plan.MySite.Equal(state.MySite) {
		planMySite := plan.MySite.ValueString()
		requestBody.SetMySite(&planMySite)
	}

	if !plan.OfficeLocation.Equal(state.OfficeLocation) {
		planOfficeLocation := plan.OfficeLocation.ValueString()
		requestBody.SetOfficeLocation(&planOfficeLocation)
	}

	if !plan.OnPremisesDistinguishedName.Equal(state.OnPremisesDistinguishedName) {
		planOnPremisesDistinguishedName := plan.OnPremisesDistinguishedName.ValueString()
		requestBody.SetOnPremisesDistinguishedName(&planOnPremisesDistinguishedName)
	}

	if !plan.OnPremisesDomainName.Equal(state.OnPremisesDomainName) {
		planOnPremisesDomainName := plan.OnPremisesDomainName.ValueString()
		requestBody.SetOnPremisesDomainName(&planOnPremisesDomainName)
	}

	if !plan.OnPremisesExtensionAttributes.Equal(state.OnPremisesExtensionAttributes) {
		onPremisesExtensionAttributes := models.NewOnPremisesExtensionAttributes()
		onPremisesExtensionAttributesModel := userOnPremisesExtensionAttributesModel{}
		plan.OnPremisesExtensionAttributes.As(ctx, &onPremisesExtensionAttributesModel, basetypes.ObjectAsOptions{})
		onPremisesExtensionAttributesState := userOnPremisesExtensionAttributesModel{}
		state.OnPremisesExtensionAttributes.As(ctx, &onPremisesExtensionAttributesState, basetypes.ObjectAsOptions{})

		if !onPremisesExtensionAttributesModel.ExtensionAttribute1.Equal(onPremisesExtensionAttributesState.ExtensionAttribute1) {
			planExtensionAttribute1 := onPremisesExtensionAttributesModel.ExtensionAttribute1.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute1(&planExtensionAttribute1)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute10.Equal(onPremisesExtensionAttributesState.ExtensionAttribute10) {
			planExtensionAttribute10 := onPremisesExtensionAttributesModel.ExtensionAttribute10.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute10(&planExtensionAttribute10)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute11.Equal(onPremisesExtensionAttributesState.ExtensionAttribute11) {
			planExtensionAttribute11 := onPremisesExtensionAttributesModel.ExtensionAttribute11.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute11(&planExtensionAttribute11)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute12.Equal(onPremisesExtensionAttributesState.ExtensionAttribute12) {
			planExtensionAttribute12 := onPremisesExtensionAttributesModel.ExtensionAttribute12.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute12(&planExtensionAttribute12)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute13.Equal(onPremisesExtensionAttributesState.ExtensionAttribute13) {
			planExtensionAttribute13 := onPremisesExtensionAttributesModel.ExtensionAttribute13.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute13(&planExtensionAttribute13)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute14.Equal(onPremisesExtensionAttributesState.ExtensionAttribute14) {
			planExtensionAttribute14 := onPremisesExtensionAttributesModel.ExtensionAttribute14.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute14(&planExtensionAttribute14)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute15.Equal(onPremisesExtensionAttributesState.ExtensionAttribute15) {
			planExtensionAttribute15 := onPremisesExtensionAttributesModel.ExtensionAttribute15.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute15(&planExtensionAttribute15)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute2.Equal(onPremisesExtensionAttributesState.ExtensionAttribute2) {
			planExtensionAttribute2 := onPremisesExtensionAttributesModel.ExtensionAttribute2.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute2(&planExtensionAttribute2)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute3.Equal(onPremisesExtensionAttributesState.ExtensionAttribute3) {
			planExtensionAttribute3 := onPremisesExtensionAttributesModel.ExtensionAttribute3.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute3(&planExtensionAttribute3)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute4.Equal(onPremisesExtensionAttributesState.ExtensionAttribute4) {
			planExtensionAttribute4 := onPremisesExtensionAttributesModel.ExtensionAttribute4.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute4(&planExtensionAttribute4)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute5.Equal(onPremisesExtensionAttributesState.ExtensionAttribute5) {
			planExtensionAttribute5 := onPremisesExtensionAttributesModel.ExtensionAttribute5.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute5(&planExtensionAttribute5)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute6.Equal(onPremisesExtensionAttributesState.ExtensionAttribute6) {
			planExtensionAttribute6 := onPremisesExtensionAttributesModel.ExtensionAttribute6.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute6(&planExtensionAttribute6)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute7.Equal(onPremisesExtensionAttributesState.ExtensionAttribute7) {
			planExtensionAttribute7 := onPremisesExtensionAttributesModel.ExtensionAttribute7.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute7(&planExtensionAttribute7)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute8.Equal(onPremisesExtensionAttributesState.ExtensionAttribute8) {
			planExtensionAttribute8 := onPremisesExtensionAttributesModel.ExtensionAttribute8.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute8(&planExtensionAttribute8)
		}

		if !onPremisesExtensionAttributesModel.ExtensionAttribute9.Equal(onPremisesExtensionAttributesState.ExtensionAttribute9) {
			planExtensionAttribute9 := onPremisesExtensionAttributesModel.ExtensionAttribute9.ValueString()
			onPremisesExtensionAttributes.SetExtensionAttribute9(&planExtensionAttribute9)
		}
		requestBody.SetOnPremisesExtensionAttributes(onPremisesExtensionAttributes)
		objectValue, _ := types.ObjectValueFrom(ctx, onPremisesExtensionAttributesModel.AttributeTypes(), onPremisesExtensionAttributesModel)
		plan.OnPremisesExtensionAttributes = objectValue
	}

	if !plan.OnPremisesImmutableId.Equal(state.OnPremisesImmutableId) {
		planOnPremisesImmutableId := plan.OnPremisesImmutableId.ValueString()
		requestBody.SetOnPremisesImmutableId(&planOnPremisesImmutableId)
	}

	if !plan.OnPremisesLastSyncDateTime.Equal(state.OnPremisesLastSyncDateTime) {
		planOnPremisesLastSyncDateTime := plan.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planOnPremisesLastSyncDateTime)
		requestBody.SetOnPremisesLastSyncDateTime(&t)
	}

	if !plan.OnPremisesProvisioningErrors.Equal(state.OnPremisesProvisioningErrors) {
		var planOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for k, i := range plan.OnPremisesProvisioningErrors.Elements() {
			onPremisesProvisioningErrors := models.NewOnPremisesProvisioningError()
			onPremisesProvisioningErrorsModel := userOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &onPremisesProvisioningErrorsModel)
			onPremisesProvisioningErrorsState := userOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, state.OnPremisesProvisioningErrors.Elements()[k].Type(ctx), &onPremisesProvisioningErrorsModel)

			if !onPremisesProvisioningErrorsModel.Category.Equal(onPremisesProvisioningErrorsState.Category) {
				planCategory := onPremisesProvisioningErrorsModel.Category.ValueString()
				onPremisesProvisioningErrors.SetCategory(&planCategory)
			}

			if !onPremisesProvisioningErrorsModel.OccurredDateTime.Equal(onPremisesProvisioningErrorsState.OccurredDateTime) {
				planOccurredDateTime := onPremisesProvisioningErrorsModel.OccurredDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planOccurredDateTime)
				onPremisesProvisioningErrors.SetOccurredDateTime(&t)
			}

			if !onPremisesProvisioningErrorsModel.PropertyCausingError.Equal(onPremisesProvisioningErrorsState.PropertyCausingError) {
				planPropertyCausingError := onPremisesProvisioningErrorsModel.PropertyCausingError.ValueString()
				onPremisesProvisioningErrors.SetPropertyCausingError(&planPropertyCausingError)
			}

			if !onPremisesProvisioningErrorsModel.Value.Equal(onPremisesProvisioningErrorsState.Value) {
				planValue := onPremisesProvisioningErrorsModel.Value.ValueString()
				onPremisesProvisioningErrors.SetValue(&planValue)
			}
		}
		requestBody.SetOnPremisesProvisioningErrors(planOnPremisesProvisioningErrors)
	}

	if !plan.OnPremisesSamAccountName.Equal(state.OnPremisesSamAccountName) {
		planOnPremisesSamAccountName := plan.OnPremisesSamAccountName.ValueString()
		requestBody.SetOnPremisesSamAccountName(&planOnPremisesSamAccountName)
	}

	if !plan.OnPremisesSecurityIdentifier.Equal(state.OnPremisesSecurityIdentifier) {
		planOnPremisesSecurityIdentifier := plan.OnPremisesSecurityIdentifier.ValueString()
		requestBody.SetOnPremisesSecurityIdentifier(&planOnPremisesSecurityIdentifier)
	}

	if !plan.OnPremisesSyncEnabled.Equal(state.OnPremisesSyncEnabled) {
		planOnPremisesSyncEnabled := plan.OnPremisesSyncEnabled.ValueBool()
		requestBody.SetOnPremisesSyncEnabled(&planOnPremisesSyncEnabled)
	}

	if !plan.OnPremisesUserPrincipalName.Equal(state.OnPremisesUserPrincipalName) {
		planOnPremisesUserPrincipalName := plan.OnPremisesUserPrincipalName.ValueString()
		requestBody.SetOnPremisesUserPrincipalName(&planOnPremisesUserPrincipalName)
	}

	if !plan.OtherMails.Equal(state.OtherMails) {
		var otherMails []string
		for _, i := range plan.OtherMails.Elements() {
			otherMails = append(otherMails, i.String())
		}
		requestBody.SetOtherMails(otherMails)
	}

	if !plan.PasswordPolicies.Equal(state.PasswordPolicies) {
		planPasswordPolicies := plan.PasswordPolicies.ValueString()
		requestBody.SetPasswordPolicies(&planPasswordPolicies)
	}

	if !plan.PasswordProfile.Equal(state.PasswordProfile) {
		passwordProfile := models.NewPasswordProfile()
		passwordProfileModel := userPasswordProfileModel{}
		plan.PasswordProfile.As(ctx, &passwordProfileModel, basetypes.ObjectAsOptions{})
		passwordProfileState := userPasswordProfileModel{}
		state.PasswordProfile.As(ctx, &passwordProfileState, basetypes.ObjectAsOptions{})

		if !passwordProfileModel.ForceChangePasswordNextSignIn.Equal(passwordProfileState.ForceChangePasswordNextSignIn) {
			planForceChangePasswordNextSignIn := passwordProfileModel.ForceChangePasswordNextSignIn.ValueBool()
			passwordProfile.SetForceChangePasswordNextSignIn(&planForceChangePasswordNextSignIn)
		}

		if !passwordProfileModel.ForceChangePasswordNextSignInWithMfa.Equal(passwordProfileState.ForceChangePasswordNextSignInWithMfa) {
			planForceChangePasswordNextSignInWithMfa := passwordProfileModel.ForceChangePasswordNextSignInWithMfa.ValueBool()
			passwordProfile.SetForceChangePasswordNextSignInWithMfa(&planForceChangePasswordNextSignInWithMfa)
		}

		if !passwordProfileModel.Password.Equal(passwordProfileState.Password) {
			planPassword := passwordProfileModel.Password.ValueString()
			passwordProfile.SetPassword(&planPassword)
		}
		requestBody.SetPasswordProfile(passwordProfile)
		objectValue, _ := types.ObjectValueFrom(ctx, passwordProfileModel.AttributeTypes(), passwordProfileModel)
		plan.PasswordProfile = objectValue
	}

	if !plan.PastProjects.Equal(state.PastProjects) {
		var pastProjects []string
		for _, i := range plan.PastProjects.Elements() {
			pastProjects = append(pastProjects, i.String())
		}
		requestBody.SetPastProjects(pastProjects)
	}

	if !plan.PostalCode.Equal(state.PostalCode) {
		planPostalCode := plan.PostalCode.ValueString()
		requestBody.SetPostalCode(&planPostalCode)
	}

	if !plan.PreferredDataLocation.Equal(state.PreferredDataLocation) {
		planPreferredDataLocation := plan.PreferredDataLocation.ValueString()
		requestBody.SetPreferredDataLocation(&planPreferredDataLocation)
	}

	if !plan.PreferredLanguage.Equal(state.PreferredLanguage) {
		planPreferredLanguage := plan.PreferredLanguage.ValueString()
		requestBody.SetPreferredLanguage(&planPreferredLanguage)
	}

	if !plan.PreferredName.Equal(state.PreferredName) {
		planPreferredName := plan.PreferredName.ValueString()
		requestBody.SetPreferredName(&planPreferredName)
	}

	if !plan.ProvisionedPlans.Equal(state.ProvisionedPlans) {
		var planProvisionedPlans []models.ProvisionedPlanable
		for k, i := range plan.ProvisionedPlans.Elements() {
			provisionedPlans := models.NewProvisionedPlan()
			provisionedPlansModel := userProvisionedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &provisionedPlansModel)
			provisionedPlansState := userProvisionedPlanModel{}
			types.ListValueFrom(ctx, state.ProvisionedPlans.Elements()[k].Type(ctx), &provisionedPlansModel)

			if !provisionedPlansModel.CapabilityStatus.Equal(provisionedPlansState.CapabilityStatus) {
				planCapabilityStatus := provisionedPlansModel.CapabilityStatus.ValueString()
				provisionedPlans.SetCapabilityStatus(&planCapabilityStatus)
			}

			if !provisionedPlansModel.ProvisioningStatus.Equal(provisionedPlansState.ProvisioningStatus) {
				planProvisioningStatus := provisionedPlansModel.ProvisioningStatus.ValueString()
				provisionedPlans.SetProvisioningStatus(&planProvisioningStatus)
			}

			if !provisionedPlansModel.Service.Equal(provisionedPlansState.Service) {
				planService := provisionedPlansModel.Service.ValueString()
				provisionedPlans.SetService(&planService)
			}
		}
		requestBody.SetProvisionedPlans(planProvisionedPlans)
	}

	if !plan.ProxyAddresses.Equal(state.ProxyAddresses) {
		var proxyAddresses []string
		for _, i := range plan.ProxyAddresses.Elements() {
			proxyAddresses = append(proxyAddresses, i.String())
		}
		requestBody.SetProxyAddresses(proxyAddresses)
	}

	if !plan.Responsibilities.Equal(state.Responsibilities) {
		var responsibilities []string
		for _, i := range plan.Responsibilities.Elements() {
			responsibilities = append(responsibilities, i.String())
		}
		requestBody.SetResponsibilities(responsibilities)
	}

	if !plan.Schools.Equal(state.Schools) {
		var schools []string
		for _, i := range plan.Schools.Elements() {
			schools = append(schools, i.String())
		}
		requestBody.SetSchools(schools)
	}

	if !plan.SecurityIdentifier.Equal(state.SecurityIdentifier) {
		planSecurityIdentifier := plan.SecurityIdentifier.ValueString()
		requestBody.SetSecurityIdentifier(&planSecurityIdentifier)
	}

	if !plan.ServiceProvisioningErrors.Equal(state.ServiceProvisioningErrors) {
		var planServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for k, i := range plan.ServiceProvisioningErrors.Elements() {
			serviceProvisioningErrors := models.NewServiceProvisioningError()
			serviceProvisioningErrorsModel := userServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &serviceProvisioningErrorsModel)
			serviceProvisioningErrorsState := userServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, state.ServiceProvisioningErrors.Elements()[k].Type(ctx), &serviceProvisioningErrorsModel)

			if !serviceProvisioningErrorsModel.CreatedDateTime.Equal(serviceProvisioningErrorsState.CreatedDateTime) {
				planCreatedDateTime := serviceProvisioningErrorsModel.CreatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, planCreatedDateTime)
				serviceProvisioningErrors.SetCreatedDateTime(&t)
			}

			if !serviceProvisioningErrorsModel.IsResolved.Equal(serviceProvisioningErrorsState.IsResolved) {
				planIsResolved := serviceProvisioningErrorsModel.IsResolved.ValueBool()
				serviceProvisioningErrors.SetIsResolved(&planIsResolved)
			}

			if !serviceProvisioningErrorsModel.ServiceInstance.Equal(serviceProvisioningErrorsState.ServiceInstance) {
				planServiceInstance := serviceProvisioningErrorsModel.ServiceInstance.ValueString()
				serviceProvisioningErrors.SetServiceInstance(&planServiceInstance)
			}
		}
		requestBody.SetServiceProvisioningErrors(planServiceProvisioningErrors)
	}

	if !plan.ShowInAddressList.Equal(state.ShowInAddressList) {
		planShowInAddressList := plan.ShowInAddressList.ValueBool()
		requestBody.SetShowInAddressList(&planShowInAddressList)
	}

	if !plan.SignInActivity.Equal(state.SignInActivity) {
		signInActivity := models.NewSignInActivity()
		signInActivityModel := userSignInActivityModel{}
		plan.SignInActivity.As(ctx, &signInActivityModel, basetypes.ObjectAsOptions{})
		signInActivityState := userSignInActivityModel{}
		state.SignInActivity.As(ctx, &signInActivityState, basetypes.ObjectAsOptions{})

		if !signInActivityModel.LastNonInteractiveSignInDateTime.Equal(signInActivityState.LastNonInteractiveSignInDateTime) {
			planLastNonInteractiveSignInDateTime := signInActivityModel.LastNonInteractiveSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, planLastNonInteractiveSignInDateTime)
			signInActivity.SetLastNonInteractiveSignInDateTime(&t)
		}

		if !signInActivityModel.LastNonInteractiveSignInRequestId.Equal(signInActivityState.LastNonInteractiveSignInRequestId) {
			planLastNonInteractiveSignInRequestId := signInActivityModel.LastNonInteractiveSignInRequestId.ValueString()
			signInActivity.SetLastNonInteractiveSignInRequestId(&planLastNonInteractiveSignInRequestId)
		}

		if !signInActivityModel.LastSignInDateTime.Equal(signInActivityState.LastSignInDateTime) {
			planLastSignInDateTime := signInActivityModel.LastSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, planLastSignInDateTime)
			signInActivity.SetLastSignInDateTime(&t)
		}

		if !signInActivityModel.LastSignInRequestId.Equal(signInActivityState.LastSignInRequestId) {
			planLastSignInRequestId := signInActivityModel.LastSignInRequestId.ValueString()
			signInActivity.SetLastSignInRequestId(&planLastSignInRequestId)
		}

		if !signInActivityModel.LastSuccessfulSignInDateTime.Equal(signInActivityState.LastSuccessfulSignInDateTime) {
			planLastSuccessfulSignInDateTime := signInActivityModel.LastSuccessfulSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, planLastSuccessfulSignInDateTime)
			signInActivity.SetLastSuccessfulSignInDateTime(&t)
		}

		if !signInActivityModel.LastSuccessfulSignInRequestId.Equal(signInActivityState.LastSuccessfulSignInRequestId) {
			planLastSuccessfulSignInRequestId := signInActivityModel.LastSuccessfulSignInRequestId.ValueString()
			signInActivity.SetLastSuccessfulSignInRequestId(&planLastSuccessfulSignInRequestId)
		}
		requestBody.SetSignInActivity(signInActivity)
		objectValue, _ := types.ObjectValueFrom(ctx, signInActivityModel.AttributeTypes(), signInActivityModel)
		plan.SignInActivity = objectValue
	}

	if !plan.SignInSessionsValidFromDateTime.Equal(state.SignInSessionsValidFromDateTime) {
		planSignInSessionsValidFromDateTime := plan.SignInSessionsValidFromDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planSignInSessionsValidFromDateTime)
		requestBody.SetSignInSessionsValidFromDateTime(&t)
	}

	if !plan.Skills.Equal(state.Skills) {
		var skills []string
		for _, i := range plan.Skills.Elements() {
			skills = append(skills, i.String())
		}
		requestBody.SetSkills(skills)
	}

	if !plan.State.Equal(state.State) {
		planState := plan.State.ValueString()
		requestBody.SetState(&planState)
	}

	if !plan.StreetAddress.Equal(state.StreetAddress) {
		planStreetAddress := plan.StreetAddress.ValueString()
		requestBody.SetStreetAddress(&planStreetAddress)
	}

	if !plan.Surname.Equal(state.Surname) {
		planSurname := plan.Surname.ValueString()
		requestBody.SetSurname(&planSurname)
	}

	if !plan.UsageLocation.Equal(state.UsageLocation) {
		planUsageLocation := plan.UsageLocation.ValueString()
		requestBody.SetUsageLocation(&planUsageLocation)
	}

	if !plan.UserPrincipalName.Equal(state.UserPrincipalName) {
		planUserPrincipalName := plan.UserPrincipalName.ValueString()
		requestBody.SetUserPrincipalName(&planUserPrincipalName)
	}

	if !plan.UserType.Equal(state.UserType) {
		planUserType := plan.UserType.ValueString()
		requestBody.SetUserType(&planUserType)
	}

	// Update user
	_, err := r.client.Users().ByUserId(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
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
