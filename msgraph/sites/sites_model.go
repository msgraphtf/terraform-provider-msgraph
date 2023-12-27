package sites

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sitesModel struct {
	Value []sitesValueModel `tfsdk:"value"`
}

type sitesValueModel struct {
	Id                   types.String               `tfsdk:"id"`
	CreatedBy            *sitesCreatedByModel       `tfsdk:"created_by"`
	CreatedDateTime      types.String               `tfsdk:"created_date_time"`
	Description          types.String               `tfsdk:"description"`
	ETag                 types.String               `tfsdk:"e_tag"`
	LastModifiedBy       *sitesLastModifiedByModel  `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String               `tfsdk:"last_modified_date_time"`
	Name                 types.String               `tfsdk:"name"`
	ParentReference      *sitesParentReferenceModel `tfsdk:"parent_reference"`
	WebUrl               types.String               `tfsdk:"web_url"`
	DisplayName          types.String               `tfsdk:"display_name"`
	Error                *sitesErrorModel           `tfsdk:"error"`
	IsPersonalSite       types.Bool                 `tfsdk:"is_personal_site"`
	Root                 *sitesRootModel            `tfsdk:"root"`
	SharepointIds        *sitesSharepointIdsModel   `tfsdk:"sharepoint_ids"`
	SiteCollection       *sitesSiteCollectionModel  `tfsdk:"site_collection"`
}

type sitesCreatedByModel struct {
	Application *sitesApplicationModel `tfsdk:"application"`
	Device      *sitesDeviceModel      `tfsdk:"device"`
	User        *sitesUserModel        `tfsdk:"user"`
}

type sitesApplicationModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type sitesDeviceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type sitesUserModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type sitesLastModifiedByModel struct {
	Application *sitesApplicationModel `tfsdk:"application"`
	Device      *sitesDeviceModel      `tfsdk:"device"`
	User        *sitesUserModel        `tfsdk:"user"`
}

type sitesParentReferenceModel struct {
	DriveId       types.String             `tfsdk:"drive_id"`
	DriveType     types.String             `tfsdk:"drive_type"`
	Id            types.String             `tfsdk:"id"`
	Name          types.String             `tfsdk:"name"`
	Path          types.String             `tfsdk:"path"`
	ShareId       types.String             `tfsdk:"share_id"`
	SharepointIds *sitesSharepointIdsModel `tfsdk:"sharepoint_ids"`
	SiteId        types.String             `tfsdk:"site_id"`
}

type sitesSharepointIdsModel struct {
	ListId           types.String `tfsdk:"list_id"`
	ListItemId       types.String `tfsdk:"list_item_id"`
	ListItemUniqueId types.String `tfsdk:"list_item_unique_id"`
	SiteId           types.String `tfsdk:"site_id"`
	SiteUrl          types.String `tfsdk:"site_url"`
	TenantId         types.String `tfsdk:"tenant_id"`
	WebId            types.String `tfsdk:"web_id"`
}

type sitesErrorModel struct {
	Code       types.String          `tfsdk:"code"`
	Details    []sitesDetailsModel   `tfsdk:"details"`
	InnerError *sitesInnerErrorModel `tfsdk:"inner_error"`
	Message    types.String          `tfsdk:"message"`
	Target     types.String          `tfsdk:"target"`
}

type sitesDetailsModel struct {
	Code    types.String `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

type sitesInnerErrorModel struct {
	Code    types.String        `tfsdk:"code"`
	Details []sitesDetailsModel `tfsdk:"details"`
	Message types.String        `tfsdk:"message"`
	Target  types.String        `tfsdk:"target"`
}

type sitesRootModel struct {
}

type sitesSiteCollectionModel struct {
	DataLocationCode types.String    `tfsdk:"data_location_code"`
	Hostname         types.String    `tfsdk:"hostname"`
	Root             *sitesRootModel `tfsdk:"root"`
}
