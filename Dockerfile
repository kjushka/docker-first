FROM golang:1.19-alpine
WORKDIR /docker-first
COPY / ./
RUN go mod download
RUN go build -o ./docker-first ./cmd/docker-first

EXPOSE 80

CMD [ "./docker-first" ]