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
sudo systemctl enable systemd-timesyncd
sudo systemctl start systemd-timesyncd

# nginx 설치
sudo apt install nginx -y

# nginx 서비스 등록
sudo systemctl enable nginx

# ufw 설치
sudo apt install ufw -y

# nginx 방화벽 허용
sudo ufw allow 'Nginx Full'

# SSG 방화벽 허용
sudo ufw allow 'OpenSSH'

# Pocketbase Admin 방화벽 허용
sudo ufw allow 9000

# ufw 활성화
sudo ufw enable

# ufw 규칙 확인
sudo ufw status verbose

# ufw 서비스 등록
sudo systemctl enable ufw
```

## Nginx 설정

`etc > nginx > sites-enabled > default`

[파일 참조](default)

### nginx 설정 확인

```shell
sudo nginx -t
```

## 서비스 관리

- nginx
- ufw
- pocketbase.service
- main.service
- systemd-timesyncd

```shell
sudo systemctl daemon-reload
sudo systemctl start nginx
sudo systemctl stop nginx
sudo systemctl restart nginx
sudo systemctl status nginx
sudo systemctl enable nginx
```

## 로그 파일 관리

```shell
sudo apt install logrotate -y
```

## SSL 설치

```shell
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d toy-project.n-e.kr
sudo certbot renew --dry-run
# 한달에 한번 수동 갱신 (혹시나 동작 안할때 대비)
# 0 0 1 * * certbot renew --renew-hook "sudo service nginx restart"
```
