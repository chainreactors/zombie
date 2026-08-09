#!/bin/sh
set -eu

ssh-keygen -A >/dev/null
exec /usr/sbin/sshd -D -e
