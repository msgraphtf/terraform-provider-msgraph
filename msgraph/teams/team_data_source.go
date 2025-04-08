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
		tfStateFunSettings := teamTeamFunSettingsModel{}

		if result.GetFunSettings().GetAllowCustomMemes() != nil {
			tfStateFunSettings.AllowCustomMemes = types.BoolValue(*result.GetFunSettings().GetAllowCustomMemes())
		} else {
			tfStateFunSettings.AllowCustomMemes = types.BoolNull()
		}
		if result.GetFunSettings().GetAllowGiphy() != nil {
			tfStateFunSettings.AllowGiphy = types.BoolValue(*result.GetFunSettings().GetAllowGiphy())
		} else {
			tfStateFunSettings.AllowGiphy = types.BoolNull()
		}
		if result.GetFunSettings().GetAllowStickersAndMemes() != nil {
			tfStateFunSettings.AllowStickersAndMemes = types.BoolValue(*result.GetFunSettings().GetAllowStickersAndMemes())
		} else {
			tfStateFunSettings.AllowStickersAndMemes = types.BoolNull()
		}
		if result.GetFunSettings().GetGiphyContentRating() != nil {
			tfStateFunSettings.GiphyContentRating = types.StringValue(result.GetFunSettings().GetGiphyContentRating().String())
		} else {
			tfStateFunSettings.GiphyContentRating = types.StringNull()
		}

		tfStateTeam.FunSettings, _ = types.ObjectValueFrom(ctx, tfStateFunSettings.AttributeTypes(), tfStateFunSettings)
	}
	if result.GetGuestSettings() != nil {
		tfStateGuestSettings := teamTeamGuestSettingsModel{}

		if result.GetGuestSettings().GetAllowCreateUpdateChannels() != nil {
			tfStateGuestSettings.AllowCreateUpdateChannels = types.BoolValue(*result.GetGuestSettings().GetAllowCreateUpdateChannels())
		} else {
			tfStateGuestSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		if result.GetGuestSettings().GetAllowDeleteChannels() != nil {
			tfStateGuestSettings.AllowDeleteChannels = types.BoolValue(*result.GetGuestSettings().GetAllowDeleteChannels())
		} else {
			tfStateGuestSettings.AllowDeleteChannels = types.BoolNull()
		}

		tfStateTeam.GuestSettings, _ = types.ObjectValueFrom(ctx, tfStateGuestSettings.AttributeTypes(), tfStateGuestSettings)
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
		tfStateMemberSettings := teamTeamMemberSettingsModel{}

		if result.GetMemberSettings().GetAllowAddRemoveApps() != nil {
			tfStateMemberSettings.AllowAddRemoveApps = types.BoolValue(*result.GetMemberSettings().GetAllowAddRemoveApps())
		} else {
			tfStateMemberSettings.AllowAddRemoveApps = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreatePrivateChannels() != nil {
			tfStateMemberSettings.AllowCreatePrivateChannels = types.BoolValue(*result.GetMemberSettings().GetAllowCreatePrivateChannels())
		} else {
			tfStateMemberSettings.AllowCreatePrivateChannels = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreateUpdateChannels() != nil {
			tfStateMemberSettings.AllowCreateUpdateChannels = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateChannels())
		} else {
			tfStateMemberSettings.AllowCreateUpdateChannels = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreateUpdateRemoveConnectors() != nil {
			tfStateMemberSettings.AllowCreateUpdateRemoveConnectors = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateRemoveConnectors())
		} else {
			tfStateMemberSettings.AllowCreateUpdateRemoveConnectors = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowCreateUpdateRemoveTabs() != nil {
			tfStateMemberSettings.AllowCreateUpdateRemoveTabs = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateRemoveTabs())
		} else {
			tfStateMemberSettings.AllowCreateUpdateRemoveTabs = types.BoolNull()
		}
		if result.GetMemberSettings().GetAllowDeleteChannels() != nil {
			tfStateMemberSettings.AllowDeleteChannels = types.BoolValue(*result.GetMemberSettings().GetAllowDeleteChannels())
		} else {
			tfStateMemberSettings.AllowDeleteChannels = types.BoolNull()
		}

		tfStateTeam.MemberSettings, _ = types.ObjectValueFrom(ctx, tfStateMemberSettings.AttributeTypes(), tfStateMemberSettings)
	}
	if result.GetMessagingSettings() != nil {
		tfStateMessagingSettings := teamTeamMessagingSettingsModel{}

		if result.GetMessagingSettings().GetAllowChannelMentions() != nil {
			tfStateMessagingSettings.AllowChannelMentions = types.BoolValue(*result.GetMessagingSettings().GetAllowChannelMentions())
		} else {
			tfStateMessagingSettings.AllowChannelMentions = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowOwnerDeleteMessages() != nil {
			tfStateMessagingSettings.AllowOwnerDeleteMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowOwnerDeleteMessages())
		} else {
			tfStateMessagingSettings.AllowOwnerDeleteMessages = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowTeamMentions() != nil {
			tfStateMessagingSettings.AllowTeamMentions = types.BoolValue(*result.GetMessagingSettings().GetAllowTeamMentions())
		} else {
			tfStateMessagingSettings.AllowTeamMentions = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowUserDeleteMessages() != nil {
			tfStateMessagingSettings.AllowUserDeleteMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowUserDeleteMessages())
		} else {
			tfStateMessagingSettings.AllowUserDeleteMessages = types.BoolNull()
		}
		if result.GetMessagingSettings().GetAllowUserEditMessages() != nil {
			tfStateMessagingSettings.AllowUserEditMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowUserEditMessages())
		} else {
			tfStateMessagingSettings.AllowUserEditMessages = types.BoolNull()
		}

		tfStateTeam.MessagingSettings, _ = types.ObjectValueFrom(ctx, tfStateMessagingSettings.AttributeTypes(), tfStateMessagingSettings)
	}
	if result.GetSpecialization() != nil {
		tfStateTeam.Specialization = types.StringValue(result.GetSpecialization().String())
	} else {
		tfStateTeam.Specialization = types.StringNull()
	}
	if result.GetSummary() != nil {
		tfStateSummary := teamTeamSummaryModel{}

		tfStateTeam.Summary, _ = types.ObjectValueFrom(ctx, tfStateSummary.AttributeTypes(), tfStateSummary)
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
