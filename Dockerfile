FROM golang:1.19.7-alpine3.17

WORKDIR src/
COPY go.* .
RUN go mod download

COPY . .
RUN go build -o klaxn-api main.go

EXPOSE 8080
ENTRYPOINT ["./klaxn-api"]
