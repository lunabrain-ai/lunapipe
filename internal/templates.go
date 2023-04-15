package internal

import (
	"fmt"
	"github.com/lunabrain-ai/lunapipe/prompts"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"
	"text/template/parse"
)

func loadPromptFromTemplate(context *cli.Context, tmplName string) (string, error) {
	paramLookup, err := loadParams(context)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load params")
	}

	promptTmplDir := context.String("prompts")

	lookup, err := loadTemplates(promptTmplDir)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load templates")
	}
	var (
		tmpl *template.Template
		ok   bool
	)
	if tmpl, ok = lookup[tmplName]; !ok {
		return "", fmt.Errorf("template %s not found", tmplName)
	}

	params := findIndexCalls(tmpl)
	for _, param := range params {
		if _, ok := paramLookup[param]; ok {
			continue
		}

		if !context.Bool("interact") {
			return "", fmt.Errorf("parameter \"%s\" not set. Please set it to use the template \"%s\"", param, tmplName)
		}

		var p string
		fmt.Printf("Enter value for \"%s\": ", param)
		_, err = fmt.Scanf("%s", &p)
		println()

		if err != nil {
			return "", errors.Wrapf(err, "failed to read param %s", param)
		}
		paramLookup[param] = p
	}

	var writer = &strings.Builder{}
	err = tmpl.Execute(writer, PromptInput{
		Params: paramLookup,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute template")
	}
	return writer.String(), nil
}

func loadTemplates(promptTmplDir string) (map[string]*template.Template, error) {
	templateLookup := map[string]*template.Template{}

	funcMap := map[string]any{
		"readDir": func(dir string) []string {
			matches, err := fs.Glob(os.DirFS(dir), "*")
			if err != nil {
				log.Warn().Err(err).Msg("failed to read dir")
				return []string{}
			}
			return matches
		},
	}
	tmpl := template.New("base").Funcs(funcMap)

	// TODO breadchris duplicate code
	var tmpls []*template.Template
	if promptTmplDir != "" {
		tmpl, err := tmpl.ParseFS(os.DirFS(promptTmplDir), "**/*.tmpl")
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse templates")
		}
		tmpls = append(tmpls, tmpl.Templates()...)
	}

	builtIns, err := tmpl.ParseFS(prompts.Templates, "*.tmpl")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse built-in templates")
	}
	tmpls = append(tmpls, builtIns.Templates()...)

	for _, t := range tmpls {
		// remove extension from template name
		tmplName := t.Name()
		ext := path.Ext(tmplName)
		baseTmplName := tmplName[:len(tmplName)-len(ext)]
		templateLookup[baseTmplName] = t
	}
	return templateLookup, nil
}

func findIndexCalls(tpl *template.Template) []string {
	var indexCalls []string
	ast := tpl.Tree

	var crawlNodes func(node parse.Node)
	crawlNodes = func(node parse.Node) {
		switch n := node.(type) {
		case *parse.ListNode:
			for _, listItem := range n.Nodes {
				crawlNodes(listItem)
			}
		case *parse.ActionNode:
			if n.Pipe != nil {
				for _, cmd := range n.Pipe.Cmds {
					if cmd.Args[0].String() == "index" {
						param := cmd.Args[2].String()
						indexCalls = append(indexCalls, strings.ReplaceAll(param, "\"", ""))
					}
				}
			}
		}
	}
	crawlNodes(ast.Root)
	return indexCalls
}
