# Stage 1: Build
FROM golang:1.20-alpine AS build

# Install tzdata for Timezone related stuff
RUN apk add --no-cache tzdata

# Set working directory
WORKDIR /app

# Copy source code into the container
COPY . /app

# Download dependencies
## Doing this since vendor is not commited for this project
## To have reproducible builds, vendor would be committed
## & therefore this won't be needed
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bite ./cmd/

# Stage 2: Run
FROM alpine:latest

# Install tzdata for Timezone related stuff
RUN apk add --no-cache tzdata

# Set working directory
WORKDIR /app/

# Copy the pre-built binary from the build stage
COPY --from=build /app/bite /app/bite
COPY --from=build /app/config /app/config

# # Create a new user (good practice)
# RUN groupadd -r appgroup && useradd -r -g appgroup appuser

# Expose ports
EXPOSE 8088

# Run the binary
CMD ["/app/bite"]