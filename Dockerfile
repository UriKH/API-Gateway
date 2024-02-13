# syntax=docker/dockerfile:1

FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Copy the source code
COPY . .

# Download Go modules
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /api-gateway -buildvcs=false

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
EXPOSE 8080

# Run
CMD ["/api-gateway"]