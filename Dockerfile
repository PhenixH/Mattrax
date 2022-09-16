FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/github.com/mattrax/Mattrax/
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/mattrax ./cmd/mattrax

FROM scratch
COPY --from=builder /go/bin/mattrax /go/bin/mattrax
ENTRYPOINT ["/go/bin/mattrax"]