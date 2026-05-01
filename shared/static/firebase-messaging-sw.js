importScripts(
  "https://www.gstatic.com/firebasejs/11.0.2/firebase-app-compat.js",
)
importScripts(
  "https://www.gstatic.com/firebasejs/11.0.2/firebase-messaging-compat.js",
)

// Firebase 초기화
const firebaseConfig = {
  apiKey: "AIzaSyCWIebyvcBiwWchfYGUegHf22c9nlBEOWQ",
  authDomain: "warm-braid-383411.firebaseapp.com",
  projectId: "warm-braid-383411",
  storageBucket: "warm-braid-383411.firebasestorage.app",
  messagingSenderId: "1001293129594",
  appId: "1:1001293129594:web:a579e07714a18ec3b598c3",
}

// 초기화
firebase.initializeApp(firebaseConfig)

// FCM Messaging 초기화
const messaging = firebase.messaging()

const OFFLINE_NAVIGATION_HTML = `<!doctype html>
<html lang="ko">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>Deario 오프라인</title>
    <style>
      :root {
        color-scheme: light;
      }

      body {
        margin: 0;
        min-height: 100vh;
        display: grid;
        place-items: center;
        background: #fff8f8;
        color: #211a1a;
        font-family:
          system-ui,
          -apple-system,
          BlinkMacSystemFont,
          "Segoe UI",
          sans-serif;
      }

      main {
        width: min(360px, calc(100vw - 48px));
      }

      h1 {
        margin: 0 0 24px;
        font-size: 40px;
        font-weight: 500;
      }

      p {
        margin: 0 0 10px;
        font-size: 16px;
        line-height: 1.6;
      }

      button {
        min-height: 44px;
        margin-top: 24px;
        padding: 0 20px;
        border: 0;
        border-radius: 22px;
        background: #a6383b;
        color: #ffffff;
        font: inherit;
        font-weight: 700;
      }
    </style>
  </head>
  <body>
    <main>
      <h1>Deario</h1>
      <p>오프라인 상태입니다.</p>
      <p>일기 내용 보호를 위해 저장된 화면을 표시하지 않습니다.</p>
      <button type="button" onclick="location.reload()">새로고침</button>
    </main>
  </body>
</html>`

self.addEventListener("install", function () {
  self.skipWaiting()
})

self.addEventListener("activate", function (event) {
  event.waitUntil(self.clients.claim())
})

self.addEventListener("fetch", function (event) {
  if (event.request.method !== "GET" || event.request.mode !== "navigate") {
    return
  }

  event.respondWith(
    fetch(event.request).catch(function () {
      return new Response(OFFLINE_NAVIGATION_HTML, {
        status: 200,
        headers: {
          "Content-Type": "text/html; charset=utf-8",
          "Cache-Control": "no-store",
        },
      })
    }),
  )
})

// 백그라운드 메시지 수신
messaging.onBackgroundMessage(function (payload) {
  console.log(
    "[firebase-messaging-sw.js] Received background message:",
    payload,
  )
  const notificationTitle = payload.data.title || "Default Title"
  const notificationOptions = {
    body: payload.data.body || "Default body content",
    data: payload.data,
    // icon: '/your-icon.png'  // 알림 아이콘 (선택사항)
  }

  self.registration.showNotification(notificationTitle, notificationOptions)
})

// 알림 클릭 시 PWA 앱으로 진입
self.addEventListener("notificationclick", function (event) {
  event.notification.close()

  // 알림 데이터에서 URL 가져오기 (URL이 없으면 기본값 사용)
  const urlToOpen = event.notification.data?.url || self.registration.scope

  event.waitUntil(
    clients
      .matchAll({ type: "window", includeUncontrolled: true })
      .then(function (clientList) {
        for (const client of clientList) {
          // PWA가 이미 열려 있는 경우 포커스
          if (
            client.url.includes(self.registration.scope) &&
            "focus" in client
          ) {
            return client.focus()
          }
        }

        // PWA가 열려있지 않은 경우 새 창으로 열기
        if (clients.openWindow) {
          return clients.openWindow(urlToOpen)
        }
      }),
  )
})
