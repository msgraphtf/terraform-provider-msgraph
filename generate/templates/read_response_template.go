{{- /* Define templates for mapping each response{{.ParentName}}. type to state */}}
{{- define "ReadStringAttribute" }}
if response{{.ParentName}}.Get{{.GetMethod}}() != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(*response{{.ParentName}}.Get{{.GetMethod}}())
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringBase64Attribute" }}
if response{{.ParentName}}.Get{{.GetMethod}}() != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(string(response{{.ParentName}}.Get{{.GetMethod}}()[:]))
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringFormattedAttribute" }}
if response{{.ParentName}}.Get{{.GetMethod}}() != nil {
	tfState{{.ParentName}}.{{.Name}} = types.StringValue(response{{.ParentName}}.Get{{.GetMethod}}().String())
} else {
	tfState{{.ParentName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadInt64Attribute" }}
if response{{.ParentName}}.Get{{.GetMethod}}() != nil {
	{{.ParentName}}.{{.Name}} = types.Int64Value(int64(*response{{.ParentName}}.Get{{.GetMethod}}()))
} else {
	tfState{{.ParentName}}.{{.Name}} = types.Int64Null()
}
{{- end}}

{{- define "ReadBoolAttribute" }}
if response{{.ParentName}}.Get{{.GetMethod}}() != nil {
	tfState{{ .ParentName}}.{{.Name}} = types.BoolValue(*response{{.ParentName}}.Get{{.GetMethod}}())
} else {
	tfState{{.ParentName}}.{{.Name}} = types.BoolNull()
}
{{- end}}

{{- define "ReadSingleNestedAttribute" }}
if response{{.ParentName}}.Get{{.GetMethod}}() != nil {
	tfState{{.ObjectOf}} := {{.TfModelName}}Model{}
	response{{.ObjectOf}} := response{{.ParentName}}.Get{{.GetMethod}}()
	{{template "generate_read" .NestedRead}}

	tfState{{.ParentName}}.{{.Name}}, _ = types.ObjectValueFrom(ctx, tfState{{.ObjectOf}}.AttributeTypes(), tfState{{.ObjectOf}})
}
{{- end}}

{{- define "ReadListStringAttribute" }}
if len(response{{.ParentName}}.Get{{.GetMethod}}()) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, response{{.Name}} := range response{{.ParentName}}.Get{{.GetMethod}}() {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(response{{.Name}}))
	}
	listValue, _ := types.ListValue(types.StringType, valueArray{{.Name}})
	tfState{{.ParentName}}.{{.Name}} = listValue
} else {
	tfState{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListStringFormattedAttribute" }}
if len(response{{.ParentName}}.Get{{.GetMethod}}()) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, response{{.Name}} := range response{{.ParentName}}.Get{{.GetMethod}}() {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(response{{.Name}}.String()))
	}
	tfState{{.ParentName}}.{{.Name}}, _ = types.ListValue(types.StringType, valueArray{{.Name}})
} else {
	tfState{{.ParentName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListNestedAttribute" }}
if len(response{{.ParentName}}.Get{{.GetMethod}}()) > 0 {
	objectValues := []basetypes.ObjectValue{}
	for _, response{{.ObjectOf}} := range response{{.ParentName}}.Get{{.GetMethod}}() {
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
{{- if eq .Type "ReadStringAttribute"}}
{{- template "ReadStringAttribute" .}}
{{- else if eq .Type "ReadStringBase64Attribute"}}
{{- template "ReadStringBase64Attribute" .}}
{{- else if eq .Type "ReadStringFormattedAttribute"}}
{{- template "ReadStringFormattedAttribute" .}}
{{- else if eq .Type "ReadInt64Attribute"}}
{{- template "ReadInt64Attribute" .}}
{{- else if eq .Type "ReadBoolAttribute"}}
{{- template "ReadBoolAttribute" .}}
{{- else if eq .Type "ReadListStringAttribute"}}
{{- template "ReadListStringAttribute" .}}
{{- else if eq .Type "ReadListStringFormattedAttribute"}}
{{- template "ReadListStringFormattedAttribute" .}}
{{- else if eq .Type "ReadSingleNestedAttribute"}}
{{- template "ReadSingleNestedAttribute" .}}
{{- else if eq .Type "ReadListNestedAttribute"}}
{{- template "ReadListNestedAttribute" .}}
{{- end}}
{{- end}}
{{- end}}
