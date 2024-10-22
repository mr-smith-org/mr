package helpers

import (
	"strings"
	"text/template"

	"github.com/arthurbcp/kuma/v2/pkg/functions"
)

func ReplaceVars(text string, vars interface{}, funcs template.FuncMap) (string, error) {
	t, err := template.New("").Funcs(functions.GetFuncMap()).Parse(text)
	if err != nil {
		return "", err
	}
	var buf strings.Builder
	err = t.Execute(&buf, vars)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
