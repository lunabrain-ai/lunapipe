package prompt

import (
	"bufio"
	"io/fs"
	"os"
	"os/exec"
)

func NewFuncMap(params map[string]string) map[string]any {
	// TODO breadchris what other useful funcs should be here?
	return map[string]any{
		"param": func(key string) string {
			return params[key]
		},
		"shell": func(cmd string) ([]string, error) {
			return runCmdAndGetOutput(cmd)
		},
		"readDir": func(dir string) ([]string, error) {
			return fs.Glob(os.DirFS(dir), "*")
		},
	}
}

func runCmdAndGetOutput(command string) ([]string, error) {
	cmd := exec.Command("sh", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(stdout)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return lines, nil
}
