package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/maxneuvians/copilot-french-tutor/pkg/ui/consts"
	proxy "github.com/maxneuvians/go-copilot-proxy/pkg"
)

type Question struct {
	Question    string
	Answer      string
	Complete    string
	Translation string
	Valid       bool `json:"omitempty"`
}

func main() {
	// Get the authentication token
	accessToken := getAuthToken()

	// Get the session token
	sessionToken := getSessionToken(accessToken)

	// Read entire file into buffer
	buffer, err := os.ReadFile("./questions/indirect_object_pronouns.json")

	if err != nil {
		panic(err)
	}

	questions := []Question{}

	// Parse the questions
	err = json.Unmarshal(buffer, &questions)

	if err != nil {
		panic(err)
	}

	// Print the questions
	for i, q := range questions {
		println("Question " + string(i))
		valid := validateSentence(sessionToken, q.Complete)
		q.Valid = valid
	}

	// Write the questions
	questionsJSON, err := json.Marshal(questions)

	if err != nil {
		panic(err)
	}

	writeToFile("./questions/indirect_object_pronouns.json", string(questionsJSON))
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

func validateSentence(sessionToken string, sentence string) bool {
	// Set the prompt
	var chatMessages []proxy.Message

	chatMessages = append(chatMessages, proxy.Message{
		Role:    "system",
		Content: `You are a system that validates if french grammar questions containe the correct answer. The questions will be provided in the following json schema: { "question": "string", "answer": "string", "complete": "string", "translation": "string", "valid": "bool" }. The "question" field contains the question, the "answer" field contains the answer, the "complete" field contains the complete sentence, the "translation" field contains the translation of the sentence, and the "valid" field contains the boolean value of the validation. If the complete sentence is valid, return "true", otherwise return the corrected version of the sentence and update the answer field with the corrected version of the sentence.`,
	})

	chatMessages = append(chatMessages, proxy.Message{
		Role:    "user",
		Content: `Please validate that the following sentence is gramatically correct in french: ` + sentence + `. If it is correct return "true", otherwise return the corrected version of the sentence".`,
	})

	var resp string
	err := proxy.Chat(sessionToken, chatMessages, false, func(cr proxy.CompletionResponse) error {
		resp = cr.Choices[0].Message.Content
		return nil
	})

	if err != nil {
		panic(err)
	}

	if resp == "true" {
		return true
	} else {

		fmt.Println("The sentence is not valid: " + sentence)
		fmt.Println("Explanation: " + resp)
		return false
	}
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
