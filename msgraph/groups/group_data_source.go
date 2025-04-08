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
				Description: "The list of sensitivity label pairs (label ID, label name) associated with a Microsoft 365 group. Returned only on $select. This property can be updated only in delegated scenarios where the caller requires both the Microsoft Graph permission and a supported administrator role.",
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
							Description: "A collection of the unique identifiers for plans that have been disabled. IDs are available in servicePlans > servicePlanId in the tenant's subscribedSkus or serviceStatus > servicePlanId in the tenant's companySubscription.",
							Computed:    true,
							ElementType: types.StringType,
						},
						"sku_id": schema.StringAttribute{
							Description: "The unique identifier for the SKU. Corresponds to the skuId from subscribedSkus or companySubscription.",
							Computed:    true,
						},
					},
				},
			},
			"classification": schema.StringAttribute{
				Description: "Describes a classification for the group (such as low, medium, or high business impact). Valid values for this property are defined by creating a ClassificationList setting value, based on the template definition.Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith).",
				Computed:    true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group was created. The value can't be modified and is automatically populated when the group is created. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "An optional description for the group. Returned by default. Supports $filter (eq, ne, not, ge, le, startsWith) and $search.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name for the group. This property is required when a group is created and can't be cleared during updates. Maximum length is 256 characters. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Computed:    true,
			},
			"expiration_date_time": schema.StringAttribute{
				Description: "Timestamp of when the group is set to expire. It's null for security groups, but for Microsoft 365 groups, it represents when the group is set to expire as defined in the groupLifecyclePolicy. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
				Computed:    true,
			},
			"group_types": schema.ListAttribute{
				Description: "Specifies the group type and its membership. If the collection contains Unified, the group is a Microsoft 365 group; otherwise, it's either a security group or a distribution group. For details, see groups overview.If the collection includes DynamicMembership, the group has dynamic membership; otherwise, membership is static. Returned by default. Supports $filter (eq, not).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"is_assignable_to_role": schema.BoolAttribute{
				Description: "Indicates whether this group can be assigned to a Microsoft Entra role. Optional. This property can only be set while creating the group and is immutable. If set to true, the securityEnabled property must also be set to true, visibility must be Hidden, and the group can't be a dynamic group (that is, groupTypes can't contain DynamicMembership). Only callers with at least the Privileged Role Administrator role can set this property. The caller must also be assigned the RoleManagement.ReadWrite.Directory permission to set this property or update the membership of such groups. For more, see Using a group to manage Microsoft Entra role assignmentsUsing this feature requires a Microsoft Entra ID P1 license. Returned by default. Supports $filter (eq, ne, not).",
				Computed:    true,
			},
			"is_management_restricted": schema.BoolAttribute{
				Description: "",
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
				Description: "The SMTP address for the group, for example, 'serviceadmins@contoso.com'. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
				Computed:    true,
			},
			"mail_enabled": schema.BoolAttribute{
				Description: "Specifies whether the group is mail-enabled. Required. Returned by default. Supports $filter (eq, ne, not).",
				Computed:    true,
			},
			"mail_nickname": schema.StringAttribute{
				Description: "The mail alias for the group, unique for Microsoft 365 groups in the organization. Maximum length is 64 characters. This property can contain only characters in the ASCII character set 0 - 127 except the following characters: @ () / [] ' ; : <> , SPACE. Required. Returned by default. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values).",
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
				Description: "Contains the on-premises domain FQDN, also called dnsDomainName synchronized from the on-premises directory. The property is only populated for customers synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.Returned by default. Read-only.",
				Computed:    true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "Indicates the last time at which the group was synced with the on-premises directory. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Read-only. Supports $filter (eq, ne, not, ge, le, in).",
				Computed:    true,
			},
			"on_premises_net_bios_name": schema.StringAttribute{
				Description: "Contains the on-premises netBios name synchronized from the on-premises directory. The property is only populated for customers synchronizing their on-premises directory to Microsoft Entra ID via Microsoft Entra Connect.Returned by default. Read-only.",
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
				Description: "Contains the on-premises security identifier (SID) for the group synchronized from on-premises to the cloud. Read-only. Returned by default. Supports $filter (eq including on null values).",
				Computed:    true,
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Description: "true if this group is synced from an on-premises directory; false if this group was originally synced from an on-premises directory but is no longer synced; null if this object has never synced from an on-premises directory (default). Returned by default. Read-only. Supports $filter (eq, ne, not, in, and eq on null values).",
				Computed:    true,
			},
			"preferred_data_location": schema.StringAttribute{
				Description: "The preferred data location for the Microsoft 365 group. By default, the group inherits the group creator's preferred data location. To set this property, the calling app must be granted the Directory.ReadWrite.All permission and the user be assigned at least one of the following Microsoft Entra roles: User Account Administrator Directory Writer  Exchange Administrator  SharePoint Administrator  For more information about this property, see OneDrive Online Multi-Geo. Nullable. Returned by default.",
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
				Description: "Timestamp of when the group was last renewed. This value can't be modified directly and is only updated via the renew service action. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC. For example, midnight UTC on January 1, 2014 is 2014-01-01T00:00:00Z. Returned by default. Supports $filter (eq, ne, not, ge, le, in). Read-only.",
				Computed:    true,
			},
			"security_enabled": schema.BoolAttribute{
				Description: "Specifies whether the group is a security group. Required. Returned by default. Supports $filter (eq, ne, not, in).",
				Computed:    true,
			},
			"security_identifier": schema.StringAttribute{
				Description: "Security identifier of the group, used in Windows scenarios. Read-only. Returned by default.",
				Computed:    true,
			},
			"service_provisioning_errors": schema.ListNestedAttribute{
				Description: "Errors published by a federated service describing a nontransient, service-specific error regarding the properties or link from a group object.  Supports $filter (eq, not, for isResolved and serviceInstance).",
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
				Description: "Specifies a Microsoft 365 group's color theme. Possible values are Teal, Purple, Green, Blue, Pink, Orange, or Red. Returned by default.",
				Computed:    true,
			},
			"unique_name": schema.StringAttribute{
				Description: "The unique identifier that can be assigned to a group and used as an alternate key. Immutable. Read-only.",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or HiddenMembership. HiddenMembership can be set only for Microsoft 365 groups when the groups are created. It can't be updated later. Other values of visibility can be updated after group creation. If visibility value isn't specified during group creation on Microsoft Graph, a security group is created as Private by default, and the Microsoft 365 group is Public. Groups assignable to roles are always Private. To learn more, see group visibility options. Returned by default. Nullable.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfStateGroup groupModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateGroup)...)
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

	if !tfStateGroup.Id.IsNull() {
		result, err = d.client.Groups().ByGroupId(tfStateGroup.Id.ValueString()).Get(context.Background(), &qparams)
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
		tfStateGroup.Id = types.StringValue(*result.GetId())
	} else {
		tfStateGroup.Id = types.StringNull()
	}
	if result.GetDeletedDateTime() != nil {
		tfStateGroup.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	} else {
		tfStateGroup.DeletedDateTime = types.StringNull()
	}
	if len(result.GetAssignedLabels()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedLabels() {
			tfStateAssignedLabel := groupAssignedLabelModel{}

			if v.GetDisplayName() != nil {
				tfStateAssignedLabel.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				tfStateAssignedLabel.DisplayName = types.StringNull()
			}
			if v.GetLabelId() != nil {
				tfStateAssignedLabel.LabelId = types.StringValue(*v.GetLabelId())
			} else {
				tfStateAssignedLabel.LabelId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedLabel.AttributeTypes(), tfStateAssignedLabel)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.AssignedLabels, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if len(result.GetAssignedLicenses()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAssignedLicenses() {
			tfStateAssignedLicense := groupAssignedLicenseModel{}

			if len(v.GetDisabledPlans()) > 0 {
				var valueArrayDisabledPlans []attr.Value
				for _, resultDisabledPlans := range v.GetDisabledPlans() {
					valueArrayDisabledPlans = append(valueArrayDisabledPlans, types.StringValue(resultDisabledPlans.String()))
				}
				tfStateAssignedLicense.DisabledPlans, _ = types.ListValue(types.StringType, valueArrayDisabledPlans)
			} else {
				tfStateAssignedLicense.DisabledPlans = types.ListNull(types.StringType)
			}
			if v.GetSkuId() != nil {
				tfStateAssignedLicense.SkuId = types.StringValue(v.GetSkuId().String())
			} else {
				tfStateAssignedLicense.SkuId = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAssignedLicense.AttributeTypes(), tfStateAssignedLicense)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.AssignedLicenses, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetClassification() != nil {
		tfStateGroup.Classification = types.StringValue(*result.GetClassification())
	} else {
		tfStateGroup.Classification = types.StringNull()
	}
	if result.GetCreatedDateTime() != nil {
		tfStateGroup.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		tfStateGroup.CreatedDateTime = types.StringNull()
	}
	if result.GetDescription() != nil {
		tfStateGroup.Description = types.StringValue(*result.GetDescription())
	} else {
		tfStateGroup.Description = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		tfStateGroup.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		tfStateGroup.DisplayName = types.StringNull()
	}
	if result.GetExpirationDateTime() != nil {
		tfStateGroup.ExpirationDateTime = types.StringValue(result.GetExpirationDateTime().String())
	} else {
		tfStateGroup.ExpirationDateTime = types.StringNull()
	}
	if len(result.GetGroupTypes()) > 0 {
		var valueArrayGroupTypes []attr.Value
		for _, resultGroupTypes := range result.GetGroupTypes() {
			valueArrayGroupTypes = append(valueArrayGroupTypes, types.StringValue(resultGroupTypes))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayGroupTypes)
		tfStateGroup.GroupTypes = listValue
	} else {
		tfStateGroup.GroupTypes = types.ListNull(types.StringType)
	}
	if result.GetIsAssignableToRole() != nil {
		tfStateGroup.IsAssignableToRole = types.BoolValue(*result.GetIsAssignableToRole())
	} else {
		tfStateGroup.IsAssignableToRole = types.BoolNull()
	}
	if result.GetIsManagementRestricted() != nil {
		tfStateGroup.IsManagementRestricted = types.BoolValue(*result.GetIsManagementRestricted())
	} else {
		tfStateGroup.IsManagementRestricted = types.BoolNull()
	}
	if result.GetLicenseProcessingState() != nil {
		tfStateLicenseProcessingState := groupLicenseProcessingStateModel{}

		if result.GetLicenseProcessingState().GetState() != nil {
			tfStateLicenseProcessingState.State = types.StringValue(*result.GetLicenseProcessingState().GetState())
		} else {
			tfStateLicenseProcessingState.State = types.StringNull()
		}

		tfStateGroup.LicenseProcessingState, _ = types.ObjectValueFrom(ctx, tfStateLicenseProcessingState.AttributeTypes(), tfStateLicenseProcessingState)
	}
	if result.GetMail() != nil {
		tfStateGroup.Mail = types.StringValue(*result.GetMail())
	} else {
		tfStateGroup.Mail = types.StringNull()
	}
	if result.GetMailEnabled() != nil {
		tfStateGroup.MailEnabled = types.BoolValue(*result.GetMailEnabled())
	} else {
		tfStateGroup.MailEnabled = types.BoolNull()
	}
	if result.GetMailNickname() != nil {
		tfStateGroup.MailNickname = types.StringValue(*result.GetMailNickname())
	} else {
		tfStateGroup.MailNickname = types.StringNull()
	}
	if result.GetMembershipRule() != nil {
		tfStateGroup.MembershipRule = types.StringValue(*result.GetMembershipRule())
	} else {
		tfStateGroup.MembershipRule = types.StringNull()
	}
	if result.GetMembershipRuleProcessingState() != nil {
		tfStateGroup.MembershipRuleProcessingState = types.StringValue(*result.GetMembershipRuleProcessingState())
	} else {
		tfStateGroup.MembershipRuleProcessingState = types.StringNull()
	}
	if result.GetOnPremisesDomainName() != nil {
		tfStateGroup.OnPremisesDomainName = types.StringValue(*result.GetOnPremisesDomainName())
	} else {
		tfStateGroup.OnPremisesDomainName = types.StringNull()
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		tfStateGroup.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	} else {
		tfStateGroup.OnPremisesLastSyncDateTime = types.StringNull()
	}
	if result.GetOnPremisesNetBiosName() != nil {
		tfStateGroup.OnPremisesNetBiosName = types.StringValue(*result.GetOnPremisesNetBiosName())
	} else {
		tfStateGroup.OnPremisesNetBiosName = types.StringNull()
	}
	if len(result.GetOnPremisesProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetOnPremisesProvisioningErrors() {
			tfStateOnPremisesProvisioningError := groupOnPremisesProvisioningErrorModel{}

			if v.GetCategory() != nil {
				tfStateOnPremisesProvisioningError.Category = types.StringValue(*v.GetCategory())
			} else {
				tfStateOnPremisesProvisioningError.Category = types.StringNull()
			}
			if v.GetOccurredDateTime() != nil {
				tfStateOnPremisesProvisioningError.OccurredDateTime = types.StringValue(v.GetOccurredDateTime().String())
			} else {
				tfStateOnPremisesProvisioningError.OccurredDateTime = types.StringNull()
			}
			if v.GetPropertyCausingError() != nil {
				tfStateOnPremisesProvisioningError.PropertyCausingError = types.StringValue(*v.GetPropertyCausingError())
			} else {
				tfStateOnPremisesProvisioningError.PropertyCausingError = types.StringNull()
			}
			if v.GetValue() != nil {
				tfStateOnPremisesProvisioningError.Value = types.StringValue(*v.GetValue())
			} else {
				tfStateOnPremisesProvisioningError.Value = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateOnPremisesProvisioningError.AttributeTypes(), tfStateOnPremisesProvisioningError)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.OnPremisesProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetOnPremisesSamAccountName() != nil {
		tfStateGroup.OnPremisesSamAccountName = types.StringValue(*result.GetOnPremisesSamAccountName())
	} else {
		tfStateGroup.OnPremisesSamAccountName = types.StringNull()
	}
	if result.GetOnPremisesSecurityIdentifier() != nil {
		tfStateGroup.OnPremisesSecurityIdentifier = types.StringValue(*result.GetOnPremisesSecurityIdentifier())
	} else {
		tfStateGroup.OnPremisesSecurityIdentifier = types.StringNull()
	}
	if result.GetOnPremisesSyncEnabled() != nil {
		tfStateGroup.OnPremisesSyncEnabled = types.BoolValue(*result.GetOnPremisesSyncEnabled())
	} else {
		tfStateGroup.OnPremisesSyncEnabled = types.BoolNull()
	}
	if result.GetPreferredDataLocation() != nil {
		tfStateGroup.PreferredDataLocation = types.StringValue(*result.GetPreferredDataLocation())
	} else {
		tfStateGroup.PreferredDataLocation = types.StringNull()
	}
	if result.GetPreferredLanguage() != nil {
		tfStateGroup.PreferredLanguage = types.StringValue(*result.GetPreferredLanguage())
	} else {
		tfStateGroup.PreferredLanguage = types.StringNull()
	}
	if len(result.GetProxyAddresses()) > 0 {
		var valueArrayProxyAddresses []attr.Value
		for _, resultProxyAddresses := range result.GetProxyAddresses() {
			valueArrayProxyAddresses = append(valueArrayProxyAddresses, types.StringValue(resultProxyAddresses))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayProxyAddresses)
		tfStateGroup.ProxyAddresses = listValue
	} else {
		tfStateGroup.ProxyAddresses = types.ListNull(types.StringType)
	}
	if result.GetRenewedDateTime() != nil {
		tfStateGroup.RenewedDateTime = types.StringValue(result.GetRenewedDateTime().String())
	} else {
		tfStateGroup.RenewedDateTime = types.StringNull()
	}
	if result.GetSecurityEnabled() != nil {
		tfStateGroup.SecurityEnabled = types.BoolValue(*result.GetSecurityEnabled())
	} else {
		tfStateGroup.SecurityEnabled = types.BoolNull()
	}
	if result.GetSecurityIdentifier() != nil {
		tfStateGroup.SecurityIdentifier = types.StringValue(*result.GetSecurityIdentifier())
	} else {
		tfStateGroup.SecurityIdentifier = types.StringNull()
	}
	if len(result.GetServiceProvisioningErrors()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetServiceProvisioningErrors() {
			tfStateServiceProvisioningError := groupServiceProvisioningErrorModel{}

			if v.GetCreatedDateTime() != nil {
				tfStateServiceProvisioningError.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			} else {
				tfStateServiceProvisioningError.CreatedDateTime = types.StringNull()
			}
			if v.GetIsResolved() != nil {
				tfStateServiceProvisioningError.IsResolved = types.BoolValue(*v.GetIsResolved())
			} else {
				tfStateServiceProvisioningError.IsResolved = types.BoolNull()
			}
			if v.GetServiceInstance() != nil {
				tfStateServiceProvisioningError.ServiceInstance = types.StringValue(*v.GetServiceInstance())
			} else {
				tfStateServiceProvisioningError.ServiceInstance = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateServiceProvisioningError.AttributeTypes(), tfStateServiceProvisioningError)
			objectValues = append(objectValues, objectValue)
		}
		tfStateGroup.ServiceProvisioningErrors, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetTheme() != nil {
		tfStateGroup.Theme = types.StringValue(*result.GetTheme())
	} else {
		tfStateGroup.Theme = types.StringNull()
	}
	if result.GetUniqueName() != nil {
		tfStateGroup.UniqueName = types.StringValue(*result.GetUniqueName())
	} else {
		tfStateGroup.UniqueName = types.StringNull()
	}
	if result.GetVisibility() != nil {
		tfStateGroup.Visibility = types.StringValue(*result.GetVisibility())
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
