FROM golang:1.18
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o wow-client client.go
CMD ["/app/wow-client", "-config", "./resources/client-config.yaml"]