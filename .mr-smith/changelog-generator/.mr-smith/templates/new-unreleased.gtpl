
{{- $u := groupByKey .data.unreleased.data "type" }}
{{ range $k, $c := $u }}
### {{$k}}
    {{ range $c -}}
    - {{.description}}
    {{end -}}
{{end}}
