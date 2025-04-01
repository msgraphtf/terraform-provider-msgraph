// Create creates the resource and sets the initial Terraform state.
func (r *{{.BlockName.LowerCamel}}Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from Terraform plan
	var tfPlan {{.BlockName.LowerCamel}}Model
	diags := req.Plan.Get(ctx, &tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from Terraform plan
	sdkModel{{.BlockName.UpperCamel}} := models.New{{.BlockName.UpperCamel}}()

	{{- define "CreateStringAttribute" }}
	if !{{.PlanVar}}.IsUnknown(){
	tfPlan{{.Name}} := {{.PlanVar}}.ValueString()
	{{.SdkModelVarName}}.Set{{.SetModelMethod}}(&tfPlan{{.Name}})
	} else {
		{{.PlanVar}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringEnumAttribute" }}
	if !{{.PlanVar}}.IsUnknown(){
	tfPlan{{.Name}} := {{.PlanVar}}.ValueString()
	parsed{{.Name}}, _ := models.Parse{{.SdkModelName}}(tfPlan{{.Name}})
	asserted{{.Name}} := parsed{{.Name}}.(models.{{.SdkModelName}})
	{{.SdkModelVarName}}.Set{{.Name}}(&asserted{{.Name}})
	} else {
		{{.PlanVar}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringTimeAttribute" }}
	if !{{.PlanVar}}.IsUnknown(){
	tfPlan{{.Name}} := {{.PlanVar}}.ValueString()
	t, _ := time.Parse(time.RFC3339, tfPlan{{.Name}})
	{{.SdkModelVarName}}.Set{{.Name}}(&t)
	} else {
		{{.PlanVar}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringUuidAttribute" }}
	if !{{.PlanVar}}.IsUnknown(){
	tfPlan{{.Name}} := {{.PlanVar}}.ValueString()
	u, _ := uuid.Parse(tfPlan{{.Name}})
	{{.SdkModelVarName}}.Set{{.Name}}(&u)
	} else {
		{{.PlanVar}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringBase64UrlAttribute" }}
	if !{{.PlanVar}}.IsUnknown(){
	tfPlan{{.Name}} := {{.PlanVar}}.ValueString()
	{{.SdkModelVarName}}.Set{{.SetModelMethod}}([]byte(tfPlan{{.Name}}))
	} else {
		{{.PlanVar}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateInt64Attribute" }}
	if !{{.PlanVar}}.IsUnknown(){
	tfPlan{{.Name}} := {{.PlanVar}}.ValueInt64()
	{{.SdkModelVarName}}.Set{{.Name}}(&tfPlan{{.Name}})
	} else {
		{{.PlanVar}} = types.Int64Null()
	}
	{{- end}}

	{{- define "CreateBoolAttribute" }}
	if !{{.PlanVar}}.IsUnknown(){
	tfPlan{{.Name}} := {{.PlanVar}}.ValueBool()
	{{.SdkModelVarName}}.Set{{.Name}}(&tfPlan{{.Name}})
	} else {
		{{.PlanVar}} = types.BoolNull()
	}
	{{- end}}

	{{- define "CreateArrayStringAttribute" }}
	if len({{.PlanVar}}.Elements()) > 0 {
		var stringArray{{.Name}} []string
		for _, i := range {{.PlanVar}}.Elements() {
			stringArray{{.Name}} = append(stringArray{{.Name}}, i.String())
		}
		{{.SdkModelVarName}}.Set{{.Name}}(stringArray{{.Name}})
	} else {
		{{.PlanVar}} = types.ListNull(types.StringType)
	}
	{{- end}}

	{{- define "CreateArrayUuidAttribute" }}
	if len({{.PlanVar}}.Elements()) > 0 {
		var uuidArray{{.Name}} []uuid.UUID
		for _, i := range {{.PlanVar}}.Elements() {
			u, _ := uuid.Parse(i.String())
			uuidArray{{.Name}} = append(uuidArray{{.Name}}, u)
		}
		{{.SdkModelVarName}}.Set{{.Name}}(uuidArray{{.Name}})
	} else {
		{{.PlanVar}} = types.ListNull(types.StringType)
	}
	{{- end}}

	{{- define "CreateArrayObjectAttribute" }}
	if len({{.PlanVar}}.Elements()) > 0 {
		var tfPlan{{.Name}} []models.{{.SdkModelName}}able
		for _, i := range {{.PlanVar}}.Elements() {
			{{.SdkModelVarName}} := models.New{{.SdkModelName}}()
			{{.TfPlanVarName}} := {{.TfModelName}}Model{}
			types.ListValueFrom(ctx, i.Type(ctx), &{{.TfPlanVarName}})
			{{template "generate_create" .NestedCreate}}
		}
		{{.ParentSdkModelVarName}}.Set{{.Name}}(tfPlan{{.Name}})
	} else {
		{{.PlanVar}} = types.ListNull({{.PlanVar}}.ElementType(ctx))
	}
	{{- end}}

	{{- define "CreateObjectAttribute" }}
	if !{{.PlanVar}}.IsUnknown(){
		{{.SdkModelVarName}} := models.New{{.SdkModelName}}()
		{{.TfPlanVarName}} := {{.TfModelName}}Model{}
		{{.PlanVar}}.As(ctx, &{{.TfPlanVarName}}, basetypes.ObjectAsOptions{})
		{{template "generate_create" .NestedCreate}}
		{{.ParentSdkModelVarName}}.Set{{.Name}}({{.SdkModelVarName}})
		{{.ParentPlanVar}}, _ = types.ObjectValueFrom(ctx, {{.TfPlanVarName}}.AttributeTypes(), {{.SdkModelVarName}})
	} else {
		{{.PlanVar}} = types.ObjectNull({{.PlanVar}}.AttributeTypes(ctx))
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
	tfPlan.Id = types.StringValue(*result.GetId())

	// Set state to fully populated data
	diags = resp.State.Set(ctx, tfPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

