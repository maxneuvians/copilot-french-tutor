package exercisepane

type prompt struct {
	name   string
	prompt string
}

var (
	prompts = []prompt{
		{
			name: "Adjectifs démonstratifs",
			prompt: `
				Generate ten questions that test the students understanding of the demonstrative adjectives. 
				The format should look like this: {"question": "Wahou ! Tu as vu ___ robe ?", "answer": "cette", "translation": "Wow! Have you seen that dress?"}.
				`,
		},
		{
			name: "Imparfait",
			prompt: `
				Generate ten questions that test the students understanding of the imparfait. The format should look like this: {"question": "Je (manger) un sandwich.", "answer": "mangeais", "translation": "I was eating a sandwich"}.
				Please make sure to include a variety of verbs and to include the correct form of the verb in the answer.
			`,
		},
		{
			name: "Indirect Object Pronouns",
			prompt: `
				Generate ten questions that test the students understanding of indirect object pronouns, please vary the gender and number of the pronouns including me, te, lui, nos, vos, and leur.
				The format could follow something like this: {"question": "Je donne le livre à mon frère. Je ___ donne le livre.", "answer": "lui", "translation": "I gave the book to my brother"}.
				However, the types of sentences and pronouns should vary, please make sure to include quiestion for vos and nos.
			`,
		},
		{
			name: "Passé Composé",
			prompt: `
				Generate ten questions that test the students understanding of the passé composé. The format should look like this: {"question": "Je (manger) un sandwich.", "answer": "ai mangé", "translation": "I ate a sandwich"}.
				Please make sure to include a variety of verbs and to include the correct form of the verb in the answer.
			`,
		},
		{
			name: "Personal Pronouns",
			prompt: `
				Generate five questions that test the students understanding of personal pronouns.
				The format should look like this: {"question": "__ suis allé a la plage.", "answer": "Je", "translation": "I went to the beach"} 
			`,
		},
		{
			name: "Verb Aller (Future simple)",
			prompt: `
				Generate five questions that test the students understanding of the verb aller in the future simple tense.
				The format should look like this: {"question": "je (aller) à la plage.", "answer": "irai", "translation": "I'll go to the beach"}
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
		{
			name: "Verb Etre (Future Simple)",
			prompt: `
				Generate five questions that test the students understanding of the verb être in the future simple tense.
				The format should look like this: {"question": "je (être) à la plage.", "answer": "serai", "translation": "I'll be at the beach"}
			`,
		},
	}
)
