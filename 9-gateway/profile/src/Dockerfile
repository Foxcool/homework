FROM golang:1.13-alpine AS builder

# Install git
RUN apk add --no-cache --update \
	git

ENV GO111MODULE=off

# Set enviroment variable for Go
ENV GOPATH=/go \
	PATH="/go/bin:$PATH"

# Copy the local package files to the container's workspace. Add to GOPATH
ADD . /go/src/github.com/foxcool/homework/7-prometheus/app

# Set workdir
WORKDIR /go/src/github.com/foxcool/homework/7-prometheus/app

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=develop" -o /go/bin/homework ./cmd/homework


# Init new lightweight container
FROM alpine:3.11
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/bin/homework .

# Run the compiled bin by default when the container start.
CMD /app/homework

# Service listens on port 80.
EXPOSE 80
