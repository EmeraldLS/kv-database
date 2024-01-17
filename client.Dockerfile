FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY /client .

RUN go build -o client .

EXPOSE 3031

CMD ["./app/client"]