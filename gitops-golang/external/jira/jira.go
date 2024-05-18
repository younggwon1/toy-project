package jira

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type JiraConnection struct {
	JiraUrl  string
	Username string
	Token    string
}

func GetTicketStatus(ticket string, j *JiraConnection) error {
	url, err := url.JoinPath(j.JiraUrl, "/rest/api/3/issue", ticket)
	if err != nil {
		return err
	}

	// prepare request
	req, err := http.NewRequest(
		"GET",
		url,
		nil,
	)

	// set headers
	req.SetBasicAuth(j.Username, j.Token)
	req.Header.Set("Accept", "application/json")
	if err != nil {
		return err
	}

	// send request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to get ticket status: %s", res.Status)
	}
	if res != nil {
		defer res.Body.Close()
	}

	return nil
}

func TicketStatusCheck(ticket string) error {
	// separate jira ticket based '_'
	tickets := strings.Split(ticket, "_")

	// init goroutine
	wg := sync.WaitGroup{}

	// goroutine for checking ticket status
	for _, t := range tickets {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			err := GetTicketStatus(t, &JiraConnection{
				JiraUrl:  "https://jira.com",
				Username: "",
				Token:    "",
			})
			if err != nil {
				fmt.Printf("failed to get ticket status: %s", err)
			}
		}(t)
	}

	return nil
}
