# MB: sms-processor

Toy project to play around with go channels, API rate limitation and the MB sms delivery sdk.

## Install and run

Install dependencies: `make ensure`
Run the tests: `make test`
Run the project: `make run`

Then send requests to your `localhost:5000` socket. For example:

A valid request:
```
curl -X POST \
  http://localhost:5000/sms \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: BigBird' \
  -d '{"recipient":31633450007,"originator":"MB","message":"This is a test message."}'
  ```
Give the following response:
```
{
    "message": "Your SMS is being handled",
    "status": "ok"
}
```
  
An invalid request (missing originator):
  ```
curl -X POST \
  http://localhost:5000/sms \
  -H 'Content-Type: application/json' \
  -H 'X-API-Key: BigBird' \
  -d '{"recipient":31633450007,"message":"This is a test message."}'
  ```
Give the following response:
```
{
    "message": "'originator' field '' is invalid: 'Must be a phone number or an alpha numeric shorter or equal to 11 characters'",
    "status": "error"
}
```


## Design choices

There are 2 endpoints:
- A public healthcheck : `/health`
- A private (just behind an API key) post sms feature : `/sms`

The call to the downstream MB service to 1 RPS.
When this rate is exceeded, instead of just throttling the requests made to our API, I store them in a channel which I process in the background, still paying attention not to break the requirement of 1 RPS max made to the MB service.
This channel has a size of a 100 messages. In the case of exceeding this size then we start dropping requests and logging the problem. In the real world, instead of just dropping the requests I would output them to a stream (either Kinesis or Kafka) where the requests would be reprocessed until success to handle the backpressure accumulated due to the downstream service.
Last point about this feature, here the errors produced by the background job are just stored in an error channel and ignored. In a production application, the way to go would be to send the error through the reporting system of choice (Sentry, Datadog, ...) and back them up somewhere like S3 in order to investigate and maybe reprocess them later.

