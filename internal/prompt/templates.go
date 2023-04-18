package prompt

import (
	"fmt"
	"github.com/lunabrain-ai/lunapipe/internal/util"
	"github.com/lunabrain-ai/lunapipe/prompts"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"strings"
	"text/template"
)

type PromptInput struct {
	Params map[string]string
}

func loadPromptFromTemplate(
	tmplName string,
	promptTmplDir string,
	params []string,
	interact bool,
) (string, error) {
	paramLookup, err := loadParamLookup(params)
	if err != nil {
		return "", errors.Wrapf(err, "failed to load params")
	}

	tmpl, err := locateTemplate(tmplName, promptTmplDir, paramLookup)
	if err != nil {
		return "", err
	}

	err = initParamLookupForTmpl(paramLookup, interact, tmpl)
	if err != nil {
		return "", err
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

func initParamLookupForTmpl(
	paramLookup map[string]string,
	interact bool,
	tmpl *template.Template,
) error {
	tmplParams := util.FindIndexCalls(tmpl)
	for _, param := range tmplParams {
		if _, ok := paramLookup[param]; ok {
			continue
		}

		if !interact {
			return fmt.Errorf("parameter \"%s\" not set. Please set it to use the template \"%s\"", param, tmpl.Name())
		}

		var p string
		fmt.Printf("Enter value for \"%s\": ", param)
		_, err := fmt.Scanf("%s", &p)
		if err != nil {
			return errors.Wrapf(err, "failed to read param %s", param)
		}
		paramLookup[param] = p
	}
	return nil
}

func loadParamLookup(params []string) (map[string]string, error) {
	paramLookup := map[string]string{}
	for _, param := range params {
		splitParam := strings.Split(param, "=")
		if len(splitParam) != 2 {
			return nil, fmt.Errorf("invalid parameter %s", param)
		}
		paramName := strings.ToLower(splitParam[0])
		paramLookup[paramName] = splitParam[1]
	}
	return paramLookup, nil
}

func locateTemplate(
	tmplName string,
	promptTmplDir string,
	paramLookup map[string]string,
) (*template.Template, error) {
	lookup, err := loadTemplates(paramLookup, promptTmplDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load templates")
	}
	var (
		tmpl *template.Template
		ok   bool
	)
	if tmpl, ok = lookup[tmplName]; !ok {
		// see if tmplName is a path to a template
		tmpl, err = template.New(tmplName).ParseFiles(tmplName)
		if err != nil {
			log.Debug().Err(err).Str("tmplName", tmplName).Msg("failed to parse template %s")
			return nil, fmt.Errorf("template %s not found", tmplName)
		}
	}
	return tmpl, nil
}

func loadTemplates(paramLookup map[string]string, promptTmplDir string) (map[string]*template.Template, error) {
	templateLookup := map[string]*template.Template{}

	tmpl := template.New("base").Funcs(NewFuncMap(paramLookup))

	patterns := []string{"*.tmpl", "**/*.tmpl"}

	var tmpls []*template.Template
	if promptTmplDir != "" {
		tmpl, err := tmpl.ParseFS(os.DirFS(promptTmplDir), patterns...)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse templates")
		}
		tmpls = append(tmpls, tmpl.Templates()...)
	}

	builtIns, err := tmpl.ParseFS(prompts.Templates, patterns...)
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
