# Build the source code with Lang container
FROM golang:alpine as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /main

# Copy the exec file to a smaller base image
FROM alpine
COPY --from=build /main /main
CMD /main
