FROM alpine:3.20

RUN apk add --no-cache openssh-server \
    && adduser -D -s /bin/sh zombie \
    && echo "zombie:zombie_ssh_pass" | chpasswd \
    && mkdir -p /run/sshd \
    && printf "\nPasswordAuthentication yes\nPermitRootLogin no\nUsePAM no\n" >> /etc/ssh/sshd_config

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 22
CMD ["/entrypoint.sh"]
