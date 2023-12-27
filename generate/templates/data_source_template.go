{{- template "data_source_preamble.go" .}}

// Schema defines the schema for the data source.
func (d *{{.DataSourceName.LowerCamel}}DataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			{{- template "schema_template.go" .}}
		},
	}
}

{{/* Generate data models from provided data */}}
{{- range .Model}}
type {{.ModelName}} struct {
{{- range .Fields}}
{{.FieldName}} {{.FieldType}} `tfsdk:"{{.AttributeName}}"`
{{- end}}
}
{{end}}

// Read refreshes the Terraform state with the latest data.
func (d *{{.DataSourceName.LowerCamel}}DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state {{.DataSourceName.LowerCamel}}DataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}


	qparams := {{.PackageName}}.{{.ReadQueryConfiguration}}RequestBuilderGetRequestConfiguration{
		QueryParameters: &{{.PackageName}}.{{.ReadQueryConfiguration}}RequestBuilderGetQueryParameters{
			Select: []string {
				{{- range .ReadQuerySelectParameters}}
				"{{.}}",
				{{- end }}
			},
		},
	}

	{{ define "ReadQueryZeroParameters" }}
	result, err := d.client.{{range .ReadQueryGetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
	{{- end}}

	{{ define "ReadQueryNonZeroParameters" }}
	var result models.{{.DataSourceName.UpperCamel}}able
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.{{range .ReadQueryGetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
	} {{range .ReadQueryAltGetMethod}} else if !state.{{.if}}.IsNull() {
		result, err = d.client.{{.method}}.Get(context.Background(), &qparams)
	} {{end}}else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"`{{.ReadQueryErrorAttribute}}` {{range .ReadQueryErrorExtraAttributes}}or `{{.}}` {{end}}must be supplied.",
		)
		return
	}
	{{- end}}

	{{- if eq .ReadQueryGetMethodParametersCount 0}}
	{{- template "ReadQueryZeroParameters" .}}
	{{- else if gt .ReadQueryGetMethodParametersCount 0 }}
	{{- template "ReadQueryNonZeroParameters" .}}
	{{- end}}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting {{.DataSourceName.Snake}}",
			err.Error(),
		)
		return
	}

	{{- /* Define templates for mapping each response type to state */}}
	{{- define "ReadStringAttribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.StringValue(*{{.GetMethod}})}
	{{- end}}

	{{- define "ReadStringBase64Attribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.StringValue(string({{.GetMethod}}[:]))}
	{{- end}}

	{{- define "ReadStringFormattedAttribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.StringValue({{.GetMethod}}.String())}
	{{- end}}

	{{- define "ReadInt64Attribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.Int64Value(int64(*{{.GetMethod}}))}
	{{- end}}

	{{- define "ReadBoolAttribute" }}
	if {{.GetMethod}}  != nil { {{- .StateVarName}} = types.BoolValue(*{{.GetMethod}})}
	{{- end}}

	{{- define "ReadSingleNestedAttribute" }}
	if {{.GetMethod}} != nil {
		{{.StateVarName}} = new({{.ModelName}})
		{{template "generate_read" .NestedRead}}
	}
	{{- end}}

	{{- define "ReadListStringAttribute" }}
	for _, v := range {{.GetMethod}} {
		{{.StateVarName}} = append({{.StateVarName}}, types.StringValue(v))
	}
	{{- end}}

	{{- define "ReadListStringFormattedAttribute" }}
	for _, v := range {{.GetMethod}} {
		{{.StateVarName}} = append({{.StateVarName}}, types.StringValue(v.String()))
	}
	{{- end}}

	{{- define "ReadListNestedAttribute" }}
	for _, v := range {{.GetMethod}} {
		{{.ModelVarName}} := new({{.ModelName}})
			{{template "generate_read" .NestedRead}}
		{{.StateVarName}} = append({{.StateVarName}}, *{{.ModelVarName}})
	}
	{{- end}}


	{{/* Generate statements to map response to state */}}
	{{- block "generate_read" .Read}}
	{{- range .}}
	{{- if eq .AttributeType "ReadStringAttribute"}}
	{{- template "ReadStringAttribute" .}}
	{{- else if eq .AttributeType "ReadStringBase64Attribute"}}
	{{- template "ReadStringBase64Attribute" .}}
	{{- else if eq .AttributeType "ReadStringFormattedAttribute"}}
	{{- template "ReadStringFormattedAttribute" .}}
	{{- else if eq .AttributeType "ReadInt64Attribute"}}
	{{- template "ReadInt64Attribute" .}}
	{{- else if eq .AttributeType "ReadBoolAttribute"}}
	{{- template "ReadBoolAttribute" .}}
	{{- else if eq .AttributeType "ReadListStringAttribute"}}
	{{- template "ReadListStringAttribute" .}}
	{{- else if eq .AttributeType "ReadListStringFormattedAttribute"}}
	{{- template "ReadListStringFormattedAttribute" .}}
	{{- else if eq .AttributeType "ReadSingleNestedAttribute"}}
	{{- template "ReadSingleNestedAttribute" .}}
	{{- else if eq .AttributeType "ReadListNestedAttribute"}}
	{{- template "ReadListNestedAttribute" .}}
	{{- end}}
	{{- end}}
	{{- end}}


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


}
