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
	var tfPlanUser userModel
	diags := req.Plan.Get(ctx, &tfPlanUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	requestBodyUser := models.NewUser()
	// START Id | CreateStringAttribute
	if !tfPlanUser.Id.IsUnknown() {
		tfPlanId := tfPlanUser.Id.ValueString()
		requestBodyUser.SetId(&tfPlanId)
	} else {
		tfPlanUser.Id = types.StringNull()
	}
	// END Id | CreateStringAttribute

	// START DeletedDateTime | CreateStringTimeAttribute
	if !tfPlanUser.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlanUser.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyUser.SetDeletedDateTime(&t)
	} else {
		tfPlanUser.DeletedDateTime = types.StringNull()
	}
	// END DeletedDateTime | CreateStringTimeAttribute

	// START AboutMe | CreateStringAttribute
	if !tfPlanUser.AboutMe.IsUnknown() {
		tfPlanAboutMe := tfPlanUser.AboutMe.ValueString()
		requestBodyUser.SetAboutMe(&tfPlanAboutMe)
	} else {
		tfPlanUser.AboutMe = types.StringNull()
	}
	// END AboutMe | CreateStringAttribute

	// START AccountEnabled | CreateBoolAttribute
	if !tfPlanUser.AccountEnabled.IsUnknown() {
		tfPlanAccountEnabled := tfPlanUser.AccountEnabled.ValueBool()
		requestBodyUser.SetAccountEnabled(&tfPlanAccountEnabled)
	} else {
		tfPlanUser.AccountEnabled = types.BoolNull()
	}
	// END AccountEnabled | CreateBoolAttribute

	// START AgeGroup | CreateStringAttribute
	if !tfPlanUser.AgeGroup.IsUnknown() {
		tfPlanAgeGroup := tfPlanUser.AgeGroup.ValueString()
		requestBodyUser.SetAgeGroup(&tfPlanAgeGroup)
	} else {
		tfPlanUser.AgeGroup = types.StringNull()
	}
	// END AgeGroup | CreateStringAttribute

	// START AssignedLicenses | CreateArrayObjectAttribute
	if len(tfPlanUser.AssignedLicenses.Elements()) > 0 {
		var requestBodyAssignedLicenses []models.AssignedLicenseable
		for _, i := range tfPlanUser.AssignedLicenses.Elements() {
			requestBodyAssignedLicense := models.NewAssignedLicense()
			tfPlanAssignedLicense := userAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLicense)

			// START DisabledPlans | CreateArrayUuidAttribute
			if len(tfPlanAssignedLicense.DisabledPlans.Elements()) > 0 {
				var uuidArrayDisabledPlans []uuid.UUID
				for _, i := range tfPlanAssignedLicense.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					uuidArrayDisabledPlans = append(uuidArrayDisabledPlans, u)
				}
				requestBodyAssignedLicense.SetDisabledPlans(uuidArrayDisabledPlans)
			} else {
				tfPlanAssignedLicense.DisabledPlans = types.ListNull(types.StringType)
			}

			// END DisabledPlans | CreateArrayUuidAttribute

			// START SkuId | CreateStringUuidAttribute
			if !tfPlanAssignedLicense.SkuId.IsUnknown() {
				tfPlanSkuId := tfPlanAssignedLicense.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				requestBodyAssignedLicense.SetSkuId(&u)
			} else {
				tfPlanAssignedLicense.SkuId = types.StringNull()
			}
			// END SkuId | CreateStringUuidAttribute

		}
		requestBodyUser.SetAssignedLicenses(requestBodyAssignedLicenses)
	} else {
		tfPlanUser.AssignedLicenses = types.ListNull(tfPlanUser.AssignedLicenses.ElementType(ctx))
	}
	// END AssignedLicenses | CreateArrayObjectAttribute

	// START AssignedPlans | CreateArrayObjectAttribute
	if len(tfPlanUser.AssignedPlans.Elements()) > 0 {
		var requestBodyAssignedPlans []models.AssignedPlanable
		for _, i := range tfPlanUser.AssignedPlans.Elements() {
			requestBodyAssignedPlan := models.NewAssignedPlan()
			tfPlanAssignedPlan := userAssignedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedPlan)

			// START AssignedDateTime | CreateStringTimeAttribute
			if !tfPlanAssignedPlan.AssignedDateTime.IsUnknown() {
				tfPlanAssignedDateTime := tfPlanAssignedPlan.AssignedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanAssignedDateTime)
				requestBodyAssignedPlan.SetAssignedDateTime(&t)
			} else {
				tfPlanAssignedPlan.AssignedDateTime = types.StringNull()
			}
			// END AssignedDateTime | CreateStringTimeAttribute

			// START CapabilityStatus | CreateStringAttribute
			if !tfPlanAssignedPlan.CapabilityStatus.IsUnknown() {
				tfPlanCapabilityStatus := tfPlanAssignedPlan.CapabilityStatus.ValueString()
				requestBodyAssignedPlan.SetCapabilityStatus(&tfPlanCapabilityStatus)
			} else {
				tfPlanAssignedPlan.CapabilityStatus = types.StringNull()
			}
			// END CapabilityStatus | CreateStringAttribute

			// START Service | CreateStringAttribute
			if !tfPlanAssignedPlan.Service.IsUnknown() {
				tfPlanService := tfPlanAssignedPlan.Service.ValueString()
				requestBodyAssignedPlan.SetService(&tfPlanService)
			} else {
				tfPlanAssignedPlan.Service = types.StringNull()
			}
			// END Service | CreateStringAttribute

			// START ServicePlanId | CreateStringUuidAttribute
			if !tfPlanAssignedPlan.ServicePlanId.IsUnknown() {
				tfPlanServicePlanId := tfPlanAssignedPlan.ServicePlanId.ValueString()
				u, _ := uuid.Parse(tfPlanServicePlanId)
				requestBodyAssignedPlan.SetServicePlanId(&u)
			} else {
				tfPlanAssignedPlan.ServicePlanId = types.StringNull()
			}
			// END ServicePlanId | CreateStringUuidAttribute

		}
		requestBodyUser.SetAssignedPlans(requestBodyAssignedPlans)
	} else {
		tfPlanUser.AssignedPlans = types.ListNull(tfPlanUser.AssignedPlans.ElementType(ctx))
	}
	// END AssignedPlans | CreateArrayObjectAttribute

	// START AuthorizationInfo | CreateObjectAttribute
	if !tfPlanUser.AuthorizationInfo.IsUnknown() {
		requestBodyAuthorizationInfo := models.NewAuthorizationInfo()
		tfPlanAuthorizationInfo := userAuthorizationInfoModel{}
		tfPlanUser.AuthorizationInfo.As(ctx, &tfPlanAuthorizationInfo, basetypes.ObjectAsOptions{})

		// START CertificateUserIds | CreateArrayStringAttribute
		if len(tfPlanAuthorizationInfo.CertificateUserIds.Elements()) > 0 {
			var stringArrayCertificateUserIds []string
			for _, i := range tfPlanAuthorizationInfo.CertificateUserIds.Elements() {
				stringArrayCertificateUserIds = append(stringArrayCertificateUserIds, i.String())
			}
			requestBodyAuthorizationInfo.SetCertificateUserIds(stringArrayCertificateUserIds)
		} else {
			tfPlanAuthorizationInfo.CertificateUserIds = types.ListNull(types.StringType)
		}
		// END CertificateUserIds | CreateArrayStringAttribute

		requestBodyUser.SetAuthorizationInfo(requestBodyAuthorizationInfo)
		tfPlanUser.AuthorizationInfo, _ = types.ObjectValueFrom(ctx, tfPlanAuthorizationInfo.AttributeTypes(), requestBodyAuthorizationInfo)
	} else {
		tfPlanUser.AuthorizationInfo = types.ObjectNull(tfPlanUser.AuthorizationInfo.AttributeTypes(ctx))
	}
	// END AuthorizationInfo | CreateObjectAttribute

	// START Birthday | CreateStringTimeAttribute
	if !tfPlanUser.Birthday.IsUnknown() {
		tfPlanBirthday := tfPlanUser.Birthday.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanBirthday)
		requestBodyUser.SetBirthday(&t)
	} else {
		tfPlanUser.Birthday = types.StringNull()
	}
	// END Birthday | CreateStringTimeAttribute

	// START BusinessPhones | CreateArrayStringAttribute
	if len(tfPlanUser.BusinessPhones.Elements()) > 0 {
		var stringArrayBusinessPhones []string
		for _, i := range tfPlanUser.BusinessPhones.Elements() {
			stringArrayBusinessPhones = append(stringArrayBusinessPhones, i.String())
		}
		requestBodyUser.SetBusinessPhones(stringArrayBusinessPhones)
	} else {
		tfPlanUser.BusinessPhones = types.ListNull(types.StringType)
	}
	// END BusinessPhones | CreateArrayStringAttribute

	// START City | CreateStringAttribute
	if !tfPlanUser.City.IsUnknown() {
		tfPlanCity := tfPlanUser.City.ValueString()
		requestBodyUser.SetCity(&tfPlanCity)
	} else {
		tfPlanUser.City = types.StringNull()
	}
	// END City | CreateStringAttribute

	// START CompanyName | CreateStringAttribute
	if !tfPlanUser.CompanyName.IsUnknown() {
		tfPlanCompanyName := tfPlanUser.CompanyName.ValueString()
		requestBodyUser.SetCompanyName(&tfPlanCompanyName)
	} else {
		tfPlanUser.CompanyName = types.StringNull()
	}
	// END CompanyName | CreateStringAttribute

	// START ConsentProvidedForMinor | CreateStringAttribute
	if !tfPlanUser.ConsentProvidedForMinor.IsUnknown() {
		tfPlanConsentProvidedForMinor := tfPlanUser.ConsentProvidedForMinor.ValueString()
		requestBodyUser.SetConsentProvidedForMinor(&tfPlanConsentProvidedForMinor)
	} else {
		tfPlanUser.ConsentProvidedForMinor = types.StringNull()
	}
	// END ConsentProvidedForMinor | CreateStringAttribute

	// START Country | CreateStringAttribute
	if !tfPlanUser.Country.IsUnknown() {
		tfPlanCountry := tfPlanUser.Country.ValueString()
		requestBodyUser.SetCountry(&tfPlanCountry)
	} else {
		tfPlanUser.Country = types.StringNull()
	}
	// END Country | CreateStringAttribute

	// START CreatedDateTime | CreateStringTimeAttribute
	if !tfPlanUser.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlanUser.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyUser.SetCreatedDateTime(&t)
	} else {
		tfPlanUser.CreatedDateTime = types.StringNull()
	}
	// END CreatedDateTime | CreateStringTimeAttribute

	// START CreationType | CreateStringAttribute
	if !tfPlanUser.CreationType.IsUnknown() {
		tfPlanCreationType := tfPlanUser.CreationType.ValueString()
		requestBodyUser.SetCreationType(&tfPlanCreationType)
	} else {
		tfPlanUser.CreationType = types.StringNull()
	}
	// END CreationType | CreateStringAttribute

	// START Department | CreateStringAttribute
	if !tfPlanUser.Department.IsUnknown() {
		tfPlanDepartment := tfPlanUser.Department.ValueString()
		requestBodyUser.SetDepartment(&tfPlanDepartment)
	} else {
		tfPlanUser.Department = types.StringNull()
	}
	// END Department | CreateStringAttribute

	// START DisplayName | CreateStringAttribute
	if !tfPlanUser.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanUser.DisplayName.ValueString()
		requestBodyUser.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanUser.DisplayName = types.StringNull()
	}
	// END DisplayName | CreateStringAttribute

	// START EmployeeHireDate | CreateStringTimeAttribute
	if !tfPlanUser.EmployeeHireDate.IsUnknown() {
		tfPlanEmployeeHireDate := tfPlanUser.EmployeeHireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanEmployeeHireDate)
		requestBodyUser.SetEmployeeHireDate(&t)
	} else {
		tfPlanUser.EmployeeHireDate = types.StringNull()
	}
	// END EmployeeHireDate | CreateStringTimeAttribute

	// START EmployeeId | CreateStringAttribute
	if !tfPlanUser.EmployeeId.IsUnknown() {
		tfPlanEmployeeId := tfPlanUser.EmployeeId.ValueString()
		requestBodyUser.SetEmployeeId(&tfPlanEmployeeId)
	} else {
		tfPlanUser.EmployeeId = types.StringNull()
	}
	// END EmployeeId | CreateStringAttribute

	// START EmployeeLeaveDateTime | CreateStringTimeAttribute
	if !tfPlanUser.EmployeeLeaveDateTime.IsUnknown() {
		tfPlanEmployeeLeaveDateTime := tfPlanUser.EmployeeLeaveDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanEmployeeLeaveDateTime)
		requestBodyUser.SetEmployeeLeaveDateTime(&t)
	} else {
		tfPlanUser.EmployeeLeaveDateTime = types.StringNull()
	}
	// END EmployeeLeaveDateTime | CreateStringTimeAttribute

	// START EmployeeOrgData | CreateObjectAttribute
	if !tfPlanUser.EmployeeOrgData.IsUnknown() {
		requestBodyEmployeeOrgData := models.NewEmployeeOrgData()
		tfPlanEmployeeOrgData := userEmployeeOrgDataModel{}
		tfPlanUser.EmployeeOrgData.As(ctx, &tfPlanEmployeeOrgData, basetypes.ObjectAsOptions{})

		// START CostCenter | CreateStringAttribute
		if !tfPlanEmployeeOrgData.CostCenter.IsUnknown() {
			tfPlanCostCenter := tfPlanEmployeeOrgData.CostCenter.ValueString()
			requestBodyEmployeeOrgData.SetCostCenter(&tfPlanCostCenter)
		} else {
			tfPlanEmployeeOrgData.CostCenter = types.StringNull()
		}
		// END CostCenter | CreateStringAttribute

		// START Division | CreateStringAttribute
		if !tfPlanEmployeeOrgData.Division.IsUnknown() {
			tfPlanDivision := tfPlanEmployeeOrgData.Division.ValueString()
			requestBodyEmployeeOrgData.SetDivision(&tfPlanDivision)
		} else {
			tfPlanEmployeeOrgData.Division = types.StringNull()
		}
		// END Division | CreateStringAttribute

		requestBodyUser.SetEmployeeOrgData(requestBodyEmployeeOrgData)
		tfPlanUser.EmployeeOrgData, _ = types.ObjectValueFrom(ctx, tfPlanEmployeeOrgData.AttributeTypes(), requestBodyEmployeeOrgData)
	} else {
		tfPlanUser.EmployeeOrgData = types.ObjectNull(tfPlanUser.EmployeeOrgData.AttributeTypes(ctx))
	}
	// END EmployeeOrgData | CreateObjectAttribute

	// START EmployeeType | CreateStringAttribute
	if !tfPlanUser.EmployeeType.IsUnknown() {
		tfPlanEmployeeType := tfPlanUser.EmployeeType.ValueString()
		requestBodyUser.SetEmployeeType(&tfPlanEmployeeType)
	} else {
		tfPlanUser.EmployeeType = types.StringNull()
	}
	// END EmployeeType | CreateStringAttribute

	// START ExternalUserState | CreateStringAttribute
	if !tfPlanUser.ExternalUserState.IsUnknown() {
		tfPlanExternalUserState := tfPlanUser.ExternalUserState.ValueString()
		requestBodyUser.SetExternalUserState(&tfPlanExternalUserState)
	} else {
		tfPlanUser.ExternalUserState = types.StringNull()
	}
	// END ExternalUserState | CreateStringAttribute

	// START ExternalUserStateChangeDateTime | CreateStringTimeAttribute
	if !tfPlanUser.ExternalUserStateChangeDateTime.IsUnknown() {
		tfPlanExternalUserStateChangeDateTime := tfPlanUser.ExternalUserStateChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanExternalUserStateChangeDateTime)
		requestBodyUser.SetExternalUserStateChangeDateTime(&t)
	} else {
		tfPlanUser.ExternalUserStateChangeDateTime = types.StringNull()
	}
	// END ExternalUserStateChangeDateTime | CreateStringTimeAttribute

	// START FaxNumber | CreateStringAttribute
	if !tfPlanUser.FaxNumber.IsUnknown() {
		tfPlanFaxNumber := tfPlanUser.FaxNumber.ValueString()
		requestBodyUser.SetFaxNumber(&tfPlanFaxNumber)
	} else {
		tfPlanUser.FaxNumber = types.StringNull()
	}
	// END FaxNumber | CreateStringAttribute

	// START GivenName | CreateStringAttribute
	if !tfPlanUser.GivenName.IsUnknown() {
		tfPlanGivenName := tfPlanUser.GivenName.ValueString()
		requestBodyUser.SetGivenName(&tfPlanGivenName)
	} else {
		tfPlanUser.GivenName = types.StringNull()
	}
	// END GivenName | CreateStringAttribute

	// START HireDate | CreateStringTimeAttribute
	if !tfPlanUser.HireDate.IsUnknown() {
		tfPlanHireDate := tfPlanUser.HireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanHireDate)
		requestBodyUser.SetHireDate(&t)
	} else {
		tfPlanUser.HireDate = types.StringNull()
	}
	// END HireDate | CreateStringTimeAttribute

	// START Identities | CreateArrayObjectAttribute
	if len(tfPlanUser.Identities.Elements()) > 0 {
		var requestBodyIdentities []models.ObjectIdentityable
		for _, i := range tfPlanUser.Identities.Elements() {
			requestBodyObjectIdentity := models.NewObjectIdentity()
			tfPlanObjectIdentity := userObjectIdentityModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanObjectIdentity)

			// START Issuer | CreateStringAttribute
			if !tfPlanObjectIdentity.Issuer.IsUnknown() {
				tfPlanIssuer := tfPlanObjectIdentity.Issuer.ValueString()
				requestBodyObjectIdentity.SetIssuer(&tfPlanIssuer)
			} else {
				tfPlanObjectIdentity.Issuer = types.StringNull()
			}
			// END Issuer | CreateStringAttribute

			// START IssuerAssignedId | CreateStringAttribute
			if !tfPlanObjectIdentity.IssuerAssignedId.IsUnknown() {
				tfPlanIssuerAssignedId := tfPlanObjectIdentity.IssuerAssignedId.ValueString()
				requestBodyObjectIdentity.SetIssuerAssignedId(&tfPlanIssuerAssignedId)
			} else {
				tfPlanObjectIdentity.IssuerAssignedId = types.StringNull()
			}
			// END IssuerAssignedId | CreateStringAttribute

			// START SignInType | CreateStringAttribute
			if !tfPlanObjectIdentity.SignInType.IsUnknown() {
				tfPlanSignInType := tfPlanObjectIdentity.SignInType.ValueString()
				requestBodyObjectIdentity.SetSignInType(&tfPlanSignInType)
			} else {
				tfPlanObjectIdentity.SignInType = types.StringNull()
			}
			// END SignInType | CreateStringAttribute

		}
		requestBodyUser.SetIdentities(requestBodyIdentities)
	} else {
		tfPlanUser.Identities = types.ListNull(tfPlanUser.Identities.ElementType(ctx))
	}
	// END Identities | CreateArrayObjectAttribute

	// START ImAddresses | CreateArrayStringAttribute
	if len(tfPlanUser.ImAddresses.Elements()) > 0 {
		var stringArrayImAddresses []string
		for _, i := range tfPlanUser.ImAddresses.Elements() {
			stringArrayImAddresses = append(stringArrayImAddresses, i.String())
		}
		requestBodyUser.SetImAddresses(stringArrayImAddresses)
	} else {
		tfPlanUser.ImAddresses = types.ListNull(types.StringType)
	}
	// END ImAddresses | CreateArrayStringAttribute

	// START Interests | CreateArrayStringAttribute
	if len(tfPlanUser.Interests.Elements()) > 0 {
		var stringArrayInterests []string
		for _, i := range tfPlanUser.Interests.Elements() {
			stringArrayInterests = append(stringArrayInterests, i.String())
		}
		requestBodyUser.SetInterests(stringArrayInterests)
	} else {
		tfPlanUser.Interests = types.ListNull(types.StringType)
	}
	// END Interests | CreateArrayStringAttribute

	// START IsManagementRestricted | CreateBoolAttribute
	if !tfPlanUser.IsManagementRestricted.IsUnknown() {
		tfPlanIsManagementRestricted := tfPlanUser.IsManagementRestricted.ValueBool()
		requestBodyUser.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	} else {
		tfPlanUser.IsManagementRestricted = types.BoolNull()
	}
	// END IsManagementRestricted | CreateBoolAttribute

	// START IsResourceAccount | CreateBoolAttribute
	if !tfPlanUser.IsResourceAccount.IsUnknown() {
		tfPlanIsResourceAccount := tfPlanUser.IsResourceAccount.ValueBool()
		requestBodyUser.SetIsResourceAccount(&tfPlanIsResourceAccount)
	} else {
		tfPlanUser.IsResourceAccount = types.BoolNull()
	}
	// END IsResourceAccount | CreateBoolAttribute

	// START JobTitle | CreateStringAttribute
	if !tfPlanUser.JobTitle.IsUnknown() {
		tfPlanJobTitle := tfPlanUser.JobTitle.ValueString()
		requestBodyUser.SetJobTitle(&tfPlanJobTitle)
	} else {
		tfPlanUser.JobTitle = types.StringNull()
	}
	// END JobTitle | CreateStringAttribute

	// START LastPasswordChangeDateTime | CreateStringTimeAttribute
	if !tfPlanUser.LastPasswordChangeDateTime.IsUnknown() {
		tfPlanLastPasswordChangeDateTime := tfPlanUser.LastPasswordChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanLastPasswordChangeDateTime)
		requestBodyUser.SetLastPasswordChangeDateTime(&t)
	} else {
		tfPlanUser.LastPasswordChangeDateTime = types.StringNull()
	}
	// END LastPasswordChangeDateTime | CreateStringTimeAttribute

	// START LegalAgeGroupClassification | CreateStringAttribute
	if !tfPlanUser.LegalAgeGroupClassification.IsUnknown() {
		tfPlanLegalAgeGroupClassification := tfPlanUser.LegalAgeGroupClassification.ValueString()
		requestBodyUser.SetLegalAgeGroupClassification(&tfPlanLegalAgeGroupClassification)
	} else {
		tfPlanUser.LegalAgeGroupClassification = types.StringNull()
	}
	// END LegalAgeGroupClassification | CreateStringAttribute

	// START LicenseAssignmentStates | CreateArrayObjectAttribute
	if len(tfPlanUser.LicenseAssignmentStates.Elements()) > 0 {
		var requestBodyLicenseAssignmentStates []models.LicenseAssignmentStateable
		for _, i := range tfPlanUser.LicenseAssignmentStates.Elements() {
			requestBodyLicenseAssignmentState := models.NewLicenseAssignmentState()
			tfPlanLicenseAssignmentState := userLicenseAssignmentStateModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanLicenseAssignmentState)

			// START AssignedByGroup | CreateStringAttribute
			if !tfPlanLicenseAssignmentState.AssignedByGroup.IsUnknown() {
				tfPlanAssignedByGroup := tfPlanLicenseAssignmentState.AssignedByGroup.ValueString()
				requestBodyLicenseAssignmentState.SetAssignedByGroup(&tfPlanAssignedByGroup)
			} else {
				tfPlanLicenseAssignmentState.AssignedByGroup = types.StringNull()
			}
			// END AssignedByGroup | CreateStringAttribute

			// START DisabledPlans | CreateArrayUuidAttribute
			if len(tfPlanLicenseAssignmentState.DisabledPlans.Elements()) > 0 {
				var uuidArrayDisabledPlans []uuid.UUID
				for _, i := range tfPlanLicenseAssignmentState.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					uuidArrayDisabledPlans = append(uuidArrayDisabledPlans, u)
				}
				requestBodyLicenseAssignmentState.SetDisabledPlans(uuidArrayDisabledPlans)
			} else {
				tfPlanLicenseAssignmentState.DisabledPlans = types.ListNull(types.StringType)
			}

			// END DisabledPlans | CreateArrayUuidAttribute

			// START Error | CreateStringAttribute
			if !tfPlanLicenseAssignmentState.Error.IsUnknown() {
				tfPlanError := tfPlanLicenseAssignmentState.Error.ValueString()
				requestBodyLicenseAssignmentState.SetError(&tfPlanError)
			} else {
				tfPlanLicenseAssignmentState.Error = types.StringNull()
			}
			// END Error | CreateStringAttribute

			// START LastUpdatedDateTime | CreateStringTimeAttribute
			if !tfPlanLicenseAssignmentState.LastUpdatedDateTime.IsUnknown() {
				tfPlanLastUpdatedDateTime := tfPlanLicenseAssignmentState.LastUpdatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanLastUpdatedDateTime)
				requestBodyLicenseAssignmentState.SetLastUpdatedDateTime(&t)
			} else {
				tfPlanLicenseAssignmentState.LastUpdatedDateTime = types.StringNull()
			}
			// END LastUpdatedDateTime | CreateStringTimeAttribute

			// START SkuId | CreateStringUuidAttribute
			if !tfPlanLicenseAssignmentState.SkuId.IsUnknown() {
				tfPlanSkuId := tfPlanLicenseAssignmentState.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				requestBodyLicenseAssignmentState.SetSkuId(&u)
			} else {
				tfPlanLicenseAssignmentState.SkuId = types.StringNull()
			}
			// END SkuId | CreateStringUuidAttribute

			// START State | CreateStringAttribute
			if !tfPlanLicenseAssignmentState.State.IsUnknown() {
				tfPlanState := tfPlanLicenseAssignmentState.State.ValueString()
				requestBodyLicenseAssignmentState.SetState(&tfPlanState)
			} else {
				tfPlanLicenseAssignmentState.State = types.StringNull()
			}
			// END State | CreateStringAttribute

		}
		requestBodyUser.SetLicenseAssignmentStates(requestBodyLicenseAssignmentStates)
	} else {
		tfPlanUser.LicenseAssignmentStates = types.ListNull(tfPlanUser.LicenseAssignmentStates.ElementType(ctx))
	}
	// END LicenseAssignmentStates | CreateArrayObjectAttribute

	// START Mail | CreateStringAttribute
	if !tfPlanUser.Mail.IsUnknown() {
		tfPlanMail := tfPlanUser.Mail.ValueString()
		requestBodyUser.SetMail(&tfPlanMail)
	} else {
		tfPlanUser.Mail = types.StringNull()
	}
	// END Mail | CreateStringAttribute

	// START MailNickname | CreateStringAttribute
	if !tfPlanUser.MailNickname.IsUnknown() {
		tfPlanMailNickname := tfPlanUser.MailNickname.ValueString()
		requestBodyUser.SetMailNickname(&tfPlanMailNickname)
	} else {
		tfPlanUser.MailNickname = types.StringNull()
	}
	// END MailNickname | CreateStringAttribute

	// START MobilePhone | CreateStringAttribute
	if !tfPlanUser.MobilePhone.IsUnknown() {
		tfPlanMobilePhone := tfPlanUser.MobilePhone.ValueString()
		requestBodyUser.SetMobilePhone(&tfPlanMobilePhone)
	} else {
		tfPlanUser.MobilePhone = types.StringNull()
	}
	// END MobilePhone | CreateStringAttribute

	// START MySite | CreateStringAttribute
	if !tfPlanUser.MySite.IsUnknown() {
		tfPlanMySite := tfPlanUser.MySite.ValueString()
		requestBodyUser.SetMySite(&tfPlanMySite)
	} else {
		tfPlanUser.MySite = types.StringNull()
	}
	// END MySite | CreateStringAttribute

	// START OfficeLocation | CreateStringAttribute
	if !tfPlanUser.OfficeLocation.IsUnknown() {
		tfPlanOfficeLocation := tfPlanUser.OfficeLocation.ValueString()
		requestBodyUser.SetOfficeLocation(&tfPlanOfficeLocation)
	} else {
		tfPlanUser.OfficeLocation = types.StringNull()
	}
	// END OfficeLocation | CreateStringAttribute

	// START OnPremisesDistinguishedName | CreateStringAttribute
	if !tfPlanUser.OnPremisesDistinguishedName.IsUnknown() {
		tfPlanOnPremisesDistinguishedName := tfPlanUser.OnPremisesDistinguishedName.ValueString()
		requestBodyUser.SetOnPremisesDistinguishedName(&tfPlanOnPremisesDistinguishedName)
	} else {
		tfPlanUser.OnPremisesDistinguishedName = types.StringNull()
	}
	// END OnPremisesDistinguishedName | CreateStringAttribute

	// START OnPremisesDomainName | CreateStringAttribute
	if !tfPlanUser.OnPremisesDomainName.IsUnknown() {
		tfPlanOnPremisesDomainName := tfPlanUser.OnPremisesDomainName.ValueString()
		requestBodyUser.SetOnPremisesDomainName(&tfPlanOnPremisesDomainName)
	} else {
		tfPlanUser.OnPremisesDomainName = types.StringNull()
	}
	// END OnPremisesDomainName | CreateStringAttribute

	// START OnPremisesExtensionAttributes | CreateObjectAttribute
	if !tfPlanUser.OnPremisesExtensionAttributes.IsUnknown() {
		requestBodyOnPremisesExtensionAttributes := models.NewOnPremisesExtensionAttributes()
		tfPlanOnPremisesExtensionAttributes := userOnPremisesExtensionAttributesModel{}
		tfPlanUser.OnPremisesExtensionAttributes.As(ctx, &tfPlanOnPremisesExtensionAttributes, basetypes.ObjectAsOptions{})

		// START ExtensionAttribute1 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute1.IsUnknown() {
			tfPlanExtensionAttribute1 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute1.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute1(&tfPlanExtensionAttribute1)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute1 = types.StringNull()
		}
		// END ExtensionAttribute1 | CreateStringAttribute

		// START ExtensionAttribute10 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute10.IsUnknown() {
			tfPlanExtensionAttribute10 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute10.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute10(&tfPlanExtensionAttribute10)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute10 = types.StringNull()
		}
		// END ExtensionAttribute10 | CreateStringAttribute

		// START ExtensionAttribute11 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute11.IsUnknown() {
			tfPlanExtensionAttribute11 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute11.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute11(&tfPlanExtensionAttribute11)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute11 = types.StringNull()
		}
		// END ExtensionAttribute11 | CreateStringAttribute

		// START ExtensionAttribute12 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute12.IsUnknown() {
			tfPlanExtensionAttribute12 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute12.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute12(&tfPlanExtensionAttribute12)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute12 = types.StringNull()
		}
		// END ExtensionAttribute12 | CreateStringAttribute

		// START ExtensionAttribute13 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute13.IsUnknown() {
			tfPlanExtensionAttribute13 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute13.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute13(&tfPlanExtensionAttribute13)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute13 = types.StringNull()
		}
		// END ExtensionAttribute13 | CreateStringAttribute

		// START ExtensionAttribute14 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute14.IsUnknown() {
			tfPlanExtensionAttribute14 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute14.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute14(&tfPlanExtensionAttribute14)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute14 = types.StringNull()
		}
		// END ExtensionAttribute14 | CreateStringAttribute

		// START ExtensionAttribute15 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute15.IsUnknown() {
			tfPlanExtensionAttribute15 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute15.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute15(&tfPlanExtensionAttribute15)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute15 = types.StringNull()
		}
		// END ExtensionAttribute15 | CreateStringAttribute

		// START ExtensionAttribute2 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute2.IsUnknown() {
			tfPlanExtensionAttribute2 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute2.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute2(&tfPlanExtensionAttribute2)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute2 = types.StringNull()
		}
		// END ExtensionAttribute2 | CreateStringAttribute

		// START ExtensionAttribute3 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute3.IsUnknown() {
			tfPlanExtensionAttribute3 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute3.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute3(&tfPlanExtensionAttribute3)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute3 = types.StringNull()
		}
		// END ExtensionAttribute3 | CreateStringAttribute

		// START ExtensionAttribute4 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute4.IsUnknown() {
			tfPlanExtensionAttribute4 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute4.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute4(&tfPlanExtensionAttribute4)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute4 = types.StringNull()
		}
		// END ExtensionAttribute4 | CreateStringAttribute

		// START ExtensionAttribute5 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute5.IsUnknown() {
			tfPlanExtensionAttribute5 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute5.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute5(&tfPlanExtensionAttribute5)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute5 = types.StringNull()
		}
		// END ExtensionAttribute5 | CreateStringAttribute

		// START ExtensionAttribute6 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute6.IsUnknown() {
			tfPlanExtensionAttribute6 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute6.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute6(&tfPlanExtensionAttribute6)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute6 = types.StringNull()
		}
		// END ExtensionAttribute6 | CreateStringAttribute

		// START ExtensionAttribute7 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute7.IsUnknown() {
			tfPlanExtensionAttribute7 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute7.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute7(&tfPlanExtensionAttribute7)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute7 = types.StringNull()
		}
		// END ExtensionAttribute7 | CreateStringAttribute

		// START ExtensionAttribute8 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute8.IsUnknown() {
			tfPlanExtensionAttribute8 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute8.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute8(&tfPlanExtensionAttribute8)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute8 = types.StringNull()
		}
		// END ExtensionAttribute8 | CreateStringAttribute

		// START ExtensionAttribute9 | CreateStringAttribute
		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute9.IsUnknown() {
			tfPlanExtensionAttribute9 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute9.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute9(&tfPlanExtensionAttribute9)
		} else {
			tfPlanOnPremisesExtensionAttributes.ExtensionAttribute9 = types.StringNull()
		}
		// END ExtensionAttribute9 | CreateStringAttribute

		requestBodyUser.SetOnPremisesExtensionAttributes(requestBodyOnPremisesExtensionAttributes)
		tfPlanUser.OnPremisesExtensionAttributes, _ = types.ObjectValueFrom(ctx, tfPlanOnPremisesExtensionAttributes.AttributeTypes(), requestBodyOnPremisesExtensionAttributes)
	} else {
		tfPlanUser.OnPremisesExtensionAttributes = types.ObjectNull(tfPlanUser.OnPremisesExtensionAttributes.AttributeTypes(ctx))
	}
	// END OnPremisesExtensionAttributes | CreateObjectAttribute

	// START OnPremisesImmutableId | CreateStringAttribute
	if !tfPlanUser.OnPremisesImmutableId.IsUnknown() {
		tfPlanOnPremisesImmutableId := tfPlanUser.OnPremisesImmutableId.ValueString()
		requestBodyUser.SetOnPremisesImmutableId(&tfPlanOnPremisesImmutableId)
	} else {
		tfPlanUser.OnPremisesImmutableId = types.StringNull()
	}
	// END OnPremisesImmutableId | CreateStringAttribute

	// START OnPremisesLastSyncDateTime | CreateStringTimeAttribute
	if !tfPlanUser.OnPremisesLastSyncDateTime.IsUnknown() {
		tfPlanOnPremisesLastSyncDateTime := tfPlanUser.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		requestBodyUser.SetOnPremisesLastSyncDateTime(&t)
	} else {
		tfPlanUser.OnPremisesLastSyncDateTime = types.StringNull()
	}
	// END OnPremisesLastSyncDateTime | CreateStringTimeAttribute

	// START OnPremisesProvisioningErrors | CreateArrayObjectAttribute
	if len(tfPlanUser.OnPremisesProvisioningErrors.Elements()) > 0 {
		var requestBodyOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for _, i := range tfPlanUser.OnPremisesProvisioningErrors.Elements() {
			requestBodyOnPremisesProvisioningError := models.NewOnPremisesProvisioningError()
			tfPlanOnPremisesProvisioningError := userOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOnPremisesProvisioningError)

			// START Category | CreateStringAttribute
			if !tfPlanOnPremisesProvisioningError.Category.IsUnknown() {
				tfPlanCategory := tfPlanOnPremisesProvisioningError.Category.ValueString()
				requestBodyOnPremisesProvisioningError.SetCategory(&tfPlanCategory)
			} else {
				tfPlanOnPremisesProvisioningError.Category = types.StringNull()
			}
			// END Category | CreateStringAttribute

			// START OccurredDateTime | CreateStringTimeAttribute
			if !tfPlanOnPremisesProvisioningError.OccurredDateTime.IsUnknown() {
				tfPlanOccurredDateTime := tfPlanOnPremisesProvisioningError.OccurredDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanOccurredDateTime)
				requestBodyOnPremisesProvisioningError.SetOccurredDateTime(&t)
			} else {
				tfPlanOnPremisesProvisioningError.OccurredDateTime = types.StringNull()
			}
			// END OccurredDateTime | CreateStringTimeAttribute

			// START PropertyCausingError | CreateStringAttribute
			if !tfPlanOnPremisesProvisioningError.PropertyCausingError.IsUnknown() {
				tfPlanPropertyCausingError := tfPlanOnPremisesProvisioningError.PropertyCausingError.ValueString()
				requestBodyOnPremisesProvisioningError.SetPropertyCausingError(&tfPlanPropertyCausingError)
			} else {
				tfPlanOnPremisesProvisioningError.PropertyCausingError = types.StringNull()
			}
			// END PropertyCausingError | CreateStringAttribute

			// START Value | CreateStringAttribute
			if !tfPlanOnPremisesProvisioningError.Value.IsUnknown() {
				tfPlanValue := tfPlanOnPremisesProvisioningError.Value.ValueString()
				requestBodyOnPremisesProvisioningError.SetValue(&tfPlanValue)
			} else {
				tfPlanOnPremisesProvisioningError.Value = types.StringNull()
			}
			// END Value | CreateStringAttribute

		}
		requestBodyUser.SetOnPremisesProvisioningErrors(requestBodyOnPremisesProvisioningErrors)
	} else {
		tfPlanUser.OnPremisesProvisioningErrors = types.ListNull(tfPlanUser.OnPremisesProvisioningErrors.ElementType(ctx))
	}
	// END OnPremisesProvisioningErrors | CreateArrayObjectAttribute

	// START OnPremisesSamAccountName | CreateStringAttribute
	if !tfPlanUser.OnPremisesSamAccountName.IsUnknown() {
		tfPlanOnPremisesSamAccountName := tfPlanUser.OnPremisesSamAccountName.ValueString()
		requestBodyUser.SetOnPremisesSamAccountName(&tfPlanOnPremisesSamAccountName)
	} else {
		tfPlanUser.OnPremisesSamAccountName = types.StringNull()
	}
	// END OnPremisesSamAccountName | CreateStringAttribute

	// START OnPremisesSecurityIdentifier | CreateStringAttribute
	if !tfPlanUser.OnPremisesSecurityIdentifier.IsUnknown() {
		tfPlanOnPremisesSecurityIdentifier := tfPlanUser.OnPremisesSecurityIdentifier.ValueString()
		requestBodyUser.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	} else {
		tfPlanUser.OnPremisesSecurityIdentifier = types.StringNull()
	}
	// END OnPremisesSecurityIdentifier | CreateStringAttribute

	// START OnPremisesSyncEnabled | CreateBoolAttribute
	if !tfPlanUser.OnPremisesSyncEnabled.IsUnknown() {
		tfPlanOnPremisesSyncEnabled := tfPlanUser.OnPremisesSyncEnabled.ValueBool()
		requestBodyUser.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	} else {
		tfPlanUser.OnPremisesSyncEnabled = types.BoolNull()
	}
	// END OnPremisesSyncEnabled | CreateBoolAttribute

	// START OnPremisesUserPrincipalName | CreateStringAttribute
	if !tfPlanUser.OnPremisesUserPrincipalName.IsUnknown() {
		tfPlanOnPremisesUserPrincipalName := tfPlanUser.OnPremisesUserPrincipalName.ValueString()
		requestBodyUser.SetOnPremisesUserPrincipalName(&tfPlanOnPremisesUserPrincipalName)
	} else {
		tfPlanUser.OnPremisesUserPrincipalName = types.StringNull()
	}
	// END OnPremisesUserPrincipalName | CreateStringAttribute

	// START OtherMails | CreateArrayStringAttribute
	if len(tfPlanUser.OtherMails.Elements()) > 0 {
		var stringArrayOtherMails []string
		for _, i := range tfPlanUser.OtherMails.Elements() {
			stringArrayOtherMails = append(stringArrayOtherMails, i.String())
		}
		requestBodyUser.SetOtherMails(stringArrayOtherMails)
	} else {
		tfPlanUser.OtherMails = types.ListNull(types.StringType)
	}
	// END OtherMails | CreateArrayStringAttribute

	// START PasswordPolicies | CreateStringAttribute
	if !tfPlanUser.PasswordPolicies.IsUnknown() {
		tfPlanPasswordPolicies := tfPlanUser.PasswordPolicies.ValueString()
		requestBodyUser.SetPasswordPolicies(&tfPlanPasswordPolicies)
	} else {
		tfPlanUser.PasswordPolicies = types.StringNull()
	}
	// END PasswordPolicies | CreateStringAttribute

	// START PasswordProfile | CreateObjectAttribute
	if !tfPlanUser.PasswordProfile.IsUnknown() {
		requestBodyPasswordProfile := models.NewPasswordProfile()
		tfPlanPasswordProfile := userPasswordProfileModel{}
		tfPlanUser.PasswordProfile.As(ctx, &tfPlanPasswordProfile, basetypes.ObjectAsOptions{})

		// START ForceChangePasswordNextSignIn | CreateBoolAttribute
		if !tfPlanPasswordProfile.ForceChangePasswordNextSignIn.IsUnknown() {
			tfPlanForceChangePasswordNextSignIn := tfPlanPasswordProfile.ForceChangePasswordNextSignIn.ValueBool()
			requestBodyPasswordProfile.SetForceChangePasswordNextSignIn(&tfPlanForceChangePasswordNextSignIn)
		} else {
			tfPlanPasswordProfile.ForceChangePasswordNextSignIn = types.BoolNull()
		}
		// END ForceChangePasswordNextSignIn | CreateBoolAttribute

		// START ForceChangePasswordNextSignInWithMfa | CreateBoolAttribute
		if !tfPlanPasswordProfile.ForceChangePasswordNextSignInWithMfa.IsUnknown() {
			tfPlanForceChangePasswordNextSignInWithMfa := tfPlanPasswordProfile.ForceChangePasswordNextSignInWithMfa.ValueBool()
			requestBodyPasswordProfile.SetForceChangePasswordNextSignInWithMfa(&tfPlanForceChangePasswordNextSignInWithMfa)
		} else {
			tfPlanPasswordProfile.ForceChangePasswordNextSignInWithMfa = types.BoolNull()
		}
		// END ForceChangePasswordNextSignInWithMfa | CreateBoolAttribute

		// START Password | CreateStringAttribute
		if !tfPlanPasswordProfile.Password.IsUnknown() {
			tfPlanPassword := tfPlanPasswordProfile.Password.ValueString()
			requestBodyPasswordProfile.SetPassword(&tfPlanPassword)
		} else {
			tfPlanPasswordProfile.Password = types.StringNull()
		}
		// END Password | CreateStringAttribute

		requestBodyUser.SetPasswordProfile(requestBodyPasswordProfile)
		tfPlanUser.PasswordProfile, _ = types.ObjectValueFrom(ctx, tfPlanPasswordProfile.AttributeTypes(), requestBodyPasswordProfile)
	} else {
		tfPlanUser.PasswordProfile = types.ObjectNull(tfPlanUser.PasswordProfile.AttributeTypes(ctx))
	}
	// END PasswordProfile | CreateObjectAttribute

	// START PastProjects | CreateArrayStringAttribute
	if len(tfPlanUser.PastProjects.Elements()) > 0 {
		var stringArrayPastProjects []string
		for _, i := range tfPlanUser.PastProjects.Elements() {
			stringArrayPastProjects = append(stringArrayPastProjects, i.String())
		}
		requestBodyUser.SetPastProjects(stringArrayPastProjects)
	} else {
		tfPlanUser.PastProjects = types.ListNull(types.StringType)
	}
	// END PastProjects | CreateArrayStringAttribute

	// START PostalCode | CreateStringAttribute
	if !tfPlanUser.PostalCode.IsUnknown() {
		tfPlanPostalCode := tfPlanUser.PostalCode.ValueString()
		requestBodyUser.SetPostalCode(&tfPlanPostalCode)
	} else {
		tfPlanUser.PostalCode = types.StringNull()
	}
	// END PostalCode | CreateStringAttribute

	// START PreferredDataLocation | CreateStringAttribute
	if !tfPlanUser.PreferredDataLocation.IsUnknown() {
		tfPlanPreferredDataLocation := tfPlanUser.PreferredDataLocation.ValueString()
		requestBodyUser.SetPreferredDataLocation(&tfPlanPreferredDataLocation)
	} else {
		tfPlanUser.PreferredDataLocation = types.StringNull()
	}
	// END PreferredDataLocation | CreateStringAttribute

	// START PreferredLanguage | CreateStringAttribute
	if !tfPlanUser.PreferredLanguage.IsUnknown() {
		tfPlanPreferredLanguage := tfPlanUser.PreferredLanguage.ValueString()
		requestBodyUser.SetPreferredLanguage(&tfPlanPreferredLanguage)
	} else {
		tfPlanUser.PreferredLanguage = types.StringNull()
	}
	// END PreferredLanguage | CreateStringAttribute

	// START PreferredName | CreateStringAttribute
	if !tfPlanUser.PreferredName.IsUnknown() {
		tfPlanPreferredName := tfPlanUser.PreferredName.ValueString()
		requestBodyUser.SetPreferredName(&tfPlanPreferredName)
	} else {
		tfPlanUser.PreferredName = types.StringNull()
	}
	// END PreferredName | CreateStringAttribute

	// START ProvisionedPlans | CreateArrayObjectAttribute
	if len(tfPlanUser.ProvisionedPlans.Elements()) > 0 {
		var requestBodyProvisionedPlans []models.ProvisionedPlanable
		for _, i := range tfPlanUser.ProvisionedPlans.Elements() {
			requestBodyProvisionedPlan := models.NewProvisionedPlan()
			tfPlanProvisionedPlan := userProvisionedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanProvisionedPlan)

			// START CapabilityStatus | CreateStringAttribute
			if !tfPlanProvisionedPlan.CapabilityStatus.IsUnknown() {
				tfPlanCapabilityStatus := tfPlanProvisionedPlan.CapabilityStatus.ValueString()
				requestBodyProvisionedPlan.SetCapabilityStatus(&tfPlanCapabilityStatus)
			} else {
				tfPlanProvisionedPlan.CapabilityStatus = types.StringNull()
			}
			// END CapabilityStatus | CreateStringAttribute

			// START ProvisioningStatus | CreateStringAttribute
			if !tfPlanProvisionedPlan.ProvisioningStatus.IsUnknown() {
				tfPlanProvisioningStatus := tfPlanProvisionedPlan.ProvisioningStatus.ValueString()
				requestBodyProvisionedPlan.SetProvisioningStatus(&tfPlanProvisioningStatus)
			} else {
				tfPlanProvisionedPlan.ProvisioningStatus = types.StringNull()
			}
			// END ProvisioningStatus | CreateStringAttribute

			// START Service | CreateStringAttribute
			if !tfPlanProvisionedPlan.Service.IsUnknown() {
				tfPlanService := tfPlanProvisionedPlan.Service.ValueString()
				requestBodyProvisionedPlan.SetService(&tfPlanService)
			} else {
				tfPlanProvisionedPlan.Service = types.StringNull()
			}
			// END Service | CreateStringAttribute

		}
		requestBodyUser.SetProvisionedPlans(requestBodyProvisionedPlans)
	} else {
		tfPlanUser.ProvisionedPlans = types.ListNull(tfPlanUser.ProvisionedPlans.ElementType(ctx))
	}
	// END ProvisionedPlans | CreateArrayObjectAttribute

	// START ProxyAddresses | CreateArrayStringAttribute
	if len(tfPlanUser.ProxyAddresses.Elements()) > 0 {
		var stringArrayProxyAddresses []string
		for _, i := range tfPlanUser.ProxyAddresses.Elements() {
			stringArrayProxyAddresses = append(stringArrayProxyAddresses, i.String())
		}
		requestBodyUser.SetProxyAddresses(stringArrayProxyAddresses)
	} else {
		tfPlanUser.ProxyAddresses = types.ListNull(types.StringType)
	}
	// END ProxyAddresses | CreateArrayStringAttribute

	// START Responsibilities | CreateArrayStringAttribute
	if len(tfPlanUser.Responsibilities.Elements()) > 0 {
		var stringArrayResponsibilities []string
		for _, i := range tfPlanUser.Responsibilities.Elements() {
			stringArrayResponsibilities = append(stringArrayResponsibilities, i.String())
		}
		requestBodyUser.SetResponsibilities(stringArrayResponsibilities)
	} else {
		tfPlanUser.Responsibilities = types.ListNull(types.StringType)
	}
	// END Responsibilities | CreateArrayStringAttribute

	// START Schools | CreateArrayStringAttribute
	if len(tfPlanUser.Schools.Elements()) > 0 {
		var stringArraySchools []string
		for _, i := range tfPlanUser.Schools.Elements() {
			stringArraySchools = append(stringArraySchools, i.String())
		}
		requestBodyUser.SetSchools(stringArraySchools)
	} else {
		tfPlanUser.Schools = types.ListNull(types.StringType)
	}
	// END Schools | CreateArrayStringAttribute

	// START SecurityIdentifier | CreateStringAttribute
	if !tfPlanUser.SecurityIdentifier.IsUnknown() {
		tfPlanSecurityIdentifier := tfPlanUser.SecurityIdentifier.ValueString()
		requestBodyUser.SetSecurityIdentifier(&tfPlanSecurityIdentifier)
	} else {
		tfPlanUser.SecurityIdentifier = types.StringNull()
	}
	// END SecurityIdentifier | CreateStringAttribute

	// START ServiceProvisioningErrors | CreateArrayObjectAttribute
	if len(tfPlanUser.ServiceProvisioningErrors.Elements()) > 0 {
		var requestBodyServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for _, i := range tfPlanUser.ServiceProvisioningErrors.Elements() {
			requestBodyServiceProvisioningError := models.NewServiceProvisioningError()
			tfPlanServiceProvisioningError := userServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanServiceProvisioningError)

			// START CreatedDateTime | CreateStringTimeAttribute
			if !tfPlanServiceProvisioningError.CreatedDateTime.IsUnknown() {
				tfPlanCreatedDateTime := tfPlanServiceProvisioningError.CreatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
				requestBodyServiceProvisioningError.SetCreatedDateTime(&t)
			} else {
				tfPlanServiceProvisioningError.CreatedDateTime = types.StringNull()
			}
			// END CreatedDateTime | CreateStringTimeAttribute

			// START IsResolved | CreateBoolAttribute
			if !tfPlanServiceProvisioningError.IsResolved.IsUnknown() {
				tfPlanIsResolved := tfPlanServiceProvisioningError.IsResolved.ValueBool()
				requestBodyServiceProvisioningError.SetIsResolved(&tfPlanIsResolved)
			} else {
				tfPlanServiceProvisioningError.IsResolved = types.BoolNull()
			}
			// END IsResolved | CreateBoolAttribute

			// START ServiceInstance | CreateStringAttribute
			if !tfPlanServiceProvisioningError.ServiceInstance.IsUnknown() {
				tfPlanServiceInstance := tfPlanServiceProvisioningError.ServiceInstance.ValueString()
				requestBodyServiceProvisioningError.SetServiceInstance(&tfPlanServiceInstance)
			} else {
				tfPlanServiceProvisioningError.ServiceInstance = types.StringNull()
			}
			// END ServiceInstance | CreateStringAttribute

		}
		requestBodyUser.SetServiceProvisioningErrors(requestBodyServiceProvisioningErrors)
	} else {
		tfPlanUser.ServiceProvisioningErrors = types.ListNull(tfPlanUser.ServiceProvisioningErrors.ElementType(ctx))
	}
	// END ServiceProvisioningErrors | CreateArrayObjectAttribute

	// START ShowInAddressList | CreateBoolAttribute
	if !tfPlanUser.ShowInAddressList.IsUnknown() {
		tfPlanShowInAddressList := tfPlanUser.ShowInAddressList.ValueBool()
		requestBodyUser.SetShowInAddressList(&tfPlanShowInAddressList)
	} else {
		tfPlanUser.ShowInAddressList = types.BoolNull()
	}
	// END ShowInAddressList | CreateBoolAttribute

	// START SignInActivity | CreateObjectAttribute
	if !tfPlanUser.SignInActivity.IsUnknown() {
		requestBodySignInActivity := models.NewSignInActivity()
		tfPlanSignInActivity := userSignInActivityModel{}
		tfPlanUser.SignInActivity.As(ctx, &tfPlanSignInActivity, basetypes.ObjectAsOptions{})

		// START LastNonInteractiveSignInDateTime | CreateStringTimeAttribute
		if !tfPlanSignInActivity.LastNonInteractiveSignInDateTime.IsUnknown() {
			tfPlanLastNonInteractiveSignInDateTime := tfPlanSignInActivity.LastNonInteractiveSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastNonInteractiveSignInDateTime)
			requestBodySignInActivity.SetLastNonInteractiveSignInDateTime(&t)
		} else {
			tfPlanSignInActivity.LastNonInteractiveSignInDateTime = types.StringNull()
		}
		// END LastNonInteractiveSignInDateTime | CreateStringTimeAttribute

		// START LastNonInteractiveSignInRequestId | CreateStringAttribute
		if !tfPlanSignInActivity.LastNonInteractiveSignInRequestId.IsUnknown() {
			tfPlanLastNonInteractiveSignInRequestId := tfPlanSignInActivity.LastNonInteractiveSignInRequestId.ValueString()
			requestBodySignInActivity.SetLastNonInteractiveSignInRequestId(&tfPlanLastNonInteractiveSignInRequestId)
		} else {
			tfPlanSignInActivity.LastNonInteractiveSignInRequestId = types.StringNull()
		}
		// END LastNonInteractiveSignInRequestId | CreateStringAttribute

		// START LastSignInDateTime | CreateStringTimeAttribute
		if !tfPlanSignInActivity.LastSignInDateTime.IsUnknown() {
			tfPlanLastSignInDateTime := tfPlanSignInActivity.LastSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastSignInDateTime)
			requestBodySignInActivity.SetLastSignInDateTime(&t)
		} else {
			tfPlanSignInActivity.LastSignInDateTime = types.StringNull()
		}
		// END LastSignInDateTime | CreateStringTimeAttribute

		// START LastSignInRequestId | CreateStringAttribute
		if !tfPlanSignInActivity.LastSignInRequestId.IsUnknown() {
			tfPlanLastSignInRequestId := tfPlanSignInActivity.LastSignInRequestId.ValueString()
			requestBodySignInActivity.SetLastSignInRequestId(&tfPlanLastSignInRequestId)
		} else {
			tfPlanSignInActivity.LastSignInRequestId = types.StringNull()
		}
		// END LastSignInRequestId | CreateStringAttribute

		// START LastSuccessfulSignInDateTime | CreateStringTimeAttribute
		if !tfPlanSignInActivity.LastSuccessfulSignInDateTime.IsUnknown() {
			tfPlanLastSuccessfulSignInDateTime := tfPlanSignInActivity.LastSuccessfulSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastSuccessfulSignInDateTime)
			requestBodySignInActivity.SetLastSuccessfulSignInDateTime(&t)
		} else {
			tfPlanSignInActivity.LastSuccessfulSignInDateTime = types.StringNull()
		}
		// END LastSuccessfulSignInDateTime | CreateStringTimeAttribute

		// START LastSuccessfulSignInRequestId | CreateStringAttribute
		if !tfPlanSignInActivity.LastSuccessfulSignInRequestId.IsUnknown() {
			tfPlanLastSuccessfulSignInRequestId := tfPlanSignInActivity.LastSuccessfulSignInRequestId.ValueString()
			requestBodySignInActivity.SetLastSuccessfulSignInRequestId(&tfPlanLastSuccessfulSignInRequestId)
		} else {
			tfPlanSignInActivity.LastSuccessfulSignInRequestId = types.StringNull()
		}
		// END LastSuccessfulSignInRequestId | CreateStringAttribute

		requestBodyUser.SetSignInActivity(requestBodySignInActivity)
		tfPlanUser.SignInActivity, _ = types.ObjectValueFrom(ctx, tfPlanSignInActivity.AttributeTypes(), requestBodySignInActivity)
	} else {
		tfPlanUser.SignInActivity = types.ObjectNull(tfPlanUser.SignInActivity.AttributeTypes(ctx))
	}
	// END SignInActivity | CreateObjectAttribute

	// START SignInSessionsValidFromDateTime | CreateStringTimeAttribute
	if !tfPlanUser.SignInSessionsValidFromDateTime.IsUnknown() {
		tfPlanSignInSessionsValidFromDateTime := tfPlanUser.SignInSessionsValidFromDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanSignInSessionsValidFromDateTime)
		requestBodyUser.SetSignInSessionsValidFromDateTime(&t)
	} else {
		tfPlanUser.SignInSessionsValidFromDateTime = types.StringNull()
	}
	// END SignInSessionsValidFromDateTime | CreateStringTimeAttribute

	// START Skills | CreateArrayStringAttribute
	if len(tfPlanUser.Skills.Elements()) > 0 {
		var stringArraySkills []string
		for _, i := range tfPlanUser.Skills.Elements() {
			stringArraySkills = append(stringArraySkills, i.String())
		}
		requestBodyUser.SetSkills(stringArraySkills)
	} else {
		tfPlanUser.Skills = types.ListNull(types.StringType)
	}
	// END Skills | CreateArrayStringAttribute

	// START State | CreateStringAttribute
	if !tfPlanUser.State.IsUnknown() {
		tfPlanState := tfPlanUser.State.ValueString()
		requestBodyUser.SetState(&tfPlanState)
	} else {
		tfPlanUser.State = types.StringNull()
	}
	// END State | CreateStringAttribute

	// START StreetAddress | CreateStringAttribute
	if !tfPlanUser.StreetAddress.IsUnknown() {
		tfPlanStreetAddress := tfPlanUser.StreetAddress.ValueString()
		requestBodyUser.SetStreetAddress(&tfPlanStreetAddress)
	} else {
		tfPlanUser.StreetAddress = types.StringNull()
	}
	// END StreetAddress | CreateStringAttribute

	// START Surname | CreateStringAttribute
	if !tfPlanUser.Surname.IsUnknown() {
		tfPlanSurname := tfPlanUser.Surname.ValueString()
		requestBodyUser.SetSurname(&tfPlanSurname)
	} else {
		tfPlanUser.Surname = types.StringNull()
	}
	// END Surname | CreateStringAttribute

	// START UsageLocation | CreateStringAttribute
	if !tfPlanUser.UsageLocation.IsUnknown() {
		tfPlanUsageLocation := tfPlanUser.UsageLocation.ValueString()
		requestBodyUser.SetUsageLocation(&tfPlanUsageLocation)
	} else {
		tfPlanUser.UsageLocation = types.StringNull()
	}
	// END UsageLocation | CreateStringAttribute

	// START UserPrincipalName | CreateStringAttribute
	if !tfPlanUser.UserPrincipalName.IsUnknown() {
		tfPlanUserPrincipalName := tfPlanUser.UserPrincipalName.ValueString()
		requestBodyUser.SetUserPrincipalName(&tfPlanUserPrincipalName)
	} else {
		tfPlanUser.UserPrincipalName = types.StringNull()
	}
	// END UserPrincipalName | CreateStringAttribute

	// START UserType | CreateStringAttribute
	if !tfPlanUser.UserType.IsUnknown() {
		tfPlanUserType := tfPlanUser.UserType.ValueString()
		requestBodyUser.SetUserType(&tfPlanUserType)
	} else {
		tfPlanUser.UserType = types.StringNull()
	}
	// END UserType | CreateStringAttribute

	// Create new user
	result, err := r.client.Users().Post(context.Background(), requestBodyUser, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating user",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlanUser.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlanUser)
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
	// Retrieve values from Terraform plan
	var tfPlanUser userModel
	diags := req.Plan.Get(ctx, &tfPlanUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfStateUser userModel
	diags = req.State.Get(ctx, &tfStateUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBodyUser := models.NewUser()

	if !tfPlanUser.Id.Equal(tfStateUser.Id) {
		tfPlanId := tfPlanUser.Id.ValueString()
		requestBodyUser.SetId(&tfPlanId)
	}

	if !tfPlanUser.DeletedDateTime.Equal(tfStateUser.DeletedDateTime) {
		tfPlanDeletedDateTime := tfPlanUser.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyUser.SetDeletedDateTime(&t)
	}

	if !tfPlanUser.AboutMe.Equal(tfStateUser.AboutMe) {
		tfPlanAboutMe := tfPlanUser.AboutMe.ValueString()
		requestBodyUser.SetAboutMe(&tfPlanAboutMe)
	}

	if !tfPlanUser.AccountEnabled.Equal(tfStateUser.AccountEnabled) {
		tfPlanAccountEnabled := tfPlanUser.AccountEnabled.ValueBool()
		requestBodyUser.SetAccountEnabled(&tfPlanAccountEnabled)
	}

	if !tfPlanUser.AgeGroup.Equal(tfStateUser.AgeGroup) {
		tfPlanAgeGroup := tfPlanUser.AgeGroup.ValueString()
		requestBodyUser.SetAgeGroup(&tfPlanAgeGroup)
	}

	if !tfPlanUser.AssignedLicenses.Equal(tfStateUser.AssignedLicenses) {
		var tfPlanAssignedLicenses []models.AssignedLicenseable
		for k, i := range tfPlanUser.AssignedLicenses.Elements() {
			requestBodyAssignedLicense := models.NewAssignedLicense()
			tfPlanAssignedLicense := userAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLicense)
			tfStateAssignedLicense := userAssignedLicenseModel{}
			types.ListValueFrom(ctx, tfStateUser.AssignedLicenses.Elements()[k].Type(ctx), &tfPlanAssignedLicense)

			if !tfPlanAssignedLicense.DisabledPlans.Equal(tfStateAssignedLicense.DisabledPlans) {
				var DisabledPlans []uuid.UUID
				for _, i := range tfPlanAssignedLicense.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				requestBodyAssignedLicense.SetDisabledPlans(DisabledPlans)
			}

			if !tfPlanAssignedLicense.SkuId.Equal(tfStateAssignedLicense.SkuId) {
				tfPlanSkuId := tfPlanAssignedLicense.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				requestBodyAssignedLicense.SetSkuId(&u)
			}
		}
		requestBodyUser.SetAssignedLicenses(tfPlanAssignedLicenses)
	}

	if !tfPlanUser.AssignedPlans.Equal(tfStateUser.AssignedPlans) {
		var tfPlanAssignedPlans []models.AssignedPlanable
		for k, i := range tfPlanUser.AssignedPlans.Elements() {
			requestBodyAssignedPlan := models.NewAssignedPlan()
			tfPlanAssignedPlan := userAssignedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedPlan)
			tfStateAssignedPlan := userAssignedPlanModel{}
			types.ListValueFrom(ctx, tfStateUser.AssignedPlans.Elements()[k].Type(ctx), &tfPlanAssignedPlan)

			if !tfPlanAssignedPlan.AssignedDateTime.Equal(tfStateAssignedPlan.AssignedDateTime) {
				tfPlanAssignedDateTime := tfPlanAssignedPlan.AssignedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanAssignedDateTime)
				requestBodyAssignedPlan.SetAssignedDateTime(&t)
			}

			if !tfPlanAssignedPlan.CapabilityStatus.Equal(tfStateAssignedPlan.CapabilityStatus) {
				tfPlanCapabilityStatus := tfPlanAssignedPlan.CapabilityStatus.ValueString()
				requestBodyAssignedPlan.SetCapabilityStatus(&tfPlanCapabilityStatus)
			}

			if !tfPlanAssignedPlan.Service.Equal(tfStateAssignedPlan.Service) {
				tfPlanService := tfPlanAssignedPlan.Service.ValueString()
				requestBodyAssignedPlan.SetService(&tfPlanService)
			}

			if !tfPlanAssignedPlan.ServicePlanId.Equal(tfStateAssignedPlan.ServicePlanId) {
				tfPlanServicePlanId := tfPlanAssignedPlan.ServicePlanId.ValueString()
				u, _ := uuid.Parse(tfPlanServicePlanId)
				requestBodyAssignedPlan.SetServicePlanId(&u)
			}
		}
		requestBodyUser.SetAssignedPlans(tfPlanAssignedPlans)
	}

	if !tfPlanUser.AuthorizationInfo.Equal(tfStateUser.AuthorizationInfo) {
		requestBodyAuthorizationInfo := models.NewAuthorizationInfo()
		tfPlanAuthorizationInfo := userAuthorizationInfoModel{}
		tfPlanUser.AuthorizationInfo.As(ctx, &tfPlanAuthorizationInfo, basetypes.ObjectAsOptions{})
		tfStateAuthorizationInfo := userAuthorizationInfoModel{}
		tfStateUser.AuthorizationInfo.As(ctx, &tfStateAuthorizationInfo, basetypes.ObjectAsOptions{})

		if !tfPlanAuthorizationInfo.CertificateUserIds.Equal(tfStateAuthorizationInfo.CertificateUserIds) {
			var stringArrayCertificateUserIds []string
			for _, i := range tfPlanAuthorizationInfo.CertificateUserIds.Elements() {
				stringArrayCertificateUserIds = append(stringArrayCertificateUserIds, i.String())
			}
			requestBodyAuthorizationInfo.SetCertificateUserIds(stringArrayCertificateUserIds)
		}
		requestBodyUser.SetAuthorizationInfo(requestBodyAuthorizationInfo)
		tfPlanUser.AuthorizationInfo, _ = types.ObjectValueFrom(ctx, tfPlanAuthorizationInfo.AttributeTypes(), tfPlanAuthorizationInfo)
	}

	if !tfPlanUser.Birthday.Equal(tfStateUser.Birthday) {
		tfPlanBirthday := tfPlanUser.Birthday.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanBirthday)
		requestBodyUser.SetBirthday(&t)
	}

	if !tfPlanUser.BusinessPhones.Equal(tfStateUser.BusinessPhones) {
		var stringArrayBusinessPhones []string
		for _, i := range tfPlanUser.BusinessPhones.Elements() {
			stringArrayBusinessPhones = append(stringArrayBusinessPhones, i.String())
		}
		requestBodyUser.SetBusinessPhones(stringArrayBusinessPhones)
	}

	if !tfPlanUser.City.Equal(tfStateUser.City) {
		tfPlanCity := tfPlanUser.City.ValueString()
		requestBodyUser.SetCity(&tfPlanCity)
	}

	if !tfPlanUser.CompanyName.Equal(tfStateUser.CompanyName) {
		tfPlanCompanyName := tfPlanUser.CompanyName.ValueString()
		requestBodyUser.SetCompanyName(&tfPlanCompanyName)
	}

	if !tfPlanUser.ConsentProvidedForMinor.Equal(tfStateUser.ConsentProvidedForMinor) {
		tfPlanConsentProvidedForMinor := tfPlanUser.ConsentProvidedForMinor.ValueString()
		requestBodyUser.SetConsentProvidedForMinor(&tfPlanConsentProvidedForMinor)
	}

	if !tfPlanUser.Country.Equal(tfStateUser.Country) {
		tfPlanCountry := tfPlanUser.Country.ValueString()
		requestBodyUser.SetCountry(&tfPlanCountry)
	}

	if !tfPlanUser.CreatedDateTime.Equal(tfStateUser.CreatedDateTime) {
		tfPlanCreatedDateTime := tfPlanUser.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyUser.SetCreatedDateTime(&t)
	}

	if !tfPlanUser.CreationType.Equal(tfStateUser.CreationType) {
		tfPlanCreationType := tfPlanUser.CreationType.ValueString()
		requestBodyUser.SetCreationType(&tfPlanCreationType)
	}

	if !tfPlanUser.Department.Equal(tfStateUser.Department) {
		tfPlanDepartment := tfPlanUser.Department.ValueString()
		requestBodyUser.SetDepartment(&tfPlanDepartment)
	}

	if !tfPlanUser.DisplayName.Equal(tfStateUser.DisplayName) {
		tfPlanDisplayName := tfPlanUser.DisplayName.ValueString()
		requestBodyUser.SetDisplayName(&tfPlanDisplayName)
	}

	if !tfPlanUser.EmployeeHireDate.Equal(tfStateUser.EmployeeHireDate) {
		tfPlanEmployeeHireDate := tfPlanUser.EmployeeHireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanEmployeeHireDate)
		requestBodyUser.SetEmployeeHireDate(&t)
	}

	if !tfPlanUser.EmployeeId.Equal(tfStateUser.EmployeeId) {
		tfPlanEmployeeId := tfPlanUser.EmployeeId.ValueString()
		requestBodyUser.SetEmployeeId(&tfPlanEmployeeId)
	}

	if !tfPlanUser.EmployeeLeaveDateTime.Equal(tfStateUser.EmployeeLeaveDateTime) {
		tfPlanEmployeeLeaveDateTime := tfPlanUser.EmployeeLeaveDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanEmployeeLeaveDateTime)
		requestBodyUser.SetEmployeeLeaveDateTime(&t)
	}

	if !tfPlanUser.EmployeeOrgData.Equal(tfStateUser.EmployeeOrgData) {
		requestBodyEmployeeOrgData := models.NewEmployeeOrgData()
		tfPlanEmployeeOrgData := userEmployeeOrgDataModel{}
		tfPlanUser.EmployeeOrgData.As(ctx, &tfPlanEmployeeOrgData, basetypes.ObjectAsOptions{})
		tfStateEmployeeOrgData := userEmployeeOrgDataModel{}
		tfStateUser.EmployeeOrgData.As(ctx, &tfStateEmployeeOrgData, basetypes.ObjectAsOptions{})

		if !tfPlanEmployeeOrgData.CostCenter.Equal(tfStateEmployeeOrgData.CostCenter) {
			tfPlanCostCenter := tfPlanEmployeeOrgData.CostCenter.ValueString()
			requestBodyEmployeeOrgData.SetCostCenter(&tfPlanCostCenter)
		}

		if !tfPlanEmployeeOrgData.Division.Equal(tfStateEmployeeOrgData.Division) {
			tfPlanDivision := tfPlanEmployeeOrgData.Division.ValueString()
			requestBodyEmployeeOrgData.SetDivision(&tfPlanDivision)
		}
		requestBodyUser.SetEmployeeOrgData(requestBodyEmployeeOrgData)
		tfPlanUser.EmployeeOrgData, _ = types.ObjectValueFrom(ctx, tfPlanEmployeeOrgData.AttributeTypes(), tfPlanEmployeeOrgData)
	}

	if !tfPlanUser.EmployeeType.Equal(tfStateUser.EmployeeType) {
		tfPlanEmployeeType := tfPlanUser.EmployeeType.ValueString()
		requestBodyUser.SetEmployeeType(&tfPlanEmployeeType)
	}

	if !tfPlanUser.ExternalUserState.Equal(tfStateUser.ExternalUserState) {
		tfPlanExternalUserState := tfPlanUser.ExternalUserState.ValueString()
		requestBodyUser.SetExternalUserState(&tfPlanExternalUserState)
	}

	if !tfPlanUser.ExternalUserStateChangeDateTime.Equal(tfStateUser.ExternalUserStateChangeDateTime) {
		tfPlanExternalUserStateChangeDateTime := tfPlanUser.ExternalUserStateChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanExternalUserStateChangeDateTime)
		requestBodyUser.SetExternalUserStateChangeDateTime(&t)
	}

	if !tfPlanUser.FaxNumber.Equal(tfStateUser.FaxNumber) {
		tfPlanFaxNumber := tfPlanUser.FaxNumber.ValueString()
		requestBodyUser.SetFaxNumber(&tfPlanFaxNumber)
	}

	if !tfPlanUser.GivenName.Equal(tfStateUser.GivenName) {
		tfPlanGivenName := tfPlanUser.GivenName.ValueString()
		requestBodyUser.SetGivenName(&tfPlanGivenName)
	}

	if !tfPlanUser.HireDate.Equal(tfStateUser.HireDate) {
		tfPlanHireDate := tfPlanUser.HireDate.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanHireDate)
		requestBodyUser.SetHireDate(&t)
	}

	if !tfPlanUser.Identities.Equal(tfStateUser.Identities) {
		var tfPlanIdentities []models.ObjectIdentityable
		for k, i := range tfPlanUser.Identities.Elements() {
			requestBodyObjectIdentity := models.NewObjectIdentity()
			tfPlanObjectIdentity := userObjectIdentityModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanObjectIdentity)
			tfStateObjectIdentity := userObjectIdentityModel{}
			types.ListValueFrom(ctx, tfStateUser.Identities.Elements()[k].Type(ctx), &tfPlanObjectIdentity)

			if !tfPlanObjectIdentity.Issuer.Equal(tfStateObjectIdentity.Issuer) {
				tfPlanIssuer := tfPlanObjectIdentity.Issuer.ValueString()
				requestBodyObjectIdentity.SetIssuer(&tfPlanIssuer)
			}

			if !tfPlanObjectIdentity.IssuerAssignedId.Equal(tfStateObjectIdentity.IssuerAssignedId) {
				tfPlanIssuerAssignedId := tfPlanObjectIdentity.IssuerAssignedId.ValueString()
				requestBodyObjectIdentity.SetIssuerAssignedId(&tfPlanIssuerAssignedId)
			}

			if !tfPlanObjectIdentity.SignInType.Equal(tfStateObjectIdentity.SignInType) {
				tfPlanSignInType := tfPlanObjectIdentity.SignInType.ValueString()
				requestBodyObjectIdentity.SetSignInType(&tfPlanSignInType)
			}
		}
		requestBodyUser.SetIdentities(tfPlanIdentities)
	}

	if !tfPlanUser.ImAddresses.Equal(tfStateUser.ImAddresses) {
		var stringArrayImAddresses []string
		for _, i := range tfPlanUser.ImAddresses.Elements() {
			stringArrayImAddresses = append(stringArrayImAddresses, i.String())
		}
		requestBodyUser.SetImAddresses(stringArrayImAddresses)
	}

	if !tfPlanUser.Interests.Equal(tfStateUser.Interests) {
		var stringArrayInterests []string
		for _, i := range tfPlanUser.Interests.Elements() {
			stringArrayInterests = append(stringArrayInterests, i.String())
		}
		requestBodyUser.SetInterests(stringArrayInterests)
	}

	if !tfPlanUser.IsManagementRestricted.Equal(tfStateUser.IsManagementRestricted) {
		tfPlanIsManagementRestricted := tfPlanUser.IsManagementRestricted.ValueBool()
		requestBodyUser.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	}

	if !tfPlanUser.IsResourceAccount.Equal(tfStateUser.IsResourceAccount) {
		tfPlanIsResourceAccount := tfPlanUser.IsResourceAccount.ValueBool()
		requestBodyUser.SetIsResourceAccount(&tfPlanIsResourceAccount)
	}

	if !tfPlanUser.JobTitle.Equal(tfStateUser.JobTitle) {
		tfPlanJobTitle := tfPlanUser.JobTitle.ValueString()
		requestBodyUser.SetJobTitle(&tfPlanJobTitle)
	}

	if !tfPlanUser.LastPasswordChangeDateTime.Equal(tfStateUser.LastPasswordChangeDateTime) {
		tfPlanLastPasswordChangeDateTime := tfPlanUser.LastPasswordChangeDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanLastPasswordChangeDateTime)
		requestBodyUser.SetLastPasswordChangeDateTime(&t)
	}

	if !tfPlanUser.LegalAgeGroupClassification.Equal(tfStateUser.LegalAgeGroupClassification) {
		tfPlanLegalAgeGroupClassification := tfPlanUser.LegalAgeGroupClassification.ValueString()
		requestBodyUser.SetLegalAgeGroupClassification(&tfPlanLegalAgeGroupClassification)
	}

	if !tfPlanUser.LicenseAssignmentStates.Equal(tfStateUser.LicenseAssignmentStates) {
		var tfPlanLicenseAssignmentStates []models.LicenseAssignmentStateable
		for k, i := range tfPlanUser.LicenseAssignmentStates.Elements() {
			requestBodyLicenseAssignmentState := models.NewLicenseAssignmentState()
			tfPlanLicenseAssignmentState := userLicenseAssignmentStateModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanLicenseAssignmentState)
			tfStateLicenseAssignmentState := userLicenseAssignmentStateModel{}
			types.ListValueFrom(ctx, tfStateUser.LicenseAssignmentStates.Elements()[k].Type(ctx), &tfPlanLicenseAssignmentState)

			if !tfPlanLicenseAssignmentState.AssignedByGroup.Equal(tfStateLicenseAssignmentState.AssignedByGroup) {
				tfPlanAssignedByGroup := tfPlanLicenseAssignmentState.AssignedByGroup.ValueString()
				requestBodyLicenseAssignmentState.SetAssignedByGroup(&tfPlanAssignedByGroup)
			}

			if !tfPlanLicenseAssignmentState.DisabledPlans.Equal(tfStateLicenseAssignmentState.DisabledPlans) {
				var DisabledPlans []uuid.UUID
				for _, i := range tfPlanLicenseAssignmentState.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				requestBodyLicenseAssignmentState.SetDisabledPlans(DisabledPlans)
			}

			if !tfPlanLicenseAssignmentState.Error.Equal(tfStateLicenseAssignmentState.Error) {
				tfPlanError := tfPlanLicenseAssignmentState.Error.ValueString()
				requestBodyLicenseAssignmentState.SetError(&tfPlanError)
			}

			if !tfPlanLicenseAssignmentState.LastUpdatedDateTime.Equal(tfStateLicenseAssignmentState.LastUpdatedDateTime) {
				tfPlanLastUpdatedDateTime := tfPlanLicenseAssignmentState.LastUpdatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanLastUpdatedDateTime)
				requestBodyLicenseAssignmentState.SetLastUpdatedDateTime(&t)
			}

			if !tfPlanLicenseAssignmentState.SkuId.Equal(tfStateLicenseAssignmentState.SkuId) {
				tfPlanSkuId := tfPlanLicenseAssignmentState.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				requestBodyLicenseAssignmentState.SetSkuId(&u)
			}

			if !tfPlanLicenseAssignmentState.State.Equal(tfStateLicenseAssignmentState.State) {
				tfPlanState := tfPlanLicenseAssignmentState.State.ValueString()
				requestBodyLicenseAssignmentState.SetState(&tfPlanState)
			}
		}
		requestBodyUser.SetLicenseAssignmentStates(tfPlanLicenseAssignmentStates)
	}

	if !tfPlanUser.Mail.Equal(tfStateUser.Mail) {
		tfPlanMail := tfPlanUser.Mail.ValueString()
		requestBodyUser.SetMail(&tfPlanMail)
	}

	if !tfPlanUser.MailNickname.Equal(tfStateUser.MailNickname) {
		tfPlanMailNickname := tfPlanUser.MailNickname.ValueString()
		requestBodyUser.SetMailNickname(&tfPlanMailNickname)
	}

	if !tfPlanUser.MobilePhone.Equal(tfStateUser.MobilePhone) {
		tfPlanMobilePhone := tfPlanUser.MobilePhone.ValueString()
		requestBodyUser.SetMobilePhone(&tfPlanMobilePhone)
	}

	if !tfPlanUser.MySite.Equal(tfStateUser.MySite) {
		tfPlanMySite := tfPlanUser.MySite.ValueString()
		requestBodyUser.SetMySite(&tfPlanMySite)
	}

	if !tfPlanUser.OfficeLocation.Equal(tfStateUser.OfficeLocation) {
		tfPlanOfficeLocation := tfPlanUser.OfficeLocation.ValueString()
		requestBodyUser.SetOfficeLocation(&tfPlanOfficeLocation)
	}

	if !tfPlanUser.OnPremisesDistinguishedName.Equal(tfStateUser.OnPremisesDistinguishedName) {
		tfPlanOnPremisesDistinguishedName := tfPlanUser.OnPremisesDistinguishedName.ValueString()
		requestBodyUser.SetOnPremisesDistinguishedName(&tfPlanOnPremisesDistinguishedName)
	}

	if !tfPlanUser.OnPremisesDomainName.Equal(tfStateUser.OnPremisesDomainName) {
		tfPlanOnPremisesDomainName := tfPlanUser.OnPremisesDomainName.ValueString()
		requestBodyUser.SetOnPremisesDomainName(&tfPlanOnPremisesDomainName)
	}

	if !tfPlanUser.OnPremisesExtensionAttributes.Equal(tfStateUser.OnPremisesExtensionAttributes) {
		requestBodyOnPremisesExtensionAttributes := models.NewOnPremisesExtensionAttributes()
		tfPlanOnPremisesExtensionAttributes := userOnPremisesExtensionAttributesModel{}
		tfPlanUser.OnPremisesExtensionAttributes.As(ctx, &tfPlanOnPremisesExtensionAttributes, basetypes.ObjectAsOptions{})
		tfStateOnPremisesExtensionAttributes := userOnPremisesExtensionAttributesModel{}
		tfStateUser.OnPremisesExtensionAttributes.As(ctx, &tfStateOnPremisesExtensionAttributes, basetypes.ObjectAsOptions{})

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute1.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute1) {
			tfPlanExtensionAttribute1 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute1.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute1(&tfPlanExtensionAttribute1)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute10.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute10) {
			tfPlanExtensionAttribute10 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute10.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute10(&tfPlanExtensionAttribute10)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute11.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute11) {
			tfPlanExtensionAttribute11 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute11.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute11(&tfPlanExtensionAttribute11)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute12.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute12) {
			tfPlanExtensionAttribute12 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute12.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute12(&tfPlanExtensionAttribute12)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute13.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute13) {
			tfPlanExtensionAttribute13 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute13.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute13(&tfPlanExtensionAttribute13)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute14.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute14) {
			tfPlanExtensionAttribute14 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute14.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute14(&tfPlanExtensionAttribute14)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute15.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute15) {
			tfPlanExtensionAttribute15 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute15.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute15(&tfPlanExtensionAttribute15)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute2.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute2) {
			tfPlanExtensionAttribute2 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute2.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute2(&tfPlanExtensionAttribute2)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute3.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute3) {
			tfPlanExtensionAttribute3 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute3.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute3(&tfPlanExtensionAttribute3)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute4.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute4) {
			tfPlanExtensionAttribute4 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute4.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute4(&tfPlanExtensionAttribute4)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute5.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute5) {
			tfPlanExtensionAttribute5 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute5.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute5(&tfPlanExtensionAttribute5)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute6.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute6) {
			tfPlanExtensionAttribute6 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute6.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute6(&tfPlanExtensionAttribute6)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute7.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute7) {
			tfPlanExtensionAttribute7 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute7.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute7(&tfPlanExtensionAttribute7)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute8.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute8) {
			tfPlanExtensionAttribute8 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute8.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute8(&tfPlanExtensionAttribute8)
		}

		if !tfPlanOnPremisesExtensionAttributes.ExtensionAttribute9.Equal(tfStateOnPremisesExtensionAttributes.ExtensionAttribute9) {
			tfPlanExtensionAttribute9 := tfPlanOnPremisesExtensionAttributes.ExtensionAttribute9.ValueString()
			requestBodyOnPremisesExtensionAttributes.SetExtensionAttribute9(&tfPlanExtensionAttribute9)
		}
		requestBodyUser.SetOnPremisesExtensionAttributes(requestBodyOnPremisesExtensionAttributes)
		tfPlanUser.OnPremisesExtensionAttributes, _ = types.ObjectValueFrom(ctx, tfPlanOnPremisesExtensionAttributes.AttributeTypes(), tfPlanOnPremisesExtensionAttributes)
	}

	if !tfPlanUser.OnPremisesImmutableId.Equal(tfStateUser.OnPremisesImmutableId) {
		tfPlanOnPremisesImmutableId := tfPlanUser.OnPremisesImmutableId.ValueString()
		requestBodyUser.SetOnPremisesImmutableId(&tfPlanOnPremisesImmutableId)
	}

	if !tfPlanUser.OnPremisesLastSyncDateTime.Equal(tfStateUser.OnPremisesLastSyncDateTime) {
		tfPlanOnPremisesLastSyncDateTime := tfPlanUser.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		requestBodyUser.SetOnPremisesLastSyncDateTime(&t)
	}

	if !tfPlanUser.OnPremisesProvisioningErrors.Equal(tfStateUser.OnPremisesProvisioningErrors) {
		var tfPlanOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for k, i := range tfPlanUser.OnPremisesProvisioningErrors.Elements() {
			requestBodyOnPremisesProvisioningError := models.NewOnPremisesProvisioningError()
			tfPlanOnPremisesProvisioningError := userOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOnPremisesProvisioningError)
			tfStateOnPremisesProvisioningError := userOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, tfStateUser.OnPremisesProvisioningErrors.Elements()[k].Type(ctx), &tfPlanOnPremisesProvisioningError)

			if !tfPlanOnPremisesProvisioningError.Category.Equal(tfStateOnPremisesProvisioningError.Category) {
				tfPlanCategory := tfPlanOnPremisesProvisioningError.Category.ValueString()
				requestBodyOnPremisesProvisioningError.SetCategory(&tfPlanCategory)
			}

			if !tfPlanOnPremisesProvisioningError.OccurredDateTime.Equal(tfStateOnPremisesProvisioningError.OccurredDateTime) {
				tfPlanOccurredDateTime := tfPlanOnPremisesProvisioningError.OccurredDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanOccurredDateTime)
				requestBodyOnPremisesProvisioningError.SetOccurredDateTime(&t)
			}

			if !tfPlanOnPremisesProvisioningError.PropertyCausingError.Equal(tfStateOnPremisesProvisioningError.PropertyCausingError) {
				tfPlanPropertyCausingError := tfPlanOnPremisesProvisioningError.PropertyCausingError.ValueString()
				requestBodyOnPremisesProvisioningError.SetPropertyCausingError(&tfPlanPropertyCausingError)
			}

			if !tfPlanOnPremisesProvisioningError.Value.Equal(tfStateOnPremisesProvisioningError.Value) {
				tfPlanValue := tfPlanOnPremisesProvisioningError.Value.ValueString()
				requestBodyOnPremisesProvisioningError.SetValue(&tfPlanValue)
			}
		}
		requestBodyUser.SetOnPremisesProvisioningErrors(tfPlanOnPremisesProvisioningErrors)
	}

	if !tfPlanUser.OnPremisesSamAccountName.Equal(tfStateUser.OnPremisesSamAccountName) {
		tfPlanOnPremisesSamAccountName := tfPlanUser.OnPremisesSamAccountName.ValueString()
		requestBodyUser.SetOnPremisesSamAccountName(&tfPlanOnPremisesSamAccountName)
	}

	if !tfPlanUser.OnPremisesSecurityIdentifier.Equal(tfStateUser.OnPremisesSecurityIdentifier) {
		tfPlanOnPremisesSecurityIdentifier := tfPlanUser.OnPremisesSecurityIdentifier.ValueString()
		requestBodyUser.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	}

	if !tfPlanUser.OnPremisesSyncEnabled.Equal(tfStateUser.OnPremisesSyncEnabled) {
		tfPlanOnPremisesSyncEnabled := tfPlanUser.OnPremisesSyncEnabled.ValueBool()
		requestBodyUser.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	}

	if !tfPlanUser.OnPremisesUserPrincipalName.Equal(tfStateUser.OnPremisesUserPrincipalName) {
		tfPlanOnPremisesUserPrincipalName := tfPlanUser.OnPremisesUserPrincipalName.ValueString()
		requestBodyUser.SetOnPremisesUserPrincipalName(&tfPlanOnPremisesUserPrincipalName)
	}

	if !tfPlanUser.OtherMails.Equal(tfStateUser.OtherMails) {
		var stringArrayOtherMails []string
		for _, i := range tfPlanUser.OtherMails.Elements() {
			stringArrayOtherMails = append(stringArrayOtherMails, i.String())
		}
		requestBodyUser.SetOtherMails(stringArrayOtherMails)
	}

	if !tfPlanUser.PasswordPolicies.Equal(tfStateUser.PasswordPolicies) {
		tfPlanPasswordPolicies := tfPlanUser.PasswordPolicies.ValueString()
		requestBodyUser.SetPasswordPolicies(&tfPlanPasswordPolicies)
	}

	if !tfPlanUser.PasswordProfile.Equal(tfStateUser.PasswordProfile) {
		requestBodyPasswordProfile := models.NewPasswordProfile()
		tfPlanPasswordProfile := userPasswordProfileModel{}
		tfPlanUser.PasswordProfile.As(ctx, &tfPlanPasswordProfile, basetypes.ObjectAsOptions{})
		tfStatePasswordProfile := userPasswordProfileModel{}
		tfStateUser.PasswordProfile.As(ctx, &tfStatePasswordProfile, basetypes.ObjectAsOptions{})

		if !tfPlanPasswordProfile.ForceChangePasswordNextSignIn.Equal(tfStatePasswordProfile.ForceChangePasswordNextSignIn) {
			tfPlanForceChangePasswordNextSignIn := tfPlanPasswordProfile.ForceChangePasswordNextSignIn.ValueBool()
			requestBodyPasswordProfile.SetForceChangePasswordNextSignIn(&tfPlanForceChangePasswordNextSignIn)
		}

		if !tfPlanPasswordProfile.ForceChangePasswordNextSignInWithMfa.Equal(tfStatePasswordProfile.ForceChangePasswordNextSignInWithMfa) {
			tfPlanForceChangePasswordNextSignInWithMfa := tfPlanPasswordProfile.ForceChangePasswordNextSignInWithMfa.ValueBool()
			requestBodyPasswordProfile.SetForceChangePasswordNextSignInWithMfa(&tfPlanForceChangePasswordNextSignInWithMfa)
		}

		if !tfPlanPasswordProfile.Password.Equal(tfStatePasswordProfile.Password) {
			tfPlanPassword := tfPlanPasswordProfile.Password.ValueString()
			requestBodyPasswordProfile.SetPassword(&tfPlanPassword)
		}
		requestBodyUser.SetPasswordProfile(requestBodyPasswordProfile)
		tfPlanUser.PasswordProfile, _ = types.ObjectValueFrom(ctx, tfPlanPasswordProfile.AttributeTypes(), tfPlanPasswordProfile)
	}

	if !tfPlanUser.PastProjects.Equal(tfStateUser.PastProjects) {
		var stringArrayPastProjects []string
		for _, i := range tfPlanUser.PastProjects.Elements() {
			stringArrayPastProjects = append(stringArrayPastProjects, i.String())
		}
		requestBodyUser.SetPastProjects(stringArrayPastProjects)
	}

	if !tfPlanUser.PostalCode.Equal(tfStateUser.PostalCode) {
		tfPlanPostalCode := tfPlanUser.PostalCode.ValueString()
		requestBodyUser.SetPostalCode(&tfPlanPostalCode)
	}

	if !tfPlanUser.PreferredDataLocation.Equal(tfStateUser.PreferredDataLocation) {
		tfPlanPreferredDataLocation := tfPlanUser.PreferredDataLocation.ValueString()
		requestBodyUser.SetPreferredDataLocation(&tfPlanPreferredDataLocation)
	}

	if !tfPlanUser.PreferredLanguage.Equal(tfStateUser.PreferredLanguage) {
		tfPlanPreferredLanguage := tfPlanUser.PreferredLanguage.ValueString()
		requestBodyUser.SetPreferredLanguage(&tfPlanPreferredLanguage)
	}

	if !tfPlanUser.PreferredName.Equal(tfStateUser.PreferredName) {
		tfPlanPreferredName := tfPlanUser.PreferredName.ValueString()
		requestBodyUser.SetPreferredName(&tfPlanPreferredName)
	}

	if !tfPlanUser.ProvisionedPlans.Equal(tfStateUser.ProvisionedPlans) {
		var tfPlanProvisionedPlans []models.ProvisionedPlanable
		for k, i := range tfPlanUser.ProvisionedPlans.Elements() {
			requestBodyProvisionedPlan := models.NewProvisionedPlan()
			tfPlanProvisionedPlan := userProvisionedPlanModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanProvisionedPlan)
			tfStateProvisionedPlan := userProvisionedPlanModel{}
			types.ListValueFrom(ctx, tfStateUser.ProvisionedPlans.Elements()[k].Type(ctx), &tfPlanProvisionedPlan)

			if !tfPlanProvisionedPlan.CapabilityStatus.Equal(tfStateProvisionedPlan.CapabilityStatus) {
				tfPlanCapabilityStatus := tfPlanProvisionedPlan.CapabilityStatus.ValueString()
				requestBodyProvisionedPlan.SetCapabilityStatus(&tfPlanCapabilityStatus)
			}

			if !tfPlanProvisionedPlan.ProvisioningStatus.Equal(tfStateProvisionedPlan.ProvisioningStatus) {
				tfPlanProvisioningStatus := tfPlanProvisionedPlan.ProvisioningStatus.ValueString()
				requestBodyProvisionedPlan.SetProvisioningStatus(&tfPlanProvisioningStatus)
			}

			if !tfPlanProvisionedPlan.Service.Equal(tfStateProvisionedPlan.Service) {
				tfPlanService := tfPlanProvisionedPlan.Service.ValueString()
				requestBodyProvisionedPlan.SetService(&tfPlanService)
			}
		}
		requestBodyUser.SetProvisionedPlans(tfPlanProvisionedPlans)
	}

	if !tfPlanUser.ProxyAddresses.Equal(tfStateUser.ProxyAddresses) {
		var stringArrayProxyAddresses []string
		for _, i := range tfPlanUser.ProxyAddresses.Elements() {
			stringArrayProxyAddresses = append(stringArrayProxyAddresses, i.String())
		}
		requestBodyUser.SetProxyAddresses(stringArrayProxyAddresses)
	}

	if !tfPlanUser.Responsibilities.Equal(tfStateUser.Responsibilities) {
		var stringArrayResponsibilities []string
		for _, i := range tfPlanUser.Responsibilities.Elements() {
			stringArrayResponsibilities = append(stringArrayResponsibilities, i.String())
		}
		requestBodyUser.SetResponsibilities(stringArrayResponsibilities)
	}

	if !tfPlanUser.Schools.Equal(tfStateUser.Schools) {
		var stringArraySchools []string
		for _, i := range tfPlanUser.Schools.Elements() {
			stringArraySchools = append(stringArraySchools, i.String())
		}
		requestBodyUser.SetSchools(stringArraySchools)
	}

	if !tfPlanUser.SecurityIdentifier.Equal(tfStateUser.SecurityIdentifier) {
		tfPlanSecurityIdentifier := tfPlanUser.SecurityIdentifier.ValueString()
		requestBodyUser.SetSecurityIdentifier(&tfPlanSecurityIdentifier)
	}

	if !tfPlanUser.ServiceProvisioningErrors.Equal(tfStateUser.ServiceProvisioningErrors) {
		var tfPlanServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for k, i := range tfPlanUser.ServiceProvisioningErrors.Elements() {
			requestBodyServiceProvisioningError := models.NewServiceProvisioningError()
			tfPlanServiceProvisioningError := userServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanServiceProvisioningError)
			tfStateServiceProvisioningError := userServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, tfStateUser.ServiceProvisioningErrors.Elements()[k].Type(ctx), &tfPlanServiceProvisioningError)

			if !tfPlanServiceProvisioningError.CreatedDateTime.Equal(tfStateServiceProvisioningError.CreatedDateTime) {
				tfPlanCreatedDateTime := tfPlanServiceProvisioningError.CreatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
				requestBodyServiceProvisioningError.SetCreatedDateTime(&t)
			}

			if !tfPlanServiceProvisioningError.IsResolved.Equal(tfStateServiceProvisioningError.IsResolved) {
				tfPlanIsResolved := tfPlanServiceProvisioningError.IsResolved.ValueBool()
				requestBodyServiceProvisioningError.SetIsResolved(&tfPlanIsResolved)
			}

			if !tfPlanServiceProvisioningError.ServiceInstance.Equal(tfStateServiceProvisioningError.ServiceInstance) {
				tfPlanServiceInstance := tfPlanServiceProvisioningError.ServiceInstance.ValueString()
				requestBodyServiceProvisioningError.SetServiceInstance(&tfPlanServiceInstance)
			}
		}
		requestBodyUser.SetServiceProvisioningErrors(tfPlanServiceProvisioningErrors)
	}

	if !tfPlanUser.ShowInAddressList.Equal(tfStateUser.ShowInAddressList) {
		tfPlanShowInAddressList := tfPlanUser.ShowInAddressList.ValueBool()
		requestBodyUser.SetShowInAddressList(&tfPlanShowInAddressList)
	}

	if !tfPlanUser.SignInActivity.Equal(tfStateUser.SignInActivity) {
		requestBodySignInActivity := models.NewSignInActivity()
		tfPlanSignInActivity := userSignInActivityModel{}
		tfPlanUser.SignInActivity.As(ctx, &tfPlanSignInActivity, basetypes.ObjectAsOptions{})
		tfStateSignInActivity := userSignInActivityModel{}
		tfStateUser.SignInActivity.As(ctx, &tfStateSignInActivity, basetypes.ObjectAsOptions{})

		if !tfPlanSignInActivity.LastNonInteractiveSignInDateTime.Equal(tfStateSignInActivity.LastNonInteractiveSignInDateTime) {
			tfPlanLastNonInteractiveSignInDateTime := tfPlanSignInActivity.LastNonInteractiveSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastNonInteractiveSignInDateTime)
			requestBodySignInActivity.SetLastNonInteractiveSignInDateTime(&t)
		}

		if !tfPlanSignInActivity.LastNonInteractiveSignInRequestId.Equal(tfStateSignInActivity.LastNonInteractiveSignInRequestId) {
			tfPlanLastNonInteractiveSignInRequestId := tfPlanSignInActivity.LastNonInteractiveSignInRequestId.ValueString()
			requestBodySignInActivity.SetLastNonInteractiveSignInRequestId(&tfPlanLastNonInteractiveSignInRequestId)
		}

		if !tfPlanSignInActivity.LastSignInDateTime.Equal(tfStateSignInActivity.LastSignInDateTime) {
			tfPlanLastSignInDateTime := tfPlanSignInActivity.LastSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastSignInDateTime)
			requestBodySignInActivity.SetLastSignInDateTime(&t)
		}

		if !tfPlanSignInActivity.LastSignInRequestId.Equal(tfStateSignInActivity.LastSignInRequestId) {
			tfPlanLastSignInRequestId := tfPlanSignInActivity.LastSignInRequestId.ValueString()
			requestBodySignInActivity.SetLastSignInRequestId(&tfPlanLastSignInRequestId)
		}

		if !tfPlanSignInActivity.LastSuccessfulSignInDateTime.Equal(tfStateSignInActivity.LastSuccessfulSignInDateTime) {
			tfPlanLastSuccessfulSignInDateTime := tfPlanSignInActivity.LastSuccessfulSignInDateTime.ValueString()
			t, _ := time.Parse(time.RFC3339, tfPlanLastSuccessfulSignInDateTime)
			requestBodySignInActivity.SetLastSuccessfulSignInDateTime(&t)
		}

		if !tfPlanSignInActivity.LastSuccessfulSignInRequestId.Equal(tfStateSignInActivity.LastSuccessfulSignInRequestId) {
			tfPlanLastSuccessfulSignInRequestId := tfPlanSignInActivity.LastSuccessfulSignInRequestId.ValueString()
			requestBodySignInActivity.SetLastSuccessfulSignInRequestId(&tfPlanLastSuccessfulSignInRequestId)
		}
		requestBodyUser.SetSignInActivity(requestBodySignInActivity)
		tfPlanUser.SignInActivity, _ = types.ObjectValueFrom(ctx, tfPlanSignInActivity.AttributeTypes(), tfPlanSignInActivity)
	}

	if !tfPlanUser.SignInSessionsValidFromDateTime.Equal(tfStateUser.SignInSessionsValidFromDateTime) {
		tfPlanSignInSessionsValidFromDateTime := tfPlanUser.SignInSessionsValidFromDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanSignInSessionsValidFromDateTime)
		requestBodyUser.SetSignInSessionsValidFromDateTime(&t)
	}

	if !tfPlanUser.Skills.Equal(tfStateUser.Skills) {
		var stringArraySkills []string
		for _, i := range tfPlanUser.Skills.Elements() {
			stringArraySkills = append(stringArraySkills, i.String())
		}
		requestBodyUser.SetSkills(stringArraySkills)
	}

	if !tfPlanUser.State.Equal(tfStateUser.State) {
		tfPlanState := tfPlanUser.State.ValueString()
		requestBodyUser.SetState(&tfPlanState)
	}

	if !tfPlanUser.StreetAddress.Equal(tfStateUser.StreetAddress) {
		tfPlanStreetAddress := tfPlanUser.StreetAddress.ValueString()
		requestBodyUser.SetStreetAddress(&tfPlanStreetAddress)
	}

	if !tfPlanUser.Surname.Equal(tfStateUser.Surname) {
		tfPlanSurname := tfPlanUser.Surname.ValueString()
		requestBodyUser.SetSurname(&tfPlanSurname)
	}

	if !tfPlanUser.UsageLocation.Equal(tfStateUser.UsageLocation) {
		tfPlanUsageLocation := tfPlanUser.UsageLocation.ValueString()
		requestBodyUser.SetUsageLocation(&tfPlanUsageLocation)
	}

	if !tfPlanUser.UserPrincipalName.Equal(tfStateUser.UserPrincipalName) {
		tfPlanUserPrincipalName := tfPlanUser.UserPrincipalName.ValueString()
		requestBodyUser.SetUserPrincipalName(&tfPlanUserPrincipalName)
	}

	if !tfPlanUser.UserType.Equal(tfStateUser.UserType) {
		tfPlanUserType := tfPlanUser.UserType.ValueString()
		requestBodyUser.SetUserType(&tfPlanUserType)
	}

	// Update user
	_, err := r.client.Users().ByUserId(tfStateUser.Id.ValueString()).Patch(context.Background(), requestBodyUser, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating user",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlanUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfStateUser userModel
	diags := req.State.Get(ctx, &tfStateUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete user
	err := r.client.Users().ByUserId(tfStateUser.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting user",
			err.Error(),
		)
		return
	}

}
