# Stage 1: Build the application binary
FROM golang:1.26.0-alpine3.23 as build

WORKDIR /app

# Install build dependencies for CGO-dependent packages
RUN apk update && apk add build-base

# Leverage Docker layer caching by downloading dependencies separately
COPY go.mod go.sum ./
RUN go mod download

# Build the production-ready binary
COPY . .
RUN go build -o nexus-tasks

# Stage 2: Create a lightweight runtime image
FROM alpine:3.23

WORKDIR /app

# Copy the compiled binary and static assets from the build stage
COPY --from=build /app/nexus-tasks .
COPY --from=build /app/database ./database
COPY --from=build /app/assets ./assets
