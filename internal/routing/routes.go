package routing

import (
	"github.com/gin-gonic/gin"
	"github.com/wneessen/go-mail"
	"log"
	"time"
)

func Add(e *gin.Engine) {
	initializeDependencies()

	authRoutes(e)
	accountRoutes(e)

	e.POST("/send-email", func(ctx *gin.Context) {
		msg := mail.NewMsg()

		if err := msg.From("your-email@test.com"); err != nil {
			log.Fatalf("failed to set FROM address: %s", err)
		}
		if err := msg.To(ctx.PostForm("email")); err != nil {
			log.Fatalf("failed to set TO address: %s", err)
		}

		msg.Subject(ctx.PostForm("subject"))
		msg.SetBodyString(mail.TypeTextPlain, ctx.PostForm("body"))

		client, err := mail.NewClient(
			"localhost",
			mail.WithPort(1025),
			mail.WithSMTPAuth(mail.SMTPAuthNoAuth),
			mail.WithTimeout(10*time.Second),
			mail.WithTLSPolicy(mail.NoTLS),
		)

		if err != nil {
			log.Fatalf("failed to create new mail delivery client: %s", err)
		}
		if err := client.DialAndSend(msg); err != nil {
			log.Fatalf("failed to deliver mail: %s", err)
		}
		log.Printf("Test mail successfully delivered.")
	})
}
