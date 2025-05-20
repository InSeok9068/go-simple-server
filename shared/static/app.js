/**
 * 앱 기본 모듈
 * Alpine.js 스토어 설정 및 유틸리티 함수를 제공합니다.
 */
const AppModule = (function() {
    // Alpine.js 스토어 설정
    function _setupAlpineStores() {
        document.addEventListener('alpine:init', () => {
            // 인증 스토어
            Alpine.store('auth', {
                isAuthed: false,
                user: null,

                login(user) {
                    this.isAuthed = true;
                    this.user = user;
                },
                
                logout() {
                    this.isAuthed = false;
                    this.user = null;
                }
            });
            
            // 알림 스토어는 notification-module.js에서 설정됨
            // Alpine.store('notification', {...});
        });
    }
    
    // HTMX 응답 오류 처리
    function _setupHtmxErrorHandler() {
        htmx.on("htmx:afterRequest", (event) => {
            const contentType = event.detail.xhr.getResponseHeader("Content-Type");
            if (contentType !== 'application/json') {
                return;
            }

            const responseData = event.detail.xhr.responseText;
            if (responseData === '') {
                return;
            }

            const isResponseError = event.detail.xhr.status >= 400;
            if (isResponseError) {
                const parsedResponse = JSON.parse(responseData);
                if (parsedResponse.message === undefined || parsedResponse.message === '') {
                    return;
                }
                alert(parsedResponse.message);
            }
        });
    }
    
    // 공개 API
    return {
        /**
         * 앱 모듈 초기화
         */
        init: function() {
            _setupAlpineStores();
            _setupHtmxErrorHandler();
        },
        
        /**
         * 모달 열기
         * @param {string} querySelector - 모달 요소 선택자
         */
        showModal: function(querySelector) {
            document.querySelector(querySelector).showModal();
        },
        
        /**
         * 모달 닫기
         * @param {string} querySelector - 모달 요소 선택자
         */
        closeModal: function(querySelector) {
            document.querySelector(querySelector).close();
        }
    };
})();

// 모듈 초기화
document.addEventListener('DOMContentLoaded', () => {
    AppModule.init();
});

// 전역 함수 노출
window.showModal = AppModule.showModal;
window.closeModal = AppModule.closeModal;