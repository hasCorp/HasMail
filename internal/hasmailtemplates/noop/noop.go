package noop

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hascorp/hasmail/internal/hasmailtemplates"
)

type template struct{}

func (t *template) Name() string {
	return "noop"
}

func (t *template) Body() string {
	return `This is an email`
}

func (t *template) Subject() string {
	return "No Op"
}

func NoOpHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	vars := mux.Vars(r)
	log.Println("Accepted request from", r.Host)
	log.Println("Headers:")
	for k, v := range r.Header {
		log.Printf("\tKey: %v, Value: %v\n", k, v)
	}
	log.Println("mux Vars:")
	for k, v := range vars {
		log.Printf("\tKey: %v, Value: %v\n", k, v)
	}

	// read request body
	var body map[string]string
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request: %v", err)
		hasmailtemplates.SendResponse(w, 400)
	}
	log.Printf("Request body: %s\n", string(b))
	if err = json.Unmarshal(b, &body); err != nil {
		log.Printf("Failed to unmarshal request as JSON: %v", err)
		hasmailtemplates.SendResponse(w, 400)
	}

	// use the template
	mailBody, err := hasmailtemplates.InjectVars(&template{}, body)
	if err != nil {
		log.Printf("Failed to inject vars into email template: %v", err)
		hasmailtemplates.SendResponse(w, 400)
	}
	log.Printf("Parsed body: %s\n", mailBody)

	resp := hasmailtemplates.MailResponse{
		Success: true,
	}
	bytes, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusAccepted)
	_, err = w.Write(bytes)
	if err != nil {
		log.Println("failed to write response back", err)
	}
}
