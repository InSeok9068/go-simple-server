console.log("Closet에 오신걸 환영합니다.");

document.addEventListener("htmx:responseError", (event) => {
  const message =
    event.detail?.xhr?.responseText || "요청 처리 중 오류가 발생했어요.";
  if (typeof showError === "function") {
    showError(message);
  } else {
    console.error(message);
  }
});
