package main

import (
	"fmt"

	jira "github.com/andygrunwald/go-jira"
)

func main() {
	base := "https://thebilityengineer.atlassian.net"

	// TODO: Convert this to an ENV / Secrets Manager (depending on the infrastructure picked)
	tp := jira.BasicAuthTransport{}
	client, err := jira.NewClient(tp.Client(), base)
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

	issue, _, err := client.Issue.Create(&i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v", issue.Key, i.Fields.Summary)
}
