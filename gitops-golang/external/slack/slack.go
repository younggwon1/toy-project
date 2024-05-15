package slack

import (
	"bytes"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v3"
)

func SendSlackMessage(webhookUrl string, data interface{}) error {
	body, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	res, err := http.Post(
		webhookUrl,
		"application/json",
		bytes.NewBuffer(body),
	)
	fmt.Println(res)
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
