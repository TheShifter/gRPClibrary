FROM alpine:latest

RUN apk add --no-cache \
        libc6-compat

RUN mkdir -p /app
WORKDIR /app

ADD svc-books /app
EXPOSE 9091
CMD ["./svc-books"]