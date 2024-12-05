
## [{{.data.releaseVersion}}] - {{now | date "2006-01-02"}}
{{- $u := groupByKey .data.releaseData.data "type" }}
{{ range $k, $c := $u }}
### {{$k}}
    {{ range $c -}}
    - {{.description}}
    {{end -}}
{{end -}}
