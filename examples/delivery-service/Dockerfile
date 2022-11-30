FROM golang:alpine as builder
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

FROM alpine
COPY --from=builder /build/main /delivery-service/
WORKDIR /delivery-service
EXPOSE 12345
CMD ["./main"]