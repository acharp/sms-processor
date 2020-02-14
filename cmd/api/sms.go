package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
	"github.com/messagebird/go-rest-api/sms"
	"golang.org/x/time/rate"
)

// SMSRequest is the type matching the body of a POST to /sms
type SMSRequest struct {
	Recipient  int    `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
}

// HandleSMS process an sms request and send it through the MessageBird API
func HandleSMS(w http.ResponseWriter, body []byte, MBClient *messagebird.Client) {
	log.Println("Handling request on /sms endpoint")

	// Unmarshal body
	var smsRequest SMSRequest
	err := json.Unmarshal(body, &smsRequest)
	if err != nil {
		log.Println(err)
		err = NewInternalError(err.Error())
		JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate inputs
	log.Printf("Validating smsRequest input : '%v'", smsRequest)
	err = Validate(smsRequest)
	if err != nil {
		log.Printf("Input validation failed: %+v", err)
		JSONResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call MessageBird service
	log.Println("Sending SMS through messagebird service")
	err = sendSMS(MBClient, smsRequest)
	if err != nil {
		log.Println(err)
		JSONResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Finished to handle request with body '%s'", string(body))
	JSONResponse(w, "Your SMS is being handled", http.StatusOK)
}

// HandleSMSAsync process an sms request body in the background and send it through the MessageBird API
// This one deserves a bit of love to avoid duplicating code with HandleSMS but I'm running out of time
func HandleSMSAsync(bodyChannel chan []byte, errChannel chan<- error,
	MBClient *messagebird.Client, limiter *rate.Limiter) {

	for body := range bodyChannel {

		// Don't carry on if we are already exceeding the requests to the messagebird service
		if !limiter.Allow() {
			bodyChannel <- body
			time.Sleep(500 * time.Millisecond)
			continue
		}

		log.Printf("Processing body '%s' picked from bodyChannel", string(body))
		// Unmarshal body
		var smsRequest SMSRequest
		err := json.Unmarshal(body, &smsRequest)
		if err != nil {
			log.Println(err)
			err = NewInternalError(err.Error())
			errChannel <- err
		}

		// Validate inputs
		log.Printf("Validating smsRequest input : '%v'", smsRequest)
		err = Validate(smsRequest)
		if err != nil {
			log.Printf("Input validation failed: %+v", err)
			errChannel <- err
		}

		// Call MessageBird service
		log.Println("Sending SMS through messagebird service")
		err = sendSMS(MBClient, smsRequest)
		if err != nil {
			log.Println(err)
			errChannel <- err
		}

		log.Printf("Finished to process body '%v' picked from bodyChannel", string(body))
	}
}

// Call MessageBird service to send an SMS
func sendSMS(client *messagebird.Client, smsRequest SMSRequest) error {
	msg, err := sms.Create(
		client,
		smsRequest.Originator,
		[]string{strconv.Itoa(smsRequest.Recipient)},
		smsRequest.Message,
		nil,
	)
	if err != nil {
		err = NewInternalError(fmt.Sprintf("Failed sending sms through messagebird : '%s'", err.Error()))
		return err
	}

	log.Printf("SMS successfully sent to messagebird service '%v'", msg)
	return nil
}
