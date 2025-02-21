{{- /* Define templates for different Attribute types */}}
{{- define "StringAttribute" }}
"{{.Name}}": schema.StringAttribute{
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
"{{.Name}}": schema.Int64Attribute{
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
"{{.Name}}": schema.BoolAttribute{
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
"{{.Name}}": schema.ListAttribute{
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
"{{.Name}}": schema.SingleNestedAttribute{
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
	{{- template "generate_schema" .NestedAttribute}}
	},
},
{{- end }}

{{- define "ListNestedAttribute" }}
"{{.Name}}": schema.ListNestedAttribute{
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
			{{- template "generate_schema" .NestedAttribute}}
		},
	},
},
{{- end }}

{{- /* Generate our Attributes from our defined templates above */}}
{{- block "generate_schema" .Schema.Attributes}}
{{- range .}}
{{- if eq .Type "StringAttribute" }}
{{- template "StringAttribute" .}}
{{- else if eq .Type "Int64Attribute" }}
{{- template "Int64Attribute" .}}
{{- else if eq .Type "BoolAttribute" }}
{{- template "BoolAttribute" .}}
{{- else if eq .Type "ListAttribute" }}
{{- template "ListAttribute" .}}
{{- else if eq .Type "SingleNestedAttribute" }}
{{- template "SingleNestedAttribute" .}}
{{- else if eq .Type "ListNestedAttribute" }}
{{- template "ListNestedAttribute" .}}
{{- end }}
{{- end}}
{{- end}}
