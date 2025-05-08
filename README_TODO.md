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
