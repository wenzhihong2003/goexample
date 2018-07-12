package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"regexp"
)

// super simple email regex
const emailRegex = `(\w+)\@(\w+)\.[a-zA-Z]`

func main() {
	http.HandleFunc("/", handler)
	log.Printf("listening on localhost:8080")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Path[1:]

	if isEmail(email) {
		log.Printf("%s looks like an email\n", email)
		fmt.Fprintf(w, "%s looks like an email", email)
		return
	}
	log.Printf("%s doesn't look like an email\n", email)
	fmt.Fprintf(w, "%s doesn't look an email", email)
}

func isEmail(email string) bool {
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
