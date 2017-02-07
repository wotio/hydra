FROM golang:1.7

RUN go get github.com/Masterminds/glide
WORKDIR /go/src/github.com/ory/hydra

ADD ./glide.yaml ./glide.yaml
ADD ./glide.lock ./glide.lock
RUN glide install

ADD . .
RUN go install .

ENTRYPOINT /go/bin/hydra host

EXPOSE 4444