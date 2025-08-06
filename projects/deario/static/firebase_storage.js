import { getAuth } from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";
import {
  getStorage,
  ref,
  uploadBytes,
  getDownloadURL,
} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-storage.js";

window.uploadDiaryImage = async function (date) {
  const input = document.getElementById("diary-image-file");
  if (!input || input.files.length === 0) {
    return;
  }

  const auth = getAuth();
  if (!auth.currentUser) {
    showError("로그인이 필요합니다.");
    return;
  }

  const file = input.files[0];
  const storage = getStorage();
  const path = `diary/${auth.currentUser.uid}/${date}_${Date.now()}`;

  try {
    const snapshot = await uploadBytes(ref(storage, path), file);
    const url = await getDownloadURL(snapshot.ref);
    htmx.ajax("POST", "/diary/image", {
      target: "#diary-image-content",
      swap: "outerHTML",
      values: { date: date, url: url },
    });
  } catch (err) {
    console.error("업로드 실패:", err);
    showError("업로드 실패");
  }
};
