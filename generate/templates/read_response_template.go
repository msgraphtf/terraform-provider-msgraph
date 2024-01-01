{{- /* Define templates for mapping each response type to state */}}
{{- define "ReadStringAttribute" }}
if {{.GetMethod}} != nil {
	{{.StateVarName}} = types.StringValue(*{{.GetMethod}})
} else {
	{{.StateVarName}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringBase64Attribute" }}
if {{.GetMethod}} != nil {
	{{.StateVarName}} = types.StringValue(string({{.GetMethod}}[:]))
} else {
	{{.StateVarName}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringFormattedAttribute" }}
if {{.GetMethod}} != nil {
	{{.StateVarName}} = types.StringValue({{.GetMethod}}.String())
} else {
	{{.StateVarName}} = types.StringNull()
}
{{- end}}

{{- define "ReadInt64Attribute" }}
if {{.GetMethod}} != nil {
	{{ .StateVarName}} = types.Int64Value(int64(*{{.GetMethod}}))
} else {
	{{.StateVarName}} = types.Int64Null()
}
{{- end}}

{{- define "ReadBoolAttribute" }}
if {{.GetMethod}} != nil {
	{{ .StateVarName}} = types.BoolValue(*{{.GetMethod}})
} else {
	{{.StateVarName}} = types.BoolNull()
}
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
if len({{.GetMethod}}) > 0 {
	var {{.ModelVarName}} []attr.Value
	for _, v := range {{.GetMethod}} {
		{{.ModelVarName}} = append({{.ModelVarName}}, types.StringValue(v))
	}
	listValue, _ := types.ListValue(types.StringType, {{.ModelVarName}})
	{{.StateVarName}} = listValue
} else {
	{{.StateVarName}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListStringFormattedAttribute" }}
if len({{.GetMethod}}) > 0 {
	var {{.ModelVarName}} []attr.Value
	for _, v := range {{.GetMethod}} {
		{{.ModelVarName}} = append({{.ModelVarName}}, types.StringValue(v.String()))
	}
	listValue, _ := types.ListValue(types.StringType, {{.ModelVarName}})
	{{.StateVarName}} = listValue
} else {
	{{.StateVarName}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListNestedAttribute" }}
if len({{.GetMethod}}) > 0 {
	objectValues := []basetypes.ObjectValue{}
	for _, v := range {{.GetMethod}} {
		{{.ModelVarName}} := new({{.ModelName}})
			{{template "generate_read" .NestedRead}}
		objectValue, _ := types.ObjectValueFrom(ctx, {{.ModelVarName}}.AttributeTypes(), {{.ModelVarName}})
		objectValues = append(objectValues, objectValue)
	}
{{.StateVarName}}, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
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
