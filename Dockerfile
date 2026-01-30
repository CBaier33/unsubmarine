FROM golang:1.25-alpine AS builder 

WORKDIR /src 

COPY go.mod go.sum* ./ 
RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/unsubmarine .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/unsubmarine /app/unsubmarine

COPY unsubmarine.html /app/unsubmarine.html

RUN adduser -D appuser 
USER appuser 

EXPOSE 8080
CMD ["/app/unsubmarine"]
