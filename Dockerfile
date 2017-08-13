FROM golang:1.8

WORKDIR /go/src/app
COPY myproxy.go  .

RUN go-wrapper download   # "go get -d -v ./..."
RUN go-wrapper install    # "go install -v ./..."

ENTRYPOINT ["go-wrapper", "run"] # ["app"]
