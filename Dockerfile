FROM golang:1.14.7

WORKDIR /gopath/
ENV GOPATH=/gopath/
ENV GOOS=linux

COPY . /gopath/src/github.com/lixiangyun/https-gateway

WORKDIR /gopath/src/github.com/lixiangyun/https-gateway/console
RUN go build -ldflags="-w -s" .

FROM ubuntu:18.04
MAINTAINER lixiangyun linimbus@126.com

RUN apt update && apt install -y curl net-tools nginx etcd-server software-properties-common
RUN add-apt-repository ppa:certbot/certbot
RUN apt update && apt install -y certbot

WORKDIR /opt/bin

COPY --from=0 /gopath/src/github.com/lixiangyun/https-gateway/start.sh ./start.sh
COPY --from=0 /gopath/src/github.com/lixiangyun/https-gateway/nginx.conf.template ./nginx.conf.template

COPY --from=0 /gopath/src/github.com/lixiangyun/https-gateway/console/static ./static
COPY --from=0 /gopath/src/github.com/lixiangyun/https-gateway/console/console ./console

RUN chmod +x console start.sh

HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://127.0.0.1:18001/ || exit 1

ENV HTTPS_GATEWAY_HOME /opt/home

EXPOSE 80 443 18000

ENTRYPOINT ["./start.sh"]