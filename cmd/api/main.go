package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	messagebird "github.com/messagebird/go-rest-api"
	"golang.org/x/time/rate"
)

// This should be moved to a proper log pkg but it is out of scope of this project
const (
	Ldate      = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                      // the time in the local time zone: 01:23:23
	Lshortfile                 // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                       // if Ldate or Ltime is set, use UTC rather than the local time zone
	LstdFlags  = Ldate | Ltime // initial values for the standard logger
)

// Singleton instance of a Message Bird client.
// Simplest implementation thus not concurrent safe: out of the scope for this project
var MBClient *messagebird.Client

func init() {
	// Testing key:
	MBClient = messagebird.New("")
	// Production key :
	// MBClient = messagebird.New("")
}

// Handle stopping the service gracefully
func handleInterrupt() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
	log.Println("^C received, closing app in 2 seconds")
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	go handleInterrupt()

	var limiter = rate.NewLimiter(1, 1)
	// When we reach the rate limitation we queue up to 100 request bodies in a channel
	// which we process in the background.
	// The errChanel gathers the errors happening during this background processing.
	// Here we don't do anything with the content of this error channel but in the real
	// world we obviously should.
	bodyChannel := make(chan []byte, 100)
	errChannel := make(chan error)
	go HandleSMSAsync(bodyChannel, errChannel, MBClient, limiter)
	defer close(bodyChannel)
	defer close(errChannel)

	// Send sms through messagebird
	http.HandleFunc("/sms",
		authAPIKey(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// Triage request types
			if r.Method != "POST" {
				msg := "Only POST method allowed on /sms endpoint"
				log.Println(msg)
				JSONResponse(w, msg, http.StatusBadRequest)
				return
			}

			// Read and unmarshal body to avoid working with an http.Request type asap
			body, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			if err != nil {
				log.Println(err)
				err = NewInternalError(err.Error())
				JSONResponse(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Handle request:
			// We are not exceeding the messagebird limit rate
			if limiter.Allow() {
				HandleSMS(w, body, MBClient)
			} else {
				if err != nil {
					log.Println(err)
					err = NewInternalError(err.Error())
					JSONResponse(w, err.Error(), http.StatusBadRequest)
				}
				select {
				// We are exceeding the messagebird limit rate but have room to queue the request
				case bodyChannel <- body:
					log.Printf("Request '%v' queued and will be processed later", r)
					JSONResponse(w, "Request queued, will be processed later", http.StatusOK)
				// We are exceeding the messagebird limit rate and don't room to queue the request
				default:
					log.Printf("Downstream service rate exceeded and queue full, "+
						"dropping the request '%v'", r)
					JSONResponse(w,
						"Downstream service rate exceeded and queue full, request droppped",
						http.StatusTooManyRequests,
					)
				}
			}
		}),
		),
	)

	// Healthcheck
	http.HandleFunc("/health",
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		},
	)

	log.Println("Server started listening for requests on port 5000")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
