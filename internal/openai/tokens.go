package openai

import (
	"fmt"
	"github.com/pkg/errors"
	tokenizer "github.com/samber/go-gpt-3-encoder"
	"github.com/sashabaranov/go-openai"
)

// validateChatCtx ensures that the chat context is not too long based on the max model tokens https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
func validateChatCtx(chatCtx []openai.ChatCompletionMessage, modelDetails ModelDetails) (int, error) {
	tokenCount, err := numTokensFromMessages(chatCtx, modelDetails)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to encode prompt data")
	}
	diff := modelDetails.MaxTokens - 4 - tokenCount
	if diff < 0 {
		return 0, fmt.Errorf("chat context is too long")
	}
	return diff, nil
}

func numTokensFromMessages(chatCtx []openai.ChatCompletionMessage, modelDetails ModelDetails) (int, error) {
	encoder, err := tokenizer.NewEncoder()
	if err != nil {
		return 0, err
	}

	numTokens := 0
	for _, msg := range chatCtx {
		numTokens += modelDetails.TokensPerMsg
		encoded, err := encoder.Encode(msg.Content)
		if err != nil {
			return 0, err
		}
		numTokens += len(encoded)
		if msg.Name != "" {
			numTokens += modelDetails.TokensPerName
		}
	}
	numTokens += 3
	return numTokens, nil
}
