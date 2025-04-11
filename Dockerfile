# # Use the official Go base image for building
# FROM golang:1.24 as builder

# # Set environment for Go modules
# WORKDIR /app
# COPY . .

# # Build the Go application (static binary for portability)
# RUN go mod tidy
# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go

# FROM gcr.io/distroless/base:latest

# COPY --from=builder /app/action-app /app/action-app

# # Define the entrypoint to run the app
# ENTRYPOINT ["/app/action-app"]

# FROM debian:bookworm-slim AS builder

# # Install git and necessary certificates
# RUN apt-get update && apt-get install -y \
#     git \
#     ca-certificates \
#     && apt-get clean && rm -rf /var/lib/apt/lists/*

# # Use 'base:nonroot' instead of 'static:nonroot' for compatibility
# FROM gcr.io/distroless/base:nonroot

# # LABEL maintainer=Cloudbees \
# #     email=engineering@cloudbees.iotw

# # Copy Git binary and dependencies
# COPY --from=builder /usr/bin/git /usr/bin/git
# COPY --from=builder /lib/ /lib/
# COPY --from=builder /usr/lib/ /usr/lib/
# COPY --from=builder /etc/ssl/certs/ /etc/ssl/certs/
# # Copy the service binary
# COPY external-ci-service /

# ENTRYPOINT ["/external-ci-service"]

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY gha_run_cbp_workflow_app /app/gha_run_cbp_workflow_app

CMD ["/app/gha_run_cbp_workflow_app"]