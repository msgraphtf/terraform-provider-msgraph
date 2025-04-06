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
	var state deviceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
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

	if !state.Id.IsNull() {
		result, err = d.client.Devices().ByDeviceId(state.Id.ValueString()).Get(context.Background(), &qparams)
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
		state.Id = types.StringValue(*result.GetId())
	} else {
		state.Id = types.StringNull()
	}
	if result.GetDeletedDateTime() != nil {
		state.DeletedDateTime = types.StringValue(result.GetDeletedDateTime().String())
	} else {
		state.DeletedDateTime = types.StringNull()
	}
	if result.GetAccountEnabled() != nil {
		state.AccountEnabled = types.BoolValue(*result.GetAccountEnabled())
	} else {
		state.AccountEnabled = types.BoolNull()
	}
	if len(result.GetAlternativeSecurityIds()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, v := range result.GetAlternativeSecurityIds() {
			alternativeSecurityIds := new(deviceAlternativeSecurityIdModel)

			if v.GetIdentityProvider() != nil {
				alternativeSecurityIds.IdentityProvider = types.StringValue(*v.GetIdentityProvider())
			} else {
				alternativeSecurityIds.IdentityProvider = types.StringNull()
			}
			if v.GetKey() != nil {
				alternativeSecurityIds.Key = types.StringValue(string(v.GetKey()[:]))
			} else {
				alternativeSecurityIds.Key = types.StringNull()
			}
			objectValue, _ := types.ObjectValueFrom(ctx, alternativeSecurityIds.AttributeTypes(), alternativeSecurityIds)
			objectValues = append(objectValues, objectValue)
		}
		state.AlternativeSecurityIds, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}
	if result.GetApproximateLastSignInDateTime() != nil {
		state.ApproximateLastSignInDateTime = types.StringValue(result.GetApproximateLastSignInDateTime().String())
	} else {
		state.ApproximateLastSignInDateTime = types.StringNull()
	}
	if result.GetComplianceExpirationDateTime() != nil {
		state.ComplianceExpirationDateTime = types.StringValue(result.GetComplianceExpirationDateTime().String())
	} else {
		state.ComplianceExpirationDateTime = types.StringNull()
	}
	if result.GetDeviceCategory() != nil {
		state.DeviceCategory = types.StringValue(*result.GetDeviceCategory())
	} else {
		state.DeviceCategory = types.StringNull()
	}
	if result.GetDeviceId() != nil {
		state.DeviceId = types.StringValue(*result.GetDeviceId())
	} else {
		state.DeviceId = types.StringNull()
	}
	if result.GetDeviceMetadata() != nil {
		state.DeviceMetadata = types.StringValue(*result.GetDeviceMetadata())
	} else {
		state.DeviceMetadata = types.StringNull()
	}
	if result.GetDeviceOwnership() != nil {
		state.DeviceOwnership = types.StringValue(*result.GetDeviceOwnership())
	} else {
		state.DeviceOwnership = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		state.DisplayName = types.StringNull()
	}
	if result.GetEnrollmentProfileName() != nil {
		state.EnrollmentProfileName = types.StringValue(*result.GetEnrollmentProfileName())
	} else {
		state.EnrollmentProfileName = types.StringNull()
	}
	if result.GetEnrollmentType() != nil {
		state.EnrollmentType = types.StringValue(*result.GetEnrollmentType())
	} else {
		state.EnrollmentType = types.StringNull()
	}
	if result.GetIsCompliant() != nil {
		state.IsCompliant = types.BoolValue(*result.GetIsCompliant())
	} else {
		state.IsCompliant = types.BoolNull()
	}
	if result.GetIsManaged() != nil {
		state.IsManaged = types.BoolValue(*result.GetIsManaged())
	} else {
		state.IsManaged = types.BoolNull()
	}
	if result.GetIsManagementRestricted() != nil {
		state.IsManagementRestricted = types.BoolValue(*result.GetIsManagementRestricted())
	} else {
		state.IsManagementRestricted = types.BoolNull()
	}
	if result.GetIsRooted() != nil {
		state.IsRooted = types.BoolValue(*result.GetIsRooted())
	} else {
		state.IsRooted = types.BoolNull()
	}
	if result.GetManagementType() != nil {
		state.ManagementType = types.StringValue(*result.GetManagementType())
	} else {
		state.ManagementType = types.StringNull()
	}
	if result.GetManufacturer() != nil {
		state.Manufacturer = types.StringValue(*result.GetManufacturer())
	} else {
		state.Manufacturer = types.StringNull()
	}
	if result.GetMdmAppId() != nil {
		state.MdmAppId = types.StringValue(*result.GetMdmAppId())
	} else {
		state.MdmAppId = types.StringNull()
	}
	if result.GetModel() != nil {
		state.Model = types.StringValue(*result.GetModel())
	} else {
		state.Model = types.StringNull()
	}
	if result.GetOnPremisesLastSyncDateTime() != nil {
		state.OnPremisesLastSyncDateTime = types.StringValue(result.GetOnPremisesLastSyncDateTime().String())
	} else {
		state.OnPremisesLastSyncDateTime = types.StringNull()
	}
	if result.GetOnPremisesSecurityIdentifier() != nil {
		state.OnPremisesSecurityIdentifier = types.StringValue(*result.GetOnPremisesSecurityIdentifier())
	} else {
		state.OnPremisesSecurityIdentifier = types.StringNull()
	}
	if result.GetOnPremisesSyncEnabled() != nil {
		state.OnPremisesSyncEnabled = types.BoolValue(*result.GetOnPremisesSyncEnabled())
	} else {
		state.OnPremisesSyncEnabled = types.BoolNull()
	}
	if result.GetOperatingSystem() != nil {
		state.OperatingSystem = types.StringValue(*result.GetOperatingSystem())
	} else {
		state.OperatingSystem = types.StringNull()
	}
	if result.GetOperatingSystemVersion() != nil {
		state.OperatingSystemVersion = types.StringValue(*result.GetOperatingSystemVersion())
	} else {
		state.OperatingSystemVersion = types.StringNull()
	}
	if len(result.GetPhysicalIds()) > 0 {
		var physicalIds []attr.Value
		for _, v := range result.GetPhysicalIds() {
			physicalIds = append(physicalIds, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, physicalIds)
		state.PhysicalIds = listValue
	} else {
		state.PhysicalIds = types.ListNull(types.StringType)
	}
	if result.GetProfileType() != nil {
		state.ProfileType = types.StringValue(*result.GetProfileType())
	} else {
		state.ProfileType = types.StringNull()
	}
	if result.GetRegistrationDateTime() != nil {
		state.RegistrationDateTime = types.StringValue(result.GetRegistrationDateTime().String())
	} else {
		state.RegistrationDateTime = types.StringNull()
	}
	if len(result.GetSystemLabels()) > 0 {
		var systemLabels []attr.Value
		for _, v := range result.GetSystemLabels() {
			systemLabels = append(systemLabels, types.StringValue(v))
		}
		listValue, _ := types.ListValue(types.StringType, systemLabels)
		state.SystemLabels = listValue
	} else {
		state.SystemLabels = types.ListNull(types.StringType)
	}
	if result.GetTrustType() != nil {
		state.TrustType = types.StringValue(*result.GetTrustType())
	} else {
		state.TrustType = types.StringNull()
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan deviceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state deviceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewDevice()

	if !plan.Id.Equal(state.Id) {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	}

	if !plan.DeletedDateTime.Equal(state.DeletedDateTime) {
		planDeletedDateTime := plan.DeletedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planDeletedDateTime)
		requestBody.SetDeletedDateTime(&t)
	}

	if !plan.AccountEnabled.Equal(state.AccountEnabled) {
		planAccountEnabled := plan.AccountEnabled.ValueBool()
		requestBody.SetAccountEnabled(&planAccountEnabled)
	}

	if !plan.AlternativeSecurityIds.Equal(state.AlternativeSecurityIds) {
		var planAlternativeSecurityIds []models.AlternativeSecurityIdable
		for k, i := range plan.AlternativeSecurityIds.Elements() {
			alternativeSecurityIds := models.NewAlternativeSecurityId()
			alternativeSecurityIdsModel := deviceAlternativeSecurityIdModel{}
			types.ListValueFrom(ctx, i.Type(ctx), &alternativeSecurityIdsModel)
			alternativeSecurityIdsState := deviceAlternativeSecurityIdModel{}
			types.ListValueFrom(ctx, state.AlternativeSecurityIds.Elements()[k].Type(ctx), &alternativeSecurityIdsModel)

			if !alternativeSecurityIdsModel.IdentityProvider.Equal(alternativeSecurityIdsState.IdentityProvider) {
				planIdentityProvider := alternativeSecurityIdsModel.IdentityProvider.ValueString()
				alternativeSecurityIds.SetIdentityProvider(&planIdentityProvider)
			}

			if !alternativeSecurityIdsModel.Key.Equal(alternativeSecurityIdsState.Key) {
				planKey := alternativeSecurityIdsModel.Key.ValueString()
				alternativeSecurityIds.SetKey([]byte(planKey))
			}
		}
		requestBody.SetAlternativeSecurityIds(planAlternativeSecurityIds)
	}

	if !plan.ApproximateLastSignInDateTime.Equal(state.ApproximateLastSignInDateTime) {
		planApproximateLastSignInDateTime := plan.ApproximateLastSignInDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planApproximateLastSignInDateTime)
		requestBody.SetApproximateLastSignInDateTime(&t)
	}

	if !plan.ComplianceExpirationDateTime.Equal(state.ComplianceExpirationDateTime) {
		planComplianceExpirationDateTime := plan.ComplianceExpirationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planComplianceExpirationDateTime)
		requestBody.SetComplianceExpirationDateTime(&t)
	}

	if !plan.DeviceCategory.Equal(state.DeviceCategory) {
		planDeviceCategory := plan.DeviceCategory.ValueString()
		requestBody.SetDeviceCategory(&planDeviceCategory)
	}

	if !plan.DeviceId.Equal(state.DeviceId) {
		planDeviceId := plan.DeviceId.ValueString()
		requestBody.SetDeviceId(&planDeviceId)
	}

	if !plan.DeviceMetadata.Equal(state.DeviceMetadata) {
		planDeviceMetadata := plan.DeviceMetadata.ValueString()
		requestBody.SetDeviceMetadata(&planDeviceMetadata)
	}

	if !plan.DeviceOwnership.Equal(state.DeviceOwnership) {
		planDeviceOwnership := plan.DeviceOwnership.ValueString()
		requestBody.SetDeviceOwnership(&planDeviceOwnership)
	}

	if !plan.DisplayName.Equal(state.DisplayName) {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	}

	if !plan.EnrollmentProfileName.Equal(state.EnrollmentProfileName) {
		planEnrollmentProfileName := plan.EnrollmentProfileName.ValueString()
		requestBody.SetEnrollmentProfileName(&planEnrollmentProfileName)
	}

	if !plan.EnrollmentType.Equal(state.EnrollmentType) {
		planEnrollmentType := plan.EnrollmentType.ValueString()
		requestBody.SetEnrollmentType(&planEnrollmentType)
	}

	if !plan.IsCompliant.Equal(state.IsCompliant) {
		planIsCompliant := plan.IsCompliant.ValueBool()
		requestBody.SetIsCompliant(&planIsCompliant)
	}

	if !plan.IsManaged.Equal(state.IsManaged) {
		planIsManaged := plan.IsManaged.ValueBool()
		requestBody.SetIsManaged(&planIsManaged)
	}

	if !plan.IsManagementRestricted.Equal(state.IsManagementRestricted) {
		planIsManagementRestricted := plan.IsManagementRestricted.ValueBool()
		requestBody.SetIsManagementRestricted(&planIsManagementRestricted)
	}

	if !plan.IsRooted.Equal(state.IsRooted) {
		planIsRooted := plan.IsRooted.ValueBool()
		requestBody.SetIsRooted(&planIsRooted)
	}

	if !plan.ManagementType.Equal(state.ManagementType) {
		planManagementType := plan.ManagementType.ValueString()
		requestBody.SetManagementType(&planManagementType)
	}

	if !plan.Manufacturer.Equal(state.Manufacturer) {
		planManufacturer := plan.Manufacturer.ValueString()
		requestBody.SetManufacturer(&planManufacturer)
	}

	if !plan.MdmAppId.Equal(state.MdmAppId) {
		planMdmAppId := plan.MdmAppId.ValueString()
		requestBody.SetMdmAppId(&planMdmAppId)
	}

	if !plan.Model.Equal(state.Model) {
		planModel := plan.Model.ValueString()
		requestBody.SetModel(&planModel)
	}

	if !plan.OnPremisesLastSyncDateTime.Equal(state.OnPremisesLastSyncDateTime) {
		planOnPremisesLastSyncDateTime := plan.OnPremisesLastSyncDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planOnPremisesLastSyncDateTime)
		requestBody.SetOnPremisesLastSyncDateTime(&t)
	}

	if !plan.OnPremisesSecurityIdentifier.Equal(state.OnPremisesSecurityIdentifier) {
		planOnPremisesSecurityIdentifier := plan.OnPremisesSecurityIdentifier.ValueString()
		requestBody.SetOnPremisesSecurityIdentifier(&planOnPremisesSecurityIdentifier)
	}

	if !plan.OnPremisesSyncEnabled.Equal(state.OnPremisesSyncEnabled) {
		planOnPremisesSyncEnabled := plan.OnPremisesSyncEnabled.ValueBool()
		requestBody.SetOnPremisesSyncEnabled(&planOnPremisesSyncEnabled)
	}

	if !plan.OperatingSystem.Equal(state.OperatingSystem) {
		planOperatingSystem := plan.OperatingSystem.ValueString()
		requestBody.SetOperatingSystem(&planOperatingSystem)
	}

	if !plan.OperatingSystemVersion.Equal(state.OperatingSystemVersion) {
		planOperatingSystemVersion := plan.OperatingSystemVersion.ValueString()
		requestBody.SetOperatingSystemVersion(&planOperatingSystemVersion)
	}

	if !plan.PhysicalIds.Equal(state.PhysicalIds) {
		var stringArrayPhysicalIds []string
		for _, i := range plan.PhysicalIds.Elements() {
			stringArrayPhysicalIds = append(stringArrayPhysicalIds, i.String())
		}
		requestBody.SetPhysicalIds(stringArrayPhysicalIds)
	}

	if !plan.ProfileType.Equal(state.ProfileType) {
		planProfileType := plan.ProfileType.ValueString()
		requestBody.SetProfileType(&planProfileType)
	}

	if !plan.RegistrationDateTime.Equal(state.RegistrationDateTime) {
		planRegistrationDateTime := plan.RegistrationDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planRegistrationDateTime)
		requestBody.SetRegistrationDateTime(&t)
	}

	if !plan.SystemLabels.Equal(state.SystemLabels) {
		var stringArraySystemLabels []string
		for _, i := range plan.SystemLabels.Elements() {
			stringArraySystemLabels = append(stringArraySystemLabels, i.String())
		}
		requestBody.SetSystemLabels(stringArraySystemLabels)
	}

	if !plan.TrustType.Equal(state.TrustType) {
		planTrustType := plan.TrustType.ValueString()
		requestBody.SetTrustType(&planTrustType)
	}

	// Update device
	_, err := r.client.Devices().ByDeviceId(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating device",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deviceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete device
	err := r.client.Devices().ByDeviceId(state.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting device",
			err.Error(),
		)
		return
	}

}
