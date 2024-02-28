package exercisepane

type prompt struct {
	name   string
	prompt string
}

var (
	prompts = []prompt{
		{
			name: "Personal Pronouns",
			prompt: `
				Generate five questions that test the students understanding of personal pronouns.
				The format should look like this: {"question": "__ suis allé a la plage.", "answer": "Je", "translation": "I went to the beach"} 
			`,
		},
		{
			name: "Verb Aller (Present Tense)",
			prompt: `
				Generate five questions that test the students understanding of the verb aller in the present tense. 
				The format should look like this: {"question": "je (aller) à la plage.", "answer": "suis allé", "translation": "I went to the beach"}
			`,
		},
		{
			name: "Verb Aller (Past Tense)",
			prompt: `
				Generate five questions that test the students understanding of the verb aller in the past tense". 
				The format should look like this: {"question": "je (aller) à la plage.", "answer": "suis allé", "translation": "I went to the beach"} 
			`,
		},
	}
)
