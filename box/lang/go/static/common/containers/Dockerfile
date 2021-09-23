FROM golang:1.17.0-bullseye

LABEL maintainer="readyGo team. JitenP@Outlook.Com"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]
