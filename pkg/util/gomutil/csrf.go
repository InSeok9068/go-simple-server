package gomutil

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// CSRFInput returns a hidden input element with the given CSRF token.
func CSRFInput(token string) Node {
	return Input(Type("hidden"), Name("csrf"), Value(token))
}
