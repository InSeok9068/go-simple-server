# 스캐폴딩

### ✅ 서비스 추가 (로컬 저장소)

- [ ] `projects/서비스이름` 디렉토리 생성
- [ ] `cmd/서비스이름/main.go` 기본 파일 생성
- [ ] `embed.go` 파일에 새 서비스의 `static` 경로 추가 (`//go:embed projects/서비스이름/static`)
- [ ] `embed.go` 파일에 새 서비스의 `migrations` 경로 추가 (`//go:embed projects/서비스이름/migrations/*.sql`)
- [ ] `scripts/change.sh` 스크립트에 새 서비스 이름 추가 (스크립트가 서비스를 인식하도록)

### 🔥 서비스 제거 (로컬 저장소)

- [ ] `projects/서비스이름` 디렉토리 전체 삭제
- [ ] `cmd/서비스이름` 디렉토리 전체 삭제
- [ ] `embed.go` 파일에서 해당 서비스의 `static` 경로 제거
- [ ] `embed.go` 파일에서 해당 서비스의 `migrations` 경로 제거
- [ ] `scripts/change.sh` 스크립트에서 해당 서비스 이름 제거

### 🚀 서비스 배포 (서버)

- [ ] 서비스의 리눅스용 바이너리 빌드 (`scripts/build-linux.sh 서비스이름`)
- [ ] 빌드된 바이너리 파일을 서버의 배포 경로로 전송 (`/app/서비스이름` 또는 별도 경로)
- [ ] 서버로 전송된 바이너리 파일에 실행 권한 부여 (`chmod +x /app/서비스이름`)
- [ ] `.linux/systemctl/` 경로에 `서비스이름.service` 파일 생성 (`deario.service` 파일 참고)
- [ ] 생성된 `서비스이름.service` 파일을 서버의 `/etc/systemd/system/` 경로로 전송
- [ ] `.linux/caddy/Caddyfile` 파일에 새 서비스를 위한 리버스 프록시 설정 추가
- [ ] 수정된 `Caddyfile`을 서버의 `/etc/caddy/` 경로로 전송 후 Caddy 서비스 재시작 (`sudo systemctl reload caddy`)
- [ ] `/srv/서비스이름/data` DB 디렉토리 생성
- [ ] `sudo chown -R www-data:www-data /srv`
- [ ] `sudo chmod -R 755 /srv`
- [ ] systemd 데몬 리로드 (`sudo systemctl daemon-reload`)
- [ ] 서버에서 새 서비스 시작 (`sudo systemctl start 서비스이름.service`)
- [ ] 서버 부팅 시 서비스가 자동으로 시작되도록 활성화 (`sudo systemctl enable 서비스이름.service`)

### 🗑️ 서비스 배포 회수 (서버)

- [ ] 서버에서 실행 중인 서비스 중지 (`sudo systemctl stop 서비스이름.service`)
- [ ] 서버 부팅 시 서비스가 자동 시작되지 않도록 비활성화 (`sudo systemctl disable 서비스이름.service`)
- [ ] 서버의 `/etc/systemd/system/` 경로에서 `서비스이름.service` 파일 삭제
- [ ] systemd 데몬 리로드 (`sudo systemctl daemon-reload`)
- [ ] `/srv/서비스이름/data` DB 디렉토리 삭제
- [ ] 서버의 `/etc/caddy/Caddyfile`에서 해당 서비스의 리버스 프록시 설정 제거 후 Caddy 서비스 재시작 (`sudo systemctl reload caddy`)
- [ ] 서버에 배포된 바이너리 파일 삭제

---

### 스캐폴드 서비스 생성, 삭제 스크립트 동작성 상세 설명

**로컬 동작성 우선**

1. 서비스명 입력 [!필수]
2. 포트 입력 [!필수]
3. `projects/서비스이름` 디렉토리 생성
4. `cmd/서비스이름/main.go` 실행 파일 생성
5. `scripts/change.sh` 스크립트에 서비스명에 따른 포트 설정 추가
6. `embed.go` 파일에 새 서비스의 `static` 경로 추가 (`//go:embed projects/서비스이름/static`)
7. DB가 필요한 프로젝트인지 확인 후 필요하다면 아래의 작업을 수행
   - `embed.go` 파일에 새 서비스의 `migrations` 경로 추가 (`//go:embed projects/서비스이름/migrations/*.sql`)
   - `db` 디렉토리, `sqlc.yaml` 파일, `query.sql` 파일 생성

### Github Action을 통한 서비스 첫 배포 스크립트 동작성 상세 설명

**Github Action을 통한 서비스 첫 배포 후 URL 접근 가능**

1. 서비스 빌드
