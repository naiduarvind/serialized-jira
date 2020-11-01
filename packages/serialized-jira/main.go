package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/apex/gateway"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/aws/aws-xray-sdk-go/xraylog"
	"github.com/gorilla/mux"
	"github.com/secrethub/secrethub-go/pkg/secrethub"
)

var (
	jiraUsername string
	jiraPassword string
	sha1ver      string
	buildTime    string

	baseURL = "https://thebilityengineer.atlassian.net"
	jql 	= "project = TBE and type = Task and Status IN ('In Progress')"
)

type TicketData struct {
	TicketSummary     string
	TicketDescription string
	TicketLabel       string
	TicketProgress    int
}

func init() {
	var err error

	client := secrethub.Must(secrethub.NewClient())
	jiraUsername, err = client.Secrets().ReadString("naiduarvind/serializedjira/username")
	if err != nil {
		panic(err)
	}
	jiraPassword, err = client.Secrets().ReadString("naiduarvind/serializedjira/password")
	if err != nil {
		panic(err)
	}

	err = xray.Configure(xray.Config{
		ServiceVersion: sha1ver,
	})
	if err != nil {
		panic(err)
	}
	xray.SetLogger(xraylog.NewDefaultLogger(os.Stderr, xraylog.LogLevelError))
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

	issues, _, err := establishClient().Issue.Search(jql, nil)
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

	w.Header().Set("Content-Type", "text/html")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
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

func handleDebug(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("url: %s %s", r.Method, r.RequestURI)
	a := []string{s}

	a = append(a, "Headers:")
	for k, v := range r.Header {
		if len(v) == 0 {
			a = append(a, k)
		} else if len(v) == 1 {
			s = fmt.Sprintf("  %s: %v", k, v[0])
			a = append(a, s)
		} else {
			a = append(a, "  "+k+":")
			for _, v2 := range v {
				a = append(a, "    "+v2)
			}
		}
	}

	a = append(a, "")
	a = append(a, fmt.Sprintf("ver: https://github.com/naiduarvind/serialized-jira/commit/%s", sha1ver))
	a = append(a, fmt.Sprintf("built on: %s", buildTime))

	s = strings.Join(a, "\n")
	servePlainText(w, s)
}

func servePlainText(w http.ResponseWriter, s string) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(s)))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

// TODO: Possibility of moving this block into func init()
func establishClient() *jira.Client {

	tp := jira.BasicAuthTransport{
		Username: jiraUsername,
		Password: jiraPassword,
	}
	jiraClient, err := jira.NewClient(tp.Client(), baseURL)
	if err != nil {
		log.Panic(err, "Unable to establish a connection to Jira service.")
	}

	return jiraClient
}
