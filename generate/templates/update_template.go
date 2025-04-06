// Update updates the resource and sets the updated Terraform state on success.
func (r *{{.BlockName.LowerCamel}}Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from Terraform plan
	var tfPlan{{.BlockName.UpperCamel}} {{.BlockName.LowerCamel}}Model
	diags := req.Plan.Get(ctx, &tfPlan{{.BlockName.UpperCamel}})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current Terraform state
	var tfState{{.BlockName.UpperCamel}} {{.BlockName.LowerCamel}}Model
	diags = req.State.Get(ctx, &tfState{{.BlockName.UpperCamel}})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	requestBody{{.BlockName.UpperCamel}} := models.New{{.BlockName.UpperCamel}}()

	{{- define "UpdateStringAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	requestBody{{.ParentName}}.Set{{.SetModelMethod}}(&tfPlan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateStringEnumAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	parsed{{.Name}}, _ := models.Parse{{.ObjectOf}}(tfPlan{{.Name}})
	asserted{{.Name}} := parsed{{.Name}}.(models.{{.ObjectOf}})
	requestBody{{.ParentName}}.Set{{.Name}}(&asserted{{.Name}})
	}
	{{- end}}

	{{- define "UpdateStringTimeAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	t, _ := time.Parse(time.RFC3339, tfPlan{{.Name}})
	requestBody{{.ParentName}}.Set{{.Name}}(&t)
	}
	{{- end}}

	{{- define "UpdateStringUuidAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	u, _ := uuid.Parse(tfPlan{{.Name}})
	requestBody{{.ParentName}}.Set{{.Name}}(&u)
	}
	{{- end}}

	{{- define "UpdateStringBase64UrlAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	requestBody{{.ParentName}}.Set{{.SetModelMethod}}([]byte(tfPlan{{.Name}}))
	}
	{{- end}}

	{{- define "UpdateInt64Attribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueInt64()
	requestBody{{.ParentName}}.Set{{.Name}}(&tfPlan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateInt32Attribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := int32(tfPlan{{.ParentName}}.{{.Name}}.ValueInt64())
	requestBody{{.ParentName}}.Set{{.Name}}(&tfPlan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateBoolAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueBool()
	requestBody{{.ParentName}}.Set{{.Name}}(&tfPlan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateArrayStringAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}) {
		var stringArray{{.Name}} []string
		for _, i := range tfPlan{{.ParentName}}.{{.Name}}.Elements() {
			stringArray{{.Name}} = append(stringArray{{.Name}}, i.String())
		}
		requestBody{{.ParentName}}.Set{{.Name}}(stringArray{{.Name}})
	}
	{{- end}}

	{{- define "UpdateArrayUuidAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}) {
		var {{.Name}} []uuid.UUID
		for _, i := range tfPlan{{.ParentName}}.{{.Name}}.Elements() {
			u, _ := uuid.Parse(i.String())
			{{.Name}} = append({{.Name}}, u)
		}
		requestBody{{.ParentName}}.Set{{.Name}}({{.Name}})
	}
	{{- end}}

	{{- define "UpdateArrayObjectAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}) {
		var tfPlan{{.Name}} []models.{{.ObjectOf}}able
		for k, i := range tfPlan{{.ParentName}}.{{.Name}}.Elements() {
			requestBody{{.ObjectOf}} := models.New{{.ObjectOf}}()
			tfPlan{{.ObjectOf}} := {{.TfModelName}}Model{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlan{{.ObjectOf}})
			tfState{{.ObjectOf}} := {{.TfModelName}}Model{}
			types.ListValueFrom(ctx, tfState{{.ParentName}}.{{.Name}}.Elements()[k].Type(ctx), &tfPlan{{.ObjectOf}})
			{{template "generate_update" .NestedUpdate}}
		}
		requestBody{{.ParentName}}.Set{{.Name}}(tfPlan{{.Name}})
	}
	{{- end}}

	{{- define "UpdateObjectAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.Equal(tfState{{.ParentName}}.{{.Name}}){
		requestBody{{.ObjectOf}} := models.New{{.ObjectOf}}()
		tfPlan{{.ObjectOf}} := {{.TfModelName}}Model{}
		tfPlan{{.ParentName}}.{{.Name}}.As(ctx, &tfPlan{{.ObjectOf}}, basetypes.ObjectAsOptions{})
		tfState{{.ObjectOf}} := {{.TfModelName}}Model{}
		tfState{{.ParentName}}.{{.Name}}.As(ctx, &tfState{{.ObjectOf}}, basetypes.ObjectAsOptions{})
		{{template "generate_update" .NestedUpdate}}
		requestBody{{.ParentName}}.Set{{.Name}}(requestBody{{.ObjectOf}})
		tfPlan{{.ParentName}}.{{.Name}}, _ = types.ObjectValueFrom(ctx, tfPlan{{.ObjectOf}}.AttributeTypes(), tfPlan{{.ObjectOf}})
	}
	{{- end}}

	{{- block "generate_update" .Attributes}}
	{{- range .}}
	{{- if eq .Type "UpdateStringAttribute"}}
	{{ template "UpdateStringAttribute" .}}
	{{- else if eq .Type "UpdateStringEnumAttribute"}}
	{{ template "UpdateStringEnumAttribute" .}}
	{{- else if eq .Type "UpdateStringTimeAttribute"}}
	{{ template "UpdateStringTimeAttribute" .}}
	{{- else if eq .Type "UpdateStringUuidAttribute"}}
	{{ template "UpdateStringUuidAttribute" .}}
	{{- else if eq .Type "UpdateStringBase64UrlAttribute"}}
	{{ template "UpdateStringBase64UrlAttribute" .}}
	{{- else if eq .Type "UpdateInt64Attribute"}}
	{{ template "UpdateInt64Attribute" .}}
	{{- else if eq .Type "UpdateInt32Attribute"}}
	{{ template "UpdateInt32Attribute" .}}
	{{- else if eq .Type "UpdateBoolAttribute"}}
	{{ template "UpdateBoolAttribute" .}}
	{{- else if eq .Type "UpdateArrayStringAttribute"}}
	{{ template "UpdateArrayStringAttribute" .}}
	{{- else if eq .Type "UpdateArrayUuidAttribute"}}
	{{ template "UpdateArrayUuidAttribute" .}}
	{{ else if eq .Type "UpdateArrayObjectAttribute" }}
	{{ template "UpdateArrayObjectAttribute" . }}
	{{- else if eq .Type "UpdateObjectAttribute"}}
	{{ template "UpdateObjectAttribute" .}}
	{{- end}}
	{{- end}}
	{{- end}}


	// Update {{.BlockName.LowerCamel}}
	_, err := r.client.{{range .PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Patch(context.Background(), requestBody{{.BlockName.UpperCamel}}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating {{.BlockName.Snake}}",
			err.Error(),
		)
		return
	}

	// Update resource state with Computed values
	diags = resp.State.Set(ctx, tfPlan{{.BlockName.UpperCamel}})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
