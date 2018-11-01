FROM golang:latest

WORKDIR /app

ADD . /app

RUN chmod 600 /app/src/id_rsa.pub

RUN chmod 700 /app/src/id_rsa

EXPOSE 8181

ENV GOPATH=/app

CMD [ "go","run","server.go" ]