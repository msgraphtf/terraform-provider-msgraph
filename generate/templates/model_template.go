package {{.PackageName}}

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

{{- range .Model}}
type {{.ModelName}} struct {
{{- range .ModelFields}}
{{.FieldName}} {{.FieldType}} `tfsdk:"{{.AttributeName}}"`
{{- end}}
}


func (m {{.ModelName}}) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		{{- range .ModelFields}}
		{{- if .IfPrimitiveType }}
		"{{.AttributeName}}": {{.AttributeType}},
		{{- end}}
		{{- end}}
	}
}

{{end}}
