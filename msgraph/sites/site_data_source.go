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
		tfStateIdentitySet := siteIdentitySetModel{}

		if result.GetCreatedBy().GetApplication() != nil {
			tfStateIdentity := siteIdentityModel{}

			if result.GetCreatedBy().GetApplication().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*result.GetCreatedBy().GetApplication().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if result.GetCreatedBy().GetApplication().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*result.GetCreatedBy().GetApplication().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if result.GetCreatedBy().GetDevice() != nil {
			tfStateIdentity := siteIdentityModel{}

			if result.GetCreatedBy().GetDevice().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*result.GetCreatedBy().GetDevice().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if result.GetCreatedBy().GetDevice().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*result.GetCreatedBy().GetDevice().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if result.GetCreatedBy().GetUser() != nil {
			tfStateIdentity := siteIdentityModel{}

			if result.GetCreatedBy().GetUser().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*result.GetCreatedBy().GetUser().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if result.GetCreatedBy().GetUser().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*result.GetCreatedBy().GetUser().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}

		tfStateSite.CreatedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
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
		tfStateIdentitySet := siteIdentitySetModel{}

		if result.GetLastModifiedBy().GetApplication() != nil {
			tfStateIdentity := siteIdentityModel{}

			if result.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetApplication().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if result.GetLastModifiedBy().GetApplication().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*result.GetLastModifiedBy().GetApplication().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if result.GetLastModifiedBy().GetDevice() != nil {
			tfStateIdentity := siteIdentityModel{}

			if result.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetDevice().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if result.GetLastModifiedBy().GetDevice().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*result.GetLastModifiedBy().GetDevice().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if result.GetLastModifiedBy().GetUser() != nil {
			tfStateIdentity := siteIdentityModel{}

			if result.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetUser().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if result.GetLastModifiedBy().GetUser().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*result.GetLastModifiedBy().GetUser().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}

		tfStateSite.LastModifiedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
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
		tfStateItemReference := siteItemReferenceModel{}

		if result.GetParentReference().GetDriveId() != nil {
			tfStateItemReference.DriveId = types.StringValue(*result.GetParentReference().GetDriveId())
		} else {
			tfStateItemReference.DriveId = types.StringNull()
		}
		if result.GetParentReference().GetDriveType() != nil {
			tfStateItemReference.DriveType = types.StringValue(*result.GetParentReference().GetDriveType())
		} else {
			tfStateItemReference.DriveType = types.StringNull()
		}
		if result.GetParentReference().GetId() != nil {
			tfStateItemReference.Id = types.StringValue(*result.GetParentReference().GetId())
		} else {
			tfStateItemReference.Id = types.StringNull()
		}
		if result.GetParentReference().GetName() != nil {
			tfStateItemReference.Name = types.StringValue(*result.GetParentReference().GetName())
		} else {
			tfStateItemReference.Name = types.StringNull()
		}
		if result.GetParentReference().GetPath() != nil {
			tfStateItemReference.Path = types.StringValue(*result.GetParentReference().GetPath())
		} else {
			tfStateItemReference.Path = types.StringNull()
		}
		if result.GetParentReference().GetShareId() != nil {
			tfStateItemReference.ShareId = types.StringValue(*result.GetParentReference().GetShareId())
		} else {
			tfStateItemReference.ShareId = types.StringNull()
		}
		if result.GetParentReference().GetSharepointIds() != nil {
			tfStateSharepointIds := siteSharepointIdsModel{}

			if result.GetParentReference().GetSharepointIds().GetListId() != nil {
				tfStateSharepointIds.ListId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListId())
			} else {
				tfStateSharepointIds.ListId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetListItemId() != nil {
				tfStateSharepointIds.ListItemId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListItemId())
			} else {
				tfStateSharepointIds.ListItemId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
				tfStateSharepointIds.ListItemUniqueId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListItemUniqueId())
			} else {
				tfStateSharepointIds.ListItemUniqueId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetSiteId() != nil {
				tfStateSharepointIds.SiteId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetSiteId())
			} else {
				tfStateSharepointIds.SiteId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
				tfStateSharepointIds.SiteUrl = types.StringValue(*result.GetParentReference().GetSharepointIds().GetSiteUrl())
			} else {
				tfStateSharepointIds.SiteUrl = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetTenantId() != nil {
				tfStateSharepointIds.TenantId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetTenantId())
			} else {
				tfStateSharepointIds.TenantId = types.StringNull()
			}
			if result.GetParentReference().GetSharepointIds().GetWebId() != nil {
				tfStateSharepointIds.WebId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetWebId())
			} else {
				tfStateSharepointIds.WebId = types.StringNull()
			}

			tfStateItemReference.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
		}
		if result.GetParentReference().GetSiteId() != nil {
			tfStateItemReference.SiteId = types.StringValue(*result.GetParentReference().GetSiteId())
		} else {
			tfStateItemReference.SiteId = types.StringNull()
		}

		tfStateSite.ParentReference, _ = types.ObjectValueFrom(ctx, tfStateItemReference.AttributeTypes(), tfStateItemReference)
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
		tfStatePublicError := sitePublicErrorModel{}

		if result.GetError().GetCode() != nil {
			tfStatePublicError.Code = types.StringValue(*result.GetError().GetCode())
		} else {
			tfStatePublicError.Code = types.StringNull()
		}
		if len(result.GetError().GetDetails()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, resultDetails := range result.GetError().GetDetails() {
				tfStatePublicErrorDetail := sitePublicErrorDetailModel{}

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
		if result.GetError().GetInnerError() != nil {
			tfStatePublicInnerError := sitePublicInnerErrorModel{}

			if result.GetError().GetInnerError().GetCode() != nil {
				tfStatePublicInnerError.Code = types.StringValue(*result.GetError().GetInnerError().GetCode())
			} else {
				tfStatePublicInnerError.Code = types.StringNull()
			}
			if len(result.GetError().GetInnerError().GetDetails()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, resultDetails := range result.GetError().GetInnerError().GetDetails() {
					tfStatePublicErrorDetail := sitePublicErrorDetailModel{}

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
			if result.GetError().GetInnerError().GetMessage() != nil {
				tfStatePublicInnerError.Message = types.StringValue(*result.GetError().GetInnerError().GetMessage())
			} else {
				tfStatePublicInnerError.Message = types.StringNull()
			}
			if result.GetError().GetInnerError().GetTarget() != nil {
				tfStatePublicInnerError.Target = types.StringValue(*result.GetError().GetInnerError().GetTarget())
			} else {
				tfStatePublicInnerError.Target = types.StringNull()
			}

			tfStatePublicError.InnerError, _ = types.ObjectValueFrom(ctx, tfStatePublicInnerError.AttributeTypes(), tfStatePublicInnerError)
		}
		if result.GetError().GetMessage() != nil {
			tfStatePublicError.Message = types.StringValue(*result.GetError().GetMessage())
		} else {
			tfStatePublicError.Message = types.StringNull()
		}
		if result.GetError().GetTarget() != nil {
			tfStatePublicError.Target = types.StringValue(*result.GetError().GetTarget())
		} else {
			tfStatePublicError.Target = types.StringNull()
		}

		tfStateSite.Error, _ = types.ObjectValueFrom(ctx, tfStatePublicError.AttributeTypes(), tfStatePublicError)
	}
	if result.GetIsPersonalSite() != nil {
		tfStateSite.IsPersonalSite = types.BoolValue(*result.GetIsPersonalSite())
	} else {
		tfStateSite.IsPersonalSite = types.BoolNull()
	}
	if result.GetRoot() != nil {
		tfStateRoot := siteRootModel{}

		tfStateSite.Root, _ = types.ObjectValueFrom(ctx, tfStateRoot.AttributeTypes(), tfStateRoot)
	}
	if result.GetSharepointIds() != nil {
		tfStateSharepointIds := siteSharepointIdsModel{}

		if result.GetSharepointIds().GetListId() != nil {
			tfStateSharepointIds.ListId = types.StringValue(*result.GetSharepointIds().GetListId())
		} else {
			tfStateSharepointIds.ListId = types.StringNull()
		}
		if result.GetSharepointIds().GetListItemId() != nil {
			tfStateSharepointIds.ListItemId = types.StringValue(*result.GetSharepointIds().GetListItemId())
		} else {
			tfStateSharepointIds.ListItemId = types.StringNull()
		}
		if result.GetSharepointIds().GetListItemUniqueId() != nil {
			tfStateSharepointIds.ListItemUniqueId = types.StringValue(*result.GetSharepointIds().GetListItemUniqueId())
		} else {
			tfStateSharepointIds.ListItemUniqueId = types.StringNull()
		}
		if result.GetSharepointIds().GetSiteId() != nil {
			tfStateSharepointIds.SiteId = types.StringValue(*result.GetSharepointIds().GetSiteId())
		} else {
			tfStateSharepointIds.SiteId = types.StringNull()
		}
		if result.GetSharepointIds().GetSiteUrl() != nil {
			tfStateSharepointIds.SiteUrl = types.StringValue(*result.GetSharepointIds().GetSiteUrl())
		} else {
			tfStateSharepointIds.SiteUrl = types.StringNull()
		}
		if result.GetSharepointIds().GetTenantId() != nil {
			tfStateSharepointIds.TenantId = types.StringValue(*result.GetSharepointIds().GetTenantId())
		} else {
			tfStateSharepointIds.TenantId = types.StringNull()
		}
		if result.GetSharepointIds().GetWebId() != nil {
			tfStateSharepointIds.WebId = types.StringValue(*result.GetSharepointIds().GetWebId())
		} else {
			tfStateSharepointIds.WebId = types.StringNull()
		}

		tfStateSite.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
	}
	if result.GetSiteCollection() != nil {
		tfStateSiteCollection := siteSiteCollectionModel{}

		if result.GetSiteCollection().GetArchivalDetails() != nil {
			tfStateSiteArchivalDetails := siteSiteArchivalDetailsModel{}

			if result.GetSiteCollection().GetArchivalDetails().GetArchiveStatus() != nil {
				tfStateSiteArchivalDetails.ArchiveStatus = types.StringValue(result.GetSiteCollection().GetArchivalDetails().GetArchiveStatus().String())
			} else {
				tfStateSiteArchivalDetails.ArchiveStatus = types.StringNull()
			}

			tfStateSiteCollection.ArchivalDetails, _ = types.ObjectValueFrom(ctx, tfStateSiteArchivalDetails.AttributeTypes(), tfStateSiteArchivalDetails)
		}
		if result.GetSiteCollection().GetDataLocationCode() != nil {
			tfStateSiteCollection.DataLocationCode = types.StringValue(*result.GetSiteCollection().GetDataLocationCode())
		} else {
			tfStateSiteCollection.DataLocationCode = types.StringNull()
		}
		if result.GetSiteCollection().GetHostname() != nil {
			tfStateSiteCollection.Hostname = types.StringValue(*result.GetSiteCollection().GetHostname())
		} else {
			tfStateSiteCollection.Hostname = types.StringNull()
		}
		if result.GetSiteCollection().GetRoot() != nil {
			tfStateRoot := siteRootModel{}

			tfStateSiteCollection.Root, _ = types.ObjectValueFrom(ctx, tfStateRoot.AttributeTypes(), tfStateRoot)
		}

		tfStateSite.SiteCollection, _ = types.ObjectValueFrom(ctx, tfStateSiteCollection.AttributeTypes(), tfStateSiteCollection)
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &tfStateSite)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
