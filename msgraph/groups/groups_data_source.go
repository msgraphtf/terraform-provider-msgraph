package groups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/groups"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &groupsDataSource{}
	_ datasource.DataSourceWithConfigure = &groupsDataSource{}
)

// NewGroupsDataSource is a helper function to simplify the provider implementation.
func NewGroupsDataSource() datasource.DataSource {
	return &groupsDataSource{}
}

// groupsDataSource is the data source implementation.
type groupsDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *groupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groups"
}

// Configure adds the provider configured client to the data source.
func (d *groupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *groupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"value": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier for an entity. Read-only.",
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
				},
			},
		},
	}
}

type groupsDataSourceModel struct {
	Value []groupsValueDataSourceModel `tfsdk:"value"`
}

type groupsValueDataSourceModel struct {
	Id                            types.String                                        `tfsdk:"id"`
	DeletedDateTime               types.String                                        `tfsdk:"deleted_date_time"`
	AssignedLabels                []groupsAssignedLabelsDataSourceModel               `tfsdk:"assigned_labels"`
	AssignedLicenses              []groupsAssignedLicensesDataSourceModel             `tfsdk:"assigned_licenses"`
	Classification                types.String                                        `tfsdk:"classification"`
	CreatedDateTime               types.String                                        `tfsdk:"created_date_time"`
	Description                   types.String                                        `tfsdk:"description"`
	DisplayName                   types.String                                        `tfsdk:"display_name"`
	ExpirationDateTime            types.String                                        `tfsdk:"expiration_date_time"`
	GroupTypes                    []types.String                                      `tfsdk:"group_types"`
	IsAssignableToRole            types.Bool                                          `tfsdk:"is_assignable_to_role"`
	LicenseProcessingState        *groupsLicenseProcessingStateDataSourceModel        `tfsdk:"license_processing_state"`
	Mail                          types.String                                        `tfsdk:"mail"`
	MailEnabled                   types.Bool                                          `tfsdk:"mail_enabled"`
	MailNickname                  types.String                                        `tfsdk:"mail_nickname"`
	MembershipRule                types.String                                        `tfsdk:"membership_rule"`
	MembershipRuleProcessingState types.String                                        `tfsdk:"membership_rule_processing_state"`
	OnPremisesDomainName          types.String                                        `tfsdk:"on_premises_domain_name"`
	OnPremisesLastSyncDateTime    types.String                                        `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesNetBiosName         types.String                                        `tfsdk:"on_premises_net_bios_name"`
	OnPremisesProvisioningErrors  []groupsOnPremisesProvisioningErrorsDataSourceModel `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName      types.String                                        `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier  types.String                                        `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled         types.Bool                                          `tfsdk:"on_premises_sync_enabled"`
	PreferredDataLocation         types.String                                        `tfsdk:"preferred_data_location"`
	PreferredLanguage             types.String                                        `tfsdk:"preferred_language"`
	ProxyAddresses                []types.String                                      `tfsdk:"proxy_addresses"`
	RenewedDateTime               types.String                                        `tfsdk:"renewed_date_time"`
	SecurityEnabled               types.Bool                                          `tfsdk:"security_enabled"`
	SecurityIdentifier            types.String                                        `tfsdk:"security_identifier"`
	ServiceProvisioningErrors     []groupsServiceProvisioningErrorsDataSourceModel    `tfsdk:"service_provisioning_errors"`
	Theme                         types.String                                        `tfsdk:"theme"`
	Visibility                    types.String                                        `tfsdk:"visibility"`
}

type groupsAssignedLabelsDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	LabelId     types.String `tfsdk:"label_id"`
}

type groupsAssignedLicensesDataSourceModel struct {
	DisabledPlans []types.String `tfsdk:"disabled_plans"`
	SkuId         types.String   `tfsdk:"sku_id"`
}

type groupsLicenseProcessingStateDataSourceModel struct {
	State types.String `tfsdk:"state"`
}

type groupsOnPremisesProvisioningErrorsDataSourceModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

type groupsServiceProvisioningErrorsDataSourceModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

