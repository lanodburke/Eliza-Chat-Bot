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

var initialMessages = []string{
	"Hello How are you feeling today?",
	"How do you do. Please tell me your problem.",
	"Please tell me what's been bothering you.",
}

var quitMessages = []string{
	"goodbye",
	"I have to leave",
	"Goodbye",
	"See ya latar alligator!",
	"Our time is up, if you would like to continue that will be another $150.",
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

	input = strings.TrimRight(input, "\n.!")
	input = strings.ToLower(input)

	for _, responses := range list.Keyword {
		re := regexp.MustCompile(responses.Word)
		matches := re.FindStringSubmatch(input)

		if len(matches) > 0 {
			var frag string
			if len(matches) > 1 {
				frag = reflection(matches[1])
			}

			response := randChoice(responses.Responses)
			if strings.Contains(response, "%s") {
				response = fmt.Sprintf(response, frag)
			}

			return response
		}
	}

	return randChoice(defualtResponses)
}
func randChoice(list []string) string {
	randIndex := rand.Intn(len(list))
	return list[randIndex]
}

func reflection(input string) string {
	s := strings.Split(input, " ")
	for i, str := range s {
		if reflected, ok := reflections[str]; ok {
			s[i] = reflected
		}
	}
	return strings.Join(s, " ")
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
