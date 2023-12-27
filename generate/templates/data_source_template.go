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


	qparams := {{.PackageName}}.{{.ReadQuery.Configuration}}RequestBuilderGetRequestConfiguration{
		QueryParameters: &{{.PackageName}}.{{.ReadQuery.Configuration}}RequestBuilderGetQueryParameters{
			Select: []string {
				{{- range .ReadQuery.SelectParameters}}
				"{{.}}",
				{{- end }}
			},
		},
	}

	{{ define "ReadQuery.ZeroParameters" }}
	result, err := d.client.{{range .ReadQuery.GetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
	{{- end}}

	{{ define "ReadQuery.NonZeroParameters" }}
	var result models.{{.DataSourceName.UpperCamel}}able
	var err error

	if !state.Id.IsNull() {
		result, err = d.client.{{range .ReadQuery.GetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
	} {{range .ReadQuery.AltGetMethod}} else if !state.{{.if}}.IsNull() {
		result, err = d.client.{{.method}}.Get(context.Background(), &qparams)
	} {{end}}else {
		resp.Diagnostics.AddError(
			"Missing argument",
			"`{{.ReadQuery.ErrorAttribute}}` {{range .ReadQuery.ErrorExtraAttributes}}or `{{.}}` {{end}}must be supplied.",
		)
		return
	}
	{{- end}}

	{{- if not .ReadQuery.MultipleGetMethodParameters }}
	{{- template "ReadQuery.ZeroParameters" .}}
	{{- else }}
	{{- template "ReadQuery.NonZeroParameters" .}}
	{{- end}}

	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting {{.DataSourceName.Snake}}",
			err.Error(),
		)
		return
	}

	{{- template "read_response_template.go" .ReadResponse}}


	// Overwrite items with refreshed state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


}
