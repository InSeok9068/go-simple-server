package views

import (
	"simple-server/pkg/util/gomutil"
	shared "simple-server/shared/views"

	x "github.com/glsubri/gomponents-alpine"
	h "maragu.dev/gomponents-htmx"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func Index(title string, groups map[string][]ClosetItem) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "ko",
		Head: gomutil.MergeHeads(
			shared.HeadsWithBeer(title),
			shared.HeadWithFirebaseAuth(),
			[]Node{
				Meta(Name("description"), Content("옷장 서비스")),
				Link(Rel("manifest"), Href("/manifest.json")),
				Script(Src("/static/closet.js"), Defer()),
			},
		),
		Body: []Node{
			shared.Snackbar(),

			Header(Class("primary"),
				Nav(
					A(Class("max"), Href("/"),
						H3(Text("Closet")),
					),
					A(Class("a-login"), ID("login"), Href("/login"), x.Data(""),
						x.Show("!$store.auth.isAuthed"),
						Text("로그인"),
					),
					A(Class("a-login"), ID("logout"), Href("#"), x.Data(""),
						x.Show("$store.auth.isAuthed"),
						Attr("onclick", "logoutUser()"),
						Text("로그아웃"),
					),
				),
			),

			Main(ID("closet-main"), Class("responsive"),
				Section(Class("page stack gap-lg"),
					Div(Class("grid gap-lg"),
						Div(Class("s12 m6 l4"), uploadCard()),
						Div(Class("s12 m6 l8"), filterCard()),
					),
					ItemsSection(groups),
				),
			),

			Nav(Class("bottom")),
		},
	})
}

func uploadCard() Node {
	return Article(Class("card closet-upload"), x.Data(""),
		Div(Class("padding stack gap-sm"),
			H3(Class("title"), Text("내 옷장에 추가")),
			P(Class("caption"), Text("타입과 태그를 입력하면 검색과 추천에 활용돼요.")),
			Div(x.Show("$store.auth.isAuthed"),
				Form(Class("stack gap-sm"),
					h.Post("/items"),
					h.Target("#items-list"),
					h.Swap("outerHTML"),
					Attr("hx-encoding", "multipart/form-data"),
					EncType("multipart/form-data"),
					h.On("htmx:after-request", "if (event.detail.successful) { this.reset(); showInfo('옷장에 새 옷을 추가했어요. 임베딩을 준비 중입니다.'); }"),
					h.On("htmx:response-error", "showError(event.detail.xhr.responseText || '업로드에 실패했어요. 다시 시도해주세요.');"),
					Div(Class("field label border"),
						Select(Name("kind"), Required(),
							Option(Value(""), Text("종류 선택"), Selected()),
							Option(Value("top"), Text("상의")),
							Option(Value("bottom"), Text("하의")),
							Option(Value("shoes"), Text("신발")),
							Option(Value("accessory"), Text("악세사리")),
						),
						Label(Text("종류")),
					),
					Div(Class("field label border"),
						Input(Name("tags"), Placeholder("예: 캐주얼, 스트릿")),
						Label(Text("태그")),
						Small(Class("caption"), Text("콤마(,) 또는 줄바꿈으로 여러 태그를 입력할 수 있어요.")),
					),
					Div(Class("field label border"),
						Input(Type("file"), Name("image"), Accept("image/*"), Required()),
						Label(Text("이미지")),
					),
					Div(Class("right-align"),
						Button(Class("button primary"), Text("업로드")),
					),
				),
			),
			Div(x.Show("!$store.auth.isAuthed"), Class("stack gap-sm"),
				P(Class("caption"), Text("로그인하면 나만의 옷장을 만들 수 있어요.")),
				A(Class("button outline"), Href("/login"), Text("로그인하기")),
			),
		),
	)
}

func filterCard() Node {
	return Article(Class("card closet-filter"),
		Div(Class("padding stack gap-sm"),
			H3(Class("title"), Text("옷장 찾기")),
			P(Class("caption"), Text("원하는 태그와 종류로 빠르게 찾아보세요.")),
			Form(Class("grid gap-sm"),
				h.Get("/items"),
				h.Target("#items-list"),
				h.Swap("outerHTML"),
				h.Trigger("change delay:300ms, keyup changed delay:500ms, submit"),
				h.On("htmx:response-error", "showError(event.detail.xhr.responseText || '검색에 실패했어요.');"),
				Div(Class("s12 m6"),
					Div(Class("field label border"),
						Select(Name("kind"),
							Option(Value(""), Text("전체"), Selected()),
							Option(Value("top"), Text("상의")),
							Option(Value("bottom"), Text("하의")),
							Option(Value("shoes"), Text("신발")),
							Option(Value("accessory"), Text("악세사리")),
						),
						Label(Text("종류")),
					),
				),
				Div(Class("s12"),
					Div(Class("field label border"),
						Input(Name("tags"), Placeholder("예: 캐주얼, 스트릿")),
						Label(Text("태그")),
					),
				),
				Div(Class("s12 right-align"),
					Nav(Class("chip"),
						Button(Class("button tonal"), Type("submit"), Text("검색")),
						Button(Class("button ghost"), Type("button"),
							Attr("onclick", "const form=this.form; form.reset(); form.dispatchEvent(new Event('submit',{bubbles:true}));"),
							Text("초기화"),
						),
					),
				),
				Div(Class("s12"),
					Img(Class("htmx-indicator closet-indicator"), Src("/shared/static/spinner.svg"), Alt("로딩 중")),
				),
			),
		),
	)
}
