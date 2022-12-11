// Code generated by codegen1; DO NOT EDIT.
package recipes

import "{{.PkgPath}}"

func init() {
	tables := []Table{
    {{- range .Tables}}
		{
			Service: "{{$.BaseImport}}",
      Name: "{{.Name}}",
      Struct: &{{$.BaseImport}}.{{.Struct}}{},
      ResponseStruct: &{{$.BaseImport}}.{{.ResponseStruct}}{},
      Client: &{{$.BaseImport}}.{{.Client}}{},
      ListFunc: (&{{$.BaseImport}}.{{.Client}}{}).{{.ListFunc}},
			NewFunc: {{$.BaseImport}}.{{.NewFunc}},
			URL: "{{.URL}}",
			Multiplex: `{{.Multiplex}}`,
      {{- if .ExtraColumns}}
      ExtraColumns: {{.ExtraColumns}},
      {{- end }}
		},
    {{- end}}
	}
  Tables = append(Tables, tables...)
}