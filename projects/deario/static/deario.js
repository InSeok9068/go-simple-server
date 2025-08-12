console.log("Deario에 오신걸 환영합니다.");

document.addEventListener("alpine:init", () => {
  Alpine.store("save", {
    isOk: true,

    ok() {
      this.isOk = true;
    },

    unok() {
      this.isOk = false;
    },
  });
});

document.addEventListener("htmx:afterSwap", (e) => {
  if (e.detail.target.id === "ai-feedback-content") {
    const mdEl = document.getElementById("ai-feedback-markdown");
    if (mdEl) {
      mdEl.innerHTML = marked.parse(mdEl.textContent);
    }
  }
});
