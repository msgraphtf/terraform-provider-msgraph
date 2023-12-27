package sites

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sitesDataSourceModel struct {
	Value []sitesValueDataSourceModel `tfsdk:"value"`
}

type sitesValueDataSourceModel struct {
	Id                   types.String                         `tfsdk:"id"`
	CreatedBy            *sitesCreatedByDataSourceModel       `tfsdk:"created_by"`
	CreatedDateTime      types.String                         `tfsdk:"created_date_time"`
	Description          types.String                         `tfsdk:"description"`
	ETag                 types.String                         `tfsdk:"e_tag"`
	LastModifiedBy       *sitesLastModifiedByDataSourceModel  `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String                         `tfsdk:"last_modified_date_time"`
	Name                 types.String                         `tfsdk:"name"`
	ParentReference      *sitesParentReferenceDataSourceModel `tfsdk:"parent_reference"`
	WebUrl               types.String                         `tfsdk:"web_url"`
	DisplayName          types.String                         `tfsdk:"display_name"`
	Error                *sitesErrorDataSourceModel           `tfsdk:"error"`
	IsPersonalSite       types.Bool                           `tfsdk:"is_personal_site"`
	Root                 *sitesRootDataSourceModel            `tfsdk:"root"`
	SharepointIds        *sitesSharepointIdsDataSourceModel   `tfsdk:"sharepoint_ids"`
	SiteCollection       *sitesSiteCollectionDataSourceModel  `tfsdk:"site_collection"`
}

type sitesCreatedByDataSourceModel struct {
	Application *sitesApplicationDataSourceModel `tfsdk:"application"`
	Device      *sitesDeviceDataSourceModel      `tfsdk:"device"`
	User        *sitesUserDataSourceModel        `tfsdk:"user"`
}

type sitesApplicationDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type sitesDeviceDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type sitesUserDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type sitesLastModifiedByDataSourceModel struct {
	Application *sitesApplicationDataSourceModel `tfsdk:"application"`
	Device      *sitesDeviceDataSourceModel      `tfsdk:"device"`
	User        *sitesUserDataSourceModel        `tfsdk:"user"`
}

type sitesParentReferenceDataSourceModel struct {
	DriveId       types.String                       `tfsdk:"drive_id"`
	DriveType     types.String                       `tfsdk:"drive_type"`
	Id            types.String                       `tfsdk:"id"`
	Name          types.String                       `tfsdk:"name"`
	Path          types.String                       `tfsdk:"path"`
	ShareId       types.String                       `tfsdk:"share_id"`
	SharepointIds *sitesSharepointIdsDataSourceModel `tfsdk:"sharepoint_ids"`
	SiteId        types.String                       `tfsdk:"site_id"`
}

type sitesSharepointIdsDataSourceModel struct {
	ListId           types.String `tfsdk:"list_id"`
	ListItemId       types.String `tfsdk:"list_item_id"`
	ListItemUniqueId types.String `tfsdk:"list_item_unique_id"`
	SiteId           types.String `tfsdk:"site_id"`
	SiteUrl          types.String `tfsdk:"site_url"`
	TenantId         types.String `tfsdk:"tenant_id"`
	WebId            types.String `tfsdk:"web_id"`
}

type sitesErrorDataSourceModel struct {
	Code       types.String                    `tfsdk:"code"`
	Details    []sitesDetailsDataSourceModel   `tfsdk:"details"`
	InnerError *sitesInnerErrorDataSourceModel `tfsdk:"inner_error"`
	Message    types.String                    `tfsdk:"message"`
	Target     types.String                    `tfsdk:"target"`
}

type sitesDetailsDataSourceModel struct {
	Code    types.String `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

type sitesInnerErrorDataSourceModel struct {
	Code    types.String                  `tfsdk:"code"`
	Details []sitesDetailsDataSourceModel `tfsdk:"details"`
	Message types.String                  `tfsdk:"message"`
	Target  types.String                  `tfsdk:"target"`
}

type sitesRootDataSourceModel struct {
}

type sitesSiteCollectionDataSourceModel struct {
	DataLocationCode types.String              `tfsdk:"data_location_code"`
	Hostname         types.String              `tfsdk:"hostname"`
	Root             *sitesRootDataSourceModel `tfsdk:"root"`
}
