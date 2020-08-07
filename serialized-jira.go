package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	jira "github.com/andygrunwald/go-jira"
	"golang.org/x/crypto/ssh/terminal"
)

func main() {
	r := bufio.NewReader(os.Stdin)

	fmt.Print("Jira URL: ")
	jiraURL, _ := r.ReadString('\n')

	fmt.Print("Jira Username: ")
	username, _ := r.ReadString('\n')
	fmt.Print("Jira username is ", username)

	fmt.Print("Jira Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)

	tp := jira.BasicAuthTransport{
		Username: strings.TrimSpace(username),
		Password: strings.TrimSpace(password),
	}

	client, err := jira.NewClient(tp.Client(), strings.TrimSpace(jiraURL))
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	u, _, err := client.User.Get("f5c8ce98-1fd2-4368-89f1-0d34217d61e4")

	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return
	}

	fmt.Printf("\nEmail: %v\nSuccess!\n", u.EmailAddress)

}