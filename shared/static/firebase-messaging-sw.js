importScripts('https://www.gstatic.com/firebasejs/11.0.2/firebase-app-compat.js');
importScripts('https://www.gstatic.com/firebasejs/11.0.2/firebase-messaging-compat.js');

// Firebase 초기화
const firebaseConfig = {
    apiKey: "AIzaSyCWIebyvcBiwWchfYGUegHf22c9nlBEOWQ",
    authDomain: "warm-braid-383411.firebaseapp.com",
    projectId: "warm-braid-383411",
    storageBucket: "warm-braid-383411.firebasestorage.app",
    messagingSenderId: "1001293129594",
    appId: "1:1001293129594:web:a579e07714a18ec3b598c3"
};

// 초기화
firebase.initializeApp(firebaseConfig);

// FCM Messaging 초기화
const messaging = firebase.messaging();

// 백그라운드 메시지 수신
// messaging.onBackgroundMessage(function (payload) {
//     console.log('[firebase-messaging-sw.js] Received background message:', payload);
//     const notificationTitle = payload.notification?.title || 'Default Title';
//     const notificationOptions = {
//         body: payload.notification?.body || 'Default body content',
//         // icon: '/your-icon.png'  // 알림 아이콘 (선택사항)
//     };
//
//     self.registration.showNotification(notificationTitle, notificationOptions);
// });

// 알림 클릭 시 PWA 앱으로 진입
self.addEventListener('notificationclick', function (event) {
    event.notification.close();

    event.waitUntil(
        clients.matchAll({type: 'window', includeUncontrolled: true}).then(function (clientList) {
            for (const client of clientList) {
                // PWA가 이미 열려 있는 경우 포커스
                if (client.url.includes(self.registration.scope) && 'focus' in client) {
                    return client.focus();
                }
            }

            // PWA가 열려있지 않은 경우 새 창으로 열기
            if (clients.openWindow) {
                return clients.openWindow(self.registration.scope);  // PWA 경로로 설정
            }
        })
    );
});