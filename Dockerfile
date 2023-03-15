FROM alpine:3.17 as final
FROM golang:1.19.7-alpine3.17 as build

WORKDIR src/
COPY go.* .
RUN go mod download

COPY . .
RUN go build -o klaxn-api main.go

FROM final
COPY --from=build /go/src/klaxn-api klaxn-api
EXPOSE 8080
ENTRYPOINT ["./klaxn-api"]
