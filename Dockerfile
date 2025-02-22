FROM golang:1.24-alpine AS builder

WORKDIR /app/

COPY . .

RUN go mod download
RUN go mod verify

RUN CGO_ENABLED=0 go build -o /medguardian/app ./cmd/medguardian/

FROM scratch

COPY --from=builder /medguardian /medguardian
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/medguardian/app" ]