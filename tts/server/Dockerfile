FROM golang:1.21 as builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 go build -o /app/main

FROM debian:bullseye-slim

COPY --from=builder /app/main /app/main

EXPOSE 80
CMD ["/app/main", "-l", ":80", "-p", "/piper/piper", "-m", "/voice.onnx"]
