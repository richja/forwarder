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

func forwardEmail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Email forwarded")

	/*if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}*/
	fmt.Printf("Post from website! subject = %v\n", r.FormValue("Subject"))
	fmt.Printf("Post from website! body-html = %v\n", r.FormValue("body-html"))
	fmt.Printf("Post from website! sender = %v\n", r.FormValue("sender"))
	fmt.Printf("Post from website! From = %v\n", r.FormValue("From"))
	//fmt.Printf("Post from website! r.Form = %v\n", r.Form)
	//fmt.Printf("Post from website! r.PostForm = %v\n", r.PostForm)

	//fmt.Printf("%v", r)

	domain := os.Getenv("MAILGUN_DOMAINS")
	secretKey := os.Getenv("MAILGUN_PRIVATE_KEY")
	emailSender := os.Getenv("EMAIL_SENDER")
	emailRecipient := os.Getenv("EMAIL_RECIPIENT")

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(domain, secretKey)

	sender := emailSender
	subject := "Test!"
	body := "Hello from Mailgun Go!"
	recipient := emailRecipient

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
