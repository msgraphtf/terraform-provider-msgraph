{{- /* Define templates for different Attribute types */}}
{{- define "StringAttribute" }}
"{{.AttributeName}}": schema.StringAttribute{
	Description: "{{.Description}}",
	{{- if .Required}}
	Required: true,
	{{- end}}
	{{- if .Optional}}
	Optional: true,
	{{- end}}
	{{- if .Computed}}
	Computed: true,
	{{- end}}
	{{- if .PlanModifiers}}
	PlanModifiers: []planmodifier.String{
		stringplanmodifiers.UseStateForUnconfigured(),
	},
	{{- end}}
},
{{- end }}

{{- define "Int64Attribute" }}
"{{.AttributeName}}": schema.Int64Attribute{
	Description: "{{.Description}}",
	{{- if .Required}}
	Required: true,
	{{- end}}
	{{- if .Optional}}
	Optional: true,
	{{- end}}
	{{- if .Computed}}
	Computed: true,
	{{- end}}
},
{{- end }}

{{- define "BoolAttribute" }}
"{{.AttributeName}}": schema.BoolAttribute{
	Description: "{{.Description}}",
	{{- if .Required}}
	Required: true,
	{{- end}}
	{{- if .Optional}}
	Optional: true,
	{{- end}}
	{{- if .Computed}}
	Computed: true,
	{{- end}}
	{{- if .PlanModifiers}}
	PlanModifiers: []planmodifier.Bool{
		boolplanmodifiers.UseStateForUnconfigured(),
	},
	{{- end}}
},
{{- end }}

{{- define "ListAttribute" }}
"{{.AttributeName}}": schema.ListAttribute{
	Description: "{{.Description}}",
	{{- if .Required}}
	Required: true,
	{{- end}}
	{{- if .Optional}}
	Optional: true,
	{{- end}}
	{{- if .Computed}}
	Computed: true,
	{{- end}}
	{{- if .PlanModifiers}}
	PlanModifiers: []planmodifier.List{
		listplanmodifiers.UseStateForUnconfigured(),
	},
	{{- end}}
	ElementType: {{.ElementType}},
},
{{- end }}

{{- define "SingleNestedAttribute" }}
"{{.AttributeName}}": schema.SingleNestedAttribute{
	Description: "{{.Description}}",
	{{- if .Required}}
	Required: true,
	{{- end}}
	{{- if .Optional}}
	Optional: true,
	{{- end}}
	{{- if .Computed}}
	Computed: true,
	{{- end}}
	{{- if .PlanModifiers}}
	PlanModifiers: []planmodifier.Object{
		objectplanmodifiers.UseStateForUnconfigured(),
	},
	{{- end}}
	Attributes: map[string]schema.Attribute{
	{{- template "generate_schema" .Attributes}}
	},
},
{{- end }}

{{- define "ListNestedAttribute" }}
"{{.AttributeName}}": schema.ListNestedAttribute{
	Description: "{{.Description}}",
	{{- if .Required}}
	Required: true,
	{{- end}}
	{{- if .Optional}}
	Optional: true,
	{{- end}}
	{{- if .Computed}}
	Computed: true,
	{{- end}}
	{{- if .PlanModifiers}}
	PlanModifiers: []planmodifier.List{
		listplanmodifiers.UseStateForUnconfigured(),
	},
	{{- end}}
	NestedObject: schema.NestedAttributeObject{
		Attributes: map[string]schema.Attribute{
			{{- template "generate_schema" .NestedObject}}
		},
	},
},
{{- end }}

{{- /* Generate our Attributes from our defined templates above */}}
{{- block "generate_schema" .Schema}}
{{- range .}}
{{- if eq .AttributeType "StringAttribute" }}
{{- template "StringAttribute" .}}
{{- else if eq .AttributeType "Int64Attribute" }}
{{- template "Int64Attribute" .}}
{{- else if eq .AttributeType "BoolAttribute" }}
{{- template "BoolAttribute" .}}
{{- else if eq .AttributeType "ListAttribute" }}
{{- template "ListAttribute" .}}
{{- else if eq .AttributeType "SingleNestedAttribute" }}
{{- template "SingleNestedAttribute" .}}
{{- else if eq .AttributeType "ListNestedAttribute" }}
{{- template "ListNestedAttribute" .}}
{{- end }}
{{- end}}
{{- end}}
