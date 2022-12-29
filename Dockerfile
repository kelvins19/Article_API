FROM golang:latest

COPY . /app

WORKDIR /app

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go get -d -v

RUN go build -o main .

EXPOSE 8080

CMD ["./main"]