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