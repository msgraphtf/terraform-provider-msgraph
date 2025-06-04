package devices

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type devicesModel struct {
	Value types.List `tfsdk:"value"`
}

func (m devicesModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.ListType{ElemType: types.ObjectType{AttrTypes: devicesDeviceModel{}.AttributeTypes()}},
	}
}

type devicesDeviceModel struct {
	AccountEnabled                types.Bool   `tfsdk:"account_enabled"`
	AlternativeSecurityIds        types.List   `tfsdk:"alternative_security_ids"`
	ApproximateLastSignInDateTime types.String `tfsdk:"approximate_last_sign_in_date_time"`
	ComplianceExpirationDateTime  types.String `tfsdk:"compliance_expiration_date_time"`
	DeletedDateTime               types.String `tfsdk:"deleted_date_time"`
	DeviceCategory                types.String `tfsdk:"device_category"`
	DeviceId                      types.String `tfsdk:"device_id"`
	DeviceMetadata                types.String `tfsdk:"device_metadata"`
	DeviceOwnership               types.String `tfsdk:"device_ownership"`
	DeviceVersion                 types.Int64  `tfsdk:"device_version"`
	DisplayName                   types.String `tfsdk:"display_name"`
	EnrollmentProfileName         types.String `tfsdk:"enrollment_profile_name"`
	EnrollmentType                types.String `tfsdk:"enrollment_type"`
	Id                            types.String `tfsdk:"id"`
	IsCompliant                   types.Bool   `tfsdk:"is_compliant"`
	IsManaged                     types.Bool   `tfsdk:"is_managed"`
	IsManagementRestricted        types.Bool   `tfsdk:"is_management_restricted"`
	IsRooted                      types.Bool   `tfsdk:"is_rooted"`
	ManagementType                types.String `tfsdk:"management_type"`
	Manufacturer                  types.String `tfsdk:"manufacturer"`
	MdmAppId                      types.String `tfsdk:"mdm_app_id"`
	Model                         types.String `tfsdk:"model"`
	OnPremisesLastSyncDateTime    types.String `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesSecurityIdentifier  types.String `tfsdk:"on_premises_security_identifier"`
	OnPremisesSyncEnabled         types.Bool   `tfsdk:"on_premises_sync_enabled"`
	OperatingSystem               types.String `tfsdk:"operating_system"`
	OperatingSystemVersion        types.String `tfsdk:"operating_system_version"`
	PhysicalIds                   types.List   `tfsdk:"physical_ids"`
	ProfileType                   types.String `tfsdk:"profile_type"`
	RegistrationDateTime          types.String `tfsdk:"registration_date_time"`
	SystemLabels                  types.List   `tfsdk:"system_labels"`
	TrustType                     types.String `tfsdk:"trust_type"`
}

func (m devicesDeviceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"account_enabled":                    types.BoolType,
		"alternative_security_ids":           types.ListType{ElemType: types.ObjectType{AttrTypes: devicesAlternativeSecurityIdModel{}.AttributeTypes()}},
		"approximate_last_sign_in_date_time": types.StringType,
		"compliance_expiration_date_time":    types.StringType,
		"deleted_date_time":                  types.StringType,
		"device_category":                    types.StringType,
		"device_id":                          types.StringType,
		"device_metadata":                    types.StringType,
		"device_ownership":                   types.StringType,
		"device_version":                     types.Int64Type,
		"display_name":                       types.StringType,
		"enrollment_profile_name":            types.StringType,
		"enrollment_type":                    types.StringType,
		"id":                                 types.StringType,
		"is_compliant":                       types.BoolType,
		"is_managed":                         types.BoolType,
		"is_management_restricted":           types.BoolType,
		"is_rooted":                          types.BoolType,
		"management_type":                    types.StringType,
		"manufacturer":                       types.StringType,
		"mdm_app_id":                         types.StringType,
		"model":                              types.StringType,
		"on_premises_last_sync_date_time":    types.StringType,
		"on_premises_security_identifier":    types.StringType,
		"on_premises_sync_enabled":           types.BoolType,
		"operating_system":                   types.StringType,
		"operating_system_version":           types.StringType,
		"physical_ids":                       types.ListType{ElemType: types.StringType},
		"profile_type":                       types.StringType,
		"registration_date_time":             types.StringType,
		"system_labels":                      types.ListType{ElemType: types.StringType},
		"trust_type":                         types.StringType,
	}
}

type devicesAlternativeSecurityIdModel struct {
	IdentityProvider types.String `tfsdk:"identity_provider"`
	Key              types.String `tfsdk:"key"`
	Type             types.Int64  `tfsdk:"type"`
}

func (m devicesAlternativeSecurityIdModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"identity_provider": types.StringType,
		"key":               types.StringType,
		"type":              types.Int64Type,
	}
}
