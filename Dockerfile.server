# Use an official Golang runtime as a base image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the server application code into the container
COPY . .

# Build the Golang server application
RUN go build -o recurb-api ./cmd/recurb-api

# Expose the port your Golang server will run on
EXPOSE 8080

# Define the command to run your Golang server
CMD ["./recurb-api"]
