package groups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
