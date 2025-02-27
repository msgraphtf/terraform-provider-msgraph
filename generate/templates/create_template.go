// Create creates the resource and sets the initial Terraform state.
func (r *{{.BlockName.LowerCamel}}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan {{.BlockName.LowerCamel}}Model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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

	{{- define "CreateStringEnumAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	parsed{{.AttributeName.UpperCamel}}, _ := models.Parse{{.NewModelMethod}}(plan{{.AttributeName.UpperCamel}})
	asserted{{.AttributeName.UpperCamel}} := parsed{{.AttributeName.UpperCamel}}.(models.{{.NewModelMethod}})
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&asserted{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringTimeAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	t, _ := time.Parse(time.RFC3339, plan{{.AttributeName.UpperCamel}})
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&t)
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringUuidAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	u, _ := uuid.Parse(plan{{.AttributeName.UpperCamel}})
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
			u, _ := uuid.Parse(i.String())
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
			{{.RequestBodyVar}}Model := {{.ModelName}}{}
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
		{{.RequestBodyVar}}Model := {{.ModelName}}{}
		plan.{{.AttributeName.UpperCamel}}.As(ctx, &{{.RequestBodyVar}}Model, basetypes.ObjectAsOptions{})
		{{template "generate_create" .NestedCreate}}
		requestBody.Set{{.AttributeName.UpperCamel}}({{.RequestBodyVar}})
		objectValue, _ := types.ObjectValueFrom(ctx, {{.RequestBodyVar}}Model.AttributeTypes(), {{.RequestBodyVar}}Model)
		plan.{{.AttributeName.UpperCamel}} = objectValue
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.ObjectNull({{.PlanVar}}{{.AttributeName.UpperCamel}}.AttributeTypes(ctx))
	}
	{{- end}}

	{{- block "generate_create" .Attributes}}
	{{- range .}}
	{{- if eq .AttributeType "CreateStringAttribute"}}
	{{ template "CreateStringAttribute" .}}
	{{- else if eq .AttributeType "CreateStringEnumAttribute"}}
	{{ template "CreateStringEnumAttribute" .}}
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
	result, err := r.client.{{range .PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Post(context.Background(), requestBody, nil)
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

