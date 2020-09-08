FROM ubuntu:18.04

MAINTAINER lixiangyun linimbus@126.com

RUN apt update && apt install -y curl net-tools software-properties-common
RUN add-apt-repository ppa:certbot/certbot
RUN apt update && apt install -y certbot nginx

HEALTHCHECK --interval=5m --timeout=3s CMD curl -f http://127.0.0.1:18001/ || exit 1



ENTRYPOINT ["/bin/bash"]