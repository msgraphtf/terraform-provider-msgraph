// Create creates the resource and sets the initial Terraform state.
func (r *{{.BlockName.LowerCamel}}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from Terraform plan
	var tfPlan{{.BlockName.UpperCamel}} {{.BlockName.LowerCamel}}Model
	diags := req.Plan.Get(ctx, &tfPlan{{.BlockName.UpperCamel}})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	sdkModel{{.BlockName.UpperCamel}} := models.New{{.BlockName.UpperCamel}}()

	{{- define "CreateStringAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	sdkModel{{.ParentName}}.Set{{.SetModelMethod}}(&tfPlan{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringEnumAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	parsed{{.Name}}, _ := models.Parse{{.SdkModelName}}(tfPlan{{.Name}})
	asserted{{.Name}} := parsed{{.Name}}.(models.{{.SdkModelName}})
	sdkModel{{.ParentName}}.Set{{.Name}}(&asserted{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringTimeAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	t, _ := time.Parse(time.RFC3339, tfPlan{{.Name}})
	sdkModel{{.ParentName}}.Set{{.Name}}(&t)
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringUuidAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	u, _ := uuid.Parse(tfPlan{{.Name}})
	sdkModel{{.ParentName}}.Set{{.Name}}(&u)
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringBase64UrlAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueString()
	sdkModel{{.ParentName}}.Set{{.SetModelMethod}}([]byte(tfPlan{{.Name}}))
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateInt64Attribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueInt64()
	sdkModel{{.ParentName}}.Set{{.Name}}(&tfPlan{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.Int64Null()
	}
	{{- end}}

	{{- define "CreateBoolAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
	tfPlan{{.Name}} := tfPlan{{.ParentName}}.{{.Name}}.ValueBool()
	sdkModel{{.ParentName}}.Set{{.Name}}(&tfPlan{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.BoolNull()
	}
	{{- end}}

	{{- define "CreateArrayStringAttribute" }}
	if len(tfPlan{{.ParentName}}.{{.Name}}.Elements()) > 0 {
		var stringArray{{.Name}} []string
		for _, i := range tfPlan{{.ParentName}}.{{.Name}}.Elements() {
			stringArray{{.Name}} = append(stringArray{{.Name}}, i.String())
		}
		sdkModel{{.ParentName}}.Set{{.Name}}(stringArray{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
	}
	{{- end}}

	{{- define "CreateArrayUuidAttribute" }}
	if len(tfPlan{{.ParentName}}.{{.Name}}.Elements()) > 0 {
		var uuidArray{{.Name}} []uuid.UUID
		for _, i := range tfPlan{{.ParentName}}.{{.Name}}.Elements() {
			u, _ := uuid.Parse(i.String())
			uuidArray{{.Name}} = append(uuidArray{{.Name}}, u)
		}
		sdkModel{{.ParentName}}.Set{{.Name}}(uuidArray{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
	}
	{{- end}}

	{{- define "CreateArrayObjectAttribute" }}
	if len(tfPlan{{.ParentName}}.{{.Name}}.Elements()) > 0 {
		var tfPlan{{.Name}} []models.{{.SdkModelName}}able
		for _, i := range tfPlan{{.ParentName}}.{{.Name}}.Elements() {
			sdkModel{{.Name}} := models.New{{.SdkModelName}}()
			tfPlan{{.Name}} := {{.TfModelName}}Model{}
			types.ListValueFrom(ctx, i.Type(ctx), &tfPlan{{.Name}})
			{{template "generate_create" .NestedCreate}}
		}
		sdkModel{{.ParentSdkModelVarName}}.Set{{.Name}}(tfPlan{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.ListNull(tfPlan{{.ParentName}}.{{.Name}}.ElementType(ctx))
	}
	{{- end}}

	{{- define "CreateObjectAttribute" }}
	if !tfPlan{{.ParentName}}.{{.Name}}.IsUnknown(){
		sdkModel{{.Name}} := models.New{{.SdkModelName}}()
		tfPlan{{.Name}} := {{.TfModelName}}Model{}
		tfPlan{{.ParentName}}.{{.Name}}.As(ctx, &tfPlan{{.Name}}, basetypes.ObjectAsOptions{})
		{{template "generate_create" .NestedCreate}}
		sdkModel{{.ParentSdkModelVarName}}.Set{{.Name}}(sdkModel{{.Name}})
		tfPlan{{.ParentName}}.{{.Name}}, _ = types.ObjectValueFrom(ctx, tfPlan{{.Name}}.AttributeTypes(), sdkModel{{.Name}})
	} else {
		tfPlan{{.ParentName}}.{{.Name}} = types.ObjectNull(tfPlan{{.ParentName}}.{{.Name}}.AttributeTypes(ctx))
	}
	{{- end}}

	{{- block "generate_create" .Attributes}}
	{{- range .}}
	// START {{.Name}} | {{.Type -}}
	{{- if eq .Type "CreateStringAttribute"}}
	{{- template "CreateStringAttribute" .}}
	{{- else if eq .Type "CreateStringEnumAttribute"}}
	{{- template "CreateStringEnumAttribute" .}}
	{{- else if eq .Type "CreateStringTimeAttribute"}}
	{{- template "CreateStringTimeAttribute" .}}
	{{- else if eq .Type "CreateStringUuidAttribute"}}
	{{- template "CreateStringUuidAttribute" .}}
	{{- else if eq .Type "CreateStringBase64UrlAttribute"}}
	{{- template "CreateStringBase64UrlAttribute" .}}
	{{- else if eq .Type "CreateInt64Attribute"}}
	{{- template "CreateInt64Attribute" .}}
	{{- else if eq .Type "CreateInt32Attribute"}}
	{{- template "CreateInt32Attribute" .}}
	{{- else if eq .Type "CreateBoolAttribute"}}
	{{- template "CreateBoolAttribute" .}}
	{{- else if eq .Type "CreateArrayStringAttribute"}}
	{{- template "CreateArrayStringAttribute" .}}
	{{- else if eq .Type "CreateArrayUuidAttribute"}}
	{{- template "CreateArrayUuidAttribute" .}}
	{{ else if eq .Type "CreateArrayObjectAttribute" }}
	{{- template "CreateArrayObjectAttribute" . }}
	{{- else if eq .Type "CreateObjectAttribute"}}
	{{- template "CreateObjectAttribute" .}}
	{{- end}}
	// END {{.Name}} | {{.Type}}
	{{end}}
	{{- end}}

	// Create new {{.BlockName.LowerCamel}}
	result, err := r.client.{{range .PostMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Post(context.Background(), sdkModel{{.BlockName.UpperCamel}}, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating {{.BlockName.Snake}}",
			err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute value
	// TODO: Add support for other Computed values
	tfPlan{{.BlockName.UpperCamel}}.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlan{{.BlockName.UpperCamel}})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

