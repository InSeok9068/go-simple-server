# Adminer 설치

링크 : https://www.adminer.org/en/

### 설치

```bash
sudo apt install -y php php-fpm php-sqlite3
sudo mkdir /var/www/html/adminer
sudo wget https://www.adminer.org/latest.php -O /var/www/html/adminer/index.php
```

### 플러그인 설치 (SQLite 로그인 비밀번호 없이)

```bash
sudo mkdir /var/www/html/adminer/adminer-plugins
sudo wget https://www.adminer.org/download/v5.3.0/plugins/login-password-less.php -O /var/www/html/adminer/adminer-plugins/login-password-less.php
sudo wget https://www.adminer.org/download/v5.3.0/adminer/sqlite.php -O /var/www/html/adminer/sqlite.php
```

### sqlite.php 수정

```php
...
include_once "../plugins/login-password-less.php"; => include_once "./adminer-plugins/login-password-less.php";
...
```

### 데이터베이스 파일 접근권한 부여

sudo chown -R www-data:www-data /srv
sudo chmod -R 755 /srv

### Caddyfile 수정

```Caddyfile
db.toy-project.n-e.kr {
    root * /var/www/html/adminer

    php_fastcgi unix//run/php/php-fpm.sock
    file_server
}
```

✅ **아래 링크 접속 확인**

https://db.toy-project.n-e.kr/sqlite.php
