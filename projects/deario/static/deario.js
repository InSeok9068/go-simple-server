;(function () {
  const DATA_DOT_CLASS = "diary-data-dot"
  const AI_FEEDBACK_DIALOG = "#ai-feedback-dialog"
  const EMPTY_AI_FEEDBACK_HTML =
    '<div id="ai-feedback-markdown"></div><textarea name="ai-feedback" hidden></textarea>'

  const htmxAfterHandlers = {
    "ai-feedback:show": showAiFeedback,
    "ai-feedback:saved": onAiFeedbackSaved,
    "ai-feedback:deleted": onAiFeedbackDeleted,
    toast: showToastFromElement,
  }

  document.addEventListener("alpine:init", registerSaveStore)
  document.addEventListener("click", handleLogoutClick)
  document.addEventListener("input", handleDiaryInput)
  document.addEventListener("htmx:afterRequest", handleDiarySaveRequest)
  document.addEventListener("htmx:afterOnLoad", handleHtmxAfterOnLoad)
  document.addEventListener("htmx:afterSwap", handleHtmxAfterSwap)
  document.addEventListener("DOMContentLoaded", initSwipeNavigation)

  function registerSaveStore() {
    Alpine.store("save", {
      isOk: true,

      setSaved(isSaved) {
        this.isOk = isSaved
      },

      ok() {
        this.setSaved(true)
      },

      unok() {
        this.setSaved(false)
      },
    })
  }

  function handleLogoutClick(event) {
    const trigger = event.target.closest?.("[data-deario-logout]")
    if (!trigger) return

    event.preventDefault()
    window.logoutUser?.()
  }

  function handleDiaryInput(event) {
    if (event.target.matches?.("[data-deario-diary-input]")) {
      setDiarySaved(false)
    }
  }

  function handleDiarySaveRequest(event) {
    const source = htmxSourceElement(event)
    if (
      source instanceof Element &&
      source.matches("[data-deario-diary-input]")
    ) {
      setDiarySaved(true)
    }
  }

  function setDiarySaved(isSaved) {
    if (!window.Alpine) return

    const store = Alpine.store("save")
    if (store) {
      store.setSaved(isSaved)
    }
  }

  function handleHtmxAfterOnLoad(event) {
    if (!isHtmxSuccessful(event)) return

    const actionEl = htmxActionElement(event, "dearioAfter")
    const action = actionEl?.dataset.dearioAfter
    const handler = htmxAfterHandlers[action]

    if (handler) {
      handler(actionEl, event)
    }
  }

  function handleHtmxAfterSwap(event) {
    const target = event.detail.target
    if (
      target?.id === "diary-image-content" ||
      target?.querySelector?.("#diary-image-content")
    ) {
      syncDiaryImageState()
    }
  }

  function isHtmxSuccessful(event) {
    if (typeof event.detail.successful === "boolean") {
      return event.detail.successful
    }

    const status = event.detail.xhr?.status || 0
    return status >= 200 && status < 400
  }

  function htmxActionElement(event, datasetKey) {
    const elt = htmxSourceElement(event)
    if (!(elt instanceof Element)) return null

    if (elt.dataset?.[datasetKey]) {
      return elt
    }

    return elt.closest(`[data-${kebabCase(datasetKey)}]`)
  }

  function htmxSourceElement(event) {
    return event.detail.elt || event.detail.requestConfig?.elt || event.target
  }

  function kebabCase(value) {
    return value.replace(/[A-Z]/g, (ch) => `-${ch.toLowerCase()}`)
  }

  function showToastFromElement(el) {
    showInfo(el.dataset.dearioMessage || "저장되었습니다.")
  }

  function showAiFeedback() {
    const mdEl = document.getElementById("ai-feedback-markdown")
    if (mdEl && mdEl.textContent && window.marked) {
      mdEl.innerHTML = marked.parse(mdEl.textContent)
    }
    showModal(AI_FEEDBACK_DIALOG)
  }

  function onAiFeedbackSaved() {
    setAiFeedbackDataState(true)
    closeModal(AI_FEEDBACK_DIALOG)
    showInfo("저장 되었습니다.")
  }

  function onAiFeedbackDeleted() {
    setAiFeedbackDataState(false)
    resetAiFeedbackContent()
    closeModal(AI_FEEDBACK_DIALOG)
    showInfo("삭제되었습니다.")
  }

  function setAiFeedbackDataState(hasData) {
    setDiaryDataDot("#ai-feedback-action", hasData)

    const deleteAction = document.getElementById("ai-feedback-delete-action")
    if (deleteAction) {
      deleteAction.hidden = !hasData
    }
  }

  function resetAiFeedbackContent() {
    const content = document.getElementById("ai-feedback-content")
    if (!content) return

    content.innerHTML = EMPTY_AI_FEEDBACK_HTML
  }

  function syncDiaryImageState() {
    const content = document.getElementById("diary-image-content")
    if (!content) return

    setDiaryDataDot(
      "#diary-image-action",
      Boolean(content.querySelector("img")),
    )
  }

  function setDiaryDataDot(selector, hasData) {
    const action = document.querySelector(selector)
    if (!action) return

    action.classList.toggle(DATA_DOT_CLASS, hasData)
  }

  function initSwipeNavigation() {
    const mainEl = document.getElementById("diary-main")
    if (!mainEl || !window.Hammer) return

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
  }
})()
