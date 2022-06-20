FROM golang:1.15.5-alpine3.12 as builder

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /src

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build  -o /go/bin/mono_bot .

FROM alpine:3.12

RUN apk update
RUN apk add tzdata

COPY --from=builder /go/bin/mono_bot /bin/mono_bot

ENTRYPOINT sh -c "/bin/mono_bot --token=$TOKEN --timeout=$TIMEOUT --redis_url=$REDIS_URL"