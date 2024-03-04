package main

import (
	"bufio"
	"os"

	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/consts"
	proxy "github.com/maxneuvians/go-copilot-proxy/pkg"
)

func main() {
	// Get the authentication token
	accessToken := getAuthToken()

	// Get the session token
	sessionToken := getSessionToken(accessToken)

	// Set the temperature
	proxy.Completion_temperature = 0.5

	// Set the prompt
	var chatMessages []proxy.Message

	chatMessages = append(chatMessages, proxy.Message{
		Role: "system",
		Content: `You are a system that generates simple quiz questions for french grammar in the context of working in the Canadian government. 
			For example you could be asked, 'Generate five questions that test the students understanding of the verb aller in the past tense'
			You provide your question in a json format. 
			If you are returning multiple questions, please use a JSON array. 
			Your response should be valid JSON that is minized. 
			Please validate that the answers form gramatically correct sentences when inserted into the sentences.
			`,
	})

	chatMessages = append(chatMessages, proxy.Message{
		Role: "user",
		Content: `
		Generate 20 questions that test the students understanding of indirect object pronouns, please vary the gender and number of the pronouns including me, te, lui, nos, vos, and leur.
		The format could follow something like this: {"question": "Je ____ donne le livre (à mon frère).", "answer": "lui", "translation": "I gave the book to my brother", "complete": "Je lui donne le livre"}.
		However, the types of sentences and pronouns should vary, please make sure to include quiestion for vos and nos.
		`,
	})

	var resp string
	err := proxy.Chat(sessionToken, chatMessages, false, func(cr proxy.CompletionResponse) error {
		resp = cr.Choices[0].Message.Content
		return nil
	})

	if err != nil {
		panic(err)
	}

	// Print the response
	writeToFile("./questions/indirect_object_pronouns.json", resp)
}

func getAuthToken() string {

	// Open .github_copilot_token file and read the first line

	//Check if .github_copilot_token file exists
	_, err := os.Stat(consts.TokenFile)

	if err != nil {
		panic(err)
	}

	// Get the authentication token
	file, err := os.Open(consts.TokenFile)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	r := bufio.NewReader(file)
	buffer, _, err := r.ReadLine()

	if err != nil {
		panic(err)
	}

	return string(buffer)
}

func getSessionToken(accessToken string) string {
	// Get the session token
	sessionResponse, err := proxy.GetSessionToken(accessToken)

	if err != nil {
		panic(err)
	}

	return sessionResponse.Token
}

func writeToFile(fileName string, data string) {
	// Write the data
	file, err := os.Create(fileName)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	_, err = file.WriteString(data)

	if err != nil {
		panic(err)
	}

	file.Sync()
}
