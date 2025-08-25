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

document.addEventListener("DOMContentLoaded", () => {
  const textarea = document.querySelector("textarea[name='content']");
  if (!textarea) return;
  const editor = pell.init({
    element: document.getElementById("editor"),
    onChange: (html) => {
      textarea.value = html;
      textarea.dispatchEvent(new Event("input", { bubbles: true }));
    },
  });
  editor.content.innerHTML = textarea.value;
  editor.content.setAttribute(
    "data-placeholder",
    textarea.getAttribute("placeholder"),
  );
  editor.content.focus();
});

function showAiFeedback() {
  const mdEl = document.getElementById("ai-feedback-markdown");
  if (mdEl && mdEl.textContent) {
    mdEl.innerHTML = marked.parse(mdEl.textContent);
  }
  showModal("#ai-feedback-dialog");
}
