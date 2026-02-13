신규 프로젝트 생성은 AI에게 맡겨서 로컬에서 파일구조 잡는건 어렵지 않은데 이제 프로젝트가 만들어지고 나서 운영에 릴리스 시키는 과정이 반복적이라 이과정을 자동화 시키고 싶다.
보통 실제 서비스를 배포하려면 하는 작업은 아래와 같다.

1. 리눅스 실행 바이너리를 만든다.
2. 서버에 접속해서 바이너리 파일을 /home/ubuntu/app 폴더에 바이너리 파일을 넣는다.
3. 아래와 같은 서비스 파일을 만든다.

```{프로젝트명}.service
[Unit]
Description={프로젝트명} Service
After=network.target

[Service]
ExecStart=/home/ubuntu/app/{프로젝트명}
WorkingDirectory=/home/ubuntu/app
Restart=always
User=root
Group=root

[Install]
WantedBy=multi-user.target
```

4. /etc/systemd/system 폴더에 서비스 파일을 넣는다.
5. 아래 명령어를 실행한다.

```shell
sudo systemctl daemon-reload
sudo systemctl enable {프로젝트명}
sudo systemctl start {프로젝트명}
sudo systemctl status {프로젝트명}
```

6. 서비스가 실행이된지 확인이되면 /srv/{프로젝트} 폴더가 생겼는지 확인 한 뒤
7. 아래의 명령어로 해당 /srv 프로젝트에 소유자를 변경한다.

```shell
sudo chown -R www-data:www-data /srv
sudo chmod -R 755 /srv
```

8. Caddy 파일을 업데이트 한다.
   .linux/caddy/Caddyfile (로컬) => /etc/caddy/Caddyfile (리눅스)

9. Caddy를 리로드 한다.

```shell
sudo systemctl reload caddy
```

10. 자연스럽게 역순으로 프로젝트를 서버에서 회수할 수 있다.
