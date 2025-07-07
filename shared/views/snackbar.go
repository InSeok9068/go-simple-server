package shared

import (
	x "github.com/glsubri/gomponents-alpine"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Snackbar() Node {
	return Div(
		x.Data("$store.snackbar"),
		x.Show("visible"),
		x.Class("`${type} ${visible ? 'active' : ''}`"),
		Class("snackbar"),
		Span(x.Text("message")),
	)
}
