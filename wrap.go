package mailopen

import (
	"os"

	"github.com/gobuffalo/buffalo/mail"
)

// Wrap other sender to be used if env is not development
func Wrap(sender mail.Sender) mail.Sender {
	env := os.Getenv("GO_ENV")

	if env == "" {
		env = "development"
	}

	if env != "development" {
		return sender
	}

	return FileSender{
		Open:    true,
		TempDir: os.TempDir(),
	}
}
