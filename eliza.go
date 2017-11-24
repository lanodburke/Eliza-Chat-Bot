package main

// Donal Burke - G00337729
/* sources
---------------
https://bootsnipp.com/snippets/WaEvr - chat template
---------------
*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
)

type jsonobject struct {
	Keyword      []keyword      `json:"responses"`
	Substitution []substitution `json:"substitutions"`
}

type substitution struct {
	Word     string `json:"word"`
	Response string `json:"response"`
}

type keyword struct {
	Word      string   `json:"word"`
	Responses []string `json:"responses"`
}

var defualtResponses = []string{
	"Please tell me more.",
	"Let's change focus a bit... Tell me about your family.",
	"Can you elaborate on that?",
	"I see.",
	"Ohh, man! Ohh, ohhh geez! Ohh... ",
	"Aww geez Rick, not again!",
	"Very interesting.",
	"I see. And what does that tell you?",
	"How does that make you feel?",
	"How do you feel when you say that?",
}

var initialMessages = []string{
	"Hello How are you feeling today?",
	"How do you do. Please tell me your problem.",
	"Please tell me what's been bothering you.",
}

var quitMessages = []string{
	"goodbye",
	"I have to leave",
	"Goodbye",
	"See ya later alligator!",
	"Our time is up, if you would like to continue that will be another 150 schmeckls.",
}

func parseKeywords() jsonobject {
	file, err := ioutil.ReadFile("./eliza.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}

	list := jsonobject{}
	_ = json.Unmarshal(file, &list)

	return list
}

var firstQuestion = false

func askEliza(input string) string {
	list := parseKeywords()

	if firstQuestion == false {
		randChoice(initialMessages)
		firstQuestion = true
	}

	if isQuit(input) {
		return randChoice(quitMessages)
	}
	// Look for a possible response.
	for _, response := range list.Keyword {
		re := regexp.MustCompile(response.Word)
		// Check if the user input matches the original, capturing any groups.
		if matches := re.FindStringSubmatch(input); matches != nil {
			// Select a random response.
			output := response.Responses[rand.Intn(len(response.Responses))]
			// We'll tokenise the captured groups using the following regular expression.
			boundaries := regexp.MustCompile(`[\s,.?!]+`)
			// Fill the response with each captured group from the input.
			// This is a bit complex, because we have to reflect the pronouns.
			for _, match := range matches[1:] {
				// First split the captured group into tokens.
				tokens := boundaries.Split(match, -1)
				// Loop through the tokens.
				for t, token := range tokens {
					// Loop through the potential substitutions.
					for _, substitution := range list.Substitution {
						// Check if the original of the current substitution matches the token.
						sub := regexp.MustCompile(substitution.Word)
						if sub.MatchString(token) {
							// If it matches, replace the token with one of the replacements (at random).
							// Then break.
							tokens[t] = substitution.Response
							break
						}
					}
				}
				output = strings.Replace(output, "%s", strings.Join(tokens, " "), -1)
			}
			// Send the filled answer back.
			return output
		}
	}
	// If there are no matches, then return this generic response.
	return randChoice(defualtResponses)
}

func isQuit(input string) bool {
	input = strings.TrimRight(input, "\n.!")
	input = strings.ToLower(input)
	for _, quit := range quitMessages {
		if input == quit {
			return true
		}
	}
	return false
}

func randChoice(list []string) string {
	randIndex := rand.Intn(len(list))
	return list[randIndex]
}

func chatWindow(w http.ResponseWriter, r *http.Request) {
	input := r.URL.Query().Get("value")
	fmt.Fprintf(w, askEliza(input))
}

func main() {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)
	http.HandleFunc("/user-input", chatWindow)
	http.ListenAndServe(":8080", nil)
}
