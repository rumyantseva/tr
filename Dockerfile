FROM golang:1.13 AS modules

ENV GOPROXY https://gocenter.io

ADD go.mod go.sum /m/
RUN cd /m && go mod download

FROM golang:1.13 AS build

COPY --from=modules /go/pkg/mod /go/pkg/mod

RUN mkdir -p /rumyantseva/tr
ADD . /rumyantseva/tr
WORKDIR /rumyantseva/tr

RUN GOOS=linux GOARCH=amd64 make build

FROM scratch

COPY --from=build /rumyantseva/tr/bin/linux-amd64/tr /tr

ENTRYPOINT ["/tr"]
