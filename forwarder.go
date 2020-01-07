package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/mailgun/mailgun-go/v3"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func sendEmailWithMailgun(recipient string, body string) {
	domain := os.Getenv("MAILGUN_DOMAINS")
	secretKey := os.Getenv("MAILGUN_PRIVATE_KEY")

	sender := os.Getenv("EMAIL_SENDER")
	subject := "Ahoj, máš zprávu od Ježíška"

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, secretKey)

	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(sender, subject, body, recipient)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := mg.Send(ctx, message)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

func forwardEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Email forwarded")

	fmt.Printf("Post from website! subject = %v\n", r.FormValue("Subject"))
	fmt.Printf("Post from website! body-html = %v\n", r.FormValue("body-html"))
	fmt.Printf("Post from website! sender = %v\n", r.FormValue("sender"))
	fmt.Printf("Post from website! To = %v\n", r.FormValue("To"))

	to := r.FormValue("Subject")
	body := r.FormValue("body-html")

	if to != "" && body != "" {
		sendEmailWithMailgun(to, body)
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
