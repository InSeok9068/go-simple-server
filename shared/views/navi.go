package shared

import (
	x "github.com/glsubri/gomponents-alpine"
	b "github.com/willoma/bulma-gomponents"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func NaviShare(title string) Node {
	return b.Navbar(
		b.Transparent,
		x.Data("{ open : false }"),
		b.NavbarBrand(
			b.NavbarAHref("/"),
		),
		b.NavbarEnd(
			b.NavbarItem(
				P(ID("username")),
			),
			b.NavbarItem(
				b.Buttons(
					b.ButtonA(Href("/login"), b.Hidden, Text("Login")),
					b.ButtonA(Href("/logout"), b.Hidden, Text("Logout")),
				),
			),
		),
	)
}

/*
package shared

templ Navi(title string, menus templ.Component) {
	<nav class="navbar is-transparent" x-data="{ open: false }">
		<div class="navbar-brand">
			<a class="navbar-item" href="/">
				<h1 class="title">{ title }</h1>
			</a>
			<div class="navbar-burger js-burger" :class="open ? 'is-active' : ''" @click="open = !open">
				<span></span>
				<span></span>
				<span></span>
				<span></span>
			</div>
		</div>
		<div class="navbar-menu" :class="open ? 'is-active' : ''">
			<div class="navbar-end">
				<div class="navbar-item">
					<p id="username"></p>
				</div>
				<div class="navbar-item">
					<div class="buttons">
						<a id="login" class="button is-hidden" href="/login">Login</a>
						<a id="logout" class="button is-hidden">Logout</a>
					</div>
				</div>
				if menus != nil {
					@menus
				}
			</div>
		</div>
	</nav>
}

*/
