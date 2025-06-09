1. 불필요한 복잡성보다는 직관적인 코드 선호

## 프로젝트 구성

- 언어 : Go
- 서버 프레임워크 : Echo
- 템플릿 라이브러리 : Gomponents
- 데이터베이스 : SQLite + SQLC
- 프론트엔드 라이브러리 : HTMX + Alpinejs
- CSS 프레임워크 : BeerCSS

## 폴더 구조

- cmd/{프로젝트명}/main.go : 해당 프로젝트 서버 실행 파일
- projects/{프로젝트명}/ : 해당 프로젝트 폴더
- projects/{프로젝트명}/static : 해당 프로젝트 정적소스 폴더
- internal : 프로젝트 의존성이 존재하는 공통/공유 서버 패키지
- shared : 프로젝트 의존성이 존재하는 공통/공유 프론트엔드 패키지
- pkg : 프로젝트 의존성 없는 공통 패키지

## 기타

- Gomponents 특성상 자바스크립트는 별도의 JS 파일로 분리 후 사용