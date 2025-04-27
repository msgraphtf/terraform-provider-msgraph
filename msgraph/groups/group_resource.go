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
	if !tfPlanGroup.Id.IsUnknown() {
		tfPlanId := tfPlanGroup.Id.ValueString()
		requestBodyGroup.SetId(&tfPlanId)
	} else {
		tfPlanGroup.Id = types.StringNull()
	}

	if !tfPlanGroup.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlanGroup.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyGroup.SetDeletedDateTime(&t)
	} else {
		tfPlanGroup.DeletedDateTime = types.StringNull()
	}

	if len(tfPlanGroup.AssignedLabels.Elements()) > 0 {
		var requestBodyAssignedLabels []models.AssignedLabelable
		for _, i := range tfPlanGroup.AssignedLabels.Elements() {
			requestBodyAssignedLabel := models.NewAssignedLabel()
			tfPlanAssignedLabel := groupAssignedLabelModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLabel)

			if !tfPlanAssignedLabel.DisplayName.IsUnknown() {
				tfPlanDisplayName := tfPlanAssignedLabel.DisplayName.ValueString()
				requestBodyAssignedLabel.SetDisplayName(&tfPlanDisplayName)
			} else {
				tfPlanAssignedLabel.DisplayName = types.StringNull()
			}

			if !tfPlanAssignedLabel.LabelId.IsUnknown() {
				tfPlanLabelId := tfPlanAssignedLabel.LabelId.ValueString()
				requestBodyAssignedLabel.SetLabelId(&tfPlanLabelId)
			} else {
				tfPlanAssignedLabel.LabelId = types.StringNull()
			}

		}
		requestBodyGroup.SetAssignedLabels(requestBodyAssignedLabels)
	} else {
		tfPlanGroup.AssignedLabels = types.ListNull(tfPlanGroup.AssignedLabels.ElementType(ctx))
	}

	if len(tfPlanGroup.AssignedLicenses.Elements()) > 0 {
		var requestBodyAssignedLicenses []models.AssignedLicenseable
		for _, i := range tfPlanGroup.AssignedLicenses.Elements() {
			requestBodyAssignedLicense := models.NewAssignedLicense()
			tfPlanAssignedLicense := groupAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLicense)

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

			if !tfPlanAssignedLicense.SkuId.IsUnknown() {
				tfPlanSkuId := tfPlanAssignedLicense.SkuId.ValueString()
				u, _ := uuid.Parse(tfPlanSkuId)
				requestBodyAssignedLicense.SetSkuId(&u)
			} else {
				tfPlanAssignedLicense.SkuId = types.StringNull()
			}

		}
		requestBodyGroup.SetAssignedLicenses(requestBodyAssignedLicenses)
	} else {
		tfPlanGroup.AssignedLicenses = types.ListNull(tfPlanGroup.AssignedLicenses.ElementType(ctx))
	}

	if !tfPlanGroup.Classification.IsUnknown() {
		tfPlanClassification := tfPlanGroup.Classification.ValueString()
		requestBodyGroup.SetClassification(&tfPlanClassification)
	} else {
		tfPlanGroup.Classification = types.StringNull()
	}

	if !tfPlanGroup.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlanGroup.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyGroup.SetCreatedDateTime(&t)
	} else {
		tfPlanGroup.CreatedDateTime = types.StringNull()
	}

	if !tfPlanGroup.Description.IsUnknown() {
		tfPlanDescription := tfPlanGroup.Description.ValueString()
		requestBodyGroup.SetDescription(&tfPlanDescription)
	} else {
		tfPlanGroup.Description = types.StringNull()
	}

	if !tfPlanGroup.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanGroup.DisplayName.ValueString()
		requestBodyGroup.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanGroup.DisplayName = types.StringNull()
	}

	if !tfPlanGroup.ExpirationDateTime.IsUnknown() {
		tfPlanExpirationDateTime := tfPlanGroup.ExpirationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanExpirationDateTime)
		requestBodyGroup.SetExpirationDateTime(&t)
	} else {
		tfPlanGroup.ExpirationDateTime = types.StringNull()
	}

	if len(tfPlanGroup.GroupTypes.Elements()) > 0 {
		var stringArrayGroupTypes []string
		for _, i := range tfPlanGroup.GroupTypes.Elements() {
			stringArrayGroupTypes = append(stringArrayGroupTypes, i.String())
		}
		requestBodyGroup.SetGroupTypes(stringArrayGroupTypes)
	} else {
		tfPlanGroup.GroupTypes = types.ListNull(types.StringType)
	}

	if !tfPlanGroup.IsAssignableToRole.IsUnknown() {
		tfPlanIsAssignableToRole := tfPlanGroup.IsAssignableToRole.ValueBool()
		requestBodyGroup.SetIsAssignableToRole(&tfPlanIsAssignableToRole)
	} else {
		tfPlanGroup.IsAssignableToRole = types.BoolNull()
	}

	if !tfPlanGroup.IsManagementRestricted.IsUnknown() {
		tfPlanIsManagementRestricted := tfPlanGroup.IsManagementRestricted.ValueBool()
		requestBodyGroup.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	} else {
		tfPlanGroup.IsManagementRestricted = types.BoolNull()
	}

	if !tfPlanGroup.LicenseProcessingState.IsUnknown() {
		requestBodyLicenseProcessingState := models.NewLicenseProcessingState()
		tfPlanLicenseProcessingState := groupLicenseProcessingStateModel{}
		tfPlanGroup.LicenseProcessingState.As(ctx, &tfPlanLicenseProcessingState, basetypes.ObjectAsOptions{})

		if !tfPlanLicenseProcessingState.State.IsUnknown() {
			tfPlanState := tfPlanLicenseProcessingState.State.ValueString()
			requestBodyLicenseProcessingState.SetState(&tfPlanState)
		} else {
			tfPlanLicenseProcessingState.State = types.StringNull()
		}

		requestBodyGroup.SetLicenseProcessingState(requestBodyLicenseProcessingState)
		tfPlanGroup.LicenseProcessingState, _ = types.ObjectValueFrom(ctx, tfPlanLicenseProcessingState.AttributeTypes(), requestBodyLicenseProcessingState)
	} else {
		tfPlanGroup.LicenseProcessingState = types.ObjectNull(tfPlanGroup.LicenseProcessingState.AttributeTypes(ctx))
	}

	if !tfPlanGroup.Mail.IsUnknown() {
		tfPlanMail := tfPlanGroup.Mail.ValueString()
		requestBodyGroup.SetMail(&tfPlanMail)
	} else {
		tfPlanGroup.Mail = types.StringNull()
	}

	if !tfPlanGroup.MailEnabled.IsUnknown() {
		tfPlanMailEnabled := tfPlanGroup.MailEnabled.ValueBool()
		requestBodyGroup.SetMailEnabled(&tfPlanMailEnabled)
	} else {
		tfPlanGroup.MailEnabled = types.BoolNull()
	}

	if !tfPlanGroup.MailNickname.IsUnknown() {
		tfPlanMailNickname := tfPlanGroup.MailNickname.ValueString()
		requestBodyGroup.SetMailNickname(&tfPlanMailNickname)
	} else {
		tfPlanGroup.MailNickname = types.StringNull()
	}

	if !tfPlanGroup.MembershipRule.IsUnknown() {
		tfPlanMembershipRule := tfPlanGroup.MembershipRule.ValueString()
		requestBodyGroup.SetMembershipRule(&tfPlanMembershipRule)
	} else {
		tfPlanGroup.MembershipRule = types.StringNull()
	}

	if !tfPlanGroup.MembershipRuleProcessingState.IsUnknown() {
		tfPlanMembershipRuleProcessingState := tfPlanGroup.MembershipRuleProcessingState.ValueString()
		requestBodyGroup.SetMembershipRuleProcessingState(&tfPlanMembershipRuleProcessingState)
	} else {
		tfPlanGroup.MembershipRuleProcessingState = types.StringNull()
	}

	if !tfPlanGroup.OnPremisesDomainName.IsUnknown() {
		tfPlanOnPremisesDomainName := tfPlanGroup.OnPremisesDomainName.ValueString()
		requestBodyGroup.SetOnPremisesDomainName(&tfPlanOnPremisesDomainName)
	} else {
		tfPlanGroup.OnPremisesDomainName = types.StringNull()
	}

	if !tfPlanGroup.OnPremisesLastSyncDateTime.IsUnknown() {
		tfPlanOnPremisesLastSyncDateTime := tfPlanGroup.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		requestBodyGroup.SetOnPremisesLastSyncDateTime(&t)
	} else {
		tfPlanGroup.OnPremisesLastSyncDateTime = types.StringNull()
	}

	if !tfPlanGroup.OnPremisesNetBiosName.IsUnknown() {
		tfPlanOnPremisesNetBiosName := tfPlanGroup.OnPremisesNetBiosName.ValueString()
		requestBodyGroup.SetOnPremisesNetBiosName(&tfPlanOnPremisesNetBiosName)
	} else {
		tfPlanGroup.OnPremisesNetBiosName = types.StringNull()
	}

	if len(tfPlanGroup.OnPremisesProvisioningErrors.Elements()) > 0 {
		var requestBodyOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for _, i := range tfPlanGroup.OnPremisesProvisioningErrors.Elements() {
			requestBodyOnPremisesProvisioningError := models.NewOnPremisesProvisioningError()
			tfPlanOnPremisesProvisioningError := groupOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOnPremisesProvisioningError)

			if !tfPlanOnPremisesProvisioningError.Category.IsUnknown() {
				tfPlanCategory := tfPlanOnPremisesProvisioningError.Category.ValueString()
				requestBodyOnPremisesProvisioningError.SetCategory(&tfPlanCategory)
			} else {
				tfPlanOnPremisesProvisioningError.Category = types.StringNull()
			}

			if !tfPlanOnPremisesProvisioningError.OccurredDateTime.IsUnknown() {
				tfPlanOccurredDateTime := tfPlanOnPremisesProvisioningError.OccurredDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanOccurredDateTime)
				requestBodyOnPremisesProvisioningError.SetOccurredDateTime(&t)
			} else {
				tfPlanOnPremisesProvisioningError.OccurredDateTime = types.StringNull()
			}

			if !tfPlanOnPremisesProvisioningError.PropertyCausingError.IsUnknown() {
				tfPlanPropertyCausingError := tfPlanOnPremisesProvisioningError.PropertyCausingError.ValueString()
				requestBodyOnPremisesProvisioningError.SetPropertyCausingError(&tfPlanPropertyCausingError)
			} else {
				tfPlanOnPremisesProvisioningError.PropertyCausingError = types.StringNull()
			}

			if !tfPlanOnPremisesProvisioningError.Value.IsUnknown() {
				tfPlanValue := tfPlanOnPremisesProvisioningError.Value.ValueString()
				requestBodyOnPremisesProvisioningError.SetValue(&tfPlanValue)
			} else {
				tfPlanOnPremisesProvisioningError.Value = types.StringNull()
			}

		}
		requestBodyGroup.SetOnPremisesProvisioningErrors(requestBodyOnPremisesProvisioningErrors)
	} else {
		tfPlanGroup.OnPremisesProvisioningErrors = types.ListNull(tfPlanGroup.OnPremisesProvisioningErrors.ElementType(ctx))
	}

	if !tfPlanGroup.OnPremisesSamAccountName.IsUnknown() {
		tfPlanOnPremisesSamAccountName := tfPlanGroup.OnPremisesSamAccountName.ValueString()
		requestBodyGroup.SetOnPremisesSamAccountName(&tfPlanOnPremisesSamAccountName)
	} else {
		tfPlanGroup.OnPremisesSamAccountName = types.StringNull()
	}

	if !tfPlanGroup.OnPremisesSecurityIdentifier.IsUnknown() {
		tfPlanOnPremisesSecurityIdentifier := tfPlanGroup.OnPremisesSecurityIdentifier.ValueString()
		requestBodyGroup.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	} else {
		tfPlanGroup.OnPremisesSecurityIdentifier = types.StringNull()
	}

	if !tfPlanGroup.OnPremisesSyncEnabled.IsUnknown() {
		tfPlanOnPremisesSyncEnabled := tfPlanGroup.OnPremisesSyncEnabled.ValueBool()
		requestBodyGroup.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	} else {
		tfPlanGroup.OnPremisesSyncEnabled = types.BoolNull()
	}

	if !tfPlanGroup.PreferredDataLocation.IsUnknown() {
		tfPlanPreferredDataLocation := tfPlanGroup.PreferredDataLocation.ValueString()
		requestBodyGroup.SetPreferredDataLocation(&tfPlanPreferredDataLocation)
	} else {
		tfPlanGroup.PreferredDataLocation = types.StringNull()
	}

	if !tfPlanGroup.PreferredLanguage.IsUnknown() {
		tfPlanPreferredLanguage := tfPlanGroup.PreferredLanguage.ValueString()
		requestBodyGroup.SetPreferredLanguage(&tfPlanPreferredLanguage)
	} else {
		tfPlanGroup.PreferredLanguage = types.StringNull()
	}

	if len(tfPlanGroup.ProxyAddresses.Elements()) > 0 {
		var stringArrayProxyAddresses []string
		for _, i := range tfPlanGroup.ProxyAddresses.Elements() {
			stringArrayProxyAddresses = append(stringArrayProxyAddresses, i.String())
		}
		requestBodyGroup.SetProxyAddresses(stringArrayProxyAddresses)
	} else {
		tfPlanGroup.ProxyAddresses = types.ListNull(types.StringType)
	}

	if !tfPlanGroup.RenewedDateTime.IsUnknown() {
		tfPlanRenewedDateTime := tfPlanGroup.RenewedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanRenewedDateTime)
		requestBodyGroup.SetRenewedDateTime(&t)
	} else {
		tfPlanGroup.RenewedDateTime = types.StringNull()
	}

	if !tfPlanGroup.SecurityEnabled.IsUnknown() {
		tfPlanSecurityEnabled := tfPlanGroup.SecurityEnabled.ValueBool()
		requestBodyGroup.SetSecurityEnabled(&tfPlanSecurityEnabled)
	} else {
		tfPlanGroup.SecurityEnabled = types.BoolNull()
	}

	if !tfPlanGroup.SecurityIdentifier.IsUnknown() {
		tfPlanSecurityIdentifier := tfPlanGroup.SecurityIdentifier.ValueString()
		requestBodyGroup.SetSecurityIdentifier(&tfPlanSecurityIdentifier)
	} else {
		tfPlanGroup.SecurityIdentifier = types.StringNull()
	}

	if len(tfPlanGroup.ServiceProvisioningErrors.Elements()) > 0 {
		var requestBodyServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for _, i := range tfPlanGroup.ServiceProvisioningErrors.Elements() {
			requestBodyServiceProvisioningError := models.NewServiceProvisioningError()
			tfPlanServiceProvisioningError := groupServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanServiceProvisioningError)

			if !tfPlanServiceProvisioningError.CreatedDateTime.IsUnknown() {
				tfPlanCreatedDateTime := tfPlanServiceProvisioningError.CreatedDateTime.ValueString()
				t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
				requestBodyServiceProvisioningError.SetCreatedDateTime(&t)
			} else {
				tfPlanServiceProvisioningError.CreatedDateTime = types.StringNull()
			}

			if !tfPlanServiceProvisioningError.IsResolved.IsUnknown() {
				tfPlanIsResolved := tfPlanServiceProvisioningError.IsResolved.ValueBool()
				requestBodyServiceProvisioningError.SetIsResolved(&tfPlanIsResolved)
			} else {
				tfPlanServiceProvisioningError.IsResolved = types.BoolNull()
			}

			if !tfPlanServiceProvisioningError.ServiceInstance.IsUnknown() {
				tfPlanServiceInstance := tfPlanServiceProvisioningError.ServiceInstance.ValueString()
				requestBodyServiceProvisioningError.SetServiceInstance(&tfPlanServiceInstance)
			} else {
				tfPlanServiceProvisioningError.ServiceInstance = types.StringNull()
			}

		}
		requestBodyGroup.SetServiceProvisioningErrors(requestBodyServiceProvisioningErrors)
	} else {
		tfPlanGroup.ServiceProvisioningErrors = types.ListNull(tfPlanGroup.ServiceProvisioningErrors.ElementType(ctx))
	}

	if !tfPlanGroup.Theme.IsUnknown() {
		tfPlanTheme := tfPlanGroup.Theme.ValueString()
		requestBodyGroup.SetTheme(&tfPlanTheme)
	} else {
		tfPlanGroup.Theme = types.StringNull()
	}

	if !tfPlanGroup.UniqueName.IsUnknown() {
		tfPlanUniqueName := tfPlanGroup.UniqueName.ValueString()
		requestBodyGroup.SetUniqueName(&tfPlanUniqueName)
	} else {
		tfPlanGroup.UniqueName = types.StringNull()
	}

	if !tfPlanGroup.Visibility.IsUnknown() {
		tfPlanVisibility := tfPlanGroup.Visibility.ValueString()
		requestBodyGroup.SetVisibility(&tfPlanVisibility)
	} else {
		tfPlanGroup.Visibility = types.StringNull()
	}

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
	var tfStateGroup groupModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfStateGroup)...)
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

	var responseGroup models.Groupable
	var err error

	if !tfStateGroup.Id.IsNull() {
		responseGroup, err = d.client.Groups().ByGroupId(tfStateGroup.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Group",
			err.Error(),
		)
		return
	}

	if responseGroup.GetId() != nil {
		tfStateGroup.Id = types.StringValue(*responseGroup.GetId())
	} else {
		tfStateGroup.Id = types.StringNull()
	}
	if responseGroup.GetDeletedDateTime() != nil {
		tfStateGroup.DeletedDateTime = types.StringValue(responseGroup.GetDeletedDateTime().String())
	} else {
		tfStateGroup.DeletedDateTime = types.StringNull()
	}
	if len(responseGroup.GetAssignedLabels()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseAssignedLabel := range responseGroup.GetAssignedLabels() {
			tfStateAssignedLabel := groupAssignedLabelModel{}

			if responseAssignedLabel.GetDisplayName() != nil {
				tfStateAssignedLabel.DisplayName = types.StringValue(*responseAssignedLabel.GetDisplayName())
			} else {
				tfStateAssignedLabel.DisplayName = types.StringNull()
			}
			if responseAssignedLabel.GetLabelId() != nil {
				tfStateAssignedLabel.LabelId = types.StringValue(*responseAssignedLabel.GetLabelId())
			} else {
				tfStateAssignedLabel.LabelId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedLabel.AttributeTypes(), tfStateAssignedLabel)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.AssignedLabels, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(responseGroup.GetAssignedLicenses()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseAssignedLicense := range responseGroup.GetAssignedLicenses() {
			tfStateAssignedLicense := groupAssignedLicenseModel{}

			if len(responseAssignedLicense.GetDisabledPlans()) > 0 {
				var valueArrayDisabledPlans []attr.Value
				for _, responseDisabledPlans := range responseAssignedLicense.GetDisabledPlans() {
					valueArrayDisabledPlans = append(valueArrayDisabledPlans, types.StringValue(responseDisabledPlans.String()))
				}
				tfStateAssignedLicense.DisabledPlans, _ = types.ListValue(types.StringType, valueArrayDisabledPlans)
			} else {
				tfStateAssignedLicense.DisabledPlans = types.ListNull(types.StringType)
			}
			if responseAssignedLicense.GetSkuId() != nil {
				tfStateAssignedLicense.SkuId = types.StringValue(responseAssignedLicense.GetSkuId().String())
			} else {
				tfStateAssignedLicense.SkuId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedLicense.AttributeTypes(), tfStateAssignedLicense)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.AssignedLicenses, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if responseGroup.GetClassification() != nil {
		tfStateGroup.Classification = types.StringValue(*responseGroup.GetClassification())
	} else {
		tfStateGroup.Classification = types.StringNull()
	}
	if responseGroup.GetCreatedDateTime() != nil {
		tfStateGroup.CreatedDateTime = types.StringValue(responseGroup.GetCreatedDateTime().String())
	} else {
		tfStateGroup.CreatedDateTime = types.StringNull()
	}
	if responseGroup.GetDescription() != nil {
		tfStateGroup.Description = types.StringValue(*responseGroup.GetDescription())
	} else {
		tfStateGroup.Description = types.StringNull()
	}
	if responseGroup.GetDisplayName() != nil {
		tfStateGroup.DisplayName = types.StringValue(*responseGroup.GetDisplayName())
	} else {
		tfStateGroup.DisplayName = types.StringNull()
	}
	if responseGroup.GetExpirationDateTime() != nil {
		tfStateGroup.ExpirationDateTime = types.StringValue(responseGroup.GetExpirationDateTime().String())
	} else {
		tfStateGroup.ExpirationDateTime = types.StringNull()
	}
	if len(responseGroup.GetGroupTypes()) > 0 {
		var valueArrayGroupTypes []attr.Value
		for _, responseGroupTypes := range responseGroup.GetGroupTypes() {
			valueArrayGroupTypes = append(valueArrayGroupTypes, types.StringValue(responseGroupTypes))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayGroupTypes)
		tfStateGroup.GroupTypes = listValue
	} else {
		tfStateGroup.GroupTypes = types.ListNull(types.StringType)
	}
	if responseGroup.GetIsAssignableToRole() != nil {
		tfStateGroup.IsAssignableToRole = types.BoolValue(*responseGroup.GetIsAssignableToRole())
	} else {
		tfStateGroup.IsAssignableToRole = types.BoolNull()
	}
	if responseGroup.GetIsManagementRestricted() != nil {
		tfStateGroup.IsManagementRestricted = types.BoolValue(*responseGroup.GetIsManagementRestricted())
	} else {
		tfStateGroup.IsManagementRestricted = types.BoolNull()
	}
	if responseGroup.GetLicenseProcessingState() != nil {
		tfStateLicenseProcessingState := groupLicenseProcessingStateModel{}
		responseLicenseProcessingState := responseGroup.GetLicenseProcessingState()

		if responseLicenseProcessingState.GetState() != nil {
			tfStateLicenseProcessingState.State = types.StringValue(*responseLicenseProcessingState.GetState())
		} else {
			tfStateLicenseProcessingState.State = types.StringNull()
		}

		tfStateGroup.LicenseProcessingState, _ = types.ObjectValueFrom(ctx, tfStateLicenseProcessingState.AttributeTypes(), tfStateLicenseProcessingState)
	}
	if responseGroup.GetMail() != nil {
		tfStateGroup.Mail = types.StringValue(*responseGroup.GetMail())
	} else {
		tfStateGroup.Mail = types.StringNull()
	}
	if responseGroup.GetMailEnabled() != nil {
		tfStateGroup.MailEnabled = types.BoolValue(*responseGroup.GetMailEnabled())
	} else {
		tfStateGroup.MailEnabled = types.BoolNull()
	}
	if responseGroup.GetMailNickname() != nil {
		tfStateGroup.MailNickname = types.StringValue(*responseGroup.GetMailNickname())
	} else {
		tfStateGroup.MailNickname = types.StringNull()
	}
	if responseGroup.GetMembershipRule() != nil {
		tfStateGroup.MembershipRule = types.StringValue(*responseGroup.GetMembershipRule())
	} else {
		tfStateGroup.MembershipRule = types.StringNull()
	}
	if responseGroup.GetMembershipRuleProcessingState() != nil {
		tfStateGroup.MembershipRuleProcessingState = types.StringValue(*responseGroup.GetMembershipRuleProcessingState())
	} else {
		tfStateGroup.MembershipRuleProcessingState = types.StringNull()
	}
	if responseGroup.GetOnPremisesDomainName() != nil {
		tfStateGroup.OnPremisesDomainName = types.StringValue(*responseGroup.GetOnPremisesDomainName())
	} else {
		tfStateGroup.OnPremisesDomainName = types.StringNull()
	}
	if responseGroup.GetOnPremisesLastSyncDateTime() != nil {
		tfStateGroup.OnPremisesLastSyncDateTime = types.StringValue(responseGroup.GetOnPremisesLastSyncDateTime().String())
	} else {
		tfStateGroup.OnPremisesLastSyncDateTime = types.StringNull()
	}
	if responseGroup.GetOnPremisesNetBiosName() != nil {
		tfStateGroup.OnPremisesNetBiosName = types.StringValue(*responseGroup.GetOnPremisesNetBiosName())
	} else {
		tfStateGroup.OnPremisesNetBiosName = types.StringNull()
	}
	if len(responseGroup.GetOnPremisesProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseOnPremisesProvisioningError := range responseGroup.GetOnPremisesProvisioningErrors() {
			tfStateOnPremisesProvisioningError := groupOnPremisesProvisioningErrorModel{}

			if responseOnPremisesProvisioningError.GetCategory() != nil {
				tfStateOnPremisesProvisioningError.Category = types.StringValue(*responseOnPremisesProvisioningError.GetCategory())
			} else {
				tfStateOnPremisesProvisioningError.Category = types.StringNull()
			}
			if responseOnPremisesProvisioningError.GetOccurredDateTime() != nil {
				tfStateOnPremisesProvisioningError.OccurredDateTime = types.StringValue(responseOnPremisesProvisioningError.GetOccurredDateTime().String())
			} else {
				tfStateOnPremisesProvisioningError.OccurredDateTime = types.StringNull()
			}
			if responseOnPremisesProvisioningError.GetPropertyCausingError() != nil {
				tfStateOnPremisesProvisioningError.PropertyCausingError = types.StringValue(*responseOnPremisesProvisioningError.GetPropertyCausingError())
			} else {
				tfStateOnPremisesProvisioningError.PropertyCausingError = types.StringNull()
			}
			if responseOnPremisesProvisioningError.GetValue() != nil {
				tfStateOnPremisesProvisioningError.Value = types.StringValue(*responseOnPremisesProvisioningError.GetValue())
			} else {
				tfStateOnPremisesProvisioningError.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateOnPremisesProvisioningError.AttributeTypes(), tfStateOnPremisesProvisioningError)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.OnPremisesProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if responseGroup.GetOnPremisesSamAccountName() != nil {
		tfStateGroup.OnPremisesSamAccountName = types.StringValue(*responseGroup.GetOnPremisesSamAccountName())
	} else {
		tfStateGroup.OnPremisesSamAccountName = types.StringNull()
	}
	if responseGroup.GetOnPremisesSecurityIdentifier() != nil {
		tfStateGroup.OnPremisesSecurityIdentifier = types.StringValue(*responseGroup.GetOnPremisesSecurityIdentifier())
	} else {
		tfStateGroup.OnPremisesSecurityIdentifier = types.StringNull()
	}
	if responseGroup.GetOnPremisesSyncEnabled() != nil {
		tfStateGroup.OnPremisesSyncEnabled = types.BoolValue(*responseGroup.GetOnPremisesSyncEnabled())
	} else {
		tfStateGroup.OnPremisesSyncEnabled = types.BoolNull()
	}
	if responseGroup.GetPreferredDataLocation() != nil {
		tfStateGroup.PreferredDataLocation = types.StringValue(*responseGroup.GetPreferredDataLocation())
	} else {
		tfStateGroup.PreferredDataLocation = types.StringNull()
	}
	if responseGroup.GetPreferredLanguage() != nil {
		tfStateGroup.PreferredLanguage = types.StringValue(*responseGroup.GetPreferredLanguage())
	} else {
		tfStateGroup.PreferredLanguage = types.StringNull()
	}
	if len(responseGroup.GetProxyAddresses()) > 0 {
		var valueArrayProxyAddresses []attr.Value
		for _, responseProxyAddresses := range responseGroup.GetProxyAddresses() {
			valueArrayProxyAddresses = append(valueArrayProxyAddresses, types.StringValue(responseProxyAddresses))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayProxyAddresses)
		tfStateGroup.ProxyAddresses = listValue
	} else {
		tfStateGroup.ProxyAddresses = types.ListNull(types.StringType)
	}
	if responseGroup.GetRenewedDateTime() != nil {
		tfStateGroup.RenewedDateTime = types.StringValue(responseGroup.GetRenewedDateTime().String())
	} else {
		tfStateGroup.RenewedDateTime = types.StringNull()
	}
	if responseGroup.GetSecurityEnabled() != nil {
		tfStateGroup.SecurityEnabled = types.BoolValue(*responseGroup.GetSecurityEnabled())
	} else {
		tfStateGroup.SecurityEnabled = types.BoolNull()
	}
	if responseGroup.GetSecurityIdentifier() != nil {
		tfStateGroup.SecurityIdentifier = types.StringValue(*responseGroup.GetSecurityIdentifier())
	} else {
		tfStateGroup.SecurityIdentifier = types.StringNull()
	}
	if len(responseGroup.GetServiceProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseServiceProvisioningError := range responseGroup.GetServiceProvisioningErrors() {
			tfStateServiceProvisioningError := groupServiceProvisioningErrorModel{}

			if responseServiceProvisioningError.GetCreatedDateTime() != nil {
				tfStateServiceProvisioningError.CreatedDateTime = types.StringValue(responseServiceProvisioningError.GetCreatedDateTime().String())
			} else {
				tfStateServiceProvisioningError.CreatedDateTime = types.StringNull()
			}
			if responseServiceProvisioningError.GetIsResolved() != nil {
				tfStateServiceProvisioningError.IsResolved = types.BoolValue(*responseServiceProvisioningError.GetIsResolved())
			} else {
				tfStateServiceProvisioningError.IsResolved = types.BoolNull()
			}
			if responseServiceProvisioningError.GetServiceInstance() != nil {
				tfStateServiceProvisioningError.ServiceInstance = types.StringValue(*responseServiceProvisioningError.GetServiceInstance())
			} else {
				tfStateServiceProvisioningError.ServiceInstance = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateServiceProvisioningError.AttributeTypes(), tfStateServiceProvisioningError)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.ServiceProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if responseGroup.GetTheme() != nil {
		tfStateGroup.Theme = types.StringValue(*responseGroup.GetTheme())
	} else {
		tfStateGroup.Theme = types.StringNull()
	}
	if responseGroup.GetUniqueName() != nil {
		tfStateGroup.UniqueName = types.StringValue(*responseGroup.GetUniqueName())
	} else {
		tfStateGroup.UniqueName = types.StringNull()
	}
	if responseGroup.GetVisibility() != nil {
		tfStateGroup.Visibility = types.StringValue(*responseGroup.GetVisibility())
	} else {
		tfStateGroup.Visibility = types.StringNull()
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanGroup groupModel
	diags := req.Plan.Get(ctx, &tfPlanGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfStateGroup groupModel
	diags = req.State.Get(ctx, &tfStateGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBodyGroup := models.NewGroup()

	if !tfPlanGroup.Id.Equal(tfStateGroup.Id) {
		tfPlanId := tfPlanGroup.Id.ValueString()
		requestBodyGroup.SetId(&tfPlanId)
	}

	if !tfPlanGroup.DeletedDateTime.Equal(tfStateGroup.DeletedDateTime) {
		tfPlanDeletedDateTime := tfPlanGroup.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyGroup.SetDeletedDateTime(&t)
	}

	if !tfPlanGroup.AssignedLabels.Equal(tfStateGroup.AssignedLabels) {
		var tfPlanAssignedLabels []models.AssignedLabelable
		for k, i := range tfPlanGroup.AssignedLabels.Elements() {
			requestBodyAssignedLabel := models.NewAssignedLabel()
			tfPlanAssignedLabel := groupAssignedLabelModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLabel)
			tfStateAssignedLabel := groupAssignedLabelModel{}
			types.ListValueFrom(ctx, tfStateGroup.AssignedLabels.Elements()[k].Type(ctx), &tfPlanAssignedLabel)

			if !tfPlanAssignedLabel.DisplayName.Equal(tfStateAssignedLabel.DisplayName) {
				tfPlanDisplayName := tfPlanAssignedLabel.DisplayName.ValueString()
				requestBodyAssignedLabel.SetDisplayName(&tfPlanDisplayName)
			}

			if !tfPlanAssignedLabel.LabelId.Equal(tfStateAssignedLabel.LabelId) {
				tfPlanLabelId := tfPlanAssignedLabel.LabelId.ValueString()
				requestBodyAssignedLabel.SetLabelId(&tfPlanLabelId)
			}
		}
		requestBodyGroup.SetAssignedLabels(tfPlanAssignedLabels)
	}

	if !tfPlanGroup.AssignedLicenses.Equal(tfStateGroup.AssignedLicenses) {
		var tfPlanAssignedLicenses []models.AssignedLicenseable
		for k, i := range tfPlanGroup.AssignedLicenses.Elements() {
			requestBodyAssignedLicense := models.NewAssignedLicense()
			tfPlanAssignedLicense := groupAssignedLicenseModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAssignedLicense)
			tfStateAssignedLicense := groupAssignedLicenseModel{}
			types.ListValueFrom(ctx, tfStateGroup.AssignedLicenses.Elements()[k].Type(ctx), &tfPlanAssignedLicense)

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
		requestBodyGroup.SetAssignedLicenses(tfPlanAssignedLicenses)
	}

	if !tfPlanGroup.Classification.Equal(tfStateGroup.Classification) {
		tfPlanClassification := tfPlanGroup.Classification.ValueString()
		requestBodyGroup.SetClassification(&tfPlanClassification)
	}

	if !tfPlanGroup.CreatedDateTime.Equal(tfStateGroup.CreatedDateTime) {
		tfPlanCreatedDateTime := tfPlanGroup.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyGroup.SetCreatedDateTime(&t)
	}

	if !tfPlanGroup.Description.Equal(tfStateGroup.Description) {
		tfPlanDescription := tfPlanGroup.Description.ValueString()
		requestBodyGroup.SetDescription(&tfPlanDescription)
	}

	if !tfPlanGroup.DisplayName.Equal(tfStateGroup.DisplayName) {
		tfPlanDisplayName := tfPlanGroup.DisplayName.ValueString()
		requestBodyGroup.SetDisplayName(&tfPlanDisplayName)
	}

	if !tfPlanGroup.ExpirationDateTime.Equal(tfStateGroup.ExpirationDateTime) {
		tfPlanExpirationDateTime := tfPlanGroup.ExpirationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanExpirationDateTime)
		requestBodyGroup.SetExpirationDateTime(&t)
	}

	if !tfPlanGroup.GroupTypes.Equal(tfStateGroup.GroupTypes) {
		var stringArrayGroupTypes []string
		for _, i := range tfPlanGroup.GroupTypes.Elements() {
			stringArrayGroupTypes = append(stringArrayGroupTypes, i.String())
		}
		requestBodyGroup.SetGroupTypes(stringArrayGroupTypes)
	}

	if !tfPlanGroup.IsAssignableToRole.Equal(tfStateGroup.IsAssignableToRole) {
		tfPlanIsAssignableToRole := tfPlanGroup.IsAssignableToRole.ValueBool()
		requestBodyGroup.SetIsAssignableToRole(&tfPlanIsAssignableToRole)
	}

	if !tfPlanGroup.IsManagementRestricted.Equal(tfStateGroup.IsManagementRestricted) {
		tfPlanIsManagementRestricted := tfPlanGroup.IsManagementRestricted.ValueBool()
		requestBodyGroup.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	}

	if !tfPlanGroup.LicenseProcessingState.Equal(tfStateGroup.LicenseProcessingState) {
		requestBodyLicenseProcessingState := models.NewLicenseProcessingState()
		tfPlanLicenseProcessingState := groupLicenseProcessingStateModel{}
		tfPlanGroup.LicenseProcessingState.As(ctx, &tfPlanLicenseProcessingState, basetypes.ObjectAsOptions{})
		tfStateLicenseProcessingState := groupLicenseProcessingStateModel{}
		tfStateGroup.LicenseProcessingState.As(ctx, &tfStateLicenseProcessingState, basetypes.ObjectAsOptions{})

		if !tfPlanLicenseProcessingState.State.Equal(tfStateLicenseProcessingState.State) {
			tfPlanState := tfPlanLicenseProcessingState.State.ValueString()
			requestBodyLicenseProcessingState.SetState(&tfPlanState)
		}
		requestBodyGroup.SetLicenseProcessingState(requestBodyLicenseProcessingState)
		tfPlanGroup.LicenseProcessingState, _ = types.ObjectValueFrom(ctx, tfPlanLicenseProcessingState.AttributeTypes(), tfPlanLicenseProcessingState)
	}

	if !tfPlanGroup.Mail.Equal(tfStateGroup.Mail) {
		tfPlanMail := tfPlanGroup.Mail.ValueString()
		requestBodyGroup.SetMail(&tfPlanMail)
	}

	if !tfPlanGroup.MailEnabled.Equal(tfStateGroup.MailEnabled) {
		tfPlanMailEnabled := tfPlanGroup.MailEnabled.ValueBool()
		requestBodyGroup.SetMailEnabled(&tfPlanMailEnabled)
	}

	if !tfPlanGroup.MailNickname.Equal(tfStateGroup.MailNickname) {
		tfPlanMailNickname := tfPlanGroup.MailNickname.ValueString()
		requestBodyGroup.SetMailNickname(&tfPlanMailNickname)
	}

	if !tfPlanGroup.MembershipRule.Equal(tfStateGroup.MembershipRule) {
		tfPlanMembershipRule := tfPlanGroup.MembershipRule.ValueString()
		requestBodyGroup.SetMembershipRule(&tfPlanMembershipRule)
	}

	if !tfPlanGroup.MembershipRuleProcessingState.Equal(tfStateGroup.MembershipRuleProcessingState) {
		tfPlanMembershipRuleProcessingState := tfPlanGroup.MembershipRuleProcessingState.ValueString()
		requestBodyGroup.SetMembershipRuleProcessingState(&tfPlanMembershipRuleProcessingState)
	}

	if !tfPlanGroup.OnPremisesDomainName.Equal(tfStateGroup.OnPremisesDomainName) {
		tfPlanOnPremisesDomainName := tfPlanGroup.OnPremisesDomainName.ValueString()
		requestBodyGroup.SetOnPremisesDomainName(&tfPlanOnPremisesDomainName)
	}

	if !tfPlanGroup.OnPremisesLastSyncDateTime.Equal(tfStateGroup.OnPremisesLastSyncDateTime) {
		tfPlanOnPremisesLastSyncDateTime := tfPlanGroup.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		requestBodyGroup.SetOnPremisesLastSyncDateTime(&t)
	}

	if !tfPlanGroup.OnPremisesNetBiosName.Equal(tfStateGroup.OnPremisesNetBiosName) {
		tfPlanOnPremisesNetBiosName := tfPlanGroup.OnPremisesNetBiosName.ValueString()
		requestBodyGroup.SetOnPremisesNetBiosName(&tfPlanOnPremisesNetBiosName)
	}

	if !tfPlanGroup.OnPremisesProvisioningErrors.Equal(tfStateGroup.OnPremisesProvisioningErrors) {
		var tfPlanOnPremisesProvisioningErrors []models.OnPremisesProvisioningErrorable
		for k, i := range tfPlanGroup.OnPremisesProvisioningErrors.Elements() {
			requestBodyOnPremisesProvisioningError := models.NewOnPremisesProvisioningError()
			tfPlanOnPremisesProvisioningError := groupOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanOnPremisesProvisioningError)
			tfStateOnPremisesProvisioningError := groupOnPremisesProvisioningErrorModel{}
			types.ListValueFrom(ctx, tfStateGroup.OnPremisesProvisioningErrors.Elements()[k].Type(ctx), &tfPlanOnPremisesProvisioningError)

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
		requestBodyGroup.SetOnPremisesProvisioningErrors(tfPlanOnPremisesProvisioningErrors)
	}

	if !tfPlanGroup.OnPremisesSamAccountName.Equal(tfStateGroup.OnPremisesSamAccountName) {
		tfPlanOnPremisesSamAccountName := tfPlanGroup.OnPremisesSamAccountName.ValueString()
		requestBodyGroup.SetOnPremisesSamAccountName(&tfPlanOnPremisesSamAccountName)
	}

	if !tfPlanGroup.OnPremisesSecurityIdentifier.Equal(tfStateGroup.OnPremisesSecurityIdentifier) {
		tfPlanOnPremisesSecurityIdentifier := tfPlanGroup.OnPremisesSecurityIdentifier.ValueString()
		requestBodyGroup.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	}

	if !tfPlanGroup.OnPremisesSyncEnabled.Equal(tfStateGroup.OnPremisesSyncEnabled) {
		tfPlanOnPremisesSyncEnabled := tfPlanGroup.OnPremisesSyncEnabled.ValueBool()
		requestBodyGroup.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	}

	if !tfPlanGroup.PreferredDataLocation.Equal(tfStateGroup.PreferredDataLocation) {
		tfPlanPreferredDataLocation := tfPlanGroup.PreferredDataLocation.ValueString()
		requestBodyGroup.SetPreferredDataLocation(&tfPlanPreferredDataLocation)
	}

	if !tfPlanGroup.PreferredLanguage.Equal(tfStateGroup.PreferredLanguage) {
		tfPlanPreferredLanguage := tfPlanGroup.PreferredLanguage.ValueString()
		requestBodyGroup.SetPreferredLanguage(&tfPlanPreferredLanguage)
	}

	if !tfPlanGroup.ProxyAddresses.Equal(tfStateGroup.ProxyAddresses) {
		var stringArrayProxyAddresses []string
		for _, i := range tfPlanGroup.ProxyAddresses.Elements() {
			stringArrayProxyAddresses = append(stringArrayProxyAddresses, i.String())
		}
		requestBodyGroup.SetProxyAddresses(stringArrayProxyAddresses)
	}

	if !tfPlanGroup.RenewedDateTime.Equal(tfStateGroup.RenewedDateTime) {
		tfPlanRenewedDateTime := tfPlanGroup.RenewedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanRenewedDateTime)
		requestBodyGroup.SetRenewedDateTime(&t)
	}

	if !tfPlanGroup.SecurityEnabled.Equal(tfStateGroup.SecurityEnabled) {
		tfPlanSecurityEnabled := tfPlanGroup.SecurityEnabled.ValueBool()
		requestBodyGroup.SetSecurityEnabled(&tfPlanSecurityEnabled)
	}

	if !tfPlanGroup.SecurityIdentifier.Equal(tfStateGroup.SecurityIdentifier) {
		tfPlanSecurityIdentifier := tfPlanGroup.SecurityIdentifier.ValueString()
		requestBodyGroup.SetSecurityIdentifier(&tfPlanSecurityIdentifier)
	}

	if !tfPlanGroup.ServiceProvisioningErrors.Equal(tfStateGroup.ServiceProvisioningErrors) {
		var tfPlanServiceProvisioningErrors []models.ServiceProvisioningErrorable
		for k, i := range tfPlanGroup.ServiceProvisioningErrors.Elements() {
			requestBodyServiceProvisioningError := models.NewServiceProvisioningError()
			tfPlanServiceProvisioningError := groupServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanServiceProvisioningError)
			tfStateServiceProvisioningError := groupServiceProvisioningErrorModel{}
			types.ListValueFrom(ctx, tfStateGroup.ServiceProvisioningErrors.Elements()[k].Type(ctx), &tfPlanServiceProvisioningError)

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
		requestBodyGroup.SetServiceProvisioningErrors(tfPlanServiceProvisioningErrors)
	}

	if !tfPlanGroup.Theme.Equal(tfStateGroup.Theme) {
		tfPlanTheme := tfPlanGroup.Theme.ValueString()
		requestBodyGroup.SetTheme(&tfPlanTheme)
	}

	if !tfPlanGroup.UniqueName.Equal(tfStateGroup.UniqueName) {
		tfPlanUniqueName := tfPlanGroup.UniqueName.ValueString()
		requestBodyGroup.SetUniqueName(&tfPlanUniqueName)
	}

	if !tfPlanGroup.Visibility.Equal(tfStateGroup.Visibility) {
		tfPlanVisibility := tfPlanGroup.Visibility.ValueString()
		requestBodyGroup.SetVisibility(&tfPlanVisibility)
	}

	// Update group
	_, err := r.client.Groups().ByGroupId(tfStateGroup.Id.ValueString()).Patch(context.Background(), requestBodyGroup, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating group",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlanGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfStateGroup groupModel
	diags := req.State.Get(ctx, &tfStateGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete group
	err := r.client.Groups().ByGroupId(tfStateGroup.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting group",
			err.Error(),
		)
		return
	}

}
