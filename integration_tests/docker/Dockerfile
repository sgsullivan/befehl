FROM alpine:latest

RUN apk add --update --no-cache openssh bash

RUN mkdir -p /root/.ssh
RUN chmod 700 /root/.ssh
COPY ssh/id_rsa.pub /root/.ssh/authorized_keys
RUN chmod 600 /root/.ssh/authorized_keys

ENTRYPOINT ["/entrypoint.sh"]

EXPOSE 22

COPY entrypoint.sh /
