const firebaseConfig = {
    apiKey: "AIzaSyCWIebyvcBiwWchfYGUegHf22c9nlBEOWQ",
    authDomain: "warm-braid-383411.firebaseapp.com",
    projectId: "warm-braid-383411",
    storageBucket: "warm-braid-383411.firebasestorage.app",
    messagingSenderId: "1001293129594",
    appId: "1:1001293129594:web:a579e07714a18ec3b598c3"
};

firebase.initializeApp(firebaseConfig);

const uiConfig = {
    // 로그인 성공 시 이동할 URL, 혹은 자동 리다이렉트를 막을 수도 있음
    signInFlow: "popup",
    signInSuccessUrl: "/",
    signInOptions: [
        // 이메일 + 비밀번호 로그인
        firebase.auth.EmailAuthProvider.PROVIDER_ID,
        // 소셜 로그인 예시 (Google)
        firebase.auth.GoogleAuthProvider.PROVIDER_ID,
    ],
    callbacks: {
        signInSuccessWithAuthResult: function (authResult, redirectUrl) {
            authResult.user.getIdToken(true).then((idToken) => {
                return fetch('/create-session', {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({token: idToken})
                })
            }).then(response => {
                if (response.ok) {
                    window.location.href = "/";
                } else {
                    alert("세션 생성 실패");
                }
            }).catch((err) => {
                console.error("세션 생성 중 에러:", err);
            });

            // return false로 하면 signInSuccessUrl로 리다이렉트 안 함
            return false;
        },
        signInFailure: function (error) {
            // Some unrecoverable error occurred during sign-in.
            // Return a promise when error handling is completed and FirebaseUI
            // will reset, clearing any UI. This commonly occurs for error code
            // 'firebaseui/anonymous-upgrade-merge-conflict' when merge conflict
            // occurs. Check below for more details on this.
            return handleUIError(error);
        },
        uiShown: function () {
            // The widget is rendered.
            // Hide the loader.
            document.getElementById('loader').style.display = 'none';
        }
    },
};

// (7) FirebaseUI 인스턴스 생성 및 초기화
const ui = new firebaseui.auth.AuthUI(firebase.auth());
ui.start("#firebaseui-auth-container", uiConfig);