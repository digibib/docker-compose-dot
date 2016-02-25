FROM alpine:3.3
WORKDIR /app
RUN apk add --no-cache openssl ca-certificates
ADD ./build/app /app

CMD ["/app/app"]
