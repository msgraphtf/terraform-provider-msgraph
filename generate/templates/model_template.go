package {{.PackageName}}

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

{{- range .Model}}
type {{.ModelName}} struct {
{{- range .ModelFields}}
{{.FieldName}} {{.FieldType}} `tfsdk:"{{.AttributeName}}"`
{{- end}}
}
{{end}}
