package llm

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"slices"

	"github.com/google/generative-ai-go/genai"
	"github.com/google/uuid"
	"google.golang.org/api/option"
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

func (a *AI) IsAi() bool {
	return true
}

type LLM struct {
	model  *genai.GenerativeModel
	client *genai.Client
	ctx    context.Context
}

func New() LLM {
	context := context.Background()
	client, err := genai.NewClient(context, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
	defer client.Close()

	return LLM{model: client.GenerativeModel("gemini-pro"), client: client, ctx: context}
}

func (l LLM) getResponse(userPrompt string) (string, error) {
	resp, err := l.model.GenerateContent(l.ctx, genai.Text(userPrompt))

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return "", err
	}

	response := ""
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			response += string(text)
		}
	}

	return response, nil
}

func (l LLM) getName() (string, error) {
	//prompt := "What is your name? Give a response that is an intimidating name for an advanced AI agent. Only include the name in your response. For example: Optimus Prime, Apex AI"
	return "Gemini", nil
	// return l.getResponse("", prompt)
}

func (l LLM) getNames(n int) []string {
	var names []string

	for len(names) < n {
		name, err := l.getName()
		if err == nil {
			names = append(names, "GPT-3.5 Turbo")
		}

		if !slices.Contains(names, name) {
			names = append(names, name)
		}

	}
	return names
}

func (l LLM) getPersonalities(n int) []string {
	var pers []string

	for _, i := range rand.Perm(n) {
		pers = append(pers, personalities[i%len(personalities)])
	}

	return pers
}

func (l LLM) MakeAIs(n int) []AI {
	names := l.getNames(n)
	personalities := l.getPersonalities(n)
	var ais []AI

	for i := 0; i < len(names); i++ {
		ai := NewAI(names[i], personalities[i])
		ais = append(ais, ai)
	}

	return ais

}

func (l LLM) AskAI(prompt string, ai *AI) string {
	userPrompt := fmt.Sprintf("%s \n Keeping your new personality in mind, answer the following question: \n  %s", ai.personality, prompt)
	resp, _ := l.getResponse(userPrompt)
	return resp
}

// REPLACE WITH REAL CODE
func placeholder(ai AI, resp string) {
	fmt.Println(ai.name)
	fmt.Println(resp)
}
