FROM ubuntu:latest

EXPOSE 8080

WORKDIR dist

RUN apt-get update

RUN apt-get install -y ca-certificates

COPY .env s2p-api ./

CMD chmod +x s2p-api

CMD ./s2p-api
