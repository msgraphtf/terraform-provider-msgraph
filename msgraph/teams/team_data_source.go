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
				Attributes: map[string]schema.Attribute{
					"guests_count": schema.Int64Attribute{
						Description: "Count of guests in a team.",
						Computed:    true,
					},
					"members_count": schema.Int64Attribute{
						Description: "Count of members in a team.",
						Computed:    true,
					},
					"owners_count": schema.Int64Attribute{
						Description: "Count of owners in a team.",
						Computed:    true,
					},
				},
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

type teamDataSourceModel struct {
	Id                types.String                          `tfsdk:"id"`
	Classification    types.String                          `tfsdk:"classification"`
	CreatedDateTime   types.String                          `tfsdk:"created_date_time"`
	Description       types.String                          `tfsdk:"description"`
	DisplayName       types.String                          `tfsdk:"display_name"`
	FunSettings       *teamFunSettingsDataSourceModel       `tfsdk:"fun_settings"`
	GuestSettings     *teamGuestSettingsDataSourceModel     `tfsdk:"guest_settings"`
	InternalId        types.String                          `tfsdk:"internal_id"`
	IsArchived        types.Bool                            `tfsdk:"is_archived"`
	MemberSettings    *teamMemberSettingsDataSourceModel    `tfsdk:"member_settings"`
	MessagingSettings *teamMessagingSettingsDataSourceModel `tfsdk:"messaging_settings"`
	Specialization    types.String                          `tfsdk:"specialization"`
	Summary           *teamSummaryDataSourceModel           `tfsdk:"summary"`
	TenantId          types.String                          `tfsdk:"tenant_id"`
	Visibility        types.String                          `tfsdk:"visibility"`
	WebUrl            types.String                          `tfsdk:"web_url"`
}

type teamFunSettingsDataSourceModel struct {
	AllowCustomMemes      types.Bool   `tfsdk:"allow_custom_memes"`
	AllowGiphy            types.Bool   `tfsdk:"allow_giphy"`
	AllowStickersAndMemes types.Bool   `tfsdk:"allow_stickers_and_memes"`
	GiphyContentRating    types.String `tfsdk:"giphy_content_rating"`
}

type teamGuestSettingsDataSourceModel struct {
	AllowCreateUpdateChannels types.Bool `tfsdk:"allow_create_update_channels"`
	AllowDeleteChannels       types.Bool `tfsdk:"allow_delete_channels"`
}

type teamMemberSettingsDataSourceModel struct {
	AllowAddRemoveApps                types.Bool `tfsdk:"allow_add_remove_apps"`
	AllowCreatePrivateChannels        types.Bool `tfsdk:"allow_create_private_channels"`
	AllowCreateUpdateChannels         types.Bool `tfsdk:"allow_create_update_channels"`
	AllowCreateUpdateRemoveConnectors types.Bool `tfsdk:"allow_create_update_remove_connectors"`
	AllowCreateUpdateRemoveTabs       types.Bool `tfsdk:"allow_create_update_remove_tabs"`
	AllowDeleteChannels               types.Bool `tfsdk:"allow_delete_channels"`
}

type teamMessagingSettingsDataSourceModel struct {
	AllowChannelMentions     types.Bool `tfsdk:"allow_channel_mentions"`
	AllowOwnerDeleteMessages types.Bool `tfsdk:"allow_owner_delete_messages"`
	AllowTeamMentions        types.Bool `tfsdk:"allow_team_mentions"`
	AllowUserDeleteMessages  types.Bool `tfsdk:"allow_user_delete_messages"`
	AllowUserEditMessages    types.Bool `tfsdk:"allow_user_edit_messages"`
}

type teamSummaryDataSourceModel struct {
	GuestsCount  types.Int64 `tfsdk:"guests_count"`
	MembersCount types.Int64 `tfsdk:"members_count"`
	OwnersCount  types.Int64 `tfsdk:"owners_count"`
}

