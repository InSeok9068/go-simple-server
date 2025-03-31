console.log("Deario에 오신걸 환영합니다.")

history.pushState(null, document.title, location.href);
window.addEventListener('popstate', function (event) {
    // 뒤로가기 누르면 여기 실행됨
    if (confirm('앱을 종료하시겠습니까?')) {
        // 실제 브라우저에서는 종료 못하니까 종료처럼 보이게 만들기
        window.close(); // 일부 브라우저에서만 동작함
        // 대안: 홈 화면으로 리디렉션
        location.href = 'about:blank';
    } else {
        // 사용자가 종료 원하지 않으면 다시 push
        history.pushState(null, document.title, location.href);
    }
});