const firebaseConfig = {
  apiKey: "AIzaSyCWIebyvcBiwWchfYGUegHf22c9nlBEOWQ",
  authDomain: "warm-braid-383411.firebaseapp.com",
  projectId: "warm-braid-383411",
  storageBucket: "warm-braid-383411.firebasestorage.app",
  messagingSenderId: "1001293129594",
  appId: "1:1001293129594:web:a579e07714a18ec3b598c3",
};

firebase.initializeApp(firebaseConfig);

const uiConfig = {
  // 로그인 성공 시 이동할 URL, 혹은 자동 리다이렉트를 막을 수도 있음
  signInFlow: "popup",
  signInSuccessUrl: "/",

  // Google YOLO (You Only Log-in Once)는 사용자가 브라우저에 구글 계정으로 로그인되어 있을 경우,
  // 클릭 한 번으로 바로 로그인할 수 있게 도와주는 매우 편리한 기능입니다.
  credentialHelper: firebaseui.auth.CredentialHelper.GOOGLE_YOLO,

  // 로그인 화면에 서비스 약관과 개인정보처리방침 링크를 추가합니다.
  // 실제 서비스에서는 법적 고지 및 사용자 신뢰를 위해 제공하는 것이 좋습니다.
  // 주석을 해제하고 실제 URL로 교체하여 사용하세요.
  // tosUrl: "/terms-of-service",
  privacyPolicyUrl: "/privacy",

  signInOptions: [
    // 소셜 로그인 (Google)
    firebase.auth.GoogleAuthProvider.PROVIDER_ID,
    // 소셜 로그인 (Facebook)
    firebase.auth.FacebookAuthProvider.PROVIDER_ID,
    // 소셜 로그인 (Github)
    firebase.auth.GithubAuthProvider.PROVIDER_ID,
    // 핸드폰 로그인
    firebase.auth.PhoneAuthProvider.PROVIDER_ID,
    // 이메일 + 비밀번호 로그인
    firebase.auth.EmailAuthProvider.PROVIDER_ID,
  ],
  callbacks: {
    signInSuccessWithAuthResult: function (authResult, redirectUrl) {
      if (authResult.additionalUserInfo.isNewUser) {
        alert("환영합니다! 회원가입이 완료되었습니다.");
      }

      loader.classList.remove("hidden");
      loader.querySelector("span").textContent = "로그인 중...";

      authResult.user
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
            window.location.href = "/";
          } else {
            showError("세션 생성 실패");
          }
        })
        .catch((err) => {
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
      loader.classList.add("hidden");
    },
  },
};

// (7) FirebaseUI 인스턴스 생성 및 초기화
const loader = document.getElementById("loader");
const ui = new firebaseui.auth.AuthUI(firebase.auth());
ui.start("#firebaseui-auth-container", uiConfig);
