package main

import (
	"errors"
	"regexp"

	"github.com/microcosm-cc/bluemonday"

	"github.com/austinbisharat/brightwheel/emailservice"
)

func validateEmailRequest(req emailRequest) (sendReq emailservice.EmailSendRequest, err error) {

	if len(req.ToName) < 0 {
		return sendReq, errors.New(`"to_name" field is empty`)
	} else if len(req.ToName) > 256 {
		return sendReq, errors.New(`"to_name" field is > 256 characters`)
	}

	if len(req.FromName) < 0 {
		return sendReq, errors.New(`"from_name" field is empty`)
	} else if len(req.FromName) > 256 {
		return sendReq, errors.New(`"from_name" field is > 256 characters`)
	}

	if len(req.Subject) < 0 {
		return sendReq, errors.New(`"subject" field is empty`)
	} else if len(req.Subject) > 512 {
		return sendReq, errors.New(`"subject" field is > 512 characters`)
	}

	if !isValidEmailFormat(req.To) {
		return sendReq, errors.New(`"to" field appears to be an invalid email address`)

	}

	if !isValidEmailFormat(req.From) {
		return sendReq, errors.New(`"from" field appears to be an invalid email address`)
	}

	// This sanitizes the request body of all html
	rawBody := bluemonday.StrictPolicy().Sanitize(req.Body)

	sendReq = emailservice.EmailSendRequest{
		To:         req.To,
		ToName:     req.ToName,
		From:       req.From,
		FromName:   req.FromName,
		RawSubject: req.Subject,
		RawBody:    rawBody,
	}
	return sendReq, nil
}

// This monster regex is stolen straight from google, but does a pretty decent sanity check that an email looks valid
var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func isValidEmailFormat(email string) bool {
	// emails also cannot be longer than 254 bytes long
	return len(email) <= 254 && rxEmail.MatchString(email)
}
