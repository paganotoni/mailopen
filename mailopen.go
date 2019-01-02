package mailopen

import (
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gofrs/uuid"
	"github.com/pkg/browser"
)

const (
	htmlHeaderTmpl = `
<div class="email-information" style="background-color:white; padding: 10px; border-bottom: 1px solid #333;">
	<p style="margin-bottom: 0;"><span style="font-weight: bold;">From:</span> %v </p>
	<p style="margin-bottom: 0;"><span style="font-weight: bold;">To:</span> %v </p>
	<p style="margin-bottom: 0;"><span style="font-weight: bold;">Cc:</span> %v </p>
	<p style="margin-bottom: 0;"><span style="font-weight: bold;">Bcc:</span> %v </p>
	<p style="margin-bottom: 0;"><span style="font-weight: bold;">Subject:</span> %v</p>
</div>
`
	plainHeaderTmpl = `
From: %v
To: %v
Cc: %v
Bcc: %v 
Subject: %v
----------------------------
`
)

//FileSender implements the Sender interface to be used
//within buffalo mailer generated package.
type FileSender struct{}

//Send sends an email to Sendgrid for delivery, it assumes
//bodies[0] is HTML body and bodies[1] is text.
func (ps FileSender) Send(m mail.Message) error {
	if len(m.Bodies) < 2 {
		return errors.New("you must specify at least 2 bodies HTML and plain text")
	}

	htmlContent := ps.addHTMLHeader(m.Bodies[0].Content, m)
	htmlPath, err := ps.saveEmailBody(htmlContent, "html")
	if err != nil {
		return err
	}

	plainContent := ps.addPlainHeader(fmt.Sprintf("<html><head></head><body><pre>%v</pre></body></html>", m.Bodies[1].Content), m)
	plainPath, err := ps.saveEmailBody(plainContent, "txt")
	if err != nil {
		return err
	}

	time.Sleep(100 * time.Millisecond)

	if err := browser.OpenFile(plainPath); err != nil {
		return err
	}

	if err := browser.OpenFile(htmlPath); err != nil {
		return err
	}

	return nil
}

func (ps FileSender) saveEmailBody(content, ctype string) (string, error) {
	path := fmt.Sprintf(path.Join("tmp", "mailopen-%v-%s.html"), uuid.Must(uuid.NewV4()), ctype)
	err := ioutil.WriteFile(path, []byte(content), 0644)
	return path, err
}

func (ps FileSender) addHTMLHeader(content string, m mail.Message) string {
	header := fmt.Sprintf(htmlHeaderTmpl, html.EscapeString(m.From), strings.Join(m.To, ","), strings.Join(m.CC, ","), strings.Join(m.Bcc, ","), html.EscapeString(m.Subject))
	var re = regexp.MustCompile(`(.*<body[^>]*>)((.|[\n\r])*)(<\/body>.*)`)
	return re.ReplaceAllString(content, fmt.Sprintf(`$1%v$2$3`, header))
}

func (ps FileSender) addPlainHeader(content string, m mail.Message) string {
	header := fmt.Sprintf(plainHeaderTmpl, html.EscapeString(m.From), strings.Join(m.To, ","), strings.Join(m.CC, ","), strings.Join(m.Bcc, ","), html.EscapeString(m.Subject))
	var re = regexp.MustCompile(`(.*<pre[^>]*>)((.|[\n\r])*)(<\/pre>.*)`)
	return re.ReplaceAllString(content, fmt.Sprintf(`$1%v$2$3`, header))
}

// New creates a sender that writes emails into disk
func New() FileSender {
	return FileSender{}
}
