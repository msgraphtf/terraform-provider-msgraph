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

	{{- define "UpdateStringAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := {{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.SetModelMethod}}(&plan{{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateStringEnumAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := {{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}()
	parsed{{.AttributeName}}, _ := models.Parse{{.NewModelMethod}}(plan{{.AttributeName}})
	asserted{{.AttributeName}} := parsed{{.AttributeName}}.(models.{{.NewModelMethod}})
	{{.RequestBodyVar}}.Set{{.AttributeName}}(&asserted{{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateStringTimeAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := {{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}()
	t, _ := time.Parse(time.RFC3339, plan{{.AttributeName}})
	{{.RequestBodyVar}}.Set{{.AttributeName}}(&t)
	}
	{{- end}}

	{{- define "UpdateStringUuidAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := {{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}()
	u, _ := uuid.Parse(plan{{.AttributeName}})
	{{.RequestBodyVar}}.Set{{.AttributeName}}(&u)
	}
	{{- end}}

	{{- define "UpdateStringBase64UrlAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := {{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.SetModelMethod}}([]byte(plan{{.AttributeName}}))
	}
	{{- end}}

	{{- define "UpdateInt64Attribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := {{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName}}(&plan{{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateInt32Attribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := int32({{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}())
	{{.RequestBodyVar}}.Set{{.AttributeName}}(&plan{{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateBoolAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
	plan{{.AttributeName}} := {{.PlanVar}}{{.AttributeName}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.AttributeName}}(&plan{{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateArrayStringAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}) {
		var stringArray{{.AttributeName}} []string
		for _, i := range {{.PlanVar}}{{.AttributeName}}.Elements() {
			stringArray{{.AttributeName}} = append(stringArray{{.AttributeName}}, i.String())
		}
		{{.RequestBodyVar}}.Set{{.AttributeName}}(stringArray{{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateArrayUuidAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}) {
		var {{.AttributeName}} []uuid.UUID
		for _, i := range {{.PlanVar}}{{.AttributeName}}.Elements() {
			u, _ := uuid.Parse(i.String())
			{{.AttributeName}} = append({{.AttributeName}}, u)
		}
		{{.RequestBodyVar}}.Set{{.AttributeName}}({{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateArrayObjectAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}) {
		var plan{{.AttributeName}} []models.{{.NewModelMethod}}able
		for k, i := range {{.PlanVar}}{{.AttributeName}}.Elements() {
			{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
			{{.RequestBodyVar}}Model := {{.ModelName}}{}
			types.ListValueFrom(ctx, i.Type(ctx), &{{.RequestBodyVar}}Model)
			{{.RequestBodyVar}}State := {{.ModelName}}{}
			types.ListValueFrom(ctx, {{.StateVar}}{{.AttributeName}}.Elements()[k].Type(ctx), &{{.RequestBodyVar}}Model)
			{{template "generate_update" .NestedUpdate}}
		}
		{{.ParentRequestBodyVar}}.Set{{.AttributeName}}(plan{{.AttributeName}})
	}
	{{- end}}

	{{- define "UpdateObjectAttribute" }}
	if !{{.PlanVar}}{{.AttributeName}}.Equal({{.StateVar}}{{.AttributeName}}){
		{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
		{{.RequestBodyVar}}Model := {{.ModelName}}{}
		{{.NestedPlan}}.As(ctx, &{{.RequestBodyVar}}Model, basetypes.ObjectAsOptions{})
		{{.RequestBodyVar}}State := {{.ModelName}}{}
		{{.NestedState}}.As(ctx, &{{.RequestBodyVar}}State, basetypes.ObjectAsOptions{})
		{{template "generate_update" .NestedUpdate}}
		{{.ParentRequestBodyVar}}.Set{{.AttributeName}}({{.RequestBodyVar}})
		objectValue, _ := types.ObjectValueFrom(ctx, {{.RequestBodyVar}}Model.AttributeTypes(), {{.RequestBodyVar}}Model)
		{{.ParentPlanVar}} = objectValue
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
	{{- else if eq .AttributeType "UpdateStringBase64UrlAttribute"}}
	{{ template "UpdateStringBase64UrlAttribute" .}}
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
