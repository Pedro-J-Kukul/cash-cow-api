FROM golang:1.25 AS builder

# setting the working directory
WORKDIR /api

# Installing git
RUN apk add --no-cache git

# copy go mod and sum files
COPY go.mod go.sum ./

# download dependencies
RUN go mod download

# copy source code
COPY . .

# building the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage
FROM alpine:latest

# Running IDK
RUN apk --no-cache add ca-certificates tzdata

# Setting hte working directory
WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy migration files
COPY --from=builder /app/migrations ./migrations

# exposing the main port
EXPOSE 8080

CMD ["./main"]