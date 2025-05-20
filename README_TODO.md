[TOC]

## 0. 홈페이지 (homepage)

- 서비스 소개 사이트

**TODO**

- [ ]

---

## 1. AI 공부 도우미 (ai-study)

- 단어를 제공해주면 10가지 공부 주제를 제공

**TODO**

- [ ]

---

## 2. 나만의 일기장 (deario)

- 일기를 작성하고 AI에게 공감 피드백을 제공

**TODO**

- [ ]

---

## ✅ 서비스 추가

**TODO**

- [ ] `projects/서비스` 폴더 생성
- [ ] `cmd/서비스/main.go` 파일 생성
- [ ] [embed.go](embed.go) 서비스 static 경로 추가
- [ ] [change.sh](change.sh) 서비스 추가
    - 서버경로 : /etc/caddy/`Caddyfile`
- [ ] `서비스.service` 생성
    - 빌드 파일 서버 전송 (초기 배포)
    - chmod +x 서비스
    - [deario.service](.linux/systemctl/deario.service) 참고
    - 서버경로 : /etc/systemd/system/`서비스.service`
    - sudo systemctl start 서비스.service
    - sudo systemctl enable 서비스.service
- [ ] `서비스-admin.service` 생성 [선택]
    - [pocketbase](pocketbase) => `서비스-admin` 파일 서버 전송
    - chmod +x 서비스-admin
    - [deario-admin.service](.linux/systemctl/deario-admin.service) 참고
    - 서버경로 : /etc/systemd/system/`서비스-admin.service`
    - 어드민 계정 생성
        - /home/ubuntu/app/서비스-admin serve --dir /home/ubuntu/app/projects/서비스/pb_data --http=127.0.0.1:?
        - 접속 후 계정 생성
    - sudo systemctl start 서비스-admin.service
    - sudo systemctl enable 서비스-admin.service
- [ ] [Caddyfile](.linux/caddy/Caddyfile) 서비스 프록시 추가

## 프로젝트 개선 사항

# 현재 프로젝트 구조 분석 및 개선 제안

## 현재 구조의 장점

1. `cmd`, `internal`, `pkg`, `shared` 등의 표준 디렉토리 구조를 활용하고 있음
2. 관심사 분리가 명확히 되어 있음
3. Go의 `internal` 패키지 규칙을 활용하여 비공개 코드를 적절히 관리함

## 개선 가능한 부분

### 1. aiclient 패키지 위치 재고

현재 코드에서 보이는 `aiclient` 패키지가 루트 레벨에 있는 것 같습니다. 이는 일반적인 Go 프로젝트 구조에 맞지 않습니다. AI 클라이언트가 프로젝트의 핵심 기능이라면 `internal/ai/client`
또는 재사용 가능한 패키지라면 `pkg/aiclient`로 이동하는 것이 좋습니다.

### 2. 도메인 중심 구조 도입 고려

현재 프로젝트는 기술적 관심사별로 디렉토리가 구분되어 있습니다. 이는 매우 일반적인 방식이지만, 프로젝트 규모가 커지면 다음과 같은 도메인 중심 구조를 고려해볼 수 있습니다:

``` 
/internal/
  /user/      # 사용자 관련 모든 기능
    /handler/   # HTTP 핸들러
    /service/   # 비즈니스 로직
    /repository/ # 데이터 접근
  /auth/      # 인증 관련 모든 기능
    /handler/
    /service/
    /repository/
  /ai/        # AI 관련 모든 기능
    /client/
    /service/
  /common/    # 여러 도메인에서 공유하는 코드
```

이렇게 구성하면 특정 기능을 수정할 때 관련 코드를 찾기 쉽고, 책임 소재가 명확해집니다.

### 3. `pkg` 디렉토리 활용 강화

`pkg` 디렉토리는 외부에서 재사용 가능한 패키지를 위한 곳입니다. 현재는 `util` 패키지가 있지만, 유틸리티성 코드가 `pkg`에 적절히 모여있는지 확인하고, 필요하다면 더 세분화할 수 있습니다:

``` 
/pkg/
  /httputil/     # HTTP 관련 유틸리티
  /timeutil/     # 시간 관련 유틸리티 (현재 date.go)
  /secutil/      # 보안 관련 유틸리티 (현재 auth.go)
  /viewutil/     # 뷰 관련 유틸리티 (현재 gomponents.go)
```

### 4. `shared` 디렉토리 네이밍 재고

Go 언어 커뮤니티에서는 일반적으로 `shared`보다 `web`, `ui`, 또는 `frontend` 같은 이름을 더 흔히 사용합니다. 현재 `shared`는 프론트엔드 코드를 주로 포함하므로 이름을 변경하는 것을
고려해볼 수 있습니다.

### 5. 설정 파일 관리 개선

현재 구조에서는 환경 변수 설정(`internal/config/env.go`)이 보입니다. 프로덕션/개발/테스트 환경에 따라 설정을 구분하는 구조를 고려해볼 수 있습니다:

``` 
/configs/           # 설정 파일들
  /development/     # 개발 환경 설정
  /production/      # 프로덕션 환경 설정
  /test/            # 테스트 환경 설정
```

### 6. API 버전 관리 구조 도입

API를 제공하는 서버라면, API 버전 관리를 위한 구조를 고려해볼 수 있습니다:

``` 
/internal/api/
  /v1/        # API 버전 1
  /v2/        # API 버전 2
```

## 결론

현재의 프로젝트 구조는 이미 Go 언어의 관례를 상당 부분 따르고 있어 큰 문제가 없습니다. 제안된 개선 사항들은 프로젝트의 규모와 복잡성에 따라 선택적으로 적용하면 됩니다.
특히 현재 구조가 작은 프로젝트라면 오히려 현재의 단순한 구조가 더 적합할 수 있습니다. 프로젝트가 성장함에 따라 점진적으로 도메인 중심 구조로 전환하는 것을 고려해볼 수 있습니다.
현재 디렉토리 구조는 이미 합리적이고 Go 언어의 관례를 잘 따르고 있으므로, 당장 큰 변경이 필요하지는 않습니다. 제안한 개선 사항들은 프로젝트의 요구사항과 개발 팀의 선호도에 따라 검토하시면 됩니다.

