package llm

func examples() {

	llm := New()
	// get list of 10 AI names
	names := llm.GetNames(10)

	// ask ai a question
	ai := NewAI(names[0], "Pretend you are a very rude person.")
	prompt := "How many days are there in a year?"
	llm.AskAI(prompt, ai)
}
