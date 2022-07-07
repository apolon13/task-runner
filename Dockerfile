FROM golang:1.16 AS task_runner

WORKDIR /go/src/app
COPY . .

RUN go build -o bin/task-runner

FROM php:8.0.17-fpm

ENV CPPFLAGS "-Wno-maybe-uninitialized"

RUN apt-get update && apt-get upgrade -y

RUN apt-get install -y --no-install-recommends \
    libmcrypt-dev \
    libicu-dev \
    libfreetype6-dev \
    libjpeg62-turbo-dev \
    libldb-dev \
    libldap2-dev \
    libxml2-dev \
    libssl-dev \
    libxslt-dev \
    libpq-dev \
    libsqlite3-dev \
    libsqlite3-0 \
    libc-client-dev \
    libkrb5-dev \
    libtidy-dev \
    libsnmp-dev \
    librecode-dev \
    libpspell-dev \
    libmagickwand-dev \
    libmagickcore-dev \
    libzip-dev \
    libonig-dev \
    zlib1g-dev \
    libzmq3-dev \
    git \
    wget \
    curl \
    librecode0 \
    firebird-dev \
    aspell-en \
    aspell-de \
    cmake \
    unzip

RUN docker-php-ext-install -j$(nproc) \
    intl \
    gd \
    opcache \
    ftp \
    xsl \
    bcmath \
    calendar \
    ctype \
    dom \
    session \
    sockets \
    pdo \
    mbstring \
    tokenizer \
    pdo_mysql \
    mysqli \
    fileinfo \
    gettext \
    iconv \
    opcache \
    pcntl \
    phar \
    posix \
    simplexml \
    sysvmsg \
    sysvsem \
    sysvshm \
    xmlwriter \
    shmop \
    zip

RUN pecl install mongodb \
    && docker-php-ext-enable mongodb

RUN pecl install grpc-1.36.0 \
    && docker-php-ext-enable grpc

RUN pecl install redis \
    && docker-php-ext-enable redis

RUN pecl install igbinary \
    && docker-php-ext-enable igbinary

RUN pecl install imagick \
	&& docker-php-ext-enable imagick

RUN git clone https://github.com/mkoppanen/php-zmq.git \
    && cd php-zmq \
    && phpize \
    && ./configure \
    && make \
    && make install

RUN docker-php-ext-enable zmq

RUN wget https://getcomposer.org/installer -O - -q \
    | php -- --install-dir=/bin --filename=composer --version=2.3.5 --quiet

COPY 'install_node.sh' '/root/install_node.sh'
ENV NVM_DIR /usr/local/nvm
RUN bash /root/install_node.sh
ENV NODE_VERSION v16.13.1
RUN /bin/bash -c "source $NVM_DIR/nvm.sh && nvm install $NODE_VERSION && nvm use --delete-prefix $NODE_VERSION"

ENV NODE_PATH $NVM_DIR/versions/node/$NODE_VERSION/lib/node_modules
ENV PATH      $NVM_DIR/versions/node/$NODE_VERSION/bin:$PATH

RUN npm install -g yarn
RUN npm install -g less
RUN npm install -g less-plugin-clean-css

COPY --from=task_runner /go/src/app/bin/task-runner /usr/local/bin  
