package sites

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type siteModel struct {
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
	SharepointIds        types.Object `tfsdk:"sharepoint_ids"`
	SiteCollection       types.Object `tfsdk:"site_collection"`
}

func (m siteModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                      types.StringType,
		"created_by":              types.ObjectType{AttrTypes: siteIdentitySetModel{}.AttributeTypes()},
		"created_date_time":       types.StringType,
		"description":             types.StringType,
		"e_tag":                   types.StringType,
		"last_modified_by":        types.ObjectType{AttrTypes: siteIdentitySetModel{}.AttributeTypes()},
		"last_modified_date_time": types.StringType,
		"name":                    types.StringType,
		"parent_reference":        types.ObjectType{AttrTypes: siteItemReferenceModel{}.AttributeTypes()},
		"web_url":                 types.StringType,
		"display_name":            types.StringType,
		"error":                   types.ObjectType{AttrTypes: sitePublicErrorModel{}.AttributeTypes()},
		"is_personal_site":        types.BoolType,
		"sharepoint_ids":          types.ObjectType{AttrTypes: siteSharepointIdsModel{}.AttributeTypes()},
		"site_collection":         types.ObjectType{AttrTypes: siteSiteCollectionModel{}.AttributeTypes()},
	}
}

type siteIdentitySetModel struct {
	Application types.Object `tfsdk:"application"`
	Device      types.Object `tfsdk:"device"`
	User        types.Object `tfsdk:"user"`
}

func (m siteIdentitySetModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"application": types.ObjectType{AttrTypes: siteIdentityModel{}.AttributeTypes()},
		"device":      types.ObjectType{AttrTypes: siteIdentityModel{}.AttributeTypes()},
		"user":        types.ObjectType{AttrTypes: siteIdentityModel{}.AttributeTypes()},
	}
}

type siteIdentityModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

func (m siteIdentityModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"display_name": types.StringType,
		"id":           types.StringType,
	}
}

type siteItemReferenceModel struct {
	DriveId       types.String `tfsdk:"drive_id"`
	DriveType     types.String `tfsdk:"drive_type"`
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Path          types.String `tfsdk:"path"`
	ShareId       types.String `tfsdk:"share_id"`
	SharepointIds types.Object `tfsdk:"sharepoint_ids"`
	SiteId        types.String `tfsdk:"site_id"`
}

func (m siteItemReferenceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"drive_id":       types.StringType,
		"drive_type":     types.StringType,
		"id":             types.StringType,
		"name":           types.StringType,
		"path":           types.StringType,
		"share_id":       types.StringType,
		"sharepoint_ids": types.ObjectType{AttrTypes: siteSharepointIdsModel{}.AttributeTypes()},
		"site_id":        types.StringType,
	}
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

func (m siteSharepointIdsModel) AttributeTypes() map[string]attr.Type {
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

type sitePublicErrorModel struct {
	Code       types.String `tfsdk:"code"`
	Details    types.List   `tfsdk:"details"`
	InnerError types.Object `tfsdk:"inner_error"`
	Message    types.String `tfsdk:"message"`
	Target     types.String `tfsdk:"target"`
}

func (m sitePublicErrorModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"code":        types.StringType,
		"details":     types.ListType{ElemType: types.ObjectType{AttrTypes: sitePublicErrorDetailModel{}.AttributeTypes()}},
		"inner_error": types.ObjectType{AttrTypes: sitePublicInnerErrorModel{}.AttributeTypes()},
		"message":     types.StringType,
		"target":      types.StringType,
	}
}

type sitePublicErrorDetailModel struct {
	Code    types.String `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

func (m sitePublicErrorDetailModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"code":    types.StringType,
		"message": types.StringType,
		"target":  types.StringType,
	}
}

type sitePublicInnerErrorModel struct {
	Code    types.String `tfsdk:"code"`
	Details types.List   `tfsdk:"details"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

func (m sitePublicInnerErrorModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"code":    types.StringType,
		"details": types.ListType{ElemType: types.ObjectType{AttrTypes: sitePublicErrorDetailModel{}.AttributeTypes()}},
		"message": types.StringType,
		"target":  types.StringType,
	}
}

type siteSiteCollectionModel struct {
	ArchivalDetails  types.Object `tfsdk:"archival_details"`
	DataLocationCode types.String `tfsdk:"data_location_code"`
	Hostname         types.String `tfsdk:"hostname"`
}

func (m siteSiteCollectionModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"archival_details":   types.ObjectType{AttrTypes: siteSiteArchivalDetailsModel{}.AttributeTypes()},
		"data_location_code": types.StringType,
		"hostname":           types.StringType,
	}
}

type siteSiteArchivalDetailsModel struct {
	ArchiveStatus types.String `tfsdk:"archive_status"`
}

func (m siteSiteArchivalDetailsModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"archive_status": types.StringType,
	}
}
