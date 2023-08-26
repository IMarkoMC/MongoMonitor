package pushover

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type PushoverClient struct {
	token string
	user  string
}

func New(t, u string) PushoverClient {
	return PushoverClient{token: t, user: u}
}

// * Handles sending the POST request to the Pushover API
func (pc PushoverClient) send(t, m string) error {
	params := url.Values{}

	params.Add("token", pc.token)
	params.Add("user", pc.user)

	params.Add("title", t)
	params.Add("message", m)

	req, err := http.NewRequest("POST", "https://api.pushover.net/1/messages.json", strings.NewReader(params.Encode()))

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected response code. Expected 200 got %d", resp.StatusCode)
	}

	return nil
}

func (pc PushoverClient) SendNotification(title, message string) error {
	return pc.send(title, message)
}
