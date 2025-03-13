FROM golang:1.24-bookworm AS builder

WORKDIR /app

ARG BEGIN_YEAR
ARG END_YEAR

ENV BEGIN_YEAR=${BEGIN_YEAR}
ENV END_YEAR=${END_YEAR}

COPY go.mod ./
RUN go mod download

COPY . /app

RUN go build -o /initializer

FROM debian:12-slim 

WORKDIR /

COPY --from=builder /initializer /initializer

ENTRYPOINT [ "/initializer" ]
