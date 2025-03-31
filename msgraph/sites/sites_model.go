package sites

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type sitesModel struct {
	Value types.List `tfsdk:"value"`
}

func (m sitesModel) AttributeTypes() map[string]attr.Type {
	sitesValue := sitesSiteModel{}
	return map[string]attr.Type{
		"value": types.ListType{ElemType: types.ObjectType{AttrTypes: sitesValue.AttributeTypes()}},
	}
}

type sitesSiteModel struct {
	Id                   types.String `tfsdk:"id"`
	CreatedBy            types.Object `tfsdk:"created_by"`
	CreatedDateTime      types.String `tfsdk:"created_date_time"`
	Description          types.String `tfsdk:"description"`
	ETag                 types.String `tfsdk:"e_tag"`
	LastModifiedBy       types.Object `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String `tfsdk:"last_modified_date_time"`
	Name                 types.String `tfsdk:"name"`
	ParentReference      types.Object `tfsdk:"parent_reference"`
	WebUrl               types.String `tfsdk:"web_url"`
	DisplayName          types.String `tfsdk:"display_name"`
	Error                types.Object `tfsdk:"error"`
	IsPersonalSite       types.Bool   `tfsdk:"is_personal_site"`
	Root                 types.Object `tfsdk:"root"`
	SharepointIds        types.Object `tfsdk:"sharepoint_ids"`
	SiteCollection       types.Object `tfsdk:"site_collection"`
}

func (m sitesSiteModel) AttributeTypes() map[string]attr.Type {
	sitesCreatedBy := sitesIdentitySetModel{}
	sitesLastModifiedBy := sitesIdentitySetModel{}
	sitesParentReference := sitesItemReferenceModel{}
	sitesError := sitesPublicErrorModel{}
	sitesRoot := sitesRootModel{}
	sitesSharepointIds := sitesSharepointIdsModel{}
	sitesSiteCollection := sitesSiteCollectionModel{}
	return map[string]attr.Type{
		"id":                      types.StringType,
		"created_by":              types.ObjectType{AttrTypes: sitesCreatedBy.AttributeTypes()},
		"created_date_time":       types.StringType,
		"description":             types.StringType,
		"e_tag":                   types.StringType,
		"last_modified_by":        types.ObjectType{AttrTypes: sitesLastModifiedBy.AttributeTypes()},
		"last_modified_date_time": types.StringType,
		"name":                    types.StringType,
		"parent_reference":        types.ObjectType{AttrTypes: sitesParentReference.AttributeTypes()},
		"web_url":                 types.StringType,
		"display_name":            types.StringType,
		"error":                   types.ObjectType{AttrTypes: sitesError.AttributeTypes()},
		"is_personal_site":        types.BoolType,
		"root":                    types.ObjectType{AttrTypes: sitesRoot.AttributeTypes()},
		"sharepoint_ids":          types.ObjectType{AttrTypes: sitesSharepointIds.AttributeTypes()},
		"site_collection":         types.ObjectType{AttrTypes: sitesSiteCollection.AttributeTypes()},
	}
}

type sitesIdentitySetModel struct {
	Application types.Object `tfsdk:"application"`
	Device      types.Object `tfsdk:"device"`
	User        types.Object `tfsdk:"user"`
}

func (m sitesIdentitySetModel) AttributeTypes() map[string]attr.Type {
	sitesApplication := sitesIdentityModel{}
	sitesDevice := sitesIdentityModel{}
	sitesUser := sitesIdentityModel{}
	return map[string]attr.Type{
		"application": types.ObjectType{AttrTypes: sitesApplication.AttributeTypes()},
		"device":      types.ObjectType{AttrTypes: sitesDevice.AttributeTypes()},
		"user":        types.ObjectType{AttrTypes: sitesUser.AttributeTypes()},
	}
}

type sitesIdentityModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

func (m sitesIdentityModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"display_name": types.StringType,
		"id":           types.StringType,
	}
}

type sitesItemReferenceModel struct {
	DriveId       types.String `tfsdk:"drive_id"`
	DriveType     types.String `tfsdk:"drive_type"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Path          types.String `tfsdk:"path"`
	ShareId       types.String `tfsdk:"share_id"`
	SharepointIds types.Object `tfsdk:"sharepoint_ids"`
	SiteId        types.String `tfsdk:"site_id"`
}

