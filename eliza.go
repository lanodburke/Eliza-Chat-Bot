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
	"net/http"
)

var defualtResponses = []string{
	"Please tell me more.",
	"Let's change focus a bit... Tell me about your family.",
	"Can you elaborate on that?",
	"I see.",
	"Very interesting.",
	"I see. And what does that tell you?",
	"How does that make you fell?",
	"How do you fell when you say that?",
}

var reflections = map[string]string{
	"am":     "are",
	"was":    "were",
	"i":      "you",
	"i'd":    "you would",
	"i've":   "you have",
	"i'll":   "you will",
	"my":     "your",
	"are":    "am",
	"you've": "I have",
	"you'll": "I will",
	"your":   "my",
	"yours":  "mine",
	"you":    "me",
	"me":     "you",
}

type keywords struct {
	Keyword []keyword `json:"keywords"`
}

type keyword struct {
	Word      string   `json:"word"`
	Weight    int      `json:"weight"`
	Responses []string `json:"responses"`
}

func readFile() {
	file, err := ioutil.ReadFile("./eliza.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}

	var jsontype keywords
	json.Unmarshal(file, &jsontype)
	fmt.Printf("Results: %v\n", jsontype)
}

// func askEliza(input string) string {

// }

// ReplyTo will construct a reply for a given statement using ELIZA's rules.
// func ReplyTo(statement string) string {
// 	// First, preprocess the statement for more effective matching
// 	statement = preprocess(statement)

// 	// Next, we try to match the statement to a statement that ELIZA can
// 	// recognize, and construct a pre-determined, appropriate response.
// 	for pattern, responses := range Psychobabble {
// 		re := regexp.MustCompile(pattern)
// 		matches := re.FindStringSubmatch(statement)

// 		// If the statement matched any recognizable statements.
// 		if len(matches) > 0 {
// 			// If we matched a regex group in parentheses, get the first match.
// 			// The matched regex group will match a "fragment" that will form
// 			// part of the response, for added realism.
// 			var fragment string
// 			if len(matches) > 1 {
// 				fragment = reflect(matches[1])
// 			}

// 			// Choose a random appropriate response, and format it with the
// 			// fragment, if needed.
// 			response := randChoice(responses)
// 			if strings.Contains(response, "%s") {
// 				response = fmt.Sprintf(response, fragment)
// 			}
// 			return response
// 		}
// 	}

// 	// If no patterns were matched, return a default response.
// 	return randChoice(DefaultResponses)
// }

// // preprocess will do some normalization on a statement for better regex matching
// func preprocess(statement string) string {
// 	statement = strings.TrimRight(statement, "\n.!")
// 	statement = strings.ToLower(statement)
// 	return statement
// }

// // reflect flips a few words in an input fragment (such as "I" -> "you").
// func reflect(fragment string) string {
// 	words := strings.Split(fragment, " ")
// 	for i, word := range words {
// 		if reflectedWord, ok := Reflections[word]; ok {
// 			words[i] = reflectedWord
// 		}
// 	}
// 	return strings.Join(words, " ")
// }

// // randChoice returns a random element in an (string) array.
// func randChoice(list []string) string {
// 	randIndex := rand.Intn(len(list))
// 	return list[randIndex]
// }

func chatWindow(w http.ResponseWriter, r *http.Request) {
	//get the string from the input box
	// userSent := r.Header.Get("userAskEliza")

	// fmt.Fprintf(w, ReplyTo(userSent))
}

func main() {
	readFile()
	http.HandleFunc("/askEliza", chatWindow)
	http.Handle("/", http.FileServer(http.Dir("./page")))
	http.ListenAndServe(":8080", nil)
}
