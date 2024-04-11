FROM golang:1.22-alpine as builder

RUN apk add --no-cache make && mkdir /build 

COPY . /build

RUN cd /build && go mod download

RUN cd /build && make build

FROM alpine:3.19 

ARG ID=1000

RUN apk add --no-cache curl mysql-client zip bash postgresql-client mongodb-tools

RUN addgroup -g $ID devops && adduser -u $ID -G devops -D devops

RUN mkdir /data

# Latest releases available at https://github.com/aptible/supercronic/releases
ENV SUPERCRONIC_URL=https://github.com/aptible/supercronic/releases/download/v0.2.29/supercronic-linux-amd64 \
    SUPERCRONIC=supercronic-linux-amd64 \
    SUPERCRONIC_SHA1SUM=cd48d45c4b10f3f0bfdd3a57d054cd05ac96812b

RUN curl -fsSLO "$SUPERCRONIC_URL" \
    && echo "${SUPERCRONIC_SHA1SUM}  ${SUPERCRONIC}" | sha1sum -c - \
    && chmod +x "$SUPERCRONIC" \
    && mv "$SUPERCRONIC" "/usr/local/bin/${SUPERCRONIC}" \
    && ln -s "/usr/local/bin/${SUPERCRONIC}" /usr/local/bin/supercronic

WORKDIR /home/devops

USER devops

COPY --from=builder /build/bin/bekup /usr/local/bin/

COPY --from=builder /build/configs/crontab /data/



CMD [ "supercronic" ,"/data/crontab"]