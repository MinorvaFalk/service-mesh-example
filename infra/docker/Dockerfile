FROM golang:1.22.5-alpine AS base

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .


FROM base as producer
RUN CGO_ENABLED=0 GOOS=linux go build -o producer ./cmd/producer/main.go


FROM base as consumer
RUN CGO_ENABLED=0 GOOS=linux go build -o consumer ./cmd/consumer/main.go

FROM alpine:3.19 as final

RUN apk add --no-cache tzdata words-en
RUN ln -s /usr/share/dict/american-english /usr/share/dict/words

ENV TZ=Asia/Jakarta

COPY --from=producer /build/producer /usr/local/bin
COPY --from=consumer /build/consumer /usr/local/bin

RUN ln -s /usr/local/bin/producer / \
    && ln -s /usr/local/bin/producer /bin/

RUN ln -s /usr/local/bin/consumer / \
    && ln -s /usr/local/bin/consumer /bin/

RUN adduser go --disabled-password --home /home/go/
USER go