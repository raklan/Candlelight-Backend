FROM golang:1.22 AS builder

RUN mkdir /logs
WORKDIR /app

RUN mkdir ./candlelight-ruleengine
RUN mkdir ./candlelight-models
RUN mkdir ./candlelight-api

COPY ./candlelight-ruleengine/ ./candlelight-ruleengine
COPY ./candlelight-models/ ./candlelight-models
COPY ./candlelight-api/ ./candlelight-api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C ./candlelight-api -o myapp.exe

FROM scratch
COPY --from=builder /app/candlelight-api/myapp.exe /app/
COPY --from=builder /logs /logs
EXPOSE 10000

ENTRYPOINT ["/app/myapp.exe"]
