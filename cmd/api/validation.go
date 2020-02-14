package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// Validate an smsRequest object
func Validate(smsRequest SMSRequest) (err error) {
	err = validateOriginator(smsRequest.Originator)
	if err != nil {
		return
	}
	err = validateMessage(smsRequest.Message)
	if err != nil {
		return
	}
	err = validateRecipient(smsRequest.Recipient)
	if err != nil {
		return
	}

	return
}

// validate the originator field of an SMSRequest object
func validateOriginator(originator string) error {
	// To make it simpler we assume an originator phone number is simple int with 6 to 15 digits
	isPhoneNumber, err := regexp.MatchString("^\\d{6,15}$", originator)
	if err != nil {
		return NewInternalError(fmt.Sprintf("Failed checking originator regexp: '%s'", err.Error()))
	}
	if isPhoneNumber {
		return nil
	}

	isAlphaNum, err := regexp.MatchString("^[[:alnum:]]+$", originator)
	if err != nil {
		return NewInternalError(fmt.Sprintf("Failed checking originator regexp: '%s'", err.Error()))
	}
	if isAlphaNum && len(originator) <= 11 {
		return nil
	}

	return NewInvalidInputError("originator", originator,
		"Must be a phone number or an alpha numeric shorter or equal to 11 characters")

}

// validate the message field of an SMSRequest object
func validateMessage(message string) error {
	if len(message) <= 160 && message != "" {
		return nil
	}
	return NewInvalidInputError("message", message, "Must not be empty and shorter or equal to 160")

}

// validate the recipient field of an SMSRequest object
func validateRecipient(recipient int) error {
	// Not sure this is the right regexp of an MSISDN but this is not the point of the project
	isMSISDN, err := regexp.MatchString("^+[1-9]{1}[0-9]{3,14}$", strconv.Itoa(recipient))
	if err != nil {
		return NewInternalError(fmt.Sprintf("Failed checking recipient regexp: '%s'", err.Error()))
	}
	if isMSISDN {
		return nil
	}
	return NewInvalidInputError("recipient", strconv.Itoa(recipient), "Must be a valid MSISDN")
}
