package devices

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/devices"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deviceDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceDataSource{}
)

// NewDeviceDataSource is a helper function to simplify the provider implementation.
func NewDeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

// deviceDataSource is the data source implementation.
type deviceDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Configure adds the provider configured client to the data source.
func (d *deviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
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
			"account_enabled": schema.BoolAttribute{
				Description: "true if the account is enabled; otherwise, false. Required. Default is true.  Supports $filter (eq, ne, not, in). Only callers in Global Administrator and Cloud Device Administrator roles can set this property.",
				Computed:    true,
			},
			"alternative_security_ids": schema.ListNestedAttribute{
				Description: "For internal use only. Not nullable. Supports $filter (eq, not, ge, le).",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identity_provider": schema.StringAttribute{
							Description: "For internal use only.",
							Computed:    true,
						},
						"key": schema.StringAttribute{
							Description: "For internal use only.",
							Computed:    true,
						},
						"type": schema.Int64Attribute{
							Description: "For internal use only.",
							Computed:    true,
						},
					},
				},
			},
			"approximate_last_sign_in_date_time": schema.StringAttribute{
				Description: "The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports $filter (eq, ne, not, ge, le, and eq on null values) and $orderby.",
				Computed:    true,
			},
			"compliance_expiration_date_time": schema.StringAttribute{
				Description: "The timestamp when the device is no longer deemed compliant. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
				Computed:    true,
			},
			"device_category": schema.StringAttribute{
				Description: "User-defined property set by Intune to automatically add devices to groups and simplify managing devices.",
				Computed:    true,
			},
			"device_id": schema.StringAttribute{
				Description: "Unique identifier set by Azure Device Registration Service at the time of registration. This is an alternate key that can be used to reference the device object. Supports $filter (eq, ne, not, startsWith).",
				Computed:    true,
			},
			"device_metadata": schema.StringAttribute{
				Description: "For internal use only. Set to null.",
				Computed:    true,
			},
			"device_ownership": schema.StringAttribute{
				Description: "Ownership of the device. This property is set by Intune. Possible values are: unknown, company, personal.",
				Computed:    true,
			},
			"device_version": schema.Int64Attribute{
				Description: "For internal use only.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The display name for the device. Required. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Computed:    true,
			},
			"enrollment_profile_name": schema.StringAttribute{
				Description: "Enrollment profile applied to the device. For example, Apple Device Enrollment Profile, Device enrollment - Corporate device identifiers, or Windows Autopilot profile name. This property is set by Intune.",
				Computed:    true,
			},
			"is_compliant": schema.BoolAttribute{
				Description: "true if the device complies with Mobile Device Management (MDM) policies; otherwise, false. Read-only. This can only be updated by Intune for any device OS type or by an approved MDM app for Windows OS devices. Supports $filter (eq, ne, not).",
				Computed:    true,
			},
			"is_managed": schema.BoolAttribute{
				Description: "true if the device is managed by a Mobile Device Management (MDM) app; otherwise, false. This can only be updated by Intune for any device OS type or by an approved MDM app for Windows OS devices. Supports $filter (eq, ne, not).",
				Computed:    true,
			},
			"mdm_app_id": schema.StringAttribute{
				Description: "Application identifier used to register device into MDM. Read-only. Supports $filter (eq, ne, not, startsWith).",
				Computed:    true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "The last time at which the object was synced with the on-premises directory. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z Read-only. Supports $filter (eq, ne, not, ge, le, in).",
				Computed:    true,
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Description: "true if this object is synced from an on-premises directory; false if this object was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default). Read-only. Supports $filter (eq, ne, not, in, and eq on null values).",
				Computed:    true,
			},
			"operating_system": schema.StringAttribute{
				Description: "The type of operating system on the device. Required. Supports $filter (eq, ne, not, ge, le, startsWith, and eq on null values).",
				Computed:    true,
			},
			"operating_system_version": schema.StringAttribute{
				Description: "The version of the operating system on the device. Required. Supports $filter (eq, ne, not, ge, le, startsWith, and eq on null values).",
				Computed:    true,
			},
			"physical_ids": schema.ListAttribute{
				Description: "For internal use only. Not nullable. Supports $filter (eq, not, ge, le, startsWith,/$count eq 0, /$count ne 0).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"profile_type": schema.StringAttribute{
				Description: "The profile type of the device. Possible values: RegisteredDevice (default), SecureVM, Printer, Shared, IoT.",
				Computed:    true,
			},
			"registration_date_time": schema.StringAttribute{
				Description: "Date and time of when the device was registered. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
				Computed:    true,
			},
			"system_labels": schema.ListAttribute{
				Description: "List of labels applied to the device by the system. Supports $filter (/$count eq 0, /$count ne 0).",
				Computed:    true,
				ElementType: types.StringType,
			},
			"trust_type": schema.StringAttribute{
				Description: "Type of trust for the joined device. Read-only. Possible values:  Workplace (indicates bring your own personal devices), AzureAd (Cloud only joined devices), ServerAd (on-premises domain joined devices joined to Microsoft Entra ID). For more details, see Introduction to device management in Microsoft Entra ID.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := devices.DeviceItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &devices.DeviceItemRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"deletedDateTime",
				"accountEnabled",
				"alternativeSecurityIds",
				"approximateLastSignInDateTime",
				"complianceExpirationDateTime",
				"deviceCategory",
				"deviceId",
				"deviceMetadata",
				"deviceOwnership",
				"deviceVersion",
				"displayName",
				"enrollmentProfileName",
				"isCompliant",
				"isManaged",
				"mdmAppId",
				"onPremisesLastSyncDateTime",
				"onPremisesSyncEnabled",
				"operatingSystem",
				"operatingSystemVersion",
				"physicalIds",
				"profileType",
				"registrationDateTime",
				"systemLabels",
				"trustType",
				"memberOf",
				"registeredOwners",
				"registeredUsers",
				"transitiveMemberOf",
				"extensions",
			},
		},
	}

	var result models.Deviceable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.Devices().ByDeviceId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"`id` must be supplied.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting device",
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
	if result.GetAccountEnabled() != nil {
		state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	}
	for _, v := range result.GetAlternativeSecurityIds() {
		alternativeSecurityIds := new(deviceAlternativeSecurityIdsDataSourceModel)

		if v.GetIdentityProvider() != nil {
			alternativeSecurityIds.IdentityProvider = types.StringValue(*v.GetIdentityProvider())
		}
		if v.GetKey() != nil {
			alternativeSecurityIds.Key = types.StringValue(string(v.GetKey()[:]))
		}
		if v.GetTypeEscaped() != nil {
			alternativeSecurityIds.Type = types.Int64Value(int64(*v.GetTypeEscaped()))
		}
		state.AlternativeSecurityIds = append(state.AlternativeSecurityIds, *alternativeSecurityIds)
	}
	if result.GetApproximateLastSignInDateTime() != nil {
		state.ApproximateLastSignInDateTime = types.StringValue(result.GetApproximateLastSignInDateTime().String())
	}
	if result.GetComplianceExpirationDateTime() != nil {
		state.ComplianceExpirationDateTime = types.StringValue(result.GetComplianceExpirationDateTime().String())
	}
	if result.GetDeviceCategory() != nil {
		state.DeviceCategory = types.StringValue(*result.GetDeviceCategory())
	}
	if result.GetDeviceId() != nil {
		state.DeviceId = types.StringValue(*result.GetDeviceId())
	}
	if result.GetDeviceMetadata() != nil {
		state.DeviceMetadata = types.StringValue(*result.GetDeviceMetadata())
	}
	if result.GetDeviceOwnership() != nil {
		state.DeviceOwnership = types.StringValue(*result.GetDeviceOwnership())
	}
	if result.GetDeviceVersion() != nil {
		state.DeviceVersion = types.Int64Value(int64(*result.GetDeviceVersion()))
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	}
	if result.GetEnrollmentProfileName() != nil {
		state.EnrollmentProfileName = types.StringValue(*result.GetEnrollmentProfileName())
	}
	if result.GetIsCompliant() != nil {
		state.IsCompliant = types.BoolValue(*result.GetIsCompliant())
	}
	if result.GetIsManaged() != nil {
		state.IsManaged = types.BoolValue(*result.GetIsManaged())
	}
	if result.GetMdmAppId() != nil {
		state.MdmAppId = types.StringValue(*result.GetMdmAppId())
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	}
	if result.GetOnPremisesSyncEnabled() != nil {
		state.OnPremisesSyncEnabled = types.BoolValue(*result.GetOnPremisesSyncEnabled())
	}
	if result.GetOperatingSystem() != nil {
		state.OperatingSystem = types.StringValue(*result.GetOperatingSystem())
	}
	if result.GetOperatingSystemVersion() != nil {
		state.OperatingSystemVersion = types.StringValue(*result.GetOperatingSystemVersion())
	}
	for _, v := range result.GetPhysicalIds() {
		state.PhysicalIds = append(state.PhysicalIds, types.StringValue(v))
	}
	if result.GetProfileType() != nil {
		state.ProfileType = types.StringValue(*result.GetProfileType())
	}
	if result.GetRegistrationDateTime() != nil {
		state.RegistrationDateTime = types.StringValue(result.GetRegistrationDateTime().String())
	}
	for _, v := range result.GetSystemLabels() {
		state.SystemLabels = append(state.SystemLabels, types.StringValue(v))
	}
	if result.GetTrustType() != nil {
		state.TrustType = types.StringValue(*result.GetTrustType())
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
