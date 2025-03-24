import {initializeApp} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-app.js";

// Add Firebase products that you want to use
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

// 2. 토큰을 저장할 변수와 Promise
// - authPromise는 "토큰이 준비되기를 기다리는" Promise
let authToken = null;
let authPromise = null;

// 3. onAuthStateChanged로 로그인 / 로그아웃 감지
onAuthStateChanged(auth, (user) => {
    if (user) {
        console.log("로그인됨:", user);
        // document.getElementById("username").textContent = `${user.displayName} 님 환영합니다!`;
        // document.getElementById("login").classList.add("is-hidden");
        // document.getElementById("logout").classList.remove("is-hidden");

        // user가 존재하면, 토큰 가져오는 Promise를 만들어 둠
        authPromise = user
            .getIdToken(/* forceRefresh */ false)
            .then((token) => {
                authToken = token; // 이후 htmx 요청 때 이 token을 쓰면 됨
                return token;
            })
            .catch((err) => {
                console.error("토큰 가져오기 실패:", err);
                throw err;
            });
    } else {
        // 로그아웃 상태
        console.log("로그아웃 상태");
        // document.getElementById("username").textContent = "";
        // document.getElementById("login").classList.remove("is-hidden");
        // document.getElementById("logout").classList.add("is-hidden");

        // token/Promise 초기화
        authToken = null;
        authPromise = null;
    }
})

/*
4. htmx:confirm 이벤트:
- HTMX가 요청을 보내기 직전(사용자 액션
)
에 발생하며,
요청을 계속할지(확인) 여부를 결정.
- 여기서 "토큰이 아직 준비되지 않았다면
" 요청을 잠시 중단했다가,
토큰이 준비된 뒤에 issueRequest()
로 재개.
*/
htmx.on("htmx:confirm", (e) => {
    // // authPromise가 없거나, 아직 user가 null이면
    // if (!authPromise) {
    //   console.warn("아직 로그인 안 됐으므로 HTMX 요청 중단");
    //   e.preventDefault();
    //   return;
    // }

    // // authPromise가 완료될 때까지 대기
    // // (이 시점에서 토큰이 준비됨)
    // if (authToken === null) {
    //   // 이미 Promise는 존재하지만, 토큰이 아직 안 왔을 수도 있으니
    //   e.preventDefault();
    //   authPromise.then(() => {
    //     console.log("토큰이 준비되었으므로 요청 재개");
    //     e.detail.issueRequest(); // 다시 요청을 보냄
    //   });
    // }
})

/*
5. htmx:configRequest 이벤트:
- 실제로 요청을 구성할 때 발생
- 여기에 "Authorization: Bearer <토큰>" 헤더를 추가
*/
htmx.on("htmx:configRequest", (e) => {
    // 토큰이 있다면 헤더에 실어 보냄
    if (authToken) {
        e.detail.headers["Authorization"] = "Bearer " + authToken;
    }
});

document.getElementById("logout").addEventListener("click", () => {
    auth.signOut();
});
