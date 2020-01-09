# Email forwarder
*Note this is just an early MVP and WIP.*

Forward email sent on set email address to email address specificed in the subject using Mailgun.

## Prerequisites
1) Have domain verified in Mailgun
2) Make sure you are also set for receiving emails
3) [Configure new rules in Routes](https://app.mailgun.com/app/receiving/routes) on given domain to forward incoming emails to URL of your deployed app with */forward* endpoint (eg. example.com/forward)
4) Generate private key in Mailgun

## Setup
1) Have [Go](https://golang.org/) installed on your machine
2) Rename *.env.example* file to *.env*
3) Fill in details in *.env* file

## To run
1) Run `go run forwarder` to start HTTP server
2) Go to http://localhost:8080

## To run tests
1) Run `go test -v` to build and run tests