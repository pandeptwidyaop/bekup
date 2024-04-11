FROM golang:1.22-alpine

ARG ID=501

RUN apk add --no-cache mysql-client zip bash postgresql-client mongodb-tools

RUN addgroup -g $ID devops && adduser -u $ID -G devops -D devops

WORKDIR /home/devops

USER devops

CMD [ "tail", "-f", "/dev/null" ]