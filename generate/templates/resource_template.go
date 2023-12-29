package {{.PackageName}}

import (
    "context"
	"time"
	"uuid"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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

	var t time.Time
	var u uuid.UUID

	// Generate API request body from Plan
	requestBody := models.New{{.BlockName.UpperCamel}}()

	{{- define "CreateStringAttribute" }}
	{{.PlanValueVar}} := {{.PlanVar}}{{.PlanFields}}.ValueString()
	{{.RequestBodyVar}}.Set{{.PlanSetMethod}}(&{{.PlanValueVar}})
	{{- end}}

	{{- define "CreateStringTimeAttribute" }}
	{{.PlanValueVar}} := {{.PlanVar}}{{.PlanFields}}.ValueString()
	t, _ = time.Parse(time.RFC3339, {{.PlanValueVar}})
	{{.RequestBodyVar}}.Set{{.PlanSetMethod}}(&t)
	{{- end}}

	{{- define "CreateStringUuidAttribute" }}
	{{.PlanValueVar}} := {{.PlanVar}}{{.PlanFields}}.ValueString()
	u, _ = uuid.Parse({{.PlanValueVar}})
	{{.RequestBodyVar}}.Set{{.PlanSetMethod}}(&u)
	{{- end}}

	{{- define "CreateInt64Attribute" }}
	{{.PlanValueVar}} := {{.PlanVar}}{{.PlanFields}}.ValueInt64()
	{{.RequestBodyVar}}.Set{{.PlanSetMethod}}(&{{.PlanValueVar}})
	{{- end}}

	{{- define "CreateBoolAttribute" }}
	{{.PlanValueVar}} := {{.PlanVar}}{{.PlanFields}}.ValueBool()
	{{.RequestBodyVar}}.Set{{.PlanSetMethod}}(&{{.PlanValueVar}})
	{{- end}}

	{{- define "CreateArrayStringAttribute" }}
	var {{.PlanValueVar}} []string
	for _, i := range {{.PlanVar}}{{.PlanFields}} {
		{{.PlanValueVar}} = append({{.PlanValueVar}}, i.ValueString())
	}
	{{.RequestBodyVar}}.Set{{.PlanSetMethod}}({{.PlanValueVar}})
	{{- end}}

	{{- define "CreateArrayUuidAttribute" }}
	var {{.PlanValueVar}} []uuid.UUID
	for _, i := range {{.PlanVar}}{{.PlanFields}} {
		{{.PlanValueVar}} = append({{.PlanValueVar}}, i.ValueString())
	}
	{{.RequestBodyVar}}.Set{{.PlanSetMethod}}({{.PlanValueVar}})
	{{- end}}

	{{- define "CreateArrayObjectAttribute" }}
	var {{.PlanValueVar}} []models.{{.NewModelMethod}}able
	for _, i := range {{.PlanVar}}{{.PlanFields}} {
		{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
		{{template "generate_create" .NestedCreate}}
	}
	requestBody.Set{{.PlanSetMethod}}({{.PlanValueVar}})
	{{- end}}

	{{- define "CreateObjectAttribute" }}
	{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
	{{template "generate_create" .NestedCreate}}
	requestBody.Set{{.SetModelMethod}}({{.RequestBodyVar}})
	{{- end}}

	{{- block "generate_create" .CreateRequestBody}}
	{{- range .}}
	{{- if eq .AttributeType "CreateStringAttribute"}}
	{{ template "CreateStringAttribute" .}}
	{{- else if eq .AttributeType "CreateStringTimeAttribute"}}
	{{ template "CreateStringTimeAttribute" .}}
	{{- else if eq .AttributeType "CreateStringUuidAttribute"}}
	{{ template "CreateStringUuidAttribute" .}}
	{{- else if eq .AttributeType "CreateInt64Attribute"}}
	{{ template "CreateInt64Attribute" .}}
	{{- else if eq .AttributeType "CreateBoolAttribute"}}
	{{ template "CreateBoolAttribute" .}}
	{{- else if eq .AttributeType "CreateArrayStringAttribute"}}
	{{ template "CreateArrayStringAttribute" .}}
	{{- else if eq .AttributeType "CreateArrayUuidAttribute"}}
	{{ template "CreateArrayUuidAttribute" .}}
	{{- else if eq .AttributeType "CreateArrayObjectAttribute"}}
	{{ template "CreateArrayObjectAttribute" .}}
	{{- else if eq .AttributeType "CreateObjectAttribute"}}
	{{ template "CreateObjectAttribute" .}}
	{{- end}}
	{{- end}}
	{{- end}}

	// Create new {{.BlockName.LowerCamel}}
	result, err := r.client.{{range .CreateRequest.PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Post(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating {{.BlockName.Snake}}",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	plan.Id = types.StringValue(*result.GetId())

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
	requestBody := models.New{{.BlockName.UpperCamel}}()
	var t time.Time
	var u uuid.UUID

	{{template "generate_create" .CreateRequestBody}}

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
