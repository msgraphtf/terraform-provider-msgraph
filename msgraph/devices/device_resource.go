package devices

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/devices"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"terraform-provider-msgraph/planmodifiers/boolplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/listplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/stringplanmodifiers"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &deviceResource{}
	_ resource.ResourceWithConfigure = &deviceResource{}
)

// NewDeviceResource is a helper function to simplify the provider implementation.
func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

// deviceResource is the resource implementation.
type deviceResource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (d *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

// Configure adds the provider configured client to the resource.
func (d *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the resource.
func (d *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"deleted_date_time": schema.StringAttribute{
				Description: "Date and time when this object was deleted. Always null when the object hasn't been deleted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"account_enabled": schema.BoolAttribute{
				Description: "true if the account is enabled; otherwise, false. Required. Default is true.  Supports $filter (eq, ne, not, in). Only callers with at least the Cloud Device Administrator role can set this property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"alternative_security_ids": schema.ListNestedAttribute{
				Description: "For internal use only. Not nullable. Supports $filter (eq, not, ge, le).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"identity_provider": schema.StringAttribute{
							Description: "For internal use only.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
						"key": schema.StringAttribute{
							Description: "For internal use only.",
							Optional:    true,
							Computed:    true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifiers.UseStateForUnconfigured(),
							},
						},
					},
				},
			},
			"approximate_last_sign_in_date_time": schema.StringAttribute{
				Description: "The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only. Supports $filter (eq, ne, not, ge, le, and eq on null values) and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"compliance_expiration_date_time": schema.StringAttribute{
				Description: "The timestamp when the device is no longer deemed compliant. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"device_category": schema.StringAttribute{
				Description: "User-defined property set by Intune to automatically add devices to groups and simplify managing devices.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"device_id": schema.StringAttribute{
				Description: "Unique identifier set by Azure Device Registration Service at the time of registration. This alternate key can be used to reference the device object. Supports $filter (eq, ne, not, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"device_metadata": schema.StringAttribute{
				Description: "For internal use only. Set to null.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"device_ownership": schema.StringAttribute{
				Description: "Ownership of the device. Intune sets this property. Possible values are: unknown, company, personal.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The display name for the device. Required. Supports $filter (eq, ne, not, ge, le, in, startsWith, and eq on null values), $search, and $orderby.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"enrollment_profile_name": schema.StringAttribute{
				Description: "Enrollment profile applied to the device. For example, Apple Device Enrollment Profile, Device enrollment - Corporate device identifiers, or Windows Autopilot profile name. This property is set by Intune.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"enrollment_type": schema.StringAttribute{
				Description: "Enrollment type of the device. Intune sets this property. Possible values are: unknown, userEnrollment, deviceEnrollmentManager, appleBulkWithUser, appleBulkWithoutUser, windowsAzureADJoin, windowsBulkUserless, windowsAutoEnrollment, windowsBulkAzureDomainJoin, windowsCoManagement, windowsAzureADJoinUsingDeviceAuth,appleUserEnrollment, appleUserEnrollmentWithServiceAccount. NOTE: This property might return other values apart from those listed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_compliant": schema.BoolAttribute{
				Description: "true if the device complies with Mobile Device Management (MDM) policies; otherwise, false. Read-only. This can only be updated by Intune for any device OS type or by an approved MDM app for Windows OS devices. Supports $filter (eq, ne, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_managed": schema.BoolAttribute{
				Description: "true if the device is managed by a Mobile Device Management (MDM) app; otherwise, false. This can only be updated by Intune for any device OS type or by an approved MDM app for Windows OS devices. Supports $filter (eq, ne, not).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_management_restricted": schema.BoolAttribute{
				Description: "",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_rooted": schema.BoolAttribute{
				Description: "true if the device is rooted or jail-broken. This property can only be updated by Intune.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"management_type": schema.StringAttribute{
				Description: "The management channel of the device. This property is set by Intune. Possible values are: eas, mdm, easMdm, intuneClient, easIntuneClient, configurationManagerClient, configurationManagerClientMdm, configurationManagerClientMdmEas, unknown, jamf, googleCloudDevicePolicyController.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"manufacturer": schema.StringAttribute{
				Description: "Manufacturer of the device. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"mdm_app_id": schema.StringAttribute{
				Description: "Application identifier used to register device into MDM. Read-only. Supports $filter (eq, ne, not, startsWith).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"model": schema.StringAttribute{
				Description: "Model of the device. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_last_sync_date_time": schema.StringAttribute{
				Description: "The last time at which the object was synced with the on-premises directory. The Timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z Read-only. Supports $filter (eq, ne, not, ge, le, in).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_security_identifier": schema.StringAttribute{
				Description: "The on-premises security identifier (SID) for the user who was synchronized from on-premises to the cloud. Read-only. Returned only on $select. Supports $filter (eq).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"on_premises_sync_enabled": schema.BoolAttribute{
				Description: "true if this object is synced from an on-premises directory; false if this object was originally synced from an on-premises directory but is no longer synced; null if this object has never been synced from an on-premises directory (default). Read-only. Supports $filter (eq, ne, not, in, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"operating_system": schema.StringAttribute{
				Description: "The type of operating system on the device. Required. Supports $filter (eq, ne, not, ge, le, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"operating_system_version": schema.StringAttribute{
				Description: "The version of the operating system on the device. Required. Supports $filter (eq, ne, not, ge, le, startsWith, and eq on null values).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"physical_ids": schema.ListAttribute{
				Description: "For internal use only. Not nullable. Supports $filter (eq, not, ge, le, startsWith,/$count eq 0, /$count ne 0).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"profile_type": schema.StringAttribute{
				Description: "The profile type of the device. Possible values: RegisteredDevice (default), SecureVM, Printer, Shared, IoT.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"registration_date_time": schema.StringAttribute{
				Description: "Date and time of when the device was registered. The timestamp type represents date and time information using ISO 8601 format and is always in UTC time. For example, midnight UTC on Jan 1, 2014 is 2014-01-01T00:00:00Z. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"system_labels": schema.ListAttribute{
				Description: "List of labels applied to the device by the system. Supports $filter (/$count eq 0, /$count ne 0).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.List{
					listplanmodifiers.UseStateForUnconfigured(),
				},
				ElementType: types.StringType,
			},
			"trust_type": schema.StringAttribute{
				Description: "Type of trust for the joined device. Read-only. Possible values:  Workplace (indicates bring your own personal devices), AzureAd (Cloud-only joined devices), ServerAd (on-premises domain joined devices joined to Microsoft Entra ID). For more information, see Introduction to device management in Microsoft Entra ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanDevice deviceModel
	diags := req.Plan.Get(ctx, &tfPlanDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	requestBodyDevice := models.NewDevice()
	// START Id | CreateStringAttribute
	if !tfPlanDevice.Id.IsUnknown() {
		tfPlanId := tfPlanDevice.Id.ValueString()
		requestBodyDevice.SetId(&tfPlanId)
	} else {
		tfPlanDevice.Id = types.StringNull()
	}
	// END Id | CreateStringAttribute

	// START DeletedDateTime | CreateStringTimeAttribute
	if !tfPlanDevice.DeletedDateTime.IsUnknown() {
		tfPlanDeletedDateTime := tfPlanDevice.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyDevice.SetDeletedDateTime(&t)
	} else {
		tfPlanDevice.DeletedDateTime = types.StringNull()
	}
	// END DeletedDateTime | CreateStringTimeAttribute

	// START AccountEnabled | CreateBoolAttribute
	if !tfPlanDevice.AccountEnabled.IsUnknown() {
		tfPlanAccountEnabled := tfPlanDevice.AccountEnabled.ValueBool()
		requestBodyDevice.SetAccountEnabled(&tfPlanAccountEnabled)
	} else {
		tfPlanDevice.AccountEnabled = types.BoolNull()
	}
	// END AccountEnabled | CreateBoolAttribute

	// START AlternativeSecurityIds | CreateArrayObjectAttribute
	if len(tfPlanDevice.AlternativeSecurityIds.Elements()) > 0 {
		var requestBodyAlternativeSecurityIds []models.AlternativeSecurityIdable
		for _, i := range tfPlanDevice.AlternativeSecurityIds.Elements() {
			requestBodyAlternativeSecurityId := models.NewAlternativeSecurityId()
			tfPlanAlternativeSecurityId := deviceAlternativeSecurityIdModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAlternativeSecurityId)

			// START IdentityProvider | CreateStringAttribute
			if !tfPlanAlternativeSecurityId.IdentityProvider.IsUnknown() {
				tfPlanIdentityProvider := tfPlanAlternativeSecurityId.IdentityProvider.ValueString()
				requestBodyAlternativeSecurityId.SetIdentityProvider(&tfPlanIdentityProvider)
			} else {
				tfPlanAlternativeSecurityId.IdentityProvider = types.StringNull()
			}
			// END IdentityProvider | CreateStringAttribute

			// START Key | CreateStringBase64UrlAttribute
			if !tfPlanAlternativeSecurityId.Key.IsUnknown() {
				tfPlanKey := tfPlanAlternativeSecurityId.Key.ValueString()
				requestBodyAlternativeSecurityId.SetKey([]byte(tfPlanKey))
			} else {
				tfPlanAlternativeSecurityId.Key = types.StringNull()
			}
			// END Key | CreateStringBase64UrlAttribute

			// START Type | UNKNOWN
			// END Type | UNKNOWN

		}
		requestBodyDevice.SetAlternativeSecurityIds(requestBodyAlternativeSecurityIds)
	} else {
		tfPlanDevice.AlternativeSecurityIds = types.ListNull(tfPlanDevice.AlternativeSecurityIds.ElementType(ctx))
	}
	// END AlternativeSecurityIds | CreateArrayObjectAttribute

	// START ApproximateLastSignInDateTime | CreateStringTimeAttribute
	if !tfPlanDevice.ApproximateLastSignInDateTime.IsUnknown() {
		tfPlanApproximateLastSignInDateTime := tfPlanDevice.ApproximateLastSignInDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanApproximateLastSignInDateTime)
		requestBodyDevice.SetApproximateLastSignInDateTime(&t)
	} else {
		tfPlanDevice.ApproximateLastSignInDateTime = types.StringNull()
	}
	// END ApproximateLastSignInDateTime | CreateStringTimeAttribute

	// START ComplianceExpirationDateTime | CreateStringTimeAttribute
	if !tfPlanDevice.ComplianceExpirationDateTime.IsUnknown() {
		tfPlanComplianceExpirationDateTime := tfPlanDevice.ComplianceExpirationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanComplianceExpirationDateTime)
		requestBodyDevice.SetComplianceExpirationDateTime(&t)
	} else {
		tfPlanDevice.ComplianceExpirationDateTime = types.StringNull()
	}
	// END ComplianceExpirationDateTime | CreateStringTimeAttribute

	// START DeviceCategory | CreateStringAttribute
	if !tfPlanDevice.DeviceCategory.IsUnknown() {
		tfPlanDeviceCategory := tfPlanDevice.DeviceCategory.ValueString()
		requestBodyDevice.SetDeviceCategory(&tfPlanDeviceCategory)
	} else {
		tfPlanDevice.DeviceCategory = types.StringNull()
	}
	// END DeviceCategory | CreateStringAttribute

	// START DeviceId | CreateStringAttribute
	if !tfPlanDevice.DeviceId.IsUnknown() {
		tfPlanDeviceId := tfPlanDevice.DeviceId.ValueString()
		requestBodyDevice.SetDeviceId(&tfPlanDeviceId)
	} else {
		tfPlanDevice.DeviceId = types.StringNull()
	}
	// END DeviceId | CreateStringAttribute

	// START DeviceMetadata | CreateStringAttribute
	if !tfPlanDevice.DeviceMetadata.IsUnknown() {
		tfPlanDeviceMetadata := tfPlanDevice.DeviceMetadata.ValueString()
		requestBodyDevice.SetDeviceMetadata(&tfPlanDeviceMetadata)
	} else {
		tfPlanDevice.DeviceMetadata = types.StringNull()
	}
	// END DeviceMetadata | CreateStringAttribute

	// START DeviceOwnership | CreateStringAttribute
	if !tfPlanDevice.DeviceOwnership.IsUnknown() {
		tfPlanDeviceOwnership := tfPlanDevice.DeviceOwnership.ValueString()
		requestBodyDevice.SetDeviceOwnership(&tfPlanDeviceOwnership)
	} else {
		tfPlanDevice.DeviceOwnership = types.StringNull()
	}
	// END DeviceOwnership | CreateStringAttribute

	// START DeviceVersion | UNKNOWN
	// END DeviceVersion | UNKNOWN

	// START DisplayName | CreateStringAttribute
	if !tfPlanDevice.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanDevice.DisplayName.ValueString()
		requestBodyDevice.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanDevice.DisplayName = types.StringNull()
	}
	// END DisplayName | CreateStringAttribute

	// START EnrollmentProfileName | CreateStringAttribute
	if !tfPlanDevice.EnrollmentProfileName.IsUnknown() {
		tfPlanEnrollmentProfileName := tfPlanDevice.EnrollmentProfileName.ValueString()
		requestBodyDevice.SetEnrollmentProfileName(&tfPlanEnrollmentProfileName)
	} else {
		tfPlanDevice.EnrollmentProfileName = types.StringNull()
	}
	// END EnrollmentProfileName | CreateStringAttribute

	// START EnrollmentType | CreateStringAttribute
	if !tfPlanDevice.EnrollmentType.IsUnknown() {
		tfPlanEnrollmentType := tfPlanDevice.EnrollmentType.ValueString()
		requestBodyDevice.SetEnrollmentType(&tfPlanEnrollmentType)
	} else {
		tfPlanDevice.EnrollmentType = types.StringNull()
	}
	// END EnrollmentType | CreateStringAttribute

	// START IsCompliant | CreateBoolAttribute
	if !tfPlanDevice.IsCompliant.IsUnknown() {
		tfPlanIsCompliant := tfPlanDevice.IsCompliant.ValueBool()
		requestBodyDevice.SetIsCompliant(&tfPlanIsCompliant)
	} else {
		tfPlanDevice.IsCompliant = types.BoolNull()
	}
	// END IsCompliant | CreateBoolAttribute

	// START IsManaged | CreateBoolAttribute
	if !tfPlanDevice.IsManaged.IsUnknown() {
		tfPlanIsManaged := tfPlanDevice.IsManaged.ValueBool()
		requestBodyDevice.SetIsManaged(&tfPlanIsManaged)
	} else {
		tfPlanDevice.IsManaged = types.BoolNull()
	}
	// END IsManaged | CreateBoolAttribute

	// START IsManagementRestricted | CreateBoolAttribute
	if !tfPlanDevice.IsManagementRestricted.IsUnknown() {
		tfPlanIsManagementRestricted := tfPlanDevice.IsManagementRestricted.ValueBool()
		requestBodyDevice.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	} else {
		tfPlanDevice.IsManagementRestricted = types.BoolNull()
	}
	// END IsManagementRestricted | CreateBoolAttribute

	// START IsRooted | CreateBoolAttribute
	if !tfPlanDevice.IsRooted.IsUnknown() {
		tfPlanIsRooted := tfPlanDevice.IsRooted.ValueBool()
		requestBodyDevice.SetIsRooted(&tfPlanIsRooted)
	} else {
		tfPlanDevice.IsRooted = types.BoolNull()
	}
	// END IsRooted | CreateBoolAttribute

	// START ManagementType | CreateStringAttribute
	if !tfPlanDevice.ManagementType.IsUnknown() {
		tfPlanManagementType := tfPlanDevice.ManagementType.ValueString()
		requestBodyDevice.SetManagementType(&tfPlanManagementType)
	} else {
		tfPlanDevice.ManagementType = types.StringNull()
	}
	// END ManagementType | CreateStringAttribute

	// START Manufacturer | CreateStringAttribute
	if !tfPlanDevice.Manufacturer.IsUnknown() {
		tfPlanManufacturer := tfPlanDevice.Manufacturer.ValueString()
		requestBodyDevice.SetManufacturer(&tfPlanManufacturer)
	} else {
		tfPlanDevice.Manufacturer = types.StringNull()
	}
	// END Manufacturer | CreateStringAttribute

	// START MdmAppId | CreateStringAttribute
	if !tfPlanDevice.MdmAppId.IsUnknown() {
		tfPlanMdmAppId := tfPlanDevice.MdmAppId.ValueString()
		requestBodyDevice.SetMdmAppId(&tfPlanMdmAppId)
	} else {
		tfPlanDevice.MdmAppId = types.StringNull()
	}
	// END MdmAppId | CreateStringAttribute

	// START Model | CreateStringAttribute
	if !tfPlanDevice.Model.IsUnknown() {
		tfPlanModel := tfPlanDevice.Model.ValueString()
		requestBodyDevice.SetModel(&tfPlanModel)
	} else {
		tfPlanDevice.Model = types.StringNull()
	}
	// END Model | CreateStringAttribute

	// START OnPremisesLastSyncDateTime | CreateStringTimeAttribute
	if !tfPlanDevice.OnPremisesLastSyncDateTime.IsUnknown() {
		tfPlanOnPremisesLastSyncDateTime := tfPlanDevice.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		requestBodyDevice.SetOnPremisesLastSyncDateTime(&t)
	} else {
		tfPlanDevice.OnPremisesLastSyncDateTime = types.StringNull()
	}
	// END OnPremisesLastSyncDateTime | CreateStringTimeAttribute

	// START OnPremisesSecurityIdentifier | CreateStringAttribute
	if !tfPlanDevice.OnPremisesSecurityIdentifier.IsUnknown() {
		tfPlanOnPremisesSecurityIdentifier := tfPlanDevice.OnPremisesSecurityIdentifier.ValueString()
		requestBodyDevice.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	} else {
		tfPlanDevice.OnPremisesSecurityIdentifier = types.StringNull()
	}
	// END OnPremisesSecurityIdentifier | CreateStringAttribute

	// START OnPremisesSyncEnabled | CreateBoolAttribute
	if !tfPlanDevice.OnPremisesSyncEnabled.IsUnknown() {
		tfPlanOnPremisesSyncEnabled := tfPlanDevice.OnPremisesSyncEnabled.ValueBool()
		requestBodyDevice.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	} else {
		tfPlanDevice.OnPremisesSyncEnabled = types.BoolNull()
	}
	// END OnPremisesSyncEnabled | CreateBoolAttribute

	// START OperatingSystem | CreateStringAttribute
	if !tfPlanDevice.OperatingSystem.IsUnknown() {
		tfPlanOperatingSystem := tfPlanDevice.OperatingSystem.ValueString()
		requestBodyDevice.SetOperatingSystem(&tfPlanOperatingSystem)
	} else {
		tfPlanDevice.OperatingSystem = types.StringNull()
	}
	// END OperatingSystem | CreateStringAttribute

	// START OperatingSystemVersion | CreateStringAttribute
	if !tfPlanDevice.OperatingSystemVersion.IsUnknown() {
		tfPlanOperatingSystemVersion := tfPlanDevice.OperatingSystemVersion.ValueString()
		requestBodyDevice.SetOperatingSystemVersion(&tfPlanOperatingSystemVersion)
	} else {
		tfPlanDevice.OperatingSystemVersion = types.StringNull()
	}
	// END OperatingSystemVersion | CreateStringAttribute

	// START PhysicalIds | CreateArrayStringAttribute
	if len(tfPlanDevice.PhysicalIds.Elements()) > 0 {
		var stringArrayPhysicalIds []string
		for _, i := range tfPlanDevice.PhysicalIds.Elements() {
			stringArrayPhysicalIds = append(stringArrayPhysicalIds, i.String())
		}
		requestBodyDevice.SetPhysicalIds(stringArrayPhysicalIds)
	} else {
		tfPlanDevice.PhysicalIds = types.ListNull(types.StringType)
	}
	// END PhysicalIds | CreateArrayStringAttribute

	// START ProfileType | CreateStringAttribute
	if !tfPlanDevice.ProfileType.IsUnknown() {
		tfPlanProfileType := tfPlanDevice.ProfileType.ValueString()
		requestBodyDevice.SetProfileType(&tfPlanProfileType)
	} else {
		tfPlanDevice.ProfileType = types.StringNull()
	}
	// END ProfileType | CreateStringAttribute

	// START RegistrationDateTime | CreateStringTimeAttribute
	if !tfPlanDevice.RegistrationDateTime.IsUnknown() {
		tfPlanRegistrationDateTime := tfPlanDevice.RegistrationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanRegistrationDateTime)
		requestBodyDevice.SetRegistrationDateTime(&t)
	} else {
		tfPlanDevice.RegistrationDateTime = types.StringNull()
	}
	// END RegistrationDateTime | CreateStringTimeAttribute

	// START SystemLabels | CreateArrayStringAttribute
	if len(tfPlanDevice.SystemLabels.Elements()) > 0 {
		var stringArraySystemLabels []string
		for _, i := range tfPlanDevice.SystemLabels.Elements() {
			stringArraySystemLabels = append(stringArraySystemLabels, i.String())
		}
		requestBodyDevice.SetSystemLabels(stringArraySystemLabels)
	} else {
		tfPlanDevice.SystemLabels = types.ListNull(types.StringType)
	}
	// END SystemLabels | CreateArrayStringAttribute

	// START TrustType | CreateStringAttribute
	if !tfPlanDevice.TrustType.IsUnknown() {
		tfPlanTrustType := tfPlanDevice.TrustType.ValueString()
		requestBodyDevice.SetTrustType(&tfPlanTrustType)
	} else {
		tfPlanDevice.TrustType = types.StringNull()
	}
	// END TrustType | CreateStringAttribute

	// Create new device
	result, err := r.client.Devices().Post(context.Background(), requestBodyDevice, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating device",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlanDevice.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlanDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tfStateDevice deviceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfStateDevice)...)
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

	var result models.Deviceable
	var err error

	if !tfStateDevice.Id.IsNull() {
		result, err = d.client.Devices().ByDeviceId(tfStateDevice.Id.ValueString()).Get(context.Background(), &qparams)
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

	if result.GetId() != nil {
		tfStateDevice.Id = types.StringValue(*result.GetId())
	} else {
		tfStateDevice.Id = types.StringNull()
	}
	if result.GetDeletedDateTime() != nil {
		tfStateDevice.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	} else {
		tfStateDevice.DeletedDateTime = types.StringNull()
	}
	if result.GetAccountEnabled() != nil {
		tfStateDevice.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	} else {
		tfStateDevice.AccountEnabled = types.BoolNull()
	}
	if len(result.GetAlternativeSecurityIds()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAlternativeSecurityIds() {
			tfStateAlternativeSecurityIds := deviceAlternativeSecurityIdModel{}

			if v.GetIdentityProvider() != nil {
				tfStateAlternativeSecurityIds.IdentityProvider = types.StringValue(*v.GetIdentityProvider())
			} else {
				tfStateAlternativeSecurityIds.IdentityProvider = types.StringNull()
			}
			if v.GetKey() != nil {
				tfStateAlternativeSecurityIds.Key = types.StringValue(string(v.GetKey()[:]))
			} else {
				tfStateAlternativeSecurityIds.Key = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateAlternativeSecurityIds.AttributeTypes(), tfStateAlternativeSecurityIds)
			objectValues = append(objectValues, objectValue)
		}
		tfStateDevice.AlternativeSecurityIds, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetApproximateLastSignInDateTime() != nil {
		tfStateDevice.ApproximateLastSignInDateTime = types.StringValue(result.GetApproximateLastSignInDateTime().String())
	} else {
		tfStateDevice.ApproximateLastSignInDateTime = types.StringNull()
	}
	if result.GetComplianceExpirationDateTime() != nil {
		tfStateDevice.ComplianceExpirationDateTime = types.StringValue(result.GetComplianceExpirationDateTime().String())
	} else {
		tfStateDevice.ComplianceExpirationDateTime = types.StringNull()
	}
	if result.GetDeviceCategory() != nil {
		tfStateDevice.DeviceCategory = types.StringValue(*result.GetDeviceCategory())
	} else {
		tfStateDevice.DeviceCategory = types.StringNull()
	}
	if result.GetDeviceId() != nil {
		tfStateDevice.DeviceId = types.StringValue(*result.GetDeviceId())
	} else {
		tfStateDevice.DeviceId = types.StringNull()
	}
	if result.GetDeviceMetadata() != nil {
		tfStateDevice.DeviceMetadata = types.StringValue(*result.GetDeviceMetadata())
	} else {
		tfStateDevice.DeviceMetadata = types.StringNull()
	}
	if result.GetDeviceOwnership() != nil {
		tfStateDevice.DeviceOwnership = types.StringValue(*result.GetDeviceOwnership())
	} else {
		tfStateDevice.DeviceOwnership = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		tfStateDevice.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		tfStateDevice.DisplayName = types.StringNull()
	}
	if result.GetEnrollmentProfileName() != nil {
		tfStateDevice.EnrollmentProfileName = types.StringValue(*result.GetEnrollmentProfileName())
	} else {
		tfStateDevice.EnrollmentProfileName = types.StringNull()
	}
	if result.GetEnrollmentType() != nil {
		tfStateDevice.EnrollmentType = types.StringValue(*result.GetEnrollmentType())
	} else {
		tfStateDevice.EnrollmentType = types.StringNull()
	}
	if result.GetIsCompliant() != nil {
		tfStateDevice.IsCompliant = types.BoolValue(*result.GetIsCompliant())
	} else {
		tfStateDevice.IsCompliant = types.BoolNull()
	}
	if result.GetIsManaged() != nil {
		tfStateDevice.IsManaged = types.BoolValue(*result.GetIsManaged())
	} else {
		tfStateDevice.IsManaged = types.BoolNull()
	}
	if result.GetIsManagementRestricted() != nil {
		tfStateDevice.IsManagementRestricted = types.BoolValue(*result.GetIsManagementRestricted())
	} else {
		tfStateDevice.IsManagementRestricted = types.BoolNull()
	}
	if result.GetIsRooted() != nil {
		tfStateDevice.IsRooted = types.BoolValue(*result.GetIsRooted())
	} else {
		tfStateDevice.IsRooted = types.BoolNull()
	}
	if result.GetManagementType() != nil {
		tfStateDevice.ManagementType = types.StringValue(*result.GetManagementType())
	} else {
		tfStateDevice.ManagementType = types.StringNull()
	}
	if result.GetManufacturer() != nil {
		tfStateDevice.Manufacturer = types.StringValue(*result.GetManufacturer())
	} else {
		tfStateDevice.Manufacturer = types.StringNull()
	}
	if result.GetMdmAppId() != nil {
		tfStateDevice.MdmAppId = types.StringValue(*result.GetMdmAppId())
	} else {
		tfStateDevice.MdmAppId = types.StringNull()
	}
	if result.GetModel() != nil {
		tfStateDevice.Model = types.StringValue(*result.GetModel())
	} else {
		tfStateDevice.Model = types.StringNull()
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		tfStateDevice.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	} else {
		tfStateDevice.OnPremisesLastSyncDateTime = types.StringNull()
	}
	if result.GetOnPremisesSecurityIdentifier() != nil {
		tfStateDevice.OnPremisesSecurityIdentifier = types.StringValue(*result.GetOnPremisesSecurityIdentifier())
	} else {
		tfStateDevice.OnPremisesSecurityIdentifier = types.StringNull()
	}
	if result.GetOnPremisesSyncEnabled() != nil {
		tfStateDevice.OnPremisesSyncEnabled = types.BoolValue(*result.GetOnPremisesSyncEnabled())
	} else {
		tfStateDevice.OnPremisesSyncEnabled = types.BoolNull()
	}
	if result.GetOperatingSystem() != nil {
		tfStateDevice.OperatingSystem = types.StringValue(*result.GetOperatingSystem())
	} else {
		tfStateDevice.OperatingSystem = types.StringNull()
	}
	if result.GetOperatingSystemVersion() != nil {
		tfStateDevice.OperatingSystemVersion = types.StringValue(*result.GetOperatingSystemVersion())
	} else {
		tfStateDevice.OperatingSystemVersion = types.StringNull()
	}
	if len(result.GetPhysicalIds()) > 0 {
		var valueArrayPhysicalIds []attr.Value
		for _, v := range result.GetPhysicalIds() {
			valueArrayPhysicalIds = append(valueArrayPhysicalIds, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, valueArrayPhysicalIds)
		tfStateDevice.PhysicalIds = listValue
	} else {
		tfStateDevice.PhysicalIds = types.ListNull(types.StringType)
	}
	if result.GetProfileType() != nil {
		tfStateDevice.ProfileType = types.StringValue(*result.GetProfileType())
	} else {
		tfStateDevice.ProfileType = types.StringNull()
	}
	if result.GetRegistrationDateTime() != nil {
		tfStateDevice.RegistrationDateTime = types.StringValue(result.GetRegistrationDateTime().String())
	} else {
		tfStateDevice.RegistrationDateTime = types.StringNull()
	}
	if len(result.GetSystemLabels()) > 0 {
		var valueArraySystemLabels []attr.Value
		for _, v := range result.GetSystemLabels() {
			valueArraySystemLabels = append(valueArraySystemLabels, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, valueArraySystemLabels)
		tfStateDevice.SystemLabels = listValue
	} else {
		tfStateDevice.SystemLabels = types.ListNull(types.StringType)
	}
	if result.GetTrustType() != nil {
		tfStateDevice.TrustType = types.StringValue(*result.GetTrustType())
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

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanDevice deviceModel
	diags := req.Plan.Get(ctx, &tfPlanDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfStateDevice deviceModel
	diags = req.State.Get(ctx, &tfStateDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBodyDevice := models.NewDevice()

	if !tfPlanDevice.Id.Equal(tfStateDevice.Id) {
		tfPlanId := tfPlanDevice.Id.ValueString()
		requestBodyDevice.SetId(&tfPlanId)
	}

	if !tfPlanDevice.DeletedDateTime.Equal(tfStateDevice.DeletedDateTime) {
		tfPlanDeletedDateTime := tfPlanDevice.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanDeletedDateTime)
		requestBodyDevice.SetDeletedDateTime(&t)
	}

	if !tfPlanDevice.AccountEnabled.Equal(tfStateDevice.AccountEnabled) {
		tfPlanAccountEnabled := tfPlanDevice.AccountEnabled.ValueBool()
		requestBodyDevice.SetAccountEnabled(&tfPlanAccountEnabled)
	}

	if !tfPlanDevice.AlternativeSecurityIds.Equal(tfStateDevice.AlternativeSecurityIds) {
		var tfPlanAlternativeSecurityIds []models.AlternativeSecurityIdable
		for k, i := range tfPlanDevice.AlternativeSecurityIds.Elements() {
			requestBodyAlternativeSecurityId := models.NewAlternativeSecurityId()
			tfPlanAlternativeSecurityId := deviceAlternativeSecurityIdModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlanAlternativeSecurityId)
			tfStateAlternativeSecurityId := deviceAlternativeSecurityIdModel{}
			types.ListValueFrom(ctx, tfStateDevice.AlternativeSecurityIds.Elements()[k].Type(ctx), &tfPlanAlternativeSecurityId)

			if !tfPlanAlternativeSecurityId.IdentityProvider.Equal(tfStateAlternativeSecurityId.IdentityProvider) {
				tfPlanIdentityProvider := tfPlanAlternativeSecurityId.IdentityProvider.ValueString()
				requestBodyAlternativeSecurityId.SetIdentityProvider(&tfPlanIdentityProvider)
			}

			if !tfPlanAlternativeSecurityId.Key.Equal(tfStateAlternativeSecurityId.Key) {
				tfPlanKey := tfPlanAlternativeSecurityId.Key.ValueString()
				requestBodyAlternativeSecurityId.SetKey([]byte(tfPlanKey))
			}
		}
		requestBodyDevice.SetAlternativeSecurityIds(tfPlanAlternativeSecurityIds)
	}

	if !tfPlanDevice.ApproximateLastSignInDateTime.Equal(tfStateDevice.ApproximateLastSignInDateTime) {
		tfPlanApproximateLastSignInDateTime := tfPlanDevice.ApproximateLastSignInDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanApproximateLastSignInDateTime)
		requestBodyDevice.SetApproximateLastSignInDateTime(&t)
	}

	if !tfPlanDevice.ComplianceExpirationDateTime.Equal(tfStateDevice.ComplianceExpirationDateTime) {
		tfPlanComplianceExpirationDateTime := tfPlanDevice.ComplianceExpirationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanComplianceExpirationDateTime)
		requestBodyDevice.SetComplianceExpirationDateTime(&t)
	}

	if !tfPlanDevice.DeviceCategory.Equal(tfStateDevice.DeviceCategory) {
		tfPlanDeviceCategory := tfPlanDevice.DeviceCategory.ValueString()
		requestBodyDevice.SetDeviceCategory(&tfPlanDeviceCategory)
	}

	if !tfPlanDevice.DeviceId.Equal(tfStateDevice.DeviceId) {
		tfPlanDeviceId := tfPlanDevice.DeviceId.ValueString()
		requestBodyDevice.SetDeviceId(&tfPlanDeviceId)
	}

	if !tfPlanDevice.DeviceMetadata.Equal(tfStateDevice.DeviceMetadata) {
		tfPlanDeviceMetadata := tfPlanDevice.DeviceMetadata.ValueString()
		requestBodyDevice.SetDeviceMetadata(&tfPlanDeviceMetadata)
	}

	if !tfPlanDevice.DeviceOwnership.Equal(tfStateDevice.DeviceOwnership) {
		tfPlanDeviceOwnership := tfPlanDevice.DeviceOwnership.ValueString()
		requestBodyDevice.SetDeviceOwnership(&tfPlanDeviceOwnership)
	}

	if !tfPlanDevice.DisplayName.Equal(tfStateDevice.DisplayName) {
		tfPlanDisplayName := tfPlanDevice.DisplayName.ValueString()
		requestBodyDevice.SetDisplayName(&tfPlanDisplayName)
	}

	if !tfPlanDevice.EnrollmentProfileName.Equal(tfStateDevice.EnrollmentProfileName) {
		tfPlanEnrollmentProfileName := tfPlanDevice.EnrollmentProfileName.ValueString()
		requestBodyDevice.SetEnrollmentProfileName(&tfPlanEnrollmentProfileName)
	}

	if !tfPlanDevice.EnrollmentType.Equal(tfStateDevice.EnrollmentType) {
		tfPlanEnrollmentType := tfPlanDevice.EnrollmentType.ValueString()
		requestBodyDevice.SetEnrollmentType(&tfPlanEnrollmentType)
	}

	if !tfPlanDevice.IsCompliant.Equal(tfStateDevice.IsCompliant) {
		tfPlanIsCompliant := tfPlanDevice.IsCompliant.ValueBool()
		requestBodyDevice.SetIsCompliant(&tfPlanIsCompliant)
	}

	if !tfPlanDevice.IsManaged.Equal(tfStateDevice.IsManaged) {
		tfPlanIsManaged := tfPlanDevice.IsManaged.ValueBool()
		requestBodyDevice.SetIsManaged(&tfPlanIsManaged)
	}

	if !tfPlanDevice.IsManagementRestricted.Equal(tfStateDevice.IsManagementRestricted) {
		tfPlanIsManagementRestricted := tfPlanDevice.IsManagementRestricted.ValueBool()
		requestBodyDevice.SetIsManagementRestricted(&tfPlanIsManagementRestricted)
	}

	if !tfPlanDevice.IsRooted.Equal(tfStateDevice.IsRooted) {
		tfPlanIsRooted := tfPlanDevice.IsRooted.ValueBool()
		requestBodyDevice.SetIsRooted(&tfPlanIsRooted)
	}

	if !tfPlanDevice.ManagementType.Equal(tfStateDevice.ManagementType) {
		tfPlanManagementType := tfPlanDevice.ManagementType.ValueString()
		requestBodyDevice.SetManagementType(&tfPlanManagementType)
	}

	if !tfPlanDevice.Manufacturer.Equal(tfStateDevice.Manufacturer) {
		tfPlanManufacturer := tfPlanDevice.Manufacturer.ValueString()
		requestBodyDevice.SetManufacturer(&tfPlanManufacturer)
	}

	if !tfPlanDevice.MdmAppId.Equal(tfStateDevice.MdmAppId) {
		tfPlanMdmAppId := tfPlanDevice.MdmAppId.ValueString()
		requestBodyDevice.SetMdmAppId(&tfPlanMdmAppId)
	}

	if !tfPlanDevice.Model.Equal(tfStateDevice.Model) {
		tfPlanModel := tfPlanDevice.Model.ValueString()
		requestBodyDevice.SetModel(&tfPlanModel)
	}

	if !tfPlanDevice.OnPremisesLastSyncDateTime.Equal(tfStateDevice.OnPremisesLastSyncDateTime) {
		tfPlanOnPremisesLastSyncDateTime := tfPlanDevice.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanOnPremisesLastSyncDateTime)
		requestBodyDevice.SetOnPremisesLastSyncDateTime(&t)
	}

	if !tfPlanDevice.OnPremisesSecurityIdentifier.Equal(tfStateDevice.OnPremisesSecurityIdentifier) {
		tfPlanOnPremisesSecurityIdentifier := tfPlanDevice.OnPremisesSecurityIdentifier.ValueString()
		requestBodyDevice.SetOnPremisesSecurityIdentifier(&tfPlanOnPremisesSecurityIdentifier)
	}

	if !tfPlanDevice.OnPremisesSyncEnabled.Equal(tfStateDevice.OnPremisesSyncEnabled) {
		tfPlanOnPremisesSyncEnabled := tfPlanDevice.OnPremisesSyncEnabled.ValueBool()
		requestBodyDevice.SetOnPremisesSyncEnabled(&tfPlanOnPremisesSyncEnabled)
	}

	if !tfPlanDevice.OperatingSystem.Equal(tfStateDevice.OperatingSystem) {
		tfPlanOperatingSystem := tfPlanDevice.OperatingSystem.ValueString()
		requestBodyDevice.SetOperatingSystem(&tfPlanOperatingSystem)
	}

	if !tfPlanDevice.OperatingSystemVersion.Equal(tfStateDevice.OperatingSystemVersion) {
		tfPlanOperatingSystemVersion := tfPlanDevice.OperatingSystemVersion.ValueString()
		requestBodyDevice.SetOperatingSystemVersion(&tfPlanOperatingSystemVersion)
	}

	if !tfPlanDevice.PhysicalIds.Equal(tfStateDevice.PhysicalIds) {
		var stringArrayPhysicalIds []string
		for _, i := range tfPlanDevice.PhysicalIds.Elements() {
			stringArrayPhysicalIds = append(stringArrayPhysicalIds, i.String())
		}
		requestBodyDevice.SetPhysicalIds(stringArrayPhysicalIds)
	}

	if !tfPlanDevice.ProfileType.Equal(tfStateDevice.ProfileType) {
		tfPlanProfileType := tfPlanDevice.ProfileType.ValueString()
		requestBodyDevice.SetProfileType(&tfPlanProfileType)
	}

	if !tfPlanDevice.RegistrationDateTime.Equal(tfStateDevice.RegistrationDateTime) {
		tfPlanRegistrationDateTime := tfPlanDevice.RegistrationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanRegistrationDateTime)
		requestBodyDevice.SetRegistrationDateTime(&t)
	}

	if !tfPlanDevice.SystemLabels.Equal(tfStateDevice.SystemLabels) {
		var stringArraySystemLabels []string
		for _, i := range tfPlanDevice.SystemLabels.Elements() {
			stringArraySystemLabels = append(stringArraySystemLabels, i.String())
		}
		requestBodyDevice.SetSystemLabels(stringArraySystemLabels)
	}

	if !tfPlanDevice.TrustType.Equal(tfStateDevice.TrustType) {
		tfPlanTrustType := tfPlanDevice.TrustType.ValueString()
		requestBodyDevice.SetTrustType(&tfPlanTrustType)
	}

	// Update device
	_, err := r.client.Devices().ByDeviceId(tfStateDevice.Id.ValueString()).Patch(context.Background(), requestBodyDevice, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlanDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfStateDevice deviceModel
	diags := req.State.Get(ctx, &tfStateDevice)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete device
	err := r.client.Devices().ByDeviceId(tfStateDevice.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting device",
			err.Error(),
		)
		return
	}

}
