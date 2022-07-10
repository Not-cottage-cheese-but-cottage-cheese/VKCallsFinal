FROM golang:1.18.3-alpine3.16

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
COPY .env ./
RUN go mod download

ADD server ./server
COPY cmd/10+20/main.go ./

RUN go build -o bot .

CMD ["./bot"]