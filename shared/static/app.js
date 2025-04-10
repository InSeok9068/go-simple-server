document.addEventListener('alpine:init', () => {
    Alpine.store('auth', {
        isAuthed: false,
        user: null,

        login(user) {
            this.isAuthed = true;
            this.user = user;
        },
        logout() {
            this.isAuthed = false;
            this.user = null;
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

function showModal(querySelector) {
    document.querySelector(querySelector).showModal()
}

function closeModal(querySelector) {
    document.querySelector(querySelector).close()
}