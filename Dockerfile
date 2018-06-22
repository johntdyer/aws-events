FROM golang:1.8.3 as builder
WORKDIR /go/src/github.com/johntdyer/aws-events
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
ADD . /go/src/github.com/johntdyer/aws-events/

RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o packages/aws-events .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
RUN mkdir /root/database
# COPY config.toml  .
COPY --from=builder /go/src/github.com/johntdyer/aws-events/packages/aws-events .
CMD ["./aws-events"]

