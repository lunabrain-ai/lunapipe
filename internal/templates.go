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
			return "", fmt.Errorf("param %s not found", param)
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

	// TODO breadchris duplicate code
	if promptTmplDir != "" {
		providedTmpls, err := os.ReadDir(promptTmplDir)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read provided template dir: %s", promptTmplDir)
		}

		var (
			files  []string
			blocks []string
		)
		for _, promptTmpl := range providedTmpls {
			if promptTmpl.IsDir() {
				if promptTmpl.Name() == "blocks" {
					blocks = append(blocks, promptTmpl.Name())
				}
				continue
			}
			files = append(files, promptTmpl.Name())
		}
		loadAndParseTmpls(os.DirFS(promptTmplDir), files, templateLookup)
	}

	promptTmpls, err := prompts.Templates.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var (
		files  []string
		blocks []string
	)
	for _, promptTmpl := range promptTmpls {
		if promptTmpl.IsDir() {
			if promptTmpl.Name() == "blocks" {
				blocks = append(blocks, promptTmpl.Name())
			}
			continue
		}
		files = append(files, promptTmpl.Name())
	}
	loadAndParseTmpls(prompts.Templates, files, templateLookup)

	return templateLookup, nil
}

func loadAndParseTmpls(filesys fs.FS, files []string, templateLookup map[string]*template.Template) {
	for _, file := range files {
		tmplData, err := fs.ReadFile(filesys, file)
		if err != nil {
			log.Warn().
				Err(err).
				Str("template", file).
				Msg("failed to read template")
			continue
		}
		t, err := template.New(file).Parse(string(tmplData))
		if err != nil {
			log.Warn().
				Err(err).
				Str("template", file).
				Msg("failed to parse template")
			continue
		}

		baseTmplName := file[:len(file)-len(path.Ext(file))]

		templateLookup[baseTmplName] = t
	}
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
