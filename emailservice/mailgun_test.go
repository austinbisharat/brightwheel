package emailservice

import (
	"net/url"
	"reflect"
	"testing"
)

func TestConvertGenericReqToMailgunFormData(t *testing.T) {
	genericReq := EmailSendRequest{
		To:         "fake@example.com",
		ToName:     "Foo",
		From:       "fake2@example.com",
		FromName:   "Bar",
		RawSubject: "some subject",
		RawBody:    "some body",
	}

	expectedFormData := url.Values{
		"from":    {"Bar <fake2@example.com>"},
		"to":      {"Foo <fake@example.com>"},
		"subject": {"some subject"},
		"text":    {"some body"},
	}

	formData := genericRequestToMailgunRequestFormData(genericReq)

	if !reflect.DeepEqual(formData, expectedFormData) {
		t.Errorf("Expected %+v but got %+v", expectedFormData, formData)
	}
}
