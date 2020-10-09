package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/apex/gateway"
	"github.com/gorilla/mux"
)

type TicketData struct {
	TicketSummary     string
	TicketDescription string
	TicketLabel       string
	TicketProgress    int
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/", send).Methods("POST")
	r.HandleFunc("/confirmation", confirmation).Methods("GET")

	http.Handle("/", r)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func home(w http.ResponseWriter, r *http.Request) {

	var td []TicketData

	jql := "project = TBE and type = Task and Status IN ('In Progress') AND createdDate <= startOfWeek()"

	issues, _, err := establishClient().Issue.Search(jql, nil)
	checkError(err)

	for _, issue := range issues {
		td = append(td, TicketData{
			issue.Fields.Summary,
			issue.Fields.Description,
			strings.Trim(fmt.Sprint(issue.Fields.Labels), "[]"),
			issue.Fields.Progress.Percent})
		checkError(err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	render(w, "templates/index.html", td)
}

func confirmation(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	render(w, "templates/confirmation.html", nil)
}

func send(w http.ResponseWriter, r *http.Request) {
	tickInfo := &ticketInformation{
		Summary:     r.PostFormValue("summary"),
		Description: r.PostFormValue("description"),
		Type:        r.PostFormValue("type"),
	}

	if err := tickInfo.createTicket(); err != nil {
		log.Println(err)
		http.Error(w, "Sorry, something went wrong", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Redirect(w, r, "/prod/confirmation", http.StatusMovedPermanently)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
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

func establishClient() *jira.Client {
	base := "https://thebilityengineer.atlassian.net"

	// TODO: temporary creds -- to be replaced with env variables
	tp := jira.BasicAuthTransport{
		Username: "techmaxed.net@gmail.com",
		Password: "a0jf3hW8TtJmSxc7JBQi7281",
	}
	jiraClient, err := jira.NewClient(tp.Client(), base)
	checkError(err)

	return jiraClient
}

//
func checkError(err error) {
	if err != nil {
		log.Panic(err)
	}
}
