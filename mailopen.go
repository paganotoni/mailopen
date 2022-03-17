package mailopen

import (
	"fmt"
	"os"
)

var Testing = false

// New creates a sender that writes emails into disk
func New() FileSender {
	fmt.Println("Deprecated: use WithOptions instead.")

	return WithOptions()
}

// WithOptions creates a sender that writes emails into disk
// And applies the passed options.
func WithOptions(options ...Option) FileSender {
	s := FileSender{
		Open:    true,
		TempDir: os.TempDir(),
	}

	for _, option := range options {
		option(&s)
	}

	return s
}
