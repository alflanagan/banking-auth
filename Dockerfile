FROM golang:1.25-alpine
LABEL authors="adrianflanagan"

RUN apk update; apk add make

WORKDIR /go

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8081

CMD ["/bin/ash", "-c", "make start"]
