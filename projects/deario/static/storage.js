import { getAuth } from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";
import {
  getStorage,
  ref,
  uploadBytes,
  getDownloadURL,
} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-storage.js";

window.previewDiaryImage = function (input) {
  const preview = document.getElementById("diary-image-preview");
  if (!preview) return;
  preview.innerHTML = "";
  if (!input || !input.files || input.files.length === 0) return;
  const file = input.files[0];
  if (file.type.startsWith("image/")) {
    const img = document.createElement("img");
    img.className = "responsive";
    img.src = URL.createObjectURL(file);
    img.onload = () => URL.revokeObjectURL(img.src);
    const spaceDiv = document.createElement("div");
    spaceDiv.classList.add("space");
    preview.appendChild(img);
    preview.appendChild(spaceDiv);
  } else {
    preview.textContent = `${input.files.length}개 선택됨`;
  }
};

window.uploadDiaryImage = async function (date) {
  const input = document.getElementById("diary-image-file");
  const preview = document.getElementById("diary-image-preview");
  const loading = document.getElementById("diary-image-loading");
  if (!input || input.files.length === 0) {
    alert("파일이 필요합니다.");
    return;
  }

  const auth = getAuth();
  if (!auth.currentUser) {
    alert("로그인이 필요합니다.");
    return;
  }

  const file = input.files[0];
  const storage = getStorage();
  const path = `diary/${auth.currentUser.uid}/${date}_${Date.now()}`;

  try {
    if (loading) loading.style.display = "block";
    const snapshot = await uploadBytes(ref(storage, path), file);
    const url = await getDownloadURL(snapshot.ref);
    htmx.ajax("POST", "/diary/image", {
      target: "#diary-image-content",
      swap: "outerHTML",
      values: { date: date, url: url },
    });
    if (input) input.value = "";
    if (preview) preview.innerHTML = "";
  } catch (err) {
    console.error("업로드 실패:", err);
    showError("업로드 실패");
  } finally {
    if (loading) loading.style.display = "none";
  }
};

window.viewDiaryImage = function (url) {
  const img = document.getElementById("diary-image-viewer-img");
  if (!img) return;
  img.src = url;
  showModal("#diary-image-viewer-dialog");
};
