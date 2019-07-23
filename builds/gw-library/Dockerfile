FROM alpine:latest

RUN apk add --no-cache \
        libc6-compat

RUN mkdir -p /app
WORKDIR /app

ADD gw-library /app
EXPOSE 8081
CMD ["./gw-library"]