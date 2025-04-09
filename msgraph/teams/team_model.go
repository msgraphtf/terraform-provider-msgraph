package teams

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type teamModel struct {
	Id                types.String `tfsdk:"id"`
	Classification    types.String `tfsdk:"classification"`
	CreatedDateTime   types.String `tfsdk:"created_date_time"`
	Description       types.String `tfsdk:"description"`
	DisplayName       types.String `tfsdk:"display_name"`
	FunSettings       types.Object `tfsdk:"fun_settings"`
	GuestSettings     types.Object `tfsdk:"guest_settings"`
	InternalId        types.String `tfsdk:"internal_id"`
	IsArchived        types.Bool   `tfsdk:"is_archived"`
	MemberSettings    types.Object `tfsdk:"member_settings"`
	MessagingSettings types.Object `tfsdk:"messaging_settings"`
	Specialization    types.String `tfsdk:"specialization"`
	TenantId          types.String `tfsdk:"tenant_id"`
	Visibility        types.String `tfsdk:"visibility"`
	WebUrl            types.String `tfsdk:"web_url"`
}

func (m teamModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                 types.StringType,
		"classification":     types.StringType,
		"created_date_time":  types.StringType,
		"description":        types.StringType,
		"display_name":       types.StringType,
		"fun_settings":       types.ObjectType{AttrTypes: teamTeamFunSettingsModel{}.AttributeTypes()},
		"guest_settings":     types.ObjectType{AttrTypes: teamTeamGuestSettingsModel{}.AttributeTypes()},
		"internal_id":        types.StringType,
		"is_archived":        types.BoolType,
		"member_settings":    types.ObjectType{AttrTypes: teamTeamMemberSettingsModel{}.AttributeTypes()},
		"messaging_settings": types.ObjectType{AttrTypes: teamTeamMessagingSettingsModel{}.AttributeTypes()},
		"specialization":     types.StringType,
		"tenant_id":          types.StringType,
		"visibility":         types.StringType,
		"web_url":            types.StringType,
	}
}

type teamTeamFunSettingsModel struct {
	AllowCustomMemes      types.Bool   `tfsdk:"allow_custom_memes"`
	AllowGiphy            types.Bool   `tfsdk:"allow_giphy"`
	AllowStickersAndMemes types.Bool   `tfsdk:"allow_stickers_and_memes"`
	GiphyContentRating    types.String `tfsdk:"giphy_content_rating"`
}

func (m teamTeamFunSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allow_custom_memes":       types.BoolType,
		"allow_giphy":              types.BoolType,
		"allow_stickers_and_memes": types.BoolType,
		"giphy_content_rating":     types.StringType,
	}
}

type teamTeamGuestSettingsModel struct {
	AllowCreateUpdateChannels types.Bool `tfsdk:"allow_create_update_channels"`
	AllowDeleteChannels       types.Bool `tfsdk:"allow_delete_channels"`
}

func (m teamTeamGuestSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allow_create_update_channels": types.BoolType,
		"allow_delete_channels":        types.BoolType,
	}
}

type teamTeamMemberSettingsModel struct {
	AllowAddRemoveApps                types.Bool `tfsdk:"allow_add_remove_apps"`
	AllowCreatePrivateChannels        types.Bool `tfsdk:"allow_create_private_channels"`
	AllowCreateUpdateChannels         types.Bool `tfsdk:"allow_create_update_channels"`
	AllowCreateUpdateRemoveConnectors types.Bool `tfsdk:"allow_create_update_remove_connectors"`
	AllowCreateUpdateRemoveTabs       types.Bool `tfsdk:"allow_create_update_remove_tabs"`
	AllowDeleteChannels               types.Bool `tfsdk:"allow_delete_channels"`
}

func (m teamTeamMemberSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allow_add_remove_apps":                 types.BoolType,
		"allow_create_private_channels":         types.BoolType,
		"allow_create_update_channels":          types.BoolType,
		"allow_create_update_remove_connectors": types.BoolType,
		"allow_create_update_remove_tabs":       types.BoolType,
		"allow_delete_channels":                 types.BoolType,
	}
}

type teamTeamMessagingSettingsModel struct {
	AllowChannelMentions     types.Bool `tfsdk:"allow_channel_mentions"`
	AllowOwnerDeleteMessages types.Bool `tfsdk:"allow_owner_delete_messages"`
	AllowTeamMentions        types.Bool `tfsdk:"allow_team_mentions"`
	AllowUserDeleteMessages  types.Bool `tfsdk:"allow_user_delete_messages"`
	AllowUserEditMessages    types.Bool `tfsdk:"allow_user_edit_messages"`
}

func (m teamTeamMessagingSettingsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allow_channel_mentions":      types.BoolType,
		"allow_owner_delete_messages": types.BoolType,
		"allow_team_mentions":         types.BoolType,
		"allow_user_delete_messages":  types.BoolType,
		"allow_user_edit_messages":    types.BoolType,
	}
}
