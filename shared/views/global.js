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