package config

import (
	"testing"

	openai "github.com/sashabaranov/go-openai"
)

func TestParseConfig(t *testing.T) {

	tests := []struct {
		name         string
		filename     string
		expectedMsgs []openai.ChatCompletionMessage
		expectedErrs []string
	}{
		{
			name:     "valid config",
			filename: "testdata/config/valid_messages.yaml",
			expectedMsgs: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "Hi there! How can I assist you?",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Can you help me find a nearby restaurant?",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "Sure thing! What type of cuisine are you in the mood for?",
				},
			},
			expectedErrs: nil,
		},
		{
			name:         "invalid config",
			filename:     "testdata/config/invalid_messages.yaml",
			expectedMsgs: nil,
			expectedErrs: []string{
				"unknown role 'invalid', skipping message 'This message has an invalid role.'",
				"unknown role 'another_invalid', skipping message 'This message also has an invalid role.'",
			},
		},
		{
			name:     "mixed config",
			filename: "testdata/config/mixed_messages.yaml",
			expectedMsgs: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a helpful assistant.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "Hi there! How can I assist you?",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Can you help me find a nearby restaurant?",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "Sure thing! What type of cuisine are you in the mood for?",
				},
			},
			expectedErrs: []string{
				"unknown role 'invalid', skipping message 'This message has an invalid role.'",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// get the full test file path
			// filename := filepath.Join(projectpath.Root, test.filename)

			// msgs, errs := LoadConfig(filename)

			// if len(errs) != len(test.expectedErrs) {
			// 	t.Errorf("Expected %d errors, but got %d", len(test.expectedErrs), len(errs))
			// 	for _, err := range errs {
			// 		t.Errorf("  %s", err.Error())
			// 	}
			// } else {
			// 	for i, expectedErr := range test.expectedErrs {
			// 		if errs[i].Error() != expectedErr {
			// 			t.Errorf("Expected error '%s', but got '%s'", expectedErr, errs[i].Error())
			// 		}
			// 	}
			// }

			// if len(msgs) != len(test.expectedMsgs) {
			// 	t.Errorf("Expected %d messages, but got %d", len(test.expectedMsgs), len(msgs))
			// } else {
			// 	for i, expectedMsg := range test.expectedMsgs {
			// 		if msgs[i].Role != expectedMsg.Role {
			// 			t.Errorf("Expected message %d role to be '%s', but got '%s'", i, expectedMsg.Role, msgs[i].Role)
			// 		}

			// 		if msgs[i].Content != expectedMsg.Content {
			// 			t.Errorf("Expected message %d content to be '%s', but got '%s'", i, expectedMsg.Content, msgs[i].Content)
			// 		}
			// 	}
			// }
		})
	}
}
