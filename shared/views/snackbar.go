package shared

import (
	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Snackbar returns a node representing a snackbar, which is a small notification that
// appears at the bottom of the screen. It is controlled by the $store.snackbar state.
//
// It will be rendered with a class of "primary", "secondary", "success", "warning", or
// "danger" depending on the type of the notification.
//
// The snackbar will be shown or hidden based on the "visible" property of the
// $store.snackbar state. When visible is true, the snackbar will be rendered with an
// "active" class.
//
// The content of the snackbar is set by the "message" property of the
// $store.snackbar state.
func Snackbar() Node {
	return Div(
		x.Data("$store.snackbar"),
		x.Show("visible"),
		x.Class("`${type} ${visible ? 'active' : ''}`"),
		Class("snackbar"),
		Span(x.Text("message")),
	)
}
