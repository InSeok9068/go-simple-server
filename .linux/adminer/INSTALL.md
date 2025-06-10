sudo apt install -y php php-fpm php-sqlite3
sudo mkdir /var/www/html/adminer
sudo wget https://www.adminer.org/latest.php -O /var/www/html/adminer/index.php

```Caddyfile
your.domain.com {
    root * /var/www/html

    php_fastcgi unix//run/php/php-fpm.sock {
    }

    file_server
}
```

sudo apt install -y php php-fpm php-sqlite3
sudo mkdir /var/www/html/adminer
sudo wget https://www.adminer.org/latest.php -O /var/www/html/adminer/index.php

sudo chown -R www-data:www-data /var/www/html/adminer
sudo chmod -R 755 /var/www/html/adminer
sudo chown -R www-data:www-data /home/ubuntu/app/projects
sudo chmod -R 755 /home/ubuntu/app/projects

sudo systemctl restart php8.3-fpm
sudo systemctl reload caddy

sudo systemctl status php8.3-fpm

ls -ld /var/www/html/adminer
ls -ld /home/ubuntu/app/projects

sudo tail -f /var/log/php8.3-fpm.log

sudo wget https://www.adminer.org/latest.php -O /var/www/html/adminer/index.php
sudo wget https://github.com/vrana/adminer/releases/download/v4.6.2/adminer-4.6.2.php -O /var/www/html/adminer/index.php
