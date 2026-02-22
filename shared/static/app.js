document.addEventListener("alpine:init", () => {
  Alpine.store("auth", {
    isAuthed: false,
    user: null,

    login(user) {
      this.isAuthed = true
      this.user = user
      try {
        localStorage.setItem("authUser", JSON.stringify(user))
      } catch (e) {
        console.error("localStorage save error", e)
      }
    },
    logout() {
      this.isAuthed = false
      this.user = null
      localStorage.removeItem("authUser")
    },
  })

  const saved = localStorage.getItem("authUser")
  if (saved) {
    try {
      Alpine.store("auth").login(JSON.parse(saved))
    } catch (e) {
      console.error("localStorage parse error", e)
      localStorage.removeItem("authUser")
    }
  }

  Alpine.store("notification", {
    permission: Notification.permission === "granted",
  })

  Alpine.store("snackbar", {
    visible: false,
    message: "",
    type: "primary",
    queue: [],
    isShowing: false,
    hideTimer: null,
    show(msg, type = "primary", ms = 3000) {
      const duration =
        typeof ms === "number" && ms >= 0 ? Math.floor(ms) : 3000

      return new Promise((resolve) => {
        this.queue.push({
          message: msg,
          type,
          duration,
          resolve,
        })
        this.processQueue()
      })
    },
    async processQueue() {
      if (this.isShowing) {
        return
      }
      this.isShowing = true

      while (this.queue.length > 0) {
        const next = this.queue.shift()

        this.message = next.message
        this.type = next.type
        this.visible = true

        await new Promise((resolve) => {
          this.hideTimer = setTimeout(resolve, next.duration)
        })

        this.visible = false
        this.hideTimer = null
        next.resolve()

        await new Promise((resolve) => {
          setTimeout(resolve, 120)
        })
      }

      this.isShowing = false
    },
    info(msg, ms) {
      return this.show(msg, "primary", ms)
    },
    error(msg, ms) {
      return this.show(msg, "error", ms)
    },
  })

  Alpine.store("font", {
    value: "gamja",
    init() {
      const saved = localStorage.getItem("font")
      if (saved) {
        this.value = saved
      }
      this.apply()
    },
    set(font) {
      this.value = font
      this.apply()
      try {
        localStorage.setItem("font", font)
      } catch (e) {
        console.error("localStorage save error", e)
      }
    },
    apply() {
      const fonts = {
        gamja: '"Gamja Flower"',
        humanist: "var(--font-humanist)",
        neo_grotesque: "var(--font-neo-grotesque)",
        monospace_code: "var(--font-monospace-code)",
        industrial: "var(--font-industrial)",
        rounded_sans: "var(--font-rounded-sans)",
      }
      document.documentElement.style.setProperty(
        "--font-family",
        fonts[this.value] || "Gamja Flower",
      )
    },
  })

  Alpine.store("theme", {
    value: "light",
    color: "#6200ee",
    init() {
      const saved = localStorage.getItem("theme")
      if (saved) {
        this.value = saved
      }
      const savedColor = localStorage.getItem("themeColor")
      if (savedColor) {
        this.color = savedColor
      }
      window.ui("mode", this.value)
      window.ui("theme", this.color)
      this.applyFlatpickrTheme(this.value)
    },
    set(theme) {
      this.value = theme
      window.ui("mode", theme)
      try {
        localStorage.setItem("theme", theme)
      } catch (e) {
        console.error("localStorage save error", e)
      }
      this.applyFlatpickrTheme(theme)
    },
    setColor(color) {
      this.color = color
      window.ui("theme", color)
      try {
        localStorage.setItem("themeColor", color)
      } catch (e) {
        console.error("set color error", e)
      }
    },
    toggle() {
      this.set(this.value === "dark" ? "light" : "dark")
    },
    applyFlatpickrTheme(theme) {
      const id = "flatpickr-dark"
      const href =
        "https://cdn.jsdelivr.net/npm/flatpickr@4.6.13/dist/themes/dark.css"
      const onError =
        "this.onerror=null;this.href='shared/static/lib/dark.css';"
      const link = document.getElementById(id)
      if (theme === "dark") {
        if (!link) {
          const tag = document.createElement("link")
          tag.id = id
          tag.rel = "stylesheet"
          tag.href = href
          tag.onerror = onError
          document.head.appendChild(tag)
        }
      } else if (link) {
        link.remove()
      }
    },
  })

  Alpine.store("theme").init()
  Alpine.store("font").init()
})

htmx.on("htmx:afterRequest", (event) => {
  const contentType = (
    event.detail.xhr.getResponseHeader("Content-Type") || ""
  ).toLowerCase()
  if (!contentType.includes("application/json")) {
    return
  }

  const responseData = event.detail.xhr.responseText
  if (responseData === "") {
    return
  }

  const isResponseError = event.detail.xhr.status >= 400
  if (isResponseError) {
    if (
      event.detail.xhr.status === 401 &&
      typeof window.isUnauthorizedHandled === "function" &&
      window.isUnauthorizedHandled()
    ) {
      return
    }

    const parsedResponse = JSON.parse(responseData)
    if (parsedResponse.message === undefined || parsedResponse.message === "") {
      return
    }
    showError(parsedResponse.message)
  }
})

htmx.on("htmx:configRequest", (event) => {
  const match = getCookie("_csrf")
  if (match) {
    event.detail.headers["X-CSRF-Token"] = match
  }
})

function showInfo(msg, ms) {
  return Alpine.store("snackbar").info(msg, ms)
}

function showError(msg, ms) {
  return Alpine.store("snackbar").error(msg, ms)
}

function showModal(querySelector, isScrollTop = false) {
  const modal = document.querySelector(querySelector)
  modal?.showModal()

  if (isScrollTop) {
    modal.scrollTo({ top: 0, behavior: "auto" })
  }
}

function closeModal(querySelector, isScrollTop = false) {
  const modal = document.querySelector(querySelector)
  modal?.close()

  if (isScrollTop) {
    modal.scrollTo({ top: 0, behavior: "auto" })
  }
}

function getCookie(name) {
  return (
    document.cookie
      .split("; ")
      .find((v) => v.startsWith(name + "="))
      ?.split("=")[1] || ""
  )
}
