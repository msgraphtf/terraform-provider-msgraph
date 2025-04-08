package sites

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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
											Description: "The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.",
											Computed:    true,
										},
									},
								},
								"device": schema.SingleNestedAttribute{
									Description: "Optional. The device associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.",
											Computed:    true,
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description: "Optional. The user associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.",
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
											Description: "The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.",
											Computed:    true,
										},
									},
								},
								"device": schema.SingleNestedAttribute{
									Description: "Optional. The device associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.",
											Computed:    true,
										},
									},
								},
								"user": schema.SingleNestedAttribute{
									Description: "Optional. The user associated with this action.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"display_name": schema.StringAttribute{
											Description: "The display name of the identity.For drive items, the display name might not always be available or up to date. For example, if a user changes their display name the API might show the new value in a future response, but the items associated with the user don't show up as changed when using delta.",
											Computed:    true,
										},
										"id": schema.StringAttribute{
											Description: "Unique identifier for the identity or actor. For example, in the access reviews decisions API, this property might record the id of the principal, that is, the group, user, or application that's subject to review.",
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
									Description: "Unique identifier of the drive instance that contains the driveItem. Only returned if the item is located in a drive. Read-only.",
									Computed:    true,
								},
								"drive_type": schema.StringAttribute{
									Description: "Identifies the type of drive. Only returned if the item is located in a drive. See drive resource for values.",
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
									Description: "Percent-encoded path that can be used to navigate to the item. Read-only.",
									Computed:    true,
								},
								"share_id": schema.StringAttribute{
									Description: "A unique identifier for a shared resource that can be accessed via the Shares API.",
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
									Description: "For OneDrive for Business and SharePoint, this property represents the ID of the site that contains the parent document library of the driveItem resource or the parent list of the listItem resource. The value is the same as the id property of that site resource. It is an opaque string that consists of three identifiers of the site. For OneDrive, this property is not populated.",
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
							Description: "Identifies whether the site is personal or not. Read-only.",
							Computed:    true,
						},
						"root": schema.SingleNestedAttribute{
							Description: "If present, provides the root site in the site collection. Read-only.",
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
								"archival_details": schema.SingleNestedAttribute{
									Description: "Represents whether the site collection is recently archived, fully archived, or reactivating. Possible values are: recentlyArchived, fullyArchived, reactivating, unknownFutureValue.",
									Computed:    true,
									Attributes: map[string]schema.Attribute{
										"archive_status": schema.StringAttribute{
											Description: "Represents the current archive status of the site collection. Returned only on $select. The possible values are: recentlyArchived, fullyArchived, reactivating, unknownFutureValue.",
											Computed:    true,
										},
									},
								},
								"data_location_code": schema.StringAttribute{
									Description: "The geographic region code for where this site collection resides. Only present for multi-geo tenants. Read-only.",
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
	var tfStateSites sitesModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateSites)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := sites.SitesRequestBuilderGetRequestConfiguration{
		QueryParameters: &sites.SitesRequestBuilderGetQueryParameters{
			Select: []string{
				"value",
			},
		},
	}

	response, err := d.client.Sites().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting sites",
			err.Error(),
		)
		return
	}

	if len(response.GetValue()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, responseValue := range response.GetValue() {
			tfStateSite := sitesSiteModel{}

			if responseValue.GetId() != nil {
				tfStateSite.Id = types.StringValue(*responseValue.GetId())
			} else {
				tfStateSite.Id = types.StringNull()
			}
			if responseValue.GetCreatedBy() != nil {
				tfStateIdentitySet := sitesIdentitySetModel{}

				if responseValue.GetCreatedBy().GetApplication() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if responseValue.GetCreatedBy().GetApplication().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*responseValue.GetCreatedBy().GetApplication().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if responseValue.GetCreatedBy().GetApplication().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*responseValue.GetCreatedBy().GetApplication().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if responseValue.GetCreatedBy().GetDevice() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if responseValue.GetCreatedBy().GetDevice().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*responseValue.GetCreatedBy().GetDevice().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if responseValue.GetCreatedBy().GetDevice().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*responseValue.GetCreatedBy().GetDevice().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if responseValue.GetCreatedBy().GetUser() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if responseValue.GetCreatedBy().GetUser().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*responseValue.GetCreatedBy().GetUser().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if responseValue.GetCreatedBy().GetUser().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*responseValue.GetCreatedBy().GetUser().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}

				tfStateSite.CreatedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
			}
			if responseValue.GetCreatedDateTime() != nil {
				tfStateSite.CreatedDateTime = types.StringValue(responseValue.GetCreatedDateTime().String())
			} else {
				tfStateSite.CreatedDateTime = types.StringNull()
			}
			if responseValue.GetDescription() != nil {
				tfStateSite.Description = types.StringValue(*responseValue.GetDescription())
			} else {
				tfStateSite.Description = types.StringNull()
			}
			if responseValue.GetETag() != nil {
				tfStateSite.ETag = types.StringValue(*responseValue.GetETag())
			} else {
				tfStateSite.ETag = types.StringNull()
			}
			if responseValue.GetLastModifiedBy() != nil {
				tfStateIdentitySet := sitesIdentitySetModel{}

				if responseValue.GetLastModifiedBy().GetApplication() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if responseValue.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*responseValue.GetLastModifiedBy().GetApplication().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if responseValue.GetLastModifiedBy().GetApplication().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*responseValue.GetLastModifiedBy().GetApplication().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if responseValue.GetLastModifiedBy().GetDevice() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if responseValue.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*responseValue.GetLastModifiedBy().GetDevice().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if responseValue.GetLastModifiedBy().GetDevice().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*responseValue.GetLastModifiedBy().GetDevice().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if responseValue.GetLastModifiedBy().GetUser() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if responseValue.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*responseValue.GetLastModifiedBy().GetUser().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if responseValue.GetLastModifiedBy().GetUser().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*responseValue.GetLastModifiedBy().GetUser().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}

				tfStateSite.LastModifiedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
			}
			if responseValue.GetLastModifiedDateTime() != nil {
				tfStateSite.LastModifiedDateTime = types.StringValue(responseValue.GetLastModifiedDateTime().String())
			} else {
				tfStateSite.LastModifiedDateTime = types.StringNull()
			}
			if responseValue.GetName() != nil {
				tfStateSite.Name = types.StringValue(*responseValue.GetName())
			} else {
				tfStateSite.Name = types.StringNull()
			}
			if responseValue.GetParentReference() != nil {
				tfStateItemReference := sitesItemReferenceModel{}

				if responseValue.GetParentReference().GetDriveId() != nil {
					tfStateItemReference.DriveId = types.StringValue(*responseValue.GetParentReference().GetDriveId())
				} else {
					tfStateItemReference.DriveId = types.StringNull()
				}
				if responseValue.GetParentReference().GetDriveType() != nil {
					tfStateItemReference.DriveType = types.StringValue(*responseValue.GetParentReference().GetDriveType())
				} else {
					tfStateItemReference.DriveType = types.StringNull()
				}
				if responseValue.GetParentReference().GetId() != nil {
					tfStateItemReference.Id = types.StringValue(*responseValue.GetParentReference().GetId())
				} else {
					tfStateItemReference.Id = types.StringNull()
				}
				if responseValue.GetParentReference().GetName() != nil {
					tfStateItemReference.Name = types.StringValue(*responseValue.GetParentReference().GetName())
				} else {
					tfStateItemReference.Name = types.StringNull()
				}
				if responseValue.GetParentReference().GetPath() != nil {
					tfStateItemReference.Path = types.StringValue(*responseValue.GetParentReference().GetPath())
				} else {
					tfStateItemReference.Path = types.StringNull()
				}
				if responseValue.GetParentReference().GetShareId() != nil {
					tfStateItemReference.ShareId = types.StringValue(*responseValue.GetParentReference().GetShareId())
				} else {
					tfStateItemReference.ShareId = types.StringNull()
				}
				if responseValue.GetParentReference().GetSharepointIds() != nil {
					tfStateSharepointIds := sitesSharepointIdsModel{}

					if responseValue.GetParentReference().GetSharepointIds().GetListId() != nil {
						tfStateSharepointIds.ListId = types.StringValue(*responseValue.GetParentReference().GetSharepointIds().GetListId())
					} else {
						tfStateSharepointIds.ListId = types.StringNull()
					}
					if responseValue.GetParentReference().GetSharepointIds().GetListItemId() != nil {
						tfStateSharepointIds.ListItemId = types.StringValue(*responseValue.GetParentReference().GetSharepointIds().GetListItemId())
					} else {
						tfStateSharepointIds.ListItemId = types.StringNull()
					}
					if responseValue.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
						tfStateSharepointIds.ListItemUniqueId = types.StringValue(*responseValue.GetParentReference().GetSharepointIds().GetListItemUniqueId())
					} else {
						tfStateSharepointIds.ListItemUniqueId = types.StringNull()
					}
					if responseValue.GetParentReference().GetSharepointIds().GetSiteId() != nil {
						tfStateSharepointIds.SiteId = types.StringValue(*responseValue.GetParentReference().GetSharepointIds().GetSiteId())
					} else {
						tfStateSharepointIds.SiteId = types.StringNull()
					}
					if responseValue.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
						tfStateSharepointIds.SiteUrl = types.StringValue(*responseValue.GetParentReference().GetSharepointIds().GetSiteUrl())
					} else {
						tfStateSharepointIds.SiteUrl = types.StringNull()
					}
					if responseValue.GetParentReference().GetSharepointIds().GetTenantId() != nil {
						tfStateSharepointIds.TenantId = types.StringValue(*responseValue.GetParentReference().GetSharepointIds().GetTenantId())
					} else {
						tfStateSharepointIds.TenantId = types.StringNull()
					}
					if responseValue.GetParentReference().GetSharepointIds().GetWebId() != nil {
						tfStateSharepointIds.WebId = types.StringValue(*responseValue.GetParentReference().GetSharepointIds().GetWebId())
					} else {
						tfStateSharepointIds.WebId = types.StringNull()
					}

					tfStateItemReference.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
				}
				if responseValue.GetParentReference().GetSiteId() != nil {
					tfStateItemReference.SiteId = types.StringValue(*responseValue.GetParentReference().GetSiteId())
				} else {
					tfStateItemReference.SiteId = types.StringNull()
				}

				tfStateSite.ParentReference, _ = types.ObjectValueFrom(ctx, tfStateItemReference.AttributeTypes(), tfStateItemReference)
			}
			if responseValue.GetWebUrl() != nil {
				tfStateSite.WebUrl = types.StringValue(*responseValue.GetWebUrl())
			} else {
				tfStateSite.WebUrl = types.StringNull()
			}
			if responseValue.GetDisplayName() != nil {
				tfStateSite.DisplayName = types.StringValue(*responseValue.GetDisplayName())
			} else {
				tfStateSite.DisplayName = types.StringNull()
			}
			if responseValue.GetError() != nil {
				tfStatePublicError := sitesPublicErrorModel{}

				if responseValue.GetError().GetCode() != nil {
					tfStatePublicError.Code = types.StringValue(*responseValue.GetError().GetCode())
				} else {
					tfStatePublicError.Code = types.StringNull()
				}
				if len(responseValue.GetError().GetDetails()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, responseDetails := range responseValue.GetError().GetDetails() {
						tfStatePublicErrorDetail := sitesPublicErrorDetailModel{}

						if responseDetails.GetCode() != nil {
							tfStatePublicErrorDetail.Code = types.StringValue(*responseDetails.GetCode())
						} else {
							tfStatePublicErrorDetail.Code = types.StringNull()
						}
						if responseDetails.GetMessage() != nil {
							tfStatePublicErrorDetail.Message = types.StringValue(*responseDetails.GetMessage())
						} else {
							tfStatePublicErrorDetail.Message = types.StringNull()
						}
						if responseDetails.GetTarget() != nil {
							tfStatePublicErrorDetail.Target = types.StringValue(*responseDetails.GetTarget())
						} else {
							tfStatePublicErrorDetail.Target = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStatePublicErrorDetail.AttributeTypes(), tfStatePublicErrorDetail)
						objectValues = append(objectValues, objectValue)
					}
					tfStatePublicError.Details, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}
				if responseValue.GetError().GetInnerError() != nil {
					tfStatePublicInnerError := sitesPublicInnerErrorModel{}

					if responseValue.GetError().GetInnerError().GetCode() != nil {
						tfStatePublicInnerError.Code = types.StringValue(*responseValue.GetError().GetInnerError().GetCode())
					} else {
						tfStatePublicInnerError.Code = types.StringNull()
					}
					if len(responseValue.GetError().GetInnerError().GetDetails()) > 0 {
						objectValues := []basetypes.ObjectValue{}
						for _, responseDetails := range responseValue.GetError().GetInnerError().GetDetails() {
							tfStatePublicErrorDetail := sitesPublicErrorDetailModel{}

							if responseDetails.GetCode() != nil {
								tfStatePublicErrorDetail.Code = types.StringValue(*responseDetails.GetCode())
							} else {
								tfStatePublicErrorDetail.Code = types.StringNull()
							}
							if responseDetails.GetMessage() != nil {
								tfStatePublicErrorDetail.Message = types.StringValue(*responseDetails.GetMessage())
							} else {
								tfStatePublicErrorDetail.Message = types.StringNull()
							}
							if responseDetails.GetTarget() != nil {
								tfStatePublicErrorDetail.Target = types.StringValue(*responseDetails.GetTarget())
							} else {
								tfStatePublicErrorDetail.Target = types.StringNull()
							}
							objectValue, _ := types.ObjectValueFrom(ctx, tfStatePublicErrorDetail.AttributeTypes(), tfStatePublicErrorDetail)
							objectValues = append(objectValues, objectValue)
						}
						tfStatePublicInnerError.Details, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
					}
					if responseValue.GetError().GetInnerError().GetMessage() != nil {
						tfStatePublicInnerError.Message = types.StringValue(*responseValue.GetError().GetInnerError().GetMessage())
					} else {
						tfStatePublicInnerError.Message = types.StringNull()
					}
					if responseValue.GetError().GetInnerError().GetTarget() != nil {
						tfStatePublicInnerError.Target = types.StringValue(*responseValue.GetError().GetInnerError().GetTarget())
					} else {
						tfStatePublicInnerError.Target = types.StringNull()
					}

					tfStatePublicError.InnerError, _ = types.ObjectValueFrom(ctx, tfStatePublicInnerError.AttributeTypes(), tfStatePublicInnerError)
				}
				if responseValue.GetError().GetMessage() != nil {
					tfStatePublicError.Message = types.StringValue(*responseValue.GetError().GetMessage())
				} else {
					tfStatePublicError.Message = types.StringNull()
				}
				if responseValue.GetError().GetTarget() != nil {
					tfStatePublicError.Target = types.StringValue(*responseValue.GetError().GetTarget())
				} else {
					tfStatePublicError.Target = types.StringNull()
				}

				tfStateSite.Error, _ = types.ObjectValueFrom(ctx, tfStatePublicError.AttributeTypes(), tfStatePublicError)
			}
			if responseValue.GetIsPersonalSite() != nil {
				tfStateSite.IsPersonalSite = types.BoolValue(*responseValue.GetIsPersonalSite())
			} else {
				tfStateSite.IsPersonalSite = types.BoolNull()
			}
			if responseValue.GetRoot() != nil {
				tfStateRoot := sitesRootModel{}

				tfStateSite.Root, _ = types.ObjectValueFrom(ctx, tfStateRoot.AttributeTypes(), tfStateRoot)
			}
			if responseValue.GetSharepointIds() != nil {
				tfStateSharepointIds := sitesSharepointIdsModel{}

				if responseValue.GetSharepointIds().GetListId() != nil {
					tfStateSharepointIds.ListId = types.StringValue(*responseValue.GetSharepointIds().GetListId())
				} else {
					tfStateSharepointIds.ListId = types.StringNull()
				}
				if responseValue.GetSharepointIds().GetListItemId() != nil {
					tfStateSharepointIds.ListItemId = types.StringValue(*responseValue.GetSharepointIds().GetListItemId())
				} else {
					tfStateSharepointIds.ListItemId = types.StringNull()
				}
				if responseValue.GetSharepointIds().GetListItemUniqueId() != nil {
					tfStateSharepointIds.ListItemUniqueId = types.StringValue(*responseValue.GetSharepointIds().GetListItemUniqueId())
				} else {
					tfStateSharepointIds.ListItemUniqueId = types.StringNull()
				}
				if responseValue.GetSharepointIds().GetSiteId() != nil {
					tfStateSharepointIds.SiteId = types.StringValue(*responseValue.GetSharepointIds().GetSiteId())
				} else {
					tfStateSharepointIds.SiteId = types.StringNull()
				}
				if responseValue.GetSharepointIds().GetSiteUrl() != nil {
					tfStateSharepointIds.SiteUrl = types.StringValue(*responseValue.GetSharepointIds().GetSiteUrl())
				} else {
					tfStateSharepointIds.SiteUrl = types.StringNull()
				}
				if responseValue.GetSharepointIds().GetTenantId() != nil {
					tfStateSharepointIds.TenantId = types.StringValue(*responseValue.GetSharepointIds().GetTenantId())
				} else {
					tfStateSharepointIds.TenantId = types.StringNull()
				}
				if responseValue.GetSharepointIds().GetWebId() != nil {
					tfStateSharepointIds.WebId = types.StringValue(*responseValue.GetSharepointIds().GetWebId())
				} else {
					tfStateSharepointIds.WebId = types.StringNull()
				}

				tfStateSite.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
			}
			if responseValue.GetSiteCollection() != nil {
				tfStateSiteCollection := sitesSiteCollectionModel{}

				if responseValue.GetSiteCollection().GetArchivalDetails() != nil {
					tfStateSiteArchivalDetails := sitesSiteArchivalDetailsModel{}

					if responseValue.GetSiteCollection().GetArchivalDetails().GetArchiveStatus() != nil {
						tfStateSiteArchivalDetails.ArchiveStatus = types.StringValue(responseValue.GetSiteCollection().GetArchivalDetails().GetArchiveStatus().String())
					} else {
						tfStateSiteArchivalDetails.ArchiveStatus = types.StringNull()
					}

					tfStateSiteCollection.ArchivalDetails, _ = types.ObjectValueFrom(ctx, tfStateSiteArchivalDetails.AttributeTypes(), tfStateSiteArchivalDetails)
				}
				if responseValue.GetSiteCollection().GetDataLocationCode() != nil {
					tfStateSiteCollection.DataLocationCode = types.StringValue(*responseValue.GetSiteCollection().GetDataLocationCode())
				} else {
					tfStateSiteCollection.DataLocationCode = types.StringNull()
				}
				if responseValue.GetSiteCollection().GetHostname() != nil {
					tfStateSiteCollection.Hostname = types.StringValue(*responseValue.GetSiteCollection().GetHostname())
				} else {
					tfStateSiteCollection.Hostname = types.StringNull()
				}
				if responseValue.GetSiteCollection().GetRoot() != nil {
					tfStateRoot := sitesRootModel{}

					tfStateSiteCollection.Root, _ = types.ObjectValueFrom(ctx, tfStateRoot.AttributeTypes(), tfStateRoot)
				}

				tfStateSite.SiteCollection, _ = types.ObjectValueFrom(ctx, tfStateSiteCollection.AttributeTypes(), tfStateSiteCollection)
			}
			objectValue, _ := types.ObjectValueFrom(ctx, tfStateSite.AttributeTypes(), tfStateSite)
			objectValues = append(objectValues, objectValue)
		}
		tfStateSites.Value, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateSites)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
