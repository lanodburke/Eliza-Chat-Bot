package main

// Donal Burke - G00337729

// list of imports
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// object that json data will be read into
type jsonobject struct {
	Keyword      []keyword      `json:"responses"`
	Substitution []substitution `json:"substitutions"`
}

// substitution struct to store regex pattern and a response to replace it with
type substitution struct {
	Word     string `json:"word"`
	Response string `json:"response"`
}

// keyword struct to hold regex pattern and list of responses
type keyword struct {
	Word      string   `json:"word"`
	Responses []string `json:"responses"`
}

// list of default responses
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

// if user enters one of these messages exits the program
var quitMessages = []string{
	"goodbye",
	"I have to leave",
	"Goodbye",
	"See ya later alligator!",
	"Our time is up, if you would like to continue that will be another 150 schmeckls.",
}

// function to parse json file and store its contents into the jsonboject struct
func parseKeywords() jsonobject {
	file, err := ioutil.ReadFile("./eliza.json") // file to read from
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}

	list := jsonobject{}
	_ = json.Unmarshal(file, &list)

	return list
}

// I adapted this function from the youtube video tutorials to work with a custom struct
func askEliza(input string) string {
	// clean up input
	input = strings.TrimRight(input, "\n.!")
	input = strings.ToLower(input)
	// instantiate new variable with parsed json file values from structs
	list := parseKeywords()

	// if they enter a quit statement then select random quit message from a list
	if isQuit(input) {
		os.Exit(1)
	}
	// Look for a possible response.
	for _, response := range list.Keyword {
		// compile regex pattern from string
		re := regexp.MustCompile(response.Word)
		// Check if the user input matches the pattern, capturing any groups.
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
						// compile regex pattern from string
						sub := regexp.MustCompile(substitution.Word)
						// if it matches one of the tokens replace with a substitution response
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
	// If there are no matches, then return random respone from a list of default responses.
	return randChoice(defualtResponses)
}

// check to see if the user has quit the program
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

// get a random number from 0 - 1
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
