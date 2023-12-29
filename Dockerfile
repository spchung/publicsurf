FROM golang:1.21.4 AS builder

WORKDIR /go/src/app

# Install BIMG and reqs
RUN apt-get update
RUN apt-get install -y libvips libvips-dev

# Copy the local package files to the container's workspace
COPY . .

# Install any dependencies if needed
RUN go mod download
RUN go get -u gopkg.in/h2non/bimg.v1

# Build the Go binary0
RUN go build -o backend cmd/gin/main.go

FROM golang:1.21.4

RUN apt-get update
RUN apt-get install -y libvips libvips-dev

# Expose the port on which the service will run
EXPOSE 8080

WORKDIR /go/src/app
# copy build file
COPY --from=builder /go/src/app/configurations ./configurations
COPY --from=builder /go/src/app/backend .

# Command to run your application
CMD ["./backend"]