FROM golang:1.21 AS builder

WORKDIR $GOPATH/src/main

COPY ./pkg ./pkg

COPY go.mod .
COPY go.sum .

COPY main.go .

RUN go get -d -v 

RUN CGO_ENABLED=0 go build -o /mongomonitor main.go

FROM gcr.io/distroless/base:latest

COPY --from=builder /mongomonitor /mongomonitor

CMD [ "/mongomonitor" ]