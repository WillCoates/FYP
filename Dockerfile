FROM golang:1.13.6-alpine

RUN apk add --no-cache gcc musl-dev tzdata
ENV TZ Europe/London

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
