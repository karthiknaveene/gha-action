
FROM gcr.io/distroless/base:latest


RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0  go build -o action-app main.go

WORKDIR /app
COPY action-app /app/action-app

CMD ["/app/action-app"]