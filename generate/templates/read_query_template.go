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
response{{.BlockName.UpperCamel}}, err := d.client.{{range .GetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
{{- end}}

{{ define "NonZeroParameters" }}
var response{{.BlockName.UpperCamel}} models.{{.BlockName.UpperCamel}}able
var err error

if !tfState{{.BlockName.UpperCamel}}.Id.IsNull() {
	response{{.BlockName.UpperCamel}}, err = d.client.{{range .GetMethod}}{{.MethodName}}({{.Parameter}}).{{end}}Get(context.Background(), &qparams)
} {{range .AltGetMethod}} else if !tfState{{.BlockName.UpperCamel}}.{{.if}}.IsNull() {
	response{{.BlockName.UpperCamel}}, err = d.client.{{.method}}.Get(context.Background(), &qparams)
} {{end}}else {
	resp.Diagnostics.AddError(
		"Missing argument",
		"TODO: Specify required parameters",
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
