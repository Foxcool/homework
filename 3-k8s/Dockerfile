FROM golang:1.13-alpine as builder

RUN apk add --no-cache --update \
	git

ENV _PATH=github.com/foxcool/homework
ADD . /go/src/$_PATH
WORKDIR /go/src/$_PATH

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X main.version=develop" -o /homework .



FROM alpine:3.11
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /homework .
USER 997

CMD "/app/homework"

EXPOSE 8000