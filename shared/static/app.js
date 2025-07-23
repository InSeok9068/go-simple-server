document.addEventListener("alpine:init", () => {
  Alpine.store("auth", {
    isAuthed: false,
    user: null,

    login(user) {
      this.isAuthed = true;
      this.user = user;
      try {
        localStorage.setItem("authUser", JSON.stringify(user));
      } catch (e) {
        console.error("localStorage save error", e);
      }
    },
    logout() {
      this.isAuthed = false;
      this.user = null;
      localStorage.removeItem("authUser");
    },
  });

  const saved = localStorage.getItem("authUser");
  if (saved) {
    try {
      Alpine.store("auth").login(JSON.parse(saved));
    } catch (e) {
      console.error("localStorage parse error", e);
      localStorage.removeItem("authUser");
    }
  }

  Alpine.store("notification", {
    permission: Notification.permission === "granted",
  });

  Alpine.store("snackbar", {
    visible: false,
    message: "",
    type: "primary",
    show(msg, type = "primary", ms = 3000) {
      this.message = msg;
      this.type = type;
      this.visible = true;
      setTimeout(() => {
        this.visible = false;
      }, ms);
    },
    info(msg, ms) {
      this.show(msg, "primary", ms);
    },
    error(msg, ms) {
      this.show(msg, "error", ms);
    },
  });

  Alpine.store("theme", {
    value: "light",
    init() {
      const saved = localStorage.getItem("theme");
      if (saved) {
        this.value = saved;
      }
      document.documentElement.setAttribute("data-theme", this.value);
    },
    set(theme) {
      this.value = theme;
      document.documentElement.setAttribute("data-theme", theme);
      try {
        localStorage.setItem("theme", theme);
      } catch (e) {
        console.error("localStorage save error", e);
      }
    },
    toggle() {
      this.set(this.value === "dark" ? "light" : "dark");
    },
  });

  Alpine.store("theme").init();
});

htmx.on("htmx:afterRequest", (event) => {
  const contentType = event.detail.xhr.getResponseHeader("Content-Type");
  if (contentType !== "application/json") {
    return;
  }

  const responseData = event.detail.xhr.responseText;
  if (responseData === "") {
    return;
  }

  const isResponseError = event.detail.xhr.status >= 400;
  if (isResponseError) {
    const parsedResponse = JSON.parse(responseData);
    if (parsedResponse.message === undefined || parsedResponse.message === "") {
      return;
    }
    showError(parsedResponse.message);
  }
});

htmx.on("htmx:configRequest", (event) => {
  const match = getCookie("_csrf");
  if (match) {
    event.detail.headers["X-CSRF-Token"] = match;
  }
});

function showInfo(msg, ms) {
  Alpine.store("snackbar").info(msg, ms);
}

function showError(msg, ms) {
  Alpine.store("snackbar").error(msg, ms);
}

function showModal(querySelector) {
  document.querySelector(querySelector).showModal();
}

function closeModal(querySelector) {
  document.querySelector(querySelector).close();
}

function getCookie(name) {
  return (
    document.cookie
      .split("; ")
      .find((v) => v.startsWith(name + "="))
      ?.split("=")[1] || ""
  );
}
