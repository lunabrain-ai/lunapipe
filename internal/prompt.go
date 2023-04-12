package internal

import (
	"bufio"
	"fmt"
	"github.com/lunabrain-ai/aicli/prompts"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"
)

func readStdin() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	var input string

	for scanner.Scan() {
		input += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to read from stdin: %s", err)
	}

	return input, nil
}

func readPipedData() (string, error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}

	if fi.Mode()&os.ModeNamedPipe != 0 {
		log.Debug().Msg("reading piped data")
		return readStdin()
	}
	return "", nil
}

type PromptInput struct {
	Params map[string]string
}

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

	var writer = &strings.Builder{}
	err = tmpl.Execute(writer, PromptInput{
		Params: paramLookup,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to execute template")
	}
	return writer.String(), nil
}

func getPrompt(context *cli.Context, flags Flags) (string, error) {
	var (
		prompt string
	)

	tmplName := context.String("template")
	if tmplName != "" {
		tmplPrompt, err := loadPromptFromTemplate(context, tmplName)
		if err != nil {
			return "", err
		}
		prompt += tmplPrompt
	}

	stdinData, err := readPipedData()
	if err != nil {
		return "", err
	}

	prompt += context.Args().First()
	if prompt == "" {
		if stdinData != "" {
			return "", fmt.Errorf("TODO use piped stdinData and stdin at the same time")
		}

		if !flags.Quiet {
			println("Reading from stdin, close with ctrl+D...")
		}

		var err error
		prompt, err = readStdin()
		if err != nil {
			return "", err
		}
	}

	if stdinData != "" {
		prompt += "\n" + stdinData
	}
	return prompt, nil
}

func loadParams(context *cli.Context) (map[string]string, error) {
	params := context.StringSlice("param")
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

func loadTemplates(promptTmplDir string) (map[string]*template.Template, error) {
	templateLookup := map[string]*template.Template{}

	if promptTmplDir != "" {
		providedTmpls, err := os.ReadDir(promptTmplDir)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read provided template dir: %s", promptTmplDir)
		}

		var files []string
		for _, promptTmpl := range providedTmpls {
			files = append(files, promptTmpl.Name())
		}
		loadAndParseTmpls(os.DirFS(promptTmplDir), files, templateLookup)
	}

	promptTmpls, err := prompts.Templates.ReadDir(".")
	if err != nil {
		return nil, err
	}

	var files []string
	for _, promptTmpl := range promptTmpls {
		files = append(files, promptTmpl.Name())
	}
	loadAndParseTmpls(prompts.Templates, files, templateLookup)

	return templateLookup, nil
}
