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
