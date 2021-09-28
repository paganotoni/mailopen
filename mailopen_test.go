package mailopen_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/flect"
	"github.com/paganotoni/mailopen"

	"github.com/stretchr/testify/require"
)

func init() {
	mailopen.Testing = true
}

func openFile(name string, r *require.Assertions) *os.File {
	f, err := os.Open(name)
	r.NoError(err)

	return f
}

//TODO: options open only
func Test_Send(t *testing.T) {
	r := require.New(t)
	sender := mailopen.New()
	sender.Open = false

	m := mail.NewMessage()
	m.From = "testing@testing.com"
	m.To = []string{"testing@other.com"}
	m.CC = []string{"aa@other.com"}
	m.Bcc = []string{"aax@other.com"}
	m.Subject = "something"
	m.Bodies = []mail.Body{
		{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
		{ContentType: "text/plain", Content: "Same message"},
	}
	m.Attachments = []mail.Attachment{
		{Name: "csv_test", Reader: openFile("test_files/csv_sample.csv", r), ContentType: "text/csv", Embedded: false},
		{Name: "img_test", Reader: openFile("test_files/img_sample.jpeg", r), ContentType: "image/jpeg", Embedded: false},
		{Name: "pdf_test", Reader: openFile("test_files/pdf_sample.pdf", r), ContentType: "application/pdf", Embedded: false},
		{Name: "zip_test", Reader: openFile("test_files/zip_sample.zip", r), ContentType: "application/zip", Embedded: false},
	}

	r.NoError(sender.Send(m))
	htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_%s.html", flect.Underscore(m.Subject), "html"))
	txtFile := path.Join(sender.TempDir, fmt.Sprintf("%s_%s.html", flect.Underscore(m.Subject), "txt"))

	r.FileExists(htmlFile)
	r.FileExists(txtFile)

	dat, err := ioutil.ReadFile(htmlFile)
	r.NoError(err)

	r.Contains(string(dat), m.From)
	r.Contains(string(dat), m.To[0])
	r.Contains(string(dat), m.CC[0])
	r.Contains(string(dat), m.Bcc[0])
	r.Contains(string(dat), m.Subject)

	for _, a := range m.Attachments {
		r.Contains(string(dat), a.Name)
	}

	dat, err = ioutil.ReadFile(txtFile)
	r.NoError(err)

	r.Contains(string(dat), m.From)
	r.Contains(string(dat), m.To[0])
	r.Contains(string(dat), m.CC[0])
	r.Contains(string(dat), m.Bcc[0])
	r.Contains(string(dat), m.Subject)
}

func Test_SendWithOneBody(t *testing.T) {
	r := require.New(t)
	sender := mailopen.New()
	sender.Open = false

	m := mail.NewMessage()
	m.From = "testing@testing.com"
	m.To = []string{"testing@other.com"}
	m.CC = []string{"aa@other.com"}
	m.Bcc = []string{"aax@other.com"}
	m.Subject = "something"
	m.Bodies = []mail.Body{
		{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
	}

	r.Error(sender.Send(m))
}

func Test_SendWithoutAttachments(t *testing.T) {
	r := require.New(t)
	sender := mailopen.New()
	sender.Open = false

	m := mail.NewMessage()
	m.From = "testing@testing.com"
	m.To = []string{"testing@other.com"}
	m.CC = []string{"aa@other.com"}
	m.Bcc = []string{"aax@other.com"}
	m.Subject = "something"
	m.Bodies = []mail.Body{
		{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
		{ContentType: "text/plain", Content: "Same message"},
	}

	r.NoError(sender.Send(m))

	htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_%s.html", flect.Underscore(m.Subject), "html"))

	dat, err := ioutil.ReadFile(htmlFile)
	r.NoError(err)
	r.NotContains(string(dat), "Attachment:")
}

func Test_Wrap(t *testing.T) {
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

type falseSender struct{}

func (ps falseSender) Send(m mail.Message) error {
	return nil
}
