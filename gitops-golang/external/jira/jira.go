package jira

import (
	"encoding/json"
	"fmt"
	"io"
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

type TicketStatus struct {
	Fields struct {
		Status struct {
			Name string `json:"name"`
		} `json:"status"`
	} `json:"fields"`
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
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("ticket %s is not found", ticket)
	} else if res.StatusCode != http.StatusOK {
		return err
	}
	if res != nil {
		defer res.Body.Close()
	}

	// read response body and unmarshal ticket status
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	var ticketStatus TicketStatus
	err = json.Unmarshal(body, &ticketStatus)
	if err != nil {
		return err
	}

	// check if ticket is ready to deploy
	if ticketStatus.Fields.Status.Name == "Backlog" || ticketStatus.Fields.Status.Name == "Request Review" {
		return fmt.Errorf("ticket %s is not ready to deploy", ticket)
	}

	return nil
}

func TicketStatusCheck(ticket string) error {
	// separate jira ticket based '_'
	tickets := strings.Split(ticket, "_")

	// init goroutine
	wg := sync.WaitGroup{}
	// mutex for error slice
	mu := sync.Mutex{}

	// slice for collecting errors
	var errors []error

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
				mu.Lock()
				errors = append(errors, fmt.Errorf("%s ,", err))
				mu.Unlock()
			}
		}(t)
	}

	// wait for all goroutines to finish
	wg.Wait()

	// if there are any errors, return them
	if len(errors) > 0 {
		return fmt.Errorf("occurred errors : %v : please check the tickets with issues", errors)
	}

	return nil
}
