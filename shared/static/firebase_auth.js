import {initializeApp} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-app.js";
import {getAuth, onAuthStateChanged,} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";

// 1. Firebase 초기화
const firebaseConfig = {
    apiKey: "AIzaSyCWIebyvcBiwWchfYGUegHf22c9nlBEOWQ",
    authDomain: "warm-braid-383411.firebaseapp.com",
    projectId: "warm-braid-383411",
    storageBucket: "warm-braid-383411.firebasestorage.app",
    messagingSenderId: "1001293129594",
    appId: "1:1001293129594:web:a579e07714a18ec3b598c3"
};
const app = initializeApp(firebaseConfig);
const auth = getAuth(app);

// 2. onAuthStateChanged로 로그인 / 로그아웃 감지
onAuthStateChanged(auth, (user) => {
    if (user) {
        console.log("✅ 로그인됨:", user);
        Alpine.store('auth').login(user);
        // const el = document.querySelector('[hx-trigger="firebase:authed"]');
        // htmx.trigger(el, 'firebase:authed');
    } else {
        // 로그아웃 상태
        console.log("🚪 로그아웃됨");
        Alpine.store('auth').logout();
        // const el = document.querySelector('[hx-trigger="firebase:unauthed"]');
        // htmx.trigger(el, 'firebase:unauthed');
    }
})

htmx.on("htmx:afterRequest", (event) => {
    const contentType = event.detail.xhr.getResponseHeader("Content-Type");
    if (contentType !== 'application/json') {
        return;
    }

    const responseData = event.detail.xhr.responseText;
    if (responseData === '') {
        return;
    }

    const isResponseError = event.detail.xhr.status === 401;
    if (isResponseError) {
        auth.authStateReady().then(() => {
            if (auth.currentUser === undefined) {
                location.href = "/login";
            }

            auth.currentUser.getIdToken(true).then((idToken) => {
                return fetch('/create-session', {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({token: idToken})
                })
            }).then(response => {
                if (response.ok) {
                    location.reload()
                } else {
                    alert("재 로그인 실패");
                }
            }).catch((err) => {
                console.error("세션 생성 중 에러:", err);
            });
        })
    }
});

// document.getElementById("logout").addEventListener("click", () => {
//     auth.signOut();
// });
