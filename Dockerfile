FROM golang:1.17.3-alpine

ENV listenPort=2113
ENV remootioIP=192.168.1.1:8080
ENV scrapeInterval=30

WORKDIR /app

COPY * /app/

RUN go mod download
RUN go build -o /coinspot-prom

EXPOSE $listenPort

CMD [ "/coinspot-prom" ]
