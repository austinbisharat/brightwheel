package emailservice

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type mailgunEmailService struct {
	apiKey        string
	mailgunDomain string
}

func NewMailgunEmailService() (EmailService, error) {
	apiKey := os.Getenv("MAILGUN_API_KEY")
	if len(apiKey) == 0 {
		return nil, errors.New("mailgun api key not set")
	}

	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	if len(mailgunDomain) == 0 {
		return nil, errors.New("mailgun domain not set")
	}
	return &mailgunEmailService{apiKey, mailgunDomain}, nil
}

func (meg *mailgunEmailService) SendEmail(sendReq EmailSendRequest) error {
	mailgunFormData := genericRequestToMailgunRequestFormData(sendReq)

	url := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages", meg.mailgunDomain)
	post, err := http.NewRequest("POST", url, strings.NewReader(mailgunFormData.Encode()))
	if err != nil {
		return fmt.Errorf("could not make http request object: %s", err)
	}

	post.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	post.SetBasicAuth("api", meg.apiKey)

	resp, err := http.DefaultClient.Do(post)
	if err != nil {
		return fmt.Errorf("got error from mailgun: %s", err)
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading resp body from mailgun: %s", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("error from mailgun api: %d, %s", resp.StatusCode, respBody)
	}

	log.Printf("Successfully sent email to %s via mailgun", sendReq.To)
	return nil
}

func genericRequestToMailgunRequestFormData(sendReq EmailSendRequest) url.Values {
	return url.Values{
		"from":    {fmt.Sprintf("%s <%s>", sendReq.FromName, sendReq.From)},
		"to":      {fmt.Sprintf("%s <%s>", sendReq.ToName, sendReq.To)},
		"subject": {sendReq.RawSubject},
		"text":    {sendReq.RawBody},
	}
}
