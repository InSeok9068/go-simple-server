https://github.com/IHATEGIVINGAUSERNAME/theProtector

## 의존성 설치

```shell
# Ubuntu/Debian
sudo apt update
sudo apt install yara jq inotify-tools bpfcc-tools netcat-openbsd python3 \
    linux-headers-$(uname -r)
```

## 스크립트 다운로드

```shell
cd /opt
git clone https://github.com/IHATEGIVINGAUSERNAME/theprotector.git
cd theProtector/
chmod +x theprotector.sh

# Install systemd service (recommended for servers)
sudo ./theprotector.sh systemd
```

## 대시보드 실행

```shell
nohup bash -lc 'tail -f /dev/null | ./theprotector.sh api' \
  >~/theprotector-dashboard.log 2>&1 &
```

```shell
# tmux 세션 실행
tmux new -d -s protector 'sudo -E ./theprotector.sh dashboard'
```

## 대시보드 종료

```shell
# tmux 세션 확인
tmux ls

# tmux 세션 종료
tmux kill-session -t protector
```

```shell
# tmux 세션 확인
tmux attach -t protector
```
