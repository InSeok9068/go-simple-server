console.log("Deario에 오신걸 환영합니다.")

document.addEventListener("alpine:init", () => {
  Alpine.store("save", {
    isOk: true,

    ok() {
      this.isOk = true
    },

    unok() {
      this.isOk = false
    },
  })
})

/**
 * AI 피드백 마크다운으로 랜더링 후 다이얼로그로 띄웁니다.
 */
function showAiFeedback() {
  const mdEl = document.getElementById("ai-feedback-markdown")
  if (mdEl && mdEl.textContent) {
    mdEl.innerHTML = marked.parse(mdEl.textContent)
  }
  showModal("#ai-feedback-dialog")
}

function setDiaryDataDot(selector, hasData) {
  const action = document.querySelector(selector)
  if (!action) return

  action.classList.toggle("diary-data-dot", hasData)
}

function setAiFeedbackDataState(hasData) {
  setDiaryDataDot("#ai-feedback-action", hasData)

  const deleteAction = document.getElementById("ai-feedback-delete-action")
  if (deleteAction) {
    deleteAction.hidden = !hasData
  }
}

function onAiFeedbackSaved() {
  setAiFeedbackDataState(true)
  closeModal("#ai-feedback-dialog")
  showInfo("저장 되었습니다.")
}

function onAiFeedbackDeleted() {
  setAiFeedbackDataState(false)
  resetAiFeedbackContent()
  closeModal("#ai-feedback-dialog")
  showInfo("삭제되었습니다.")
}

function resetAiFeedbackContent() {
  const content = document.getElementById("ai-feedback-content")
  if (!content) return

  content.innerHTML =
    '<div id="ai-feedback-markdown"></div><textarea name="ai-feedback" hidden></textarea>'
}

function syncDiaryImageState() {
  const content = document.getElementById("diary-image-content")
  if (!content) return

  setDiaryDataDot("#diary-image-action", Boolean(content.querySelector("img")))
}

document.addEventListener("htmx:afterSwap", (event) => {
  const target = event.detail.target
  if (
    target?.id === "diary-image-content" ||
    target?.querySelector?.("#diary-image-content")
  ) {
    syncDiaryImageState()
  }
})

/**
 * 이전/다음 페이지 스와이프 핸들러를 등록합니다.
 */
document.addEventListener("DOMContentLoaded", () => {
  const mainEl = document.getElementById("diary-main")
  if (!mainEl) return

  const prev = document.getElementById("prev-day")
  const next = document.getElementById("next-day")

  const hammer = new Hammer(mainEl)
  hammer
    .get("swipe")
    .set({ direction: Hammer.DIRECTION_HORIZONTAL, threshold: 60 })

  hammer.on("swipeleft", () => {
    if (next) location.href = next.href
  })

  hammer.on("swiperight", () => {
    if (prev) location.href = prev.href
  })
})
