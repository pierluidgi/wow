FROM golang:1.18
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o wow-server server.go
EXPOSE 9011/tcp
CMD ["/app/wow-server", "-config", "./resources/server-config.yaml"]