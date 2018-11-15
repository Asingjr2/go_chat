package main

import (
	"flag"
	"log"
	"net/http"
	"text/template"

	oth "chat_go/other"
)

// Need to investigate this code format.
func homeHandler(tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		tpl.Execute(w, req)
	})
}

func main() {
	// Check what flag package does.
	flag.Parse()
	tpl := template.Must(template.ParseFiles("index.html"))

	h := oth.NewHub()
	router := http.NewServeMux()				// multiplexer
	router.Handle("/", homeHandler(tpl))
	router.Handle("/ws", oth.WSHandler{h: h})		// need to check what this function does
	log.Printf("serving on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
