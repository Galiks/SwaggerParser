# syntax=docker/dockerfile:1
FROM golang:1.18.1-alpine
ADD . /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *go ./
RUN go build -o swagger-parser .

RUN apt-get update && apt-get install -y wkhtmltopdf && rm -r /var/lib/apt/lists/*

EXPOSE 8001
CMD ["/docker-gs-ping"]