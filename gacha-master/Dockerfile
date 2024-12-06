# Stage 1: Build the Go binary
FROM golang:1.21.0 AS build

# Set the working directory inside the container
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Create the final image
FROM alpine:latest AS release

# Create a non-root user and group
RUN addgroup -S golanggroup && adduser -S golanguser -G golanggroup

WORKDIR /app

COPY --from=build /app/main .

RUN apk --no-cache add ca-certificates tzdata

RUN chown -R golanguser:golanggroup /app && chmod +x /app/main
USER golanguser

ENV SERVER_PORT=8001

EXPOSE ${SERVER_PORT}

ENTRYPOINT ["/app/main"]
