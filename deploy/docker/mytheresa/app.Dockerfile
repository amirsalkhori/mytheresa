# Use an official Go image as the base image
FROM golang:1.23.4 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the application code
COPY . .

# Build the Go application
RUN echo "start build"
WORKDIR /app/cmd/mytheresa
RUN go build -o /app/main .
WORKDIR /app/cmd/performance
RUN go build -o /app/performance .
WORKDIR /app/cmd/seeder
RUN go build -o /app/seeder .
RUN GOBIN=/app/go go install github.com/onsi/ginkgo/v2/ginkgo@latest  # Run Ginkgo tests with coverage

# Copy migration files to runtime directory
COPY ./migrations /app/migrations

# Verify the binary is created
RUN ls -la $BIN
RUN echo "Build successful"

# Use a smaller image for the runtime
FROM debian:bookworm-slim

# Set the working directory for the runtime container
WORKDIR /app



# Copy the compiled binary and migrations from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/performance .
COPY --from=builder /app/seeder .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/go/ginkgo /usr/local/bin

# Expose the port the application will run on
EXPOSE 8080

# Command to run the application
CMD ["./main"]
