package groups

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
			"allow_external_senders": schema.BoolAttribute{
				Description: "Indicates if people external to the organization can send messages to the group. The default value is false. Returned only on $select. Supported only on the Get group API (GET /groups/{ID}).",
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
			"auto_subscribe_new_members": schema.BoolAttribute{
				Description: "Indicates if new members added to the group will be auto-subscribed to receive email notifications. You can set this property in a PATCH request for the group; do not set it in the initial POST request that creates the group. Default value is false. Returned only on $select. Supported only on the Get group API (GET /groups/{ID}).",
				Computed:    true,
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
			"hide_from_address_lists": schema.BoolAttribute{
				Description: "True if the group is not displayed in certain parts of the Outlook UI: the Address Book, address lists for selecting message recipients, and the Browse Groups dialog for searching groups; otherwise, false. Default value is false. Returned only on $select. Supported only on the Get group API (GET /groups/{ID}).",
				Computed:    true,
			},
			"hide_from_outlook_clients": schema.BoolAttribute{
				Description: "True if the group is not displayed in Outlook clients, such as Outlook for Windows and Outlook on the web; otherwise, false. The default value is false. Returned only on $select. Supported only on the Get group API (GET /groups/{ID}).",
				Computed:    true,
			},
			"is_assignable_to_role": schema.BoolAttribute{
				Description: "Indicates whether this group can be assigned to a Microsoft Entra role. Optional. This property can only be set while creating the group and is immutable. If set to true, the securityEnabled property must also be set to true, visibility must be Hidden, and the group cannot be a dynamic group (that is, groupTypes cannot contain DynamicMembership). Only callers in Global Administrator and Privileged Role Administrator roles can set this property. The caller must also be assigned the RoleManagement.ReadWrite.Directory permission to set this property or update the membership of such groups. For more, see Using a group to manage Microsoft Entra role assignmentsUsing this feature requires a Microsoft Entra ID P1 license. Returned by default. Supports $filter (eq, ne, not).",
				Computed:    true,
			},
			"is_subscribed_by_mail": schema.BoolAttribute{
				Description: "Indicates whether the signed-in user is subscribed to receive email conversations. The default value is true. Returned only on $select. Supported only on the Get group API (GET /groups/{ID}).",
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
			"theme": schema.StringAttribute{
				Description: "Specifies a Microsoft 365 group's color theme. Possible values are Teal, Purple, Green, Blue, Pink, Orange or Red. Returned by default.",
				Computed:    true,
			},
			"unseen_count": schema.Int64Attribute{
				Description: "Count of conversations that have received new posts since the signed-in user last visited the group. Returned only on $select. Supported only on the Get group API (GET /groups/{ID}).",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "Specifies the group join policy and group content visibility for groups. Possible values are: Private, Public, or HiddenMembership. HiddenMembership can be set only for Microsoft 365 groups when the groups are created. It can't be updated later. Other values of visibility can be updated after group creation. If visibility value is not specified during group creation on Microsoft Graph, a security group is created as Private by default, and the Microsoft 365 group is Public. Groups assignable to roles are always Private. To learn more, see group visibility options. Returned by default. Nullable.",
				Computed:    true,
			},
		},
	}
}

