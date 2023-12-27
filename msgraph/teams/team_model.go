package teams

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type teamModel struct {
	Id                types.String                `tfsdk:"id"`
	Classification    types.String                `tfsdk:"classification"`
	CreatedDateTime   types.String                `tfsdk:"created_date_time"`
	Description       types.String                `tfsdk:"description"`
	DisplayName       types.String                `tfsdk:"display_name"`
	FunSettings       *teamFunSettingsModel       `tfsdk:"fun_settings"`
	GuestSettings     *teamGuestSettingsModel     `tfsdk:"guest_settings"`
	InternalId        types.String                `tfsdk:"internal_id"`
	IsArchived        types.Bool                  `tfsdk:"is_archived"`
	MemberSettings    *teamMemberSettingsModel    `tfsdk:"member_settings"`
	MessagingSettings *teamMessagingSettingsModel `tfsdk:"messaging_settings"`
	Specialization    types.String                `tfsdk:"specialization"`
	Summary           *teamSummaryModel           `tfsdk:"summary"`
	TenantId          types.String                `tfsdk:"tenant_id"`
	Visibility        types.String                `tfsdk:"visibility"`
	WebUrl            types.String                `tfsdk:"web_url"`
}

type teamFunSettingsModel struct {
	AllowCustomMemes      types.Bool   `tfsdk:"allow_custom_memes"`
	AllowGiphy            types.Bool   `tfsdk:"allow_giphy"`
	AllowStickersAndMemes types.Bool   `tfsdk:"allow_stickers_and_memes"`
	GiphyContentRating    types.String `tfsdk:"giphy_content_rating"`
}

type teamGuestSettingsModel struct {
	AllowCreateUpdateChannels types.Bool `tfsdk:"allow_create_update_channels"`
	AllowDeleteChannels       types.Bool `tfsdk:"allow_delete_channels"`
}

type teamMemberSettingsModel struct {
	AllowAddRemoveApps                types.Bool `tfsdk:"allow_add_remove_apps"`
	AllowCreatePrivateChannels        types.Bool `tfsdk:"allow_create_private_channels"`
	AllowCreateUpdateChannels         types.Bool `tfsdk:"allow_create_update_channels"`
	AllowCreateUpdateRemoveConnectors types.Bool `tfsdk:"allow_create_update_remove_connectors"`
	AllowCreateUpdateRemoveTabs       types.Bool `tfsdk:"allow_create_update_remove_tabs"`
	AllowDeleteChannels               types.Bool `tfsdk:"allow_delete_channels"`
}

type teamMessagingSettingsModel struct {
	AllowChannelMentions     types.Bool `tfsdk:"allow_channel_mentions"`
	AllowOwnerDeleteMessages types.Bool `tfsdk:"allow_owner_delete_messages"`
	AllowTeamMentions        types.Bool `tfsdk:"allow_team_mentions"`
	AllowUserDeleteMessages  types.Bool `tfsdk:"allow_user_delete_messages"`
	AllowUserEditMessages    types.Bool `tfsdk:"allow_user_edit_messages"`
}

type teamSummaryModel struct {
	GuestsCount  types.Int64 `tfsdk:"guests_count"`
	MembersCount types.Int64 `tfsdk:"members_count"`
	OwnersCount  types.Int64 `tfsdk:"owners_count"`
}
