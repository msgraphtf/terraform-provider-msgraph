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
		for _, v := range result.GetValue() {
			value := new(sitesSiteModel)

			if v.GetId() != nil {
				value.Id = types.StringValue(*v.GetId())
			} else {
				value.Id = types.StringNull()
			}
			if v.GetCreatedBy() != nil {
				createdBy := new(sitesIdentitySetModel)

				if v.GetCreatedBy().GetApplication() != nil {
					application := new(sitesIdentityModel)

					if v.GetCreatedBy().GetApplication().GetDisplayName() != nil {
						application.DisplayName = types.StringValue(*v.GetCreatedBy().GetApplication().GetDisplayName())
					} else {
						application.DisplayName = types.StringNull()
					}
					if v.GetCreatedBy().GetApplication().GetId() != nil {
						application.Id = types.StringValue(*v.GetCreatedBy().GetApplication().GetId())
					} else {
						application.Id = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, application.AttributeTypes(), application)
					createdBy.Application = objectValue
				}
				if v.GetCreatedBy().GetDevice() != nil {
					device := new(sitesIdentityModel)

					if v.GetCreatedBy().GetDevice().GetDisplayName() != nil {
						device.DisplayName = types.StringValue(*v.GetCreatedBy().GetDevice().GetDisplayName())
					} else {
						device.DisplayName = types.StringNull()
					}
					if v.GetCreatedBy().GetDevice().GetId() != nil {
						device.Id = types.StringValue(*v.GetCreatedBy().GetDevice().GetId())
					} else {
						device.Id = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, device.AttributeTypes(), device)
					createdBy.Device = objectValue
				}
				if v.GetCreatedBy().GetUser() != nil {
					user := new(sitesIdentityModel)

					if v.GetCreatedBy().GetUser().GetDisplayName() != nil {
						user.DisplayName = types.StringValue(*v.GetCreatedBy().GetUser().GetDisplayName())
					} else {
						user.DisplayName = types.StringNull()
					}
					if v.GetCreatedBy().GetUser().GetId() != nil {
						user.Id = types.StringValue(*v.GetCreatedBy().GetUser().GetId())
					} else {
						user.Id = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, user.AttributeTypes(), user)
					createdBy.User = objectValue
				}

				objectValue, _ := types.ObjectValueFrom(ctx, createdBy.AttributeTypes(), createdBy)
				value.CreatedBy = objectValue
			}
			if v.GetCreatedDateTime() != nil {
				value.CreatedDateTime = types.StringValue(v.GetCreatedDateTime().String())
			} else {
				value.CreatedDateTime = types.StringNull()
			}
			if v.GetDescription() != nil {
				value.Description = types.StringValue(*v.GetDescription())
			} else {
				value.Description = types.StringNull()
			}
			if v.GetETag() != nil {
				value.ETag = types.StringValue(*v.GetETag())
			} else {
				value.ETag = types.StringNull()
			}
			if v.GetLastModifiedBy() != nil {
				lastModifiedBy := new(sitesIdentitySetModel)

				if v.GetLastModifiedBy().GetApplication() != nil {
					application := new(sitesIdentityModel)

					if v.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
						application.DisplayName = types.StringValue(*v.GetLastModifiedBy().GetApplication().GetDisplayName())
					} else {
						application.DisplayName = types.StringNull()
					}
					if v.GetLastModifiedBy().GetApplication().GetId() != nil {
						application.Id = types.StringValue(*v.GetLastModifiedBy().GetApplication().GetId())
					} else {
						application.Id = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, application.AttributeTypes(), application)
					lastModifiedBy.Application = objectValue
				}
				if v.GetLastModifiedBy().GetDevice() != nil {
					device := new(sitesIdentityModel)

					if v.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
						device.DisplayName = types.StringValue(*v.GetLastModifiedBy().GetDevice().GetDisplayName())
					} else {
						device.DisplayName = types.StringNull()
					}
					if v.GetLastModifiedBy().GetDevice().GetId() != nil {
						device.Id = types.StringValue(*v.GetLastModifiedBy().GetDevice().GetId())
					} else {
						device.Id = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, device.AttributeTypes(), device)
					lastModifiedBy.Device = objectValue
				}
				if v.GetLastModifiedBy().GetUser() != nil {
					user := new(sitesIdentityModel)

					if v.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
						user.DisplayName = types.StringValue(*v.GetLastModifiedBy().GetUser().GetDisplayName())
					} else {
						user.DisplayName = types.StringNull()
					}
					if v.GetLastModifiedBy().GetUser().GetId() != nil {
						user.Id = types.StringValue(*v.GetLastModifiedBy().GetUser().GetId())
					} else {
						user.Id = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, user.AttributeTypes(), user)
					lastModifiedBy.User = objectValue
				}

				objectValue, _ := types.ObjectValueFrom(ctx, lastModifiedBy.AttributeTypes(), lastModifiedBy)
				value.LastModifiedBy = objectValue
			}
			if v.GetLastModifiedDateTime() != nil {
				value.LastModifiedDateTime = types.StringValue(v.GetLastModifiedDateTime().String())
			} else {
				value.LastModifiedDateTime = types.StringNull()
			}
			if v.GetName() != nil {
				value.Name = types.StringValue(*v.GetName())
			} else {
				value.Name = types.StringNull()
			}
			if v.GetParentReference() != nil {
				parentReference := new(sitesItemReferenceModel)

				if v.GetParentReference().GetDriveId() != nil {
					parentReference.DriveId = types.StringValue(*v.GetParentReference().GetDriveId())
				} else {
					parentReference.DriveId = types.StringNull()
				}
				if v.GetParentReference().GetDriveType() != nil {
					parentReference.DriveType = types.StringValue(*v.GetParentReference().GetDriveType())
				} else {
					parentReference.DriveType = types.StringNull()
				}
				if v.GetParentReference().GetId() != nil {
					parentReference.Id = types.StringValue(*v.GetParentReference().GetId())
				} else {
					parentReference.Id = types.StringNull()
				}
				if v.GetParentReference().GetName() != nil {
					parentReference.Name = types.StringValue(*v.GetParentReference().GetName())
				} else {
					parentReference.Name = types.StringNull()
				}
				if v.GetParentReference().GetPath() != nil {
					parentReference.Path = types.StringValue(*v.GetParentReference().GetPath())
				} else {
					parentReference.Path = types.StringNull()
				}
				if v.GetParentReference().GetShareId() != nil {
					parentReference.ShareId = types.StringValue(*v.GetParentReference().GetShareId())
				} else {
					parentReference.ShareId = types.StringNull()
				}
				if v.GetParentReference().GetSharepointIds() != nil {
					sharepointIds := new(sitesSharepointIdsModel)

					if v.GetParentReference().GetSharepointIds().GetListId() != nil {
						sharepointIds.ListId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetListId())
					} else {
						sharepointIds.ListId = types.StringNull()
					}
					if v.GetParentReference().GetSharepointIds().GetListItemId() != nil {
						sharepointIds.ListItemId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetListItemId())
					} else {
						sharepointIds.ListItemId = types.StringNull()
					}
					if v.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
						sharepointIds.ListItemUniqueId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetListItemUniqueId())
					} else {
						sharepointIds.ListItemUniqueId = types.StringNull()
					}
					if v.GetParentReference().GetSharepointIds().GetSiteId() != nil {
						sharepointIds.SiteId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetSiteId())
					} else {
						sharepointIds.SiteId = types.StringNull()
					}
					if v.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
						sharepointIds.SiteUrl = types.StringValue(*v.GetParentReference().GetSharepointIds().GetSiteUrl())
					} else {
						sharepointIds.SiteUrl = types.StringNull()
					}
					if v.GetParentReference().GetSharepointIds().GetTenantId() != nil {
						sharepointIds.TenantId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetTenantId())
					} else {
						sharepointIds.TenantId = types.StringNull()
					}
					if v.GetParentReference().GetSharepointIds().GetWebId() != nil {
						sharepointIds.WebId = types.StringValue(*v.GetParentReference().GetSharepointIds().GetWebId())
					} else {
						sharepointIds.WebId = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, sharepointIds.AttributeTypes(), sharepointIds)
					parentReference.SharepointIds = objectValue
				}
				if v.GetParentReference().GetSiteId() != nil {
					parentReference.SiteId = types.StringValue(*v.GetParentReference().GetSiteId())
				} else {
					parentReference.SiteId = types.StringNull()
				}

				objectValue, _ := types.ObjectValueFrom(ctx, parentReference.AttributeTypes(), parentReference)
				value.ParentReference = objectValue
			}
			if v.GetWebUrl() != nil {
				value.WebUrl = types.StringValue(*v.GetWebUrl())
			} else {
				value.WebUrl = types.StringNull()
			}
			if v.GetDisplayName() != nil {
				value.DisplayName = types.StringValue(*v.GetDisplayName())
			} else {
				value.DisplayName = types.StringNull()
			}
			if v.GetError() != nil {
				error := new(sitesPublicErrorModel)

				if v.GetError().GetCode() != nil {
					error.Code = types.StringValue(*v.GetError().GetCode())
				} else {
					error.Code = types.StringNull()
				}
				if len(v.GetError().GetDetails()) > 0 {
					objectValues := []basetypes.ObjectValue{}
					for _, v := range v.GetError().GetDetails() {
						details := new(sitesPublicErrorDetailModel)

						if v.GetCode() != nil {
							details.Code = types.StringValue(*v.GetCode())
						} else {
							details.Code = types.StringNull()
						}
						if v.GetMessage() != nil {
							details.Message = types.StringValue(*v.GetMessage())
						} else {
							details.Message = types.StringNull()
						}
						if v.GetTarget() != nil {
							details.Target = types.StringValue(*v.GetTarget())
						} else {
							details.Target = types.StringNull()
						}
						objectValue, _ := types.ObjectValueFrom(ctx, details.AttributeTypes(), details)
						objectValues = append(objectValues, objectValue)
					}
					error.Details, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
				}
				if v.GetError().GetInnerError() != nil {
					innerError := new(sitesPublicInnerErrorModel)

					if v.GetError().GetInnerError().GetCode() != nil {
						innerError.Code = types.StringValue(*v.GetError().GetInnerError().GetCode())
					} else {
						innerError.Code = types.StringNull()
					}
					if len(v.GetError().GetInnerError().GetDetails()) > 0 {
						objectValues := []basetypes.ObjectValue{}
						for _, v := range v.GetError().GetInnerError().GetDetails() {
							details := new(sitesPublicErrorDetailModel)

							if v.GetCode() != nil {
								details.Code = types.StringValue(*v.GetCode())
							} else {
								details.Code = types.StringNull()
							}
							if v.GetMessage() != nil {
								details.Message = types.StringValue(*v.GetMessage())
							} else {
								details.Message = types.StringNull()
							}
							if v.GetTarget() != nil {
								details.Target = types.StringValue(*v.GetTarget())
							} else {
								details.Target = types.StringNull()
							}
							objectValue, _ := types.ObjectValueFrom(ctx, details.AttributeTypes(), details)
							objectValues = append(objectValues, objectValue)
						}
						innerError.Details, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
					}
					if v.GetError().GetInnerError().GetMessage() != nil {
						innerError.Message = types.StringValue(*v.GetError().GetInnerError().GetMessage())
					} else {
						innerError.Message = types.StringNull()
					}
					if v.GetError().GetInnerError().GetTarget() != nil {
						innerError.Target = types.StringValue(*v.GetError().GetInnerError().GetTarget())
					} else {
						innerError.Target = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, innerError.AttributeTypes(), innerError)
					error.InnerError = objectValue
				}
				if v.GetError().GetMessage() != nil {
					error.Message = types.StringValue(*v.GetError().GetMessage())
				} else {
					error.Message = types.StringNull()
				}
				if v.GetError().GetTarget() != nil {
					error.Target = types.StringValue(*v.GetError().GetTarget())
				} else {
					error.Target = types.StringNull()
				}

				objectValue, _ := types.ObjectValueFrom(ctx, error.AttributeTypes(), error)
				value.Error = objectValue
			}
			if v.GetIsPersonalSite() != nil {
				value.IsPersonalSite = types.BoolValue(*v.GetIsPersonalSite())
			} else {
				value.IsPersonalSite = types.BoolNull()
			}
			if v.GetRoot() != nil {
				root := new(sitesRootModel)

				objectValue, _ := types.ObjectValueFrom(ctx, root.AttributeTypes(), root)
				value.Root = objectValue
			}
			if v.GetSharepointIds() != nil {
				sharepointIds := new(sitesSharepointIdsModel)

				if v.GetSharepointIds().GetListId() != nil {
					sharepointIds.ListId = types.StringValue(*v.GetSharepointIds().GetListId())
				} else {
					sharepointIds.ListId = types.StringNull()
				}
				if v.GetSharepointIds().GetListItemId() != nil {
					sharepointIds.ListItemId = types.StringValue(*v.GetSharepointIds().GetListItemId())
				} else {
					sharepointIds.ListItemId = types.StringNull()
				}
				if v.GetSharepointIds().GetListItemUniqueId() != nil {
					sharepointIds.ListItemUniqueId = types.StringValue(*v.GetSharepointIds().GetListItemUniqueId())
				} else {
					sharepointIds.ListItemUniqueId = types.StringNull()
				}
				if v.GetSharepointIds().GetSiteId() != nil {
					sharepointIds.SiteId = types.StringValue(*v.GetSharepointIds().GetSiteId())
				} else {
					sharepointIds.SiteId = types.StringNull()
				}
				if v.GetSharepointIds().GetSiteUrl() != nil {
					sharepointIds.SiteUrl = types.StringValue(*v.GetSharepointIds().GetSiteUrl())
				} else {
					sharepointIds.SiteUrl = types.StringNull()
				}
				if v.GetSharepointIds().GetTenantId() != nil {
					sharepointIds.TenantId = types.StringValue(*v.GetSharepointIds().GetTenantId())
				} else {
					sharepointIds.TenantId = types.StringNull()
				}
				if v.GetSharepointIds().GetWebId() != nil {
					sharepointIds.WebId = types.StringValue(*v.GetSharepointIds().GetWebId())
				} else {
					sharepointIds.WebId = types.StringNull()
				}

				objectValue, _ := types.ObjectValueFrom(ctx, sharepointIds.AttributeTypes(), sharepointIds)
				value.SharepointIds = objectValue
			}
			if v.GetSiteCollection() != nil {
				siteCollection := new(sitesSiteCollectionModel)

				if v.GetSiteCollection().GetArchivalDetails() != nil {
					archivalDetails := new(sitesSiteArchivalDetailsModel)

					if v.GetSiteCollection().GetArchivalDetails().GetArchiveStatus() != nil {
						archivalDetails.ArchiveStatus = types.StringValue(v.GetSiteCollection().GetArchivalDetails().GetArchiveStatus().String())
					} else {
						archivalDetails.ArchiveStatus = types.StringNull()
					}

					objectValue, _ := types.ObjectValueFrom(ctx, archivalDetails.AttributeTypes(), archivalDetails)
					siteCollection.ArchivalDetails = objectValue
				}
				if v.GetSiteCollection().GetDataLocationCode() != nil {
					siteCollection.DataLocationCode = types.StringValue(*v.GetSiteCollection().GetDataLocationCode())
				} else {
					siteCollection.DataLocationCode = types.StringNull()
				}
				if v.GetSiteCollection().GetHostname() != nil {
					siteCollection.Hostname = types.StringValue(*v.GetSiteCollection().GetHostname())
				} else {
					siteCollection.Hostname = types.StringNull()
				}
				if v.GetSiteCollection().GetRoot() != nil {
					root := new(sitesRootModel)

					objectValue, _ := types.ObjectValueFrom(ctx, root.AttributeTypes(), root)
					siteCollection.Root = objectValue
				}

				objectValue, _ := types.ObjectValueFrom(ctx, siteCollection.AttributeTypes(), siteCollection)
				value.SiteCollection = objectValue
			}
			objectValue, _ := types.ObjectValueFrom(ctx, value.AttributeTypes(), value)
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
