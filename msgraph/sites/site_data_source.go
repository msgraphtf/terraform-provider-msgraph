package sites

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	}
}

type siteDataSourceModel struct {
	Id                   types.String                        `tfsdk:"id"`
	CreatedBy            *siteCreatedByDataSourceModel       `tfsdk:"created_by"`
	CreatedDateTime      types.String                        `tfsdk:"created_date_time"`
	Description          types.String                        `tfsdk:"description"`
	ETag                 types.String                        `tfsdk:"e_tag"`
	LastModifiedBy       *siteLastModifiedByDataSourceModel  `tfsdk:"last_modified_by"`
	LastModifiedDateTime types.String                        `tfsdk:"last_modified_date_time"`
	Name                 types.String                        `tfsdk:"name"`
	ParentReference      *siteParentReferenceDataSourceModel `tfsdk:"parent_reference"`
	WebUrl               types.String                        `tfsdk:"web_url"`
	DisplayName          types.String                        `tfsdk:"display_name"`
	Error                *siteErrorDataSourceModel           `tfsdk:"error"`
	IsPersonalSite       types.Bool                          `tfsdk:"is_personal_site"`
	Root                 *siteRootDataSourceModel            `tfsdk:"root"`
	SharepointIds        *siteSharepointIdsDataSourceModel   `tfsdk:"sharepoint_ids"`
	SiteCollection       *siteSiteCollectionDataSourceModel  `tfsdk:"site_collection"`
}

type siteCreatedByDataSourceModel struct {
	Application *siteApplicationDataSourceModel `tfsdk:"application"`
	Device      *siteDeviceDataSourceModel      `tfsdk:"device"`
	User        *siteUserDataSourceModel        `tfsdk:"user"`
}

type siteApplicationDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteDeviceDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteUserDataSourceModel struct {
	DisplayName types.String `tfsdk:"display_name"`
	Id          types.String `tfsdk:"id"`
}

type siteLastModifiedByDataSourceModel struct {
	Application *siteApplicationDataSourceModel `tfsdk:"application"`
	Device      *siteDeviceDataSourceModel      `tfsdk:"device"`
	User        *siteUserDataSourceModel        `tfsdk:"user"`
}

type siteParentReferenceDataSourceModel struct {
	DriveId       types.String                      `tfsdk:"drive_id"`
	DriveType     types.String                      `tfsdk:"drive_type"`
	Id            types.String                      `tfsdk:"id"`
	Name          types.String                      `tfsdk:"name"`
	Path          types.String                      `tfsdk:"path"`
	ShareId       types.String                      `tfsdk:"share_id"`
	SharepointIds *siteSharepointIdsDataSourceModel `tfsdk:"sharepoint_ids"`
	SiteId        types.String                      `tfsdk:"site_id"`
}

type siteSharepointIdsDataSourceModel struct {
	ListId           types.String `tfsdk:"list_id"`
	ListItemId       types.String `tfsdk:"list_item_id"`
	ListItemUniqueId types.String `tfsdk:"list_item_unique_id"`
	SiteId           types.String `tfsdk:"site_id"`
	SiteUrl          types.String `tfsdk:"site_url"`
	TenantId         types.String `tfsdk:"tenant_id"`
	WebId            types.String `tfsdk:"web_id"`
}

type siteErrorDataSourceModel struct {
	Code       types.String                   `tfsdk:"code"`
	Details    []siteDetailsDataSourceModel   `tfsdk:"details"`
	InnerError *siteInnerErrorDataSourceModel `tfsdk:"inner_error"`
	Message    types.String                   `tfsdk:"message"`
	Target     types.String                   `tfsdk:"target"`
}

type siteDetailsDataSourceModel struct {
	Code    types.String `tfsdk:"code"`
	Message types.String `tfsdk:"message"`
	Target  types.String `tfsdk:"target"`
}

type siteInnerErrorDataSourceModel struct {
	Code    types.String                 `tfsdk:"code"`
	Details []siteDetailsDataSourceModel `tfsdk:"details"`
	Message types.String                 `tfsdk:"message"`
	Target  types.String                 `tfsdk:"target"`
}

