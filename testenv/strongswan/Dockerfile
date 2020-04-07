FROM alpine
MAINTAINER cmj <cmj@cmj.tw>

# install necessary library
RUN apk add --no-cache strongswan nftables 

WORKDIR /srv

COPY kick-start.sh .

ENV UNIQUE=replace
ENV DEVICE=eth0
env DOMAIN=vpn.cmj.tw
ENV LAN=192.168.0.0/24
ENV VLAN=10.10.0.0/24
ENV DNS=8.8.8.8,8.8.4.4

EXPOSE 500/udp 4500/udp

ENTRYPOINT	["/srv/kick-start.sh"]
