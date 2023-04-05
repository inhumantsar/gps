package gpt

import (
	openai "github.com/sashabaranov/go-openai"
)

type GptConfig struct {
	Messages  []openai.ChatCompletionMessage `json:"messages,omitempty"`
	ApiKey    string                         `json:"api_key,omitempty"`
	MaxTokens int                            `json:"max_tokens,omitempty"`
	ModelName string                         `json:"model_name,omitempty"`
	// Temperature      float32                 `json:"temperature,omitempty"`
	// TopP             float32                 `json:"top_p,omitempty"`
	// N                int                     `json:"n,omitempty"`
	// PresencePenalty  float32                 `json:"presence_penalty,omitempty"`
	// FrequencyPenalty float32                 `json:"frequency_penalty,omitempty"`
	// LogitBias        map[string]int          `json:"logit_bias,omitempty"`
}
