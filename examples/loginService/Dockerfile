FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

FROM alpine
COPY --from=builder /build/main /auth-service/
WORKDIR /auth-service
EXPOSE 8081
CMD ["./main"]