;(function () {
  const LOCK_AFTER_MS = 5000
  const HIDDEN_AT_KEY = "deario_app_hidden_at"
  const LOCK_REQUESTED_KEY = "deario_app_lock_requested"
  const TEMPORARY_SUPPRESS_MS = 60000

  let suppressLockUntil = 0
  let lockTimer = null
  let lockInFlight = false

  function now() {
    return Date.now()
  }

  function csrfToken() {
    return (
      document.cookie
        .split("; ")
        .find((value) => value.startsWith("_csrf="))
        ?.split("=")[1] || ""
    )
  }

  function suppressLockTemporarily() {
    suppressLockUntil = now() + TEMPORARY_SUPPRESS_MS
  }

  function shouldSkipLock() {
    const path = location.pathname
    return (
      path === "/login" ||
      path === "/app-lock" ||
      now() < suppressLockUntil
    )
  }

  function applyPrivacyMask() {
    document.body?.classList.add("deario-privacy-mask")
  }

  function clearPrivacyMask() {
    document.body?.classList.remove("deario-privacy-mask")
  }

  async function requestLock(redirectOnLock) {
    if (lockInFlight) {
      if (redirectOnLock) {
        setTimeout(() => requestLock(true), 100)
      }
      return
    }

    lockInFlight = true
    sessionStorage.setItem(LOCK_REQUESTED_KEY, "1")

    try {
      const response = await fetch("/app-lock/lock", {
        method: "POST",
        headers: { "X-CSRF-Token": csrfToken() },
        credentials: "same-origin",
        keepalive: true,
      })
      const redirectTo = response.headers.get("Hx-Redirect")

      if (redirectOnLock && response.status === 401) {
        location.href = "/login"
        return
      }

      if (redirectOnLock && redirectTo) {
        location.href = redirectTo
        return
      }

      if (redirectOnLock) {
        sessionStorage.removeItem(LOCK_REQUESTED_KEY)
        clearPrivacyMask()
      }
    } catch {
      if (redirectOnLock) {
        location.href = "/app-lock"
      }
    } finally {
      lockInFlight = false
    }
  }

  function scheduleLock() {
    if (shouldSkipLock()) {
      return
    }

    sessionStorage.setItem(HIDDEN_AT_KEY, String(now()))
    applyPrivacyMask()
    clearTimeout(lockTimer)

    lockTimer = setTimeout(() => {
      if (document.visibilityState === "hidden" && !shouldSkipLock()) {
        requestLock(false)
      }
    }, LOCK_AFTER_MS)
  }

  function lockIfNeededOnReturn() {
    clearTimeout(lockTimer)

    const hiddenAt = Number(sessionStorage.getItem(HIDDEN_AT_KEY) || "0")
    const alreadyRequested =
      sessionStorage.getItem(LOCK_REQUESTED_KEY) === "1"

    sessionStorage.removeItem(HIDDEN_AT_KEY)

    if (shouldSkipLock()) {
      clearPrivacyMask()
      return
    }

    if (alreadyRequested || (hiddenAt > 0 && now() - hiddenAt >= LOCK_AFTER_MS)) {
      applyPrivacyMask()
      requestLock(true)
      return
    }

    clearPrivacyMask()
  }

  document.addEventListener(
    "click",
    (event) => {
      const target = event.target
      if (!(target instanceof Element)) {
        return
      }

      if (target.closest("input[type='file'], .microphone")) {
        suppressLockTemporarily()
        return
      }

    },
    true,
  )

  document.addEventListener("visibilitychange", () => {
    if (document.visibilityState === "hidden") {
      scheduleLock()
      return
    }

    lockIfNeededOnReturn()
  })

  window.addEventListener("pagehide", () => {
    if (shouldSkipLock()) {
      return
    }
    sessionStorage.setItem(HIDDEN_AT_KEY, String(now()))
    applyPrivacyMask()
  })

  window.addEventListener("pageshow", () => {
    if (document.visibilityState === "visible") {
      lockIfNeededOnReturn()
    }
  })
})()