// Read refreshes the Terraform state with the latest data.
func (d *teamDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state teamDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
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
				"allChannels",
				"channels",
				"group",
				"incomingChannels",
				"installedApps",
				"members",
				"operations",
				"permissionGrants",
				"photo",
				"primaryChannel",
				"tags",
				"template",
				"schedule",
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
			"`id` must be supplied.",
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
	}
	if result.GetClassification() != nil {
		state.Classification = types.StringValue(*result.GetClassification())
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	}
	if result.GetDescription() != nil {
		state.Description = types.StringValue(*result.GetDescription())
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	}
	if result.GetFunSettings() != nil {
		state.FunSettings = new(teamFunSettingsDataSourceModel)

		if result.GetFunSettings().GetAllowCustomMemes() != nil {
			state.FunSettings.AllowCustomMemes = types.BoolValue(*result.GetFunSettings().GetAllowCustomMemes())
		}
		if result.GetFunSettings().GetAllowGiphy() != nil {
			state.FunSettings.AllowGiphy = types.BoolValue(*result.GetFunSettings().GetAllowGiphy())
		}
		if result.GetFunSettings().GetAllowStickersAndMemes() != nil {
			state.FunSettings.AllowStickersAndMemes = types.BoolValue(*result.GetFunSettings().GetAllowStickersAndMemes())
		}
		if result.GetFunSettings().GetGiphyContentRating() != nil {
			state.FunSettings.GiphyContentRating = types.StringValue(result.GetFunSettings().GetGiphyContentRating().String())
		}
	}
	if result.GetGuestSettings() != nil {
		state.GuestSettings = new(teamGuestSettingsDataSourceModel)

		if result.GetGuestSettings().GetAllowCreateUpdateChannels() != nil {
			state.GuestSettings.AllowCreateUpdateChannels = types.BoolValue(*result.GetGuestSettings().GetAllowCreateUpdateChannels())
		}
		if result.GetGuestSettings().GetAllowDeleteChannels() != nil {
			state.GuestSettings.AllowDeleteChannels = types.BoolValue(*result.GetGuestSettings().GetAllowDeleteChannels())
		}
	}
	if result.GetInternalId() != nil {
		state.InternalId = types.StringValue(*result.GetInternalId())
	}
	if result.GetIsArchived() != nil {
		state.IsArchived = types.BoolValue(*result.GetIsArchived())
	}
	if result.GetMemberSettings() != nil {
		state.MemberSettings = new(teamMemberSettingsDataSourceModel)

		if result.GetMemberSettings().GetAllowAddRemoveApps() != nil {
			state.MemberSettings.AllowAddRemoveApps = types.BoolValue(*result.GetMemberSettings().GetAllowAddRemoveApps())
		}
		if result.GetMemberSettings().GetAllowCreatePrivateChannels() != nil {
			state.MemberSettings.AllowCreatePrivateChannels = types.BoolValue(*result.GetMemberSettings().GetAllowCreatePrivateChannels())
		}
		if result.GetMemberSettings().GetAllowCreateUpdateChannels() != nil {
			state.MemberSettings.AllowCreateUpdateChannels = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateChannels())
		}
		if result.GetMemberSettings().GetAllowCreateUpdateRemoveConnectors() != nil {
			state.MemberSettings.AllowCreateUpdateRemoveConnectors = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateRemoveConnectors())
		}
		if result.GetMemberSettings().GetAllowCreateUpdateRemoveTabs() != nil {
			state.MemberSettings.AllowCreateUpdateRemoveTabs = types.BoolValue(*result.GetMemberSettings().GetAllowCreateUpdateRemoveTabs())
		}
		if result.GetMemberSettings().GetAllowDeleteChannels() != nil {
			state.MemberSettings.AllowDeleteChannels = types.BoolValue(*result.GetMemberSettings().GetAllowDeleteChannels())
		}
	}
	if result.GetMessagingSettings() != nil {
		state.MessagingSettings = new(teamMessagingSettingsDataSourceModel)

		if result.GetMessagingSettings().GetAllowChannelMentions() != nil {
			state.MessagingSettings.AllowChannelMentions = types.BoolValue(*result.GetMessagingSettings().GetAllowChannelMentions())
		}
		if result.GetMessagingSettings().GetAllowOwnerDeleteMessages() != nil {
			state.MessagingSettings.AllowOwnerDeleteMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowOwnerDeleteMessages())
		}
		if result.GetMessagingSettings().GetAllowTeamMentions() != nil {
			state.MessagingSettings.AllowTeamMentions = types.BoolValue(*result.GetMessagingSettings().GetAllowTeamMentions())
		}
		if result.GetMessagingSettings().GetAllowUserDeleteMessages() != nil {
			state.MessagingSettings.AllowUserDeleteMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowUserDeleteMessages())
		}
		if result.GetMessagingSettings().GetAllowUserEditMessages() != nil {
			state.MessagingSettings.AllowUserEditMessages = types.BoolValue(*result.GetMessagingSettings().GetAllowUserEditMessages())
		}
	}
	if result.GetSpecialization() != nil {
		state.Specialization = types.StringValue(result.GetSpecialization().String())
	}
	if result.GetSummary() != nil {
		state.Summary = new(teamSummaryDataSourceModel)

		if result.GetSummary().GetGuestsCount() != nil {
			state.Summary.GuestsCount = types.Int64Value(int64(*result.GetSummary().GetGuestsCount()))
		}
		if result.GetSummary().GetMembersCount() != nil {
			state.Summary.MembersCount = types.Int64Value(int64(*result.GetSummary().GetMembersCount()))
		}
		if result.GetSummary().GetOwnersCount() != nil {
			state.Summary.OwnersCount = types.Int64Value(int64(*result.GetSummary().GetOwnersCount()))
		}
	}
	if result.GetTenantId() != nil {
		state.TenantId = types.StringValue(*result.GetTenantId())
	}
	if result.GetVisibility() != nil {
		state.Visibility = types.StringValue(result.GetVisibility().String())
	}
	if result.GetWebUrl() != nil {
		state.WebUrl = types.StringValue(*result.GetWebUrl())
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
