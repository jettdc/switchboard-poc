FROM golang:alpine as builder
RUN apk add build-base
RUN mkdir /build
ADD . /build/
WORKDIR /build
RUN go build -o main .

FROM alpine
RUN adduser -S -D -H -h /switchboard appuser
USER appuser
COPY --from=builder /build/main /switchboard/
WORKDIR /switchboard
EXPOSE 8080
CMD ["./main"]