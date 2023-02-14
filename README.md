[![Build Status](https://travis-ci.org/paganotoni/mailopen.svg?branch=master)](https://travis-ci.org/paganotoni/mailopen)

## Mailopen

Mailopen is a buffalo mailer that allows to see sent emails in the browser instead of sending these using SMTP or other sender used in production environments.

### Usage

Mailopen is only intended for development purposes, the way you use it is by simply initializing your mailer to be a mailopen instance instead of your regular sender, p.e:

```go
import (
    ...
    "github.com/gobuffalo/buffalo/mail"
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

Internally `Wrap` function returns `mailopen.FileSender` instance only if GO_ENV is `development`, otherwise it will return passed sender.

You can always write it yourself in case your conditions to switch sender are not only to be in the `development` environment.

```go
import (
    ...
    "github.com/gobuffalo/buffalo/mail"
    sendgrid "github.com/paganotoni/sendgrid-sender"
    ...
)

func init() {
    if envy.Get("GO_ENV", "development") == "development" {
        Sender = mailopen.WithOptions(mailopen.Only("text/html"))
        
		return
    }

    Sender = sendgrid.NewSendgridSender(envy.Get("SENDGRID_API_KEY", ""))
}
```

Then you will use your `Sender` instance as usual by calling `Sender.Send(m)` with the variant that in development it will open your emails in the browser.

By default, mailopen will save the emails and attachments in a temporary directory. You can customize this by passing options to the [`mailopen.WithOptions`](#options) or setting the `MAILOPEN_DIR` env variable in your machine.

### Options

You can pass options to the `mailopen.WithOptions` function to customize the way it work.

- `Only` allows you to specify which content types you want to open in the browser, p.e: `mailopen.Only("text/html")`.

- `Directory` allows you to specify the directory where the emails and attachments will be saved, p.e: `mailopen.Directory("/tmp")`.

