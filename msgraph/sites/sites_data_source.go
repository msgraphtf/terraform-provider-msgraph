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

	result, err := d.client.Sites().Get(context.Background(), &qparams)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting sites",
			err.Error(),
		)
		return
	}

	if len(result.GetValue()) > 0 {
		objectValues := []basetypes.ObjectValue{}
		for _, resultValue := range result.GetValue() {
			tfStateSite := sitesSiteModel{}

			if resultValue.GetId() != nil {
				tfStateSite.Id = types.StringValue(*resultValue.GetId())
			} else {
				tfStateSite.Id = types.StringNull()
			}
			if resultValue.GetCreatedBy() != nil {
				tfStateIdentitySet := sitesIdentitySetModel{}

				if resultValue.GetCreatedBy().GetApplication() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if resultValue.GetCreatedBy().GetApplication().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*resultValue.GetCreatedBy().GetApplication().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if resultValue.GetCreatedBy().GetApplication().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*resultValue.GetCreatedBy().GetApplication().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if resultValue.GetCreatedBy().GetDevice() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if resultValue.GetCreatedBy().GetDevice().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*resultValue.GetCreatedBy().GetDevice().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if resultValue.GetCreatedBy().GetDevice().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*resultValue.GetCreatedBy().GetDevice().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if resultValue.GetCreatedBy().GetUser() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if resultValue.GetCreatedBy().GetUser().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*resultValue.GetCreatedBy().GetUser().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if resultValue.GetCreatedBy().GetUser().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*resultValue.GetCreatedBy().GetUser().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}

				tfStateSite.CreatedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
			}
			if resultValue.GetCreatedDateTime() != nil {
				tfStateSite.CreatedDateTime = types.StringValue(resultValue.GetCreatedDateTime().String())
			} else {
				tfStateSite.CreatedDateTime = types.StringNull()
			}
			if resultValue.GetDescription() != nil {
				tfStateSite.Description = types.StringValue(*resultValue.GetDescription())
			} else {
				tfStateSite.Description = types.StringNull()
			}
			if resultValue.GetETag() != nil {
				tfStateSite.ETag = types.StringValue(*resultValue.GetETag())
			} else {
				tfStateSite.ETag = types.StringNull()
			}
			if resultValue.GetLastModifiedBy() != nil {
				tfStateIdentitySet := sitesIdentitySetModel{}

				if resultValue.GetLastModifiedBy().GetApplication() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if resultValue.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*resultValue.GetLastModifiedBy().GetApplication().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if resultValue.GetLastModifiedBy().GetApplication().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*resultValue.GetLastModifiedBy().GetApplication().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if resultValue.GetLastModifiedBy().GetDevice() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if resultValue.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*resultValue.GetLastModifiedBy().GetDevice().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if resultValue.GetLastModifiedBy().GetDevice().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*resultValue.GetLastModifiedBy().GetDevice().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}
				if resultValue.GetLastModifiedBy().GetUser() != nil {
					tfStateIdentity := sitesIdentityModel{}

					if resultValue.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
						tfStateIdentity.DisplayName = types.StringValue(*resultValue.GetLastModifiedBy().GetUser().GetDisplayName())
					} else {
						tfStateIdentity.DisplayName = types.StringNull()
					}
					if resultValue.GetLastModifiedBy().GetUser().GetId() != nil {
						tfStateIdentity.Id = types.StringValue(*resultValue.GetLastModifiedBy().GetUser().GetId())
					} else {
						tfStateIdentity.Id = types.StringNull()
					}

					tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
				}

				tfStateSite.LastModifiedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
			}
			if resultValue.GetLastModifiedDateTime() != nil {
				tfStateSite.LastModifiedDateTime = types.StringValue(resultValue.GetLastModifiedDateTime().String())
			} else {
				tfStateSite.LastModifiedDateTime = types.StringNull()
			}
			if resultValue.GetName() != nil {
				tfStateSite.Name = types.StringValue(*resultValue.GetName())
			} else {
				tfStateSite.Name = types.StringNull()
			}
			if resultValue.GetParentReference() != nil {
				tfStateItemReference := sitesItemReferenceModel{}

				if resultValue.GetParentReference().GetDriveId() != nil {
					tfStateItemReference.DriveId = types.StringValue(*resultValue.GetParentReference().GetDriveId())
				} else {
					tfStateItemReference.DriveId = types.StringNull()
				}
				if resultValue.GetParentReference().GetDriveType() != nil {
					tfStateItemReference.DriveType = types.StringValue(*resultValue.GetParentReference().GetDriveType())
				} else {
					tfStateItemReference.DriveType = types.StringNull()
				}
				if resultValue.GetParentReference().GetId() != nil {
					tfStateItemReference.Id = types.StringValue(*resultValue.GetParentReference().GetId())
				} else {
					tfStateItemReference.Id = types.StringNull()
				}
				if resultValue.GetParentReference().GetName() != nil {
					tfStateItemReference.Name = types.StringValue(*resultValue.GetParentReference().GetName())
				} else {
					tfStateItemReference.Name = types.StringNull()
				}
				if resultValue.GetParentReference().GetPath() != nil {
					tfStateItemReference.Path = types.StringValue(*resultValue.GetParentReference().GetPath())
				} else {
					tfStateItemReference.Path = types.StringNull()
				}
				if resultValue.GetParentReference().GetShareId() != nil {
					tfStateItemReference.ShareId = types.StringValue(*resultValue.GetParentReference().GetShareId())
				} else {
					tfStateItemReference.ShareId = types.StringNull()
				}
				if resultValue.GetParentReference().GetSharepointIds() != nil {
					tfStateSharepointIds := sitesSharepointIdsModel{}

					if resultValue.GetParentReference().GetSharepointIds().GetListId() != nil {
						tfStateSharepointIds.ListId = types.StringValue(*resultValue.GetParentReference().GetSharepointIds().GetListId())
					} else {
						tfStateSharepointIds.ListId = types.StringNull()
					}
					if resultValue.GetParentReference().GetSharepointIds().GetListItemId() != nil {
						tfStateSharepointIds.ListItemId = types.StringValue(*resultValue.GetParentReference().GetSharepointIds().GetListItemId())
					} else {
						tfStateSharepointIds.ListItemId = types.StringNull()
					}
					if resultValue.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
						tfStateSharepointIds.ListItemUniqueId = types.StringValue(*resultValue.GetParentReference().GetSharepointIds().GetListItemUniqueId())
					} else {
						tfStateSharepointIds.ListItemUniqueId = types.StringNull()
					}
					if resultValue.GetParentReference().GetSharepointIds().GetSiteId() != nil {
						tfStateSharepointIds.SiteId = types.StringValue(*resultValue.GetParentReference().GetSharepointIds().GetSiteId())
					} else {
						tfStateSharepointIds.SiteId = types.StringNull()
					}
					if resultValue.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
						tfStateSharepointIds.SiteUrl = types.StringValue(*resultValue.GetParentReference().GetSharepointIds().GetSiteUrl())
					} else {
						tfStateSharepointIds.SiteUrl = types.StringNull()
					}
					if resultValue.GetParentReference().GetSharepointIds().GetTenantId() != nil {
						tfStateSharepointIds.TenantId = types.StringValue(*resultValue.GetParentReference().GetSharepointIds().GetTenantId())
					} else {
						tfStateSharepointIds.TenantId = types.StringNull()
					}
					if resultValue.GetParentReference().GetSharepointIds().GetWebId() != nil {
						tfStateSharepointIds.WebId = types.StringValue(*resultValue.GetParentReference().GetSharepointIds().GetWebId())
					} else {
						tfStateSharepointIds.WebId = types.StringNull()
					}

					tfStateItemReference.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
				}
				if resultValue.GetParentReference().GetSiteId() != nil {
					tfStateItemReference.SiteId = types.StringValue(*resultValue.GetParentReference().GetSiteId())
				} else {
					tfStateItemReference.SiteId = types.StringNull()
				}

				tfStateSite.ParentReference, _ = types.ObjectValueFrom(ctx, tfStateItemReference.AttributeTypes(), tfStateItemReference)
			}
			if resultValue.GetWebUrl() != nil {
				tfStateSite.WebUrl = types.StringValue(*resultValue.GetWebUrl())
			} else {
				tfStateSite.WebUrl = types.StringNull()
			}
			if resultValue.GetDisplayName() != nil {
				tfStateSite.DisplayName = types.StringValue(*resultValue.GetDisplayName())
			} else {
				tfStateSite.DisplayName = types.StringNull()
			}
			if resultValue.GetError() != nil {
				tfStatePublicError := sitesPublicErrorModel{}

				if resultValue.GetError().GetCode() != nil {
					tfStatePublicError.Code = types.StringValue(*resultValue.GetError().GetCode())
				} else {
					tfStatePublicError.Code = types.StringNull()
				}
				if len(resultValue.GetError().GetDetails()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, resultDetails := range resultValue.GetError().GetDetails() {
						tfStatePublicErrorDetail := sitesPublicErrorDetailModel{}

						if resultDetails.GetCode() != nil {
							tfStatePublicErrorDetail.Code = types.StringValue(*resultDetails.GetCode())
						} else {
							tfStatePublicErrorDetail.Code = types.StringNull()
						}
						if resultDetails.GetMessage() != nil {
							tfStatePublicErrorDetail.Message = types.StringValue(*resultDetails.GetMessage())
						} else {
							tfStatePublicErrorDetail.Message = types.StringNull()
						}
						if resultDetails.GetTarget() != nil {
							tfStatePublicErrorDetail.Target = types.StringValue(*resultDetails.GetTarget())
						} else {
							tfStatePublicErrorDetail.Target = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, tfStatePublicErrorDetail.AttributeTypes(), tfStatePublicErrorDetail)
						objectValues = append(objectValues, objectValue)
					}
					tfStatePublicError.Details, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}
				if resultValue.GetError().GetInnerError() != nil {
					tfStatePublicInnerError := sitesPublicInnerErrorModel{}

					if resultValue.GetError().GetInnerError().GetCode() != nil {
						tfStatePublicInnerError.Code = types.StringValue(*resultValue.GetError().GetInnerError().GetCode())
					} else {
						tfStatePublicInnerError.Code = types.StringNull()
					}
					if len(resultValue.GetError().GetInnerError().GetDetails()) > 0 {
						objectValues := []basetypes.ObjectValue{}
						for _, resultDetails := range resultValue.GetError().GetInnerError().GetDetails() {
							tfStatePublicErrorDetail := sitesPublicErrorDetailModel{}

							if resultDetails.GetCode() != nil {
								tfStatePublicErrorDetail.Code = types.StringValue(*resultDetails.GetCode())
							} else {
								tfStatePublicErrorDetail.Code = types.StringNull()
							}
							if resultDetails.GetMessage() != nil {
								tfStatePublicErrorDetail.Message = types.StringValue(*resultDetails.GetMessage())
							} else {
								tfStatePublicErrorDetail.Message = types.StringNull()
							}
							if resultDetails.GetTarget() != nil {
								tfStatePublicErrorDetail.Target = types.StringValue(*resultDetails.GetTarget())
							} else {
								tfStatePublicErrorDetail.Target = types.StringNull()
							}
							objectValue, _ := types.ObjectValueFrom(ctx, tfStatePublicErrorDetail.AttributeTypes(), tfStatePublicErrorDetail)
							objectValues = append(objectValues, objectValue)
						}
						tfStatePublicInnerError.Details, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
					}
					if resultValue.GetError().GetInnerError().GetMessage() != nil {
						tfStatePublicInnerError.Message = types.StringValue(*resultValue.GetError().GetInnerError().GetMessage())
					} else {
						tfStatePublicInnerError.Message = types.StringNull()
					}
					if resultValue.GetError().GetInnerError().GetTarget() != nil {
						tfStatePublicInnerError.Target = types.StringValue(*resultValue.GetError().GetInnerError().GetTarget())
					} else {
						tfStatePublicInnerError.Target = types.StringNull()
					}

					tfStatePublicError.InnerError, _ = types.ObjectValueFrom(ctx, tfStatePublicInnerError.AttributeTypes(), tfStatePublicInnerError)
				}
				if resultValue.GetError().GetMessage() != nil {
					tfStatePublicError.Message = types.StringValue(*resultValue.GetError().GetMessage())
				} else {
					tfStatePublicError.Message = types.StringNull()
				}
				if resultValue.GetError().GetTarget() != nil {
					tfStatePublicError.Target = types.StringValue(*resultValue.GetError().GetTarget())
				} else {
					tfStatePublicError.Target = types.StringNull()
				}

				tfStateSite.Error, _ = types.ObjectValueFrom(ctx, tfStatePublicError.AttributeTypes(), tfStatePublicError)
			}
			if resultValue.GetIsPersonalSite() != nil {
				tfStateSite.IsPersonalSite = types.BoolValue(*resultValue.GetIsPersonalSite())
			} else {
				tfStateSite.IsPersonalSite = types.BoolNull()
			}
			if resultValue.GetRoot() != nil {
				tfStateRoot := sitesRootModel{}

				tfStateSite.Root, _ = types.ObjectValueFrom(ctx, tfStateRoot.AttributeTypes(), tfStateRoot)
			}
			if resultValue.GetSharepointIds() != nil {
				tfStateSharepointIds := sitesSharepointIdsModel{}

				if resultValue.GetSharepointIds().GetListId() != nil {
					tfStateSharepointIds.ListId = types.StringValue(*resultValue.GetSharepointIds().GetListId())
				} else {
					tfStateSharepointIds.ListId = types.StringNull()
				}
				if resultValue.GetSharepointIds().GetListItemId() != nil {
					tfStateSharepointIds.ListItemId = types.StringValue(*resultValue.GetSharepointIds().GetListItemId())
				} else {
					tfStateSharepointIds.ListItemId = types.StringNull()
				}
				if resultValue.GetSharepointIds().GetListItemUniqueId() != nil {
					tfStateSharepointIds.ListItemUniqueId = types.StringValue(*resultValue.GetSharepointIds().GetListItemUniqueId())
				} else {
					tfStateSharepointIds.ListItemUniqueId = types.StringNull()
				}
				if resultValue.GetSharepointIds().GetSiteId() != nil {
					tfStateSharepointIds.SiteId = types.StringValue(*resultValue.GetSharepointIds().GetSiteId())
				} else {
					tfStateSharepointIds.SiteId = types.StringNull()
				}
				if resultValue.GetSharepointIds().GetSiteUrl() != nil {
					tfStateSharepointIds.SiteUrl = types.StringValue(*resultValue.GetSharepointIds().GetSiteUrl())
				} else {
					tfStateSharepointIds.SiteUrl = types.StringNull()
				}
				if resultValue.GetSharepointIds().GetTenantId() != nil {
					tfStateSharepointIds.TenantId = types.StringValue(*resultValue.GetSharepointIds().GetTenantId())
				} else {
					tfStateSharepointIds.TenantId = types.StringNull()
				}
				if resultValue.GetSharepointIds().GetWebId() != nil {
					tfStateSharepointIds.WebId = types.StringValue(*resultValue.GetSharepointIds().GetWebId())
				} else {
					tfStateSharepointIds.WebId = types.StringNull()
				}

				tfStateSite.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
			}
			if resultValue.GetSiteCollection() != nil {
				tfStateSiteCollection := sitesSiteCollectionModel{}

				if resultValue.GetSiteCollection().GetArchivalDetails() != nil {
					tfStateSiteArchivalDetails := sitesSiteArchivalDetailsModel{}

					if resultValue.GetSiteCollection().GetArchivalDetails().GetArchiveStatus() != nil {
						tfStateSiteArchivalDetails.ArchiveStatus = types.StringValue(resultValue.GetSiteCollection().GetArchivalDetails().GetArchiveStatus().String())
					} else {
						tfStateSiteArchivalDetails.ArchiveStatus = types.StringNull()
					}

					tfStateSiteCollection.ArchivalDetails, _ = types.ObjectValueFrom(ctx, tfStateSiteArchivalDetails.AttributeTypes(), tfStateSiteArchivalDetails)
				}
				if resultValue.GetSiteCollection().GetDataLocationCode() != nil {
					tfStateSiteCollection.DataLocationCode = types.StringValue(*resultValue.GetSiteCollection().GetDataLocationCode())
				} else {
					tfStateSiteCollection.DataLocationCode = types.StringNull()
				}
				if resultValue.GetSiteCollection().GetHostname() != nil {
					tfStateSiteCollection.Hostname = types.StringValue(*resultValue.GetSiteCollection().GetHostname())
				} else {
					tfStateSiteCollection.Hostname = types.StringNull()
				}
				if resultValue.GetSiteCollection().GetRoot() != nil {
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
