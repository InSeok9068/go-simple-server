let recorder;
let chunks = [];

async function toggleRecord(btn) {
  if (recorder && recorder.state === "recording") {
    recorder.stop();
    btn.classList.remove("primary");
    btn.innerHTML = "<i>mic</i>";
    return;
  }
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
    recorder = new MediaRecorder(stream);
    chunks = [];
    recorder.ondataavailable = (e) => {
      if (e.data.size > 0) chunks.push(e.data);
    };
    recorder.onstop = async () => {
      const blob = new Blob(chunks, { type: "audio/webm" });
      const fd = new FormData();
      fd.append("audio", blob, "recording.webm");
      try {
        const res = await fetch("/diary/transcribe", {
          method: "POST",
          headers: { "X-CSRF-Token": getCookie("_csrf") },
          body: fd,
        });
        if (res.ok) {
          const text = await res.text();
          const textarea = document.querySelector(
            "#diary textarea[name='content']"
          );
          if (textarea) {
            textarea.value += (textarea.value ? "\n" : "") + text;
            textarea.dispatchEvent(new Event("input"));
          }
        } else {
          showInfo("음성 인식 실패");
        }
      } catch {
        showInfo("음성 인식 실패");
      }
    };
    recorder.start();
    btn.classList.add("primary");
    btn.innerHTML = "<i>stop</i>";
  } catch {
    showInfo("마이크 접근 실패");
  }
}
