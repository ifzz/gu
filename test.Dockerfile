FROM golang:latest


# Run go get
RUN go get -u -v github.com/gu-io/gu/...


# Install gometalinter
RUN go get -u -v github.com/alecthomas/gometalinter

# Install missing lint tools
RUN gometalinter --install

# Switch to gu-io/gu directory
WORKDIR /go/src/github.com/gu-io/gu

# Run go tests
RUN go test -v ./drivers/...
RUN go test -v ./eventx/...
RUN go test -v ./notifications/...
RUN go test -v ./parsers/...
RUN go test -v ./router/...
RUN go test -v ./shell/...
RUN go test -v ./tests/...
RUN go test -v ./utils/...

# Run go linters
RUN gometalinter --deadline 4m --errors  .
RUN gometalinter --deadline 4m --errors  ./drivers/...
RUN gometalinter --deadline 4m --errors  ./eventx/...
RUN gometalinter --deadline 4m --errors  ./notifications/...
RUN gometalinter --deadline 4m --errors  ./parsers/...
RUN gometalinter --deadline 4m --errors  ./router/...
RUN gometalinter --deadline 4m --errors  ./shell/...
RUN gometalinter --deadline 4m --errors  ./utils/...
