import {initializeApp} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-app.js";
import {getAuth, onAuthStateChanged} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";
import {getMessaging, getToken} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-messaging.js";

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
const messaging = getMessaging(app);

// 앱 초기화 시 생체 인증 확인
async function checkBiometricAuthAndSetup() {
    try {
        // 생체 인증이 지원되는지 확인
        if (!window.PublicKeyCredential) {
            console.log('생체 인증이 지원되지 않는 디바이스입니다. 일반 모드로 진입합니다.');
            await setupAppServices();
            return;
        }

        const available = await PublicKeyCredential.isUserVerifyingPlatformAuthenticatorAvailable();
        if (!available) {
            console.log('생체 인증 하드웨어를 사용할 수 없습니다. 일반 모드로 진입합니다.');
            await setupAppServices();
            return;
        }

        // 이미 인증된 상태인지 확인
        const isAuthenticated = sessionStorage.getItem('biometricAuthenticated') === 'true';
        if (isAuthenticated) {
            console.log('이미 생체 인증된 세션입니다.');
            await setupAppServices();
            return;
        }

        // 인증 옵션 설정
        const publicKeyCredentialRequestOptions = {
            challenge: new Uint8Array([21, 31, 105]), // 보안을 위해 서버에서 생성된 챌린지 사용 권장
            rpId: window.location.hostname,
            userVerification: 'required',
            timeout: 60000
        };

        try {
            // 생체 인증 요청
            const credential = await navigator.credentials.get({
                publicKey: publicKeyCredentialRequestOptions
            });

            // 인증 성공
            console.log('생체 인증 성공:', credential);
            sessionStorage.setItem('biometricAuthenticated', 'true');
        } catch (authError) {
            console.error('생체 인증 실패:', authError);
            // 인증 실패해도 앱은 계속 진행
        }

        // 인증 성공 또는 실패 후 앱 초기화
        await setupAppServices();
    } catch (error) {
        console.error('생체 인증 프로세스 중 오류 발생:', error);
        // 오류 발생해도 앱은 계속 진행
        await setupAppServices();
    }
}

// 앱 서비스 설정 함수 (Firebase 함수와 겹치지 않도록 이름 변경)
async function setupAppServices() {
    await requestPermissionAndGetToken();
}

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

async function requestPermissionAndGetToken() {
    try {
        const permission = await Notification.requestPermission();
        if (permission !== 'granted') {
            throw new Error('Permission not granted.');
        }

        const registration = await navigator.serviceWorker.register('/firebase-messaging-sw.js');
        console.log('Service Worker registered:', registration);

        const token = await getToken(messaging, {
            vapidKey: 'BFTAfRBfcOTDygKFWmR1PlFincyIeDa4jC-_6VfLUx-ZvlfBOiM7Wx3VbkpY_jAngZz2MqSsZBp0bpiuRzcJ_G4',  // FCM 콘솔에서 발급한 Web Push 인증키
            serviceWorkerRegistration: registration,
        });

        if (token) {
            console.log('FCM Token:', token);
            // 이 토큰을 서버에 저장해두고, 나중에 이걸로 푸시 보냄
            fetch('save-pushToken', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    token: token
                })
            })
        } else {
            console.log('No token available.');
        }
    } catch (error) {
        console.error('An error occurred while getting permission or token:', error);
    }
}

// 앱 시작 - 생체 인증 확인 후 앱 서비스 설정 실행
checkBiometricAuthAndSetup();