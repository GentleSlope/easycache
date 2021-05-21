FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/easycache
COPY . $GOPATH/src/easycache
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./go-gin-example"]
