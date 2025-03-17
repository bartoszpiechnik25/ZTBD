FROM golang:1.24-bookworm AS builder

WORKDIR /app

ARG BEGIN_YEAR
ARG END_YEAR
ARG VOLUME_DIR

ENV BEGIN_YEAR=${BEGIN_YEAR}
ENV END_YEAR=${END_YEAR}
ENV VOLUME_DIR=${VOLUME_DIR}

COPY go.mod ./
RUN go mod download

COPY . /app

RUN go build -o /initializer

FROM debian:12-slim 

WORKDIR /app

RUN apt update && apt install ca-certificates -y

COPY --from=builder /initializer /app/initializer

RUN mkdir /app/data

ENTRYPOINT [ "/app/initializer" ]
