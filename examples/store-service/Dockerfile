FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

FROM alpine
COPY --from=builder /build/main /store-service/
WORKDIR /store-service
EXPOSE 54321
CMD ["./main"]