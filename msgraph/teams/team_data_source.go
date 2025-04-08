package teams

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/teams"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &teamDataSource{}
	_ datasource.DataSourceWithConfigure = &teamDataSource{}
)

// NewTeamDataSource is a helper function to simplify the provider implementation.
func NewTeamDataSource() datasource.DataSource {
	return &teamDataSource{}
}

// teamDataSource is the data source implementation.
type teamDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *teamDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_team"
}

// Configure adds the provider configured client to the data source.
func (d *teamDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *teamDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
			},
			"classification": schema.StringAttribute{
				Description: "An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.",
				Computed:    true,
			},
			"created_date_time": schema.StringAttribute{
				Description: "Timestamp at which the team was created.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "An optional description for the team. Maximum length: 1024 characters.",
				Computed:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "The name of the team.",
				Computed:    true,
			},
			"fun_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure use of Giphy, memes, and stickers in the team.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allow_custom_memes": schema.BoolAttribute{
						Description: "If set to true, enables users to include custom memes.",
						Computed:    true,
					},
					"allow_giphy": schema.BoolAttribute{
						Description: "If set to true, enables Giphy use.",
						Computed:    true,
					},
					"allow_stickers_and_memes": schema.BoolAttribute{
						Description: "If set to true, enables users to include stickers and memes.",
						Computed:    true,
					},
					"giphy_content_rating": schema.StringAttribute{
						Description: "Giphy content rating. Possible values are: moderate, strict.",
						Computed:    true,
					},
				},
			},
			"guest_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure whether guests can create, update, or delete channels in the team.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allow_create_update_channels": schema.BoolAttribute{
						Description: "If set to true, guests can add and update channels.",
						Computed:    true,
					},
					"allow_delete_channels": schema.BoolAttribute{
						Description: "If set to true, guests can delete channels.",
						Computed:    true,
					},
				},
			},
			"internal_id": schema.StringAttribute{
				Description: "A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.",
				Computed:    true,
			},
			"is_archived": schema.BoolAttribute{
				Description: "Whether this team is in read-only mode.",
				Computed:    true,
			},
			"member_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allow_add_remove_apps": schema.BoolAttribute{
						Description: "If set to true, members can add and remove apps.",
						Computed:    true,
					},
					"allow_create_private_channels": schema.BoolAttribute{
						Description: "If set to true, members can add and update private channels.",
						Computed:    true,
					},
					"allow_create_update_channels": schema.BoolAttribute{
						Description: "If set to true, members can add and update channels.",
						Computed:    true,
					},
					"allow_create_update_remove_connectors": schema.BoolAttribute{
						Description: "If set to true, members can add, update, and remove connectors.",
						Computed:    true,
					},
					"allow_create_update_remove_tabs": schema.BoolAttribute{
						Description: "If set to true, members can add, update, and remove tabs.",
						Computed:    true,
					},
					"allow_delete_channels": schema.BoolAttribute{
						Description: "If set to true, members can delete channels.",
						Computed:    true,
					},
				},
			},
			"messaging_settings": schema.SingleNestedAttribute{
				Description: "Settings to configure messaging and mentions in the team.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"allow_channel_mentions": schema.BoolAttribute{
						Description: "If set to true, @channel mentions are allowed.",
						Computed:    true,
					},
					"allow_owner_delete_messages": schema.BoolAttribute{
						Description: "If set to true, owners can delete any message.",
						Computed:    true,
					},
					"allow_team_mentions": schema.BoolAttribute{
						Description: "If set to true, @team mentions are allowed.",
						Computed:    true,
					},
					"allow_user_delete_messages": schema.BoolAttribute{
						Description: "If set to true, users can delete their messages.",
						Computed:    true,
					},
					"allow_user_edit_messages": schema.BoolAttribute{
						Description: "If set to true, users can edit their messages.",
						Computed:    true,
					},
				},
			},
			"specialization": schema.StringAttribute{
				Description: "Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.",
				Computed:    true,
			},
			"summary": schema.SingleNestedAttribute{
				Description: "Contains summary information about the team, including number of owners, members, and guests.",
				Computed:    true,
				Attributes:  map[string]schema.Attribute{},
			},
			"tenant_id": schema.StringAttribute{
				Description: "The ID of the Microsoft Entra tenant.",
				Computed:    true,
			},
			"visibility": schema.StringAttribute{
				Description: "The visibility of the group and team. Defaults to Public.",
				Computed:    true,
			},
			"web_url": schema.StringAttribute{
				Description: "A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfStateTeam teamModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateTeam)...)
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

		objectValue, _ := types.ObjectValueFrom(ctx, funSettings.AttributeTypes(), funSettings)
		tfStateTeam.FunSettings = objectValue
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
		tfStateTeam.GuestSettings = objectValue
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

		objectValue, _ := types.ObjectValueFrom(ctx, memberSettings.AttributeTypes(), memberSettings)
		tfStateTeam.MemberSettings = objectValue
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
		tfStateTeam.MessagingSettings = objectValue
	}
	if result.GetSpecialization() != nil {
		tfStateTeam.Specialization = types.StringValue(result.GetSpecialization().String())
	} else {
		tfStateTeam.Specialization = types.StringNull()
	}
	if result.GetSummary() != nil {
		summary := new(teamTeamSummaryModel)

		objectValue, _ := types.ObjectValueFrom(ctx, summary.AttributeTypes(), summary)
		tfStateTeam.Summary = objectValue
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
