package exercisepane

type prompt struct {
	name   string
	prompt string
}

var (
	prompts = []prompt{
		{
			name:   "Verb Aller (Present Tense)",
			prompt: "Generate five questions that test the students understanding of the verb aller in the present tense",
		},
		{
			name:   "Verb Aller (Past Tense)",
			prompt: "Generate five questions that test the students understanding of the verb aller in the past tense",
		},
	}
)
