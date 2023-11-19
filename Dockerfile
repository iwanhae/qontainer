FROM golang:1.21-alpine3.18 as builder
WORKDIR /app
COPY . .
RUN go build -o /qontainer main.go


FROM alpine:3.18

# Need root to use with KVM
USER root

# Default dir
WORKDIR /data
RUN mkdir -p /data

# Install qemu
RUN apk add qemu qemu-system-x86_64 qemu-img tcpdump iptables
RUN sh -c "echo \"allow br0\" >> /etc/qemu/bridge.conf"

# qontainer
COPY --from=builder /qontainer /usr/local/bin/qontainer

VOLUME [ "/data" ]
ENTRYPOINT [ "qontainer" ]
