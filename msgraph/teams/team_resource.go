package teams

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"time"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/teams"

	"terraform-provider-msgraph/planmodifiers/boolplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/objectplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/stringplanmodifiers"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &teamResource{}
	_ resource.ResourceWithConfigure = &teamResource{}
)

// NewTeamResource is a helper function to simplify the provider implementation.
func NewTeamResource() resource.Resource {
	return &teamResource{}
}

// teamResource is the resource implementation.
type teamResource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (d *teamResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Configure adds the provider configured client to the resource.
func (d *teamResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the resource.
func (d *teamResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"classification": schema.StringAttribute{
				Description: "An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"created_date_time": schema.StringAttribute{
				Description: "Timestamp at which the team was created.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"description": schema.StringAttribute{
				Description: "An optional description for the team. Maximum length: 1024 characters.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"display_name": schema.StringAttribute{
				Description: "The name of the team.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"fun_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure use of Giphy, memes, and stickers in the team.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"allow_custom_memes": schema.BoolAttribute{
						Description: "If set to true, enables users to include custom memes.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_giphy": schema.BoolAttribute{
						Description: "If set to true, enables Giphy use.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_stickers_and_memes": schema.BoolAttribute{
						Description: "If set to true, enables users to include stickers and memes.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"giphy_content_rating": schema.StringAttribute{
						Description: "Giphy content rating. Possible values are: moderate, strict.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"guest_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure whether guests can create, update, or delete channels in the team.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"allow_create_update_channels": schema.BoolAttribute{
						Description: "If set to true, guests can add and update channels.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_delete_channels": schema.BoolAttribute{
						Description: "If set to true, guests can delete channels.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"internal_id": schema.StringAttribute{
				Description: "A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"is_archived": schema.BoolAttribute{
				Description: "Whether this team is in read-only mode.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"member_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"allow_add_remove_apps": schema.BoolAttribute{
						Description: "If set to true, members can add and remove apps.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_create_private_channels": schema.BoolAttribute{
						Description: "If set to true, members can add and update private channels.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_create_update_channels": schema.BoolAttribute{
						Description: "If set to true, members can add and update channels.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_create_update_remove_connectors": schema.BoolAttribute{
						Description: "If set to true, members can add, update, and remove connectors.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_create_update_remove_tabs": schema.BoolAttribute{
						Description: "If set to true, members can add, update, and remove tabs.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_delete_channels": schema.BoolAttribute{
						Description: "If set to true, members can delete channels.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"messaging_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure messaging and mentions in the team.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{
					"allow_channel_mentions": schema.BoolAttribute{
						Description: "If set to true, @channel mentions are allowed.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_owner_delete_messages": schema.BoolAttribute{
						Description: "If set to true, owners can delete any message.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_team_mentions": schema.BoolAttribute{
						Description: "If set to true, @team mentions are allowed.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_user_delete_messages": schema.BoolAttribute{
						Description: "If set to true, users can delete their messages.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
					"allow_user_edit_messages": schema.BoolAttribute{
						Description: "If set to true, users can edit their messages.",
						Optional:    true,
						Computed:    true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifiers.UseStateForUnconfigured(),
						},
					},
				},
			},
			"specialization": schema.StringAttribute{
				Description: "Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the Microsoft Entra tenant.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"visibility": schema.StringAttribute{
				Description: "The visibility of the group and team. Defaults to Public.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
			"web_url": schema.StringAttribute{
				Description: "A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.",
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
func (r *teamResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanTeam teamModel
	diags := req.Plan.Get(ctx, &tfPlanTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	requestBodyTeam := models.NewTeam()
	if !tfPlanTeam.Classification.IsUnknown() {
		tfPlanClassification := tfPlanTeam.Classification.ValueString()
		requestBodyTeam.SetClassification(&tfPlanClassification)
	} else {
		tfPlanTeam.Classification = types.StringNull()
	}

	if !tfPlanTeam.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlanTeam.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyTeam.SetCreatedDateTime(&t)
	} else {
		tfPlanTeam.CreatedDateTime = types.StringNull()
	}

	if !tfPlanTeam.Description.IsUnknown() {
		tfPlanDescription := tfPlanTeam.Description.ValueString()
		requestBodyTeam.SetDescription(&tfPlanDescription)
	} else {
		tfPlanTeam.Description = types.StringNull()
	}

	if !tfPlanTeam.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanTeam.DisplayName.ValueString()
		requestBodyTeam.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanTeam.DisplayName = types.StringNull()
	}

	if !tfPlanTeam.FunSettings.IsUnknown() {
		requestBodyTeamFunSettings := models.NewTeamFunSettings()
		tfPlanTeamFunSettings := teamTeamFunSettingsModel{}
		tfPlanTeam.FunSettings.As(ctx, &tfPlanTeamFunSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamFunSettings.AllowCustomMemes.IsUnknown() {
			tfPlanAllowCustomMemes := tfPlanTeamFunSettings.AllowCustomMemes.ValueBool()
			requestBodyTeamFunSettings.SetAllowCustomMemes(&tfPlanAllowCustomMemes)
		} else {
			tfPlanTeamFunSettings.AllowCustomMemes = types.BoolNull()
		}

		if !tfPlanTeamFunSettings.AllowGiphy.IsUnknown() {
			tfPlanAllowGiphy := tfPlanTeamFunSettings.AllowGiphy.ValueBool()
			requestBodyTeamFunSettings.SetAllowGiphy(&tfPlanAllowGiphy)
		} else {
			tfPlanTeamFunSettings.AllowGiphy = types.BoolNull()
		}

		if !tfPlanTeamFunSettings.AllowStickersAndMemes.IsUnknown() {
			tfPlanAllowStickersAndMemes := tfPlanTeamFunSettings.AllowStickersAndMemes.ValueBool()
			requestBodyTeamFunSettings.SetAllowStickersAndMemes(&tfPlanAllowStickersAndMemes)
		} else {
			tfPlanTeamFunSettings.AllowStickersAndMemes = types.BoolNull()
		}

		if !tfPlanTeamFunSettings.GiphyContentRating.IsUnknown() {
			tfPlanGiphyContentRating := tfPlanTeamFunSettings.GiphyContentRating.ValueString()
			parsedGiphyContentRating, _ := models.ParseGiphyRatingType(tfPlanGiphyContentRating)
			assertedGiphyContentRating := parsedGiphyContentRating.(models.GiphyRatingType)
			requestBodyTeamFunSettings.SetGiphyContentRating(&assertedGiphyContentRating)
		} else {
			tfPlanTeamFunSettings.GiphyContentRating = types.StringNull()
		}

		requestBodyTeam.SetFunSettings(requestBodyTeamFunSettings)
		tfPlanTeam.FunSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamFunSettings.AttributeTypes(), requestBodyTeamFunSettings)
	} else {
		tfPlanTeam.FunSettings = types.ObjectNull(tfPlanTeam.FunSettings.AttributeTypes(ctx))
	}

	if !tfPlanTeam.GuestSettings.IsUnknown() {
		requestBodyTeamGuestSettings := models.NewTeamGuestSettings()
		tfPlanTeamGuestSettings := teamTeamGuestSettingsModel{}
		tfPlanTeam.GuestSettings.As(ctx, &tfPlanTeamGuestSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamGuestSettings.AllowCreateUpdateChannels.IsUnknown() {
			tfPlanAllowCreateUpdateChannels := tfPlanTeamGuestSettings.AllowCreateUpdateChannels.ValueBool()
			requestBodyTeamGuestSettings.SetAllowCreateUpdateChannels(&tfPlanAllowCreateUpdateChannels)
		} else {
			tfPlanTeamGuestSettings.AllowCreateUpdateChannels = types.BoolNull()
		}

		if !tfPlanTeamGuestSettings.AllowDeleteChannels.IsUnknown() {
			tfPlanAllowDeleteChannels := tfPlanTeamGuestSettings.AllowDeleteChannels.ValueBool()
			requestBodyTeamGuestSettings.SetAllowDeleteChannels(&tfPlanAllowDeleteChannels)
		} else {
			tfPlanTeamGuestSettings.AllowDeleteChannels = types.BoolNull()
		}

		requestBodyTeam.SetGuestSettings(requestBodyTeamGuestSettings)
		tfPlanTeam.GuestSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamGuestSettings.AttributeTypes(), requestBodyTeamGuestSettings)
	} else {
		tfPlanTeam.GuestSettings = types.ObjectNull(tfPlanTeam.GuestSettings.AttributeTypes(ctx))
	}

	if !tfPlanTeam.Id.IsUnknown() {
		tfPlanId := tfPlanTeam.Id.ValueString()
		requestBodyTeam.SetId(&tfPlanId)
	} else {
		tfPlanTeam.Id = types.StringNull()
	}

	if !tfPlanTeam.InternalId.IsUnknown() {
		tfPlanInternalId := tfPlanTeam.InternalId.ValueString()
		requestBodyTeam.SetInternalId(&tfPlanInternalId)
	} else {
		tfPlanTeam.InternalId = types.StringNull()
	}

	if !tfPlanTeam.IsArchived.IsUnknown() {
		tfPlanIsArchived := tfPlanTeam.IsArchived.ValueBool()
		requestBodyTeam.SetIsArchived(&tfPlanIsArchived)
	} else {
		tfPlanTeam.IsArchived = types.BoolNull()
	}

	if !tfPlanTeam.MemberSettings.IsUnknown() {
		requestBodyTeamMemberSettings := models.NewTeamMemberSettings()
		tfPlanTeamMemberSettings := teamTeamMemberSettingsModel{}
		tfPlanTeam.MemberSettings.As(ctx, &tfPlanTeamMemberSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamMemberSettings.AllowAddRemoveApps.IsUnknown() {
			tfPlanAllowAddRemoveApps := tfPlanTeamMemberSettings.AllowAddRemoveApps.ValueBool()
			requestBodyTeamMemberSettings.SetAllowAddRemoveApps(&tfPlanAllowAddRemoveApps)
		} else {
			tfPlanTeamMemberSettings.AllowAddRemoveApps = types.BoolNull()
		}

		if !tfPlanTeamMemberSettings.AllowCreatePrivateChannels.IsUnknown() {
			tfPlanAllowCreatePrivateChannels := tfPlanTeamMemberSettings.AllowCreatePrivateChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreatePrivateChannels(&tfPlanAllowCreatePrivateChannels)
		} else {
			tfPlanTeamMemberSettings.AllowCreatePrivateChannels = types.BoolNull()
		}

		if !tfPlanTeamMemberSettings.AllowCreateUpdateChannels.IsUnknown() {
			tfPlanAllowCreateUpdateChannels := tfPlanTeamMemberSettings.AllowCreateUpdateChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateChannels(&tfPlanAllowCreateUpdateChannels)
		} else {
			tfPlanTeamMemberSettings.AllowCreateUpdateChannels = types.BoolNull()
		}

		if !tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors.IsUnknown() {
			tfPlanAllowCreateUpdateRemoveConnectors := tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateRemoveConnectors(&tfPlanAllowCreateUpdateRemoveConnectors)
		} else {
			tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors = types.BoolNull()
		}

		if !tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs.IsUnknown() {
			tfPlanAllowCreateUpdateRemoveTabs := tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateRemoveTabs(&tfPlanAllowCreateUpdateRemoveTabs)
		} else {
			tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs = types.BoolNull()
		}

		if !tfPlanTeamMemberSettings.AllowDeleteChannels.IsUnknown() {
			tfPlanAllowDeleteChannels := tfPlanTeamMemberSettings.AllowDeleteChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowDeleteChannels(&tfPlanAllowDeleteChannels)
		} else {
			tfPlanTeamMemberSettings.AllowDeleteChannels = types.BoolNull()
		}

		requestBodyTeam.SetMemberSettings(requestBodyTeamMemberSettings)
		tfPlanTeam.MemberSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamMemberSettings.AttributeTypes(), requestBodyTeamMemberSettings)
	} else {
		tfPlanTeam.MemberSettings = types.ObjectNull(tfPlanTeam.MemberSettings.AttributeTypes(ctx))
	}

	if !tfPlanTeam.MessagingSettings.IsUnknown() {
		requestBodyTeamMessagingSettings := models.NewTeamMessagingSettings()
		tfPlanTeamMessagingSettings := teamTeamMessagingSettingsModel{}
		tfPlanTeam.MessagingSettings.As(ctx, &tfPlanTeamMessagingSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamMessagingSettings.AllowChannelMentions.IsUnknown() {
			tfPlanAllowChannelMentions := tfPlanTeamMessagingSettings.AllowChannelMentions.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowChannelMentions(&tfPlanAllowChannelMentions)
		} else {
			tfPlanTeamMessagingSettings.AllowChannelMentions = types.BoolNull()
		}

		if !tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages.IsUnknown() {
			tfPlanAllowOwnerDeleteMessages := tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowOwnerDeleteMessages(&tfPlanAllowOwnerDeleteMessages)
		} else {
			tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages = types.BoolNull()
		}

		if !tfPlanTeamMessagingSettings.AllowTeamMentions.IsUnknown() {
			tfPlanAllowTeamMentions := tfPlanTeamMessagingSettings.AllowTeamMentions.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowTeamMentions(&tfPlanAllowTeamMentions)
		} else {
			tfPlanTeamMessagingSettings.AllowTeamMentions = types.BoolNull()
		}

		if !tfPlanTeamMessagingSettings.AllowUserDeleteMessages.IsUnknown() {
			tfPlanAllowUserDeleteMessages := tfPlanTeamMessagingSettings.AllowUserDeleteMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowUserDeleteMessages(&tfPlanAllowUserDeleteMessages)
		} else {
			tfPlanTeamMessagingSettings.AllowUserDeleteMessages = types.BoolNull()
		}

		if !tfPlanTeamMessagingSettings.AllowUserEditMessages.IsUnknown() {
			tfPlanAllowUserEditMessages := tfPlanTeamMessagingSettings.AllowUserEditMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowUserEditMessages(&tfPlanAllowUserEditMessages)
		} else {
			tfPlanTeamMessagingSettings.AllowUserEditMessages = types.BoolNull()
		}

		requestBodyTeam.SetMessagingSettings(requestBodyTeamMessagingSettings)
		tfPlanTeam.MessagingSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamMessagingSettings.AttributeTypes(), requestBodyTeamMessagingSettings)
	} else {
		tfPlanTeam.MessagingSettings = types.ObjectNull(tfPlanTeam.MessagingSettings.AttributeTypes(ctx))
	}

	if !tfPlanTeam.Specialization.IsUnknown() {
		tfPlanSpecialization := tfPlanTeam.Specialization.ValueString()
		parsedSpecialization, _ := models.ParseTeamSpecialization(tfPlanSpecialization)
		assertedSpecialization := parsedSpecialization.(models.TeamSpecialization)
		requestBodyTeam.SetSpecialization(&assertedSpecialization)
	} else {
		tfPlanTeam.Specialization = types.StringNull()
	}

	if !tfPlanTeam.TenantId.IsUnknown() {
		tfPlanTenantId := tfPlanTeam.TenantId.ValueString()
		requestBodyTeam.SetTenantId(&tfPlanTenantId)
	} else {
		tfPlanTeam.TenantId = types.StringNull()
	}

	if !tfPlanTeam.Visibility.IsUnknown() {
		tfPlanVisibility := tfPlanTeam.Visibility.ValueString()
		parsedVisibility, _ := models.ParseTeamVisibilityType(tfPlanVisibility)
		assertedVisibility := parsedVisibility.(models.TeamVisibilityType)
		requestBodyTeam.SetVisibility(&assertedVisibility)
	} else {
		tfPlanTeam.Visibility = types.StringNull()
	}

	if !tfPlanTeam.WebUrl.IsUnknown() {
		tfPlanWebUrl := tfPlanTeam.WebUrl.ValueString()
		requestBodyTeam.SetWebUrl(&tfPlanWebUrl)
	} else {
		tfPlanTeam.WebUrl = types.StringNull()
	}

	// Create new Team
	result, err := r.client.Teams().Post(context.Background(), requestBodyTeam, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Team",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlanTeam.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlanTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var tfStateTeam teamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &tfStateTeam)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := teams.TeamItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &teams.TeamItemRequestBuilderGetQueryParameters{
			Select: []string{
				"classification",
				"createdDateTime",
				"description",
				"displayName",
				"funSettings",
				"guestSettings",
				"id",
				"internalId",
				"isArchived",
				"memberSettings",
				"messagingSettings",
				"specialization",
				"tenantId",
				"visibility",
				"webUrl",
			},
		},
	}

	var responseTeam models.Teamable
	var err error

	if !tfStateTeam.Id.IsNull() {
		responseTeam, err = d.client.Teams().ByTeamId(tfStateTeam.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Team",
			err.Error(),
		)
		return
	}

	if responseTeam.GetClassification() != nil {
		tfStateTeam.Classification = types.StringValue(*responseTeam.GetClassification())
	} else {
		tfStateTeam.Classification = types.StringNull()
	}
	if responseTeam.GetCreatedDateTime() != nil {
		tfStateTeam.CreatedDateTime = types.StringValue(responseTeam.GetCreatedDateTime().String())
	} else {
		tfStateTeam.CreatedDateTime = types.StringNull()
	}
	if responseTeam.GetDescription() != nil {
		tfStateTeam.Description = types.StringValue(*responseTeam.GetDescription())
	} else {
		tfStateTeam.Description = types.StringNull()
	}
	if responseTeam.GetDisplayName() != nil {
		tfStateTeam.DisplayName = types.StringValue(*responseTeam.GetDisplayName())
	} else {
		tfStateTeam.DisplayName = types.StringNull()
	}
	if responseTeam.GetFunSettings() != nil {
		tfStateTeamFunSettings := teamTeamFunSettingsModel{}
		responseTeamFunSettings := responseTeam.GetFunSettings()

		if responseTeamFunSettings.GetAllowCustomMemes() != nil {
			tfStateTeamFunSettings.AllowCustomMemes = types.BoolValue(*responseTeamFunSettings.GetAllowCustomMemes())
		} else {
			tfStateTeamFunSettings.AllowCustomMemes = types.BoolNull()
		}
		if responseTeamFunSettings.GetAllowGiphy() != nil {
			tfStateTeamFunSettings.AllowGiphy = types.BoolValue(*responseTeamFunSettings.GetAllowGiphy())
		} else {
			tfStateTeamFunSettings.AllowGiphy = types.BoolNull()
		}
		if responseTeamFunSettings.GetAllowStickersAndMemes() != nil {
			tfStateTeamFunSettings.AllowStickersAndMemes = types.BoolValue(*responseTeamFunSettings.GetAllowStickersAndMemes())
		} else {
			tfStateTeamFunSettings.AllowStickersAndMemes = types.BoolNull()
		}
		if responseTeamFunSettings.GetGiphyContentRating() != nil {
			tfStateTeamFunSettings.GiphyContentRating = types.StringValue(responseTeamFunSettings.GetGiphyContentRating().String())
		} else {
			tfStateTeamFunSettings.GiphyContentRating = types.StringNull()
		}

		tfStateTeam.FunSettings, _ = types.ObjectValueFrom(ctx, tfStateTeamFunSettings.AttributeTypes(), tfStateTeamFunSettings)
	}
	if responseTeam.GetGuestSettings() != nil {
		tfStateTeamGuestSettings := teamTeamGuestSettingsModel{}
		responseTeamGuestSettings := responseTeam.GetGuestSettings()

		if responseTeamGuestSettings.GetAllowCreateUpdateChannels() != nil {
			tfStateTeamGuestSettings.AllowCreateUpdateChannels = types.BoolValue(*responseTeamGuestSettings.GetAllowCreateUpdateChannels())
		} else {
			tfStateTeamGuestSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		if responseTeamGuestSettings.GetAllowDeleteChannels() != nil {
			tfStateTeamGuestSettings.AllowDeleteChannels = types.BoolValue(*responseTeamGuestSettings.GetAllowDeleteChannels())
		} else {
			tfStateTeamGuestSettings.AllowDeleteChannels = types.BoolNull()
		}

		tfStateTeam.GuestSettings, _ = types.ObjectValueFrom(ctx, tfStateTeamGuestSettings.AttributeTypes(), tfStateTeamGuestSettings)
	}
	if responseTeam.GetId() != nil {
		tfStateTeam.Id = types.StringValue(*responseTeam.GetId())
	} else {
		tfStateTeam.Id = types.StringNull()
	}
	if responseTeam.GetInternalId() != nil {
		tfStateTeam.InternalId = types.StringValue(*responseTeam.GetInternalId())
	} else {
		tfStateTeam.InternalId = types.StringNull()
	}
	if responseTeam.GetIsArchived() != nil {
		tfStateTeam.IsArchived = types.BoolValue(*responseTeam.GetIsArchived())
	} else {
		tfStateTeam.IsArchived = types.BoolNull()
	}
	if responseTeam.GetMemberSettings() != nil {
		tfStateTeamMemberSettings := teamTeamMemberSettingsModel{}
		responseTeamMemberSettings := responseTeam.GetMemberSettings()

		if responseTeamMemberSettings.GetAllowAddRemoveApps() != nil {
			tfStateTeamMemberSettings.AllowAddRemoveApps = types.BoolValue(*responseTeamMemberSettings.GetAllowAddRemoveApps())
		} else {
			tfStateTeamMemberSettings.AllowAddRemoveApps = types.BoolNull()
		}
		if responseTeamMemberSettings.GetAllowCreatePrivateChannels() != nil {
			tfStateTeamMemberSettings.AllowCreatePrivateChannels = types.BoolValue(*responseTeamMemberSettings.GetAllowCreatePrivateChannels())
		} else {
			tfStateTeamMemberSettings.AllowCreatePrivateChannels = types.BoolNull()
		}
		if responseTeamMemberSettings.GetAllowCreateUpdateChannels() != nil {
			tfStateTeamMemberSettings.AllowCreateUpdateChannels = types.BoolValue(*responseTeamMemberSettings.GetAllowCreateUpdateChannels())
		} else {
			tfStateTeamMemberSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		if responseTeamMemberSettings.GetAllowCreateUpdateRemoveConnectors() != nil {
			tfStateTeamMemberSettings.AllowCreateUpdateRemoveConnectors = types.BoolValue(*responseTeamMemberSettings.GetAllowCreateUpdateRemoveConnectors())
		} else {
			tfStateTeamMemberSettings.AllowCreateUpdateRemoveConnectors = types.BoolNull()
		}
		if responseTeamMemberSettings.GetAllowCreateUpdateRemoveTabs() != nil {
			tfStateTeamMemberSettings.AllowCreateUpdateRemoveTabs = types.BoolValue(*responseTeamMemberSettings.GetAllowCreateUpdateRemoveTabs())
		} else {
			tfStateTeamMemberSettings.AllowCreateUpdateRemoveTabs = types.BoolNull()
		}
		if responseTeamMemberSettings.GetAllowDeleteChannels() != nil {
			tfStateTeamMemberSettings.AllowDeleteChannels = types.BoolValue(*responseTeamMemberSettings.GetAllowDeleteChannels())
		} else {
			tfStateTeamMemberSettings.AllowDeleteChannels = types.BoolNull()
		}

		tfStateTeam.MemberSettings, _ = types.ObjectValueFrom(ctx, tfStateTeamMemberSettings.AttributeTypes(), tfStateTeamMemberSettings)
	}
	if responseTeam.GetMessagingSettings() != nil {
		tfStateTeamMessagingSettings := teamTeamMessagingSettingsModel{}
		responseTeamMessagingSettings := responseTeam.GetMessagingSettings()

		if responseTeamMessagingSettings.GetAllowChannelMentions() != nil {
			tfStateTeamMessagingSettings.AllowChannelMentions = types.BoolValue(*responseTeamMessagingSettings.GetAllowChannelMentions())
		} else {
			tfStateTeamMessagingSettings.AllowChannelMentions = types.BoolNull()
		}
		if responseTeamMessagingSettings.GetAllowOwnerDeleteMessages() != nil {
			tfStateTeamMessagingSettings.AllowOwnerDeleteMessages = types.BoolValue(*responseTeamMessagingSettings.GetAllowOwnerDeleteMessages())
		} else {
			tfStateTeamMessagingSettings.AllowOwnerDeleteMessages = types.BoolNull()
		}
		if responseTeamMessagingSettings.GetAllowTeamMentions() != nil {
			tfStateTeamMessagingSettings.AllowTeamMentions = types.BoolValue(*responseTeamMessagingSettings.GetAllowTeamMentions())
		} else {
			tfStateTeamMessagingSettings.AllowTeamMentions = types.BoolNull()
		}
		if responseTeamMessagingSettings.GetAllowUserDeleteMessages() != nil {
			tfStateTeamMessagingSettings.AllowUserDeleteMessages = types.BoolValue(*responseTeamMessagingSettings.GetAllowUserDeleteMessages())
		} else {
			tfStateTeamMessagingSettings.AllowUserDeleteMessages = types.BoolNull()
		}
		if responseTeamMessagingSettings.GetAllowUserEditMessages() != nil {
			tfStateTeamMessagingSettings.AllowUserEditMessages = types.BoolValue(*responseTeamMessagingSettings.GetAllowUserEditMessages())
		} else {
			tfStateTeamMessagingSettings.AllowUserEditMessages = types.BoolNull()
		}

		tfStateTeam.MessagingSettings, _ = types.ObjectValueFrom(ctx, tfStateTeamMessagingSettings.AttributeTypes(), tfStateTeamMessagingSettings)
	}
	if responseTeam.GetSpecialization() != nil {
		tfStateTeam.Specialization = types.StringValue(responseTeam.GetSpecialization().String())
	} else {
		tfStateTeam.Specialization = types.StringNull()
	}
	if responseTeam.GetTenantId() != nil {
		tfStateTeam.TenantId = types.StringValue(*responseTeam.GetTenantId())
	} else {
		tfStateTeam.TenantId = types.StringNull()
	}
	if responseTeam.GetVisibility() != nil {
		tfStateTeam.Visibility = types.StringValue(responseTeam.GetVisibility().String())
	} else {
		tfStateTeam.Visibility = types.StringNull()
	}
	if responseTeam.GetWebUrl() != nil {
		tfStateTeam.WebUrl = types.StringValue(*responseTeam.GetWebUrl())
	} else {
		tfStateTeam.WebUrl = types.StringNull()
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlanTeam teamModel
	diags := req.Plan.Get(ctx, &tfPlanTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfStateTeam teamModel
	diags = req.State.Get(ctx, &tfStateTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBodyTeam := models.NewTeam()

	if !tfPlanTeam.Classification.Equal(tfStateTeam.Classification) {
		tfPlanClassification := tfPlanTeam.Classification.ValueString()
		requestBodyTeam.SetClassification(&tfPlanClassification)
	}

	if !tfPlanTeam.CreatedDateTime.Equal(tfStateTeam.CreatedDateTime) {
		tfPlanCreatedDateTime := tfPlanTeam.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyTeam.SetCreatedDateTime(&t)
	}

	if !tfPlanTeam.Description.Equal(tfStateTeam.Description) {
		tfPlanDescription := tfPlanTeam.Description.ValueString()
		requestBodyTeam.SetDescription(&tfPlanDescription)
	}

	if !tfPlanTeam.DisplayName.Equal(tfStateTeam.DisplayName) {
		tfPlanDisplayName := tfPlanTeam.DisplayName.ValueString()
		requestBodyTeam.SetDisplayName(&tfPlanDisplayName)
	}

	if !tfPlanTeam.FunSettings.Equal(tfStateTeam.FunSettings) {
		requestBodyTeamFunSettings := models.NewTeamFunSettings()
		tfPlanTeamFunSettings := teamTeamFunSettingsModel{}
		tfPlanTeam.FunSettings.As(ctx, &tfPlanTeamFunSettings, basetypes.ObjectAsOptions{})
		tfStateTeamFunSettings := teamTeamFunSettingsModel{}
		tfStateTeam.FunSettings.As(ctx, &tfStateTeamFunSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamFunSettings.AllowCustomMemes.Equal(tfStateTeamFunSettings.AllowCustomMemes) {
			tfPlanAllowCustomMemes := tfPlanTeamFunSettings.AllowCustomMemes.ValueBool()
			requestBodyTeamFunSettings.SetAllowCustomMemes(&tfPlanAllowCustomMemes)
		}

		if !tfPlanTeamFunSettings.AllowGiphy.Equal(tfStateTeamFunSettings.AllowGiphy) {
			tfPlanAllowGiphy := tfPlanTeamFunSettings.AllowGiphy.ValueBool()
			requestBodyTeamFunSettings.SetAllowGiphy(&tfPlanAllowGiphy)
		}

		if !tfPlanTeamFunSettings.AllowStickersAndMemes.Equal(tfStateTeamFunSettings.AllowStickersAndMemes) {
			tfPlanAllowStickersAndMemes := tfPlanTeamFunSettings.AllowStickersAndMemes.ValueBool()
			requestBodyTeamFunSettings.SetAllowStickersAndMemes(&tfPlanAllowStickersAndMemes)
		}

		if !tfPlanTeamFunSettings.GiphyContentRating.Equal(tfStateTeamFunSettings.GiphyContentRating) {
			tfPlanGiphyContentRating := tfPlanTeamFunSettings.GiphyContentRating.ValueString()
			parsedGiphyContentRating, _ := models.ParseGiphyRatingType(tfPlanGiphyContentRating)
			assertedGiphyContentRating := parsedGiphyContentRating.(models.GiphyRatingType)
			requestBodyTeamFunSettings.SetGiphyContentRating(&assertedGiphyContentRating)
		}
		requestBodyTeam.SetFunSettings(requestBodyTeamFunSettings)
		tfPlanTeam.FunSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamFunSettings.AttributeTypes(), tfPlanTeamFunSettings)
	}

	if !tfPlanTeam.GuestSettings.Equal(tfStateTeam.GuestSettings) {
		requestBodyTeamGuestSettings := models.NewTeamGuestSettings()
		tfPlanTeamGuestSettings := teamTeamGuestSettingsModel{}
		tfPlanTeam.GuestSettings.As(ctx, &tfPlanTeamGuestSettings, basetypes.ObjectAsOptions{})
		tfStateTeamGuestSettings := teamTeamGuestSettingsModel{}
		tfStateTeam.GuestSettings.As(ctx, &tfStateTeamGuestSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamGuestSettings.AllowCreateUpdateChannels.Equal(tfStateTeamGuestSettings.AllowCreateUpdateChannels) {
			tfPlanAllowCreateUpdateChannels := tfPlanTeamGuestSettings.AllowCreateUpdateChannels.ValueBool()
			requestBodyTeamGuestSettings.SetAllowCreateUpdateChannels(&tfPlanAllowCreateUpdateChannels)
		}

		if !tfPlanTeamGuestSettings.AllowDeleteChannels.Equal(tfStateTeamGuestSettings.AllowDeleteChannels) {
			tfPlanAllowDeleteChannels := tfPlanTeamGuestSettings.AllowDeleteChannels.ValueBool()
			requestBodyTeamGuestSettings.SetAllowDeleteChannels(&tfPlanAllowDeleteChannels)
		}
		requestBodyTeam.SetGuestSettings(requestBodyTeamGuestSettings)
		tfPlanTeam.GuestSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamGuestSettings.AttributeTypes(), tfPlanTeamGuestSettings)
	}

	if !tfPlanTeam.Id.Equal(tfStateTeam.Id) {
		tfPlanId := tfPlanTeam.Id.ValueString()
		requestBodyTeam.SetId(&tfPlanId)
	}

	if !tfPlanTeam.InternalId.Equal(tfStateTeam.InternalId) {
		tfPlanInternalId := tfPlanTeam.InternalId.ValueString()
		requestBodyTeam.SetInternalId(&tfPlanInternalId)
	}

	if !tfPlanTeam.IsArchived.Equal(tfStateTeam.IsArchived) {
		tfPlanIsArchived := tfPlanTeam.IsArchived.ValueBool()
		requestBodyTeam.SetIsArchived(&tfPlanIsArchived)
	}

	if !tfPlanTeam.MemberSettings.Equal(tfStateTeam.MemberSettings) {
		requestBodyTeamMemberSettings := models.NewTeamMemberSettings()
		tfPlanTeamMemberSettings := teamTeamMemberSettingsModel{}
		tfPlanTeam.MemberSettings.As(ctx, &tfPlanTeamMemberSettings, basetypes.ObjectAsOptions{})
		tfStateTeamMemberSettings := teamTeamMemberSettingsModel{}
		tfStateTeam.MemberSettings.As(ctx, &tfStateTeamMemberSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamMemberSettings.AllowAddRemoveApps.Equal(tfStateTeamMemberSettings.AllowAddRemoveApps) {
			tfPlanAllowAddRemoveApps := tfPlanTeamMemberSettings.AllowAddRemoveApps.ValueBool()
			requestBodyTeamMemberSettings.SetAllowAddRemoveApps(&tfPlanAllowAddRemoveApps)
		}

		if !tfPlanTeamMemberSettings.AllowCreatePrivateChannels.Equal(tfStateTeamMemberSettings.AllowCreatePrivateChannels) {
			tfPlanAllowCreatePrivateChannels := tfPlanTeamMemberSettings.AllowCreatePrivateChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreatePrivateChannels(&tfPlanAllowCreatePrivateChannels)
		}

		if !tfPlanTeamMemberSettings.AllowCreateUpdateChannels.Equal(tfStateTeamMemberSettings.AllowCreateUpdateChannels) {
			tfPlanAllowCreateUpdateChannels := tfPlanTeamMemberSettings.AllowCreateUpdateChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateChannels(&tfPlanAllowCreateUpdateChannels)
		}

		if !tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors.Equal(tfStateTeamMemberSettings.AllowCreateUpdateRemoveConnectors) {
			tfPlanAllowCreateUpdateRemoveConnectors := tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateRemoveConnectors(&tfPlanAllowCreateUpdateRemoveConnectors)
		}

		if !tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs.Equal(tfStateTeamMemberSettings.AllowCreateUpdateRemoveTabs) {
			tfPlanAllowCreateUpdateRemoveTabs := tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateRemoveTabs(&tfPlanAllowCreateUpdateRemoveTabs)
		}

		if !tfPlanTeamMemberSettings.AllowDeleteChannels.Equal(tfStateTeamMemberSettings.AllowDeleteChannels) {
			tfPlanAllowDeleteChannels := tfPlanTeamMemberSettings.AllowDeleteChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowDeleteChannels(&tfPlanAllowDeleteChannels)
		}
		requestBodyTeam.SetMemberSettings(requestBodyTeamMemberSettings)
		tfPlanTeam.MemberSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamMemberSettings.AttributeTypes(), tfPlanTeamMemberSettings)
	}

	if !tfPlanTeam.MessagingSettings.Equal(tfStateTeam.MessagingSettings) {
		requestBodyTeamMessagingSettings := models.NewTeamMessagingSettings()
		tfPlanTeamMessagingSettings := teamTeamMessagingSettingsModel{}
		tfPlanTeam.MessagingSettings.As(ctx, &tfPlanTeamMessagingSettings, basetypes.ObjectAsOptions{})
		tfStateTeamMessagingSettings := teamTeamMessagingSettingsModel{}
		tfStateTeam.MessagingSettings.As(ctx, &tfStateTeamMessagingSettings, basetypes.ObjectAsOptions{})

		if !tfPlanTeamMessagingSettings.AllowChannelMentions.Equal(tfStateTeamMessagingSettings.AllowChannelMentions) {
			tfPlanAllowChannelMentions := tfPlanTeamMessagingSettings.AllowChannelMentions.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowChannelMentions(&tfPlanAllowChannelMentions)
		}

		if !tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages.Equal(tfStateTeamMessagingSettings.AllowOwnerDeleteMessages) {
			tfPlanAllowOwnerDeleteMessages := tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowOwnerDeleteMessages(&tfPlanAllowOwnerDeleteMessages)
		}

		if !tfPlanTeamMessagingSettings.AllowTeamMentions.Equal(tfStateTeamMessagingSettings.AllowTeamMentions) {
			tfPlanAllowTeamMentions := tfPlanTeamMessagingSettings.AllowTeamMentions.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowTeamMentions(&tfPlanAllowTeamMentions)
		}

		if !tfPlanTeamMessagingSettings.AllowUserDeleteMessages.Equal(tfStateTeamMessagingSettings.AllowUserDeleteMessages) {
			tfPlanAllowUserDeleteMessages := tfPlanTeamMessagingSettings.AllowUserDeleteMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowUserDeleteMessages(&tfPlanAllowUserDeleteMessages)
		}

		if !tfPlanTeamMessagingSettings.AllowUserEditMessages.Equal(tfStateTeamMessagingSettings.AllowUserEditMessages) {
			tfPlanAllowUserEditMessages := tfPlanTeamMessagingSettings.AllowUserEditMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowUserEditMessages(&tfPlanAllowUserEditMessages)
		}
		requestBodyTeam.SetMessagingSettings(requestBodyTeamMessagingSettings)
		tfPlanTeam.MessagingSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamMessagingSettings.AttributeTypes(), tfPlanTeamMessagingSettings)
	}

	if !tfPlanTeam.Specialization.Equal(tfStateTeam.Specialization) {
		tfPlanSpecialization := tfPlanTeam.Specialization.ValueString()
		parsedSpecialization, _ := models.ParseTeamSpecialization(tfPlanSpecialization)
		assertedSpecialization := parsedSpecialization.(models.TeamSpecialization)
		requestBodyTeam.SetSpecialization(&assertedSpecialization)
	}

	if !tfPlanTeam.TenantId.Equal(tfStateTeam.TenantId) {
		tfPlanTenantId := tfPlanTeam.TenantId.ValueString()
		requestBodyTeam.SetTenantId(&tfPlanTenantId)
	}

	if !tfPlanTeam.Visibility.Equal(tfStateTeam.Visibility) {
		tfPlanVisibility := tfPlanTeam.Visibility.ValueString()
		parsedVisibility, _ := models.ParseTeamVisibilityType(tfPlanVisibility)
		assertedVisibility := parsedVisibility.(models.TeamVisibilityType)
		requestBodyTeam.SetVisibility(&assertedVisibility)
	}

	if !tfPlanTeam.WebUrl.Equal(tfStateTeam.WebUrl) {
		tfPlanWebUrl := tfPlanTeam.WebUrl.ValueString()
		requestBodyTeam.SetWebUrl(&tfPlanWebUrl)
	}

	// Update team
	_, err := r.client.Teams().ByTeamId(tfStateTeam.Id.ValueString()).Patch(context.Background(), requestBodyTeam, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating team",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlanTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfStateTeam teamModel
	diags := req.State.Get(ctx, &tfStateTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete team
	err := r.client.Teams().ByTeamId(tfStateTeam.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting team",
			err.Error(),
		)
		return
	}

}
