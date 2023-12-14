# This directive instructs the Docker builder what syntax to use when parsing the Dockerfile. 
# syntax=docker/dockerfile:1

# Based on official golang alpine image
FROM golang:1.21-alpine as build

# Create a directory inside the build image and set it as default destination for all subsequent commands.
WORKDIR /app

# Copy only necessary files and directories into the build image
COPY internal/ internal/
COPY go.mod go.sum main.go ./

# Compile our application, create the binary and save it under the root path as 'indexer'
# CGO_ENABLED=0 to disable dynamically links for a few dependencies
# Remove DWARF tables and annotations with -s -w. The DWARF tables are needed for debuggers and the annotations for stack traces.
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /indexer .

# Now copy our binary program inside our debin production image
FROM debian:bullseye-slim 

# Install 7z needed for extract dblp archive
RUN apt-get update && apt-get install -y p7zip-full ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=build /indexer /indexer

# Set the entrypoint as the binary file
ENTRYPOINT ["/indexer"]
