package mail

import (
	"context"
	"s2p-api/config"
	"time"

	"github.com/mailgun/mailgun-go/v4"
)

var (
	mg      *mailgun.MailgunImpl
	senders = map[string]string{
		"noreply": "noreply@start2play.games",
	}
)

func Initialize() {
	mg = mailgun.NewMailgun(
		config.Mail.Domain,
		config.Mail.Key)
}

func Send(message *mailgun.Message) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mg.Send(ctx, message)
}

func SendVerificationTo(email string, name string, token string) {
	message := mg.NewMessage(
		senders["noreply"],
		"Confirme Seu Endereço De Email",
		"",
		email)

	message.SetTemplate("verify")

	message.AddTemplateVariable("first_name", name)
	message.AddTemplateVariable("token", token)

	Send(message)
}

func SendRecoveryTo(email string, name string, token string) {
	message := mg.NewMessage(
		senders["noreply"],
		"Recuperar Senha",
		"",
		email)

	message.SetTemplate("recovery")

	message.AddTemplateVariable("first_name", name)
	message.AddTemplateVariable("token", token)

	Send(message)
}
