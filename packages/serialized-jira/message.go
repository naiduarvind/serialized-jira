package main

import (
	"github.com/andygrunwald/go-jira"
	"log"
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

	_, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		log.Println(err, "Unable to create issue in Jira.")
	}

	return err
}
