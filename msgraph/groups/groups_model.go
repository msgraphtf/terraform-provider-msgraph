package groups

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupsModel struct {
	Value []groupsValueModel `tfsdk:"value"`
}

type groupsValueModel struct {
	Id                            types.String                              `tfsdk:"id"`
	DeletedDateTime               types.String                              `tfsdk:"deleted_date_time"`
	AssignedLabels                []groupsAssignedLabelsModel               `tfsdk:"assigned_labels"`
	AssignedLicenses              []groupsAssignedLicensesModel             `tfsdk:"assigned_licenses"`
	Classification                types.String                              `tfsdk:"classification"`
	CreatedDateTime               types.String                              `tfsdk:"created_date_time"`
	Description                   types.String                              `tfsdk:"description"`
	DisplayName                   types.String                              `tfsdk:"display_name"`
	ExpirationDateTime            types.String                              `tfsdk:"expiration_date_time"`
	GroupTypes                    []types.String                            `tfsdk:"group_types"`
	IsAssignableToRole            types.Bool                                `tfsdk:"is_assignable_to_role"`
	LicenseProcessingState        *groupsLicenseProcessingStateModel        `tfsdk:"license_processing_state"`
	Mail                          types.String                              `tfsdk:"mail"`
	MailEnabled                   types.Bool                                `tfsdk:"mail_enabled"`
	MailNickname                  types.String                              `tfsdk:"mail_nickname"`
	MembershipRule                types.String                              `tfsdk:"membership_rule"`
	MembershipRuleProcessingState types.String                              `tfsdk:"membership_rule_processing_state"`
	OnPremisesDomainName          types.String                              `tfsdk:"on_premises_domain_name"`
	OnPremisesLastSyncDateTime    types.String                              `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesNetBiosName         types.String                              `tfsdk:"on_premises_net_bios_name"`
	OnPremisesProvisioningErrors  []groupsOnPremisesProvisioningErrorsModel `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName      types.String                              `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier  types.String                              `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled         types.Bool                                `tfsdk:"on_premises_sync_enabled"`
	PreferredDataLocation         types.String                              `tfsdk:"preferred_data_location"`
	PreferredLanguage             types.String                              `tfsdk:"preferred_language"`
	ProxyAddresses                []types.String                            `tfsdk:"proxy_addresses"`
	RenewedDateTime               types.String                              `tfsdk:"renewed_date_time"`
	SecurityEnabled               types.Bool                                `tfsdk:"security_enabled"`
	SecurityIdentifier            types.String                              `tfsdk:"security_identifier"`
	ServiceProvisioningErrors     []groupsServiceProvisioningErrorsModel    `tfsdk:"service_provisioning_errors"`
	Theme                         types.String                              `tfsdk:"theme"`
	Visibility                    types.String                              `tfsdk:"visibility"`
}

type groupsAssignedLabelsModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	LabelId     types.String `tfsdk:"label_id"`
}

type groupsAssignedLicensesModel struct {
	DisabledPlans []types.String `tfsdk:"disabled_plans"`
	SkuId         types.String   `tfsdk:"sku_id"`
}

type groupsLicenseProcessingStateModel struct {
	State types.String `tfsdk:"state"`
}

type groupsOnPremisesProvisioningErrorsModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

type groupsServiceProvisioningErrorsModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}
