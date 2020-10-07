FROM redis:5.0-alpine

RUN apk add openrc --no-cache

RUN apk add --no-cache \
    stunnel~=5.56 \
    python3~=3.8

COPY ./certs/key.pem /etc/stunnel/cert/key.pem
COPY ./certs/cert.pem /etc/stunnel/cert/cert.pem
RUN chmod 640 /etc/stunnel/cert/key.pem
RUN chmod 640 /etc/stunnel/cert/cert.pem

WORKDIR /etc/

EXPOSE 6379 6378
