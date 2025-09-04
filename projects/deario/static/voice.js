let recorder;
let chunks = [];
let micStream; // 스트림을 전역으로 보관

async function toggleRecord(btn) {
  // 녹음 중이면 정지
  if (recorder && recorder.state === "recording") {
    try {
      recorder.stop(); // 녹음 마감
    } finally {
      // 마이크 점유 해제
      if (recorder?.stream) {
        recorder.stream.getTracks().forEach((t) => t.stop());
      }
      if (micStream) {
        micStream.getTracks().forEach((t) => t.stop());
        micStream = null;
      }
      btn.classList.remove("primary");
      btn.innerHTML = "<i>mic</i>";
    }
    return;
  }

  try {
    micStream = await navigator.mediaDevices.getUserMedia({
      audio: { echoCancellation: true, noiseSuppression: true },
    });

    recorder = new MediaRecorder(micStream, { mimeType: "audio/webm" });
    chunks = [];

    recorder.ondataavailable = (e) => {
      if (e.data.size > 0) chunks.push(e.data);
    };

    recorder.onstop = async () => {
      // 혹시 남은 데이터가 있으면 플러시 (브라우저별 안전망)
      try {
        recorder.requestData?.();
      } catch {}

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
      } finally {
        // 이중 안전장치: onstop에서도 트랙 정리
        if (recorder?.stream) {
          recorder.stream.getTracks().forEach((t) => t.stop());
        }
        if (micStream) {
          micStream.getTracks().forEach((t) => t.stop());
          micStream = null;
        }
        recorder = null;
        chunks = [];
      }
    };

    recorder.start();
    btn.classList.add("primary");
    btn.innerHTML = "<i>stop</i>";
  } catch {
    showInfo("마이크 접근 실패");
  }
}
