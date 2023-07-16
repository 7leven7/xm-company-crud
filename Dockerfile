# Start from golang base image
FROM golang:alpine

ENV CGO_ENABLED=1
ENV GIN_MODE=release

RUN apk update && apk add --no-cache git bash build-base gcc musl-dev lvm2-dev gpgme-dev btrfs-progs-dev

RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
RUN ./bin/golangci-lint --version

# Setup folders
RUN mkdir /app
WORKDIR /app

# Copy the source from the current directory to the working Directory inside the container
COPY . .

# Download all the dependencies
RUN go get -d -v ./...

CMD bash -c "sleep 10 && go run main.go"
