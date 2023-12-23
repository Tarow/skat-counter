FROM golang:1.21.5-alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY main.go .
COPY internal/ internal/
COPY static/ static/
COPY templates/ templates/

RUN go build -o skatcounter

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/skatcounter /usr/bin/skatcounter

CMD ["skatcounter"]

