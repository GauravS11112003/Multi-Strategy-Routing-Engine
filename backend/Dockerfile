FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/worker cmd/worker/main.go

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /bin/api /bin/api
COPY --from=builder /bin/worker /bin/worker

EXPOSE 8080

CMD ["/bin/api"]
