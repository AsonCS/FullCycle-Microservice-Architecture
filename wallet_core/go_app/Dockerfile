FROM golang:1.24.5

WORKDIR /app/

RUN apt-get update && apt-get install -y librdkafka-dev

# CMD ["tail", "-f", "/dev/null"]
CMD ["go", "run", "cmd/walletcore/main.go"]