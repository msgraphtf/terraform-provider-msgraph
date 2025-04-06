package groups

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type groupModel struct {
	Id                            types.String `tfsdk:"id"`
	DeletedDateTime               types.String `tfsdk:"deleted_date_time"`
	AllowExternalSenders          types.Bool   `tfsdk:"allow_external_senders"`
	AssignedLabels                types.List   `tfsdk:"assigned_labels"`
	AssignedLicenses              types.List   `tfsdk:"assigned_licenses"`
	AutoSubscribeNewMembers       types.Bool   `tfsdk:"auto_subscribe_new_members"`
	Classification                types.String `tfsdk:"classification"`
	CreatedDateTime               types.String `tfsdk:"created_date_time"`
	Description                   types.String `tfsdk:"description"`
	DisplayName                   types.String `tfsdk:"display_name"`
	ExpirationDateTime            types.String `tfsdk:"expiration_date_time"`
	GroupTypes                    types.List   `tfsdk:"group_types"`
	HasMembersWithLicenseErrors   types.Bool   `tfsdk:"has_members_with_license_errors"`
	HideFromAddressLists          types.Bool   `tfsdk:"hide_from_address_lists"`
	HideFromOutlookClients        types.Bool   `tfsdk:"hide_from_outlook_clients"`
	IsArchived                    types.Bool   `tfsdk:"is_archived"`
	IsAssignableToRole            types.Bool   `tfsdk:"is_assignable_to_role"`
	IsManagementRestricted        types.Bool   `tfsdk:"is_management_restricted"`
	IsSubscribedByMail            types.Bool   `tfsdk:"is_subscribed_by_mail"`
	LicenseProcessingState        types.Object `tfsdk:"license_processing_state"`
	Mail                          types.String `tfsdk:"mail"`
	MailEnabled                   types.Bool   `tfsdk:"mail_enabled"`
	MailNickname                  types.String `tfsdk:"mail_nickname"`
	MembershipRule                types.String `tfsdk:"membership_rule"`
	MembershipRuleProcessingState types.String `tfsdk:"membership_rule_processing_state"`
	OnPremisesDomainName          types.String `tfsdk:"on_premises_domain_name"`
	OnPremisesLastSyncDateTime    types.String `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesNetBiosName         types.String `tfsdk:"on_premises_net_bios_name"`
	OnPremisesProvisioningErrors  types.List   `tfsdk:"on_premises_provisioning_errors"`
	OnPremisesSamAccountName      types.String `tfsdk:"on_premises_sam_account_name"`
	OnPremisesSecurityIdentifier  types.String `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled         types.Bool   `tfsdk:"on_premises_sync_enabled"`
	PreferredDataLocation         types.String `tfsdk:"preferred_data_location"`
	PreferredLanguage             types.String `tfsdk:"preferred_language"`
	ProxyAddresses                types.List   `tfsdk:"proxy_addresses"`
	RenewedDateTime               types.String `tfsdk:"renewed_date_time"`
	SecurityEnabled               types.Bool   `tfsdk:"security_enabled"`
	SecurityIdentifier            types.String `tfsdk:"security_identifier"`
	ServiceProvisioningErrors     types.List   `tfsdk:"service_provisioning_errors"`
	Theme                         types.String `tfsdk:"theme"`
	UniqueName                    types.String `tfsdk:"unique_name"`
	UnseenCount                   types.Int64  `tfsdk:"unseen_count"`
	Visibility                    types.String `tfsdk:"visibility"`
}

