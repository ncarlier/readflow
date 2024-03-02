package test

import (
	"context"
	"fmt"
	"net"
	"net/smtp"
	"testing"

	"github.com/brianvoe/gofakeit"
	smtpd "github.com/emersion/go-smtp"
	"github.com/ncarlier/readflow/internal/mail"
	"github.com/ncarlier/readflow/internal/model"
	service "github.com/ncarlier/readflow/internal/service"
	serviceT "github.com/ncarlier/readflow/internal/service/test"
	"github.com/stretchr/testify/require"
)

var listener net.Listener

func setupMailServer(t *testing.T) func(t *testing.T) {
	t.Log("setup test mail server")
	var err error
	listener, err = net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	smtpServer := smtpd.NewServer(mail.NewBackend())
	smtpServer.Domain = "localhost"
	smtpServer.AllowInsecureAuth = true

	go smtpServer.Serve(listener)

	return func(t *testing.T) {
		t.Log("teardown test mail server")
		defer smtpServer.Shutdown(context.Background())
	}
}

func buildFakeMail(to string) (msg []byte, from string) {
	title := gofakeit.Sentence(5)
	text := gofakeit.Paragraph(2, 2, 5, ".")
	from = gofakeit.Email()
	msg = []byte(
		"Content-Type: text/plain \r\n" +
			"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + title + "\r\n" +
			"\r\n" + text + "\r\n")
	return
}

func createIncomingWebhook(t *testing.T, alias string) func(t *testing.T) {
	// test context
	ctx := serviceT.GetTestContext()

	// create new webhook
	builder := model.NewIncomingWebhookCreateFormBuilder()
	builder.Alias(alias)
	webhook, err := service.Lookup().CreateIncomingWebhook(ctx, *builder.Build())
	require.Nil(t, err)
	require.Equal(t, alias, webhook.Alias)
	require.NotEmpty(t, webhook.Token)
	return func(t *testing.T) {
		t.Log("teardown test mail server")
		_, err := service.Lookup().DeleteIncomingWebhook(ctx, *webhook.ID)
		require.Nil(t, err)
	}
}

func TestSendMail(t *testing.T) {
	teardownTestCase := serviceT.SetupTestCase(t)
	defer teardownTestCase(t)

	teardownTestMailServer := setupMailServer(t)
	defer teardownTestMailServer(t)

	// create incoming webhook
	alias := "test"
	cleanup := createIncomingWebhook(t, alias)
	defer cleanup(t)

	// build destination
	user := serviceT.GetTestUser()
	to := fmt.Sprintf("%s-%s@example.com", alias, service.Lookup().GetUserHashID(*user.ID))
	// build mail
	msg, from := buildFakeMail(to)

	// Sending email.
	err := smtp.SendMail(listener.Addr().String(), nil, from, []string{to}, msg)
	require.Nil(t, err)
}
