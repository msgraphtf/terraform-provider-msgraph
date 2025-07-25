---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "msgraph_site Data Source - msgraph"
subcategory: ""
description: |-
  
---

# msgraph_site (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `id` (String) The unique identifier for an entity. Read-only.

### Read-Only

- `created_by` (Attributes) Identity of the user, device, or application that created the item. Read-only. (see [below for nested schema](#nestedatt--created_by))
- `created_date_time` (String) Date and time of item creation. Read-only.
- `description` (String) Provides a user-visible description of the item. Optional.
- `display_name` (String) The full title for the site. Read-only.
- `e_tag` (String) ETag for the item. Read-only.
- `error` (Attributes) (see [below for nested schema](#nestedatt--error))
- `is_personal_site` (Boolean) Identifies whether the site is personal or not. Read-only.
- `last_modified_by` (Attributes) Identity of the user, device, and application that last modified the item. Read-only. (see [below for nested schema](#nestedatt--last_modified_by))
- `last_modified_date_time` (String) Date and time the item was last modified. Read-only.
- `name` (String) The name of the item. Read-write.
- `parent_reference` (Attributes) Parent information, if the item has a parent. Read-write. (see [below for nested schema](#nestedatt--parent_reference))
- `sharepoint_ids` (Attributes) Returns identifiers useful for SharePoint REST compatibility. Read-only. (see [below for nested schema](#nestedatt--sharepoint_ids))
- `site_collection` (Attributes) Provides details about the site's site collection. Available only on the root site. Read-only. (see [below for nested schema](#nestedatt--site_collection))
- `web_url` (String) URL that either displays the resource in the browser (for Office file formats), or is a direct link to the file (for other formats). Read-only.

<a id="nestedatt--created_by"></a>
### Nested Schema for `created_by`

Read-Only:

- `application` (Attributes) Optional. The application associated with this action. (see [below for nested schema](#nestedatt--created_by--application))
- `device` (Attributes) Optional. The device associated with this action. (see [below for nested schema](#nestedatt--created_by--device))
- `user` (Attributes) Optional. The user associated with this action. (see [below for nested schema](#nestedatt--created_by--user))

<a id="nestedatt--created_by--application"></a>
### Nested Schema for `created_by.application`

Optional:

- `id` (String) Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.

Read-Only:

- `display_name` (String) The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.


<a id="nestedatt--created_by--device"></a>
### Nested Schema for `created_by.device`

Optional:

- `id` (String) Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.

Read-Only:

- `display_name` (String) The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.


<a id="nestedatt--created_by--user"></a>
### Nested Schema for `created_by.user`

Optional:

- `id` (String) Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.

Read-Only:

- `display_name` (String) The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.



<a id="nestedatt--error"></a>
### Nested Schema for `error`

Read-Only:

- `code` (String) Represents the error code.
- `details` (Attributes List) Details of the error. (see [below for nested schema](#nestedatt--error--details))
- `inner_error` (Attributes) Details of the inner error. (see [below for nested schema](#nestedatt--error--inner_error))
- `message` (String) A non-localized message for the developer.
- `target` (String) The target of the error.

<a id="nestedatt--error--details"></a>
### Nested Schema for `error.details`

Read-Only:

- `code` (String) The error code.
- `message` (String) The error message.
- `target` (String) The target of the error.


<a id="nestedatt--error--inner_error"></a>
### Nested Schema for `error.inner_error`

Read-Only:

- `code` (String) The error code.
- `details` (Attributes List) A collection of error details. (see [below for nested schema](#nestedatt--error--inner_error--details))
- `message` (String) The error message.
- `target` (String) The target of the error.

<a id="nestedatt--error--inner_error--details"></a>
### Nested Schema for `error.inner_error.details`

Read-Only:

- `code` (String) The error code.
- `message` (String) The error message.
- `target` (String) The target of the error.




<a id="nestedatt--last_modified_by"></a>
### Nested Schema for `last_modified_by`

Read-Only:

- `application` (Attributes) Optional. The application associated with this action. (see [below for nested schema](#nestedatt--last_modified_by--application))
- `device` (Attributes) Optional. The device associated with this action. (see [below for nested schema](#nestedatt--last_modified_by--device))
- `user` (Attributes) Optional. The user associated with this action. (see [below for nested schema](#nestedatt--last_modified_by--user))

<a id="nestedatt--last_modified_by--application"></a>
### Nested Schema for `last_modified_by.application`

Optional:

- `id` (String) Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.

Read-Only:

- `display_name` (String) The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.


<a id="nestedatt--last_modified_by--device"></a>
### Nested Schema for `last_modified_by.device`

Optional:

- `id` (String) Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.

Read-Only:

- `display_name` (String) The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.


<a id="nestedatt--last_modified_by--user"></a>
### Nested Schema for `last_modified_by.user`

Optional:

- `id` (String) Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.

Read-Only:

- `display_name` (String) The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.



<a id="nestedatt--parent_reference"></a>
### Nested Schema for `parent_reference`

Optional:

- `id` (String) Unique identifier of the driveItem in the drive or a listItem in a list. Read-only.

Read-Only:

- `drive_id` (String) Unique identifier of the drive instance that contains the driveItem. Only returned if the item is located in a drive. Read-only.
- `drive_type` (String) Identifies the type of drive. Only returned if the item is located in a drive. See drive resource for values.
- `name` (String) The name of the item being referenced. Read-only.
- `path` (String) Percent-encoded path that can be used to navigate to the item. Read-only.
- `share_id` (String) A unique identifier for a shared resource that can be accessed via the Shares API.
- `sharepoint_ids` (Attributes) Returns identifiers useful for SharePoint REST compatibility. Read-only. (see [below for nested schema](#nestedatt--parent_reference--sharepoint_ids))
- `site_id` (String) For OneDrive for Business and SharePoint, this property represents the ID of the site that contains the parent document library of the driveItem resource or the parent list of the listItem resource. The value is the same as the id property of that site resource. It is an opaque string that consists of three identifiers of the site. For OneDrive, this property is not populated.

<a id="nestedatt--parent_reference--sharepoint_ids"></a>
### Nested Schema for `parent_reference.sharepoint_ids`

Read-Only:

- `list_id` (String) The unique identifier (guid) for the item's list in SharePoint.
- `list_item_id` (String) An integer identifier for the item within the containing list.
- `list_item_unique_id` (String) The unique identifier (guid) for the item within OneDrive for Business or a SharePoint site.
- `site_id` (String) The unique identifier (guid) for the item's site collection (SPSite).
- `site_url` (String) The SharePoint URL for the site that contains the item.
- `tenant_id` (String) The unique identifier (guid) for the tenancy.
- `web_id` (String) The unique identifier (guid) for the item's site (SPWeb).



<a id="nestedatt--sharepoint_ids"></a>
### Nested Schema for `sharepoint_ids`

Read-Only:

- `list_id` (String) The unique identifier (guid) for the item's list in SharePoint.
- `list_item_id` (String) An integer identifier for the item within the containing list.
- `list_item_unique_id` (String) The unique identifier (guid) for the item within OneDrive for Business or a SharePoint site.
- `site_id` (String) The unique identifier (guid) for the item's site collection (SPSite).
- `site_url` (String) The SharePoint URL for the site that contains the item.
- `tenant_id` (String) The unique identifier (guid) for the tenancy.
- `web_id` (String) The unique identifier (guid) for the item's site (SPWeb).


<a id="nestedatt--site_collection"></a>
### Nested Schema for `site_collection`

Read-Only:

- `archival_details` (Attributes) Represents whether the site collection is recently archived, fully archived, or reactivating. Possible values are: recentlyArchived, fullyArchived, reactivating, unknownFutureValue. (see [below for nested schema](#nestedatt--site_collection--archival_details))
- `data_location_code` (String) The geographic region code for where this site collection resides. Only present for multi-geo tenants. Read-only.
- `hostname` (String) The hostname for the site collection. Read-only.

<a id="nestedatt--site_collection--archival_details"></a>
### Nested Schema for `site_collection.archival_details`

Read-Only:

- `archive_status` (String) Represents the current archive status of the site collection. Returned only on $select. The possible values are: recentlyArchived, fullyArchived, reactivating, unknownFutureValue.
