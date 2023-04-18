package openai

import (
	"bufio"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/sashabaranov/go-openai"
	"io"
	"os"
	"strings"
)

type ModelDetails struct {
	MaxTokens     int
	TokensPerMsg  int
	TokensPerName int
}

var (
	models = map[string]ModelDetails{
		openai.GPT3Dot5Turbo: {
			MaxTokens:     4096,
			TokensPerMsg:  4,
			TokensPerName: -1,
		},
		openai.GPT4: {
			MaxTokens:     8000,
			TokensPerMsg:  3,
			TokensPerName: 1,
		},
	}
)

type QAClient interface {
	Ask(prompt string, stream bool) (string, error)
	Chat(stream bool) error
}

type OpenAIQAClient struct {
	client       *openai.Client
	config       Config
	model        string
	modelDetails ModelDetails
}

var _ QAClient = &OpenAIQAClient{}

func (c *OpenAIQAClient) Chat(stream bool) error {
	var chatCtx []openai.ChatCompletionMessage

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		chatCtx = append(chatCtx, openai.ChatCompletionMessage{
			Role:    "user",
			Content: msg,
		})

		text, err := c.AskWithContext(chatCtx, stream)
		if err != nil {
			log.Debug().Err(err).Msg("failed to get response from OpenAI")
			println("\n--- response did not complete, type \"continue\" for more ---")
			continue
		}
		println()

		chatCtx = append(chatCtx, openai.ChatCompletionMessage{
			Role:    "system",
			Content: text,
		})
	}
	return nil
}

func (c *OpenAIQAClient) AskWithContext(chatCtx []openai.ChatCompletionMessage, stream bool) (string, error) {
	respTokenCount, err := validateChatCtx(chatCtx, c.modelDetails)
	if err != nil {
		return "", err
	}

	var data string
	log.Debug().
		Int("respTokenCount", respTokenCount).
		Bool("stream", stream).
		Msg("Sending request to OpenAI")

	req := openai.ChatCompletionRequest{
		Model:       c.model,
		Temperature: float32(0),
		MaxTokens:   respTokenCount,
		Stream:      stream,
		Messages:    chatCtx,
	}

	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), c.config.Timeout)
	defer cancel()

	chatStream, err := c.client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return "", errors.Wrapf(err, "failed to send request to OpenAI")
	}

	for {
		response, err := chatStream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			log.Error().Err(err).Msg("failed to get response from OpenAI")
			break
		}

		fmt.Printf(response.Choices[0].Delta.Content)
	}
	return data, nil
}

func (c *OpenAIQAClient) Ask(prompt string, stream bool) (string, error) {
	chatCtx := []openai.ChatCompletionMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}
	return c.AskWithContext(chatCtx, stream)
}

func NewOpenAIQAClient(c Config) (*OpenAIQAClient, error) {
	if c.APIKey == "" {
		return nil, errors.New("OpenAI client is not configured correctly. Make sure OPENAI_API_KEY is set.")
	}

	if _, ok := models[c.Model]; !ok {
		var availableModels []string
		for model := range models {
			availableModels = append(availableModels, model)
		}
		return nil, fmt.Errorf("model not supported: %s, available models: %s", c.Model, strings.Join(availableModels, ", "))
	}

	client := openai.NewClient(c.APIKey)
	return &OpenAIQAClient{
		client:       client,
		config:       c,
		model:        c.Model,
		modelDetails: models[c.Model],
	}, nil
}
