FROM golang:latest as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="-w -s" main.go

FROM scratch AS runner

WORKDIR /app

COPY --from=builder /app/main .

ENTRYPOINT ["./main"]