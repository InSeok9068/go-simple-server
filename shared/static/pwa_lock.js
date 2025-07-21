document.addEventListener("DOMContentLoaded", () => {
  const storedPin = localStorage.getItem("pwaPin");
  const unlocked = sessionStorage.getItem("pwaUnlocked") === "true";
  if (storedPin && unlocked) {
    return;
  }

  const overlay = document.createElement("div");
  overlay.id = "pwa-lock-overlay";
  overlay.style.cssText =
    "position:fixed;inset:0;z-index:9999;background:white;display:flex;flex-direction:column;justify-content:center;align-items:center;";

  const msg = storedPin ? "PIN 입력" : "새 PIN 설정";
  overlay.innerHTML = `
    <div style="display:flex;flex-direction:column;gap:0.5rem;">
      <input id="pwa-lock-pin" type="password" placeholder="${msg}" style="text-align:center;" />
      <button id="pwa-lock-btn">${storedPin ? "확인" : "저장"}</button>
      <p id="pwa-lock-msg" style="color:red;display:none;"></p>
    </div>
  `;
  document.body.appendChild(overlay);

  const input = overlay.querySelector("#pwa-lock-pin");
  const btn = overlay.querySelector("#pwa-lock-btn");
  const msgEl = overlay.querySelector("#pwa-lock-msg");

  function showError(text) {
    msgEl.textContent = text;
    msgEl.style.display = "block";
  }

  btn.addEventListener("click", () => {
    const val = input.value.trim();
    if (!storedPin) {
      if (val.length < 4) {
        showError("PIN은 4자리 이상 입력해주세요");
        return;
      }
      localStorage.setItem("pwaPin", val);
      sessionStorage.setItem("pwaUnlocked", "true");
      overlay.remove();
      if (window.showInfo) {
        showInfo("PIN이 저장되었습니다.");
      }
    } else {
      if (val === storedPin) {
        sessionStorage.setItem("pwaUnlocked", "true");
        overlay.remove();
      } else {
        showError("PIN이 일치하지 않습니다");
      }
    }
  });

  input.addEventListener("keydown", (e) => {
    if (e.key === "Enter") {
      btn.click();
    }
  });
});
