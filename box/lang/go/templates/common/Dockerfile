FROM golang:latest

LABEL maintainer="readyGo team. JitenP@Outlook.Com"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE {{$.config.APISpec.Port}}

CMD ["./main"]
