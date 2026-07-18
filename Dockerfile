FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o main ./cmd

FROM alpine

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/docs ./docs

RUN chmod +x ./main

EXPOSE 3382

CMD ["./main"]