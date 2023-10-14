FROM golang:latest
WORKDIR /app
COPY ./app .
RUN apt-get update && \
    apt-get install -y proxychains && \
    go mod download && go mod verify && \
    go build -o crazybird crazybird.go && \
    chmod +x crazybird
COPY proxychains.conf /etc/proxychains.conf
CMD ["tail", "-f", "/dev/null"]