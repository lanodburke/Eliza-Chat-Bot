package main

// Donal Burke - G00337729
/* sources
---------------
https://bootsnipp.com/snippets/WaEvr - chat template
---------------
*/

import (
	"fmt"
	"net/http"
)

func chatWindow(w http.ResponseWriter, r *http.Request) {
	//get the string from the input box
	//userSent := r.Header.Get("userAskEliza")

	fmt.Fprintf(w, "Hi im Eliza")
}

func main() {
	http.HandleFunc("/askEliza", chatWindow)
	http.Handle("/", http.FileServer(http.Dir("./page")))
	http.ListenAndServe(":8080", nil)
}
