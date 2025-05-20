/**
 * 알림 관리 모듈
 * 브라우저 알림 권한 요청 및 Firebase 푸시 토큰을 관리합니다.
 */
const NotificationModule = (function() {
    // 개인 변수와 함수
    let _messaging = null;
    let _initialized = false;

    /**
     * Firebase 메시징 인스턴스 초기화
     * @param {Object} firebaseMessaging - Firebase 메시징 인스턴스
     */
    function _initMessaging(firebaseMessaging) {
        _messaging = firebaseMessaging;
        _initialized = true;
    }

    /**
     * 알림 권한 상태 확인
     * @returns {string} 'granted', 'denied', 'default' 중 하나
     */
    function _checkPermissionStatus() {
        return Notification.permission;
    }

    /**
     * 알림 권한 요청
     * @returns {Promise<boolean>} 권한 승인 여부
     */
    async function _requestPermission() {
        try {
            const permission = await Notification.requestPermission();
            return permission === 'granted';
        } catch (error) {
            console.error('알림 권한 요청 중 오류:', error);
            return false;
        }
    }

    /**
     * Firebase 푸시 토큰 요청 및 저장
     * @returns {Promise<boolean>} 토큰 발급 및 저장 성공 여부
     */
    async function _getFirebaseToken() {
        if (!_initialized || !_messaging) {
            console.error('Firebase 메시징이 초기화되지 않았습니다.');
            return false;
        }

        try {
            const registration = await navigator.serviceWorker.register('/firebase-messaging-sw.js');
            console.log('Service Worker registered:', registration);

            const token = await getToken(_messaging, {
                vapidKey: 'BFTAfRBfcOTDygKFWmR1PlFincyIeDa4jC-_6VfLUx-ZvlfBOiM7Wx3VbkpY_jAngZz2MqSsZBp0bpiuRzcJ_G4',
                serviceWorkerRegistration: registration,
            });

            if (token) {
                console.log('FCM Token:', token);
                await _saveTokenToServer(token);
                return true;
            } else {
                console.log('토큰을 발급받지 못했습니다.');
                return false;
            }
        } catch (error) {
            console.error('푸시 토큰 발급 중 오류:', error);
            return false;
        }
    }

    /**
     * 토큰을 서버에 저장
     * @param {string} token - Firebase 푸시 토큰
     * @returns {Promise<boolean>} 서버 저장 성공 여부
     */
    async function _saveTokenToServer(token) {
        try {
            const response = await fetch('save-pushToken', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ token })
            });
            return response.ok;
        } catch (error) {
            console.error('토큰 서버 저장 중 오류:', error);
            return false;
        }
    }

    // 공개 API
    return {
        /**
         * 모듈 초기화
         * @param {Object} firebaseMessaging - Firebase 메시징 인스턴스
         */
        init: function(firebaseMessaging) {
            _initMessaging(firebaseMessaging);
            
            // Alpine 스토어 통합
            if (window.Alpine) {
                document.addEventListener('alpine:init', () => {
                    Alpine.store('notification', {
                        permission: _checkPermissionStatus() === 'granted',
                        
                        async requestPermission() {
                            const granted = await _requestPermission();
                            this.permission = granted;
                            
                            if (granted) {
                                return await NotificationModule.setupPushNotification();
                            }
                            return false;
                        },
                        
                        updatePermissionStatus() {
                            this.permission = _checkPermissionStatus() === 'granted';
                        }
                    });
                });
            }
            
            // 페이지 로드 시 권한이 있으면 자동으로 토큰 요청
            document.addEventListener('DOMContentLoaded', () => {
                if (_checkPermissionStatus() === 'granted') {
                    this.setupPushNotification();
                }
                
                // Alpine 스토어 업데이트
                if (window.Alpine && Alpine.store('notification')) {
                    Alpine.store('notification').updatePermissionStatus();
                }
            });
        },
        
        /**
         * 알림 권한 및 푸시 설정 요청
         * @returns {Promise<boolean>} 설정 성공 여부
         */
        setupPushNotification: async function() {
            if (_checkPermissionStatus() !== 'granted') {
                const granted = await _requestPermission();
                if (!granted) return false;
            }
            
            return await _getFirebaseToken();
        },
        
        /**
         * 현재 알림 권한 상태 확인
         * @returns {boolean} 알림 권한 있음 여부
         */
        hasPermission: function() {
            return _checkPermissionStatus() === 'granted';
        }
    };
})();

// 전역으로 노출
window.NotificationModule = NotificationModule;
