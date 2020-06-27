FROM golang:1.14.4

WORKDIR /go/src/app
COPY . /go/src/app
COPY ./config.test.yml /go/src/app/config.yml

RUN go get -v -t -d ./...
RUN go build -v .

RUN openssl genrsa -out test_rsa 1024
RUN openssl rsa -in test_rsa -pubout > test_rsa.pub

CMD go test -v -coverprofile=./reports/coverage.txt -covermode=atomic -coverpkg=./... ./...
