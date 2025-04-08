package devices

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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
				Description: "true if the account is enabled; otherwise, false. Required. Default is true.  Supports $filter (eq, ne, not, in). Only callers with at least the Cloud Device Administrator role can set this property.",
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
				Description: "Unique identifier set by Azure Device Registration Service at the time of registration. This alternate key can be used to reference the device object. Supports $filter (eq, ne, not, startsWith).",
				Computed:    true,
			},
			"device_metadata": schema.StringAttribute{
				Description: "For internal use only. Set to null.",
				Computed:    true,
			},
			"device_ownership": schema.StringAttribute{
				Description: "Ownership of the device. Intune sets this property. Possible values are: unknown, company, personal.",
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
			"enrollment_type": schema.StringAttribute{
				Description: "Enrollment type of the device. Intune sets this property. Possible values are: unknown, userEnrollment, deviceEnrollmentManager, appleBulkWithUser, appleBulkWithoutUser, windowsAzureADJoin, windowsBulkUserless, windowsAutoEnrollment, windowsBulkAzureDomainJoin, windowsCoManagement, windowsAzureADJoinUsingDeviceAuth,appleUserEnrollment, appleUserEnrollmentWithServiceAccount. NOTE: This property might return other values apart from those listed.",
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
			"is_management_restricted": schema.BoolAttribute{
				Description: "",
				Computed:    true,
			},
			"is_rooted": schema.BoolAttribute{
				Description: "true if the device is rooted or jail-broken. This property can only be updated by Intune.",
				Computed:    true,
			},
			"management_type": schema.StringAttribute{
				Description: "The management channel of the device. This property is set by Intune. Possible values are: eas, mdm, easMdm, intuneClient, easIntuneClient, configurationManagerClient, configurationManagerClientMdm, configurationManagerClientMdmEas, unknown, jamf, googleCloudDevicePolicyController.",
				Computed:    true,
			},
			"manufacturer": schema.StringAttribute{
				Description: "Manufacturer of the device. Read-only.",
				Computed:    true,
			},
			"mdm_app_id": schema.StringAttribute{
				Description: "Application identifier used to register device into MDM. Read-only. Supports $filter (eq, ne, not, startsWith).",
				Computed:    true,
			},
			"model": schema.StringAttribute{
				Description: "Model of the device. Read-only.",
				Computed:    true,
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "The last time at which the object was synced with the on-premises directory. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z Read-only. Supports $filter (eq, ne, not, ge, le, in).",
				Computed:    true,
			},
			"on_premises_security_identifier": schema.StringAttribute{
				Description: "The on-premises security identifier (SID) for the user who was synchronized from on-premises to the cloud. Read-only. Returned only on $select. Supports $filter (eq).",
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
				Description: "Type of trust for the joined device. Read-only. Possible values:  Workplace (indicates bring your own personal devices), AzureAd (Cloud-only joined devices), ServerAd (on-premises domain joined devices joined to Microsoft Entra ID). For more information, see Introduction to device management in Microsoft Entra ID.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfStateDevice deviceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateDevice)...)
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
				"enrollmentType",
				"isCompliant",
				"isManaged",
				"isManagementRestricted",
				"isRooted",
				"managementType",
				"manufacturer",
				"mdmAppId",
				"model",
				"onPremisesLastSyncDateTime",
				"onPremisesSecurityIdentifier",
				"onPremisesSyncEnabled",
				"operatingSystem",
				"operatingSystemVersion",
				"physicalIds",
				"profileType",
				"registrationDateTime",
				"systemLabels",
				"trustType",
			},
		},
	}

	var response models.Deviceable
	var err error

	if !tfStateDevice.Id.IsNull() {
		response, err = d.client.Devices().ByDeviceId(tfStateDevice.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
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

	if response.GetId() != nil {
		tfStateDevice.Id = types.StringValue(*response.GetId())
	} else {
		tfStateDevice.Id = types.StringNull()
	}
	if response.GetDeletedDateTime() != nil {
		tfStateDevice.DeletedDateTime = types.StringValue(response.GetDeletedDateTime().String())
	} else {
		tfStateDevice.DeletedDateTime = types.StringNull()
	}
	if response.GetAccountEnabled() != nil {
		tfStateDevice.AccountEnabled = types.BoolValue(*response.GetAccountEnabled())
	} else {
		tfStateDevice.AccountEnabled = types.BoolNull()
	}
	if len(response.GetAlternativeSecurityIds()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseAlternativeSecurityIds := range response.GetAlternativeSecurityIds() {
			tfStateAlternativeSecurityId := deviceAlternativeSecurityIdModel{}

			if responseAlternativeSecurityIds.GetIdentityProvider() != nil {
				tfStateAlternativeSecurityId.IdentityProvider = types.StringValue(*responseAlternativeSecurityIds.GetIdentityProvider())
			} else {
				tfStateAlternativeSecurityId.IdentityProvider = types.StringNull()
			}
			if responseAlternativeSecurityIds.GetKey() != nil {
				tfStateAlternativeSecurityId.Key = types.StringValue(string(responseAlternativeSecurityIds.GetKey()[:]))
			} else {
				tfStateAlternativeSecurityId.Key = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAlternativeSecurityId.AttributeTypes(), tfStateAlternativeSecurityId)
			objectValues = append(objectValues, objectValue)
		}
		tfStateDevice.AlternativeSecurityIds, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if response.GetApproximateLastSignInDateTime() != nil {
		tfStateDevice.ApproximateLastSignInDateTime = types.StringValue(response.GetApproximateLastSignInDateTime().String())
	} else {
		tfStateDevice.ApproximateLastSignInDateTime = types.StringNull()
	}
	if response.GetComplianceExpirationDateTime() != nil {
		tfStateDevice.ComplianceExpirationDateTime = types.StringValue(response.GetComplianceExpirationDateTime().String())
	} else {
		tfStateDevice.ComplianceExpirationDateTime = types.StringNull()
	}
	if response.GetDeviceCategory() != nil {
		tfStateDevice.DeviceCategory = types.StringValue(*response.GetDeviceCategory())
	} else {
		tfStateDevice.DeviceCategory = types.StringNull()
	}
	if response.GetDeviceId() != nil {
		tfStateDevice.DeviceId = types.StringValue(*response.GetDeviceId())
	} else {
		tfStateDevice.DeviceId = types.StringNull()
	}
	if response.GetDeviceMetadata() != nil {
		tfStateDevice.DeviceMetadata = types.StringValue(*response.GetDeviceMetadata())
	} else {
		tfStateDevice.DeviceMetadata = types.StringNull()
	}
	if response.GetDeviceOwnership() != nil {
		tfStateDevice.DeviceOwnership = types.StringValue(*response.GetDeviceOwnership())
	} else {
		tfStateDevice.DeviceOwnership = types.StringNull()
	}
	if response.GetDisplayName() != nil {
		tfStateDevice.DisplayName = types.StringValue(*response.GetDisplayName())
	} else {
		tfStateDevice.DisplayName = types.StringNull()
	}
	if response.GetEnrollmentProfileName() != nil {
		tfStateDevice.EnrollmentProfileName = types.StringValue(*response.GetEnrollmentProfileName())
	} else {
		tfStateDevice.EnrollmentProfileName = types.StringNull()
	}
	if response.GetEnrollmentType() != nil {
		tfStateDevice.EnrollmentType = types.StringValue(*response.GetEnrollmentType())
	} else {
		tfStateDevice.EnrollmentType = types.StringNull()
	}
	if response.GetIsCompliant() != nil {
		tfStateDevice.IsCompliant = types.BoolValue(*response.GetIsCompliant())
	} else {
		tfStateDevice.IsCompliant = types.BoolNull()
	}
	if response.GetIsManaged() != nil {
		tfStateDevice.IsManaged = types.BoolValue(*response.GetIsManaged())
	} else {
		tfStateDevice.IsManaged = types.BoolNull()
	}
	if response.GetIsManagementRestricted() != nil {
		tfStateDevice.IsManagementRestricted = types.BoolValue(*response.GetIsManagementRestricted())
	} else {
		tfStateDevice.IsManagementRestricted = types.BoolNull()
	}
	if response.GetIsRooted() != nil {
		tfStateDevice.IsRooted = types.BoolValue(*response.GetIsRooted())
	} else {
		tfStateDevice.IsRooted = types.BoolNull()
	}
	if response.GetManagementType() != nil {
		tfStateDevice.ManagementType = types.StringValue(*response.GetManagementType())
	} else {
		tfStateDevice.ManagementType = types.StringNull()
	}
	if response.GetManufacturer() != nil {
		tfStateDevice.Manufacturer = types.StringValue(*response.GetManufacturer())
	} else {
		tfStateDevice.Manufacturer = types.StringNull()
	}
	if response.GetMdmAppId() != nil {
		tfStateDevice.MdmAppId = types.StringValue(*response.GetMdmAppId())
	} else {
		tfStateDevice.MdmAppId = types.StringNull()
	}
	if response.GetModel() != nil {
		tfStateDevice.Model = types.StringValue(*response.GetModel())
	} else {
		tfStateDevice.Model = types.StringNull()
	}
	if response.GetOnPremisesLastSyncDateTime() != nil {
		tfStateDevice.OnPremisesLastSyncDateTime = types.StringValue(response.GetOnPremisesLastSyncDateTime().String())
	} else {
		tfStateDevice.OnPremisesLastSyncDateTime = types.StringNull()
	}
	if response.GetOnPremisesSecurityIdentifier() != nil {
		tfStateDevice.OnPremisesSecurityIdentifier = types.StringValue(*response.GetOnPremisesSecurityIdentifier())
	} else {
		tfStateDevice.OnPremisesSecurityIdentifier = types.StringNull()
	}
	if response.GetOnPremisesSyncEnabled() != nil {
		tfStateDevice.OnPremisesSyncEnabled = types.BoolValue(*response.GetOnPremisesSyncEnabled())
	} else {
		tfStateDevice.OnPremisesSyncEnabled = types.BoolNull()
	}
	if response.GetOperatingSystem() != nil {
		tfStateDevice.OperatingSystem = types.StringValue(*response.GetOperatingSystem())
	} else {
		tfStateDevice.OperatingSystem = types.StringNull()
	}
	if response.GetOperatingSystemVersion() != nil {
		tfStateDevice.OperatingSystemVersion = types.StringValue(*response.GetOperatingSystemVersion())
	} else {
		tfStateDevice.OperatingSystemVersion = types.StringNull()
	}
	if len(response.GetPhysicalIds()) > 0 {
		var valueArrayPhysicalIds []attr.Value
		for _, responsePhysicalIds := range response.GetPhysicalIds() {
			valueArrayPhysicalIds = append(valueArrayPhysicalIds, types.StringValue(responsePhysicalIds))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayPhysicalIds)
		tfStateDevice.PhysicalIds = listValue
	} else {
		tfStateDevice.PhysicalIds = types.ListNull(types.StringType)
	}
	if response.GetProfileType() != nil {
		tfStateDevice.ProfileType = types.StringValue(*response.GetProfileType())
	} else {
		tfStateDevice.ProfileType = types.StringNull()
	}
	if response.GetRegistrationDateTime() != nil {
		tfStateDevice.RegistrationDateTime = types.StringValue(response.GetRegistrationDateTime().String())
	} else {
		tfStateDevice.RegistrationDateTime = types.StringNull()
	}
	if len(response.GetSystemLabels()) > 0 {
		var valueArraySystemLabels []attr.Value
		for _, responseSystemLabels := range response.GetSystemLabels() {
			valueArraySystemLabels = append(valueArraySystemLabels, types.StringValue(responseSystemLabels))
		}
		listValue, _ := types.ListValue(types.StringType, valueArraySystemLabels)
		tfStateDevice.SystemLabels = listValue
	} else {
		tfStateDevice.SystemLabels = types.ListNull(types.StringType)
	}
	if response.GetTrustType() != nil {
		tfStateDevice.TrustType = types.StringValue(*response.GetTrustType())
	} else {
		tfStateDevice.TrustType = types.StringNull()
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
