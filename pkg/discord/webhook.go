package discord

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type DiscordWebhook struct {
	url string
}

func New(u string) DiscordWebhook {
	return DiscordWebhook{url: u}
}

func (dw DiscordWebhook) send(payload Message) error {
	p := new(bytes.Buffer)

	err := json.NewEncoder(p).Encode(payload)

	if err != nil {
		log.Warnf("An error occurred while encoding the webhook. Error %s", err)
		return err
	}

	resp, err := http.Post(dw.url, "application/json", p)

	if err != nil {
		log.Warnf("An error occurred while sending the webhook. Error %s", err)
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 204 {

		responseBody, err := io.ReadAll(resp.Body)

		if err != nil {
			log.Warnf("An error occurred while reading the response body. Error %s", err)
			return err
		}

		log.Warnf("An error occurred while sending the webhook, Error %s", string(responseBody))
		return err
	}

	return nil
}

func (dw DiscordWebhook) SendEmbed(e Embed) error {
	return dw.send(Message{Embeds: &[]Embed{e}})
}

func (dw DiscordWebhook) SendMessage(m string) error {
	return dw.send(Message{Content: &m})
}
