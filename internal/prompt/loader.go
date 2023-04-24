package prompt

import (
	"fmt"
)

type Loader struct {
	prompt   string
	params   []string
	tmplName string
	tmplDir  string

	// Synchronous output, not streaming
	sync bool

	// Do not print help text
	quiet bool

	// For a template, interactively prompt for parameters
	interact bool
}

func NewLoader(
	prompt string,
	params []string,
	tmplName string,
	tmplDir string,
	sync bool,
	quiet bool,
	interact bool,
) *Loader {
	return &Loader{
		prompt:   prompt,
		params:   params,
		tmplName: tmplName,
		tmplDir:  tmplDir,
		sync:     sync,
		quiet:    quiet,
		interact: interact,
	}
}

func (s *Loader) Create() (string, error) {
	var (
		createdPrompt string
	)

	if s.tmplName != "" {
		tmplPrompt, err := loadPromptFromTemplate(
			s.tmplName,
			s.tmplDir,
			s.params,
			s.interact,
		)
		if err != nil {
			return "", err
		}
		createdPrompt += tmplPrompt
	}

	stdinData, err := readPipedData()
	if err != nil {
		return "", err
	}

	createdPrompt += s.prompt
	if s.prompt == "" {
		if stdinData != "" {
			return "", fmt.Errorf("Pass instructions to lunapipe so that it knows how to interpret what you piped in.")
		}

		if !s.quiet {
			println("Reading from stdin, close with ctrl+D...")
		}

		var err error
		createdPrompt, err = readStdin()
		if err != nil {
			return "", err
		}
	}

	if stdinData != "" {
		createdPrompt += "\n" + stdinData
	}
	return createdPrompt, nil
}
