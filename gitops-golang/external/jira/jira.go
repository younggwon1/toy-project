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
		Assignee struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"assignee"`
		Participants []struct {
			EmailAddress string `json:"emailAddress"`
		} `json:"customfield_10243"` // customfield_10243 is a custom field for participants
	} `json:"fields"`
}

func GetTicketStatus(email, ticket string, j *JiraConnection) error {
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
		return fmt.Errorf("ticket %s is not found, check whether the entered ticket exists", ticket)
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

	/*
		배포 조건 정리
		1. jira ticket 상태는 "Reviewed" or "Deployed and Monitoring" 이어야 한다.
		2. executor는 assignee 이거나 participant 중 한 명이어야한다.
		3. 위 조건이 맞지 않으면 배포는 불가능하다.
	*/
	// check if ticket status is ready to deploy
	if ticketStatus.Fields.Status.Name == "Reviewed" || ticketStatus.Fields.Status.Name == "Deployed and Monitoring" {
		if ticketStatus.Fields.Assignee.EmailAddress == email {
			return nil
		} else if ticketStatus.Fields.Participants != nil {
			found := false
			for _, participant := range ticketStatus.Fields.Participants {
				if participant.EmailAddress == email {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("executor %s does not belong to the participant", email)
			}
		} else {
			return fmt.Errorf("%s is not executor, deployment is possible only if the executor is an assignee", email)
		}
	} else {
		return fmt.Errorf("ticket %s status is not ready to deploy, possible only if the ticket status is 'Reviewed' and 'Deployed and Monitoring'", ticket)
	}

	return nil
}

func TicketStatusCheck(j *JiraConnection, email, ticket string) error {
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
			err := GetTicketStatus(email, t, j)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("%s,", err))
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
