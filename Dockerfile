FROM golang as builder

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/Kalachevskyi/mono-chat

# Changing working directory.
WORKDIR /go/src/github.com/Kalachevskyi/mono-chat

# Download dependencies
RUN export GO111MODULE=on && go mod download && go mod vendor

# Building application.
RUN go build -o mono-chat main.go

######### Start a new stage from alpine #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN apk add tzdata

# Changing working directory.
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/github.com/Kalachevskyi/mono-chat/mono-chat .

# Run service
ENTRYPOINT [ "sh", "-c", "./mono-chat --token=$TOKEN --timeout=$ TIMEOUT --redis_url=$REDIS_URL" ]