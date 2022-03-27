FROM golang:1.18-alpine
WORKDIR /Gateway
COPY ./ ./
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o upserver ./cmd/app/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=0 /Gateway ./
EXPOSE 8081
CMD ["./upserver"]