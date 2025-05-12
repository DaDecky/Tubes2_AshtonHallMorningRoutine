# backend/Dockerfile
FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o server .

EXPOSE 8081

CMD ["./server"]
