package groups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

// NewGroupDataSource is a helper function to simplify the provider implementation.
func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

// groupDataSource is the data source implementation.
type groupDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Configure adds the provider configured client to the data source.
func (d *groupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"assigned_labels": schema.ListNestedAttribute{
				Description: "The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group. Returned only on $select.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"display_name": schema.StringAttribute{
							Description: "The display name of the label. Read-only.",
							Computed:    true,
						},
						"label_id": schema.StringAttribute{
							Description: "The unique identifier of the label.",
							Computed:    true,
						},
					},
				},
			},
			"assigned_licenses": schema.ListNestedAttribute{
				Description: "The licenses that are assigned to the group. Returned only on $select. Supports $filter (eq).Read-only.",
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
			"classification": schema.StringAttribute{
				Description: "Describes a classification for the group (such as low, medium or high business impact). Valid values for this property are defined by creating a ClassificationList setting value, based on the template definition.Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith).",
				Computed:    true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group was created. The value cannot be modified and is automatically populated when the group is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "An optional description for the group. Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name for the group. This property is required when a group is created and cannot be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Computed:    true,
			},
			"expiration_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group is set to expire. It is null for security groups, but for Microsoft 365 groups, it represents when the group is set to expire as defined in the groupLifecyclePolicy. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
				Computed:    true,
			},
			"group_types": schema.ListAttribute{
				Description: "Specifies the group type and its membership. If the collection contains Unified, the group is a Microsoft 365 group; otherwise, it's either a security group or a distribution group. For details, see groups overview.If the collection includes DynamicMembership, the group has dynamic membership; otherwise, membership is static. Returned by default. Supports $filter (eq, not).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_assignable_to_role": schema.BoolAttribute{
				Description: "Indicates whether this group can be assigned to a Microsoft Entra role. Optional. This property can only be set while creating the group and is immutable. If set to true, the securityEnabled property must also be set to true, visibility must be Hidden, and the group cannot be a dynamic group (that is, groupTypes cannot contain DynamicMembership). Only callers in Global Administrator and Privileged Role Administrator roles can set this property. The caller must also be assigned the RoleManagement.ReadWrite.Directory permission to set this property or update the membership of such groups. For more, see Using a group to manage Microsoft Entra role assignmentsUsing this feature requires a Microsoft Entra ID P1 license. Returned by default. Supports $filter (eq, ne, not).",
				Computed:    true,
			},
			"license_processing_state": schema.SingleNestedAttribute{
				Description: "Indicates the status of the group license assignment to all group members. The default value is false. Read-only. Possible values: QueuedForProcessing, ProcessingInProgress, and ProcessingComplete.Returned only on $select. Read-only.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"state": schema.StringAttribute{
						Description: "",
						Computed:    true,
					},
				},
			},
			"mail": schema.StringAttribute{
				Description: "The SMTP address for the group, for example, 'serviceadmins@contoso.onmicrosoft.com'. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Computed:    true,
			},
			"mail_enabled": schema.BoolAttribute{
				Description: "Specifies whether the group is mail-enabled. Required. Returned by default. Supports $filter (eq, ne, not).",
				Computed:    true,
			},
			"mail_nickname": schema.StringAttribute{
				Description: "The mail alias for the group, unique for Microsoft 365 groups in the organization. Maximum length is 64 characters. This property can contain only characters in the ASCII character set 0 - 127 except the following: @ () / [] ' ; : <> , SPACE. Required. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Computed:    true,
			},
			"membership_rule": schema.StringAttribute{
				Description: "The rule that determines members for this group if the group is a dynamic group (groupTypes contains DynamicMembership). For more information about the syntax of the membership rule, see Membership Rules syntax. Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith).",
				Computed:    true,
			},
			"membership_rule_processing_state": schema.StringAttribute{
				Description: "Indicates whether the dynamic membership processing is on or paused. Possible values are On or Paused. Returned by default. Supports $filter (eq, ne, not, in).",
				Computed:    true,
			},
			"on_premises_domain_name": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "Indicates the last time at which the group was synced with the on-premises directory.The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in).",
				Computed:    true,
			},
			"on_premises_net_bios_name": schema.StringAttribute{
				Description: "",
				Computed:    true,
			},
			"on_premises_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors when using Microsoft synchronization product during provisioning. Returned by default. Supports $filter (eq, not).",
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
				Description: "Contains the on-premises SAM account name synchronized from the on-premises directory. The property is only populated for customers synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith). Read-only.",
				Computed:    true,
			},
			"on_premises_security_identifier": schema.StringAttribute{
				Description: "Contains the on-premises security identifier (SID) for the group synchronized from on-premises to the cloud. Returned by default. Supports $filter (eq including on null values). Read-only.",
				Computed:    true,
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Description: "true if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default). Returned by default. Read-only. Supports $filter (eq, ne, not, in, and eq on null values).",
				Computed:    true,
			},
			"preferred_data_location": schema.StringAttribute{
				Description: "The preferred data location for the Microsoft 365 group. By default, the group inherits the group creator's preferred data location. To set this property, the calling app must be granted the Directory.ReadWrite.All permission and the user be assigned one of the following Microsoft Entra roles:  Global Administrator  User Account Administrator Directory Writer  Exchange Administrator  SharePoint Administrator  For more information about this property, see OneDrive Online Multi-Geo. Nullable. Returned by default.",
				Computed:    true,
			},
			"preferred_language": schema.StringAttribute{
				Description: "The preferred language for a Microsoft 365 group. Should follow ISO 639-1 Code; for example, en-US. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Computed:    true,
			},
			"proxy_addresses": schema.ListAttribute{
				Description: "Email addresses for the group that direct to the same group mailbox. For example: ['SMTP: bob@contoso.com', 'smtp: bob@sales.contoso.com']. The any operator is required to filter expressions on multi-valued properties. Returned by default. Read-only. Not nullable. Supports $filter (eq, not, ge, le, startsWith, endsWith, /$count eq 0, /$count ne 0).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"renewed_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group was last renewed. This cannot be modified directly and is only updated via the renew service action. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
				Computed:    true,
			},
			"security_enabled": schema.BoolAttribute{
				Description: "Specifies whether the group is a security group. Required. Returned by default. Supports $filter (eq, ne, not, in).",
				Computed:    true,
			},
			"security_identifier": schema.StringAttribute{
				Description: "Security identifier of the group, used in Windows scenarios. Returned by default.",
				Computed:    true,
			},
			"service_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors published by a federated service describing a non-transient, service-specific error regarding the properties or link from a group object .  Supports $filter (eq, not, for isResolved and serviceInstance).",
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
			"theme": schema.StringAttribute{
				Description: "Specifies a Microsoft 365 group's color theme. Possible values are Teal, Purple, Green, Blue, Pink, Orange or Red. Returned by default.",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or HiddenMembership. HiddenMembership can be set only for Microsoft 365 groups when the groups are created. It can't be updated later. Other values of visibility can be updated after group creation. If visibility value is not specified during group creation on Microsoft Graph, a security group is created as Private by default, and the Microsoft 365 group is Public. Groups assignable to roles are always Private. To learn more, see group visibility options. Returned by default. Nullable.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
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
