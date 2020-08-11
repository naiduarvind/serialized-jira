package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
)

var rxURL = regexp.MustCompile("(?:(?:https?|ftp):\\/\\/)?[\\w/\\-?=%.]+\\.[\\w/\\-?=%.]+")

type Message struct {
	Summary   string
	Description string
	Type string
	Errors  map[string]string
}

func (msg *Message) Validate() bool {
	msg.Errors = make(map[string]string)

	if strings.TrimSpace(msg.Summary) == "" {
		msg.Errors["Summary"] = "Please enter a ticket summary"
	}

	match := rxURL.Match([]byte(msg.Description))
	if match == false {
		msg.Errors["Description"] = "Please enter a valid URL"
	}

	if strings.TrimSpace(msg.Type) == "" {
		msg.Errors["Type"] = "Please enter the correct issue type"
	}

	return len(msg.Errors) == 0
}

func (msg *Message) Deliver() error {
	base := "https://thebilityengineer.atlassian.net"

	tp := jira.BasicAuthTransport{
		Username: "techmaxed.net@gmail.com",
		Password: "a0jf3hW8TtJmSxc7JBQi7281",
	}
	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: msg.Description,
			Type: jira.IssueType{
				Name: msg.Type,
			},
			Project: jira.Project{
				Key: "TBE",
			},
			Summary: msg.Summary,
		},
	}

	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: %+v\n", issue.Key, i.Fields.Summary)

	return err
}

func Present()  {
	base := "https://thebilityengineer.atlassian.net"

	tp := jira.BasicAuthTransport{
		Username: "techmaxed.net@gmail.com",
		Password: "a0jf3hW8TtJmSxc7JBQi7281",
	}
	jiraClient, err := jira.NewClient(tp.Client(), base)
	if err != nil {
		panic(err)
	}

	jql := "project = TBE and type = Task and Status IN ('In Progress')"

	issues, _, err := jiraClient.Issue.Search(jql, nil)
	if err != nil {
		panic(err)
	}

	for _, i := range issues {
		fmt.Printf("(%s) - %+v : %s\n", i.Key, i.Fields.Summary, i.Fields.Description)
	}
}
