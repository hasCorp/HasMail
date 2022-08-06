package hasmailtemplates

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var (
	clientInit     sync.Once
	sendgridClient MailSender
	MailSenderFunc func(string) MailSender = defaultSendGridClient
	fromEmailAddr  string
	fromName       string
)

var defaultSendGridClient = func(apiKey string) MailSender {
	return sendgrid.NewSendClient(apiKey)
}

type MailSender interface {
	SendWithContext(ctx context.Context, email *mail.SGMailV3) (*rest.Response, error)
}

type MailTemplate interface {
	Body() string
	Name() string
	Subject() string
	// TODO: not sure what else is needed
}

type MailRequest struct {
	Recipient string            `json:"recipient"`
	Name      string            `json:"name,omitempty"`
	Vars      map[string]string `json:"vars"`
}

type MailResponse struct {
	Success bool `json:"success"`
}

// InjectVars takes a MailTemplate and a key-value map that
// is used to map template variables in the mail body to resolved
// values from the input map.
// This function expects the vars in the template to be formatted as
// "{var}", where `var` is the key, and it is surrounded by curly braces.
// This makes it so that the request body can contain a simple key-value
// pair, e.g.:
// {
//   "var": "foo"
// }
func InjectVars(t MailTemplate, vars map[string]string) (string, error) {
	parsed := t.Body()
	log.Printf("Mail body template: %s\n", parsed)

	for k, v := range vars {
		parsed = strings.ReplaceAll(parsed, fmt.Sprintf("{%s}", k), v)
	}

	// TODO: handle scenario where the vars does not satisfy the
	//       full set of required variables to replace in the template
	return parsed, nil
}

// SendResponse writes a response back to the HTTP client
func SendResponse(w http.ResponseWriter, code int) {
	resp := MailResponse{
		Success: code < 400,
	}
	b, _ := json.Marshal(resp)
	w.WriteHeader(code)
	_, err := w.Write(b)
	if err != nil {
		log.Println("failed to write response data for request", err)
	}
}

func SendEmail(ctx context.Context,
	subject string,
	to *mail.Email,
	plainTextContent string,
	htmlContent string) (*rest.Response, error) {
	clientInit.Do(loadEmailClient)
	from := mail.NewEmail(fromName, fromEmailAddr)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	return sendgridClient.SendWithContext(ctx, message)
}

func loadEmailClient() {
	log.Println("Init mail client auth")

	var apiKey string

	// first try to read from the credentials.json file
	f, err := os.Open("credentials.json")

	if err == nil {
		// file exists, so read from it
		fromEmailAddr, fromName, apiKey = readConfigFromFile(f)
		if err = f.Close(); err != nil {
			log.Fatal("failed to close config file", err)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		// read from environment
		fromEmailAddr, fromName, apiKey = readConfigFromEnv()
	} else {
		log.Fatal("unexpected error checking credential file", err)
	}

	sendgridClient = MailSenderFunc(apiKey)
	log.Println("SendGrid client initialized")
}

func readConfigFromEnv() (addr string, name string, apiKey string) {
	log.Println("Reading configs from env")

	var ok bool
	addr, ok = os.LookupEnv("FROM_ADDR")
	if !ok {
		log.Fatal("from email address not found in environment.")
	}

	name, ok = os.LookupEnv("FROM_NAME")
	if !ok {
		log.Fatal("from friendly name not found in environment.")
	}

	apiKey, ok = os.LookupEnv("SENDGRID_API_KEY")
	if !ok {
		log.Fatal("API Key not found in environment.")
	}
	return
}

func readConfigFromFile(f *os.File) (addr string, name string, apiKey string) {
	log.Println("Reading configs from credential file")
	bytes, err := io.ReadAll(f)
	if err != nil {
		log.Fatal("unexpected error occurred reading credential file", err)
	}

	var creds map[string]string
	err = json.Unmarshal(bytes, &creds)
	if err != nil {
		log.Fatal("unexpected error occurred reading credential file", err)
	}

	var ok bool
	addr, ok = creds["FROM_ADDR"]
	if !ok {
		log.Fatal("Missing FROM_ADDR from credential file")
	}
	name, ok = creds["FROM_NAME"]
	if !ok {
		log.Fatal("Missing FROM_NAME from credential file")
	}
	apiKey, ok = creds["SENDGRID_API_KEY"]
	if !ok {
		log.Fatal("Missing SENDGRID_API_KEY from credential file")
	}
	return
}
