### Stage 1: Build Go Application ###
FROM golang:latest AS build

WORKDIR /app

# Copy Go source code to the container
COPY . .

# Build the Go application
RUN go build -o platform

### Stage 2: Final Image ###
FROM postgres:latest

# Set environment variables for PostgreSQL
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD postgres
ENV POSTGRES_DB platform

# Expose PostgreSQL port
EXPOSE 5432

# Copy the built Go binary from the first stage
COPY --from=build /app/platform /usr/local/bin/platform

# Start the Go application
CMD ["platform"]
