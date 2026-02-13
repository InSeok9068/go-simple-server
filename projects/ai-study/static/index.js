document.addEventListener("DOMContentLoaded", () => {
  const copyButton = document.getElementById("copy")
  if (!copyButton) {
    return
  }

  copyButton.addEventListener("click", async () => {
    const result = document.getElementById("result")
    if (!result) {
      return
    }

    try {
      await navigator.clipboard.writeText(result.innerText || "")
      if (typeof showInfo === "function") {
        showInfo("복사 되었습니다.")
      }
    } catch (_error) {
      if (typeof showError === "function") {
        showError("복사에 실패했습니다.")
      }
    }
  })
})
