FROM golang:1.19-alpine
WORKDIR /docker-first
COPY . /docker-first
RUN go mod download
RUN go build -o ./docker-first ./main.go

EXPOSE 80

CMD [ "./docker-first" ]