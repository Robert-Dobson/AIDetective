package llm

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/google/uuid"
	openai "github.com/sashabaranov/go-openai"
)

type AI struct {
	uuid        string
	name        string
	personality string
	eliminated  bool
}

func NewAI(name string, personality string) AI {
	return AI{
		uuid:        uuid.NewString(),
		name:        name,
		personality: personality,
		eliminated:  false,
	}
}

func (a *AI) UUID() string {
	return a.uuid
}

func (a *AI) Name() string {
	return a.name
}

func (a *AI) Eliminated() bool {
	return a.eliminated
}

func (a *AI) Eliminate() {
	a.eliminated = true
}

type LLM struct {
	client *openai.Client
}

func New() LLM {
	return LLM{openai.NewClient(os.Getenv("OPENAI_KEY"))}
}

func (l LLM) getResponse(systemPrompt string, userPrompt string) (string, error) {
	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: userPrompt,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "", err
	}

	return resp.Choices[0].Message.Content, nil

}

func (l LLM) getName() string {
	prompt := "What is your name? Give a response that is an intimidating name for an advanced AI agent. Only include the name in your response. For example: Optimus Prime, Apex AI"
	response, _ := l.getResponse("", prompt)
	return response
}

func (l LLM) GetNames(n int) []string {
	var names []string

	for len(names) < n {
		name := l.getName()

		if !slices.Contains(names, name) {
			names = append(names, name)
		}

	}

	return names
}

func (l LLM) AskAI(prompt string, ai AI, callback func(ai AI, resp string)) {
	userPrompt := fmt.Sprintf("%s \n Keeping your new personality in mind, answer the following question: \n  %s", ai.personality, prompt)
	resp, _ := l.getResponse("", userPrompt)
	callback(ai, resp)
}

// REPLACE WITH REAL CODE
func placeholder(ai AI, resp string) {
	fmt.Println(ai.name)
	fmt.Println(resp)
}
