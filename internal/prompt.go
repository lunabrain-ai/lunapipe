package internal

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
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
	if context.Args().First() == "" {
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
