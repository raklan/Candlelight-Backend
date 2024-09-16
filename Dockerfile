FROM golang:1.22

RUN mkdir /candlelight-ruleengine
RUN mkdir /candlelight-models

RUN cd ..

COPY /candlelight-ruleengine/ /candlelight-ruleengine
COPY /candlelight-models/ /candlelight-models

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /myapp

USER 1000

CMD ["/myapp"]