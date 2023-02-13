package mailopen

import (
	"bytes"
	_ "embed"
	"fmt"
	"html"
	"html/template"
	"io"
	"mime"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/gobuffalo/buffalo/mail"
	"github.com/gofrs/uuid"
	"github.com/pkg/browser"
)

var (

	//go:embed html-header.html
	htmlHeader string

	//go:embed plain-header.txt
	plainHeader string

	// config used to write the email body to files
	// depending on the type we need to do a few things differently.
	config = map[string]contentConfig{
		"text/html": {
			headerTemplate: htmlHeader,
			replaceRegexp:  `(.*<body[^>]*>)((.|[\n\r])*)(<\/body>.*)`,
		},

		"text/plain": {
			headerTemplate: plainHeader,
			replaceRegexp:  `(.*<pre[^>]*>)((.|[\n\r])*)(<\/pre>.*)`,

			preformatter: func(s string) string {
				return fmt.Sprintf("<html><head></head><body><pre>%v</pre></body></html>", s)
			},
		},
	}
)

type contentConfig struct {
	headerTemplate string
	replaceRegexp  string
	preformatter   func(string) string
}

// FileSender implements the Sender interface to be used
// within buffalo mailer generated package.
type FileSender struct {
	Open    bool
	TempDir string

	// openContentTypes are those content types to open in browser
	openContentTypes []string
}

// AttFile is a file to be attached to the email
type AttFile struct {
	Path string
	Name string
}

func (ps FileSender) shouldOpen(contentType string) bool {
	if len(ps.openContentTypes) == 0 {
		return true
	}

	for _, v := range ps.openContentTypes {
		if v == contentType {
			return true
		}
	}

	return false
}

// Send sends an email to Sendgrid for delivery, it assumes
// bodies[0] is HTML body and bodies[1] is text.
func (ps FileSender) Send(m mail.Message) error {
	if len(m.Bodies) < 2 {
		return fmt.Errorf("mailopen: expected at least 2 bodies, got %d", len(m.Bodies))
	}

	for _, v := range m.Bodies {
		if !ps.shouldOpen(v.ContentType) {
			continue
		}

		cc := config[v.ContentType]
		if cc.preformatter != nil {
			v.Content = cc.preformatter(v.Content)
		}

		header := fmt.Sprintf(
			cc.headerTemplate,

			html.EscapeString(m.From),
			strings.Join(m.To, ","),
			strings.Join(m.CC, ","),
			strings.Join(m.Bcc, ","),

			html.EscapeString(m.Subject),
		)

		var re = regexp.MustCompile(cc.replaceRegexp)
		content := re.ReplaceAllString(v.Content, fmt.Sprintf("$1\n%v\n$2$3", header))
		tmpName := strings.ReplaceAll(v.ContentType, "/", "_") + "_body"

		path, err := ps.saveEmailBody(content, tmpName, m)
		if err != nil {
			return err
		}

		if Testing {
			continue
		}

		if err := browser.OpenFile(path); err != nil {
			return err
		}
	}

	return nil
}

func (ps FileSender) saveEmailBody(content, tmpName string, m mail.Message) (string, error) {
	id := uuid.Must(uuid.NewV4())

	afs, err := ps.saveAttachmentFiles(m.Attachments)
	if err != nil {
		return "", fmt.Errorf("mailopen: failed to save attachments: %w", err)
	}

	tmpl := template.Must(template.New("mail").Parse(content))
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, afs)
	if err != nil {
		return "", fmt.Errorf("mailopen: failed to execute template: %w", err)
	}

	filePath := fmt.Sprintf("%s_%s.html", tmpName, id)
	if Testing {
		filePath = fmt.Sprintf("%s.html", tmpName)
	}

	path := path.Join(ps.TempDir, filePath)
	err = os.WriteFile(path, tpl.Bytes(), 0644)

	return path, err
}

func (ps FileSender) saveAttachmentFiles(Attachments []mail.Attachment) ([]AttFile, error) {
	var afs []AttFile

	for _, a := range Attachments {
		if len(a.Name) > 50 {
			a.Name = a.Name[:50]
		}

		exts, err := mime.ExtensionsByType(a.ContentType)
		if err != nil {
			return []AttFile{}, fmt.Errorf("mailopen: failed to get extension for content type %s: %w", a.ContentType, err)
		}

		filePath := path.Join(ps.TempDir, fmt.Sprintf("%s_%s%s", uuid.Must(uuid.NewV4()), a.Name, exts[0]))
		if Testing {
			filePath = path.Join(ps.TempDir, fmt.Sprintf("%s%s", a.Name, exts[0]))
		}

		b, err := io.ReadAll(a.Reader)
		if err != nil {
			return []AttFile{}, fmt.Errorf("mailopen: failed to read attachment %s: %w", a.Name, err)
		}

		err = os.WriteFile(filePath, b, 0644)
		if err != nil {
			return []AttFile{}, fmt.Errorf("mailopen: failed to write attachment %s: %w", a.Name, err)
		}

		afs = append(afs, AttFile{
			Path: filePath,
			Name: a.Name,
		})
	}

	return afs, nil
}
