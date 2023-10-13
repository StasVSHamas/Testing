FROM golang:latest
WORKDIR /app
COPY ./app .
RUN apt-get update && \
    apt-get install -y proxychains
COPY proxychains.conf /etc/proxychains.conf
RUN go mod download && go mod verify
RUN go build -o crazybird crazybird.go
RUN chmod +x crazybird
CMD ["tail", "-f", "/dev/null"]
