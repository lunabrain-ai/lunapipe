package prompt

import (
	"bufio"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
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
