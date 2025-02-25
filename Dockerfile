# syntax=docker/dockerfile:1
# https://dev.to/heroku/deploying-your-first-golang-webapp-11b3
# FROM jrottenberg/ffmpeg:4.4-alpine as ffmpeg
# COPY --from=ffmpeg /usr/local /usr/local
FROM golang:1.16-alpine AS builder
RUN mkdir /build
ADD go.mod go.sum main.go /build/
WORKDIR /build
COPY . .
RUN chmod +x /build
RUN go build

FROM alpine
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/kithli-api /app/
WORKDIR /app
CMD ["./kithli-api"]

# FROM golang:1.16-alpine AS builder
# RUN mkdir /build
# ADD go.mod go.sum main.go /build/
# WORKDIR /build
# COPY . .
# RUN chmod +x /build
# RUN go build

# # FROM alpine
# # RUN apk add --no-cache ca-certificates && update-ca-certificates
# # RUN adduser -S -D -H -h /app appuser
# # USER appuser
# # COPY --from=builder /build/kithli-api /usr/bin/kithli-api
# # EXPOSE 5000 5020
# # WORKDIR /app
# # # CMD ["./kithli-api"]
# # ENTRYPOINT [ "/usr/bin/kithli-api" ]