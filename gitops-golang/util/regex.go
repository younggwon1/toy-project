package util

import "regexp"

func ValidateTicket(ticket string) bool {
	// regex for jira ticket
	regex := regexp.MustCompile(`^[A-Z]+-\d+(_[A-Z]+-\d+)*$`)
	result := regex.MatchString(ticket)

	return result
}
