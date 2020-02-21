FROM golang:1.13.6-alpine

RUN apk add --no-cache gcc musl-dev

WORKDIR /go/src/app
COPY go.mod go.sum ./
RUN go get -d -v ./...

COPY . .
RUN go install -v ./...

CMD ["auth", "/data/auth.toml"]
