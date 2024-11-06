{{.data.type}}{{if .data.scope}}({{.data.scope}}){{end}}{{if .data.breakingChange}}!{{end}}: {{.data.description}}{{if .data.breakingMessage}}

BREAKING CHANGE: {{.data.breakingMessage}}{{end}}