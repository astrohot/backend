# Get the proper golang image to build the app.
FROM golang:1.12-alpine as builder

# Add non-root user to be used in the second step.
RUN adduser -D scratchuser

# Get the proper dependencies.
RUN apk update && apk add git ca-certificates

# Set the current working directory inside the container.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod go.sum ./

# Install all the dependencies.
RUN go mod download

# Copy the source code.
COPY . .

# Run go:generate to generate code using gqlgen if needed.
RUN CGO_ENABLED=0 go generate ./internal/api/...

# Build the go app.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build cmd/api/main.go

# Use a non-root user.
USER scratchuser

# CMD to run the executable.
CMD ["./main"]
