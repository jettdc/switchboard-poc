FROM golang:alpine as builder

WORKDIR /
ADD delivery.go .

RUN go mod init "github.com/jettdc/switchboard/examples/delivery-service"
RUN go mod tidy
RUN go build -o main .

EXPOSE 12345

CMD ["./main"]