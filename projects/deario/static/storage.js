import { getAuth } from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";
import {
  getStorage,
  ref,
  uploadBytes,
  getDownloadURL,
} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-storage.js";

// 한국시간 YYYY-MM-DD
function ymdKST(date = new Date()) {
  return new Intl.DateTimeFormat("en-CA", {
    timeZone: "Asia/Seoul",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(date);
}

function getExt(file) {
  // contentType 우선 매핑, 없으면 파일명에서 추출
  const map = {
    "image/jpeg": "jpg",
    "image/png": "png",
    "image/webp": "webp",
    "image/avif": "avif",
  };
  const byType = map[file.type];
  if (byType) return byType;
  const m = file.name.match(/\.([A-Za-z0-9]+)$/);
  return (m?.[1] || "jpg").toLowerCase();
}

window.previewDiaryImage = function (input) {
  const preview = document.getElementById("diary-image-content");
  if (!preview) return;
  if (!input || !input.files || input.files.length === 0) return;
  const file = input.files[0];
  if (file.type.startsWith("image/")) {
    const img = document.createElement("img");
    img.className = "small-width small-height";
    img.src = URL.createObjectURL(file);
    img.onload = () => URL.revokeObjectURL(img.src);
    preview.insertBefore(img, preview.firstChild);
  } else {
    preview.textContent = `${input.files.length}개 선택됨`;
  }
};

window.uploadDiaryImage = async function (date) {
  const input = document.getElementById("diary-image-file");
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

  const uid = auth.currentUser.uid;
  const uploadYMD = ymdKST(new Date()); // 업로드 “오늘” (KST)
  const diaryYMD = ymdKST(new Date(date)); // 일기 날짜 (사용자가 선택)
  const ts = Date.now(); // 충돌 방지용 타임스탬프
  const ext = getExt(file);

  // ✅ 업로드일자 → 사용자 → 일기일자 → 타임스탬프
  const path = `diary/${uploadYMD}/${uid}/${diaryYMD}/${ts}.${ext}`;

  try {
    if (loading) loading.style.display = "block";
    const snapshot = await uploadBytes(ref(storage, path), file, {
      contentType: file.type,
      customMetadata: { uploadYMD, diaryYMD, uid },
    });
    const url = await getDownloadURL(snapshot.ref);
    htmx.ajax("POST", "/diary/image", {
      target: "#diary-image-content",
      swap: "outerHTML",
      values: { date: date, url: url },
    });
    if (input) input.value = "";
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
