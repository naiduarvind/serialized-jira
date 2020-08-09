package main

import (
	"fmt"

	"github.com/andygrunwald/go-jira"
)

func main() {
	base := "https://thebilityengineer.atlassian.net"

	// TODO: Convert this to an ENV / Secrets Manager (depending on the infrastructure picked)
	tp := jira.BasicAuthTransport{

	}
	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: "Test Issue",
			Type: jira.IssueType{
				Name: "Task",
			},
			Project: jira.Project{
				Key: "TBE",
			},
			Summary: "Just a demo issue",
		},
	}

	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, i.Fields.Summary)

	jql := "project = TBE and type = Task and Status IN ('In Progress')"

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
		fmt.Printf("%s: %+v\n", i.Key, i.Fields.Summary, i.Fields.Description)
	}
}
