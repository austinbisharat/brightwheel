package main

import (
	"strings"
	"testing"

	"github.com/austinbisharat/brightwheel/emailservice"
)

type validationTestCase struct {
	apiReq                 emailRequest
	expectedGenericRequest emailservice.EmailSendRequest
	doesExpectError        bool
}

func TestValidateEmailRequest(t *testing.T) {
	var str strings.Builder

	for i := 0; i < 513; i++ {
		str.WriteString("a")
	}

	veryLongString := str.String()

	testCases := []validationTestCase{
		{
			// Empty input causes an error
			apiReq:                 emailRequest{},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// invalid to email
			apiReq: emailRequest{
				To:       "not an email",
				ToName:   "Foo",
				From:     "fake@example.com",
				FromName: "Fake",
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// invalid to email pt 2
			apiReq: emailRequest{
				To:       "not@an@email",
				ToName:   "Foo",
				From:     "fake@example.com",
				FromName: "Fake",
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// invalid from email
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "Foo",
				From:     "fake@  example.com",
				FromName: "Fake",
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// long name
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   veryLongString,
				From:     "fake@  example.com",
				FromName: "Fake",
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// long from name
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "fake",
				From:     "fake@  example.com",
				FromName: veryLongString,
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// long subject
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "Fake",
				From:     "fake@  example.com",
				FromName: "Fake",
				Subject:  veryLongString,
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// empty name
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "",
				From:     "fake@  example.com",
				FromName: "Fake",
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// empty from name
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "fake",
				From:     "fake@  example.com",
				FromName: "",
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// empty subject
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "Fake",
				From:     "fake@  example.com",
				FromName: "Fake",
				Subject:  "",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{},
			doesExpectError:        true,
		},
		{
			// basic correct case
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "Foo",
				From:     "fake@example.com",
				FromName: "Fake",
				Subject:  "some subject",
				Body:     "some body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{
				To:         "fake@example.com",
				ToName:     "Foo",
				From:       "fake@example.com",
				FromName:   "Fake",
				RawSubject: "some subject",
				RawBody:    "some body",
			},
			doesExpectError: false,
		},
		{
			// html tags case
			apiReq: emailRequest{
				To:       "fake@example.com",
				ToName:   "Foo",
				From:     "fake@example.com",
				FromName: "Fake",
				Subject:  "some subject",
				Body:     "<h1>some</h1> body",
			},
			expectedGenericRequest: emailservice.EmailSendRequest{
				To:         "fake@example.com",
				ToName:     "Foo",
				From:       "fake@example.com",
				FromName:   "Fake",
				RawSubject: "some subject",
				RawBody:    "some body",
			},
			doesExpectError: false,
		},
	}

	for i, testCase := range testCases {
		genericReq, err := validateEmailRequest(testCase.apiReq)

		if testCase.doesExpectError && err == nil {
			t.Errorf("Test case %d expected a validation error but got none", i)
		} else if !testCase.doesExpectError && err != nil {
			t.Errorf("Test case %d expected no validation error but (%s)", i, err)
		} else if testCase.expectedGenericRequest != genericReq {
			t.Errorf("Test case %d expected generic request %+v but got %+v", i, testCase.expectedGenericRequest, genericReq)
		}
	}
}
