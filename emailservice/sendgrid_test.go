package emailservice

import (
	"reflect"
	"testing"
)

func TestConvertGenericReqToSendGridReq(t *testing.T) {
	genericReq := EmailSendRequest{
		To:         "fake@example.com",
		ToName:     "Foo",
		From:       "fake2@example.com",
		FromName:   "Bar",
		RawSubject: "some subject",
		RawBody:    "some body",
	}

	expectedSendgridRequest := sendgridEmailRequest{
		Personalizations: []sendgridPersonalization{
			{
				To: []sendgridEmail{{"Foo <fake@example.com>"}},
			},
		},
		From:    sendgridEmail{"Bar <fake2@example.com>"},
		Subject: "some subject",
		Content: []sendgridContent{
			{
				ContentType: "text/plain",
				Value:       "some body",
			},
		},
	}

	sgr := genericRequestToSendgridRequest(genericReq)

	if !reflect.DeepEqual(sgr, expectedSendgridRequest) {
		t.Errorf("Expected %+v but got %+v", expectedSendgridRequest, sgr)
	}
}
