FROM golang:1.25-alpine AS builder

RUN apk add --no-cache make && mkdir /build 

COPY . /build

RUN cd /build && go mod download

RUN cd /build && make build

# alpine 3.23: postgresql18-client tersedia. pg_dump 18 bisa men-dump server PostgreSQL lama
# (9.2+), jadi tetap kompatibel dengan instalasi PG 14/16; sebaliknya pg_dump 16 menolak server 18.
FROM alpine:3.23

ARG TARGETARCH
ARG TARGETOS
ARG TARGETPLATFORM
ARG ID=1000

RUN apk add --no-cache curl mysql-client zip bash postgresql18-client mongodb-tools redis \
    && pg_dump --version | grep -q ' 18\.'

RUN addgroup -g $ID devops && adduser -u $ID -G devops -D devops

RUN mkdir /data

# Latest releases available at https://github.com/aptible/supercronic/releases
ENV SUPERCRONIC_URL=https://github.com/aptible/supercronic/releases/download/v0.2.33/supercronic-linux-$TARGETARCH
ENV SUPERCRONIC=supercronic-linux-$TARGETARCH

RUN curl -fsSLO "$SUPERCRONIC_URL" \
    && chmod +x "$SUPERCRONIC" \
    && mv "$SUPERCRONIC" "/usr/local/bin/${SUPERCRONIC}" \
    && ln -s "/usr/local/bin/${SUPERCRONIC}" /usr/local/bin/supercronic

WORKDIR /home/devops

USER devops

COPY --from=builder /build/bin/bekup /usr/local/bin/

COPY --from=builder /build/configs/crontab /data/



CMD [ "supercronic" ,"/data/crontab"]