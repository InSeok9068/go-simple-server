import {initializeApp} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-app.js";
import {getAuth, onAuthStateChanged} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-auth.js";
import {getMessaging, getToken} from "https://www.gstatic.com/firebasejs/11.0.2/firebase-messaging.js";

/**
 * Firebase 인증 및 알림 모듈
 * 인증 상태 관리 및 알림 기능을 처리합니다.
 */
const FirebaseModule = (function() {
    // 개인 변수
    let _app = null;
    let _auth = null;
    let _messaging = null;
    
    // Firebase 설정
    const _firebaseConfig = {
        apiKey: "AIzaSyCWIebyvcBiwWchfYGUegHf22c9nlBEOWQ",
        authDomain: "warm-braid-383411.firebaseapp.com",
        projectId: "warm-braid-383411",
        storageBucket: "warm-braid-383411.firebasestorage.app",
        messagingSenderId: "1001293129594",
        appId: "1:1001293129594:web:a579e07714a18ec3b598c3"
    };
    
    // 401 오류 처리
    function _setupAuthErrorHandler() {
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
                _auth.authStateReady().then(() => {
                    if (_auth.currentUser === undefined) {
                        location.href = "/login";
                    }

                    _auth.currentUser.getIdToken(true).then((idToken) => {
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
    }
    
    // 인증 상태 변화 감지 설정
    function _setupAuthStateChanged() {
        onAuthStateChanged(_auth, (user) => {
            if (user) {
                console.log("✅ 로그인됨:", user);
                if (window.Alpine && Alpine.store('auth')) {
                    Alpine.store('auth').login(user);
                }
            } else {
                console.log("🚪 로그아웃됨");
                if (window.Alpine && Alpine.store('auth')) {
                    Alpine.store('auth').logout();
                }
            }
        });
    }
    
    // 공개 API
    return {
        /**
         * Firebase 모듈 초기화
         */
        init: function() {
            // Firebase 초기화
            _app = initializeApp(_firebaseConfig);
            _auth = getAuth(_app);
            _messaging = getMessaging(_app);
            
            // 인증 상태 변화 감지 설정
            _setupAuthStateChanged();
            
            // 인증 오류 처리 설정
            _setupAuthErrorHandler();
            
            // 알림 모듈 초기화
            if (window.NotificationModule) {
                window.NotificationModule.init(_messaging);
            } else {
                console.warn('알림 모듈이 로드되지 않았습니다.');
            }
            
            return this;
        },
        
        /**
         * 현재 인증된 사용자 반환
         * @returns {Object|null} 인증된 사용자 또는 null
         */
        getCurrentUser: function() {
            return _auth?.currentUser || null;
        },
        
        /**
         * Firebase 인증 인스턴스 반환
         * @returns {Object} Firebase 인증 인스턴스
         */
        getAuth: function() {
            return _auth;
        },
        
        /**
         * Firebase 메시징 인스턴스 반환
         * @returns {Object} Firebase 메시징 인스턴스
         */
        getMessaging: function() {
            return _messaging;
        }
    };
})();

// 모듈 초기화
document.addEventListener('DOMContentLoaded', () => {
    FirebaseModule.init();
});

// 전역으로 노출
window.FirebaseModule = FirebaseModule;
