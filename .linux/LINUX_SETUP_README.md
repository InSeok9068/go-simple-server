```shell
# 루트 비밀번호 설정
sudo passwd root

# 패키지 업데이트
sudo apt update

# nginx 설치
sudo apt install nginx -y

# nginx 서비스 등록
sudo systemctl enable nginx

# nginx 설정 파일 수정

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

sudo systemctl daemon-reload
enabled
start
restart
stop
status

```shell
sudo systemctl status nginx
sudo systemctl status ufw
sudo systemctl status pocketbase.service
sudo systemctl status pocketbase.service
```
