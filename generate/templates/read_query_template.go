qparams := {{.Configuration}}RequestBuilderGetRequestConfiguration{
	QueryParameters: &{{.Configuration}}RequestBuilderGetQueryParameters{
		Select: []string {
			{{- range .SelectParameters}}
			"{{.}}",
			{{- end }}
		},
	},
}

{{ define "ZeroParameters" }}
result, err := d.client.{{range .GetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
{{- end}}

{{ define "NonZeroParameters" }}
var result models.{{.BlockName.UpperCamel}}able
var err error

if !state.Id.IsNull() {
	result, err = d.client.{{range .GetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
} {{range .AltGetMethod}} else if !state.{{.if}}.IsNull() {
	result, err = d.client.{{.method}}.Get(context.Background(), &qparams)
} {{end}}else {
	resp.Diagnostics.AddError(
		"Missing argument",
		"`{{.ErrorAttribute}}` {{range .ErrorExtraAttributes}}or `{{.}}` {{end}}must be supplied.",
	)
	return
}
{{- end}}

{{- if not .MultipleGetMethodParameters }}
{{- template "ZeroParameters" .}}
{{- else }}
{{- template "NonZeroParameters" .}}
{{- end}}

if err != nil {
	resp.Diagnostics.AddError(
		"Error getting {{.BlockName.Snake}}",
		err.Error(),
	)
	return
}
