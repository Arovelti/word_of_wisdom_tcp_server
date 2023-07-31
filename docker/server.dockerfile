FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd/server/main.go


FROM alpine:latest
COPY --from=builder /app/app /app/app
WORKDIR /app

CMD ["./app"]