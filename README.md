[![Build Status](https://travis-ci.org/paganotoni/mailopen.svg?branch=master)](https://travis-ci.org/paganotoni/mailopen)

## Mailopen

Mailopen is a buffalo mailer that allows to see sent emails in the browser instead of sending these using SMTP or other sender used in production environments.

### Usage

Mailopen is only intended for development purposes, the way you use it is by simply initialyzing your mailer to be a mailopen instance instead of your regular sender, p.e:

```go
import (
    ...
    "github.com/paganotoni/gonbuffalo"
    sendgrid "github.com/paganotoni/sendgrid-sender"
    ...
)

//Sender allows us to send emails
var Sender mail.Sender

func init() {
    sgSender := sendgrid.NewSendgridSender(envy.Get("SENDGRID_API_KEY", ""))
    Sender = mailopen.Wrap(sgSender)
}
```

Internally Wrap function returns mailopen.FileSender instance only if GO_ENV is `development`, otherwise it will return passed sender.

You can always write it yourself in case your conditions to switch sender are not only to be in the development environment.

```go
import (
    ...
    "github.com/paganotoni/gonbuffalo"
    sendgrid "github.com/paganotoni/sendgrid-sender"
    ...
)

func init() {
    if envy.Get("GO_ENV", "development") == "development" {
        Sender = mailopen.New()
		return
    }

    Sender = sendgrid.NewSendgridSender(envy.Get("SENDGRID_API_KEY", ""))
}
```

Then you will use your `Sender` instance as usual by calling `Sender.Send(m)` with the variant that in development it will open your emails in the browser.