type siteRootDataSourceModel struct {
}

type siteSiteCollectionDataSourceModel struct {
	DataLocationCode types.String             `tfsdk:"data_location_code"`
	Hostname         types.String             `tfsdk:"hostname"`
	Root             *siteRootDataSourceModel `tfsdk:"root"`
}

// Read refreshes the Terraform state with the latest data.
func (d *siteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state siteDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
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

	var result models.Siteable
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.Sites().BySiteId(state.Id.ValueString()).Get(context.Background(), &qparams)
	} else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"`id` must be supplied.",
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
		state.Id = types.StringValue(*result.GetId())
	}
	if result.GetCreatedBy() != nil {
		state.CreatedBy = new(siteCreatedByDataSourceModel)

		if result.GetCreatedBy().GetApplication() != nil {
			state.CreatedBy.Application = new(siteApplicationDataSourceModel)

			if result.GetCreatedBy().GetApplication().GetDisplayName() != nil {
				state.CreatedBy.Application.DisplayName = types.StringValue(*result.GetCreatedBy().GetApplication().GetDisplayName())
			}
			if result.GetCreatedBy().GetApplication().GetId() != nil {
				state.CreatedBy.Application.Id = types.StringValue(*result.GetCreatedBy().GetApplication().GetId())
			}
		}
		if result.GetCreatedBy().GetDevice() != nil {
			state.CreatedBy.Device = new(siteDeviceDataSourceModel)

			if result.GetCreatedBy().GetDevice().GetDisplayName() != nil {
				state.CreatedBy.Device.DisplayName = types.StringValue(*result.GetCreatedBy().GetDevice().GetDisplayName())
			}
			if result.GetCreatedBy().GetDevice().GetId() != nil {
				state.CreatedBy.Device.Id = types.StringValue(*result.GetCreatedBy().GetDevice().GetId())
			}
		}
		if result.GetCreatedBy().GetUser() != nil {
			state.CreatedBy.User = new(siteUserDataSourceModel)

			if result.GetCreatedBy().GetUser().GetDisplayName() != nil {
				state.CreatedBy.User.DisplayName = types.StringValue(*result.GetCreatedBy().GetUser().GetDisplayName())
			}
			if result.GetCreatedBy().GetUser().GetId() != nil {
				state.CreatedBy.User.Id = types.StringValue(*result.GetCreatedBy().GetUser().GetId())
			}
		}
	}
	if result.GetCreatedDateTime() != nil {
		state.CreatedDateTime = types.StringValue(result.GetCreatedDateTime().String())
	}
	if result.GetDescription() != nil {
		state.Description = types.StringValue(*result.GetDescription())
	}
	if result.GetETag() != nil {
		state.ETag = types.StringValue(*result.GetETag())
	}
	if result.GetLastModifiedBy() != nil {
		state.LastModifiedBy = new(siteLastModifiedByDataSourceModel)

		if result.GetLastModifiedBy().GetApplication() != nil {
			state.LastModifiedBy.Application = new(siteApplicationDataSourceModel)

			if result.GetLastModifiedBy().GetApplication().GetDisplayName() != nil {
				state.LastModifiedBy.Application.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetApplication().GetDisplayName())
			}
			if result.GetLastModifiedBy().GetApplication().GetId() != nil {
				state.LastModifiedBy.Application.Id = types.StringValue(*result.GetLastModifiedBy().GetApplication().GetId())
			}
		}
		if result.GetLastModifiedBy().GetDevice() != nil {
			state.LastModifiedBy.Device = new(siteDeviceDataSourceModel)

			if result.GetLastModifiedBy().GetDevice().GetDisplayName() != nil {
				state.LastModifiedBy.Device.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetDevice().GetDisplayName())
			}
			if result.GetLastModifiedBy().GetDevice().GetId() != nil {
				state.LastModifiedBy.Device.Id = types.StringValue(*result.GetLastModifiedBy().GetDevice().GetId())
			}
		}
		if result.GetLastModifiedBy().GetUser() != nil {
			state.LastModifiedBy.User = new(siteUserDataSourceModel)

			if result.GetLastModifiedBy().GetUser().GetDisplayName() != nil {
				state.LastModifiedBy.User.DisplayName = types.StringValue(*result.GetLastModifiedBy().GetUser().GetDisplayName())
			}
			if result.GetLastModifiedBy().GetUser().GetId() != nil {
				state.LastModifiedBy.User.Id = types.StringValue(*result.GetLastModifiedBy().GetUser().GetId())
			}
		}
	}
	if result.GetLastModifiedDateTime() != nil {
		state.LastModifiedDateTime = types.StringValue(result.GetLastModifiedDateTime().String())
	}
	if result.GetName() != nil {
		state.Name = types.StringValue(*result.GetName())
	}
	if result.GetParentReference() != nil {
		state.ParentReference = new(siteParentReferenceDataSourceModel)

		if result.GetParentReference().GetDriveId() != nil {
			state.ParentReference.DriveId = types.StringValue(*result.GetParentReference().GetDriveId())
		}
		if result.GetParentReference().GetDriveType() != nil {
			state.ParentReference.DriveType = types.StringValue(*result.GetParentReference().GetDriveType())
		}
		if result.GetParentReference().GetId() != nil {
			state.ParentReference.Id = types.StringValue(*result.GetParentReference().GetId())
		}
		if result.GetParentReference().GetName() != nil {
			state.ParentReference.Name = types.StringValue(*result.GetParentReference().GetName())
		}
		if result.GetParentReference().GetPath() != nil {
			state.ParentReference.Path = types.StringValue(*result.GetParentReference().GetPath())
		}
		if result.GetParentReference().GetShareId() != nil {
			state.ParentReference.ShareId = types.StringValue(*result.GetParentReference().GetShareId())
		}
		if result.GetParentReference().GetSharepointIds() != nil {
			state.ParentReference.SharepointIds = new(siteSharepointIdsDataSourceModel)

			if result.GetParentReference().GetSharepointIds().GetListId() != nil {
				state.ParentReference.SharepointIds.ListId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListId())
			}
			if result.GetParentReference().GetSharepointIds().GetListItemId() != nil {
				state.ParentReference.SharepointIds.ListItemId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListItemId())
			}
			if result.GetParentReference().GetSharepointIds().GetListItemUniqueId() != nil {
				state.ParentReference.SharepointIds.ListItemUniqueId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetListItemUniqueId())
			}
			if result.GetParentReference().GetSharepointIds().GetSiteId() != nil {
				state.ParentReference.SharepointIds.SiteId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetSiteId())
			}
			if result.GetParentReference().GetSharepointIds().GetSiteUrl() != nil {
				state.ParentReference.SharepointIds.SiteUrl = types.StringValue(*result.GetParentReference().GetSharepointIds().GetSiteUrl())
			}
			if result.GetParentReference().GetSharepointIds().GetTenantId() != nil {
				state.ParentReference.SharepointIds.TenantId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetTenantId())
			}
			if result.GetParentReference().GetSharepointIds().GetWebId() != nil {
				state.ParentReference.SharepointIds.WebId = types.StringValue(*result.GetParentReference().GetSharepointIds().GetWebId())
			}
		}
		if result.GetParentReference().GetSiteId() != nil {
			state.ParentReference.SiteId = types.StringValue(*result.GetParentReference().GetSiteId())
		}
	}
	if result.GetWebUrl() != nil {
		state.WebUrl = types.StringValue(*result.GetWebUrl())
	}
	if result.GetDisplayName() != nil {
		state.DisplayName = types.StringValue(*result.GetDisplayName())
	}
	if result.GetError() != nil {
		state.Error = new(siteErrorDataSourceModel)

		if result.GetError().GetCode() != nil {
			state.Error.Code = types.StringValue(*result.GetError().GetCode())
		}
		for _, value := range result.GetError().GetDetails() {
			details := new(siteDetailsDataSourceModel)

			if value.GetCode() != nil {
				details.Code = types.StringValue(*value.GetCode())
			}
			if value.GetMessage() != nil {
				details.Message = types.StringValue(*value.GetMessage())
			}
			if value.GetTarget() != nil {
				details.Target = types.StringValue(*value.GetTarget())
			}
			state.Error.Details = append(state.Error.Details, *details)
		}
		if result.GetError().GetInnerError() != nil {
			state.Error.InnerError = new(siteInnerErrorDataSourceModel)

			if result.GetError().GetInnerError().GetCode() != nil {
				state.Error.InnerError.Code = types.StringValue(*result.GetError().GetInnerError().GetCode())
			}
			for _, value := range result.GetError().GetInnerError().GetDetails() {
				details := new(siteDetailsDataSourceModel)

				if value.GetCode() != nil {
					details.Code = types.StringValue(*value.GetCode())
				}
				if value.GetMessage() != nil {
					details.Message = types.StringValue(*value.GetMessage())
				}
				if value.GetTarget() != nil {
					details.Target = types.StringValue(*value.GetTarget())
				}
				state.Error.InnerError.Details = append(state.Error.InnerError.Details, *details)
			}
			if result.GetError().GetInnerError().GetMessage() != nil {
				state.Error.InnerError.Message = types.StringValue(*result.GetError().GetInnerError().GetMessage())
			}
			if result.GetError().GetInnerError().GetTarget() != nil {
				state.Error.InnerError.Target = types.StringValue(*result.GetError().GetInnerError().GetTarget())
			}
		}
		if result.GetError().GetMessage() != nil {
			state.Error.Message = types.StringValue(*result.GetError().GetMessage())
		}
		if result.GetError().GetTarget() != nil {
			state.Error.Target = types.StringValue(*result.GetError().GetTarget())
		}
	}
	if result.GetIsPersonalSite() != nil {
		state.IsPersonalSite = types.BoolValue(*result.GetIsPersonalSite())
	}
	if result.GetRoot() != nil {
		state.Root = new(siteRootDataSourceModel)

	}
	if result.GetSharepointIds() != nil {
		state.SharepointIds = new(siteSharepointIdsDataSourceModel)

		if result.GetSharepointIds().GetListId() != nil {
			state.SharepointIds.ListId = types.StringValue(*result.GetSharepointIds().GetListId())
		}
		if result.GetSharepointIds().GetListItemId() != nil {
			state.SharepointIds.ListItemId = types.StringValue(*result.GetSharepointIds().GetListItemId())
		}
		if result.GetSharepointIds().GetListItemUniqueId() != nil {
			state.SharepointIds.ListItemUniqueId = types.StringValue(*result.GetSharepointIds().GetListItemUniqueId())
		}
		if result.GetSharepointIds().GetSiteId() != nil {
			state.SharepointIds.SiteId = types.StringValue(*result.GetSharepointIds().GetSiteId())
		}
		if result.GetSharepointIds().GetSiteUrl() != nil {
			state.SharepointIds.SiteUrl = types.StringValue(*result.GetSharepointIds().GetSiteUrl())
		}
		if result.GetSharepointIds().GetTenantId() != nil {
			state.SharepointIds.TenantId = types.StringValue(*result.GetSharepointIds().GetTenantId())
		}
		if result.GetSharepointIds().GetWebId() != nil {
			state.SharepointIds.WebId = types.StringValue(*result.GetSharepointIds().GetWebId())
		}
	}
	if result.GetSiteCollection() != nil {
		state.SiteCollection = new(siteSiteCollectionDataSourceModel)

		if result.GetSiteCollection().GetDataLocationCode() != nil {
			state.SiteCollection.DataLocationCode = types.StringValue(*result.GetSiteCollection().GetDataLocationCode())
		}
		if result.GetSiteCollection().GetHostname() != nil {
			state.SiteCollection.Hostname = types.StringValue(*result.GetSiteCollection().GetHostname())
		}
		if result.GetSiteCollection().GetRoot() != nil {
			state.SiteCollection.Root = new(siteRootDataSourceModel)

		}
	}

	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
