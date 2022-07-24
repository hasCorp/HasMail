package noop

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hascorp/hasmail/internal/hasmailtemplates"
)

type NoOpTemplate struct{}

func (t *NoOpTemplate) Name() string {
	return "noop"
}

func (t *NoOpTemplate) Body() string {
	return `This is an email`
}

func (t *NoOpTemplate) Subject() string {
	return "No Op"
}

func NoOpHandler(w http.ResponseWriter, r *http.Request) {
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
		// TODO: handle this
	}
	log.Printf("Request body: %s\n", string(b))
	if err := json.Unmarshal(b, &body); err != nil {
		// TODO: handle this
	}

	// use the template
	mailBody, err := hasmailtemplates.InjectVars(&NoOpTemplate{}, body)
	if err != nil {
		// TODO: handle this
	}
	log.Printf("Parsed body: %s\n", mailBody)

	resp := hasmailtemplates.MailResponse{
		Success: true,
	}
	bytes, err := json.Marshal(resp)
	if err != nil {
		// no real good reason this should ever happen
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(make([]byte, 0))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(bytes)
}
