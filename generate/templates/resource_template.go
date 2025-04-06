package {{.PackageName}}

import (
    "context"
	{{- if .CreateRequest.IfUuidImportNeeded }}
	"github.com/google/uuid"
	{{- end}}
	"time"

	{{- if .ReadResponse.IfAttrImportNeeded }}
	"github.com/hashicorp/terraform-plugin-framework/attr"
	{{- end}}
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	{{- if .ReadQuery.MultipleGetMethodParameters }}
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	{{- end}}
	"github.com/microsoftgraph/msgraph-sdk-go/{{.PackageName}}"

	"terraform-provider-msgraph/planmodifiers/boolplanmodifiers"
	{{- if .Schema.IfListPlanModifiersImportNeeded }}
	"terraform-provider-msgraph/planmodifiers/listplanmodifiers"
	{{- end}}
	{{- if .Schema.IfSingleNestedAttributeUsed nil }}
	"terraform-provider-msgraph/planmodifiers/objectplanmodifiers"
	{{- end}}
	"terraform-provider-msgraph/planmodifiers/stringplanmodifiers"
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

{{template "create_template.go" .CreateRequest}}

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

{{template "update_template.go" .UpdateRequest}}

// Delete deletes the resource and removes the Terraform state on success.
func (r *{{.BlockName.LowerCamel}}Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from Terraform state
	var tfState {{.BlockName.LowerCamel}}Model
	diags := req.State.Get(ctx, &tfState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: Delete {{.BlockName.LowerCamel}}
	err := r.client.{{range .UpdateRequest.PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting {{.BlockName.Snake}}",
			err.Error(),
		)
		return
	}

}