// Read refreshes the Terraform state with the latest data.
func (d *groupsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupsRequestBuilderGetQueryParameters{
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

	result, err := d.client.Groups().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting groups",
			err.Error(),
		)
		return
	}

	for _, v := range result.GetValue() {
		value := new(groupsValueDataSourceModel)

		if v.GetId() != nil {
			value.Id = types.StringValue(*v.GetId())
		}
		if v.GetDeletedDateTime() != nil {
			value.DeletedDateTime = types.StringValue(v.GetDeletedDateTime().String())
		}
		for _, v := range v.GetAssignedLabels() {
			assignedLabels := new(groupsAssignedLabelsDataSourceModel)

			if v.GetDisplayName() != nil {
				assignedLabels.DisplayName = types.StringValue(*v.GetDisplayName())
			}
			if v.GetLabelId() != nil {
				assignedLabels.LabelId = types.StringValue(*v.GetLabelId())
			}
			value.AssignedLabels = append(value.AssignedLabels, *assignedLabels)
		}
		for _, v := range v.GetAssignedLicenses() {
			assignedLicenses := new(groupsAssignedLicensesDataSourceModel)

			for _, v := range v.GetDisabledPlans() {
				assignedLicenses.DisabledPlans = append(assignedLicenses.DisabledPlans, types.StringValue(v.String()))
			}
			if v.GetSkuId() != nil {
				assignedLicenses.SkuId = types.StringValue(v.GetSkuId().String())
			}
			value.AssignedLicenses = append(value.AssignedLicenses, *assignedLicenses)
		}
		if v.GetClassification() != nil {
			value.Classification = types.StringValue(*v.GetClassification())
		}
		if v.GetCreatedDateTime() != nil {
			value.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
		}
		if v.GetDescription() != nil {
			value.Description = types.StringValue(*v.GetDescription())
		}
		if v.GetDisplayName() != nil {
			value.DisplayName = types.StringValue(*v.GetDisplayName())
		}
		if v.GetExpirationDateTime() != nil {
			value.ExpirationDateTime = types.StringValue(v.GetExpirationDateTime().String())
		}
		for _, v := range v.GetGroupTypes() {
			value.GroupTypes = append(value.GroupTypes, types.StringValue(v))
		}
		if v.GetIsAssignableToRole() != nil {
			value.IsAssignableToRole = types.BoolValue(*v.GetIsAssignableToRole())
		}
		if v.GetLicenseProcessingState() != nil {
			value.LicenseProcessingState = new(groupsLicenseProcessingStateDataSourceModel)

			if v.GetLicenseProcessingState().GetState() != nil {
				value.LicenseProcessingState.State = types.StringValue(*v.GetLicenseProcessingState().GetState())
			}
		}
		if v.GetMail() != nil {
			value.Mail = types.StringValue(*v.GetMail())
		}
		if v.GetMailEnabled() != nil {
			value.MailEnabled = types.BoolValue(*v.GetMailEnabled())
		}
		if v.GetMailNickname() != nil {
			value.MailNickname = types.StringValue(*v.GetMailNickname())
		}
		if v.GetMembershipRule() != nil {
			value.MembershipRule = types.StringValue(*v.GetMembershipRule())
		}
		if v.GetMembershipRuleProcessingState() != nil {
			value.MembershipRuleProcessingState = types.StringValue(*v.GetMembershipRuleProcessingState())
		}
		if v.GetOnPremisesDomainName() != nil {
			value.OnPremisesDomainName = types.StringValue(*v.GetOnPremisesDomainName())
		}
		if v.GetOnPremisesLastSyncDateTime() != nil {
			value.OnPremisesLastSyncDateTime = types.StringValue(v.GetOnPremisesLastSyncDateTime().String())
		}
		if v.GetOnPremisesNetBiosName() != nil {
			value.OnPremisesNetBiosName = types.StringValue(*v.GetOnPremisesNetBiosName())
		}
		for _, v := range v.GetOnPremisesProvisioningErrors() {
			onPremisesProvisioningErrors := new(groupsOnPremisesProvisioningErrorsDataSourceModel)

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
			value.OnPremisesProvisioningErrors = append(value.OnPremisesProvisioningErrors, *onPremisesProvisioningErrors)
		}
		if v.GetOnPremisesSamAccountName() != nil {
			value.OnPremisesSamAccountName = types.StringValue(*v.GetOnPremisesSamAccountName())
		}
		if v.GetOnPremisesSecurityIdentifier() != nil {
			value.OnPremisesSecurityIdentifier = types.StringValue(*v.GetOnPremisesSecurityIdentifier())
		}
		if v.GetOnPremisesSyncEnabled() != nil {
			value.OnPremisesSyncEnabled = types.BoolValue(*v.GetOnPremisesSyncEnabled())
		}
		if v.GetPreferredDataLocation() != nil {
			value.PreferredDataLocation = types.StringValue(*v.GetPreferredDataLocation())
		}
		if v.GetPreferredLanguage() != nil {
			value.PreferredLanguage = types.StringValue(*v.GetPreferredLanguage())
		}
		for _, v := range v.GetProxyAddresses() {
			value.ProxyAddresses = append(value.ProxyAddresses, types.StringValue(v))
		}
		if v.GetRenewedDateTime() != nil {
			value.RenewedDateTime = types.StringValue(v.GetRenewedDateTime().String())
		}
		if v.GetSecurityEnabled() != nil {
			value.SecurityEnabled = types.BoolValue(*v.GetSecurityEnabled())
		}
		if v.GetSecurityIdentifier() != nil {
			value.SecurityIdentifier = types.StringValue(*v.GetSecurityIdentifier())
		}
		for _, v := range v.GetServiceProvisioningErrors() {
			serviceProvisioningErrors := new(groupsServiceProvisioningErrorsDataSourceModel)

			if v.GetCreatedDateTime() != nil {
				serviceProvisioningErrors.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			}
			if v.GetIsResolved() != nil {
				serviceProvisioningErrors.IsResolved = types.BoolValue(*v.GetIsResolved())
			}
			if v.GetServiceInstance() != nil {
				serviceProvisioningErrors.ServiceInstance = types.StringValue(*v.GetServiceInstance())
			}
			value.ServiceProvisioningErrors = append(value.ServiceProvisioningErrors, *serviceProvisioningErrors)
		}
		if v.GetTheme() != nil {
			value.Theme = types.StringValue(*v.GetTheme())
		}
		if v.GetVisibility() != nil {
			value.Visibility = types.StringValue(*v.GetVisibility())
		}
		state.Value = append(state.Value, *value)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
