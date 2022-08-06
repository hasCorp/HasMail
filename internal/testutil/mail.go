package testutil

import (
	"context"
	"fmt"
	"os"

	"github.com/hascorp/hasmail/internal/hasmailtemplates"
	"github.com/sendgrid/rest"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type NoOpMailSender struct{}

type FailWithCode struct{}

func (s *NoOpMailSender) SendWithContext(ctx context.Context, email *mail.SGMailV3) (*rest.Response, error) {
	var err error
	mockCode := 202 // happy path
	v := ctx.Value(FailWithCode{})
	if v != nil {
		mockCode = v.(int)
		err = fmt.Errorf("mocked error with code %d", mockCode)
	}

	res := rest.Response{
		StatusCode: mockCode,
	}
	return &res, err
}

func MockEmailClient() {
	os.Setenv("FROM_ADDR", "foo@example.com")
	os.Setenv("FROM_NAME", "bar")
	os.Setenv("SENDGRID_API_KEY", "abc123")
	hasmailtemplates.MailSenderFunc = func(s string) hasmailtemplates.MailSender {
		return &NoOpMailSender{}
	}
}
