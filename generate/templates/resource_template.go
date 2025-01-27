package {{.PackageName}}

import (
    "context"
	"github.com/google/uuid"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	"terraform-provider-msgraph/planmodifiers/listplanmodifiers"
	"terraform-provider-msgraph/planmodifiers/objectplanmodifiers"
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
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringTimeAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	t, _ = time.Parse(time.RFC3339, plan{{.AttributeName.UpperCamel}})
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&t)
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringUuidAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	u, _ = uuid.Parse(plan{{.AttributeName.UpperCamel}})
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&u)
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateInt64Attribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.Int64Null()
	}
	{{- end}}

	{{- define "CreateInt32Attribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := int32({{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}())
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.Int64Null()
	}
	{{- end}}

	{{- define "CreateBoolAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.BoolNull()
	}
	{{- end}}

	{{- define "CreateArrayStringAttribute" }}
	if len({{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements()) > 0 {
		var {{.AttributeName.LowerCamel}} []string
		for _, i := range {{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements() {
			{{.AttributeName.LowerCamel}} = append({{.AttributeName.LowerCamel}}, i.String())
		}
		{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}({{.AttributeName.LowerCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.ListNull(types.StringType)
	}
	{{- end}}

	{{- define "CreateArrayUuidAttribute" }}
	if len({{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements()) > 0 {
		var {{.AttributeName.UpperCamel}} []uuid.UUID
		for _, i := range {{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements() {
			u, _ = uuid.Parse(i.String())
			{{.AttributeName.UpperCamel}} = append({{.AttributeName.UpperCamel}}, u)
		}
		{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}({{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.ListNull(types.StringType)
	}
	{{- end}}

	{{- define "CreateArrayObjectAttribute" }}
	if len({{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements()) > 0 {
		var plan{{.AttributeName.UpperCamel}} []models.{{.NewModelMethod}}able
		for _, i := range {{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements() {
			{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
			{{.RequestBodyVar}}Model := {{.BlockName}}{{.AttributeName.UpperCamel}}Model{}
			types.ListValueFrom(ctx, i.Type(ctx), &{{.RequestBodyVar}}Model)
			{{template "generate_create" .NestedCreate}}
		}
		requestBody.Set{{.AttributeName.UpperCamel}}(plan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.ListNull({{.PlanVar}}{{.AttributeName.UpperCamel}}.ElementType(ctx))
	}
	{{- end}}

	{{- define "CreateObjectAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
		{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
		{{.RequestBodyVar}}Model := {{.BlockName}}{{.AttributeName.UpperCamel}}Model{}
		plan.{{.AttributeName.UpperCamel}}.As(ctx, &{{.RequestBodyVar}}Model, basetypes.ObjectAsOptions{})
		{{template "generate_create" .NestedCreate}}
		requestBody.Set{{.AttributeName.UpperCamel}}({{.RequestBodyVar}})
		objectValue, _ := types.ObjectValueFrom(ctx, {{.RequestBodyVar}}Model.AttributeTypes(), {{.RequestBodyVar}}Model)
		plan.{{.AttributeName.UpperCamel}} = objectValue
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.ObjectNull({{.PlanVar}}{{.AttributeName.UpperCamel}}.AttributeTypes(ctx))
	}
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
	{{- else if eq .AttributeType "CreateInt32Attribute"}}
	{{ template "CreateInt32Attribute" .}}
	{{- else if eq .AttributeType "CreateBoolAttribute"}}
	{{ template "CreateBoolAttribute" .}}
	{{- else if eq .AttributeType "CreateArrayStringAttribute"}}
	{{ template "CreateArrayStringAttribute" .}}
	{{- else if eq .AttributeType "CreateArrayUuidAttribute"}}
	{{ template "CreateArrayUuidAttribute" .}}
	{{ else if eq .AttributeType "CreateArrayObjectAttribute" }}
	{{ template "CreateArrayObjectAttribute" . }}
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

	// Generate API request body from plan
	requestBody := models.New{{.BlockName.UpperCamel}}()
	var t time.Time
	var u uuid.UUID

	{{- define "UpdateStringAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	}
	{{- end}}

	{{- define "UpdateStringTimeAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	t, _ = time.Parse(time.RFC3339, plan{{.AttributeName.UpperCamel}})
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&t)
	}
	{{- end}}

	{{- define "UpdateStringUuidAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	u, _ = uuid.Parse(plan{{.AttributeName.UpperCamel}})
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&u)
	}
	{{- end}}

	{{- define "UpdateInt64Attribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	}
	{{- end}}

	{{- define "UpdateInt32Attribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
	plan{{.AttributeName.UpperCamel}} := int32({{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}())
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	}
	{{- end}}

	{{- define "UpdateBoolAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&plan{{.AttributeName.UpperCamel}})
	}
	{{- end}}

	{{- define "UpdateArrayStringAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}) {
		var {{.AttributeName.LowerCamel}} []string
		for _, i := range {{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements() {
			{{.AttributeName.LowerCamel}} = append({{.AttributeName.LowerCamel}}, i.String())
		}
		{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}({{.AttributeName.LowerCamel}})
	}
	{{- end}}

	{{- define "UpdateArrayUuidAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}) {
		var {{.AttributeName.UpperCamel}} []uuid.UUID
		for _, i := range {{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements() {
			u, _ = uuid.Parse(i.String())
			{{.AttributeName.UpperCamel}} = append({{.AttributeName.UpperCamel}}, u)
		}
		{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}({{.AttributeName.UpperCamel}})
	}
	{{- end}}

	{{- define "UpdateArrayObjectAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}) {
		var plan{{.AttributeName.UpperCamel}} []models.{{.NewModelMethod}}able
		for k, i := range {{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements() {
			{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
			{{.RequestBodyVar}}Model := {{.BlockName}}{{.AttributeName.UpperCamel}}Model{}
			types.ListValueFrom(ctx, i.Type(ctx), &{{.RequestBodyVar}}Model)
			{{.RequestBodyVar}}State := {{.BlockName}}{{.AttributeName.UpperCamel}}Model{}
			types.ListValueFrom(ctx, {{.StateVar}}{{.AttributeName.UpperCamel}}.Elements()[k].Type(ctx), &{{.RequestBodyVar}}Model)
			{{template "generate_update" .NestedUpdate}}
		}
		requestBody.Set{{.AttributeName.UpperCamel}}(plan{{.AttributeName.UpperCamel}})
	}
	{{- end}}

	{{- define "UpdateObjectAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
		{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
		{{.RequestBodyVar}}Model := {{.BlockName}}{{.AttributeName.UpperCamel}}Model{}
		plan.{{.AttributeName.UpperCamel}}.As(ctx, &{{.RequestBodyVar}}Model, basetypes.ObjectAsOptions{})
		{{.RequestBodyVar}}State := {{.BlockName}}{{.AttributeName.UpperCamel}}Model{}
		state.{{.AttributeName.UpperCamel}}.As(ctx, &{{.RequestBodyVar}}State, basetypes.ObjectAsOptions{})
		{{template "generate_update" .NestedUpdate}}
		requestBody.Set{{.AttributeName.UpperCamel}}({{.RequestBodyVar}})
		objectValue, _ := types.ObjectValueFrom(ctx, {{.RequestBodyVar}}Model.AttributeTypes(), {{.RequestBodyVar}}Model)
		plan.{{.AttributeName.UpperCamel}} = objectValue
	}
	{{- end}}

	{{- block "generate_update" .UpdateRequestBody}}
	{{- range .}}
	{{- if eq .AttributeType "UpdateStringAttribute"}}
	{{ template "UpdateStringAttribute" .}}
	{{- else if eq .AttributeType "UpdateStringTimeAttribute"}}
	{{ template "UpdateStringTimeAttribute" .}}
	{{- else if eq .AttributeType "UpdateStringUuidAttribute"}}
	{{ template "UpdateStringUuidAttribute" .}}
	{{- else if eq .AttributeType "UpdateInt64Attribute"}}
	{{ template "UpdateInt64Attribute" .}}
	{{- else if eq .AttributeType "UpdateInt32Attribute"}}
	{{ template "UpdateInt32Attribute" .}}
	{{- else if eq .AttributeType "UpdateBoolAttribute"}}
	{{ template "UpdateBoolAttribute" .}}
	{{- else if eq .AttributeType "UpdateArrayStringAttribute"}}
	{{ template "UpdateArrayStringAttribute" .}}
	{{- else if eq .AttributeType "UpdateArrayUuidAttribute"}}
	{{ template "UpdateArrayUuidAttribute" .}}
	{{ else if eq .AttributeType "UpdateArrayObjectAttribute" }}
	{{ template "UpdateArrayObjectAttribute" . }}
	{{- else if eq .AttributeType "UpdateObjectAttribute"}}
	{{ template "UpdateObjectAttribute" .}}
	{{- end}}
	{{- end}}
	{{- end}}


	// Update {{.BlockName.LowerCamel}}
	_, err := r.client.{{range .UpdateRequest.PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Patch(context.Background(), requestBody, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating {{.BlockName.Snake}}",
			err.Error(),
		)
		return
	}

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
	err := r.client.{{range .UpdateRequest.PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Delete(context.Background(), nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting {{.BlockName.Snake}}",
			err.Error(),
		)
		return
	}

}
