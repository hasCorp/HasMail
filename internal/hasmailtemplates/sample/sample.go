package sample

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/hascorp/hasmail/internal/hasmailtemplates"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type template struct{}

func (t *template) Name() string {
	return "sample"
}

func (t *template) Body() string {
	return "This is a test email. Congrats!\nThe value of foo is: {foo}"
}

func (t *template) Subject() string {
	return "Test Email"
}

func SampleHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	log.Println("Headers:")
	for k, v := range r.Header {
		log.Printf("\tKey: %v, Value: %v\n", k, v)
	}

	// read request body
	var req hasmailtemplates.MailRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil || len(b) < 1 {
		// missing request body
		hasmailtemplates.SendResponse(w, http.StatusBadRequest)
		return
	}
	log.Printf("Request body: %s\n", string(b))
	if err = json.Unmarshal(b, &req); err != nil {
		// request body didn't marshal into a map properly
		hasmailtemplates.SendResponse(w, http.StatusBadRequest)
		return
	}

	to, err := mail.ParseEmail(req.Recipient)
	if err != nil {
		log.Printf("Couldn't parse recipient %s as email", req.Recipient)
		hasmailtemplates.SendResponse(w, http.StatusBadRequest)
		return
	}
	if len(req.Name) > 0 {
		to.Name = req.Name
	}

	// use the template
	// TODO: this stuff
	t := template{}
	mailBody, err := hasmailtemplates.InjectVars(&t, req.Vars)
	if err != nil {
		hasmailtemplates.SendResponse(w, http.StatusBadRequest)
		return
	}
	log.Printf("Parsed body: %s\n", mailBody)
	response, err := hasmailtemplates.SendEmail(r.Context(), t.Subject(), to, mailBody, mailBody)

	if response != nil {
		log.Printf("SendGrid responded with code %d\n", response.StatusCode)
	}

	if err != nil {
		log.Printf("error occurred: %v", err)
		code := http.StatusInternalServerError
		if response != nil {
			code = response.StatusCode
		}
		hasmailtemplates.SendResponse(w, code)
		return
	}

	log.Println("successfully sent email")
	hasmailtemplates.SendResponse(w, http.StatusAccepted) // assumes response is not nil when an error occurs
}
