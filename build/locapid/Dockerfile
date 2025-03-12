# Stage 1
FROM golang:1.24-alpine3.21 AS gobuilder

WORKDIR /app/locapid/
COPY . .
ENV CGO_ENABLED=1
ENV GOFLAGS=-mod=vendor
ENV GOOS=linux
ENV GOARCH=amd64
RUN apk -U add ca-certificates
RUN apk update && apk add pkgconf git bash build-base sudo
RUN go mod download && go mod vendor
RUN go build -tags musl -o locapid ./cmd/locapid/main.go

# Stage 2
FROM alpine:3.21
WORKDIR /app
COPY --from=gobuilder /app/locapid/locapid .
USER 1001
CMD [ "/app/locapid" ]
