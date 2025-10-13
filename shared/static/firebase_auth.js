import { initializeApp } from "https://www.gstatic.com/firebasejs/11.0.2/firebase-app.js";
import {
  getAuth,
  onAuthStateChanged,
  signOut,
} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";
import {
  getMessaging,
  getToken,
} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-messaging.js";

// 1. Firebase 초기화
const firebaseConfig = {
  apiKey: "AIzaSyCWIebyvcBiwWchfYGUegHf22c9nlBEOWQ",
  authDomain: "warm-braid-383411.firebaseapp.com",
  projectId: "warm-braid-383411",
  storageBucket: "warm-braid-383411.firebasestorage.app",
  messagingSenderId: "1001293129594",
  appId: "1:1001293129594:web:a579e07714a18ec3b598c3",
};
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);
const messaging = getMessaging(app);

window.logoutUser = async function () {
  try {
    await signOut(auth);
    await fetch("/logout", {
      method: "POST",
      headers: { "X-CSRF-Token": getCookie("_csrf") },
    });
    location.reload();
  } catch (err) {
    console.error("로그아웃 실패:", err);
  }
};

// 2. onAuthStateChanged로 로그인 / 로그아웃 감지
onAuthStateChanged(auth, (user) => {
  if (user) {
    console.log("✅ 로그인됨:", user);
    Alpine.store("auth").login(user);
    // const el = document.querySelector('[hx-trigger="firebase:authed"]');
    // htmx.trigger(el, 'firebase:authed');
  } else {
    // 로그아웃 상태
    console.log("🚪 로그아웃됨");
    Alpine.store("auth").logout();
    // const el = document.querySelector('[hx-trigger="firebase:unauthed"]');
    // htmx.trigger(el, 'firebase:unauthed');
  }
});

let reauthInProgress = false;
htmx.on("htmx:afterRequest", (event) => {
  const contentType = (
    event.detail.xhr.getResponseHeader("Content-Type") || ""
  ).toLowerCase();
  if (!contentType.includes("application/json")) {
    return;
  }

  const responseData = event.detail.xhr.responseText;
  if (responseData === "") {
    return;
  }

  const isResponseError = event.detail.xhr.status === 401;
  if (isResponseError && !reauthInProgress) {
    showInfo("자동 로그인 재시도중입니다.");
    reauthInProgress = true;
    auth.authStateReady().then(() => {
      if (auth.currentUser === undefined) {
        location.href = "/login";
      }

      auth.currentUser
        .getIdToken(true)
        .then((idToken) => {
          return fetch("/create-session", {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
              "X-CSRF-Token": getCookie("_csrf"),
            },
            body: JSON.stringify({ token: idToken }),
          });
        })
        .then((response) => {
          if (response.ok) {
            showInfo("로그인이 완료되었습니다.");
            setTimeout(() => {
              location.reload();
            }, 500);
          } else {
            showError("자동 로그인에 실패했습니다.");
          }
        })
        .catch((err) => {
          console.error("세션 생성 중 에러:", err);
        })
        .finally(() => {
          reauthInProgress = false;
        });
    });
  }
});

async function requestPermissionAndGetToken() {
  try {
    const permission = await Notification.requestPermission();
    if (permission !== "granted") {
      throw new Error("Permission not granted.");
    }

    const registration = await navigator.serviceWorker.register(
      "/firebase-messaging-sw.js"
    );
    console.log("Service Worker registered:", registration);

    const token = await getToken(messaging, {
      vapidKey:
        "BFTAfRBfcOTDygKFWmR1PlFincyIeDa4jC-_6VfLUx-ZvlfBOiM7Wx3VbkpY_jAngZz2MqSsZBp0bpiuRzcJ_G4", // FCM 콘솔에서 발급한 Web Push 인증키
      serviceWorkerRegistration: registration,
    });

    if (token) {
      console.log("FCM Token:", token);
      // 이 토큰을 서버에 저장해두고, 나중에 이걸로 푸시 보냄
      fetch("save-pushToken", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          "X-CSRF-Token": getCookie("_csrf"),
        },
        body: JSON.stringify({
          token: token,
        }),
      });
    } else {
      console.log("No token available.");
    }
  } catch (error) {
    console.error(
      "An error occurred while getting permission or token:",
      error
    );
  }
}

requestPermissionAndGetToken();
