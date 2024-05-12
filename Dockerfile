# Official Golang runtime as a base image
FROM golang:1.22.2 AS builder

# Set the working directory in the container
WORKDIR /service-catalogue

# Copy the local package files to the container's workspace
COPY go.mod go.sum ./

RUN go mod download

COPY . .

# Install PostgreSQL client tools
RUN apt-get update && apt-get install -y postgresql-client

RUN go get -u github.com/pressly/goose/cmd/goose
# RUN goose -dir database/migrations postgres "postgres://postgres:postgres@db:5432/catalogue_service?sslmode=disable" up

ENV PATH=$PATH:/go/bin

# Build the Go app
RUN go build -o service-catalogue .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./service-catalogue"]
