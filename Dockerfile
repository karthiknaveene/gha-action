
# FROM gcr.io/distroless/base:latest


# RUN go mod tidy
# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go

# WORKDIR /app
# COPY action-app /app/action-app

# CMD ["/app/action-app"]

FROM golang:1.24.1-alpine AS GOLANG

# LABEL maintainer=Cloudbees-pod-7 \
#     email=engineering@cloudbees.io

RUN apk "upgrade" libssl3 libcrypto3

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go


WORKDIR /app
COPY action-app /app/action-app

# RUN mkdir /app
# WORKDIR /app
CMD ["/app/action-app"]

# COPY jenkins_actions_app /app

# CMD ["/app/jenkins_actions_app"]

# ENTRYPOINT [ "go", "run" ,"main.go" ]
