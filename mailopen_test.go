package mailopen_test

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gobuffalo/flect"
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
	mailopen.Testing = true

	r := require.New(t)
	sender := mailopen.WithOptions()
	sender.Open = false

	m := mail.NewMessage()
	m.From = "testing@testing.com"
	m.To = []string{"testing@other.com"}
	m.CC = []string{"aa@other.com"}
	m.Bcc = []string{"aax@other.com"}
	m.Subject = "something"

	m.Bodies = []mail.Body{
		{ContentType: "text/html", Content: "<html><head></head><body><div>Some Message</div></body></html>"},
		{ContentType: "text/plain", Content: "Some message"},
	}

	m.Attachments = []mail.Attachment{
		{Name: "csv_test", Reader: openFile(filepath.Join("test_files", "csv_sample.csv"), r), ContentType: "text/csv", Embedded: false},
		{Name: "img_test", Reader: openFile(filepath.Join("test_files", "img_sample.jpeg"), r), ContentType: "image/jpeg", Embedded: false},
		{Name: "pdf_test", Reader: openFile(filepath.Join("test_files", "pdf_sample.pdf"), r), ContentType: "application/pdf", Embedded: false},
		{Name: "zip_test", Reader: openFile(filepath.Join("test_files", "zip_sample.zip"), r), ContentType: "application/zip", Embedded: false},
	}

	r.NoError(sender.Send(m))

	htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_0.html", flect.Underscore(m.Subject)))
	txtFile := path.Join(sender.TempDir, fmt.Sprintf("%s_1.html", flect.Underscore(m.Subject)))

	r.FileExists(htmlFile)
	r.FileExists(txtFile)

	htmlContent, err := ioutil.ReadFile(htmlFile)
	r.NoError(err)

	r.Contains(string(htmlContent), m.From)
	r.Contains(string(htmlContent), m.To[0])
	r.Contains(string(htmlContent), m.CC[0])
	r.Contains(string(htmlContent), m.Bcc[0])
	r.Contains(string(htmlContent), m.Subject)
	r.Contains(string(htmlContent), "Some Message")

	for _, a := range m.Attachments {
		r.Contains(string(htmlContent), a.Name)
	}

	txtContent, err := ioutil.ReadFile(txtFile)
	r.NoError(err)

	r.Contains(string(txtContent), m.From)
	r.Contains(string(txtContent), m.To[0])
	r.Contains(string(txtContent), m.CC[0])
	r.Contains(string(txtContent), m.Bcc[0])
	r.Contains(string(txtContent), m.Subject)
	r.Contains(string(txtContent), "Some message")

	r.Contains(string(txtContent), fmt.Sprintf(txtFormat, m.From, m.To[0], m.CC[0], m.Bcc[0], m.Subject))
}

func Test_SendWithOptionsOnlyHTML(t *testing.T) {
	r := require.New(t)

	mailopen.Testing = true

	sender := mailopen.WithOptions(mailopen.Only("text/html"))
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

	htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_0.html", flect.Underscore(m.Subject)))
	txtFile := path.Join(sender.TempDir, fmt.Sprintf("%s_1.html", flect.Underscore(m.Subject)))

	os.Remove(htmlFile)
	os.Remove(txtFile)

	r.NoError(sender.Send(m))

	r.FileExists(htmlFile)
	r.NoFileExists(txtFile)
}

func Test_SendWithOptionsOnlyTXT(t *testing.T) {
	r := require.New(t)

	mailopen.Testing = true

	sender := mailopen.WithOptions(mailopen.Only("text/plain"))
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

	htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_0.html", flect.Underscore(m.Subject)))
	txtFile := path.Join(sender.TempDir, fmt.Sprintf("%s_1.html", flect.Underscore(m.Subject)))

	os.Remove(htmlFile)
	os.Remove(txtFile)

	r.NoError(sender.Send(m))

	r.NoFileExists(htmlFile)
	r.FileExists(txtFile)
}

func Test_SendWithOneBody(t *testing.T) {
	mailopen.Testing = true

	r := require.New(t)
	sender := mailopen.WithOptions()
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
	mailopen.Testing = true

	r := require.New(t)
	sender := mailopen.WithOptions()
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

	htmlFile := path.Join(sender.TempDir, fmt.Sprintf("%s_%s.html", flect.Underscore(m.Subject), "0"))

	dat, err := ioutil.ReadFile(htmlFile)
	r.NoError(err)
	r.NotContains(string(dat), "Attachment:")
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

func openFile(name string, r *require.Assertions) *os.File {
	f, err := os.Open(name)
	r.NoError(err)

	return f
}
