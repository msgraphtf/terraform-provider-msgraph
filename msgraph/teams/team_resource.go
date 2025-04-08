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
	// Retrieve values from Terraform plan
	var tfPlanTeam teamModel
	diags := req.Plan.Get(ctx, &tfPlanTeam)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	requestBodyTeam := models.NewTeam()
	// START Id | CreateStringAttribute
	if !tfPlanTeam.Id.IsUnknown() {
		tfPlanId := tfPlanTeam.Id.ValueString()
		requestBodyTeam.SetId(&tfPlanId)
	} else {
		tfPlanTeam.Id = types.StringNull()
	}
	// END Id | CreateStringAttribute

	// START Classification | CreateStringAttribute
	if !tfPlanTeam.Classification.IsUnknown() {
		tfPlanClassification := tfPlanTeam.Classification.ValueString()
		requestBodyTeam.SetClassification(&tfPlanClassification)
	} else {
		tfPlanTeam.Classification = types.StringNull()
	}
	// END Classification | CreateStringAttribute

	// START CreatedDateTime | CreateStringTimeAttribute
	if !tfPlanTeam.CreatedDateTime.IsUnknown() {
		tfPlanCreatedDateTime := tfPlanTeam.CreatedDateTime.ValueString()
		t, _ := time.Parse(time.RFC3339, tfPlanCreatedDateTime)
		requestBodyTeam.SetCreatedDateTime(&t)
	} else {
		tfPlanTeam.CreatedDateTime = types.StringNull()
	}
	// END CreatedDateTime | CreateStringTimeAttribute

	// START Description | CreateStringAttribute
	if !tfPlanTeam.Description.IsUnknown() {
		tfPlanDescription := tfPlanTeam.Description.ValueString()
		requestBodyTeam.SetDescription(&tfPlanDescription)
	} else {
		tfPlanTeam.Description = types.StringNull()
	}
	// END Description | CreateStringAttribute

	// START DisplayName | CreateStringAttribute
	if !tfPlanTeam.DisplayName.IsUnknown() {
		tfPlanDisplayName := tfPlanTeam.DisplayName.ValueString()
		requestBodyTeam.SetDisplayName(&tfPlanDisplayName)
	} else {
		tfPlanTeam.DisplayName = types.StringNull()
	}
	// END DisplayName | CreateStringAttribute

	// START FunSettings | CreateObjectAttribute
	if !tfPlanTeam.FunSettings.IsUnknown() {
		requestBodyTeamFunSettings := models.NewTeamFunSettings()
		tfPlanTeamFunSettings := teamTeamFunSettingsModel{}
		tfPlanTeam.FunSettings.As(ctx, &tfPlanTeamFunSettings, basetypes.ObjectAsOptions{})

		// START AllowCustomMemes | CreateBoolAttribute
		if !tfPlanTeamFunSettings.AllowCustomMemes.IsUnknown() {
			tfPlanAllowCustomMemes := tfPlanTeamFunSettings.AllowCustomMemes.ValueBool()
			requestBodyTeamFunSettings.SetAllowCustomMemes(&tfPlanAllowCustomMemes)
		} else {
			tfPlanTeamFunSettings.AllowCustomMemes = types.BoolNull()
		}
		// END AllowCustomMemes | CreateBoolAttribute

		// START AllowGiphy | CreateBoolAttribute
		if !tfPlanTeamFunSettings.AllowGiphy.IsUnknown() {
			tfPlanAllowGiphy := tfPlanTeamFunSettings.AllowGiphy.ValueBool()
			requestBodyTeamFunSettings.SetAllowGiphy(&tfPlanAllowGiphy)
		} else {
			tfPlanTeamFunSettings.AllowGiphy = types.BoolNull()
		}
		// END AllowGiphy | CreateBoolAttribute

		// START AllowStickersAndMemes | CreateBoolAttribute
		if !tfPlanTeamFunSettings.AllowStickersAndMemes.IsUnknown() {
			tfPlanAllowStickersAndMemes := tfPlanTeamFunSettings.AllowStickersAndMemes.ValueBool()
			requestBodyTeamFunSettings.SetAllowStickersAndMemes(&tfPlanAllowStickersAndMemes)
		} else {
			tfPlanTeamFunSettings.AllowStickersAndMemes = types.BoolNull()
		}
		// END AllowStickersAndMemes | CreateBoolAttribute

		// START GiphyContentRating | CreateStringEnumAttribute
		if !tfPlanTeamFunSettings.GiphyContentRating.IsUnknown() {
			tfPlanGiphyContentRating := tfPlanTeamFunSettings.GiphyContentRating.ValueString()
			parsedGiphyContentRating, _ := models.ParseGiphyRatingType(tfPlanGiphyContentRating)
			assertedGiphyContentRating := parsedGiphyContentRating.(models.GiphyRatingType)
			requestBodyTeamFunSettings.SetGiphyContentRating(&assertedGiphyContentRating)
		} else {
			tfPlanTeamFunSettings.GiphyContentRating = types.StringNull()
		}
		// END GiphyContentRating | CreateStringEnumAttribute

		requestBodyTeam.SetFunSettings(requestBodyTeamFunSettings)
		tfPlanTeam.FunSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamFunSettings.AttributeTypes(), requestBodyTeamFunSettings)
	} else {
		tfPlanTeam.FunSettings = types.ObjectNull(tfPlanTeam.FunSettings.AttributeTypes(ctx))
	}
	// END FunSettings | CreateObjectAttribute

	// START GuestSettings | CreateObjectAttribute
	if !tfPlanTeam.GuestSettings.IsUnknown() {
		requestBodyTeamGuestSettings := models.NewTeamGuestSettings()
		tfPlanTeamGuestSettings := teamTeamGuestSettingsModel{}
		tfPlanTeam.GuestSettings.As(ctx, &tfPlanTeamGuestSettings, basetypes.ObjectAsOptions{})

		// START AllowCreateUpdateChannels | CreateBoolAttribute
		if !tfPlanTeamGuestSettings.AllowCreateUpdateChannels.IsUnknown() {
			tfPlanAllowCreateUpdateChannels := tfPlanTeamGuestSettings.AllowCreateUpdateChannels.ValueBool()
			requestBodyTeamGuestSettings.SetAllowCreateUpdateChannels(&tfPlanAllowCreateUpdateChannels)
		} else {
			tfPlanTeamGuestSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		// END AllowCreateUpdateChannels | CreateBoolAttribute

		// START AllowDeleteChannels | CreateBoolAttribute
		if !tfPlanTeamGuestSettings.AllowDeleteChannels.IsUnknown() {
			tfPlanAllowDeleteChannels := tfPlanTeamGuestSettings.AllowDeleteChannels.ValueBool()
			requestBodyTeamGuestSettings.SetAllowDeleteChannels(&tfPlanAllowDeleteChannels)
		} else {
			tfPlanTeamGuestSettings.AllowDeleteChannels = types.BoolNull()
		}
		// END AllowDeleteChannels | CreateBoolAttribute

		requestBodyTeam.SetGuestSettings(requestBodyTeamGuestSettings)
		tfPlanTeam.GuestSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamGuestSettings.AttributeTypes(), requestBodyTeamGuestSettings)
	} else {
		tfPlanTeam.GuestSettings = types.ObjectNull(tfPlanTeam.GuestSettings.AttributeTypes(ctx))
	}
	// END GuestSettings | CreateObjectAttribute

	// START InternalId | CreateStringAttribute
	if !tfPlanTeam.InternalId.IsUnknown() {
		tfPlanInternalId := tfPlanTeam.InternalId.ValueString()
		requestBodyTeam.SetInternalId(&tfPlanInternalId)
	} else {
		tfPlanTeam.InternalId = types.StringNull()
	}
	// END InternalId | CreateStringAttribute

	// START IsArchived | CreateBoolAttribute
	if !tfPlanTeam.IsArchived.IsUnknown() {
		tfPlanIsArchived := tfPlanTeam.IsArchived.ValueBool()
		requestBodyTeam.SetIsArchived(&tfPlanIsArchived)
	} else {
		tfPlanTeam.IsArchived = types.BoolNull()
	}
	// END IsArchived | CreateBoolAttribute

	// START MemberSettings | CreateObjectAttribute
	if !tfPlanTeam.MemberSettings.IsUnknown() {
		requestBodyTeamMemberSettings := models.NewTeamMemberSettings()
		tfPlanTeamMemberSettings := teamTeamMemberSettingsModel{}
		tfPlanTeam.MemberSettings.As(ctx, &tfPlanTeamMemberSettings, basetypes.ObjectAsOptions{})

		// START AllowAddRemoveApps | CreateBoolAttribute
		if !tfPlanTeamMemberSettings.AllowAddRemoveApps.IsUnknown() {
			tfPlanAllowAddRemoveApps := tfPlanTeamMemberSettings.AllowAddRemoveApps.ValueBool()
			requestBodyTeamMemberSettings.SetAllowAddRemoveApps(&tfPlanAllowAddRemoveApps)
		} else {
			tfPlanTeamMemberSettings.AllowAddRemoveApps = types.BoolNull()
		}
		// END AllowAddRemoveApps | CreateBoolAttribute

		// START AllowCreatePrivateChannels | CreateBoolAttribute
		if !tfPlanTeamMemberSettings.AllowCreatePrivateChannels.IsUnknown() {
			tfPlanAllowCreatePrivateChannels := tfPlanTeamMemberSettings.AllowCreatePrivateChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreatePrivateChannels(&tfPlanAllowCreatePrivateChannels)
		} else {
			tfPlanTeamMemberSettings.AllowCreatePrivateChannels = types.BoolNull()
		}
		// END AllowCreatePrivateChannels | CreateBoolAttribute

		// START AllowCreateUpdateChannels | CreateBoolAttribute
		if !tfPlanTeamMemberSettings.AllowCreateUpdateChannels.IsUnknown() {
			tfPlanAllowCreateUpdateChannels := tfPlanTeamMemberSettings.AllowCreateUpdateChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateChannels(&tfPlanAllowCreateUpdateChannels)
		} else {
			tfPlanTeamMemberSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		// END AllowCreateUpdateChannels | CreateBoolAttribute

		// START AllowCreateUpdateRemoveConnectors | CreateBoolAttribute
		if !tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors.IsUnknown() {
			tfPlanAllowCreateUpdateRemoveConnectors := tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateRemoveConnectors(&tfPlanAllowCreateUpdateRemoveConnectors)
		} else {
			tfPlanTeamMemberSettings.AllowCreateUpdateRemoveConnectors = types.BoolNull()
		}
		// END AllowCreateUpdateRemoveConnectors | CreateBoolAttribute

		// START AllowCreateUpdateRemoveTabs | CreateBoolAttribute
		if !tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs.IsUnknown() {
			tfPlanAllowCreateUpdateRemoveTabs := tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs.ValueBool()
			requestBodyTeamMemberSettings.SetAllowCreateUpdateRemoveTabs(&tfPlanAllowCreateUpdateRemoveTabs)
		} else {
			tfPlanTeamMemberSettings.AllowCreateUpdateRemoveTabs = types.BoolNull()
		}
		// END AllowCreateUpdateRemoveTabs | CreateBoolAttribute

		// START AllowDeleteChannels | CreateBoolAttribute
		if !tfPlanTeamMemberSettings.AllowDeleteChannels.IsUnknown() {
			tfPlanAllowDeleteChannels := tfPlanTeamMemberSettings.AllowDeleteChannels.ValueBool()
			requestBodyTeamMemberSettings.SetAllowDeleteChannels(&tfPlanAllowDeleteChannels)
		} else {
			tfPlanTeamMemberSettings.AllowDeleteChannels = types.BoolNull()
		}
		// END AllowDeleteChannels | CreateBoolAttribute

		requestBodyTeam.SetMemberSettings(requestBodyTeamMemberSettings)
		tfPlanTeam.MemberSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamMemberSettings.AttributeTypes(), requestBodyTeamMemberSettings)
	} else {
		tfPlanTeam.MemberSettings = types.ObjectNull(tfPlanTeam.MemberSettings.AttributeTypes(ctx))
	}
	// END MemberSettings | CreateObjectAttribute

	// START MessagingSettings | CreateObjectAttribute
	if !tfPlanTeam.MessagingSettings.IsUnknown() {
		requestBodyTeamMessagingSettings := models.NewTeamMessagingSettings()
		tfPlanTeamMessagingSettings := teamTeamMessagingSettingsModel{}
		tfPlanTeam.MessagingSettings.As(ctx, &tfPlanTeamMessagingSettings, basetypes.ObjectAsOptions{})

		// START AllowChannelMentions | CreateBoolAttribute
		if !tfPlanTeamMessagingSettings.AllowChannelMentions.IsUnknown() {
			tfPlanAllowChannelMentions := tfPlanTeamMessagingSettings.AllowChannelMentions.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowChannelMentions(&tfPlanAllowChannelMentions)
		} else {
			tfPlanTeamMessagingSettings.AllowChannelMentions = types.BoolNull()
		}
		// END AllowChannelMentions | CreateBoolAttribute

		// START AllowOwnerDeleteMessages | CreateBoolAttribute
		if !tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages.IsUnknown() {
			tfPlanAllowOwnerDeleteMessages := tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowOwnerDeleteMessages(&tfPlanAllowOwnerDeleteMessages)
		} else {
			tfPlanTeamMessagingSettings.AllowOwnerDeleteMessages = types.BoolNull()
		}
		// END AllowOwnerDeleteMessages | CreateBoolAttribute

		// START AllowTeamMentions | CreateBoolAttribute
		if !tfPlanTeamMessagingSettings.AllowTeamMentions.IsUnknown() {
			tfPlanAllowTeamMentions := tfPlanTeamMessagingSettings.AllowTeamMentions.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowTeamMentions(&tfPlanAllowTeamMentions)
		} else {
			tfPlanTeamMessagingSettings.AllowTeamMentions = types.BoolNull()
		}
		// END AllowTeamMentions | CreateBoolAttribute

		// START AllowUserDeleteMessages | CreateBoolAttribute
		if !tfPlanTeamMessagingSettings.AllowUserDeleteMessages.IsUnknown() {
			tfPlanAllowUserDeleteMessages := tfPlanTeamMessagingSettings.AllowUserDeleteMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowUserDeleteMessages(&tfPlanAllowUserDeleteMessages)
		} else {
			tfPlanTeamMessagingSettings.AllowUserDeleteMessages = types.BoolNull()
		}
		// END AllowUserDeleteMessages | CreateBoolAttribute

		// START AllowUserEditMessages | CreateBoolAttribute
		if !tfPlanTeamMessagingSettings.AllowUserEditMessages.IsUnknown() {
			tfPlanAllowUserEditMessages := tfPlanTeamMessagingSettings.AllowUserEditMessages.ValueBool()
			requestBodyTeamMessagingSettings.SetAllowUserEditMessages(&tfPlanAllowUserEditMessages)
		} else {
			tfPlanTeamMessagingSettings.AllowUserEditMessages = types.BoolNull()
		}
		// END AllowUserEditMessages | CreateBoolAttribute

		requestBodyTeam.SetMessagingSettings(requestBodyTeamMessagingSettings)
		tfPlanTeam.MessagingSettings, _ = types.ObjectValueFrom(ctx, tfPlanTeamMessagingSettings.AttributeTypes(), requestBodyTeamMessagingSettings)
	} else {
		tfPlanTeam.MessagingSettings = types.ObjectNull(tfPlanTeam.MessagingSettings.AttributeTypes(ctx))
	}
	// END MessagingSettings | CreateObjectAttribute

	// START Specialization | CreateStringEnumAttribute
	if !tfPlanTeam.Specialization.IsUnknown() {
		tfPlanSpecialization := tfPlanTeam.Specialization.ValueString()
		parsedSpecialization, _ := models.ParseTeamSpecialization(tfPlanSpecialization)
		assertedSpecialization := parsedSpecialization.(models.TeamSpecialization)
		requestBodyTeam.SetSpecialization(&assertedSpecialization)
	} else {
		tfPlanTeam.Specialization = types.StringNull()
	}
	// END Specialization | CreateStringEnumAttribute

	// START Summary | CreateObjectAttribute
	if !tfPlanTeam.Summary.IsUnknown() {
		requestBodyTeamSummary := models.NewTeamSummary()
		tfPlanTeamSummary := teamTeamSummaryModel{}
		tfPlanTeam.Summary.As(ctx, &tfPlanTeamSummary, basetypes.ObjectAsOptions{})

		// START GuestsCount | UNKNOWN
		// END GuestsCount | UNKNOWN

		// START MembersCount | UNKNOWN
		// END MembersCount | UNKNOWN

		// START OwnersCount | UNKNOWN
		// END OwnersCount | UNKNOWN

		requestBodyTeam.SetSummary(requestBodyTeamSummary)
		tfPlanTeam.Summary, _ = types.ObjectValueFrom(ctx, tfPlanTeamSummary.AttributeTypes(), requestBodyTeamSummary)
	} else {
		tfPlanTeam.Summary = types.ObjectNull(tfPlanTeam.Summary.AttributeTypes(ctx))
	}
	// END Summary | CreateObjectAttribute

	// START TenantId | CreateStringAttribute
	if !tfPlanTeam.TenantId.IsUnknown() {
		tfPlanTenantId := tfPlanTeam.TenantId.ValueString()
		requestBodyTeam.SetTenantId(&tfPlanTenantId)
	} else {
		tfPlanTeam.TenantId = types.StringNull()
	}
	// END TenantId | CreateStringAttribute

	// START Visibility | CreateStringEnumAttribute
	if !tfPlanTeam.Visibility.IsUnknown() {
		tfPlanVisibility := tfPlanTeam.Visibility.ValueString()
		parsedVisibility, _ := models.ParseTeamVisibilityType(tfPlanVisibility)
		assertedVisibility := parsedVisibility.(models.TeamVisibilityType)
		requestBodyTeam.SetVisibility(&assertedVisibility)
	} else {
		tfPlanTeam.Visibility = types.StringNull()
	}
	// END Visibility | CreateStringEnumAttribute

	// START WebUrl | CreateStringAttribute
	if !tfPlanTeam.WebUrl.IsUnknown() {
		tfPlanWebUrl := tfPlanTeam.WebUrl.ValueString()
		requestBodyTeam.SetWebUrl(&tfPlanWebUrl)
	} else {
		tfPlanTeam.WebUrl = types.StringNull()
	}
	// END WebUrl | CreateStringAttribute

	// Create new team
	result, err := r.client.Teams().Post(context.Background(), requestBodyTeam, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating team",
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

	if !tfStateTeam.Id.IsNull() {
		result, err = d.client.Teams().ByTeamId(tfStateTeam.Id.ValueString()).Get(context.Background(), &qparams)
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
		tfStateTeam.Id = types.StringValue(*result.GetId())
	} else {
		tfStateTeam.Id = types.StringNull()
	}
	if result.GetClassification() != nil {
		tfStateTeam.Classification = types.StringValue(*result.GetClassification())
	} else {
		tfStateTeam.Classification = types.StringNull()
	}
	if result.GetCreatedDateTime() != nil {
		tfStateTeam.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		tfStateTeam.CreatedDateTime = types.StringNull()
	}
	if result.GetDescription() != nil {
		tfStateTeam.Description = types.StringValue(*result.GetDescription())
	} else {
		tfStateTeam.Description = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		tfStateTeam.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		tfStateTeam.DisplayName = types.StringNull()
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

		tfStateTeam.FunSettings, _ = types.ObjectValueFrom(ctx, funSettings.AttributeTypes(), funSettings)
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

		tfStateTeam.GuestSettings, _ = types.ObjectValueFrom(ctx, guestSettings.AttributeTypes(), guestSettings)
	}
	if result.GetInternalId() != nil {
		tfStateTeam.InternalId = types.StringValue(*result.GetInternalId())
	} else {
		tfStateTeam.InternalId = types.StringNull()
	}
	if result.GetIsArchived() != nil {
		tfStateTeam.IsArchived = types.BoolValue(*result.GetIsArchived())
	} else {
		tfStateTeam.IsArchived = types.BoolNull()
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

		tfStateTeam.MemberSettings, _ = types.ObjectValueFrom(ctx, memberSettings.AttributeTypes(), memberSettings)
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

		tfStateTeam.MessagingSettings, _ = types.ObjectValueFrom(ctx, messagingSettings.AttributeTypes(), messagingSettings)
	}
	if result.GetSpecialization() != nil {
		tfStateTeam.Specialization = types.StringValue(result.GetSpecialization().String())
	} else {
		tfStateTeam.Specialization = types.StringNull()
	}
	if result.GetSummary() != nil {
		summary := new(teamTeamSummaryModel)

		tfStateTeam.Summary, _ = types.ObjectValueFrom(ctx, summary.AttributeTypes(), summary)
	}
	if result.GetTenantId() != nil {
		tfStateTeam.TenantId = types.StringValue(*result.GetTenantId())
	} else {
		tfStateTeam.TenantId = types.StringNull()
	}
	if result.GetVisibility() != nil {
		tfStateTeam.Visibility = types.StringValue(result.GetVisibility().String())
	} else {
		tfStateTeam.Visibility = types.StringNull()
	}
	if result.GetWebUrl() != nil {
		tfStateTeam.WebUrl = types.StringValue(*result.GetWebUrl())
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

	if !tfPlanTeam.Id.Equal(tfStateTeam.Id) {
		tfPlanId := tfPlanTeam.Id.ValueString()
		requestBodyTeam.SetId(&tfPlanId)
	}

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

	if !tfPlanTeam.Summary.Equal(tfStateTeam.Summary) {
		requestBodyTeamSummary := models.NewTeamSummary()
		tfPlanTeamSummary := teamTeamSummaryModel{}
		tfPlanTeam.Summary.As(ctx, &tfPlanTeamSummary, basetypes.ObjectAsOptions{})
		tfStateTeamSummary := teamTeamSummaryModel{}
		tfStateTeam.Summary.As(ctx, &tfStateTeamSummary, basetypes.ObjectAsOptions{})

		requestBodyTeam.SetSummary(requestBodyTeamSummary)
		tfPlanTeam.Summary, _ = types.ObjectValueFrom(ctx, tfPlanTeamSummary.AttributeTypes(), tfPlanTeamSummary)
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
