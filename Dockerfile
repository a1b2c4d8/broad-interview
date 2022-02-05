FROM golang:1.17

WORKDIR /go/src/github.com/a1b2c4d8/broad-interview
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go get github.com/onsi/ginkgo/v2/ginkgo
RUN go get github.com/onsi/gomega/...
RUN go get github.com/onsi/ginkgo/v2/internal@v2.1.1
RUN make clean build
RUN mv _output/main /usr/local/bin/broad-interview

CMD [ "/usr/local/bin/broad-interview" ]
