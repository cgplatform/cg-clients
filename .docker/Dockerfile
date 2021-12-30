FROM ubuntu:latest

EXPOSE 8080

WORKDIR dist

COPY .env s2p-api ./

CMD chmod +x s2p-api

CMD ./s2p-api
