package mailopen_test

import (
	_ "embed"
	"fmt"
	"mime"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/paganotoni/mailopen/v2"

	"github.com/stretchr/testify/require"
)

type falseSender struct{}

func (ps falseSender) Send(m mail.Message) error {
	return nil
}

var (
	//go:embed plain-header.txt
	txtFormat string
)

func Test_Send(t *testing.T) {
	r := require.New(t)

	mailopen.Testing = true

	m := mail.NewMessage()
	m.From = "testing@testing.com"
	m.To = []string{"testing@other.com"}
	m.CC = []string{"aa@other.com"}
	m.Bcc = []string{"aax@other.com"}
	m.Subject = "something"

	const testHTMLcontent = `<html><head></head><body><div>Some Message</div></body></html>`

	t.Run("html and plain with attachments", func(t *testing.T) {
		tmpDir := t.TempDir()
		sender := mailopen.WithOptions(mailopen.Directory(tmpDir))
		sender.Open = false

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: testHTMLcontent},
			{ContentType: "text/plain", Content: "Same message"},
		}

		m.Attachments = []mail.Attachment{
			{Name: "txt_test", Reader: strings.NewReader(""), ContentType: "text/plain", Embedded: false},
			{Name: "csv_test", Reader: strings.NewReader(""), ContentType: "text/csv", Embedded: false},
			{Name: "img_test", Reader: strings.NewReader(""), ContentType: "image/jpeg", Embedded: false},
			{Name: "pdf_test", Reader: strings.NewReader(""), ContentType: "application/pdf", Embedded: false},
			{Name: "zip_test", Reader: strings.NewReader(""), ContentType: "application/zip", Embedded: false},
		}

		r.NoError(sender.Send(m))

		htmlFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[0].ContentType, "/", "_")))
		txtFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[1].ContentType, "/", "_")))

		r.FileExists(htmlFile)
		r.FileExists(txtFile)

		htmlHeader, err := os.ReadFile(htmlFile)
		r.NoError(err)

		r.Contains(string(htmlHeader), m.From)
		r.Contains(string(htmlHeader), m.To[0])
		r.Contains(string(htmlHeader), m.CC[0])
		r.Contains(string(htmlHeader), m.Bcc[0])
		r.Contains(string(htmlHeader), m.Subject)

		for _, a := range m.Attachments {
			r.Contains(string(htmlHeader), a.Name)

			ext, err := mime.ExtensionsByType(a.ContentType)
			r.NoError(err)

			filePath := path.Join(tmpDir, fmt.Sprintf("%s%s", a.Name, ext[0]))
			r.FileExists(filePath)
		}

		txtContent, err := os.ReadFile(txtFile)
		r.NoError(err)

		r.Contains(string(txtContent), m.From)
		r.Contains(string(txtContent), m.To[0])
		r.Contains(string(txtContent), m.CC[0])
		r.Contains(string(txtContent), m.Bcc[0])
		r.Contains(string(txtContent), m.Subject)
		r.Contains(string(txtContent), "Same message")

		r.Contains(string(txtContent), fmt.Sprintf(txtFormat, m.From, m.To[0], m.CC[0], m.Bcc[0], m.Subject))
	})

	t.Run("html only", func(t *testing.T) {
		tmpDir := t.TempDir()
		sender := mailopen.WithOptions(mailopen.Only("text/html"), mailopen.Directory(tmpDir))
		sender.Open = false

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: testHTMLcontent},
			{ContentType: "text/plain", Content: "Same message"},
		}

		htmlFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[0].ContentType, "/", "_")))
		txtFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[1].ContentType, "/", "_")))

		r.NoError(sender.Send(m))

		r.FileExists(htmlFile)
		r.NoFileExists(txtFile)

		dat, err := os.ReadFile(htmlFile)
		r.NoError(err)

		r.NotContains(string(dat), "Attachment:")
	})

	t.Run("plain only", func(t *testing.T) {
		tmpDir := t.TempDir()
		sender := mailopen.WithOptions(mailopen.Only("text/plain"), mailopen.Directory(tmpDir))
		sender.Open = false

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: testHTMLcontent},
			{ContentType: "text/plain", Content: "Same message"},
		}

		htmlFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[0].ContentType, "/", "_")))
		txtFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[1].ContentType, "/", "_")))

		r.NoError(sender.Send(m))

		r.NoFileExists(htmlFile)
		r.FileExists(txtFile)
	})

	t.Run("long subject and long file name`", func(t *testing.T) {
		tmpDir := t.TempDir()
		sender := mailopen.WithOptions(mailopen.Directory(tmpDir))
		sender.Open = false

		m := mail.NewMessage()

		m.Subject = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam nec leo tellus. Aliquam ac facilisis est, condimentum pellentesque velit. In quis erat turpis. Morbi accumsan ante nec nunc dapibus, quis lacinia mi ornare. Vivamus venenatis accumsan dolor ac placerat. Sed pulvinar sem eu est accumsan, ut commodo mi viverra. Quisque turpis metus, ultrices id mauris vel, suscipit sollicitudin erat. Vivamus eget quam non sem volutpat eleifend eget in lacus. In vulputate, justo fringilla lacinia lobortis, neque turpis dignissim tellus, in placerat eros justo nec massa. Duis ex enim, convallis ut leo nec, condimentum consectetur mi. Vestibulum imperdiet pharetra ipsum. Etiam venenatis tincidunt odio, sed feugiat quam blandit sit amet. Donec eget nulla dui."
		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: testHTMLcontent},
			{ContentType: "text/plain", Content: "Same message"},
		}

		m.Attachments = []mail.Attachment{
			{Name: "123456789-123456789-123456789-123456789-123456789-1", Reader: strings.NewReader(""), ContentType: "text/plain", Embedded: false},
		}

		r.NoError(sender.Send(m))

		att := m.Attachments[0]

		exts, err := mime.ExtensionsByType(att.ContentType)
		r.NoError(err)

		filePath := path.Join(tmpDir, fmt.Sprintf("%s%s", att.Name[0:50], exts[0]))
		r.FileExists(filePath)
	})

	t.Run("only one body", func(t *testing.T) {
		sender := mailopen.WithOptions(mailopen.Directory(t.TempDir()))
		sender.Open = false

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: testHTMLcontent},
		}

		r.Error(sender.Send(m))
	})

	t.Run("with custom folder", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.Setenv(mailopen.MailOpenDirKey, tmpDir)
		sender := mailopen.WithOptions()
		sender.Open = false

		m.Bodies = []mail.Body{
			{ContentType: "text/html", Content: testHTMLcontent},
			{ContentType: "text/plain", Content: "Same message"},
		}

		r.NoError(sender.Send(m))

		htmlFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[0].ContentType, "/", "_")))
		txtFile := path.Join(tmpDir, fmt.Sprintf("%s_body.html", strings.ReplaceAll(m.Bodies[1].ContentType, "/", "_")))

		r.FileExists(htmlFile)
		r.FileExists(txtFile)
	})
}

func Test_Wrap(t *testing.T) {
	mailopen.Testing = true

	r := require.New(t)

	os.Setenv("GO_ENV", "development")
	s := mailopen.Wrap(falseSender{})
	r.IsType(mailopen.FileSender{}, s)

	os.Setenv("GO_ENV", "")
	s = mailopen.Wrap(falseSender{})
	r.IsType(mailopen.FileSender{}, s)

	os.Setenv("GO_ENV", "staging")
	s = mailopen.Wrap(falseSender{})
	r.IsType(falseSender{}, s)

	os.Setenv("GO_ENV", "production")
	s = mailopen.Wrap(falseSender{})
	r.IsType(falseSender{}, s)

}
