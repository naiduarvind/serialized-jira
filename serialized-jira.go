package main

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func main() {
	base := "https://thebilityengineer.atlassian.net"

	// TODO: Convert this to an ENV / Secrets Manager (depending on the infrastructure picked)
	tp := jira.BasicAuthTransport{}
	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	jql := "project = TBE and type = Task and Status IN ('In Progress')"
	fmt.Printf("Usecase: Running a JQL query '%s'\n", jql)
	issues, resp, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		panic(err)
	}

	outputResponse(issues, resp)
}

// TODO: Remove in favour of web forms through handlers
func outputResponse(issues []jira.Issue, resp *jira.Response) {
	// fmt.Printf("Call to %s\n", resp.Request.URL)
	// fmt.Printf("Response Code: %d\n", resp.StatusCode)
	fmt.Println("==================================")
	for _, i := range issues {
		fmt.Printf("%s (%s/%s): %+v\n", i.Key, i.Fields.Type.Name, i.Fields.Priority.Name, i.Fields.Summary)
	}
}
