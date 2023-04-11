package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	tokenizer "github.com/samber/go-gpt-3-encoder"
	"github.com/urfave/cli/v2"
	"gorm.io/gorm"
)

func CallGPT(
	client gpt3.Client,
	prompt string,
	stream bool,
) (string, error) {
	var (
		maxModelTokens = 4000
	)

	encoder, err := tokenizer.NewEncoder()
	if err != nil {
		return "", err
	}

	encodedPromptData, err := encoder.Encode(prompt)
	if err != nil {
		return "", err
	}

	diff := maxModelTokens - len(encodedPromptData)
	if diff < 0 {
		return "", errors.New("prompt is too long")
	}

	var data string
	onData := func(resp *gpt3.ChatCompletionStreamResponse) {
		if len(resp.Choices) == 0 {
			return
		}
		text := resp.Choices[0].Delta.Content
		if stream {
			fmt.Print(text)
		}
		data += text
	}

	err = client.ChatCompletionStream(context.Background(), gpt3.ChatCompletionRequest{
		Temperature: float32(0),
		MaxTokens:   diff,
		Stream:      true,
		Messages: []gpt3.ChatCompletionRequestMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}, onData)
	if err != nil {
		return "", err
	}
	return data, nil
}

func NewCLI(
	config Config,
	db *gorm.DB,
) *cli.App {
	return &cli.App{
		Name:  "aicli",
		Usage: "AI for your CLI!",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "sync",
				Aliases: []string{"s"},
				Usage:   "do not stream output",
			},
		},
		Action: func(context *cli.Context) error {
			prompt := context.Args().First()

			client := gpt3.NewClient(config.APIKey)
			stream := !context.Bool("sync")

			data, err := CallGPT(client, prompt, stream)
			if err != nil {
				return err
			}
			if !stream {
				fmt.Println(data)
			}
			return nil
		},
	}
}
