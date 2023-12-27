package sites

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type siteDataSourceModel struct {
	Id                   types.String                        `tfsdk:"id"`
	CreatedBy            *siteCreatedByDataSourceModel       `tfsdk:"created_by"`
	CreatedDateTime      types.String                        `tfsdk:"created_date_time"`
	Description          types.String                        `tfsdk:"description"`
	ETag                 types.String                        `tfsdk:"e_tag"`
	LastModifiedBy       *siteLastModifiedByDataSourceModel  `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String                        `tfsdk:"last_modified_date_time"`
	Name                 types.String                        `tfsdk:"name"`
	ParentReference      *siteParentReferenceDataSourceModel `tfsdk:"parent_reference"`
	WebUrl               types.String                        `tfsdk:"web_url"`
	DisplayName          types.String                        `tfsdk:"display_name"`
	Error                *siteErrorDataSourceModel           `tfsdk:"error"`
	IsPersonalSite       types.Bool                          `tfsdk:"is_personal_site"`
	Root                 *siteRootDataSourceModel            `tfsdk:"root"`
	SharepointIds        *siteSharepointIdsDataSourceModel   `tfsdk:"sharepoint_ids"`
	SiteCollection       *siteSiteCollectionDataSourceModel  `tfsdk:"site_collection"`
}

type siteCreatedByDataSourceModel struct {
	Application *siteApplicationDataSourceModel `tfsdk:"application"`
	Device      *siteDeviceDataSourceModel      `tfsdk:"device"`
	User        *siteUserDataSourceModel        `tfsdk:"user"`
}

type siteApplicationDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteDeviceDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteUserDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteLastModifiedByDataSourceModel struct {
	Application *siteApplicationDataSourceModel `tfsdk:"application"`
	Device      *siteDeviceDataSourceModel      `tfsdk:"device"`
	User        *siteUserDataSourceModel        `tfsdk:"user"`
}

type siteParentReferenceDataSourceModel struct {
	DriveId       types.String                      `tfsdk:"drive_id"`
	DriveType     types.String                      `tfsdk:"drive_type"`
	Id            types.String                      `tfsdk:"id"`
	Name          types.String                      `tfsdk:"name"`
	Path          types.String                      `tfsdk:"path"`
	ShareId       types.String                      `tfsdk:"share_id"`
	SharepointIds *siteSharepointIdsDataSourceModel `tfsdk:"sharepoint_ids"`
	SiteId        types.String                      `tfsdk:"site_id"`
}

type siteSharepointIdsDataSourceModel struct {
	ListId           types.String `tfsdk:"list_id"`
	ListItemId       types.String `tfsdk:"list_item_id"`
	ListItemUniqueId types.String `tfsdk:"list_item_unique_id"`
	SiteId           types.String `tfsdk:"site_id"`
	SiteUrl          types.String `tfsdk:"site_url"`
	TenantId         types.String `tfsdk:"tenant_id"`
	WebId            types.String `tfsdk:"web_id"`
}

type siteErrorDataSourceModel struct {
	Code       types.String                   `tfsdk:"code"`
	Details    []siteDetailsDataSourceModel   `tfsdk:"details"`
	InnerError *siteInnerErrorDataSourceModel `tfsdk:"inner_error"`
	Message    types.String                   `tfsdk:"message"`
	Target     types.String                   `tfsdk:"target"`
}

type siteDetailsDataSourceModel struct {
	Code    types.String `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

type siteInnerErrorDataSourceModel struct {
	Code    types.String                 `tfsdk:"code"`
	Details []siteDetailsDataSourceModel `tfsdk:"details"`
	Message types.String                 `tfsdk:"message"`
	Target  types.String                 `tfsdk:"target"`
}

type siteRootDataSourceModel struct {
}

type siteSiteCollectionDataSourceModel struct {
	DataLocationCode types.String             `tfsdk:"data_location_code"`
	Hostname         types.String             `tfsdk:"hostname"`
	Root             *siteRootDataSourceModel `tfsdk:"root"`
}
