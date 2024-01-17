FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY /server .

RUN go build -o server .

EXPOSE 3031

CMD ["./app/server"]