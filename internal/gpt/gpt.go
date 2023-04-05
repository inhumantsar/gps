package gpt

import (
	"context"

	openai "github.com/sashabaranov/go-openai"
)

func ChatCompletionStream(cfg GptConfig, prompt string, additionalContext []openai.ChatCompletionMessage) (*openai.ChatCompletionStream, error) {
	c := openai.NewClient(cfg.ApiKey)
	ctx := context.Background()
	msgs := append(cfg.Messages, additionalContext...)
	msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: prompt, Name: ""})

	req := openai.ChatCompletionRequest{
		Model:     cfg.ModelName,
		MaxTokens: cfg.MaxTokens,
		Messages:  msgs,
		Stream:    true,
	}

	return c.CreateChatCompletionStream(ctx, req)
}

func ChatCompletion(cfg GptConfig, prompt string, additionalContext []openai.ChatCompletionMessage) (string, error) {
	c := openai.NewClient(cfg.ApiKey)
	ctx := context.Background()
	msgs := append(cfg.Messages, additionalContext...)
	msgs = append(msgs, openai.ChatCompletionMessage{Role: openai.ChatMessageRoleUser, Content: prompt, Name: ""})

	// for _, msg := range msgs {
	// 	fmt.Printf("msg: %s, role: %s, name: %s", msg.Content, msg.Role, msg.Name)
	// }

	req := openai.ChatCompletionRequest{
		Model:     cfg.ModelName,
		MaxTokens: cfg.MaxTokens,
		Messages:  msgs,
		Stream:    false,
	}

	resp, err := c.CreateChatCompletion(ctx, req)

	if err != nil || resp.Choices == nil || len(resp.Choices) == 0 {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
