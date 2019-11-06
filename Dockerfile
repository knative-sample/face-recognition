# Build the manager binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/knative-sample/face-recognition
COPY pkg/    pkg/
COPY config/   config/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o face-recognition github.com/knative-sample/face-recognition/cmd

# Copy the face-recognition into a thin image
FROM alpine:3.7
WORKDIR /
COPY --from=builder /go/src/github.com/knative-sample/face-recognition/face-recognition app/
COPY --from=builder /go/src/github.com/knative-sample/face-recognition/config app/config
RUN apk upgrade && apk add --no-cache ca-certificates
#ENTRYPOINT ["/app/face-recognition"]