func (m sitesItemReferenceModel) AttributeTypes() map[string]attr.Type {
	sitesSharepointIds := sitesSharepointIdsModel{}
	return map[string]attr.Type{
		"drive_id":       types.StringType,
		"drive_type":     types.StringType,
		"id":             types.StringType,
		"name":           types.StringType,
		"path":           types.StringType,
		"share_id":       types.StringType,
		"sharepoint_ids": types.ObjectType{AttrTypes: sitesSharepointIds.AttributeTypes()},
		"site_id":        types.StringType,
	}
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

func (m sitesSharepointIdsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"list_id":             types.StringType,
		"list_item_id":        types.StringType,
		"list_item_unique_id": types.StringType,
		"site_id":             types.StringType,
		"site_url":            types.StringType,
		"tenant_id":           types.StringType,
		"web_id":              types.StringType,
	}
}

type sitesPublicErrorModel struct {
	Code       types.String `tfsdk:"code"`
	Details    types.List   `tfsdk:"details"`
	InnerError types.Object `tfsdk:"inner_error"`
	Message    types.String `tfsdk:"message"`
	Target     types.String `tfsdk:"target"`
}

func (m sitesPublicErrorModel) AttributeTypes() map[string]attr.Type {
	sitesDetails := sitesPublicErrorDetailModel{}
	sitesInnerError := sitesPublicInnerErrorModel{}
	return map[string]attr.Type{
		"code":        types.StringType,
		"details":     types.ListType{ElemType: types.ObjectType{AttrTypes: sitesDetails.AttributeTypes()}},
		"inner_error": types.ObjectType{AttrTypes: sitesInnerError.AttributeTypes()},
		"message":     types.StringType,
		"target":      types.StringType,
	}
}

type sitesPublicErrorDetailModel struct {
	Code    types.String `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

func (m sitesPublicErrorDetailModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"code":    types.StringType,
		"message": types.StringType,
		"target":  types.StringType,
	}
}

type sitesPublicInnerErrorModel struct {
	Code    types.String `tfsdk:"code"`
	Details types.List   `tfsdk:"details"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

func (m sitesPublicInnerErrorModel) AttributeTypes() map[string]attr.Type {
	sitesDetails := sitesPublicErrorDetailModel{}
	return map[string]attr.Type{
		"code":    types.StringType,
		"details": types.ListType{ElemType: types.ObjectType{AttrTypes: sitesDetails.AttributeTypes()}},
		"message": types.StringType,
		"target":  types.StringType,
	}
}

type sitesRootModel struct {
}

func (m sitesRootModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{}
}

type sitesSiteCollectionModel struct {
	ArchivalDetails  types.Object `tfsdk:"archival_details"`
	DataLocationCode types.String `tfsdk:"data_location_code"`
	Hostname         types.String `tfsdk:"hostname"`
	Root             types.Object `tfsdk:"root"`
}

func (m sitesSiteCollectionModel) AttributeTypes() map[string]attr.Type {
	sitesArchivalDetails := sitesSiteArchivalDetailsModel{}
	sitesRoot := sitesRootModel{}
	return map[string]attr.Type{
		"archival_details":   types.ObjectType{AttrTypes: sitesArchivalDetails.AttributeTypes()},
		"data_location_code": types.StringType,
		"hostname":           types.StringType,
		"root":               types.ObjectType{AttrTypes: sitesRoot.AttributeTypes()},
	}
}

type sitesSiteArchivalDetailsModel struct {
	ArchiveStatus types.String `tfsdk:"archive_status"`
}

func (m sitesSiteArchivalDetailsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"archive_status": types.StringType,
	}
}
