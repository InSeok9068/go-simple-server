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

function showAiFeedback() {
  const mdEl = document.getElementById("ai-feedback-markdown");
  if (mdEl && mdEl.textContent) {
    mdEl.innerHTML = marked.parse(mdEl.textContent);
  }
  showModal("#ai-feedback-dialog");
}

document.addEventListener("DOMContentLoaded", () => {
  const mainEl = document.getElementById("diary-main");
  if (!mainEl) return;

  const prev = document.getElementById("prev-day");
  const next = document.getElementById("next-day");

  const hammer = new Hammer(mainEl);
  hammer
    .get("swipe")
    .set({ direction: Hammer.DIRECTION_HORIZONTAL, threshold: 60 });

  hammer.on("swipeleft", () => {
    if (next) location.href = next.href;
  });

  hammer.on("swiperight", () => {
    if (prev) location.href = prev.href;
  });
});
