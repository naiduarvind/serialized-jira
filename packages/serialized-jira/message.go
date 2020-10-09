package main

import (
	"github.com/andygrunwald/go-jira"
)

type ticketInformation struct {
	Summary     string
	Description string
	Type        string
	Errors      map[string]string
}

func (tickInfo *ticketInformation) createTicket() error {

	i := jira.Issue{
		Fields: &jira.IssueFields{
			Description: tickInfo.Description,
			Type: jira.IssueType{
				Name: tickInfo.Type,
			},
			Project: jira.Project{
				Key: "TBE",
			},
			Summary: tickInfo.Summary,
		},
	}

	_, _, err := establishClient().Issue.Create(&i)
	checkError(err)

	return err
}
