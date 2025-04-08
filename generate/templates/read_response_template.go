{{- /* Define templates for mapping each response type to state */}}
{{- define "ReadStringAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(*{{.GetMethod}})
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringBase64Attribute" }}
if {{.GetMethod}} != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(string({{.GetMethod}}[:]))
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringFormattedAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue({{.GetMethod}}.String())
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadInt64Attribute" }}
if {{.GetMethod}} != nil {
	{{.ParentName}}.{{.Name}} = types.Int64Value(int64(*{{.GetMethod}}))
} else {
	tfState{{.ParentName}}.{{.Name}} = types.Int64Null()
}
{{- end}}

{{- define "ReadBoolAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{ .ParentName}}.{{.Name}} = types.BoolValue(*{{.GetMethod}})
} else {
	tfState{{.ParentName}}.{{.Name}} = types.BoolNull()
}
{{- end}}

{{- define "ReadSingleNestedAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{.ObjectOf}} := {{.TfModelName}}Model{}
	{{template "generate_read" .NestedRead}}

	tfState{{.ParentName}}.{{.Name}}, _ = types.ObjectValueFrom(ctx, tfState{{.ObjectOf}}.AttributeTypes(), tfState{{.ObjectOf}})
}
{{- end}}

{{- define "ReadListStringAttribute" }}
if len({{.GetMethod}}) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, result{{.Name}} := range {{.GetMethod}} {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(result{{.Name}}))
	}
	listValue, _ := types.ListValue(types.StringType, valueArray{{.Name}})
	tfState{{.ParentName}}.{{.Name}} = listValue
} else {
	tfState{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListStringFormattedAttribute" }}
if len({{.GetMethod}}) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, result{{.Name}} := range {{.GetMethod}} {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(result{{.Name}}.String()))
	}
	tfState{{.ParentName}}.{{.Name}}, _ = types.ListValue(types.StringType, valueArray{{.Name}})
} else {
	tfState{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListNestedAttribute" }}
if len({{.GetMethod}}) > 0 {
	objectValues := []basetypes.ObjectValue{}
	for _, result{{.Name}} := range {{.GetMethod}} {
		tfState{{.ObjectOf}} := {{.TfModelName}}Model{}
			{{template "generate_read" .NestedRead}}
		objectValue, _ := types.ObjectValueFrom(ctx, tfState{{.ObjectOf}}.AttributeTypes(), tfState{{.ObjectOf}})
		objectValues = append(objectValues, objectValue)
	}
tfState{{.ParentName}}.{{.Name}}, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
}
{{- end}}


{{/* Generate statements to map response to state */}}
{{- block "generate_read" .Attributes}}
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
