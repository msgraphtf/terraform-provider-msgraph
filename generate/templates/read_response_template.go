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
	{{.ModelVarName}} := new({{.ModelName}})
	{{template "generate_read" .NestedRead}}

	objectValue, _ := types.ObjectValueFrom(ctx, {{.ModelVarName}}.AttributeTypes(), {{.ModelVarName}})
	{{.StateVarName}} = objectValue
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
	objectValue, _ := types.ObjectValueFrom(ctx, {{.ModelVarName}}.AttributeTypes(), {{.ModelVarName}})
	{{.StateVarName}} = append({{.StateVarName}}, objectValue)
}
{{- end}}


{{/* Generate statements to map response to state */}}
{{- block "generate_read" .}}
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
