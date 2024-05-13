FROM golang:1.22.3-alpine3.19 as builder

ARG SERVICE

WORKDIR /build

COPY services/${SERVICE} .
COPY lib/ /lib

RUN go mod download && \
    go build -tags timetzdata -o ${SERVICE}

FROM alpine:3.19.1

ARG SERVICE
ENV SERVICE ${SERVICE}

WORKDIR /app

COPY --from=builder /build/${SERVICE} ${SERVICE}

CMD ["/bin/sh", "-c", "/app/${SERVICE}"]