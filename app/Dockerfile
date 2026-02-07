# Stage 1
FROM golang:1.25.7-alpine3.23 AS builder

WORKDIR /app
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o testecargak6 .

# Stage 2
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/testecargak6 ./

CMD ["./testecargak6"]
