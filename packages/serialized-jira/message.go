package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/andygrunwald/go-jira"
	"github.com/sony/gobreaker"
)

var cb *gobreaker.CircuitBreaker
var rxURL = regexp.MustCompile("(?:(?:https?|ftp):\\/\\/)?[\\w/\\-?=%.]+\\.[\\w/\\-?=%.]+")

type Message struct {
	Summary   string
	Description string
	Type string
	Errors  map[string]string
}

// TODO: Move validation of fields to frontend
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

	//issue, _, err := establishClient().Issue.Create(&i)
	//checkError(err)
	//
	//// TODO: Remove printing to console
	//fmt.Printf("%s: %+v\n", issue.Key, i.Fields.Summary)
	//
	//return err

	_, err := cb.Execute(func() (interface{}, error) {
		issue, _, err := establishClient().Issue.Create(&i)
		checkError(err)
		
		fmt.Printf("%s: %+v\n", issue.Key, i.Fields.Summary)

		return err, nil
	})

	return err
}
