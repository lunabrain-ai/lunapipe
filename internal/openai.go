package internal

import (
	"bufio"
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	tokenizer "github.com/samber/go-gpt-3-encoder"
	"os"
	"time"
)

// TODO breadchris is the max tokens?
const (
	maxModelTokens = 4096
)

type QAClient interface {
	Ask(prompt string, stream bool) (string, error)
	Chat(stream bool) error
}

type OpenAIQAClient struct {
	client gpt3.Client
}

var QAClientProviderSet = wire.NewSet(
	NewOpenAIQAClient,
	wire.Bind(new(QAClient), new(*OpenAIQAClient)),
)

var _ QAClient = &OpenAIQAClient{}

func (c *OpenAIQAClient) Chat(stream bool) error {
	var chatCtx []gpt3.ChatCompletionRequestMessage

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		chatCtx = append(chatCtx, gpt3.ChatCompletionRequestMessage{
			Role:    "user",
			Content: msg,
		})

		text, err := c.AskWithContext(chatCtx, stream)
		if err != nil {
			return err
		}
		println()

		chatCtx = append(chatCtx, gpt3.ChatCompletionRequestMessage{
			Role:    "system",
			Content: text,
		})
	}
	return nil
}

func (c *OpenAIQAClient) AskWithContext(chatCtx []gpt3.ChatCompletionRequestMessage, stream bool) (string, error) {
	respTokenCount, err := validateChatCtx(chatCtx)
	if err != nil {
		return "", err
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

	log.Debug().
		Int("respTokenCount", respTokenCount).
		Bool("stream", stream).
		Msg("Sending request to OpenAI")

	req := gpt3.ChatCompletionRequest{
		Temperature: float32(0),
		MaxTokens:   respTokenCount,
		Stream:      stream,
		Messages:    chatCtx,
	}

	// context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = c.client.ChatCompletionStream(ctx, req, onData)
	if err != nil {
		return "", err
	}
	return data, nil
}

func (c *OpenAIQAClient) Ask(prompt string, stream bool) (string, error) {
	chatCtx := []gpt3.ChatCompletionRequestMessage{
		{
			Role:    "user",
			Content: prompt,
		},
	}
	return c.AskWithContext(chatCtx, stream)
}

// validateChatCtx ensures that the chat context is not too long based on the max model tokens https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
func validateChatCtx(chatCtx []gpt3.ChatCompletionRequestMessage) (int, error) {
	tokenCount, err := numTokensFromMessages(chatCtx, "gpt-3.5-turbo")
	if err != nil {
		return 0, errors.Wrapf(err, "failed to encode prompt data")
	}
	diff := maxModelTokens - 4 - tokenCount
	if diff < 0 {
		return 0, fmt.Errorf("chat context is too long")
	}
	return diff, nil
}

func numTokensFromMessages(chatCtx []gpt3.ChatCompletionRequestMessage, model string) (int, error) {
	encoder, err := tokenizer.NewEncoder()
	if err != nil {
		return 0, err
	}

	var (
		tokensPerMessage int
		//tokensPerName    int
	)

	if model == "gpt-3.5-turbo" {
		tokensPerMessage = 4
		//tokensPerName = -1
	} else if model == "gpt-4" {
		tokensPerMessage = 3
		//tokensPerName = 1
	} else {
		return 0, fmt.Errorf("model not supported: %s", model)
	}

	numTokens := 0
	for _, msg := range chatCtx {
		numTokens += tokensPerMessage
		encoded, err := encoder.Encode(msg.Content)
		if err != nil {
			return 0, err
		}
		numTokens += len(encoded)
		// TODO breadchris there is no chat message name at the moment.
		//if key == "name" {
		//	numTokens += tokensPerName
		//}
	}
	numTokens += 3
	return numTokens, nil
}

func NewOpenAIQAClient(config Config) *OpenAIQAClient {
	client := gpt3.NewClient(config.APIKey)
	return &OpenAIQAClient{
		client: client,
	}
}
