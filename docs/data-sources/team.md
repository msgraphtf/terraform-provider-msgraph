---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "msgraph_team Data Source - msgraph"
subcategory: ""
description: |-
  
---

# msgraph_team (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) The unique identifier for an entity. Read-only.

### Read-Only

- `classification` (String) An optional label. Typically describes the data or business sensitivity of the team. Must match one of a pre-configured set in the tenant's directory.
- `created_date_time` (String) Timestamp at which the team was created.
- `description` (String) An optional description for the team. Maximum length: 1024 characters.
- `display_name` (String) The name of the team.
- `fun_settings` (Attributes) Settings to configure use of Giphy, memes, and stickers in the team. (see [below for nested schema](#nestedatt--fun_settings))
- `guest_settings` (Attributes) Settings to configure whether guests can create, update, or delete channels in the team. (see [below for nested schema](#nestedatt--guest_settings))
- `internal_id` (String) A unique ID for the team that has been used in a few places such as the audit log/Office 365 Management Activity API.
- `is_archived` (Boolean) Whether this team is in read-only mode.
- `member_settings` (Attributes) Settings to configure whether members can perform certain actions, for example, create channels and add bots, in the team. (see [below for nested schema](#nestedatt--member_settings))
- `messaging_settings` (Attributes) Settings to configure messaging and mentions in the team. (see [below for nested schema](#nestedatt--messaging_settings))
- `specialization` (String) Optional. Indicates whether the team is intended for a particular use case.  Each team specialization has access to unique behaviors and experiences targeted to its use case.
- `tenant_id` (String) The ID of the Microsoft Entra tenant.
- `visibility` (String) The visibility of the group and team. Defaults to Public.
- `web_url` (String) A hyperlink that will go to the team in the Microsoft Teams client. This is the URL that you get when you right-click a team in the Microsoft Teams client and select Get link to team. This URL should be treated as an opaque blob, and not parsed.

<a id="nestedatt--fun_settings"></a>
### Nested Schema for `fun_settings`

Read-Only:

- `allow_custom_memes` (Boolean) If set to true, enables users to include custom memes.
- `allow_giphy` (Boolean) If set to true, enables Giphy use.
- `allow_stickers_and_memes` (Boolean) If set to true, enables users to include stickers and memes.
- `giphy_content_rating` (String) Giphy content rating. Possible values are: moderate, strict.


<a id="nestedatt--guest_settings"></a>
### Nested Schema for `guest_settings`

Read-Only:

- `allow_create_update_channels` (Boolean) If set to true, guests can add and update channels.
- `allow_delete_channels` (Boolean) If set to true, guests can delete channels.


<a id="nestedatt--member_settings"></a>
### Nested Schema for `member_settings`

Read-Only:

- `allow_add_remove_apps` (Boolean) If set to true, members can add and remove apps.
- `allow_create_private_channels` (Boolean) If set to true, members can add and update private channels.
- `allow_create_update_channels` (Boolean) If set to true, members can add and update channels.
- `allow_create_update_remove_connectors` (Boolean) If set to true, members can add, update, and remove connectors.
- `allow_create_update_remove_tabs` (Boolean) If set to true, members can add, update, and remove tabs.
- `allow_delete_channels` (Boolean) If set to true, members can delete channels.


<a id="nestedatt--messaging_settings"></a>
### Nested Schema for `messaging_settings`

Read-Only:

- `allow_channel_mentions` (Boolean) If set to true, @channel mentions are allowed.
- `allow_owner_delete_messages` (Boolean) If set to true, owners can delete any message.
- `allow_team_mentions` (Boolean) If set to true, @team mentions are allowed.
- `allow_user_delete_messages` (Boolean) If set to true, users can delete their messages.
- `allow_user_edit_messages` (Boolean) If set to true, users can edit their messages.
