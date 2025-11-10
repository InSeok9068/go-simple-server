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
				Meta(Name("description"), Content("개인 옷장 관리 서비스")),
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
					A(ID("login"), Href("/login"), x.Data(""),
						x.Show("!$store.auth.isAuthed"),
						Text("로그인"),
					),
					A(ID("logout"), Href("#"), x.Data(""),
						x.Show("$store.auth.isAuthed"),
						Attr("onclick", "logoutUser()"),
						Text("로그아웃"),
					),
				),
			),

			Main(ID("closet-main"),
				Section(
					Div(Class("row"),
						Button(Class("button"), Attr("data-ui", "#upload-dialog"), Text("옷장에 추가")),
						Button(Class("button"), Attr("data-ui", "#search-dialog"), Text("옷장에서 찾기")),
					),
					recommendCard(),
					Br(),
					ItemsSection(groups),
				),
			),

			uploadDialog(),
			searchDialog(),
			RecommendDialogContainer(),
		},
	})
}

func uploadDialog() Node {
	return Dialog(Class("top"), ID("upload-dialog"), x.Data(""),
		Div(
			H5(Text("내 옷장에 추가")),
			P(Class("caption"), Text("이미지와 태그를 입력하면 검색과 추천에 활용돼요.")),
			Div(x.Show("$store.auth.isAuthed"),
				Form(Class("upload-form"), x.Data("{ mode: 'file', fileName: '' }"),
					h.Post("/items"),
					h.Target("#items-list"),
					h.Swap("outerHTML"),
					Attr("hx-encoding", "multipart/form-data"),
					EncType("multipart/form-data"),
					Attr("hx-on::after-request", "if (event.detail.successful) { this.reset(); if (this.__x) { this.__x.$data.fileName = ''; } closeModal('#upload-dialog'); showInfo('새 옷을 추가했어요. 임베딩을 준비 중입니다.'); }"),
					h.On("htmx:response-error", "showError(event.detail.xhr.responseText || '업로드에 실패했어요. 다시 시도해주세요.');"),
					Div(Class("field label border"),
						Select(Name("kind"), Required(),
							Option(Value(""), Text("종류 선택"), Selected()),
							Option(Value("top"), Text("상의")),
							Option(Value("bottom"), Text("하의")),
							Option(Value("shoes"), Text("신발")),
							Option(Value("accessory"), Text("액세서리")),
						),
						Label(Text("종류")),
					),
					Div(
						Div(Class("field label border"),
							Input(Type("text"), Name("tags")),
							Label(Text("태그")),
							Span(Class("helper"), Text("콤마(,)로 여러 태그를 입력해 주세요.")),
						),
					),
					Div(
						Div(Class("padding"),
							Nav(
								Label(Class("radio"),
									Input(Type("radio"), Name("input_mode"), Value("file"), x.Model("mode"), Checked()),
									Span(Text("파일 업로드")),
								),
								Label(Class("radio"),
									Input(Type("radio"), Name("input_mode"), Value("url"), x.Model("mode"),
										x.On("change", "if ($event.target.checked && $refs.fileInput) { $refs.fileInput.value = ''; fileName = ''; }")),
									Span(Text("이미지 URL")),
								),
							),
						),
					),
					Div(x.Show("mode === 'file'"),
						Label(Class("chip primary"),
							I(Class("icon"), Text("attach_file")),
							Span(Text("이미지 선택")),
							Input(Type("file"), Name("image"), Accept("image/*"), x.Bind("required", "mode === 'file'"),
								x.Ref("fileInput"),
								x.On("change", "fileName = $event.target.files?.[0]?.name || ''")),
						),
						Div(Class("caption"), x.Show("fileName"), x.Text("'선택한 파일: ' + fileName")),
					),
					Div(x.Show("mode === 'url'"),
						Div(Class("field label"),
							Input(Type("url"), Name("image_url"), x.Bind("required", "mode === 'url'"),
								Style("border:none;border-bottom:2px solid #6c5ce7;border-radius:0;background:transparent;padding-left:0;")),
							Label(Text("이미지 URL")),
						),
						Small(Class("caption"), Text("URL도 20MB 이하만 허용돼요.")),
					),
					Div(Class("row"),
						Button(Class("button"), Text("업로드")),
						Button(Class("button"), Type("button"), Attr("data-ui", "#upload-dialog"), Text("닫기")),
					),
					Div(Class("upload-overlay htmx-indicator"),
						DataAttr("indicator-mode", "overlay"),
						Div(Class("upload-overlay__content"),
							I(Class("icon"), Text("hourglass_empty")),
							P(Class("upload-overlay__text"), Text("이미지를 분석 중이에요...")),
							Small(Class("caption"), Text("잠시만 기다려 주세요.")),
						),
					),
				),
			),
			Div(x.Show("!$store.auth.isAuthed"),
				P(Class("caption"), Text("로그인해야 옷장을 만들 수 있어요.")),
				A(Class("button"), Href("/login"), Text("로그인하기")),
			),
		),
	)
}

