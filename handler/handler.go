package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	messagebird "github.com/messagebird/go-rest-api/v6"
	"github.com/messagebird/go-rest-api/v6/sms"
)

// SMSRequest is the type matching the body of a POST to /sms
type SMSRequest struct {
	Recipient  int    `json:"recipient"`
	Originator string `json:"originator"`
	Message    string `json:"message"`
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
		return fmt.Errorf("failed sending sms through messagebird : '%s'", err.Error())
	}

	log.Printf("SMS successfully sent to messagebird service '%v'", msg)
	return nil
}

// Singleton instance of a Message Bird client.
var MBClient *messagebird.Client

func init() {
	// Test API key
	// MBClient = messagebird.New(os.Getenv("MBIRD_TEST"))
	// Live API key
	MBClient = messagebird.New(os.Getenv("MBIRD_LIVE"))
}

func HandleLambdaEvent() error {
	log.Println("Go HB!")

	birthdates := map[string]time.Time{
		"Anais":   time.Date(1987, time.February, 5, 0, 0, 0, 0, time.UTC),
		"Karine":  time.Date(1984, time.April, 1, 0, 0, 0, 0, time.UTC),
		"Mama":    time.Date(1960, time.May, 1, 0, 0, 0, 0, time.UTC),
		"Mathieu": time.Date(1987, time.June, 13, 0, 0, 0, 0, time.UTC),
		"Papa":    time.Date(1955, time.July, 12, 0, 0, 0, 0, time.UTC),
		"Gabitou": time.Date(1985, time.July, 28, 0, 0, 0, 0, time.UTC),
		"Rebecca": time.Date(1991, time.December, 10, 0, 0, 0, 0, time.UTC),

		"La douce": time.Date(1991, time.February, 7, 0, 0, 0, 0, time.UTC),
		// "Maz":      time.Date(1991, time.February, 10, 0, 0, 0, 0, time.UTC),
		"La dinde": time.Date(1992, time.February, 17, 0, 0, 0, 0, time.UTC),
		"Elod":     time.Date(1990, time.February, 18, 0, 0, 0, 0, time.UTC),
		"Thanos":   time.Date(1989, time.March, 15, 0, 0, 0, 0, time.UTC),
		"La garce": time.Date(1992, time.March, 16, 0, 0, 0, 0, time.UTC),
		// "Vidoule":  time.Date(1991, time.March, 17, 0, 0, 0, 0, time.UTC),
		// "Fritz":    time.Date(1990, time.April, 2, 0, 0, 0, 0, time.UTC),
		// "Guyor":    time.Date(1991, time.April, 4, 0, 0, 0, 0, time.UTC),
		"Rabbi": time.Date(1991, time.April, 23, 0, 0, 0, 0, time.UTC),
		"Pam":   time.Date(1991, time.May, 1, 0, 0, 0, 0, time.UTC),
		"Joubs": time.Date(1991, time.June, 24, 0, 0, 0, 0, time.UTC),
		// "Mike":     time.Date(1993, time.June, 26, 0, 0, 0, 0, time.UTC),
		"Pique":   time.Date(1991, time.July, 11, 0, 0, 0, 0, time.UTC),
		"La bete": time.Date(1991, time.July, 30, 0, 0, 0, 0, time.UTC),
		"Monfe":   time.Date(1991, time.August, 30, 0, 0, 0, 0, time.UTC),
		// "Chive":    time.Date(1991, time.November, 6, 0, 0, 0, 0, time.UTC),
		"Guede":    time.Date(1991, time.November, 20, 0, 0, 0, 0, time.UTC),
		"Jop":      time.Date(1983, time.December, 2, 0, 0, 0, 0, time.UTC),
		"Feral":    time.Date(1991, time.December, 4, 0, 0, 0, 0, time.UTC),
		"La lope":  time.Date(1991, time.December, 13, 0, 0, 0, 0, time.UTC),
		"La delbe": time.Date(1991, time.December, 26, 0, 0, 0, 0, time.UTC),
	}
	now := time.Now()
	smsSent := false
	phoneNumber, err := strconv.Atoi(os.Getenv("PHONE_NUMBER"))
	if err != nil {
		return err
	}

	for name, birthdate := range birthdates {
		if now.Day() == birthdate.Day() && now.Month() == birthdate.Month() {
			age := now.Year() - birthdate.Year()
			log.Printf("It is %s's birthday! %d years of greatness :) \n", name, age)
			smsRequest := SMSRequest{
				Recipient:  phoneNumber,
				Originator: "go-HB app",
				Message:    fmt.Sprintf("It is %s's birthday! %d years of greatness, wish it :) ", name, age),
			}
			err := sendSMS(MBClient, smsRequest)
			if err != nil {
				return err
			}
			smsSent = true
		} else {
			// log.Printf("It is not %s's birthday\n", name)
		}
	}

	if now.Day() == 1 && now.Month()%2 == 0 && !smsSent {
		smsRequest := SMSRequest{
			Recipient:  phoneNumber,
			Originator: "go-HB app",
			Message:    "go-hb is still standing strong!",
		}
		err := sendSMS(MBClient, smsRequest)
		if err != nil {
			return err
		}
	}

	log.Println("All the birthdays from the crew have been checked, bye!")
	return nil
}

func main() {
	lambda.Start(HandleLambdaEvent)
}
