FROM golang:1.16

WORKDIR /go/src/app
COPY . .

RUN go build -o bin/task-runner

RUN cp bin/task-runner /usr/local/bin/

RUN mkdir /var/www/merged
WORKDIR /var/www/merged