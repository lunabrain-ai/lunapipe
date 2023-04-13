package internal

import (
	"testing"
	"text/template"
)

func TestFindIndexCalls(t *testing.T) {
	tmplText := `
{{ $lang := index .Params "language" }}
{{ $type := index .Params "type" }}
Return a {{ if $type}}{{$type}}{{end}}{{if $lang }}, written in {{ $lang }},{{end}} based on the following requirements:
`
	tmpl, err := template.New("test").Parse(tmplText)
	if err != nil {
		t.Fatal(err)
	}
	calls := FindIndexCalls(tmpl)
	println(calls)
}
