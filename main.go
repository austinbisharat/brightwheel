package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/austinbisharat/brightwheel/emailservice"
)

// Represents that json that the /email endpoint recieves
type emailRequest struct {
	To       string `json:"to"`
	ToName   string `json:"to_name"`
	From     string `json:"from"`
	FromName string `json:"from_name"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
}

func main() {
	emailService, err := emailservice.NewSendgridEmailService()
	if err != nil {
		log.Fatalf("Cannot create new email service (%s)", err)
	}

	emailRequestHandler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		decoder := json.NewDecoder(r.Body)

		var emailRequest emailRequest
		err := decoder.Decode(&emailRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid json body (%s)\n", err)))
			return
		}

		sendReq, err := validateEmailRequest(emailRequest)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("Invalid email request (%s)\n", err)))
			return
		}

		err = emailService.SendEmail(sendReq)
		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(fmt.Sprintf("Error from email provider (%s)\n", err)))
			return
		}
	}

	http.HandleFunc("/email", emailRequestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
