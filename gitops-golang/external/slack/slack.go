package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendMessage(webhookUrl, tmpl string, data interface{}) error {
	// parse template
	t, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	// execute template
	var tpl bytes.Buffer
	err = t.Execute(&tpl, data)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	// send slack message
	msg := tpl.String()
	body, _ := json.Marshal(SlackRequestBody{Text: msg})
	res, err := http.Post(
		webhookUrl,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send slack message: %s", res.Status)
	}
	if res != nil {
		defer res.Body.Close()
	}

	return nil
}
