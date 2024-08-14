# Base image.
FROM golang:1.22.6-alpine as builder

# Declare Environment variable.
ENV APP_ID=7f2f0a47c0554a0d930074f826599ade

# Set the working directory.
WORKDIR /app

# Copy the Go module files.
COPY go.mod go.sum ./

#  Download the dependencies.
RUN go mod download

# Copy otcher app files.
COPY . .

# Buid executable file.
RUN go build -o main .

# Create a minimal final image.
FROM alpine:3.20

# Set the working directory.
WORKDIR /app

# Copy the compiled binary from the builder stage.
COPY --from=builder /app/main .

# Specify the command to run the binary.
CMD ["./main"]

# Describe which ports application is listening on.
EXPOSE 8080

