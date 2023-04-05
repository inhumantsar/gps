# Use an official Go runtime as a parent image
FROM golang:1.17-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any necessary dependencies
RUN apk update && apk add --no-cache git && go get -u github.com/stretchr/testify

# Create a non-root user with the UID 1000
RUN adduser -D -u 1000 golang

# Set the user to the non-root user
USER golang

# Set the owner and permissions of the Go module cache
RUN mkdir -p /home/golang/go/{pkg,src,bin} && \
    chown -R golang:golang /home/golang/go && \
    chmod -R 775 /home/golang/go

# Set the GOPATH to the non-root user's home directory
ENV GOPATH /home/golang/go

# Build the executable and run the tests
CMD ["go", "mod", "tidy", "&&", "go", "test", "./..."]
