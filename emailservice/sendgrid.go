package emailservice

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type sendgridEmailService struct {
	apiKey string
}

type sendgridEmailRequest struct {
	Personalizations []sendgridPersonalization `json:"personalizations"`
	From             sendgridEmail             `json:"from"`
	Subject          string                    `json:"subject"`
	Content          []sendgridContent         `json:"content"`
}

type sendgridPersonalization struct {
	To []sendgridEmail `json:"to"`
	// we could probably flush this out with other fields,
	// but no need for this application
}

type sendgridEmail struct {
	Email string `json:"email"`
}

type sendgridContent struct {
	ContentType string `json:"type"`
	Value       string `json:"value"`
}

func NewSendgridEmailService() (EmailService, error) {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("sendgrid api key not set")
	}
	return &sendgridEmailService{apiKey}, nil
}

func (seg *sendgridEmailService) SendEmail(sendReq EmailSendRequest) error {
	sendgridReq := genericRequestToSendgridRequest(sendReq)
	bodyBytes, err := json.Marshal(sendgridReq)
	if err != nil {
		return fmt.Errorf("could not marshall sendgrid request json: %s", err)
	}

	post, err := http.NewRequest("POST", "https://api.sendgrid.com/v3/mail/send", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("could not make http request object: %s", err)
	}

	post.Header.Set("Content-Type", "application/json")
	post.Header.Set("Authorization", fmt.Sprintf("Bearer %s", seg.apiKey))

	resp, err := http.DefaultClient.Do(post)
	if err != nil {
		return fmt.Errorf("got error from sendgrid: %s", err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resp body from sendgrid: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error from sendgrid api: %d, %s", resp.StatusCode, respBody)
	}

	log.Printf("Successfully send email to %s via sendgrid %s", sendReq.To, respBody)
	return nil
}

func genericRequestToSendgridRequest(sendReq EmailSendRequest) sendgridEmailRequest {
	return sendgridEmailRequest{
		Personalizations: []sendgridPersonalization{
			{
				To: []sendgridEmail{{fmt.Sprintf("%s <%s>", sendReq.ToName, sendReq.To)}},
			},
		},
		From:    sendgridEmail{fmt.Sprintf("%s <%s>", sendReq.FromName, sendReq.From)},
		Subject: sendReq.RawSubject,
		Content: []sendgridContent{
			{
				ContentType: "text/plain",
				Value:       sendReq.RawBody,
			},
		},
	}
}
