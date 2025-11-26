FROM golang:1.23.3-alpine3.20 AS build

WORKDIR /app

# Modules layer
COPY go.mod go.sum ./
RUN go mod download

# Build layer
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /noway ./cmd/app

FROM alpine:3.20 AS run

COPY --from=build /noway /noway
COPY --from=build /app/.env /.env

EXPOSE 8080

CMD ["/noway"]