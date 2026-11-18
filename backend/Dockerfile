# Stage 1: Build binary
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy dependency definition
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary dengan flag optimasi:
# CGO_ENABLED=0 menghasilkan standalone static binary tanpa dependensi C library
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Stage 2: Final lightweight image
FROM alpine:3.20

WORKDIR /app

# Tambahkan sertifikat SSL & timezone data
RUN apk --no-cache add ca-certificates tzdata

# Copy binary dari stage builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]