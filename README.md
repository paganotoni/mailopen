## Mailopen

Mailopen is a buffalo mailer that allows to see sent emails in the browser instead of sending these using SMTP or other sender used in production environments.

### Usage

Mailopen is only intended for development purposes, the way you use it is by simply initialyzing your mailer to be a mailopen instance instead of your regular sender, p.e:

```go
//Sender allows us to send emails
var Sender mail.Sender

func init() {
	...

	if envy.Get("GO_ENV", "development") == "development" {
        Sender = mailopen.New()
		return
    }
    
    Sender = sendgrid.NewSendgridSender(envy.Get("SENDGRID_API_KEY", ""))

	
}
```

By default mailopen writes files to the `tmp` folder. 