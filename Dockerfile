# syntax=docker/dockerfile:1
FROM golang:1.18.1-alpine
RUN apk add git
ADD . /app
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *go ./

RUN apk add --no-cache \
        wkhtmltopdf \
        xvfb \
        ttf-dejavu ttf-droid ttf-freefont ttf-liberation \
    ;

RUN ln -s /usr/bin/wkhtmltopdf /usr/local/bin/wkhtmltopdf;
RUN chmod +x /usr/local/bin/wkhtmltopdf;

RUN go build -o swagger-parser .

EXPOSE 8001
CMD ["/docker-gs-ping"]