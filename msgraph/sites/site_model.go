package sites

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type siteModel struct {
	Id                   types.String              `tfsdk:"id"`
	CreatedBy            *siteCreatedByModel       `tfsdk:"created_by"`
	CreatedDateTime      types.String              `tfsdk:"created_date_time"`
	Description          types.String              `tfsdk:"description"`
	ETag                 types.String              `tfsdk:"e_tag"`
	LastModifiedBy       *siteLastModifiedByModel  `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String              `tfsdk:"last_modified_date_time"`
	Name                 types.String              `tfsdk:"name"`
	ParentReference      *siteParentReferenceModel `tfsdk:"parent_reference"`
	WebUrl               types.String              `tfsdk:"web_url"`
	DisplayName          types.String              `tfsdk:"display_name"`
	Error                *siteErrorModel           `tfsdk:"error"`
	IsPersonalSite       types.Bool                `tfsdk:"is_personal_site"`
	Root                 *siteRootModel            `tfsdk:"root"`
	SharepointIds        *siteSharepointIdsModel   `tfsdk:"sharepoint_ids"`
	SiteCollection       *siteSiteCollectionModel  `tfsdk:"site_collection"`
}

type siteCreatedByModel struct {
	Application *siteApplicationModel `tfsdk:"application"`
	Device      *siteDeviceModel      `tfsdk:"device"`
	User        *siteUserModel        `tfsdk:"user"`
}

type siteApplicationModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteDeviceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteUserModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteLastModifiedByModel struct {
	Application *siteApplicationModel `tfsdk:"application"`
	Device      *siteDeviceModel      `tfsdk:"device"`
	User        *siteUserModel        `tfsdk:"user"`
}

type siteParentReferenceModel struct {
	DriveId       types.String            `tfsdk:"drive_id"`
	DriveType     types.String            `tfsdk:"drive_type"`
	Id            types.String            `tfsdk:"id"`
	Name          types.String            `tfsdk:"name"`
	Path          types.String            `tfsdk:"path"`
	ShareId       types.String            `tfsdk:"share_id"`
	SharepointIds *siteSharepointIdsModel `tfsdk:"sharepoint_ids"`
	SiteId        types.String            `tfsdk:"site_id"`
}

type siteSharepointIdsModel struct {
	ListId           types.String `tfsdk:"list_id"`
	ListItemId       types.String `tfsdk:"list_item_id"`
	ListItemUniqueId types.String `tfsdk:"list_item_unique_id"`
	SiteId           types.String `tfsdk:"site_id"`
	SiteUrl          types.String `tfsdk:"site_url"`
	TenantId         types.String `tfsdk:"tenant_id"`
	WebId            types.String `tfsdk:"web_id"`
}

type siteErrorModel struct {
	Code       types.String         `tfsdk:"code"`
	Details    []siteDetailsModel   `tfsdk:"details"`
	InnerError *siteInnerErrorModel `tfsdk:"inner_error"`
	Message    types.String         `tfsdk:"message"`
	Target     types.String         `tfsdk:"target"`
}

type siteDetailsModel struct {
	Code    types.String `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

type siteInnerErrorModel struct {
	Code    types.String       `tfsdk:"code"`
	Details []siteDetailsModel `tfsdk:"details"`
	Message types.String       `tfsdk:"message"`
	Target  types.String       `tfsdk:"target"`
}

type siteRootModel struct {
}

type siteSiteCollectionModel struct {
	DataLocationCode types.String   `tfsdk:"data_location_code"`
	Hostname         types.String   `tfsdk:"hostname"`
	Root             *siteRootModel `tfsdk:"root"`
}
