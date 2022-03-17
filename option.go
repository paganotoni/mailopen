package mailopen

type Option func(*FileSender)

// Only opens specified content types, without this
// mailopen opens all of the bodies in the message.
func Only(contentTypes ...string) Option {
	return func(fo *FileSender) {
		fo.Open = true
		fo.openContentTypes = contentTypes
	}
}
