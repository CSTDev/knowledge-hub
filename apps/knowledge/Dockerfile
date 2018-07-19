FROM golang:1.9.7 as builder

ADD . /go/src/github.com/cstdev/knowledge-hub/apps/knowledge/
RUN go get -t -v ./...
RUN go test -v ./...
WORKDIR /go/src/github.com/cstdev/knowledge-hub/apps/knowledge/
RUN go build ./cmd/server/main.go


FROM alpine:3.7
WORKDIR /root/
COPY --from=builder /go/src/github.com/cstdev/knowledge-hub/apps/knowledge/main .
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 
ENV PORT=8000
CMD ["./main"]

EXPOSE 8000