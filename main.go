package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

func main() {
	mux := pat.New()
	mux.Get("/", http.HandlerFunc(home))
	mux.Post("/", http.HandlerFunc(send))
	mux.Get("/confirmation", http.HandlerFunc(confirmation))

	log.Println("Listening...")
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/index.html", nil)
}

func send(w http.ResponseWriter, r *http.Request) {
	// Step 1: Validate form
	msg := &Message{
		Summary:   r.PostFormValue("summary"),
		Description: r.PostFormValue("description"),
		Type: r.PostFormValue("type"),
		//Label: r.PostFormValue("label"),
	}

	if msg.Validate() == false {
		render(w, "templates/index.html", msg)
		return
	}

	// Step 2: Create ticket in Jira
	if err := msg.Deliver(); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	// Step 3: Redirect to confirmation page
	http.Redirect(w, r, "/confirmation", http.StatusSeeOther)
}

func confirmation(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/confirmation.html", nil)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	tmpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}

	if err := tmpl.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
	}
}