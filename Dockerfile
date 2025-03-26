FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o main .

EXPOSE 8080

ENV PORT=8080
CMD ["./main"] 