type groupDataSourceModel struct {
	Id                            types.String                                       `tfsdk:"id"`
	DeletedDateTime               types.String                                       `tfsdk:"deleted_date_time"`
	AllowExternalSenders          types.Bool                                         `tfsdk:"allow_external_senders"`
	AssignedLabels                []groupAssignedLabelsDataSourceModel               `tfsdk:"assigned_labels"`
	AssignedLicenses              []groupAssignedLicensesDataSourceModel             `tfsdk:"assigned_licenses"`
	AutoSubscribeNewMembers       types.Bool                                         `tfsdk:"auto_subscribe_new_members"`
	Classification                types.String                                       `tfsdk:"classification"`
	CreatedDateTime               types.String                                       `tfsdk:"created_date_time"`
	Description                   types.String                                       `tfsdk:"description"`
	DisplayName                   types.String                                       `tfsdk:"display_name"`
	ExpirationDateTime            types.String                                       `tfsdk:"expiration_date_time"`
	GroupTypes                    []types.String                                     `tfsdk:"group_types"`
	HideFromAddressLists          types.Bool                                         `tfsdk:"hide_from_address_lists"`
	HideFromOutlookClients        types.Bool                                         `tfsdk:"hide_from_outlook_clients"`
	IsAssignableToRole            types.Bool                                         `tfsdk:"is_assignable_to_role"`
	IsSubscribedByMail            types.Bool                                         `tfsdk:"is_subscribed_by_mail"`
	LicenseProcessingState        *groupLicenseProcessingStateDataSourceModel        `tfsdk:"license_processing_state"`
	Mail                          types.String                                       `tfsdk:"mail"`
	MailEnabled                   types.Bool                                         `tfsdk:"mail_enabled"`
	MailNickname                  types.String                                       `tfsdk:"mail_nickname"`
	MembershipRule                types.String                                       `tfsdk:"membership_rule"`
	MembershipRuleProcessingState types.String                                       `tfsdk:"membership_rule_processing_state"`
	OnPremisesDomainName          types.String                                       `tfsdk:"on_premises_domain_name"`
	OnPremisesLastSyncDateTime    types.String                                       `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesNetBiosName         types.String                                       `tfsdk:"on_premises_net_bios_name"`
	OnPremisesProvisioningErrors  []groupOnPremisesProvisioningErrorsDataSourceModel `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName      types.String                                       `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier  types.String                                       `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled         types.Bool                                         `tfsdk:"on_premises_sync_enabled"`
	PreferredDataLocation         types.String                                       `tfsdk:"preferred_data_location"`
	PreferredLanguage             types.String                                       `tfsdk:"preferred_language"`
	ProxyAddresses                []types.String                                     `tfsdk:"proxy_addresses"`
	RenewedDateTime               types.String                                       `tfsdk:"renewed_date_time"`
	SecurityEnabled               types.Bool                                         `tfsdk:"security_enabled"`
	SecurityIdentifier            types.String                                       `tfsdk:"security_identifier"`
	ServiceProvisioningErrors     []groupServiceProvisioningErrorsDataSourceModel    `tfsdk:"service_provisioning_errors"`
	Theme                         types.String                                       `tfsdk:"theme"`
	UnseenCount                   types.Int64                                        `tfsdk:"unseen_count"`
	Visibility                    types.String                                       `tfsdk:"visibility"`
}

type groupAssignedLabelsDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	LabelId     types.String `tfsdk:"label_id"`
}

type groupAssignedLicensesDataSourceModel struct {
	DisabledPlans []types.String `tfsdk:"disabled_plans"`
	SkuId         types.String   `tfsdk:"sku_id"`
}

type groupLicenseProcessingStateDataSourceModel struct {
	State types.String `tfsdk:"state"`
}

type groupOnPremisesProvisioningErrorsDataSourceModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

type groupServiceProvisioningErrorsDataSourceModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

