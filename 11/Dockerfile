FROM golang:1.17 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o ingress-manager main.go

FROM alpine:3.15.3

WORKDIR /app

COPY --from=builder /app/ingress-manager .

CMD ["./ingress-manager"]