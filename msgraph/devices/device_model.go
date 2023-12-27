package devices

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type deviceModel struct {
	Id                            types.String                        `tfsdk:"id"`
	DeletedDateTime               types.String                        `tfsdk:"deleted_date_time"`
	AccountEnabled                types.Bool                          `tfsdk:"account_enabled"`
	AlternativeSecurityIds        []deviceAlternativeSecurityIdsModel `tfsdk:"alternative_security_ids"`
	ApproximateLastSignInDateTime types.String                        `tfsdk:"approximate_last_sign_in_date_time"`
	ComplianceExpirationDateTime  types.String                        `tfsdk:"compliance_expiration_date_time"`
	DeviceCategory                types.String                        `tfsdk:"device_category"`
	DeviceId                      types.String                        `tfsdk:"device_id"`
	DeviceMetadata                types.String                        `tfsdk:"device_metadata"`
	DeviceOwnership               types.String                        `tfsdk:"device_ownership"`
	DeviceVersion                 types.Int64                         `tfsdk:"device_version"`
	DisplayName                   types.String                        `tfsdk:"display_name"`
	EnrollmentProfileName         types.String                        `tfsdk:"enrollment_profile_name"`
	IsCompliant                   types.Bool                          `tfsdk:"is_compliant"`
	IsManaged                     types.Bool                          `tfsdk:"is_managed"`
	MdmAppId                      types.String                        `tfsdk:"mdm_app_id"`
	OnPremisesLastSyncDateTime    types.String                        `tfsdk:"on_premises_last_sync_date_time"`
	OnPremisesSyncEnabled         types.Bool                          `tfsdk:"on_premises_sync_enabled"`
	OperatingSystem               types.String                        `tfsdk:"operating_system"`
	OperatingSystemVersion        types.String                        `tfsdk:"operating_system_version"`
	PhysicalIds                   []types.String                      `tfsdk:"physical_ids"`
	ProfileType                   types.String                        `tfsdk:"profile_type"`
	RegistrationDateTime          types.String                        `tfsdk:"registration_date_time"`
	SystemLabels                  []types.String                      `tfsdk:"system_labels"`
	TrustType                     types.String                        `tfsdk:"trust_type"`
}

type deviceAlternativeSecurityIdsModel struct {
	IdentityProvider types.String `tfsdk:"identity_provider"`
	Key              types.String `tfsdk:"key"`
	Type             types.Int64  `tfsdk:"type"`
}