// Read refreshes the Terraform state with the latest data.
func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state groupDataSourceModel
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
				"allowExternalSenders",
				"autoSubscribeNewMembers",
				"hideFromAddressLists",
				"hideFromOutlookClients",
				"isSubscribedByMail",
				"unseenCount",
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
			"`id` or `user_principal_name` must be supplied.",
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
	}
	if result.GetDeletedDateTime() != nil {
		state.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	}
	if result.GetAllowExternalSenders() != nil {
		state.AllowExternalSenders = types.BoolValue(*result.GetAllowExternalSenders())
	}
	for _, value := range result.GetAssignedLabels() {
		assignedLabels := new(groupAssignedLabelsDataSourceModel)

		if value.GetDisplayName() != nil {
			assignedLabels.DisplayName = types.StringValue(*value.GetDisplayName())
		}
		if value.GetLabelId() != nil {
			assignedLabels.LabelId = types.StringValue(*value.GetLabelId())
		}
		state.AssignedLabels = append(state.AssignedLabels, *assignedLabels)
	}
	for _, value := range result.GetAssignedLicenses() {
		assignedLicenses := new(groupAssignedLicensesDataSourceModel)

		for _, value := range value.GetDisabledPlans() {
			assignedLicenses.DisabledPlans = append(assignedLicenses.DisabledPlans, types.StringValue(value.String()))
		}
		if value.GetSkuId() != nil {
			assignedLicenses.SkuId = types.StringValue(value.GetSkuId().String())
		}
		state.AssignedLicenses = append(state.AssignedLicenses, *assignedLicenses)
	}
	if result.GetAutoSubscribeNewMembers() != nil {
		state.AutoSubscribeNewMembers = types.BoolValue(*result.GetAutoSubscribeNewMembers())
	}
	if result.GetClassification() != nil {
		state.Classification = types.StringValue(*result.GetClassification())
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	}
	if result.GetDescription() != nil {
		state.Description = types.StringValue(*result.GetDescription())
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	}
	if result.GetExpirationDateTime() != nil {
		state.ExpirationDateTime = types.StringValue(result.GetExpirationDateTime().String())
	}
	for _, value := range result.GetGroupTypes() {
		state.GroupTypes = append(state.GroupTypes, types.StringValue(value))
	}
	if result.GetHideFromAddressLists() != nil {
		state.HideFromAddressLists = types.BoolValue(*result.GetHideFromAddressLists())
	}
	if result.GetHideFromOutlookClients() != nil {
		state.HideFromOutlookClients = types.BoolValue(*result.GetHideFromOutlookClients())
	}
	if result.GetIsAssignableToRole() != nil {
		state.IsAssignableToRole = types.BoolValue(*result.GetIsAssignableToRole())
	}
	if result.GetIsSubscribedByMail() != nil {
		state.IsSubscribedByMail = types.BoolValue(*result.GetIsSubscribedByMail())
	}
	if result.GetLicenseProcessingState() != nil {
		state.LicenseProcessingState = new(groupLicenseProcessingStateDataSourceModel)

		if result.GetLicenseProcessingState().GetState() != nil {
			state.LicenseProcessingState.State = types.StringValue(*result.GetLicenseProcessingState().GetState())
		}
	}
	if result.GetMail() != nil {
		state.Mail = types.StringValue(*result.GetMail())
	}
	if result.GetMailEnabled() != nil {
		state.MailEnabled = types.BoolValue(*result.GetMailEnabled())
	}
	if result.GetMailNickname() != nil {
		state.MailNickname = types.StringValue(*result.GetMailNickname())
	}
	if result.GetMembershipRule() != nil {
		state.MembershipRule = types.StringValue(*result.GetMembershipRule())
	}
	if result.GetMembershipRuleProcessingState() != nil {
		state.MembershipRuleProcessingState = types.StringValue(*result.GetMembershipRuleProcessingState())
	}
	if result.GetOnPremisesDomainName() != nil {
		state.OnPremisesDomainName = types.StringValue(*result.GetOnPremisesDomainName())
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	}
	if result.GetOnPremisesNetBiosName() != nil {
		state.OnPremisesNetBiosName = types.StringValue(*result.GetOnPremisesNetBiosName())
	}
	for _, value := range result.GetOnPremisesProvisioningErrors() {
		onPremisesProvisioningErrors := new(groupOnPremisesProvisioningErrorsDataSourceModel)

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
	if result.GetPreferredDataLocation() != nil {
		state.PreferredDataLocation = types.StringValue(*result.GetPreferredDataLocation())
	}
	if result.GetPreferredLanguage() != nil {
		state.PreferredLanguage = types.StringValue(*result.GetPreferredLanguage())
	}
	for _, value := range result.GetProxyAddresses() {
		state.ProxyAddresses = append(state.ProxyAddresses, types.StringValue(value))
	}
	if result.GetRenewedDateTime() != nil {
		state.RenewedDateTime = types.StringValue(result.GetRenewedDateTime().String())
	}
	if result.GetSecurityEnabled() != nil {
		state.SecurityEnabled = types.BoolValue(*result.GetSecurityEnabled())
	}
	if result.GetSecurityIdentifier() != nil {
		state.SecurityIdentifier = types.StringValue(*result.GetSecurityIdentifier())
	}
	for _, value := range result.GetServiceProvisioningErrors() {
		serviceProvisioningErrors := new(groupServiceProvisioningErrorsDataSourceModel)

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
	if result.GetTheme() != nil {
		state.Theme = types.StringValue(*result.GetTheme())
	}
	if result.GetVisibility() != nil {
		state.Visibility = types.StringValue(*result.GetVisibility())
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
