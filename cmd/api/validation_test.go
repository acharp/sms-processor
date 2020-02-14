package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		smsRequest    SMSRequest
		expectedError error
		description   string
	}{
		{
			smsRequest: SMSRequest{
				Recipient: 123456, Originator: "AngryBirds", Message: "A bird has no name",
			},
			expectedError: nil,
			description:   "Valid input",
		},
		{
			smsRequest: SMSRequest{
				Recipient: 123456, Originator: "", Message: "A bird has no name",
			},
			expectedError: NewInvalidInputError("originator", "", ""),
			description:   "Empty Originator",
		},
		{
			smsRequest: SMSRequest{
				Recipient: 123456, Originator: "abcdefghijkl", Message: "A bird has no name",
			},
			expectedError: NewInvalidInputError("originator", "", ""),
			description:   "Invalid Originator",
		},
		{
			smsRequest: SMSRequest{
				Recipient: 123456, Originator: "AngryBirds", Message: "",
			},
			expectedError: NewInvalidInputError("message", "", ""),
			description:   "Empty Message",
		},
		{
			smsRequest: SMSRequest{
				Recipient: 123456, Originator: "AngryBirds",
				Message: "A waaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaay too long msg",
			},
			expectedError: NewInvalidInputError("message", "", ""),
			description:   "Invalid Message",
		},
		{
			smsRequest: SMSRequest{
				Recipient: 0, Originator: "AngryBirds", Message: "A bird has no name",
			},
			expectedError: NewInvalidInputError("recipient", "", ""),
			description:   "Empty Recipient",
		},
		{
			smsRequest: SMSRequest{
				Recipient: 130549632359034612, Originator: "AngryBirds", Message: "A bird has no name",
			},
			expectedError: NewInvalidInputError("recipient", "", ""),
			description:   "Invalid Recipient",
		},
	}

	for _, test := range testCases {
		err := Validate(test.smsRequest)
		// We only test error type since the error message is likely to change
		assert.IsType(t, test.expectedError, err, test.description)
	}
}