func (m groupModel) AttributeTypes() map[string]attr.Type {
	groupAssignedLabels := groupAssignedLabelModel{}
	groupAssignedLicenses := groupAssignedLicenseModel{}
	groupLicenseProcessingState := groupLicenseProcessingStateModel{}
	groupOnPremisesProvisioningErrors := groupOnPremisesProvisioningErrorModel{}
	groupServiceProvisioningErrors := groupServiceProvisioningErrorModel{}
	return map[string]attr.Type{
		"id":                               types.StringType,
		"deleted_date_time":                types.StringType,
		"allow_external_senders":           types.BoolType,
		"assigned_labels":                  types.ListType{ElemType: types.ObjectType{AttrTypes: groupAssignedLabels.AttributeTypes()}},
		"assigned_licenses":                types.ListType{ElemType: types.ObjectType{AttrTypes: groupAssignedLicenses.AttributeTypes()}},
		"auto_subscribe_new_members":       types.BoolType,
		"classification":                   types.StringType,
		"created_date_time":                types.StringType,
		"description":                      types.StringType,
		"display_name":                     types.StringType,
		"expiration_date_time":             types.StringType,
		"group_types":                      types.ListType{ElemType: types.StringType},
		"has_members_with_license_errors":  types.BoolType,
		"hide_from_address_lists":          types.BoolType,
		"hide_from_outlook_clients":        types.BoolType,
		"is_archived":                      types.BoolType,
		"is_assignable_to_role":            types.BoolType,
		"is_management_restricted":         types.BoolType,
		"is_subscribed_by_mail":            types.BoolType,
		"license_processing_state":         types.ObjectType{AttrTypes: groupLicenseProcessingState.AttributeTypes()},
		"mail":                             types.StringType,
		"mail_enabled":                     types.BoolType,
		"mail_nickname":                    types.StringType,
		"membership_rule":                  types.StringType,
		"membership_rule_processing_state": types.StringType,
		"on_premises_domain_name":          types.StringType,
		"on_premises_last_sync_date_time":  types.StringType,
		"on_premises_net_bios_name":        types.StringType,
		"on_premises_provisioning_errors":  types.ListType{ElemType: types.ObjectType{AttrTypes: groupOnPremisesProvisioningErrors.AttributeTypes()}},
		"on_premises_sam_account_name":     types.StringType,
		"on_premises_security_identifier":  types.StringType,
		"on_premises_sync_enabled":         types.BoolType,
		"preferred_data_location":          types.StringType,
		"preferred_language":               types.StringType,
		"proxy_addresses":                  types.ListType{ElemType: types.StringType},
		"renewed_date_time":                types.StringType,
		"security_enabled":                 types.BoolType,
		"security_identifier":              types.StringType,
		"service_provisioning_errors":      types.ListType{ElemType: types.ObjectType{AttrTypes: groupServiceProvisioningErrors.AttributeTypes()}},
		"theme":                            types.StringType,
		"unique_name":                      types.StringType,
		"unseen_count":                     types.Int64Type,
		"visibility":                       types.StringType,
	}
}

type groupAssignedLabelModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	LabelId     types.String `tfsdk:"label_id"`
}

func (m groupAssignedLabelModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"display_name": types.StringType,
		"label_id":     types.StringType,
	}
}

type groupAssignedLicenseModel struct {
	DisabledPlans types.List   `tfsdk:"disabled_plans"`
	SkuId         types.String `tfsdk:"sku_id"`
}

func (m groupAssignedLicenseModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"disabled_plans": types.ListType{ElemType: types.StringType},
		"sku_id":         types.StringType,
	}
}

type groupLicenseProcessingStateModel struct {
	State types.String `tfsdk:"state"`
}

func (m groupLicenseProcessingStateModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"state": types.StringType,
	}
}

type groupOnPremisesProvisioningErrorModel struct {
	Category             types.String `tfsdk:"category"`
	OccurredDateTime     types.String `tfsdk:"occurred_date_time"`
	PropertyCausingError types.String `tfsdk:"property_causing_error"`
	Value                types.String `tfsdk:"value"`
}

func (m groupOnPremisesProvisioningErrorModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"category":               types.StringType,
		"occurred_date_time":     types.StringType,
		"property_causing_error": types.StringType,
		"value":                  types.StringType,
	}
}

type groupServiceProvisioningErrorModel struct {
	CreatedDateTime types.String `tfsdk:"created_date_time"`
	IsResolved      types.Bool   `tfsdk:"is_resolved"`
	ServiceInstance types.String `tfsdk:"service_instance"`
}

func (m groupServiceProvisioningErrorModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"created_date_time": types.StringType,
		"is_resolved":       types.BoolType,
		"service_instance":  types.StringType,
	}
}
