package {{.PackageName}}

import (
    "context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	{{- if .ReadQuery.MultipleGetMethodParameters }}
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	{{- end}}
	"github.com/microsoftgraph/msgraph-sdk-go/{{.PackageName}}"
)

// Ensure the implementation satisfies the expected interfaces.
var (
    _ resource.Resource = &{{.BlockName.LowerCamel}}Resource{}
    _ resource.ResourceWithConfigure = &{{.BlockName.LowerCamel}}Resource{}
)

// New{{.BlockName.UpperCamel}}Resource is a helper function to simplify the provider implementation.
func New{{.BlockName.UpperCamel}}Resource() resource.Resource {
    return &{{.BlockName.LowerCamel}}Resource{}
}

// {{.BlockName.LowerCamel}}Resource is the resource implementation.
type {{.BlockName.LowerCamel}}Resource struct{
	client *msgraphsdk.GraphServiceClient
}

// Metadata returns the resource type name.
func (d *{{.BlockName.LowerCamel}}Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_{{.BlockName.Snake}}"
}

// Configure adds the provider configured client to the resource.
func (d *{{.BlockName.LowerCamel}}Resource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*msgraphsdk.GraphServiceClient)
}

// Schema defines the schema for the resource.
func (d *{{.BlockName.LowerCamel}}Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			{{- template "schema_template.go" .}}
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *{{.BlockName.LowerCamel}}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan {{.BlockName.LowerCamel}}Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Generate API request body from Plan

	// TODO: Create new {{.BlockName.LowerCamel}}

	// TODO: Map response body to schema and populate Computed attribute value

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (d *{{.BlockName.LowerCamel}}Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state {{.BlockName.LowerCamel}}Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	{{ template "read_query_template.go" .ReadQuery}}

	{{ template "read_response_template.go" .ReadResponse}}


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


}

// Update updates the resource and sets the updated Terraform state on success.
func (r *{{.BlockName.LowerCamel}}Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan {{.BlockName.LowerCamel}}Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state {{.BlockName.LowerCamel}}Model
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Generate API request body from plan

	// TODO: Update {{.BlockName.LowerCamel}}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *{{.BlockName.LowerCamel}}Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state {{.BlockName.LowerCamel}}Model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete {{.BlockName.LowerCamel}}
}
