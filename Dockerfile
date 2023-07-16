# Start from golang base image
FROM golang:alpine as builder


# Install git. (alpine image does not have git in it)
RUN apk update && apk add --no-cache git

# Set current working directory
WORKDIR /go/src/job_posting_retreiver

COPY . .

RUN go get -u -v

# Build the application.
RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/

# Finally our multi-stage to build a small image
# Start a new stage from scratch
#FROM scratch
#FROM alpine:latest

# Copy the Pre-built binary file
#COPY --from=builder /vanir/bin /go/bin

#Expose port and start application
EXPOSE 8080

# Run executable
#CMD ["./vanir"]
ENTRYPOINT /go/bin/job_posting_retreiver