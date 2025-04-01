package groups

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
				Description: "The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group. Returned only on $select. This property can be updated only in delegated scenarios where the caller requires both the Microsoft Graph permission and a supported administrator role.",
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
			"classification": schema.StringAttribute{
				Description: "Describes a classification for the group (such as low, medium, or high business impact). Valid values for this property are defined by creating a ClassificationList setting value, based on the template definition.Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"created_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group was created. The value can't be modified and is automatically populated when the group is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only.",
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
				Description: "The display name for the group. This property is required when a group is created and can't be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"expiration_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group is set to expire. It's null for security groups, but for Microsoft 365 groups, it represents when the group is set to expire as defined in the groupLifecyclePolicy. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
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
				Description: "Indicates whether this group can be assigned to a Microsoft Entra role. Optional. This property can only be set while creating the group and is immutable. If set to true, the securityEnabled property must also be set to true, visibility must be Hidden, and the group can't be a dynamic group (that is, groupTypes can't contain DynamicMembership). Only callers with at least the Privileged Role Administrator role can set this property. The caller must also be assigned the RoleManagement.ReadWrite.Directory permission to set this property or update the membership of such groups. For more, see Using a group to manage Microsoft Entra role assignmentsUsing this feature requires a Microsoft Entra ID P1 license. Returned by default. Supports $filter (eq, ne, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_management_restricted": schema.BoolAttribute{
				Description: "",
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
				Description: "The SMTP address for the group, for example, 'serviceadmins@contoso.com'. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
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
				Description: "The mail alias for the group, unique for Microsoft 365 groups in the organization. Maximum length is 64 characters. This property can contain only characters in the ASCII character set 0 - 127 except the following characters: @ () / [] ' ; : <> , SPACE. Required. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
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
				Description: "Contains the on-premises domain FQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.Returned by default. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "Indicates the last time at which the group was synced with the on-premises directory. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_net_bios_name": schema.StringAttribute{
				Description: "Contains the on-premises netBios name synchronized from the on-premises directory. The property is only populated for customers synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.Returned by default. Read-only.",
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
				Description: "Contains the on-premises security identifier (SID) for the group synchronized from on-premises to the cloud. Read-only. Returned by default. Supports $filter (eq including on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Description: "true if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never synced from an on-premises directory (default). Returned by default. Read-only. Supports $filter (eq, ne, not, in, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"preferred_data_location": schema.StringAttribute{
				Description: "The preferred data location for the Microsoft 365 group. By default, the group inherits the group creator's preferred data location. To set this property, the calling app must be granted the Directory.ReadWrite.All permission and the user be assigned at least one of the following Microsoft Entra roles: User Account Administrator Directory Writer  Exchange Administrator  SharePoint Administrator  For more information about this property, see OneDrive Online Multi-Geo. Nullable. Returned by default.",
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
				Description: "Timestamp of when the group was last renewed. This value can't be modified directly and is only updated via the renew service action. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
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
				Description: "Security identifier of the group, used in Windows scenarios. Read-only. Returned by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"service_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors published by a federated service describing a nontransient, service-specific error regarding the properties or link from a group object.  Supports $filter (eq, not, for isResolved and serviceInstance).",
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
				Description: "Specifies a Microsoft 365 group's color theme. Possible values are Teal, Purple, Green, Blue, Pink, Orange, or Red. Returned by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"unique_name": schema.StringAttribute{
				Description: "The unique identifier that can be assigned to a group and used as an alternate key. Immutable. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"visibility": schema.StringAttribute{
				Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or HiddenMembership. HiddenMembership can be set only for Microsoft 365 groups when the groups are created. It can't be updated later. Other values of visibility can be updated after group creation. If visibility value isn't specified during group creation on Microsoft Graph, a security group is created as Private by default, and the Microsoft 365 group is Public. Groups assignable to roles are always Private. To learn more, see group visibility options. Returned by default. Nullable.",
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
	// Retrieve values from Terraform plan
	var tfPlanGroup groupModel
	diags := req.Plan.Get(ctx, &tfPlanGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	requestBodyGroup := models.NewGroup()
	// START Id | CreateStringAttribute
	if !tfPlanGroup.Id.IsUnknown() {
		tfPlanId := tfPlanGroup.Id.ValueString()
		requestBodyGroup.SetId(&tfPlanId)
	} else {
		tfPlanGroup.Id = types.StringNull()
	}
	// END Id | CreateStringAttribute

	// START DeletedDateTime | CreateStringTimeAttribute
	if !tfPlanGroup.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlanGroup.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyGroup.SetDeletedDateTime(&t)
	} else {
		tfPlanGroup.DeletedDateTime = types.StringNull()
	}
	// END DeletedDateTime | CreateStringTimeAttribute

	// START AssignedLabels | CreateArrayObjectAttribute
	if len(tfPlanGroup.AssignedLabels.Elements()) > 0 {
		var requestBodyAssignedLabels []models.AssignedLabelable
		for _, i := range tfPlanGroup.AssignedLabels.Elements() {
			requestBodyAssignedLabel := models.NewAssignedLabel()
			tfPlanAssignedLabels := groupAssignedLabelModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLabels)

			// START DisplayName | CreateStringAttribute
			if !tfPlanAssignedLabels.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanAssignedLabels.DisplayName.ValueString()
				requestBodyAssignedLabels.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanAssignedLabels.DisplayName = types.StringNull()
			}
			// END DisplayName | CreateStringAttribute

			// START LabelId | CreateStringAttribute
			if !tfPlanAssignedLabels.LabelId.IsUnknown() {
				tfPlanLabelId := tfPlanAssignedLabels.LabelId.ValueString()
				requestBodyAssignedLabels.SetLabelId(&tfPlanLabelId)
			} else {
				tfPlanAssignedLabels.LabelId = types.StringNull()
			}
			// END LabelId | CreateStringAttribute

		}
		requestBodyGroup.SetAssignedLabels(requestBodyAssignedLabels)
	} else {
		tfPlanGroup.AssignedLabels = types.ListNull(tfPlanGroup.AssignedLabels.ElementType(ctx))
	}
	// END AssignedLabels | CreateArrayObjectAttribute

	// START AssignedLicenses | CreateArrayObjectAttribute
	if len(tfPlanGroup.AssignedLicenses.Elements()) > 0 {
		var requestBodyAssignedLicenses []models.AssignedLicenseable
		for _, i := range tfPlanGroup.AssignedLicenses.Elements() {
			requestBodyAssignedLicense := models.NewAssignedLicense()
			tfPlanAssignedLicenses := groupAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLicenses)

			// START DisabledPlans | CreateArrayUuidAttribute
			if len(tfPlanAssignedLicenses.DisabledPlans.Elements()) > 0 {
				var uuidArrayDisabledPlans []uuid.UUID
				for _, i := range tfPlanAssignedLicenses.DisabledPlans.Elements() {
					u, _ := uuid.Parse(i.String())
					uuidArrayDisabledPlans = append(uuidArrayDisabledPlans, u)
				}
				requestBodyAssignedLicenses.SetDisabledPlans(uuidArrayDisabledPlans)
			} else {
				tfPlanAssignedLicenses.DisabledPlans = types.ListNull(types.StringType)
			}

			// END DisabledPlans | CreateArrayUuidAttribute

			// START SkuId | CreateStringUuidAttribute
			if !tfPlanAssignedLicenses.SkuId.IsUnknown() {
				tfPlanSkuId := tfPlanAssignedLicenses.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				requestBodyAssignedLicenses.SetSkuId(&u)
			} else {
				tfPlanAssignedLicenses.SkuId = types.StringNull()
			}
			// END SkuId | CreateStringUuidAttribute

		}
		requestBodyGroup.SetAssignedLicenses(requestBodyAssignedLicenses)
	} else {
		tfPlanGroup.AssignedLicenses = types.ListNull(tfPlanGroup.AssignedLicenses.ElementType(ctx))
	}
	// END AssignedLicenses | CreateArrayObjectAttribute

	// START Classification | CreateStringAttribute
	if !tfPlanGroup.Classification.IsUnknown() {
		tfPlanClassification := tfPlanGroup.Classification.ValueString()
		requestBodyGroup.SetClassification(&tfPlanClassification)
	} else {
		tfPlanGroup.Classification = types.StringNull()
	}
	// END Classification | CreateStringAttribute

	// START CreatedDateTime | CreateStringTimeAttribute
	if !tfPlanGroup.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlanGroup.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyGroup.SetCreatedDateTime(&t)
	} else {
		tfPlanGroup.CreatedDateTime = types.StringNull()
	}
	// END CreatedDateTime | CreateStringTimeAttribute

	// START Description | CreateStringAttribute
	if !tfPlanGroup.Description.IsUnknown() {
		tfPlanDescription := tfPlanGroup.Description.ValueString()
		requestBodyGroup.SetDescription(&tfPlanDescription)
	} else {
		tfPlanGroup.Description = types.StringNull()
	}
	// END Description | CreateStringAttribute

	// START DisplayName | CreateStringAttribute
	if !tfPlanGroup.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanGroup.DisplayName.ValueString()
		requestBodyGroup.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanGroup.DisplayName = types.StringNull()
	}
	// END DisplayName | CreateStringAttribute

	// START ExpirationDateTime | CreateStringTimeAttribute
	if !tfPlanGroup.ExpirationDateTime.IsUnknown() {
		tfPlanExpirationDateTime := tfPlanGroup.ExpirationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanExpirationDateTime)
		requestBodyGroup.SetExpirationDateTime(&t)
	} else {
		tfPlanGroup.ExpirationDateTime = types.StringNull()
	}
	// END ExpirationDateTime | CreateStringTimeAttribute

	// START GroupTypes | CreateArrayStringAttribute
	if len(tfPlanGroup.GroupTypes.Elements()) > 0 {
		var stringArrayGroupTypes []string
		for _, i := range tfPlanGroup.GroupTypes.Elements() {
			stringArrayGroupTypes = append(stringArrayGroupTypes, i.String())
		}
		requestBodyGroup.SetGroupTypes(stringArrayGroupTypes)
	} else {
		tfPlanGroup.GroupTypes = types.ListNull(types.StringType)
	}
	// END GroupTypes | CreateArrayStringAttribute

	// START IsAssignableToRole | CreateBoolAttribute
	if !tfPlanGroup.IsAssignableToRole.IsUnknown() {
		tfPlanIsAssignableToRole := tfPlanGroup.IsAssignableToRole.ValueBool()
		requestBodyGroup.SetIsAssignableToRole(&tfPlanIsAssignableToRole)
	} else {
		tfPlanGroup.IsAssignableToRole = types.BoolNull()
	}
	// END IsAssignableToRole | CreateBoolAttribute

	// START IsManagementRestricted | CreateBoolAttribute
	if !tfPlanGroup.IsManagementRestricted.IsUnknown() {
		tfPlanIsManagementRestricted := tfPlanGroup.IsManagementRestricted.ValueBool()
		requestBodyGroup.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	} else {
		tfPlanGroup.IsManagementRestricted = types.BoolNull()
	}
	// END IsManagementRestricted | CreateBoolAttribute

	// START LicenseProcessingState | CreateObjectAttribute
	if !tfPlanGroup.LicenseProcessingState.IsUnknown() {
		requestBodyLicenseProcessingState := models.NewLicenseProcessingState()
		tfPlanLicenseProcessingState := groupLicenseProcessingStateModel{}
		tfPlanGroup.LicenseProcessingState.As(ctx, &tfPlanLicenseProcessingState, basetypes.ObjectAsOptions{})

		// START State | CreateStringAttribute
		if !tfPlanLicenseProcessingState.State.IsUnknown() {
			tfPlanState := tfPlanLicenseProcessingState.State.ValueString()
			requestBodyLicenseProcessingState.SetState(&tfPlanState)
		} else {
			tfPlanLicenseProcessingState.State = types.StringNull()
		}
		// END State | CreateStringAttribute

		requestBodyGroup.SetLicenseProcessingState(requestBodyLicenseProcessingState)
		tfPlanGroup.LicenseProcessingState, _ = types.ObjectValueFrom(ctx, tfPlanLicenseProcessingState.AttributeTypes(), requestBodyLicenseProcessingState)
	} else {
		tfPlanGroup.LicenseProcessingState = types.ObjectNull(tfPlanGroup.LicenseProcessingState.AttributeTypes(ctx))
	}
	// END LicenseProcessingState | CreateObjectAttribute

	// START Mail | CreateStringAttribute
	if !tfPlanGroup.Mail.IsUnknown() {
		tfPlanMail := tfPlanGroup.Mail.ValueString()
		requestBodyGroup.SetMail(&tfPlanMail)
	} else {
		tfPlanGroup.Mail = types.StringNull()
	}
	// END Mail | CreateStringAttribute

	// START MailEnabled | CreateBoolAttribute
	if !tfPlanGroup.MailEnabled.IsUnknown() {
		tfPlanMailEnabled := tfPlanGroup.MailEnabled.ValueBool()
		requestBodyGroup.SetMailEnabled(&tfPlanMailEnabled)
	} else {
		tfPlanGroup.MailEnabled = types.BoolNull()
	}
	// END MailEnabled | CreateBoolAttribute

	// START MailNickname | CreateStringAttribute
	if !tfPlanGroup.MailNickname.IsUnknown() {
		tfPlanMailNickname := tfPlanGroup.MailNickname.ValueString()
		requestBodyGroup.SetMailNickname(&tfPlanMailNickname)
	} else {
		tfPlanGroup.MailNickname = types.StringNull()
	}
	// END MailNickname | CreateStringAttribute

	// START MembershipRule | CreateStringAttribute
	if !tfPlanGroup.MembershipRule.IsUnknown() {
		tfPlanMembershipRule := tfPlanGroup.MembershipRule.ValueString()
		requestBodyGroup.SetMembershipRule(&tfPlanMembershipRule)
	} else {
		tfPlanGroup.MembershipRule = types.StringNull()
	}
	// END MembershipRule | CreateStringAttribute

	// START MembershipRuleProcessingState | CreateStringAttribute
	if !tfPlanGroup.MembershipRuleProcessingState.IsUnknown() {
		tfPlanMembershipRuleProcessingState := tfPlanGroup.MembershipRuleProcessingState.ValueString()
		requestBodyGroup.SetMembershipRuleProcessingState(&tfPlanMembershipRuleProcessingState)
	} else {
		tfPlanGroup.MembershipRuleProcessingState = types.StringNull()
	}
	// END MembershipRuleProcessingState | CreateStringAttribute

	// START OnPremisesDomainName | CreateStringAttribute
	if !tfPlanGroup.OnPremisesDomainName.IsUnknown() {
		tfPlanOnPremisesDomainName := tfPlanGroup.OnPremisesDomainName.ValueString()
		requestBodyGroup.SetOnPremisesDomainName(&tfPlanOnPremisesDomainName)
	} else {
		tfPlanGroup.OnPremisesDomainName = types.StringNull()
	}
	// END OnPremisesDomainName | CreateStringAttribute

	// START OnPremisesLastSyncDateTime | CreateStringTimeAttribute
	if !tfPlanGroup.OnPremisesLastSyncDateTime.IsUnknown() {
		tfPlanOnPremisesLastSyncDateTime := tfPlanGroup.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		requestBodyGroup.SetOnPremisesLastSyncDateTime(&t)
	} else {
		tfPlanGroup.OnPremisesLastSyncDateTime = types.StringNull()
	}
	// END OnPremisesLastSyncDateTime | CreateStringTimeAttribute

	// START OnPremisesNetBiosName | CreateStringAttribute
	if !tfPlanGroup.OnPremisesNetBiosName.IsUnknown() {
		tfPlanOnPremisesNetBiosName := tfPlanGroup.OnPremisesNetBiosName.ValueString()
		requestBodyGroup.SetOnPremisesNetBiosName(&tfPlanOnPremisesNetBiosName)
	} else {
		tfPlanGroup.OnPremisesNetBiosName = types.StringNull()
	}
	// END OnPremisesNetBiosName | CreateStringAttribute

	// START OnPremisesProvisioningErrors | CreateArrayObjectAttribute
	if len(tfPlanGroup.OnPremisesProvisioningErrors.Elements()) > 0 {
		var requestBodyOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for _, i := range tfPlanGroup.OnPremisesProvisioningErrors.Elements() {
			requestBodyOnPremisesProvisioningError := models.NewOnPremisesProvisioningError()
			tfPlanOnPremisesProvisioningErrors := groupOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOnPremisesProvisioningErrors)

			// START Category | CreateStringAttribute
			if !tfPlanOnPremisesProvisioningErrors.Category.IsUnknown() {
				tfPlanCategory := tfPlanOnPremisesProvisioningErrors.Category.ValueString()
				requestBodyOnPremisesProvisioningErrors.SetCategory(&tfPlanCategory)
			} else {
				tfPlanOnPremisesProvisioningErrors.Category = types.StringNull()
			}
			// END Category | CreateStringAttribute

			// START OccurredDateTime | CreateStringTimeAttribute
			if !tfPlanOnPremisesProvisioningErrors.OccurredDateTime.IsUnknown() {
				tfPlanOccurredDateTime := tfPlanOnPremisesProvisioningErrors.OccurredDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanOccurredDateTime)
				requestBodyOnPremisesProvisioningErrors.SetOccurredDateTime(&t)
			} else {
				tfPlanOnPremisesProvisioningErrors.OccurredDateTime = types.StringNull()
			}
			// END OccurredDateTime | CreateStringTimeAttribute

			// START PropertyCausingError | CreateStringAttribute
			if !tfPlanOnPremisesProvisioningErrors.PropertyCausingError.IsUnknown() {
				tfPlanPropertyCausingError := tfPlanOnPremisesProvisioningErrors.PropertyCausingError.ValueString()
				requestBodyOnPremisesProvisioningErrors.SetPropertyCausingError(&tfPlanPropertyCausingError)
			} else {
				tfPlanOnPremisesProvisioningErrors.PropertyCausingError = types.StringNull()
			}
			// END PropertyCausingError | CreateStringAttribute

			// START Value | CreateStringAttribute
			if !tfPlanOnPremisesProvisioningErrors.Value.IsUnknown() {
				tfPlanValue := tfPlanOnPremisesProvisioningErrors.Value.ValueString()
				requestBodyOnPremisesProvisioningErrors.SetValue(&tfPlanValue)
			} else {
				tfPlanOnPremisesProvisioningErrors.Value = types.StringNull()
			}
			// END Value | CreateStringAttribute

		}
		requestBodyGroup.SetOnPremisesProvisioningErrors(requestBodyOnPremisesProvisioningErrors)
	} else {
		tfPlanGroup.OnPremisesProvisioningErrors = types.ListNull(tfPlanGroup.OnPremisesProvisioningErrors.ElementType(ctx))
	}
	// END OnPremisesProvisioningErrors | CreateArrayObjectAttribute

	// START OnPremisesSamAccountName | CreateStringAttribute
	if !tfPlanGroup.OnPremisesSamAccountName.IsUnknown() {
		tfPlanOnPremisesSamAccountName := tfPlanGroup.OnPremisesSamAccountName.ValueString()
		requestBodyGroup.SetOnPremisesSamAccountName(&tfPlanOnPremisesSamAccountName)
	} else {
		tfPlanGroup.OnPremisesSamAccountName = types.StringNull()
	}
	// END OnPremisesSamAccountName | CreateStringAttribute

	// START OnPremisesSecurityIdentifier | CreateStringAttribute
	if !tfPlanGroup.OnPremisesSecurityIdentifier.IsUnknown() {
		tfPlanOnPremisesSecurityIdentifier := tfPlanGroup.OnPremisesSecurityIdentifier.ValueString()
		requestBodyGroup.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	} else {
		tfPlanGroup.OnPremisesSecurityIdentifier = types.StringNull()
	}
	// END OnPremisesSecurityIdentifier | CreateStringAttribute

	// START OnPremisesSyncEnabled | CreateBoolAttribute
	if !tfPlanGroup.OnPremisesSyncEnabled.IsUnknown() {
		tfPlanOnPremisesSyncEnabled := tfPlanGroup.OnPremisesSyncEnabled.ValueBool()
		requestBodyGroup.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	} else {
		tfPlanGroup.OnPremisesSyncEnabled = types.BoolNull()
	}
	// END OnPremisesSyncEnabled | CreateBoolAttribute

	// START PreferredDataLocation | CreateStringAttribute
	if !tfPlanGroup.PreferredDataLocation.IsUnknown() {
		tfPlanPreferredDataLocation := tfPlanGroup.PreferredDataLocation.ValueString()
		requestBodyGroup.SetPreferredDataLocation(&tfPlanPreferredDataLocation)
	} else {
		tfPlanGroup.PreferredDataLocation = types.StringNull()
	}
	// END PreferredDataLocation | CreateStringAttribute

	// START PreferredLanguage | CreateStringAttribute
	if !tfPlanGroup.PreferredLanguage.IsUnknown() {
		tfPlanPreferredLanguage := tfPlanGroup.PreferredLanguage.ValueString()
		requestBodyGroup.SetPreferredLanguage(&tfPlanPreferredLanguage)
	} else {
		tfPlanGroup.PreferredLanguage = types.StringNull()
	}
	// END PreferredLanguage | CreateStringAttribute

	// START ProxyAddresses | CreateArrayStringAttribute
	if len(tfPlanGroup.ProxyAddresses.Elements()) > 0 {
		var stringArrayProxyAddresses []string
		for _, i := range tfPlanGroup.ProxyAddresses.Elements() {
			stringArrayProxyAddresses = append(stringArrayProxyAddresses, i.String())
		}
		requestBodyGroup.SetProxyAddresses(stringArrayProxyAddresses)
	} else {
		tfPlanGroup.ProxyAddresses = types.ListNull(types.StringType)
	}
	// END ProxyAddresses | CreateArrayStringAttribute

	// START RenewedDateTime | CreateStringTimeAttribute
	if !tfPlanGroup.RenewedDateTime.IsUnknown() {
		tfPlanRenewedDateTime := tfPlanGroup.RenewedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanRenewedDateTime)
		requestBodyGroup.SetRenewedDateTime(&t)
	} else {
		tfPlanGroup.RenewedDateTime = types.StringNull()
	}
	// END RenewedDateTime | CreateStringTimeAttribute

	// START SecurityEnabled | CreateBoolAttribute
	if !tfPlanGroup.SecurityEnabled.IsUnknown() {
		tfPlanSecurityEnabled := tfPlanGroup.SecurityEnabled.ValueBool()
		requestBodyGroup.SetSecurityEnabled(&tfPlanSecurityEnabled)
	} else {
		tfPlanGroup.SecurityEnabled = types.BoolNull()
	}
	// END SecurityEnabled | CreateBoolAttribute

	// START SecurityIdentifier | CreateStringAttribute
	if !tfPlanGroup.SecurityIdentifier.IsUnknown() {
		tfPlanSecurityIdentifier := tfPlanGroup.SecurityIdentifier.ValueString()
		requestBodyGroup.SetSecurityIdentifier(&tfPlanSecurityIdentifier)
	} else {
		tfPlanGroup.SecurityIdentifier = types.StringNull()
	}
	// END SecurityIdentifier | CreateStringAttribute

	// START ServiceProvisioningErrors | CreateArrayObjectAttribute
	if len(tfPlanGroup.ServiceProvisioningErrors.Elements()) > 0 {
		var requestBodyServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for _, i := range tfPlanGroup.ServiceProvisioningErrors.Elements() {
			requestBodyServiceProvisioningError := models.NewServiceProvisioningError()
			tfPlanServiceProvisioningErrors := groupServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanServiceProvisioningErrors)

			// START CreatedDateTime | CreateStringTimeAttribute
			if !tfPlanServiceProvisioningErrors.CreatedDateTime.IsUnknown() {
				tfPlanCreatedDateTime := tfPlanServiceProvisioningErrors.CreatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
				requestBodyServiceProvisioningErrors.SetCreatedDateTime(&t)
			} else {
				tfPlanServiceProvisioningErrors.CreatedDateTime = types.StringNull()
			}
			// END CreatedDateTime | CreateStringTimeAttribute

			// START IsResolved | CreateBoolAttribute
			if !tfPlanServiceProvisioningErrors.IsResolved.IsUnknown() {
				tfPlanIsResolved := tfPlanServiceProvisioningErrors.IsResolved.ValueBool()
				requestBodyServiceProvisioningErrors.SetIsResolved(&tfPlanIsResolved)
			} else {
				tfPlanServiceProvisioningErrors.IsResolved = types.BoolNull()
			}
			// END IsResolved | CreateBoolAttribute

			// START ServiceInstance | CreateStringAttribute
			if !tfPlanServiceProvisioningErrors.ServiceInstance.IsUnknown() {
				tfPlanServiceInstance := tfPlanServiceProvisioningErrors.ServiceInstance.ValueString()
				requestBodyServiceProvisioningErrors.SetServiceInstance(&tfPlanServiceInstance)
			} else {
				tfPlanServiceProvisioningErrors.ServiceInstance = types.StringNull()
			}
			// END ServiceInstance | CreateStringAttribute

		}
		requestBodyGroup.SetServiceProvisioningErrors(requestBodyServiceProvisioningErrors)
	} else {
		tfPlanGroup.ServiceProvisioningErrors = types.ListNull(tfPlanGroup.ServiceProvisioningErrors.ElementType(ctx))
	}
	// END ServiceProvisioningErrors | CreateArrayObjectAttribute

	// START Theme | CreateStringAttribute
	if !tfPlanGroup.Theme.IsUnknown() {
		tfPlanTheme := tfPlanGroup.Theme.ValueString()
		requestBodyGroup.SetTheme(&tfPlanTheme)
	} else {
		tfPlanGroup.Theme = types.StringNull()
	}
	// END Theme | CreateStringAttribute

	// START UniqueName | CreateStringAttribute
	if !tfPlanGroup.UniqueName.IsUnknown() {
		tfPlanUniqueName := tfPlanGroup.UniqueName.ValueString()
		requestBodyGroup.SetUniqueName(&tfPlanUniqueName)
	} else {
		tfPlanGroup.UniqueName = types.StringNull()
	}
	// END UniqueName | CreateStringAttribute

	// START Visibility | CreateStringAttribute
	if !tfPlanGroup.Visibility.IsUnknown() {
		tfPlanVisibility := tfPlanGroup.Visibility.ValueString()
		requestBodyGroup.SetVisibility(&tfPlanVisibility)
	} else {
		tfPlanGroup.Visibility = types.StringNull()
	}
	// END Visibility | CreateStringAttribute

	// Create new group
	result, err := r.client.Groups().Post(context.Background(), requestBodyGroup, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlanGroup.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlanGroup)
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
				"isManagementRestricted",
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
				"uniqueName",
				"visibility",
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
			assignedLabels := new(groupAssignedLabelModel)

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
			assignedLicenses := new(groupAssignedLicenseModel)

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
	if result.GetIsManagementRestricted() != nil {
		state.IsManagementRestricted = types.BoolValue(*result.GetIsManagementRestricted())
	} else {
		state.IsManagementRestricted = types.BoolNull()
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
			onPremisesProvisioningErrors := new(groupOnPremisesProvisioningErrorModel)

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
			serviceProvisioningErrors := new(groupServiceProvisioningErrorModel)

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
	if result.GetUniqueName() != nil {
		state.UniqueName = types.StringValue(*result.GetUniqueName())
	} else {
		state.UniqueName = types.StringNull()
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

	if !plan.Id.Equal(state.Id) {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	}

	if !plan.DeletedDateTime.Equal(state.DeletedDateTime) {
		planDeletedDateTime := plan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !plan.AssignedLabels.Equal(state.AssignedLabels) {
		var planAssignedLabels []models.AssignedLabelable
		for k, i := range plan.AssignedLabels.Elements() {
			assignedLabels := models.NewAssignedLabel()
			assignedLabelsModel := groupAssignedLabelModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedLabelsModel)
			assignedLabelsState := groupAssignedLabelModel{}
			types.ListValueFrom(ctx, state.AssignedLabels.Elements()[k].Type(ctx), &assignedLabelsModel)

			if !assignedLabelsModel.DisplayName.Equal(assignedLabelsState.DisplayName) {
				planDisplayName := assignedLabelsModel.DisplayName.ValueString()
				assignedLabels.SetDisplayName(&planDisplayName)
			}

			if !assignedLabelsModel.LabelId.Equal(assignedLabelsState.LabelId) {
				planLabelId := assignedLabelsModel.LabelId.ValueString()
				assignedLabels.SetLabelId(&planLabelId)
			}
		}
		requestBody.SetAssignedLabels(planAssignedLabels)
	}

	if !plan.AssignedLicenses.Equal(state.AssignedLicenses) {
		var planAssignedLicenses []models.AssignedLicenseable
		for k, i := range plan.AssignedLicenses.Elements() {
			assignedLicenses := models.NewAssignedLicense()
			assignedLicensesModel := groupAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &assignedLicensesModel)
			assignedLicensesState := groupAssignedLicenseModel{}
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

	if !plan.Classification.Equal(state.Classification) {
		planClassification := plan.Classification.ValueString()
		requestBody.SetClassification(&planClassification)
	}

	if !plan.CreatedDateTime.Equal(state.CreatedDateTime) {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planCreatedDateTime)
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
		t, _ := time.Parse(time.RFC3339, planExpirationDateTime)
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

	if !plan.IsManagementRestricted.Equal(state.IsManagementRestricted) {
		planIsManagementRestricted := plan.IsManagementRestricted.ValueBool()
		requestBody.SetIsManagementRestricted(&planIsManagementRestricted)
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
		t, _ := time.Parse(time.RFC3339, planOnPremisesLastSyncDateTime)
		requestBody.SetOnPremisesLastSyncDateTime(&t)
	}

	if !plan.OnPremisesNetBiosName.Equal(state.OnPremisesNetBiosName) {
		planOnPremisesNetBiosName := plan.OnPremisesNetBiosName.ValueString()
		requestBody.SetOnPremisesNetBiosName(&planOnPremisesNetBiosName)
	}

	if !plan.OnPremisesProvisioningErrors.Equal(state.OnPremisesProvisioningErrors) {
		var planOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for k, i := range plan.OnPremisesProvisioningErrors.Elements() {
			onPremisesProvisioningErrors := models.NewOnPremisesProvisioningError()
			onPremisesProvisioningErrorsModel := groupOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &onPremisesProvisioningErrorsModel)
			onPremisesProvisioningErrorsState := groupOnPremisesProvisioningErrorModel{}
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
		t, _ := time.Parse(time.RFC3339, planRenewedDateTime)
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
			serviceProvisioningErrors := models.NewServiceProvisioningError()
			serviceProvisioningErrorsModel := groupServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &serviceProvisioningErrorsModel)
			serviceProvisioningErrorsState := groupServiceProvisioningErrorModel{}
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

	if !plan.Theme.Equal(state.Theme) {
		planTheme := plan.Theme.ValueString()
		requestBody.SetTheme(&planTheme)
	}

	if !plan.UniqueName.Equal(state.UniqueName) {
		planUniqueName := plan.UniqueName.ValueString()
		requestBody.SetUniqueName(&planUniqueName)
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
