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
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifiers.UseStateForUnconfigured(),
				},
			},
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
			"summary": schema.SingleNestedAttribute{
				Description: "Contains summary information about the team, including number of owners, members, and guests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifiers.UseStateForUnconfigured(),
				},
				Attributes: map[string]schema.Attribute{},
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
	// Retrieve values from plan
	var plan teamModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Plan
	requestBody := models.NewTeam()

	if !plan.Id.IsUnknown() {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	} else {
		plan.Id = types.StringNull()
	}

	if !plan.Classification.IsUnknown() {
		planClassification := plan.Classification.ValueString()
		requestBody.SetClassification(&planClassification)
	} else {
		plan.Classification = types.StringNull()
	}

	if !plan.CreatedDateTime.IsUnknown() {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	} else {
		plan.CreatedDateTime = types.StringNull()
	}

	if !plan.Description.IsUnknown() {
		planDescription := plan.Description.ValueString()
		requestBody.SetDescription(&planDescription)
	} else {
		plan.Description = types.StringNull()
	}

	if !plan.DisplayName.IsUnknown() {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	} else {
		plan.DisplayName = types.StringNull()
	}

	if !plan.FunSettings.IsUnknown() {
		funSettings := models.NewTeamFunSettings()
		funSettingsModel := teamTeamFunSettingsModel{}
		plan.FunSettings.As(ctx, &funSettingsModel, basetypes.ObjectAsOptions{})

		if !funSettingsModel.AllowCustomMemes.IsUnknown() {
			planAllowCustomMemes := funSettingsModel.AllowCustomMemes.ValueBool()
			funSettings.SetAllowCustomMemes(&planAllowCustomMemes)
		} else {
			funSettingsModel.AllowCustomMemes = types.BoolNull()
		}

		if !funSettingsModel.AllowGiphy.IsUnknown() {
			planAllowGiphy := funSettingsModel.AllowGiphy.ValueBool()
			funSettings.SetAllowGiphy(&planAllowGiphy)
		} else {
			funSettingsModel.AllowGiphy = types.BoolNull()
		}

		if !funSettingsModel.AllowStickersAndMemes.IsUnknown() {
			planAllowStickersAndMemes := funSettingsModel.AllowStickersAndMemes.ValueBool()
			funSettings.SetAllowStickersAndMemes(&planAllowStickersAndMemes)
		} else {
			funSettingsModel.AllowStickersAndMemes = types.BoolNull()
		}

		if !funSettingsModel.GiphyContentRating.IsUnknown() {
			planGiphyContentRating := funSettingsModel.GiphyContentRating.ValueString()
			parsedGiphyContentRating, _ := models.ParseGiphyRatingType(planGiphyContentRating)
			assertedGiphyContentRating := parsedGiphyContentRating.(models.GiphyRatingType)
			funSettings.SetGiphyContentRating(&assertedGiphyContentRating)
		} else {
			funSettingsModel.GiphyContentRating = types.StringNull()
		}
		requestBody.SetFunSettings(funSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, funSettingsModel.AttributeTypes(), funSettingsModel)
		plan.FunSettings = objectValue
	} else {
		plan.FunSettings = types.ObjectNull(plan.FunSettings.AttributeTypes(ctx))
	}

	if !plan.GuestSettings.IsUnknown() {
		guestSettings := models.NewTeamGuestSettings()
		guestSettingsModel := teamTeamGuestSettingsModel{}
		plan.GuestSettings.As(ctx, &guestSettingsModel, basetypes.ObjectAsOptions{})

		if !guestSettingsModel.AllowCreateUpdateChannels.IsUnknown() {
			planAllowCreateUpdateChannels := guestSettingsModel.AllowCreateUpdateChannels.ValueBool()
			guestSettings.SetAllowCreateUpdateChannels(&planAllowCreateUpdateChannels)
		} else {
			guestSettingsModel.AllowCreateUpdateChannels = types.BoolNull()
		}

		if !guestSettingsModel.AllowDeleteChannels.IsUnknown() {
			planAllowDeleteChannels := guestSettingsModel.AllowDeleteChannels.ValueBool()
			guestSettings.SetAllowDeleteChannels(&planAllowDeleteChannels)
		} else {
			guestSettingsModel.AllowDeleteChannels = types.BoolNull()
		}
		requestBody.SetGuestSettings(guestSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, guestSettingsModel.AttributeTypes(), guestSettingsModel)
		plan.GuestSettings = objectValue
	} else {
		plan.GuestSettings = types.ObjectNull(plan.GuestSettings.AttributeTypes(ctx))
	}

	if !plan.InternalId.IsUnknown() {
		planInternalId := plan.InternalId.ValueString()
		requestBody.SetInternalId(&planInternalId)
	} else {
		plan.InternalId = types.StringNull()
	}

	if !plan.IsArchived.IsUnknown() {
		planIsArchived := plan.IsArchived.ValueBool()
		requestBody.SetIsArchived(&planIsArchived)
	} else {
		plan.IsArchived = types.BoolNull()
	}

	if !plan.MemberSettings.IsUnknown() {
		memberSettings := models.NewTeamMemberSettings()
		memberSettingsModel := teamTeamMemberSettingsModel{}
		plan.MemberSettings.As(ctx, &memberSettingsModel, basetypes.ObjectAsOptions{})

		if !memberSettingsModel.AllowAddRemoveApps.IsUnknown() {
			planAllowAddRemoveApps := memberSettingsModel.AllowAddRemoveApps.ValueBool()
			memberSettings.SetAllowAddRemoveApps(&planAllowAddRemoveApps)
		} else {
			memberSettingsModel.AllowAddRemoveApps = types.BoolNull()
		}

		if !memberSettingsModel.AllowCreatePrivateChannels.IsUnknown() {
			planAllowCreatePrivateChannels := memberSettingsModel.AllowCreatePrivateChannels.ValueBool()
			memberSettings.SetAllowCreatePrivateChannels(&planAllowCreatePrivateChannels)
		} else {
			memberSettingsModel.AllowCreatePrivateChannels = types.BoolNull()
		}

		if !memberSettingsModel.AllowCreateUpdateChannels.IsUnknown() {
			planAllowCreateUpdateChannels := memberSettingsModel.AllowCreateUpdateChannels.ValueBool()
			memberSettings.SetAllowCreateUpdateChannels(&planAllowCreateUpdateChannels)
		} else {
			memberSettingsModel.AllowCreateUpdateChannels = types.BoolNull()
		}

		if !memberSettingsModel.AllowCreateUpdateRemoveConnectors.IsUnknown() {
			planAllowCreateUpdateRemoveConnectors := memberSettingsModel.AllowCreateUpdateRemoveConnectors.ValueBool()
			memberSettings.SetAllowCreateUpdateRemoveConnectors(&planAllowCreateUpdateRemoveConnectors)
		} else {
			memberSettingsModel.AllowCreateUpdateRemoveConnectors = types.BoolNull()
		}

		if !memberSettingsModel.AllowCreateUpdateRemoveTabs.IsUnknown() {
			planAllowCreateUpdateRemoveTabs := memberSettingsModel.AllowCreateUpdateRemoveTabs.ValueBool()
			memberSettings.SetAllowCreateUpdateRemoveTabs(&planAllowCreateUpdateRemoveTabs)
		} else {
			memberSettingsModel.AllowCreateUpdateRemoveTabs = types.BoolNull()
		}

		if !memberSettingsModel.AllowDeleteChannels.IsUnknown() {
			planAllowDeleteChannels := memberSettingsModel.AllowDeleteChannels.ValueBool()
			memberSettings.SetAllowDeleteChannels(&planAllowDeleteChannels)
		} else {
			memberSettingsModel.AllowDeleteChannels = types.BoolNull()
		}
		requestBody.SetMemberSettings(memberSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, memberSettingsModel.AttributeTypes(), memberSettingsModel)
		plan.MemberSettings = objectValue
	} else {
		plan.MemberSettings = types.ObjectNull(plan.MemberSettings.AttributeTypes(ctx))
	}

	if !plan.MessagingSettings.IsUnknown() {
		messagingSettings := models.NewTeamMessagingSettings()
		messagingSettingsModel := teamTeamMessagingSettingsModel{}
		plan.MessagingSettings.As(ctx, &messagingSettingsModel, basetypes.ObjectAsOptions{})

		if !messagingSettingsModel.AllowChannelMentions.IsUnknown() {
			planAllowChannelMentions := messagingSettingsModel.AllowChannelMentions.ValueBool()
			messagingSettings.SetAllowChannelMentions(&planAllowChannelMentions)
		} else {
			messagingSettingsModel.AllowChannelMentions = types.BoolNull()
		}

		if !messagingSettingsModel.AllowOwnerDeleteMessages.IsUnknown() {
			planAllowOwnerDeleteMessages := messagingSettingsModel.AllowOwnerDeleteMessages.ValueBool()
			messagingSettings.SetAllowOwnerDeleteMessages(&planAllowOwnerDeleteMessages)
		} else {
			messagingSettingsModel.AllowOwnerDeleteMessages = types.BoolNull()
		}

		if !messagingSettingsModel.AllowTeamMentions.IsUnknown() {
			planAllowTeamMentions := messagingSettingsModel.AllowTeamMentions.ValueBool()
			messagingSettings.SetAllowTeamMentions(&planAllowTeamMentions)
		} else {
			messagingSettingsModel.AllowTeamMentions = types.BoolNull()
		}

		if !messagingSettingsModel.AllowUserDeleteMessages.IsUnknown() {
			planAllowUserDeleteMessages := messagingSettingsModel.AllowUserDeleteMessages.ValueBool()
			messagingSettings.SetAllowUserDeleteMessages(&planAllowUserDeleteMessages)
		} else {
			messagingSettingsModel.AllowUserDeleteMessages = types.BoolNull()
		}

		if !messagingSettingsModel.AllowUserEditMessages.IsUnknown() {
			planAllowUserEditMessages := messagingSettingsModel.AllowUserEditMessages.ValueBool()
			messagingSettings.SetAllowUserEditMessages(&planAllowUserEditMessages)
		} else {
			messagingSettingsModel.AllowUserEditMessages = types.BoolNull()
		}
		requestBody.SetMessagingSettings(messagingSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, messagingSettingsModel.AttributeTypes(), messagingSettingsModel)
		plan.MessagingSettings = objectValue
	} else {
		plan.MessagingSettings = types.ObjectNull(plan.MessagingSettings.AttributeTypes(ctx))
	}

	if !plan.Specialization.IsUnknown() {
		planSpecialization := plan.Specialization.ValueString()
		parsedSpecialization, _ := models.ParseTeamSpecialization(planSpecialization)
		assertedSpecialization := parsedSpecialization.(models.TeamSpecialization)
		requestBody.SetSpecialization(&assertedSpecialization)
	} else {
		plan.Specialization = types.StringNull()
	}

	if !plan.Summary.IsUnknown() {
		summary := models.NewTeamSummary()
		summaryModel := teamTeamSummaryModel{}
		plan.Summary.As(ctx, &summaryModel, basetypes.ObjectAsOptions{})

		requestBody.SetSummary(summary)
		objectValue, _ := types.ObjectValueFrom(ctx, summaryModel.AttributeTypes(), summaryModel)
		plan.Summary = objectValue
	} else {
		plan.Summary = types.ObjectNull(plan.Summary.AttributeTypes(ctx))
	}

	if !plan.TenantId.IsUnknown() {
		planTenantId := plan.TenantId.ValueString()
		requestBody.SetTenantId(&planTenantId)
	} else {
		plan.TenantId = types.StringNull()
	}

	if !plan.Visibility.IsUnknown() {
		planVisibility := plan.Visibility.ValueString()
		parsedVisibility, _ := models.ParseTeamVisibilityType(planVisibility)
		assertedVisibility := parsedVisibility.(models.TeamVisibilityType)
		requestBody.SetVisibility(&assertedVisibility)
	} else {
		plan.Visibility = types.StringNull()
	}

	if !plan.WebUrl.IsUnknown() {
		planWebUrl := plan.WebUrl.ValueString()
		requestBody.SetWebUrl(&planWebUrl)
	} else {
		plan.WebUrl = types.StringNull()
	}

	// Create new team
	result, err := r.client.Teams().Post(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating team",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	plan.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *teamResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state teamModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := teams.TeamItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &teams.TeamItemRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"classification",
				"createdDateTime",
				"description",
				"displayName",
				"funSettings",
				"guestSettings",
				"internalId",
				"isArchived",
				"memberSettings",
				"messagingSettings",
				"specialization",
				"summary",
				"tenantId",
				"visibility",
				"webUrl",
			},
		},
	}

	var result models.Teamable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.Teams().ByTeamId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting team",
			err.Error(),
		)
		return
	}

	if result.GetId() != nil {
		state.Id = types.StringValue(*result.GetId())
	} else {
		state.Id = types.StringNull()
	}
	if result.GetClassification() != nil {
		state.Classification = types.StringValue(*result.GetClassification())
	} else {
		state.Classification = types.StringNull()
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		state.CreatedDateTime = types.StringNull()
	}
	if result.GetDescription() != nil {
		state.Description = types.StringValue(*result.GetDescription())
	} else {
		state.Description = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		state.DisplayName = types.StringNull()
	}
	if result.GetFunSettings() != nil {
		funSettings := new(teamTeamFunSettingsModel)

		if result.GetFunSettings().GetAllowCustomMemes() != nil {
			funSettings.AllowCustomMemes = types.BoolValue(*result.GetFunSettings().GetAllowCustomMemes())
		} else {
			funSettings.AllowCustomMemes = types.BoolNull()
		}
		if result.GetFunSettings().GetAllowGiphy() != nil {
			funSettings.AllowGiphy = types.BoolValue(*result.GetFunSettings().GetAllowGiphy())
		} else {
			funSettings.AllowGiphy = types.BoolNull()
		}
		if result.GetFunSettings().GetAllowStickersAndMemes() != nil {
			funSettings.AllowStickersAndMemes = types.BoolValue(*result.GetFunSettings().GetAllowStickersAndMemes())
		} else {
			funSettings.AllowStickersAndMemes = types.BoolNull()
		}
		if result.GetFunSettings().GetGiphyContentRating() != nil {
			funSettings.GiphyContentRating = types.StringValue(result.GetFunSettings().GetGiphyContentRating().String())
		} else {
			funSettings.GiphyContentRating = types.StringNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, funSettings.AttributeTypes(), funSettings)
		state.FunSettings = objectValue
	}
	if result.GetGuestSettings() != nil {
		guestSettings := new(teamTeamGuestSettingsModel)

		if result.GetGuestSettings().GetAllowCreateUpdateChannels() != nil {
			guestSettings.AllowCreateUpdateChannels = types.BoolValue(*result.GetGuestSettings().GetAllowCreateUpdateChannels())
		} else {
			guestSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		if result.GetGuestSettings().GetAllowDeleteChannels() != nil {
			guestSettings.AllowDeleteChannels = types.BoolValue(*result.GetGuestSettings().GetAllowDeleteChannels())
		} else {
			guestSettings.AllowDeleteChannels = types.BoolNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, guestSettings.AttributeTypes(), guestSettings)
		state.GuestSettings = objectValue
	}
	if result.GetInternalId() != nil {
		state.InternalId = types.StringValue(*result.GetInternalId())
	} else {
		state.InternalId = types.StringNull()
	}
	if result.GetIsArchived() != nil {
		state.IsArchived = types.BoolValue(*result.GetIsArchived())
	} else {
		state.IsArchived = types.BoolNull()
	}
	if result.GetMemberSettings() != nil {
		memberSettings := new(teamTeamMemberSettingsModel)

		if result.GetMemberSettings().GetAllowAddRemoveApps() != nil {
			memberSettings.AllowAddRemoveApps = types.BoolValue(*result.GetMemberSettings().GetAllowAddRemoveApps())
		} else {
			memberSettings.AllowAddRemoveApps = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreatePrivateChannels() != nil {
			memberSettings.AllowCreatePrivateChannels = types.BoolValue(*result.GetMemberSettings().GetAllowCreatePrivateChannels())
		} else {
			memberSettings.AllowCreatePrivateChannels = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreateUpdateChannels() != nil {
			memberSettings.AllowCreateUpdateChannels = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateChannels())
		} else {
			memberSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreateUpdateRemoveConnectors() != nil {
			memberSettings.AllowCreateUpdateRemoveConnectors = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateRemoveConnectors())
		} else {
			memberSettings.AllowCreateUpdateRemoveConnectors = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreateUpdateRemoveTabs() != nil {
			memberSettings.AllowCreateUpdateRemoveTabs = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateRemoveTabs())
		} else {
			memberSettings.AllowCreateUpdateRemoveTabs = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowDeleteChannels() != nil {
			memberSettings.AllowDeleteChannels = types.BoolValue(*result.GetMemberSettings().GetAllowDeleteChannels())
		} else {
			memberSettings.AllowDeleteChannels = types.BoolNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, memberSettings.AttributeTypes(), memberSettings)
		state.MemberSettings = objectValue
	}
	if result.GetMessagingSettings() != nil {
		messagingSettings := new(teamTeamMessagingSettingsModel)

		if result.GetMessagingSettings().GetAllowChannelMentions() != nil {
			messagingSettings.AllowChannelMentions = types.BoolValue(*result.GetMessagingSettings().GetAllowChannelMentions())
		} else {
			messagingSettings.AllowChannelMentions = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowOwnerDeleteMessages() != nil {
			messagingSettings.AllowOwnerDeleteMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowOwnerDeleteMessages())
		} else {
			messagingSettings.AllowOwnerDeleteMessages = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowTeamMentions() != nil {
			messagingSettings.AllowTeamMentions = types.BoolValue(*result.GetMessagingSettings().GetAllowTeamMentions())
		} else {
			messagingSettings.AllowTeamMentions = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowUserDeleteMessages() != nil {
			messagingSettings.AllowUserDeleteMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowUserDeleteMessages())
		} else {
			messagingSettings.AllowUserDeleteMessages = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowUserEditMessages() != nil {
			messagingSettings.AllowUserEditMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowUserEditMessages())
		} else {
			messagingSettings.AllowUserEditMessages = types.BoolNull()
		}

		objectValue, _ := types.ObjectValueFrom(ctx, messagingSettings.AttributeTypes(), messagingSettings)
		state.MessagingSettings = objectValue
	}
	if result.GetSpecialization() != nil {
		state.Specialization = types.StringValue(result.GetSpecialization().String())
	} else {
		state.Specialization = types.StringNull()
	}
	if result.GetSummary() != nil {
		summary := new(teamTeamSummaryModel)

		objectValue, _ := types.ObjectValueFrom(ctx, summary.AttributeTypes(), summary)
		state.Summary = objectValue
	}
	if result.GetTenantId() != nil {
		state.TenantId = types.StringValue(*result.GetTenantId())
	} else {
		state.TenantId = types.StringNull()
	}
	if result.GetVisibility() != nil {
		state.Visibility = types.StringValue(result.GetVisibility().String())
	} else {
		state.Visibility = types.StringNull()
	}
	if result.GetWebUrl() != nil {
		state.WebUrl = types.StringValue(*result.GetWebUrl())
	} else {
		state.WebUrl = types.StringNull()
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *teamResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan teamModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state teamModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody := models.NewTeam()

	if !plan.Id.Equal(state.Id) {
		planId := plan.Id.ValueString()
		requestBody.SetId(&planId)
	}

	if !plan.Classification.Equal(state.Classification) {
		planClassification := plan.Classification.ValueString()
		requestBody.SetClassification(&planClassification)
	}

	if !plan.CreatedDateTime.Equal(state.CreatedDateTime) {
		planCreatedDateTime := plan.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, planCreatedDateTime)
		requestBody.SetCreatedDateTime(&t)
	}

	if !plan.Description.Equal(state.Description) {
		planDescription := plan.Description.ValueString()
		requestBody.SetDescription(&planDescription)
	}

	if !plan.DisplayName.Equal(state.DisplayName) {
		planDisplayName := plan.DisplayName.ValueString()
		requestBody.SetDisplayName(&planDisplayName)
	}

	if !plan.FunSettings.Equal(state.FunSettings) {
		funSettings := models.NewTeamFunSettings()
		funSettingsModel := teamTeamFunSettingsModel{}
		plan.FunSettings.As(ctx, &funSettingsModel, basetypes.ObjectAsOptions{})
		funSettingsState := teamTeamFunSettingsModel{}
		state.FunSettings.As(ctx, &funSettingsState, basetypes.ObjectAsOptions{})

		if !funSettingsModel.AllowCustomMemes.Equal(funSettingsState.AllowCustomMemes) {
			planAllowCustomMemes := funSettingsModel.AllowCustomMemes.ValueBool()
			funSettings.SetAllowCustomMemes(&planAllowCustomMemes)
		}

		if !funSettingsModel.AllowGiphy.Equal(funSettingsState.AllowGiphy) {
			planAllowGiphy := funSettingsModel.AllowGiphy.ValueBool()
			funSettings.SetAllowGiphy(&planAllowGiphy)
		}

		if !funSettingsModel.AllowStickersAndMemes.Equal(funSettingsState.AllowStickersAndMemes) {
			planAllowStickersAndMemes := funSettingsModel.AllowStickersAndMemes.ValueBool()
			funSettings.SetAllowStickersAndMemes(&planAllowStickersAndMemes)
		}

		if !funSettingsModel.GiphyContentRating.Equal(funSettingsState.GiphyContentRating) {
			planGiphyContentRating := funSettingsModel.GiphyContentRating.ValueString()
			parsedGiphyContentRating, _ := models.ParseGiphyRatingType(planGiphyContentRating)
			assertedGiphyContentRating := parsedGiphyContentRating.(models.GiphyRatingType)
			funSettings.SetGiphyContentRating(&assertedGiphyContentRating)
		}
		requestBody.SetFunSettings(funSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, funSettingsModel.AttributeTypes(), funSettingsModel)
		plan.FunSettings = objectValue
	}

	if !plan.GuestSettings.Equal(state.GuestSettings) {
		guestSettings := models.NewTeamGuestSettings()
		guestSettingsModel := teamTeamGuestSettingsModel{}
		plan.GuestSettings.As(ctx, &guestSettingsModel, basetypes.ObjectAsOptions{})
		guestSettingsState := teamTeamGuestSettingsModel{}
		state.GuestSettings.As(ctx, &guestSettingsState, basetypes.ObjectAsOptions{})

		if !guestSettingsModel.AllowCreateUpdateChannels.Equal(guestSettingsState.AllowCreateUpdateChannels) {
			planAllowCreateUpdateChannels := guestSettingsModel.AllowCreateUpdateChannels.ValueBool()
			guestSettings.SetAllowCreateUpdateChannels(&planAllowCreateUpdateChannels)
		}

		if !guestSettingsModel.AllowDeleteChannels.Equal(guestSettingsState.AllowDeleteChannels) {
			planAllowDeleteChannels := guestSettingsModel.AllowDeleteChannels.ValueBool()
			guestSettings.SetAllowDeleteChannels(&planAllowDeleteChannels)
		}
		requestBody.SetGuestSettings(guestSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, guestSettingsModel.AttributeTypes(), guestSettingsModel)
		plan.GuestSettings = objectValue
	}

	if !plan.InternalId.Equal(state.InternalId) {
		planInternalId := plan.InternalId.ValueString()
		requestBody.SetInternalId(&planInternalId)
	}

	if !plan.IsArchived.Equal(state.IsArchived) {
		planIsArchived := plan.IsArchived.ValueBool()
		requestBody.SetIsArchived(&planIsArchived)
	}

	if !plan.MemberSettings.Equal(state.MemberSettings) {
		memberSettings := models.NewTeamMemberSettings()
		memberSettingsModel := teamTeamMemberSettingsModel{}
		plan.MemberSettings.As(ctx, &memberSettingsModel, basetypes.ObjectAsOptions{})
		memberSettingsState := teamTeamMemberSettingsModel{}
		state.MemberSettings.As(ctx, &memberSettingsState, basetypes.ObjectAsOptions{})

		if !memberSettingsModel.AllowAddRemoveApps.Equal(memberSettingsState.AllowAddRemoveApps) {
			planAllowAddRemoveApps := memberSettingsModel.AllowAddRemoveApps.ValueBool()
			memberSettings.SetAllowAddRemoveApps(&planAllowAddRemoveApps)
		}

		if !memberSettingsModel.AllowCreatePrivateChannels.Equal(memberSettingsState.AllowCreatePrivateChannels) {
			planAllowCreatePrivateChannels := memberSettingsModel.AllowCreatePrivateChannels.ValueBool()
			memberSettings.SetAllowCreatePrivateChannels(&planAllowCreatePrivateChannels)
		}

		if !memberSettingsModel.AllowCreateUpdateChannels.Equal(memberSettingsState.AllowCreateUpdateChannels) {
			planAllowCreateUpdateChannels := memberSettingsModel.AllowCreateUpdateChannels.ValueBool()
			memberSettings.SetAllowCreateUpdateChannels(&planAllowCreateUpdateChannels)
		}

		if !memberSettingsModel.AllowCreateUpdateRemoveConnectors.Equal(memberSettingsState.AllowCreateUpdateRemoveConnectors) {
			planAllowCreateUpdateRemoveConnectors := memberSettingsModel.AllowCreateUpdateRemoveConnectors.ValueBool()
			memberSettings.SetAllowCreateUpdateRemoveConnectors(&planAllowCreateUpdateRemoveConnectors)
		}

		if !memberSettingsModel.AllowCreateUpdateRemoveTabs.Equal(memberSettingsState.AllowCreateUpdateRemoveTabs) {
			planAllowCreateUpdateRemoveTabs := memberSettingsModel.AllowCreateUpdateRemoveTabs.ValueBool()
			memberSettings.SetAllowCreateUpdateRemoveTabs(&planAllowCreateUpdateRemoveTabs)
		}

		if !memberSettingsModel.AllowDeleteChannels.Equal(memberSettingsState.AllowDeleteChannels) {
			planAllowDeleteChannels := memberSettingsModel.AllowDeleteChannels.ValueBool()
			memberSettings.SetAllowDeleteChannels(&planAllowDeleteChannels)
		}
		requestBody.SetMemberSettings(memberSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, memberSettingsModel.AttributeTypes(), memberSettingsModel)
		plan.MemberSettings = objectValue
	}

	if !plan.MessagingSettings.Equal(state.MessagingSettings) {
		messagingSettings := models.NewTeamMessagingSettings()
		messagingSettingsModel := teamTeamMessagingSettingsModel{}
		plan.MessagingSettings.As(ctx, &messagingSettingsModel, basetypes.ObjectAsOptions{})
		messagingSettingsState := teamTeamMessagingSettingsModel{}
		state.MessagingSettings.As(ctx, &messagingSettingsState, basetypes.ObjectAsOptions{})

		if !messagingSettingsModel.AllowChannelMentions.Equal(messagingSettingsState.AllowChannelMentions) {
			planAllowChannelMentions := messagingSettingsModel.AllowChannelMentions.ValueBool()
			messagingSettings.SetAllowChannelMentions(&planAllowChannelMentions)
		}

		if !messagingSettingsModel.AllowOwnerDeleteMessages.Equal(messagingSettingsState.AllowOwnerDeleteMessages) {
			planAllowOwnerDeleteMessages := messagingSettingsModel.AllowOwnerDeleteMessages.ValueBool()
			messagingSettings.SetAllowOwnerDeleteMessages(&planAllowOwnerDeleteMessages)
		}

		if !messagingSettingsModel.AllowTeamMentions.Equal(messagingSettingsState.AllowTeamMentions) {
			planAllowTeamMentions := messagingSettingsModel.AllowTeamMentions.ValueBool()
			messagingSettings.SetAllowTeamMentions(&planAllowTeamMentions)
		}

		if !messagingSettingsModel.AllowUserDeleteMessages.Equal(messagingSettingsState.AllowUserDeleteMessages) {
			planAllowUserDeleteMessages := messagingSettingsModel.AllowUserDeleteMessages.ValueBool()
			messagingSettings.SetAllowUserDeleteMessages(&planAllowUserDeleteMessages)
		}

		if !messagingSettingsModel.AllowUserEditMessages.Equal(messagingSettingsState.AllowUserEditMessages) {
			planAllowUserEditMessages := messagingSettingsModel.AllowUserEditMessages.ValueBool()
			messagingSettings.SetAllowUserEditMessages(&planAllowUserEditMessages)
		}
		requestBody.SetMessagingSettings(messagingSettings)
		objectValue, _ := types.ObjectValueFrom(ctx, messagingSettingsModel.AttributeTypes(), messagingSettingsModel)
		plan.MessagingSettings = objectValue
	}

	if !plan.Specialization.Equal(state.Specialization) {
		planSpecialization := plan.Specialization.ValueString()
		parsedSpecialization, _ := models.ParseTeamSpecialization(planSpecialization)
		assertedSpecialization := parsedSpecialization.(models.TeamSpecialization)
		requestBody.SetSpecialization(&assertedSpecialization)
	}

	if !plan.Summary.Equal(state.Summary) {
		summary := models.NewTeamSummary()
		summaryModel := teamTeamSummaryModel{}
		plan.Summary.As(ctx, &summaryModel, basetypes.ObjectAsOptions{})
		summaryState := teamTeamSummaryModel{}
		state.Summary.As(ctx, &summaryState, basetypes.ObjectAsOptions{})

		requestBody.SetSummary(summary)
		objectValue, _ := types.ObjectValueFrom(ctx, summaryModel.AttributeTypes(), summaryModel)
		plan.Summary = objectValue
	}

	if !plan.TenantId.Equal(state.TenantId) {
		planTenantId := plan.TenantId.ValueString()
		requestBody.SetTenantId(&planTenantId)
	}

	if !plan.Visibility.Equal(state.Visibility) {
		planVisibility := plan.Visibility.ValueString()
		parsedVisibility, _ := models.ParseTeamVisibilityType(planVisibility)
		assertedVisibility := parsedVisibility.(models.TeamVisibilityType)
		requestBody.SetVisibility(&assertedVisibility)
	}

	if !plan.WebUrl.Equal(state.WebUrl) {
		planWebUrl := plan.WebUrl.ValueString()
		requestBody.SetWebUrl(&planWebUrl)
	}

	// Update team
	_, err := r.client.Teams().ByTeamId(state.Id.ValueString()).Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating team",
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
func (r *teamResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state teamModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete team
	err := r.client.Teams().ByTeamId(state.Id.ValueString()).Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting team",
			err.Error(),
		)
		return
	}

}
