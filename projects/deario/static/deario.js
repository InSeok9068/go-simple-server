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

window.resetPin = async function () {
  try {
    const token = await window.getIdToken();
    const res = await fetch("/pin/reset", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        "X-CSRF-Token": getCookie("_csrf"),
      },
      body: JSON.stringify({ token }),
    });
    if (res.ok) {
      showInfo("핀번호가 초기화되었습니다.");
      location.reload();
    } else {
      showError("핀번호 초기화 실패");
    }
  } catch (e) {
    showError("재인증 실패");
  }
};
