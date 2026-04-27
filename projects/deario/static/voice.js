;(function () {
  const TOGGLE_SELECTOR = "[data-deario-voice-toggle]"

  let recorder
  let chunks = []
  let micStream

  document.addEventListener("click", (event) => {
    const button = event.target.closest?.(TOGGLE_SELECTOR)
    if (!button) return

    toggleRecord(button)
  })

  async function toggleRecord(button) {
    if (recorder?.state === "recording") {
      stopRecording(button)
      return
    }

    try {
      micStream = await navigator.mediaDevices.getUserMedia({
        audio: { echoCancellation: true, noiseSuppression: true },
      })
      startRecording(button)
    } catch {
      showInfo("마이크 접근 실패")
    }
  }

  function startRecording(button) {
    recorder = new MediaRecorder(micStream, { mimeType: "audio/webm" })
    chunks = []

    recorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        chunks.push(event.data)
      }
    }
    recorder.onstop = transcribeRecording

    recorder.start()
    setButtonRecording(button, true)
  }

  function stopRecording(button) {
    try {
      recorder.stop()
    } finally {
      releaseMic()
      setButtonRecording(button, false)
    }
  }

  async function transcribeRecording() {
    try {
      recorder?.requestData?.()
    } catch {}

    const formData = new FormData()
    formData.append(
      "audio",
      new Blob(chunks, { type: "audio/webm" }),
      "recording.webm",
    )

    try {
      const response = await fetch("/diary/transcribe", {
        method: "POST",
        headers: { "X-CSRF-Token": getCookie("_csrf") },
        body: formData,
      })

      if (!response.ok) {
        showInfo("음성 인식 실패")
        return
      }

      appendTranscribedText(await response.text())
    } catch {
      showInfo("음성 인식 실패")
    } finally {
      releaseMic()
      recorder = null
      chunks = []
    }
  }

  function appendTranscribedText(text) {
    const textarea = document.querySelector("#diary textarea[name='content']")
    if (!textarea) return

    textarea.value += (textarea.value ? "\n" : "") + text
    textarea.dispatchEvent(new Event("input", { bubbles: true }))
  }

  function releaseMic() {
    recorder?.stream?.getTracks().forEach((track) => track.stop())
    micStream?.getTracks().forEach((track) => track.stop())
    micStream = null
  }

  function setButtonRecording(button, recording) {
    button.classList.toggle("primary", recording)
    button.innerHTML = recording ? "<i>stop</i>" : "<i>mic</i>"
  }
})()
