package main

import (
	"errors"
	"regexp"

	"github.com/microcosm-cc/bluemonday"

	"github.com/austinbisharat/brightwheel/emailservice"
)

func validateEmailRequest(req emailRequest) (sendReq emailservice.EmailSendRequest, err error) {

	if len(req.ToName) > 256 {
		return sendReq, errors.New(`"to_name" field is > 256 characters`)
	}

	if len(req.FromName) > 256 {
		return sendReq, errors.New(`"from_name" field is > 256 characters`)
	}

	if len(req.Subject) > 512 {
		return sendReq, errors.New(`"subject" field is > 512 characters`)
	}

	if !isValidEmailFormat(req.To) {
		return sendReq, errors.New(`"to" field appears to be an invalid email address`)

	}

	if !isValidEmailFormat(req.From) {
		return sendReq, errors.New(`"from" field appears to be an invalid email address`)
	}

	rawSubject := bluemonday.StrictPolicy().Sanitize(req.Subject)
	rawBody := bluemonday.StrictPolicy().Sanitize(req.Body)

	sendReq = emailservice.EmailSendRequest{
		To:         req.To,
		ToName:     req.ToName,
		From:       req.From,
		FromName:   req.FromName,
		RawSubject: rawSubject,
		RawBody:    rawBody,
	}
	return sendReq, nil
}

var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidEmailFormat(email string) bool {
	return len(email) <= 254 && rxEmail.MatchString(email)
}