func searchDialog() Node {
	return Dialog(Class("top"), ID("search-dialog"),
		Div(Class("padding"),
			H5(Text("옷장 찾기")),
			P(Class("caption"), Text("원하는 태그와 종류로 빠르게 찾아보세요.")),
			Form(
				h.Get("/items"),
				h.Target("#items-list"),
				h.Swap("outerHTML"),
				h.Trigger("change delay:300ms, keyup changed delay:500ms, submit"),
				h.On("htmx:response-error", "showError(event.detail.xhr.responseText || '검색에 실패했어요.');"),
				Div(Class("field label border"),
					Select(Name("kind"),
						Option(Value(""), Text("전체"), Selected()),
						Option(Value("top"), Text("상의")),
						Option(Value("bottom"), Text("하의")),
						Option(Value("shoes"), Text("신발")),
						Option(Value("accessory"), Text("액세서리")),
					),
					Label(Text("종류")),
				),
				Div(Class("field label border"),
					Input(Type("text"), Name("tags")),
					Label(Text("태그")),
				),
				Div(Class("row"),
					Button(Class("button"), Type("submit"), Text("검색")),
					Button(Class("button"), Type("button"),
						Attr("onclick", "const form=this.form; form.reset(); closeModal('#search-dialog');"),
						Text("닫기"),
					),
				),
				Img(Class("htmx-indicator"), Src("/shared/static/spinner.svg"), Alt("로딩 중")),
			),
		),
	)
}

func recommendCard() Node {
	return Article(Class("padding"), x.Data(""),
		Div(
			Div(x.Show("$store.auth.isAuthed"),
				Form(
					h.Post("/recommend"),
					h.Target("#recommend-dialog-body"),
					h.Swap("innerHTML"),
					h.On("htmx:after-request", "if (event.detail.successful) { document.getElementById('recommend-dialog-trigger')?.click(); }"),
					h.On("htmx:response-error", "showError(event.detail.xhr.responseText || '추천에 실패했어요.');"),
					Div(Class("field label border"),
						Input(Type("text"), Name("weather")),
						Label(Text("날씨")),
					),
					Div(Class("field label border"),
						Input(Type("text"), Name("style")),
						Label(Text("스타일")),
					),
					Div(Class("row"),
						Button(Class("button"), Type("submit"), Text("추천받기")),
					),
				),
				Button(Type("button"), ID("recommend-dialog-trigger"), Attr("data-ui", "#recommend-dialog"), Style("display:none")),
			),
			Div(x.Show("!$store.auth.isAuthed"),
				P(Class("caption"), Text("로그인하면 AI 추천을 받아볼 수 있어요.")),
				A(Class("button"), Href("/login"), Text("로그인하기")),
			),
		),
	)
}
