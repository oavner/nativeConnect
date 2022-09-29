FROM golang:1.16-alpine
WORKDIR /app
COPY ./image/go.mod ./
RUN go mod download

COPY ./image/*.go ./
RUN go build -o ./connectionsLogger
CMD [ "./connectionsLogger" ]
