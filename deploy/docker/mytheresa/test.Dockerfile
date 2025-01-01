FROM golang:1.23.4 

ENV GO111MODULE=on
ENV CGO_ENABLED=1

ADD . /app

WORKDIR /app

RUN go install github.com/onsi/ginkgo/v2/ginkgo@latest 

CMD ["ginkgo", "-r"]

