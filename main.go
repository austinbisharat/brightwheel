package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/austinbisharat/brightwheel/emailservice"
)

// Represents the json that the /email endpoint recieves
type emailRequest struct {
	To       string `json:"to"`
	ToName   string `json:"to_name"`
	From     string `json:"from"`
	FromName string `json:"from_name"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

func main() {
	var emailServiceName string
	flag.StringVar(&emailServiceName, "email_service", "mailgun", "which email service to use")
	flag.Parse()

	var emailService emailservice.EmailService
	var err error
	if emailServiceName == "mailgun" {
		emailService, err = emailservice.NewMailgunEmailService()
	} else if emailServiceName == "sendgrid" {
		emailService, err = emailservice.NewSendgridEmailService()
	} else {
		log.Fatalf("Unknown email service option (%s)", emailServiceName)
	}

	if err != nil {
		log.Fatalf("Cannot create new email service (%s)", err)
	}

	emailRequestHandler := func(w http.ResponseWriter, r *http.Request) {
		// Disallow anything other than posts
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Attempt to parse the body into emailRequest struct
		decoder := json.NewDecoder(r.Body)
		var emailRequest emailRequest
		err := decoder.Decode(&emailRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid json body (%s)\n", err)))
			return
		}

		// Validate our input
		sendReq, err := validateEmailRequest(emailRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid email request (%s)\n", err)))
			return
		}

		// Use our EmailService to attempt to send the email
		err = emailService.SendEmail(sendReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf("Error from email provider: %s\n", err)))
			return
		}
	}

	http.HandleFunc("/email", emailRequestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
