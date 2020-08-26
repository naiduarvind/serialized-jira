package main

import (
	"context"

	"github.com/andygrunwald/go-jira"
	"github.com/slok/goresilience/circuitbreaker"
)

type ticketInformation struct {
	Summary   string
	Description string
	Type string
	Errors  map[string]string
}

func (tickInfo *ticketInformation) createTicket() error {
	runner := circuitbreaker.New(circuitbreaker.Config{})

	err := runner.Run(context.Background(), func(ctx context.Context) error {
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

		_,_, err := establishClient().Issue.Create(&i)
		checkError(err)

		return nil
	})

	return err
}
