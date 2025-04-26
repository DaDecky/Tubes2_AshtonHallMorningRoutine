# use official Golang image
FROM golang:1.16.3-alpine3.13

# set working directory
WORKDIR /backend-server

# Copy the source code
COPY . .

# Download and install the dependencies
RUN go get -d -v ./...

# Build the Go app
RUN go build -o backend-server .

#EXPOSE the port
EXPOSE 8000

# Run the executable
CMD ["./backend-server"]

# betak dari website, gk tau work apa kgk. gw masih manual