FROM golang:1.25-alpine3.23 AS builder

WORKDIR /src

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o external-scaler .

FROM alpine:3.23

WORKDIR /

COPY --from=builder /src/external-scaler .

EXPOSE 6000

ENTRYPOINT ["/external-scaler"]
