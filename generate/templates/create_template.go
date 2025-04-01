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
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.SdkModelVarName}}.Set{{.SetModelMethod}}(&tfPlan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringEnumAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	parsed{{.AttributeName.UpperCamel}}, _ := models.Parse{{.NewModelMethod}}(tfPlan{{.AttributeName.UpperCamel}})
	asserted{{.AttributeName.UpperCamel}} := parsed{{.AttributeName.UpperCamel}}.(models.{{.NewModelMethod}})
	{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}(&asserted{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringTimeAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	t, _ := time.Parse(time.RFC3339, tfPlan{{.AttributeName.UpperCamel}})
	{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}(&t)
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringUuidAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	u, _ := uuid.Parse(tfPlan{{.AttributeName.UpperCamel}})
	{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}(&u)
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateStringBase64UrlAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.SdkModelVarName}}.Set{{.SetModelMethod}}([]byte(tfPlan{{.AttributeName.UpperCamel}}))
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.StringNull()
	}
	{{- end}}

	{{- define "CreateInt64Attribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}(&tfPlan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.Int64Null()
	}
	{{- end}}

	{{- define "CreateInt32Attribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := int32({{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}())
	{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}(&tfPlan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.Int64Null()
	}
	{{- end}}

	{{- define "CreateBoolAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
	tfPlan{{.AttributeName.UpperCamel}} := {{.PlanVar}}{{.AttributeName.UpperCamel}}.{{.PlanValueMethod}}()
	{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}(&tfPlan{{.AttributeName.UpperCamel}})
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
		{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}({{.AttributeName.LowerCamel}})
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
		{{.SdkModelVarName}}.Set{{.AttributeName.UpperCamel}}({{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.ListNull(types.StringType)
	}
	{{- end}}

	{{- define "CreateArrayObjectAttribute" }}
	if len({{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements()) > 0 {
		var tfPlan{{.AttributeName.UpperCamel}} []models.{{.NewModelMethod}}able
		for _, i := range {{.PlanVar}}{{.AttributeName.UpperCamel}}.Elements() {
			{{.SdkModelVarName}} := models.New{{.NewModelMethod}}()
			{{.SdkModelVarName}}Model := {{.ModelName}}{}
			types.ListValueFrom(ctx, i.Type(ctx), &{{.SdkModelVarName}}Model)
			{{template "generate_create" .NestedCreate}}
		}
		{{.ParentSdkModelVarName}}.Set{{.AttributeName.UpperCamel}}(tfPlan{{.AttributeName.UpperCamel}})
	} else {
		{{.PlanVar}}{{.AttributeName.UpperCamel}} = types.ListNull({{.PlanVar}}{{.AttributeName.UpperCamel}}.ElementType(ctx))
	}
	{{- end}}

	{{- define "CreateObjectAttribute" }}
	if !{{.PlanVar}}{{.AttributeName.UpperCamel}}.IsUnknown(){
		{{.SdkModelVarName}} := models.New{{.NewModelMethod}}()
		{{.SdkModelVarName}}Model := {{.ModelName}}{}
		{{.NestedPlan}}.As(ctx, &{{.SdkModelVarName}}Model, basetypes.ObjectAsOptions{})
		{{template "generate_create" .NestedCreate}}
		{{.ParentSdkModelVarName}}.Set{{.AttributeName.UpperCamel}}({{.SdkModelVarName}})
		objectValue, _ := types.ObjectValueFrom(ctx, {{.SdkModelVarName}}Model.AttributeTypes(), {{.SdkModelVarName}}Model)
		{{.ParentPlanVar}} = objectValue
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
	{{- else if eq .AttributeType "CreateStringBase64UrlAttribute"}}
	{{ template "CreateStringBase64UrlAttribute" .}}
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

