FROM golang:1.22.0-alpine3.19 as builder

RUN mkdir "$GOPATH/src/outbox"
WORKDIR "$GOPATH/src/outbox"
COPY . .

RUN go mod download && go mod verify && go mod tidy

RUN go build -o /go/bin/outbox

FROM golang:1.22.0-alpine3.19 as runner

RUN apk add --no-cache --no-interactive curl

COPY --from=builder /go/bin/outbox /go/bin/outbox

CMD ["/go/bin/outbox"]
