
FROM golang:1.18

WORKDIR /app

COPY go.mod .
COPY go.sum .


RUN go get github.com/gofiber/fiber/v2

RUN go mod download



COPY . .

RUN go build -o main .


EXPOSE 3000

CMD ["./main"]
