FROM alpine:latest

RUN apk add --no-cache \
        libc6-compat

RUN mkdir -p /app
WORKDIR /app

ADD svc-library /app
EXPOSE 9090
CMD ["./svc-library"]