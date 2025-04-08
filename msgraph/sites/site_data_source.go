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

	var response models.Siteable
	var err error

	if !tfStateSite.Id.IsNull() {
		response, err = d.client.Sites().BySiteId(tfStateSite.Id.ValueString()).Get(context.Background(), &qparams)
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

	if response.GetId() != nil {
		tfStateSite.Id = types.StringValue(*response.GetId())
	} else {
		tfStateSite.Id = types.StringNull()
	}
	if response.GetCreatedBy() != nil {
		tfStateIdentitySet := siteIdentitySetModel{}

		if response.GetCreatedBy().GetApplication() != nil {
			tfStateIdentity := siteIdentityModel{}

			if response.GetCreatedBy().GetApplication().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*response.GetCreatedBy().GetApplication().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if response.GetCreatedBy().GetApplication().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*response.GetCreatedBy().GetApplication().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if response.GetCreatedBy().GetDevice() != nil {
			tfStateIdentity := siteIdentityModel{}

			if response.GetCreatedBy().GetDevice().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*response.GetCreatedBy().GetDevice().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if response.GetCreatedBy().GetDevice().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*response.GetCreatedBy().GetDevice().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if response.GetCreatedBy().GetUser() != nil {
			tfStateIdentity := siteIdentityModel{}

			if response.GetCreatedBy().GetUser().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*response.GetCreatedBy().GetUser().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if response.GetCreatedBy().GetUser().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*response.GetCreatedBy().GetUser().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}

		tfStateSite.CreatedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
	}
	if response.GetCreatedDateTime() != nil {
		tfStateSite.CreatedDateTime = types.StringValue(response.GetCreatedDateTime().String())
	} else {
		tfStateSite.CreatedDateTime = types.StringNull()
	}
	if response.GetDescription() != nil {
		tfStateSite.Description = types.StringValue(*response.GetDescription())
	} else {
		tfStateSite.Description = types.StringNull()
	}
	if response.GetETag() != nil {
		tfStateSite.ETag = types.StringValue(*response.GetETag())
	} else {
		tfStateSite.ETag = types.StringNull()
	}
	if response.GetLastModifiedBy() != nil {
		tfStateIdentitySet := siteIdentitySetModel{}

		if response.GetLastModifiedBy().GetApplication() != nil {
			tfStateIdentity := siteIdentityModel{}

			if response.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*response.GetLastModifiedBy().GetApplication().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if response.GetLastModifiedBy().GetApplication().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*response.GetLastModifiedBy().GetApplication().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Application, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if response.GetLastModifiedBy().GetDevice() != nil {
			tfStateIdentity := siteIdentityModel{}

			if response.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*response.GetLastModifiedBy().GetDevice().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if response.GetLastModifiedBy().GetDevice().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*response.GetLastModifiedBy().GetDevice().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.Device, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}
		if response.GetLastModifiedBy().GetUser() != nil {
			tfStateIdentity := siteIdentityModel{}

			if response.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
				tfStateIdentity.DisplayName = types.StringValue(*response.GetLastModifiedBy().GetUser().GetDisplayName())
			} else {
				tfStateIdentity.DisplayName = types.StringNull()
			}
			if response.GetLastModifiedBy().GetUser().GetId() != nil {
				tfStateIdentity.Id = types.StringValue(*response.GetLastModifiedBy().GetUser().GetId())
			} else {
				tfStateIdentity.Id = types.StringNull()
			}

			tfStateIdentitySet.User, _ = types.ObjectValueFrom(ctx, tfStateIdentity.AttributeTypes(), tfStateIdentity)
		}

		tfStateSite.LastModifiedBy, _ = types.ObjectValueFrom(ctx, tfStateIdentitySet.AttributeTypes(), tfStateIdentitySet)
	}
	if response.GetLastModifiedDateTime() != nil {
		tfStateSite.LastModifiedDateTime = types.StringValue(response.GetLastModifiedDateTime().String())
	} else {
		tfStateSite.LastModifiedDateTime = types.StringNull()
	}
	if response.GetName() != nil {
		tfStateSite.Name = types.StringValue(*response.GetName())
	} else {
		tfStateSite.Name = types.StringNull()
	}
	if response.GetParentReference() != nil {
		tfStateItemReference := siteItemReferenceModel{}

		if response.GetParentReference().GetDriveId() != nil {
			tfStateItemReference.DriveId = types.StringValue(*response.GetParentReference().GetDriveId())
		} else {
			tfStateItemReference.DriveId = types.StringNull()
		}
		if response.GetParentReference().GetDriveType() != nil {
			tfStateItemReference.DriveType = types.StringValue(*response.GetParentReference().GetDriveType())
		} else {
			tfStateItemReference.DriveType = types.StringNull()
		}
		if response.GetParentReference().GetId() != nil {
			tfStateItemReference.Id = types.StringValue(*response.GetParentReference().GetId())
		} else {
			tfStateItemReference.Id = types.StringNull()
		}
		if response.GetParentReference().GetName() != nil {
			tfStateItemReference.Name = types.StringValue(*response.GetParentReference().GetName())
		} else {
			tfStateItemReference.Name = types.StringNull()
		}
		if response.GetParentReference().GetPath() != nil {
			tfStateItemReference.Path = types.StringValue(*response.GetParentReference().GetPath())
		} else {
			tfStateItemReference.Path = types.StringNull()
		}
		if response.GetParentReference().GetShareId() != nil {
			tfStateItemReference.ShareId = types.StringValue(*response.GetParentReference().GetShareId())
		} else {
			tfStateItemReference.ShareId = types.StringNull()
		}
		if response.GetParentReference().GetSharepointIds() != nil {
			tfStateSharepointIds := siteSharepointIdsModel{}

			if response.GetParentReference().GetSharepointIds().GetListId() != nil {
				tfStateSharepointIds.ListId = types.StringValue(*response.GetParentReference().GetSharepointIds().GetListId())
			} else {
				tfStateSharepointIds.ListId = types.StringNull()
			}
			if response.GetParentReference().GetSharepointIds().GetListItemId() != nil {
				tfStateSharepointIds.ListItemId = types.StringValue(*response.GetParentReference().GetSharepointIds().GetListItemId())
			} else {
				tfStateSharepointIds.ListItemId = types.StringNull()
			}
			if response.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
				tfStateSharepointIds.ListItemUniqueId = types.StringValue(*response.GetParentReference().GetSharepointIds().GetListItemUniqueId())
			} else {
				tfStateSharepointIds.ListItemUniqueId = types.StringNull()
			}
			if response.GetParentReference().GetSharepointIds().GetSiteId() != nil {
				tfStateSharepointIds.SiteId = types.StringValue(*response.GetParentReference().GetSharepointIds().GetSiteId())
			} else {
				tfStateSharepointIds.SiteId = types.StringNull()
			}
			if response.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
				tfStateSharepointIds.SiteUrl = types.StringValue(*response.GetParentReference().GetSharepointIds().GetSiteUrl())
			} else {
				tfStateSharepointIds.SiteUrl = types.StringNull()
			}
			if response.GetParentReference().GetSharepointIds().GetTenantId() != nil {
				tfStateSharepointIds.TenantId = types.StringValue(*response.GetParentReference().GetSharepointIds().GetTenantId())
			} else {
				tfStateSharepointIds.TenantId = types.StringNull()
			}
			if response.GetParentReference().GetSharepointIds().GetWebId() != nil {
				tfStateSharepointIds.WebId = types.StringValue(*response.GetParentReference().GetSharepointIds().GetWebId())
			} else {
				tfStateSharepointIds.WebId = types.StringNull()
			}

			tfStateItemReference.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
		}
		if response.GetParentReference().GetSiteId() != nil {
			tfStateItemReference.SiteId = types.StringValue(*response.GetParentReference().GetSiteId())
		} else {
			tfStateItemReference.SiteId = types.StringNull()
		}

		tfStateSite.ParentReference, _ = types.ObjectValueFrom(ctx, tfStateItemReference.AttributeTypes(), tfStateItemReference)
	}
	if response.GetWebUrl() != nil {
		tfStateSite.WebUrl = types.StringValue(*response.GetWebUrl())
	} else {
		tfStateSite.WebUrl = types.StringNull()
	}
	if response.GetDisplayName() != nil {
		tfStateSite.DisplayName = types.StringValue(*response.GetDisplayName())
	} else {
		tfStateSite.DisplayName = types.StringNull()
	}
	if response.GetError() != nil {
		tfStatePublicError := sitePublicErrorModel{}

		if response.GetError().GetCode() != nil {
			tfStatePublicError.Code = types.StringValue(*response.GetError().GetCode())
		} else {
			tfStatePublicError.Code = types.StringNull()
		}
		if len(response.GetError().GetDetails()) > 0 {
			objectValues := []basetypes.ObjectValue{}
			for _, responseDetails := range response.GetError().GetDetails() {
				tfStatePublicErrorDetail := sitePublicErrorDetailModel{}

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
		if response.GetError().GetInnerError() != nil {
			tfStatePublicInnerError := sitePublicInnerErrorModel{}

			if response.GetError().GetInnerError().GetCode() != nil {
				tfStatePublicInnerError.Code = types.StringValue(*response.GetError().GetInnerError().GetCode())
			} else {
				tfStatePublicInnerError.Code = types.StringNull()
			}
			if len(response.GetError().GetInnerError().GetDetails()) > 0 {
				objectValues := []basetypes.ObjectValue{}
				for _, responseDetails := range response.GetError().GetInnerError().GetDetails() {
					tfStatePublicErrorDetail := sitePublicErrorDetailModel{}

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
			if response.GetError().GetInnerError().GetMessage() != nil {
				tfStatePublicInnerError.Message = types.StringValue(*response.GetError().GetInnerError().GetMessage())
			} else {
				tfStatePublicInnerError.Message = types.StringNull()
			}
			if response.GetError().GetInnerError().GetTarget() != nil {
				tfStatePublicInnerError.Target = types.StringValue(*response.GetError().GetInnerError().GetTarget())
			} else {
				tfStatePublicInnerError.Target = types.StringNull()
			}

			tfStatePublicError.InnerError, _ = types.ObjectValueFrom(ctx, tfStatePublicInnerError.AttributeTypes(), tfStatePublicInnerError)
		}
		if response.GetError().GetMessage() != nil {
			tfStatePublicError.Message = types.StringValue(*response.GetError().GetMessage())
		} else {
			tfStatePublicError.Message = types.StringNull()
		}
		if response.GetError().GetTarget() != nil {
			tfStatePublicError.Target = types.StringValue(*response.GetError().GetTarget())
		} else {
			tfStatePublicError.Target = types.StringNull()
		}

		tfStateSite.Error, _ = types.ObjectValueFrom(ctx, tfStatePublicError.AttributeTypes(), tfStatePublicError)
	}
	if response.GetIsPersonalSite() != nil {
		tfStateSite.IsPersonalSite = types.BoolValue(*response.GetIsPersonalSite())
	} else {
		tfStateSite.IsPersonalSite = types.BoolNull()
	}
	if response.GetRoot() != nil {
		tfStateRoot := siteRootModel{}

		tfStateSite.Root, _ = types.ObjectValueFrom(ctx, tfStateRoot.AttributeTypes(), tfStateRoot)
	}
	if response.GetSharepointIds() != nil {
		tfStateSharepointIds := siteSharepointIdsModel{}

		if response.GetSharepointIds().GetListId() != nil {
			tfStateSharepointIds.ListId = types.StringValue(*response.GetSharepointIds().GetListId())
		} else {
			tfStateSharepointIds.ListId = types.StringNull()
		}
		if response.GetSharepointIds().GetListItemId() != nil {
			tfStateSharepointIds.ListItemId = types.StringValue(*response.GetSharepointIds().GetListItemId())
		} else {
			tfStateSharepointIds.ListItemId = types.StringNull()
		}
		if response.GetSharepointIds().GetListItemUniqueId() != nil {
			tfStateSharepointIds.ListItemUniqueId = types.StringValue(*response.GetSharepointIds().GetListItemUniqueId())
		} else {
			tfStateSharepointIds.ListItemUniqueId = types.StringNull()
		}
		if response.GetSharepointIds().GetSiteId() != nil {
			tfStateSharepointIds.SiteId = types.StringValue(*response.GetSharepointIds().GetSiteId())
		} else {
			tfStateSharepointIds.SiteId = types.StringNull()
		}
		if response.GetSharepointIds().GetSiteUrl() != nil {
			tfStateSharepointIds.SiteUrl = types.StringValue(*response.GetSharepointIds().GetSiteUrl())
		} else {
			tfStateSharepointIds.SiteUrl = types.StringNull()
		}
		if response.GetSharepointIds().GetTenantId() != nil {
			tfStateSharepointIds.TenantId = types.StringValue(*response.GetSharepointIds().GetTenantId())
		} else {
			tfStateSharepointIds.TenantId = types.StringNull()
		}
		if response.GetSharepointIds().GetWebId() != nil {
			tfStateSharepointIds.WebId = types.StringValue(*response.GetSharepointIds().GetWebId())
		} else {
			tfStateSharepointIds.WebId = types.StringNull()
		}

		tfStateSite.SharepointIds, _ = types.ObjectValueFrom(ctx, tfStateSharepointIds.AttributeTypes(), tfStateSharepointIds)
	}
	if response.GetSiteCollection() != nil {
		tfStateSiteCollection := siteSiteCollectionModel{}

		if response.GetSiteCollection().GetArchivalDetails() != nil {
			tfStateSiteArchivalDetails := siteSiteArchivalDetailsModel{}

			if response.GetSiteCollection().GetArchivalDetails().GetArchiveStatus() != nil {
				tfStateSiteArchivalDetails.ArchiveStatus = types.StringValue(response.GetSiteCollection().GetArchivalDetails().GetArchiveStatus().String())
			} else {
				tfStateSiteArchivalDetails.ArchiveStatus = types.StringNull()
			}

			tfStateSiteCollection.ArchivalDetails, _ = types.ObjectValueFrom(ctx, tfStateSiteArchivalDetails.AttributeTypes(), tfStateSiteArchivalDetails)
		}
		if response.GetSiteCollection().GetDataLocationCode() != nil {
			tfStateSiteCollection.DataLocationCode = types.StringValue(*response.GetSiteCollection().GetDataLocationCode())
		} else {
			tfStateSiteCollection.DataLocationCode = types.StringNull()
		}
		if response.GetSiteCollection().GetHostname() != nil {
			tfStateSiteCollection.Hostname = types.StringValue(*response.GetSiteCollection().GetHostname())
		} else {
			tfStateSiteCollection.Hostname = types.StringNull()
		}
		if response.GetSiteCollection().GetRoot() != nil {
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
