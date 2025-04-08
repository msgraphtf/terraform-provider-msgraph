package sites

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/sites"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &siteDataSource{}
	_ datasource.DataSourceWithConfigure = &siteDataSource{}
)

// NewSiteDataSource is a helper function to simplify the provider implementation.
func NewSiteDataSource() datasource.DataSource {
	return &siteDataSource{}
}

// siteDataSource is the data source implementation.
type siteDataSource struct {
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the data source type name.
func (d *siteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

// Configure adds the provider configured client to the data source.
func (d *siteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the data source.
func (d *siteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for an entity. Read-only.",
				Optional:    true,
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
								Optional:    true,
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
								Optional:    true,
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
								Optional:    true,
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
								Optional:    true,
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
								Optional:    true,
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
								Optional:    true,
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
						Optional:    true,
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
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *siteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfStateSite siteModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &tfStateSite)...)
	if resp.Diagnostics.HasError() {
		return
	}

	qparams := sites.SiteItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &sites.SiteItemRequestBuilderGetQueryParameters{
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
			},
		},
	}

	var result models.Siteable
	var err error

	if !tfStateSite.Id.IsNull() {
		result, err = d.client.Sites().BySiteId(tfStateSite.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"TODO: Specify required parameters",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting site",
			err.Error(),
		)
		return
	}

	if result.GetId() != nil {
		tfStateSite.Id = types.StringValue(*result.GetId())
	} else {
		tfStateSite.Id = types.StringNull()
	}
	if result.GetCreatedBy() != nil {
		createdBy := new(siteIdentitySetModel)

		if result.GetCreatedBy().GetApplication() != nil {
			application := new(siteIdentityModel)

			if result.GetCreatedBy().GetApplication().GetDisplayName() != nil {
				application.DisplayName = types.StringValue(*result.GetCreatedBy().GetApplication().GetDisplayName())
			} else {
				application.DisplayName = types.StringNull()
			}
			if result.GetCreatedBy().GetApplication().GetId() != nil {
				application.Id = types.StringValue(*result.GetCreatedBy().GetApplication().GetId())
			} else {
				application.Id = types.StringNull()
			}

			createdBy.Application, _ = types.ObjectValueFrom(ctx, application.AttributeTypes(), application)
		}
		if result.GetCreatedBy().GetDevice() != nil {
			device := new(siteIdentityModel)

			if result.GetCreatedBy().GetDevice().GetDisplayName() != nil {
				device.DisplayName = types.StringValue(*result.GetCreatedBy().GetDevice().GetDisplayName())
			} else {
				device.DisplayName = types.StringNull()
			}
			if result.GetCreatedBy().GetDevice().GetId() != nil {
				device.Id = types.StringValue(*result.GetCreatedBy().GetDevice().GetId())
			} else {
				device.Id = types.StringNull()
			}

			createdBy.Device, _ = types.ObjectValueFrom(ctx, device.AttributeTypes(), device)
		}
		if result.GetCreatedBy().GetUser() != nil {
			user := new(siteIdentityModel)

			if result.GetCreatedBy().GetUser().GetDisplayName() != nil {
				user.DisplayName = types.StringValue(*result.GetCreatedBy().GetUser().GetDisplayName())
			} else {
				user.DisplayName = types.StringNull()
			}
			if result.GetCreatedBy().GetUser().GetId() != nil {
				user.Id = types.StringValue(*result.GetCreatedBy().GetUser().GetId())
			} else {
				user.Id = types.StringNull()
			}

			createdBy.User, _ = types.ObjectValueFrom(ctx, user.AttributeTypes(), user)
		}

		tfStateSite.CreatedBy, _ = types.ObjectValueFrom(ctx, createdBy.AttributeTypes(), createdBy)
	}
	if result.GetCreatedDateTime() != nil {
		tfStateSite.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	} else {
		tfStateSite.CreatedDateTime = types.StringNull()
	}
	if result.GetDescription() != nil {
		tfStateSite.Description = types.StringValue(*result.GetDescription())
	} else {
		tfStateSite.Description = types.StringNull()
	}
	if result.GetETag() != nil {
		tfStateSite.ETag = types.StringValue(*result.GetETag())
	} else {
		tfStateSite.ETag = types.StringNull()
	}
	if result.GetLastModifiedBy() != nil {
		lastModifiedBy := new(siteIdentitySetModel)

		if result.GetLastModifiedBy().GetApplication() != nil {
			application := new(siteIdentityModel)

			if result.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
				application.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetApplication().GetDisplayName())
			} else {
				application.DisplayName = types.StringNull()
			}
			if result.GetLastModifiedBy().GetApplication().GetId() != nil {
				application.Id = types.StringValue(*result.GetLastModifiedBy().GetApplication().GetId())
			} else {
				application.Id = types.StringNull()
			}

			lastModifiedBy.Application, _ = types.ObjectValueFrom(ctx, application.AttributeTypes(), application)
		}
		if result.GetLastModifiedBy().GetDevice() != nil {
			device := new(siteIdentityModel)

			if result.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
				device.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetDevice().GetDisplayName())
			} else {
				device.DisplayName = types.StringNull()
			}
			if result.GetLastModifiedBy().GetDevice().GetId() != nil {
				device.Id = types.StringValue(*result.GetLastModifiedBy().GetDevice().GetId())
			} else {
				device.Id = types.StringNull()
			}

			lastModifiedBy.Device, _ = types.ObjectValueFrom(ctx, device.AttributeTypes(), device)
		}
		if result.GetLastModifiedBy().GetUser() != nil {
			user := new(siteIdentityModel)

			if result.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
				user.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetUser().GetDisplayName())
			} else {
				user.DisplayName = types.StringNull()
			}
			if result.GetLastModifiedBy().GetUser().GetId() != nil {
				user.Id = types.StringValue(*result.GetLastModifiedBy().GetUser().GetId())
			} else {
				user.Id = types.StringNull()
			}

			lastModifiedBy.User, _ = types.ObjectValueFrom(ctx, user.AttributeTypes(), user)
		}

		tfStateSite.LastModifiedBy, _ = types.ObjectValueFrom(ctx, lastModifiedBy.AttributeTypes(), lastModifiedBy)
	}
	if result.GetLastModifiedDateTime() != nil {
		tfStateSite.LastModifiedDateTime = types.StringValue(result.GetLastModifiedDateTime().String())
	} else {
		tfStateSite.LastModifiedDateTime = types.StringNull()
	}
	if result.GetName() != nil {
		tfStateSite.Name = types.StringValue(*result.GetName())
	} else {
		tfStateSite.Name = types.StringNull()
	}
	if result.GetParentReference() != nil {
		parentReference := new(siteItemReferenceModel)

		if result.GetParentReference().GetDriveId() != nil {
			parentReference.DriveId = types.StringValue(*result.GetParentReference().GetDriveId())
		} else {
			parentReference.DriveId = types.StringNull()
		}
		if result.GetParentReference().GetDriveType() != nil {
			parentReference.DriveType = types.StringValue(*result.GetParentReference().GetDriveType())
		} else {
			parentReference.DriveType = types.StringNull()
		}
		if result.GetParentReference().GetId() != nil {
			parentReference.Id = types.StringValue(*result.GetParentReference().GetId())
		} else {
			parentReference.Id = types.StringNull()
		}
		if result.GetParentReference().GetName() != nil {
			parentReference.Name = types.StringValue(*result.GetParentReference().GetName())
		} else {
			parentReference.Name = types.StringNull()
		}
		if result.GetParentReference().GetPath() != nil {
			parentReference.Path = types.StringValue(*result.GetParentReference().GetPath())
		} else {
			parentReference.Path = types.StringNull()
		}
		if result.GetParentReference().GetShareId() != nil {
			parentReference.ShareId = types.StringValue(*result.GetParentReference().GetShareId())
		} else {
			parentReference.ShareId = types.StringNull()
		}
		if result.GetParentReference().GetSharepointIds() != nil {
			sharepointIds := new(siteSharepointIdsModel)

			if result.GetParentReference().GetSharepointIds().GetListId() != nil {
				sharepointIds.ListId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListId())
			} else {
				sharepointIds.ListId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetListItemId() != nil {
				sharepointIds.ListItemId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListItemId())
			} else {
				sharepointIds.ListItemId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
				sharepointIds.ListItemUniqueId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListItemUniqueId())
			} else {
				sharepointIds.ListItemUniqueId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetSiteId() != nil {
				sharepointIds.SiteId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetSiteId())
			} else {
				sharepointIds.SiteId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
				sharepointIds.SiteUrl = types.StringValue(*result.GetParentReference().GetSharepointIds().GetSiteUrl())
			} else {
				sharepointIds.SiteUrl = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetTenantId() != nil {
				sharepointIds.TenantId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetTenantId())
			} else {
				sharepointIds.TenantId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetWebId() != nil {
				sharepointIds.WebId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetWebId())
			} else {
				sharepointIds.WebId = types.StringNull()
			}

			parentReference.SharepointIds, _ = types.ObjectValueFrom(ctx, sharepointIds.AttributeTypes(), sharepointIds)
		}
		if result.GetParentReference().GetSiteId() != nil {
			parentReference.SiteId = types.StringValue(*result.GetParentReference().GetSiteId())
		} else {
			parentReference.SiteId = types.StringNull()
		}

		tfStateSite.ParentReference, _ = types.ObjectValueFrom(ctx, parentReference.AttributeTypes(), parentReference)
	}
	if result.GetWebUrl() != nil {
		tfStateSite.WebUrl = types.StringValue(*result.GetWebUrl())
	} else {
		tfStateSite.WebUrl = types.StringNull()
	}
	if result.GetDisplayName() != nil {
		tfStateSite.DisplayName = types.StringValue(*result.GetDisplayName())
	} else {
		tfStateSite.DisplayName = types.StringNull()
	}
	if result.GetError() != nil {
		error := new(sitePublicErrorModel)

		if result.GetError().GetCode() != nil {
			error.Code = types.StringValue(*result.GetError().GetCode())
		} else {
			error.Code = types.StringNull()
		}
		if len(result.GetError().GetDetails()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, v := range result.GetError().GetDetails() {
				details := new(sitePublicErrorDetailModel)

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
		if result.GetError().GetInnerError() != nil {
			innerError := new(sitePublicInnerErrorModel)

			if result.GetError().GetInnerError().GetCode() != nil {
				innerError.Code = types.StringValue(*result.GetError().GetInnerError().GetCode())
			} else {
				innerError.Code = types.StringNull()
			}
			if len(result.GetError().GetInnerError().GetDetails()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, v := range result.GetError().GetInnerError().GetDetails() {
					details := new(sitePublicErrorDetailModel)

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
			if result.GetError().GetInnerError().GetMessage() != nil {
				innerError.Message = types.StringValue(*result.GetError().GetInnerError().GetMessage())
			} else {
				innerError.Message = types.StringNull()
			}
			if result.GetError().GetInnerError().GetTarget() != nil {
				innerError.Target = types.StringValue(*result.GetError().GetInnerError().GetTarget())
			} else {
				innerError.Target = types.StringNull()
			}

			error.InnerError, _ = types.ObjectValueFrom(ctx, innerError.AttributeTypes(), innerError)
		}
		if result.GetError().GetMessage() != nil {
			error.Message = types.StringValue(*result.GetError().GetMessage())
		} else {
			error.Message = types.StringNull()
		}
		if result.GetError().GetTarget() != nil {
			error.Target = types.StringValue(*result.GetError().GetTarget())
		} else {
			error.Target = types.StringNull()
		}

		tfStateSite.Error, _ = types.ObjectValueFrom(ctx, error.AttributeTypes(), error)
	}
	if result.GetIsPersonalSite() != nil {
		tfStateSite.IsPersonalSite = types.BoolValue(*result.GetIsPersonalSite())
	} else {
		tfStateSite.IsPersonalSite = types.BoolNull()
	}
	if result.GetRoot() != nil {
		root := new(siteRootModel)

		tfStateSite.Root, _ = types.ObjectValueFrom(ctx, root.AttributeTypes(), root)
	}
	if result.GetSharepointIds() != nil {
		sharepointIds := new(siteSharepointIdsModel)

		if result.GetSharepointIds().GetListId() != nil {
			sharepointIds.ListId = types.StringValue(*result.GetSharepointIds().GetListId())
		} else {
			sharepointIds.ListId = types.StringNull()
		}
		if result.GetSharepointIds().GetListItemId() != nil {
			sharepointIds.ListItemId = types.StringValue(*result.GetSharepointIds().GetListItemId())
		} else {
			sharepointIds.ListItemId = types.StringNull()
		}
		if result.GetSharepointIds().GetListItemUniqueId() != nil {
			sharepointIds.ListItemUniqueId = types.StringValue(*result.GetSharepointIds().GetListItemUniqueId())
		} else {
			sharepointIds.ListItemUniqueId = types.StringNull()
		}
		if result.GetSharepointIds().GetSiteId() != nil {
			sharepointIds.SiteId = types.StringValue(*result.GetSharepointIds().GetSiteId())
		} else {
			sharepointIds.SiteId = types.StringNull()
		}
		if result.GetSharepointIds().GetSiteUrl() != nil {
			sharepointIds.SiteUrl = types.StringValue(*result.GetSharepointIds().GetSiteUrl())
		} else {
			sharepointIds.SiteUrl = types.StringNull()
		}
		if result.GetSharepointIds().GetTenantId() != nil {
			sharepointIds.TenantId = types.StringValue(*result.GetSharepointIds().GetTenantId())
		} else {
			sharepointIds.TenantId = types.StringNull()
		}
		if result.GetSharepointIds().GetWebId() != nil {
			sharepointIds.WebId = types.StringValue(*result.GetSharepointIds().GetWebId())
		} else {
			sharepointIds.WebId = types.StringNull()
		}

		tfStateSite.SharepointIds, _ = types.ObjectValueFrom(ctx, sharepointIds.AttributeTypes(), sharepointIds)
	}
	if result.GetSiteCollection() != nil {
		siteCollection := new(siteSiteCollectionModel)

		if result.GetSiteCollection().GetArchivalDetails() != nil {
			archivalDetails := new(siteSiteArchivalDetailsModel)

			if result.GetSiteCollection().GetArchivalDetails().GetArchiveStatus() != nil {
				archivalDetails.ArchiveStatus = types.StringValue(result.GetSiteCollection().GetArchivalDetails().GetArchiveStatus().String())
			} else {
				archivalDetails.ArchiveStatus = types.StringNull()
			}

			siteCollection.ArchivalDetails, _ = types.ObjectValueFrom(ctx, archivalDetails.AttributeTypes(), archivalDetails)
		}
		if result.GetSiteCollection().GetDataLocationCode() != nil {
			siteCollection.DataLocationCode = types.StringValue(*result.GetSiteCollection().GetDataLocationCode())
		} else {
			siteCollection.DataLocationCode = types.StringNull()
		}
		if result.GetSiteCollection().GetHostname() != nil {
			siteCollection.Hostname = types.StringValue(*result.GetSiteCollection().GetHostname())
		} else {
			siteCollection.Hostname = types.StringNull()
		}
		if result.GetSiteCollection().GetRoot() != nil {
			root := new(siteRootModel)

			siteCollection.Root, _ = types.ObjectValueFrom(ctx, root.AttributeTypes(), root)
		}

		tfStateSite.SiteCollection, _ = types.ObjectValueFrom(ctx, siteCollection.AttributeTypes(), siteCollection)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateSite)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
