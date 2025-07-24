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
		Description: "",
		Attributes: map[string]schema.Attribute{
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
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
				Computed:    true,
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
