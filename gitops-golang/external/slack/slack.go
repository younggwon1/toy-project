package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

// func SendSlackMessage(webhookUrl string, data interface{}) error {
// 	body, err := yaml.Marshal(data)
// 	if err != nil {
// 		return err
// 	}
// 	res, err := http.Post(
// 		webhookUrl,
// 		"application/json",
// 		bytes.NewBuffer(body),
// 	)
// 	fmt.Println(res)
// 	if err != nil {
// 		return err
// 	}
// 	if res.StatusCode != http.StatusOK {
// 		return fmt.Errorf("failed to send slack message: %s", res.Status)
// 	}
// 	if res != nil {
// 		defer res.Body.Close()
// 	}

//		return nil
//	}
type SlackRequestBody struct {
	Text string `json:"text"`
}

func SendSlackMessage(webhookUrl string, tmpl string, data interface{}) error {
	t, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		return fmt.Errorf("failed to parse template: %w", err)
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}

	msg := tpl.String()
	body, _ := json.Marshal(SlackRequestBody{Text: msg})
	res, err := http.Post(
		webhookUrl,
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if buf.String() != "ok" {
		return fmt.Errorf("non-ok response returned from Slack: %s", buf.String())
	}

	return nil
}
