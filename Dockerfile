FROM golang:1.12-alpine AS builder

RUN apk add git

WORKDIR /opt/src
ADD . .
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -v -o toggly-app .

FROM alpine:3.9

RUN apk --no-cache add ca-certificates

WORKDIR /opt

COPY --from=builder /opt/src/toggly-app /opt/src/configs ./

RUN chmod +x toggly-app

ENTRYPOINT "./toggly-app"