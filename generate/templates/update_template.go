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
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := {{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.SetModelMethod}}(&plan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateStringEnumAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := {{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}()
	parsed{{.Name}}, _ := models.Parse{{.NewModelMethod}}(plan{{.Name}})
	asserted{{.Name}} := parsed{{.Name}}.(models.{{.NewModelMethod}})
	{{.RequestBodyVar}}.Set{{.Name}}(&asserted{{.Name}})
	}
	{{- end}}

	{{- define "UpdateStringTimeAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := {{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}()
	t, _ := time.Parse(time.RFC3339, plan{{.Name}})
	{{.RequestBodyVar}}.Set{{.Name}}(&t)
	}
	{{- end}}

	{{- define "UpdateStringUuidAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := {{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}()
	u, _ := uuid.Parse(plan{{.Name}})
	{{.RequestBodyVar}}.Set{{.Name}}(&u)
	}
	{{- end}}

	{{- define "UpdateStringBase64UrlAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := {{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.SetModelMethod}}([]byte(plan{{.Name}}))
	}
	{{- end}}

	{{- define "UpdateInt64Attribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := {{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.Name}}(&plan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateInt32Attribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := int32({{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}())
	{{.RequestBodyVar}}.Set{{.Name}}(&plan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateBoolAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
	plan{{.Name}} := {{.PlanVar}}{{.Name}}.{{.PlanValueMethod}}()
	{{.RequestBodyVar}}.Set{{.Name}}(&plan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateArrayStringAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}) {
		var stringArray{{.Name}} []string
		for _, i := range {{.PlanVar}}{{.Name}}.Elements() {
			stringArray{{.Name}} = append(stringArray{{.Name}}, i.String())
		}
		{{.RequestBodyVar}}.Set{{.Name}}(stringArray{{.Name}})
	}
	{{- end}}

	{{- define "UpdateArrayUuidAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}) {
		var {{.Name}} []uuid.UUID
		for _, i := range {{.PlanVar}}{{.Name}}.Elements() {
			u, _ := uuid.Parse(i.String())
			{{.Name}} = append({{.Name}}, u)
		}
		{{.RequestBodyVar}}.Set{{.Name}}({{.Name}})
	}
	{{- end}}

	{{- define "UpdateArrayObjectAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}) {
		var plan{{.Name}} []models.{{.NewModelMethod}}able
		for k, i := range {{.PlanVar}}{{.Name}}.Elements() {
			{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
			{{.RequestBodyVar}}Model := {{.ModelName}}{}
			types.ListValueFrom(ctx, i.Type(ctx), &{{.RequestBodyVar}}Model)
			{{.RequestBodyVar}}State := {{.ModelName}}{}
			types.ListValueFrom(ctx, {{.StateVar}}{{.Name}}.Elements()[k].Type(ctx), &{{.RequestBodyVar}}Model)
			{{template "generate_update" .NestedUpdate}}
		}
		{{.ParentRequestBodyVar}}.Set{{.Name}}(plan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateObjectAttribute" }}
	if !{{.PlanVar}}{{.Name}}.Equal({{.StateVar}}{{.Name}}){
		{{.RequestBodyVar}} := models.New{{.NewModelMethod}}()
		{{.RequestBodyVar}}Model := {{.ModelName}}{}
		{{.NestedPlan}}.As(ctx, &{{.RequestBodyVar}}Model, basetypes.ObjectAsOptions{})
		{{.RequestBodyVar}}State := {{.ModelName}}{}
		{{.NestedState}}.As(ctx, &{{.RequestBodyVar}}State, basetypes.ObjectAsOptions{})
		{{template "generate_update" .NestedUpdate}}
		{{.ParentRequestBodyVar}}.Set{{.Name}}({{.RequestBodyVar}})
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
