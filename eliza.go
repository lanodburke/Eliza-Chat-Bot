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
)

type keywords struct {
	Keyword []keyword `json:"keywords"`
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
	"Very interesting.",
	"I see. And what does that tell you?",
	"How does that make you feel?",
	"How do you feel when you say that?",
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

func parseKeywords() keywords {
	file, err := ioutil.ReadFile("./eliza.json")
	if err != nil {
		fmt.Printf("File error: %v\n", err)
	}

	listKeywords := keywords{}
	_ = json.Unmarshal(file, &listKeywords)

	return listKeywords
}

func askEliza(input string) string {
	list := parseKeywords()

	for _, responses := range list.Keyword {
		re := regexp.MustCompile(responses.Word)
		matches := re.FindStringSubmatch(input)

		if len(matches) > 0 {
			response := randChoice(responses.Responses)

			return response
		}
	}

	return randChoice(defualtResponses)
}
func randChoice(list []string) string {
	randIndex := rand.Intn(len(list))
	return list[randIndex]
}

func chatWindow(w http.ResponseWriter, r *http.Request) {
	userSent := r.Header.Get("userAskEliza")

	fmt.Fprintf(w, askEliza(userSent))
}

func main() {
	http.HandleFunc("/askEliza", chatWindow)
	http.Handle("/", http.FileServer(http.Dir("./page")))
	http.ListenAndServe(":8080", nil)
}
