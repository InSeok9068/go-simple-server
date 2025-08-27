## 서버 초기 설정

```shell
# 루트 비밀번호 설정
sudo passwd root

# 패키지 업데이트
sudo apt update

# 시간대 설정
sudo timedatectl set-timezone Asia/Seoul

# 시간 동기화 설정
sudo apt install systemd-timesyncd -y

# UFW 설치
sudo apt install ufw -y

# SSH 방화벽 허용
sudo ufw allow 'OpenSSH'

# Caddy 설치
sudo apt install caddy -y

# Caddy 방화벽 허용
sudo ufw allow 'Caddy Full'

# fail2ban 설치
sudo apt install fail2ban -y

# ufw 활성화
sudo ufw enable

# ufw 규칙 확인
sudo ufw status verbose

# 서비스 등록 확인 (systemd-timesyncd, caddy, ufw, fail2ban)
sudo systemctl list-unit-files | grep caddy
sudo systemctl list-unit-files | grep ufw
sudo systemctl list-unit-files | grep systemd-timesyncd
sudo systemctl list-unit-files | grep fail2ban
sudo systemctl status caddy
sudo systemctl status ufw
sudo systemctl status systemd-timesyncd
sudo systemctl status fail2ban
```

## Caddy 설정

서버 경로 : `/etc/caddy/Caddyfile` </br>
파일 경로 : `./linux/caddy/Caddyfile`

## 로그 파일 관리

```shell
sudo apt install logrotate -y
sudo systemctl list-timers | grep logrotate
sudo systemctl status logrotate.timer
```

## [rod 사용 시 필수!] Chromium 설치

```shell
sudo apt install chromium-browser -y
```

## [선택] 시스템 모니터링 도구 설치

htop : 향상된 top </br>
iotop : 디스크 I/O 모니터링 </br>
iftop : 네트워크 트래픽 모니터링 </br>
nmon : 시스템 성능 모니터링 </br>
ncdu : 디스크 사용량 분석 </br>

```shell
sudo apt install -y htop iotop iftop nmon ncdu
```

## [선택] 자동 보안 업데이트 설정

```shell
sudo apt install unattended-upgrades -y
sudo dpkg-reconfigure unattended-upgrades


# 보안 업데이트 누락 확인
sudo unattended-upgrade --dry-run
```

## [선택] litestream 설치

```shell
wget https://github.com/benbjohnson/litestream/releases/download/v0.3.13/litestream-v0.3.13-linux-amd64.deb
sudo dpkg -i litestream-v0.3.13-linux-amd64.deb
litestream version

sudo systemctl enable litestream
sudo systemctl start litestream
sudo journalctl -u litestream -f

# If you make changes to Litestream configuration file, you’ll need to restart the service
# /etc/litestream.yml
sudo systemctl restart litestream

sudo nano /etc/systemd/system/litestream.service.d/override.conf
sudo systemctl daemon-reload
sudo systemctl restart litestream
sudo journalctl -u litestream -f
```

### override.conf

```ini
[Service]
Environment="GOOGLE_APPLICATION_CREDENTIALS=/etc/secrets/litestream.json"
ExecStart=
ExecStart=/usr/bin/litestream replicate -config /etc/litestream.yml
```

### litestream.yml

```yml
# /etc/litestream.yml
dbs:
  - path: /srv/deario/data/data.db
    replicas:
      - url: gcs://warm-braid-383411.firebasestorage.app/litestream/deario/prod
        snapshot-interval: 24h # 하루마다 스냅샷
        retention: 168h # 7일 보관
```
