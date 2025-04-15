
# # FROM gcr.io/distroless/base:latest


# # RUN go mod tidy
# # RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go

# # WORKDIR /app
# # COPY action-app /app/action-app

# # CMD ["/app/action-app"]

# FROM golang:1.24.1-alpine AS GOLANG

# # LABEL maintainer=Cloudbees-pod-7 \
# #     email=engineering@cloudbees.io

# RUN apk "upgrade" libssl3 libcrypto3

# RUN go mod tidy
# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go


# WORKDIR /app
# COPY . .
# RUN go mod tidy
# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go

# # FROM gcr.io/distroless/base:latest

# COPY --from=builder /app/action-app /app/action-app

# # Define the entrypoint to run the app
# ENTRYPOINT ["/app/action-app"]

# # COPY jenkins_actions_app /app

# # CMD ["/app/jenkins_actions_app"]

# # ENTRYPOINT [ "go", "run" ,"main.go" ]


# Use the official Go base image for building
FROM golang:1.24 as builder

# Set environment for Go modules
WORKDIR /app
COPY . .

# Build the Go application (static binary for portability)
RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go

FROM gcr.io/distroless/base:latest

COPY --from=builder /app/action-app /app/action-app

# Define the entrypoint to run the app
ENTRYPOINT ["/app/action-app"]
