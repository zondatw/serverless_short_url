FROM golang:1.15-alpine as builder

WORKDIR /go/app
ADD . .

# fix $GOPATH/go.mod exists but should not
ENV GOPATH=""

ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

EXPOSE 5000


FROM debian

ENV GIN_MODE=release

RUN apt update && apt upgrade -y
RUN apt install -y ca-certificates

WORKDIR /usr/local/
COPY --from=builder /go/app/statisticsserver /usr/local/statisticsserver

CMD ["./statisticsserver"]