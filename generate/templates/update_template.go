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

	{{- define "UpdateStringEnumAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
	plan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	parsed{{.AttributeName.UpperCamel}}, _ := models.Parse{{.NewModelMethod}}(plan{{.AttributeName.UpperCamel}})
	asserted{{.AttributeName.UpperCamel}} := parsed{{.AttributeName.UpperCamel}}.(models.{{.NewModelMethod}})
	{{.RequestBodyVar}}.Set{{.AttributeName.UpperCamel}}(&asserted{{.AttributeName.UpperCamel}})
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
			{{.RequestBodyVar}}Model := {{.ModelName}}{}
			types.ListValueFrom(ctx, i.Type(ctx), &{{.RequestBodyVar}}Model)
			{{.RequestBodyVar}}State := {{.ModelName}}{}
			types.ListValueFrom(ctx, {{.StateVar}}{{.AttributeName.UpperCamel}}.Elements()[k].Type(ctx), &{{.RequestBodyVar}}Model)
			{{template "generate_update" .NestedUpdate}}
		}
		requestBody.Set{{.AttributeName.UpperCamel}}(plan{{.AttributeName.UpperCamel}})
	}
	{{- end}}

	{{- define "UpdateObjectAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.Equal({{.StateVar}}{{.AttributeName.UpperCamel}}){
		{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
		{{.RequestBodyVar}}Model := {{.ModelName}}{}
		plan.{{.AttributeName.UpperCamel}}.As(ctx, &{{.RequestBodyVar}}Model, basetypes.ObjectAsOptions{})
		{{.RequestBodyVar}}State := {{.ModelName}}{}
		state.{{.AttributeName.UpperCamel}}.As(ctx, &{{.RequestBodyVar}}State, basetypes.ObjectAsOptions{})
		{{template "generate_update" .NestedUpdate}}
		requestBody.Set{{.AttributeName.UpperCamel}}({{.RequestBodyVar}})
		objectValue, _ := types.ObjectValueFrom(ctx, {{.RequestBodyVar}}Model.AttributeTypes(), {{.RequestBodyVar}}Model)
		plan.{{.AttributeName.UpperCamel}} = objectValue
	}
	{{- end}}

	{{- block "generate_update" .Attributes}}
	{{- range .}}
	{{- if eq .AttributeType "UpdateStringAttribute"}}
	{{ template "UpdateStringAttribute" .}}
	{{- else if eq .AttributeType "UpdateStringEnumAttribute"}}
	{{ template "UpdateStringEnumAttribute" .}}
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
	_, err := r.client.{{range .PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Patch(context.Background(), requestBody, nil)
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
