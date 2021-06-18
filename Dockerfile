FROM golang:1.16 AS task_runner

WORKDIR /go/src/app
COPY . .

RUN go build -o bin/task-runner

FROM php:7.1-fpm

RUN apt update
RUN apt upgrade -y
RUN apt install -y apt-utils
RUN apt install -y libmcrypt-dev
RUN apt install -y libicu-dev
RUN docker-php-ext-install -j$(nproc) intl
RUN apt-get install -y libfreetype6-dev libjpeg62-turbo-dev
RUN docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/
RUN docker-php-ext-install -j$(nproc) gd
RUN apt install -y libldb-dev
RUN apt install -y libldap2-dev
RUN apt install -y libxml2-dev
RUN apt install -y libssl-dev
RUN apt install -y libxslt-dev
RUN apt install -y libpq-dev
RUN apt install -y libsqlite3-dev
RUN apt install -y libsqlite3-0
RUN apt install -y libc-client-dev
RUN apt install -y libkrb5-dev
RUN apt install -y curl
RUN apt install -y firebird-dev
RUN apt-get install -y libpspell-dev
RUN apt-get install -y aspell-en
RUN apt-get install -y aspell-de
RUN apt install -y libtidy-dev
RUN apt install -y libsnmp-dev
RUN apt install -y librecode0
RUN apt install -y librecode-dev
RUN docker-php-ext-install opcache
RUN docker-php-ext-install soap
RUN docker-php-ext-install ftp
RUN docker-php-ext-install xsl
RUN docker-php-ext-install bcmath
RUN docker-php-ext-install calendar
RUN docker-php-ext-install ctype
RUN docker-php-ext-install dba
RUN docker-php-ext-install dom
RUN docker-php-ext-install zip
RUN docker-php-ext-install session
RUN docker-php-ext-configure ldap --with-libdir=lib/x86_64-linux-gnu
RUN docker-php-ext-install ldap
RUN docker-php-ext-install json
RUN docker-php-ext-install hash
RUN docker-php-ext-install sockets
RUN docker-php-ext-install pdo
RUN docker-php-ext-install mbstring
RUN docker-php-ext-install tokenizer
RUN docker-php-ext-install pgsql
RUN docker-php-ext-install pdo_pgsql
RUN docker-php-ext-install pdo_mysql
RUN docker-php-ext-install pdo_sqlite
RUN docker-php-ext-install intl
RUN docker-php-ext-install mcrypt
RUN docker-php-ext-install mysqli
RUN docker-php-ext-configure imap --with-kerberos --with-imap-ssl
RUN docker-php-ext-install imap
RUN docker-php-ext-install gd
RUN docker-php-ext-install exif
RUN docker-php-ext-install fileinfo
RUN docker-php-ext-install gettext
RUN docker-php-ext-install iconv
RUN docker-php-ext-install interbase
RUN docker-php-ext-install pdo_firebird
RUN docker-php-ext-install opcache
RUN docker-php-ext-install pcntl
RUN docker-php-ext-install phar
RUN docker-php-ext-install posix
RUN docker-php-ext-install pspell
RUN docker-php-ext-install recode
RUN docker-php-ext-install shmop
RUN docker-php-ext-install simplexml
RUN docker-php-ext-install snmp
RUN docker-php-ext-install sysvmsg
RUN docker-php-ext-install sysvsem
RUN docker-php-ext-install sysvshm
RUN docker-php-ext-install tidy
RUN docker-php-ext-install wddx
RUN docker-php-ext-install xmlrpc
RUN docker-php-ext-install xmlwriter
RUN pecl install grpc-1.36.0
RUN docker-php-ext-enable grpc
RUN apt-get update && apt-get install -y \
    libmagickwand-dev --no-install-recommends \
    && pecl install imagick \
	&& docker-php-ext-enable imagick

RUN apt-get update && apt-get install -y zlib1g-dev libzmq3-dev wget git \
    && pecl install zmq-beta \
    && docker-php-ext-enable zmq

RUN wget https://getcomposer.org/installer -O - -q \
    | php -- --install-dir=/bin --filename=composer --version=1.4.2 --quiet
RUN composer global require "fxp/composer-asset-plugin:^1.4.1"

ENV NVM_DIR /usr/local/nvm
RUN curl -o- https://raw.githubusercontent.com/creationix/nvm/v0.33.1/install.sh | bash
ENV NODE_VERSION v10.24.1
RUN /bin/bash -c "source $NVM_DIR/nvm.sh && nvm install $NODE_VERSION && nvm use --delete-prefix $NODE_VERSION"

ENV NODE_PATH $NVM_DIR/versions/node/$NODE_VERSION/lib/node_modules
ENV PATH      $NVM_DIR/versions/node/$NODE_VERSION/bin:$PATH
RUN npm install less -g
RUN npm install -g yarn
RUN yarn upgrade
RUN npm install -g cross-env

COPY --from=task_runner /go/src/app/bin/task-runner /usr/local/bin  
