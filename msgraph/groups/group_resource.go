package groups

import (
	"context"
	"github.com/google/uuid"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"terraform-provider-msgraph/planmodifiers/boolplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/listplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/objectplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/stringplanmodifiers"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &groupResource{}
	_ resource.ResourceWithConfigure = &groupResource{}
)

// NewGroupResource is a helper function to simplify the provider implementation.
func NewGroupResource() resource.Resource {
	return &groupResource{}
}

// groupResource is the resource implementation.
type groupResource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (d *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Configure adds the provider configured client to the resource.
func (d *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the resource.
func (d *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"assigned_labels": schema.ListNestedAttribute{
				Description: "The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group. Returned only on $select.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Description: "The display name of the label. Read-only.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"label_id": schema.StringAttribute{
							Description: "The unique identifier of the label.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Description: "The licenses that are assigned to the group. Returned only on $select. Supports $filter (eq).Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"disabled_plans": schema.ListAttribute{
							Description: "A collection of the unique identifiers for plans that have been disabled.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.List{
								listplanmodifiers.UseStateForUnconfigured(),
							},
							ElementType: types.StringType,
						},
						"sku_id": schema.StringAttribute{
							Description: "The unique identifier for the SKU.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"classification": schema.StringAttribute{
				Description: "Describes a classification for the group (such as low, medium or high business impact). Valid values for this property are defined by creating a ClassificationList setting value, based on the template definition.Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"created_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group was created. The value cannot be modified and is automatically populated when the group is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"description": schema.StringAttribute{
				Description: "An optional description for the group. Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name for the group. This property is required when a group is created and cannot be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"expiration_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group is set to expire. It is null for security groups, but for Microsoft 365 groups, it represents when the group is set to expire as defined in the groupLifecyclePolicy. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"group_types": schema.ListAttribute{
				Description: "Specifies the group type and its membership. If the collection contains Unified, the group is a Microsoft 365 group; otherwise, it's either a security group or a distribution group. For details, see groups overview.If the collection includes DynamicMembership, the group has dynamic membership; otherwise, membership is static. Returned by default. Supports $filter (eq, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"is_assignable_to_role": schema.BoolAttribute{
				Description: "Indicates whether this group can be assigned to a Microsoft Entra role. Optional. This property can only be set while creating the group and is immutable. If set to true, the securityEnabled property must also be set to true, visibility must be Hidden, and the group cannot be a dynamic group (that is, groupTypes cannot contain DynamicMembership). Only callers in Global Administrator and Privileged Role Administrator roles can set this property. The caller must also be assigned the RoleManagement.ReadWrite.Directory permission to set this property or update the membership of such groups. For more, see Using a group to manage Microsoft Entra role assignmentsUsing this feature requires a Microsoft Entra ID P1 license. Returned by default. Supports $filter (eq, ne, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"license_processing_state": schema.SingleNestedAttribute{
				Description: "Indicates the status of the group license assignment to all group members. The default value is false. Read-only. Possible values: QueuedForProcessing, ProcessingInProgress, and ProcessingComplete.Returned only on $select. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
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
			"mail": schema.StringAttribute{
				Description: "The SMTP address for the group, for example, 'serviceadmins@contoso.onmicrosoft.com'. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"mail_enabled": schema.BoolAttribute{
				Description: "Specifies whether the group is mail-enabled. Required. Returned by default. Supports $filter (eq, ne, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"mail_nickname": schema.StringAttribute{
				Description: "The mail alias for the group, unique for Microsoft 365 groups in the organization. Maximum length is 64 characters. This property can contain only characters in the ASCII character set 0 - 127 except the following: @ () / [] ' ; : <> , SPACE. Required. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"membership_rule": schema.StringAttribute{
				Description: "The rule that determines members for this group if the group is a dynamic group (groupTypes contains DynamicMembership). For more information about the syntax of the membership rule, see Membership Rules syntax. Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"membership_rule_processing_state": schema.StringAttribute{
				Description: "Indicates whether the dynamic membership processing is on or paused. Possible values are On or Paused. Returned by default. Supports $filter (eq, ne, not, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_domain_name": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "Indicates the last time at which the group was synced with the on-premises directory.The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_net_bios_name": schema.StringAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors when using Microsoft synchronization product during provisioning. Returned by default. Supports $filter (eq, not).",
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
				Description: "Contains the on-premises SAM account name synchronized from the on-premises directory. The property is only populated for customers synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith). Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_security_identifier": schema.StringAttribute{
				Description: "Contains the on-premises security identifier (SID) for the group synchronized from on-premises to the cloud. Returned by default. Supports $filter (eq including on null values). Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Description: "true if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default). Returned by default. Read-only. Supports $filter (eq, ne, not, in, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"preferred_data_location": schema.StringAttribute{
				Description: "The preferred data location for the Microsoft 365 group. By default, the group inherits the group creator's preferred data location. To set this property, the calling app must be granted the Directory.ReadWrite.All permission and the user be assigned one of the following Microsoft Entra roles:  Global Administrator  User Account Administrator Directory Writer  Exchange Administrator  SharePoint Administrator  For more information about this property, see OneDrive Online Multi-Geo. Nullable. Returned by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"preferred_language": schema.StringAttribute{
				Description: "The preferred language for a Microsoft 365 group. Should follow ISO 639-1 Code; for example, en-US. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"proxy_addresses": schema.ListAttribute{
				Description: "Email addresses for the group that direct to the same group mailbox. For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. The any operator is required to filter expressions on multi-valued properties. Returned by default. Read-only. Not nullable. Supports $filter (eq, not, ge, le, startsWith, endsWith, /$count eq 0, /$count ne 0).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"renewed_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group was last renewed. This cannot be modified directly and is only updated via the renew service action. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"security_enabled": schema.BoolAttribute{
				Description: "Specifies whether the group is a security group. Required. Returned by default. Supports $filter (eq, ne, not, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"security_identifier": schema.StringAttribute{
				Description: "Security identifier of the group, used in Windows scenarios. Returned by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"service_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors published by a federated service describing a non-transient, service-specific error regarding the properties or link from a group object .  Supports $filter (eq, not, for isResolved and serviceInstance).",
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
			"theme": schema.StringAttribute{
				Description: "Specifies a Microsoft 365 group's color theme. Possible values are Teal, Purple, Green, Blue, Pink, Orange or Red. Returned by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"visibility": schema.StringAttribute{
				Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or HiddenMembership. HiddenMembership can be set only for Microsoft 365 groups when the groups are created. It can't be updated later. Other values of visibility can be updated after group creation. If visibility value is not specified during group creation on Microsoft Graph, a security group is created as Private by default, and the Microsoft 365 group is Public. Groups assignable to roles are always Private. To learn more, see group visibility options. Returned by default. Nullable.",
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
func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var t time.Time
	var u uuid.UUID

	// Generate API request body from Plan
	requestBody := models.NewGroup()

	if !plan.Id.IsUnknown() {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	} else {
		plan.Id = types.StringNull()
	}

	if !plan.DeletedDateTime.IsUnknown() {
		planDeletedDateTime := plan.DeletedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	} else {
		plan.DeletedDateTime = types.StringNull()
	}

	if len(plan.AssignedLabels.Elements()) > 0 {
		var planAssignedLabels []models.AssignedLabelable
		for _, i := range plan.AssignedLabels.Elements() {
			assignedLabel := models.NewAssignedLabel()
			assignedLabelModel := groupAssignedLabelsModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedLabelModel)

			if !assignedLabelModel.DisplayName.IsUnknown() {
				planDisplayName := assignedLabelModel.DisplayName.ValueString()
				assignedLabel.SetDisplayName(&planDisplayName)
			} else {
				assignedLabelModel.DisplayName = types.StringNull()
			}

			if !assignedLabelModel.LabelId.IsUnknown() {
				planLabelId := assignedLabelModel.LabelId.ValueString()
				assignedLabel.SetLabelId(&planLabelId)
			} else {
				assignedLabelModel.LabelId = types.StringNull()
			}
		}
		requestBody.SetAssignedLabels(planAssignedLabels)
	} else {
		plan.AssignedLabels = types.ListNull(plan.AssignedLabels.ElementType(ctx))
	}

	if len(plan.AssignedLicenses.Elements()) > 0 {
		var planAssignedLicenses []models.AssignedLicenseable
		for _, i := range plan.AssignedLicenses.Elements() {
			assignedLicense := models.NewAssignedLicense()
			assignedLicenseModel := groupAssignedLicensesModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedLicenseModel)

			if len(assignedLicenseModel.DisabledPlans.Elements()) > 0 {
				var DisabledPlans []uuid.UUID
				for _, i := range assignedLicenseModel.DisabledPlans.Elements() {
					u, _ = uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				assignedLicense.SetDisabledPlans(DisabledPlans)
			} else {
				assignedLicenseModel.DisabledPlans = types.ListNull(types.StringType)
			}

			if !assignedLicenseModel.SkuId.IsUnknown() {
				planSkuId := assignedLicenseModel.SkuId.ValueString()
				u, _ = uuid.Parse(planSkuId)
				assignedLicense.SetSkuId(&u)
			} else {
				assignedLicenseModel.SkuId = types.StringNull()
			}
		}
		requestBody.SetAssignedLicenses(planAssignedLicenses)
	} else {
		plan.AssignedLicenses = types.ListNull(plan.AssignedLicenses.ElementType(ctx))
	}

	if !plan.Classification.IsUnknown() {
		planClassification := plan.Classification.ValueString()
		requestBody.SetClassification(&planClassification)
	} else {
		plan.Classification = types.StringNull()
	}

	if !plan.CreatedDateTime.IsUnknown() {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	} else {
		plan.CreatedDateTime = types.StringNull()
	}

	if !plan.Description.IsUnknown() {
		planDescription := plan.Description.ValueString()
		requestBody.SetDescription(&planDescription)
	} else {
		plan.Description = types.StringNull()
	}

	if !plan.DisplayName.IsUnknown() {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	} else {
		plan.DisplayName = types.StringNull()
	}

	if !plan.ExpirationDateTime.IsUnknown() {
		planExpirationDateTime := plan.ExpirationDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planExpirationDateTime)
		requestBody.SetExpirationDateTime(&t)
	} else {
		plan.ExpirationDateTime = types.StringNull()
	}

	if len(plan.GroupTypes.Elements()) > 0 {
		var groupTypes []string
		for _, i := range plan.GroupTypes.Elements() {
			groupTypes = append(groupTypes, i.String())
		}
		requestBody.SetGroupTypes(groupTypes)
	} else {
		plan.GroupTypes = types.ListNull(types.StringType)
	}

	if !plan.IsAssignableToRole.IsUnknown() {
		planIsAssignableToRole := plan.IsAssignableToRole.ValueBool()
		requestBody.SetIsAssignableToRole(&planIsAssignableToRole)
	} else {
		plan.IsAssignableToRole = types.BoolNull()
	}

	if !plan.LicenseProcessingState.IsUnknown() {
		licenseProcessingState := models.NewLicenseProcessingState()
		licenseProcessingStateModel := groupLicenseProcessingStateModel{}
		plan.LicenseProcessingState.As(ctx, &licenseProcessingStateModel, basetypes.ObjectAsOptions{})

		if !licenseProcessingStateModel.State.IsUnknown() {
			planState := licenseProcessingStateModel.State.ValueString()
			licenseProcessingState.SetState(&planState)
		} else {
			licenseProcessingStateModel.State = types.StringNull()
		}
		requestBody.SetLicenseProcessingState(licenseProcessingState)
		objectValue, _ := types.ObjectValueFrom(ctx, licenseProcessingStateModel.AttributeTypes(), licenseProcessingStateModel)
		plan.LicenseProcessingState = objectValue
	} else {
		plan.LicenseProcessingState = types.ObjectNull(plan.LicenseProcessingState.AttributeTypes(ctx))
	}

	if !plan.Mail.IsUnknown() {
		planMail := plan.Mail.ValueString()
		requestBody.SetMail(&planMail)
	} else {
		plan.Mail = types.StringNull()
	}

	if !plan.MailEnabled.IsUnknown() {
		planMailEnabled := plan.MailEnabled.ValueBool()
		requestBody.SetMailEnabled(&planMailEnabled)
	} else {
		plan.MailEnabled = types.BoolNull()
	}

	if !plan.MailNickname.IsUnknown() {
		planMailNickname := plan.MailNickname.ValueString()
		requestBody.SetMailNickname(&planMailNickname)
	} else {
		plan.MailNickname = types.StringNull()
	}

	if !plan.MembershipRule.IsUnknown() {
		planMembershipRule := plan.MembershipRule.ValueString()
		requestBody.SetMembershipRule(&planMembershipRule)
	} else {
		plan.MembershipRule = types.StringNull()
	}

	if !plan.MembershipRuleProcessingState.IsUnknown() {
		planMembershipRuleProcessingState := plan.MembershipRuleProcessingState.ValueString()
		requestBody.SetMembershipRuleProcessingState(&planMembershipRuleProcessingState)
	} else {
		plan.MembershipRuleProcessingState = types.StringNull()
	}

	if !plan.OnPremisesDomainName.IsUnknown() {
		planOnPremisesDomainName := plan.OnPremisesDomainName.ValueString()
		requestBody.SetOnPremisesDomainName(&planOnPremisesDomainName)
	} else {
		plan.OnPremisesDomainName = types.StringNull()
	}

	if !plan.OnPremisesLastSyncDateTime.IsUnknown() {
		planOnPremisesLastSyncDateTime := plan.OnPremisesLastSyncDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planOnPremisesLastSyncDateTime)
		requestBody.SetOnPremisesLastSyncDateTime(&t)
	} else {
		plan.OnPremisesLastSyncDateTime = types.StringNull()
	}

	if !plan.OnPremisesNetBiosName.IsUnknown() {
		planOnPremisesNetBiosName := plan.OnPremisesNetBiosName.ValueString()
		requestBody.SetOnPremisesNetBiosName(&planOnPremisesNetBiosName)
	} else {
		plan.OnPremisesNetBiosName = types.StringNull()
	}

	if len(plan.OnPremisesProvisioningErrors.Elements()) > 0 {
		var planOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for _, i := range plan.OnPremisesProvisioningErrors.Elements() {
			onPremisesProvisioningError := models.NewOnPremisesProvisioningError()
			onPremisesProvisioningErrorModel := groupOnPremisesProvisioningErrorsModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &onPremisesProvisioningErrorModel)

			if !onPremisesProvisioningErrorModel.Category.IsUnknown() {
				planCategory := onPremisesProvisioningErrorModel.Category.ValueString()
				onPremisesProvisioningError.SetCategory(&planCategory)
			} else {
				onPremisesProvisioningErrorModel.Category = types.StringNull()
			}

			if !onPremisesProvisioningErrorModel.OccurredDateTime.IsUnknown() {
				planOccurredDateTime := onPremisesProvisioningErrorModel.OccurredDateTime.ValueString()
				t, _ = time.Parse(time.RFC3339, planOccurredDateTime)
				onPremisesProvisioningError.SetOccurredDateTime(&t)
			} else {
				onPremisesProvisioningErrorModel.OccurredDateTime = types.StringNull()
			}

			if !onPremisesProvisioningErrorModel.PropertyCausingError.IsUnknown() {
				planPropertyCausingError := onPremisesProvisioningErrorModel.PropertyCausingError.ValueString()
				onPremisesProvisioningError.SetPropertyCausingError(&planPropertyCausingError)
			} else {
				onPremisesProvisioningErrorModel.PropertyCausingError = types.StringNull()
			}

			if !onPremisesProvisioningErrorModel.Value.IsUnknown() {
				planValue := onPremisesProvisioningErrorModel.Value.ValueString()
				onPremisesProvisioningError.SetValue(&planValue)
			} else {
				onPremisesProvisioningErrorModel.Value = types.StringNull()
			}
		}
		requestBody.SetOnPremisesProvisioningErrors(planOnPremisesProvisioningErrors)
	} else {
		plan.OnPremisesProvisioningErrors = types.ListNull(plan.OnPremisesProvisioningErrors.ElementType(ctx))
	}

	if !plan.OnPremisesSamAccountName.IsUnknown() {
		planOnPremisesSamAccountName := plan.OnPremisesSamAccountName.ValueString()
		requestBody.SetOnPremisesSamAccountName(&planOnPremisesSamAccountName)
	} else {
		plan.OnPremisesSamAccountName = types.StringNull()
	}

	if !plan.OnPremisesSecurityIdentifier.IsUnknown() {
		planOnPremisesSecurityIdentifier := plan.OnPremisesSecurityIdentifier.ValueString()
		requestBody.SetOnPremisesSecurityIdentifier(&planOnPremisesSecurityIdentifier)
	} else {
		plan.OnPremisesSecurityIdentifier = types.StringNull()
	}

	if !plan.OnPremisesSyncEnabled.IsUnknown() {
		planOnPremisesSyncEnabled := plan.OnPremisesSyncEnabled.ValueBool()
		requestBody.SetOnPremisesSyncEnabled(&planOnPremisesSyncEnabled)
	} else {
		plan.OnPremisesSyncEnabled = types.BoolNull()
	}

	if !plan.PreferredDataLocation.IsUnknown() {
		planPreferredDataLocation := plan.PreferredDataLocation.ValueString()
		requestBody.SetPreferredDataLocation(&planPreferredDataLocation)
	} else {
		plan.PreferredDataLocation = types.StringNull()
	}

	if !plan.PreferredLanguage.IsUnknown() {
		planPreferredLanguage := plan.PreferredLanguage.ValueString()
		requestBody.SetPreferredLanguage(&planPreferredLanguage)
	} else {
		plan.PreferredLanguage = types.StringNull()
	}

	if len(plan.ProxyAddresses.Elements()) > 0 {
		var proxyAddresses []string
		for _, i := range plan.ProxyAddresses.Elements() {
			proxyAddresses = append(proxyAddresses, i.String())
		}
		requestBody.SetProxyAddresses(proxyAddresses)
	} else {
		plan.ProxyAddresses = types.ListNull(types.StringType)
	}

	if !plan.RenewedDateTime.IsUnknown() {
		planRenewedDateTime := plan.RenewedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planRenewedDateTime)
		requestBody.SetRenewedDateTime(&t)
	} else {
		plan.RenewedDateTime = types.StringNull()
	}

	if !plan.SecurityEnabled.IsUnknown() {
		planSecurityEnabled := plan.SecurityEnabled.ValueBool()
		requestBody.SetSecurityEnabled(&planSecurityEnabled)
	} else {
		plan.SecurityEnabled = types.BoolNull()
	}

	if !plan.SecurityIdentifier.IsUnknown() {
		planSecurityIdentifier := plan.SecurityIdentifier.ValueString()
		requestBody.SetSecurityIdentifier(&planSecurityIdentifier)
	} else {
		plan.SecurityIdentifier = types.StringNull()
	}

	if len(plan.ServiceProvisioningErrors.Elements()) > 0 {
		var planServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for _, i := range plan.ServiceProvisioningErrors.Elements() {
			serviceProvisioningError := models.NewServiceProvisioningError()
			serviceProvisioningErrorModel := groupServiceProvisioningErrorsModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &serviceProvisioningErrorModel)

			if !serviceProvisioningErrorModel.CreatedDateTime.IsUnknown() {
				planCreatedDateTime := serviceProvisioningErrorModel.CreatedDateTime.ValueString()
				t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
				serviceProvisioningError.SetCreatedDateTime(&t)
			} else {
				serviceProvisioningErrorModel.CreatedDateTime = types.StringNull()
			}

			if !serviceProvisioningErrorModel.IsResolved.IsUnknown() {
				planIsResolved := serviceProvisioningErrorModel.IsResolved.ValueBool()
				serviceProvisioningError.SetIsResolved(&planIsResolved)
			} else {
				serviceProvisioningErrorModel.IsResolved = types.BoolNull()
			}

			if !serviceProvisioningErrorModel.ServiceInstance.IsUnknown() {
				planServiceInstance := serviceProvisioningErrorModel.ServiceInstance.ValueString()
				serviceProvisioningError.SetServiceInstance(&planServiceInstance)
			} else {
				serviceProvisioningErrorModel.ServiceInstance = types.StringNull()
			}
		}
		requestBody.SetServiceProvisioningErrors(planServiceProvisioningErrors)
	} else {
		plan.ServiceProvisioningErrors = types.ListNull(plan.ServiceProvisioningErrors.ElementType(ctx))
	}

	if !plan.Theme.IsUnknown() {
		planTheme := plan.Theme.ValueString()
		requestBody.SetTheme(&planTheme)
	} else {
		plan.Theme = types.StringNull()
	}

	if !plan.Visibility.IsUnknown() {
		planVisibility := plan.Visibility.ValueString()
		requestBody.SetVisibility(&planVisibility)
	} else {
		plan.Visibility = types.StringNull()
	}

	// Create new group
	result, err := r.client.Groups().Post(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
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
func (d *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := groups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"deletedDateTime",
				"assignedLabels",
				"assignedLicenses",
				"classification",
				"createdDateTime",
				"description",
				"displayName",
				"expirationDateTime",
				"groupTypes",
				"isAssignableToRole",
				"licenseProcessingState",
				"mail",
				"mailEnabled",
				"mailNickname",
				"membershipRule",
				"membershipRuleProcessingState",
				"onPremisesDomainName",
				"onPremisesLastSyncDateTime",
				"onPremisesNetBiosName",
				"onPremisesProvisioningErrors",
				"onPremisesSamAccountName",
				"onPremisesSecurityIdentifier",
				"onPremisesSyncEnabled",
				"preferredDataLocation",
				"preferredLanguage",
				"proxyAddresses",
				"renewedDateTime",
				"securityEnabled",
				"securityIdentifier",
				"serviceProvisioningErrors",
				"theme",
				"visibility",
				"appRoleAssignments",
				"createdOnBehalfOf",
				"memberOf",
				"members",
				"membersWithLicenseErrors",
				"owners",
				"permissionGrants",
				"settings",
				"transitiveMemberOf",
				"transitiveMembers",
				"acceptedSenders",
				"calendar",
				"calendarView",
				"conversations",
				"events",
				"rejectedSenders",
				"threads",
				"drive",
				"drives",
				"sites",
				"extensions",
				"groupLifecyclePolicies",
				"planner",
				"onenote",
				"photo",
				"photos",
				"team",
			},
		},
	}

	var result models.Groupable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.Groups().ByGroupId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting group",
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
	if len(result.GetAssignedLabels()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedLabels() {
			assignedLabels := new(groupAssignedLabelsModel)

			if v.GetDisplayName() != nil {
				assignedLabels.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				assignedLabels.DisplayName = types.StringNull()
			}
			if v.GetLabelId() != nil {
				assignedLabels.LabelId = types.StringValue(*v.GetLabelId())
			} else {
				assignedLabels.LabelId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, assignedLabels.AttributeTypes(), assignedLabels)
			objectValues = append(objectValues, objectValue)
		}
		state.AssignedLabels, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(result.GetAssignedLicenses()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedLicenses() {
			assignedLicenses := new(groupAssignedLicensesModel)

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
	if result.GetClassification() != nil {
		state.Classification = types.StringValue(*result.GetClassification())
	} else {
		state.Classification = types.StringNull()
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		state.CreatedDateTime = types.StringNull()
	}
	if result.GetDescription() != nil {
		state.Description = types.StringValue(*result.GetDescription())
	} else {
		state.Description = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		state.DisplayName = types.StringNull()
	}
	if result.GetExpirationDateTime() != nil {
		state.ExpirationDateTime = types.StringValue(result.GetExpirationDateTime().String())
	} else {
		state.ExpirationDateTime = types.StringNull()
	}
	if len(result.GetGroupTypes()) > 0 {
		var groupTypes []attr.Value
		for _, v := range result.GetGroupTypes() {
			groupTypes = append(groupTypes, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, groupTypes)
		state.GroupTypes = listValue
	} else {
		state.GroupTypes = types.ListNull(types.StringType)
	}
	if result.GetIsAssignableToRole() != nil {
		state.IsAssignableToRole = types.BoolValue(*result.GetIsAssignableToRole())
	} else {
		state.IsAssignableToRole = types.BoolNull()
	}
	if result.GetLicenseProcessingState() != nil {
		licenseProcessingState := new(groupLicenseProcessingStateModel)

		if result.GetLicenseProcessingState().GetState() != nil {
			licenseProcessingState.State = types.StringValue(*result.GetLicenseProcessingState().GetState())
		} else {
			licenseProcessingState.State = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, licenseProcessingState.AttributeTypes(), licenseProcessingState)
		state.LicenseProcessingState = objectValue
	}
	if result.GetMail() != nil {
		state.Mail = types.StringValue(*result.GetMail())
	} else {
		state.Mail = types.StringNull()
	}
	if result.GetMailEnabled() != nil {
		state.MailEnabled = types.BoolValue(*result.GetMailEnabled())
	} else {
		state.MailEnabled = types.BoolNull()
	}
	if result.GetMailNickname() != nil {
		state.MailNickname = types.StringValue(*result.GetMailNickname())
	} else {
		state.MailNickname = types.StringNull()
	}
	if result.GetMembershipRule() != nil {
		state.MembershipRule = types.StringValue(*result.GetMembershipRule())
	} else {
		state.MembershipRule = types.StringNull()
	}
	if result.GetMembershipRuleProcessingState() != nil {
		state.MembershipRuleProcessingState = types.StringValue(*result.GetMembershipRuleProcessingState())
	} else {
		state.MembershipRuleProcessingState = types.StringNull()
	}
	if result.GetOnPremisesDomainName() != nil {
		state.OnPremisesDomainName = types.StringValue(*result.GetOnPremisesDomainName())
	} else {
		state.OnPremisesDomainName = types.StringNull()
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	} else {
		state.OnPremisesLastSyncDateTime = types.StringNull()
	}
	if result.GetOnPremisesNetBiosName() != nil {
		state.OnPremisesNetBiosName = types.StringValue(*result.GetOnPremisesNetBiosName())
	} else {
		state.OnPremisesNetBiosName = types.StringNull()
	}
	if len(result.GetOnPremisesProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetOnPremisesProvisioningErrors() {
			onPremisesProvisioningErrors := new(groupOnPremisesProvisioningErrorsModel)

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
	if result.GetRenewedDateTime() != nil {
		state.RenewedDateTime = types.StringValue(result.GetRenewedDateTime().String())
	} else {
		state.RenewedDateTime = types.StringNull()
	}
	if result.GetSecurityEnabled() != nil {
		state.SecurityEnabled = types.BoolValue(*result.GetSecurityEnabled())
	} else {
		state.SecurityEnabled = types.BoolNull()
	}
	if result.GetSecurityIdentifier() != nil {
		state.SecurityIdentifier = types.StringValue(*result.GetSecurityIdentifier())
	} else {
		state.SecurityIdentifier = types.StringNull()
	}
	if len(result.GetServiceProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetServiceProvisioningErrors() {
			serviceProvisioningErrors := new(groupServiceProvisioningErrorsModel)

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
	if result.GetTheme() != nil {
		state.Theme = types.StringValue(*result.GetTheme())
	} else {
		state.Theme = types.StringNull()
	}
	if result.GetVisibility() != nil {
		state.Visibility = types.StringValue(*result.GetVisibility())
	} else {
		state.Visibility = types.StringNull()
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan groupModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state groupModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewGroup()
	var t time.Time
	var u uuid.UUID

	if !plan.Id.Equal(state.Id) {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	}

	if !plan.DeletedDateTime.Equal(state.DeletedDateTime) {
		planDeletedDateTime := plan.DeletedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !plan.AssignedLabels.Equal(state.AssignedLabels) {
		var planAssignedLabels []models.AssignedLabelable
		for k, i := range plan.AssignedLabels.Elements() {
			assignedLabel := models.NewAssignedLabel()
			assignedLabelModel := groupAssignedLabelsModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedLabelModel)
			assignedLabelState := groupAssignedLabelsModel{}
			types.ListValueFrom(ctx, state.AssignedLabels.Elements()[k].Type(ctx), &assignedLabelModel)

			if !assignedLabelModel.DisplayName.Equal(assignedLabelState.DisplayName) {
				planDisplayName := assignedLabelModel.DisplayName.ValueString()
				assignedLabel.SetDisplayName(&planDisplayName)
			}

			if !assignedLabelModel.LabelId.Equal(assignedLabelState.LabelId) {
				planLabelId := assignedLabelModel.LabelId.ValueString()
				assignedLabel.SetLabelId(&planLabelId)
			}
		}
		requestBody.SetAssignedLabels(planAssignedLabels)
	}

	if !plan.AssignedLicenses.Equal(state.AssignedLicenses) {
		var planAssignedLicenses []models.AssignedLicenseable
		for k, i := range plan.AssignedLicenses.Elements() {
			assignedLicense := models.NewAssignedLicense()
			assignedLicenseModel := groupAssignedLicensesModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedLicenseModel)
			assignedLicenseState := groupAssignedLicensesModel{}
			types.ListValueFrom(ctx, state.AssignedLicenses.Elements()[k].Type(ctx), &assignedLicenseModel)

			if !assignedLicenseModel.DisabledPlans.Equal(assignedLicenseState.DisabledPlans) {
				var DisabledPlans []uuid.UUID
				for _, i := range assignedLicenseModel.DisabledPlans.Elements() {
					u, _ = uuid.Parse(i.String())
					DisabledPlans = append(DisabledPlans, u)
				}
				assignedLicense.SetDisabledPlans(DisabledPlans)
			}

			if !assignedLicenseModel.SkuId.Equal(assignedLicenseState.SkuId) {
				planSkuId := assignedLicenseModel.SkuId.ValueString()
				u, _ = uuid.Parse(planSkuId)
				assignedLicense.SetSkuId(&u)
			}
		}
		requestBody.SetAssignedLicenses(planAssignedLicenses)
	}

	if !plan.Classification.Equal(state.Classification) {
		planClassification := plan.Classification.ValueString()
		requestBody.SetClassification(&planClassification)
	}

	if !plan.CreatedDateTime.Equal(state.CreatedDateTime) {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	}

	if !plan.Description.Equal(state.Description) {
		planDescription := plan.Description.ValueString()
		requestBody.SetDescription(&planDescription)
	}

	if !plan.DisplayName.Equal(state.DisplayName) {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	}

	if !plan.ExpirationDateTime.Equal(state.ExpirationDateTime) {
		planExpirationDateTime := plan.ExpirationDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planExpirationDateTime)
		requestBody.SetExpirationDateTime(&t)
	}

	if !plan.GroupTypes.Equal(state.GroupTypes) {
		var groupTypes []string
		for _, i := range plan.GroupTypes.Elements() {
			groupTypes = append(groupTypes, i.String())
		}
		requestBody.SetGroupTypes(groupTypes)
	}

	if !plan.IsAssignableToRole.Equal(state.IsAssignableToRole) {
		planIsAssignableToRole := plan.IsAssignableToRole.ValueBool()
		requestBody.SetIsAssignableToRole(&planIsAssignableToRole)
	}

	if !plan.LicenseProcessingState.Equal(state.LicenseProcessingState) {
		licenseProcessingState := models.NewLicenseProcessingState()
		licenseProcessingStateModel := groupLicenseProcessingStateModel{}
		plan.LicenseProcessingState.As(ctx, &licenseProcessingStateModel, basetypes.ObjectAsOptions{})
		licenseProcessingStateState := groupLicenseProcessingStateModel{}
		state.LicenseProcessingState.As(ctx, &licenseProcessingStateState, basetypes.ObjectAsOptions{})

		if !licenseProcessingStateModel.State.Equal(licenseProcessingStateState.State) {
			planState := licenseProcessingStateModel.State.ValueString()
			licenseProcessingState.SetState(&planState)
		}
		requestBody.SetLicenseProcessingState(licenseProcessingState)
		objectValue, _ := types.ObjectValueFrom(ctx, licenseProcessingStateModel.AttributeTypes(), licenseProcessingStateModel)
		plan.LicenseProcessingState = objectValue
	}

	if !plan.Mail.Equal(state.Mail) {
		planMail := plan.Mail.ValueString()
		requestBody.SetMail(&planMail)
	}

	if !plan.MailEnabled.Equal(state.MailEnabled) {
		planMailEnabled := plan.MailEnabled.ValueBool()
		requestBody.SetMailEnabled(&planMailEnabled)
	}

	if !plan.MailNickname.Equal(state.MailNickname) {
		planMailNickname := plan.MailNickname.ValueString()
		requestBody.SetMailNickname(&planMailNickname)
	}

	if !plan.MembershipRule.Equal(state.MembershipRule) {
		planMembershipRule := plan.MembershipRule.ValueString()
		requestBody.SetMembershipRule(&planMembershipRule)
	}

	if !plan.MembershipRuleProcessingState.Equal(state.MembershipRuleProcessingState) {
		planMembershipRuleProcessingState := plan.MembershipRuleProcessingState.ValueString()
		requestBody.SetMembershipRuleProcessingState(&planMembershipRuleProcessingState)
	}

	if !plan.OnPremisesDomainName.Equal(state.OnPremisesDomainName) {
		planOnPremisesDomainName := plan.OnPremisesDomainName.ValueString()
		requestBody.SetOnPremisesDomainName(&planOnPremisesDomainName)
	}

	if !plan.OnPremisesLastSyncDateTime.Equal(state.OnPremisesLastSyncDateTime) {
		planOnPremisesLastSyncDateTime := plan.OnPremisesLastSyncDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planOnPremisesLastSyncDateTime)
		requestBody.SetOnPremisesLastSyncDateTime(&t)
	}

	if !plan.OnPremisesNetBiosName.Equal(state.OnPremisesNetBiosName) {
		planOnPremisesNetBiosName := plan.OnPremisesNetBiosName.ValueString()
		requestBody.SetOnPremisesNetBiosName(&planOnPremisesNetBiosName)
	}

	if !plan.OnPremisesProvisioningErrors.Equal(state.OnPremisesProvisioningErrors) {
		var planOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for k, i := range plan.OnPremisesProvisioningErrors.Elements() {
			onPremisesProvisioningError := models.NewOnPremisesProvisioningError()
			onPremisesProvisioningErrorModel := groupOnPremisesProvisioningErrorsModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &onPremisesProvisioningErrorModel)
			onPremisesProvisioningErrorState := groupOnPremisesProvisioningErrorsModel{}
			types.ListValueFrom(ctx, state.OnPremisesProvisioningErrors.Elements()[k].Type(ctx), &onPremisesProvisioningErrorModel)

			if !onPremisesProvisioningErrorModel.Category.Equal(onPremisesProvisioningErrorState.Category) {
				planCategory := onPremisesProvisioningErrorModel.Category.ValueString()
				onPremisesProvisioningError.SetCategory(&planCategory)
			}

			if !onPremisesProvisioningErrorModel.OccurredDateTime.Equal(onPremisesProvisioningErrorState.OccurredDateTime) {
				planOccurredDateTime := onPremisesProvisioningErrorModel.OccurredDateTime.ValueString()
				t, _ = time.Parse(time.RFC3339, planOccurredDateTime)
				onPremisesProvisioningError.SetOccurredDateTime(&t)
			}

			if !onPremisesProvisioningErrorModel.PropertyCausingError.Equal(onPremisesProvisioningErrorState.PropertyCausingError) {
				planPropertyCausingError := onPremisesProvisioningErrorModel.PropertyCausingError.ValueString()
				onPremisesProvisioningError.SetPropertyCausingError(&planPropertyCausingError)
			}

			if !onPremisesProvisioningErrorModel.Value.Equal(onPremisesProvisioningErrorState.Value) {
				planValue := onPremisesProvisioningErrorModel.Value.ValueString()
				onPremisesProvisioningError.SetValue(&planValue)
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

	if !plan.PreferredDataLocation.Equal(state.PreferredDataLocation) {
		planPreferredDataLocation := plan.PreferredDataLocation.ValueString()
		requestBody.SetPreferredDataLocation(&planPreferredDataLocation)
	}

	if !plan.PreferredLanguage.Equal(state.PreferredLanguage) {
		planPreferredLanguage := plan.PreferredLanguage.ValueString()
		requestBody.SetPreferredLanguage(&planPreferredLanguage)
	}

	if !plan.ProxyAddresses.Equal(state.ProxyAddresses) {
		var proxyAddresses []string
		for _, i := range plan.ProxyAddresses.Elements() {
			proxyAddresses = append(proxyAddresses, i.String())
		}
		requestBody.SetProxyAddresses(proxyAddresses)
	}

	if !plan.RenewedDateTime.Equal(state.RenewedDateTime) {
		planRenewedDateTime := plan.RenewedDateTime.ValueString()
		t, _ = time.Parse(time.RFC3339, planRenewedDateTime)
		requestBody.SetRenewedDateTime(&t)
	}

	if !plan.SecurityEnabled.Equal(state.SecurityEnabled) {
		planSecurityEnabled := plan.SecurityEnabled.ValueBool()
		requestBody.SetSecurityEnabled(&planSecurityEnabled)
	}

	if !plan.SecurityIdentifier.Equal(state.SecurityIdentifier) {
		planSecurityIdentifier := plan.SecurityIdentifier.ValueString()
		requestBody.SetSecurityIdentifier(&planSecurityIdentifier)
	}

	if !plan.ServiceProvisioningErrors.Equal(state.ServiceProvisioningErrors) {
		var planServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for k, i := range plan.ServiceProvisioningErrors.Elements() {
			serviceProvisioningError := models.NewServiceProvisioningError()
			serviceProvisioningErrorModel := groupServiceProvisioningErrorsModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &serviceProvisioningErrorModel)
			serviceProvisioningErrorState := groupServiceProvisioningErrorsModel{}
			types.ListValueFrom(ctx, state.ServiceProvisioningErrors.Elements()[k].Type(ctx), &serviceProvisioningErrorModel)

			if !serviceProvisioningErrorModel.CreatedDateTime.Equal(serviceProvisioningErrorState.CreatedDateTime) {
				planCreatedDateTime := serviceProvisioningErrorModel.CreatedDateTime.ValueString()
				t, _ = time.Parse(time.RFC3339, planCreatedDateTime)
				serviceProvisioningError.SetCreatedDateTime(&t)
			}

			if !serviceProvisioningErrorModel.IsResolved.Equal(serviceProvisioningErrorState.IsResolved) {
				planIsResolved := serviceProvisioningErrorModel.IsResolved.ValueBool()
				serviceProvisioningError.SetIsResolved(&planIsResolved)
			}

			if !serviceProvisioningErrorModel.ServiceInstance.Equal(serviceProvisioningErrorState.ServiceInstance) {
				planServiceInstance := serviceProvisioningErrorModel.ServiceInstance.ValueString()
				serviceProvisioningError.SetServiceInstance(&planServiceInstance)
			}
		}
		requestBody.SetServiceProvisioningErrors(planServiceProvisioningErrors)
	}

	if !plan.Theme.Equal(state.Theme) {
		planTheme := plan.Theme.ValueString()
		requestBody.SetTheme(&planTheme)
	}

	if !plan.Visibility.Equal(state.Visibility) {
		planVisibility := plan.Visibility.ValueString()
		requestBody.SetVisibility(&planVisibility)
	}

	// Update group
	_, err := r.client.Groups().ByGroupId(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating group",
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
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groupModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete group
	err := r.client.Groups().ByGroupId(state.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting group",
			err.Error(),
		)
		return
	}

}
