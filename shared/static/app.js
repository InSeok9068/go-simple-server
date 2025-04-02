document.addEventListener('alpine:init', () => {
    Alpine.store('auth', {
        isAuthed: false,
        user: null,
        token: null,

        login(user, token) {
            this.isAuthed = true;
            this.user = user;
            this.token = token;
        },
        logout() {
            this.isAuthed = false;
            this.user = null;
            this.token = null;
        }
    });
});

htmx.on("htmx:afterRequest", (event) => {
    const contentType = event.detail.xhr.getResponseHeader("Content-Type");
    if (contentType !== 'application/json') {
        return;
    }

    const responseData = event.detail.xhr.responseText;
    if (responseData === '') {
        return;
    }

    const isResponseError = event.detail.xhr.status >= 400;
    if (isResponseError) {
        const parsedResponse = JSON.parse(responseData);
        if (parsedResponse.message === undefined || parsedResponse.message === '') {
            return;
        }
        alert(parsedResponse.message);
    }
});

function blockBackAndExit() {
    history.pushState(null, '', location.href); // 현재 상태 복사
    window.addEventListener('popstate', () => {
        // 사용자가 뒤로가기를 시도하면 강제로 종료 페이지로 보냄
        location.replace('about:blank');
    });

    // 앱 종료 느낌 주기
    document.body.innerHTML = '';
    window.close()
    location.replace('about:blank');
}