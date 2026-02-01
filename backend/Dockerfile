FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/api

EXPOSE 8080

CMD ["go", "run", "."]
