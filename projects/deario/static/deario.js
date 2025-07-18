console.log("Deario에 오신걸 환영합니다.");

document.addEventListener("alpine:init", () => {
  Alpine.store("save", {
    isOk: true,

    ok() {
      this.isOk = true;
    },

    unok() {
      this.isOk = false;
    },
  });
});

let recorder;
let audioChunks = [];

function toggleRecord() {
  if (recorder && recorder.state === "recording") {
    recorder.stop();
    return;
  }
  navigator.mediaDevices
    .getUserMedia({ audio: true })
    .then((stream) => {
      recorder = new MediaRecorder(stream);
      audioChunks = [];
      recorder.ondataavailable = (e) => audioChunks.push(e.data);
      recorder.onstop = () => {
        const blob = new Blob(audioChunks, { type: "audio/webm" });
        sendAudio(blob);
      };
      recorder.start();
    })
    .catch(() => alert("마이크 권한이 필요합니다."));
}

function sendAudio(blob) {
  const formData = new FormData();
  formData.append("audio", blob, "record.webm");
  fetch("/ai-feedback/audio", {
    method: "POST",
    body: formData,
  })
    .then((res) => res.json())
    .then((data) => {
      document.querySelector("#diary textarea[name='content']").value =
        data.text;
    })
    .catch((err) => console.error(err));
}
