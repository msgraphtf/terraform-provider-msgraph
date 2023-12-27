package sites

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/sites"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &sitesDataSource{}
	_ datasource.DataSourceWithConfigure = &sitesDataSource{}
)

// NewSitesDataSource is a helper function to simplify the provider implementation.
func NewSitesDataSource() datasource.DataSource {
	return &sitesDataSource{}
}

// sitesDataSource is the data source implementation.
type sitesDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *sitesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sites"
}

// Configure adds the provider configured client to the data source.
func (d *sitesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *sitesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"value": schema.ListNestedAttribute{
				Description: "",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "The unique identifier for an entity. Read-only.",
							Computed:    true,
						},
						"created_by": schema.SingleNestedAttribute{
							Description: "Identity of the user, device, or application that created the item. Read-only.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"application": schema.SingleNestedAttribute{
									Description: "Optional. The application associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity. The display name might not always be available or up to date. For example, if a user changes their display name, the API might show the new value in a future response, but the items associated with the user won't show up as having changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity.",
											Computed:    true,
										},
									},
								},
								"device": schema.SingleNestedAttribute{
									Description: "Optional. The device associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity. The display name might not always be available or up to date. For example, if a user changes their display name, the API might show the new value in a future response, but the items associated with the user won't show up as having changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity.",
											Computed:    true,
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description: "Optional. The user associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity. The display name might not always be available or up to date. For example, if a user changes their display name, the API might show the new value in a future response, but the items associated with the user won't show up as having changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity.",
											Computed:    true,
										},
									},
								},
							},
						},
						"created_date_time": schema.StringAttribute{
							Description: "Date and time of item creation. Read-only.",
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: "Provides a user-visible description of the item. Optional.",
							Computed:    true,
						},
						"e_tag": schema.StringAttribute{
							Description: "ETag for the item. Read-only.",
							Computed:    true,
						},
						"last_modified_by": schema.SingleNestedAttribute{
							Description: "Identity of the user, device, and application that last modified the item. Read-only.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"application": schema.SingleNestedAttribute{
									Description: "Optional. The application associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity. The display name might not always be available or up to date. For example, if a user changes their display name, the API might show the new value in a future response, but the items associated with the user won't show up as having changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity.",
											Computed:    true,
										},
									},
								},
								"device": schema.SingleNestedAttribute{
									Description: "Optional. The device associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity. The display name might not always be available or up to date. For example, if a user changes their display name, the API might show the new value in a future response, but the items associated with the user won't show up as having changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity.",
											Computed:    true,
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description: "Optional. The user associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity. The display name might not always be available or up to date. For example, if a user changes their display name, the API might show the new value in a future response, but the items associated with the user won't show up as having changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity.",
											Computed:    true,
										},
									},
								},
							},
						},
						"last_modified_date_time": schema.StringAttribute{
							Description: "Date and time the item was last modified. Read-only.",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "The name of the item. Read-write.",
							Computed:    true,
						},
						"parent_reference": schema.SingleNestedAttribute{
							Description: "Parent information, if the item has a parent. Read-write.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"drive_id": schema.StringAttribute{
									Description: "Unique identifier of the drive instance that contains the driveItem. Only returned if the item is located in a [drive][]. Read-only.",
									Computed:    true,
								},
								"drive_type": schema.StringAttribute{
									Description: "Identifies the type of drive. Only returned if the item is located in a [drive][]. See [drive][] resource for values.",
									Computed:    true,
								},
								"id": schema.StringAttribute{
									Description: "Unique identifier of the driveItem in the drive or a listItem in a list. Read-only.",
									Computed:    true,
								},
								"name": schema.StringAttribute{
									Description: "The name of the item being referenced. Read-only.",
									Computed:    true,
								},
								"path": schema.StringAttribute{
									Description: "Path that can be used to navigate to the item. Read-only.",
									Computed:    true,
								},
								"share_id": schema.StringAttribute{
									Description: "A unique identifier for a shared resource that can be accessed via the [Shares][] API.",
									Computed:    true,
								},
								"sharepoint_ids": schema.SingleNestedAttribute{
									Description: "Returns identifiers useful for SharePoint REST compatibility. Read-only.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"list_id": schema.StringAttribute{
											Description: "The unique identifier (guid) for the item's list in SharePoint.",
											Computed:    true,
										},
										"list_item_id": schema.StringAttribute{
											Description: "An integer identifier for the item within the containing list.",
											Computed:    true,
										},
										"list_item_unique_id": schema.StringAttribute{
											Description: "The unique identifier (guid) for the item within OneDrive for Business or a SharePoint site.",
											Computed:    true,
										},
										"site_id": schema.StringAttribute{
											Description: "The unique identifier (guid) for the item's site collection (SPSite).",
											Computed:    true,
										},
										"site_url": schema.StringAttribute{
											Description: "The SharePoint URL for the site that contains the item.",
											Computed:    true,
										},
										"tenant_id": schema.StringAttribute{
											Description: "The unique identifier (guid) for the tenancy.",
											Computed:    true,
										},
										"web_id": schema.StringAttribute{
											Description: "The unique identifier (guid) for the item's site (SPWeb).",
											Computed:    true,
										},
									},
								},
								"site_id": schema.StringAttribute{
									Description: "For OneDrive for Business and SharePoint, this property represents the ID of the site that contains the parent document library of the driveItem resource or the parent list of the listItem resource. The value is the same as the id property of that [site][] resource. It is an opaque string that consists of three identifiers of the site. For OneDrive, this property is not populated.",
									Computed:    true,
								},
							},
						},
						"web_url": schema.StringAttribute{
							Description: "URL that either displays the resource in the browser (for Office file formats), or is a direct link to the file (for other formats). Read-only.",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "The full title for the site. Read-only.",
							Computed:    true,
						},
						"error": schema.SingleNestedAttribute{
							Description: "",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"code": schema.StringAttribute{
									Description: "Represents the error code.",
									Computed:    true,
								},
								"details": schema.ListNestedAttribute{
									Description: "Details of the error.",
									Computed:    true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"code": schema.StringAttribute{
												Description: "The error code.",
												Computed:    true,
											},
											"message": schema.StringAttribute{
												Description: "The error message.",
												Computed:    true,
											},
											"target": schema.StringAttribute{
												Description: "The target of the error.",
												Computed:    true,
											},
										},
									},
								},
								"inner_error": schema.SingleNestedAttribute{
									Description: "Details of the inner error.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"code": schema.StringAttribute{
											Description: "The error code.",
											Computed:    true,
										},
										"details": schema.ListNestedAttribute{
											Description: "A collection of error details.",
											Computed:    true,
											NestedObject: schema.NestedAttributeObject{
												Attributes: map[string]schema.Attribute{
													"code": schema.StringAttribute{
														Description: "The error code.",
														Computed:    true,
													},
													"message": schema.StringAttribute{
														Description: "The error message.",
														Computed:    true,
													},
													"target": schema.StringAttribute{
														Description: "The target of the error.",
														Computed:    true,
													},
												},
											},
										},
										"message": schema.StringAttribute{
											Description: "The error message.",
											Computed:    true,
										},
										"target": schema.StringAttribute{
											Description: "The target of the error.",
											Computed:    true,
										},
									},
								},
								"message": schema.StringAttribute{
									Description: "A non-localized message for the developer.",
									Computed:    true,
								},
								"target": schema.StringAttribute{
									Description: "The target of the error.",
									Computed:    true,
								},
							},
						},
						"is_personal_site": schema.BoolAttribute{
							Description: "",
							Computed:    true,
						},
						"root": schema.SingleNestedAttribute{
							Description: "If present, indicates that this is the root site in the site collection. Read-only.",
							Computed:    true,
							Attributes:  map[string]schema.Attribute{},
						},
						"sharepoint_ids": schema.SingleNestedAttribute{
							Description: "Returns identifiers useful for SharePoint REST compatibility. Read-only.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"list_id": schema.StringAttribute{
									Description: "The unique identifier (guid) for the item's list in SharePoint.",
									Computed:    true,
								},
								"list_item_id": schema.StringAttribute{
									Description: "An integer identifier for the item within the containing list.",
									Computed:    true,
								},
								"list_item_unique_id": schema.StringAttribute{
									Description: "The unique identifier (guid) for the item within OneDrive for Business or a SharePoint site.",
									Computed:    true,
								},
								"site_id": schema.StringAttribute{
									Description: "The unique identifier (guid) for the item's site collection (SPSite).",
									Computed:    true,
								},
								"site_url": schema.StringAttribute{
									Description: "The SharePoint URL for the site that contains the item.",
									Computed:    true,
								},
								"tenant_id": schema.StringAttribute{
									Description: "The unique identifier (guid) for the tenancy.",
									Computed:    true,
								},
								"web_id": schema.StringAttribute{
									Description: "The unique identifier (guid) for the item's site (SPWeb).",
									Computed:    true,
								},
							},
						},
						"site_collection": schema.SingleNestedAttribute{
							Description: "Provides details about the site's site collection. Available only on the root site. Read-only.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"data_location_code": schema.StringAttribute{
									Description: "The geographic region code for where this site collection resides. Read-only.",
									Computed:    true,
								},
								"hostname": schema.StringAttribute{
									Description: "The hostname for the site collection. Read-only.",
									Computed:    true,
								},
								"root": schema.SingleNestedAttribute{
									Description: "If present, indicates that this is a root site collection in SharePoint. Read-only.",
									Computed:    true,
									Attributes:  map[string]schema.Attribute{},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *sitesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state sitesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := sites.SitesRequestBuilderGetRequestConfiguration{
		QueryParameters: &sites.SitesRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"createdBy",
				"createdDateTime",
				"description",
				"eTag",
				"lastModifiedBy",
				"lastModifiedDateTime",
				"name",
				"parentReference",
				"webUrl",
				"displayName",
				"error",
				"isPersonalSite",
				"root",
				"sharepointIds",
				"siteCollection",
				"createdByUser",
				"lastModifiedByUser",
				"analytics",
				"columns",
				"contentTypes",
				"drive",
				"drives",
				"externalColumns",
				"items",
				"lists",
				"operations",
				"permissions",
				"sites",
				"termStore",
				"termStores",
				"onenote",
			},
		},
	}

	result, err := d.client.Sites().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting sites",
			err.Error(),
		)
		return
	}

	for _, v := range result.GetValue() {
		value := new(sitesValueDataSourceModel)

		if v.GetId() != nil {
			value.Id = types.StringValue(*v.GetId())
		}
		if v.GetCreatedBy() != nil {
			value.CreatedBy = new(sitesCreatedByDataSourceModel)

			if v.GetCreatedBy().GetApplication() != nil {
				value.CreatedBy.Application = new(sitesApplicationDataSourceModel)

				if v.GetCreatedBy().GetApplication().GetDisplayName() != nil {
					value.CreatedBy.Application.DisplayName = types.StringValue(*v.GetCreatedBy().GetApplication().GetDisplayName())
				}
				if v.GetCreatedBy().GetApplication().GetId() != nil {
					value.CreatedBy.Application.Id = types.StringValue(*v.GetCreatedBy().GetApplication().GetId())
				}
			}
			if v.GetCreatedBy().GetDevice() != nil {
				value.CreatedBy.Device = new(sitesDeviceDataSourceModel)

				if v.GetCreatedBy().GetDevice().GetDisplayName() != nil {
					value.CreatedBy.Device.DisplayName = types.StringValue(*v.GetCreatedBy().GetDevice().GetDisplayName())
				}
				if v.GetCreatedBy().GetDevice().GetId() != nil {
					value.CreatedBy.Device.Id = types.StringValue(*v.GetCreatedBy().GetDevice().GetId())
				}
			}
			if v.GetCreatedBy().GetUser() != nil {
				value.CreatedBy.User = new(sitesUserDataSourceModel)

				if v.GetCreatedBy().GetUser().GetDisplayName() != nil {
					value.CreatedBy.User.DisplayName = types.StringValue(*v.GetCreatedBy().GetUser().GetDisplayName())
				}
				if v.GetCreatedBy().GetUser().GetId() != nil {
					value.CreatedBy.User.Id = types.StringValue(*v.GetCreatedBy().GetUser().GetId())
				}
			}
		}
		if v.GetCreatedDateTime() != nil {
			value.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
		}
		if v.GetDescription() != nil {
			value.Description = types.StringValue(*v.GetDescription())
		}
		if v.GetETag() != nil {
			value.ETag = types.StringValue(*v.GetETag())
		}
		if v.GetLastModifiedBy() != nil {
			value.LastModifiedBy = new(sitesLastModifiedByDataSourceModel)

			if v.GetLastModifiedBy().GetApplication() != nil {
				value.LastModifiedBy.Application = new(sitesApplicationDataSourceModel)

				if v.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
					value.LastModifiedBy.Application.DisplayName = types.StringValue(*v.GetLastModifiedBy().GetApplication().GetDisplayName())
				}
				if v.GetLastModifiedBy().GetApplication().GetId() != nil {
					value.LastModifiedBy.Application.Id = types.StringValue(*v.GetLastModifiedBy().GetApplication().GetId())
				}
			}
			if v.GetLastModifiedBy().GetDevice() != nil {
				value.LastModifiedBy.Device = new(sitesDeviceDataSourceModel)

				if v.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
					value.LastModifiedBy.Device.DisplayName = types.StringValue(*v.GetLastModifiedBy().GetDevice().GetDisplayName())
				}
				if v.GetLastModifiedBy().GetDevice().GetId() != nil {
					value.LastModifiedBy.Device.Id = types.StringValue(*v.GetLastModifiedBy().GetDevice().GetId())
				}
			}
			if v.GetLastModifiedBy().GetUser() != nil {
				value.LastModifiedBy.User = new(sitesUserDataSourceModel)

				if v.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
					value.LastModifiedBy.User.DisplayName = types.StringValue(*v.GetLastModifiedBy().GetUser().GetDisplayName())
				}
				if v.GetLastModifiedBy().GetUser().GetId() != nil {
					value.LastModifiedBy.User.Id = types.StringValue(*v.GetLastModifiedBy().GetUser().GetId())
				}
			}
		}
		if v.GetLastModifiedDateTime() != nil {
			value.LastModifiedDateTime = types.StringValue(v.GetLastModifiedDateTime().String())
		}
		if v.GetName() != nil {
			value.Name = types.StringValue(*v.GetName())
		}
		if v.GetParentReference() != nil {
			value.ParentReference = new(sitesParentReferenceDataSourceModel)

			if v.GetParentReference().GetDriveId() != nil {
				value.ParentReference.DriveId = types.StringValue(*v.GetParentReference().GetDriveId())
			}
			if v.GetParentReference().GetDriveType() != nil {
				value.ParentReference.DriveType = types.StringValue(*v.GetParentReference().GetDriveType())
			}
			if v.GetParentReference().GetId() != nil {
				value.ParentReference.Id = types.StringValue(*v.GetParentReference().GetId())
			}
			if v.GetParentReference().GetName() != nil {
				value.ParentReference.Name = types.StringValue(*v.GetParentReference().GetName())
			}
			if v.GetParentReference().GetPath() != nil {
				value.ParentReference.Path = types.StringValue(*v.GetParentReference().GetPath())
			}
			if v.GetParentReference().GetShareId() != nil {
				value.ParentReference.ShareId = types.StringValue(*v.GetParentReference().GetShareId())
			}
			if v.GetParentReference().GetSharepointIds() != nil {
				value.ParentReference.SharepointIds = new(sitesSharepointIdsDataSourceModel)

				if v.GetParentReference().GetSharepointIds().GetListId() != nil {
					value.ParentReference.SharepointIds.ListId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetListId())
				}
				if v.GetParentReference().GetSharepointIds().GetListItemId() != nil {
					value.ParentReference.SharepointIds.ListItemId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetListItemId())
				}
				if v.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
					value.ParentReference.SharepointIds.ListItemUniqueId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetListItemUniqueId())
				}
				if v.GetParentReference().GetSharepointIds().GetSiteId() != nil {
					value.ParentReference.SharepointIds.SiteId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetSiteId())
				}
				if v.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
					value.ParentReference.SharepointIds.SiteUrl = types.StringValue(*v.GetParentReference().GetSharepointIds().GetSiteUrl())
				}
				if v.GetParentReference().GetSharepointIds().GetTenantId() != nil {
					value.ParentReference.SharepointIds.TenantId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetTenantId())
				}
				if v.GetParentReference().GetSharepointIds().GetWebId() != nil {
					value.ParentReference.SharepointIds.WebId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetWebId())
				}
			}
			if v.GetParentReference().GetSiteId() != nil {
				value.ParentReference.SiteId = types.StringValue(*v.GetParentReference().GetSiteId())
			}
		}
		if v.GetWebUrl() != nil {
			value.WebUrl = types.StringValue(*v.GetWebUrl())
		}
		if v.GetDisplayName() != nil {
			value.DisplayName = types.StringValue(*v.GetDisplayName())
		}
		if v.GetError() != nil {
			value.Error = new(sitesErrorDataSourceModel)

			if v.GetError().GetCode() != nil {
				value.Error.Code = types.StringValue(*v.GetError().GetCode())
			}
			for _, v := range v.GetError().GetDetails() {
				details := new(sitesDetailsDataSourceModel)

				if v.GetCode() != nil {
					details.Code = types.StringValue(*v.GetCode())
				}
				if v.GetMessage() != nil {
					details.Message = types.StringValue(*v.GetMessage())
				}
				if v.GetTarget() != nil {
					details.Target = types.StringValue(*v.GetTarget())
				}
				value.Error.Details = append(value.Error.Details, *details)
			}
			if v.GetError().GetInnerError() != nil {
				value.Error.InnerError = new(sitesInnerErrorDataSourceModel)

				if v.GetError().GetInnerError().GetCode() != nil {
					value.Error.InnerError.Code = types.StringValue(*v.GetError().GetInnerError().GetCode())
				}
				for _, v := range v.GetError().GetInnerError().GetDetails() {
					details := new(sitesDetailsDataSourceModel)

					if v.GetCode() != nil {
						details.Code = types.StringValue(*v.GetCode())
					}
					if v.GetMessage() != nil {
						details.Message = types.StringValue(*v.GetMessage())
					}
					if v.GetTarget() != nil {
						details.Target = types.StringValue(*v.GetTarget())
					}
					value.Error.InnerError.Details = append(value.Error.InnerError.Details, *details)
				}
				if v.GetError().GetInnerError().GetMessage() != nil {
					value.Error.InnerError.Message = types.StringValue(*v.GetError().GetInnerError().GetMessage())
				}
				if v.GetError().GetInnerError().GetTarget() != nil {
					value.Error.InnerError.Target = types.StringValue(*v.GetError().GetInnerError().GetTarget())
				}
			}
			if v.GetError().GetMessage() != nil {
				value.Error.Message = types.StringValue(*v.GetError().GetMessage())
			}
			if v.GetError().GetTarget() != nil {
				value.Error.Target = types.StringValue(*v.GetError().GetTarget())
			}
		}
		if v.GetIsPersonalSite() != nil {
			value.IsPersonalSite = types.BoolValue(*v.GetIsPersonalSite())
		}
		if v.GetRoot() != nil {
			value.Root = new(sitesRootDataSourceModel)

		}
		if v.GetSharepointIds() != nil {
			value.SharepointIds = new(sitesSharepointIdsDataSourceModel)

			if v.GetSharepointIds().GetListId() != nil {
				value.SharepointIds.ListId = types.StringValue(*v.GetSharepointIds().GetListId())
			}
			if v.GetSharepointIds().GetListItemId() != nil {
				value.SharepointIds.ListItemId = types.StringValue(*v.GetSharepointIds().GetListItemId())
			}
			if v.GetSharepointIds().GetListItemUniqueId() != nil {
				value.SharepointIds.ListItemUniqueId = types.StringValue(*v.GetSharepointIds().GetListItemUniqueId())
			}
			if v.GetSharepointIds().GetSiteId() != nil {
				value.SharepointIds.SiteId = types.StringValue(*v.GetSharepointIds().GetSiteId())
			}
			if v.GetSharepointIds().GetSiteUrl() != nil {
				value.SharepointIds.SiteUrl = types.StringValue(*v.GetSharepointIds().GetSiteUrl())
			}
			if v.GetSharepointIds().GetTenantId() != nil {
				value.SharepointIds.TenantId = types.StringValue(*v.GetSharepointIds().GetTenantId())
			}
			if v.GetSharepointIds().GetWebId() != nil {
				value.SharepointIds.WebId = types.StringValue(*v.GetSharepointIds().GetWebId())
			}
		}
		if v.GetSiteCollection() != nil {
			value.SiteCollection = new(sitesSiteCollectionDataSourceModel)

			if v.GetSiteCollection().GetDataLocationCode() != nil {
				value.SiteCollection.DataLocationCode = types.StringValue(*v.GetSiteCollection().GetDataLocationCode())
			}
			if v.GetSiteCollection().GetHostname() != nil {
				value.SiteCollection.Hostname = types.StringValue(*v.GetSiteCollection().GetHostname())
			}
			if v.GetSiteCollection().GetRoot() != nil {
				value.SiteCollection.Root = new(sitesRootDataSourceModel)

			}
		}
		state.Value = append(state.Value, *value)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
