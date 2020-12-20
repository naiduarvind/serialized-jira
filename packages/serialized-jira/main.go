package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-xray-sdk-go/xraylog"

	"github.com/andygrunwald/go-jira"
	"github.com/apex/gateway"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/gorilla/mux"
	"github.com/secrethub/secrethub-go/pkg/secrethub"
)

var (
	err          error
	jiraUsername string
	jiraPassword string
	jiraClient   *jira.Client

	baseURL = "https://thebilityengineer.atlassian.net"
	jql     = "project = TBE and type = Task and Status IN ('In Progress')"
)

// TODO: Move into message.go
type TicketData struct {
	TicketSummary     string
	TicketDescription string
	TicketLabel       string
	TicketProgress    int
}

func init() {
	// INFO: AWS X-Ray Configuration & Logger Setup
	err = xray.Configure(xray.Config{
		ServiceVersion: sha1ver,
	})
	if err != nil {
		panic(err)
	}
	xray.SetLogger(xraylog.NewDefaultLogger(os.Stderr, xraylog.LogLevelError))
}

func init() {
	// INFO: Retrieval of Secrets using SecretHub Client
	client := secrethub.Must(secrethub.NewClient())
	jiraUsername, err = client.Secrets().ReadString("naiduarvind/serializedjira/username")
	if err != nil {
		panic(err)
	}
	jiraPassword, err = client.Secrets().ReadString("naiduarvind/serializedjira/password")
	if err != nil {
		panic(err)
	}
}

func init() {
	// INFO: Create JIRA client for reuse
	jiraClient = createJiraClient()
}

func createJiraClient() *jira.Client {
	tp := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraPassword,
	}
	client, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		log.Panic(err, "Unable to establish a connection to Jira service.")
	}

	return client
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/", send).Methods("POST")
	r.HandleFunc("/app/debug", handleDebug).Methods("GET")

	http.Handle("/", xray.Handler(xray.NewDynamicSegmentNamer("SerializedJIra", "jira.thebility.engineer"), r))
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	var td []TicketData

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		log.Panic(err, "Unable to perform JQL search in Jira.")
	}

	for _, issue := range issues {
		td = append(td, TicketData{
			issue.Fields.Summary,
			issue.Fields.Description,
			strings.Trim(fmt.Sprint(issue.Fields.Labels), "[]"),
			issue.Fields.Progress.Percent})
	}

	// TODO: Abstract writing headers separately
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	render(w, "templates/index.html", td)
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

	// TODO: Abstract writing headers separately
	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	// TODO: Abstract writing headers separately
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
