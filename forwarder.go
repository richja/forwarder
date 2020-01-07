package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func sendEmailWithMailgun(recipient string, body string) (string, string, error) {
	domain := os.Getenv("MAILGUN_DOMAINS")
	secretKey := os.Getenv("MAILGUN_PRIVATE_KEY")

	sender := os.Getenv("EMAIL_SENDER")
	subject := "Hi, someone sent you a message"

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, secretKey)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	return mg.Send(ctx, message)
	//return resp, id, err
}

func isEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	fmt.Printf("error: '%s' mentioned in subject is not a valid email address\n", email)

	return len(email) > 254 || !rxEmail.MatchString(email)
}

func forwardEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Post from website! subject = %v\n", r.FormValue("Subject"))
	fmt.Printf("Post from website! body-html = %v\n", r.FormValue("body-html"))
	fmt.Printf("Post from website! sender = %v\n", r.FormValue("sender"))
	fmt.Printf("Post from website! To = %v\n", r.FormValue("To"))

	to := r.FormValue("Subject")
	body := r.FormValue("body-html")

	if to != "" && isEmail(to) && body != "" {
		resp, id, err := sendEmailWithMailgun(to, body)

		if err != nil {
			fmt.Fprintf(w, "Sorry, something wrong happened, please try again later")
			log.Fatal(err)
		}

		fmt.Fprintf(w, "Email forwarded!")
		fmt.Printf("ID: %s Resp: %s\n", id, resp)
	} else {
		fmt.Fprintf(w, "Missing or invalid recepient email")
	}
}

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/forward", forwardEmail)

	fmt.Printf("Starting server on port %s \n", os.Getenv("PORT"))

	port := ":" + os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(port, nil))
}
