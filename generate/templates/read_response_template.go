{{- /* Define templates for mapping each response type to state */}}
{{- define "ReadStringAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{.StateVarName}}.{{.Name}} = types.StringValue(*{{.GetMethod}})
} else {
	tfState{{.StateVarName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringBase64Attribute" }}
if {{.GetMethod}} != nil {
	tfState{{.StateVarName}}.{{.Name}} = types.StringValue(string({{.GetMethod}}[:]))
} else {
	tfState{{.StateVarName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadStringFormattedAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{.StateVarName}}.{{.Name}} = types.StringValue({{.GetMethod}}.String())
} else {
	tfState{{.StateVarName}}.{{.Name}} = types.StringNull()
}
{{- end}}

{{- define "ReadInt64Attribute" }}
if {{.GetMethod}} != nil {
	{{.StateVarName}}.{{.Name}} = types.Int64Value(int64(*{{.GetMethod}}))
} else {
	tfState{{.StateVarName}}.{{.Name}} = types.Int64Null()
}
{{- end}}

{{- define "ReadBoolAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{ .StateVarName}}.{{.Name}} = types.BoolValue(*{{.GetMethod}})
} else {
	tfState{{.StateVarName}}.{{.Name}} = types.BoolNull()
}
{{- end}}

{{- define "ReadSingleNestedAttribute" }}
if {{.GetMethod}} != nil {
	tfState{{.Name}} := {{.TfModelName}}Model{}
	{{template "generate_read" .NestedRead}}

	tfState{{.StateVarName}}.{{.Name}}, _ = types.ObjectValueFrom(ctx, tfState{{.Name}}.AttributeTypes(), tfState{{.Name}})
}
{{- end}}

{{- define "ReadListStringAttribute" }}
if len({{.GetMethod}}) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, v := range {{.GetMethod}} {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(v))
	}
	listValue, _ := types.ListValue(types.StringType, valueArray{{.Name}})
	tfState{{.StateVarName}}.{{.Name}} = listValue
} else {
	tfState{{.StateVarName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListStringFormattedAttribute" }}
if len({{.GetMethod}}) > 0 {
	var valueArray{{.Name}} []attr.Value
	for _, v := range {{.GetMethod}} {
		valueArray{{.Name}} = append(valueArray{{.Name}}, types.StringValue(v.String()))
	}
	tfState{{.StateVarName}}.{{.Name}}, _ = types.ListValue(types.StringType, valueArray{{.Name}})
} else {
	tfState{{.StateVarName}}.{{.Name}} = types.ListNull(types.StringType)
}
{{- end}}

{{- define "ReadListNestedAttribute" }}
if len({{.GetMethod}}) > 0 {
	objectValues := []basetypes.ObjectValue{}
	for _, v := range {{.GetMethod}} {
		tfState{{.Name}} := {{.TfModelName}}Model{}
			{{template "generate_read" .NestedRead}}
		objectValue, _ := types.ObjectValueFrom(ctx, tfState{{.Name}}.AttributeTypes(), tfState{{.Name}})
		objectValues = append(objectValues, objectValue)
	}
tfState{{.StateVarName}}.{{.Name}}, _ = types.ListValueFrom(ctx, objectValues[0].Type(ctx), objectValues)
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
