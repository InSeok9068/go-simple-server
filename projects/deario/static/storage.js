import { getAuth } from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js"
import {
  getStorage,
  ref,
  uploadBytes,
  getDownloadURL,
} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-storage.js"
;(function () {
  const PREVIEW_SELECTOR = "[data-deario-image-preview]"
  const UPLOAD_SELECTOR = "[data-deario-image-upload]"
  const PENDING_PREVIEW_SELECTOR = '[data-role="pending-preview"]'

  document.addEventListener("change", handlePreviewChange)
  document.addEventListener("click", handleUploadClick)

  function handlePreviewChange(event) {
    const input = event.target.closest?.(PREVIEW_SELECTOR)
    if (!input) return

    previewDiaryImage(input)
  }

  function handleUploadClick(event) {
    const button = event.target.closest?.(UPLOAD_SELECTOR)
    if (!button) return

    uploadDiaryImage(button.dataset.date)
  }

  function previewDiaryImage(input) {
    const preview = document.getElementById("diary-image-content")
    if (!preview || !input.files || input.files.length === 0) return

    clearPreviewDiaryImage(preview)

    const file = input.files[0]
    const previewEl = file.type.startsWith("image/")
      ? imagePreview(file)
      : textPreview(input.files.length)

    preview.insertBefore(previewEl, preview.firstChild)
  }

  async function uploadDiaryImage(date) {
    const input = document.getElementById("diary-image-file")
    const loading = document.getElementById("diary-image-loading")

    if (!input || input.files.length === 0) {
      alert("파일이 필요합니다.")
      return
    }

    const auth = getAuth()
    if (!auth.currentUser) {
      alert("로그인이 필요합니다.")
      return
    }

    const file = input.files[0]
    const uid = auth.currentUser.uid
    const uploadYMD = ymdKST(new Date())
    const diaryYMD = ymdKST(date)
    const path = `diary/${uploadYMD}/${uid}/${diaryYMD}/${Date.now()}.${getExt(file)}`

    try {
      showLoading(loading, true)

      const snapshot = await uploadBytes(ref(getStorage(), path), file, {
        contentType: file.type,
        customMetadata: { uploadYMD, diaryYMD, uid },
      })
      const url = await getDownloadURL(snapshot.ref)

      clearPreviewDiaryImage(document.getElementById("diary-image-content"))
      htmx.ajax("POST", "/diary/image", {
        target: "#diary-image-content",
        swap: "outerHTML",
        values: { date, url },
      })
      input.value = ""
    } catch (err) {
      console.error("업로드 실패:", err)
      showError("업로드 실패")
    } finally {
      showLoading(loading, false)
    }
  }

  function imagePreview(file) {
    const img = document.createElement("img")
    img.className = "small-width small-height"
    img.dataset.role = "pending-preview"
    img.src = URL.createObjectURL(file)
    img.onload = () => URL.revokeObjectURL(img.src)
    return img
  }

  function textPreview(count) {
    const text = document.createElement("p")
    text.dataset.role = "pending-preview"
    text.textContent = `${count}개 선택됨`
    return text
  }

  function clearPreviewDiaryImage(preview) {
    preview?.querySelector(PENDING_PREVIEW_SELECTOR)?.remove()
  }

  function showLoading(el, visible) {
    if (el) {
      el.style.display = visible ? "block" : "none"
    }
  }

  function ymdKST(input = new Date()) {
    const date = parseKSTInput(input)
    if (isNaN(date)) {
      throw new Error("Invalid date")
    }

    return new Intl.DateTimeFormat("en-CA", {
      timeZone: "Asia/Seoul",
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
    }).format(date)
  }

  function parseKSTInput(input) {
    if (input instanceof Date) {
      return input
    }

    const value = String(input).trim()
    if (/^\d{8}$/.test(value)) {
      const year = +value.slice(0, 4)
      const month = +value.slice(4, 6) - 1
      const day = +value.slice(6, 8)
      const date = new Date(Date.UTC(year, month, day))
      date.setUTCHours(date.getUTCHours() + 9)
      return date
    }

    if (/^\d{4}-\d{2}-\d{2}$/.test(value)) {
      return new Date(`${value}T00:00:00+09:00`)
    }

    return new Date(input)
  }

  function getExt(file) {
    const byType = {
      "image/jpeg": "jpg",
      "image/png": "png",
      "image/webp": "webp",
      "image/avif": "avif",
    }[file.type]

    if (byType) return byType

    const match = file.name.match(/\.([A-Za-z0-9]+)$/)
    return (match?.[1] || "jpg").toLowerCase()
  }
})()
