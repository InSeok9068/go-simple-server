import { getAuth } from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";
import {
  getStorage,
  ref,
  uploadBytes,
  getDownloadURL,
} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-storage.js";

// 한국시간 YYYY-MM-DD로 변환 (Date | 'YYYYMMDD' | 'YYYY-MM-DD' 지원)
function ymdKST(input = new Date()) {
  let d;

  if (input instanceof Date) {
    d = input;
  } else if (typeof input === "number" || typeof input === "string") {
    const s = String(input).trim();

    if (/^\d{8}$/.test(s)) {
      // 'YYYYMMDD' -> KST 자정으로 해석
      const y = +s.slice(0, 4);
      const m = +s.slice(4, 6) - 1;
      const day = +s.slice(6, 8);
      // KST 자정 타임스탬프 만들기
      const kst = new Date(Date.UTC(y, m, day));
      // UTC 자정에 +9시간을 더해 KST 자정으로 맞춤
      kst.setUTCHours(kst.getUTCHours() + 9);
      d = kst;
    } else if (/^\d{4}-\d{2}-\d{2}$/.test(s)) {
      // 'YYYY-MM-DD' -> KST 자정으로 고정
      d = new Date(`${s}T00:00:00+09:00`);
    } else {
      // 기타 문자열은 JS 파서에 위임 (권장하진 않음)
      d = new Date(s);
    }
  } else {
    d = new Date(input);
  }

  if (isNaN(d)) {
    throw new Error("Invalid date");
  }

  return new Intl.DateTimeFormat("en-CA", {
    timeZone: "Asia/Seoul",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
  }).format(d);
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
  const diaryYMD = ymdKST(date); // 일기 날짜 (사용자가 선택)
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
