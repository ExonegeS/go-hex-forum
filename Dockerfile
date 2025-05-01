FROM golang:1.24

WORKDIR /app

COPY  . .

RUN go build -o go-hex-forum ./cmd/myapp

CMD ["./go-hex-forum"]