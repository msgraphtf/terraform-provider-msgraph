{{- /* Define templates for mapping each response type to state */}}
{{- define "ReadStringAttribute" }}
if response{{.GetMethod}} != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(*response{{.GetMethod}})
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringBase64Attribute" }}
if response{{.GetMethod}} != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(string(response{{.GetMethod}}[:]))
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringFormattedAttribute" }}
if response{{.GetMethod}} != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(response{{.GetMethod}}.String())
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadInt64Attribute" }}
if response{{.GetMethod}} != nil {
	{{.ParentName}}.{{.Name}} = types.Int64Value(int64(*response{{.GetMethod}}))
} else {
	tfState{{.ParentName}}.{{.Name}} = types.Int64Null()
}
{{- end}}

{{- define "ReadBoolAttribute" }}
if response{{.GetMethod}} != nil {
	tfState{{ .ParentName}}.{{.Name}} = types.BoolValue(*response{{.GetMethod}})
} else {
	tfState{{.ParentName}}.{{.Name}} = types.BoolNull()
}
{{- end}}

{{- define "ReadSingleNestedAttribute" }}
if response{{.GetMethod}} != nil {
	tfState{{.ObjectOf}} := {{.TfModelName}}Model{}
	{{template "generate_read" .NestedRead}}

	tfState{{.ParentName}}.{{.Name}}, _ = types.ObjectValueFrom(ctx, tfState{{.ObjectOf}}.AttributeTypes(), tfState{{.ObjectOf}})
}
{{- end}}

{{- define "ReadListStringAttribute" }}
if len(response{{.GetMethod}}) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, response{{.Name}} := range response{{.GetMethod}} {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(response{{.Name}}))
	}
	listValue, _ := types.ListValue(types.StringType, valueArray{{.Name}})
	tfState{{.ParentName}}.{{.Name}} = listValue
} else {
	tfState{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListStringFormattedAttribute" }}
if len(response{{.GetMethod}}) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, response{{.Name}} := range response{{.GetMethod}} {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(response{{.Name}}.String()))
	}
	tfState{{.ParentName}}.{{.Name}}, _ = types.ListValue(types.StringType, valueArray{{.Name}})
} else {
	tfState{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListNestedAttribute" }}
if len(response{{.GetMethod}}) > 0 {
	objectValues := []basetypes.ObjectValue{}
	for _, response{{.ObjectOf}} := range response{{.GetMethod}} {
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
