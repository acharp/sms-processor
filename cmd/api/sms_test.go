package main

import (
	"net/http/httptest"
	"testing"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	MBClient = messagebird.New("1pNVqoQvEhC9cmFaswLEW6YUU")

	testCases := []struct {
		body          []byte
		expStatusCode int
		description   string
	}{
		{
			body:          []byte("{\"recipient\":31633450007,\"originator\":\"MB\",\"message\":\"test\"}"),
			expStatusCode: 200,
			description:   "Valid body",
		},
		{
			body:          []byte("{,\"recipient\":31633450007,\"originator\":\"MB\",\"message\":\"test\"}"),
			expStatusCode: 400,
			description:   "Wrongly formatted json",
		},
		{
			body:          []byte("{\"recipient\":\"31633450007\",\"originator\":\"MB\",\"message\":\"test\"}"),
			expStatusCode: 400,
			description:   "Recipient has wrong type",
		},
		{
			body:          []byte("{\"originator\":\"MB\",\"message\":\"test\"}"),
			expStatusCode: 400,
			description:   "Missing recipient",
		},
		{
			body:          []byte("{\"recipient\":31633450007,\"originator\":\"\",\"message\":\"test\"}"),
			expStatusCode: 400,
			description:   "Empty originator",
		},
		{
			body:          []byte("{\"recipient\":31633450007,\"originator\":\"abcdefghijklmnopq\",\"message\":\"test\"}"),
			expStatusCode: 400,
			description:   "Originator too long",
		},
		{
			body:          []byte("{\"recipient\":31633450007,\"originator\":\"\",\"message\":\"\"}"),
			expStatusCode: 400,
			description:   "Empty message",
		},
		{
			body:          []byte("{\"recipient\":31633450007,\"originator\":\"\"}"),
			expStatusCode: 400,
			description:   "Missing message",
		},
	}

	for _, test := range testCases {
		rr := httptest.NewRecorder()
		HandleSMS(rr, test.body, MBClient)
		// We only test the status code since the response body is likely to change
		assert.Equal(t, test.expStatusCode, rr.Result().StatusCode, test.description)
	}
}
