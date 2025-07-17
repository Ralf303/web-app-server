FROM golang:1.22.1

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o server ./cmd/main

CMD ["./server"]